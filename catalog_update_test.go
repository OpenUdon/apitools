//go:build linux

package apitools

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/catalogupdate"
)

func TestCatalogUpdateRunResumesAndAdoptsOnlyCompatibleFixture(t *testing.T) {
	// The maintenance DB requires a durable filesystem. The workspace is
	// file-backed even when the host defaults TMPDIR to tmpfs.
	dir, err := os.MkdirTemp(".", ".catalog-update-fixture-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Error(err)
		}
	})
	ctx := t.Context()
	cfg := CatalogUpdateConfig{
		CacheDir: filepath.Join(dir, "cache"), StateDB: filepath.Join(dir, "state", "catalog.sqlite"),
		ReviewedOperations: map[string][]string{"attio/attio-rest-api-openapi": {"#/paths/~1records/get"}},
	}
	const original = `{
	  "openapi":"3.0.3",
	  "info":{"title":"Records","version":"1.0"},
	  "servers":[{"url":"https://api.attio.com"}],
	  "security":[{"account":[]}],
	  "components":{"securitySchemes":{"account":{"type":"http","scheme":"basic"}}},
	  "paths":{"/records":{"get":{"operationId":"list","responses":{
	    "200":{"description":"ok","content":{"application/json":{"schema":{"type":"string"}}}}
	  }}}}
	}`
	provider, ok := catalog.FindBuiltInProvider("attio")
	if !ok {
		t.Fatal("missing reviewed provider fixture")
	}
	source := catalogUpdateReferenceSources(provider)[0]
	seedCatalogUpdateFixture(t, cfg, source, original)

	body := strings.Replace(original, `"1.0"`, `"2.0"`, 1)
	runContext, cancel := context.WithCancel(ctx)
	defer cancel()
	stop := true
	heads, gets := 0, 0
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != source.URL {
			t.Fatalf("unexpected request: %s", r.URL)
		}
		header := make(http.Header)
		header.Set("Content-Length", strconv.Itoa(len(body)))
		reader := ""
		if r.Method == http.MethodHead {
			heads++
			if stop {
				cancel()
				return nil, context.Canceled
			}
		} else {
			gets++
			reader = body
		}
		return &http.Response{StatusCode: 200, Header: header, ContentLength: int64(len(body)), Body: io.NopCloser(strings.NewReader(reader)), Request: r}, nil
	})}}
	first, err := client.UpdateCatalog(runContext, CatalogUpdateOptions{Provider: "attio", Config: cfg})
	if !errors.Is(err, context.Canceled) || first.RunID == "" || first.Completed {
		t.Fatalf("initial canceled run = %+v, %v", first, err)
	}
	stop = false
	resumed, err := client.UpdateCatalog(ctx, CatalogUpdateOptions{Resume: first.RunID, Config: cfg})
	if err != nil || resumed.RunID != first.RunID || !resumed.Completed ||
		len(resumed.Outcomes) != 2 || resumed.Outcomes[1].Status != "accepted-compatible" {
		t.Fatalf("resumed run = %+v, %v", resumed, err)
	}
	const matchedRequests = 3 // canceled HEAD plus successful HEAD and body
	if heads+gets != matchedRequests || gets != 1 {
		t.Fatalf("fixture requests = %d HEAD, %d GET", heads, gets)
	}
	again, err := client.UpdateCatalog(ctx, CatalogUpdateOptions{Resume: first.RunID, Config: cfg})
	if err != nil || !again.Completed || heads+gets != matchedRequests {
		t.Fatalf("completed resume repeated publisher requests: %+v, %v", again, err)
	}
	body = strings.Replace(body, `"basic"`, `"bearer"`, 1)
	changed, err := client.UpdateCatalog(ctx, CatalogUpdateOptions{Provider: "attio", Config: cfg})
	if err != nil || len(changed.Outcomes) != 2 || changed.Outcomes[1].Status != "review-required" {
		t.Fatalf("changed auth = %+v, %v", changed, err)
	}
	exported, err := client.UpdateCatalog(ctx, CatalogUpdateOptions{Export: true, Config: cfg})
	compatibleBody := strings.Replace(original, `"1.0"`, `"2.0"`, 1)
	wantRevision := fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(compatibleBody)))
	if err != nil || !strings.Contains(string(exported.Export), `"version": "2.0"`) || !strings.Contains(string(exported.Export), wantRevision) {
		t.Fatalf("export = %s, %v", exported.Export, err)
	}
}

