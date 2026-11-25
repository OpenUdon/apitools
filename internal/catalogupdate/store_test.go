//go:build linux

package catalogupdate

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func testObject(content string) Object {
	return Object{SHA256: fmt.Sprintf("%x", sha256.Sum256([]byte(content))), Bytes: int64(len(content))}
}

func openTestStore(ctx context.Context, path string, cache *os.Root) (*Store, error) {
	// A disposable fixture may use the system temp directory on tmpfs.
	return openStore(ctx, path, cache, func(*os.File) error { return nil })
}

func TestStoreRecoversStageAndReusesPublishedOrphan(t *testing.T) {
	ctx := t.Context()
	dir := t.TempDir()
	cache, err := OpenArtifactCache(filepath.Join(dir, "cache"), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	state := filepath.Join(dir, "catalog.sqlite")
	store, err := openTestStore(ctx, state, cache)
	if err != nil {
		t.Fatal(err)
	}
	lock, err := lockCatalogCache(ctx, cache)
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Close()
	const stage = ".catalog-update-stage"
	if err := cache.WriteFile(stage, []byte("part"), 0o600); err != nil {
		t.Fatal(err)
	}
	artifact := testObject("complete")
	wait, cancel := context.WithTimeout(ctx, 30*time.Millisecond)
	defer cancel()
	if err := store.PutObject(wait, artifact, strings.NewReader("complete")); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("competing writer = %v; want deadline", err)
	}
	if got, err := cache.ReadFile(stage); err != nil || string(got) != "part" {
		t.Fatalf("active stage modified: %q, %v", got, err)
	}
	if err := lock.Close(); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	if err := store.PutObject(ctx, artifact, strings.NewReader("complete")); err == nil {
		t.Fatal("published object registered despite closed database")
	}
	if _, err := cache.Lstat(stage); !os.IsNotExist(err) {
		t.Fatalf("abandoned stage still exists: %v", err)
	}
	store, err = openTestStore(ctx, state, cache)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	if err := store.PutObject(ctx, artifact, strings.NewReader("")); err != nil {
		t.Fatalf("reuse published file during retry: %v", err)
	}
}

func TestStateFilesystemRejectsVolatileAndExternalFilesystems(t *testing.T) {
	for _, kind := range []uint32{0x01021994, 0x7366746e, 0x2011bab0, 0x4d44, 0xfe534d42, 0x01021997} {
		if err := rejectStateFilesystem(kind); err == nil {
			t.Fatalf("filesystem %#x admitted for durable state", kind)
		}
	}
	if err := rejectStateFilesystem(0xef53); err != nil {
		t.Fatalf("ext4 rejected: %v", err)
	}
}

