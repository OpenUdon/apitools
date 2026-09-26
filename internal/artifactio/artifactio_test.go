package artifactio

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadFileConfinesPathsAndVerifiesIntegrity(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "source.json"), []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	good := digestBytes([]byte("source"))
	file, err := ReadFile(root, "source.json", ReadOptions{SHA256: good, Bytes: 6})
	if err != nil {
		t.Fatal(err)
	}
	if string(file.Data) != "source" || file.SHA256 != good {
		t.Fatalf("file = %#v", file)
	}
	for _, path := range []string{"../source.json", filepath.Join(root, "source.json"), "."} {
		if _, err := ReadFile(root, path, ReadOptions{}); err == nil || !strings.Contains(err.Error(), "local relative") {
			t.Errorf("path %q: expected confinement error, got %v", path, err)
		}
	}
	if _, err := ReadFile(root, "source.json", ReadOptions{SHA256: strings.Repeat("0", 64)}); err == nil || !strings.Contains(err.Error(), "SHA256") {
		t.Fatalf("expected digest mismatch, got %v", err)
	}
}

func TestReadAndWriteRejectSymlinks(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside")
	if err := os.WriteFile(outside, []byte("outside"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "linked")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if _, err := ReadFile(root, "linked", ReadOptions{}); err == nil || !strings.Contains(err.Error(), "regular file") {
		t.Fatalf("expected symlink read rejection, got %v", err)
	}
	if _, err := WriteFile(root, "linked", []byte("replacement"), WriteOptions{Force: true}); err == nil || !strings.Contains(err.Error(), "regular file") {
		t.Fatalf("expected symlink write rejection, got %v", err)
	}
	content, err := os.ReadFile(outside)
	if err != nil || string(content) != "outside" {
		t.Fatalf("outside content = %q, err %v", content, err)
	}
}

func TestArtifactRootsAndParentsRejectSymlinks(t *testing.T) {
	realRoot := t.TempDir()
	if err := os.WriteFile(filepath.Join(realRoot, "source"), []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	container := t.TempDir()
	linkedRoot := filepath.Join(container, "linked-root")
	if err := os.Symlink(realRoot, linkedRoot); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if _, err := ReadFile(linkedRoot, "source", ReadOptions{}); err == nil || !strings.Contains(err.Error(), "directory component") {
		t.Fatalf("expected symlink root rejection, got %v", err)
	}

	root := t.TempDir()
	outsideDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(outsideDir, "source"), []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outsideDir, filepath.Join(root, "linked-parent")); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadFile(root, "linked-parent/source", ReadOptions{}); err == nil || !strings.Contains(err.Error(), "directory component") {
		t.Fatalf("expected symlink parent rejection, got %v", err)
	}
	tx, err := BeginDir(filepath.Join(root, "linked-parent", "output"), true)
	if err != nil {
		t.Fatalf("expected transaction target to resolve its symlinked ancestor, got %v", err)
	}
	if want := filepath.Join(outsideDir, "output"); tx.target != want {
		t.Fatalf("transaction target = %q, want canonical path %q", tx.target, want)
	}
	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}
}

