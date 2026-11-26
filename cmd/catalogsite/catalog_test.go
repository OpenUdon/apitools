package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/OpenUdon/apitools/catalog"
)

const testCommit = "0123456789abcdef0123456789abcdef01234567"

func TestBuiltInPublicCatalog(t *testing.T) {
	base, err := parseBaseURL("https://apis.kinet.com")
	if err != nil {
		t.Fatal(err)
	}
	data, err := builtInSiteCatalog(testCommit, base)
	if err != nil {
		t.Fatal(err)
	}
	if data.SchemaVersion != 1 || data.BaseURL != "https://apis.kinet.com/" || data.SourceCommit != testCommit ||
		data.ContactURL != "https://github.com/OpenUdon/apitools/issues" {
		t.Fatalf("unexpected catalog metadata: version %d, base %q, commit %q, contact %q", data.SchemaVersion, data.BaseURL, data.SourceCommit, data.ContactURL)
	}
	providers := catalog.BuiltInProviders()
	if len(data.Providers) != len(providers) {
		t.Fatalf("got %d public providers, want %d", len(data.Providers), len(providers))
	}
	examples := map[string]string{
		"adyen": coverageOpenAPI, "dropbox": coverageMachine, "airtable": coverageOverlay, "acumatica": coverageDocs,
	}
	for i, provider := range data.Providers {
		if provider.ID != providers[i].ID {
			t.Fatalf("provider at index %d = %q, want %q", i, provider.ID, providers[i].ID)
		}
		if want, ok := examples[provider.ID]; ok {
			if provider.Coverage != want {
				t.Errorf("coverage for %s = %q, want %q", provider.ID, provider.Coverage, want)
			}
			delete(examples, provider.ID)
		}
		if provider.URL != "https://apis.kinet.com/providers/"+provider.ID+"/index.html" || len(provider.Sources) == 0 {
			t.Fatalf("provider %s missing canonical URL or provenance", provider.ID)
		}
		if provider.ID == "aws-elbv2" && (strings.Contains(provider.Description, "current") || !strings.Contains(provider.Description, "Gateway Load Balancers")) {
			t.Fatalf("AWS ELBv2 description has unexpected wording: %q", provider.Description)
		}
		for _, op := range provider.Operations {
			if op.Label == "Get CSRF crumb" || (provider.ID == "filemaker" && op.Label == "Delete session") {
				t.Fatalf("excluded operation label %q is public for %s", op.Label, provider.ID)
			}
			found := false
			for _, source := range provider.Sources {
				if source.ID == op.SourceID {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("operation %s on %s has no matching source", op.Label, provider.ID)
			}
		}
	}
	if len(examples) != 0 {
		t.Fatalf("missing coverage examples: %v", examples)
	}
	files, err := buildFiles(data)
	if err != nil {
		t.Fatal(err)
	}
	var encoded map[string]json.RawMessage
	if err := json.Unmarshal(files["catalog.json"], &encoded); err != nil {
		t.Fatal(err)
	}
	for _, unwanted := range []string{`"selector"`, `"overlay_path"`, `"security_schemes"`, `"source_note"`} {
		if bytes.Contains(files["catalog.json"], []byte(unwanted)) {
			t.Fatalf("internal field %s was published", unwanted)
		}
	}
}

func TestSourceLinksAndBaseURL(t *testing.T) {
	for _, raw := range []string{
		"http://api.example.com/", "https://user:pass@api.example.com/", "https://127.0.0.1/",
		"https://host.internal/", "https://api.example.com:443/", "https://api.example.com/?token=private",
		"https://api.example.com/#access_token=placeholder", "https://api.example.com/#access_token%3dplaceholder",
		"javascript:alert(1)", "//api.example.com/path",
	} {
		if publicSourceURL(raw) {
			t.Errorf("accepted unsafe URL %q", raw)
		}
	}
	if !publicSourceURL("https://api.example.com/docs?version=v1#details") {
		t.Error("public version selector was rejected")
	}
	base, err := parseBaseURL("https://apis.kinet.com/catalog")
	if err != nil || base.String() != "https://apis.kinet.com/catalog/" {
		t.Fatalf("parseBaseURL = %v, %v", base, err)
	}
	for _, raw := range []string{"https://apis.kinet.com/../one", "https://apis.kinet.com/#x", "https://apis.kinet.com/?version=v1"} {
		if _, err := parseBaseURL(raw); err == nil {
			t.Errorf("accepted invalid base URL %q", raw)
		}
	}
}

func TestCleanSourceRevision(t *testing.T) {
	settings := []debug.BuildSetting{{Key: "vcs", Value: "git"}, {Key: "vcs.revision", Value: testCommit}, {Key: "vcs.modified", Value: "false"}}
	if err := checkRevision(testCommit, settings); err != nil {
		t.Fatal(err)
	}
	if err := checkRevision(strings.Repeat("f", 40), settings); err == nil {
		t.Error("accepted the wrong commit")
	}
	settings[2].Value = "true"
	if err := checkRevision(testCommit, settings); err == nil {
		t.Error("accepted a dirty source tree")
	}
	if err := checkRevision(testCommit, nil); err == nil {
		t.Error("accepted an unstamped build")
	}
}

func TestOutputManifest(t *testing.T) {
	files := map[string][]byte{"sub/index.html": []byte("html"), "catalog.json": []byte("json")}
	root := t.TempDir()
	var previous []byte
	for _, name := range []string{"one", "two"} {
		directory := filepath.Join(root, name)
		if err := writeOutput(directory, files); err != nil {
			t.Fatal(err)
		}
		manifest, err := os.ReadFile(filepath.Join(directory, "SHA256SUMS"))
		if err != nil {
			t.Fatal(err)
		}
		want := fmt.Sprintf("%x  catalog.json\n%x  sub/index.html\n", sha256.Sum256(files["catalog.json"]), sha256.Sum256(files["sub/index.html"]))
		if string(manifest) != want || (previous != nil && !bytes.Equal(previous, manifest)) {
			t.Fatalf("unexpected manifest %q", manifest)
		}
		for relativePath, expected := range files {
			actual, err := os.ReadFile(filepath.Join(directory, relativePath))
			if err != nil || !bytes.Equal(expected, actual) {
				t.Fatalf("output %s: %q, %v", relativePath, actual, err)
			}
		}
		previous = manifest
		if err := writeOutput(directory, files); err == nil {
			t.Error("accepted an existing output directory")
		}
	}
	if err := writeOutput(filepath.Join(root, "invalid"), map[string][]byte{"../outside": nil}); err == nil {
		t.Error("accepted a path escaping the output directory")
	}
}
