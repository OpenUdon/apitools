package catalogupdate

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMaintenancePreflightSharesWriterSchemaChecks(t *testing.T) {
	root, err := os.MkdirTemp(".", ".catalog-monthly-schema-fixture-") // durable filesystem
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(root); err != nil {
			t.Error(err)
		}
	})
	root, err = filepath.Abs(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct{ name, schema, expected string }{
		{"legacy", schemaSQL + versionSchemaSQL, ""},
		{"missing-preferred", schemaSQL + versionSchemaSQL + `DROP TABLE update_preferred; PRAGMA user_version=4;`, "lacks version tables"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			state := filepath.Join(root, tt.name+".sqlite")
			db, err := sql.Open("sqlite", state)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := db.ExecContext(t.Context(), tt.schema); err != nil {
				t.Fatal(err)
			}
			if err := db.Close(); err != nil {
				t.Fatal(err)
			}
			before, err := os.ReadFile(state)
			if err != nil {
				t.Fatal(err)
			}
			err = CheckMaintenanceStorage(t.Context(), root, state, "")
			if tt.expected == "" && err != nil || tt.expected != "" && (err == nil || !strings.Contains(err.Error(), tt.expected)) {
				t.Fatalf("preflight returned %v, expected %q", err, tt.expected)
			}
			after, err := os.ReadFile(state)
			if err != nil || string(before) != string(after) {
				t.Fatalf("preflight changed the state database: %v", err)
			}
		})
	}
}
