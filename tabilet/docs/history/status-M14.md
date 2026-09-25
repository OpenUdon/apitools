# Retired milestone M14 - Provider Advisory Output

**Milestone.** M14
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M14.md
**Source specification.** tabilet/memory-bank/milestone.md#m14---provider-advisory-output
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M14 - Provider Advisory Output

**Goal.** Provide a directly consumable catalog advisory report for operators
and downstream authoring integrations without making catalog metadata
authoritative.

**Scope.**

- Add reusable catalog advisory report helpers that combine provider metadata,
  spec references, user OpenAPI need, auth/security status, overlay IDs,
  registered artifact paths, and manual follow-ups.
- Add CLI reporting through `apitools catalog advisory [provider] [--cache
  <path>] [--json]`, resolving provider by ID, display name, or alias when a
  provider is supplied.
- Join SQLite artifact paths when a cache exists, but do not create a new cache
  just to render advisory output.
- Keep output deterministic and metadata-only; do not fetch remote documents,
  execute operations, resolve credentials, or apply overlays to upstream specs.
- Preserve explicit downstream user/project OpenAPI inputs as higher
  precedence conceptually; advisory output is guidance, not a release gate.

**Acceptance.** Callers can obtain deterministic provider advisory summaries in
text or JSON, including cache artifact paths when available, while offline
quality checks and existing catalog commands continue to behave as before.
````

## Status record

````markdown
# Status M14 - Provider Advisory Output

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Expose catalog advisory summaries that are easier for operators and downstream
authoring tools to consume, while keeping catalog metadata optional and
non-authoritative.

## Direction

Advisory output combines existing catalog metadata, security classifications,
overlays, spec references, and optional cache artifact paths. It must not fetch
remote documents, apply overlays to upstream specs, execute API operations, or
create cache files just to render a report.

## Tasks

| Item | State | Notes |
|---|---|---|
| Advisory report shape defined | `[+]` | Added reusable provider advisory report/result types for provider ID/name, spec refs, user OpenAPI need, auth status, overlay IDs, artifact paths, source notes, resolved metadata, and manual follow-ups. |
| Advisory builder added | `[+]` | Added deterministic built-in and configurable advisory helpers for all providers or a selected provider, including provider ID/name/alias resolution. |
| Cache artifact join added | `[+]` | Advisory helpers join registered artifact paths from an existing SQLite cache and tolerate a missing cache path without creating a cache. |
| CLI command added | `[+]` | Added `apitools catalog advisory [provider] [--cache <path>] [--json]` with text and JSON output. |
| Advisory output documented | `[+]` | Updated README and memory-bank product/architecture/tech-stack notes for advisory command semantics and metadata-only boundaries. |
| Tests added | `[+]` | Added deterministic ordering, selected provider lookup, missing cache behavior, artifact path joins, JSON/text output, and unknown provider error coverage. |
| Verification completed | `[+]` | Ran `go test ./...`, `go vet ./...`, `git diff --check`, `go run ./cmd/apitools catalog check --as-of 2026-05-19`, CLI smoke checks, and sibling `../openudon` / `../udon` test suites. |

## Boundary Checks

- Do not make advisory output authoritative.
- Do not fetch remote documents.
- Do not create cache files for read-only output.
- Do not execute API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not apply overlays to upstream OpenAPI files.
- Do not change OpenUdon runtime behavior in this milestone.
````
