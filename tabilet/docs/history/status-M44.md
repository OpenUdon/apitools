# Retired milestone M44 - Google Discovery Service Expansion

**Milestone.** M44
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M44.md
**Source specification.** tabilet/memory-bank/milestone.md#m44---google-discovery-service-expansion
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M44 - Google Discovery Service Expansion

**Goal.** Expand first-class Google Discovery source coverage for high-priority
Google services represented in `../n8n`, using official Discovery documents as
cataloged native API metadata.

**Scope.**

- Freeze the service batch: BigQuery, Cloud Storage, Docs, Slides, Chat,
  YouTube, Tasks, Admin SDK, Analytics, Translate, Cloud Natural Language,
  Books, Business Profile, and Firestore where official Discovery coverage is
  available.
- Add or update provider entries, aliases, categories, official Discovery spec
  references, source notes, license notes, verified dates, quirks, and user
  OpenAPI need classifications for each service.
- Save official Discovery review artifacts under
  `catalog-openapi-cache/google-discovery/` when appropriate and register
  artifact paths for offline refresh/review.
- Exercise `github.com/OpenUdon/googlediscovery` native parsing against
  selected service fixtures covering nested resources, inherited parameters,
  media upload, OAuth scopes, and schema shapes used by the batch.
- Keep Discovery artifacts classified as `google-discovery`; do not treat them
  as official OpenAPI and do not reintroduce Discovery-to-OpenAPI conversion in
  `apitools`.
- Preserve metadata-only boundaries: no Google API execution, token fetching,
  credential resolution, project/account selection, or request signing.

**Acceptance.** The Google batch is discoverable through catalog/provider/spec
inspection with official Discovery provenance, native parser coverage protects
the batch's Discovery shapes, catalog quality passes, and no API execution,
credential resolution, token fetching, project selection, or account selection
is introduced.
````

## Status record

````markdown
# Status M44 - Google Discovery Service Expansion

| Item | State | Notes |
|---|---|---|
| Batch source review | `[+]` | Confirmed official Discovery document URLs, versions, and revisions for BigQuery, Cloud Storage, Docs, Slides, Chat, YouTube, Tasks, Admin SDK Directory, Analytics Data, Translate v2, Cloud Natural Language, Books, Business Profile Business Information and Account Management, and Firestore. |
| Provider/catalog rows | `[+]` | Added durable provider entries, aliases, categories, official Discovery spec references, source notes, license notes, verified dates, quirks, and user OpenAPI need classifications for the batch. |
| Review artifacts | `[+]` | Saved selected official Discovery documents under ignored `catalog-openapi-cache/google-discovery/` paths and registered them for offline review; Business Profile has two registered current Discovery artifacts. |
| Native parser coverage | `[+]` | Extended `googlediscovery` catalog artifact coverage so downloaded batch artifacts remain native Discovery-shaped and parse through the wrapper when present. |
| Security metadata | `[+]` | Added source-backed OAuth scope overlays as metadata only, without token fetching, credential resolution, account selection, or project selection. |
| Catalog/reporting checks | `[+]` | Verified provider lookup, `catalog specs`, `catalog advisory`, `catalog inspect`, `catalog stats`, `catalog security-report`, `catalog refresh-report`, and catalog quality output for the batch. |
| Documentation and memory bank | `[+]` | Product, architecture, and tech-stack boundaries already cover native Google Discovery parsing; this status file records the M44 service expansion without changing the broader contract. |
| Verification | `[+]` | `go test ./catalog`, `go test ./googlediscovery`, `go test ./...`, `go vet ./...`, `go run ./cmd/apitools catalog check`, catalog CLI smoke checks, and diff checks pass. |

## Boundary Notes

- Google Discovery documents remain native `google-discovery` source metadata,
  not official OpenAPI and not an OpenAPI import contract.
- `apitools` must not execute Google API operations, resolve credentials, fetch
  tokens, sign requests, or choose Google accounts or projects.
- `../n8n` is only a priority signal for service selection.
````
