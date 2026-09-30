package apitools

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/OpenUdon/apitools/catalog"
)

const catalogDiscoveryRemoteTimeout = 8 * time.Second
const catalogDiscoveryRemoteDocuments = 3

func retrieveCatalogDiscoveryRemote(ctx context.Context, evidence catalogDiscoveryLocalEvidence, options CatalogDiscoveryOptions) catalogDiscoveryLocalEvidence {
	report := &evidence.report
	if report.Scope.ProviderConstrained && len(report.Scope.ProviderIDs) == 0 {
		return evidence
	}
	report.Incomplete, report.Scope.Complete = true, false
	report.Diagnostics = append(report.Diagnostics, Diagnostic{Severity: "warning", Code: "discovery.remote_partial", Message: "bounded public-catalog lookup is partial evidence and cannot establish API absence"})
	if report.ExaminedOperations >= report.Limits.MaxOperations {
		return evidence
	}
	client := *options.RemoteClient.effective()
	client.Cache = nil // ephemeral evidence never reads or writes a source cache
	client.Timeout = catalogDiscoveryRemoteTimeout
	if client.MaxBytes <= 0 || client.MaxBytes > DefaultMaxBytes {
		client.MaxBytes = DefaultMaxBytes
	}
	if client.HTTPClient != nil {
		httpClient := *client.HTTPClient
		httpClient.Jar = nil           // installation cookies are not discovery credentials
		httpClient.CheckRedirect = nil // no installation callback may attach credentials
		client.HTTPClient = &httpClient
	}
	listURL := strings.TrimSpace(client.APIsGuruListURL)
	if listURL == "" {
		listURL = DefaultAPIsGuruListURL
	}
	if _, err := client.validateCacheURL(listURL); err != nil {
		return catalogDiscoveryRefusal(*report, "discovery.remote_unsafe", "remote catalog URL violates the configured URL safety policy")
	}
	if _, err := client.redirectSafeClient(); err != nil {
		return catalogDiscoveryRefusal(*report, "discovery.remote_transport", "remote transport violates the configured URL safety policy")
	}
	network, cancel := context.WithTimeout(ctx, catalogDiscoveryRemoteTimeout)
	defer cancel()
	search, err := client.Search(network, SearchOptions{Query: options.Request.Contract.Purpose, Limit: catalogDiscoveryRemoteDocuments, Source: SourceAPIsGuru, CacheMode: CacheModeBypass})
	if err != nil {
		return catalogDiscoveryRemoteFailure(evidence, err)
	}
	if len(search.Results) == 0 {
		report.Diagnostics = append(report.Diagnostics, Diagnostic{Severity: "warning", Code: "discovery.remote_empty", Message: "bounded public-catalog lookup returned no source leads; this does not establish absence"})
	}
	cat := catalog.BuiltInCatalog()
	if options.Index.Catalog != nil {
		cat = *options.Index.Catalog
	}
	seen := map[string]int{}
	nativeKey := func(candidate OperationCandidate) string {
		return string(candidate.Source.Kind) + "\x00" + candidate.Source.SHA256 + "\x00" + candidate.Source.Selector
	}
	for i, item := range evidence.candidates {
		seen[nativeKey(item.Candidate)] = i
	}
	for _, lead := range search.Results {
		if err := network.Err(); err != nil {
			return catalogDiscoveryRemoteFailure(evidence, err)
		}
		if report.ExaminedOperations >= report.Limits.MaxOperations {
			report.Truncated = true
			break
		}
		sourceEvidence := CatalogDiscoverySourceEvidence{Authority: "public-catalog", Advisory: true, LicenseIdentifier: "unknown", Redistribution: "unknown"}
		remoteSources := []CatalogDiscoverySourceEvidence{sourceEvidence}
		if report.Scope.ProviderConstrained {
			remoteSources = nil
			// Public catalog provider labels do not establish a canonical provider.
			// A constrained lookup requires an exact catalog source URL link.
			for _, provider := range cat.ListProviders() {
				for _, id := range report.Scope.ProviderIDs {
					if id != provider.ID {
						continue
					}
					for _, ref := range provider.SpecReferences {
						if ref.URL == lead.SpecURL {
							linked := sourceEvidence
							linked.ProviderID, linked.SpecRefID = provider.ID, ref.ID
							remoteSources = append(remoteSources, linked)
						}
					}
				}
			}
			if len(remoteSources) == 0 {
				report.Exclusions = append(report.Exclusions, CatalogDiscoveryExclusion{Reason: "remote source has no exact catalog URL link to the constrained providers"})
				continue
			}
		}
		sourceEvidence = remoteSources[0]
		if reason := catalogDiscoveryFilterExclusion(sourceEvidence, options.Request.Filters); reason != "" {
			for _, linked := range remoteSources {
				report.Exclusions = append(report.Exclusions, CatalogDiscoveryExclusion{ProviderID: linked.ProviderID, SpecRefID: linked.SpecRefID, Reason: reason})
			}
			continue
		}
		if _, err := client.validateCacheURL(lead.SpecURL); err != nil {
			return catalogDiscoveryRefusal(*report, "discovery.remote_unsafe", "remote source URL violates the configured URL safety policy")
		}
		content, finalURL, err := client.downloadBounded(network, lead.SpecURL)
		if err != nil {
			evidence = catalogDiscoveryRemoteFailure(evidence, err)
			if evidence.report.Outcome == CatalogDiscoveryBlocked {
				return evidence
			}
			report = &evidence.report
			continue
		}
		digest := sha256.Sum256(content)
		sha := hex.EncodeToString(digest[:])
		cleanURL, _, urlErr := sanitizeOperationSourceURL(finalURL.String())
		if urlErr != nil {
			return catalogDiscoveryRefusal(*report, "discovery.remote_identity", "remote provenance URL is invalid")
		}
		remote := &CatalogDiscoveryRemoteEvidence{Authority: "public-catalog", FinalURL: cleanURL, SHA256: sha, Bytes: int64(len(content))}
		budget := resolvedPromptBudget(PromptBudget{})
		remaining := report.Limits.MaxOperations - report.ExaminedOperations
		if remaining < budget.MaxOperations {
			budget.MaxOperations = remaining
		}
		source := loadedOperationSource{input: OperationSourceInput{Kind: OperationSourceOpenAPI, Name: "apis-guru-" + sha[:16], URL: cleanURL}, content: content, digest: sha}
		if err := validateOperationSourceIdentity(source.input, budget); err != nil {
			return catalogDiscoveryRefusal(*report, "discovery.remote_identity", "remote provenance exceeds source identity safety limits")
		}
		sanitized := map[string]bool{}
		inventory, inventoryErr := buildOperationInventory(network, InventoryOptions{Documents: []InventoryDocument{{Name: source.input.Name, URL: cleanURL, Content: content}}, MaxBytes: DefaultMaxBytes, MaxOperations: budget.MaxOperations}, budget)
		adapted, parseErr := openAPICandidatesFromInventory(network, source, budget, inventory, inventoryErr, 16<<20)
		if parseErr != nil {
			report.Diagnostics = append(report.Diagnostics, Diagnostic{Severity: "warning", Code: "discovery.remote_parse_partial", Message: "remote source metadata could not be completely parsed within existing safety limits"})
		}
		for _, diagnostic := range adapted.diagnostics {
			if diagnostic.Code == "prompt.operation_sanitized" {
				sanitized[diagnostic.Path] = true
			}
		}
		for _, candidate := range adapted.candidates {
			if report.ExaminedOperations >= report.Limits.MaxOperations {
				report.Truncated = true
				break
			}
			item := CatalogDiscoveryCandidate{Candidate: candidate, Sources: remoteSources, ProviderConstrained: report.Scope.ProviderConstrained, Remote: remote}
			if sanitized[candidate.Operation.Provenance] {
				item.QualificationGaps = append(item.QualificationGaps, "remote operation metadata was sanitized; review any identity or selected-field loss")
			}
			if i, duplicate := seen[nativeKey(candidate)]; duplicate {
				evidence.candidates[i].Remote = remote
				evidence.candidates[i].Sources = append(evidence.candidates[i].Sources, remoteSources...)
			} else {
				seen[nativeKey(candidate)] = len(evidence.candidates)
				evidence.candidates = append(evidence.candidates, item)
				report.ExaminedOperations++
			}
		}
		reason := "ephemeral public-catalog source needs explicit registration before local artifact export"
		if len(adapted.candidates) == 0 {
			reason = "fetched public-catalog source has no usable operation metadata"
		}
		for _, linked := range remoteSources {
			report.Leads = append(report.Leads, CatalogDiscoveryLead{ProviderID: linked.ProviderID, SpecRefID: linked.SpecRefID, Kind: "openapi", Reason: reason, Evidence: linked, Remote: remote})
		}
	}
	return evidence
}

func catalogDiscoveryRemoteFailure(evidence catalogDiscoveryLocalEvidence, err error) catalogDiscoveryLocalEvidence {
	// Guard errors may be wrapped by net/http. Never put their URLs or raw
	// response text in ordinary discovery reports.
	message := err.Error()
	for _, refusal := range []string{"refusing ", "userinfo is not allowed", "URL scheme must", "URL port", "custom HTTP transport", "unsupported protocol scheme", "too many redirects", "valid URL is required"} {
		if strings.Contains(message, refusal) {
			return catalogDiscoveryRefusal(evidence.report, "discovery.remote_unsafe", "remote lookup violated URL, redirect or dial-time safety policy")
		}
	}
	code := "discovery.remote_failed"
	if errors.Is(err, context.DeadlineExceeded) {
		code = "discovery.remote_timeout"
	} else if errors.Is(err, context.Canceled) {
		code = "discovery.remote_interrupted"
	}
	evidence.report.Diagnostics = append(evidence.report.Diagnostics, Diagnostic{Severity: "warning", Code: code, Message: "bounded remote lookup did not complete; no absence claim is supported"})
	return evidence
}
