# Retired milestone M36 - Refresh Correction Review

**Milestone.** M36
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M36.md
**Source specification.** tabilet/memory-bank/milestone.md#m36---refresh-correction-review
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M36 - Refresh Correction Review

**Goal.** Reduce known invalid official OpenAPI/Swagger refresh artifacts where
source review supports narrow deterministic correction, while leaving broader
malformed artifacts invalid.

**Scope.**

- Add provider-specific refresh validation corrections for the official
  HighLevel module specs, Strava Swagger, and Spotify Web API OpenAPI artifact.
- Keep SyncroMSP invalid unless a future source-reviewed normalization can
  handle its broader schema-shape issues without broad importer relaxation.
- Preserve raw downloaded artifacts in the ignored cache and surface correction
  notes in live refresh and offline refresh-review outputs.

**Acceptance.** HighLevel, Strava, and Spotify refresh validation reports move
from invalid to valid with correction notes; SyncroMSP remains invalid; catalog
stats reflect the reduced invalid bucket; deterministic tests and offline
quality checks pass.
````

## Status record

````markdown
# Status M36 - Refresh Correction Review

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Apply narrow source-reviewed validation corrections to selected official
OpenAPI/Swagger refresh artifacts that were previously reported as invalid,
without broadening importer behavior for malformed documents.

## Tasks

| Item | State | Notes |
|---|---|---|
| Invalid artifact causes rechecked | `[+]` | Live refresh and direct parser review confirmed HighLevel and Strava were correction candidates, Spotify needed source-reviewed schema cleanup, and SyncroMSP had broader schema-shape issues that should remain invalid. |
| HighLevel refresh corrections implemented | `[+]` | Added provider-specific validation correction for the five official HighLevel module specs: bundle official common error schemas, rewrite sibling common-schema refs, fill empty response descriptions, and add placeholder `items` for incomplete nested arrays. |
| Strava refresh correction implemented | `[+]` | Added provider-specific validation correction for the official Strava Swagger: replace official external model refs with permissive local object placeholders and remove the undeclared root `public` OAuth scope. |
| Spotify refresh correction implemented | `[+]` | Added provider-specific validation correction for the official Spotify Web API OpenAPI: remove the external policy extension ref, drop invalid schema-level boolean `required` flags, and prune `required` names not present in schema properties. |
| SyncroMSP no-overlay decision retained | `[+]` | No correction was added for SyncroMSP because its official document uses broader OpenAPI 3.0-incompatible schema shapes, including JSON Schema-style nullable type arrays and nested schema-shape issues. |
| Refresh metadata and documentation updated | `[+]` | Live refresh and offline review results now expose `correction_notes`; artifact registry statuses were updated for HighLevel, Spotify, and Strava while SyncroMSP remains invalid. Product and architecture memory now describe refresh corrections. |
| Verification and review completed | `[+]` | Full tests, vet, catalog check/stats, refresh-report spot checks, diff checks, and correction diff review passed. Catalog stats now report one invalid artifact, `syncromsp/syncromsp-api-openapi`, as intended. |

## Boundary Checks

- Do not execute provider API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not fetch secondary refs during validation; corrections must be local and
  source-reviewed.
- Do not relax validation globally for unrelated providers.
- Do not vendor raw downloaded provider specs or SQLite cache files into the
  public repository.
````
