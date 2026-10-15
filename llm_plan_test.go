package apitools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type mockTestLLMClient struct {
	completeFunc func(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
}

func (m *mockTestLLMClient) Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
	if m.completeFunc != nil {
		return m.completeFunc(ctx, req)
	}
	return CompletionResponse{}, nil
}

func TestPlanSearchQueryWithMockLLM(t *testing.T) {
	mock := &mockTestLLMClient{
		completeFunc: func(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
			return CompletionResponse{
				Content: `{"provider_keys": ["cloudflare"], "keywords": ["dns", "records"], "effect": "write", "purpose": "manage dns records"}`,
			}, nil
		},
	}

	plan, err := PlanSearchQuery(context.Background(), mock, "manage cloudflare dns records")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plan.ProviderKeys) != 1 || plan.ProviderKeys[0] != "cloudflare" {
		t.Errorf("expected provider cloudflare, got %v", plan.ProviderKeys)
	}
	if plan.Effect != OperationEffectWrite {
		t.Errorf("expected effect write, got %s", plan.Effect)
	}
	if plan.Purpose != "manage dns records" {
		t.Errorf("expected purpose 'manage dns records', got %s", plan.Purpose)
	}
}

func TestPlanSearchQueryDeterministicFallbackNilLLM(t *testing.T) {
	plan, err := PlanSearchQuery(context.Background(), nil, "create s3 bucket in aws")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if plan.Effect != OperationEffectWrite {
		t.Errorf("expected effect write, got %s", plan.Effect)
	}
	if len(plan.ProviderKeys) != 1 || plan.ProviderKeys[0] != "aws" {
		t.Errorf("expected provider aws, got %v", plan.ProviderKeys)
	}
}

func TestPlanSearchQueryDeterministicFallbackInvalidOutput(t *testing.T) {
	mock := &mockTestLLMClient{
		completeFunc: func(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
			return CompletionResponse{Content: "Here is your plan: not json"}, nil
		},
	}

	plan, err := PlanSearchQuery(context.Background(), mock, "list github repositories")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if plan.Effect != OperationEffectRead {
		t.Errorf("expected effect read, got %s", plan.Effect)
	}
	if len(plan.ProviderKeys) != 1 || plan.ProviderKeys[0] != "github" {
		t.Errorf("expected provider github, got %v", plan.ProviderKeys)
	}
}

func TestPlanSearchQueryResolvesF11IdempotentWrite(t *testing.T) {
	mock := &mockTestLLMClient{
		completeFunc: func(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
			return CompletionResponse{
				Content: `{"provider_keys": ["stripe"], "keywords": ["charges"], "effect": "idempotent_write", "purpose": "create charge"}`,
			}, nil
		},
	}

	plan, err := PlanSearchQuery(context.Background(), mock, "create stripe charge")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Finding F11: idempotent_write must be normalized to write
	if plan.Effect != OperationEffectWrite {
		t.Errorf("expected normalized effect write, got %s", plan.Effect)
	}
}

func TestPlanSearchQueryContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock := &mockTestLLMClient{
		completeFunc: func(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
			t.Fatal("mock should not be invoked when context is cancelled")
			return CompletionResponse{}, nil
		},
	}

	_, err := PlanSearchQuery(ctx, mock, "search records")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestPlanDiscoveryRequestWithMockLLM(t *testing.T) {
	mock := &mockTestLLMClient{
		completeFunc: func(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
			return CompletionResponse{
				Content: `{
					"purpose": "create dns record",
					"effect": "write",
					"inputs": {
						"zone_id": {"type": "string", "required": true},
						"name": {"type": "string", "required": true}
					}
				}`,
			}, nil
		},
	}

	draft, err := PlanDiscoveryRequest(context.Background(), mock, "create a dns record for my domain")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !draft.Ready {
		t.Fatalf("expected draft to be ready, got ready=false, diagnostics=%#v", draft.Diagnostics)
	}
	if draft.Request.Contract.Effect != OperationEffectWrite {
		t.Errorf("expected effect write, got %s", draft.Request.Contract.Effect)
	}
	if len(draft.Request.Contract.Inputs) != 2 {
		t.Errorf("expected 2 inputs, got %d", len(draft.Request.Contract.Inputs))
	}
	if in, ok := draft.Request.Contract.Inputs["zone_id"]; !ok || in.Type != "string" || in.Required == nil || !*in.Required {
		t.Errorf("unexpected zone_id input: %#v", in)
	}
}

