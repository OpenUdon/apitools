package apitools

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestRankOperationCandidatesRanksAllContractDimensions(t *testing.T) {
	candidate := rankTestCandidate("sendMessage", "Send a message", OperationEffectWrite)
	contract := StepContract{
		Purpose: "send a message",
		Inputs: map[string]ContractValue{
			"body": {Type: "object", Required: boolPointer(true), Properties: map[string]ContractValue{
				"message": {Type: "string", Required: boolPointer(true)},
			}},
		},
		Outputs: map[string]ContractValue{"message_id": {Type: "string", Required: boolPointer(true)}},
		Effect:  OperationEffectWrite,
	}
	report, err := rankOperationCandidates(context.Background(), []OperationCandidate{candidate}, contract, PromptBudget{}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Candidates) != 1 {
		t.Fatalf("ranked candidates = %#v", report.Candidates)
	}
	got := report.Candidates[0].Match
	if got.Score != 100 || got.Purpose.Status != ContractMatchCompatible || got.Inputs.Status != ContractMatchCompatible || got.Outputs.Status != ContractMatchCompatible || got.Effect.Status != ContractMatchCompatible {
		t.Fatalf("full contract match = %#v", got)
	}
	if len(got.Purpose.Reasons) == 0 || len(got.Inputs.Reasons) == 0 || len(got.Outputs.Reasons) == 0 || len(got.Effect.Reasons) == 0 {
		t.Fatalf("match dimensions do not explain their scores: %#v", got)
	}
}

