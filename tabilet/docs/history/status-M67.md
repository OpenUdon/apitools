# Retired milestone M67 - OData Source Support

**Milestone.** M67
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M67.md
**Source specification.** tabilet/memory-bank/milestone.md#m67---odata-source-support
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M67 - OData Source Support

**Goal.** Add metadata-only OData source support for UWS 1.4 `odata` source
descriptions so provider-owned or user-exported CSDL/service metadata can
participate in source-aware authoring and review without generic REST-shaped
overlays.

**Scope.**

- Create a focused OData metadata package/API boundary.
- Detect local CSDL XML/JSON and exported service metadata without calling live
  tenant `$metadata` endpoints.
- Parse entity types, entity sets, singletons, navigation properties, functions,
  actions, parameters, and source refs.
- Define stable selector metadata for entity operations, functions, and actions.
- Classify and materialize OData as its own source family and UWS source type.
- Do not execute OData operations, call tenant metadata endpoints, resolve
  credentials, log in to ERP tenants, execute `$batch`, or select
  accounts/tenants.

**Acceptance.** OData artifacts can be parsed or inspected offline into
deterministic metadata, selector summaries are available to downstream
authoring/review tools, catalog/materialization identifies OData as its own
source family, and focused parser/catalog, Go, vet, and diff checks pass.
````

## Status record

````markdown
# Status M67 - OData Source Support

State: Complete

## Scope

Add metadata-only OData source support for UWS 1.4 `odata` source descriptions.
`apitools` should parse or inspect local OData CSDL/service metadata enough to
expose entity set, singleton, action, function, navigation, query-option, and
selector metadata for downstream authoring and review. Tenant login, live
metadata calls, credential resolution, account/tenant selection, `$batch`
execution, and OData request execution remain out of scope.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| Package boundary placeholder | `[+]` | Added focused `odata` parser/model/selector package plus root summary helpers without changing OpenAPI import behavior. |
| Local artifact detection | `[+]` | Added local-only CSDL XML/JSON detection based on supplied bytes and paths; no tenant `$metadata` calls are made. |
| Metadata model and parser | `[+]` | Preserves schemas, entity types, complex types, entity sets, singletons, navigation properties, functions, actions, parameters, imports, and source refs from XML and JSON CSDL. |
| Selector metadata | `[+]` | Defines stable selectors and source refs for entity sets, singletons, navigation properties, actions, functions, and operation imports. |
| Operation summaries | `[+]` | Exposes prompt-safe OData summaries and root `apitools` operation summaries with OData parameter/query-option metadata. |
| Catalog/source classification | `[+]` | Added OData spec kind/protocol, UWS source type `odata`, `odata-only` capability, refresh validation, and `odata/` materialization. |
| Boundary checks | `[+]` | Confirmed implementation parses local bytes only and adds no live tenant metadata calls, ERP login, credential lookup, account/tenant selection, `$batch` execution, OData client, or request execution. |
| Verification | `[+]` | Focused parser/catalog/root tests passed; full verification recorded below. |

## Initial Decision Posture

- OData is now a normative UWS 1.4 source family, but `apitools` support must
  stay metadata-only.
- Prefer user-exported or provider-published CSDL/service metadata. Tenant-
  generated metadata remains operator-supplied local evidence.
- Preserve OData entity/action/function/query semantics instead of hiding them
  inside generic REST-shaped endpoint overlays.

## Completion Notes

- Public API additions are metadata-only: `catalog.SpecKindOData`,
  `catalog.SpecProtocolOData`, `ProviderArtifactCapabilityODataOnly`, the
  `odata` package, and root OData summary helpers.
- Catalog refresh accepts reviewed local/downloaded OData artifacts as
  structured metadata and saves them under `odata/`; materialization preserves
  `uws_source_type: odata`.
- No OpenAPI lowering, parser package side effects, runtime transport,
  tenant/account behavior, credential behavior, or `$batch` support was added.
- Evolution check: no bump. Existing evolution direction already covers UWS
  1.4 native source tooling; M67 completes that planned implementation rather
  than changing product direction or public/private boundaries.

## Verification

- `go test ./odata ./catalog .`
- `go test ./...`
- `go vet ./...`
- `git diff --check` in `../apitools`
- `git diff --check` in `../tofu`
- CLI smoke: `go run ./cmd/apitools search --help`,
  `go run ./cmd/apitools import --help`,
  `go run ./cmd/apitools catalog check`, and
  `go run ./cmd/apitools catalog materialize --help`
- Dependent checks after exported API additions:
  `(cd ../openudon && go test ./...)` and
  `(cd ../udon && go test ./...)`
````
