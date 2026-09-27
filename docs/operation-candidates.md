# Step-contract operation candidates

This document defines the additive `apitools.operation-candidates/v1` metadata
contract for consumers such as OpenUdon `step candidates`. The contract is
advisory source metadata: it does not bind a workflow step, approve an
operation, or authorize execution. `BuildOperationCandidates` implements the
local-only v1 producer for the source families listed below.

## Wire contract

`OperationCandidateRequest` carries `schema_version` and one `contract`.
`OperationCandidateReport` returns the same version, per-source reports,
candidates, diagnostics, and a `truncated` flag. Source bytes are supplied to
the Go API separately and are never serialized into reports. A source URL is
provenance only; APItools does not fetch it.

```go
package example

import (
	"context"
	"fmt"
	"os"

	"github.com/OpenUdon/apitools"
)

func candidates(ctx context.Context) error {
	data, err := os.ReadFile("openapi/pets.yaml")
	if err != nil {
		return err
	}
	required := true
	report, err := apitools.BuildOperationCandidates(ctx, apitools.OperationCandidateOptions{
		Sources: []apitools.OperationSourceInput{{
			Kind: apitools.OperationSourceOpenAPI,
			Name: "pets",
			Path: "openapi/pets.yaml",
			Content: data, // When non-nil, bytes take precedence over Path.
		}},
		Contract: apitools.StepContract{
			Purpose: "create a pet",
			Inputs:  map[string]apitools.ContractValue{"name": {Type: "string", Required: &required}},
			Outputs: map[string]apitools.ContractValue{"id": {Type: "string", Required: &required}},
			Effect:  apitools.OperationEffectWrite,
		},
	})
	if err != nil { // Includes incomplete source scans and unsafe truncation.
		return err
	}
	if report.Truncated {
		return fmt.Errorf("candidate report is incomplete")
	}
	// The caller still reviews ties, gaps, and approval policy before binding.
	return nil
}
```

The API accepts explicitly declared local source kinds and either non-nil
inline bytes or a regular local file path. Supplying only `URL` is rejected;
URLs are optional provenance attached to local bytes and are never fetched.
Reports strip URL userinfo, query, and fragment components so credentials and
signed query values are not echoed; the content digest remains the exact source
identity.
Final symlinks, directories, non-regular files, unreadable files, unknown
source kinds, and documents above the byte bound produce blocking diagnostics.
The default bounds are 20 MiB per document, 10,000 operations, 32 sources,
100 returned candidates, and the shared 512 KiB prompt-context budget. Zero
selects defaults; negative explicit source, operation, candidate, or byte
limits are invalid. Reaching a result or source bound is visible and sets
`truncated`; callers must not treat a partial report as a complete search.

The shared step declaration uses these field names:

| Field | Meaning |
|---|---|
| `purpose` | Bounded natural-language intent used only as ranking evidence. |
| `inputs` | Named values the step can supply. Each value may declare type, format, requiredness, description, nested object properties, and array items. |
| `outputs` | Named values the step expects. The same value metadata applies; requiredness describes whether the expected result is mandatory. |
| `effect` | Optional `read`, `write`, or `unknown` constraint. Omission imposes no effect constraint; explicit `unknown` is not a read-only constraint. |

For `required`, `true` means required, `false` means explicitly optional, and
omission means unknown. An empty type is missing type evidence, not a wildcard.
Named nested types that APItools cannot resolve remain opaque and produce an
indeterminate comparison rather than an assumed match. The richer metadata
contract intentionally does not require a downstream workflow-model package.

Each candidate pairs the unchanged public `OperationSummary` with an exact
source identity (`kind`, path or URL, SHA-256, and a source-native `selector`),
a consumer-readable summary, an evidence-bearing effect assessment, and
separate purpose/input/output/effect match results. The consumer summary has a
description plus structured input/output values; each value records name,
location, type/format when known, requiredness, description, and source
evidence. For step matching, the summary's tri-state `Required` values are
authoritative; the paired legacy `OperationSummary` retains its existing
Boolean fields for compatibility. If purpose text is missing, the fallback identifies the operation
and explicitly says that its documented purpose is unavailable. Existing
OpenAPI security alternatives remain on `OperationSummary`; alternatives and
requirements are not flattened. Evidence references are local selectors or
locations only and never trigger reference fetching.

For non-OpenAPI sources, `OperationSummary.Method` and `Path` are display
projections only. Consumers must use `Source.Kind` and `Source.Selector` as the
native protocol identity. Candidate ordering never merges equal operation IDs
from different documents. OpenAPI security requirements retain their ordered
OR-of-AND structure. Google Discovery scopes are represented as separate scope
alternatives, with `auth` capability `partial` because scopes do not describe a
complete authentication or account policy.
An OpenAPI requirement that names an undeclared or incomplete security scheme
is retained as source evidence but marks that candidate's `auth` capability
`partial` with an explicit gap.

Source reports are retained even when parsing yields no candidates. They carry
content digests, counts, capability claims, and parse/adapter diagnostics, so a
zero-result or partially supported source is distinguishable from a source
that was never examined. Capability `partial` means only that a parser can
expose a subset; it does not assert that the provider supports only that
subset.

## Ranking and uncertainty

