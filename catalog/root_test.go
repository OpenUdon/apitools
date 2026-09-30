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