func TestCatalogUpdateResumeSkipsRecordedReferenceInPartialProvider(t *testing.T) {
	dir, err := os.MkdirTemp(".", ".catalog-update-checkpoint-fixture-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Error(err)
		}
	})
	cfg := CatalogUpdateConfig{CacheDir: filepath.Join(dir, "cache"), StateDB: filepath.Join(dir, "state.sqlite")}
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	stop := true
	visits := map[string]int{}
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodHead || r.URL.Host != "nocodb.com" {
			t.Fatalf("unexpected publisher request: %s %s", r.Method, r.URL)
		}
		visits[r.URL.Path]++
		if stop && strings.Contains(r.URL.Path, "/v3/") {
			cancel()
			return nil, context.Canceled
		}
		return &http.Response{
			StatusCode: 200, ContentLength: -1, Header: make(http.Header),
			Body: io.NopCloser(strings.NewReader("")), Request: r,
		}, nil
	})}}
	partial, err := client.UpdateCatalog(ctx, CatalogUpdateOptions{Provider: "nocodb", Config: cfg})
	if !errors.Is(err, context.Canceled) || partial.Completed || partial.RunID == "" {
		t.Fatalf("partial run = %+v, %v", partial, err)
	}
	stop = false
	final, err := client.UpdateCatalog(t.Context(), CatalogUpdateOptions{Resume: partial.RunID, Config: cfg})
	if err != nil || !final.Completed || len(final.Outcomes) != 3 ||
		visits["/apis/v2/swagger-v2.json"] != 1 || visits["/apis/v3/swagger-v3.json"] != 2 {
		t.Fatalf("resume outcomes = %+v, visits = %v, error = %v", final, visits, err)
	}
}

func TestCatalogUpdateReconstructsReviewedImmutableBaseline(t *testing.T) {
	for _, scenario := range []string{"valid baseline", "nonregular cache entry", "changed locator", "unreproduced pin at cached revision"} {
		t.Run(scenario, func(t *testing.T) { testCatalogUpdateImmutableBaseline(t, scenario) })
	}
}

