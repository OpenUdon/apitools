# Status M33 - Provider Node Expansion Batch 16

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Continue broadening provider-like n8n node catalog coverage before standalone
converter work. M33 focuses on operations, webinar, and productivity services
with source-backed metadata and auth/security classification.

## Frozen Batch

Batch 16 services:

- Flow
- GoToWebinar
- HaloPSA
- LoneScale
- Lemlist
- Orbit
- Oura
- ProfitWell
- Quickbase
- SyncroMSP
- Taiga
- Tapfiliate

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
| Batch list frozen | `[+]` | Frozen Batch 16 services are Flow, GoToWebinar, HaloPSA, LoneScale, Lemlist, Orbit, Oura, ProfitWell, Quickbase, SyncroMSP, Taiga, and Tapfiliate. |
| Flow curated | `[+]` | Added candidate/provider rows, bearer auth overlay, and docs-derived overlay `flow-api-overlay.json` for tasks, access tokens, and integration webhooks. |
| GoToWebinar curated | `[+]` | Added candidate/provider rows, OAuth2 auth overlay, and docs-derived overlay `gotowebinar-api-overlay.json` for webinars, sessions, registrants, attendees, panelists, and co-organizers. |
| HaloPSA curated | `[+]` | Added candidate/provider rows, OAuth2 client-credentials auth overlay, and docs-derived overlay `halopsa-api-overlay.json` for tickets, clients, sites, and users. |
| LoneScale curated | `[+]` | Added candidate/provider rows, X-API-KEY auth overlay, and docs-derived overlay `lonescale-api-overlay.json` for lists, list items, current user, and workflows. |
| Lemlist curated | `[+]` | Added candidate/provider rows, Basic API-key auth overlay, and docs-derived overlay `lemlist-api-overlay.json` for campaigns, leads, activities, enrichment, team, and unsubscribes. |
| Orbit curated | `[+]` | Added candidate/provider rows and historical bearer auth metadata; no endpoint overlay because Orbit has shut down and official docs were not reliably reachable. |
| Oura curated | `[+]` | Added candidate/provider rows, registered official Oura API v2 OpenAPI schema, and added a bearer auth overlay because the reviewed schema references BearerAuth without defining the scheme. |
| ProfitWell curated | `[+]` | Added candidate/provider rows, Authorization-header token auth overlay, and docs-derived overlay `profitwell-api-v2-overlay.json` for company settings and metrics. |
| Quickbase curated | `[+]` | Added candidate/provider rows, bearer user-token plus realm-hostname auth overlay, and docs-derived overlay `quickbase-rest-api-overlay.json` for fields, records, reports, and files. |
| SyncroMSP curated | `[+]` | Added candidate/provider rows and registered the official SyncroMSP OpenAPI artifact; strict refresh validation reports it invalid, so it remains a review artifact. |
| Taiga curated | `[+]` | Added candidate/provider rows, bearer token auth overlay, and docs-derived overlay `taiga-api-overlay.json` for projects, epics, issues, tasks, user stories, users, and webhooks. |
| Tapfiliate curated | `[+]` | Added candidate/provider rows, Api-Key auth overlay, and docs-derived overlay `tapfiliate-api-overlay.json` for affiliates, metadata, programs, and program affiliates. |
| Batch quality review completed | `[+]` | Final review completed; verification rerun before the end-of-batch commit. |

## Boundary Checks

- Do not treat n8n node presence as runtime compatibility evidence.
- Do not vendor downloaded official specs, Discovery documents, or SQLite
  caches into the public repository.
- Do not execute API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not treat advisory overlays as provider truth.
- Do not start standalone converter work in this milestone.
