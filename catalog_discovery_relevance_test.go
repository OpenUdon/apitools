package apitools_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/apitools/catalog"
)

func TestCatalogDiscoveryEvalRelevanceBaseline(t *testing.T) {
	root := t.TempDir()
	cat := catalog.Catalog{}
	var registrations []catalog.CatalogSpecArtifact
	base := preparedCatalogDiscovery(t).Index.Catalog.Providers[0]
	for _, id := range []string{"jira", "slack", "drive"} {
		data, err := os.ReadFile(filepath.Join("testdata/catalog-discovery/relevance", id+".yaml"))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, id+".yaml"), data, 0600); err != nil {
			t.Fatal(err)
		}
		digest := sha256.Sum256(data)
		provider := base
		provider.ID, provider.DisplayName, provider.CandidateID = id, id, id
		provider.SpecReferences = []catalog.SpecReference{{ID: id, Kind: catalog.SpecKindOpenAPI, URL: "https://" + id + ".example.test/spec.yaml", SourceAuthority: catalog.SourceAuthorityOfficialProvider, SourceNote: "Repository-owned OpenUdon eval fixture; synthetic, not provider evidence."}}
		cat.Providers = append(cat.Providers, provider)
		registrations = append(registrations, catalog.CatalogSpecArtifact{ProviderID: id, SpecRefID: id, ArtifactID: id, Kind: "openapi", Path: id + ".yaml", SHA256: hex.EncodeToString(digest[:]), Bytes: int64(len(data))})
	}
	indexOptions := apitools.CatalogIndexOptions{Root: catalog.RootOptions{Directory: root}, Catalog: &cat, ReadRegistrations: func(context.Context, catalog.RootOptions) ([]catalog.CatalogSpecArtifact, error) {
		return registrations, nil
	}}
	index, err := apitools.BuildCatalogOperationIndex(context.Background(), indexOptions)
	if err != nil {
		t.Fatal(err)
	}
	if err := apitools.WriteCatalogOperationIndex(context.Background(), indexOptions, index); err != nil {
		t.Fatal(err)
	}
	queries := []struct{ purpose, inputs, operation string }{
		{"Create an incident issue", "summary,description,severity", "createIssue"},
		{"Post an incident alert to Slack", "channel,text", "postMessage"},
		{"Upload an incident report file", "name,parentId,content", "uploadFile"},
	}
	hits := 0
	for queryIndex, query := range queries {
		required := true
		inputs := map[string]apitools.ContractValue{}
		for _, name := range strings.Split(query.inputs, ",") {
			inputs[name] = apitools.ContractValue{Type: "string", Required: &required}
		}
		report, err := apitools.DiscoverCatalogOperations(context.Background(), apitools.CatalogDiscoveryOptions{Index: indexOptions, Request: apitools.CatalogDiscoveryRequest{SchemaVersion: apitools.CatalogDiscoverySchemaVersion, Contract: apitools.StepContract{Purpose: query.purpose, Inputs: inputs, Effect: apitools.OperationEffectWrite}}})
		if err != nil {
			t.Fatal(err)
		}
		if len(report.Candidates) > 0 && report.Candidates[0].Candidate.Operation.OperationID == query.operation {
			hits++
		}
		want := []apitools.CatalogDiscoveryOutcome{apitools.CatalogDiscoveryInsufficientEvidence, apitools.CatalogDiscoveryMatch, apitools.CatalogDiscoveryInsufficientEvidence}[queryIndex]
		if report.Outcome != want {
			t.Fatalf("eval query %q: outcome %s, expected %s", query.purpose, report.Outcome, want)
		}
	}
	t.Logf("repository-owned relevance baseline precision@1 = %d/%d (%.3f)", hits, len(queries), float64(hits)/float64(len(queries)))
}

// CatalogPlan selects explicit provider/artifact pairs; discovery preserves
// those provider identities without adopting its model planning or order.
func TestCatalogDiscoveryNamedCatalogPlanParity(t *testing.T) {
	for _, keys := range [][]string{{"Gmail", "OpenWeatherMap"}, {"Google Drive", "Slack", "Jira"}, {"Microsoft Outlook"}} {
		report, err := apitools.DiscoverCatalogOperations(context.Background(), apitools.CatalogDiscoveryOptions{Request: apitools.CatalogDiscoveryRequest{SchemaVersion: apitools.CatalogDiscoverySchemaVersion, ProviderKeys: keys, Contract: apitools.StepContract{Purpose: "report incident records"}}})
		if err != nil || report.Outcome != apitools.CatalogDiscoveryInsufficientEvidence || !report.Scope.ProviderConstrained || len(report.Scope.ProviderIDs) != len(keys) {
			t.Fatalf("named provider parity: %v %+v", err, report.Scope)
		}
		cat := catalog.BuiltInCatalog()
		for _, key := range keys {
			provider, ok := cat.FindProvider(key)
			if !ok {
				t.Fatalf("unknown parity provider %q", key)
			}
			found := false
			for _, id := range report.Scope.ProviderIDs {
				found = found || id == provider.ID
			}
			if !found {
				t.Fatal("named selection was broadened or split")
			}
		}
	}
}
