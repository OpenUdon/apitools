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
