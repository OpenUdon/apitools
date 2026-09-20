# Status M22 - Provider Node Expansion Batch 5

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Continue broadening provider-node catalog coverage before starting
non-OpenAPI converter work. M22 focuses on REST-oriented provider nodes with
source-backed metadata and auth/security classification.

## Frozen Batch

Batch 5 services:

- Acuity Scheduling
- Bitly
- Brevo
- Clockify
- Coda
- DeepL
- Figma
- Ghost
- Harvest
- Help Scout
- Iterable
- Mailjet

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
| Batch list frozen | `[+]` | Frozen Batch 5 services are Acuity Scheduling, Bitly, Brevo, Clockify, Coda, DeepL, Figma, Ghost, Harvest, Help Scout, Iterable, and Mailjet. |
| Acuity Scheduling curated | `[+]` | Added official-docs provider metadata, Basic-auth security overlay, and docs-derived endpoint overlay for a focused API v1 subset. |
| Bitly curated | `[+]` | Added official OpenAPI v4 metadata, saved ignored review artifact, and complete bearer-auth classification. |
| Brevo curated | `[+]` | Added official OpenAPI v3 metadata, saved ignored review artifact, and complete api-key/partner-key auth classification. |
| Clockify curated | `[+]` | Added official-docs provider metadata, X-Api-Key/X-Addon-Token security overlay, and docs-derived endpoint overlay for a focused API v1 subset. |
| Coda curated | `[+]` | Added official API v1 OpenAPI metadata, saved ignored review artifact, and complete bearer-auth classification. |
| DeepL curated | `[+]` | Added official OpenAPI repository metadata, saved ignored review artifact, and complete auth-header classification. |
| Figma curated | `[+]` | Added official rest-api-spec OpenAPI metadata, saved ignored review artifact, and complete PAT/plan-token/OAuth classification. |
| Ghost curated | `[+]` | Added official-docs provider metadata, Admin/Content API security overlay, and docs-derived Admin API endpoint overlay. |
| Harvest curated | `[+]` | Added official-docs provider metadata, bearer-auth security overlay, and docs-derived endpoint overlay with Harvest-Account-ID as request metadata. |
| Help Scout curated | `[+]` | Added official-docs provider metadata, Inbox/Docs API security overlay, and docs-derived Inbox API endpoint overlay. |
| Iterable curated | `[+]` | Added official-docs provider metadata, Api-Key security overlay, and docs-derived endpoint overlay for a focused REST subset. |
| Mailjet curated | `[+]` | Added official-docs provider metadata, Basic-auth security overlay, and docs-derived endpoint overlay for focused REST/send endpoints. |
| Human-doc endpoint overlays completed | `[+]` | Built and registered docs-derived endpoint overlays for Acuity Scheduling, Clockify, Ghost, Harvest, Help Scout, Iterable, and Mailjet; official-OpenAPI providers did not need endpoint overlays. |
| Batch quality review completed | `[+]` | Reviewed M22 diff; fixed Harvest account header so it is request metadata rather than a security credential. Verified `go test -count=1 ./...`, `go vet ./...`, `git diff --check`, and `go run ./cmd/apitools catalog check --as-of 2026-05-19`. |

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
