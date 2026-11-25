package catalog

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"
	"unicode"
)

//go:embed data/provider_search_profiles.json
var providerSearchProfilesJSON []byte

//go:embed data/provider_search_coverage.json
var providerSearchCoverageJSON []byte

// ProviderSearchProfile describes a reviewed provider for catalog search.
// Description and Keywords describe the service. Coverage describes selected
// operations supported by catalog evidence.
type ProviderSearchProfile struct {
	ProviderID  string                 `json:"provider_id"`
	DisplayName string                 `json:"display_name"`
	Category    string                 `json:"category"`
	Description string                 `json:"description"`
	Aliases     []string               `json:"aliases"`
	Keywords    []string               `json:"keywords"`
	Coverage    ProviderSearchCoverage `json:"coverage"`
}

// ProviderSearchCoverageStatus describes the available operation evidence.
type ProviderSearchCoverageStatus string

const (
	ProviderSearchSelectedSubset ProviderSearchCoverageStatus = "selected-subset"
	ProviderSearchDocsOnly       ProviderSearchCoverageStatus = "docs-only"
	ProviderSearchReferenceOnly  ProviderSearchCoverageStatus = "reference-only"
)

// ProviderSearchCoverage lists provenance for selected operations or source leads.
type ProviderSearchCoverage struct {
	Status  ProviderSearchCoverageStatus `json:"status"`
	Sources []ProviderSearchSource       `json:"sources"`
}

// ProviderSearchSourceKind distinguishes the evidence behind a coverage record.
type ProviderSearchSourceKind string

const (
	ProviderSearchAdvisoryOverlay ProviderSearchSourceKind = "advisory-overlay"
	ProviderSearchPinnedSource    ProviderSearchSourceKind = "pinned-source"
	ProviderSearchSourceReference ProviderSearchSourceKind = "source-reference"
)

// ProviderSearchSource identifies an advisory artifact, a pinned source, or a
// source reference. VerifiedAt is the catalog reference's verification date.
// The selector belongs to the artifact identified by SHA256 when operations
// are present. These provenance fields are for display and source selection.
type ProviderSearchSource struct {
	Kind       ProviderSearchSourceKind  `json:"kind"`
	SpecRefID  string                    `json:"spec_ref_id"`
	ArtifactID string                    `json:"artifact_id,omitempty"`
	Protocol   SpecProtocol              `json:"protocol"`
	Authority  SourceAuthority           `json:"source_authority"`
	SourceURL  string                    `json:"source_url"`
	Revision   string                    `json:"revision,omitempty"`
	SHA256     string                    `json:"sha256,omitempty"`
	VerifiedAt string                    `json:"verified_at,omitempty"`
	Operations []ProviderSearchOperation `json:"operations,omitempty"`
}

// ProviderSearchOperation pairs a public label with an exact source selector.
type ProviderSearchOperation struct {
	Label    string `json:"label"`
	Selector string `json:"selector"`
}