func testCatalogUpdateImmutableBaseline(t *testing.T, scenario string) {
	t.Helper()
	dir, err := os.MkdirTemp(".", ".catalog-update-baseline-fixture-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Error(err)
		}
	})
	const oldSHA = "1111111111111111111111111111111111111111"
	const newSHA = "2222222222222222222222222222222222222222"
	const prefix = "https://raw.githubusercontent.com/fixture/source/"
	const suffix = "/records.json"
	const prior = `{"openapi":"3.0.3","info":{"title":"Records","version":"1.0"},"servers":[{"url":"https://example.com"}],"security":[],"paths":{"/records":{"get":{"responses":{"200":{"description":"ok","content":{"application/json":{"schema":{"type":"string"}}}}}}}}}`
	next := strings.Replace(prior, `"1.0"`, `"2.0"`, 1)
	source := catalogUpdateSource{
		ProviderID: "attio", ReferenceID: "attio-rest-api-openapi", Service: "raw.githubusercontent.com/fixture/source/records.json",
		Kind: catalog.SpecKindOpenAPI, URL: prefix + oldSHA + suffix,
		HeldRevision: fmt.Sprintf("repo:%s sha256:%x", oldSHA, sha256.Sum256([]byte(prior))),
	}
	cfg := CatalogUpdateConfig{ReviewedOperations: map[string][]string{"attio/attio-rest-api-openapi": {"#/paths/~1records/get"}}}
	cfg.CacheDir = filepath.Join(dir, "cache")
	cfg.StateDB = filepath.Join(dir, "state.sqlite")
	var existing catalogUpdateSource
	switch scenario {
	case "changed locator":
		existing = source
		existing.Service = "example.com/records.json"
		existing.URL = "https://example.com/records.json"
		seedCatalogUpdateFixture(t, cfg, existing, prior)
	case "unreproduced pin at cached revision":
		source.URL = prefix + newSHA + suffix
		source.Revision = newSHA
		existing = source
		seedCatalogUpdateFixture(t, cfg, existing, prior)
		authored := strings.Replace(prior, `"1.0"`, `"3.0"`, 1)
		source.HeldRevision = fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(authored)))
	}
	cache, err := catalogupdate.OpenArtifactCache(filepath.Join(dir, "cache"), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	store, err := catalogupdate.OpenStore(t.Context(), filepath.Join(dir, "state.sqlite"), cache)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	run, err := store.StartRun(t.Context(), fmt.Sprintf("%x", sha256.Sum256([]byte("catalog"))), []string{"attio", "auth0"})
	if err != nil {
		t.Fatal(err)
	}
	digestPath := fmt.Sprintf("sha256-%x", sha256.Sum256([]byte(prior)))
	if scenario == "nonregular cache entry" {
		if err := cache.Mkdir(digestPath, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	oldGets := 0
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		var data string
		switch r.URL.Host {
		case "auth0.com":
			return &http.Response{StatusCode: http.StatusMethodNotAllowed, Header: make(http.Header), Body: http.NoBody, Request: r}, nil
		case "api.github.com":
			data = fmt.Sprintf(`[{"sha":"%s"}]`, newSHA)
		case "raw.githubusercontent.com":
			switch r.URL.String() {
			case source.URL:
				data = prior
				if r.Method == http.MethodGet {
					oldGets++
				}
			case prefix + newSHA + suffix:
				data = next
			default:
				t.Fatalf("unexpected file: %s", r.URL)
			}
		default:
			t.Fatalf("unexpected publisher: %s", r.URL)
		}
		length := int64(len(data))
		if r.Method == http.MethodHead {
			data = ""
		}
		header := make(http.Header)
		header.Set("Content-Length", strconv.FormatInt(length, 10))
		return &http.Response{
			StatusCode: http.StatusOK, Header: header, ContentLength: length,
			Body: io.NopCloser(strings.NewReader(data)), Request: r,
		}, nil
	})}}
	out, err := client.checkCatalogUpdateSource(t.Context(), store, run.ID, catalog.Provider{ID: "attio"}, source, cfg, nil)
	if scenario == "changed locator" || scenario == "unreproduced pin at cached revision" {
		if err != nil || out.Status != "review-required" {
			t.Fatalf("unreproduced pin outcome = %+v, %v", out, err)
		}
		preferred, err := store.ReviewedReferencePreference(t.Context(), source.ProviderID, source.ReferenceID)
		if err != nil || preferred.Service != existing.Service || preferred.Locator != existing.URL ||
			preferred.Object.SHA256 != catalogUpdateReviewedDigest(existing.HeldRevision) {
			t.Fatalf("unreproduced pin changed the preference = %+v, %v", preferred, err)
		}
		return
	}
	if scenario == "nonregular cache entry" {
		if err != nil || out.Status != "review-required" || oldGets != 1 {
			t.Fatalf("damaged baseline = %+v, baseline downloaded %d times, error %v", out, oldGets, err)
		}
		info, err := cache.Lstat(digestPath)
		if err != nil || !info.IsDir() {
			t.Fatalf("nonregular cache entry replaced: %v, %v", info, err)
		}
		report, err := client.continueCatalogUpdate(t.Context(), store, CatalogUpdateReport{RunID: run.ID, Providers: run.Providers}, false, cfg)
		if err != nil || !report.Completed || len(report.Outcomes) != 4 ||
			report.Outcomes[1].Status != "review-required" || report.Outcomes[3].Status != "size-unknown" {
			t.Fatalf("resume after damaged baseline = %+v, %v", report, err)
		}
		return
	}
	if err != nil || out.Status != "accepted-compatible" || oldGets != 1 {
		t.Fatalf("current outcome = %+v, baseline downloaded %d times, error %v", out, oldGets, err)
	}
	preferred, err := store.PreferredVersion(t.Context(), source.ProviderID, source.ReferenceID, source.Service)
	if err != nil || preferred.Revision != newSHA || preferred.URL != prefix+newSHA+suffix {
		t.Fatalf("published preference = %+v, %v", preferred, err)
	}
}

func TestCatalogUpdateRefusesVersionRegressionAndPreviews(t *testing.T) {
	for _, tc := range []struct{ held, candidate string }{
		{"2.0", "3.0-beta"},
		{"2.0", "1.9"},
		{"2026-11-18", "2026-04-01"},
		{"2026-11-18.dahlia", "2026-04-01.endive"},
	} {
		if catalogUpdateVersionCanAdvance(tc.held, tc.candidate) {
			t.Fatalf("accepted unstable or older candidate %q replacing %q", tc.candidate, tc.held)
		}
	}
}

