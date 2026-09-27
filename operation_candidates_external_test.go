package apitools_test

import (
	"context"
	"testing"

	"github.com/OpenUdon/apitools"
)

func TestBuildOperationCandidatesPublicAPI(t *testing.T) {
	content := []byte(`openapi: 3.0.3
info: {title: Public API fixture, version: 1}
paths: {/messages: {post: {operationId: sendMessage, summary: Send a message, responses: {"200": {description: accepted}}}}}
`)
	required := true
	report, err := apitools.BuildOperationCandidates(context.Background(), apitools.OperationCandidateOptions{
		Sources: []apitools.OperationSourceInput{{Kind: apitools.OperationSourceOpenAPI, Content: content}},
		Contract: apitools.StepContract{
			Purpose: "send a message",
			Inputs:  map[string]apitools.ContractValue{"message": {Type: "string", Required: &required}},
			Effect:  apitools.OperationEffectWrite,
		},
	})
	if err != nil {
		t.Fatalf("BuildOperationCandidates() error = %v; diagnostics=%#v", err, report.Diagnostics)
	}
	if report.SchemaVersion != apitools.OperationCandidatesSchemaVersion || len(report.Candidates) != 1 {
		t.Fatalf("public candidate report = %#v", report)
	}
	if report.Candidates[0].Effect.Class != apitools.OperationEffectWrite || report.Candidates[0].Match.Effect.Status != apitools.ContractMatchCompatible {
		t.Fatalf("public candidate effect contract = %#v", report.Candidates[0])
	}
}
