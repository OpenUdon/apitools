package apitools

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/catalogupdate"
)

const catalogUpdateDatedYAML = `openapi: 3.0.3
info:
  title: Test records
  version: 2024-01-01
paths:
  /records:
    get:
      responses:
        200:
          description: ok
          content:
            application/json:
              schema:
                type: string
                example: 2024-01-01
`

func TestCatalogUpdateNativeVersionHasIndependentBounds(t *testing.T) {
	expanded := strings.Repeat("a", 12<<20)
	for _, test := range []struct {
		name, artifact, version string
		valid                   bool
	}{
		{"date", catalogUpdateDatedYAML, "2024-01-01", true},
		{"alias", strings.Replace(catalogUpdateDatedYAML, "description: ok", "description: &ok ok\n          x-example: *ok", 1), "2024-01-01", true},
		{"merge alias", strings.Replace(catalogUpdateDatedYAML, "info:\n  title: Test records\n  version: 2024-01-01", "x-info: &info {title: Test records, version: 2023-01-01}\ninfo: {<<: *info, version: 2024-01-01}", 1), "2024-01-01", true},
		{"structural work above shape limit", `{"openapi":"3.0.3","info":{"title":"Example","version":"2.1"},"paths":{},"x-data":[` + strings.Repeat("0,", 100000) + "0]}", "2.1", true},
		{"alias cycle", catalogUpdateDatedYAML + "x-cycle: &cycle [*cycle]\n", "", false},
		{"alias expansion over byte budget", catalogUpdateDatedYAML + "x-data: &value " + expanded + "\nx-other: *value\n", "", false},
		{"duplicate field", catalogUpdateDatedYAML + "info: {title: duplicate, version: '1'}\n", "", false},
		{"invalid explicit bool", catalogUpdateDatedYAML + "x-invalid: !!bool maybe\n", "", false},
		{"invalid explicit null", catalogUpdateDatedYAML + "x-invalid: !!null maybe\n", "", false},
		{"invalid explicit date", catalogUpdateDatedYAML + "x-invalid: !!timestamp invalid\n", "", false},
		{"multiple documents", catalogUpdateDatedYAML + "---\n" + catalogUpdateDatedYAML, "", false},
		{"pinned external reference", catalogUpdateDatedYAML + "components: {schemas: {External: {$ref: 'https://example.com/external.json'}}}\n", "2024-01-01", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			pin := fmt.Sprintf("%x", sha256.Sum256([]byte(test.artifact)))
			got := catalogUpdateNativeOpenAPI(t.Context(), []byte(test.artifact), "https://example.com/openapi.yaml", catalog.RefreshableSpecReference{}, pin)
			if got.accepted() != test.valid || got.Version != test.version {
				t.Fatalf("native validation = %+v", got)
			}
			if test.valid {
				if _, err := decodeShapeDocument(t.Context(), []byte(test.artifact)); err == nil {
					t.Fatal("metadata support weakened the independent shape decoder")
				}
			}
		})
	}
}

