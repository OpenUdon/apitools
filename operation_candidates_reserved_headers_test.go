package apitools

import (
	"context"
	"encoding/json"
	"testing"
)

func TestOperationCandidatesReservedHeaders(t *testing.T) {
	fixture := completeShapeFixture()
	operation := mapValue(mapValue(mapValue(fixture["paths"])["/pets/{id}"])["get"])
	operation["summary"] = "Get pet"
	operation["parameters"] = append(operation["parameters"].([]any),
		map[string]any{"name": "Accept", "in": "header", "required": true, "schema": map[string]any{"type": "string"}},
		map[string]any{"name": "Content-Type", "in": "header", "required": true, "schema": map[string]any{"type": "string"}},
		map[string]any{"name": "Authorization", "in": "header", "required": true, "schema": map[string]any{"type": "string"}})
	source, err := json.Marshal(fixture)
	if err != nil {
		t.Fatal(err)
	}
	inventory, err := BuildOperationInventory(context.Background(), InventoryOptions{Documents: []InventoryDocument{{Content: source}}})
	if err != nil || len(inventory.Operations) != 1 || len(inventory.Operations[0].Parameters) != 1 {
		t.Fatalf("effective inventory: %+v %v", inventory, err)
	}
	report, err := BuildOperationCandidates(context.Background(), OperationCandidateOptions{
		Sources:  []OperationSourceInput{{Kind: OperationSourceOpenAPI, Content: source}},
		Contract: StepContract{Purpose: "Get pet", Inputs: map[string]ContractValue{"id": {Type: "string", Required: boolPointer(true)}}, Effect: OperationEffectRead},
	})
	if err != nil || len(report.Candidates) != 1 {
		t.Fatalf("candidate report: %+v %v", report, err)
	}
	candidate := report.Candidates[0]
	if len(candidate.Summary.Inputs) != 1 || len(candidate.Summary.Gaps) != 0 || candidate.Match.Inputs.Status != ContractMatchCompatible {
		t.Fatalf("reserved declaration constrained candidate: %+v", candidate)
	}
}
