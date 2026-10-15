package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/apitools/catalog"
)

func writeLifecycleCatalog(t *testing.T, root string, ids map[string][]string) {
	t.Helper()
	cat := catalog.Catalog{}
	for provider, specs := range ids {
		p := catalog.Provider{ID: provider, DisplayName: provider, ReviewState: catalog.ProviderCandidateOnly, OfficialOpenAPIAvailability: catalog.SpecAvailabilityUnknown, OfficialMachineSpecAvailability: catalog.SpecAvailabilityUnknown, UserOpenAPINeed: catalog.UserOpenAPINeedPossible}
		for _, id := range specs {
			p.SpecReferences = append(p.SpecReferences, catalog.SpecReference{ID: id, Kind: catalog.SpecKindOpenAPI, URL: "https://example.com/" + id + ".json", SourceAuthority: catalog.SourceAuthorityPublicCatalog, SourceNote: "Synthetic offline fixture; no official authority or license assertion."})
		}
		cat.Providers = append(cat.Providers, p)
	}
	if err := cat.Validate(); err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(cat)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "catalog.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
}

func TestCustomCatalogAddImportIndexSearchRoundTrip(t *testing.T) {
	root := t.TempDir()
	writeLifecycleCatalog(t, root, map[string][]string{"custom-service": {"local", "downloaded"}})
	source := filepath.Join(t.TempDir(), "local.json")
	os.WriteFile(source, []byte(lifecycleSpec), 0600)
	t.Setenv("APITOOLS_CATALOG_ROOT", root)
	t.Setenv("APITOOLS_CATALOG_REGISTRY", "state/registry.sqlite")
	t.Setenv("APITOOLS_CATALOG_INDEX", "state/operations.json")
	for _, args := range [][]string{{"catalog", "init", "--catalog", "catalog.json"}, {"catalog", "add", source, "--provider", "custom-service", "--spec", "local", "--catalog", "catalog.json"}} {
		var out, errOut bytes.Buffer
		if code := run(args, &out, &errOut); code != exitSuccess {
			t.Fatalf("%v: %d %s", args, code, errOut.String())
		}
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(lifecycleSpec)) }))
	defer server.Close()
	var out, errOut bytes.Buffer
	code := runCatalogImportWithClient([]string{server.URL + "/downloaded.json", "--provider", "custom-service", "--name", "downloaded", "--catalog", "catalog.json"}, &out, &errOut, func(string) (*apitools.Client, func(), error) {
		client := &apitools.Client{}
		client.AllowUnsafeHosts = true
		return client, func() {}, nil
	})
	if code != exitSuccess {
		t.Fatalf("import: %d %s", code, errOut.String())
	}
	for _, args := range [][]string{{"catalog", "index", "--catalog", "catalog.json"}, {"catalog", "search", "list users", "--catalog", "catalog.json", "--json"}} {
		out.Reset()
		errOut.Reset()
		if code := run(args, &out, &errOut); code != exitSuccess {
			t.Fatalf("%v: %d %s", args, code, errOut.String())
		}
	}
	var report apitools.CatalogDiscoveryReport
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if len(report.Candidates) == 0 || !strings.Contains(out.String(), "listUsers") || !strings.Contains(out.String(), "custom-service") {
		t.Fatalf("custom operation missing: %s", out.String())
	}
	// Index identity must not silently switch back to the built-in catalog.
	out.Reset()
	errOut.Reset()
	if code := run([]string{"catalog", "search", "list users", "--json"}, &out, &errOut); code != exitSuccess || !strings.Contains(out.String(), "discovery.index_invalid") || strings.Contains(out.String(), "listUsers") {
		t.Fatalf("wrong catalog identity was not refused: %s %s", out.String(), errOut.String())
	}
	// Metadata-only fallback also uses the selected custom catalog.
	os.Remove(filepath.Join(root, "state", "operations.json"))
	out.Reset()
	errOut.Reset()
	if code := run([]string{"catalog", "search", "custom-service", "--catalog", "catalog.json", "--json"}, &out, &errOut); code != exitSuccess || !strings.Contains(out.String(), "custom-service") {
		t.Fatalf("metadata fallback: %d %s %s", code, out.String(), errOut.String())
	}
}

func TestCatalogUnknownIdentityFailsBeforeWrites(t *testing.T) {
	for _, custom := range []bool{false, true} {
		for _, provider := range []string{"unknown-service", "github"} {
			t.Run(provider+map[bool]string{true: "custom", false: "builtin"}[custom], func(t *testing.T) {
				root := t.TempDir()
				flags := []string{"--root", root, "--provider", provider}
				if custom {
					writeLifecycleCatalog(t, root, map[string][]string{"custom-service": {"known"}})
					flags = append(flags, "--catalog", "catalog.json")
				}
				before, _ := os.ReadDir(root)
				source := filepath.Join(t.TempDir(), "spec.json")
				os.WriteFile(source, []byte(lifecycleSpec), 0600)
				var out, errOut bytes.Buffer
				args := append([]string{source, "--spec", "unknown-spec"}, flags...)
				if code := runCatalogAdd(args, &out, &errOut); code != exitRuntime {
					t.Fatal("accepted unknown identity")
				}
				called := false
				args = append([]string{"https://example.com/spec.json", "--name", "unknown-spec"}, flags...)
				code := runCatalogImportWithClient(args, &out, &errOut, func(string) (*apitools.Client, func(), error) { called = true; return nil, nil, nil })
				if code != exitRuntime || called {
					t.Fatal("unknown identity reached download")
				}
				after, _ := os.ReadDir(root)
				if len(before) != len(after) {
					t.Fatal("unknown identity wrote files")
				}
			})
		}
	}
}

func TestCatalogMetadataDecoderIsBoundedAndStrict(t *testing.T) {
	for _, data := range []string{"null", `{"unknown":true}`, `{} {}`, strings.Repeat(" ", 4<<20) + "{}", `{"providers":` + strings.Repeat("[", 101) + strings.Repeat("]", 101) + "}"} {
		root := t.TempDir()
		os.WriteFile(filepath.Join(root, "catalog.json"), []byte(data), 0600)
		if _, err := readCatalogMetadata(root, "catalog.json"); err == nil {
			t.Fatal("accepted malformed metadata")
		}
		var out, errOut bytes.Buffer
		if code := runCatalogInit([]string{"--root", root, "--catalog", "catalog.json"}, &out, &errOut); code != exitRuntime {
			t.Fatal("invalid catalog initialized")
		}
		entries, _ := os.ReadDir(root)
		if len(entries) != 1 {
			t.Fatal("invalid catalog mutated root")
		}
	}
}
