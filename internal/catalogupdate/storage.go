package catalogupdate

import (
	"context"
	"crypto/rand"
	"database/sql"
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
	if err := checkVolumeUUID(expectedUUID); err != nil {
		return nil, err
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return openArtifactCache(absolute, func(root *os.Root) error {
		return checkArtifactFilesystem(root, expectedUUID, requiredBytes)
	})
}

func checkVolumeUUID(uuid string) error {
	if uuid == "." || uuid == ".." || strings.ContainsAny(uuid, "/\\") || strings.TrimSpace(uuid) != uuid {
		return errors.New("invalid artifact cache filesystem UUID")
	}
	return nil
}

// CheckMaintenanceStorage checks existing directories and a state file, if
// present. Temporary writability probes are removed before returning.
func CheckMaintenanceStorage(ctx context.Context, cachePath, statePath, expectedUUID string) error {
	if !filepath.IsAbs(cachePath) || !filepath.IsAbs(statePath) {
		return errors.New("catalog maintenance needs absolute cache_dir and state_db paths")
	}
	if err := checkVolumeUUID(expectedUUID); err != nil {
		return err
	}
	cache, err := os.OpenRoot(cachePath)
	if err != nil {
		return err
	}
	defer cache.Close()
	if err := checkArtifactFilesystem(cache, expectedUUID, 0); err != nil {
		return err
	}
	if err := probeWritableRoot(cache); err != nil {
		return err
	}
	stateDir, err := os.OpenRoot(filepath.Dir(statePath))
	if err != nil {
		return err
	}
	defer stateDir.Close()
	directory, err := stateDir.Open(".")
	if err != nil {
		return err
	}
	if err := errors.Join(checkStateFilesystem(directory), directory.Close()); err != nil {
		return err
	}
	if err := probeWritableRoot(stateDir); err != nil {
		return err
	}
	name := filepath.Base(statePath)
	info, err := stateDir.Lstat(name)
	if errors.Is(err, fs.ErrNotExist) {
		return nil // The normal updater initializes the database on the first run.
	}
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return errors.New("catalog state path must be a regular file")
	}
	file, err := openStatePath(stateDir, name)
	if err != nil {
		return err
	}
	defer file.Close()
	if opened, err := file.Stat(); err != nil || !os.SameFile(info, opened) {
		return errors.New("catalog state file changed during preflight")
	}
	if err := checkStateFilesystem(file); err != nil {
		return err
	}
	if err := checkStatePath(statePath, info); err != nil {
		return err
	}
	if err := checkStateWritable(statePath); err != nil {
		return err
	}
	db, err := openReadOnlyState(statePath)
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = checkCatalogStateSchema(ctx, tx)
	return err
}

// LockMonthlyRun serializes maintenance jobs sharing a state database.
// Artifact publication retains its separate short-lived object lock.
func LockMonthlyRun(ctx context.Context, statePath string) (*os.File, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if !filepath.IsAbs(statePath) {
		return nil, errors.New("catalog maintenance needs an absolute state_db path")
	}
	root, err := os.OpenRoot(filepath.Dir(statePath))
	if err != nil {
		return nil, err
	}
	defer root.Close()
	dir, err := root.Open(".")
	if err != nil {
		return nil, err
	}
	if err := errors.Join(checkStateFilesystem(dir), dir.Close()); err != nil {
		return nil, err
	}
	return lockCatalogFile(ctx, root, filepath.Base(statePath)+".monthly.lock", false)
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

	if err := probeWritableRoot(root); err != nil {
		return nil, err
	}
	cache := root
	root = nil
	return cache, nil
}

func probeWritableRoot(root *os.Root) error {
	name := ".apitools-write-check-" + rand.Text()
	probe, err := root.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("test catalog directory writability: %w", err)
	}
	if err := errors.Join(probe.Close(), root.Remove(name)); err != nil {
		return fmt.Errorf("finish catalog directory write test: %w", err)
	}
	return nil
}

func syncDirectory(root *os.Root) error {
	dir, err := root.Open(".")
	if err != nil {
		return err
	}
	return errors.Join(dir.Sync(), dir.Close())
}
