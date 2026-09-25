# Retired milestone M63 - OpenRPC Source Support

**Milestone.** M63
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M63.md
**Source specification.** tabilet/memory-bank/milestone.md#m63---openrpc-source-support
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M63 - OpenRPC Source Support

**Goal.** Add metadata-only OpenRPC support for JSON-RPC source documents so
provider-owned method contracts can participate in source-aware authoring and
review without being flattened into OpenAPI.

**Scope.**

- Keep the `openrpc` package as the public boundary for OpenRPC parsing and
  summary APIs.
- Parse OpenRPC documents offline, preserving document info, servers, methods,
  params, results, errors, components, and source refs.
- Expose method summary and selector metadata for downstream OpenUdon/UWS
  binding and review evidence.
- Classify OpenRPC artifacts as a distinct source family in catalog and
  materialization output.
- Copy reviewed OpenRPC artifacts and provenance into output packages when
  materialized.
- Do not execute JSON-RPC methods, contact servers, resolve credentials, or
  implement runtime transport behavior in `apitools`.

**Acceptance.** OpenRPC artifacts can be parsed and summarized offline with
stable method selectors, catalog/materialization identifies OpenRPC as its own
source family, downstream authoring tools have enough metadata to bind UWS
operations, and parser, catalog, Go, and diff checks pass.
````

## Status record

````markdown
# Status M63 - OpenRPC Source Support

State: Complete

## Scope

Add metadata-only OpenRPC source support for JSON-RPC services. `apitools`
should parse and summarize provider-owned OpenRPC documents, expose method
selector metadata to downstream authoring tools, and keep runtime invocation
outside this repository.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| Package boundary placeholder | `[+]` | `openrpc` is now the metadata-only parser package boundary. |
| OpenRPC document model and parser | `[+]` | Added standard-library JSON parser preserving document info, servers, methods, params, results, errors, components, raw source maps, and `$ref` values without contacting RPC servers. |
| Method summary and selector metadata | `[+]` | Added deterministic method summaries, canonical JSON-RPC method IDs, local `#/methods/{name}` selectors, `method:{name}` aliases, and root `apitools` operation summaries with OpenRPC selector/result/param metadata. |
| Catalog/source classification | `[+]` | Added `openrpc` spec kind/protocol, UWS source type, refreshable source rows, stats display, and OpenRPC-only provider capability classification. |
| Materialization and provenance | `[+]` | OpenRPC artifacts save and materialize under `openrpc/` with existing provenance manifests; no credential, server, or RPC behavior was added. |
| OpenUdon/UWS coordination | `[+]` | Recorded OpenRPC as a first-class metadata source family and left downstream OpenUdon/UWS fixtures and runtime semantics outside `apitools`. |
| Verification | `[+]` | Focused parser/catalog tests, full Go tests, and follow-up checks run for this milestone. |

## Initial Verification

- Passed: `go test ./openrpc`
- Passed: `go test ./...`
- Passed: `git diff --check`

## Implementation Verification

- Passed: `go test ./openrpc ./catalog`
- Passed: `go test .`
- Passed: `go test ./...`
- Passed: `GOWORK=off go test ./...`
- Passed: `go vet ./...`
- Passed: `GOWORK=off go vet ./...`
- Passed: `go run ./cmd/apitools search --help`
- Passed: `go run ./cmd/apitools import --help`
- Passed: `(cd ../openudon && go test ./...)`
- Passed: `(cd ../udon && go test ./...)`
- Passed: `git diff --check`

## Acceptance Criteria

- OpenRPC artifacts can be parsed offline into deterministic metadata.
- Method selectors and parameter/result metadata are available to downstream
  authoring and review tools.
- OpenRPC remains metadata-only: no JSON-RPC method execution, server contact,
  credential resolution, or runtime transport behavior is introduced.
- Catalog/materialization paths identify OpenRPC as its own source family rather
  than flattening it into OpenAPI.
````
