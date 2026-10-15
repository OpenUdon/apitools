package apitools

import (
	"context"
	"encoding/json"
	"strings"
	"time"
)

// DefaultLLMTimeout is the default deadline for LLM completion requests.
const DefaultLLMTimeout = 15 * time.Second

// PlanSearchQuery translates freeform developer intent into a structured search plan.
// If llm is nil, times out, or returns invalid output, it falls back deterministically.
// Context cancellation stops immediately without fallback.
func PlanSearchQuery(ctx context.Context, llm LLMClient, goal string) (CatalogQueryPlan, error) {
	if ctx.Err() != nil {
		return CatalogQueryPlan{}, ctx.Err()
	}

	budget := DefaultPromptBudget()
	sanitizedGoal, _ := sanitizePromptString(goal, budget.MaxTextRunes)
	if strings.TrimSpace(sanitizedGoal) == "" {
		return CatalogQueryPlan{}, nil
	}

	if llm == nil {
		return deterministicPlanSearchQuery(sanitizedGoal, budget), nil
	}

	systemPrompt := `You are an API specification search assistant.
Given a natural language developer goal, extract:
1. provider_keys: lowercase provider or service names (e.g. "cloudflare", "aws", "stripe").
2. keywords: 2 to 5 specific search terms (e.g. "dns", "records", "zones").
3. effect: "read" (fetch/get/list) or "write" (create/update/delete) or "unknown".
4. purpose: a concise 3-6 word task description.

Respond ONLY with a valid JSON object matching this schema:
{
  "provider_keys": ["string"],
  "keywords": ["string"],
  "effect": "read" | "write" | "unknown",
  "purpose": "string"
}`

	req := CompletionRequest{
		SystemPrompt: systemPrompt,
		UserPrompt:   "Goal: " + sanitizedGoal,
		Timeout:      DefaultLLMTimeout,
	}

	resp, err := llm.Complete(ctx, req)
	if err != nil {
		if ctx.Err() != nil {
			return CatalogQueryPlan{}, ctx.Err()
		}
		return deterministicPlanSearchQuery(sanitizedGoal, budget), nil
	}

	plan, ok := parseQueryPlanResponse(resp.Content, budget)
	if !ok {
		return deterministicPlanSearchQuery(sanitizedGoal, budget), nil
	}
	return plan, nil
}

// PlanSearchQuery on *Client delegates to PlanSearchQuery.
func (c *Client) PlanSearchQuery(ctx context.Context, llm LLMClient, goal string) (CatalogQueryPlan, error) {
	return PlanSearchQuery(ctx, llm, goal)
}

// PlanDiscoveryRequest synthesizes a structured CatalogDiscoveryRequest draft from a task description.
// If llm is nil, times out, or returns invalid output, it falls back deterministically.
// Context cancellation stops immediately without fallback.
func PlanDiscoveryRequest(ctx context.Context, llm LLMClient, taskDescription string) (DiscoveryRequestDraft, error) {
	if ctx.Err() != nil {
		return DiscoveryRequestDraft{}, ctx.Err()
	}

	budget := DefaultPromptBudget()
	sanitizedTask, _ := sanitizePromptString(taskDescription, budget.MaxTextRunes)
	if strings.TrimSpace(sanitizedTask) == "" {
		return PrepareDiscoveryRequestDraft(CatalogDiscoveryRequest{}, budget), nil
	}

	if llm == nil {
		return deterministicPlanDiscoveryRequest(sanitizedTask, budget), nil
	}

	systemPrompt := `You are an API integration engineer. Given a specific API task, synthesize a structured StepContract.
Extract:
- purpose: clear summary of the operation.
- effect: "read", "write", or "unknown".
- inputs: key-value dictionary of required parameter names and types ("string", "integer", "boolean").

Respond ONLY with a valid JSON object matching:
{
  "purpose": "string",
  "effect": "read" | "write" | "unknown",
  "inputs": {
    "<param_name>": {"type": "string" | "integer" | "boolean", "required": true}
  }
}`

	req := CompletionRequest{
		SystemPrompt: systemPrompt,
		UserPrompt:   "Task: " + sanitizedTask,
		Timeout:      DefaultLLMTimeout,
	}

	resp, err := llm.Complete(ctx, req)
	if err != nil {
		if ctx.Err() != nil {
			return DiscoveryRequestDraft{}, ctx.Err()
		}
		return deterministicPlanDiscoveryRequest(sanitizedTask, budget), nil
	}

	draft, ok := parseDiscoveryRequestResponse(resp.Content, sanitizedTask, budget)
	if !ok {
		return deterministicPlanDiscoveryRequest(sanitizedTask, budget), nil
	}
	return draft, nil
}