func TestRankOperationCandidatesHonorsInputTypeAndRequirednessDirection(t *testing.T) {
	candidate := rankTestCandidate("readNumber", "Read value", OperationEffectRead)
	candidate.Summary.Inputs = []OperationValueSummary{{
		Name: "count", Location: "query", Type: "integer", Required: boolPointer(true),
		Evidence: []OperationEvidence{{Kind: "query.parameter", Reference: "#/parameters/count"}},
	}}
	contract := StepContract{Purpose: "read value", Inputs: map[string]ContractValue{"count": {Type: "number", Required: boolPointer(true)}}}
	report, err := rankOperationCandidates(context.Background(), []OperationCandidate{candidate}, contract, PromptBudget{}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if report.Candidates[0].Match.Inputs.Status != ContractMatchIncompatible || len(report.Candidates[0].Match.Inputs.Conflicts) == 0 {
		t.Fatalf("broad step number was accepted by integer API input: %#v", report.Candidates[0].Match.Inputs)
	}

	contract.Inputs["count"] = ContractValue{Type: "integer", Required: boolPointer(false)}
	report, err = rankOperationCandidates(context.Background(), []OperationCandidate{candidate}, contract, PromptBudget{}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if report.Candidates[0].Match.Inputs.Status != ContractMatchIncompatible {
		t.Fatalf("optional step input satisfied required operation input: %#v", report.Candidates[0].Match.Inputs)
	}

	contract.Inputs["count"] = ContractValue{Type: "integer"}
	report, err = rankOperationCandidates(context.Background(), []OperationCandidate{candidate}, contract, PromptBudget{}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if report.Candidates[0].Match.Inputs.Status != ContractMatchIndeterminate {
		t.Fatalf("unknown step-input requiredness was treated as compatible: %#v", report.Candidates[0].Match.Inputs)
	}
}

func TestRankOperationCandidatesHonorsOutputDirectionAndMissingEvidence(t *testing.T) {
	candidate := rankTestCandidate("readNumber", "Read value", OperationEffectRead)
	candidate.Summary.Outputs = []OperationValueSummary{{Name: "count", Location: "response", Type: "integer", Required: boolPointer(true)}}
	contract := StepContract{Purpose: "read value", Outputs: map[string]ContractValue{"count": {Type: "number", Required: boolPointer(true)}}}
	report, err := rankOperationCandidates(context.Background(), []OperationCandidate{candidate}, contract, PromptBudget{}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if report.Candidates[0].Match.Outputs.Status != ContractMatchCompatible {
		t.Fatalf("integer output could not satisfy number expectation: %#v", report.Candidates[0].Match.Outputs)
	}

	contract.Outputs["count"] = ContractValue{Type: "integer", Required: boolPointer(true)}
	candidate.Summary.Outputs[0].Type = "number"
	report, err = rankOperationCandidates(context.Background(), []OperationCandidate{candidate}, contract, PromptBudget{}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if report.Candidates[0].Match.Outputs.Status != ContractMatchIncompatible {
		t.Fatalf("number response was accepted for integer expectation: %#v", report.Candidates[0].Match.Outputs)
	}

	candidate.Capabilities[1].Status = OperationCapabilityPartial
	candidate.Summary.Outputs = nil
	report, err = rankOperationCandidates(context.Background(), []OperationCandidate{candidate}, contract, PromptBudget{}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if report.Candidates[0].Match.Outputs.Status != ContractMatchIndeterminate {
		t.Fatalf("missing partial output schema was treated as a positive or definite negative: %#v", report.Candidates[0].Match.Outputs)
	}
}

func TestRankOperationCandidatesDoesNotAssumeMatchingNamedTypes(t *testing.T) {
	candidate := rankTestCandidate("readRecord", "Read record", OperationEffectRead)
	candidate.Summary.Outputs = []OperationValueSummary{{Name: "record", Location: "response", Type: "CustomerRecord", Required: boolPointer(true)}}
	contract := StepContract{Purpose: "read record", Outputs: map[string]ContractValue{"record": {Type: "CustomerRecord", Required: boolPointer(true)}}}
	report, err := rankOperationCandidates(context.Background(), []OperationCandidate{candidate}, contract, PromptBudget{}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if report.Candidates[0].Match.Outputs.Status != ContractMatchIndeterminate {
		t.Fatalf("same unresolved named type was assumed compatible: %#v", report.Candidates[0].Match.Outputs)
	}
}

func TestRankOperationCandidatesUnknownEffectDoesNotSatisfyReadConstraint(t *testing.T) {
	candidate := rankTestCandidate("findCustomer", "Find customer", OperationEffectUnknown)
	contract := StepContract{Purpose: "find customer", Effect: OperationEffectRead}
	report, err := rankOperationCandidates(context.Background(), []OperationCandidate{candidate}, contract, PromptBudget{}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if report.Candidates[0].Match.Effect.Status != ContractMatchIndeterminate {
		t.Fatalf("unknown effect satisfied read constraint: %#v", report.Candidates[0].Match.Effect)
	}
	if !hasDiagnosticCode(report.Diagnostics, "candidate.no_compatible_match") {
		t.Fatalf("no-compatible result not surfaced: %#v", report.Diagnostics)
	}

	candidate.Effect.Class = OperationEffectWrite
	report, err = rankOperationCandidates(context.Background(), []OperationCandidate{candidate}, contract, PromptBudget{}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if report.Candidates[0].Match.Effect.Status != ContractMatchIncompatible {
		t.Fatalf("write effect satisfied read constraint: %#v", report.Candidates[0].Match.Effect)
	}
}

func TestRankOperationCandidatesReportsPurposeGapsAndTies(t *testing.T) {
	first := rankTestCandidate("sendMessageA", "Send a message", OperationEffectWrite)
	second := rankTestCandidate("sendMessageB", "Send a message", OperationEffectWrite)
	second.Source.DocumentPath = "openapi/other.yaml"
	contract := rankTestContract("send a message")
	report, err := rankOperationCandidates(context.Background(), []OperationCandidate{second, first}, contract, PromptBudget{}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if report.Candidates[0].Source.DocumentPath != "openapi/messages.yaml" || report.Candidates[0].Match.Purpose.Score != report.Candidates[1].Match.Purpose.Score {
		t.Fatalf("unexpected tie ordering/scores: %#v", report.Candidates)
	}
	if !hasDiagnosticCode(report.Diagnostics, "candidate.top_tie") {
		t.Fatalf("top tie not surfaced: %#v", report.Diagnostics)
	}

	first.Summary.Description = "No shared topic"
	first.Summary.Evidence = []OperationEvidence{{Kind: "operation.summary", Reference: first.Source.Selector}}
	report, err = rankOperationCandidates(context.Background(), []OperationCandidate{first}, contract, PromptBudget{}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if report.Candidates[0].Match.Purpose.Status != ContractMatchIndeterminate || !hasDiagnosticCode(report.Diagnostics, "candidate.no_compatible_match") {
		t.Fatalf("unmatched natural language was treated as incompatibility or success: %#v", report)
	}
}

func TestRankOperationCandidatesDoesNotDoubleCountContractAliases(t *testing.T) {
	candidate := rankTestCandidate("getRecord", "Get record", OperationEffectRead)
	candidate.Summary.Inputs = []OperationValueSummary{{Name: "id", Location: "body", Type: "string", Required: boolPointer(true)}}
	candidate.Summary.Outputs = []OperationValueSummary{{Name: "id", Location: "response", Type: "string", Required: boolPointer(true)}}
	contract := StepContract{
		Inputs: map[string]ContractValue{
			"body.id": {Type: "string", Required: boolPointer(true)},
			"id":      {Type: "number", Required: boolPointer(true)},
		},
		Outputs: map[string]ContractValue{
			"body.id": {Type: "string", Required: boolPointer(true)},
			"id":      {Type: "string", Required: boolPointer(true)},
		},
	}
	report, err := rankOperationCandidates(context.Background(), []OperationCandidate{candidate}, contract, PromptBudget{}, 0)
	if err != nil {
		t.Fatal(err)
	}
	match := report.Candidates[0].Match
	if match.Inputs.Status != ContractMatchIndeterminate || match.Inputs.Score != 0 || !strings.Contains(strings.Join(match.Inputs.Gaps, " "), "multiple possible step-input aliases") {
		t.Fatalf("conflicting input aliases were treated as a positive match: %#v", match.Inputs)
	}
	if match.Outputs.Status != ContractMatchIndeterminate || match.Outputs.Score != 15 || !strings.Contains(strings.Join(match.Outputs.Gaps, " "), "resolve to the same source value") {
		t.Fatalf("duplicate output aliases were counted as separate evidence: %#v", match.Outputs)
	}
}

func TestRankOperationCandidatesBlocksInvalidOrOverbudgetContracts(t *testing.T) {
	candidate := rankTestCandidate("getMessage", "Get message", OperationEffectRead)
	cyclic := map[string]ContractValue{}
	root := ContractValue{Type: "object", Properties: cyclic}
	cyclic["self"] = root
	_, err := rankOperationCandidates(context.Background(), []OperationCandidate{candidate}, StepContract{Inputs: map[string]ContractValue{"body": root}}, PromptBudget{}, 0)
	if err == nil {
		t.Fatal("cyclic nested contract did not fail")
	}
	_, err = rankOperationCandidates(context.Background(), []OperationCandidate{candidate}, StepContract{Effect: "mutate"}, PromptBudget{}, 0)
	if err == nil {
		t.Fatal("unsupported effect did not fail")
	}
	_, err = rankOperationCandidates(context.Background(), []OperationCandidate{candidate}, rankTestContract("read message"), PromptBudget{MaxContextBytes: 1}, 0)
	if err == nil {
		t.Fatal("over-budget result did not fail")
	}
	_, err = rankOperationCandidates(context.Background(), []OperationCandidate{candidate, candidate}, rankTestContract("read message"), PromptBudget{}, 1)
	if err == nil {
		t.Fatal("truncated shortlist did not block")
	}
}

func TestRankOperationCandidatesDoesNotMutateInputs(t *testing.T) {
	candidate := rankTestCandidate("sendMessage", "Send message", OperationEffectWrite)
	contract := StepContract{Purpose: "send message", Inputs: map[string]ContractValue{"body": {Type: "object", Properties: map[string]ContractValue{"message": {Type: "string"}}}}}
	beforeCandidate := candidate
	beforePurpose := contract.Purpose
	_, err := rankOperationCandidates(context.Background(), []OperationCandidate{candidate}, contract, PromptBudget{}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if candidate.Match.Score != beforeCandidate.Match.Score || contract.Purpose != beforePurpose {
		t.Fatal("ranking mutated candidate or contract values")
	}
}

func TestRankOperationCandidatesHonorsCancellationBeforeWork(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := rankOperationCandidates(ctx, []OperationCandidate{rankTestCandidate("readMessage", "Read message", OperationEffectRead)}, rankTestContract("read message"), PromptBudget{}, 0)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled ranking error = %v, want context.Canceled", err)
	}
}

func rankTestCandidate(operationID, summary string, effect OperationEffect) OperationCandidate {
	selector := "#/paths/~1messages/post"
	return OperationCandidate{
		Source: OperationSourceIdentity{
			Kind: OperationSourceOpenAPI, DocumentPath: "openapi/messages.yaml",
			SHA256: strings.Repeat("a", 64), Selector: selector,
		},
		Operation: OperationSummary{ID: operationID, OperationID: operationID, Method: "POST", Path: "/messages", Summary: summary},
		Summary: ConsumerOperationSummary{
			Description: summary,
			Evidence:    []OperationEvidence{{Kind: "operation.summary", Reference: selector}},
			Inputs: []OperationValueSummary{
				{Name: "body", Location: "body", Type: "object", Required: boolPointer(true), Evidence: []OperationEvidence{{Kind: "request.body", Reference: selector + "/request"}}},
				{Name: "message", Location: "body", Type: "string", Required: boolPointer(true), ContainerRequired: boolPointer(true), Evidence: []OperationEvidence{{Kind: "request.body_field", Reference: selector + "/request/fields/message"}}},
			},
			Outputs: []OperationValueSummary{
				{Name: "body", Location: "response", Type: "object", Evidence: []OperationEvidence{{Kind: "response.body", Reference: selector + "/response"}}},
				{Name: "message_id", Location: "response", Type: "string", Required: boolPointer(true), Evidence: []OperationEvidence{{Kind: "response.field", Reference: selector + "/response/fields/message_id"}}},
			},
		},
		Capabilities: []OperationCapability{
			{Dimension: "inputs", Status: OperationCapabilitySupported},
			{Dimension: "outputs", Status: OperationCapabilitySupported},
		},
		Effect: EffectAssessment{Class: effect, Evidence: []OperationEvidence{{Kind: "operation.operation_id", Reference: selector}}, Reasons: []string{"test source evidence"}},
	}
}

func rankTestContract(purpose string) StepContract {
	return StepContract{
		Purpose: purpose,
		Inputs: map[string]ContractValue{"body": {Type: "object", Required: boolPointer(true), Properties: map[string]ContractValue{
			"message": {Type: "string", Required: boolPointer(true)},
		}}},
	}
}

func hasDiagnosticCode(diagnostics []Diagnostic, code string) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic.Code == code {
			return true
		}
	}
	return false
}
