package apitools

import (
	"context"
	"strings"
	"time"
)

// LLMClient is the pluggable vendor-neutral completion interface.
type LLMClient interface {
	Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
}

// CompletionRequest contains the bounded prompt payload for an LLM completion.
type CompletionRequest struct {
	SystemPrompt string        `json:"system_prompt"`
	UserPrompt   string        `json:"user_prompt"`
	Temperature  float32       `json:"temperature,omitempty"`
	MaxTokens    int           `json:"max_tokens,omitempty"`
	Timeout      time.Duration `json:"timeout,omitempty"`
}

// CompletionResponse contains the completion text and token usage.
type CompletionResponse struct {
	Content    string `json:"content"`
	TokensUsed int    `json:"tokens_used,omitempty"`
}

// CatalogQueryPlan represents structured search keywords and provider keys
// extracted from conversational developer intent.
type CatalogQueryPlan struct {
	ProviderKeys []string        `json:"provider_keys,omitempty"`
	Keywords     []string        `json:"keywords,omitempty"`
	Effect       OperationEffect `json:"effect,omitempty"`
	Purpose      string          `json:"purpose,omitempty"`
}

// StepContractDraft represents a prepared contract draft with readiness and diagnostics.
type StepContractDraft struct {
	Contract    StepContract `json:"contract"`
	Ready       bool         `json:"ready"`
	Diagnostics []Diagnostic `json:"diagnostics,omitempty"`
}

// DiscoveryRequestDraft represents a prepared discovery request draft with readiness and diagnostics.
type DiscoveryRequestDraft struct {
	Request     CatalogDiscoveryRequest `json:"request"`
	Ready       bool                    `json:"ready"`
	Diagnostics []Diagnostic            `json:"diagnostics,omitempty"`
}

// PrepareStepContractDraft prepares a StepContract, validating prompt budgets and effects.
func PrepareStepContractDraft(contract StepContract, budget PromptBudget) StepContractDraft {
	if budget.MaxTextRunes == 0 {
		budget = DefaultPromptBudget()
	}
	prepared, _, _, diagnostics := prepareStepContract(contract, budget)
	ready := true
	if strings.TrimSpace(contract.Purpose) == "" {
		ready = false
		diagnostics = append(diagnostics, Diagnostic{
			Severity: "error",
			Code:     "contract.purpose_required",
			Message:  "a documented step purpose is required",
		})
	}
	for _, diag := range diagnostics {
		if diag.Severity == "error" {
			ready = false
			break
		}
	}
	if len(contract.Inputs) == 0 {
		diagnostics = append(diagnostics, Diagnostic{
			Severity:    "info",
			Code:        "contract.input_missing",
			Message:     "no specific input parameters defined; discovery will rely on purpose wording and produce metadata-only leads without parameter qualification",
			Remediation: "Provide explicit parameter names if specific input matching is required.",
		})
	}
	return StepContractDraft{
		Contract:    prepared.contract,
		Ready:       ready,
		Diagnostics: diagnostics,
	}
}

// ValidateStepContractDraft runs validation against DefaultPromptBudget.
func ValidateStepContractDraft(contract StepContract) []Diagnostic {
	draft := PrepareStepContractDraft(contract, DefaultPromptBudget())
	return draft.Diagnostics
}

// PrepareDiscoveryRequestDraft prepares and validates a CatalogDiscoveryRequest draft.
func PrepareDiscoveryRequestDraft(request CatalogDiscoveryRequest, budget PromptBudget) DiscoveryRequestDraft {
	if budget.MaxTextRunes == 0 {
		budget = DefaultPromptBudget()
	}
	contractDraft := PrepareStepContractDraft(request.Contract, budget)
	request.Contract = contractDraft.Contract
	if request.SchemaVersion == "" {
		request.SchemaVersion = CatalogDiscoverySchemaVersion
	}
	if request.Limits.MaxOperations == 0 {
		request.Limits.MaxOperations = DefaultCatalogDiscoveryOperations
	}
	if request.Limits.MaxResults == 0 {
		request.Limits.MaxResults = DefaultCatalogDiscoveryResults
	}
	if request.Limits.MaxContextBytes == 0 {
		request.Limits.MaxContextBytes = DefaultCatalogDiscoveryContextBytes
	}
	if request.Limits.TimeoutMillis == 0 {
		request.Limits.TimeoutMillis = int(DefaultCatalogDiscoveryTimeout.Milliseconds())
	}

	var diagnostics []Diagnostic
	diagnostics = append(diagnostics, contractDraft.Diagnostics...)
	ready := contractDraft.Ready

	if request.SchemaVersion != CatalogDiscoverySchemaVersion {
		ready = false
		diagnostics = append(diagnostics, Diagnostic{
			Severity: "error",
			Code:     "discovery.version_invalid",
			Message:  "unsupported catalog discovery schema version",
		})
	}
	if request.Limits.MaxOperations < 0 || request.Limits.MaxOperations > MaxCatalogIndexOperations ||
		request.Limits.MaxResults < 0 || request.Limits.MaxResults > MaxCatalogDiscoveryResults ||
		request.Limits.MaxContextBytes < 0 || request.Limits.MaxContextBytes > MaxCatalogDiscoveryContextBytes ||
		request.Limits.TimeoutMillis < 0 || request.Limits.TimeoutMillis > int(MaxCatalogDiscoveryTimeout.Milliseconds()) {
		ready = false
		diagnostics = append(diagnostics, Diagnostic{
			Severity: "error",
			Code:     "discovery.limit_invalid",
			Message:  "catalog discovery limits exceed supported bounds",
		})
	}

	return DiscoveryRequestDraft{
		Request:     request,
		Ready:       ready,
		Diagnostics: diagnostics,
	}
}

// ValidateDiscoveryRequestDraft runs validation on a CatalogDiscoveryRequest draft.
func ValidateDiscoveryRequestDraft(request CatalogDiscoveryRequest) []Diagnostic {
	draft := PrepareDiscoveryRequestDraft(request, DefaultPromptBudget())
	return draft.Diagnostics
}
