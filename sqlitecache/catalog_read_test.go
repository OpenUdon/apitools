package sqlitecache

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/OpenUdon/apitools/catalog"
)

func TestReadCatalogArtifactsBoundsTimestampText(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "spec.json"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	cache, err := Open(filepath.Join(dir, catalog.DefaultRegistryPath))
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	if err := cache.StoreCatalogArtifact(context.Background(), CatalogArtifact{ProviderID: "example", ArtifactID: "example", Kind: "openapi", Path: "spec.json"}); err != nil {
		t.Fatal(err)
	}
	if _, err := cache.db.Exec("UPDATE catalog_artifacts SET updated_at = ?", strings.Repeat("x", MaxCatalogRegistrationBytes+1)); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadCatalogArtifacts(context.Background(), catalog.RootOptions{Directory: dir}); err == nil || !strings.Contains(err.Error(), "snapshot bounds") {
		t.Fatalf("oversized timestamp bypassed aggregate bound: %v", err)
	}
}

func TestReadCatalogArtifactsDoesNotCreateMissingRoot(t *testing.T) {
	rows, err := ReadCatalogArtifacts(context.Background(), catalog.RootOptions{})
	if err != nil || len(rows) != 0 {
		t.Fatalf("no root: %v, %v", rows, err)
	}
	dir := t.TempDir()
	_, err = ReadCatalogArtifacts(context.Background(), catalog.RootOptions{Directory: dir})
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("missing registry: %v", err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Fatal("created missing registry")
	}
}

func TestReadCatalogArtifactsSyntheticSnapshot(t *testing.T) {
	dir := t.TempDir()
	fixture := filepath.Join("..", "testdata", "catalog-root")
	if err := os.Mkdir(filepath.Join(dir, "openapi"), 0o700); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(fixture, "openapi", "notes.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "openapi", "notes.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	var cat catalog.Catalog
	data, err = os.ReadFile(filepath.Join(fixture, "catalog.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &cat); err != nil {
		t.Fatal(err)
	}
	if err := cat.Validate(); err != nil {
		t.Fatal(err)
	}
	if _, ok := cat.FindProvider("Example Notes Shared"); !ok {
		t.Fatal("multiword provider missing")
	}
	var expected []CatalogArtifact
	data, err = os.ReadFile(filepath.Join(fixture, "registrations.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &expected); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, catalog.DefaultRegistryPath)
	cache, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, artifact := range expected {
		if err := cache.StoreCatalogArtifact(context.Background(), artifact); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := cache.db.Exec("UPDATE catalog_artifacts SET accessed_at = 123"); err != nil {
		t.Fatal(err)
	}
	if err := cache.Close(); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := ReadCatalogArtifacts(context.Background(), catalog.RootOptions{Directory: dir})
	if err != nil {
		t.Fatal(err)
	}
	if len(actual) != len(expected) {
		t.Fatalf("got %d rows", len(actual))
	}
	for i := range actual {
		actual[i].StoredAt = expected[i].StoredAt
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("snapshot mismatch: %#v", actual)
	}
	after, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("reader changed registry: %v", err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var accessed int64
	if err := db.QueryRow("SELECT accessed_at FROM catalog_artifacts LIMIT 1").Scan(&accessed); err != nil || accessed != 123 {
		t.Fatalf("access time changed: %d %v", accessed, err)
	}
	if _, err := db.Exec("UPDATE schema_meta SET version = 99"); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadCatalogArtifacts(context.Background(), catalog.RootOptions{Directory: dir}); err == nil {
		t.Fatal("accepted newer schema")
	}
	if _, err := db.Exec("UPDATE schema_meta SET version = 1"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DROP INDEX idx_catalog_artifacts_accessed_at"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("ALTER TABLE catalog_artifacts DROP COLUMN accessed_at"); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadCatalogArtifacts(context.Background(), catalog.RootOptions{Directory: dir}); err != nil {
		t.Fatalf("read legacy schema without migration: %v", err)
	}
	if _, err := db.Exec("UPDATE catalog_artifacts SET path = '../outside'"); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadCatalogArtifacts(context.Background(), catalog.RootOptions{Directory: dir}); err == nil {
		t.Fatal("accepted tampered registration path")
	}
	if _, err := db.Exec("UPDATE catalog_artifacts SET path = 'openapi/notes.json', metadata_json = zeroblob(?)", MaxCatalogRegistrationBytes); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadCatalogArtifacts(context.Background(), catalog.RootOptions{Directory: dir}); err == nil {
		t.Fatal("accepted oversized aggregate registry")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := ReadCatalogArtifacts(ctx, catalog.RootOptions{Directory: dir}); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation lost: %v", err)
	}
}
