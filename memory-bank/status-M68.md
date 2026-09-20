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
