package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func createTestDatabase(t *testing.T, path string, wal bool) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	uri := url.URL{Scheme: "file", Path: filepath.ToSlash(path), RawQuery: "mode=rwc"}
	db, err := sql.Open("sqlite", uri.String())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT); INSERT INTO test VALUES (1, 'demo');"); err != nil {
		t.Fatal(err)
	}
	if wal {
		var mode string
		if err := db.QueryRow("PRAGMA journal_mode=WAL").Scan(&mode); err != nil {
			t.Fatal(err)
		}
		// Checkpoint to ensure WAL is flushed to file but header indicates WAL mode
		if _, err := db.Exec("PRAGMA wal_checkpoint(TRUNCATE)"); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSnapshotRootSuccessAndJournalNormalization(t *testing.T) {
	srcDir := t.TempDir()
	targetDir := t.TempDir()

	// 1. Setup source root with nested paths
	regRel := filepath.Join("nested", "cache.sqlite")
	idxRel := filepath.Join("nested", "operations.v1.json")
	srcReg := filepath.Join(srcDir, regRel)
	srcIdx := filepath.Join(srcDir, idxRel)

	createTestDatabase(t, srcReg, true)
	if err := os.WriteFile(srcIdx, []byte(`{"version": 1}`), 0o600); err != nil {
		t.Fatal(err)
	}

	// External metadata
	metaSrc := filepath.Join(t.TempDir(), "meta.json")
	if err := os.WriteFile(metaSrc, []byte(`{"title": "catalog"}`), 0o600); err != nil {
		t.Fatal(err)
	}

	// Selected source inside root
	specRel := filepath.Join("specs", "api.json")
	if err := os.MkdirAll(filepath.Dir(filepath.Join(srcDir, specRel)), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, specRel), []byte(`{"openapi": "3.0.0"}`), 0o600); err != nil {
		t.Fatal(err)
	}

	paths, err := ResolveRoot(RootOptions{
		Directory:    srcDir,
		RegistryPath: regRel,
		IndexPath:    idxRel,
	})
	if err != nil {
		t.Fatal(err)
	}

	opts := SnapshotOptions{
		MetadataSource:      metaSrc,
		MetadataDestination: filepath.Join("catalog-metadata.json"),
		SelectedSourcePath:  specRel,
		NormalizeJournal:    true,
	}

	destRoot := filepath.Join(targetDir, "snapshot")
	if err := SnapshotRoot(paths, destRoot, opts); err != nil {
		t.Fatalf("SnapshotRoot failed: %v", err)
	}

	// Verify all destination files exist
	for _, expectedRel := range []string{regRel, idxRel, specRel, "catalog-metadata.json"} {
		destFile := filepath.Join(destRoot, expectedRel)
		info, err := os.Stat(destFile)
		if err != nil || !info.Mode().IsRegular() {
			t.Fatalf("expected regular file at %s, err: %v", destFile, err)
		}
	}

	// Verify destination SQLite database has journal mode normalized to DELETE
	destDB := filepath.Join(destRoot, regRel)
	uri := url.URL{Scheme: "file", Path: filepath.ToSlash(destDB), RawQuery: "mode=rw"}
	db, err := sql.Open("sqlite", uri.String())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var mode string
	if err := db.QueryRowContext(context.Background(), "PRAGMA journal_mode").Scan(&mode); err != nil {
		t.Fatal(err)
	}
	if mode != "delete" {
		t.Fatalf("expected destination journal mode delete, got %q", mode)
	}

	// Verify producer database remains untouched
	srcURI := url.URL{Scheme: "file", Path: filepath.ToSlash(srcReg), RawQuery: "mode=ro"}
	srcDB, err := sql.Open("sqlite", srcURI.String())
	if err != nil {
		t.Fatal(err)
	}
	defer srcDB.Close()
	var srcMode string
	if err := srcDB.QueryRowContext(context.Background(), "PRAGMA journal_mode").Scan(&srcMode); err != nil {
		t.Fatal(err)
	}
	if srcMode != "wal" {
		t.Fatalf("expected producer journal mode wal, got %q", srcMode)
	}
}

