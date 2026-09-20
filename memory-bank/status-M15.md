# Status M15 - Refresh Review Reports

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Make local refresh artifact review repeatable without rerunning network
downloads. M15 should help maintainers audit whether saved artifacts and cache
registrations match built-in refreshable spec metadata.

## Direction

The report is offline and read-only. It inspects built-in refreshable spec
references, existing SQLite registrations, and files under the catalog cache
directory. It reports evidence and follow-ups but does not promote metadata.

## Tasks

| Item | State | Notes |
|---|---|---|
| Refresh review shape defined | `[+]` | Added reusable report/result types for provider/spec IDs, kind, URL, registered path, saved path, bytes, SHA-256, validation status, verified date, and follow-ups. |
| Offline artifact scanner added | `[+]` | Scanner inspects existing cache registrations and saved cache files without network access or cache creation, including unregistered files whose names match provider/spec metadata. |
| Validation reuse added | `[+]` | Reuses M11 validation categories for OpenAPI/Swagger roots, structured Discovery/index artifacts, and non-OpenAPI machine specs. |
| CLI command added | `[+]` | Added `apitools catalog refresh-report [--cache-dir catalog-openapi-cache] [--cache catalog-openapi-cache/cache.sqlite] [--json]`. |
| Missing/invalid findings added | `[+]` | Reports missing registered files, unregistered refreshable specs, invalid parse/validation status, stale verification dates, unsafe local artifact paths, and manual follow-ups deterministically. |
| Refresh report documented | `[+]` | Updated README/docs/memory for offline report behavior and curation workflow. |
| Tests added | `[+]` | Covers missing cache, missing file, valid OpenAPI, valid structured Discovery/index, non-OpenAPI machine specs, JSON/text output, and deterministic ordering. |
| Verification completed | `[+]` | Ran Go tests/vet, diff check, catalog quality, CLI smoke checks, and sibling tests for exported API changes. |

## Boundary Checks

- Do not download documents.
- Do not run network probes from report generation.
- Do not create cache files for read-only output.
- Do not promote verified dates or catalog metadata automatically.
- Do not overwrite tracked advisory overlays.
- Do not execute API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