The total advisory score is the sum of purpose (20 points), input coverage (30),
expected output compatibility (30), and effect (20). An unconstrained effect
adds no points. Purpose uses deterministic token overlap and is explicitly
lexical, not proof that an operation achieves the requested outcome. Input
coverage compares the step's available values with source-declared operation
inputs, including types and requiredness direction. Expected outputs compare
source response values with the step's declared result contract. Missing,
opaque, unresolved, or named protocol types remain indeterminate; even equal
custom type names are not treated as compatible unless their shapes are
resolved. Primitive integer/number compatibility follows the data-flow
direction. A type or requiredness conflict stays visible even when other
dimensions produce a high score. OpenAPI response schemas that permit null are
reported as partial output evidence. The v1 step contract has no nullability
field, so a nullable response cannot earn output compatibility points or a
compatible output status.

Each dimension has `compatible`, `incompatible`, or `indeterminate` status,
with separate evidence, reasons, missing values, conflicts, and gaps. A
`candidate.no_compatible_match` or `candidate.top_tie` warning is not resolved
by sorting: ordering is deterministic, but equal scores remain ambiguous.
Unknown effect cannot satisfy an explicit `read` or `write` constraint.

## Determinism and safety

The producer will process explicit sources in deterministic order, preserve
source-native selectors, and order candidates by descending compatibility
score followed by source kind, source path/URL, digest, and selector. Ties are
reported, not resolved as authority. Work, bytes, candidate count, and prompt
text use explicit bounds; semantic truncation is visible and blocking rather
than silently dropping required inputs, outputs, security alternatives, or
source identity. If the assembled report itself exceeds `MaxContextBytes`, the
producer returns a bounded diagnostic-only report instead of retaining the
oversized candidate/source payload.

Effect is an evidence-based metadata assessment. HTTP methods and names alone
cannot establish `read`; method fields are not classification evidence.
Recognized operation-meaning tokens in the operation ID and leading documented
action phrases can support a classification, while negated, compound, or
conflicting meaning stays `unknown`. Once a leading documented action is found,
the bounded summary or description is scanned for later conflicting verbs,
including verbs in another sentence or separated from a connector by
request-object wording. Source-native protocol semantics may add independent
evidence, and disagreement with operation wording stays `unknown`. An
assessment never substitutes for a user-confirmed decision, downstream
approval, account binding, or runtime enforcement.

All artifacts are treated as untrusted local input. The API does not execute
operations, resolve credentials, select accounts, contact providers, fetch
external references, or make source descriptions authoritative. Summaries are
sanitized and bounded; gaps are explicit and no model chain-of-thought is
exposed.

## Planned v1 source-family coverage

The adapter matrix below is the M77 target. It describes source-native evidence
that the adapters may use, not a promise that every document contains each
dimension. Unsupported or unavailable dimensions must be reported as gaps.

| Source kind | Native operation identity and useful evidence | Implemented limitation to expose |
|---|---|---|
| `openapi` (including Swagger) | Native path/method selector, operation ID, parameters, request/response schemas, descriptions, and OR-of-AND security alternatives. | Unresolved/external references, absent response schemas, multiple media types, and multiple successful response statuses remain gaps; request/response field count/depth truncation blocks that candidate; no remote `$ref` fetch. |
| `google-discovery` | Native method ID/path, parameters, locally declared request/response schemas, and OAuth scope alternatives. | Scope hints are partial auth evidence, not a complete authentication/account policy; unresolved schema names and absent response refs remain gaps. |
| `aws-smithy` | Native operation shape selector, protocol path, top-level input/output members, and explicit `smithy.api#readonly` evidence. | Nested structures and collection item shapes are not expanded; signing traits do not establish credential/account policy; only explicit read-only trait establishes native read evidence. |
| `asyncapi` | Native operation selector/action, channel/message references, and local message payload fields. | Payload/header/binding support is partial; only the first declared payload alternative is summarized; external references are not fetched; send/receive direction does not establish business effect. |
| `graphql` | Native schema-root or operation-document selector, query/mutation kind, variables/root-field arguments, and selected output fields. | Named input/output types remain opaque unless mapped to supported scalar shapes; endpoint auth is unsupported; operation documents may omit schema detail. |
| `openrpc` | Native JSON-RPC method selector, parameter descriptors, result schema, and method summary. | Only local supported JSON Schema fields are summarized; method names do not establish effect; auth is unsupported. |
| `grpc-protobuf` | Native service/method selector, request/response top-level message fields, and streaming flags. | Nested messages and streaming contracts are partial; protobuf shape does not establish auth or read/write effect. |
| `odata` | Native entity/resource/action/function selector, parameters, return/entity type, and collection metadata. | Only action/function kinds provide effect hints; entity-set, singleton, and navigation resources remain effect-unknown because the selector does not identify an HTTP operation. Parameter requiredness and response presence are partial; nullable metadata does not prove field presence; auth is unsupported. |

The `openapi` kind intentionally covers both OpenAPI and Swagger documents; the
parsed document version remains in the existing operation/document metadata.
The other kinds preserve native protocol identity and are not lowered to
OpenAPI. Credential-shaped parameters and schema fields are omitted from
workflow data metadata with an explicit authentication-review gap.

Native-family JSON Schema traversal expands declared properties, array items,
and resolvable local references. `allOf`, `oneOf`, `anyOf`, and dynamic
`additionalProperties` values are not expanded into named fields; each is
reported as a summary gap so missing mappings remain visible.

## Compatibility

This is an additive API. Existing inventory, selection, ranking, and lifecycle
functions and JSON shapes remain unchanged. Consumers should reject an unknown
`schema_version` rather than assuming compatible semantics. A future
incompatible wire change requires a new versioned contract; additional optional
fields within v1 must not change the meaning of existing fields.

Versioned examples and request/report fixtures are stored under
[`testdata/operation-candidates/v1/`](../testdata/operation-candidates/v1/),
including one local source artifact for each supported family.
