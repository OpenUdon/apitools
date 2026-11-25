package apitools

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/catalogupdate"
)

func TestCatalogUpdateAnchorsVerifiedHeldDigest(t *testing.T) {
	for _, fixture := range []struct {
		name                  string
		migrated, sameVersion bool
	}{{name: "recorded anchor"}, {name: "migrated", migrated: true}, {name: "same native version", sameVersion: true}} {
		t.Run(fixture.name, func(t *testing.T) {
			ctx := t.Context()
			cfg := newCatalogRepairFixture(t)
			cfg.ReviewedOperations = map[string][]string{"example/records": {"#/paths/~1records/get"}}
			source := catalogUpdateSource{ProviderID: "example", ReferenceID: "records", Service: "example.com/spec.json",
				Kind: catalog.SpecKindOpenAPI, URL: "https://example.com/spec.json"}
			first := catalogRepairFixture
			second := strings.Replace(first, `"1.0"`, `"2.0"`, 1)
			third := strings.Replace(first, `"1.0"`, `"3.0"`, 1)
			if fixture.sameVersion {
				second = strings.Replace(first, `"Records"`, `"Records second revision"`, 1)
				third = strings.Replace(first, `"Records"`, `"Records third revision"`, 1)
			}
			source.HeldRevision = fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(first)))
			if fixture.migrated {
				source.HeldRevision = ""
			}
			seedCatalogUpdateFixture(t, cfg, source, first)
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
			body := second
			client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
				if r.URL.String() != source.URL {
					t.Fatalf("unexpected fixture request: %s", r.URL)
				}
				return catalogRepairResponse(r, body), nil
			})}}
			run := func(want string) {
				t.Helper()
				run, err := store.StartRun(ctx, catalog.BuiltInCatalogManifest().CatalogSHA256, []string{source.ProviderID})
				if err != nil {
					t.Fatal(err)
				}
				out, err := client.checkCatalogUpdateSource(ctx, store, run.ID, catalog.Provider{}, source, cfg, nil)
				if err != nil || out.Status != want {
					t.Fatalf("outcome = %s, %s, %v", out.Status, out.Evidence, err)
				}
			}
			run("accepted-compatible")
			newPin := fmt.Sprintf("%x", sha256.Sum256([]byte(second)))
			source.HeldRevision = "sha256:" + newPin
			body = third
			run("accepted-compatible")
			body = second
			run("review-required")
			body = first
			run("review-required")
			source.HeldRevision = fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(first)))
			body = first
			run("review-required")
			pref, err := store.PreferredVersion(ctx, source.ProviderID, source.ReferenceID, source.Service)
			if err != nil || pref.Object.SHA256 != fmt.Sprintf("%x", sha256.Sum256([]byte(third))) || pref.CatalogPin != newPin {
				t.Fatalf("preferred version = %+v, %v", pref, err)
			}
		})
	}
}

