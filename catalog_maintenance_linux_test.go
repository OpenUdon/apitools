package apitools

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/catalogupdate"
)

func TestMonthlyCatalogUpdateRestartsTheSameRunAndRejectsOverlap(t *testing.T) {
	dir, err := os.MkdirTemp(".", ".catalog-monthly-fixture-") // durable state filesystem
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
	config := CatalogUpdateConfig{CacheDir: filepath.Join(dir, "cache"), StateDB: filepath.Join(dir, "catalog.sqlite")}
	ctx, stop := context.WithTimeout(t.Context(), 30*time.Second)
	defer stop()
	firstContext, cancel := context.WithCancel(ctx)
	defer cancel()
	calls := 0
	var firstURL string
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		if calls == 1 {
			firstURL = r.URL.String()
		}
		if calls == 3 {
			cancel()
			return nil, firstContext.Err()
		}
		return &http.Response{StatusCode: http.StatusNotFound, Body: io.NopCloser(strings.NewReader("")), Request: r}, nil
	})}}
	const period = "2026-11"
	var overlapErr error
	first, err := client.runMonthlyCatalogUpdate(firstContext, config, period, "", func(id string) {
		_, overlapErr = client.runMonthlyCatalogUpdate(ctx, config, period, "", nil)
		for _, opts := range []CatalogUpdateOptions{{All: true}, {Provider: "attio"}, {Resume: id}} {
			opts.Config = config
			_, err := client.UpdateCatalog(ctx, opts)
			if err == nil || !strings.Contains(err.Error(), "already running") {
				t.Fatalf("manual writer bypassed the monthly lock: %+v, %v", opts, err)
			}
		}
	})
	if !errors.Is(err, context.Canceled) || first.RunID == "" || first.Completed || calls != 3 {
		t.Fatalf("expected interrupted checkpointed run, got %s completed %t, requests %d, error %v", first.RunID, first.Completed, calls, err)
	}
	if overlapErr == nil || !strings.Contains(overlapErr.Error(), "already running") {
		t.Fatalf("overlapping job did not stop at the run lock: %v", overlapErr)
	}
	pending, committed, err := catalogupdate.ReadExistingRun(ctx, config.StateDB, first.RunID)
	if err != nil || pending.Completed || len(first.Outcomes) == 0 || !reflect.DeepEqual(first.Outcomes, committed) {
		t.Fatalf("interrupted report omitted durable checkpoints: reported %d, stored %d, error %v", len(first.Outcomes), len(committed), err)
	}
	invalid := config
	invalid.ExpectedVolumeUUID = "apitools-fixture-missing-catalog-volume"
	before := calls
	if _, err := client.runMonthlyCatalogUpdate(ctx, invalid, period, "", nil); err == nil || calls != before {
		t.Fatalf("loss of expected volume did not prevent publisher work: %v", err)
	}
	const nextPeriod = "2026-12"
	var started []string
	resumed, err := client.runMonthlyCatalogUpdate(ctx, config, nextPeriod, "", func(id string) { started = append(started, id) })
	if err != nil || !resumed.Completed || len(started) != 2 || started[0] != first.RunID || started[1] == first.RunID ||
		resumed.RunID != started[1] || len(resumed.Providers) != len(catalog.BuiltInProviders()) {
		t.Fatalf("restart did not complete the old run and run the due month: %v, completed %t, providers %d, error %v", started, resumed.Completed, len(resumed.Providers), err)
	}
	old, _, err := catalogupdate.ReadExistingRun(ctx, config.StateDB, first.RunID)
	if err != nil || !old.Completed {
		t.Fatalf("carryover lost the first month's checkpoints: %+v, %v", old, err)
	}
	if firstURL == "" || calls <= before {
		t.Fatal("fixture did not check publishers before and after cancellation")
	}
	count := calls
	repeated, err := client.runMonthlyCatalogUpdate(ctx, config, nextPeriod, "", nil)
	if err != nil || repeated.RunID != resumed.RunID || len(repeated.Outcomes) != len(resumed.Outcomes) || calls != count {
		t.Fatalf("same month launched a duplicate: run %s, requests %d instead of %d, error %v", repeated.RunID, calls, count, err)
	}
}

