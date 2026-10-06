package apitools

import (
	"sort"
	"strings"
)

// compareVersions consumes only caller-provided baseline metadata and the
// newly fetched bytes. It never rereads or reconstructs the held source.
func (s *apiVersionSession) compareVersions() {
	baseline := s.opts.BaselineInventory
	if baseline != nil {
		sanitizeInventory(baseline, DefaultPromptBudget())
	}
	baselineOK := baseline != nil && strings.EqualFold(s.opts.BaselineInventorySHA256, s.req.Known.SHA256) && versionInventoryComplete(*baseline)
	for i := range s.report.Versions {
		row := &s.report.Versions[i]
		content, ok := s.contents[row.ID]
		if !ok || row.Validation != "valid" {
			continue
		}
		if row.SourceKind == OperationSourceOpenAPI && s.req.Known.SourceKind == OperationSourceOpenAPI && baselineOK {
			inventory, err := BuildOperationInventory(s.ctx, InventoryOptions{Documents: []InventoryDocument{{URL: row.FinalURL, Content: content}}, MaxBytes: apiVersionMaxSourceBytes})
			if err == nil && versionInventoryComplete(inventory) {
				row.Diff = diffVersionOperations(*baseline, inventory)
			} else {
				row.Diagnostics = append(row.Diagnostics, versionDiagnostic("version.diff_unexamined", "Fetched inventory is incomplete; absence cannot be established."))
			}
		} else {
			row.Diagnostics = append(row.Diagnostics, versionDiagnostic("version.diff_unexamined", "Supply complete, digest-matched baseline inventory for source-native comparison."))
		}
		if s.req.Contract != nil {
			ranked, err := BuildOperationCandidates(s.ctx, OperationCandidateOptions{Sources: []OperationSourceInput{{Kind: row.SourceKind, URL: row.FinalURL, Content: content}}, Contract: *s.req.Contract, MaxBytes: apiVersionMaxSourceBytes})
			row.Candidates = ranked.Candidates
			row.Diagnostics = append(row.Diagnostics, ranked.Diagnostics...)
			if err != nil || ranked.Truncated {
				row.Diagnostics = append(row.Diagnostics, versionDiagnostic("version.capability_unexamined", "Capability evidence is incomplete; source review remains necessary."))
			}
		}
	}
}

func versionInventoryComplete(inventory OperationInventory) bool {
	if inventory.Truncated {
		return false
	}
	for _, d := range inventory.Diagnostics {
		if d.Severity == "error" || strings.Contains(d.Code, "sanit") || strings.Contains(d.Code, "compact") {
			return false
		}
	}
	return true
}

func diffVersionOperations(a, b OperationInventory) *APIVersionOperationDiff {
	left, right := map[string]bool{}, map[string]bool{}
	idsLeft, idsRight := map[string]string{}, map[string]string{}
	for _, op := range a.Operations {
		left[strings.ToUpper(op.Method)+" "+op.Path] = true
		idsLeft[strings.ToUpper(op.Method)+" "+op.Path] = op.OperationID
	}
	for _, op := range b.Operations {
		right[strings.ToUpper(op.Method)+" "+op.Path] = true
		idsRight[strings.ToUpper(op.Method)+" "+op.Path] = op.OperationID
	}
	diff := &APIVersionOperationDiff{}
	for key := range left {
		if right[key] {
			diff.Retained = append(diff.Retained, key)
			if idsLeft[key] != idsRight[key] {
				diff.ChangedOperationIDs = append(diff.ChangedOperationIDs, key+": "+idsLeft[key]+" -> "+idsRight[key])
			}
		} else {
			diff.Removed = append(diff.Removed, key)
		}
	}
	for key := range right {
		if !left[key] {
			diff.Added = append(diff.Added, key)
		}
	}
	sort.Strings(diff.Added)
	sort.Strings(diff.Removed)
	sort.Strings(diff.Retained)
	sort.Strings(diff.ChangedOperationIDs)
	return diff
}
