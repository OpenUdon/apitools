# Retired milestone M66 - gRPC/Protobuf Source Support

**Milestone.** M66
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M66.md
**Source specification.** tabilet/memory-bank/milestone.md#m66---grpcprotobuf-source-support
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M66 - gRPC/Protobuf Source Support

**Goal.** Add metadata-only gRPC/protobuf source support for UWS 1.4
`grpc-protobuf` source descriptions so provider-owned protobuf descriptors can
participate in source-aware authoring and review without REST-shaped lowering.

**Scope.**

- Create a focused gRPC/protobuf metadata package/API boundary.
- Detect local `.proto`, descriptor set, and reviewed descriptor artifacts
  without using server reflection.
- Parse packages, services, methods, request/response messages, streaming
  shape, fields, and source refs.
- Define stable selector metadata for service/method `sourceOperationId` and
  descriptor `sourceOperationRef`.
- Classify and materialize gRPC/protobuf as its own source family and UWS source
  type.
- Do not execute RPCs, use server reflection, open network channels, negotiate
  TLS, resolve credentials, or select accounts/projects.

**Acceptance.** gRPC/protobuf artifacts can be parsed or inspected offline into
deterministic service/method metadata, selector summaries are available to
downstream authoring/review tools, catalog/materialization identifies
gRPC/protobuf as its own source family, and focused parser/catalog, Go, vet, and
diff checks pass.
````

## Status record

````markdown
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
````
