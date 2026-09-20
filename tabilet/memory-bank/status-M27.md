# Status M27 - Provider Node Expansion Batch 10

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Continue broadening provider-like n8n node catalog coverage before standalone
converter work. M27 focuses on communications and social services, then runs
the first embedded backlog checkpoint.

## Frozen Batch

Batch 10 services:

- Facebook
- Facebook Lead Ads
- LinkedIn
- Line
- Mattermost
- MessageBird
- Matrix
- Rocket.Chat
- Twake
- Twist
- Twitter
- WhatsApp

`../n8n/packages/nodes-base/nodes` is a priority signal only. It is not runtime
compatibility evidence, and no n8n behavior claims should be recorded.

## Per-Service Curation Loop

For each service, try official machine-readable sources first, save and
register ignored official review artifacts when retrieved, review
auth/security completeness, add candidate/provider metadata, add security
classification or overlays when needed, generate docs-derived advisory overlays
only for useful REST-shaped subsets, record provenance fields, and run focused
tests plus `go run ./cmd/apitools catalog check --as-of 2026-05-19` before
committing each service row.

## Tasks

| Item | State | Notes |
|---|---|---|
| Batch list frozen | `[+]` | Frozen Batch 10 services are Facebook, Facebook Lead Ads, LinkedIn, Line, Mattermost, MessageBird, Matrix, Rocket.Chat, Twake, Twist, Twitter, and WhatsApp. |
| Facebook curated | `[+]` | Added candidate/provider metadata, Graph API docs refs, bearer auth overlay, and `facebook-graph-api-overlay.json`. |
| Facebook Lead Ads curated | `[+]` | Added candidate/provider metadata, Lead Ads retrieval/webhook docs refs, bearer auth overlay, and `facebook-lead-ads-overlay.json`. |
| LinkedIn curated | `[+]` | Added candidate/provider metadata, Microsoft Learn docs refs, OAuth bearer auth overlay, and `linkedin-api-overlay.json`; no stable public official OpenAPI found. |
| Line curated | `[+]` | Added official `line/line-openapi` Messaging API OpenAPI ref, saved ignored review artifact, registry row, and complete bearer auth classification. |
| Mattermost curated | `[+]` | Added official API v4 OpenAPI ref, saved ignored review artifact, registry row, and complete bearer auth classification; saved spec is parseable but strict-invalid. |
| MessageBird curated | `[+]` | Added candidate/provider metadata, official API docs refs, AccessKey auth overlay, and `messagebird-api-overlay.json`. |
| Matrix curated | `[+]` | Added official Client-Server API OpenAPI ref, saved ignored review artifact, registry row, and complete access-token auth classification. |
| Rocket.Chat curated | `[+]` | Added official OpenAPI module refs for Authentication and Messaging, saved ignored review artifacts, registry rows, and combined `X-Auth-Token`/`X-User-Id` auth review overlay. |
| Twake curated | `[+]` | Added candidate/provider metadata, official Developers API docs refs, Basic auth overlay, and `twake-developers-api-overlay.json`. |
| Twist curated | `[+]` | Added candidate/provider metadata, official API v3 docs refs, bearer auth overlay, and `twist-api-v3-overlay.json`. |
| Twitter curated | `[+]` | Added official X API v2 OpenAPI ref, saved ignored review artifact, registry row, and complete auth classification; saved spec is parseable but strict-invalid. |
| WhatsApp curated | `[+]` | Added candidate/provider metadata, official Cloud API docs refs, bearer auth overlay, and `whatsapp-cloud-api-overlay.json`. |
| Backlog checkpoint completed | `[+]` | Confirmed all 12 frozen Batch 10 n8n priority paths exist and are now represented. Remaining provider-like roots stay queued in M28-M35 status files; TimeSaved remains the only recorded false-positive-style provider row from prior batches. |
| Batch quality review completed | `[+]` | Catalog tests, full tests, vet, diff checks, catalog check/stats, inspect loop, and final code review completed before commit. |

## Boundary Checks

- Do not treat n8n node presence as runtime compatibility evidence.
- Do not vendor downloaded official specs, Discovery documents, or SQLite
  caches into the public repository.
- Do not execute API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not treat advisory overlays as provider truth.
- Do not start standalone converter work in this milestone.