// PlanDiscoveryRequest on *Client delegates to PlanDiscoveryRequest.
func (c *Client) PlanDiscoveryRequest(ctx context.Context, llm LLMClient, taskDescription string) (DiscoveryRequestDraft, error) {
	return PlanDiscoveryRequest(ctx, llm, taskDescription)
}

func parseQueryPlanResponse(content string, budget PromptBudget) (CatalogQueryPlan, bool) {
	raw := extractJSON(content)
	if len(raw) == 0 || len(raw) > MaxCatalogDecisionBytes {
		return CatalogQueryPlan{}, false
	}

	var parsed struct {
		ProviderKeys []string `json:"provider_keys"`
		Keywords     []string `json:"keywords"`
		Effect       string   `json:"effect"`
		Purpose      string   `json:"purpose"`
	}

	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return CatalogQueryPlan{}, false
	}

	effect := normalizeEffect(parsed.Effect)
	purpose, _ := sanitizePromptString(parsed.Purpose, budget.MaxTextRunes)

	var providerKeys []string
	for _, key := range parsed.ProviderKeys {
		k := strings.ToLower(strings.TrimSpace(key))
		if k != "" && !strings.ContainsAny(k, " *?") {
			if safe, _ := sanitizePromptString(k, budget.MaxIdentifierRunes); safe != "" {
				providerKeys = append(providerKeys, safe)
			}
		}
	}

	var keywords []string
	for _, kw := range parsed.Keywords {
		w := strings.ToLower(strings.TrimSpace(kw))
		if w != "" && !strings.ContainsAny(w, " *?") {
			if safe, _ := sanitizePromptString(w, budget.MaxIdentifierRunes); safe != "" {
				keywords = append(keywords, safe)
			}
		}
	}

	return CatalogQueryPlan{
		ProviderKeys: providerKeys,
		Keywords:     keywords,
		Effect:       effect,
		Purpose:      purpose,
	}, true
}

func parseDiscoveryRequestResponse(content, originalTask string, budget PromptBudget) (DiscoveryRequestDraft, bool) {
	raw := extractJSON(content)
	if len(raw) == 0 || len(raw) > MaxCatalogDecisionBytes {
		return DiscoveryRequestDraft{}, false
	}

	var parsed struct {
		Purpose string `json:"purpose"`
		Effect  string `json:"effect"`
		Inputs  map[string]struct {
			Type     string `json:"type"`
			Required *bool  `json:"required"`
		} `json:"inputs"`
	}

	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return DiscoveryRequestDraft{}, false
	}

	purpose := strings.TrimSpace(parsed.Purpose)
	if purpose == "" {
		purpose = originalTask
	}
	purpose, _ = sanitizePromptString(purpose, budget.MaxTextRunes)

	effect := normalizeEffect(parsed.Effect)

	inputs := make(map[string]ContractValue)
	for name, in := range parsed.Inputs {
		trimmedName := strings.TrimSpace(name)
		// Disallow wildcards or invalid names (resolving F12)
		if trimmedName == "" || strings.ContainsAny(trimmedName, " *?") {
			continue
		}
		safeName, _ := sanitizePromptString(trimmedName, budget.MaxIdentifierRunes)
		if safeName == "" {
			continue
		}
		reqBool := true
		if in.Required != nil {
			reqBool = *in.Required
		}
		safeType, _ := sanitizePromptString(in.Type, budget.MaxIdentifierRunes)
		if safeType == "" {
			safeType = "string"
		}
		inputs[safeName] = ContractValue{
			Type:     safeType,
			Required: &reqBool,
		}
	}

	contract := StepContract{
		Purpose: purpose,
		Effect:  effect,
		Inputs:  inputs,
	}

	req := CatalogDiscoveryRequest{
		SchemaVersion: CatalogDiscoverySchemaVersion,
		Contract:      contract,
	}

	return PrepareDiscoveryRequestDraft(req, budget), true
}

