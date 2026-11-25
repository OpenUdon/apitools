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
		ProviderID: "fixture", ReferenceID: "records", Service: "raw.githubusercontent.com/fixture/source/records.json",
		Kind: catalog.SpecKindOpenAPI, URL: prefix + oldSHA + suffix,
		HeldRevision: fmt.Sprintf("repo:%s sha256:%x", oldSHA, sha256.Sum256([]byte(prior))),
	}
	cfg := CatalogUpdateConfig{ReviewedOperations: map[string][]string{"fixture/records": {"#/paths/~1records/get"}}}
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
	run, err := store.StartRun(t.Context(), fmt.Sprintf("%x", sha256.Sum256([]byte("catalog"))), []string{"fixture"})
	if err != nil {
		t.Fatal(err)
	}
	oldGets := 0
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		var data string
		switch r.URL.Host {
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
	out, err := client.checkCatalogUpdateSource(t.Context(), store, run.ID, catalog.Provider{ID: "fixture"}, source, cfg)
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
	if err := store.AcceptVersion(t.Context(), accepted, nil); err != nil {
		t.Fatal(err)
	}
}
