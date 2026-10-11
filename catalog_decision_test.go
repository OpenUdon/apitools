package apitools

import (
	"errors"
	"strings"
	"testing"
)

func sampleContract() StepContract {
	return StepContract{
		Purpose: "create dns record",
		Effect:  OperationEffectWrite,
	}
}

func sampleCandidate(name, method, path string, qualified bool, hasLocalRef bool) CatalogDiscoveryCandidate {
	c := CatalogDiscoveryCandidate{
		Candidate: OperationCandidate{
			Operation: OperationSummary{
				OperationID: name,
				Method:      method,
				Path:        path,
			},
		},
		Qualified: qualified,
		Sources: []CatalogDiscoverySourceEvidence{
			{
				ProviderID:        "cloudflare",
				SpecRefID:         "dns",
				Authority:         "official-provider",
				LicenseIdentifier: "MIT",
				Redistribution:    "permitted",
			},
		},
		QualificationGaps: []string{},
	}
	if hasLocalRef {
		c.References = []CatalogArtifactReference{
			{
				CatalogSHA256: "abc",
				ProviderID:    "cloudflare",
				SpecRefID:     "dns",
				ArtifactID:    "dns.json",
				Kind:          "openapi",
				SHA256:        "123",
				Bytes:         1024,
			},
		}
	}
	return c
}

