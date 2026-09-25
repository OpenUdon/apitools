# Retired milestone M21 - Provider Node Expansion Batch 4

**Milestone.** M21
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M21.md
**Source specification.** tabilet/memory-bank/milestone.md#m21---provider-node-expansion-batch-4
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M21 - Provider Node Expansion Batch 4

**Goal.** Add another frozen batch of popular provider-node services to the
core catalog before starting non-OpenAPI converter work.

**Scope.**

- Freeze Batch 4 before implementation: Intercom, Chargebee, Webflow,
  ActiveCampaign, BambooHR, Contentful, Customer.io, Eventbrite, Freshdesk,
  Mailgun, Postmark, and Todoist.
- Treat each service as a separate task and commit unit once curation begins.
- Use `../n8n/packages/nodes-base/nodes` as a priority signal only, not as
  runtime compatibility evidence.
- For each service, try official OpenAPI, Swagger, OpenAPI index, Discovery, or
  other official machine-readable sources first; save downloaded official
  artifacts under ignored catalog cache paths and register them in SQLite when
  saved.
- If no official OpenAPI exists, record official human docs and build tracked
  docs-derived endpoint overlay assets plus builders only when official docs
  expose enough endpoint instructions for a useful reviewed subset.
- Add or update candidate records, durable provider entries, spec references,
  security classifications, and security overlays with source authority,
  verified dates, source notes, source refs, license notes, quirks, and user
  OpenAPI need.
- Keep Google Discovery, Smithy, Stone, and GraphQL converter work deferred
  until provider-node catalog coverage is substantially broader.
- Avoid n8n runtime compatibility claims, provider API execution, credential
  handling, request signing, account selection, or automatic metadata
  promotion.

**Acceptance.** Batch 4 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.
````

## Status record

````markdown
# Status M21 - Provider Node Expansion Batch 4

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Continue broadening provider-node catalog coverage before starting
non-OpenAPI converter work. M21 focuses on popular REST-oriented provider nodes
with source-backed metadata and auth/security classification.

## Frozen Batch

Batch 4 services:

- Intercom
- Chargebee
- Webflow
- ActiveCampaign
- BambooHR
- Contentful
- Customer.io
- Eventbrite
- Freshdesk
- Mailgun
- Postmark
- Todoist

`../n8n/packages/nodes-base/nodes` is a priority signal only. It is not runtime
compatibility evidence, and no n8n behavior claims should be recorded.

## Per-Service Curation Loop

For each service:

1. Try official OpenAPI, Swagger, OpenAPI index, Discovery, or other official
   machine-readable sources first.
2. Save downloaded official review artifacts under ignored
   `catalog-openapi-cache/openapi/` or
   `catalog-openapi-cache/google-discovery/` when a document is retrieved.
3. Register saved official artifact paths in `catalog-openapi-cache/cache.sqlite`
   when artifacts are saved.
4. Review auth/security completeness from the official spec or docs.
5. Add a candidate record if missing.
6. Promote the provider into durable catalog metadata only after source review.
7. Add source-backed security classification or security overlay metadata when
   upstream security is incomplete, absent, non-OpenAPI, or docs-derived.
8. If no official OpenAPI exists, build tracked docs-derived advisory overlay
   assets and service-specific builder programs only when endpoint/security
   metadata is useful for OpenAPI-only consumers.
9. Record source authority, source notes, source refs, verified dates, license
   notes, quirks, user OpenAPI need, and auth/security completeness.
10. Run focused tests and `go run ./cmd/apitools catalog check --as-of
    2026-05-19` before committing each service row.

## Tasks

| Item | State | Notes |
|---|---|---|
| Batch list frozen | `[+]` | Frozen Batch 4 services are Intercom, Chargebee, Webflow, ActiveCampaign, BambooHR, Contentful, Customer.io, Eventbrite, Freshdesk, Mailgun, Postmark, and Todoist. |
| Intercom curated | `[+]` | Added Intercom candidate/provider metadata, official API v2.15 OpenAPI reference, REST/auth/tooling docs refs, present-incomplete auth review overlay for the parseable-openapi-invalid official artifact, ignored official review artifact registration, and focused checks. |
| Chargebee curated | `[+]` | Added Chargebee candidate/provider metadata, official API v2 Product Catalog v2 OpenAPI reference, API docs refs, present-incomplete Basic auth review overlay, ignored official review artifact registration, and focused checks. |
| Webflow curated | `[+]` | Added Webflow candidate/provider metadata, official API/Data API/auth docs refs, OpenAPI-unavailable classification, bearer-token advisory auth overlay, and focused checks. |
| ActiveCampaign curated | `[+]` | Added ActiveCampaign candidate/provider metadata, official API v3/auth/contact docs refs, OpenAPI-unavailable classification, Api-Token advisory auth overlay, and focused checks. |
| BambooHR curated | `[+]` | Added BambooHR candidate/provider metadata, official getting-started/API-details docs refs, OpenAPI-unavailable classification, Basic API-key advisory auth overlay, and focused checks. |
| Contentful curated | `[+]` | Added Contentful candidate/provider metadata, official Management/Delivery/auth docs refs, OpenAPI-unavailable classification, bearer-token advisory auth overlay, and focused checks. |
| Customer.io curated | `[+]` | Added Customer.io candidate/provider metadata, official App/Track/Pipelines OpenAPI references, API docs refs, present-incomplete mixed bearer/basic auth review overlay, ignored official review artifact registrations, and focused checks. |
| Eventbrite curated | `[+]` | Added Eventbrite candidate/provider metadata, official platform API/auth docs refs, OpenAPI-unavailable classification, OAuth bearer advisory auth overlay, and focused checks. |
| Freshdesk curated | `[+]` | Added Freshdesk candidate/provider metadata, official API/support docs refs, OpenAPI-unavailable classification, Basic API-key advisory auth overlay, and focused checks. |
| Mailgun curated | `[+]` | Added Mailgun candidate/provider metadata, official Mailgun API OpenAPI reference, API docs ref, present-incomplete Basic auth review overlay, ignored official review artifact registration, and focused checks. |
| Postmark curated | `[+]` | Added Postmark candidate/provider metadata, official API/API Explorer refs, OpenAPI-unavailable classification, server/account token advisory auth overlay, and focused checks. |
| Todoist curated | `[+]` | Added Todoist candidate/provider metadata, official REST/auth docs refs, OpenAPI-unavailable classification, bearer-token advisory auth overlay, and focused checks. |
| Human-doc endpoint overlays completed | `[+]` | Added tracked docs-derived endpoint advisory overlays and a deterministic batch builder for ActiveCampaign, BambooHR, Contentful, Eventbrite, Freshdesk, Postmark, Todoist, and Webflow; registered the advisory artifacts and recorded built overlay decisions. |
| Batch quality review completed | `[+]` | Ran full verification and code quality review after the final provider row: Go tests/vet, diff checks, catalog quality, catalog stats, refresh-report, advisory smoke checks, and provider inspect spot checks passed; review found and fixed M21 human-doc verification dates before commit. |

## Boundary Checks

- Do not treat n8n node presence as runtime compatibility evidence.
- Do not vendor downloaded official specs, Discovery documents, or SQLite
  caches into the public repository.
- Do not execute API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not treat advisory overlays as provider truth.
- Do not start Google Discovery, Smithy, Stone, or GraphQL converter work in
  this milestone.
````
