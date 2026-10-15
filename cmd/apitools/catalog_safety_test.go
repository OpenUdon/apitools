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

func TestCatalogNestedRegistryLifecycle(t *testing.T) {
	root := t.TempDir()
	t.Setenv("APITOOLS_CATALOG_ROOT", root)
	t.Setenv("APITOOLS_CATALOG_REGISTRY", "state/db/cache.sqlite")
	t.Setenv("APITOOLS_CATALOG_INDEX", "metadata/index.json")
	source := filepath.Join(t.TempDir(), "spec.json")
	os.WriteFile(source, []byte(lifecycleSpec), 0600)
	commands := [][]string{
		{"catalog", "init"},
		{"catalog", "add", source, "--provider", "github", "--spec", "github-rest-api-openapi"},
		{"catalog", "index"},
		{"catalog", "search", "list users", "--json"},
		{"catalog", "snapshot", "--target", filepath.Join(t.TempDir(), "snapshot")},
	}
	for _, args := range commands {
		var out, errOut bytes.Buffer
		if code := run(args, &out, &errOut); code != exitSuccess {
			t.Fatalf("%v: %d %s", args, code, errOut.String())
		}
		if args[1] == "search" && !strings.Contains(out.String(), "listUsers") {
			t.Fatalf("lost indexed operation: %s", out.String())
		}
	}
	for _, rel := range []string{"state/db/cache.sqlite", "metadata/index.json", "specs/github/spec.json"} {
		if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
			t.Fatal(err)
		}
	}
	rows, err := sqlitecache.ReadCatalogArtifacts(context.Background(), catalog.RootOptionsFromEnvironment())
	if err != nil || len(rows) != 1 || rows[0].Path != "specs/github/spec.json" {
		t.Fatalf("root-relative rows: %#v %v", rows, err)
	}
}

func TestCatalogFlagEnvironmentPrecedence(t *testing.T) {
	for _, mode := range []string{"canonical", "alias", "explicit", "empty-explicit", "defaults"} {
		t.Run(mode, func(t *testing.T) {
			for _, key := range []string{"APITOOLS_CATALOG_ROOT", "CATALOG_ROOT", "KINET_CATALOG_ROOT", "APITOOLS_CATALOG_REGISTRY", "CATALOG_REGISTRY", "KINET_CATALOG_REGISTRY", "APITOOLS_CATALOG_INDEX", "CATALOG_INDEX", "KINET_CATALOG_INDEX"} {
				t.Setenv(key, "")
			}
			root := t.TempDir()
			t.Setenv("APITOOLS_CATALOG_ROOT", root)
			wantReg, wantIdx := "cache.sqlite", "operations.v1.json"
			args := []string{"catalog", "init"}
			if mode != "defaults" {
				t.Setenv("KINET_CATALOG_REGISTRY", "legacy.sqlite")
				t.Setenv("KINET_CATALOG_INDEX", "legacy.json")
				t.Setenv("CATALOG_REGISTRY", "alias.sqlite")
				t.Setenv("CATALOG_INDEX", "alias.json")
				wantReg, wantIdx = "alias.sqlite", "alias.json"
			}
			if mode == "canonical" || mode == "explicit" || mode == "empty-explicit" {
				t.Setenv("APITOOLS_CATALOG_REGISTRY", "canonical.sqlite")
				t.Setenv("APITOOLS_CATALOG_INDEX", "canonical.json")
				wantReg, wantIdx = "canonical.sqlite", "canonical.json"
			}
			if mode == "explicit" {
				args = append(args, "--registry", "explicit.sqlite", "--index", "explicit.json")
				wantReg, wantIdx = "explicit.sqlite", "explicit.json"
			}
			if mode == "empty-explicit" {
				args = append(args, "--registry", "", "--index", "")
				wantReg, wantIdx = "cache.sqlite", "operations.v1.json"
			}
			var out, errOut bytes.Buffer
			if code := run(args, &out, &errOut); code != exitSuccess {
				t.Fatalf("%d %s", code, errOut.String())
			}
			for _, rel := range []string{wantReg, wantIdx} {
				if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
