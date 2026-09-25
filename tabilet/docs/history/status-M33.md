# Retired milestone M33 - Provider Node Expansion Batch 16

**Milestone.** M33
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M33.md
**Source specification.** tabilet/memory-bank/milestone.md#m33---provider-node-expansion-batch-16
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M33 - Provider Node Expansion Batch 16

**Goal.** Add a frozen operations, webinar, and productivity provider-node
batch to the core catalog before standalone converter work.

**Scope.**

- Freeze Batch 16 before implementation: Flow, GoToWebinar, HaloPSA,
  LoneScale, Lemlist, Orbit, Oura, ProfitWell, Quickbase, SyncroMSP, Taiga, and
  Tapfiliate.
- Apply the same per-service curation loop, source-first evidence requirements,
  priority-only n8n evidence boundary, and no-execution/no-credential rules as
  M24.

**Acceptance.** Batch 16 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.
````

## Status record

````markdown
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
````
