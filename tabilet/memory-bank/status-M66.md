# Status M66 - gRPC/Protobuf Source Support

State: Complete

## Scope

Add metadata-only gRPC/protobuf source support for UWS 1.4 `grpc-protobuf`
source descriptions. `apitools` should parse or inspect local `.proto` files
and descriptor sets enough to expose service, method, message, field, and
selector metadata for downstream authoring and review. Runtime channel
creation, server reflection, network calls, credential resolution, TLS policy,
and RPC execution remain out of scope.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| Package boundary placeholder | `[+]` | Added focused `grpcproto` metadata package/API boundary and root operation-summary helper without changing OpenAPI import behavior. |
| Local artifact detection | `[+]` | Added local-only detection for `.proto`, JSON descriptor artifacts, binary descriptor sets, and reviewed generated descriptor artifacts; no server reflection was introduced. |
| Metadata model and parser | `[+]` | Preserved packages, services, methods, request/response message names, streaming shape, message fields, enums, imports, and source refs. |
| Selector metadata | `[+]` | Stable `sourceOperationId` values use `<package.Service>/<Method>` and stable `sourceOperationRef` values use selectors such as `#/services/{service}/methods/{method}`. |
| Operation summaries | `[+]` | Added prompt-safe gRPC/protobuf operation summaries with full method name, service, request/response types, request fields, streaming flags, provenance, and source selector extensions. |
| Catalog/source classification | `[+]` | Added gRPC/protobuf spec kind/protocol, UWS source type `grpc-protobuf`, `grpc-protobuf-only` capability classification, refresh validation, and `grpc-protobuf/` materialization posture. |
| Boundary checks | `[+]` | Confirmed no RPC execution, server reflection, credential lookup, TLS negotiation, endpoint probing, channel creation, token fetching, or account/project selection. |
| Verification | `[+]` | Focused parser/catalog tests, full `go test ./...`, `go vet ./...`, CLI smoke checks, and `git diff --check` passed. |

## Initial Decision Posture

- gRPC/protobuf is now a normative UWS 1.4 source family, but `apitools`
  support must stay metadata-only.
- Prefer provider-owned `.proto` files or reviewed descriptor sets. Generated
  client libraries are not source truth unless they are accompanied by a
  reviewed descriptor artifact.
- Preserve protobuf and gRPC semantics instead of lowering service methods into
  generic REST-shaped OpenAPI.

## Completion Notes

- gRPC/protobuf artifacts are parsed as native source metadata. `.proto` files,
  JSON descriptor artifacts, and binary `FileDescriptorSet` bytes expose
  package, service, method, streaming, message, field, enum, selector, and
  source-ref metadata.
- Catalog refresh validates gRPC/protobuf artifacts with the same local parser
  and saves source-aligned artifacts under `grpc-protobuf/`. Materialization
  and provider resolution report `uws_source_type: grpc-protobuf` and
  `grpc-protobuf-only` capability when gRPC/protobuf is the strongest available
  machine source.
- M66 keeps the existing boundary: no RPC execution, no server reflection, no
  gRPC channel creation, no TLS negotiation, no endpoint probing, no credential
  or token behavior, no account/project selection, and no OpenAPI-shaped
  gRPC/protobuf lowering.
- Evolution decision: no bump. Evolution v18 already opened the UWS 1.4 source
  tooling roadmap and M66 implements that recorded direction without changing
  the product or public/private contract direction.
- Verification completed with `go test ./...`, `go vet ./...`,
  `git diff --check` in `../apitools`, `git diff --check` in `../tofu`,
  `go run ./cmd/apitools search --help`, `go run ./cmd/apitools import --help`,
  `go run ./cmd/apitools catalog check`, and
  `go run ./cmd/apitools catalog materialize --help`. Because M66 adds
  exported root helpers, dependent `go test ./...` checks also passed in
  `../openudon` and `../udon`.
