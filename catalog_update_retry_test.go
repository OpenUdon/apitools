package apitools

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/catalogupdate"
)

func TestCatalogUpdateQuotaYieldsOriginAndResumesLaterProviders(t *testing.T) {
	t.Run("resume after quota", func(t *testing.T) { testCatalogUpdateQuota(t, "0", false) })
	t.Run("defer beyond recovery horizon", func(t *testing.T) { testCatalogUpdateQuota(t, "999999999", true) })
}

func TestCatalogUpdateDefersMissingRetryArtifactAndChecksLaterProvider(t *testing.T) {
	cfg := newCatalogRepairFixture(t)
	ctx := t.Context()
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
	run, err := store.StartRun(ctx, catalog.BuiltInCatalogManifest().CatalogSHA256, []string{"attio", "stripe"})
	if err != nil {
		t.Fatal(err)
	}
	const body = "historical reviewed artifact"
	sum := sha256.Sum256([]byte(body))
	object := catalogupdate.Object{SHA256: hex.EncodeToString(sum[:]), Bytes: int64(len(body))}
	if err := store.PutObject(ctx, object, strings.NewReader(body)); err != nil {
		t.Fatal(err)
	}
	provider, _ := catalog.FindBuiltInProvider("attio")
	source := catalogUpdateReferenceSources(provider)[0]
	data, err := json.Marshal(catalogUpdateEvidence{
		ReviewedBaseline: &catalogUpdateBaselineEvidence{Preferred: catalogupdate.Version{Object: object}},
		RetryBlocks:      map[string]catalogUpdateOriginBlock{catalogUpdateArtifactOrigin(source.URL): {OutsideWindow: true}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Record(ctx, run.ID, catalogupdate.Outcome{ProviderID: "attio", ReferenceID: source.ReferenceID,
		Status: catalogupdate.RetryPendingStatus, Object: &object, Evidence: data}); err != nil {
		t.Fatal(err)
	}
	if err := cache.Remove("sha256-" + object.SHA256); err != nil {
		t.Fatal(err)
	}
	checkedStripe := false
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host == "api.attio.com" {
			t.Fatal("contacted an origin deferred beyond the retry window")
		}
		checkedStripe = true
		return &http.Response{StatusCode: http.StatusNotFound, Body: http.NoBody, Request: r}, nil
	})}}
	result, err := client.continueCatalogUpdate(ctx, store, CatalogUpdateReport{RunID: run.ID, Providers: run.Providers}, false, cfg)
	if err != nil || !result.Completed || !checkedStripe {
		t.Fatalf("missing object blocked later references: completed=%t, stripe=%t, %v", result.Completed, checkedStripe, err)
	}
	for _, row := range result.Outcomes {
		if row.ProviderID == "attio" && row.ReferenceID == source.ReferenceID {
			if row.Status != "publisher-retry-deferred" || row.Object == nil || *row.Object != object ||
				!strings.Contains(string(row.Evidence), "artifact is unavailable") || strings.Contains(string(row.Evidence), "reviewed_baseline") {
				t.Fatalf("deferral obscured missing historical artifact: %+v", row)
			}
			return
		}
	}
	t.Fatal("missing deferred Attio checkpoint")
}

