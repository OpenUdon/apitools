package apitools

import (
	"encoding/json"
	"sort"
)

// Bound each projection before constructing the complete serialized report.
// Native metadata is not silently changed to fit a prompt. Omitted evidence
// makes the scope incomplete; qualification must remain explicit.
func boundCatalogDiscoveryReport(report CatalogDiscoveryReport) CatalogDiscoveryReport {
	budget := report.Limits.MaxContextBytes
	if budget < MinCatalogDiscoveryContextBytes || budget > MaxCatalogDiscoveryContextBytes {
		budget = DefaultCatalogDiscoveryContextBytes
	}
	maxResults := report.Limits.MaxResults
	if maxResults <= 0 || maxResults > MaxCatalogDiscoveryResults {
		maxResults = DefaultCatalogDiscoveryResults
	}
	candidates, leads, coverage, exclusions, diagnostics, ties := report.Candidates, report.Leads, report.Coverage, report.Exclusions, report.Diagnostics, report.Ties
	report.Candidates = nil
	report.Leads = nil
	report.Coverage = nil
	report.Exclusions = nil
	report.Diagnostics = nil
	report.Ties = nil
	data, _ := json.Marshal(report)
	used := len(data) + 256 // reserve a generic loss diagnostic and JSON field framing
	loss := false
	if used > budget {
		report.Scope.ProviderIDs = nil
		if report.Outcome != CatalogDiscoveryBlocked {
			report.Outcome = CatalogDiscoveryInsufficientEvidence
		}
		loss = true
		data, _ = json.Marshal(report)
		used = len(data) + 256
	}
	fit := func(item any) bool {
		data, err := json.Marshal(item)
		if err != nil || used+len(data)+64 > budget {
			loss = true
			return false
		}
		used += len(data) + 64
		return true
	}
	sortDiagnostics(diagnostics)
	for _, item := range diagnostics {
		if fit(item) {
			report.Diagnostics = append(report.Diagnostics, item)
		}
	}
	for i, item := range candidates {
		if i >= maxResults {
			report.Truncated = true
			break
		}
		if len(item.References) > 32 || len(item.Sources) > 32 {
			loss = true
			continue
		}
		if fit(item) {
			report.Candidates = append(report.Candidates, item)
		}
	}
	sort.Slice(leads, func(i, j int) bool {
		if leads[i].ProviderID != leads[j].ProviderID {
			return leads[i].ProviderID < leads[j].ProviderID
		}
		return leads[i].SpecRefID < leads[j].SpecRefID
	})
	for _, item := range leads {
		if fit(item) {
			report.Leads = append(report.Leads, item)
		}
	}
	sortCatalogCoverage(coverage)
	for _, item := range coverage {
		if fit(item) {
			report.Coverage = append(report.Coverage, item)
		}
	}
	sort.Slice(exclusions, func(i, j int) bool {
		if exclusions[i].ProviderID != exclusions[j].ProviderID {
			return exclusions[i].ProviderID < exclusions[j].ProviderID
		}
		if exclusions[i].SpecRefID != exclusions[j].SpecRefID {
			return exclusions[i].SpecRefID < exclusions[j].SpecRefID
		}
		return exclusions[i].ArtifactID < exclusions[j].ArtifactID
	})
	for _, item := range exclusions {
		if fit(item) {
			report.Exclusions = append(report.Exclusions, item)
		}
	}
	if !loss && len(report.Candidates) == len(candidates) {
		for _, item := range ties {
			if fit(item) {
				report.Ties = append(report.Ties, item)
			}
		}
	}
	if loss {
		report.Truncated = true
		report.Incomplete = true
		report.Scope.Complete = false
		if report.Outcome == CatalogDiscoveryNoQualifyingAPI || (report.Outcome == CatalogDiscoveryMatch && len(report.Candidates) == 0) {
			report.Outcome = CatalogDiscoveryInsufficientEvidence
		}
		report.Diagnostics = append(report.Diagnostics, Diagnostic{Severity: "warning", Code: "discovery.context_limit", Message: "report limits omitted evidence; scope remains incomplete"})
	}
	return report
}