const catalogRepairFixture = `{"openapi":"3.0.3","info":{"title":"Records","version":"1.0"},"servers":[{"url":"https://example.com"}],"security":[{"key":[]}],"components":{"securitySchemes":{"key":{"type":"http","scheme":"basic"}}},"paths":{"/records":{"get":{"operationId":"list","responses":{"200":{"description":"ok","content":{"application/json":{"schema":{"type":"string"}}}}}}}}}`

func newCatalogRepairFixture(t *testing.T) CatalogUpdateConfig {
	t.Helper()
	dir, err := os.MkdirTemp(".", ".catalog-update-repair-fixture-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Error(err)
		}
	})
	return CatalogUpdateConfig{CacheDir: filepath.Join(dir, "cache"), StateDB: filepath.Join(dir, "state.sqlite")}
}

func catalogRepairResponse(r *http.Request, spec string) *http.Response {
	header := make(http.Header)
	length := int64(len(spec))
	header.Set("Content-Length", strconv.FormatInt(length, 10))
	if r.Method == http.MethodHead {
		spec = ""
	}
	return &http.Response{StatusCode: http.StatusOK, Header: header, ContentLength: length,
		Body: io.NopCloser(strings.NewReader(spec)), Request: r}
}

func TestCatalogUpdateAdoptsReviewedPinAfterLocatorAndAuthChange(t *testing.T) {
	ctx := t.Context()
	cfg := newCatalogRepairFixture(t)
	cfg.ReviewedOperations = map[string][]string{"example/records": {"#/paths/~1records/get"}}
	source := catalogUpdateSource{ProviderID: "example", ReferenceID: "records", Service: "example.com/old/spec.json",
		Kind: catalog.SpecKindOpenAPI, URL: "https://example.com/old/spec.json",
		HeldRevision: fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(catalogRepairFixture)))}
	provider := catalog.Provider{ID: source.ProviderID}
	seedCatalogUpdateFixture(t, cfg, source, catalogRepairFixture)
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
	checked := strings.Replace(catalogRepairFixture, `"basic"`, `"bearer"`, 1)
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		return catalogRepairResponse(r, checked), nil
	})}}
	held, err := store.PreferredVersion(ctx, source.ProviderID, source.ReferenceID, source.Service)
	if err != nil {
		t.Fatal(err)
	}
	check := func(current catalogUpdateSource) string {
		t.Helper()
		run, err := store.StartRun(ctx, fmt.Sprintf("%x", sha256.Sum256([]byte("catalog"))), []string{source.ProviderID})
		if err != nil {
			t.Fatal(err)
		}
		out, err := client.checkCatalogUpdateSource(ctx, store, run.ID, provider, current, cfg, nil)
		if err != nil {
			t.Fatal(err)
		}
		return out.Status
	}
	if status := check(source); status != "review-required" {
		t.Fatalf("unreviewed auth change = %q", status)
	}
	repinned := source
	repinned.URL = "https://example.com/new/spec.json"
	repinned.Service = "example.com/new/spec.json"
	repinned.HeldRevision = fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(checked)))
	if status := check(repinned); status != "reviewed-pin-reproduced" {
		t.Fatalf("reviewed repin = %q", status)
	}
	preferred, err := store.PreferredVersion(ctx, source.ProviderID, source.ReferenceID, source.Service)
	if err != nil || preferred.URL != repinned.URL || preferred.Locator != repinned.URL || preferred.Object == held.Object {
		t.Fatalf("repinned preference = %+v, %v", preferred, err)
	}
	d2 := preferred
	reviewedBytes := checked
	checked = strings.Replace(checked, `"1.0"`, `"2.0"`, 1)
	if status := check(repinned); status != "accepted-compatible" {
		t.Fatalf("release after locator change = %q", status)
	}
	d3, err := store.PreferredVersion(ctx, source.ProviderID, source.ReferenceID, source.Service)
	if err != nil || d3.Object == d2.Object || d3.CatalogPin != d2.Object.SHA256 {
		t.Fatalf("new preference after relocation = %+v, %v", d3, err)
	}
	checked = reviewedBytes
	if status := check(repinned); status != "review-required" {
		t.Fatalf("standing pin returned after newer local release = %q", status)
	}
	moved := repinned
	moved.URL = "https://example.com/moved/spec.json"
	moved.Service = "example.com/moved/spec.json"
	if status := check(moved); status != "review-required" {
		t.Fatalf("standing pin returned at a different locator after newer local release = %q", status)
	}
	current, err := store.ReviewedReferencePreference(ctx, source.ProviderID, source.ReferenceID)
	if err != nil || current != d3 {
		t.Fatalf("standing pin replaced the current preference: %+v, %v", current, err)
	}
	authoredNext := strings.Replace(reviewedBytes, `"1.0"`, `"4.0"`, 1)
	repinned.HeldRevision = fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(authoredNext)))
	checked = strings.Replace(reviewedBytes, `"1.0"`, `"3.0"`, 1)
	if status := check(repinned); status != "review-required" {
		t.Fatalf("update bypassed an unreproduced new catalog pin: %q", status)
	}
	checked = authoredNext
	if status := check(repinned); status != "reviewed-pin-reproduced" {
		t.Fatalf("exact reproduction of the new catalog pin = %q", status)
	}
	d2File, err := store.OpenObject(ctx, d2.Object)
	if err != nil {
		t.Fatal(err)
	}
	defer d2File.Close()
	old, err := store.OpenObject(ctx, held.Object)
	if err != nil {
		t.Fatal(err)
	}
	defer old.Close()
	if body, err := io.ReadAll(old); err != nil || string(body) != catalogRepairFixture {
		t.Fatalf("previous pin = %q, %v", body, err)
	}
}

