package apitools

// OperationCandidatesSchemaVersion identifies the JSON contract emitted by
// future step-contract candidate reports.
const OperationCandidatesSchemaVersion = "apitools.operation-candidates/v1"

// OperationSourceKind names a parser or source family supported by apitools.
type OperationSourceKind string

const (
	OperationSourceOpenAPI         OperationSourceKind = "openapi"
	OperationSourceGoogleDiscovery OperationSourceKind = "google-discovery"
	OperationSourceAWSSmithy       OperationSourceKind = "aws-smithy"
	OperationSourceAsyncAPI        OperationSourceKind = "asyncapi"
	OperationSourceGraphQL         OperationSourceKind = "graphql"
	OperationSourceOpenRPC         OperationSourceKind = "openrpc"
	OperationSourceGRPCProtobuf    OperationSourceKind = "grpc-protobuf"
	OperationSourceOData           OperationSourceKind = "odata"
)

// OperationEffect is the metadata assessment of an operation's externally
// visible effect. Unknown is explicit and is never a read-only assertion.
type OperationEffect string

const (
	OperationEffectRead    OperationEffect = "read"
	OperationEffectWrite   OperationEffect = "write"
	OperationEffectUnknown OperationEffect = "unknown"
)

// ContractMatchStatus describes one evidence-backed contract comparison.
type ContractMatchStatus string

const (
	ContractMatchCompatible    ContractMatchStatus = "compatible"
	ContractMatchIncompatible  ContractMatchStatus = "incompatible"
	ContractMatchIndeterminate ContractMatchStatus = "indeterminate"
)

// OperationSourceIdentity binds a candidate to one exact native operation in
// one source artifact. Digest is the lowercase SHA-256 of the source bytes;
// Selector remains native to the source family.
type OperationSourceIdentity struct {
	Kind         OperationSourceKind `json:"kind"`
	DocumentPath string              `json:"document_path,omitempty"`
	DocumentURL  string              `json:"document_url,omitempty"`
	SHA256       string              `json:"sha256"`
	Selector     string              `json:"selector"`
}

// ContractValue describes a declared input or expected output. An empty Type
// means the caller supplied no type evidence, not that every type is accepted.
// Properties and Items preserve the shape needed for nested compatibility.
type ContractValue struct {
	Type   string `json:"type,omitempty"`
	Format string `json:"format,omitempty"`
	// Required is nil when the source or caller has not declared requiredness.
	// False means explicitly optional; true means required.
	Required    *bool                    `json:"required,omitempty"`
	Description string                   `json:"description,omitempty"`
	Properties  map[string]ContractValue `json:"properties,omitempty"`
	Items       *ContractValue           `json:"items,omitempty"`
}

// StepContract is the shared purpose/inputs/outputs/effect declaration used
// to rank candidate operations. Its JSON field names align with the pending
// step contract under design in UWS and OpenUdon.
type StepContract struct {
	Purpose string                   `json:"purpose"`
	Inputs  map[string]ContractValue `json:"inputs,omitempty"`
	Outputs map[string]ContractValue `json:"outputs,omitempty"`
	Effect  OperationEffect          `json:"effect"`
}

// OperationCandidateRequest wraps a step contract in the versioned request
// envelope used by conformance fixtures and API consumers.
type OperationCandidateRequest struct {
	SchemaVersion string       `json:"schema_version"`
	Contract      StepContract `json:"contract"`
}

// OperationSourceInput identifies one bounded, local source artifact. URL is
// provenance only; this API never fetches it. Reports omit URL userinfo,
// query, and fragment components. Content, when supplied, takes precedence
// over Path and is not serialized in reports.
type OperationSourceInput struct {
	Kind    OperationSourceKind `json:"kind"`
	Name    string              `json:"name,omitempty"`
	Path    string              `json:"path,omitempty"`
	URL     string              `json:"url,omitempty"`
	Content []byte              `json:"-"`
}

// OperationCandidateOptions configures bounded, local candidate generation.
// Zero limits select documented defaults; negative limits are invalid.
type OperationCandidateOptions struct {
	Sources       []OperationSourceInput `json:"-"`
	Contract      StepContract           `json:"contract"`
	MaxBytes      int64                  `json:"max_bytes,omitempty"`
	MaxOperations int                    `json:"max_operations,omitempty"`
	MaxCandidates int                    `json:"max_candidates,omitempty"`
	PromptBudget  PromptBudget           `json:"prompt_budget,omitempty"`
}

// OperationEvidence points to source metadata that supports a summary,
// classification, or compatibility result. Reference uses a local selector
// or source-document location; it never causes reference fetching.
type OperationEvidence struct {
	Kind      string `json:"kind"`
	Reference string `json:"reference"`
}

// OperationCapabilityStatus records whether an adapter can expose a metadata
// dimension for a source. Partial means the adapter may expose only the
// source's declared subset; it never asserts that the provider supports only
// that subset.
type OperationCapabilityStatus string

