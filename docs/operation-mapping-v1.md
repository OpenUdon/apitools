# Exact operation eligibility and field mapping

M91 adds metadata helpers over the existing `uws.shape-table.v1` contract.
They do not choose an operation, admit a source, approve a workflow, resolve a
credential or execute a request. Consumers must independently reproduce their
shape table against the selected source bytes before trusting its claims.

## Existing behavior inventory

At the M91 integration baseline `1aa42bbbd6cce01b902e1ce612243ec0b61b928a`:

- APItools already supplies `BuildOperationCandidates`, typed contract ranking,
  `SelectOperationByHints`, `BuildOperationShapeTable`,
  `VerifyOperationShapeTable` and `operationlifecycle.Expand`. These retain
  native selectors, raw source SHA256, partial capabilities, security
  alternatives and tied choices. M91 reuses these contracts.
- Kinet's `internal/workflow/mapping.go` owns purpose/effect/authentication
  admission, exact confirmed contract fields, mapping proposals and final
  publication. `native_candidate_projection.go` separately retains nullable,
  completeness and native binding proof beyond the frozen display card.
  Product prompts, credential symbols and confirmation stay with Kinet.
- Ramen's `authoring/role_bindings.go` owns resource role/identity bindings,
  required parameters and scalar read-response projection. It excludes
  transport and credential headers. `internal/wflowservice/service.go` owns
  saving the selected role's bindings. Desired-state semantics, lifecycle
  selection, resource identity and persisted project formats stay with Ramen.

These source locations are an inventory, not imports or cross-product APIs.
M91 adds no catalog-root ingestion (reserved APItools:M90), transport,
universal planner, source-family expansion or HCL removal.

## Frozen common behavior

`AssessOperationEligibility` resolves one explicit `binding.Binding` from a
validated private snapshot of a shape table. Its result preserves the complete
native shape and exact source identity. Missing/stale identity is refused;
multiple native matches remain ambiguous. Incomplete shapes, unknown schemas,
unknown security and a caller-declared truncated report remain indeterminate.
An optional `OperationCandidate` must agree with the selected native source
kind, SHA256 and reference selector; ranking/display text cannot replace shape
evidence. Eligibility describes metadata sufficiency, never product admission.

`MatchOperationFields` assesses explicit named input/output mappings over that
same selection. Each pair names one declared native field by direction,
location and name, plus the consumer field's schema and presence guarantee.
It compares whole field schemas through the existing UWS binding validator:
consumer schema flows into an input; operation schema flows into an output.
Objects and arrays retain their required properties and item schemas. No dotted
path guessing, schema extraction, coercion, value evaluation or generated
binding expression is returned. Nested per-property authoring remains a
consumer adapter concern; an unsupported field reference stays unresolved.

Required native inputs must have a non-placeholder mapping from a guaranteed
consumer field. Optional or unknown consumer presence cannot supply a required
input. Duplicate target inputs or duplicate consumer outputs are refused;
operation aliases are resolved without choosing between ambiguous operations.
Unknown schemas and unproved schema containment remain indeterminate.
Placeholders are explicit indeterminate pairs, never successful mappings.

Credential locations named in any security alternative and reserved/credential
shaped request headers are excluded from ordinary mapping. Credential-shaped
field/property names are excluded conservatively in either direction. Excluded
required inputs are reported separately for product-owned security handling;
their exclusion never claims authentication has been satisfied. The helpers
take schemas and identity only, never credential values or security bindings.

Results use stable value-free codes, retain each requested pair and the source
shape, and distinguish incompatible from indeterminate. A compatible field
report proves only the declared schema relationships and coverage of ordinary
required inputs. Products still own purpose/effect policy, source admission,
authentication, target choice, user decisions, runtime authority and final
independent binding validation.

Inputs are bounded by existing shape-table/schema limits and 64 mapping pairs.
Malformed metadata returns the existing bounded shape-contract error;
cancellation returns the context error. Neither path exposes parser excerpts.
Versioned conformance cases live under `testdata/operation-mapping/v1` and use
disposable business-message and fake Cloudflare resource metadata only.