func TestSnapshotRootMissingOptionalEvidence(t *testing.T) {
	srcDir := t.TempDir()
	targetDir := t.TempDir()

	paths, err := ResolveRoot(RootOptions{Directory: srcDir})
	if err != nil {
		t.Fatal(err)
	}

	destRoot := filepath.Join(targetDir, "snapshot")
	if err := SnapshotRoot(paths, destRoot, SnapshotOptions{}); err != nil {
		t.Fatalf("expected success with missing optional evidence, got %v", err)
	}

	// Neither registry nor index should exist in destination
	if _, err := os.Stat(filepath.Join(destRoot, DefaultRegistryPath)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected registry to stay missing, err: %v", err)
	}
	if _, err := os.Stat(filepath.Join(destRoot, DefaultOperationIndexPath)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected index to stay missing, err: %v", err)
	}
}

func TestSnapshotRootMissingRequiredEvidence(t *testing.T) {
	srcDir := t.TempDir()
	targetDir := t.TempDir()

	paths, err := ResolveRoot(RootOptions{Directory: srcDir})
	if err != nil {
		t.Fatal(err)
	}

	destRoot := filepath.Join(targetDir, "snapshot")

	// 1. Missing selected source returns error
	opts1 := SnapshotOptions{
		SelectedSourcePath: "missing.json",
	}
	if err := SnapshotRoot(paths, destRoot, opts1); err == nil {
		t.Fatal("expected error for missing selected source")
	}

	// 2. Missing metadata source returns error
	opts2 := SnapshotOptions{
		MetadataSource:      filepath.Join(srcDir, "nonexistent-meta.json"),
		MetadataDestination: "meta.json",
	}
	if err := SnapshotRoot(paths, destRoot, opts2); err == nil {
		t.Fatal("expected error for missing metadata source")
	}
}

func TestSnapshotRootActiveJournalRefusal(t *testing.T) {
	srcDir := t.TempDir()
	targetDir := t.TempDir()

	regPath := filepath.Join(srcDir, DefaultRegistryPath)
	createTestDatabase(t, regPath, false)

	// Create fake active WAL file
	walPath := regPath + "-wal"
	if err := os.WriteFile(walPath, []byte("active wal contents"), 0o600); err != nil {
		t.Fatal(err)
	}

	paths, err := ResolveRoot(RootOptions{Directory: srcDir})
	if err != nil {
		t.Fatal(err)
	}

	destRoot := filepath.Join(targetDir, "snapshot")
	err = SnapshotRoot(paths, destRoot, SnapshotOptions{})
	if !errors.Is(err, ErrActiveJournal) {
		t.Fatalf("expected ErrActiveJournal, got %v", err)
	}
}

func TestSnapshotRootCollisionAndOverlapDefenses(t *testing.T) {
	srcDir := t.TempDir()
	targetDir := t.TempDir()

	paths, err := ResolveRoot(RootOptions{Directory: srcDir})
	if err != nil {
		t.Fatal(err)
	}

	destRoot := filepath.Join(targetDir, "snapshot")

	// 1. Root overlap (target inside source)
	if err := SnapshotRoot(paths, filepath.Join(srcDir, "child"), SnapshotOptions{}); !errors.Is(err, ErrRootOverlap) {
		t.Fatalf("expected ErrRootOverlap for child destination, got %v", err)
	}

	// 2. Metadata destination collides with index path
	opts1 := SnapshotOptions{
		MetadataSource:      filepath.Join(srcDir, "any.json"),
		MetadataDestination: DefaultOperationIndexPath,
	}
	if err := SnapshotRoot(paths, destRoot, opts1); err == nil {
		t.Fatal("expected error for metadata colliding with index")
	}

	// 3. Metadata destination overlaps registry sidecar
	opts2 := SnapshotOptions{
		MetadataSource:      filepath.Join(srcDir, "any.json"),
		MetadataDestination: DefaultRegistryPath + "-wal",
	}
	if err := SnapshotRoot(paths, destRoot, opts2); err == nil {
		t.Fatal("expected error for metadata overlapping registry sidecar")
	}
}

func TestSnapshotRootCeilingEnforcement(t *testing.T) {
	srcDir := t.TempDir()
	targetDir := t.TempDir()

	regPath := filepath.Join(srcDir, DefaultRegistryPath)
	if err := os.WriteFile(regPath, make([]byte, 1024), 0o600); err != nil {
		t.Fatal(err)
	}

	paths, err := ResolveRoot(RootOptions{Directory: srcDir})
	if err != nil {
		t.Fatal(err)
	}

	destRoot := filepath.Join(targetDir, "snapshot")

	// Limit registry to 512 bytes
	err = SnapshotRoot(paths, destRoot, SnapshotOptions{
		MaxRegistryBytes: 512,
	})
	if err == nil {
		t.Fatal("expected ceiling error when registry exceeds limit")
	}
}

