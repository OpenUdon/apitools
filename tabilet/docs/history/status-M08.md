# Retired milestone M08 - OpenUdon Intent Advisory Integration

**Milestone.** M08
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M08.md
**Source specification.** tabilet/memory-bank/milestone.md#m08---openudon-intent-advisory-integration
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M08 - OpenUdon Intent Advisory Integration

**Goal.** Consume `apitools/catalog` from `../openudon` so intent authoring can
surface optional provider/spec/security advice without making catalog metadata
authoritative.

**Scope.**

- Add an OpenUdon adapter that resolves provider catalog hints from
  `apitools/catalog`.
- Improve optional `intent.hcl` authoring or review guidance with provider ID,
  known spec reference, user OpenAPI need, auth/security status, and overlay
  source notes.
- Preserve explicit project/user OpenAPI inputs as higher precedence than
  built-in catalog metadata.
- Keep catalog advice optional and fail-open when no catalog match exists.
- Avoid runtime credential binding, account selection, release gate decisions,
  workflow execution, or trusted-runner behavior changes.
- Keep `apitools` responsible for metadata and `openudon` responsible for
  workflow authoring/review behavior.

**Acceptance.** OpenUdon can surface catalog-derived intent advice when useful,
but existing explicit user inputs and non-catalog workflows continue to work
without depending on built-in catalog data.
````

## Status record

````markdown
# Status M08 - OpenUdon Intent Advisory Integration

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Use `apitools/catalog` from `../openudon` to improve optional `intent.hcl`
authoring and review guidance. Catalog metadata should help users discover
known provider/spec/security context, but it must not become authoritative
workflow behavior.

The feature should answer these questions in OpenUdon:

- Does a workflow intent mention a provider known to the apitools catalog?
- Is there a known official spec reference or a likely need for user OpenAPI?
- What auth/security status or overlay note should authoring show?
- Did an explicit project/user OpenAPI input override catalog metadata?

## Integration Direction

Start with OpenUdon intent advisory. Reporting-only or apitools-only expansion
is deferred unless implementation discovers a blocker.

## Tasks

| Item | State | Notes |
|---|---|---|
| OpenUdon integration point identified | `[+]` | Review evidence now receives optional catalog advice after OpenAPI candidates/discovery and before inferred data flow. |
| Catalog adapter added in OpenUdon | `[+]` | `internal/workflowintent` exposes a small adapter over `apitools/catalog` without duplicating catalog records downstream. |
| Provider matching implemented | `[+]` | Intent provider and OpenAPI hints match catalog IDs, display names, aliases, and common spec filename stems deterministically. |
| Advisory output added | `[+]` | Review evidence can surface provider ID, spec reference, user OpenAPI need, auth/security status, overlay IDs, and source notes. |
| Explicit input precedence preserved | `[+]` | Intent/project OpenAPI paths are recorded as explicit overrides so catalog specs remain advisory. |
| No-match fallback implemented | `[+]` | Unknown providers and disabled advice render no catalog section and do not affect review generation. |
| OpenUdon tests added | `[+]` | Tests cover advisory presence, no-match fallback, explicit-input precedence, disabled behavior, filename matching, and review rendering. |
| Cross-repo verification completed | `[+]` | Ran OpenUdon, apitools, and udon checks required by AGENTS.md, including public-module OpenUdon tests. |

## Boundary Checks

- Do not move workflow semantics into `apitools`.
- Do not make catalog metadata mandatory for OpenUdon workflows.
- Do not bind credentials, choose accounts, sign requests, or execute workflows.
- Do not change trusted-runner handoff semantics.
- Do not treat catalog advice as a release gate unless a later milestone
  explicitly defines that policy.
````
