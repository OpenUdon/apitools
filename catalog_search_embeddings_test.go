package apitools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/OpenUdon/apitools/catalog"
)

func reviewedEmbeddingFixture(t *testing.T, query string, responseStatus int) (*httptest.Server, *atomic.Int64) {
	t.Helper()
	calls := &atomic.Int64{}
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if responseStatus != http.StatusOK {
			http.Error(w, "secret fixture body", responseStatus)
			return
		}
		if r.Method != http.MethodPost || r.URL.Path != "/v1/embeddings" || r.Header.Get("Authorization") != "" || r.Header.Get("Cookie") != "" {
			t.Errorf("unexpected request method, path or credential headers: %s %s", r.Method, r.URL.Path)
			http.Error(w, "unexpected request", http.StatusBadRequest)
			return
		}
		var body struct {
			Model string   `json:"model"`
			Input []string `json:"input"`
		}
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&body); err != nil || body.Model != "fixture-model" || len(body.Input) == 0 || len(body.Input) > reviewedEmbeddingBatch {
			t.Errorf("malformed embeddings request: %s, %d, %v", body.Model, len(body.Input), err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		data := make([]map[string]any, 0, len(body.Input))
		for i := len(body.Input) - 1; i >= 0; i-- {
			input := body.Input[i]
			if input != query && (strings.Contains(input, "://") || strings.Contains(input, "#/paths/") || strings.Contains(input, "sha256:") || strings.Contains(input, "https://private")) {
				t.Errorf("unexpected provenance in an embedding input: %q", input)
				http.Error(w, "private data", http.StatusBadRequest)
				return
			}
			vector := []float32{0, 1}
			if strings.HasPrefix(input, "Stripe\n") || input == query {
				vector = []float32{1, 0}
			}
			data = append(data, map[string]any{"index": i, "embedding": vector})
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": data})
	}))
	t.Cleanup(server.Close)
	return server, calls
}

func TestReviewedEmbeddingsUsesPublicTextPersistsProvidersAndReusesCache(t *testing.T) {
	const query = "reverse a card charge"
	server, calls := reviewedEmbeddingFixture(t, query, http.StatusOK)
	client := &Client{AllowUnsafeHosts: true, HTTPClient: server.Client()}
	options := ReviewedProviderSearchOptions{
		Query: query, IndexDB: filepath.Join(t.TempDir(), "reviewed.sqlite"), Limit: 1,
		Embeddings: &ReviewedEmbeddingConfig{
			Endpoint: server.URL + "/v1/embeddings", Model: "fixture-model",
			MaxRequests: 16, MaxInputs: 350, MaxBytes: 1 << 20,
		},
	}
	inputs := len(catalog.BuiltInProviderSearchProfiles()) + 1
	expectedRequests := int64((inputs-1+reviewedEmbeddingBatch-1)/reviewedEmbeddingBatch + 1)
	got, err := client.SearchReviewedProviders(t.Context(), options)
	if err != nil || got.Embeddings == nil || got.Embeddings.Status != "used" ||
		got.Embeddings.Inputs != inputs || calls.Load() != expectedRequests ||
		len(got.Results) != 1 || got.Results[0].Profile.ProviderID != "stripe" ||
		got.Results[0].Match != ReviewedProviderSemanticMatch || len(got.Results[0].MatchedCapabilities) != 0 {
		t.Fatalf("semantic lookup: %q, %+v, calls %d, %v", got.Results, got.Embeddings, calls.Load(), err)
	}
	options.Embeddings.MaxInputs = 1
	options.Embeddings.MaxRequests = 1
	reused, err := client.SearchReviewedProviders(t.Context(), options)
	if err != nil || reused.Embeddings.Status != "used" || reused.Embeddings.Requests != 1 ||
		reused.Embeddings.Inputs != 1 || reused.Embeddings.Bytes != len(query) || calls.Load() != expectedRequests+1 {
		t.Fatalf("cache reuse: %+v, calls %d, %v", reused.Embeddings, calls.Load(), err)
	}
	options.Query = query + "\xff"
	options.Embeddings.MaxBytes = len(options.Query)
	malformed, err := client.SearchReviewedProviders(t.Context(), options)
	if err != nil || malformed.Embeddings.Status != "invalid-query-text" || calls.Load() != expectedRequests+1 {
		t.Fatalf("malformed warm query exceeded its byte budget: %+v, calls %d, %v", malformed.Embeddings, calls.Load(), err)
	}
	options.Query = query
	options.Embeddings.MaxBytes = 1 << 20
	options.Embeddings.Model = "another-model"
	uncached, err := client.SearchReviewedProviders(t.Context(), options)
	if err != nil || uncached.Embeddings.Status != "budget-exhausted" || calls.Load() != expectedRequests+1 {
		t.Fatalf("cache reused another model: %+v, calls %d, %v", uncached.Embeddings, calls.Load(), err)
	}
}

