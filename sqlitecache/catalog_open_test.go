package sqlitecache

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/OpenUdon/apitools/catalog"
)

func TestOpenCatalogNestedRootAndLegacySemantics(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "specs"), 0700)
	os.WriteFile(filepath.Join(root, "specs", "api.json"), []byte("{}"), 0600)
	opts := catalog.RootOptions{Directory: root, RegistryPath: "state/cache.sqlite"}
	cache, err := OpenCatalog(opts)
	if err != nil {
		t.Fatal(err)
	}
	artifact := CatalogArtifact{ProviderID: "test", ArtifactID: "test", Kind: "openapi", Path: "specs/api.json"}
	if err := cache.StoreCatalogArtifact(context.Background(), artifact); err != nil {
		t.Fatal(err)
	}
	cache.Close()
	rows, err := ReadCatalogArtifacts(context.Background(), opts)
	if err != nil || len(rows) != 1 || rows[0].Path != artifact.Path {
		t.Fatalf("%#v %v", rows, err)
	}
	legacy, err := Open(filepath.Join(root, "state", "legacy.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer legacy.Close()
	if err := legacy.StoreCatalogArtifact(context.Background(), artifact); err == nil {
		t.Fatal("legacy opener changed its base directory")
	}
	if _, err := OpenCatalog(catalog.RootOptions{Directory: root, RegistryPath: "../escape.sqlite"}); err == nil {
		t.Fatal("accepted registry escape")
	}
}