func testCatalogUpdateQuota(t *testing.T, hint string, beyondWindow bool) {
	t.Helper()
	cfg := newCatalogRepairFixture(t)
	ctx := t.Context()
	cache, err := catalogupdate.OpenArtifactCache(cfg.CacheDir, "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	openStore := func() *catalogupdate.Store {
		t.Helper()
		s, err := catalogupdate.OpenStore(ctx, cfg.StateDB, cache)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = s.Close() })
		return s
	}
	store := openStore()
	run, err := store.StartRun(ctx, catalog.BuiltInCatalogManifest().CatalogSHA256, []string{"amazon-selling-partner", "datadog", "attio"})
	if err != nil {
		t.Fatal(err)
	}
	input := CatalogUpdateReport{RunID: run.ID, Providers: run.Providers}
	throttled := true
	github, unrelated := 0, 0
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		header := make(http.Header)
		code := http.StatusOK
		switch r.URL.Host {
		case "api.github.com":
			github++
			if throttled {
				code = http.StatusForbidden
				header.Set("X-RateLimit-Remaining", "0")
				header.Set("Retry-After", hint)
				break
			}
			return catalogRepairResponse(r, `[{"sha":"1111111111111111111111111111111111111111"}]`), nil
		case "raw.githubusercontent.com":
			if r.Method != http.MethodHead {
				t.Fatal("unexpected body request for unknown content length")
			}
		case "api.attio.com":
			unrelated++
		default:
			t.Fatalf("unexpected host: %s", r.URL.Host)
		}
		return &http.Response{StatusCode: code, Header: header, Body: http.NoBody, Request: r}, nil
	})}}
	first, err := client.continueCatalogUpdate(ctx, store, input, false, cfg)
	if err != nil || first.Completed != beyondWindow || github != 1 || unrelated != 1 {
		t.Fatalf("first pass: complete=%t GitHub=%d unrelated=%d error=%v", first.Completed, github, unrelated, err)
	}
	if beyondWindow {
		terminal := map[string]string{}
		for _, row := range first.Outcomes {
			if row.ReferenceID != "" {
				terminal[row.ProviderID] = row.Status
			}
		}
		if terminal["amazon-selling-partner"] != "publisher-retry-exhausted" || terminal["datadog"] != "publisher-retry-deferred" {
			t.Fatalf("missing terminal quota outcomes: %v", terminal)
		}
		return
	}
	if len(first.Outcomes) != 3 || first.Outcomes[0].Status != catalogupdate.RetryPendingStatus ||
		readCatalogUpdateRetry(first.Outcomes[0]).HTTPStatus != http.StatusForbidden {
		t.Fatalf("unrelated results or retry evidence missing: %+v", first.Outcomes)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	store = openStore()
	throttled = false
	after, err := client.continueCatalogUpdate(ctx, store, input, false, cfg)
	if err != nil || !after.Completed || github != 3 || unrelated != 1 || len(after.Outcomes) != 6 {
		t.Fatalf("resumed pass: complete=%t GitHub=%d unrelated=%d outcome count=%d error=%v", after.Completed, github, unrelated, len(after.Outcomes), err)
	}
}

func TestCatalogUpdatePublisherRetryStopsAtBound(t *testing.T) {
	cfg := newCatalogRepairFixture(t)
	ctx := t.Context()
	attempts := 0
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodHead || r.URL.Host != "api.attio.com" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL)
		}
		attempts++
		return &http.Response{StatusCode: http.StatusServiceUnavailable,
			Header: http.Header{"Retry-After": {"0"}}, Body: http.NoBody, Request: r}, nil
	})}}
	options := CatalogUpdateOptions{Config: cfg, Provider: "attio"}
	for n := 1; n <= catalogUpdatePublisherAttempts; n++ {
		result, err := client.UpdateCatalog(ctx, options)
		if err != nil || result.Completed != (n == catalogUpdatePublisherAttempts) || attempts != n {
			t.Fatalf("pass %d: completed=%t attempts=%d error=%v", n, result.Completed, attempts, err)
		}
		outcome := result.Outcomes[len(result.Outcomes)-1]
		retry := readCatalogUpdateRetry(outcome)
		wantStatus := catalogupdate.RetryPendingStatus
		if n == catalogUpdatePublisherAttempts {
			wantStatus = "publisher-retry-exhausted"
		}
		if outcome.Status != wantStatus || retry.Attempt != n || retry.Reason != "publisher-server-error" || retry.HTTPStatus != 503 {
			t.Fatalf("pass %d retry %+v", n, retry)
		}
		options.Provider, options.Resume = "", result.RunID
	}
	if _, err := client.UpdateCatalog(ctx, options); err != nil || attempts != catalogUpdatePublisherAttempts {
		t.Fatalf("completed run made another publisher request: %d %v", attempts, err)
	}
}