func TestCatalogUpdateAdoptsExactPinAfterReviewedCorrection(t *testing.T) {
	const commit = "1111111111111111111111111111111111111111"
	const body = `{"openapi":"3.0.3","info":{"title":"Spotify test fixture","version":"1.0"},"servers":[{"url":"https://example.com"}],"security":[],"paths":{"/records":{"get":{"responses":{"200":{"description":"ok","content":{"application/json":{"schema":{"type":"string"}}}}}}}},"components":{"schemas":{"Unused":{"type":"object","required":true}}}}`
	source := catalogUpdateSource{
		ProviderID: "spotify", ReferenceID: "spotify-web-api-openapi", Service: "web-api",
		Kind: catalog.SpecKindOpenAPI, URL: "https://raw.githubusercontent.com/fixture/spotify/" + commit + "/openapi.json",
		HeldRevision: fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(body))),
	}
	ctx := t.Context()
	dir, err := os.MkdirTemp(".", ".catalog-update-validation-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	cache, err := catalogupdate.OpenArtifactCache(filepath.Join(dir, "cache"), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	store, err := catalogupdate.OpenStore(ctx, filepath.Join(dir, "state.sqlite"), cache)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	run, err := store.StartRun(ctx, catalog.BuiltInCatalogManifest().CatalogSHA256, []string{"spotify"})
	if err != nil {
		t.Fatal(err)
	}
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host == "api.github.com" && r.URL.Path == "/repos/fixture/spotify/commits" {
			return catalogRepairResponse(r, fmt.Sprintf(`[{"sha":%q}]`, commit)), nil
		}
		if r.URL.String() != source.URL {
			t.Fatalf("unexpected request: %s", r.URL)
		}
		return catalogRepairResponse(r, body), nil
	})}}
	provider, _ := catalog.FindBuiltInProvider(source.ProviderID)
	out, err := client.checkCatalogUpdateSource(ctx, store, run.ID, provider, source, CatalogUpdateConfig{}, nil)
	if err != nil || out.Status != "reviewed-pin-reproduced" || out.Object == nil ||
		!strings.Contains(string(out.Evidence), `"corrected_status":"valid-openapi"`) {
		t.Fatalf("reviewed update = %s, %s, %v", out.Status, out.Evidence, err)
	}
	f, err := store.OpenObject(ctx, *out.Object)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	got, err := io.ReadAll(f)
	if err != nil || string(got) != body {
		t.Fatalf("retained bytes differ from reviewed publisher body, %v", err)
	}
	reconstructedStore, err := catalogupdate.OpenStore(ctx, filepath.Join(dir, "reconstructed.sqlite"), cache)
	if err != nil {
		t.Fatal(err)
	}
	defer reconstructedStore.Close()
	reconstructed, err := client.reconstructCatalogUpdateBaseline(ctx, reconstructedStore, source)
	if err != nil || reconstructed.Preferred.Object != *out.Object || reconstructed.Preferred.Version != "1.0" || reconstructed.Preferred.Revision != commit || reconstructed.OpenAPI == nil || reconstructed.OpenAPI.CorrectedStatus != CatalogRefreshValidOpenAPI {
		t.Fatalf("reconstructed baseline = %+v, %v", reconstructed, err)
	}
	for _, pin := range []string{"", strings.Repeat("0", 64)} {
		got := catalogUpdateNativeOpenAPI(ctx, []byte(body), source.URL, source.refreshReference(), pin)
		if got.accepted() {
			t.Fatalf("unreviewed artifact used a correction: %+v", got)
		}
	}
	corrected, _ := correctedCatalogRefreshContent(source.refreshReference(), []byte(body))
	original := OperationSourceInput{Kind: OperationSourceOpenAPI, URL: source.URL, Content: []byte(body)}
	replacement := OperationSourceInput{Kind: OperationSourceOpenAPI, URL: source.URL, Content: corrected}
	table, err := BuildOperationShapeTable(ctx, OperationShapeOptions{
		Sources: []ShapeSourceInput{{ID: "catalog-update", OperationSourceInput: original}},
	})
	if err != nil || len(table.Operations) != 1 || !knownCatalogUpdateShape(table.Operations[0]) {
		t.Fatalf("comparison fixture must have a complete raw shape, %v", err)
	}
	if got := compareCatalogUpdate(ctx, original, replacement, []string{"#/paths/~1records/get"}, source.refreshReference(), catalogUpdateReviewedDigest(source.HeldRevision)); got != catalogUpdateUnexamined {
		t.Fatalf("comparison must mark the corrected baseline unexamined, got %s", got)
	}
}

func TestCatalogUpdatePinnedNativeEvidence(t *testing.T) {
	dir := os.Getenv("APITOOLS_CATALOG_QUALIFICATION_DIR")
	if dir == "" {
		dir = "catalog-openapi-cache"
	}
	for _, id := range []string{"nvidia-dsx-air", "datadog", "digitalocean", "hugging-face-hub", "miro", "mongodb-atlas", "plaid", "vercel", "square"} {
		t.Run(id, func(t *testing.T) {
			provider, ok := catalog.FindBuiltInProvider(id)
			if !ok {
				t.Fatal("provider missing")
			}
			source := catalogUpdateReferenceSources(provider)[0]
			var raw []byte
			for _, ext := range []string{".json", ".yaml", ".yml"} {
				b, err := os.ReadFile(filepath.Join(dir, "openapi", source.ReferenceID+ext))
				if os.IsNotExist(err) && id == "nvidia-dsx-air" {
					b, err = os.ReadFile(filepath.Join("catalog-openapi-cache", "openapi", source.ReferenceID+ext))
				}
				if os.IsNotExist(err) {
					continue
				}
				if err != nil {
					t.Fatal(err)
				}
				raw = b
				break
			}
			if raw == nil {
				t.Skip("optional publisher snapshot not staged locally")
			}
			digest := sha256.Sum256(raw)
			pin := catalogUpdateReviewedDigest(source.HeldRevision)
			if hex.EncodeToString(digest[:]) != pin {
				t.Fatal("snapshot does not match the reviewed catalog pin")
			}
			got := catalogUpdateNativeOpenAPI(t.Context(), raw, source.URL, source.refreshReference(), pin)
			if id == "square" {
				if got.accepted() || got.RawStatus != CatalogRefreshInvalid {
					t.Fatalf("Square evidence = %+v", got)
				}
			} else if !got.accepted() || got.ReviewedRoute == "" {
				t.Fatalf("expected reviewed route for %q, evidence = %+v", source.ReferenceID, got)
			}
		})
	}
}