func TestSnapshotRootDriftDefense(t *testing.T) {
	srcDir := t.TempDir()
	targetDir := t.TempDir()

	regPath := filepath.Join(srcDir, DefaultRegistryPath)
	createTestDatabase(t, regPath, false)
	idxPath := filepath.Join(srcDir, DefaultOperationIndexPath)
	if err := os.WriteFile(idxPath, []byte(`{"version": 1}`), 0o600); err != nil {
		t.Fatal(err)
	}

	paths, err := ResolveRoot(RootOptions{Directory: srcDir})
	if err != nil {
		t.Fatal(err)
	}

	destRoot := filepath.Join(targetDir, "snapshot")
	if err := SnapshotRoot(paths, destRoot, SnapshotOptions{}); err != nil {
		t.Fatalf("expected snapshot to succeed, got %v", err)
	}

	// Verify active journal appearing causes refusal
	walPath := regPath + "-wal"
	if err := os.WriteFile(walPath, []byte("active wal contents"), 0o600); err != nil {
		t.Fatal(err)
	}
	destRoot2 := filepath.Join(targetDir, "snapshot2")
	err = SnapshotRoot(paths, destRoot2, SnapshotOptions{})
	if !errors.Is(err, ErrActiveJournal) {
		t.Fatalf("expected ErrActiveJournal, got %v", err)
	}
}

func TestSnapshotPreservesExistingTargets(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			src, target := t.TempDir(), t.TempDir()
			os.WriteFile(filepath.Join(src, DefaultRegistryPath), []byte("new"), 0600)
			old := filepath.Join(target, DefaultRegistryPath)
			os.WriteFile(old, []byte("old"), 0600)
			paths, err := ResolveRoot(RootOptions{Directory: src})
			if err != nil {
				t.Fatal(err)
			}
			opts := SnapshotOptions{}
			if missing {
				opts.SelectedSourcePath = "missing.json"
			}
			if err := SnapshotRoot(paths, target, opts); err == nil {
				t.Fatal("accepted differing target")
			}
			got, err := os.ReadFile(old)
			if err != nil || string(got) != "old" {
				t.Fatalf("lost existing data: %q %v", got, err)
			}
		})
	}
}

func TestSnapshotRejectsDestinationAndSourceSymlinks(t *testing.T) {
	for _, source := range []bool{false, true} {
		t.Run(fmt.Sprint(source), func(t *testing.T) {
			src, target, outside := t.TempDir(), t.TempDir(), t.TempDir()
			os.WriteFile(filepath.Join(outside, "api.json"), []byte("original"), 0600)
			os.Mkdir(filepath.Join(src, "specs"), 0700)
			os.WriteFile(filepath.Join(src, "specs", "api.json"), []byte("replacement"), 0600)
			if source {
				os.RemoveAll(filepath.Join(src, "specs"))
				os.Symlink(outside, filepath.Join(src, "specs"))
			} else {
				os.Symlink(outside, filepath.Join(target, "specs"))
			}
			paths, err := ResolveRoot(RootOptions{Directory: src})
			if err != nil {
				t.Fatal(err)
			}
			if err := SnapshotRoot(paths, target, SnapshotOptions{SelectedSourcePath: "specs/api.json"}); err == nil {
				t.Fatal("accepted symlink")
			}
			got, _ := os.ReadFile(filepath.Join(outside, "api.json"))
			if string(got) != "original" {
				t.Fatalf("external file changed: %q", got)
			}
		})
	}
}

func TestSnapshotEmptyTargetAndIdenticalReuse(t *testing.T) {
	src, target := t.TempDir(), t.TempDir()
	os.WriteFile(filepath.Join(src, DefaultOperationIndexPath), []byte("{}"), 0600)
	paths, err := ResolveRoot(RootOptions{Directory: src})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := SnapshotRoot(paths, target, SnapshotOptions{}); err != nil {
			t.Fatal(err)
		}
	}
	os.WriteFile(filepath.Join(target, "extra"), []byte("keep"), 0600)
	if err := SnapshotRoot(paths, target, SnapshotOptions{}); err == nil {
		t.Fatal("accepted differing nonempty target")
	}
}
