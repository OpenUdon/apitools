# Retired milestone M37 - Dropbox Stone-Derived OpenAPI Advisory Artifact

**Milestone.** M37
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M37.md
**Source specification.** tabilet/memory-bank/milestone.md#m37---dropbox-stone-derived-openapi-advisory-artifact
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M37 - Dropbox Stone-Derived OpenAPI Advisory Artifact

**Goal.** Formalize the existing Dropbox core advisory overlay as a narrow
Stone-derived OpenAPI-shaped subset while keeping Dropbox's official catalog
source classified as Dropbox Stone, not official OpenAPI.

**Scope.**

- Preserve `dropbox/dropbox-api-spec` as the official machine-readable source
  with `dropbox-stone` protocol classification and keep official OpenAPI
  availability unavailable.
- Tighten the existing `dropbox-core-api-overlay` metadata so it explicitly
  records `official_openapi: false`, Stone-derived provenance, and source refs
  to the Stone repository, HTTP docs, and OAuth guide.
- Keep the reviewed subset limited to account info, folder listing, metadata,
  upload, download, and shared-link operations.
- Update advisory artifact registry and reporting wording where needed so
  advisory overlays may be derived from official machine specs, not only human
  docs.
- Do not build a general Stone-to-OpenAPI converter, add runtime protocol
  support, execute Dropbox operations, resolve credentials, or change `../udon`.

**Acceptance.** The regenerated Dropbox advisory overlay is deterministic,
catalog advisory output still reports Dropbox's official protocol as
`dropbox-stone`, the OpenAPI-shaped overlay is exposed only as advisory
metadata, and deterministic tests plus catalog quality checks pass.
````

## Status record

````markdown
# Status M37 - Dropbox Stone-Derived OpenAPI Advisory Artifact

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Formalize the existing Dropbox core advisory overlay as a reviewed
Stone-derived OpenAPI-shaped subset without treating Dropbox's official Stone
source as official OpenAPI or adding runtime protocol support.

## Tasks

| Item | State | Notes |
|---|---|---|
| Dropbox source boundary retained | `[+]` | `dropbox-api-stone-spec` remains the official machine-readable source, official OpenAPI remains unavailable, and advisory tests assert protocol `dropbox-stone`. |
| Dropbox advisory overlay metadata tightened | `[+]` | Regenerated the existing overlay with explicit non-official OpenAPI, official machine-spec derivation, `dropbox-stone` source protocol, and reviewed source refs. |
| Artifact registry provenance updated | `[+]` | Registry generation now records machine-spec-derived advisory overlay metadata for Dropbox while keeping artifact kind `advisory-overlay`. |
| Advisory wording and tests reconciled | `[+]` | Advisory endpoint overlay wording now allows non-human-doc derivation, and tests preserve registered overlay metadata. |
| Verification and review completed | `[+]` | Ran the requested test, vet, diff, catalog check, advisory, inspect, and stats commands; cache-backed Dropbox advisory output shows the overlay only as `advisory-overlay` metadata with `source_protocol=dropbox-stone` and `official_openapi=false`. |

## Boundary Checks

- Do not classify Dropbox as official OpenAPI.
- Do not build a general Stone-to-OpenAPI converter.
- Do not execute Dropbox API operations.
- Do not resolve credentials, fetch tokens, sign requests, or choose runtime
  accounts.
- Do not change `../udon` or add runtime protocol support.
````
