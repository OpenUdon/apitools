# Retired milestone M41 - Google Discovery Native Model Ownership

**Milestone.** M41
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M41.md
**Source specification.** tabilet/memory-bank/milestone.md#m41---google-discovery-native-model-ownership
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M41 - Google Discovery Native Model Ownership

**Goal.** Make `apitools/googlediscovery` own native Google Discovery model
parsing and remove OpenAPI-shaped Discovery conversion as a runtime contract
source.

**Scope.**

- Add `Parse` and `ParseMap` APIs for official Google Discovery documents.
- Preserve service metadata, root/base URL handling, schemas, inherited
  root/resource/method parameters, flattened methods, OAuth scopes, and media
  upload hints without producing OpenAPI.
- Remove `Convert` and `ConvertMap`; downstream callers must use native
  Discovery parsing or their own explicit compatibility adapter.
- Do not execute Google operations, resolve credentials, fetch tokens, sign
  requests, or choose Google accounts or projects.

**Acceptance.** Native parsing is covered by focused Discovery tests,
`../udon` no longer uses converted OpenAPI-shaped Discovery metadata for
runtime-plan generation, no public Discovery-to-OpenAPI conversion API remains
in `apitools`, and required tests/vet/diff checks pass.
````

## Status record

````markdown
# Status M41 - Google Discovery Native Model Ownership

## Goal

Make `apitools/googlediscovery` the native Google Discovery parser while
removing OpenAPI-shaped conversion from the public Discovery API.

## Status

| Item | State | Notes |
|---|---|---|
| Native parser added | `[+]` | `googlediscovery.Parse` and `ParseMap` expose service, schema, inherited parameter, method, OAuth scope, and media upload metadata without producing OpenAPI. The implementation now lives in standalone `github.com/OpenUdon/googlediscovery`; `apitools/googlediscovery` is a deprecated wrapper. |
| Runtime-relevant metadata preserved | `[+]` | Parser coverage preserves flattened operation IDs, service URL/path-prefix behavior, request/response refs, all declared media upload protocols, selected simple/multipart execution hints, and per-method scopes for downstream runtimes. |
| Conversion removed | `[+]` | `Convert` and `ConvertMap` were removed; the standalone `googlediscovery` module exposes native parsing only. |
| Downstream native adoption verified | `[+]` | `../udon` Discovery runtime-plan generation uses native Discovery metadata instead of converted OpenAPI-shaped bytes. |
| Verification completed | `[+]` | Focused tests passed in `apitools`, `udon`, and `grand`; full test/vet/diff checks are part of closeout. |

## Boundary Checks

- `apitools` parses and classifies untrusted Google Discovery metadata only by
  delegating to `github.com/OpenUdon/googlediscovery`.
- OpenAPI-shaped Discovery conversion is not exposed by `apitools`; Discovery
  metadata is parsed natively and is not Google provider OpenAPI truth.
- No Google operations are executed, no credentials are resolved, no tokens are
  fetched, and no Google accounts or projects are selected in `apitools`.
````