const (
	OperationCapabilitySupported   OperationCapabilityStatus = "supported"
	OperationCapabilityPartial     OperationCapabilityStatus = "partial"
	OperationCapabilityUnsupported OperationCapabilityStatus = "unsupported"
)

// OperationCapability describes one source adapter capability and its gaps.
type OperationCapability struct {
	Dimension string                    `json:"dimension"`
	Status    OperationCapabilityStatus `json:"status"`
	Evidence  []OperationEvidence       `json:"evidence,omitempty"`
	Gaps      []string                  `json:"gaps,omitempty"`
}

// OperationSourceReport summarizes one parsed source independently of its
// candidate operations, including parse diagnostics and semantic limitations.
type OperationSourceReport struct {
	Kind           OperationSourceKind   `json:"kind"`
	Name           string                `json:"name,omitempty"`
	DocumentPath   string                `json:"document_path,omitempty"`
	DocumentURL    string                `json:"document_url,omitempty"`
	SHA256         string                `json:"sha256,omitempty"`
	Title          string                `json:"title,omitempty"`
	OperationCount int                   `json:"operation_count"`
	Capabilities   []OperationCapability `json:"capabilities,omitempty"`
	Diagnostics    []Diagnostic          `json:"diagnostics,omitempty"`
}

// ConsumerOperationSummary gives a workflow author a concise description
// grounded in operation metadata. Gaps records information the source does not
// establish; it is not hidden model reasoning.
type ConsumerOperationSummary struct {
	Description string                  `json:"description,omitempty"`
	Inputs      []OperationValueSummary `json:"inputs,omitempty"`
	Outputs     []OperationValueSummary `json:"outputs,omitempty"`
	Evidence    []OperationEvidence     `json:"evidence,omitempty"`
	Gaps        []string                `json:"gaps,omitempty"`
}

// OperationValueSummary describes one source-declared request or response
// value. Required is nil only when the source metadata does not establish the
// value's requiredness. Response Nullable is true when the value or one of its
// schema ancestors permits null.
type OperationValueSummary struct {
	Name              string              `json:"name"`
	Location          string              `json:"location,omitempty"`
	Type              string              `json:"type,omitempty"`
	Format            string              `json:"format,omitempty"`
	Required          *bool               `json:"required,omitempty"`
	Nullable          bool                `json:"nullable,omitempty"`
	ContainerRequired *bool               `json:"container_required,omitempty"`
	Description       string              `json:"description,omitempty"`
	Evidence          []OperationEvidence `json:"evidence,omitempty"`
}

// EffectAssessment reports the operation effect and the metadata supporting
// it. Reasons are concise, user-readable findings tied to Evidence.
type EffectAssessment struct {
	Class    OperationEffect     `json:"class"`
	Evidence []OperationEvidence `json:"evidence,omitempty"`
	Reasons  []string            `json:"reasons,omitempty"`
}

// ContractDimensionMatch reports one part of purpose/input/output/effect
// compatibility. Missing and Conflicts preserve why a dimension is not a
// positive match.
type ContractDimensionMatch struct {
	Status    ContractMatchStatus `json:"status"`
	Score     int                 `json:"score,omitempty"`
	Reasons   []string            `json:"reasons,omitempty"`
	Evidence  []OperationEvidence `json:"evidence,omitempty"`
	Missing   []string            `json:"missing,omitempty"`
	Conflicts []string            `json:"conflicts,omitempty"`
	Gaps      []string            `json:"gaps,omitempty"`
}

// StepContractMatch is an advisory, evidence-backed comparison. Score is a
// ranking aid only and does not prove that the intended workflow outcome is
// achievable.
type StepContractMatch struct {
	Score   int                    `json:"score"`
	Purpose ContractDimensionMatch `json:"purpose"`
	Inputs  ContractDimensionMatch `json:"inputs"`
	Outputs ContractDimensionMatch `json:"outputs"`
	Effect  ContractDimensionMatch `json:"effect"`
}

// OperationCandidate combines the existing operation summary with native
// source identity and the consumer-facing metadata used by step candidates.
type OperationCandidate struct {
	Source       OperationSourceIdentity  `json:"source"`
	Operation    OperationSummary         `json:"operation"`
	Summary      ConsumerOperationSummary `json:"summary"`
	Capabilities []OperationCapability    `json:"capabilities,omitempty"`
	Effect       EffectAssessment         `json:"effect"`
	Match        StepContractMatch        `json:"match"`
}

// OperationCandidateReport is the versioned public result envelope for
// deterministic step-contract candidate ranking.
type OperationCandidateReport struct {
	SchemaVersion string                  `json:"schema_version"`
	Sources       []OperationSourceReport `json:"sources,omitempty"`
	Candidates    []OperationCandidate    `json:"candidates,omitempty"`
	Diagnostics   []Diagnostic            `json:"diagnostics,omitempty"`
	Truncated     bool                    `json:"truncated"`
}