func TestPlanDiscoveryRequestResolvesF11IdempotentWrite(t *testing.T) {
	mock := &mockTestLLMClient{
		completeFunc: func(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
			return CompletionResponse{
				Content: `{
					"purpose": "update record",
					"effect": "idempotent_write",
					"inputs": {
						"id": {"type": "string", "required": true}
					}
				}`,
			}, nil
		},
	}

	draft, err := PlanDiscoveryRequest(context.Background(), mock, "update record")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if draft.Request.Contract.Effect != OperationEffectWrite {
		t.Errorf("expected effect write, got %s", draft.Request.Contract.Effect)
	}
	// Verify that StepContract validation passes without contract.effect_invalid diagnostic
	for _, diag := range draft.Diagnostics {
		if diag.Code == "contract.effect_invalid" {
			t.Fatalf("contract.effect_invalid found: %#v", diag)
		}
	}
}

func TestPlanDiscoveryRequestResolvesF12NoWildcardAndClarification(t *testing.T) {
	mock := &mockTestLLMClient{
		completeFunc: func(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
			return CompletionResponse{
				Content: `{
					"purpose": "find user",
					"effect": "read",
					"inputs": {
						"*": {"type": "string", "required": false}
					}
				}`,
			}, nil
		},
	}

	draft, err := PlanDiscoveryRequest(context.Background(), mock, "find user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Wildcard input "*" must be stripped out (finding F12)
	if _, ok := draft.Request.Contract.Inputs["*"]; ok {
		t.Fatalf("wildcard input '*' should have been stripped out")
	}

	// Missing inputs must produce informational diagnostic for consumer review
	foundInfo := false
	for _, diag := range draft.Diagnostics {
		if diag.Code == "contract.input_missing" {
			foundInfo = true
			break
		}
	}
	if !foundInfo {
		t.Fatalf("expected contract.input_missing diagnostic, got %#v", draft.Diagnostics)
	}
}

func TestPlanDiscoveryRequestDeterministicFallbackNilLLM(t *testing.T) {
	draft, err := PlanDiscoveryRequest(context.Background(), nil, "fetch user profile")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !draft.Ready {
		t.Fatalf("expected draft to be ready, got diagnostics=%#v", draft.Diagnostics)
	}
	if draft.Request.Contract.Effect != OperationEffectRead {
		t.Errorf("expected effect read, got %s", draft.Request.Contract.Effect)
	}
	if len(draft.Request.Contract.Inputs) != 0 {
		t.Errorf("expected no inputs from deterministic fallback, got %v", draft.Request.Contract.Inputs)
	}
}

func TestClientLLMMethods(t *testing.T) {
	client := &Client{}
	plan, err := client.PlanSearchQuery(context.Background(), nil, "get weather")
	if err != nil || plan.Effect != OperationEffectRead {
		t.Fatalf("unexpected client.PlanSearchQuery result: plan=%#v err=%v", plan, err)
	}

	draft, err := client.PlanDiscoveryRequest(context.Background(), nil, "delete account")
	if err != nil || draft.Request.Contract.Effect != OperationEffectWrite {
		t.Fatalf("unexpected client.PlanDiscoveryRequest result: draft=%#v err=%v", draft, err)
	}
}

func TestPrepareAndValidateStepContractDraft(t *testing.T) {
	// Invalid effect
	invalidEffect := StepContract{
		Purpose: "test",
		Effect:  "invalid_effect",
	}
	draft := PrepareStepContractDraft(invalidEffect, DefaultPromptBudget())
	if draft.Ready {
		t.Fatalf("expected draft with invalid effect not to be ready")
	}
	found := false
	for _, diag := range draft.Diagnostics {
		if diag.Code == "contract.effect_invalid" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected contract.effect_invalid diagnostic")
	}

	// Empty purpose
	emptyPurpose := StepContract{
		Purpose: "",
		Effect:  OperationEffectRead,
	}
	draft2 := PrepareStepContractDraft(emptyPurpose, DefaultPromptBudget())
	if draft2.Ready {
		t.Fatalf("expected draft with empty purpose not to be ready")
	}
	diags := ValidateStepContractDraft(emptyPurpose)
	if len(diags) == 0 {
		t.Fatalf("expected diagnostics from ValidateStepContractDraft")
	}
}

func TestInvalidModelPlansRetainDeterministicFallback(t *testing.T) {
	invalid := []string{"null", "{}", `{"purpose":"override"}`, `{"purpose":null,"effect":"write"}`, `{"purpose":42,"effect":"read"}`, `{"purpose":"override","effect":false}`, `{"purpose":"override","effect":"write","provider_keys":{},"keywords":[]}`, `{"purpose":"override","effect":"write","inputs":[]}`, `{"purpose":"override","effect":"write","inputs":{"id":null}}`, `{"purpose":"override","effect":"write","inputs":{"id":{"type":null}}}`, `{"purpose":"override","effect":"write","inputs":{"id":{"type":"string","required":"yes"}}}`}
	invalid = append(invalid,
		`{"purpose":"`+strings.Repeat("x", DefaultPromptTextRunes+1)+`","effect":"write","provider_keys":[],"keywords":[],"inputs":{}}`,
		strings.Repeat(" ", MaxCatalogDecisionBytes)+`{}`,
	)
	ctx := context.Background()
	wantPlan, _ := PlanSearchQuery(ctx, nil, "list github users")
	wantDraft, _ := PlanDiscoveryRequest(ctx, nil, "list github users")
	for i, content := range invalid {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			mock := &mockTestLLMClient{completeFunc: func(context.Context, CompletionRequest) (CompletionResponse, error) {
				return CompletionResponse{Content: content}, nil
			}}
			gotPlan, err := PlanSearchQuery(ctx, mock, "list github users")
			if err != nil || !reflect.DeepEqual(gotPlan, wantPlan) {
				t.Errorf("query fallback: %#v %v", gotPlan, err)
			}
			gotDraft, err := PlanDiscoveryRequest(ctx, mock, "list github users")
			if err != nil || !reflect.DeepEqual(gotDraft, wantDraft) {
				t.Errorf("discovery fallback: %#v %v", gotDraft, err)
			}
		})
	}
}