func TestCatalogUpdateReproducesAuthoredCommitBeforeUsingNewHead(t *testing.T) {
	for _, tc := range []struct {
		name                             string
		havePreference, compare          bool
		headReusesPin, wrongAuthored     bool
		headFailure                      string
		branchCache, protobuf, relocated bool
		wantStatus, preferredVersion     string
	}{
		{name: "first pin with no reviewed operations", wantStatus: "reviewed-pin-reproduced", preferredVersion: "2.0"},
		{name: "first pin with comparison", compare: true, wantStatus: "accepted-compatible", preferredVersion: "3.0"},
		{name: "new pin with no reviewed operations", havePreference: true, wantStatus: "reviewed-pin-reproduced", preferredVersion: "2.0"},
		{name: "new pin with comparison", havePreference: true, compare: true, wantStatus: "accepted-compatible", preferredVersion: "3.0"},
		{name: "head digest alone is insufficient", headReusesPin: true, wrongAuthored: true, wantStatus: "review-required"},
		{name: "head shares the reproduced digest", headReusesPin: true, wantStatus: "reviewed-pin-reproduced", preferredVersion: "2.0"},
		{name: "metadata unavailable", headFailure: "metadata", wantStatus: catalogupdate.RetryPendingStatus, preferredVersion: "2.0"},
		{name: "authored download throttled", headFailure: "authored-throttle", wantStatus: catalogupdate.RetryPendingStatus},
		{name: "head throttling skips same raw origin", headFailure: "head-throttle", wantStatus: catalogupdate.RetryPendingStatus},
		{name: "metadata and authored origins both throttle", headFailure: "dual-throttle", wantStatus: catalogupdate.RetryPendingStatus},
		{name: "branch updates to an exact locally held commit", branchCache: true, compare: true, wantStatus: "accepted-compatible", preferredVersion: "3.0"},
		{name: "new reviewed file at a different locator", relocated: true, havePreference: true, compare: true, wantStatus: "accepted-compatible", preferredVersion: "3.0"},
		{name: "authored protobuf", protobuf: true, wantStatus: "reviewed-pin-reproduced", preferredVersion: "v26"},
		{name: "authored protobuf with invalid head", protobuf: true, headFailure: "invalid", wantStatus: "publisher-spec-unexamined", preferredVersion: "v26"},
		{name: "head invalid", headFailure: "invalid", havePreference: true, compare: true, wantStatus: "publisher-spec-invalid", preferredVersion: "2.0"},
		{name: "head redirects to the authored commit", headFailure: "redirect", wantStatus: "publisher-source-unexamined", preferredVersion: "2.0"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			const beforeCommit = "1111111111111111111111111111111111111111"
			const authoredCommit = "2222222222222222222222222222222222222222"
			const headCommit = "3333333333333333333333333333333333333333"
			const base = "https://raw.githubusercontent.com/fixture/repo/"
			const path = "/spec.json"
			ctx := t.Context()
			cfg := newCatalogRepairFixture(t)
			if tc.compare {
				cfg.ReviewedOperations = map[string][]string{"example/records": {"#/paths/~1records/get"}}
			}
			source := catalogUpdateSource{ProviderID: "example", ReferenceID: "records", Service: "raw.githubusercontent.com/fixture/repo/spec.json", Kind: catalog.SpecKindOpenAPI}
			previous := catalogRepairFixture
			authored := strings.Replace(previous, `"1.0"`, `"2.0"`, 1)
			head := strings.Replace(previous, `"1.0"`, `"3.0"`, 1)
			if tc.havePreference {
				source.URL, source.Revision = base+beforeCommit+path, beforeCommit
				if tc.relocated {
					source.URL = base + beforeCommit + "/previous.json"
					source.Service = "raw.githubusercontent.com/fixture/repo/previous.json"
				}
				source.HeldRevision = fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(previous)))
				seedCatalogUpdateFixture(t, cfg, source, previous)
				if tc.relocated {
					source.Service = "raw.githubusercontent.com/fixture/repo/spec.json"
				}
			}
			if tc.protobuf {
				source.Kind = catalog.SpecKindGRPCProtobuf
				source.HeldVersion = "v26"
				authored = `syntax = "proto3"; package records.v26; service Records { rpc Get (Request) returns (Response); } message Request { string id = 1; } message Response { string value = 1; }`
				head = authored + "\n// a new commit"
			}
			source.URL = base + authoredCommit + path
			source.HeldRevision = fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(authored)))
			if tc.headReusesPin {
				head = authored
			}
			if tc.headFailure == "invalid" {
				head = `{"not":"openapi"}`
			}
			if tc.wrongAuthored {
				authored = previous
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
			if tc.branchCache {
				held := catalogupdate.Version{
					ProviderID: source.ProviderID, ReferenceID: source.ReferenceID, Service: source.Service,
					Locator: base + "main" + path, URL: source.URL, Revision: authoredCommit, Version: "2.0",
					CheckedAt: time.Now(), CatalogPin: fmt.Sprintf("%x", sha256.Sum256([]byte(previous))),
					Object: catalogupdate.Object{Bytes: int64(len(authored)), SHA256: catalogUpdateReviewedDigest(source.HeldRevision)},
				}
				if err := store.PutObject(ctx, held.Object, strings.NewReader(authored)); err != nil {
					t.Fatal(err)
				}
				if err := store.AcceptVersion(ctx, held, nil); err != nil {
					t.Fatal(err)
				}
			}
			authoredGets, authoredHeads := 0, 0
			client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
				switch r.URL.Host {
				case "api.github.com":
					if tc.headFailure == "dual-throttle" {
						return &http.Response{StatusCode: http.StatusServiceUnavailable, Header: http.Header{"Retry-After": {"120"}}, Body: http.NoBody, Request: r}, nil
					}
					if tc.headFailure == "metadata" {
						return &http.Response{StatusCode: http.StatusTooManyRequests, Header: make(http.Header), Body: http.NoBody, Request: r}, nil
					}
					return catalogRepairResponse(r, fmt.Sprintf(`[{"sha":%q}]`, headCommit)), nil
				case "raw.githubusercontent.com":
					switch r.URL.String() {
					case source.URL:
						if r.Method == http.MethodHead {
							authoredHeads++
						}
						if tc.headFailure == "dual-throttle" {
							return &http.Response{StatusCode: http.StatusTooManyRequests, Header: http.Header{"Retry-After": {"900"}}, Body: http.NoBody, Request: r}, nil
						}
						if tc.headFailure == "authored-throttle" {
							return &http.Response{StatusCode: http.StatusTooManyRequests, Header: make(http.Header), Body: http.NoBody, Request: r}, nil
						}
						if r.Method == http.MethodGet {
							authoredGets++
						}
						return catalogRepairResponse(r, authored), nil
					case base + headCommit + path:
						if tc.headFailure == "head-throttle" {
							return &http.Response{StatusCode: http.StatusTooManyRequests, Header: make(http.Header), Body: http.NoBody, Request: r}, nil
						}
						if tc.headFailure == "redirect" {
							return &http.Response{StatusCode: http.StatusFound, Header: http.Header{"Location": {source.URL}}, Body: http.NoBody, Request: r}, nil
						}
						return catalogRepairResponse(r, head), nil
					}
				}
				t.Fatalf("unexpected fixture request: %s", r.URL)
				return nil, nil
			})}}
			run, err := store.StartRun(ctx, catalog.BuiltInCatalogManifest().CatalogSHA256, []string{source.ProviderID})
			if err != nil {
				t.Fatal(err)
			}
			out, err := client.checkCatalogUpdateSource(ctx, store, run.ID, catalog.Provider{}, source, cfg, nil)
			wantGets := 1
			if tc.branchCache {
				wantGets = 0
			}
			if strings.HasSuffix(tc.headFailure, "-throttle") {
				wantGets = 0
			}
			if tc.headFailure == "redirect" {
				wantGets++
			}
			if err != nil || out.Status != tc.wantStatus || authoredGets != wantGets {
				t.Fatalf("outcome = %s, %s, authored GETs %d, %v", out.Status, out.Evidence, authoredGets, err)
			}
			pref, err := store.ReviewedReferencePreference(ctx, source.ProviderID, source.ReferenceID)
			if tc.wrongAuthored || strings.HasSuffix(tc.headFailure, "-throttle") {
				if err == nil {
					t.Fatalf("unreproduced locator accepted: %+v", pref)
				}
				if tc.headFailure == "head-throttle" && authoredHeads != 0 {
					t.Fatal("attempted authored HEAD on the same throttled origin")
				}
				if tc.headFailure == "dual-throttle" && len(readCatalogUpdateRetryBlocks(out)) != 2 {
					t.Fatalf("distinct origin throttles were lost: %s", out.Evidence)
				}
				if tc.headFailure == "authored-throttle" && readCatalogUpdateRetry(out).Reason != "rate-limit" {
					t.Fatalf("authored throttle was lost: %s", out.Evidence)
				}
				return
			}
			if err != nil || pref.Version != tc.preferredVersion || pref.Locator != source.URL ||
				pref.CatalogPin != catalogUpdateReviewedDigest(source.HeldRevision) || out.Object == nil || pref.Object != *out.Object {
				t.Fatalf("preferred version = %+v, %v, outcome = %s", pref, err, out.Evidence)
			}
			if tc.relocated && pref.Service != "raw.githubusercontent.com/fixture/repo/previous.json" {
				t.Fatalf("locator move changed the existing service key: %+v", pref)
			}
			if tc.preferredVersion == "2.0" && (pref.URL != source.URL || pref.Revision != authoredCommit) {
				t.Fatalf("reviewed locator and revision disagree: %+v", pref)
			}
			outcomes, err := store.ReadOutcomes(ctx, run.ID)
			if err != nil || len(outcomes) != 1 || outcomes[0].Object == nil || *outcomes[0].Object != pref.Object {
				t.Fatalf("missing durable acceptance checkpoint, %v", err)
			}
		})
	}
}

