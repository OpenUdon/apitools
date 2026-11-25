//go:build linux

package catalogupdate

import (
	"bytes"
	"database/sql"
	"errors"
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
