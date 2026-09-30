package apitools

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/OpenUdon/apitools/catalog"
)

type catalogDiscoveryLocalEvidence struct {
	report     CatalogDiscoveryReport
	candidates []CatalogDiscoveryCandidate
}

// DiscoverCatalogOperations retrieves index-backed metadata for the additive
// contract. Source confirmation, provisioning and execution remain downstream.
func DiscoverCatalogOperations(ctx context.Context, options CatalogDiscoveryOptions) (CatalogDiscoveryReport, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	parent := ctx
	timeout := options.Request.Limits.TimeoutMillis
	if timeout <= 0 || timeout > int(MaxCatalogDiscoveryTimeout/time.Millisecond) {
		timeout = int(DefaultCatalogDiscoveryTimeout / time.Millisecond)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
	defer cancel()
	evidence, err := retrieveCatalogDiscovery(ctx, options)
	report := evidence.report
	if err != nil || report.Outcome == CatalogDiscoveryBlocked {
		if parent.Err() == nil && ctx.Err() != nil {
			err = nil
		}
		return boundCatalogDiscoveryReport(report), err
	}
	if options.Request.RemoteLookup {
		evidence = retrieveCatalogDiscoveryRemote(ctx, evidence, options)
		if evidence.report.Outcome == CatalogDiscoveryBlocked {
			return boundCatalogDiscoveryReport(evidence.report), nil
		}
	}
	report, err = rankCatalogDiscovery(ctx, evidence, options.Request)
	if err != nil && parent.Err() == nil {
		err = nil // installation deadline is incomplete evidence, not caller cancellation
	}
	report = boundCatalogDiscoveryReport(report)
	if parent.Err() != nil {
		return boundCatalogDiscoveryReport(catalogDiscoveryInterrupted(report, parent.Err()).report), parent.Err()
	}
	if ctx.Err() != nil {
		return boundCatalogDiscoveryReport(catalogDiscoveryInterrupted(report, ctx.Err()).report), nil
	}
	return report, err
}

func retrieveCatalogDiscovery(ctx context.Context, options CatalogDiscoveryOptions) (catalogDiscoveryLocalEvidence, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	parent := ctx
	report := CatalogDiscoveryReport{SchemaVersion: CatalogDiscoverySchemaVersion, Outcome: CatalogDiscoveryInsufficientEvidence, Scope: CatalogDiscoveryScope{ProviderIDs: []string{}, ProviderConstrained: options.Request.ProviderKeys != nil}}
	result := catalogDiscoveryLocalEvidence{report: report}
	cat := catalog.BuiltInCatalog()
	if options.Index.Catalog != nil {
		cat = *options.Index.Catalog
	}
	if err := cat.Validate(); err != nil {
		return catalogDiscoveryRefusal(report, "discovery.catalog_invalid", "catalog installation metadata is invalid"), nil
	}
	limits, ids, diagnostics := prepareCatalogDiscoveryRequest(options.Request, cat)
	report.Limits = limits
	report.Scope.ProviderIDs = ids
	report.Diagnostics = diagnostics
	if len(errorDiagnostics(diagnostics)) > 0 {
		return catalogDiscoveryLocalEvidence{report: reportWithDiscoveryBlock(report)}, nil
	}
	digest, err := catalogIndexIdentity(cat, options.Index.Catalog == nil)
	if err != nil {
		return catalogDiscoveryRefusal(report, "discovery.catalog_identity", "catalog identity could not be established"), nil
	}
	report.CatalogSHA256 = digest
	if err := parent.Err(); err != nil {
		return catalogDiscoveryInterrupted(report, err), err
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(limits.TimeoutMillis)*time.Millisecond)
	defer cancel()
	if options.Request.RemoteLookup && !options.RemoteEnabled {
		return catalogDiscoveryRefusal(report, "discovery.remote_unavailable", "remote lookup requires configured supported installation capability"), nil
	}
	paths, err := catalog.ResolveRoot(options.Index.Root)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			report = missingCatalogDiscoveryEvidence(report, cat, options.Request.Filters, "discovery.root_missing", "configure an existing prepared catalog root")
			return catalogDiscoveryLocalEvidence{report: report}, nil
		}
		return catalogDiscoveryRefusal(report, "discovery.root_invalid", "catalog root configuration is unsafe or invalid"), nil
	}
	if len(ids) == 0 {
		report.Incomplete = true
		report.Diagnostics = append(report.Diagnostics, Diagnostic{Severity: "warning", Code: "discovery.scope_empty", Message: "the selected provider scope is empty"})
		return catalogDiscoveryLocalEvidence{report: report}, nil
	}
	if paths.Directory == "" {
		report = missingCatalogDiscoveryEvidence(report, cat, options.Request.Filters, "discovery.root_missing", "configure an explicit prepared catalog root and index")
		return catalogDiscoveryLocalEvidence{report: report}, nil
	}
	index, err := ReadCatalogOperationIndex(ctx, options.Index)
	if err != nil {
		if parent.Err() != nil {
			return catalogDiscoveryInterrupted(report, parent.Err()), parent.Err()
		}
		if ctx.Err() != nil {
			return catalogDiscoveryInterrupted(report, ctx.Err()), nil
		}
		if errors.Is(err, os.ErrNotExist) {
			report = missingCatalogDiscoveryEvidence(report, cat, options.Request.Filters, "discovery.index_missing", "prepare registrations and explicitly build the local catalog index")
			return catalogDiscoveryLocalEvidence{report: report}, nil
		}
		return catalogDiscoveryRefusal(report, "discovery.index_invalid", "catalog operation index or registration identity is invalid"), nil
	}
	scope := map[string]bool{}
	for _, id := range ids {
		scope[id] = true
	}
	allowed := map[string]bool{}
	workLimited := map[string]bool{}
	providerByID := map[string]catalog.Provider{}
	for _, provider := range cat.ListProviders() {
		providerByID[provider.ID] = provider
	}
	for _, artifact := range index.Artifacts {
		if err := ctx.Err(); err != nil {
			if parent.Err() != nil {
				return catalogDiscoveryInterrupted(report, parent.Err()), parent.Err()
			}
			return catalogDiscoveryInterrupted(report, err), nil
		}
		var refs []CatalogArtifactReference
		var sources []CatalogDiscoverySourceEvidence
		for _, link := range artifact.Links {
			if !scope[link.ProviderID] {
				continue
			}
			evidence := catalogDiscoverySourceEvidence(providerByID[link.ProviderID], link.SpecRefID, link.ArtifactID, link.Advisory)
			if len(evidence.EvidenceGaps) > 0 {
				report.Incomplete = true
			}
			key := link.ProviderID + "\x00" + link.SpecRefID + "\x00" + link.ArtifactID
			if reason := catalogDiscoveryFilterExclusion(evidence, options.Request.Filters); reason != "" {
				report.Exclusions = append(report.Exclusions, CatalogDiscoveryExclusion{ProviderID: link.ProviderID, SpecRefID: link.SpecRefID, ArtifactID: link.ArtifactID, Reason: reason})
				continue
			}
			allowed[key] = true
			sources = append(sources, evidence)
			ref := CatalogArtifactReference{CatalogSHA256: digest, ProviderID: link.ProviderID, SpecRefID: link.SpecRefID, ArtifactID: link.ArtifactID, Kind: artifact.Kind, SHA256: artifact.SHA256, Bytes: artifact.Bytes}
			refs = append(refs, ref)
			if len(artifact.Operations) == 0 || artifact.State != "indexed" {
				report.Leads = append(report.Leads, CatalogDiscoveryLead{ProviderID: link.ProviderID, DisplayName: providerByID[link.ProviderID].DisplayName, SpecRefID: link.SpecRefID, Kind: string(artifact.Kind), Reason: "registered source coverage is unexamined", Evidence: evidence, Reference: &ref})
			}
		}
		if len(refs) == 0 {
			continue
		}
		sanitized := map[string]bool{}
		for _, diagnostic := range artifact.Diagnostics {
			if diagnostic.Code == "prompt.operation_sanitized" {
				sanitized[diagnostic.Path] = true
			}
		}
		limited := false
		for _, candidate := range artifact.Operations {
			if err := ctx.Err(); err != nil {
				if parent.Err() != nil {
					return catalogDiscoveryInterrupted(report, parent.Err()), parent.Err()
				}
				return catalogDiscoveryInterrupted(report, err), nil
			}
			if report.ExaminedOperations >= limits.MaxOperations {
				report.Incomplete = true
				report.Truncated = true
				if !limited {
					for _, ref := range refs {
						workLimited[ref.ProviderID+"\x00"+ref.SpecRefID+"\x00"+ref.ArtifactID] = true
					}
					limited = true
				}
				continue
			}
			selected := append([]CatalogArtifactReference(nil), refs...)
			for i := range selected {
				selected[i].Selector = candidate.Source.Selector
			}
			item := CatalogDiscoveryCandidate{Candidate: candidate, References: selected, Sources: sources, ProviderConstrained: report.Scope.ProviderConstrained}
			if sanitized[candidate.Operation.Provenance] || sanitized[candidate.Source.Selector] {
				item.QualificationGaps = append(item.QualificationGaps, "source operation metadata was sanitized; review any identity or selected-field loss")
			}
			result.candidates = append(result.candidates, item)
			report.ExaminedOperations++
		}
	}
	for _, coverage := range index.Coverage {
		if !scope[coverage.ProviderID] {
			continue
		}
		key := coverage.ProviderID + "\x00" + coverage.SpecRefID + "\x00" + coverage.ArtifactID
		if coverage.ArtifactID != "" && !allowed[key] {
			continue
		}
		if coverage.ArtifactID == "" {
			evidence := catalogDiscoverySourceEvidence(providerByID[coverage.ProviderID], coverage.SpecRefID, "", false)
			if reason := catalogDiscoveryFilterExclusion(evidence, options.Request.Filters); reason != "" {
				report.Exclusions = append(report.Exclusions, CatalogDiscoveryExclusion{ProviderID: coverage.ProviderID, SpecRefID: coverage.SpecRefID, Reason: reason})
				continue
			}
			report.Leads = append(report.Leads, CatalogDiscoveryLead{ProviderID: coverage.ProviderID, DisplayName: providerByID[coverage.ProviderID].DisplayName, SpecRefID: coverage.SpecRefID, Reason: coverage.Reason, Evidence: evidence})
		}
		if workLimited[key] {
			coverage.State = "work_limit"
			coverage.Reason = "operation work limit left this registered scope unexamined"
		}
		report.Coverage = append(report.Coverage, coverage)
		if coverage.State != "indexed" {
			report.Incomplete = true
		}
	}
	if report.ExaminedOperations == 0 {
		report.Incomplete = true
		report.Diagnostics = append(report.Diagnostics, Diagnostic{Severity: "warning", Code: "discovery.operations_unexamined", Message: "no selected source-backed operations were examined"})
	}
	if report.Truncated {
		report.Diagnostics = append(report.Diagnostics, Diagnostic{Severity: "warning", Code: "discovery.work_limit", Message: "the operation work limit left relevant scope unexamined"})
	}
	report.Scope.Complete = !report.Incomplete
	sortCatalogCoverage(report.Coverage)
	sort.Slice(result.candidates, func(i, j int) bool {
		return catalogDiscoveryCandidateKey(result.candidates[i]) < catalogDiscoveryCandidateKey(result.candidates[j])
	})
	result.report = report
	return result, nil
}

