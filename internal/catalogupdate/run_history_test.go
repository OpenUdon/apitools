package catalogupdate

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadLegacyCompletedCatalogRunAndExportWithoutMigration(t *testing.T) {
	for _, version := range []int{2, 3} {
		t.Run(fmt.Sprintf("schema-v%d", version), func(t *testing.T) {
			ctx := t.Context()
			dir := t.TempDir()
			state := filepath.Join(dir, "maintenance.sqlite")
			db, err := sql.Open("sqlite", state)
			if err != nil {
				t.Fatal(err)
			}
			schemas := schemaSQL + versionSchemaSQL
			if version == 2 {
				schemas = strings.Replace(schemas, "  catalog_pin TEXT NOT NULL DEFAULT '',\n", "", 1)
				schemas = strings.Replace(schemas, "user_version = 3", "user_version = 2", 1)
			}
			if _, err := db.ExecContext(ctx, schemas); err != nil {
				t.Fatal(err)
			}
			identity := testObject("previous catalog").SHA256
			if _, err := db.ExecContext(ctx, "INSERT INTO update_run VALUES ('old-run', ?, 1)", identity); err != nil {
				t.Fatal(err)
			}
			if _, err := db.ExecContext(ctx, `
INSERT INTO update_provider VALUES ('old-run','fixture',0);
INSERT INTO update_outcome VALUES ('old-run','fixture','', 'checked', NULL, NULL);`); err != nil {
				t.Fatal(err)
			}
			if err := db.Close(); err != nil {
				t.Fatal(err)
			}
			verify := func() {
				t.Helper()
				run, outcomes, err := ReadExistingRun(ctx, state, "old-run")
				if err != nil || !run.Completed || run.Catalog != identity || run.CatalogSHA256 != "" ||
					!run.StartedAt.IsZero() || !run.CompletedAt.IsZero() || len(outcomes) != 1 {
					t.Fatalf("legacy run = %+v, outcomes=%d, error=%v", run, len(outcomes), err)
				}
			}
			verify()
			before, err := ReadVersionExport(ctx, state)
			if err != nil {
				t.Fatal(err)
			}
			if _, _, err := ReadExistingRun(ctx, state, "unknown-id"); !errors.Is(err, ErrRunNotFound) {
				t.Fatalf("unknown run = %v", err)
			}
			readVersion := func() int {
				t.Helper()
				db, err := openReadOnlyState(state)
				if err != nil {
					t.Fatal(err)
				}
				defer db.Close()
				var current int
				if err := db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&current); err != nil {
					t.Fatal(err)
				}
				return current
			}
			if got := readVersion(); got != version {
				t.Fatalf("readonly inspection migrated schema to %d", got)
			}
			cache, err := OpenArtifactCache(filepath.Join(dir, "cache"), "", 0)
			if err != nil {
				t.Fatal(err)
			}
			defer cache.Close()
			opened, err := openTestStore(ctx, state, cache)
			if err != nil {
				t.Fatal(err)
			}
			if err := opened.Close(); err != nil {
				t.Fatal(err)
			}
			if got := readVersion(); got != stateSchema {
				t.Fatalf("writer migrated to schema %d, want %d", got, stateSchema)
			}
			verify()
			if after, err := ReadVersionExport(ctx, state); err != nil || !bytes.Equal(before, after) {
				t.Fatalf("migration changed deterministic version export: %s, %v", after, err)
			}
		})
	}
}
