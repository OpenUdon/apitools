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