func catalogDiscoverySourceEvidence(provider catalog.Provider, specID, artifactID string, advisory bool) CatalogDiscoverySourceEvidence {
	evidence := CatalogDiscoverySourceEvidence{ProviderID: provider.ID, SpecRefID: specID, ArtifactID: artifactID, Authority: "unknown", Advisory: advisory, LicenseIdentifier: "unknown", Redistribution: "unknown"}
	for _, ref := range provider.SpecReferences {
		if ref.ID != specID {
			continue
		}
		evidence.LicenseNote = ref.LicenseNote
		if len(evidence.LicenseNote) > 64<<10 {
			evidence.LicenseNote = ""
			evidence.EvidenceGaps = []string{"raw license note exceeds the bounded evidence limit"}
		}
		if validCatalogDiscoveryAuthority(string(ref.SourceAuthority)) && (!advisory || ref.SourceAuthority == catalog.SourceAuthorityOfficialDocs) {
			evidence.Authority = string(ref.SourceAuthority)
		}
		break
	}
	return evidence
}

func catalogDiscoveryFilterExclusion(evidence CatalogDiscoverySourceEvidence, filters CatalogDiscoveryFilters) string {
	if len(filters.Authorities) > 0 {
		found := false
		for _, authority := range filters.Authorities {
			if authority == evidence.Authority {
				found = true
				break
			}
		}
		if !found {
			return "source authority is excluded by the requested filter"
		}
	}
	if filters.RequireLicenseEvidence && evidence.LicenseIdentifier == "unknown" {
		return "license identification evidence is unknown"
	}
	if filters.RequireRedistributionEvidence && evidence.Redistribution == "unknown" {
		return "redistribution evidence is unknown"
	}
	return ""
}

