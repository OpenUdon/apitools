package apitools

import (
	"time"

	"github.com/OpenUdon/apitools/catalog"
)

const CatalogDiscoverySchemaVersion = "apitools.catalog-discovery/v1"

type CatalogDiscoveryOutcome string

const (
	CatalogDiscoveryMatch                CatalogDiscoveryOutcome = "match"
	CatalogDiscoveryAmbiguous            CatalogDiscoveryOutcome = "ambiguous"
	CatalogDiscoveryNoQualifyingAPI      CatalogDiscoveryOutcome = "no_qualifying_api"
	CatalogDiscoveryInsufficientEvidence CatalogDiscoveryOutcome = "insufficient_evidence"
	CatalogDiscoveryBlocked              CatalogDiscoveryOutcome = "blocked"
)

const (
	DefaultCatalogDiscoveryOperations   = 50_000
	DefaultCatalogDiscoveryResults      = 20
	DefaultCatalogDiscoveryContextBytes = 512 << 10
	DefaultCatalogDiscoveryTimeout      = 30 * time.Second
	MaxCatalogDiscoveryResults          = 100
	MaxCatalogDiscoveryContextBytes     = 2 << 20
	MaxCatalogDiscoveryTimeout          = 60 * time.Second
	MaxCatalogDiscoveryRequestBytes     = 64 << 10
	MinCatalogDiscoveryContextBytes     = 4 << 10
)

type CatalogDiscoveryFilters struct {
	Authorities                   []string `json:"authorities,omitempty"`
	RequireLicenseEvidence        bool     `json:"require_license_evidence,omitempty"`
	RequireRedistributionEvidence bool     `json:"require_redistribution_evidence,omitempty"`
}

// Limits bound lookup independently of the larger offline index generation.
// Zero fields select defaults; callers cannot exceed the documented ceilings.
type CatalogDiscoveryLimits struct {
	MaxOperations   int `json:"max_operations,omitempty"`
	MaxResults      int `json:"max_results,omitempty"`
	MaxContextBytes int `json:"max_context_bytes,omitempty"`
	TimeoutMillis   int `json:"timeout_millis,omitempty"`
}

type CatalogDiscoveryRequest struct {
	SchemaVersion string       `json:"schema_version"`
	Contract      StepContract `json:"contract"`
	// nil (omitted/null) searches open scope; an explicit empty array selects
	// empty scope and must not silently broaden into catalog-wide retrieval.
	ProviderKeys []string                `json:"provider_keys"`
	Filters      CatalogDiscoveryFilters `json:"filters,omitempty"`
	Limits       CatalogDiscoveryLimits  `json:"limits,omitempty"`
	RemoteLookup bool                    `json:"remote_lookup,omitempty"`
}

// Installation configuration is separate from the untrusted request wire.
// RemoteLookup requires both explicit request opt-in and installation opt-in.
type CatalogDiscoveryOptions struct {
	Request       CatalogDiscoveryRequest `json:"-"`
	Index         CatalogIndexOptions     `json:"-"`
	RemoteEnabled bool                    `json:"-"`
	RemoteClient  *Client                 `json:"-"`
}

type CatalogDiscoveryScope struct {
	ProviderIDs         []string `json:"provider_ids"`
	ProviderConstrained bool     `json:"provider_constrained"`
	Complete            bool     `json:"complete"`
}

// LicenseNote is verbatim bounded evidence. Explicit unknowns must not be
// inferred into license identification or redistribution permission.
type CatalogDiscoverySourceEvidence struct {
	ProviderID        string   `json:"provider_id"`
	SpecRefID         string   `json:"spec_ref_id"`
	ArtifactID        string   `json:"artifact_id,omitempty"`
	Authority         string   `json:"authority"`
	Advisory          bool     `json:"advisory,omitempty"`
	LicenseNote       string   `json:"license_note,omitempty"`
	LicenseIdentifier string   `json:"license_identifier"`
	Redistribution    string   `json:"redistribution"`
	EvidenceGaps      []string `json:"evidence_gaps,omitempty"`
}

type CatalogDiscoveryLead struct {
	ProviderID  string                          `json:"provider_id"`
	DisplayName string                          `json:"display_name"`
	SpecRefID   string                          `json:"spec_ref_id,omitempty"`
	Kind        string                          `json:"kind,omitempty"`
	Reason      string                          `json:"reason"`
	Evidence    CatalogDiscoverySourceEvidence  `json:"evidence"`
	Reference   *CatalogArtifactReference       `json:"reference,omitempty"`
	Remote      *CatalogDiscoveryRemoteEvidence `json:"remote,omitempty"`
}

type CatalogDiscoveryCandidate struct {
	Candidate           OperationCandidate               `json:"candidate"`
	References          []CatalogArtifactReference       `json:"references,omitempty"`
	Sources             []CatalogDiscoverySourceEvidence `json:"sources"`
	ProviderConstrained bool                             `json:"provider_constrained"`
	Qualified           bool                             `json:"qualified"`
	QualificationGaps   []string                         `json:"qualification_gaps,omitempty"`
	Remote              *CatalogDiscoveryRemoteEvidence  `json:"remote,omitempty"`
}

type CatalogDiscoveryExclusion struct {
	ProviderID string `json:"provider_id"`
	SpecRefID  string `json:"spec_ref_id"`
	ArtifactID string `json:"artifact_id,omitempty"`
	Reason     string `json:"reason"`
}

type CatalogDiscoveryRemoteEvidence struct {
	Authority string `json:"authority"`
	FinalURL  string `json:"final_url"`
	SHA256    string `json:"sha256"`
	Bytes     int64  `json:"bytes"`
}

type CatalogDiscoveryTie struct {
	Score            int   `json:"score"`
	CandidateIndexes []int `json:"candidate_indexes"`
}

// Reports are advisory evidence, never routing, approval or execution authority.
// A positive result can coexist with unexamined scope; only complete scoped
// negative evidence yields no_qualifying_api. Truncation cannot hide ambiguity.
type CatalogDiscoveryReport struct {
	SchemaVersion       string                      `json:"schema_version"`
	Outcome             CatalogDiscoveryOutcome     `json:"outcome"`
	CatalogSHA256       string                      `json:"catalog_sha256,omitempty"`
	Scope               CatalogDiscoveryScope       `json:"scope"`
	Candidates          []CatalogDiscoveryCandidate `json:"candidates,omitempty"`
	Leads               []CatalogDiscoveryLead      `json:"leads,omitempty"`
	Coverage            []CatalogIndexCoverage      `json:"coverage,omitempty"`
	Exclusions          []CatalogDiscoveryExclusion `json:"exclusions,omitempty"`
	Diagnostics         []Diagnostic                `json:"diagnostics,omitempty"`
	Limits              CatalogDiscoveryLimits      `json:"limits"`
	ExaminedOperations  int                         `json:"examined_operations"`
	QualifiedOperations int                         `json:"qualified_operations"`
	Truncated           bool                        `json:"truncated"`
	Incomplete          bool                        `json:"incomplete"`
	Ties                []CatalogDiscoveryTie       `json:"ties,omitempty"`
}

// Keep the catalog source vocabulary separate from explicit unknown evidence.
func validCatalogDiscoveryAuthority(value string) bool {
	switch value {
	case string(catalog.SourceAuthorityOfficialProvider), string(catalog.SourceAuthorityOfficialGitHub), string(catalog.SourceAuthorityOfficialDocs), string(catalog.SourceAuthorityPublicCatalog), "unknown":
		return true
	default:
		return false
	}
}
