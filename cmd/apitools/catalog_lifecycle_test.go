package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/apitools/catalog"
	_ "modernc.org/sqlite"
)

func TestCatalogInitDispatchAndExecution(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"catalog", "init", "--help"}, &out, &errOut); code != exitSuccess || !strings.Contains(out.String(), "--root DIRECTORY") {
		t.Fatalf("help: %d %s %s", code, out.String(), errOut.String())
	}
	out.Reset()
	errOut.Reset()

	// Missing root
	if code := run([]string{"catalog", "init"}, &out, &errOut); code != exitUsage {
		t.Fatalf("expected exitUsage for missing root, got %d: %s", code, errOut.String())
	}
	out.Reset()
	errOut.Reset()

	// Positional args rejected
	if code := run([]string{"catalog", "init", "--root", t.TempDir(), "extra"}, &out, &errOut); code != exitUsage {
		t.Fatalf("expected exitUsage for extra positional arg, got %d", code)
	}
	out.Reset()
	errOut.Reset()

	// Root dir = filesystem root rejected
	if code := run([]string{"catalog", "init", "--root", "/"}, &out, &errOut); code != exitRuntime {
		t.Fatalf("expected exitRuntime for filesystem root, got %d", code)
	}
	out.Reset()
	errOut.Reset()

	// Successful init
	dir := filepath.Join(t.TempDir(), "new-catalog")
	if code := run([]string{"catalog", "init", "--root", dir}, &out, &errOut); code != exitSuccess {
		t.Fatalf("expected exitSuccess, got %d: %s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "initialized catalog root") {
		t.Fatalf("unexpected output: %s", out.String())
	}
	// Verify directory structure
	if fi, err := os.Stat(filepath.Join(dir, "specs")); err != nil || !fi.IsDir() {
		t.Fatalf("specs dir missing or not a dir: %v", err)
	}
	if fi, err := os.Stat(filepath.Join(dir, catalog.DefaultRegistryPath)); err != nil || !fi.Mode().IsRegular() {
		t.Fatalf("registry missing: %v", err)
	}
	if fi, err := os.Stat(filepath.Join(dir, catalog.DefaultOperationIndexPath)); err != nil || !fi.Mode().IsRegular() {
		t.Fatalf("index missing: %v", err)
	}

	// JSON output mode
	dirJSON := filepath.Join(t.TempDir(), "json-catalog")
	out.Reset()
	errOut.Reset()
	if code := run([]string{"catalog", "init", "--root", dirJSON, "--json"}, &out, &errOut); code != exitSuccess {
		t.Fatalf("expected exitSuccess for json init, got %d: %s", code, errOut.String())
	}
	var res struct {
		SchemaVersion string `json:"schema_version"`
		Root          string `json:"root"`
		Registry      string `json:"registry"`
		Index         string `json:"index"`
	}
	if err := json.Unmarshal(out.Bytes(), &res); err != nil {
		t.Fatalf("invalid json: %v (%s)", err, out.String())
	}
	if res.SchemaVersion != "apitools.catalog-init/v1" || res.Root != dirJSON {
		t.Fatalf("unexpected json fields: %#v", res)
	}
}

