package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/apitools"
)

func TestCatalogMaintenancePreflightAndDryRun(t *testing.T) {
	dir, err := os.MkdirTemp(".", ".catalog-maintenance-cli-fixture-") // durable state filesystem
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Error(err)
		}
	})
	config := filepath.Join(dir, "maintenance.json")
	if err := os.WriteFile(config, []byte(`{"cache_dir":"cache","state_db":"state.sqlite"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	if code := runCatalogMaintenance([]string{"run", "--config", config, "--dry-run", "--json"}, &out, &errOut); code != 0 {
		t.Fatalf("dry run exited %d: %s", code, errOut.String())
	}
	var report apitools.CatalogUpdateReport
	if err := json.Unmarshal(out.Bytes(), &report); err != nil || !report.DryRun || len(report.Providers) != 334 {
		t.Fatalf("bad dry-run scope: %d providers, %v", len(report.Providers), err)
	}
	for _, path := range []string{"state.sqlite", "state.sqlite.monthly.lock", "cache"} {
		if _, err := os.Stat(filepath.Join(dir, path)); !os.IsNotExist(err) {
			t.Fatalf("dry run touched %s: %v", path, err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "cache"), 0o700); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errOut.Reset()
	args := []string{"preflight", "--config", config, "--json"}
	if code := runCatalogMaintenance(args, &out, &errOut); code != 0 {
		t.Fatalf("preflight exited %d: %s", code, errOut.String())
	}
	var readiness catalogMaintenanceReadiness
	if err := json.Unmarshal(out.Bytes(), &readiness); err != nil || !filepath.IsAbs(readiness.CacheDir) || !filepath.IsAbs(readiness.Binary) {
		t.Fatalf("preflight did not resolve absolute paths: %+v, %v", readiness, err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 2 {
		t.Fatalf("preflight left files: %v, %v", entries, err)
	}
	if err := os.WriteFile(config, []byte(`{"cache_dir":"cache","state_db":"state.sqlite","expected_volume_uuid":"apitools-fixture-lost-volume"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errOut.Reset()
	if code := runCatalogMaintenance(args, &out, &errOut); code != exitRuntime || !strings.Contains(errOut.String(), "apitools-fixture-lost-volume") {
		t.Fatalf("missing mount passed preflight: %d, %s", code, errOut.String())
	}
}
