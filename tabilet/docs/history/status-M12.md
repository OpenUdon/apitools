# Retired milestone M12 - Existing Catalog Hardening

**Milestone.** M12
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M12.md
**Source specification.** tabilet/memory-bank/milestone.md#m12---existing-catalog-hardening
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M12 - Existing Catalog Hardening

**Goal.** Reduce unresolved status in the existing built-in catalog before
adding another provider batch.

**Scope.**

- Review providers that currently report unknown, needs-verification, or
  advisory-only status: Airtable, Calendly, Linear, Mailchimp, Monday.com,
  OpenWeatherMap, QuickBooks, Salesforce, ServiceNow, Shopify, Telegram, and
  Typeform.
- Re-run the source-backed curation loop for each targeted provider:
  official OpenAPI/Swagger/Discovery first, then official machine-readable
  alternatives, then docs-derived endpoint overlays only when no official
  OpenAPI exists and official docs expose enough endpoint instructions.
- Reconcile ignored official review artifacts, tracked docs-derived advisory
  overlays, tracked overlay builders, and SQLite artifact registrations.
- Tighten source notes, source refs, license notes, verified dates, user
  OpenAPI need, auth/security classification, and security overlays where the
  current catalog is ambiguous.
- Keep `catalog check` offline and avoid adding network probes to quality
  gates.
- Avoid API execution, credential resolution, request signing, automatic
  metadata promotion, or treating advisory overlays as provider truth.

**Acceptance.** Existing ambiguous providers have reviewed source-backed status
and artifact registration where appropriate, offline quality checks remain
clean, and unresolved cases are explicitly documented rather than silently
left as stale `unknown` or `needs-verification` metadata.
````

## Status record

````markdown
# Status M12 - Existing Catalog Hardening

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Clean up ambiguous catalog status before expanding again. M12 should reduce
`unknown`, `needs-verification`, and advisory-only ambiguity for existing
providers while preserving the metadata-only catalog boundary.

## Direction

Each targeted provider follows the same source-backed curation loop used in M6
and M10. Official OpenAPI/Swagger/Discovery is preferred. Docs-derived endpoint
overlays are tracked only when no official OpenAPI exists and official docs are
sufficient for a useful reviewed subset. Downloaded official specs remain
ignored local review artifacts.

## Tasks

| Item | State | Notes |
|---|---|---|
| Hardening set frozen | `[+]` | Target Airtable, Calendly, Linear, Mailchimp, Monday.com, OpenWeatherMap, QuickBooks, Salesforce, ServiceNow, Shopify, Telegram, and Typeform, based on current `unknown`, `needs-verification`, and advisory-only catalog status. |
| Artifact registry reconciled | `[+]` | Re-ran the tracked registry script; `catalog specs --cache catalog-openapi-cache/cache.sqlite` resolves registered official artifact paths, and tracked advisory overlays/builders are enumerated in the registry source. |
| Airtable hardened | `[+]` | Marked official OpenAPI as unavailable after M12 review, added official support/auth docs refs, retained likely user OpenAPI need for base-specific schemas, and kept advisory bearer-token overlay metadata. |
| Calendly hardened | `[+]` | Marked official OpenAPI as unavailable after M12 review, added official PAT/OAuth docs refs, retained likely user OpenAPI need, and kept advisory bearer-token overlay metadata. |
| Linear hardened | `[+]` | Re-review confirmed existing metadata: GraphQL API docs, official OpenAPI unavailable, likely user/generated OpenAPI need for REST-shaped metadata, and bearer-token advisory overlay. |
| Mailchimp hardened | `[+]` | Resolved `needs-verification` to unavailable for stable public official OpenAPI download, retained docs-derived endpoint overlay and API-key/OAuth overlay, and added official OpenAPI/Swagger docs evidence refs. |
| Monday.com hardened | `[+]` | Re-review confirmed existing metadata: GraphQL API docs, official OpenAPI unavailable, likely user/generated OpenAPI need for REST-shaped metadata, and Authorization-header token advisory overlay. |
| OpenWeatherMap hardened | `[+]` | Marked official OpenAPI as unavailable after M12 review, added One Call 3.0 and API-key docs refs, retained likely user/generated OpenAPI need, and kept appid query-key overlay metadata. |
| QuickBooks hardened | `[+]` | Resolved `needs-verification` to unavailable for stable public official OpenAPI download, retained docs-derived endpoint overlay and OAuth overlay, and kept likely user/generated OpenAPI need. |
| Salesforce hardened | `[+]` | Resolved `needs-verification` to unavailable for stable public official OpenAPI download, retained org-generated/user-provided OpenAPI need, docs-derived endpoint overlay, and OAuth bearer overlay. |
| ServiceNow hardened | `[+]` | Resolved `needs-verification` to unavailable for stable public provider-wide OpenAPI download, retained instance-export/user-provided OpenAPI need, docs-derived endpoint overlay, and Basic/OAuth overlay. |
| Shopify hardened | `[+]` | Marked official OpenAPI as unavailable after M12 review, added Admin REST/GraphQL/access-scope refs, retained likely user/generated OpenAPI need, and kept Admin token overlay metadata. |
| Telegram hardened | `[+]` | Re-review confirmed existing metadata: official Bot API docs, official OpenAPI unavailable, path-token auth limitation notes, likely user/generated OpenAPI need, and advisory token overlay. |
| Typeform hardened | `[+]` | Resolved `needs-verification` to unavailable for stable public official OpenAPI download, retained docs-derived endpoint overlay and PAT/OAuth overlay, and added endpoint-reference evidence. |
| Quality and docs updated | `[+]` | Reviewed README/docs/memory impact; no public command or curation workflow docs needed changes because M12 only tightened existing provider metadata. Catalog list now has no provider OpenAPI status left as `unknown` or `needs-verification`. |
| Verification completed | `[+]` | Ran `go test ./...`, `go vet ./...`, `git diff --check`, `catalog check --as-of 2026-05-19`, `catalog list`, `catalog specs --cache catalog-openapi-cache/cache.sqlite`, `catalog refresh --help`, `search --help`, `import --help`, and sibling `go test ./...` in `../openudon` and `../udon`. Milestone review found no boundary drift and no evolution bump was needed. |

## Boundary Checks

- Do not treat advisory overlays as provider truth.
- Do not vendor downloaded official specs or Discovery documents.
- Do not run network probes from `catalog check`.
- Do not execute API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not claim n8n runtime compatibility from node presence.
````