func TestArtifactRootsResolveSymlinkedAncestors(t *testing.T) {
	realParent := t.TempDir()
	aliasParent := filepath.Join(t.TempDir(), "workspace")
	if err := os.Symlink(realParent, aliasParent); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	rootAlias := filepath.Join(aliasParent, "new", "artifacts")
	root, err := EnsureRoot(rootAlias, 0o700)
	if err != nil {
		t.Fatal(err)
	}
	wantRoot := filepath.Join(realParent, "new", "artifacts")
	if root != wantRoot {
		t.Fatalf("EnsureRoot() = %q, want canonical root %q", root, wantRoot)
	}
	if _, err := WriteFile(rootAlias, "source.json", []byte("source"), WriteOptions{}); err != nil {
		t.Fatal(err)
	}
	file, err := ReadFile(rootAlias, "source.json", ReadOptions{})
	if err != nil || string(file.Data) != "source" {
		t.Fatalf("ReadFile() = %#v, %v", file, err)
	}

	tx, err := BeginDir(filepath.Join(aliasParent, "new", "transaction"), false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Stage("result.txt", []byte("result"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

func TestWriteFileReusesIdenticalAndRequiresForceForDifference(t *testing.T) {
	root := t.TempDir()
	first, err := WriteFile(root, "nested/artifact.json", []byte("one"), WriteOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if first.Reused {
		t.Fatal("new write reported reuse")
	}
	second, err := WriteFile(root, "nested/artifact.json", []byte("one"), WriteOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !second.Reused {
		t.Fatal("identical write was not reused")
	}
	if _, err := WriteFile(root, "nested/artifact.json", []byte("two"), WriteOptions{}); !errors.Is(err, ErrCollision) {
		t.Fatalf("expected collision, got %v", err)
	}
	if _, err := WriteFile(root, "nested/artifact.json", []byte("two"), WriteOptions{Force: true}); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(root, "nested", "artifact.json"))
	if err != nil || string(content) != "two" {
		t.Fatalf("content = %q, err %v", content, err)
	}
}

func TestDirTransactionCollisionForceAndRollback(t *testing.T) {
	parent := t.TempDir()
	target := filepath.Join(parent, "artifacts")
	tx, err := BeginDir(target, false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Stage("one.txt", []byte("one"), 0o600); err != nil {
		t.Fatal(err)
	}
	if reused, err := tx.Commit(); err != nil || reused {
		t.Fatalf("first commit reused=%v err=%v", reused, err)
	}

	tx, err = BeginDir(target, false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Stage("one.txt", []byte("one"), 0o600); err != nil {
		t.Fatal(err)
	}
	if reused, err := tx.Commit(); err != nil || !reused {
		t.Fatalf("identical commit reused=%v err=%v", reused, err)
	}

	tx, err = BeginDir(target, false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Stage("two.txt", []byte("two"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Commit(); !errors.Is(err, ErrCollision) {
		t.Fatalf("expected directory collision, got %v", err)
	}
	if _, err := os.Stat(filepath.Join(target, "one.txt")); err != nil {
		t.Fatalf("collision mutated old tree: %v", err)
	}

	tx, err = BeginDir(target, true)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Stage("two.txt", []byte("two"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(target, "one.txt")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("forced replacement retained obsolete file: %v", err)
	}

	tx, err = BeginDir(target, true)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Stage("three.txt", []byte("three"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(target, "two.txt")); err != nil {
		t.Fatalf("rollback mutated target: %v", err)
	}
}

func TestMoveAsideRestoreDiscardAndRemoveArtifact(t *testing.T) {
	root := t.TempDir()
	relative := "openapi/spec.json"
	full := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	original := []byte("original content")
	if err := os.WriteFile(full, original, 0o644); err != nil {
		t.Fatal(err)
	}

	// Moving aside a missing file is a no-op, not an error.
	backup, err := MoveAside(root, "openapi/missing.json")
	if err != nil || backup != "" {
		t.Fatalf("MoveAside(missing) = (%q, %v), want (\"\", nil)", backup, err)
	}

	backup, err = MoveAside(root, relative)
	if err != nil {
		t.Fatalf("MoveAside: %v", err)
	}
	if backup == "" {
		t.Fatal("MoveAside returned no backup path for an existing file")
	}
	if _, err := os.Stat(full); !os.IsNotExist(err) {
		t.Fatalf("original path still exists after MoveAside, stat err = %v", err)
	}

	// Simulate a refresh writing new content at the same path.
	newContent := []byte("new content")
	if _, err := WriteFile(root, relative, newContent, WriteOptions{}); err != nil {
		t.Fatal(err)
	}

	// RestoreAside discards the new content and brings the original back.
	if err := RestoreAside(root, relative, backup); err != nil {
		t.Fatalf("RestoreAside: %v", err)
	}
	got, err := os.ReadFile(full)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(original) {
		t.Fatalf("restored content = %q, want %q", got, original)
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(backup))); !os.IsNotExist(err) {
		t.Fatalf("backup still exists after RestoreAside, stat err = %v", err)
	}

	// Move aside again and discard: original content is gone for good, and the
	// path written after it (simulating a committed refresh) is untouched.
	backup, err = MoveAside(root, relative)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := WriteFile(root, relative, newContent, WriteOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := DiscardAside(root, backup); err != nil {
		t.Fatalf("DiscardAside: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(backup))); !os.IsNotExist(err) {
		t.Fatalf("backup still exists after DiscardAside, stat err = %v", err)
	}
	// Discarding an already-removed backup is not an error.
	if err := DiscardAside(root, backup); err != nil {
		t.Fatalf("DiscardAside on missing backup: %v", err)
	}
	got, err = os.ReadFile(full)
	if err != nil || string(got) != string(newContent) {
		t.Fatalf("content after discard = %q, %v, want %q", got, err, newContent)
	}

	// RemoveArtifact deletes a file that was never registered anywhere else.
	if err := RemoveArtifact(root, relative); err != nil {
		t.Fatalf("RemoveArtifact: %v", err)
	}
	if _, err := os.Stat(full); !os.IsNotExist(err) {
		t.Fatalf("artifact still exists after RemoveArtifact, stat err = %v", err)
	}
	// Removing an already-missing artifact is not an error.
	if err := RemoveArtifact(root, relative); err != nil {
		t.Fatalf("RemoveArtifact on missing file: %v", err)
	}
}

func TestMoveAsideRejectsSymlinkAndConfinement(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "target.json")
	if err := os.WriteFile(outside, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "link.json")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if _, err := MoveAside(root, "link.json"); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("MoveAside on symlink = %v, want a symlink rejection", err)
	}
	if _, err := MoveAside(root, "../outside.json"); err == nil || !strings.Contains(err.Error(), "local relative") {
		t.Fatalf("MoveAside outside root = %v, want a confinement rejection", err)
	}
}