func TestCatalogUpdateHonorsFiniteRetryWindow(t *testing.T) {
	for _, tc := range []struct {
		name, hint           string
		exhausted, saturated bool
	}{
		{"future", "600", false, false},
		{"beyond run horizon", "999999999", true, false},
		{"unrepresentable delay", "99999999999999999999999999", true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg := newCatalogRepairFixture(t)
			attempts := 0
			client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
				attempts++
				return &http.Response{StatusCode: http.StatusTooManyRequests,
					Header: http.Header{"Retry-After": {tc.hint}}, Body: http.NoBody, Request: r}, nil
			})}}
			first, err := client.UpdateCatalog(t.Context(), CatalogUpdateOptions{Config: cfg, Provider: "attio"})
			if err != nil || first.Completed != tc.exhausted {
				t.Fatalf("first pass completed=%t, error=%v", first.Completed, err)
			}
			retry := readCatalogUpdateRetry(first.Outcomes[len(first.Outcomes)-1])
			if retry.OutsideWindow != tc.exhausted || retry.HintSaturated != tc.saturated || !retry.NextAttemptAt.After(time.Now()) {
				t.Fatalf("observed future retry hint was lost: %+v", retry)
			}
			second, err := client.UpdateCatalog(t.Context(), CatalogUpdateOptions{Config: cfg, Resume: first.RunID})
			if err != nil || second.Completed != tc.exhausted || attempts != 1 {
				t.Fatalf("early resume completed=%t, attempts=%d, error=%v", second.Completed, attempts, err)
			}
		})
	}
}

func TestCatalogUpdatePublisherFailureKeepsSafeRetryHints(t *testing.T) {
	const raw = "https://api.github.com/repos/test/example/commits"
	reset := time.Now().UTC().Truncate(time.Second).Add(15 * time.Minute)
	cause := HTTPStatusError{Code: 403, Status: "private response"}
	err := catalogUpdateRetryFailure(raw, cause, 403, "1", "0", strconv.FormatInt(reset.Unix(), 10))
	var retry *catalogUpdatePublisherError
	if !errors.As(err, &retry) || !retry.notify.Equal(reset) || retry.retry.Origin != "https://api.github.com:443" || err.Error() != "retryable publisher rate-limit" {
		t.Fatalf("quota result = %v", err)
	}
	err = catalogUpdateRetryFailure(raw, cause, 403, "", "", "")
	if errors.As(err, &retry) {
		t.Fatalf("ordinary forbidden response treated as quota: %v", err)
	}
	err = catalogUpdateRetryFailure(raw, cause, 403, "1", "0", "999999999999999999")
	if !errors.As(err, &retry) || retry.retry.HintSaturated || retry.notify.After(time.Now().Add(2*time.Second)) {
		t.Fatalf("an uninterpretable reset replaced the explicit Retry-After: %v", err)
	}
	err = catalogUpdateRetryFailure(raw, context.DeadlineExceeded, 0, "", "", "")
	if !errors.As(err, &retry) || retry.retry.Reason != "timeout" || retry.retry.Shared {
		t.Fatalf("request timeout = %v", err)
	}
}

func TestNilClientCanPrepareCatalogUpdate(t *testing.T) {
	var client *Client
	report, err := client.UpdateCatalog(t.Context(), CatalogUpdateOptions{Provider: "datadog", DryRun: true})
	if err != nil || !report.DryRun || len(report.Providers) != 1 {
		t.Fatalf("nil client dry run = %+v, %v", report, err)
	}
}

func TestNilClientReadsSuccessfulPublisherMetadata(t *testing.T) {
	const commits = `[{"sha":"1111111111111111111111111111111111111111"}]`
	var client *Client
	response := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(commits)), ContentLength: int64(len(commits))}
	t.Cleanup(func() { _ = response.Body.Close() })
	got := client.readAPIVersionProbeBody(t.Context(), "https://api.github.com/repos/test/spec/commits", false, response, apiVersionProbe{})
	if got.err != nil || string(got.content) != commits {
		t.Fatalf("nil client metadata fixture: content=%q error=%v", got.content, got.err)
	}
}
