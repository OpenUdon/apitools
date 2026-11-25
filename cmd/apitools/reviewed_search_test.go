package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/apitools"
)

func TestReviewedSearchCLITextAndJSON(t *testing.T) {
	dir := t.TempDir()
	config := filepath.Join(dir, "search.json")
	if err := os.WriteFile(config, []byte(`{"index_db":"reviewed.sqlite","state_db":"missing.sqlite"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	for _, jsonMode := range []bool{false, true} {
		args := []string{"--query", "send slack message", "--source", "reviewed", "--config", config, "--limit", "1"}
		if jsonMode {
			args = append(args, "--json")
		}
		out.Reset()
		errOut.Reset()
		code := runSearchWithClient(args, &out, &errOut, func(string) (*apitools.Client, func(), error) {
			t.Fatal("explicit reviewed search opened the public directory client")
			return nil, func() {}, nil
		})
		if code != exitSuccess {
			t.Fatalf("exit %d: %s", code, errOut.String())
		}
		if jsonMode {
			var report apitools.ReviewedProviderSearchReport
			if err := json.Unmarshal(out.Bytes(), &report); err != nil {
				t.Fatal(err)
			}
			if report.Source != "reviewed-catalog" || report.PublisherState != "unavailable" ||
				len(report.Results) != 1 || report.Results[0].Profile.ProviderID != "slack" ||
				len(report.Results[0].MatchedCapabilities) != 1 {
				t.Fatalf("unexpected reviewed search: %+v", report)
			}
		} else {
			for _, want := range []string{"Slack", "covered-operation", "selected-subset", "Send Slack message",
				"chat.postMessage", "advisory-overlay", "Catalog review date",
				"Operation source slack-web-api-docs (slack-web-message-history-overlay) sha256:"} {
				if !strings.Contains(out.String(), want) {
					t.Fatalf("text result missing %q:\n%s", want, out.String())
				}
			}
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "reviewed.sqlite")); err != nil {
		t.Fatal(err)
	}
}

func TestReviewedSearchConfigRejectsUnknownFields(t *testing.T) {
	p := filepath.Join(t.TempDir(), "search.json")
	for _, invalid := range []string{`{"metadata_endpoint":"https://example.test"}`, `null`, `{}{}`} {
		if err := os.WriteFile(p, []byte(invalid), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := readReviewedSearchConfig(p); err == nil {
			t.Fatalf("accepted malformed config %s", invalid)
		}
	}
}

func TestReviewedSearchOfflineConfigDisablesEmbeddingRequests(t *testing.T) {
	path := filepath.Join(t.TempDir(), "search.json")
	data := []byte(`{
	  "embeddings": {
	    "endpoint": "https://127.0.0.1/v1/embeddings",
	    "model": "fixture",
	    "max_requests": 16,
	    "max_inputs": 350,
	    "max_bytes": 1048576
	  }
	}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	for _, mode := range [][]string{{"--offline"}, {"--cache-mode", "offline"}} {
		var out, errOut bytes.Buffer
		args := append([]string{"search", "--query", "send slack message", "--config", path, "--json"}, mode...)
		if code := run(args, &out, &errOut); code != exitSuccess {
			t.Fatalf("exit %d: %s", code, errOut.String())
		}
		var report apitools.ReviewedProviderSearchReport
		if err := json.Unmarshal(out.Bytes(), &report); err != nil {
			t.Fatal(err)
		}
		if report.Embeddings == nil || report.Embeddings.Status != "offline" || report.Embeddings.Requests != 0 ||
			len(report.Results) == 0 || report.Results[0].Profile.ProviderID != "slack" {
			t.Fatalf("offline search did not stay lexical: %+v", report.Embeddings)
		}
	}
}

func TestReviewedSearchRejectsAnInvalidCacheMode(t *testing.T) {
	var out, errOut bytes.Buffer
	code := run([]string{"search", "--source", "reviewed", "--query", "send slack message",
		"--cache-mode", "offine", "--config", "missing-config.json"}, &out, &errOut)
	if code != exitRuntime || !strings.Contains(errOut.String(), "unknown cache mode") || out.Len() != 0 {
		t.Fatalf("invalid mode allowed lookup: %d, %q, %q", code, out.String(), errOut.String())
	}
}

func TestReviewedSearchDefaultUsesReviewedEvidence(t *testing.T) {
	for _, tt := range []struct {
		query string
		match apitools.ReviewedProviderMatchKind
		want  map[string]string
	}{
		{"send slack message", apitools.ReviewedProviderOperationMatch, map[string]string{"slack": "#/paths/~1chat.postMessage/post"}},
		{"refund payment", apitools.ReviewedProviderOperationMatch, map[string]string{
			"stripe": "#/paths/~1v1~1refunds/post", "square": "#/paths/~1v2~1refunds/post", "mercado-pago": "#/paths/~1v1~1payments~1{id}~1refunds/post",
		}},
		{"slack web api", apitools.ReviewedProviderExactMatch, map[string]string{"slack": ""}},
	} {
		t.Run(tt.query, func(t *testing.T) {
			var out, errOut bytes.Buffer
			code := runSearchWithClient([]string{"--query", tt.query, "--json", "--offline"}, &out, &errOut, func(string) (*apitools.Client, func(), error) {
				t.Fatal("the default reviewed match opened a public search client")
				return nil, nil, nil
			})
			if code != exitSuccess {
				t.Fatalf("exit %d: %s", code, errOut.String())
			}
			var report apitools.ReviewedProviderSearchReport
			if err := json.Unmarshal(out.Bytes(), &report); err != nil {
				t.Fatal(err)
			}
			if report.Source != "reviewed-catalog" || len(report.Results) == 0 || report.Results[0].Match != tt.match {
				t.Fatalf("unexpected default search: %+v", report)
			}
			found := make(map[string]bool)
			for _, result := range report.Results {
				id := result.Profile.ProviderID
				selector, expected := tt.want[id]
				if !expected || result.Match != tt.match {
					continue
				}
				if selector == "" {
					found[id] = true
				}
				for _, operation := range result.MatchedCapabilities {
					if operation.Selector == selector && len(operation.SHA256) == 64 {
						found[id] = true
					}
				}
			}
			if len(found) != len(tt.want) {
				t.Fatalf("for %q got verified providers %v, expected %v", tt.query, found, tt.want)
			}
		})
	}
}

func TestReviewedSearchPreservesExplicitAndScopedPublicSources(t *testing.T) {
	var baseURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/guru":
			_, _ = w.Write([]byte(`{}`))
		case "/lap":
			w.Header().Set("Content-Type", "text/lap")
			_, _ = w.Write([]byte("@search_query slack\n@total 0\n"))
		case "/list":
			fmt.Fprintf(w, `{"entries":[{"API":"Slack","Description":"Fixture service","Link":%q,"Category":"Communication"}]}`, baseURL+"/docs")
		case "/openapi.json":
			_, _ = w.Write([]byte(`{"openapi":"3.0.0","info":{"title":"Slack Fixture","version":"1"},"paths":{}}`))
		case "/.well-known/api-catalog":
			w.Header().Set("Content-Type", "application/linkset+json")
			fmt.Fprintf(w, `{"linkset":[{"anchor":%q,"service-desc":[{"href":%q,"type":"application/json"}]}]}`, baseURL+"/api/slack", baseURL+"/openapi.json")
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	baseURL = server.URL
	for _, tt := range []struct {
		args   []string
		source apitools.Source
	}{
		{[]string{"--source=auto"}, apitools.SourcePublicAPIs},
		{[]string{"--source", "public-apis"}, apitools.SourcePublicAPIs},
		{[]string{"--provider-url", baseURL + "/docs"}, apitools.SourceRFC9727},
	} {
		var out, errOut bytes.Buffer
		args := append([]string{"--query", "slack", "--json", "--config", "unused-reviewed-config.json"}, tt.args...)
		code := runSearchWithClient(args, &out, &errOut, func(string) (*apitools.Client, func(), error) {
			return &apitools.Client{APIsGuruListURL: baseURL + "/guru", LAPSearchURL: baseURL + "/lap",
				PublicAPIsURL: baseURL + "/list", AllowUnsafeHosts: true}, func() {}, nil
		})
		if code != exitSuccess {
			t.Fatalf("%v exit %d: %s", tt.args, code, errOut.String())
		}
		var report apitools.SearchReport
		if err := json.Unmarshal(out.Bytes(), &report); err != nil {
			t.Fatal(err)
		}
		if len(report.Results) != 1 || report.Results[0].Source != string(tt.source) || report.Results[0].SpecURL != baseURL+"/openapi.json" {
			t.Fatalf("%v changed public results: %+v", tt.args, report)
		}
	}
}
