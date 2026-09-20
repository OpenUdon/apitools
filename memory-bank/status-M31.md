# Status M31 - Provider Node Expansion Batch 14

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Continue broadening provider-like n8n node catalog coverage before standalone
converter work. M31 focuses on messaging and notification services, then runs
an embedded converter/protocol blocker checkpoint.

## Frozen Batch

Batch 14 services:

- Mocean
- MSG91
- Plivo
- Pushbullet
- Pushcut
- Pushover
- SIGNL4
- sms77
- Vonage
- Zulip
- Gotify
- Emelia

`../n8n/packages/nodes-base/nodes` is a priority signal only. It is not runtime
compatibility evidence, and no n8n behavior claims should be recorded.

## Per-Service Curation Loop

For each service, try official machine-readable sources first, save and
register ignored official review artifacts when retrieved, review
auth/security completeness, add candidate/provider metadata, add security
classification or overlays when needed, generate docs-derived advisory overlays
only for useful REST-shaped subsets, record provenance fields, and run focused
tests plus `go run ./cmd/apitools catalog check --as-of 2026-05-20` before the
single end-of-batch commit.

## Tasks

| Item | State | Notes |
|---|---|---|
| Batch list frozen | `[+]` | Frozen Batch 14 services are Mocean, MSG91, Plivo, Pushbullet, Pushcut, Pushover, SIGNL4, sms77, Vonage, Zulip, Gotify, and Emelia. |
| Mocean curated | `[+]` | Added candidate/provider rows from official docs, auth overlay for `mocean-api-key` and `mocean-api-secret`, and docs-derived overlay `mocean-api-overlay.json` for SMS, voice, and balance endpoints. |
| MSG91 curated | `[+]` | Added candidate/provider rows from official MSG91 docs, auth overlay for `authkey`, and docs-derived overlay `msg91-api-overlay.json` for SMS and OTP endpoints. |
| Plivo curated | `[+]` | Added candidate/provider rows from official Messaging and Voice docs, auth overlay for Auth ID/Auth Token HTTP Basic, and docs-derived overlay `plivo-api-overlay.json` for message and call endpoints. |
| Pushbullet curated | `[+]` | Added candidate/provider rows from official API docs, OAuth bearer auth overlay, and docs-derived overlay `pushbullet-api-overlay.json` for pushes, upload requests, devices, contacts, channels, and subscriptions. |
| Pushcut curated | `[+]` | Added candidate/provider rows from official support docs, API-Key header auth overlay, and docs-derived overlay `pushcut-api-overlay.json` for notifications, devices, and subscriptions. |
| Pushover curated | `[+]` | Added candidate/provider rows from official Message and Client API docs, token/user auth overlay, and docs-derived overlay `pushover-api-overlay.json` for messages, receipts, validation, and sounds. |
| SIGNL4 curated | `[+]` | Added candidate/provider rows and registered the official `signl4-api-v1-openapi` artifact; added auth overlay because the reviewed OpenAPI declares schemes but leaves root security anonymous. |
| sms77 curated | `[+]` | Added candidate/provider rows under current seven.io branding with sms77 aliases, X-Api-Key auth overlay, and docs-derived overlay `sms77-api-overlay.json` for SMS, voice, balance, pricing, status, and contacts. |
| Vonage curated | `[+]` | Added candidate/provider rows and registered the official Vonage Messages API OpenAPI export; n8n's legacy SMS endpoint is recorded only as a priority signal, not provider truth. |
| Zulip curated | `[+]` | Added candidate/provider rows and registered the official Zulip REST API OpenAPI from the Zulip repository. |
| Gotify curated | `[+]` | Added candidate/provider rows and registered the official Gotify server Swagger/OpenAPI artifact from the Gotify repository. |
| Emelia curated | `[+]` | Added candidate/provider rows and Authorization header auth overlay from official GraphQL docs; recorded no endpoint overlay because a generic `POST /graphql` wrapper would hide operation semantics. |
| Converter/protocol checkpoint completed | `[+]` | M31 produced four official machine-readable artifacts (SIGNL4, Vonage, Zulip, Gotify), seven REST-shaped docs-derived overlays, and one GraphQL-only no-overlay decision (Emelia). No SOAP, SDK-only, or AI/vector-specific cluster in this batch justifies blocking M32 for standalone converter work. |
| Batch quality review completed | `[+]` | Final review found and fixed Pushover endpoint overlay security alignment and the sms77 overlay decision label; required tests, vet, diff checks, catalog check/stats, provider inspect loop, and endpoint overlay advisory checks passed before commit. |

## Boundary Checks

- Do not treat n8n node presence as runtime compatibility evidence.
- Do not vendor downloaded official specs, Discovery documents, or SQLite
  caches into the public repository.
- Do not execute API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not treat advisory overlays as provider truth.
- Do not start standalone converter work in this milestone.
