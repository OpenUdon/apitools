package apitools

import (
	"strings"
	"testing"
)

func TestAssessOperationEffectUsesMeaningNotHTTPMethod(t *testing.T) {
	tests := []struct {
		name      string
		operation OperationSummary
		want      OperationEffect
	}{
		{name: "POST search is a read", operation: OperationSummary{OperationID: "searchMessages", Method: "POST"}, want: OperationEffectRead},
		{name: "GET update is a write", operation: OperationSummary{OperationID: "updateCustomer", Method: "GET"}, want: OperationEffectWrite},
		{name: "summary read despite POST", operation: OperationSummary{Summary: "Retrieve project details", Method: "POST"}, want: OperationEffectRead},
		{name: "summary write despite GET", operation: OperationSummary{Summary: "Send a notification", Method: "GET"}, want: OperationEffectWrite},
		{name: "POST-style name identifies write", operation: OperationSummary{OperationID: "postMessage", Method: "GET"}, want: OperationEffectWrite},
		{name: "method alone unknown", operation: OperationSummary{Method: "GET"}, want: OperationEffectUnknown},
		{name: "neutral name unknown", operation: OperationSummary{OperationID: "processCustomer", Method: "POST"}, want: OperationEffectUnknown},
		{name: "noun after read verb", operation: OperationSummary{OperationID: "getCreatedAt"}, want: OperationEffectRead},
		{name: "unknown leading operation id before read token", operation: OperationSummary{OperationID: "markMessageRead"}, want: OperationEffectUnknown},
		{name: "unknown leading summary action before read token", operation: OperationSummary{Summary: "Mark message as read"}, want: OperationEffectUnknown},
		{name: "unknown save search operation id", operation: OperationSummary{OperationID: "saveSearch"}, want: OperationEffectUnknown},
		{name: "unknown save search summary action", operation: OperationSummary{Summary: "Save search"}, want: OperationEffectUnknown},
		{name: "conflicting operation and summary", operation: OperationSummary{OperationID: "getCustomer", Summary: "Delete customer"}, want: OperationEffectUnknown},
		{name: "compound read and write", operation: OperationSummary{OperationID: "getOrCreateCustomer"}, want: OperationEffectUnknown},
		{name: "distant compound action", operation: OperationSummary{OperationID: "getMessage", Description: "Retrieves the next pending message and deletes it from the queue."}, want: OperationEffectUnknown},
		{name: "second sentence mutation", operation: OperationSummary{Summary: "Retrieve the message. Delete it after retrieval."}, want: OperationEffectUnknown},
		{name: "later negated action", operation: OperationSummary{Summary: "Retrieve the message and do not delete it."}, want: OperationEffectUnknown},
		{name: "description conflicts with summary", operation: OperationSummary{Summary: "Get customer", Description: "Delete the customer"}, want: OperationEffectUnknown},
		{name: "negated summary write", operation: OperationSummary{Summary: "Does not delete the customer"}, want: OperationEffectUnknown},
		{name: "negated summary read", operation: OperationSummary{Summary: "Doesn't retrieve the customer"}, want: OperationEffectUnknown},
		{name: "negated operation id", operation: OperationSummary{OperationID: "notDeleteCustomer"}, want: OperationEffectUnknown},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := assessOperationEffect(test.operation, "operation:fixture")
			if got.Class != test.want {
				t.Fatalf("effect = %q, want %q; assessment=%#v", got.Class, test.want, got)
			}
			if len(got.Reasons) == 0 {
				t.Fatalf("effect assessment has no public rationale: %#v", got)
			}
			if test.operation.Method != "" && got.Class == OperationEffectUnknown && test.operation.OperationID == "" && test.operation.Summary == "" && len(got.Evidence) != 0 {
				t.Fatalf("method-only classification unexpectedly has evidence: %#v", got)
			}
		})
	}
	negated := assessOperationEffect(OperationSummary{Summary: "Does not delete the customer"}, "operation:negated")
	if negated.Class != OperationEffectUnknown || len(negated.Evidence) != 1 || negated.Evidence[0].Kind != "operation.summary" {
		t.Fatalf("negated effect did not retain its source evidence: %#v", negated)
	}
}

func TestAssessOperationEffectCombinesNativeProtocolEvidence(t *testing.T) {
	query := OperationSummary{Method: "POST"}
	got := assessOperationEffectWithNative(query, "graphql:query", []effectSignal{{
		effect: OperationEffectRead, kind: "graphql.operation_kind", reference: "graphql:query", word: "query",
	}})
	if got.Class != OperationEffectRead || len(got.Evidence) != 1 || got.Evidence[0].Kind != "graphql.operation_kind" {
		t.Fatalf("native GraphQL query evidence not preserved: %#v", got)
	}

	namespaced := OperationSummary{OperationID: "pets.messages.get"}
	got = assessOperationEffect(namespaced, "#/methods/pets.messages.get")
	if got.Class != OperationEffectRead {
		t.Fatalf("Google Discovery method suffix was not classified: %#v", got)
	}

	resourceName := OperationSummary{OperationID: "pet"}
	got = assessOperationEffectWithNative(resourceName, "graphql:query", []effectSignal{{
		effect: OperationEffectRead, kind: "graphql.operation_kind", reference: "graphql:query", word: "query",
	}})
	if got.Class != OperationEffectRead {
		t.Fatalf("native GraphQL operation kind was masked by a resource-name ID: %#v", got)
	}

	conflict := OperationSummary{OperationID: "deleteCustomer"}
	got = assessOperationEffectWithNative(conflict, "graphql:query", []effectSignal{{
		effect: OperationEffectRead, kind: "graphql.operation_kind", reference: "graphql:query", word: "query",
	}})
	if got.Class != OperationEffectUnknown {
		t.Fatalf("conflicting native/text evidence = %#v, want unknown", got)
	}
}

func TestAssessOperationEffectBoundsUntrustedText(t *testing.T) {
	assessment := assessOperationEffect(OperationSummary{Description: strings.Repeat("word ", 1<<18)}, "operation:large-description")
	if assessment.Class != OperationEffectUnknown || len(assessment.Evidence) != 1 || !strings.Contains(strings.Join(assessment.Reasons, " "), "bounded effect-analysis budget") {
		t.Fatalf("oversized source wording did not remain bounded and unknown: %#v", assessment)
	}
}
