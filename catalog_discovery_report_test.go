package apitools

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCatalogDiscoveryProjectionRemapsVisibleTies(t *testing.T) {
	report := CatalogDiscoveryReport{SchemaVersion: CatalogDiscoverySchemaVersion, Outcome: CatalogDiscoveryAmbiguous, Limits: CatalogDiscoveryLimits{MaxResults: 2, MaxContextBytes: MinCatalogDiscoveryContextBytes}, Scope: CatalogDiscoveryScope{Complete: true}, QualifiedOperations: 3}
	for i := 0; i < 3; i++ {
		report.Candidates = append(report.Candidates, CatalogDiscoveryCandidate{Qualified: true, Candidate: OperationCandidate{Match: StepContractMatch{Score: 20}}})
	}
	report.Candidates[0].Sources = []CatalogDiscoverySourceEvidence{{LicenseNote: strings.Repeat("x", 60000)}}
	report.Ties = []CatalogDiscoveryTie{{Score: 20, CandidateIndexes: []int{0, 1, 2}}}
	// MaxResults limits considered rank positions; the oversized first result
	// is omitted, so the single visible result cannot carry a fabricated tie.
	limited := boundCatalogDiscoveryReport(report)
	if len(limited.Candidates) != 1 || len(limited.Ties) != 0 || limited.Outcome != CatalogDiscoveryAmbiguous || limited.QualifiedOperations != 3 {
		t.Fatal("display limit changed qualification or fabricated tied indexes")
	}
	report.Limits.MaxResults = 3
	projected := boundCatalogDiscoveryReport(report)
	if len(projected.Candidates) != 2 || len(projected.Ties) != 1 || len(projected.Ties[0].CandidateIndexes) != 2 || projected.Ties[0].CandidateIndexes[0] != 0 || projected.Ties[0].CandidateIndexes[1] != 1 {
		t.Fatal("tie indexes did not follow the displayed projection")
	}
	encoded, err := json.Marshal(projected)
	if err != nil || len(encoded) > MinCatalogDiscoveryContextBytes || !projected.Incomplete || projected.Scope.Complete {
		t.Fatal("projection hid evidence loss or exceeded its bound")
	}
}
