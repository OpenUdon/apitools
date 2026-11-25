package apitools

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/catalogupdate"
)

func TestCatalogMaintenanceRepeatReusesArtifactAndEmbeddingInputs(t *testing.T) {
	dir, err := os.MkdirTemp(".", ".catalog-repeat-fixture-")
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
	ctx, cancel := context.WithTimeout(t.Context(), 45*time.Second)
	defer cancel()
	config := CatalogUpdateConfig{CacheDir: filepath.Join(dir, "cache"), StateDB: filepath.Join(dir, "state.sqlite")}
	const spec = `{"openapi":"3.0.3","info":{"title":"Records","version":"1.0"},"paths":{"/records":{"get":{"responses":{"204":{"description":"empty"}}}}}}`
	provider, _ := catalog.FindBuiltInProvider("attio")
	source := catalogUpdateReferenceSources(provider)[0]
	seed := source
	seed.Revision = fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(spec)))
	seedCatalogUpdateFixture(t, config, seed, spec)
	var artifactGets int
	publisher := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		response := &http.Response{StatusCode: http.StatusNotFound, Body: io.NopCloser(strings.NewReader("")), Request: r}
		if r.URL.String() == source.URL {
			response.StatusCode, response.ContentLength = http.StatusOK, int64(len(spec))
			response.Header = http.Header{"Content-Length": {strconv.Itoa(len(spec))}}
			if r.Method == http.MethodGet {
				artifactGets++
				response.Body = io.NopCloser(strings.NewReader(spec))
			}
		}
		return response, nil
	})}}
	first, err := publisher.runMonthlyCatalogUpdate(ctx, config, "2026-11", "", nil)
	if err != nil || !first.Completed || len(first.Providers) != 334 {
		t.Fatalf("first scheduled scope: %d providers, completed %t, %v", len(first.Providers), first.Completed, err)
	}
	firstExport, err := catalogupdate.ReadVersionExport(ctx, config.StateDB)
	if err != nil || !bytes.Contains(firstExport, []byte("attio")) {
		t.Fatalf("first export omitted the accepted reference: %s, %v", firstExport, err)
	}
	const query = "reverse a card charge"
	server, _ := reviewedEmbeddingFixture(t, query, http.StatusOK)
	embedder := &Client{AllowUnsafeHosts: true, HTTPClient: server.Client()}
	index := filepath.Join(dir, "search.sqlite")
	options := ReviewedProviderSearchOptions{Embeddings: &ReviewedEmbeddingConfig{
		Endpoint: server.URL + "/v1/embeddings", Model: "fixture-model", MaxRequests: 16, MaxInputs: 350, MaxBytes: 1 << 20,
	}}
	profiles := catalog.BuiltInProviderSearchProfiles()
	embed := func() ReviewedEmbeddingUsage {
		t.Helper()
		db, err := openReviewedProviderIndex(index)
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()
		_, usage := embedder.rankReviewedEmbeddings(ctx, db, options, profiles, nil, query)
		return *usage
	}
	if usage := embed(); usage.Status != "used" || usage.Inputs != 335 {
		t.Fatalf("first search did not populate the 334 provider vectors: %+v", usage)
	}
	second, err := publisher.runMonthlyCatalogUpdate(ctx, config, "2026-12", first.RunID, nil)
	if err != nil || !second.Completed || second.RunID == first.RunID || !slices.Equal(first.Providers, second.Providers) {
		t.Fatalf("repeat lost its distinct full scope: %d providers, completed %t, %v", len(second.Providers), second.Completed, err)
	}
	for _, report := range []CatalogUpdateReport{first, second} {
		var count, summaries int
		for _, outcome := range report.Outcomes {
			if outcome.ReferenceID == "" {
				summaries++
			}
			if outcome.ProviderID == "attio" && outcome.ReferenceID == source.ReferenceID && outcome.Status == "unchanged" && outcome.Object != nil {
				count++
			}
		}
		if count != 1 || summaries != 334 {
			t.Fatalf("run %s has %d unchanged source records and %d summaries", report.RunID, count, summaries)
		}
	}
	cache, err := catalogupdate.MeasureArtifactCache(ctx, config.CacheDir, "")
	if err != nil || cache.Objects != 1 || cache.Bytes != int64(len(spec)) || artifactGets != 2 {
		t.Fatalf("repeat did not deduplicate identical body bytes: %+v, GETs %d, %v", cache, artifactGets, err)
	}
	secondExport, err := catalogupdate.ReadVersionExport(ctx, config.StateDB)
	if err != nil || !bytes.Equal(firstExport, secondExport) {
		t.Fatalf("repeat changed the accepted Git export: %v", err)
	}
	profiles[0].Description += " Updated public fixture description."
	if usage := embed(); usage.Status != "used" || usage.Requests != 2 || usage.Inputs != 2 {
		t.Fatalf("repeat should send the changed public profile and query only: %+v", usage)
	}
}

func TestCatalogMaintenanceReconcilesRetainedFirstInventory(t *testing.T) {
	path := os.Getenv("APITOOLS_C09_FIRST_REPORT")
	if path == "" {
		t.Skip("set APITOOLS_C09_FIRST_REPORT to the frozen C07 completed report")
	}
	data, err := readLocalSpecFile(path, DefaultMaxBytes)
	if err != nil || fmt.Sprintf("%x", sha256.Sum256(data)) != "c712ed2d84c5fdd7e201e090ddd5cc691e4f97e866d8555d6293d55757c2e687" {
		t.Fatalf("first inventory bytes differ from the accepted C07 receipt: %v", err)
	}
	var report CatalogUpdateReport
	if err := json.Unmarshal(data, &report); err != nil || !report.Completed || report.RunID != "2OXWJLDNCITU5LOUTAZAVOZ25V" ||
		len(report.Providers) != 334 || len(report.Outcomes) != 529 {
		t.Fatalf("unexpected C07 report identity or scope: %s, completed %t, %d providers, %d outcomes, %v",
			report.RunID, report.Completed, len(report.Providers), len(report.Outcomes), err)
	}
	dry, err := (&Client{}).UpdateCatalog(t.Context(), CatalogUpdateOptions{All: true, DryRun: true})
	if err != nil || !slices.Equal(dry.Providers, report.Providers) {
		t.Fatalf("reviewed provider scope changed since C07: %v", err)
	}
	seen := make(map[string]bool)
	for _, outcome := range report.Outcomes {
		if outcome.Status == catalogupdate.RetryPendingStatus {
			t.Fatalf("closed report contains a pending reference: %s/%s", outcome.ProviderID, outcome.ReferenceID)
		}
		if outcome.ReferenceID == "" {
			if seen[outcome.ProviderID] {
				t.Fatalf("duplicate provider checkpoint %s", outcome.ProviderID)
			}
			seen[outcome.ProviderID] = true
		}
	}
	for _, provider := range dry.Providers {
		if !seen[provider] {
			t.Errorf("first inventory misses reviewed provider %s", provider)
		}
	}
}