func TestStoreLockRefusesSymlinkWithinCache(t *testing.T) {
	cache, err := OpenArtifactCache(t.TempDir(), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	if err := cache.WriteFile("ordinary-object", []byte("fixture"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := cache.Symlink("ordinary-object", ".catalog-update.lock"); err != nil {
		t.Fatal(err)
	}
	if file, err := lockCatalogCache(t.Context(), cache); err == nil {
		_ = file.Close()
		t.Fatal("acquired a lock through a symlink")
	}
}

func TestStoreRefusesPathChangedAfterFilesystemCheck(t *testing.T) {
	root := t.TempDir()
	cache, err := OpenArtifactCache(filepath.Join(root, "cache"), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	state := filepath.Join(root, "state.sqlite")
	store, err := openStore(t.Context(), state, cache, func(file *os.File) error {
		if file.Name() != "state.sqlite" {
			return nil
		}
		if err := os.Rename(state, state+".held"); err != nil {
			return err
		}
		return os.WriteFile(state, []byte("replacement"), 0o600)
	})
	if store != nil || err == nil || !strings.Contains(err.Error(), "state path changed") {
		t.Fatalf("store = %v, error = %v", store, err)
	}
	if data, err := os.ReadFile(state); err != nil || string(data) != "replacement" {
		t.Fatalf("replacement changed: %q, %v", data, err)
	}
}

func TestStoreStatePathCheckPreservesAnotherConnectionWriteLock(t *testing.T) {
	ctx := t.Context()
	root := t.TempDir()
	cache, err := OpenArtifactCache(filepath.Join(root, "cache"), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	state := filepath.Join(root, "state.sqlite")
	store, err := openTestStore(ctx, state, cache)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	observer, err := os.Open(state)
	if err != nil {
		t.Fatal(err)
	}
	defer observer.Close()
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	assertWriteLock := func() {
		t.Helper()
		// An OFD lock request sees conflicts with this process's POSIX locks.
		lock := unix.Flock_t{Type: unix.F_RDLCK}
		if err := unix.FcntlFlock(observer.Fd(), unix.F_OFD_GETLK, &lock); err != nil {
			t.Fatal(err)
		}
		if lock.Type != unix.F_WRLCK {
			t.Fatalf("SQLite write lock missing: %+v", lock)
		}
	}
	assertWriteLock()
	stop := errors.New("stop after inspecting the state file")
	second, err := openStore(ctx, state, cache, func(file *os.File) error {
		if file.Name() == filepath.Base(state) {
			return stop
		}
		return nil
	})
	if !errors.Is(err, stop) {
		if second != nil {
			_ = second.Close()
		}
		t.Fatalf("second open: %v", err)
	}
	assertWriteLock()
}

func TestStoreKeepsOverlappingRecordOutOfRollingBackTransaction(t *testing.T) {
	ctx := t.Context()
	root := t.TempDir()
	cache, err := OpenArtifactCache(filepath.Join(root, "cache"), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	store, err := openTestStore(ctx, filepath.Join(root, "db.sqlite"), cache)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	run, err := store.StartRun(ctx, testObject("catalog revision").SHA256, []string{"a"})
	if err != nil {
		t.Fatal(err)
	}
	store.mu.Lock()
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		store.mu.Unlock()
		t.Fatal(err)
	}
	started := make(chan struct{})
	recorded := make(chan error, 1)
	go func() {
		close(started)
		recorded <- store.Record(ctx, run.ID, Outcome{ProviderID: "a", Status: "checked"})
	}()
	<-started
	select {
	case err := <-recorded:
		_ = tx.Rollback()
		store.mu.Unlock()
		t.Fatalf("record ended before the transaction rolled back: %v", err)
	case <-time.After(20 * time.Millisecond):
	}
	err = tx.Rollback()
	store.mu.Unlock()
	if err != nil {
		t.Fatal(err)
	}
	if err := <-recorded; err != nil {
		t.Fatal(err)
	}
	if rows, err := store.ReadOutcomes(ctx, run.ID); err != nil || len(rows) != 1 {
		t.Fatalf("record lost during rollback: %+v, %v", rows, err)
	}
}

func TestStoreResumesRunAndDeduplicatesObjects(t *testing.T) {
	ctx := t.Context()
	root := t.TempDir()
	dir := filepath.Join(root, "cache")
	cache, err := OpenArtifactCache(dir, "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	state := filepath.Join(root, "state", "catalog.sqlite")
	store, err := openTestStore(ctx, state, cache)
	if err != nil {
		t.Fatal(err)
	}
	run, err := store.StartRun(ctx, testObject("catalog revision").SHA256, []string{"a", "b"})
	if err != nil {
		t.Fatal(err)
	}
	artifact := testObject("shared publisher bytes")
	if err := store.PutObject(ctx, artifact, strings.NewReader("shared publisher bytes")); err != nil {
		t.Fatal(err)
	}
	if err := store.Record(ctx, run.ID, Outcome{ProviderID: "a", ReferenceID: "publisher", Status: "validated", Object: &artifact}); err != nil {
		t.Fatal(err)
	}
	if err := store.Record(ctx, run.ID, Outcome{ProviderID: "a", Status: "checked"}); err != nil {
		t.Fatal(err)
	}
	if err := store.FinishRun(ctx, run.ID); err == nil {
		t.Fatal("completed a run missing provider b's checkpoint")
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	store, err = openTestStore(ctx, state, cache)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	got, err := store.ReadRun(ctx, run.ID)
	if err != nil || got.Completed || got.Catalog != run.Catalog || strings.Join(got.Providers, ",") != "a,b" {
		t.Fatalf("resumed run = %+v, %v", got, err)
	}
	if err := store.PutObject(ctx, artifact, strings.NewReader("")); err != nil {
		t.Fatalf("artifact should be reused without a body: %v", err)
	}
	if err := store.Record(ctx, run.ID, Outcome{ProviderID: "outside", Status: "checked"}); err == nil {
		t.Fatal("recorded a provider outside the run scope")
	}
	if err := store.Record(ctx, run.ID, Outcome{ProviderID: "b", Status: "checked", Object: &artifact}); err != nil {
		t.Fatal(err)
	}
	if err := store.FinishRun(ctx, run.ID); err != nil {
		t.Fatal(err)
	}
	rows, err := store.ReadOutcomes(ctx, run.ID)
	if err != nil || len(rows) != 3 || rows[1].Object == nil || rows[2].Object == nil ||
		*rows[1].Object != artifact || *rows[2].Object != artifact {
		t.Fatalf("run outcomes = %+v, %v", rows, err)
	}
	if err := store.Record(ctx, run.ID, Outcome{ProviderID: "a", Status: "replaced"}); err == nil {
		t.Fatal("modified a completed run")
	}
	items, err := os.ReadDir(dir)
	if err != nil || len(items) != 2 {
		t.Fatalf("artifact files = %v, %v", items, err)
	}
	file, err := store.OpenObject(ctx, artifact)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if body, err := io.ReadAll(file); err != nil || string(body) != "shared publisher bytes" {
		t.Fatalf("offline artifact = %q, %v", body, err)
	}
}

func TestStoreRejectsTruncatedOversizedAndSubstitutedObjects(t *testing.T) {
	ctx := t.Context()
	root := t.TempDir()
	dir := filepath.Join(root, "cache")
	cache, err := OpenArtifactCache(dir, "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	store, err := openTestStore(ctx, filepath.Join(root, "db"), cache)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	want := testObject("right")
	for _, body := range []string{"righ", "right!", "wrong"} {
		if err := store.PutObject(ctx, want, strings.NewReader(body)); err == nil {
			t.Fatalf("stored mismatched bytes %q", body)
		}
	}
	items, err := os.ReadDir(dir)
	if err != nil || len(items) != 1 || items[0].Name() != ".catalog-update.lock" {
		t.Fatalf("failed bodies left cache entries %v, %v", items, err)
	}
	name, err := want.filename()
	if err != nil {
		t.Fatal(err)
	}
	if err := cache.WriteFile(name, []byte("wrong"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := store.OpenObject(ctx, want); !errors.Is(err, ErrObjectIntegrity) {
		t.Fatal("read a corrupt cached artifact")
	}
	if err := store.PutObject(ctx, want, strings.NewReader("wrong")); !errors.Is(err, ErrObjectIntegrity) {
		t.Fatalf("unverified replacement succeeded: %v", err)
	}
	if err := store.PutObject(ctx, want, strings.NewReader("right")); err != nil {
		t.Fatalf("restore digest file with verified bytes: %v", err)
	}
	file, err := store.OpenObject(ctx, want)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if observed, err := io.ReadAll(file); err != nil || string(observed) != "right" {
		t.Fatalf("restored artifact = %q, %v", observed, err)
	}
}
