# Retired milestone M07 - Overlay Applied Inspection Views

**Milestone.** M07
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M07.md
**Source specification.** tabilet/memory-bank/milestone.md#m07---overlay-applied-inspection-views
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M07 - Overlay Applied Inspection Views

**Goal.** Provide inspection-only views that show how catalog security overlays
would supplement provider metadata without mutating or exporting upstream
OpenAPI documents.

**Scope.**

- Define provenance-labeled inspection view types under `catalog`.
- Combine base catalog security classifications, spec references, and security
  overlays into a read-only advisory view.
- Preserve source labels for overlay-provided schemes, root security,
  operation security, OAuth scopes, and API-key placement.
- Report conflicts such as duplicate scheme names, missing referenced schemes,
  overlay-only additions, and unresolved operation matches.
- Expose inspection output through library helpers and catalog CLI reporting.
- Keep merged views metadata-only; do not write overlay-applied OpenAPI files.
- Do not execute API operations, parse credentials, sign requests, or treat
  overlays as upstream provider truth.

**Acceptance.** Downstream callers and operators can inspect overlay effects and
conflicts deterministically while the upstream spec and local overlay remain
separate provenance-labeled metadata.
````

## Status record

````markdown
# Status M07 - Overlay Applied Inspection Views

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Provide deterministic inspection views that show how catalog security overlays
would supplement provider metadata while keeping upstream specs and local
overlays separate.

The feature should answer these questions without mutating OpenAPI documents:

- Which auth/security metadata comes from upstream catalog classifications?
- Which auth/security metadata comes from a catalog overlay?
- Which overlay additions or corrections would be visible to downstream review?
- Which conflicts or unresolved operation matches need human review?

## Output Shape

M7 should produce inspection-only views. It should not export or write
overlay-applied OpenAPI documents.

## Tasks

| Item | State | Notes |
|---|---|---|
| Inspection view types defined | `[+]` | Added provenance-labeled catalog types for effective security schemes, root security, operation security, source notes, classifications, and conflicts. |
| Overlay application helper added | `[+]` | Added `BuildSecurityInspectionView` and `BuiltInSecurityInspectionView` for read-only advisory views from provider metadata, security classifications, and overlays. |
| Provenance labels preserved | `[+]` | Schemes, root requirements, operation requirements, classification notes, and overlay notes carry explicit provenance and source labels. |
| Conflict reporting added | `[+]` | Reports duplicate scheme names, missing referenced schemes, overlay-only additions, and unresolved operation matches. |
| CLI inspection output added | `[+]` | Added `apitools catalog overlay-view <provider>` with text and JSON output. |
| Tests added | `[+]` | Covered deterministic views, conflict cases, provenance labels, copy safety, known operation resolution, and CLI output. |
| Documentation updated | `[+]` | Updated README plus product, architecture, and tech-stack memory-bank docs for overlay inspection views and CLI smoke command. |
| Verification completed | `[+]` | Passed `go test ./...`, `go vet ./...`, `git diff --check`, catalog/search/import CLI smoke checks, and sibling `../openudon` plus `../udon` tests. |

## Boundary Checks

- Do not write overlay-applied OpenAPI files.
- Do not mutate upstream specs or built-in overlay records.
- Do not execute API operations.
- Do not fetch credentials, tokens, or secrets.
- Do not sign requests.
- Do not treat overlays as upstream provider truth.
````