func TestCatalogUpdateMigratedPreferenceRequiresReproducedNewerPin(t *testing.T) {
	ctx := t.Context()
	cfg := newCatalogRepairFixture(t)
	cfg.ReviewedOperations = map[string][]string{"example/records": {"#/paths/~1records/get"}}
	source := catalogUpdateSource{ProviderID: "example", ReferenceID: "records", Service: "example.com/spec.json",
		Kind: catalog.SpecKindOpenAPI, URL: "https://example.com/spec.json"}
	// Persist a preference without a recorded catalog pin, as v2 migration
	// does. The held bytes kept the release version but changed the title.
	localNewer := strings.Replace(catalogRepairFixture, "Records", "Records current", 1)
	seedCatalogUpdateFixture(t, cfg, source, localNewer)
	source.HeldRevision = fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(catalogRepairFixture)))
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
	run, err := store.StartRun(ctx, fmt.Sprintf("%x", sha256.Sum256([]byte("catalog"))), []string{source.ProviderID})
	if err != nil {
		t.Fatal(err)
	}
	served := catalogRepairFixture
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		return catalogRepairResponse(r, served), nil
	})}}
	result, err := client.checkCatalogUpdateSource(ctx, store, run.ID, catalog.Provider{ID: source.ProviderID}, source, cfg, nil)
	if err != nil || result.Status != "review-required" {
		t.Fatalf("same-version replay after migration = %+v, %v", result, err)
	}
	held, err := store.PreferredVersion(ctx, source.ProviderID, source.ReferenceID, source.Service)
	if err != nil || held.Object.SHA256 != fmt.Sprintf("%x", sha256.Sum256([]byte(localNewer))) {
		t.Fatalf("migration replay changed the preference: %+v, %v", held, err)
	}
	authored := strings.Replace(catalogRepairFixture, `"1.0"`, `"2.0"`, 1)
	authored = strings.Replace(authored, `"basic"`, `"bearer"`, 1)
	source.HeldRevision = fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(authored)))
	served = strings.Replace(catalogRepairFixture, `"1.0"`, `"3.0"`, 1)
	run, err = store.StartRun(ctx, fmt.Sprintf("%x", sha256.Sum256([]byte("changed catalog"))), []string{source.ProviderID})
	if err != nil {
		t.Fatal(err)
	}
	result, err = client.checkCatalogUpdateSource(ctx, store, run.ID, catalog.Provider{ID: source.ProviderID}, source, cfg, nil)
	if err != nil || result.Status != "review-required" {
		t.Fatalf("migration accepted a compatible successor before reproducing the new pin: %+v, %v", result, err)
	}
	current, err := store.PreferredVersion(ctx, source.ProviderID, source.ReferenceID, source.Service)
	if err != nil || current != held {
		t.Fatalf("unreproduced pin changed the migrated preference: %+v, %v", current, err)
	}
	served = authored
	run, err = store.StartRun(ctx, fmt.Sprintf("%x", sha256.Sum256([]byte("changed catalog"))), []string{source.ProviderID})
	if err != nil {
		t.Fatal(err)
	}
	result, err = client.checkCatalogUpdateSource(ctx, store, run.ID, catalog.Provider{ID: source.ProviderID}, source, cfg, nil)
	if err != nil || result.Status != "reviewed-pin-reproduced" {
		t.Fatalf("reproduced newer pin after migration: %+v, %v", result, err)
	}
}

