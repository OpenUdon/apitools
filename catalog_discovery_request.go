package apitools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/sourceguard"
)

// DecodeCatalogDiscoveryRequest accepts only the bounded additive wire.
// Errors omit raw request text. Root, transport and endpoint configuration are
// deliberately absent, and unknown fields are refused at every typed level.
func DecodeCatalogDiscoveryRequest(reader io.Reader) (CatalogDiscoveryRequest, error) {
	if reader == nil {
		return CatalogDiscoveryRequest{}, fmt.Errorf("a catalog discovery request is required")
	}
	data, err := io.ReadAll(io.LimitReader(reader, MaxCatalogDiscoveryRequestBytes+1))
	if err != nil || len(data) > MaxCatalogDiscoveryRequestBytes {
		return CatalogDiscoveryRequest{}, fmt.Errorf("cannot read bounded catalog discovery request")
	}
	if err := sourceguard.CheckJSON("catalog discovery request", data); err != nil {
		return CatalogDiscoveryRequest{}, fmt.Errorf("invalid catalog discovery request JSON")
	}
	var request CatalogDiscoveryRequest
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		return request, fmt.Errorf("invalid catalog discovery request fields")
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return CatalogDiscoveryRequest{}, fmt.Errorf("trailing catalog discovery request data")
	}
	return request, nil
}

func prepareCatalogDiscoveryRequest(request CatalogDiscoveryRequest, cat catalog.Catalog) (CatalogDiscoveryLimits, []string, []Diagnostic) {
	limits := request.Limits
	var diagnostics []Diagnostic
	invalid := func(code, message string) {
		diagnostics = append(diagnostics, Diagnostic{Severity: "error", Code: code, Message: message})
	}
	if request.SchemaVersion != CatalogDiscoverySchemaVersion {
		invalid("discovery.version_invalid", "unsupported catalog discovery schema version")
	}
	if limits.MaxOperations < 0 || limits.MaxOperations > MaxCatalogIndexOperations || limits.MaxResults < 0 || limits.MaxResults > MaxCatalogDiscoveryResults || limits.MaxContextBytes < 0 || limits.MaxContextBytes > MaxCatalogDiscoveryContextBytes || limits.TimeoutMillis < 0 || limits.TimeoutMillis > int(MaxCatalogDiscoveryTimeout.Milliseconds()) {
		invalid("discovery.limit_invalid", "catalog discovery limits exceed supported bounds")
	}
	if limits.MaxContextBytes > 0 && limits.MaxContextBytes < MinCatalogDiscoveryContextBytes {
		invalid("discovery.context_limit_invalid", "context byte limit must leave room for a bounded report")
	}
	if limits.MaxOperations == 0 {
		limits.MaxOperations = DefaultCatalogDiscoveryOperations
	}
	if limits.MaxResults == 0 {
		limits.MaxResults = DefaultCatalogDiscoveryResults
	}
	if limits.MaxContextBytes == 0 {
		limits.MaxContextBytes = DefaultCatalogDiscoveryContextBytes
	}
	if limits.TimeoutMillis == 0 {
		limits.TimeoutMillis = int(DefaultCatalogDiscoveryTimeout.Milliseconds())
	}
	if strings.TrimSpace(request.Contract.Purpose) == "" {
		invalid("discovery.purpose_required", "a documented step purpose is required")
	}
	if len(request.ProviderKeys) > 32 || len(request.Filters.Authorities) > 5 || len(request.Contract.Inputs) > DefaultPromptBudget().MaxFields || len(request.Contract.Outputs) > DefaultPromptBudget().MaxFields {
		invalid("discovery.request_work_limit", "request collections exceed supported work bounds")
		return limits, nil, diagnostics
	}
	_, _, _, contractDiagnostics := prepareStepContract(request.Contract, DefaultPromptBudget())
	diagnostics = append(diagnostics, contractDiagnostics...)
	ids := []string{}
	for _, key := range request.ProviderKeys {
		if _, changed := sanitizePromptString(key, DefaultPromptIdentifierRunes); changed || strings.TrimSpace(key) == "" {
			invalid("discovery.provider_invalid", "provider keys must be bounded nonempty exact catalog keys")
			continue
		}
		provider, ok := cat.FindProvider(key)
		if !ok {
			invalid("discovery.provider_unknown", "an exact provider key has no catalog match")
			continue
		}
		ids = append(ids, provider.ID)
	}
	if request.ProviderKeys == nil {
		for _, provider := range cat.ListProviders() {
			ids = append(ids, provider.ID)
		}
	}
	ids = uniqueSortedStrings(ids)
	for _, authority := range request.Filters.Authorities {
		if !validCatalogDiscoveryAuthority(authority) {
			invalid("discovery.authority_invalid", "source authority filter is unsupported")
		}
	}
	sortDiagnostics(diagnostics)
	return limits, ids, diagnostics
}
