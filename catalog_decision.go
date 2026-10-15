package apitools

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// MaxCatalogDecisionBytes is the strict 64 KiB ceiling for decision view JSON and rendered prompt.
const MaxCatalogDecisionBytes = 64 << 10

// ErrInvalidCatalogDecision indicates an invalid or unbounded catalog decision projection.
var ErrInvalidCatalogDecision = errors.New("invalid catalog decision projection")

// CatalogDecisionCandidate wraps a discovery candidate with its original 1-based index.
type CatalogDecisionCandidate struct {
	Index    int                       `json:"index"`
	Evidence CatalogDiscoveryCandidate `json:"evidence"`
}

// CatalogDecisionView projects validated discovery report evidence into a bounded decision model.
type CatalogDecisionView struct {
	Kind        string                     `json:"kind"`
	Outcome     CatalogDiscoveryOutcome    `json:"outcome"`
	Incomplete  bool                       `json:"incomplete"`
	Truncated   bool                       `json:"truncated"`
	Candidates  []CatalogDecisionCandidate `json:"candidates"`
	Selectable  []int                      `json:"selectable"`
	Omitted     int                        `json:"omitted"`
	Explanation string                     `json:"explanation,omitempty"`
	Scope       json.RawMessage            `json:"scope,omitempty"`
	Contract    json.RawMessage            `json:"contract,omitempty"`
}

// ProjectCatalogDecision preserves producer order and original 1-based candidate indexes.
// Omitted evidence never becomes selectable; candidates that exceed the independent 64 KiB
// bounds are omitted while preserving the rollback state.
func ProjectCatalogDecision(report CatalogDiscoveryReport, request CatalogDiscoveryRequest) (CatalogDecisionView, string, error) {
	v := CatalogDecisionView{
		Kind:       "catalog",
		Outcome:    report.Outcome,
		Candidates: []CatalogDecisionCandidate{},
		Selectable: []int{},
		Incomplete: report.Incomplete,
		Truncated:  report.Truncated,
	}

	scopeBytes, err := json.Marshal(report.Scope)
	if err != nil {
		return CatalogDecisionView{}, "", ErrInvalidCatalogDecision
	}
	v.Scope = scopeBytes

	if v.Outcome == CatalogDiscoveryNoQualifyingAPI {
		contractBytes, err := json.Marshal(request.Contract)
		if err != nil {
			return CatalogDecisionView{}, "", ErrInvalidCatalogDecision
		}
		v.Contract = contractBytes
		if !report.Scope.Complete || report.Incomplete || report.Truncated || report.QualifiedOperations != 0 {
			return CatalogDecisionView{}, "", ErrInvalidCatalogDecision
		}
	} else if v.Outcome == CatalogDiscoveryMatch || v.Outcome == CatalogDiscoveryAmbiguous {
		if v.Outcome == CatalogDiscoveryMatch && report.QualifiedOperations != 1 ||
			v.Outcome == CatalogDiscoveryAmbiguous && report.QualifiedOperations < 2 {
			return CatalogDecisionView{}, "", ErrInvalidCatalogDecision
		}
		v.Omitted = len(report.Candidates)
		for i, c := range report.Candidates {
			data, err := json.Marshal(c)
			if err != nil {
				return CatalogDecisionView{}, "", ErrInvalidCatalogDecision
			}
			var display CatalogDiscoveryCandidate
			if json.Unmarshal(data, &display) != nil {
				return CatalogDecisionView{}, "", ErrInvalidCatalogDecision
			}
			previous := v
			v.Candidates = append(slices.Clone(v.Candidates), CatalogDecisionCandidate{
				Index:    i + 1,
				Evidence: display,
			})
			v.Selectable = slices.Clone(v.Selectable)
			if c.Qualified && len(c.References) > 0 {
				v.Selectable = append(v.Selectable, i+1)
			}
			v.Omitted--
			if !v.bounded() {
				v = previous
				break
			}
		}
		if len(v.Candidates) == 0 {
			v.Explanation = "Candidate evidence exceeds the question limit; nothing is selectable. Reject and narrow the request."
		}
	} else {
		return CatalogDecisionView{}, "", ErrInvalidCatalogDecision
	}

	if !v.Valid() {
		return CatalogDecisionView{}, "", ErrInvalidCatalogDecision
	}
	return v, v.Prompt(), nil
}