func missingCatalogDiscoveryEvidence(report CatalogDiscoveryReport, cat catalog.Catalog, filters CatalogDiscoveryFilters, code, message string) CatalogDiscoveryReport {
	scope := map[string]bool{}
	for _, id := range report.Scope.ProviderIDs {
		scope[id] = true
	}
	for _, provider := range cat.ListProviders() {
		if !scope[provider.ID] {
			continue
		}
		for _, ref := range provider.SpecReferences {
			evidence := catalogDiscoverySourceEvidence(provider, ref.ID, "", false)
			if reason := catalogDiscoveryFilterExclusion(evidence, filters); reason != "" {
				report.Exclusions = append(report.Exclusions, CatalogDiscoveryExclusion{ProviderID: provider.ID, SpecRefID: ref.ID, Reason: reason})
				continue
			}
			report.Leads = append(report.Leads, CatalogDiscoveryLead{ProviderID: provider.ID, DisplayName: provider.DisplayName, SpecRefID: ref.ID, Kind: string(ref.Kind), Reason: message, Evidence: evidence})
			report.Coverage = append(report.Coverage, CatalogIndexCoverage{ProviderID: provider.ID, SpecRefID: ref.ID, State: "missing", Reason: message})
		}
	}
	report.Incomplete = true
	report.Scope.Complete = false
	report.Diagnostics = append(report.Diagnostics, Diagnostic{Severity: "warning", Code: code, Message: message})
	return report
}

