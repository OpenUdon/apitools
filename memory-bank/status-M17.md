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
