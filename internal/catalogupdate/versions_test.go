//go:build linux

package catalogupdate

import (
	"bytes"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestStoreAdvancesApprovedVersionAndPreservesPreviousBytes(t *testing.T) {
	ctx := t.Context()
	dir := t.TempDir()
	cache, err := OpenArtifactCache(filepath.Join(dir, "cache"), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	state := filepath.Join(dir, "catalog.sqlite")
	legacy, err := sql.Open("sqlite", state)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := legacy.ExecContext(ctx, schemaSQL); err != nil {
		t.Fatal(err)
	}
	if err := legacy.Close(); err != nil {
		t.Fatal(err)
	}
	store, err := openTestStore(ctx, state, cache)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	previous := Version{
		ProviderID: "example", ReferenceID: "records", Service: "records", Version: "v1",
		Locator: "https://example.com/records.json", URL: "https://example.com/v1/records.json",
		CheckedAt: time.Date(2026, 11, 25, 0, 0, 0, 0, time.UTC),
		Object:    testObject("previous spec contents"),
	}
	if err := store.PutObject(ctx, previous.Object, strings.NewReader("previous spec contents")); err != nil {
		t.Fatal(err)
	}
	if err := store.AcceptVersion(ctx, previous, nil); err != nil {
		t.Fatal(err)
	}
	next := previous
	next.Version = "v2"
	next.URL = "https://example.com/v2/records.json"
	next.Object = testObject("new spec contents")
	if err := store.PutObject(ctx, next.Object, strings.NewReader("new spec contents")); err != nil {
		t.Fatal(err)
	}
	if err := store.AcceptVersion(ctx, next, nil); !errors.Is(err, ErrPreferenceChanged) {
		t.Fatalf("stale update = %v", err)
	}
	if held, err := store.PreferredVersion(ctx, "example", "records", "records"); err != nil || held.Object != previous.Object {
		t.Fatalf("previous version = %+v, %v", held, err)
	}
	if err := store.AcceptVersion(ctx, next, &previous); err != nil {
		t.Fatal(err)
	}
	if preferred, err := store.PreferredVersion(ctx, "example", "records", "records"); err != nil || preferred.Object != next.Object {
		t.Fatalf("preferred version = %+v, %v", preferred, err)
	}
	pin, err := store.OpenObject(ctx, previous.Object)
	if err != nil {
		t.Fatalf("previous pin unavailable: %v", err)
	}
	pin.Close()
	sameBytes := next
	sameBytes.Version = "v3"
	sameBytes.URL = "https://example.com/v3/records.json"
	sameBytes.CheckedAt = next.CheckedAt.Add(time.Hour)
	if err := store.AcceptVersion(ctx, sameBytes, &next); err != nil {
		t.Fatal(err)
	}
	if preferred, err := store.PreferredVersion(ctx, "example", "records", "records"); err != nil || !samePreferredVersion(preferred, sameBytes) {
		t.Fatalf("same bytes retain previous version metadata: %+v, %v", preferred, err)
	}
	otherService := sameBytes
	otherService.Service = "another service"
	if err := store.AcceptVersion(ctx, otherService, nil); err != nil {
		t.Fatal(err)
	}
	if preferred, err := store.PreferredVersion(ctx, "example", "records", "records"); err != nil || !samePreferredVersion(preferred, sameBytes) {
		t.Fatalf("another index member changed the first preference: %+v, %v", preferred, err)
	}
	if _, err := store.ReviewedReferencePreference(ctx, "example", "records"); !errors.Is(err, ErrAmbiguousReferencePreference) {
		t.Fatalf("relocated service was guessed across several index members: %v", err)
	}
	first, err := store.ExportVersions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.ExportVersions(ctx)
	if err != nil || !bytes.Equal(first, second) || bytes.Contains(first, []byte("new spec contents")) ||
		!bytes.Contains(first, []byte(sameBytes.URL)) {
		t.Fatalf("metadata export = %s / %s / %v", first, second, err)
	}
}

func TestStoreRollsBackVersionWhenReferenceCheckpointFails(t *testing.T) {
	ctx := t.Context()
	dir := t.TempDir()
	cache, err := OpenArtifactCache(filepath.Join(dir, "cache"), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	store, err := openTestStore(ctx, filepath.Join(dir, "catalog.sqlite"), cache)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	pin := Version{ProviderID: "fixture", ReferenceID: "example", Service: "service", Locator: "https://example.com", URL: "https://example.com", CheckedAt: time.Now(), Object: testObject("bytes")}
	if err := store.PutObject(ctx, pin.Object, strings.NewReader("bytes")); err != nil {
		t.Fatal(err)
	}
	checkpoint := Outcome{ProviderID: pin.ProviderID, ReferenceID: pin.ReferenceID, Status: "reviewed", Object: &pin.Object}
	if err := store.AcceptDirectVersionAndRecord(ctx, pin, nil, "missing-run", checkpoint); err == nil {
		t.Fatal("accepted a version without committing its checkpoint")
	}
	if _, err := store.PreferredVersion(ctx, pin.ProviderID, pin.ReferenceID, pin.Service); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("preference leaked after failed checkpoint: %v", err)
	}
	run, err := store.StartRun(ctx, testObject("catalog").SHA256, []string{pin.ProviderID})
	if err != nil {
		t.Fatal(err)
	}
	first := Outcome{ProviderID: pin.ProviderID, ReferenceID: pin.ReferenceID, Status: "review-required"}
	if err := store.Record(ctx, run.ID, first); err != nil {
		t.Fatal(err)
	}
	if err := store.AcceptDirectVersionAndRecord(ctx, pin, nil, run.ID, checkpoint); err == nil {
		t.Fatal("accepted version after another writer finalized its checkpoint")
	}
	if _, err := store.PreferredVersion(ctx, pin.ProviderID, pin.ReferenceID, pin.Service); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("preference leaked after checkpoint collision: %v", err)
	}
	rows, err := store.ReadOutcomes(ctx, run.ID)
	if err != nil || len(rows) != 1 || rows[0].Status != first.Status {
		t.Fatalf("checkpoint overwritten by a second writer: %+v, %v", rows, err)
	}
	competing, err := store.StartRun(ctx, testObject("revised catalog").SHA256, []string{pin.ProviderID})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.AcceptVersion(ctx, pin, nil); err != nil {
		t.Fatal(err)
	}
	newService := pin
	newService.Service = "moved service"
	if err := store.AcceptDirectVersionAndRecord(ctx, newService, nil, competing.ID, checkpoint); !errors.Is(err, ErrPreferenceChanged) {
		t.Fatalf("a stale first acceptance gained a second direct-source preference: %v", err)
	}
	if rows, err := store.ReadOutcomes(ctx, competing.ID); err != nil || len(rows) != 0 {
		t.Fatalf("checkpoint leaked after a competing preference: %+v, %v", rows, err)
	}
}

func TestStoreMigratesVersionTwoPreferences(t *testing.T) {
	ctx := t.Context()
	dir := t.TempDir()
	cache, err := OpenArtifactCache(filepath.Join(dir, "cache"), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	path := filepath.Join(dir, "catalog.sqlite")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	previousSchema := strings.Replace(versionSchemaSQL, "  catalog_pin TEXT NOT NULL DEFAULT '',\n", "", 1)
	previousSchema = strings.Replace(previousSchema, "user_version = 3", "user_version = 2", 1)
	if _, err := db.ExecContext(ctx, schemaSQL+previousSchema); err != nil {
		t.Fatal(err)
	}
	object := testObject("earlier version")
	if _, err := db.ExecContext(ctx, "INSERT INTO update_object VALUES (?, ?)", object.SHA256, object.Bytes); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `
INSERT INTO update_version VALUES (1, 'example', 'records', 'records', ?,
 'https://example.com/spec.json', 'https://example.com/spec.json', '1.0', '', '2026-11-25T00:00:00Z');
INSERT INTO update_preferred VALUES ('example', 'records', 'records', 1);`, object.SHA256); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	export, err := ReadVersionExport(ctx, path)
	if err != nil || !bytes.Contains(export, []byte(object.SHA256)) || bytes.Contains(export, []byte("catalog_pin")) {
		t.Fatalf("read-only v2 export = %s, %v", export, err)
	}
	if after, err := os.ReadFile(path); err != nil || !bytes.Equal(before, after) {
		t.Fatalf("read-only v2 export changed the database: %v", err)
	}
	store, err := openTestStore(ctx, path, cache)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	held, err := store.PreferredVersion(ctx, "example", "records", "records")
	if err != nil || held.Object != object || held.Version != "1.0" || held.CatalogPin != "" {
		t.Fatalf("migrated version = %+v, %v", held, err)
	}
}