func TestProjectCatalogDecisionMatch(t *testing.T) {
	req := CatalogDiscoveryRequest{
		SchemaVersion: CatalogDiscoverySchemaVersion,
		Contract:      sampleContract(),
	}
	rep := CatalogDiscoveryReport{
		SchemaVersion:       CatalogDiscoverySchemaVersion,
		Outcome:             CatalogDiscoveryMatch,
		QualifiedOperations: 1,
		Scope: CatalogDiscoveryScope{
			ProviderIDs: []string{"cloudflare"},
			Complete:    true,
		},
		Candidates: []CatalogDiscoveryCandidate{
			sampleCandidate("createRecord", "POST", "/zones/{zone_id}/dns_records", true, true),
			sampleCandidate("listRecords", "GET", "/zones/{zone_id}/dns_records", false, true),
		},
	}

	view, prompt, err := ProjectCatalogDecision(rep, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(view.Candidates) != 2 {
		t.Fatalf("expected 2 candidates, got %d", len(view.Candidates))
	}
	if len(view.Selectable) != 1 || view.Selectable[0] != 1 {
		t.Fatalf("expected selectable [1], got %v", view.Selectable)
	}
	if view.Omitted != 0 {
		t.Fatalf("expected 0 omitted, got %d", view.Omitted)
	}

	// Verify prompt contains detailed evidence (F09)
	if !strings.Contains(prompt, "Candidate 1; qualified=true") {
		t.Fatalf("expected prompt to contain candidate 1 evidence: %s", prompt)
	}
	if !strings.Contains(prompt, "native source evidence=") || !strings.Contains(prompt, "references=") {
		t.Fatalf("expected prompt to contain detailed source and reference evidence: %s", prompt)
	}

	// Verify allowed answers
	if !view.Allows("select 1") {
		t.Fatal("expected Allows('select 1') == true")
	}
	if view.Allows("select 2") {
		t.Fatal("expected Allows('select 2') == false")
	}
	if !view.Allows("reject") {
		t.Fatal("expected Allows('reject') == true")
	}
}

func TestProjectCatalogDecisionAmbiguous(t *testing.T) {
	req := CatalogDiscoveryRequest{
		SchemaVersion: CatalogDiscoverySchemaVersion,
		Contract:      sampleContract(),
	}
	rep := CatalogDiscoveryReport{
		SchemaVersion:       CatalogDiscoverySchemaVersion,
		Outcome:             CatalogDiscoveryAmbiguous,
		QualifiedOperations: 2,
		Scope: CatalogDiscoveryScope{
			ProviderIDs: []string{"cloudflare"},
			Complete:    true,
		},
		Candidates: []CatalogDiscoveryCandidate{
			sampleCandidate("createRecord", "POST", "/zones/{zone_id}/dns_records", true, true),
			sampleCandidate("batchRecord", "POST", "/zones/{zone_id}/dns_records/batch", true, false), // no local ref
		},
	}

	view, _, err := ProjectCatalogDecision(rep, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Candidate 2 has no local references, so only candidate 1 is selectable
	if len(view.Selectable) != 1 || view.Selectable[0] != 1 {
		t.Fatalf("expected selectable [1], got %v", view.Selectable)
	}
}

func TestProjectCatalogDecisionNoQualifyingAPI(t *testing.T) {
	req := CatalogDiscoveryRequest{
		SchemaVersion: CatalogDiscoverySchemaVersion,
		Contract:      sampleContract(),
	}
	rep := CatalogDiscoveryReport{
		SchemaVersion:       CatalogDiscoverySchemaVersion,
		Outcome:             CatalogDiscoveryNoQualifyingAPI,
		QualifiedOperations: 0,
		Scope: CatalogDiscoveryScope{
			ProviderIDs: []string{"cloudflare"},
			Complete:    true,
		},
		Candidates: []CatalogDiscoveryCandidate{},
	}

	view, prompt, err := ProjectCatalogDecision(rep, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(view.Candidates) != 0 || len(view.Selectable) != 0 {
		t.Fatalf("expected empty candidates and selectable, got %v, %v", view.Candidates, view.Selectable)
	}
	if !strings.Contains(prompt, "No qualifying API was found in the complete examined scope") {
		t.Fatalf("unexpected negative prompt: %s", prompt)
	}
	if !view.Allows("browser") || !view.Allows("pending") || !view.Allows("reject") {
		t.Fatal("expected browser, pending, and reject allowed")
	}
	if view.Allows("select 1") {
		t.Fatal("expected select disallowed for negative outcome")
	}

	// Incomplete scope must be rejected
	repIncomplete := rep
	repIncomplete.Scope.Complete = false
	if _, _, err := ProjectCatalogDecision(repIncomplete, req); !errors.Is(err, ErrInvalidCatalogDecision) {
		t.Fatalf("expected ErrInvalidCatalogDecision for incomplete scope, got %v", err)
	}
}

func TestProjectCatalogDecisionBoundingAndRollback(t *testing.T) {
	req := CatalogDiscoveryRequest{
		SchemaVersion: CatalogDiscoverySchemaVersion,
		Contract:      sampleContract(),
	}

	// Create many candidates that together exceed MaxCatalogDecisionBytes (64 KiB)
	var candidates []CatalogDiscoveryCandidate
	longStr := strings.Repeat("x", 2000)
	for i := 0; i < 50; i++ {
		c := sampleCandidate("cand"+longStr, "POST", "/path/"+longStr, true, true)
		candidates = append(candidates, c)
	}

	rep := CatalogDiscoveryReport{
		SchemaVersion:       CatalogDiscoverySchemaVersion,
		Outcome:             CatalogDiscoveryAmbiguous,
		QualifiedOperations: 50,
		Scope: CatalogDiscoveryScope{
			ProviderIDs: []string{"cloudflare"},
			Complete:    true,
		},
		Candidates: candidates,
	}

	view, prompt, err := ProjectCatalogDecision(rep, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(view.Candidates) >= 50 {
		t.Fatalf("expected some candidates to be omitted due to bounding, got %d", len(view.Candidates))
	}
	if view.Omitted <= 0 {
		t.Fatalf("expected omitted > 0, got %d", view.Omitted)
	}
	if len(prompt) > MaxCatalogDecisionBytes {
		t.Fatalf("expected prompt <= 64 KiB, got %d", len(prompt))
	}
}

func TestProjectCatalogDecisionOversizedNegativeScopeRefusal(t *testing.T) {
	// F08: verify oversized base scope on negative outcome is rejected by bounded/Valid
	req := CatalogDiscoveryRequest{
		SchemaVersion: CatalogDiscoverySchemaVersion,
		Contract:      sampleContract(),
	}
	hugeIDs := make([]string, 10000)
	for i := range hugeIDs {
		hugeIDs[i] = "provider-very-long-identifier-" + strings.Repeat("a", 20)
	}
	rep := CatalogDiscoveryReport{
		SchemaVersion:       CatalogDiscoverySchemaVersion,
		Outcome:             CatalogDiscoveryNoQualifyingAPI,
		QualifiedOperations: 0,
		Scope: CatalogDiscoveryScope{
			ProviderIDs: hugeIDs,
			Complete:    true,
		},
	}
	if _, _, err := ProjectCatalogDecision(rep, req); !errors.Is(err, ErrInvalidCatalogDecision) {
		t.Fatalf("expected ErrInvalidCatalogDecision for oversized negative scope, got %v", err)
	}
}

func TestProjectCatalogDecisionClone(t *testing.T) {
	req := CatalogDiscoveryRequest{
		SchemaVersion: CatalogDiscoverySchemaVersion,
		Contract:      sampleContract(),
	}
	rep := CatalogDiscoveryReport{
		SchemaVersion:       CatalogDiscoverySchemaVersion,
		Outcome:             CatalogDiscoveryMatch,
		QualifiedOperations: 1,
		Scope: CatalogDiscoveryScope{
			ProviderIDs: []string{"cloudflare"},
			Complete:    true,
		},
		Candidates: []CatalogDiscoveryCandidate{
			sampleCandidate("createRecord", "POST", "/zones/{zone_id}/dns_records", true, true),
		},
	}
	view, _, err := ProjectCatalogDecision(rep, req)
	if err != nil {
		t.Fatal(err)
	}
	clone := view.Clone()
	if clone == nil || len(clone.Candidates) != 1 || clone.Candidates[0].Index != 1 {
		t.Fatalf("clone failed: %#v", clone)
	}
}
