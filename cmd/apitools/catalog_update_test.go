package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/apitools/catalog"
)

func TestCatalogUpdateDryRunThroughCLI(t *testing.T) {
	dir := t.TempDir()
	config := filepath.Join(dir, "catalog.json")
	if err := os.WriteFile(config, []byte(`{"cache_dir":"artifacts","state_db":"state/db.sqlite"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		args  []string
		count int
	}{
		{[]string{"--all"}, len(catalog.BuiltInProviders())},
		{[]string{"--provider", "Attio"}, 1},
	} {
		var stdout, stderr bytes.Buffer
		args := append([]string{"update", "--config", config, "--dry-run", "--json"}, tc.args...)
		if exit := runCatalog(args, &stdout, &stderr); exit != 0 {
			t.Fatalf("exit %d: %s", exit, stderr.String())
		}
		var report apitools.CatalogUpdateReport
		if err := json.Unmarshal(stdout.Bytes(), &report); err != nil ||
			!report.DryRun || report.RunID != "" || len(report.Providers) != tc.count {
			t.Fatalf("report = %+v; %v", report, err)
		}
	}
	for _, path := range []string{"artifacts", "state"} {
		if _, err := os.Stat(filepath.Join(dir, path)); !os.IsNotExist(err) {
			t.Fatalf("dry run created %s: %v", path, err)
		}
	}
}