func TestModelPlanFencesAndEmptyInputs(t *testing.T) {
	budget := DefaultPromptBudget()
	query := "```json\n" + `{"provider_keys":["github"],"keywords":["users"],"effect":"READ","purpose":"list users"}` + "\n```"
	if plan, ok := parseQueryPlanResponse(query, budget); !ok || plan.Effect != OperationEffectRead {
		t.Fatalf("fenced query: %#v %v", plan, ok)
	}
	draft, ok := parseDiscoveryRequestResponse("```json\n"+`{"purpose":"list users","effect":"idempotent_write","inputs":{}}`+"\n```", "original", budget)
	if !ok || !draft.Ready || draft.Request.Contract.Effect != OperationEffectWrite {
		t.Fatalf("valid draft: %#v %v", draft, ok)
	}
	found := false
	for _, d := range draft.Diagnostics {
		found = found || d.Code == "contract.input_missing"
	}
	if !found {
		t.Fatal("missing clarification")
	}
}

func TestModelPlanStructuralAndFieldBudgets(t *testing.T) {
	budget := DefaultPromptBudget()
	validQuery := `{"provider_keys":[],"keywords":["users"],"purpose":"list users","effect":"read"}`
	var query map[string]any
	json.Unmarshal([]byte(validQuery), &query)
	for _, mutate := range []func(map[string]any){
		func(m map[string]any) { m["keywords"] = make([]string, budget.MaxCollectionItems+1) },
		func(m map[string]any) { m["provider_keys"] = []any{nil} },
		func(m map[string]any) { m["keywords"] = []string{strings.Repeat("x", budget.MaxIdentifierRunes+1)} },
		func(m map[string]any) { m["effect"] = "" },
		func(m map[string]any) { m["effect"] = "invalid" },
	} {
		json.Unmarshal([]byte(validQuery), &query)
		mutate(query)
		data, _ := json.Marshal(query)
		if _, ok := parseQueryPlanResponse(string(data), budget); ok {
			t.Fatalf("accepted invalid query: %s", data)
		}
	}
	inputs := map[string]any{}
	for i := 0; i <= budget.MaxFields; i++ {
		inputs[fmt.Sprint("field", i)] = map[string]any{"type": "string"}
	}
	data, _ := json.Marshal(map[string]any{"purpose": "list users", "effect": "read", "inputs": inputs})
	for _, raw := range []string{
		string(data),
		`{"purpose":"list users","effect":"read","inputs":{"id":{"type":"string","required":null}}}`,
		`{"purpose":"list users","effect":"read","inputs":{"id":{"type":"string","required":123}}}`,
		`{"purpose":"list users","effect":"read","inputs":null}`,
		"[" + validQuery + "]",
		validQuery + " " + validQuery,
	} {
		if _, ok := parseDiscoveryRequestResponse(raw, "original", budget); ok {
			t.Fatalf("accepted invalid discovery plan: %.200s", raw)
		}
	}
}
