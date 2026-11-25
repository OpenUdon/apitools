package apitools

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/catalogupdate"
)

func TestCatalogUpdateResumesAfterInterruptedSQLiteWrite(t *testing.T) {
	const childStateKey = "APITOOLS_TEST_C10_HOT_JOURNAL"
	if state := os.Getenv(childStateKey); state != "" {
		db, err := sql.Open("sqlite", state)
		if err != nil {
			t.Fatal(err)
		}
		db.SetMaxOpenConns(1)
		if _, err := db.Exec(`PRAGMA cache_size = 3; BEGIN IMMEDIATE; UPDATE crash_fixture SET payload = replace(payload, '0', '1');`); err != nil {
			t.Fatal(err)
		}
		os.Exit(23)
	}
	t.Run("manual", func(t *testing.T) { testCatalogUpdateRecoversJournal(t, childStateKey, false) })
	t.Run("monthly", func(t *testing.T) { testCatalogUpdateRecoversJournal(t, childStateKey, true) })
}

func testCatalogUpdateRecoversJournal(t *testing.T, childStateKey string, monthly bool) {
	t.Helper()
	ctx := t.Context()
	cfg := newCatalogRepairFixture(t)
	var err error
	cfg.CacheDir, err = filepath.Abs(cfg.CacheDir)
	if err != nil {
		t.Fatal(err)
	}
	cfg.StateDB, err = filepath.Abs(cfg.StateDB)
	if err != nil {
		t.Fatal(err)
	}
	cache, err := catalogupdate.OpenArtifactCache(cfg.CacheDir, "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	store, err := catalogupdate.OpenStore(ctx, cfg.StateDB, cache)
	if err != nil {
		t.Fatal(err)
	}
	catalogSHA := catalog.BuiltInCatalogManifest().CatalogSHA256
	settings, err := json.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	identity := fmt.Sprintf("%x", sha256.Sum256(append([]byte(catalogSHA+"\n"), settings...)))
	period := time.Now().Format("2006-01")
	var run catalogupdate.Run
	if monthly {
		run, err = store.MonthlyRun(ctx, period, identity, catalogSHA, []string{"attio"})
	} else {
		run, err = store.StartRunSnapshot(ctx, identity, catalogSHA, []string{"attio"})
	}
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite", cfg.StateDB)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `
CREATE TABLE crash_fixture(payload TEXT);
WITH RECURSIVE seq(i) AS (SELECT 1 UNION ALL SELECT i + 1 FROM seq WHERE i < 300)
INSERT INTO crash_fixture SELECT printf('%03500d', 0) FROM seq;`); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	child := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestCatalogUpdateResumesAfterInterruptedSQLiteWrite$")
	child.Env = append(os.Environ(), childStateKey+"="+cfg.StateDB)
	output, err := child.CombinedOutput()
	var exit *exec.ExitError
	if !errors.As(err, &exit) || exit.ExitCode() != 23 {
		t.Fatalf("could not stage interrupted SQLite writer: %s %v", output, err)
	}
	options := CatalogUpdateOptions{Resume: run.ID, Config: CatalogUpdateConfig{StateDB: cfg.StateDB}}
	if _, err := (&Client{}).UpdateCatalog(ctx, options); !errors.Is(err, catalogupdate.ErrRunJournalRecovery) {
		t.Fatalf("read-only inspection unexpectedly recovered SQLite: %v", err)
	}
	if monthly {
		if _, err := catalogupdate.ReadLatestMonthlyRunID(ctx, cfg.StateDB); !errors.Is(err, catalogupdate.ErrRunJournalRecovery) {
			t.Fatalf("monthly report could not identify the hot journal: %v", err)
		}
		if err := catalogupdate.CheckMaintenanceStorage(ctx, cfg.CacheDir, cfg.StateDB, ""); !errors.Is(err, catalogupdate.ErrRunJournalRecovery) {
			t.Fatalf("preflight could not identify the hot journal: %v", err)
		}
	}
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodHead || r.URL.Host != "api.attio.com" {
			t.Fatalf("unexpected fixture request: %s %s", r.Method, r.URL)
		}
		return &http.Response{StatusCode: http.StatusMethodNotAllowed, Header: make(http.Header), Body: http.NoBody, Request: r}, nil
	})}}
	options.Config = cfg
	var result CatalogUpdateReport
	if monthly {
		result, err = client.runMonthlyCatalogUpdate(ctx, cfg, period, "", nil)
	} else {
		result, err = client.UpdateCatalog(ctx, options)
	}
	if err != nil || !result.Completed || result.RunID != run.ID || result.RunIdentity != identity ||
		len(result.Outcomes) != 2 || result.Outcomes[1].Status != catalogUpdateSizeUnknown {
		t.Fatalf("resume after journal recovery = %t, status count %d, %v", result.Completed, len(result.Outcomes), err)
	}
}

func TestCatalogUpdateReadsCompletedRunWithoutCacheOrCurrentCatalog(t *testing.T) {
	ctx := t.Context()
	cfg := newCatalogRepairFixture(t)
	cache, err := catalogupdate.OpenArtifactCache(cfg.CacheDir, "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	store, err := catalogupdate.OpenStore(ctx, cfg.StateDB, cache)
	if err != nil {
		t.Fatal(err)
	}
	identity := fmt.Sprintf("%x", sha256.Sum256([]byte("old catalog and settings")))
	manifest := fmt.Sprintf("%x", sha256.Sum256([]byte("different historical catalog")))
	started, err := store.StartRunSnapshot(ctx, identity, manifest, []string{"attio"})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Record(ctx, started.ID, catalogupdate.Outcome{ProviderID: "attio", Status: "checked"}); err != nil {
		t.Fatal(err)
	}
	if err := store.FinishRun(ctx, started.ID); err != nil {
		t.Fatal(err)
	}
	incomplete, err := store.StartRunSnapshot(ctx, identity, manifest, []string{"auth0"})
	if err != nil {
		t.Fatal(err)
	}
	configBytes, err := json.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	legacyIdentity := fmt.Sprintf("%x", sha256.Sum256(append([]byte(catalog.BuiltInCatalogManifest().CatalogSHA256+"\n"), configBytes...)))
	legacy, err := store.StartRun(ctx, legacyIdentity, []string{"attio"})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Record(ctx, legacy.ID, catalogupdate.Outcome{ProviderID: "attio", Status: "checked"}); err != nil {
		t.Fatal(err)
	}
	if err := store.FinishRun(ctx, legacy.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(cfg.StateDB)
	if err != nil {
		t.Fatal(err)
	}
	changed := CatalogUpdateConfig{StateDB: cfg.StateDB, CacheDir: filepath.Join(filepath.Dir(cfg.CacheDir), "unused"),
		ExpectedVolumeUUID: "different-volume"}
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		t.Fatalf("historical inspection requested %s", r.URL)
		return nil, nil
	})}}
	options := CatalogUpdateOptions{Config: changed, Resume: started.ID, OnStart: func(string) { t.Fatal("read-only inspection began an update") }}
	report, err := client.UpdateCatalog(ctx, options)
	if err != nil || !report.Completed || report.CatalogSHA256 != manifest || report.RunIdentity != identity ||
		!report.StartedAt.Equal(started.StartedAt) || report.CompletedAt.IsZero() || report.CompletedAt.Before(report.StartedAt) || len(report.Outcomes) != 1 {
		t.Fatalf("completed historical report = %+v, %v", report, err)
	}
	second, err := client.UpdateCatalog(ctx, CatalogUpdateOptions{Config: CatalogUpdateConfig{StateDB: cfg.StateDB}, Resume: started.ID})
	firstJSON, _ := json.Marshal(report)
	secondJSON, _ := json.Marshal(second)
	if err != nil || !bytes.Equal(firstJSON, secondJSON) {
		t.Fatalf("completed report depended on current settings: %s / %s, %v", firstJSON, secondJSON, err)
	}
	for _, settings := range []CatalogUpdateConfig{cfg, changed, {StateDB: cfg.StateDB}} {
		report, err := client.UpdateCatalog(ctx, CatalogUpdateOptions{Config: settings, Resume: legacy.ID})
		if err != nil || report.CatalogSHA256 != "" || report.RunIdentity != legacyIdentity {
			t.Fatalf("legacy run acquired an unstored catalog digest: %+v, %v", report, err)
		}
	}
	options.Resume = incomplete.ID
	if _, err := client.UpdateCatalog(ctx, options); err == nil || !strings.Contains(err.Error(), "configuration changed") {
		t.Fatalf("incomplete historical run was resumed with different inputs: %v", err)
	}
	options.Resume = "unknown-run"
	if _, err := client.UpdateCatalog(ctx, options); !errors.Is(err, catalogupdate.ErrRunNotFound) || !strings.Contains(err.Error(), "unknown-run") {
		t.Fatalf("unknown run diagnostic = %v", err)
	}
	after, err := os.ReadFile(cfg.StateDB)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("inspection wrote SQLite state: %v", err)
	}
	if _, err := os.Stat(changed.CacheDir); !os.IsNotExist(err) {
		t.Fatalf("inspection touched the configured artifact path: %v", err)
	}
}

func TestCatalogUpdateRecordsAttemptTimesAndSafeLocatorsOnFailure(t *testing.T) {
	for _, tc := range []struct{ provider, status string }{
		{"attio", "publisher-check-failed"}, {"amazon-selling-partner", "publisher-source-unexamined"}, {"datadog", "publisher-check-failed"},
	} {
		t.Run(tc.provider, func(t *testing.T) {
			cfg := newCatalogRepairFixture(t)
			provider, _ := catalog.FindBuiltInProvider(tc.provider)
			source := catalogUpdateReferenceSources(provider)[0]
			const redirect = "https://api.attio.com/openapi/fixture-redirect.json"
			client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
				if r.URL.Host == "api.github.com" {
					if tc.provider == "datadog" {
						return catalogRepairResponse(r, `[{"sha":"1111111111111111111111111111111111111111"}]`), nil
					}
					return &http.Response{StatusCode: http.StatusNotFound, Header: make(http.Header), Body: http.NoBody, Request: r}, nil
				}
				if r.Method == http.MethodHead {
					if tc.provider == "attio" && r.URL.String() == source.URL {
						return &http.Response{StatusCode: http.StatusFound, Header: http.Header{"Location": {redirect}}, Body: http.NoBody, Request: r}, nil
					}
					return nil, errors.New("mock connection failure")
				}
				t.Fatalf("unexpected fixture request: %s %s", r.Method, r.URL)
				return nil, nil
			})}}
			before := time.Now().UTC()
			report, err := client.UpdateCatalog(t.Context(), CatalogUpdateOptions{Config: cfg, Provider: tc.provider})
			if err != nil || !report.Completed || len(report.Outcomes) != 2 || report.Outcomes[1].Status != tc.status {
				t.Fatalf("run outcome count %d, completed=%t, error=%v", len(report.Outcomes), report.Completed, err)
			}
			var evidence catalogUpdateEvidence
			if err := json.Unmarshal(report.Outcomes[1].Evidence, &evidence); err != nil {
				t.Fatal(err)
			}
			if evidence.CheckedAt.Before(evidence.AttemptedAt) || evidence.AttemptedAt.Before(before) ||
				evidence.CheckedAt.After(report.CompletedAt) || evidence.Source.URL != source.URL ||
				(evidence.Source.ResolvedURL != "") != (tc.provider == "datadog") || report.CatalogSHA256 != catalog.BuiltInCatalogManifest().CatalogSHA256 {
				t.Fatalf("incomplete failure evidence: %s", report.Outcomes[1].Evidence)
			}
			expectedTarget := firstNonEmpty(evidence.Source.ResolvedURL, source.URL)
			if tc.provider == "attio" {
				expectedTarget = redirect
			}
			if tc.status == "publisher-check-failed" && (evidence.Artifact == nil || evidence.Artifact.URL != expectedTarget ||
				evidence.Artifact.CheckedAt.Before(evidence.Artifact.AttemptedAt) || evidence.Artifact.AttemptedAt.Before(before)) {
				t.Fatalf("failed preflight lost artifact evidence: %s", report.Outcomes[1].Evidence)
			}
		})
	}
}

func TestCatalogUpdateRecordsDNSTimeoutAsAnAttempt(t *testing.T) {
	originalResolver := net.DefaultResolver
	net.DefaultResolver = &net.Resolver{PreferGo: true, Dial: func(context.Context, string, string) (net.Conn, error) {
		return nil, context.DeadlineExceeded
	}}
	t.Cleanup(func() { net.DefaultResolver = originalResolver })
	for _, provider := range []string{"attio", "amazon-selling-partner"} {
		t.Run(provider, func(t *testing.T) {
			cfg := newCatalogRepairFixture(t)
			before := time.Now().UTC()
			report, err := (&Client{}).UpdateCatalog(t.Context(), CatalogUpdateOptions{Config: cfg, Provider: provider})
			if err != nil || report.Completed || len(report.Outcomes) != 1 || report.Outcomes[0].Status != catalogupdate.RetryPendingStatus {
				t.Fatalf("DNS timeout report status = %t, outcomes=%d, %v", report.Completed, len(report.Outcomes), err)
			}
			var evidence catalogUpdateEvidence
			if err := json.Unmarshal(report.Outcomes[0].Evidence, &evidence); err != nil {
				t.Fatal(err)
			}
			if evidence.Retry == nil || evidence.Retry.Reason != "timeout" ||
				evidence.AttemptedAt.Before(before) || evidence.CheckedAt.Before(evidence.AttemptedAt) {
				t.Fatalf("missing DNS attempt evidence: %s", report.Outcomes[0].Evidence)
			}
		})
	}
}

func TestCatalogUpdateRejectsArtifactFragmentsWithoutPersistingThem(t *testing.T) {
	ctx := t.Context()
	cfg := newCatalogRepairFixture(t)
	cache, err := catalogupdate.OpenArtifactCache(cfg.CacheDir, "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	store, err := catalogupdate.OpenStore(ctx, cfg.StateDB, cache)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	run, err := store.StartRun(ctx, catalog.BuiltInCatalogManifest().CatalogSHA256, []string{"example"})
	if err != nil {
		t.Fatal(err)
	}
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		t.Fatal("invalid locator reached the publisher transport")
		return nil, nil
	})}}
	source := catalogUpdateSource{ProviderID: "example", ReferenceID: "records", Kind: catalog.SpecKindOpenAPI,
		URL: "https://example.com/spec#access_token=fixture-secret"}
	out, err := client.checkCatalogUpdateSource(ctx, store, run.ID, catalog.Provider{}, source, cfg, nil)
	if err != nil || out.Status != "publisher-source-unexamined" ||
		bytes.Contains(out.Evidence, []byte("fixture-secret")) || bytes.Contains(out.Evidence, []byte("attempted_at")) {
		t.Fatalf("unsafe skipped outcome = %s, %s, %v", out.Status, out.Evidence, err)
	}
	const index = `{"results":[{"name":"records","versions":[{"version":"1","stage":"LATEST","openApi":"https://api.hubapi.com/public/api/spec/records#access_token=fixture-secret"}]}]}`
	if _, err := hubSpotUpdateInventory("https://api.hubapi.com/public/api/spec/v1/specs", time.Now(), []byte(index)); err == nil {
		t.Fatal("index accepted an artifact link with a URL fragment")
	}
}
