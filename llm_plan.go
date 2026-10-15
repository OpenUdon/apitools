package apitools

import (
	"context"
	"encoding/json"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/OpenUdon/apitools/internal/sourceguard"
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
	raw := planningJSON(content)
	if raw == "" {
		return CatalogQueryPlan{}, false
	}

	var parsed struct {
		ProviderKeys []*string `json:"provider_keys"`
		Keywords     []*string `json:"keywords"`
		Effect       *string   `json:"effect"`
		Purpose      *string   `json:"purpose"`
	}

	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return CatalogQueryPlan{}, false
	}

	if parsed.Purpose == nil || parsed.Effect == nil || !supportedPlanEffect(parsed.Effect) || parsed.ProviderKeys == nil || parsed.Keywords == nil ||
		strings.TrimSpace(*parsed.Purpose) == "" || utf8.RuneCountInString(*parsed.Purpose) > budget.MaxTextRunes ||
		len(parsed.ProviderKeys) > budget.MaxCollectionItems || len(parsed.Keywords) > budget.MaxCollectionItems {
		return CatalogQueryPlan{}, false
	}
	for _, values := range [][]*string{parsed.ProviderKeys, parsed.Keywords} {
		for _, v := range values {
			if v == nil || utf8.RuneCountInString(*v) > budget.MaxIdentifierRunes {
				return CatalogQueryPlan{}, false
			}
		}
	}
	effect := normalizeEffect(*parsed.Effect)
	purpose, _ := sanitizePromptString(*parsed.Purpose, budget.MaxTextRunes)
	if strings.TrimSpace(purpose) == "" {
		return CatalogQueryPlan{}, false
	}

	var providerKeys []string
	for _, key := range parsed.ProviderKeys {
		k := strings.ToLower(strings.TrimSpace(*key))
		if k != "" && !strings.ContainsAny(k, " *?") {
			if safe, _ := sanitizePromptString(k, budget.MaxIdentifierRunes); safe != "" {
				providerKeys = append(providerKeys, safe)
			}
		}
	}

	var keywords []string
	for _, kw := range parsed.Keywords {
		w := strings.ToLower(strings.TrimSpace(*kw))
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
	raw := planningJSON(content)
	if raw == "" {
		return DiscoveryRequestDraft{}, false
	}

	var parsed struct {
		Purpose *string `json:"purpose"`
		Effect  *string `json:"effect"`
		Inputs  map[string]*struct {
			Type     *string         `json:"type"`
			Required json.RawMessage `json:"required"`
		} `json:"inputs"`
	}

	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return DiscoveryRequestDraft{}, false
	}

	if parsed.Purpose == nil || parsed.Effect == nil || !supportedPlanEffect(parsed.Effect) || parsed.Inputs == nil ||
		strings.TrimSpace(*parsed.Purpose) == "" || utf8.RuneCountInString(*parsed.Purpose) > budget.MaxTextRunes ||
		len(parsed.Inputs) > budget.MaxFields {
		return DiscoveryRequestDraft{}, false
	}
	purpose, _ := sanitizePromptString(*parsed.Purpose, budget.MaxTextRunes)
	if strings.TrimSpace(purpose) == "" {
		return DiscoveryRequestDraft{}, false
	}
	effect := normalizeEffect(*parsed.Effect)

	inputs := make(map[string]ContractValue)
	for name, in := range parsed.Inputs {
		if in == nil || in.Type == nil || strings.TrimSpace(*in.Type) == "" || utf8.RuneCountInString(name) > budget.MaxIdentifierRunes || utf8.RuneCountInString(*in.Type) > budget.MaxIdentifierRunes {
			return DiscoveryRequestDraft{}, false
		}
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
		if len(in.Required) != 0 {
			var required *bool
			if err := json.Unmarshal(in.Required, &required); err != nil || required == nil {
				return DiscoveryRequestDraft{}, false
			}
			reqBool = *required
		}
		safeType, _ := sanitizePromptString(*in.Type, budget.MaxIdentifierRunes)
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

	draft := PrepareDiscoveryRequestDraft(req, budget)
	return draft, draft.Ready
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

// planningJSON accepts a single JSON object, optionally inside a JSON fence.
// Bound the complete response before allocating fence lines or decoding fields.
func planningJSON(content string) string {
	if len(content) == 0 || len(content) > MaxCatalogDecisionBytes {
		return ""
	}
	raw := strings.TrimSpace(content)
	if strings.HasPrefix(raw, "```") {
		first := strings.IndexByte(raw, '\n')
		if first < 0 || !strings.HasSuffix(raw, "```") {
			return ""
		}
		fence := strings.TrimSpace(raw[:first])
		if fence != "```" && fence != "```json" {
			return ""
		}
		raw = strings.TrimSpace(raw[first+1 : len(raw)-3])
	}
	if !strings.HasPrefix(raw, "{") || !strings.HasSuffix(raw, "}") {
		return ""
	}
	if err := sourceguard.CheckJSON("model plan", []byte(raw)); err != nil {
		return ""
	}
	return raw
}

func supportedPlanEffect(effect *string) bool {
	if effect == nil {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(*effect)) {
	case "read", "write", "idempotent_write", "unknown":
		return true
	default:
		return false
	}
}