func TestMonthlyCatalogUpdateRetryDelayUsesEffectivePerOriginLimits(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	report := CatalogUpdateReport{Providers: []string{"attio", "stripe"}}
	for i, id := range report.Providers {
		provider, ok := catalog.FindBuiltInProvider(id)
		if !ok || len(catalogUpdateReferenceSources(provider)) != 1 {
			t.Fatalf("invalid pending-source fixture for %s", id)
		}
		source := catalogUpdateReferenceSources(provider)[0]
		date := now.Add(time.Duration(i) * time.Hour)
		data, err := json.Marshal(catalogUpdateEvidence{Retry: &catalogUpdateRetry{NextAttemptAt: date}})
		if err != nil {
			t.Fatal(err)
		}
		report.Outcomes = append(report.Outcomes, catalogupdate.Outcome{ProviderID: id, ReferenceID: source.ReferenceID,
			Status: catalogupdate.RetryPendingStatus, Evidence: data})
	}
	if delay := catalogMonthlyRetryDelay(report, now); delay != 0 {
		t.Fatalf("ready attio retry would wait behind stripe for %s", delay)
	}
	provider, _ := catalog.FindBuiltInProvider("attio")
	origin := catalogUpdateArtifactOrigin(catalogUpdateReferenceSources(provider)[0].URL)
	var err error
	report.Outcomes[0].Evidence, err = json.Marshal(catalogUpdateEvidence{
		Retry:       &catalogUpdateRetry{NextAttemptAt: now},
		RetryBlocks: map[string]catalogUpdateOriginBlock{origin: {Until: now.Add(3 * time.Minute)}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if delay := catalogMonthlyRetryDelay(report, now); delay != 3*time.Minute {
		t.Fatalf("origin share must set the eligible delay, got %s", delay)
	}
}

func TestMonthlyCatalogUpdateRetainsCheckpointsWhenIdentityChanges(t *testing.T) {
	config := newCatalogRepairFixture(t)
	var err error
	config.CacheDir, err = filepath.Abs(config.CacheDir)
	if err != nil {
		t.Fatal(err)
	}
	config.StateDB, err = filepath.Abs(config.StateDB)
	if err != nil {
		t.Fatal(err)
	}
	ctx := t.Context()
	firstContext, cancel := context.WithCancel(ctx)
	defer cancel()
	calls := 0
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		if calls == 3 {
			cancel()
			return nil, firstContext.Err()
		}
		return &http.Response{StatusCode: http.StatusNotFound, Body: http.NoBody, Request: r}, nil
	})}}
	first, err := client.runMonthlyCatalogUpdate(firstContext, config, "2026-11", "", nil)
	if !errors.Is(err, context.Canceled) || len(first.Outcomes) == 0 {
		t.Fatalf("could not save an interrupted scheduled run: %d checkpoints, %v", len(first.Outcomes), err)
	}
	changed := config
	changed.RetainImmutable = !changed.RetainImmutable
	refused, err := client.runMonthlyCatalogUpdate(ctx, changed, "2026-12", "", nil)
	if err == nil || !strings.Contains(err.Error(), "2026-11") || !strings.Contains(err.Error(), "supersede --config FILE --run "+first.RunID) ||
		calls != 3 || refused.RunID != first.RunID || !reflect.DeepEqual(refused.Outcomes, first.Outcomes) {
		t.Fatalf("changed identity lost checkpoints or fetched from a publisher: %s had %d checkpoints, %d requests, %v",
			refused.RunID, len(refused.Outcomes), calls, err)
	}
	if err := catalogupdate.SupersedeMonthlyRun(ctx, changed.StateDB, changed.CacheDir, "", first.RunID); err != nil {
		t.Fatal(err)
	}
	old, outcomes, err := catalogupdate.ReadExistingRun(ctx, changed.StateDB, first.RunID)
	if err != nil || old.MonthlyPeriod != "2026-11" || old.MonthlySupersededAt.IsZero() || old.Completed || !reflect.DeepEqual(outcomes, first.Outcomes) {
		t.Fatalf("superseding did not preserve the unfinished run: %+v, %d checkpoints, %v", old, len(outcomes), err)
	}
	if err := catalogupdate.SupersedeMonthlyRun(ctx, changed.StateDB, changed.CacheDir, "", first.RunID); err != nil {
		t.Fatal(err)
	}
	saved, err := client.UpdateCatalog(ctx, CatalogUpdateOptions{Config: config, Resume: first.RunID})
	if err != nil || !saved.SupersededAt.Equal(old.MonthlySupersededAt) || !reflect.DeepEqual(saved.Outcomes, first.Outcomes) || calls != 3 {
		t.Fatalf("manual resume reopened the superseded run: %d requests, %v", calls, err)
	}
	const nextPeriod = "2026-12"
	due, err := client.runMonthlyCatalogUpdate(ctx, changed, nextPeriod, "", nil)
	if err != nil || !due.Completed || due.MonthlyPeriod != nextPeriod || due.RunID == first.RunID || calls <= 3 {
		t.Fatalf("new catalog did not service its due month: %s, %t, %d requests, %v", due.MonthlyPeriod, due.Completed, calls, err)
	}
	if err := catalogupdate.SupersedeMonthlyRun(ctx, changed.StateDB, changed.CacheDir, "", due.RunID); err == nil {
		t.Fatal("superseded a completed scheduled run")
	}
}

func TestMonthlyCatalogUpdatePersistsBoundedPublisherRetries(t *testing.T) {
	config := newCatalogRepairFixture(t)
	var err error
	config.CacheDir, err = filepath.Abs(config.CacheDir)
	if err != nil {
		t.Fatal(err)
	}
	config.StateDB, err = filepath.Abs(config.StateDB)
	if err != nil {
		t.Fatal(err)
	}
	var attempts []time.Time
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		status, header := http.StatusNotFound, make(http.Header)
		if r.Method == http.MethodHead && r.URL.Host == "api.attio.com" {
			attempts = append(attempts, time.Now())
			status = http.StatusTooManyRequests
			header.Set("Retry-After", "1")
		}
		return &http.Response{StatusCode: status, Header: header, Body: http.NoBody, Request: r}, nil
	})}}
	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
	defer cancel()
	run, err := client.runMonthlyCatalogUpdate(ctx, config, "2026-12", "", nil)
	if err != nil || !run.Completed || len(attempts) != catalogUpdatePublisherAttempts {
		t.Fatalf("scheduled run did not finish bounded retries: %d attempts, complete=%t, %v", len(attempts), run.Completed, err)
	}
	for i := 1; i < len(attempts); i++ {
		if delay := attempts[i].Sub(attempts[i-1]); delay < time.Second {
			t.Fatalf("retry %d did not honor Retry-After: %v", i+1, delay)
		}
	}
	for _, row := range run.Outcomes {
		if row.ProviderID == "attio" && row.ReferenceID != "" {
			if row.Status != "publisher-retry-exhausted" || readCatalogUpdateRetry(row).Attempt != catalogUpdatePublisherAttempts {
				t.Fatalf("wrong saved terminal retry: %s, %s", row.Status, row.Evidence)
			}
			return
		}
	}
	t.Fatal("missing Attio retry checkpoint")
}
