package apitools_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/apitools/catalog"
)

func preparedCatalogDiscovery(t *testing.T) apitools.CatalogDiscoveryOptions {
	t.Helper()
	index := syntheticCatalogIndex(t)
	value, err := apitools.BuildCatalogOperationIndex(context.Background(), index)
	if err != nil {
		t.Fatal(err)
	}
	if err := apitools.WriteCatalogOperationIndex(context.Background(), index, value); err != nil {
		t.Fatal(err)
	}
	return apitools.CatalogDiscoveryOptions{Index: index, Request: apitools.CatalogDiscoveryRequest{SchemaVersion: apitools.CatalogDiscoverySchemaVersion, Contract: apitools.StepContract{Purpose: "list notes records", Effect: apitools.OperationEffectRead}}}
}

func TestCatalogDiscoveryRetrievalSharedMetadataAndRelocation(t *testing.T) {
	first, second := preparedCatalogDiscovery(t), preparedCatalogDiscovery(t)
	ctx := context.Background()
	a, err := apitools.DiscoverCatalogOperations(ctx, first)
	if err != nil {
		t.Fatal(err)
	}
	if len(a.Candidates) != 2 || len(a.Candidates[0].References) != 2 || a.ExaminedOperations != 2 || !a.Scope.Complete {
		t.Fatalf("shared retrieval: %#v", a)
	}
	for _, candidate := range a.Candidates {
		for _, source := range candidate.Sources {
			if source.LicenseNote == "" || source.LicenseIdentifier != "unknown" || source.Redistribution != "unknown" || source.Authority != "official-provider" {
				t.Fatal("source evidence lost or permission inferred")
			}
		}
	}
	if err := os.Remove(filepath.Join(first.Index.Root.Directory, "openapi/notes.json")); err != nil {
		t.Fatal(err)
	}
	after, err := apitools.DiscoverCatalogOperations(ctx, first)
	if err != nil || len(after.Candidates) != 2 {
		t.Fatal("discovery parsed source bytes")
	}
	b, err := apitools.DiscoverCatalogOperations(ctx, second)
	if err != nil {
		t.Fatal(err)
	}
	encodedA, _ := json.Marshal(a)
	encodedB, _ := json.Marshal(b)
	if !bytes.Equal(encodedA, encodedB) {
		t.Fatal("root relocation changed report")
	}
	if bytes.Contains(encodedA, []byte(first.Index.Root.Directory)) {
		t.Fatal("installation path entered evidence wire")
	}
	first.Request.ProviderKeys = []string{"Example Notes Shared"}
	selected, err := apitools.DiscoverCatalogOperations(ctx, first)
	if err != nil || len(selected.Scope.ProviderIDs) != 1 || selected.Scope.ProviderIDs[0] != "example-notes-alias" || len(selected.Candidates[0].References) != 1 {
		t.Fatal("provider constraint broadened")
	}
}

