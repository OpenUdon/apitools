package apitools

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func makeTestDecisionView(t *testing.T) CatalogDecisionView {
	t.Helper()
	op1 := OperationSummary{
		ID:          "create_record",
		OperationID: "createRecord",
		Method:      "POST",
		Path:        "/records",
		Summary:     "Create a record",
	}
	op2 := OperationSummary{
		ID:          "batch_create_records",
		OperationID: "batchCreateRecords",
		Method:      "POST",
		Path:        "/records/batch",
		Summary:     "Batch create records",
	}

	c1 := CatalogDiscoveryCandidate{
		Qualified: true,
		References: []CatalogArtifactReference{
			{ProviderID: "prov1", SpecRefID: "spec1", ArtifactID: "art1"},
		},
		Candidate: OperationCandidate{
			Operation: op1,
		},
	}
	c2 := CatalogDiscoveryCandidate{
		Qualified: true,
		References: []CatalogArtifactReference{
			{ProviderID: "prov1", SpecRefID: "spec1", ArtifactID: "art2"},
		},
		Candidate: OperationCandidate{
			Operation: op2,
		},
	}

	report := CatalogDiscoveryReport{
		Outcome:             CatalogDiscoveryAmbiguous,
		QualifiedOperations: 2,
		Candidates:          []CatalogDiscoveryCandidate{c1, c2},
		Scope: CatalogDiscoveryScope{
			Complete: true,
		},
	}

	view, _, err := ProjectCatalogDecision(report, CatalogDiscoveryRequest{
		SchemaVersion: CatalogDiscoverySchemaVersion,
		Contract:      StepContract{Purpose: "create record", Effect: OperationEffectWrite},
	})
	if err != nil {
		t.Fatalf("failed to project decision view: %v", err)
	}
	return view
}

func TestRecommendCandidateWithMockLLM(t *testing.T) {
	view := makeTestDecisionView(t)
	mock := &mockTestLLMClient{
		completeFunc: func(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
			return CompletionResponse{
				Content: `{"recommended_index": 2, "rationale": "Matches batch record creation perfectly.", "caveats": "Requires payload array."}`,
			}, nil
		},
	}

	rec, updatedView, err := RecommendCandidate(context.Background(), mock, "create multiple records", view)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.RecommendedIndex != 2 {
		t.Errorf("expected recommended index 2, got %d", rec.RecommendedIndex)
	}
	if !strings.Contains(updatedView.Explanation, "Recommended [#2]") {
		t.Errorf("expected explanation to mention [#2], got: %s", updatedView.Explanation)
	}
	if !updatedView.Valid() {
		t.Errorf("expected updated view to be valid")
	}
}

func TestRecommendCandidateResolvesF13HallucinatedOrUnqualifiedIndex(t *testing.T) {
	view := makeTestDecisionView(t)
	// Make candidate #2 unqualified in a modified view
	view.Selectable = []int{1}
	view.Candidates[1].Evidence.Qualified = false
	view.Candidates[1].Evidence.References = nil

	// Mock LLM tries to recommend non-selectable index 2 (resolving F13)
	mock := &mockTestLLMClient{
		completeFunc: func(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
			return CompletionResponse{
				Content: `{"recommended_index": 2, "rationale": "I hallucinated candidate 2."}`,
			}, nil
		},
	}

	rec, updatedView, err := RecommendCandidate(context.Background(), mock, "create record", view)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Fallback to deterministic selectable candidate #1
	if rec.RecommendedIndex != 1 {
		t.Errorf("expected fallback to selectable index 1, got %d", rec.RecommendedIndex)
	}
	if !strings.Contains(updatedView.Explanation, "Recommended [#1]") {
		t.Errorf("expected explanation to mention [#1], got: %s", updatedView.Explanation)
	}
}

func TestRecommendCandidateEmptySelectableAbstains(t *testing.T) {
	view := makeTestDecisionView(t)
	view.Selectable = []int{}

	rec, _, err := RecommendCandidate(context.Background(), nil, "create record", view)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.RecommendedIndex != 0 {
		t.Errorf("expected recommended index 0 (abstain), got %d", rec.RecommendedIndex)
	}
}

func TestRecommendCandidateDeterministicFallbackNilLLM(t *testing.T) {
	view := makeTestDecisionView(t)
	rec, updatedView, err := RecommendCandidate(context.Background(), nil, "create record", view)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.RecommendedIndex != 1 {
		t.Errorf("expected first selectable candidate 1, got %d", rec.RecommendedIndex)
	}
	if !strings.Contains(updatedView.Explanation, "Recommended [#1]") {
		t.Errorf("expected explanation to mention [#1], got: %s", updatedView.Explanation)
	}
	if !updatedView.Valid() {
		t.Errorf("expected updated view to be valid")
	}
}

func TestRecommendCandidateContextCancellation(t *testing.T) {
	view := makeTestDecisionView(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err := RecommendCandidate(ctx, nil, "create record", view)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestRecommendCandidatePreservesViewBounds(t *testing.T) {
	view := makeTestDecisionView(t)
	// Very long explanation
	mock := &mockTestLLMClient{
		completeFunc: func(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
			return CompletionResponse{
				Content: `{"recommended_index": 1, "rationale": "` + strings.Repeat("a", 70000) + `"}`,
			}, nil
		},
	}

	rec, updatedView, err := RecommendCandidate(context.Background(), mock, "create record", view)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.RecommendedIndex != 1 {
		t.Errorf("expected recommended index 1, got %d", rec.RecommendedIndex)
	}
	data, _ := json.Marshal(updatedView)
	if len(data) > MaxCatalogDecisionBytes {
		t.Errorf("updatedView JSON exceeds MaxCatalogDecisionBytes: %d", len(data))
	}
	if len(updatedView.Prompt()) > MaxCatalogDecisionBytes {
		t.Errorf("updatedView Prompt exceeds MaxCatalogDecisionBytes: %d", len(updatedView.Prompt()))
	}
}

func TestClientRecommendCandidate(t *testing.T) {
	view := makeTestDecisionView(t)
	client := &Client{}
	rec, _, err := client.RecommendCandidate(context.Background(), nil, "create record", view)
	if err != nil || rec.RecommendedIndex != 1 {
		t.Fatalf("unexpected client.RecommendCandidate: rec=%#v err=%v", rec, err)
	}
}
