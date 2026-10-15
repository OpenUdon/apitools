package apitools

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"
)

// DecisionRecommendation represents an advisory candidate recommendation and rationale.
type DecisionRecommendation struct {
	RecommendedIndex int    `json:"recommended_index,omitempty"`
	Rationale        string `json:"rationale,omitempty"`
	Caveats          string `json:"caveats,omitempty"`
}

// RecommendCandidate evaluates candidate operations in a validated CatalogDecisionView.
// Only indexes present in view.Selectable are eligible for recommendation (resolving F13).
// If view.Selectable is empty, it abstains from recommending (returns RecommendedIndex: 0).
// Advice never selects or executes. It updates view.Explanation while preserving strict
// view bounds and validation invariants.
func RecommendCandidate(ctx context.Context, llm LLMClient, goal string, view CatalogDecisionView) (DecisionRecommendation, CatalogDecisionView, error) {
	if ctx.Err() != nil {
		return DecisionRecommendation{}, view, ctx.Err()
	}

	// Empty eligibility yields no recommendation
	if len(view.Selectable) == 0 {
		rec := DecisionRecommendation{
			RecommendedIndex: 0,
			Rationale:        "No selectable candidates available in the current decision view.",
		}
		return rec, view, nil
	}

	budget := DefaultPromptBudget()
	sanitizedGoal, _ := sanitizePromptString(goal, budget.MaxTextRunes)

	// Filter down to only selectable candidates
	selectableCandidates := make([]CatalogDecisionCandidate, 0, len(view.Selectable))
	for _, c := range view.Candidates {
		if slices.Contains(view.Selectable, c.Index) {
			selectableCandidates = append(selectableCandidates, c)
		}
	}

	if len(selectableCandidates) == 0 {
		rec := DecisionRecommendation{
			RecommendedIndex: 0,
			Rationale:        "No selectable candidates available in the current decision view.",
		}
		return rec, view, nil
	}

	var rec DecisionRecommendation
	if llm == nil {
		rec = deterministicRecommendCandidate(selectableCandidates, sanitizedGoal)
	} else {
		systemPrompt := `You are an API architecture advisor. Given a user goal and a list of selectable candidate API operations, recommend the single best candidate.
Explain why it was chosen in 1-2 concise sentences, and mention any important required parameters or caveats.
You MUST choose ONLY from the provided candidate indexes.

Respond ONLY with valid JSON:
{
  "recommended_index": <integer>,
  "rationale": "<concise explanation>",
  "caveats": "<caveats or empty string>"
}`

		var userPrompt strings.Builder
		if sanitizedGoal != "" {
			fmt.Fprintf(&userPrompt, "Goal: %s\n\n", sanitizedGoal)
		}
		userPrompt.WriteString("Selectable Candidates:\n")
		for _, c := range selectableCandidates {
			op := c.Evidence.Candidate.Operation
			fmt.Fprintf(&userPrompt, "- Candidate #%d: %s %s (operationId: %s)\n", c.Index, op.Method, op.Path, op.OperationID)
			if op.Summary != "" {
				fmt.Fprintf(&userPrompt, "  Summary: %s\n", op.Summary)
			}
			if len(c.Evidence.QualificationGaps) > 0 {
				fmt.Fprintf(&userPrompt, "  Gaps: %s\n", strings.Join(c.Evidence.QualificationGaps, "; "))
			}
		}

		req := CompletionRequest{
			SystemPrompt: systemPrompt,
			UserPrompt:   userPrompt.String(),
			Timeout:      DefaultLLMTimeout,
		}

		resp, err := llm.Complete(ctx, req)
		if err != nil {
			if ctx.Err() != nil {
				return DecisionRecommendation{}, view, ctx.Err()
			}
			rec = deterministicRecommendCandidate(selectableCandidates, sanitizedGoal)
		} else {
			parsedRec, ok := parseRecommendationResponse(resp.Content, view.Selectable, budget)
			if !ok {
				rec = deterministicRecommendCandidate(selectableCandidates, sanitizedGoal)
			} else {
				rec = parsedRec
			}
		}
	}

	// Attach explanation to view and recheck bounds
	updatedView := view
	if rec.RecommendedIndex > 0 {
		var explBuilder strings.Builder
		fmt.Fprintf(&explBuilder, "Recommended [#%d]: %s", rec.RecommendedIndex, rec.Rationale)
		if rec.Caveats != "" {
			fmt.Fprintf(&explBuilder, " Note: %s", rec.Caveats)
		}
		newExpl := explBuilder.String()
		safeExpl, _ := sanitizePromptString(newExpl, budget.MaxTextRunes)

		updatedView.Explanation = safeExpl
		if !updatedView.bounded() || !updatedView.Valid() {
			// If adding explanation breached bounds, revert to prior explanation
			updatedView.Explanation = view.Explanation
		}
	}

	return rec, updatedView, nil
}

// RecommendCandidate on *Client delegates to RecommendCandidate.
func (c *Client) RecommendCandidate(ctx context.Context, llm LLMClient, goal string, view CatalogDecisionView) (DecisionRecommendation, CatalogDecisionView, error) {
	return RecommendCandidate(ctx, llm, goal, view)
}

func deterministicRecommendCandidate(candidates []CatalogDecisionCandidate, goal string) DecisionRecommendation {
	if len(candidates) == 0 {
		return DecisionRecommendation{
			RecommendedIndex: 0,
			Rationale:        "No selectable candidates available.",
		}
	}

	// Deterministic selection: takes the first candidate (which preserves existing ranking order)
	first := candidates[0]
	op := first.Evidence.Candidate.Operation

	rationale := fmt.Sprintf("Candidate [#%d] %s %s selected based on schema compatibility and qualification ranking.",
		first.Index, op.Method, op.Path)
	var caveats string
	if len(first.Evidence.QualificationGaps) > 0 {
		caveats = strings.Join(first.Evidence.QualificationGaps, "; ")
	}

	return DecisionRecommendation{
		RecommendedIndex: first.Index,
		Rationale:        rationale,
		Caveats:          caveats,
	}
}

func parseRecommendationResponse(content string, selectable []int, budget PromptBudget) (DecisionRecommendation, bool) {
	raw := planningJSON(content)
	if raw == "" {
		return DecisionRecommendation{}, false
	}

	var parsed struct {
		RecommendedIndex *int    `json:"recommended_index"`
		Rationale        *string `json:"rationale"`
		Caveats          string  `json:"caveats"`
	}

	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return DecisionRecommendation{}, false
	}

	// Must be in selectable set (resolving F13: ranking is not selection authority; cannot recommend outside selectable set)
	if parsed.RecommendedIndex == nil || parsed.Rationale == nil || strings.TrimSpace(*parsed.Rationale) == "" || utf8.RuneCountInString(*parsed.Rationale) > budget.MaxTextRunes || utf8.RuneCountInString(parsed.Caveats) > budget.MaxTextRunes || !slices.Contains(selectable, *parsed.RecommendedIndex) {
		return DecisionRecommendation{}, false
	}

	rationale, _ := sanitizePromptString(*parsed.Rationale, budget.MaxTextRunes)
	caveats, _ := sanitizePromptString(parsed.Caveats, budget.MaxTextRunes)

	return DecisionRecommendation{
		RecommendedIndex: *parsed.RecommendedIndex,
		Rationale:        rationale,
		Caveats:          caveats,
	}, true
}
