# Retired milestone M38 - Google Discovery Converter Ownership

**Milestone.** M38
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M38.md
**Source specification.** tabilet/memory-bank/milestone.md#m38---google-discovery-converter-ownership
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M38 - Google Discovery Converter Ownership

**Goal.** Move reusable Google Discovery to OpenAPI-shaped conversion ownership
into `apitools` so downstream projects share one deterministic metadata-only
converter.

**Scope.**

- Add `github.com/OpenUdon/apitools/googlediscovery` with byte and decoded-map
  conversion APIs.
- Preserve the existing bounded OpenAPI 3.0.1 shaping behavior for Discovery
  resources, methods, schemas, media upload metadata, parameter names,
  deterministic ordering, and Google OAuth scope metadata.
- Port focused Drive and Gmail fixtures and converter tests into `apitools`.
- Keep catalog Google Discovery references classified as `google-discovery` and
  structured review artifacts, not official OpenAPI.
- Do not add CLI commands, execute Google API operations, resolve credentials,
  fetch tokens, sign requests, or select Google accounts.

**Acceptance.** `googlediscovery` is covered by focused tests, full `apitools`
test/vet/diff checks pass, `../udon` delegates reusable conversion to
`apitools` while keeping native document adaptation locally, and no
`apitools` package imports `../udon`.
````

## Status record

````markdown
# Status M38 - Google Discovery Converter Ownership

## Goal

Move reusable Google Discovery to OpenAPI-shaped conversion into `apitools`
while preserving the metadata-only boundary and downstream behavior.

## Status

| Item | State | Notes |
|---|---|---|
| Converter package added | `[+]` | `googlediscovery` exposes `Convert` and `ConvertMap` with Drive/Gmail fixture coverage and no `udon` import. |
| Downstream ownership split updated | `[+]` | `../udon` delegates conversion to `apitools/googlediscovery` and retains native document adaptation. |
| Verification completed | `[+]` | Ran focused/full `apitools` and `udon` test, vet, diff, CLI help, and catalog Discovery classification checks. |
| Review hardening completed | `[+]` | Follow-up review fixes preserve root-level Google Discovery parameters on operations, reject duplicate path/method collisions instead of overwriting operations, and honor non-multipart simple media upload metadata. |

## Boundary Checks

- Google Discovery conversion remains derived metadata, not official provider
  OpenAPI truth.
- Do not execute Google API operations, resolve credentials, fetch tokens, sign
  requests, or choose Google accounts.
- Keep catalog Google Discovery artifacts classified as `google-discovery` /
  structured review evidence unless an explicit caller converts them.
- Direct `../udon` adoption preserves the same derived metadata shape through
  native parsing and generator/runtime-plan lowering.
````
