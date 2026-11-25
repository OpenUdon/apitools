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
)

type catalogFixtureTransport func(*http.Request) (*http.Response, error)

func (f catalogFixtureTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestCatalogUpdateCLIUsesFixturePublisherAndResumesState(t *testing.T) {
	dir, err := os.MkdirTemp(".", ".catalog-update-cli-fixture-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Error(err)
		}
	})
	config := filepath.Join(dir, "maintenance.json")
	if err := os.WriteFile(config, []byte(`{"cache_dir":"artifacts","state_db":"state/db.sqlite"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	const openapi = `{"openapi":"3.0.3","info":{"title":"Records","version":"1.0"},"paths":{"/records":{"get":{"responses":{"204":{"description":"empty"}}}}}}`
	requests := 0
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
				status = http.StatusServiceUnavailable
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
	run := func(args ...string) apitools.CatalogUpdateReport {
		t.Helper()
		progress.Reset()
		var stdout bytes.Buffer
		args = append([]string{"--config", config, "--json"}, args...)
		if code := runCatalogUpdateWithClient(args, &stdout, &progress, client); code != 0 {
			t.Fatalf("exit %d: %s", code, progress.String())
		}
		var report apitools.CatalogUpdateReport
		if err := json.Unmarshal(stdout.Bytes(), &report); err != nil || !report.Completed || report.RunID == "" {
			t.Fatalf("CLI report = %+v, %v", report, err)
		}
		return report
	}
	one := run("--provider", "Attio")
	if len(one.Outcomes) != 2 || one.Outcomes[1].Status != "review-required" {
		t.Fatalf("one-provider report = %+v", one)
	}
	before := requests
	resumed := run("--resume", one.RunID)
	if resumed.RunID != one.RunID || requests != before {
		t.Fatalf("resume requested publisher again or changed run: %+v", resumed)
	}
	all := run("--all")
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
}