func TestCatalogAddDispatchAndExecution(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"catalog", "add", "--help"}, &out, &errOut); code != exitSuccess || !strings.Contains(out.String(), "Usage: apitools catalog add") {
		t.Fatalf("help: %d %s %s", code, out.String(), errOut.String())
	}
	out.Reset()
	errOut.Reset()

	// Missing file argument
	if code := run([]string{"catalog", "add"}, &out, &errOut); code != exitUsage {
		t.Fatalf("expected exitUsage for missing file, got %d", code)
	}
	out.Reset()
	errOut.Reset()

	// Missing root
	validSpec := []byte(`{"openapi":"3.0.0","info":{"title":"Pet Store","version":"1.0.0"},"paths":{"/pets":{"get":{"operationId":"listPets","responses":{"200":{"description":"ok"}}}}}}`)
	specFile := filepath.Join(t.TempDir(), "pets.json")
	if err := os.WriteFile(specFile, validSpec, 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"catalog", "add", specFile}, &out, &errOut); code != exitUsage {
		t.Fatalf("expected exitUsage for missing root, got %d", code)
	}
	out.Reset()
	errOut.Reset()

	// Initialize catalog root
	catDir := filepath.Join(t.TempDir(), "cat-add")
	if code := run([]string{"catalog", "init", "--root", catDir}, &out, &errOut); code != exitSuccess {
		t.Fatalf("init failed: %d %s", code, errOut.String())
	}
	out.Reset()
	errOut.Reset()

	// Add invalid spec file (not openapi)
	invalidSpec := filepath.Join(t.TempDir(), "invalid.txt")
	if err := os.WriteFile(invalidSpec, []byte("hello world"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"catalog", "add", invalidSpec, "--root", catDir}, &out, &errOut); code != exitRuntime {
		t.Fatalf("expected exitRuntime for invalid spec, got %d", code)
	}
	out.Reset()
	errOut.Reset()

	writeLifecycleCatalog(t, catDir, map[string][]string{"petstore": {"pets-v1"}, "pet-store": {"internal"}})

	// Add valid spec from outside root
	if code := run([]string{"catalog", "add", specFile, "--root", catDir, "--provider", "petstore", "--spec", "pets-v1", "--catalog", "catalog.json"}, &out, &errOut); code != exitSuccess {
		t.Fatalf("add failed: %d %s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "added petstore/pets-v1") {
		t.Fatalf("unexpected add output: %s", out.String())
	}

	// Verify file was copied into specs/petstore/pets.json
	copiedPath := filepath.Join(catDir, "specs", "petstore", "pets.json")
	if fi, err := os.Stat(copiedPath); err != nil || !fi.Mode().IsRegular() {
		t.Fatalf("copied spec file missing: %v", err)
	}

	// Verify artifact registered in SQLite database
	db, err := sql.Open("sqlite", filepath.Join(catDir, catalog.DefaultRegistryPath))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var count int
	if err := db.QueryRow("SELECT count(*) FROM catalog_artifacts WHERE provider_id = 'petstore' AND artifact_id = 'pets-v1'").Scan(&count); err != nil || count != 1 {
		t.Fatalf("artifact not found in sqlite: count=%d, err=%v", count, err)
	}

	// Add file from inside root with JSON output
	insidePath := filepath.Join(catDir, "specs", "internal.json")
	if err := os.WriteFile(insidePath, validSpec, 0o644); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errOut.Reset()
	if code := run([]string{"catalog", "add", insidePath, "--root", catDir, "--catalog", "catalog.json", "--json"}, &out, &errOut); code != exitSuccess {
		t.Fatalf("add inside root failed: %d %s", code, errOut.String())
	}
	var artifactRes struct {
		ProviderID string `json:"provider_id"`
		ArtifactID string `json:"artifact_id"`
		Path       string `json:"path"`
		SHA256     string `json:"sha256"`
	}
	if err := json.Unmarshal(out.Bytes(), &artifactRes); err != nil {
		t.Fatalf("invalid json: %v (%s)", err, out.String())
	}
	if artifactRes.Path != "specs/internal.json" {
		t.Fatalf("expected specs/internal.json relative path, got %q", artifactRes.Path)
	}
}

