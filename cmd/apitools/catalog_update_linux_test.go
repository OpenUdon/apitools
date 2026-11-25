package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/catalogupdate"
)

type catalogFixtureTransport func(*http.Request) (*http.Response, error)

func (f catalogFixtureTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func catalogCLIStateFixture(t *testing.T) string {
	t.Helper()
	base := os.Getenv("APITOOLS_TEST_STATE_DIR")
	if base == "" {
		dir := t.TempDir()
		if err := catalogupdate.CheckMaintenanceStorage(t.Context(), dir, filepath.Join(dir, "state.sqlite"), ""); err == nil {
			return dir
		}
		base = "." // The workspace can be durable when the system temp directory is tmpfs.
	}
	dir, err := os.MkdirTemp(base, ".catalog-update-cli-fixture-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Error(err)
		}
	})
	dir, err = filepath.Abs(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := catalogupdate.CheckMaintenanceStorage(t.Context(), dir, filepath.Join(dir, "state.sqlite"), ""); err != nil {
		t.Fatalf("use APITOOLS_TEST_STATE_DIR for an existing writable native filesystem: %v", err)
	}
	return dir
}

func TestCatalogUpdateCLIUsesFixturePublisherAndResumesState(t *testing.T) {
	dir := catalogCLIStateFixture(t)
	config := filepath.Join(dir, "maintenance.json")
	if err := os.WriteFile(config, []byte(`{"cache_dir":"artifacts","state_db":"state/db.sqlite"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	const openapi = `{"openapi":"3.0.3","info":{"title":"Records","version":"1.0"},"paths":{"/records":{"get":{"responses":{"204":{"description":"empty"}}}}}}`
	requests := 0
	authRequests := 0
	var progress bytes.Buffer
	client := &apitools.Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{
		Transport: catalogFixtureTransport(func(r *http.Request) (*http.Response, error) {
			requests++
			if !strings.Contains(progress.String(), "Catalog update run") {
				t.Error("run ID was not exposed before publisher requests")
			}
			length, body, status := int64(-1), "", http.StatusOK
			header := make(http.Header)
			switch r.URL.Host {
			case "api.attio.com":
				length = int64(len(openapi))
				header.Set("Content-Length", strconv.Itoa(len(openapi)))
				if r.Method == http.MethodGet {
					body = openapi
				}
			case "api.github.com":
				status = http.StatusNotFound
			case "auth0.com":
				authRequests++
				status = http.StatusNotFound
				if authRequests == 1 {
					status = http.StatusServiceUnavailable
					header.Set("Retry-After", "0")
				}
			default:
				if r.Method != http.MethodHead {
					t.Errorf("unexpected body request for %s", r.URL.Host)
				}
			}
			return &http.Response{
				StatusCode: status, Status: http.StatusText(status), Header: header,
				ContentLength: length, Body: io.NopCloser(strings.NewReader(body)), Request: r,
			}, nil
		}),
	}}
	run := func(complete bool, args ...string) apitools.CatalogUpdateReport {
		t.Helper()
		progress.Reset()
		var stdout bytes.Buffer
		args = append([]string{"--config", config, "--json"}, args...)
		if code := runCatalogUpdateWithClient(args, &stdout, &progress, client); code != 0 {
			t.Fatalf("exit %d: %s", code, progress.String())
		}
		var report apitools.CatalogUpdateReport
		if err := json.Unmarshal(stdout.Bytes(), &report); err != nil || report.Completed != complete || report.RunID == "" {
			t.Fatalf("CLI run ID = %q, complete = %t, expected complete = %t, error = %v", report.RunID, report.Completed, complete, err)
		}
		return report
	}
	one := run(true, "--provider", "Attio")
	if len(one.Outcomes) != 2 || one.Outcomes[1].Status != "review-required" {
		t.Fatalf("one-provider report = %+v", one)
	}
	before := requests
	resumed := run(true, "--resume", one.RunID)
	if resumed.RunID != one.RunID || requests != before {
		t.Fatalf("resume requested publisher again or changed run: %+v", resumed)
	}
	pending := run(false, "--all")
	foundPending := false
	for _, result := range pending.Outcomes {
		foundPending = foundPending || result.ProviderID == "auth0" && result.Status == "publisher-retry-pending"
	}
	if !foundPending || authRequests != 1 {
		t.Fatalf("expected one pending Auth0 publisher check, got %d visits, pending=%t", authRequests, foundPending)
	}
	var inspection, diagnostics bytes.Buffer
	beforeInspection := requests
	if code := runCatalogMaintenanceWithClient([]string{"report", "--config", config, "--run", pending.RunID, "--json"}, &inspection, &diagnostics, client); code != exitSuccess {
		t.Fatalf("cannot inspect pending run: %s", diagnostics.String())
	}
	var saved catalogMaintenanceReport
	if err := json.Unmarshal(inspection.Bytes(), &saved); err != nil || saved.Run.Completed ||
		len(saved.Run.Outcomes) != len(pending.Outcomes) || saved.Summary.PendingRetries == 0 || requests != beforeInspection {
		t.Fatalf("inspection resumed the run or dropped checkpoints: %+v, %v", saved.Summary, err)
	}
	all := run(true, "--resume", pending.RunID)
	statuses := make(map[string]bool)
	checkpoints := 0
	for _, result := range all.Outcomes {
		statuses[result.Status] = true
		if result.ReferenceID == "" {
			checkpoints++
		}
	}
	if checkpoints != len(catalog.BuiltInProviders()) || !statuses["size-unknown"] || !statuses["publisher-check-failed"] {
		t.Fatalf("full run checkpoints = %d, statuses = %v", checkpoints, statuses)
	}
	state := filepath.Join(dir, "state", "db.sqlite")
	lock, err := catalogupdate.LockMonthlyRun(t.Context(), state)
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Close()
	beforeRetry := requests
	diagnostics.Reset()
	if code := runCatalogMaintenanceWithClient([]string{"retry", "--provider", "attio", "--config", config}, io.Discard, &diagnostics, client); code != exitRuntime ||
		!strings.Contains(diagnostics.String(), "already running") || requests != beforeRetry {
		t.Fatalf("overlapping retry did not stop before publisher checks: %d, %s", code, diagnostics.String())
	}
	if err := lock.Close(); err != nil {
		t.Fatal(err)
	}
	diagnostics.Reset()
	inspection.Reset()
	if code := runCatalogMaintenanceWithClient([]string{"retry", "--provider", "attio", "--config", config, "--json"}, &inspection, &diagnostics, client); code != exitSuccess {
		t.Fatalf("provider retry after lock release failed: %d, %s", code, diagnostics.String())
	}
	if err := json.Unmarshal(inspection.Bytes(), &saved); err != nil || saved.Run.RunID == all.RunID ||
		len(saved.Run.Providers) != 1 || saved.Run.Providers[0] != "attio" || requests == beforeRetry {
		t.Fatalf("provider retry lost its distinct scoped report: %+v, %v", saved.Run.Providers, err)
	}
	var firstExport, secondExport bytes.Buffer
	beforeExport := requests
	for _, export := range []*bytes.Buffer{&firstExport, &secondExport} {
		if code := runCatalogUpdateWithClient([]string{"--config", config, "--export"}, export, &diagnostics, client); code != exitSuccess {
			t.Fatalf("catalog Git export failed: %d, %s", code, diagnostics.String())
		}
	}
	if !json.Valid(firstExport.Bytes()) || !bytes.Equal(firstExport.Bytes(), secondExport.Bytes()) || requests != beforeExport {
		t.Fatal("accepted Git export was unstable or contacted a publisher")
	}
}