func jsonText(value any) string {
	data, _ := json.Marshal(value)
	return string(data)
}

// Prompt renders the operator-facing decision prompt including detailed evidence.
func (v CatalogDecisionView) Prompt() string {
	if v.Outcome == CatalogDiscoveryNoQualifyingAPI {
		return "No qualifying API was found in the complete examined scope " + string(v.Scope) + " for " + string(v.Contract) + ". Choose browser, pending, or reject; this is not global absence: "
	}
	var prompt strings.Builder
	fmt.Fprintf(&prompt, "Catalog %s. Positive results may have incomplete scope. Choose one qualified shown candidate using select N, or reject. No operation runs.\n", v.Outcome)
	fmt.Fprintf(&prompt, "Examined scope: %s; incomplete evidence=%t; truncated native report=%t.\n", string(v.Scope), v.Incomplete, v.Truncated)
	for _, c := range v.Candidates {
		e := c.Evidence
		fmt.Fprintf(&prompt, "Candidate %d; qualified=%t; gaps=%s; native source evidence=%s; references=%s; native operation=%s\n",
			c.Index, e.Qualified, jsonText(e.QualificationGaps), jsonText(e.Sources), jsonText(e.References), jsonText(e.Candidate))
	}
	fmt.Fprintf(&prompt, "Omitted candidates: %d; omitted candidates cannot be selected.\n", v.Omitted)
	if v.Explanation != "" {
		prompt.WriteString(v.Explanation + "\n")
	}
	return prompt.String()
}

func (v CatalogDecisionView) bounded() bool {
	data, err := json.Marshal(v)
	return err == nil && len(data) <= MaxCatalogDecisionBytes && len(v.Prompt()) <= MaxCatalogDecisionBytes
}

// Valid confirms the view satisfies all structural and security invariants.
func (v CatalogDecisionView) Valid() bool {
	if v.Kind != "catalog" || !slices.Contains([]string{string(CatalogDiscoveryMatch), string(CatalogDiscoveryAmbiguous), string(CatalogDiscoveryNoQualifyingAPI)}, string(v.Outcome)) ||
		v.Omitted < 0 || v.Omitted+len(v.Candidates) > 100 || len(v.Selectable) > len(v.Candidates) || !v.bounded() {
		return false
	}
	var scope CatalogDiscoveryScope
	if json.Unmarshal(v.Scope, &scope) != nil {
		return false
	}
	if v.Outcome == CatalogDiscoveryNoQualifyingAPI {
		var contract StepContract
		return len(v.Candidates) == 0 && len(v.Selectable) == 0 && v.Omitted == 0 && !v.Incomplete && !v.Truncated &&
			scope.Complete && json.Unmarshal(v.Contract, &contract) == nil && contract.Purpose != "" &&
			(contract.Effect == OperationEffectRead || contract.Effect == OperationEffectWrite)
	}
	indexes := []int{}
	last := 0
	for _, c := range v.Candidates {
		if c.Index <= last || c.Index > 100 {
			return false
		}
		last = c.Index
		if c.Evidence.Qualified && len(c.Evidence.References) > 0 {
			indexes = append(indexes, c.Index)
		}
	}
	return slices.Equal(indexes, v.Selectable)
}

// Allows checks whether the given human operator answer is accepted by this decision view.
func (v CatalogDecisionView) Allows(answer string) bool {
	if !v.Valid() {
		return false
	}
	text := strings.TrimSpace(answer)
	if text == "reject" {
		return true
	}
	if v.Outcome == CatalogDiscoveryNoQualifyingAPI {
		return text == "browser" || text == "pending"
	}
	fields := strings.Fields(text)
	if len(fields) != 2 || fields[0] != "select" {
		return false
	}
	index, err := strconv.Atoi(fields[1])
	return err == nil && slices.Contains(v.Selectable, index)
}

// Clone creates an independent deep copy of the decision view.
func (v CatalogDecisionView) Clone() *CatalogDecisionView {
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var copy CatalogDecisionView
	if json.Unmarshal(data, &copy) != nil {
		return nil
	}
	return &copy
}
