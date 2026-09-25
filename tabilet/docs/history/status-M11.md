# Retired milestone M11 - Catalog Spec Refresh Tooling

**Milestone.** M11
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M11.md
**Source specification.** tabilet/memory-bank/milestone.md#m11---catalog-spec-refresh-tooling
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M11 - Catalog Spec Refresh Tooling

**Goal.** Make known-spec refresh and review artifact registration more
repeatable for maintainers without expanding the catalog by hand-running
one-off commands.

**Scope.**

- Add a metadata-only helper or CLI command that lists known catalog spec
  references eligible for refresh, including provider ID, spec ref ID, kind,
  URL, source authority, current verified date, and saved artifact path when
  registered.
- Add a bounded refresh/download path for selected catalog spec references that
  uses existing safe HTTP(S) download rules, unsafe-host rejection, size limits,
  redirects, and validation where applicable.
- Save refreshed official OpenAPI/Swagger/Discovery artifacts under ignored
  catalog cache directories and update cache registry paths instead of storing
  duplicate document blobs.
- Emit a deterministic refresh report with changed artifact paths, validation
  status, content hash or revision evidence when available, and required manual
  catalog metadata follow-ups.
- Keep refresh tooling opt-in and local; do not run network probes from
  `catalog check`.
- Do not automatically promote metadata, overwrite tracked advisory overlays,
  execute operations, resolve credentials, sign requests, or choose provider
  accounts.

**Acceptance.** Maintainers can refresh selected known spec references through
a safe, deterministic local workflow that records ignored review artifacts and
actionable metadata follow-ups while keeping catalog quality checks offline.
````

## Status record

````markdown
# Status M11 - Catalog Spec Refresh Tooling

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Make known catalog spec refresh and artifact registration repeatable for
maintainers. Refresh tooling should help review official spec references and
local cache artifacts, but it must not replace manual catalog metadata review or
turn offline quality checks into network operations.

## Direction

`catalog check` remains offline. Refresh tooling is opt-in and local, uses the
same safe HTTP(S) fetch constraints as normal imports, and records review
artifacts under ignored cache paths rather than vendoring upstream documents.

## Tasks

| Item | State | Notes |
|---|---|---|
| Refresh inventory shape defined | `[+]` | Added `catalog.RefreshableSpecReference` and artifact join metadata for provider ID, spec ref ID, kind, URL, source authority, verified date, and registered artifact path when known. |
| Refresh report types added | `[+]` | Added `CatalogSpecRefreshReport` and result metadata for selected refs, download status, validation status, saved path, SHA-256/bytes evidence, extracted metadata, and manual follow-ups. |
| Known spec listing helper added | `[+]` | Added offline `catalog.BuiltInRefreshableSpecReferences` and `catalog.RefreshableSpecReferences` helpers that omit human-docs refs and never perform network access. |
| CLI listing command added | `[+]` | Added `apitools catalog specs` with text/JSON output and optional cache path join for registered artifact paths. |
| Selected refresh command added | `[+]` | Added opt-in `apitools catalog refresh --provider <provider> [--spec <ref>]`, using the existing safe bounded HTTP(S) download path, timeout, size limit, redirect safety, and unsafe-host rejection. |
| Artifact registration updated | `[+]` | Refresh saves official artifacts under ignored `openapi/` or `google-discovery/` cache paths and registers SQLite `content_path`/artifact rows without storing duplicate document blobs. |
| Validation behavior defined | `[+]` | Refresh validates OpenAPI/Swagger roots, validates Discovery/index artifacts as structured documents, and records non-OpenAPI machine specs as review artifacts without pretending they are OpenAPI. |
| No automatic promotion enforced | `[+]` | Refresh reports manual follow-ups and does not edit provider metadata, tracked advisory overlays, verified dates, or security classifications automatically. |
| Tests added | `[+]` | Added coverage for listing, selected refresh command, safe-download host rejection, JSON/text output, registry path updates, validation status, and unchanged offline `catalog check` behavior. |
| README and memory updated | `[+]` | Updated README, catalog curation docs, and tech-stack memory for `catalog specs`, selected `catalog refresh`, ignored artifact paths, SQLite path registration, and offline `catalog check` separation. |
| Verification completed | `[+]` | Ran `go test ./...`, `go vet ./...`, `git diff --check`, catalog CLI smoke checks (`catalog check`, `catalog specs`, `catalog refresh --help`, `search --help`, `import --help`), and sibling `go test ./...` in `../openudon` and `../udon`. Milestone review found and fixed registered-artifact path preservation for refresh selection. |

## Boundary Checks

- Do not run network probes from `catalog check`.
- Do not download every provider by default.
- Do not vendor refreshed official specs into the public repository.
- Do not overwrite tracked advisory overlays automatically.
- Do not execute API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not automatically promote metadata after a refresh.
````
