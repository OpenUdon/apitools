package catalogupdate

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestOpenArtifactCachePathOnly(t *testing.T) {
	path := filepath.Join(t.TempDir(), "new", "objects")
	cache, err := OpenArtifactCache(path, "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	if err := cache.WriteFile("artifact", []byte("contents"), 0o600); err != nil {
		t.Fatal(err)
	}
	files, err := os.ReadDir(path)
	if err != nil || len(files) != 1 || files[0].Name() != "artifact" {
		t.Fatalf("cache contents after preflight = %v, %v", files, err)
	}
}

func TestOpenArtifactCacheRejectsMissingOrSubstitutedVolume(t *testing.T) {
	for _, failure := range []error{fs.ErrNotExist, errors.New("substituted volume")} {
		t.Run(failure.Error(), func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "unmounted", "objects")
			cache, err := openArtifactCache(path, func(*os.Root) error { return failure })
			if cache != nil || !errors.Is(err, failure) {
				t.Fatalf("cache = %v, error = %v, want %v", cache, err, failure)
			}
			if _, err := os.Lstat(filepath.Dir(path)); !errors.Is(err, fs.ErrNotExist) {
				t.Fatalf("ancestor created during failed preflight: %v", err)
			}
		})
	}
}

func TestOpenArtifactCacheChecksEachParentBeforeCreatingChildren(t *testing.T) {
	parent := t.TempDir()
	path := filepath.Join(parent, "intermediate", "objects")
	wrongVolume := errors.New("substituted intermediate volume")
	calls := 0
	cache, err := openArtifactCache(path, func(root *os.Root) error {
		calls++
		if calls == 1 {
			// A new directory appears after the existing ancestor is checked.
			return root.Mkdir("intermediate", 0o700)
		}
		return wrongVolume
	})
	if cache != nil || !errors.Is(err, wrongVolume) {
		t.Fatalf("cache = %v, error = %v, want %v", cache, err, wrongVolume)
	}
	if _, err := os.Lstat(path); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("new child in substituted directory: %v", err)
	}
}