func TestReviewedEmbeddingsRefreshOnlyChangedPublicText(t *testing.T) {
	const query = "reverse a card charge"
	server, calls := reviewedEmbeddingFixture(t, query, http.StatusOK)
	client := &Client{AllowUnsafeHosts: true, HTTPClient: server.Client()}
	db, err := openReviewedProviderIndex("")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	options := ReviewedProviderSearchOptions{Embeddings: &ReviewedEmbeddingConfig{
		Endpoint: server.URL + "/v1/embeddings", Model: "fixture-model", MaxRequests: 3, MaxInputs: 3, MaxBytes: 4096,
	}}
	profiles := []catalog.ProviderSearchProfile{
		{ProviderID: "one", DisplayName: "One", Description: "Approved text", Coverage: catalog.ProviderSearchCoverage{
			Sources: []catalog.ProviderSearchSource{{SourceURL: "https://private.test/artifact", Revision: "private-evidence",
				Operations: []catalog.ProviderSearchOperation{{Label: "List records", Selector: "#/paths/~1records/get"}}}},
		}},
		{ProviderID: "two", DisplayName: "Two", Description: "More approved text"},
	}
	_, first := client.rankReviewedEmbeddings(t.Context(), db, options, profiles, nil, query)
	if first.Status != "used" || first.Inputs != 3 {
		t.Fatalf("initial embedding: %+v", first)
	}
	profiles[0].Description = "Revised approved text"
	profiles[1].Coverage.Sources = []catalog.ProviderSearchSource{{SourceURL: "https://private.test/review-only"}}
	_, changed := client.rankReviewedEmbeddings(t.Context(), db, options, profiles, nil, query)
	if changed.Status != "used" || changed.Inputs != 2 || calls.Load() != 4 {
		t.Fatalf("selective rebuild: %+v, calls %d", changed, calls.Load())
	}
}

func TestReviewedEmbeddingsFallsBackWithoutProviderOperations(t *testing.T) {
	server, calls := reviewedEmbeddingFixture(t, "send slack message", http.StatusServiceUnavailable)
	client := &Client{AllowUnsafeHosts: true, HTTPClient: server.Client()}
	valid := ReviewedEmbeddingConfig{
		Endpoint: server.URL + "/v1/embeddings", Model: "fixture-model", MaxRequests: 16, MaxInputs: 350, MaxBytes: 1 << 20,
	}
	for _, tc := range []struct {
		name, query, wantStatus string
		config                  *ReviewedEmbeddingConfig
		offline                 bool
		calls                   int64
	}{
		{name: "absent", query: "send slack message"},
		{name: "unconfigured", query: "send slack message", config: &ReviewedEmbeddingConfig{}, wantStatus: "unconfigured"},
		{name: "offline", query: "send slack message", config: &valid, offline: true, wantStatus: "offline"},
		{name: "exact", query: "slack web api", config: &valid, wantStatus: "exact-provider"},
		{name: "budget", query: "send slack message", config: &ReviewedEmbeddingConfig{
			Endpoint: valid.Endpoint, Model: valid.Model, MaxRequests: 1, MaxInputs: 1, MaxBytes: 32,
		}, wantStatus: "budget-exhausted"},
		{name: "cleartext", query: "send slack message", config: &ReviewedEmbeddingConfig{
			Endpoint: strings.Replace(valid.Endpoint, "https:", "http:", 1),
			Model:    valid.Model, MaxRequests: 16, MaxInputs: 350, MaxBytes: 1 << 20,
		}, wantStatus: "invalid-configuration"},
		{name: "query credentials", query: "send slack message", config: &ReviewedEmbeddingConfig{
			Endpoint: valid.Endpoint + "?sig=fixture-only", Model: valid.Model, MaxRequests: 16, MaxInputs: 350, MaxBytes: 1 << 20,
		}, wantStatus: "invalid-configuration"},
		{name: "invalid model text", query: "send slack message", config: &ReviewedEmbeddingConfig{
			Endpoint: valid.Endpoint, Model: valid.Model + "\xff", MaxRequests: 16, MaxInputs: 350, MaxBytes: 1 << 20,
		}, wantStatus: "invalid-configuration"},
		{name: "publisher failure", query: "send slack message", config: &valid, wantStatus: "endpoint-unavailable", calls: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			before := calls.Load()
			got, err := client.SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{
				Query: tc.query, Offline: tc.offline, Embeddings: tc.config, Limit: 1,
			})
			if err != nil || len(got.Results) != 1 || got.Results[0].Profile.ProviderID != "slack" || calls.Load()-before != tc.calls {
				t.Fatalf("fallback: %d results, %d calls, %v", len(got.Results), calls.Load()-before, err)
			}
			if tc.wantStatus == "" {
				if got.Embeddings != nil {
					t.Fatalf("unconfigured lookup has usage: %+v", got.Embeddings)
				}
			} else if got.Embeddings == nil || got.Embeddings.Status != tc.wantStatus {
				t.Fatalf("status %+v, want %s", got.Embeddings, tc.wantStatus)
			}
		})
	}
}