func TestCatalogImportDispatchAndExecution(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"catalog", "import", "--help"}, &out, &errOut); code != exitSuccess || !strings.Contains(out.String(), "Usage: apitools catalog import") {
		t.Fatalf("help: %d %s %s", code, out.String(), errOut.String())
	}
	out.Reset()
	errOut.Reset()

	// Missing URL
	if code := run([]string{"catalog", "import"}, &out, &errOut); code != exitUsage {
		t.Fatalf("expected exitUsage for missing URL, got %d", code)
	}
	out.Reset()
	errOut.Reset()

	// Invalid URL
	if code := run([]string{"catalog", "import", "ftp://example.com/spec.json"}, &out, &errOut); code != exitUsage {
		t.Fatalf("expected exitUsage for invalid scheme, got %d", code)
	}
	out.Reset()
	errOut.Reset()

	// Initialize catalog root
	catDir := filepath.Join(t.TempDir(), "cat-import")
	if code := run([]string{"catalog", "init", "--root", catDir}, &out, &errOut); code != exitSuccess {
		t.Fatalf("init failed: %d %s", code, errOut.String())
	}

	writeLifecycleCatalog(t, catDir, map[string][]string{"remote-test": {"remote-spec"}})

	// Setup local test server with valid OpenAPI spec
	specContent := `{"openapi":"3.0.0","info":{"title":"Remote API","version":"2.0.0"},"paths":{"/remote":{"get":{"operationId":"getRemote","responses":{"200":{"description":"ok"}}}}}}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(specContent))
	}))
	defer server.Close()

	out.Reset()
	errOut.Reset()
	// Import with client allowing unsafe hosts
	code := runCatalogImportWithClient(
		[]string{server.URL + "/openapi.json", "--root", catDir, "--provider", "remote-test", "--name", "remote-spec", "--catalog", "catalog.json"},
		&out,
		&errOut,
		func(cachePath string) (*apitools.Client, func(), error) {
			client, closeCache, err := clientForCache(cachePath)
			if err != nil {
				return nil, nil, err
			}
			client.AllowUnsafeHosts = true
			return client, closeCache, nil
		},
	)
	if code != exitSuccess {
		t.Fatalf("catalog import failed: %d %s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "imported remote-test/remote-spec") {
		t.Fatalf("unexpected import output: %s", out.String())
	}

	// Verify file saved in specs/remote-test/remote-spec.json
	specPath := filepath.Join(catDir, "specs", "remote-test", "remote-spec.json")
	if fi, err := os.Stat(specPath); err != nil || !fi.Mode().IsRegular() {
		t.Fatalf("imported spec file missing: %v", err)
	}

	// Verify registered in SQLite
	db, err := sql.Open("sqlite", filepath.Join(catDir, catalog.DefaultRegistryPath))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var count int
	if err := db.QueryRow("SELECT count(*) FROM catalog_artifacts WHERE provider_id = 'remote-test' AND artifact_id = 'remote-spec'").Scan(&count); err != nil || count != 1 {
		t.Fatalf("artifact not found in sqlite: count=%d, err=%v", count, err)
	}
}

func TestCatalogSnapshotDispatchAndExecution(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"catalog", "snapshot", "--help"}, &out, &errOut); code != exitSuccess || !strings.Contains(out.String(), "--target DIRECTORY") {
		t.Fatalf("help: %d %s %s", code, out.String(), errOut.String())
	}
	out.Reset()
	errOut.Reset()

	// Missing target
	if code := run([]string{"catalog", "snapshot"}, &out, &errOut); code != exitUsage {
		t.Fatalf("expected exitUsage for missing target, got %d", code)
	}
	out.Reset()
	errOut.Reset()

	// Initialize catalog root
	catDir := filepath.Join(t.TempDir(), "cat-snap-src")
	if code := run([]string{"catalog", "init", "--root", catDir}, &out, &errOut); code != exitSuccess {
		t.Fatalf("init failed: %d %s", code, errOut.String())
	}

	targetDir := filepath.Join(t.TempDir(), "cat-snap-dest")
	out.Reset()
	errOut.Reset()
	if code := run([]string{"catalog", "snapshot", "--root", catDir, "--target", targetDir}, &out, &errOut); code != exitSuccess {
		t.Fatalf("snapshot failed: %d %s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "snapshot created in") {
		t.Fatalf("unexpected snapshot output: %s", out.String())
	}

	// Verify destination registry exists and journal mode is DELETE
	destReg := filepath.Join(targetDir, catalog.DefaultRegistryPath)
	if fi, err := os.Stat(destReg); err != nil || !fi.Mode().IsRegular() {
		t.Fatalf("destination registry missing: %v", err)
	}
	db, err := sql.Open("sqlite", destReg)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var mode string
	if err := db.QueryRow("PRAGMA journal_mode").Scan(&mode); err != nil || mode != "delete" {
		t.Fatalf("expected journal_mode=delete, got %q (err=%v)", mode, err)
	}

	// Test JSON mode
	targetJSON := filepath.Join(t.TempDir(), "cat-snap-json")
	out.Reset()
	errOut.Reset()
	if code := run([]string{"catalog", "snapshot", "--root", catDir, "--target", targetJSON, "--json"}, &out, &errOut); code != exitSuccess {
		t.Fatalf("snapshot json failed: %d %s", code, errOut.String())
	}
	var snapRes struct {
		SchemaVersion string `json:"schema_version"`
		Target        string `json:"target"`
		Normalized    bool   `json:"normalized"`
	}
	if err := json.Unmarshal(out.Bytes(), &snapRes); err != nil {
		t.Fatalf("invalid json: %v (%s)", err, out.String())
	}
	if snapRes.SchemaVersion != "apitools.catalog-snapshot/v1" || snapRes.Target != targetJSON || !snapRes.Normalized {
		t.Fatalf("unexpected snapshot json: %#v", snapRes)
	}
}

func TestCatalogSearchDispatchAndExecution(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"catalog", "search", "--help"}, &out, &errOut); code != exitSuccess || !strings.Contains(out.String(), "Usage: apitools catalog search") {
		t.Fatalf("help: %d %s %s", code, out.String(), errOut.String())
	}
	out.Reset()
	errOut.Reset()

	// Missing query
	if code := run([]string{"catalog", "search"}, &out, &errOut); code != exitUsage {
		t.Fatalf("expected exitUsage for missing query, got %d", code)
	}
	out.Reset()
	errOut.Reset()

	// Search built-in providers when no root is configured
	if code := run([]string{"catalog", "search", "github"}, &out, &errOut); code != exitSuccess {
		t.Fatalf("search built-in providers failed: %d %s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "GitHub") {
		t.Fatalf("expected GitHub in provider search output, got: %s", out.String())
	}

	// Initialize catalog, add spec, index, and search indexed operations
	catDir := filepath.Join(t.TempDir(), "cat-search")
	if code := run([]string{"catalog", "init", "--root", catDir}, &out, &errOut); code != exitSuccess {
		t.Fatalf("init failed: %d %s", code, errOut.String())
	}

	validSpec := []byte(`{"openapi":"3.0.0","info":{"title":"User API","version":"1.0.0"},"paths":{"/users":{"get":{"operationId":"listUsers","summary":"List all registered users","responses":{"200":{"description":"ok"}}}}}}`)
	specFile := filepath.Join(t.TempDir(), "users.json")
	if err := os.WriteFile(specFile, validSpec, 0o644); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errOut.Reset()
	if code := run([]string{"catalog", "add", specFile, "--root", catDir, "--provider", "github", "--spec", "github-rest-api-openapi"}, &out, &errOut); code != exitSuccess {
		t.Fatalf("add failed: %d %s", code, errOut.String())
	}

	out.Reset()
	errOut.Reset()
	if code := run([]string{"catalog", "index", "--root", catDir}, &out, &errOut); code != exitSuccess {
		t.Fatalf("index failed: %d %s", code, errOut.String())
	}

	out.Reset()
	errOut.Reset()
	if code := run([]string{"catalog", "search", "list users", "--root", catDir}, &out, &errOut); code != exitSuccess {
		t.Fatalf("search indexed operations failed: %d %s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "/users") {
		t.Fatalf("expected /users in search output, got: %s", out.String())
	}

	// JSON mode search
	out.Reset()
	errOut.Reset()
	if code := run([]string{"catalog", "search", "list users", "--root", catDir, "--json"}, &out, &errOut); code != exitSuccess {
		t.Fatalf("search json failed: %d %s", code, errOut.String())
	}
	var report struct {
		SchemaVersion string `json:"schema_version"`
		Candidates    []any  `json:"candidates"`
	}
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatalf("invalid json report: %v (%s)", err, out.String())
	}
	if len(report.Candidates) == 0 {
		t.Fatalf("expected candidates in report: %s", out.String())
	}

	// Read command never implicitly creates missing index
	noIndexDir := filepath.Join(t.TempDir(), "cat-no-index")
	if err := os.MkdirAll(noIndexDir, 0o755); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errOut.Reset()
	if code := run([]string{"catalog", "search", "github", "--root", noIndexDir}, &out, &errOut); code != exitSuccess {
		t.Fatalf("search without index failed: %d", code)
	}
	// Index must NOT have been created
	if _, err := os.Stat(filepath.Join(noIndexDir, catalog.DefaultOperationIndexPath)); !os.IsNotExist(err) {
		t.Fatal("read search implicitly created index")
	}
}

func TestCatalogIndexEnvironmentFallback(t *testing.T) {
	catDir := filepath.Join(t.TempDir(), "cat-env")
	var out, errOut bytes.Buffer
	if code := run([]string{"catalog", "init", "--root", catDir}, &out, &errOut); code != exitSuccess {
		t.Fatalf("init failed: %d %s", code, errOut.String())
	}

	// Set environment variable
	t.Setenv("APITOOLS_CATALOG_ROOT", catDir)

	out.Reset()
	errOut.Reset()
	// Run catalog index without --root flag
	if code := run([]string{"catalog", "index"}, &out, &errOut); code != exitSuccess {
		t.Fatalf("index with env fallback failed: %d %s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "indexed") {
		t.Fatalf("unexpected index output: %s", out.String())
	}
}