func TestCatalogUpdateContinuesPastMissingOrCorruptCacheObject(t *testing.T) {
	for _, tc := range []string{"missing baseline", "corrupt matching object"} {
		t.Run(tc, func(t *testing.T) {
			ctx := t.Context()
			cfg := newCatalogRepairFixture(t)
			cfg.ReviewedOperations = map[string][]string{"attio/attio-rest-api-openapi": {"#/paths/~1records/get"}}
			provider, ok := catalog.FindBuiltInProvider("attio")
			if !ok {
				t.Fatal("Attio is required for the reviewed catalog fixture")
			}
			source := catalogUpdateReferenceSources(provider)[0]
			seedCatalogUpdateFixture(t, cfg, source, catalogRepairFixture)
			cached := filepath.Join(cfg.CacheDir, fmt.Sprintf("sha256-%x", sha256.Sum256([]byte(catalogRepairFixture))))
			body := catalogRepairFixture
			want := "accepted-compatible"
			if tc == "missing baseline" {
				if err := os.Remove(cached); err != nil {
					t.Fatal(err)
				}
				body = strings.Replace(body, `"1.0"`, `"2.0"`, 1)
				want = "review-required"
			} else if err := os.WriteFile(cached, []byte("truncated"), 0o600); err != nil {
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
			defer store.Close()
			run, err := store.StartRun(ctx, fmt.Sprintf("%x", sha256.Sum256([]byte("catalog"))), []string{"attio", "auth0"})
			if err != nil {
				t.Fatal(err)
			}
			client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
				if r.URL.Host == "auth0.com" && r.Method == http.MethodHead {
					return &http.Response{StatusCode: http.StatusMethodNotAllowed, Header: make(http.Header), Body: http.NoBody, Request: r}, nil
				}
				if r.URL.String() != source.URL {
					t.Fatalf("unexpected fixture request: %s", r.URL)
				}
				return catalogRepairResponse(r, body), nil
			})}}
			report, err := client.continueCatalogUpdate(ctx, store, CatalogUpdateReport{RunID: run.ID, Providers: run.Providers}, false, cfg)
			if err != nil || !report.Completed || len(report.Outcomes) != 4 ||
				report.Outcomes[1].Status != want || report.Outcomes[3].Status != "size-unknown" {
				t.Fatalf("report after cache issue = %+v, %v", report, err)
			}
			if tc == "corrupt matching object" {
				if actual, err := os.ReadFile(cached); err != nil || string(actual) != catalogRepairFixture {
					t.Fatalf("restored bytes = %q, %v", actual, err)
				}
			}
		})
	}
}

func seedCatalogUpdateFixture(t *testing.T, cfg CatalogUpdateConfig, source catalogUpdateSource, body string) {
	t.Helper()
	cache, err := catalogupdate.OpenArtifactCache(cfg.CacheDir, "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	store, err := catalogupdate.OpenStore(t.Context(), cfg.StateDB, cache)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	object := catalogupdate.Object{SHA256: fmt.Sprintf("%x", sha256.Sum256([]byte(body))), Bytes: int64(len(body))}
	if err := store.PutObject(t.Context(), object, strings.NewReader(body)); err != nil {
		t.Fatal(err)
	}
	accepted := catalogupdate.Version{
		ProviderID: source.ProviderID, ReferenceID: source.ReferenceID, Service: source.Service,
		Locator: source.URL, URL: source.URL, Version: "1.0", Revision: source.Revision,
		CheckedAt: time.Date(2026, 11, 25, 0, 0, 0, 0, time.UTC), Object: object,
	}
	if hasCatalogUpdateReviewedDigest(source.HeldRevision, object.SHA256) {
		accepted.CatalogPin = object.SHA256
	}
	if err := store.AcceptVersion(t.Context(), accepted, nil); err != nil {
		t.Fatal(err)
	}
}
