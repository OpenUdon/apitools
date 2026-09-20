# Status M65 - GraphQL Source Support

State: Complete

## Scope

Add metadata-only GraphQL source support for UWS 1.4 `graphql` source
descriptions. `apitools` should parse or inspect local GraphQL SDL,
introspection JSON, and operation documents enough to expose source-aware
operation/schema summaries and selector metadata for downstream authoring and
review. Runtime query execution, endpoint calls, introspection against live
servers, credential resolution, and account/workspace selection remain out of
scope.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| Package boundary placeholder | `[+]` | Added focused `graphql` metadata package/API boundary and root operation-summary helper without changing OpenAPI import behavior. |
| Local artifact detection | `[+]` | Added local-only detection for SDL, introspection JSON, and operation documents; no live GraphQL endpoint probing was introduced. |
| Metadata model and parser | `[+]` | Preserved schema/type/field metadata, schema-root operations, operation documents, variables, operation kind, selection names, and source refs without executing queries. |
| Selector metadata | `[+]` | Stable `sourceOperationId` values use `<kind>.<name>` and stable `sourceOperationRef` values use GraphQL source selectors such as `#/operations/{id}` or `#/schema/{rootType}/fields/{field}`. |
| Operation summaries | `[+]` | Added prompt-safe GraphQL operation summaries with operation kind, GraphQL variable parameters, field/root metadata, provenance, and source selector extensions. |
| Catalog/source classification | `[+]` | Added GraphQL spec kind/protocol, UWS source type `graphql`, `graphql-only` capability classification, refresh validation, and `graphql/` materialization posture. |
| Boundary checks | `[+]` | Confirmed implementation has no GraphQL execution, live introspection, credential lookup, token fetching, endpoint selection, variable resolution, or account/workspace selection. |
| Verification | `[+]` | Focused parser/catalog tests, full `go test ./...`, `go vet ./...`, CLI smoke checks, and `git diff --check` passed. |

## Initial Decision Posture

- GraphQL is now a normative UWS 1.4 source family, but `apitools` support must
  stay metadata-only.
- Prefer provider-owned SDL, introspection JSON, or checked-in operation
  documents. Human docs that merely describe GraphQL remain advisory until a
  machine-readable schema or operation artifact is supplied.
- Do not create REST-shaped `POST /graphql` overlays that hide operation,
  variable, type, query, mutation, or subscription semantics.

## Completion Notes

- GraphQL artifacts are parsed as native source metadata. SDL and
  introspection schemas expose root query/mutation/subscription fields as
  source operations; operation documents preserve named operation kind,
  variables, and top-level selection names.
- Catalog refresh validates GraphQL artifacts with the same local parser and
  saves source-aligned artifacts under `graphql/`. Materialization and provider
  resolution report `uws_source_type: graphql` and `graphql-only` capability
  when GraphQL is the strongest available machine source.
- M65 keeps the existing boundary: no endpoint calls, no remote/live
  introspection, no credential or token behavior, no account/workspace
  selection, and no OpenAPI-shaped GraphQL lowering.
- Evolution decision: no bump. Evolution v18 already opened the UWS 1.4 source
  tooling roadmap and M65 implements that recorded direction without changing
  the product or public/private contract direction.
- Verification completed with `go test ./...`, `go vet ./...`,
  `git diff --check` in `../apitools`, `git diff --check` in `../tofu`,
  `go run ./cmd/apitools search --help`, `go run ./cmd/apitools import --help`,
  `go run ./cmd/apitools catalog check`, and
  `go run ./cmd/apitools catalog materialize --help`.
