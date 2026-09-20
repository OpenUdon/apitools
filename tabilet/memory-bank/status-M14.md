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
