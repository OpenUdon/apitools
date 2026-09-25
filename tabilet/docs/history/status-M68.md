# Retired milestone M68 - UWS 1.4 Source Artifact Corpus

**Milestone.** M68
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M68.md
**Source specification.** tabilet/memory-bank/milestone.md#m68---uws-14-source-artifact-corpus
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M68 - UWS 1.4 Source Artifact Corpus

**Goal.** Add a small reviewed local artifact corpus for the four UWS 1.4 source
families so refresh-review validation covers real official examples without
expanding parser or runtime behavior.

**Scope.**

- Add two official example artifacts each for GraphQL, OpenRPC,
  gRPC/protobuf, and OData under source-aligned `catalog-openapi-cache/`
  directories.
- Register the examples as curation-only artifact rows with provider IDs
  `graphql-examples`, `openrpc-examples`, `grpc-protobuf-examples`, and
  `odata-examples`.
- Add root refresh-review coverage that validates each committed artifact as
  `valid-structured`, asserts the expected protocol, records byte/SHA-256
  evidence, and confirms nonzero native operation summaries.
- Keep the examples out of the durable built-in provider catalog.
- Do not add new source types, OpenAPI lowering, parser expansion, runtime
  transport behavior, credentials, endpoint calls, or provider account
  behavior.

**Acceptance.** All eight selected official artifacts are tracked, registered
as curation-only local artifacts, validated through the refresh-review path with
nonzero operation summaries, documented as a local corpus, and focused
parser/catalog/root, full Go, vet, CLI smoke, and diff checks pass.
````

## Status record

````markdown
# Status M68 - UWS 1.4 Source Artifact Corpus

State: Complete

## Scope

Add a small reviewed local artifact corpus for the four UWS 1.4 source families
using two official example artifacts per family: GraphQL, OpenRPC,
gRPC/protobuf, and OData. This milestone is artifact curation and validation
coverage only. It must not add source types, OpenAPI lowering, parser behavior
expansion, runtime transport behavior, credentials, endpoint calls, or provider
account behavior.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| Artifact selection | `[+]` | Selected exactly two official upstream example artifacts each for GraphQL, OpenRPC, gRPC/protobuf, and OData. |
| Source-aligned cache corpus | `[+]` | Added tracked files under `catalog-openapi-cache/graphql/`, `openrpc/`, `grpc-protobuf/`, and `odata/`, with `.gitignore` narrowed to those reviewed examples. |
| Curation-only registry rows | `[+]` | Added artifact registry entries for `graphql-examples`, `openrpc-examples`, `grpc-protobuf-examples`, and `odata-examples` with `valid-structured` status. |
| Refresh-review coverage | `[+]` | Added a root corpus test that reads all eight committed artifacts, validates them as structured source metadata, asserts protocol/byte/SHA-256 evidence, checks nonzero operation summaries, and runs `BuildCatalogSpecRefreshReviewReport`. |
| Public docs | `[+]` | Updated `docs/non-openapi-protocols.md` with a short note about the reviewed UWS 1.4 local artifact corpus. |
| Boundary checks | `[+]` | Confirmed no built-in provider catalog rows, new source types, OpenAPI lowering, parser expansion, runtime behavior, credential handling, endpoint calls, or account selection were added. |
| Verification | `[+]` | Focused parser/catalog/root checks passed; full verification recorded below. |

## Completion Notes

- The corpus consists of official upstream examples only:
  `graphiql-starwars-schema.graphql`, `graphiql-schema-kitchen-sink.graphql`,
  `openrpc-petstore.json`, `openrpc-simple-math.json`,
  `google-pubsub-v1.proto`, `opentelemetry-trace-service-v1.proto`,
  `odata-csdl-16-1.json`, and `odata-operations-service.xml`.
- Registry rows are local curation entries only and do not promote these
  examples into the durable provider catalog.
- Evolution check: no bump. M68 adds validation corpus coverage for the
  already-established UWS 1.4 source-family direction and does not change
  product direction, architecture boundaries, milestone sequencing, or
  public/private contracts.

## Verification

- `go test ./graphql ./openrpc ./grpcproto ./odata .`
- `go test ./catalog .`
- `go test ./...`
- `go vet ./...`
- `git diff --check` in `../apitools`
- `go run ./cmd/apitools catalog check`
- `go run catalog-openapi-cache/artifact-registry/register_catalog_artifacts.go`
- `go run ./cmd/apitools catalog refresh-report`
- `git diff --check` in `../tofu`
````
