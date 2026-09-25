# Retired milestone M09 - Catalog Quality Gates

**Milestone.** M09
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M09.md
**Source specification.** tabilet/memory-bank/milestone.md#m09---catalog-quality-gates
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M09 - Catalog Quality Gates

**Goal.** Add static, offline quality checks that help maintain catalog
metadata as provider count, spec references, and overlays grow.

**Scope.**

- Define static quality rules for missing source notes, duplicate aliases,
  stale verification dates, missing auth/security status, unresolved overlay
  provider/spec references, and inconsistent user OpenAPI need.
- Add reusable catalog quality report helpers with deterministic output.
- Add CLI reporting, likely `apitools catalog check`, with text and JSON
  output.
- Keep M9 checks offline by default; do not probe remote URLs or fetch provider
  documents.
- Add tests for passing catalog data, failing records, deterministic report
  order, and CLI exit codes.
- Update README development commands if a new check command is introduced.
- Preserve existing catalog APIs unless any additive exported API is needed for
  quality reports.

**Acceptance.** Maintainers can run a deterministic offline catalog quality
check that catches common metadata regressions without network access, provider
API calls, credential handling, or runtime compatibility claims.
````

## Status record

````markdown
# Status M09 - Catalog Quality Gates

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Add static, offline quality checks that catch common catalog metadata
regressions as provider count, source references, and security overlays grow.

The feature should answer these questions without network access:

- Which catalog records are missing required source or provenance notes?
- Which providers have stale verification dates or missing security status?
- Which overlays reference missing providers, specs, schemes, or operations?
- Which duplicate aliases or inconsistent classifications could confuse
  downstream resolution?

## Quality Gate Direction

M9 checks are offline by default. Remote URL probing, freshness fetching, or
provider document downloads are out of scope for this milestone.

## Tasks

| Item | State | Notes |
|---|---|---|
| Quality rules defined | `[+]` | Static rules cover source notes, duplicate lookup conflicts, stale verification dates, security status, overlay references, operation targets when known operations are supplied, and user OpenAPI need consistency. |
| Quality report types added | `[+]` | Added deterministic catalog quality report helpers with severity, provider/candidate/overlay/spec references, field, code, and message metadata. |
| Built-in catalog check added | `[+]` | Added `BuiltInCatalogQualityReport` to check built-in providers, candidates, classifications, and overlays together. |
| CLI check command added | `[+]` | Added `apitools catalog check` with text and JSON output. |
| CLI exit policy implemented | `[+]` | Error-level findings return nonzero; warning-only reports return zero for inspection. |
| Tests added | `[+]` | Tests cover clean built-in catalog, failing records, stale dates, overlay mismatch, unresolved operation matches, deterministic report order, JSON output, and CLI exit codes. |
| README updated | `[+]` | Added the quality gate command to CLI, catalog, Go usage, and development docs. |
| Verification completed | `[+]` | Ran apitools tests/vet/diff checks, catalog CLI smoke checks, and sibling openudon/udon tests for exported API changes. |

## Boundary Checks

- Do not perform network URL probes by default.
- Do not download provider documents.
- Do not execute API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests.
- Do not turn quality warnings into runtime policy decisions.
````
