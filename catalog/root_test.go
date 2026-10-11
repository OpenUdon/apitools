package catalog

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestResolveCatalogRoot(t *testing.T) {
	paths, err := ResolveRoot(RootOptions{})
	if err != nil || paths != (RootPaths{}) {
		t.Fatalf("no root: %#v, %v", paths, err)
	}
	dir := t.TempDir()
	paths, err = ResolveRoot(RootOptions{Directory: dir})
	if err != nil || paths.RegistryPath != filepath.Join(dir, DefaultRegistryPath) || paths.IndexPath != filepath.Join(dir, DefaultOperationIndexPath) {
		t.Fatalf("explicit root: %#v, %v", paths, err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 0 {
		t.Fatalf("resolution wrote files: %v, %v", entries, err)
	}
	if _, err := ResolveRoot(RootOptions{Directory: filepath.Join(dir, "missing")}); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("missing root: %v", err)
	}
	for _, path := range []string{"../outside", "a/../cache.sqlite", "/tmp/cache.sqlite", ".", "./cache.sqlite", "a//b"} {
		if _, err := ResolveRoot(RootOptions{Directory: dir, RegistryPath: path}); err == nil {
			t.Errorf("accepted unsafe registry %q", path)
		}
	}
	if _, err := ResolveRoot(RootOptions{RegistryPath: "cache.sqlite"}); err == nil {
		t.Fatal("accepted paths without a root")
	}
	if _, err := ResolveRoot(RootOptions{Directory: dir, RegistryPath: "same", IndexPath: "same/child"}); err == nil {
		t.Fatal("accepted overlapping paths")
	}
}

func TestResolveCatalogRootSymlinksAndSpecialFiles(t *testing.T) {
	dir := t.TempDir()
	actual := filepath.Join(dir, "actual")
	if err := os.MkdirAll(filepath.Join(actual, "root"), 0o700); err != nil {
		t.Fatal(err)
	}
	alias := filepath.Join(dir, "alias")
	if err := os.Symlink(actual, alias); err != nil {
		t.Skip(err)
	}
	paths, err := ResolveRoot(RootOptions{Directory: filepath.Join(alias, "root")})
	if err != nil || paths.Directory != filepath.Join(actual, "root") {
		t.Fatalf("symlinked ancestor: %#v, %v", paths, err)
	}
	if _, err := ResolveRoot(RootOptions{Directory: alias}); err == nil {
		t.Fatal("accepted symlink at root")
	}
	if err := os.Symlink(dir, filepath.Join(paths.Directory, "link")); err != nil {
		t.Fatal(err)
	}
	if _, err := ResolveRoot(RootOptions{Directory: paths.Directory, RegistryPath: "link/cache.sqlite"}); err == nil {
		t.Fatal("accepted descendant symlink")
	}
	if err := os.Mkdir(filepath.Join(paths.Directory, "cache.sqlite"), 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := ResolveRoot(RootOptions{Directory: paths.Directory}); err == nil {
		t.Fatal("accepted directory as registry")
	}
}

func TestRootOptionsFromEnvironmentPrecedence(t *testing.T) {
	for _, k := range []string{
		"APITOOLS_CATALOG_ROOT", "CATALOG_ROOT", "KINET_CATALOG_ROOT",
		"APITOOLS_CATALOG_REGISTRY", "CATALOG_REGISTRY", "KINET_CATALOG_REGISTRY",
		"APITOOLS_CATALOG_INDEX", "CATALOG_INDEX", "KINET_CATALOG_INDEX",
	} {
		t.Setenv(k, "")
	}

	// 1. All empty -> empty RootOptions
	opts := RootOptionsFromEnvironment()
	if opts != (RootOptions{}) {
		t.Fatalf("expected empty RootOptions, got %#v", opts)
	}

	// 2. KINET_CATALOG_* legacy fallback
	t.Setenv("KINET_CATALOG_ROOT", "/kinet/root")
	t.Setenv("KINET_CATALOG_REGISTRY", "kinet.sqlite")
	t.Setenv("KINET_CATALOG_INDEX", "kinet.json")
	opts = RootOptionsFromEnvironment()
	if opts.Directory != "/kinet/root" || opts.RegistryPath != "kinet.sqlite" || opts.IndexPath != "kinet.json" {
		t.Fatalf("expected kinet fallback, got %#v", opts)
	}

	// 3. CATALOG_* overrides KINET_CATALOG_*
	t.Setenv("CATALOG_ROOT", "/catalog/root")
	t.Setenv("CATALOG_REGISTRY", "catalog.sqlite")
	t.Setenv("CATALOG_INDEX", "catalog.json")
	opts = RootOptionsFromEnvironment()
	if opts.Directory != "/catalog/root" || opts.RegistryPath != "catalog.sqlite" || opts.IndexPath != "catalog.json" {
		t.Fatalf("expected catalog override, got %#v", opts)
	}

	// 4. APITOOLS_CATALOG_* overrides CATALOG_*
	t.Setenv("APITOOLS_CATALOG_ROOT", "/apitools/root")
	t.Setenv("APITOOLS_CATALOG_REGISTRY", "apitools.sqlite")
	t.Setenv("APITOOLS_CATALOG_INDEX", "apitools.json")
	opts = RootOptionsFromEnvironment()
	if opts.Directory != "/apitools/root" || opts.RegistryPath != "apitools.sqlite" || opts.IndexPath != "apitools.json" {
		t.Fatalf("expected apitools override, got %#v", opts)
	}

	// 5. Whitespace-only values fall back to next alias
	t.Setenv("APITOOLS_CATALOG_ROOT", "   ")
	t.Setenv("APITOOLS_CATALOG_REGISTRY", "\t\n ")
	t.Setenv("APITOOLS_CATALOG_INDEX", " ")
	opts = RootOptionsFromEnvironment()
	if opts.Directory != "/catalog/root" || opts.RegistryPath != "catalog.sqlite" || opts.IndexPath != "catalog.json" {
		t.Fatalf("expected whitespace fallback to catalog, got %#v", opts)
	}

	// 6. Explicit configuration in RootOptionsWithEnvironment overrides everything
	explicit := RootOptions{
		Directory:    "/explicit/root",
		RegistryPath: "explicit.sqlite",
		IndexPath:    "explicit.json",
	}
	opts = RootOptionsWithEnvironment(explicit)
	if opts.Directory != "/explicit/root" || opts.RegistryPath != "explicit.sqlite" || opts.IndexPath != "explicit.json" {
		t.Fatalf("expected explicit override, got %#v", opts)
	}
}

func TestResolveDefaultRootReadOnlyAndMissing(t *testing.T) {
	for _, k := range []string{
		"APITOOLS_CATALOG_ROOT", "CATALOG_ROOT", "KINET_CATALOG_ROOT",
		"APITOOLS_CATALOG_REGISTRY", "CATALOG_REGISTRY", "KINET_CATALOG_REGISTRY",
		"APITOOLS_CATALOG_INDEX", "CATALOG_INDEX", "KINET_CATALOG_INDEX",
	} {
		t.Setenv(k, "")
	}

	// 1. Missing root environment returns empty RootPaths without error or file writes
	paths, err := ResolveDefaultRoot()
	if err != nil || paths != (RootPaths{}) {
		t.Fatalf("expected empty RootPaths without error, got %#v, %v", paths, err)
	}

	// 2. Explicit directory in environment resolves cleanly and writes no files
	dir := t.TempDir()
	t.Setenv("APITOOLS_CATALOG_ROOT", dir)
	paths, err = ResolveDefaultRoot()
	if err != nil {
		t.Fatalf("unexpected error resolving default root: %v", err)
	}
	if paths.Directory != dir || paths.RegistryPath != filepath.Join(dir, DefaultRegistryPath) || paths.IndexPath != filepath.Join(dir, DefaultOperationIndexPath) {
		t.Fatalf("unexpected paths resolved: %#v", paths)
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 0 {
		t.Fatalf("resolution wrote files: %v, %v", entries, err)
	}

	// 3. Missing directory returns os.ErrNotExist
	t.Setenv("APITOOLS_CATALOG_ROOT", filepath.Join(dir, "nonexistent"))
	if _, err := ResolveDefaultRoot(); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected ErrNotExist, got %v", err)
	}

	// 4. Registry set without root returns error
	t.Setenv("APITOOLS_CATALOG_ROOT", "")
	t.Setenv("APITOOLS_CATALOG_REGISTRY", "cache.sqlite")
	if _, err := ResolveDefaultRoot(); err == nil {
		t.Fatal("expected error when registry configured without root")
	}
}