type providerSearchProfileData struct {
	ID          string   `json:"id"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Aliases     []string `json:"aliases,omitempty"`
	Keywords    []string `json:"keywords"`
}

type providerSearchCoverageData struct {
	ProviderID string                    `json:"provider_id"`
	SpecRefID  string                    `json:"spec_ref_id"`
	ArtifactID string                    `json:"artifact_id,omitempty"`
	SHA256     string                    `json:"sha256"`
	Operations []ProviderSearchOperation `json:"operations"`
}

var builtInProviderSearchProfiles = sync.OnceValue(func() []ProviderSearchProfile {
	profiles, err := loadProviderSearchProfiles(providerSearchProfilesJSON, providerSearchCoverageJSON)
	if err != nil {
		panic(err)
	}
	return profiles
})

// BuiltInProviderSearchProfiles returns independent copies of the reviewed
// built-in provider search profiles in exact provider ID order. The profiles
// are embedded and do not require a catalog cache or working tree.
func BuiltInProviderSearchProfiles() []ProviderSearchProfile {
	profiles := slices.Clone(builtInProviderSearchProfiles())
	for i := range profiles {
		profile := &profiles[i]
		profile.Aliases = slices.Clone(profile.Aliases)
		profile.Keywords = slices.Clone(profile.Keywords)
		profile.Coverage.Sources = slices.Clone(profile.Coverage.Sources)
		for j := range profile.Coverage.Sources {
			source := &profile.Coverage.Sources[j]
			source.Operations = slices.Clone(source.Operations)
		}
	}
	return profiles
}

func loadProviderSearchProfiles(profileJSON, coverageJSON []byte) ([]ProviderSearchProfile, error) {
	var profileData []providerSearchProfileData
	if err := decodeProviderSearchData(profileJSON, &profileData); err != nil {
		return nil, fmt.Errorf("provider search profiles: %w", err)
	}
	var coverageData []providerSearchCoverageData
	if err := decodeProviderSearchData(coverageJSON, &coverageData); err != nil {
		return nil, fmt.Errorf("provider search coverage: %w", err)
	}
	providers := BuiltInProviders()
	if len(profileData) != len(providers) {
		return nil, fmt.Errorf("provider search profiles: got %d profiles for %d providers", len(profileData), len(providers))
	}
	coverageByProvider := make(map[string][]providerSearchCoverageData)
	for _, row := range coverageData {
		if _, ok := generatedProviderIndexByID[row.ProviderID]; !ok {
			return nil, fmt.Errorf("provider search coverage: unknown provider %q", row.ProviderID)
		}
		coverageByProvider[row.ProviderID] = append(coverageByProvider[row.ProviderID], row)
	}
	artifacts := make(map[string]CatalogSpecArtifact)
	for _, artifact := range BuiltInEndpointOverlayArtifacts() {
		artifacts[artifact.ArtifactID] = artifact
	}
	addedAliases := make(map[string]bool)
	profiles := make([]ProviderSearchProfile, 0, len(providers))
	for i, provider := range providers {
		row := profileData[i]
		if row.ID != provider.ID || provider.ReviewState != ProviderReviewedCatalogEntry {
			return nil, fmt.Errorf("provider search profile %q does not match reviewed provider %q", row.ID, provider.ID)
		}
		aliases := append(provider.Aliases, row.Aliases...)
		if !validID(row.Category) || row.Category == "" || !validProviderSearchText(provider.DisplayName, 100) ||
			!validProviderSearchText(row.Description, 220) ||
			!validProviderSearchTerms(aliases, 1, 12) || !validProviderSearchTerms(row.Keywords, 3, 8) {
			return nil, fmt.Errorf("provider search profile %q has invalid public text", row.ID)
		}
		for _, alias := range row.Aliases {
			key := normalizeKey(alias)
			if _, exists := builtInProviderLookupIndex[key]; exists || addedAliases[key] {
				return nil, fmt.Errorf("provider search profile %q repeats identity %q", row.ID, alias)
			}
			addedAliases[key] = true
		}
		coverage, err := loadProviderSearchCoverage(provider, coverageByProvider[provider.ID], artifacts)
		if err != nil {
			return nil, fmt.Errorf("provider search coverage for %q: %w", provider.ID, err)
		}
		profiles = append(profiles, ProviderSearchProfile{
			ProviderID: provider.ID, DisplayName: provider.DisplayName, Category: row.Category,
			Description: row.Description, Aliases: aliases, Keywords: row.Keywords, Coverage: coverage,
		})
	}
	return profiles, nil
}

func loadProviderSearchCoverage(provider Provider, rows []providerSearchCoverageData, artifacts map[string]CatalogSpecArtifact) (ProviderSearchCoverage, error) {
	coverage := ProviderSearchCoverage{Status: ProviderSearchSelectedSubset}
	refs := make(map[string]ProviderSearchSource, len(provider.SpecReferences))
	for _, ref := range provider.SpecReferences {
		source := ProviderSearchSource{
			Kind: ProviderSearchSourceReference, SpecRefID: ref.ID, Protocol: ref.ProtocolClassification().Protocol,
			Authority: ref.SourceAuthority, SourceURL: ref.URL, Revision: ref.Revision,
			SHA256: providerSearchRevisionDigest(ref.Revision), VerifiedAt: ref.VerifiedAt,
		}
		refs[ref.ID] = source
	}
	if len(rows) == 0 {
		coverage.Status = ProviderSearchDocsOnly
		for _, ref := range provider.SpecReferences {
			if ref.Kind != SpecKindHumanDocs {
				coverage.Status = ProviderSearchReferenceOnly
			}
			coverage.Sources = append(coverage.Sources, refs[ref.ID])
		}
		return coverage, nil
	}
	seen := make(map[[2]string]bool, len(rows))
	for _, row := range rows {
		source, ok := refs[row.SpecRefID]
		if !ok {
			return ProviderSearchCoverage{}, fmt.Errorf("unknown reference %q", row.SpecRefID)
		}
		key := [2]string{row.SpecRefID, row.ArtifactID}
		if seen[key] || !validProviderSearchDigest(row.SHA256) || !validProviderSearchOperations(row.Operations) {
			return ProviderSearchCoverage{}, fmt.Errorf("invalid source or operations for %q", row.SpecRefID)
		}
		seen[key] = true
		source.Kind = ProviderSearchPinnedSource
		if row.ArtifactID != "" {
			artifact, exists := artifacts[row.ArtifactID]
			if !exists || artifact.ProviderID != provider.ID || !providerSearchAdvisoryMatches(artifact, source) {
				return ProviderSearchCoverage{}, fmt.Errorf("advisory artifact %q does not match the provider and reference", row.ArtifactID)
			}
			source.Kind = ProviderSearchAdvisoryOverlay
			source.ArtifactID = row.ArtifactID
			source.Protocol = SpecProtocolOpenAPI
			source.SourceURL = artifact.SourceURL
			source.Revision = ""
		} else if source.SHA256 != row.SHA256 || source.Protocol == SpecProtocolHumanDocs || source.Protocol == SpecProtocolOpenAPIIndex ||
			source.Authority == SourceAuthorityPublicCatalog {
			return ProviderSearchCoverage{}, fmt.Errorf("pinned source %q does not match accepted source identity", row.SpecRefID)
		}
		source.SHA256 = row.SHA256
		source.Operations = row.Operations
		coverage.Sources = append(coverage.Sources, source)
	}
	return coverage, nil
}

func providerSearchAdvisoryMatches(artifact CatalogSpecArtifact, source ProviderSearchSource) bool {
	if artifact.SpecRefID != artifact.ArtifactID {
		return artifact.SpecRefID == source.SpecRefID
	}
	if artifact.ArtifactID == "webflow-data-api-v2-overlay" {
		// The registration links the reference index. The reviewed catalog
		// reference links its REST introduction.
		return source.SpecRefID == "webflow-data-api-docs"
	}
	return strings.TrimRight(artifact.SourceURL, "/") == strings.TrimRight(source.SourceURL, "/")
}

func decodeProviderSearchData(data []byte, dst any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if dec.Decode(new(any)) != io.EOF {
		return fmt.Errorf("unexpected data after the JSON array")
	}
	return nil
}

func validProviderSearchText(value string, limit int) bool {
	return value != "" && len(value) <= limit && value == strings.TrimSpace(value) &&
		!strings.Contains(value, "://") && strings.IndexFunc(value, unicode.IsControl) == -1
}

func validProviderSearchTerms(terms []string, min, max int) bool {
	if len(terms) < min || len(terms) > max {
		return false
	}
	seen := make(map[string]bool, len(terms))
	for _, term := range terms {
		key := strings.ToLower(term)
		if !validProviderSearchText(term, 64) || seen[key] {
			return false
		}
		seen[key] = true
	}
	return true
}

func validProviderSearchOperations(operations []ProviderSearchOperation) bool {
	if len(operations) == 0 || len(operations) > 64 {
		return false
	}
	seen := make(map[string]bool, len(operations))
	for _, operation := range operations {
		if !validProviderSearchText(operation.Label, 120) || !validProviderSearchText(operation.Selector, 256) ||
			!strings.HasPrefix(operation.Selector, "#/") || len(operation.Selector) <= 2 || seen[operation.Selector] {
			return false
		}
		seen[operation.Selector] = true
	}
	return true
}

func providerSearchRevisionDigest(revision string) string {
	var digest string
	for _, part := range strings.Fields(revision) {
		value, ok := strings.CutPrefix(part, "sha256:")
		if !ok {
			continue
		}
		if !validProviderSearchDigest(value) || (digest != "" && digest != value) {
			return ""
		}
		digest = value
	}
	return digest
}

func validProviderSearchDigest(value string) bool {
	if len(value) != sha256.Size*2 || value != strings.ToLower(value) {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}
