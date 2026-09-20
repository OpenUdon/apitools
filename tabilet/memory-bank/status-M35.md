# Status M35 - Provider Node Expansion Batch 18 And Reconciliation

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Finish the frozen uncovered provider-like n8n node backlog and reconcile
catalog coverage before any standalone converter milestone is opened.

## Frozen Batch

Batch 18 services:

- Segment
- Strava
- Zoho
- uProc

`../n8n/packages/nodes-base/nodes` is a priority signal only. It is not runtime
compatibility evidence, and no n8n behavior claims should be recorded.

## Per-Service Curation Loop

For each service, try official machine-readable sources first, save and
register ignored official review artifacts when retrieved, review
auth/security completeness, add candidate/provider metadata, add security
classification or overlays when needed, generate docs-derived advisory overlays
only for useful REST-shaped subsets, record provenance fields, and run focused
tests plus `go run ./cmd/apitools catalog check --as-of 2026-05-20` before the
final M35 commit.

## Tasks

| Item | State | Notes |
|---|---|---|
| Batch list frozen | `[+]` | Frozen Batch 18 services are Segment, Strava, Zoho, and uProc. |
| Segment curated | `[+]` | Added source-backed candidate/provider rows from official Segment Public API docs and SDK evidence. Segment's Redocly docs expose embedded OpenAPI 3.0.3 metadata, but no stable standalone official refresh URL was found, so the catalog records human-docs provenance, auth complete, and a no-overlay decision instead of vendoring a docs-derived endpoint overlay. |
| Strava curated | `[+]` | Added source-backed candidate/provider rows and registered the official Strava API v3 Swagger artifact. The artifact is tracked in the registry as strict-refresh invalid while preserving official provenance; auth is complete from the Swagger OAuth2 security definition and Strava authorization docs. |
| Zoho curated | `[+]` | Added source-backed candidate/provider rows and registered official Zoho CRM v8 OpenAPI artifacts from `zoho/crm-oas`. Four artifacts refresh as valid OpenAPI and the users artifact is parseable OpenAPI invalid; auth is complete from the OAuth2 security schemes and operation security metadata. |
| uProc curated | `[+]` | Added source-backed candidate/provider rows from official uProc service and tool pages. No stable official OpenAPI or sufficiently complete endpoint text was found for a useful reviewed REST overlay; auth is present-incomplete because official pages reference API use and email/API-key authentication without enough stable transport detail. |
| Final provider-like backlog reconciliation completed | `[+]` | Reconciled the original 136-root provider-like backlog: M24-M34 account for 132 rows and M35 accounts for the final four frozen rows. False positives or deferred protocol families remain documented rather than promoted: TimeSaved is workflow metadata without a direct external API row, Odoo XML-RPC/JSON-RPC and Emelia GraphQL have no REST overlay, Linear and Monday are older GraphQL no-overlay precedents, Orbit is retained as a historical shutdown row, and PostBin/uProc lack stable official endpoint text for reviewed overlays. |
| Batch quality review completed | `[+]` | Final M35 diff review completed after catalog/provider/security baselines were updated. Focused catalog checks passed, full offline verification passed, ignored downloaded OpenAPI/cache artifacts remained uncommitted, and no standalone converter work was started. |

## Boundary Checks

- Do not treat n8n node presence as runtime compatibility evidence.
- Do not vendor downloaded official specs, Discovery documents, or SQLite
  caches into the public repository.
- Do not execute API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not treat advisory overlays as provider truth.
- Do not start standalone converter work in this milestone.
