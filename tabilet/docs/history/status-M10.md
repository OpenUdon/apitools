# Retired milestone M10 - Catalog Expansion Batch 2

**Milestone.** M10
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M10.md
**Source specification.** tabilet/memory-bank/milestone.md#m10---catalog-expansion-batch-2
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M10 - Catalog Expansion Batch 2

**Goal.** Add the next frozen batch of popular workflow providers to the core
catalog with source-backed spec references, auth/security classification, and
quality-gate coverage.

**Scope.**

- Freeze Batch 2 before implementation: Zendesk, Twilio, SendGrid, Mailchimp,
  Zoom, Okta, ServiceNow, PayPal, Xero, QuickBooks, Linear, Monday.com,
  Pipedrive, Typeform, and Telegram.
- Treat each service as a separate task and commit unit once curation begins.
- Use `../n8n/packages/nodes-base/nodes` as a priority signal only, not as
  runtime compatibility evidence.
- For each service, try official OpenAPI, Swagger, OpenAPI index, or Discovery
  download first; save downloaded review artifacts under ignored
  `catalog-openapi-cache/openapi/` or `catalog-openapi-cache/google-discovery/`
  and register paths in `catalog-openapi-cache/cache.sqlite` when artifacts are
  saved.
- If no official OpenAPI source exists, record official human docs and, when
  those docs expose enough endpoint instructions, generate tracked
  service-specific docs-derived endpoint overlay assets plus builder programs in
  addition to auth/security overlay metadata.
- Add or update candidate records, durable provider entries, spec references,
  security classifications, and security overlays with source authority,
  verified dates, license notes, quirks, source refs, and source notes.
- Run `apitools catalog check`, catalog tests, and required repository checks
  after the batch.
- Avoid vendoring downloaded official specs, executing API operations,
  credential handling, request signing, account selection, or n8n runtime
  compatibility claims.

**Acceptance.** Batch 2 providers are represented in the catalog with durable
source-backed metadata, quality checks pass offline, generated advisory assets
are tracked only when appropriate, downloaded official artifacts remain ignored,
and each service can be inspected without live provider calls or credential
resolution.
````

## Status record

````markdown
# Status M10 - Catalog Expansion Batch 2

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Expand the catalog with a second frozen batch of popular workflow providers.
Each service must follow the same source-backed curation loop used in M6:
official machine-readable sources first, ignored local review artifacts for
downloaded documents, durable catalog metadata only after review, and advisory
overlays only when official OpenAPI is absent or auth/security metadata is
incomplete.

## Frozen Batch

Batch 2 services:

- Zendesk
- Twilio
- SendGrid
- Mailchimp
- Zoom
- Okta
- ServiceNow
- PayPal
- Xero
- QuickBooks
- Linear
- Monday.com
- Pipedrive
- Typeform
- Telegram

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
   upstream security is incomplete, absent, or non-OpenAPI.
8. If no official OpenAPI exists, build tracked docs-derived advisory overlay
   assets and service-specific builder programs only when endpoint/security
   metadata is useful for OpenAPI-only consumers.
9. Record source authority, source notes, source refs, verified dates, license
   notes, quirks, user OpenAPI need, and auth/security completeness.
10. Run focused tests and `go run ./cmd/apitools catalog check --as-of
    2026-05-18` before committing each service row.

## Tasks

| Item | State | Notes |
|---|---|---|
| Batch list frozen | `[+]` | Frozen Batch 2 services are Zendesk, Twilio, SendGrid, Mailchimp, Zoom, Okta, ServiceNow, PayPal, Xero, QuickBooks, Linear, Monday.com, Pipedrive, Typeform, and Telegram. |
| Zendesk curated | `[+]` | Added source-backed Zendesk candidate/provider metadata, official Sunshine Conversations OpenAPI reference, Support API docs references, security classification, support auth review overlay, ignored official review artifact registration, and focused checks. |
| Twilio curated | `[+]` | Added source-backed Twilio candidate/provider metadata, official twilio-oai v2010 OpenAPI reference, complete Basic auth classification, and ignored official review artifact registration. |
| SendGrid curated | `[+]` | Added source-backed SendGrid candidate/provider metadata, official twilio/sendgrid-oai Mail v3 OpenAPI reference, complete bearer API-key auth classification, and ignored official review artifact registration. |
| Mailchimp curated | `[+]` | Added source-backed Mailchimp candidate/provider metadata, official Marketing API docs references, no recorded stable downloadable OpenAPI artifact, API-key/OAuth advisory auth overlay, docs-derived endpoint overlay asset and builder, and focused checks. |
| Zoom curated | `[+]` | Added source-backed Zoom candidate/provider metadata, official zoom/api Swagger 2.0 review artifact registration, current docs references, and present-incomplete OAuth auth review overlay. |
| Okta curated | `[+]` | Added source-backed Okta candidate/provider metadata, official Management OpenAPI current minimal artifact registration, docs references, and complete SSWS/OAuth security classification. |
| ServiceNow curated | `[+]` | Added source-backed ServiceNow candidate/provider metadata, official REST API and per-instance OpenAPI export docs references, no stable public downloadable OpenAPI artifact, Basic/OAuth advisory auth overlay, and docs-derived endpoint overlay asset and builder. |
| PayPal curated | `[+]` | Added source-backed PayPal candidate/provider metadata, official Checkout Orders v2 OpenAPI review artifact registration, docs references, and complete OAuth2 security classification. |
| Xero curated | `[+]` | Added source-backed Xero candidate/provider metadata, official Accounting OpenAPI review artifact registration, docs references, and complete OAuth2 security classification. |
| QuickBooks curated | `[+]` | Added source-backed QuickBooks candidate/provider metadata, official Intuit REST API and OAuth docs references, no stable public downloadable OpenAPI artifact, OAuth bearer advisory auth overlay, and docs-derived endpoint overlay asset and builder. |
| Linear curated | `[+]` | Added source-backed Linear candidate/provider metadata, official GraphQL docs references, OpenAPI-unavailable classification, and bearer-token advisory auth overlay. |
| Monday.com curated | `[+]` | Added source-backed Monday.com candidate/provider metadata, official GraphQL docs references, OpenAPI-unavailable classification, and Authorization-header token advisory auth overlay. |
| Pipedrive curated | `[+]` | Added source-backed Pipedrive candidate/provider metadata, official API v2 OpenAPI review artifact registration, docs reference, and complete API-token/OAuth security classification. |
| Typeform curated | `[+]` | Added source-backed Typeform candidate/provider metadata, official REST API and token/OAuth docs references, no stable public downloadable OpenAPI artifact, bearer-token advisory auth overlay, and docs-derived endpoint overlay asset and builder. |
| Telegram curated | `[+]` | Added source-backed Telegram candidate/provider metadata, official Bot API docs reference, OpenAPI-unavailable classification, and advisory bot-token auth overlay noting path-token limitations. |
| Batch quality review completed | `[+]` | Ran `go test ./...`, `go vet ./...`, `git diff --check`, `go run ./cmd/apitools catalog check --as-of 2026-05-18`, catalog list/inspect smoke checks, search/import help smoke checks, and regenerated registered docs-derived endpoint overlay assets; no exported API changes required sibling tests. |

## Boundary Checks

- Do not treat n8n node presence as runtime compatibility evidence.
- Do not vendor downloaded official specs or Discovery documents into the public
  repository.
- Do not execute API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not treat advisory overlays as provider truth.
````
