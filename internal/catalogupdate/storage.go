package catalogupdate

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// OpenArtifactCache prepares the configured directory for artifact writes.
// The returned root pins the prepared directory and the caller must close it.
// An expected filesystem UUID is optional. requiredBytes estimates the
// space needed for the next artifact write.
func OpenArtifactCache(path, expectedUUID string, requiredBytes uint64) (*os.Root, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("artifact cache directory is required")
	}
	if expectedUUID != "" && (expectedUUID == "." || expectedUUID == ".." ||
		strings.ContainsAny(expectedUUID, "/\\") || strings.TrimSpace(expectedUUID) != expectedUUID) {
		return nil, errors.New("invalid artifact cache filesystem UUID")
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return openArtifactCache(absolute, func(root *os.Root) error {
		return checkArtifactFilesystem(root, expectedUUID, requiredBytes)
	})
}

// Check the deepest existing ancestor before creating a missing directory.
// Creating through the opened root prevents an unmount from moving a write
// onto the filesystem beneath the mount point.
func openArtifactCache(path string, check func(*os.Root) error) (*os.Root, error) {
	parent := path
	var root *os.Root
	for {
		var err error
		root, err = os.OpenRoot(parent)
		if err == nil {
			break
		}
		if !errors.Is(err, fs.ErrNotExist) || parent == filepath.Dir(parent) {
			return nil, fmt.Errorf("open artifact cache directory %q: %w", parent, err)
		}
		parent = filepath.Dir(parent)
	}
	defer func() {
		if root != nil {
			_ = root.Close()
		}
	}()
	if err := check(root); err != nil {
		return nil, fmt.Errorf("check artifact cache %q: %w", path, err)
	}
	if parent != path {
		relative, err := filepath.Rel(parent, path)
		if err != nil {
			return nil, err
		}
		for _, component := range strings.Split(relative, string(filepath.Separator)) {
			err := root.Mkdir(component, 0o700)
			if err != nil && !errors.Is(err, fs.ErrExist) {
				return nil, err
			}
			if err == nil {
				if err := syncDirectory(root); err != nil {
					return nil, err
				}
			}
			nested, err := root.OpenRoot(component)
			if err != nil {
				return nil, err
			}
			_ = root.Close()
			root = nested
			if err := check(root); err != nil {
				return nil, fmt.Errorf("check artifact cache %q: %w", path, err)
			}
		}
	}

	name := ".apitools-write-check-" + rand.Text()
	probe, err := root.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("test artifact cache writability: %w", err)
	}
	if err := errors.Join(probe.Close(), root.Remove(name)); err != nil {
		return nil, fmt.Errorf("finish artifact cache write test: %w", err)
	}
	cache := root
	root = nil
	return cache, nil
}

func syncDirectory(root *os.Root) error {
	dir, err := root.Open(".")
	if err != nil {
		return err
	}
	return errors.Join(dir.Sync(), dir.Close())
}
