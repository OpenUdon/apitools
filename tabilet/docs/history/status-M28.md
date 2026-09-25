# Retired milestone M28 - Provider Node Expansion Batch 11

**Milestone.** M28
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M28.md
**Source specification.** tabilet/memory-bank/milestone.md#m28---provider-node-expansion-batch-11
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M28 - Provider Node Expansion Batch 11

**Goal.** Add a frozen commerce, operations, and business provider-node batch
to the core catalog before standalone converter work.

**Scope.**

- Freeze Batch 11 before implementation: Magento, Paddle, Gumroad, Invoice
  Ninja, Odoo, ERPNext, WooCommerce, Wise, DHL, Onfleet, Unleashed Software,
  and Workable.
- Apply the same per-service curation loop, source-first evidence requirements,
  priority-only n8n evidence boundary, and no-execution/no-credential rules as
  M24.

**Acceptance.** Batch 11 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.
````

## Status record

````markdown
# Status M28 - Provider Node Expansion Batch 11

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Continue broadening provider-like n8n node catalog coverage before standalone
converter work. M28 focuses on commerce, operations, and business services with
source-backed metadata and auth/security classification.

## Frozen Batch

Batch 11 services:

- Magento
- Paddle
- Gumroad
- Invoice Ninja
- Odoo
- ERPNext
- WooCommerce
- Wise
- DHL
- Onfleet
- Unleashed Software
- Workable

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
| Batch list frozen | `[+]` | Frozen Batch 11 services are Magento, Paddle, Gumroad, Invoice Ninja, Odoo, ERPNext, WooCommerce, Wise, DHL, Onfleet, Unleashed Software, and Workable. |
| Magento curated | `[+]` | Added candidate/provider metadata, Adobe Commerce REST docs refs, bearer auth overlay, and `magento-rest-api-overlay.json`; stable public generic OpenAPI not found because Swagger generation is instance-local. |
| Paddle curated | `[+]` | Added official `PaddleHQ/paddle-openapi` API v1 OpenAPI ref, saved ignored review artifact, registry row, and complete bearer auth classification. |
| Gumroad curated | `[+]` | Added candidate/provider metadata, official API/OAuth docs refs, bearer auth overlay, and `gumroad-api-overlay.json`. |
| Invoice Ninja curated | `[+]` | Added candidate/provider metadata, official OpenAPI-rendered docs refs, X-API-TOKEN auth overlay, and `invoice-ninja-api-overlay.json`; no stable standalone downloadable OpenAPI URL found. |
| Odoo curated | `[+]` | Added candidate/provider metadata, official external API docs refs, RPC auth overlay, and no-overlay decision because the official surface is XML-RPC/JSON-RPC rather than REST-shaped. |
| ERPNext curated | `[+]` | Added candidate/provider metadata, official ERPNext/Frappe REST docs refs, token auth overlay, and `erpnext-rest-api-overlay.json`. |
| WooCommerce curated | `[+]` | Added candidate/provider metadata, official REST API docs refs, Basic/query consumer-key auth overlay, and `woocommerce-rest-api-overlay.json`. |
| Wise curated | `[+]` | Added candidate/provider metadata, official Platform API docs refs, bearer auth overlay, and `wise-platform-api-overlay.json`. |
| DHL curated | `[+]` | Added candidate/provider metadata, official Shipment Tracking docs refs, `DHL-API-Key` auth overlay, and `dhl-shipment-tracking-overlay.json`; no stable public downloadable OpenAPI URL found. |
| Onfleet curated | `[+]` | Added candidate/provider metadata, official API docs refs, Basic API-key auth overlay, and `onfleet-api-overlay.json`. |
| Unleashed Software curated | `[+]` | Added candidate/provider metadata, official API docs refs, API ID/signature/client-type auth overlay, and `unleashed-software-api-overlay.json`; apitools records signature metadata only. |
| Workable curated | `[+]` | Added candidate/provider metadata, official API docs refs, bearer auth overlay, and `workable-api-overlay.json`. |
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
````
