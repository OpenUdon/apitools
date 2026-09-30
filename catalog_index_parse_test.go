package apitools

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/OpenUdon/apitools/catalog"
	_ "modernc.org/sqlite"
)

func TestRegisteredCatalogArtifactLargeAndIdentity(t *testing.T) {
	root := t.TempDir()
	content, err := os.ReadFile("testdata/catalog-root/openapi/notes.json")
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(content, &doc); err != nil {
		t.Fatal(err)
	}
	doc["x-synthetic-padding"] = strings.Repeat("x", int(DefaultMaxBytes))
	content, err = json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(content) <= int(DefaultMaxBytes) {
		t.Fatal("fixture is not large")
	}
	if err := os.WriteFile(filepath.Join(root, "notes.json"), content, 0o600); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(content)
	registration := catalog.CatalogSpecArtifact{ProviderID: "example-notes", SpecRefID: "synthetic-notes-openapi", ArtifactID: "synthetic-notes-openapi", Kind: "openapi", Path: "notes.json", SHA256: hex.EncodeToString(digest[:]), Bytes: int64(len(content))}
	result, err := parseRegisteredCatalogArtifact(context.Background(), root, registration)
	if err != nil || result.operationCount != 2 || len(result.candidates) != 2 {
		t.Fatalf("registered large metadata: operations=%d candidates=%d err=%v", result.operationCount, len(result.candidates), err)
	}
	if _, err := parseInventoryDocument(content); err == nil {
		t.Fatal("ordinary inventory parser accepted >20 MiB")
	}
	if _, err := BuildOperationCandidates(context.Background(), OperationCandidateOptions{Sources: []OperationSourceInput{{Kind: OperationSourceOpenAPI, Content: content}}}); err == nil {
		t.Fatal("ordinary candidates accepted >20 MiB")
	}
	wrong := registration
	wrong.SHA256 = strings.Repeat("0", 64)
	if _, err := parseRegisteredCatalogArtifact(context.Background(), root, wrong); err == nil {
		t.Fatal("accepted digest mismatch")
	}
	wrong = registration
	wrong.SpecRefID = ""
	if _, err := parseRegisteredCatalogArtifact(context.Background(), root, wrong); err == nil {
		t.Fatal("accepted missing registration identity")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := parseRegisteredCatalogArtifact(ctx, root, registration); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation: %v", err)
	}
	file, err := os.Create(filepath.Join(root, "oversize.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Truncate(maxCatalogArtifactBytes + 1); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	wrong = registration
	wrong.Path = "oversize.json"
	wrong.Bytes = maxCatalogArtifactBytes + 1
	if _, err := parseRegisteredCatalogArtifact(context.Background(), root, wrong); err == nil {
		t.Fatal("accepted >128 MiB")
	}
}

func TestCatalogOpenAPIStructuralRefusal(t *testing.T) {
	for _, content := range [][]byte{
		[]byte("openapi: 3.0.3\npaths: &paths {}\nx-copy: *paths\n"),
		[]byte(`{"openapi":"3.0.3","paths":{}} {"ignored":true}`),
		[]byte(`{"openapi":"4.0.0","paths":{}}`),
		append(append([]byte(`{"openapi":"3.0.3","x":`), bytes.Repeat([]byte("["), 101)...), append([]byte("0"), bytes.Repeat([]byte("]"), 101)...)...),
	} {
		if _, err := parseCatalogOpenAPI(context.Background(), content); err == nil {
			t.Fatal("accepted unsafe or unsupported structure")
		}
	}
}

// This is explicitly opt-in, uses already registered local files only and
// prints measurements/counts, never source content or runtime credentials.
func TestCatalogRegisteredArtifactMeasurement(t *testing.T) {
	root := os.Getenv("APITOOLS_CATALOG_MEASURE_ROOT")
	if root == "" {
		t.Skip("explicit local artifact measurement only")
	}
	paths, err := catalog.ResolveRoot(catalog.RootOptions{Directory: root})
	if err != nil {
		t.Fatal(err)
	}
	uri := url.URL{Scheme: "file", Path: filepath.ToSlash(paths.RegistryPath), RawQuery: "mode=ro"}
	db, err := sql.Open("sqlite", uri.String())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for _, id := range []string{"cloudflare-api-openapi", "microsoft-graph-v1-openapi"} {
		t.Run(id, func(t *testing.T) {
			var artifact catalog.CatalogSpecArtifact
			err := db.QueryRow(`SELECT provider_id,artifact_id,kind,path,source_url,sha256,bytes FROM catalog_artifacts WHERE artifact_id = ? ORDER BY provider_id LIMIT 1`, id).Scan(&artifact.ProviderID, &artifact.ArtifactID, &artifact.Kind, &artifact.Path, &artifact.SourceURL, &artifact.SHA256, &artifact.Bytes)
			if err != nil {
				t.Fatal(err)
			}
			artifact.SpecRefID = artifact.ArtifactID
			start := time.Now()
			result, parseErr := parseRegisteredCatalogArtifact(context.Background(), paths.Directory, artifact)
			t.Logf("bytes=%d operations=%d candidates=%d diagnostics=%d elapsed=%s metadata_complete=%t", artifact.Bytes, result.operationCount, len(result.candidates), len(result.diagnostics), time.Since(start), parseErr == nil)
			// Unsupported operation summaries are explicit incomplete coverage,
			// not an artifact size refusal or a claim that no API exists.
			var diagnostics DiagnosticError
			if parseErr != nil && !errors.As(parseErr, &diagnostics) {
				t.Fatal(parseErr)
			}
			if result.operationCount == 0 {
				t.Fatal("large registered artifact was not inspected")
			}
		})
	}
}
