# Retired milestone M17 - Refresh Validation Reconciliation

**Milestone.** M17
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M17.md
**Source specification.** tabilet/memory-bank/milestone.md#m17---refresh-validation-reconciliation
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M17 - Refresh Validation Reconciliation

**Goal.** Preserve official catalog refresh artifacts that are syntactically
parseable OpenAPI/Swagger but fail the current strict semantic validator, while
keeping import behavior strict and review-only.

**Scope.**

- Add curation-only validation statuses for parseable OpenAPI 3.x and Swagger
  2.0 artifacts that fail strict validation.
- Include validation diagnostics in selected refresh and offline refresh-report
  results so maintainers can distinguish strict importer limits from missing or
  malformed documents.
- Continue rejecting empty, malformed, unsupported, or external-ref documents as
  invalid refresh artifacts.
- Register newly saved parseable-invalid official artifacts only as local
  review artifacts; do not promote them as import-ready metadata or provider
  truth.
- Re-run the M16 missing/unregistered refresh set and reconcile any newly saved
  artifact paths into the tracked artifact registry source.
- Preserve boundaries: no provider API execution, credential resolution,
  request signing, account selection, automatic overlay promotion, or automatic
  metadata acceptance.

**Acceptance.** Refresh and refresh-report output classify parseable-but-strict
invalid official specs explicitly, newly saved review artifacts are
reproducibly registered where appropriate, truly invalid artifacts remain
blocked, and the full verification suite remains clean.
````

## Status record

````markdown
# Status M17 - Refresh Validation Reconciliation

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Distinguish official OpenAPI/Swagger refresh artifacts that are parseable but
fail strict semantic validation from documents that are missing, malformed, or
unsupported. M17 keeps normal imports strict while preserving these official
documents as explicit review artifacts.

## Direction

Selected catalog refresh and offline refresh-report may classify saved
OpenAPI/Swagger artifacts as `parseable-openapi-invalid` or
`parseable-swagger-invalid`. Those statuses are curation evidence only; they do
not make artifacts import-ready and do not promote provider metadata.

## Tasks

| Item | State | Notes |
|---|---|---|
| Parseable-invalid statuses added | `[+]` | Added curation-only OpenAPI 3.x and Swagger 2.0 statuses plus refresh result validation diagnostics. |
| Strict import boundary preserved | `[+]` | Kept `specMetadata` and normal import validation strict; only selected refresh/review paths save parseable-invalid artifacts. |
| Refresh review updated | `[+]` | Offline refresh-report now preserves parseable-invalid status and validation errors instead of collapsing these artifacts to plain `invalid`. |
| Tests added | `[+]` | Covered selected refresh and offline review behavior for parseable-invalid OpenAPI/Swagger plus hard rejection of non-spec OpenAPI refresh content. |
| M16 missing refs rerun | `[+]` | Bitbucket, CircleCI, Cloudflare, and Supabase now save as parseable-invalid review artifacts; Elastic still fails the parseable threshold and remains unregistered. |
| Registry source reconciled | `[+]` | Added Bitbucket, CircleCI, Cloudflare, and Supabase artifact rows to the tracked artifact registry source. |
| Verification completed | `[+]` | Ran registry rebuild, refresh report, Go tests/vet, diff check, catalog quality, CLI smoke checks, and sibling `openudon`/`udon` tests. |

## Boundary Checks

- Do not execute provider API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not promote parseable-invalid artifacts as import-ready metadata.
- Do not vendor downloaded official specs or SQLite cache files into the public repository.
- Do not treat validation-failed artifacts as accepted provider truth.
````
