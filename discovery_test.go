package apitools

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiscoveryExtractURLsImportsAndSelects(t *testing.T) {
	urls := ExtractURLs("Use https://api.example.test/openapi.yaml, then https://api.example.test/openapi.yaml.")
	if len(urls) != 2 || urls[0] != "https://api.example.test/openapi.yaml" {
		t.Fatalf("urls = %#v", urls)
	}

	base := t.TempDir()
	openAPIDir := filepath.Join(base, "openapi")
	if err := os.MkdirAll(openAPIDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(openAPIDir, "weather.yaml"), []byte(`openapi: 3.0.0
info:
  title: Weather API
  version: 1.0.0
  description: Forecasts and observations.
paths: {}
`), 0o644); err != nil {
		t.Fatal(err)
	}
	candidates, err := DiscoverOpenAPI(context.Background(), openAPIDir, base, "weather forecast")
	if err != nil {
		t.Fatal(err)
	}
	if len(candidates) != 1 || candidates[0].RelativePath != "openapi/weather.yaml" || candidates[0].Score == 0 {
		t.Fatalf("candidates = %#v", candidates)
	}
	primary, err := SelectPrimaryDiscoveryCandidate([]DiscoveryCandidate{
		{RelativePath: "openapi/low.yaml", Score: 1},
		{RelativePath: "openapi/high.yaml", Score: 10},
	})
	if err != nil {
		t.Fatal(err)
	}
	if primary.RelativePath != "openapi/high.yaml" {
		t.Fatalf("primary = %#v", primary)
	}
}

func TestDiscoveryAPIsGuruRejectsPrivateListURLBeforeRequest(t *testing.T) {
	var called bool
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		called = true
		return nil, nil
	})}
	_, err := (&Discoverer{
		HTTPClient:      client,
		APIsGuruListURL: "http://127.0.0.1/list.json",
	}).ImportBestAPIsGuruMatch(context.Background(), t.TempDir(), t.TempDir(), "weather")
	if err == nil || !strings.Contains(err.Error(), "private") {
		t.Fatalf("err = %v", err)
	}
	if called {
		t.Fatalf("HTTP client was called before private host rejection")
	}
}

func TestImportProjectURLsTruncatesWithoutFailingDiscovery(t *testing.T) {
	var text strings.Builder
	for i := 0; i < DefaultMaxProjectURLs+1; i++ {
		text.WriteString(" http://127.0.0.1/spec-")
		text.WriteRune(rune('a' + i))
		text.WriteString(".yaml")
	}
	openAPIDir := t.TempDir()
	report, err := (&Discoverer{}).ImportProjectURLsReport(context.Background(), openAPIDir, t.TempDir(), text.String())
	if err != nil {
		t.Fatalf("ImportProjectURLsReport() error = %v", err)
	}
	if !report.Truncated || len(report.Candidates) != 0 || len(report.Attempts) != DefaultMaxProjectURLs || len(report.Diagnostics) != DefaultMaxProjectURLs+1 || report.Diagnostics[0].Code != "discovery.limit.project_urls" || report.Diagnostics[0].Severity != "warning" {
		t.Fatalf("report = %#v", report)
	}
	for _, attempt := range report.Attempts {
		if attempt.Status != "fail" || !strings.Contains(attempt.Source, "spec-") {
			t.Fatalf("unexpected bounded attempt = %#v", attempt)
		}
	}

	exampleDir := t.TempDir()
	openAPIDir = filepath.Join(exampleDir, "openapi")
	if err := os.MkdirAll(openAPIDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeLocalFile(t, filepath.Join(openAPIDir, "local.yaml"), validLocalSpec("Local API", "Local"))
	candidates, discoveryReport, err := (&Discoverer{}).DiscoverWithReport(context.Background(), exampleDir, text.String())
	if err != nil {
		t.Fatalf("DiscoverWithReport() error = %v", err)
	}
	if !discoveryReport.Truncated || len(discoveryReport.Attempts) != DefaultMaxProjectURLs+1 || len(discoveryReport.Diagnostics) != DefaultMaxProjectURLs+1 || discoveryReport.Diagnostics[0].Code != "discovery.limit.project_urls" || len(candidates) != 1 || candidates[0].Title != "Local API" {
		t.Fatalf("partial discovery = candidates %#v, report %#v", candidates, discoveryReport)
	}
}

func TestProjectURLImportMethodsRemainSourceCompatible(t *testing.T) {
	var discoverer Discoverer
	var importURLs func(context.Context, string, string, string) ([]DiscoveryCandidate, error) = discoverer.ImportProjectURLs
	var importURLsWithReport func(context.Context, string, string, string) ([]DiscoveryCandidate, []DiscoveryAttempt) = discoverer.ImportProjectURLsWithReport

	candidates, err := importURLs(context.Background(), t.TempDir(), t.TempDir(), "no URLs")
	if err != nil || len(candidates) != 0 {
		t.Fatalf("ImportProjectURLs() = %#v, %v", candidates, err)
	}
	candidates, attempts := importURLsWithReport(context.Background(), t.TempDir(), t.TempDir(), "no URLs")
	if len(candidates) != 0 || len(attempts) != 0 {
		t.Fatalf("ImportProjectURLsWithReport() = %#v, %#v", candidates, attempts)
	}
}

func TestDiscoveryCandidatesDeduplicateByDigest(t *testing.T) {
	candidates, digests := dedupeDiscoveryCandidatesByDigest([]DiscoveryCandidate{
		{RelativePath: "openapi/a.yaml", Source: "local"},
		{RelativePath: "openapi/b.yaml", Source: "url"},
		{RelativePath: "openapi/c.yaml", Source: "url"},
		{RelativePath: "openapi/no-digest.yaml", Source: "local"},
	}, []string{"abc", "ABC", "def", ""})
	if len(candidates) != 3 || len(digests) != 3 || candidates[0].RelativePath != "openapi/a.yaml" || digests[0] != "abc" || candidates[1].RelativePath != "openapi/c.yaml" || digests[1] != "def" || candidates[2].RelativePath != "openapi/no-digest.yaml" {
		t.Fatalf("deduplicated candidates = %#v with digests %#v", candidates, digests)
	}
}

func TestDiscoveryCandidateKeepsLegacyFieldShape(t *testing.T) {
	// Unkeyed construction protects source compatibility for downstream users.
	_ = DiscoveryCandidate{"path", "relative", "title", "description", "source", 1}
	_ = LocalOptions{"dir", "base", "query", 1024}
	_ = LocalResult{"path", "relative", "title", "description", 1, SpecMetadata{}}
	_ = ImportedSpec{"name", "path", "title", "description", "url", "sha256", 1, SpecMetadata{}}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