func TestReviewedEmbeddingsRedirectAndMalformedResponse(t *testing.T) {
	target, calls := reviewedEmbeddingFixture(t, "send slack message", http.StatusOK)
	for _, fail := range []struct {
		name    string
		handler http.HandlerFunc
	}{
		{"redirect", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, target.URL+"/v1/embeddings", http.StatusTemporaryRedirect)
		}},
		{"invalid vector", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `{"data":[{"index":0,"embedding":[0,0]}]}`)
		}},
	} {
		t.Run(fail.name, func(t *testing.T) {
			server := httptest.NewTLSServer(fail.handler)
			defer server.Close()
			client := &Client{AllowUnsafeHosts: true, HTTPClient: server.Client()}
			got, err := client.SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{
				Query: "send slack message", Limit: 1,
				Embeddings: &ReviewedEmbeddingConfig{Endpoint: server.URL + "/v1/embeddings", Model: "fixture-model",
					MaxRequests: 16, MaxInputs: 350, MaxBytes: 1 << 20},
			})
			if err != nil || got.Embeddings.Status != "endpoint-unavailable" || got.Embeddings.Requests != 1 ||
				got.Results[0].Match != ReviewedProviderOperationMatch || calls.Load() != 0 {
				t.Fatalf("unsafe response handling: %+v, redirected calls %d, %v", got.Embeddings, calls.Load(), err)
			}
		})
	}
}