func TestCatalogUpdateRevalidatesHeldBytesUnderCurrentSourceKind(t *testing.T) {
	ctx := t.Context()
	cfg := newCatalogRepairFixture(t)
	source := catalogUpdateSource{ProviderID: "example", ReferenceID: "records", Service: "spec",
		Kind: catalog.SpecKindOpenAPI, URL: "https://example.com/spec"}
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
	held, err := store.PreferredVersion(ctx, source.ProviderID, source.ReferenceID, source.Service)
	if err != nil {
		t.Fatal(err)
	}
	source.Kind = catalog.SpecKindGRPCProtobuf
	source.HeldRevision = "sha256:" + held.Object.SHA256
	baseline, err := (&Client{}).catalogUpdateReviewedBaseline(ctx, store, source, held, source.URL, held.Object.SHA256)
	if err != nil || baseline.Preferred.Object.SHA256 != "" {
		t.Fatalf("accepted cached OpenAPI as protobuf: %+v, %v", baseline, err)
	}
}

func TestCatalogUpdateDiscoveryBaselineUsesNativeVersion(t *testing.T) {
	const body = `{"kind":"discovery#restDescription","discoveryVersion":"v1","name":"fixture","version":"v1","revision":"20261101","rootUrl":"https://example.com/","methods":{"get":{"id":"fixture.get","path":"items","httpMethod":"GET"}}}`
	source := catalogUpdateSource{Kind: catalog.SpecKindGoogleDiscovery}
	prior := catalogupdate.Version{URL: "https://example.com/spec.json", Version: "v3", Revision: strings.Repeat("1", 40)}
	got, valid := qualifyCatalogUpdateBaseline(t.Context(), source, prior, []byte(body))
	if !valid || got.Preferred.Version != "v1" || got.Preferred.Revision != prior.Revision {
		t.Fatalf("native version and immutable provenance = %+v, valid %t", got, valid)
	}
}

