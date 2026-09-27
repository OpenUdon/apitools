package apitools

import (
	"reflect"
	"strings"
	"testing"
)

func TestSummarizeOperationForConsumerUsesSourceEvidence(t *testing.T) {
	operation := OperationSummary{
		ID:          "lookupCustomer",
		OperationID: "lookupCustomer",
		Method:      "GET",
		Path:        "/customers/{customerId}",
		Summary:     "Look up a customer",
		Parameters: []ParameterSummary{
			{Name: "customerId", In: "path", Type: "string", Required: true},
			{Name: "api_key", In: "query", Type: "string", Required: true},
		},
		RequestBody: &RequestBodySummary{Required: false, Fields: []RequestFieldSummary{{Path: "filter.active", Type: "boolean"}}},
		ResponseBody: &ResponseBodySummary{Fields: []RequestFieldSummary{
			{Path: "name", Type: "string", Required: true},
			{Path: "createdAt", Type: "string", Format: "date-time"},
		}},
	}
	before := operation
	summary, diagnostics, err := summarizeOperationForConsumer(operation, "#/paths/~1customers~1{customerId}/get", PromptBudget{})
	if err != nil {
		t.Fatal(err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("summary diagnostics = %#v", diagnostics)
	}
	if summary.Description != "Look up a customer" || len(summary.Inputs) != 3 || len(summary.Outputs) != 3 {
		t.Fatalf("consumer summary = %#v", summary)
	}
	if summary.Inputs[2].Name != "customerId" || summary.Inputs[2].Location != "path" || summary.Inputs[2].Required == nil || !*summary.Inputs[2].Required {
		t.Fatalf("required path input = %#v", summary.Inputs[2])
	}
	if summary.Outputs[1].Name != "createdAt" || summary.Outputs[1].Required == nil || *summary.Outputs[1].Required {
		t.Fatalf("optional response output = %#v", summary.Outputs[1])
	}
	if summary.Inputs[0].Name != "body" || summary.Inputs[0].Required == nil || *summary.Inputs[0].Required {
		t.Fatalf("optional request-body container requirement = %#v", summary.Inputs[0])
	}
	if summary.Inputs[1].Name != "filter.active" || summary.Inputs[1].ContainerRequired == nil || *summary.Inputs[1].ContainerRequired {
		t.Fatalf("optional nested request-body requirement = %#v", summary.Inputs[1])
	}
	if len(summary.Evidence) == 0 || summary.Evidence[0].Kind != "operation.summary" || len(summary.Inputs[2].Evidence) == 0 {
		t.Fatalf("summary evidence missing: %#v", summary)
	}
	if strings.Count(summary.Inputs[2].Evidence[0].Reference, "#") != 1 || !strings.HasSuffix(summary.Inputs[2].Evidence[0].Reference, "/parameters/path/customerId") {
		t.Fatalf("source location reference is malformed: %#v", summary.Inputs[2].Evidence[0])
	}
	if !strings.Contains(strings.Join(summary.Gaps, " "), "Credential-shaped parameter api_key is omitted") {
		t.Fatalf("credential-shaped input was not kept out of workflow data: %#v", summary.Gaps)
	}
	if !reflect.DeepEqual(operation, before) {
		t.Fatalf("summary mutated caller operation: %#v", operation)
	}
}

func TestSummarizeOperationForConsumerUsesHonestFallbackAndGaps(t *testing.T) {
	operation := OperationSummary{ID: "opaqueCall", Method: "POST", Path: "/things"}
	summary, _, err := summarizeOperationForConsumer(operation, "rpc:opaqueCall", PromptBudget{})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(summary.Description, "documented purpose is unavailable") || len(summary.Inputs) != 0 || len(summary.Outputs) != 0 {
		t.Fatalf("fallback summary invents operation semantics: %#v", summary)
	}
	if len(summary.Gaps) < 3 {
		t.Fatalf("missing operation details did not produce visible gaps: %#v", summary.Gaps)
	}
}

func TestSummarizeOperationForConsumerSanitizesUntrustedText(t *testing.T) {
	operation := OperationSummary{ID: "notify", Summary: "\x1b[31mSend update to user\x1b[0m"}
	summary, diagnostics, err := summarizeOperationForConsumer(operation, "op", PromptBudget{MaxTextRunes: 12})
	if err != nil {
		t.Fatal(err)
	}
	if strings.ContainsRune(summary.Description, '\x1b') || len([]rune(summary.Description)) > 12 {
		t.Fatalf("summary text not bounded/sanitized: %q", summary.Description)
	}
	if len(diagnostics) == 0 || !strings.HasPrefix(diagnostics[0].Code, "prompt.") {
		t.Fatalf("sanitization was not reported: %#v", diagnostics)
	}
}

func TestSummarizeOperationForConsumerBlocksSemanticFieldTruncation(t *testing.T) {
	operation := OperationSummary{ID: "lookup", Parameters: []ParameterSummary{
		{Name: "a", In: "query"}, {Name: "b", In: "query"}, {Name: "c", In: "query"},
	}}
	_, diagnostics, err := summarizeOperationForConsumer(operation, "lookup", PromptBudget{MaxFields: 2})
	if err == nil || len(diagnostics) != 1 || diagnostics[0].Code != "operation_summary.input_budget" {
		t.Fatalf("over-budget required inputs were not blocked: diagnostics=%#v err=%v", diagnostics, err)
	}

	operation.Parameters = operation.Parameters[:2]
	_, diagnostics, err = summarizeOperationForConsumer(operation, "lookup", PromptBudget{MaxCollectionItems: 1})
	if err == nil || len(diagnostics) != 1 || diagnostics[0].Code != "operation_summary.parameter_budget" {
		t.Fatalf("over-budget parameter list was not blocked: diagnostics=%#v err=%v", diagnostics, err)
	}
}

func TestSummarizeOperationForConsumerRejectsUnsafeSelector(t *testing.T) {
	_, diagnostics, err := summarizeOperationForConsumer(OperationSummary{ID: "lookup"}, "lookup\nforged", PromptBudget{})
	if err == nil || len(diagnostics) != 1 || diagnostics[0].Code != "operation_summary.selector_unsafe" {
		t.Fatalf("unsafe selector was not rejected: diagnostics=%#v err=%v", diagnostics, err)
	}
}
