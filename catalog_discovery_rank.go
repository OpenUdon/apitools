package apitools

import (
	"context"
	"sort"
	"strings"
)

// Rank metadata already checked by the index reader. Scores order evidence;
// only the dimension statuses establish qualification.
func rankCatalogDiscovery(ctx context.Context, evidence catalogDiscoveryLocalEvidence, request CatalogDiscoveryRequest) (CatalogDiscoveryReport, error) {
	report := evidence.report
	budget := resolvedPromptBudget(PromptBudget{})
	budgetLeads := map[string]bool{}
	operationLimited := false
	contract, inputs, outputs, _ := prepareStepContract(request.Contract, budget)
	unresolved := report.Incomplete || len(evidence.candidates) == 0
	for _, item := range evidence.candidates {
		if err := ctx.Err(); err != nil {
			return catalogDiscoveryInterrupted(report, err).report, err
		}
		candidate := &item.Candidate
		metadataLoss := len(item.QualificationGaps) > 0
		candidate.Match = compareCandidateToContract(*candidate, contract, inputs, outputs)
		purpose := &candidate.Match.Purpose
		query := purposeTokens(contract.contract.Purpose)
		available := map[string]bool{}
		for _, term := range purposeTokens(candidate.Summary.Description) {
			available[term] = true
		}
		shared := 0
		for _, term := range query {
			if available[term] {
				shared++
			}
		}
		if purpose.Status == ContractMatchCompatible && (shared < 2 || shared*2 < len(query)) {
			purpose.Status = ContractMatchIndeterminate
			purpose.Missing = append(purpose.Missing, "Qualification requires at least two distinct documented purpose terms and half of the requested terms; clarify the intended operation.")
		}
		if purpose.Status != ContractMatchCompatible {
			item.QualificationGaps = append(item.QualificationGaps, "documented purpose compatibility is insufficient")
		}
		for _, gap := range candidate.Match.Inputs.Gaps {
			if candidate.Match.Inputs.Status == ContractMatchCompatible && strings.Contains(gap, "is not used by the operation") {
				candidate.Match.Inputs.Status = ContractMatchIndeterminate
			}
		}
		if candidate.Match.Inputs.Status != ContractMatchCompatible {
			item.QualificationGaps = append(item.QualificationGaps, "input compatibility is not established")
		}
		if candidate.Match.Outputs.Status != ContractMatchCompatible {
			item.QualificationGaps = append(item.QualificationGaps, "output compatibility is not established")
		}
		if (request.Contract.Effect == OperationEffectRead || request.Contract.Effect == OperationEffectWrite) && candidate.Match.Effect.Status != ContractMatchCompatible {
			item.QualificationGaps = append(item.QualificationGaps, "requested effect compatibility is not established")
		}
		if candidateCapability(*candidate, "auth").Status != OperationCapabilitySupported {
			item.QualificationGaps = append(item.QualificationGaps, "authentication alternatives are incomplete")
		}
		for _, issue := range candidate.Operation.ReadinessIssues {
			if issue.Code == "prompt.operation_budget" {
				metadataLoss = true
				item.QualificationGaps = append(item.QualificationGaps, "operation metadata was compacted; review the source")
			}
		}
		for _, gap := range candidate.Summary.Gaps {
			if strings.Contains(gap, "omitted") || strings.Contains(gap, "no usable") {
				metadataLoss = true
				item.QualificationGaps = append(item.QualificationGaps, "source fields were omitted or lack usable identity; review the source")
				break
			}
		}
		item.QualificationGaps = uniqueSortedStrings(item.QualificationGaps)
		if size, err := marshalSize(*candidate); err != nil || size > budget.MaxOperationBytes {
			unresolved, operationLimited = true, true
			report.Incomplete, report.Truncated, report.Scope.Complete = true, true, false
			for _, source := range item.Sources {
				key := source.ProviderID + "\x00" + source.SpecRefID + "\x00" + source.ArtifactID + "\x00" + candidate.Source.SHA256
				if budgetLeads[key] {
					continue
				}
				budgetLeads[key] = true
				lead := CatalogDiscoveryLead{ProviderID: source.ProviderID, SpecRefID: source.SpecRefID, Kind: string(candidate.Source.Kind), Evidence: source, Remote: item.Remote, Reason: "ranked operation exceeds the 32-KiB prompt operation budget; narrow the contract or review the exact source"}
				for _, ref := range item.References {
					if ref.ProviderID == source.ProviderID && ref.SpecRefID == source.SpecRefID && ref.ArtifactID == source.ArtifactID {
						ref.Selector = "" // reference-only source, never an operation selection
						lead.Reference = &ref
						break
					}
				}
				report.Leads = append(report.Leads, lead)
			}
			continue
		}
		item.Qualified = len(item.QualificationGaps) == 0
		if metadataLoss {
			unresolved = true
			report.Incomplete, report.Scope.Complete = true, false
		}
		if item.Qualified {
			report.QualifiedOperations++
		} else if !hasIncompatibleDimension(candidate.Match) {
			unresolved = true
		}
		report.Candidates = append(report.Candidates, item)
	}
	if operationLimited {
		report.Diagnostics = append(report.Diagnostics, Diagnostic{Severity: "warning", Code: "discovery.operation_budget", Message: "ranked operations exceed the existing prompt operation budget; source evidence remains unexamined"})
	}
	sort.Slice(report.Candidates, func(i, j int) bool {
		left, right := report.Candidates[i], report.Candidates[j]
		if left.Candidate.Match.Score != right.Candidate.Match.Score {
			return left.Candidate.Match.Score > right.Candidate.Match.Score
		}
		return catalogDiscoveryCandidateKey(left) < catalogDiscoveryCandidateKey(right)
	})
	for start := 0; start < len(report.Candidates); {
		end := start + 1
		for end < len(report.Candidates) && report.Candidates[end].Candidate.Match.Score == report.Candidates[start].Candidate.Match.Score {
			end++
		}
		if end-start > 1 {
			tie := CatalogDiscoveryTie{Score: report.Candidates[start].Candidate.Match.Score}
			for i := start; i < end; i++ {
				tie.CandidateIndexes = append(tie.CandidateIndexes, i)
			}
			report.Ties = append(report.Ties, tie)
		}
		start = end
	}
	switch {
	case report.QualifiedOperations > 1:
		report.Outcome = CatalogDiscoveryAmbiguous
	case report.QualifiedOperations == 1:
		report.Outcome = CatalogDiscoveryMatch
	case !unresolved:
		report.Outcome = CatalogDiscoveryNoQualifyingAPI
	default:
		report.Outcome = CatalogDiscoveryInsufficientEvidence
	}
	return report, nil
}
