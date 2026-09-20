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