func reportWithDiscoveryBlock(report CatalogDiscoveryReport) CatalogDiscoveryReport {
	report.Outcome = CatalogDiscoveryBlocked
	report.Incomplete = true
	report.Scope.Complete = false
	report.Candidates = nil
	return report
}

func catalogDiscoveryRefusal(report CatalogDiscoveryReport, code, message string) catalogDiscoveryLocalEvidence {
	report.Diagnostics = append(report.Diagnostics, Diagnostic{Severity: "error", Code: code, Message: message})
	return catalogDiscoveryLocalEvidence{report: reportWithDiscoveryBlock(report)}
}

func catalogDiscoveryInterrupted(report CatalogDiscoveryReport, err error) catalogDiscoveryLocalEvidence {
	report.Candidates = nil
	report.QualifiedOperations = 0
	report.Outcome = CatalogDiscoveryInsufficientEvidence
	report.Incomplete = true
	report.Scope.Complete = false
	report.Diagnostics = append(report.Diagnostics, Diagnostic{Severity: "warning", Code: "discovery.interrupted", Message: fmt.Sprintf("local catalog search is incomplete: %v", err)})
	return catalogDiscoveryLocalEvidence{report: report}
}

func catalogDiscoveryCandidateKey(candidate CatalogDiscoveryCandidate) string {
	if len(candidate.References) == 0 {
		return string(candidate.Candidate.Source.Kind) + "\x00" + candidate.Candidate.Source.SHA256 + "\x00" + candidate.Candidate.Source.Selector
	}
	return catalogArtifactReferenceKey(candidate.References[0])
}