func TestCatalogDiscoveryMissingStaleInvalidAndEmptyScope(t *testing.T) {
	for _, name := range []string{"no-root", "missing-index", "stale", "corrupt", "unknown-provider", "empty-scope", "unsafe-root", "cancel"} {
		t.Run(name, func(t *testing.T) {
			options := preparedCatalogDiscovery(t)
			ctx := context.Background()
			want := apitools.CatalogDiscoveryInsufficientEvidence
			switch name {
			case "no-root":
				options.Index.Root = catalog.RootOptions{}
			case "missing-index":
				if err := os.Remove(filepath.Join(options.Index.Root.Directory, catalog.DefaultOperationIndexPath)); err != nil {
					t.Fatal(err)
				}
			case "stale":
				reader := options.Index.ReadRegistrations
				options.Index.ReadRegistrations = func(ctx context.Context, root catalog.RootOptions) ([]catalog.CatalogSpecArtifact, error) {
					rows, err := reader(ctx, root)
					rows[0].SHA256 = strings.Repeat("0", 64)
					return rows, err
				}
			case "corrupt":
				if err := os.WriteFile(filepath.Join(options.Index.Root.Directory, catalog.DefaultOperationIndexPath), []byte("{}"), 0o600); err != nil {
					t.Fatal(err)
				}
				want = apitools.CatalogDiscoveryBlocked
			case "unknown-provider":
				options.Request.ProviderKeys = []string{"unknown provider"}
				want = apitools.CatalogDiscoveryBlocked
			case "empty-scope":
				options.Request.ProviderKeys = []string{}
			case "unsafe-root":
				options.Index.Root.IndexPath = "../outside"
				want = apitools.CatalogDiscoveryBlocked
			case "cancel":
				cancelled, cancel := context.WithCancel(ctx)
				cancel()
				ctx = cancelled
			}
			report, err := apitools.DiscoverCatalogOperations(ctx, options)
			if name == "cancel" {
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("cancellation: %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			if report.Outcome != want || len(report.Candidates) != 0 || report.Scope.Complete {
				t.Fatalf("false outcome for %s: %#v", name, report)
			}
			if name == "no-root" && len(report.Leads) != 2 {
				t.Fatal("metadata-only leads missing")
			}
			if name == "stale" {
				for _, coverage := range report.Coverage {
					if coverage.State != "stale" {
						t.Fatal("stale scope became examined")
					}
				}
			}
			if name == "empty-scope" && len(report.Scope.ProviderIDs) != 0 {
				t.Fatal("empty scope broadened")
			}
		})
	}
}

func TestCatalogDiscoveryFiltersAndWorkReportBounds(t *testing.T) {
	options := preparedCatalogDiscovery(t)
	for _, filters := range []apitools.CatalogDiscoveryFilters{{RequireLicenseEvidence: true}, {RequireRedistributionEvidence: true}, {Authorities: []string{"official-docs"}}} {
		options.Request.Filters = filters
		report, err := apitools.DiscoverCatalogOperations(context.Background(), options)
		if err != nil || len(report.Candidates) != 0 || len(report.Exclusions) != 2 || report.Outcome == apitools.CatalogDiscoveryNoQualifyingAPI {
			t.Fatalf("filter exclusions: %v %#v", err, report)
		}
	}
	options.Request.Filters = apitools.CatalogDiscoveryFilters{}
	options.Request.Limits.MaxOperations = 1
	report, err := apitools.DiscoverCatalogOperations(context.Background(), options)
	if err != nil || report.ExaminedOperations != 1 || !report.Incomplete || !report.Truncated || report.Scope.Complete {
		t.Fatalf("work bound hidden: %v %#v", err, report)
	}
	for _, coverage := range report.Coverage {
		if coverage.State != "work_limit" {
			t.Fatal("limited scope stayed completely examined")
		}
	}
	options.Request.Limits.MaxContextBytes = apitools.MinCatalogDiscoveryContextBytes
	report, err = apitools.DiscoverCatalogOperations(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	encoded, _ := json.Marshal(report)
	if len(encoded) > options.Request.Limits.MaxContextBytes {
		t.Fatalf("report budget: %d", len(encoded))
	}
	// Oversized raw notes remain an explicit gap instead of inferred permission
	// or truncation that presents altered text as verbatim evidence.
	options = preparedCatalogDiscovery(t)
	options.Index.Catalog.Providers[0].SpecReferences[0].LicenseNote = strings.Repeat("x", 70<<10)
	index, err := apitools.BuildCatalogOperationIndex(context.Background(), options.Index)
	if err != nil {
		t.Fatal(err)
	}
	if err := apitools.WriteCatalogOperationIndex(context.Background(), options.Index, index); err != nil {
		t.Fatal(err)
	}
	report, err = apitools.DiscoverCatalogOperations(context.Background(), options)
	if err != nil || !report.Incomplete {
		t.Fatalf("oversized note gap hidden: %v", err)
	}
	for _, candidate := range report.Candidates {
		for _, source := range candidate.Sources {
			if source.ProviderID == "example-notes" && (source.LicenseNote != "" || len(source.EvidenceGaps) == 0) {
				t.Fatal("oversized raw evidence altered or hidden")
			}
		}
	}
}

func TestCatalogDiscoveryIndexPermutationAndBoundedTraversal(t *testing.T) {
	for _, limit := range []int{0, 1} {
		options := preparedCatalogDiscovery(t)
		options.Request.Limits.MaxOperations = limit
		before, err := apitools.DiscoverCatalogOperations(context.Background(), options)
		if err != nil {
			t.Fatal(err)
		}
		index, err := apitools.ReadCatalogOperationIndex(context.Background(), options.Index)
		if err != nil {
			t.Fatal(err)
		}
		artifact := &index.Artifacts[0]
		artifact.Links[0], artifact.Links[1] = artifact.Links[1], artifact.Links[0]
		artifact.Operations[0], artifact.Operations[1] = artifact.Operations[1], artifact.Operations[0]
		index.Coverage[0], index.Coverage[1] = index.Coverage[1], index.Coverage[0]
		if err := apitools.WriteCatalogOperationIndex(context.Background(), options.Index, index); err != nil {
			t.Fatal(err)
		}
		after, err := apitools.DiscoverCatalogOperations(context.Background(), options)
		if err != nil {
			t.Fatal(err)
		}
		a, _ := json.Marshal(before)
		b, _ := json.Marshal(after)
		if !bytes.Equal(a, b) {
			t.Fatalf("valid index permutation changed discovery under operation limit %d", limit)
		}
	}
}

func TestCatalogDiscoveryProviderWithoutReferencesRemainsUnexamined(t *testing.T) {
	options := preparedCatalogDiscovery(t)
	provider := options.Index.Catalog.Providers[0]
	provider.ID, provider.DisplayName, provider.CandidateID = "unreferenced", "Unreferenced Service", "unreferenced"
	provider.SpecReferences = nil
	provider.OfficialOpenAPIAvailability = catalog.SpecAvailabilityUnavailable
	options.Index.Catalog.Providers = append(options.Index.Catalog.Providers, provider)
	index, err := apitools.BuildCatalogOperationIndex(context.Background(), options.Index)
	if err != nil {
		t.Fatal(err)
	}
	if err := apitools.WriteCatalogOperationIndex(context.Background(), options.Index, index); err != nil {
		t.Fatal(err)
	}
	options.Request.Contract.Effect = apitools.OperationEffectWrite
	report, err := apitools.DiscoverCatalogOperations(context.Background(), options)
	if err != nil || report.Outcome != apitools.CatalogDiscoveryInsufficientEvidence || !report.Incomplete || report.Scope.Complete {
		t.Fatal("curated no-reference provider established complete API absence")
	}
	found := false
	for _, item := range report.Coverage {
		found = found || item.ProviderID == provider.ID && item.State == "unsupported"
	}
	if !found {
		t.Fatal("provider-level unexamined coverage missing")
	}
	options.Index.Root = catalog.RootOptions{}
	options.Request.ProviderKeys = []string{provider.DisplayName}
	report, err = apitools.DiscoverCatalogOperations(context.Background(), options)
	if err != nil || len(report.Leads) != 1 || report.Leads[0].ProviderID != provider.ID || report.Outcome != apitools.CatalogDiscoveryInsufficientEvidence {
		t.Fatal("no-reference metadata lead missing")
	}
}

func TestCatalogDiscoverySharedLinkLimitBeforeOperationAllocation(t *testing.T) {
	options := preparedCatalogDiscovery(t)
	reader := options.Index.ReadRegistrations
	options.Index.ReadRegistrations = func(ctx context.Context, root catalog.RootOptions) ([]catalog.CatalogSpecArtifact, error) {
		rows, err := reader(ctx, root)
		if err != nil {
			return nil, err
		}
		for _, provider := range options.Index.Catalog.Providers[2:] {
			row := rows[0]
			row.ProviderID = provider.ID
			rows = append(rows, row)
		}
		return rows, nil
	}
	for i := 0; i < 31; i++ {
		provider := options.Index.Catalog.Providers[0]
		provider.ID = fmt.Sprintf("shared-%02d", i)
		provider.CandidateID = provider.ID
		provider.DisplayName = fmt.Sprintf("Shared Notes %02d", i)
		options.Index.Catalog.Providers = append(options.Index.Catalog.Providers, provider)
	}
	index, err := apitools.BuildCatalogOperationIndex(context.Background(), options.Index)
	if err != nil {
		t.Fatal(err)
	}
	if err := apitools.WriteCatalogOperationIndex(context.Background(), options.Index, index); err != nil {
		t.Fatal(err)
	}
	report, err := apitools.DiscoverCatalogOperations(context.Background(), options)
	if err != nil || report.Outcome != apitools.CatalogDiscoveryInsufficientEvidence || len(report.Candidates) != 0 || report.ExaminedOperations != 0 || !report.Incomplete || !report.Truncated {
		t.Fatal("shared references were allocated/evaluated before the query link limit")
	}
	if len(report.Leads) != 33 || len(report.Coverage) != 33 {
		t.Fatal("shared provider provenance was discarded")
	}
	for _, coverage := range report.Coverage {
		if coverage.State != "link_limit" {
			t.Fatal("link-limited scope was not labeled")
		}
	}
	options.Request.ProviderKeys = []string{"Example Notes"}
	report, err = apitools.DiscoverCatalogOperations(context.Background(), options)
	if err != nil || report.Outcome != apitools.CatalogDiscoveryMatch || len(report.Candidates[0].References) != 1 {
		t.Fatal("narrow provider scope did not recover supported retrieval")
	}
}