func TestReviewedEmbeddingsDoesNotCacheIncompatibleDimensions(t *testing.T) {
	const query = "reverse charge"
	var calls atomic.Int64
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Input []string `json:"input"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode fixture: %v", err)
			return
		}
		vector := []float32{1, 0}
		if len(body.Input) == 1 && body.Input[0] == query {
			vector = []float32{1, 0, 0}
		}
		calls.Add(1)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []any{map[string]any{"index": 0, "embedding": vector}},
		})
	}))
	defer server.Close()
	client := &Client{AllowUnsafeHosts: true, HTTPClient: server.Client()}
	db, err := openReviewedProviderIndex("")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	options := ReviewedProviderSearchOptions{Embeddings: &ReviewedEmbeddingConfig{
		Endpoint: server.URL + "/v1/embeddings", Model: "fixture-model", MaxRequests: 2, MaxInputs: 2, MaxBytes: 4096,
	}}
	profiles := []catalog.ProviderSearchProfile{{ProviderID: "one", DisplayName: "One", Description: "Approved text"}}
	_, usage := client.rankReviewedEmbeddings(t.Context(), db, options, profiles, nil, query)
	var cached int
	if err := db.QueryRow("SELECT count(*) FROM search_embedding").Scan(&cached); err != nil {
		t.Fatal(err)
	}
	if usage.Status != "inconsistent-dimensions" || usage.Requests != 2 || calls.Load() != 2 || cached != 0 {
		t.Fatalf("incompatible response was cached or used: %+v, calls %d, cached %d", usage, calls.Load(), cached)
	}
}

func TestReviewedEmbeddingsResumesValidatedBatchesAndReplacesOldDimensions(t *testing.T) {
	const query = "reverse charge"
	var calls, failAt, dimensions atomic.Int64
	failAt.Store(3)
	dimensions.Store(2)
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == failAt.Load() {
			http.Error(w, "temporary fixture failure", http.StatusServiceUnavailable)
			return
		}
		var body struct {
			Input []string `json:"input"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode fixture: %v", err)
			return
		}
		data := make([]map[string]any, len(body.Input))
		for i, text := range body.Input {
			vector := make([]float32, dimensions.Load())
			vector[0] = 1
			if text != query {
				vector[1] = 1
			}
			data[i] = map[string]any{"index": i, "embedding": vector}
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": data})
	}))
	defer server.Close()
	client := &Client{AllowUnsafeHosts: true, HTTPClient: server.Client()}
	db, err := openReviewedProviderIndex("")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	cfg := &ReviewedEmbeddingConfig{Endpoint: server.URL, Model: "fixture-model", MaxRequests: 3, MaxInputs: 34, MaxBytes: 4096}
	options := ReviewedProviderSearchOptions{Embeddings: cfg}
	profiles := make([]catalog.ProviderSearchProfile, reviewedEmbeddingBatch+1)
	for i := range profiles {
		id := fmt.Sprintf("p%d", i)
		profiles[i] = catalog.ProviderSearchProfile{ProviderID: id, DisplayName: id}
	}
	profiles[0].Coverage.Status = catalog.ProviderSearchSelectedSubset
	_, initial := client.rankReviewedEmbeddings(t.Context(), db, options, profiles, nil, query)
	if initial.Status != "endpoint-unavailable" || initial.Requests != 3 {
		t.Fatalf("fixture did not fail the second provider batch: %+v", initial)
	}
	cfg.MaxRequests, cfg.MaxInputs = 2, 2
	ranked, resumed := client.rankReviewedEmbeddings(t.Context(), db, options, profiles, nil, query)
	if resumed.Status != "used" || resumed.Requests != 2 || resumed.Inputs != 2 || len(ranked) != 0 {
		t.Fatalf("did not resume from the complete batch with the default similarity: %+v, %v", resumed, ranked)
	}
	cfg.MinSimilarity = 0.7
	cfg.MaxRequests, cfg.MaxInputs = 1, 1
	ranked, adjusted := client.rankReviewedEmbeddings(t.Context(), db, options, profiles, nil, query)
	if adjusted.Status != "used" || adjusted.Requests != 1 || len(ranked) != 1 || ranked[0].result.Match != ReviewedProviderSemanticMatch {
		t.Fatalf("configured similarity did not retain the mid-range match: %+v, %v", adjusted, ranked)
	}
	dimensions.Store(3)
	_, limited := client.rankReviewedEmbeddings(t.Context(), db, options, profiles, nil, query)
	if limited.Status != "budget-exhausted" || limited.Requests != 1 || limited.Inputs != 1 {
		t.Fatalf("dimension change exceeded the budget: %+v", limited)
	}
	cfg.MaxRequests, cfg.MaxInputs = 3, 34
	_, replaced := client.rankReviewedEmbeddings(t.Context(), db, options, profiles, nil, query)
	if replaced.Status != "used" || replaced.Requests != 3 || replaced.Inputs != 34 {
		t.Fatalf("dimension refresh did not replace cached vectors: %+v", replaced)
	}
	cfg.MaxRequests, cfg.MaxInputs = 1, 1
	_, warm := client.rankReviewedEmbeddings(t.Context(), db, options, profiles, nil, query)
	if warm.Status != "used" || warm.Requests != 1 {
		t.Fatalf("replaced vectors did not remain reusable: %+v", warm)
	}
}

func TestReviewedEmbeddingsDoesNotReuseAmbientTransport(t *testing.T) {
	const query = "send slack message"
	server, serviceCalls := reviewedEmbeddingFixture(t, query, http.StatusOK)
	var ambientCalls atomic.Int64
	client := &Client{
		AllowUnsafeHosts: true,
		HTTPClient: &http.Client{
			Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
				ambientCalls.Add(1)
				r.Header.Set("Authorization", "Bearer ambient-fixture-only")
				return server.Client().Transport.RoundTrip(r)
			}),
		},
	}
	got, err := client.SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{
		Query: query, Limit: 1,
		Embeddings: &ReviewedEmbeddingConfig{
			Endpoint: server.URL + "/v1/embeddings", Model: "fixture-model",
			MaxRequests: 16, MaxInputs: 350, MaxBytes: 1 << 20,
		},
	})
	if err != nil || got.Embeddings == nil || got.Embeddings.Status != "endpoint-unavailable" ||
		len(got.Results) != 1 || ambientCalls.Load() != 0 || serviceCalls.Load() != 0 {
		t.Fatalf("embedding reused an ambient transport: %+v, ambient calls %d, service calls %d, %v",
			got.Embeddings, ambientCalls.Load(), serviceCalls.Load(), err)
	}
}