func normalizeEffect(effectStr string) OperationEffect {
	s := strings.ToLower(strings.TrimSpace(effectStr))
	switch s {
	case "read":
		return OperationEffectRead
	case "write", "idempotent_write":
		// Normalizes idempotent_write to write (resolving F11)
		return OperationEffectWrite
	case "unknown":
		return OperationEffectUnknown
	default:
		return OperationEffectUnknown
	}
}

func deterministicPlanSearchQuery(goal string, budget PromptBudget) CatalogQueryPlan {
	words := strings.Fields(strings.ToLower(goal))
	stopWords := map[string]bool{
		"i": true, "want": true, "to": true, "manage": true, "the": true,
		"a": true, "an": true, "for": true, "with": true, "in": true,
		"of": true, "and": true, "or": true, "is": true, "api": true, "apis": true,
	}

	knownProviders := map[string]bool{
		"cloudflare": true, "aws": true, "github": true, "stripe": true,
		"slack": true, "google": true, "microsoft": true, "twilio": true,
		"datadog": true, "pagerduty": true, "sendgrid": true, "openai": true,
	}

	var providerKeys []string
	var keywords []string
	var effect OperationEffect = OperationEffectUnknown

	readVerbs := map[string]bool{
		"get": true, "list": true, "show": true, "fetch": true,
		"describe": true, "find": true, "query": true, "search": true, "read": true,
	}
	writeVerbs := map[string]bool{
		"create": true, "delete": true, "update": true, "set": true,
		"post": true, "put": true, "patch": true, "remove": true,
		"add": true, "write": true, "mutate": true,
	}

	for _, w := range words {
		clean := strings.Trim(w, ",.?!:;\"'()")
		if clean == "" || stopWords[clean] {
			continue
		}
		if readVerbs[clean] {
			effect = OperationEffectRead
			continue
		}
		if writeVerbs[clean] {
			effect = OperationEffectWrite
			continue
		}
		if knownProviders[clean] {
			providerKeys = append(providerKeys, clean)
			continue
		}
		if len(keywords) < 5 {
			keywords = append(keywords, clean)
		}
	}

	purpose, _ := sanitizePromptString(goal, budget.MaxTextRunes)

	return CatalogQueryPlan{
		ProviderKeys: providerKeys,
		Keywords:     keywords,
		Effect:       effect,
		Purpose:      purpose,
	}
}

func deterministicPlanDiscoveryRequest(taskDescription string, budget PromptBudget) DiscoveryRequestDraft {
	words := strings.Fields(strings.ToLower(taskDescription))
	var effect OperationEffect = OperationEffectUnknown

	readVerbs := map[string]bool{
		"get": true, "list": true, "show": true, "fetch": true,
		"describe": true, "find": true, "query": true, "search": true, "read": true,
	}
	writeVerbs := map[string]bool{
		"create": true, "delete": true, "update": true, "set": true,
		"post": true, "put": true, "patch": true, "remove": true,
		"add": true, "write": true, "mutate": true,
	}

	for _, w := range words {
		clean := strings.Trim(w, ",.?!:;\"'()")
		if readVerbs[clean] {
			effect = OperationEffectRead
			break
		}
		if writeVerbs[clean] {
			effect = OperationEffectWrite
			break
		}
	}

	purpose, _ := sanitizePromptString(taskDescription, budget.MaxTextRunes)

	// No wildcard inputs (resolving F12)
	contract := StepContract{
		Purpose: purpose,
		Effect:  effect,
		Inputs:  nil,
	}

	req := CatalogDiscoveryRequest{
		SchemaVersion: CatalogDiscoverySchemaVersion,
		Contract:      contract,
	}

	return PrepareDiscoveryRequestDraft(req, budget)
}

func extractJSON(content string) string {
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "```") {
		lines := strings.Split(content, "\n")
		if len(lines) >= 2 {
			start := 1
			end := len(lines) - 1
			for end > start && !strings.HasPrefix(strings.TrimSpace(lines[end]), "```") {
				end--
			}
			if end > start {
				content = strings.TrimSpace(strings.Join(lines[start:end], "\n"))
			}
		}
	}
	firstOpen := strings.Index(content, "{")
	lastClose := strings.LastIndex(content, "}")
	if firstOpen >= 0 && lastClose > firstOpen {
		content = content[firstOpen : lastClose+1]
	}
	return content
}
