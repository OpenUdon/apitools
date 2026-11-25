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
	first, err := client.runMonthlyCatalogUpdate(firstContext, config, period, "", func(string) {
		_, overlapErr = client.runMonthlyCatalogUpdate(ctx, config, period, "", nil)
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
