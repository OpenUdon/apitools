# Retired milestone M19 - Spec Protocol Classification

**Milestone.** M19
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M19.md
**Source specification.** tabilet/memory-bank/milestone.md#m19---spec-protocol-classification
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M19 - Spec Protocol Classification

**Goal.** Represent OpenAPI-adjacent server-side API description sources as
first-class catalog metadata without treating them as OpenAPI import inputs.

**Scope.**

- Add a protocol/model classification separate from catalog source reference
  kind.
- Classify OpenAPI, Swagger, Smithy JSON, Google Discovery, Dropbox Stone,
  OpenAPI indexes, human docs, and unknown references.
- Preserve protocol version hints when known, including OpenAPI/Swagger
  versions from metadata or source notes.
- Expose classification in advisory, specs, refresh, refresh-report, and
  inspect output.
- Keep existing OpenAPI import and validation behavior strict; do not add
  provider API execution, credential resolution, request signing, or automatic
  converters for non-OpenAPI sources.

**Acceptance.** Downstream callers can distinguish source kind from API
description protocol/model family in library and CLI output, while refresh and
review workflows still save non-OpenAPI machine specs only as review artifacts
and all catalog quality checks remain clean.
````

## Status record

````markdown
# Status M19 - Spec Protocol Classification

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Represent official server-side API description sources beyond OpenAPI as
catalog metadata that downstream tools can understand without broadening strict
OpenAPI import behavior.

## Direction

Catalog source kind identifies where a reference sits in the curation model;
spec protocol identifies the API description model family. OpenAPI/Swagger
remain the only strict import targets today, while Smithy, Google Discovery,
Dropbox Stone, OpenAPI indexes, and human docs remain review or advisory
metadata until an explicit importer/converter exists.

## Tasks

| Item | State | Notes |
|---|---|---|
| Protocol model added | `[+]` | Added `SpecProtocol` and protocol classification helpers for OpenAPI, Swagger, Smithy, Google Discovery, Dropbox Stone, OpenAPI indexes, human docs, and unknown references. |
| Reports expose protocol | `[+]` | Advisory, specs, inspect, refresh, and refresh-report result structs include protocol and version fields where available. |
| Version hints preserved | `[+]` | Static catalog rows infer conservative OpenAPI/Swagger version hints; refresh and refresh-report use parsed OpenAPI/Swagger metadata when available. |
| Tests updated | `[+]` | Added protocol classification coverage for catalog rows, advisory output, refresh results, refresh-review results, and CLI specs/inspect output. |
| Verification completed | `[+]` | Ran Go tests/vet, diff check, catalog quality, specs/advisory/inspect/refresh-report CLI smoke checks, and sibling `openudon`/`udon` tests. |

## Boundary Checks

- Do not execute provider API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not treat Smithy, Discovery, Stone, OpenAPI indexes, or human docs as
  import-ready OpenAPI documents without an explicit importer/converter.
- Do not treat advisory overlays as provider truth.
````