func TestCatalogUpdateRecognizesIdenticalBytesAtAnotherCommit(t *testing.T) {
	ctx := t.Context()
	cfg := newCatalogRepairFixture(t)
	const base = "https://raw.githubusercontent.com/v1/v2/"
	const file = "/specs/v3/openapi-v4.json"
	oldCommit, head := strings.Repeat("1", 40), strings.Repeat("2", 40)
	source := catalogUpdateSource{ProviderID: "example", ReferenceID: "records", Service: "records",
		Kind: catalog.SpecKindOpenAPI, URL: base + oldCommit + file, Revision: oldCommit,
		HeldRevision: fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(catalogRepairFixture)))}
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
	run, err := store.StartRun(ctx, catalog.BuiltInCatalogManifest().CatalogSHA256, []string{source.ProviderID})
	if err != nil {
		t.Fatal(err)
	}
	headURL := base + head + file
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host == "api.github.com" {
			return catalogRepairResponse(r, fmt.Sprintf(`[{"sha":%q}]`, head)), nil
		}
		if r.URL.String() != headURL {
			t.Fatalf("unexpected fixture request: %s", r.URL)
		}
		return catalogRepairResponse(r, catalogRepairFixture), nil
	})}}
	out, err := client.checkCatalogUpdateSource(ctx, store, run.ID, catalog.Provider{}, source, cfg, nil)
	if err != nil || out.Status != "unchanged" || out.Object == nil || !strings.Contains(string(out.Evidence), headURL) {
		t.Fatalf("identical publisher revision = %s, %s, %v", out.Status, out.Evidence, err)
	}
	held, err := store.PreferredVersion(ctx, source.ProviderID, source.ReferenceID, source.Service)
	if err != nil || held.URL != source.URL || held.Revision != oldCommit || held.Object != *out.Object {
		t.Fatalf("identical observation changed preference: %+v, %v", held, err)
	}
}

func TestCatalogUpdateIdenticalBytesAtDifferentRedirectDestinationNeedReview(t *testing.T) {
	ctx := t.Context()
	cfg := newCatalogRepairFixture(t)
	source := catalogUpdateSource{ProviderID: "example", ReferenceID: "records", Service: "records",
		Kind: catalog.SpecKindOpenAPI, URL: "https://example.com/docs/latest.json"}
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
	held, err := store.PreferredVersion(ctx, source.ProviderID, source.ReferenceID, source.Service)
	if err != nil {
		t.Fatal(err)
	}
	previous := held
	previous.URL = "https://example.com/docs/orders-v1.json"
	if err := store.AcceptVersion(ctx, previous, &held); err != nil {
		t.Fatal(err)
	}
	const next = "https://example.com/docs/orders-v2.json"
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		switch r.URL.String() {
		case source.URL:
			return &http.Response{StatusCode: http.StatusFound, Header: http.Header{"Location": {next}}, Body: http.NoBody, Request: r}, nil
		case next:
			return catalogRepairResponse(r, catalogRepairFixture), nil
		default:
			t.Fatalf("unexpected fixture request: %s", r.URL)
			return nil, nil
		}
	})}}
	for _, reviewed := range []bool{false, true} {
		if reviewed {
			cfg.ReviewedOperations = map[string][]string{"example/records": {"#/paths/~1records/get"}}
		}
		run, err := store.StartRun(ctx, catalog.BuiltInCatalogManifest().CatalogSHA256, []string{source.ProviderID})
		if err != nil {
			t.Fatal(err)
		}
		out, err := client.checkCatalogUpdateSource(ctx, store, run.ID, catalog.Provider{}, source, cfg, nil)
		if err != nil || out.Status != "review-required" || !strings.Contains(string(out.Evidence), next) {
			t.Fatalf("different redirect destination (reviewed=%t) = %s, %s, %v", reviewed, out.Status, out.Evidence, err)
		}
	}
	actual, err := store.PreferredVersion(ctx, source.ProviderID, source.ReferenceID, source.Service)
	if err != nil || actual.URL != previous.URL {
		t.Fatalf("unrelated redirect replaced the accepted destination: %+v, %v", actual, err)
	}
}
