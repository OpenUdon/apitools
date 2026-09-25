# Retired milestone M20 - Human-Docs Overlay Completion And Catalog Stats

**Milestone.** M20
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M20.md
**Source specification.** tabilet/memory-bank/milestone.md#m20---human-docs-overlay-completion-and-catalog-stats
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M20 - Human-Docs Overlay Completion And Catalog Stats

**Goal.** Close the known human-docs endpoint-overlay gap and make catalog
coverage statistics directly reproducible from the CLI.

**Scope.**

- Triage the remaining human-docs-primary providers without endpoint overlays:
  Databricks, Jenkins, Linear, Monday.com, Sentry, Splunk, and Telegram.
- Build tracked OpenAPI-shaped advisory endpoint overlays for REST-shaped
  services where official docs support a useful reviewed subset.
- Explicitly document no-overlay decisions where OpenAPI-shaped output would
  misrepresent the provider API model.
- Add `catalog stats` output for primary provider protocol counts, local
  artifact registry counts, refresh validation buckets, and byte totals.
- Document non-OpenAPI converter priority for Smithy JSON, Google Discovery,
  Dropbox Stone, OpenAPI indexes, human docs, and GraphQL-first services.
- Document the next catalog expansion queue and source-first acceptance rules
  before adding another provider batch.
- Preserve boundaries: no provider API execution, credential resolution,
  request signing, account selection, or treating advisory overlays as provider
  truth.

**Acceptance.** Maintainers can see catalog coverage and refresh artifact
status with one offline command, every human-docs-primary provider has either a
tracked advisory endpoint overlay or an explicit no-overlay decision, and
future non-OpenAPI conversion and provider expansion work has source-first
guidance.
````

## Status record

````markdown
# Status M20 - Human-Docs Overlay Completion And Catalog Stats

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Finish the known human-docs overlay triage, make catalog coverage statistics
repeatable, and record the next non-OpenAPI and expansion direction before
adding more providers.

## Direction

REST-shaped human-docs providers can receive small OpenAPI-shaped advisory
endpoint overlays when official docs expose enough endpoint instructions. The
overlays remain review artifacts only. GraphQL-first providers should wait for
GraphQL-aware metadata instead of receiving a misleading single-endpoint REST
wrapper.

## Tasks

| Item | State | Notes |
|---|---|---|
| Catalog stats added | `[+]` | Added `apitools catalog stats` with provider protocol counts, artifact registry kind counts, refresh validation buckets, and byte totals. |
| REST-shaped overlay gaps closed | `[+]` | Added tracked advisory overlays for Databricks, Jenkins, Sentry, Splunk, and Telegram, with registry entries and a shared M20 overlay builder. |
| No-overlay decisions recorded | `[+]` | Recorded Linear and Monday.com as GraphQL-first no-overlay cases until GraphQL-aware classification/introspection work exists. |
| Overlay structural tests added | `[+]` | Added a test that loads tracked advisory overlay artifacts into the operation index and rejects empty/invalid operation coverage. |
| Non-OpenAPI converter notes documented | `[+]` | Added public notes for Discovery, Smithy, Stone, OpenAPI index, human docs, and GraphQL-first follow-up direction. |
| Expansion queue documented | `[+]` | Added source-first criteria for the next provider expansion batch instead of adding unfrozen services ad hoc. |
| Verification completed | `[+]` | Ran Go tests/vet, diff check, catalog quality, catalog stats/specs/help smoke checks, overlay structure checks, and sibling `openudon`/`udon` tests. |

## Boundary Checks

- Do not execute provider API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not treat advisory endpoint overlays as official provider OpenAPI.
- Do not force GraphQL APIs into REST-shaped endpoint overlays.
````
