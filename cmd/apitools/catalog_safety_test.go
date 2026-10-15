package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/sqlitecache"
)

const lifecycleSpec = `{"openapi":"3.0.0","info":{"title":"Test","version":"1"},"paths":{"/users":{"get":{"operationId":"listUsers","responses":{"200":{"description":"ok"}}}}}}`

func TestCatalogInitRejectsUnsafePathsBeforeWrites(t *testing.T) {
	for _, flag := range []string{"--registry", "--index", "--catalog"} {
		for _, value := range []string{"../escape/file", "a/../../escape", "a/../file", "/absolute", "a\\b"} {
			t.Run(flag+value, func(t *testing.T) {
				parent := t.TempDir()
				root := filepath.Join(parent, "root")
				var out, errOut bytes.Buffer
				if code := runCatalogInit([]string{"--root", root, flag, value}, &out, &errOut); code != exitRuntime {
					t.Fatalf("code %d: %s", code, errOut.String())
				}
				entries, _ := os.ReadDir(parent)
				if len(entries) != 0 {
					t.Fatalf("invalid configuration wrote files: %v", entries)
				}
			})
		}
	}
	for _, path := range []string{"specs", "cache.sqlite", "operations.v1.json"} {
		t.Run("symlink-"+path, func(t *testing.T) {
			root, external := t.TempDir(), t.TempDir()
			if err := os.Symlink(external, filepath.Join(root, path)); err != nil {
				t.Fatal(err)
			}
			var out, errOut bytes.Buffer
			if code := runCatalogInit([]string{"--root", root}, &out, &errOut); code != exitRuntime {
				t.Fatalf("accepted symlink: %s", errOut.String())
			}
			entries, _ := os.ReadDir(external)
			if len(entries) != 0 {
				t.Fatal("external mutation")
			}
			entries, _ = os.ReadDir(root)
			if len(entries) != 1 {
				t.Fatal("wrote before validation")
			}
		})
	}
}

func TestCatalogProviderTraversalRejectedBeforeWrites(t *testing.T) {
	for _, provider := range []string{"../../escape", "..", ".", "bad/name", "bad\\name"} {
		t.Run(provider, func(t *testing.T) {
			root := t.TempDir()
			source := filepath.Join(t.TempDir(), "spec.json")
			os.WriteFile(source, []byte(lifecycleSpec), 0600)
			var out, errOut bytes.Buffer
			if code := runCatalogAdd([]string{source, "--root", root, "--provider", provider, "--spec", "spec"}, &out, &errOut); code != exitRuntime {
				t.Fatalf("accepted %q: %s", provider, errOut.String())
			}
			called := false
			code := runCatalogImportWithClient([]string{"https://example.com/spec.json", "--root", root, "--provider", provider}, &out, &errOut, func(string) (*apitools.Client, func(), error) { called = true; return nil, nil, nil })
			if code != exitRuntime || called {
				t.Fatal("invalid provider reached download")
			}
			entries, _ := os.ReadDir(root)
			if len(entries) != 0 {
				t.Fatalf("wrote files: %v", entries)
			}
		})
	}
}

func TestCatalogSourceCollisionAndRegistrationRollback(t *testing.T) {
	root := t.TempDir()
	paths, err := catalog.ResolveRoot(catalog.RootOptions{Directory: root})
	if err != nil {
		t.Fatal(err)
	}
	data := []byte(lifecycleSpec)
	artifact := sqlitecache.CatalogArtifact{ProviderID: "test", ArtifactID: "test", Kind: "openapi", Path: "specs/test/spec.json"}
	ctx := context.Background()
	if err := storeCatalogSource(ctx, paths, data, artifact); err != nil {
		t.Fatal(err)
	}
	if err := storeCatalogSource(ctx, paths, data, artifact); err != nil {
		t.Fatal("identical reuse:", err)
	}
	if err := storeCatalogSource(ctx, paths, []byte("different"), artifact); err == nil {
		t.Fatal("accepted differing collision")
	}
	artifact.SourceURL = strings.Repeat("x", 2<<20)
	if err := storeCatalogSource(ctx, paths, data, artifact); err == nil {
		t.Fatal("expected registration budget error")
	}
	got, err := os.ReadFile(filepath.Join(root, artifact.Path))
	if err != nil || !bytes.Equal(got, data) {
		t.Fatal("lost existing source")
	}
	artifact.Path = "specs/test/new.json"
	if err := storeCatalogSource(ctx, paths, data, artifact); err == nil {
		t.Fatal("expected registration error")
	}
	if _, err := os.Lstat(filepath.Join(root, artifact.Path)); !os.IsNotExist(err) {
		t.Fatal("left newly copied source on failure")
	}
}
