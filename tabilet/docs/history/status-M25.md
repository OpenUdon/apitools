# Retired milestone M25 - Provider Node Expansion Batch 8

**Milestone.** M25
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M25.md
**Source specification.** tabilet/memory-bank/milestone.md#m25---provider-node-expansion-batch-8
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M25 - Provider Node Expansion Batch 8

**Goal.** Add a frozen CRM, sales, and marketing provider-node batch to the
core catalog before standalone converter work.

**Scope.**

- Freeze Batch 8 before implementation: Autopilot, Drift, Freshworks CRM,
  GetResponse, HighLevel, Keap, MailerLite, Mautic, Monica CRM, Salesmate,
  Sendy, and Vero.
- Apply the same per-service curation loop, source-first evidence requirements,
  priority-only n8n evidence boundary, and no-execution/no-credential rules as
  M24.

**Acceptance.** Batch 8 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.
````

## Status record

````markdown
# Status M25 - Provider Node Expansion Batch 8

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Continue broadening provider-like n8n node catalog coverage before standalone
converter work. M25 focuses on CRM, sales, and marketing services with
source-backed metadata and auth/security classification.

## Frozen Batch

Batch 8 services:

- Autopilot
- Drift
- Freshworks CRM
- GetResponse
- HighLevel
- Keap
- MailerLite
- Mautic
- Monica CRM
- Salesmate
- Sendy
- Vero

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
| Batch list frozen | `[+]` | Frozen Batch 8 services are Autopilot, Drift, Freshworks CRM, GetResponse, HighLevel, Keap, MailerLite, Mautic, Monica CRM, Salesmate, Sendy, and Vero. |
| Autopilot curated | `[+]` | Added candidate/provider metadata, docs source refs, API-key security overlay, endpoint advisory overlay, artifact registry row, and deterministic coverage. |
| Drift curated | `[+]` | Added candidate/provider metadata, docs source refs, OAuth bearer security overlay, endpoint advisory overlay, artifact registry row, and deterministic coverage. |
| Freshworks CRM curated | `[+]` | Added candidate/provider metadata, account-host docs source refs, token security overlay, endpoint advisory overlay, artifact registry row, and deterministic coverage. |
| GetResponse curated | `[+]` | Added candidate/provider metadata, API v3 docs source refs, API-key/OAuth security overlay, endpoint advisory overlay, artifact registry row, and deterministic coverage. |
| HighLevel curated | `[+]` | Added candidate/provider metadata, official OpenAPI refs for Contacts, Opportunities, Calendars, Users, and OAuth from the official highlevel-api-docs repository, artifact registry rows, and complete auth classification. |
| Keap curated | `[+]` | Added candidate/provider metadata, REST docs source refs, OAuth bearer security overlay, endpoint advisory overlay, artifact registry row, and deterministic coverage. |
| MailerLite curated | `[+]` | Added candidate/provider metadata, current and Classic docs source refs, bearer/API-key security overlay, endpoint advisory overlay, artifact registry row, and deterministic coverage. |
| Mautic curated | `[+]` | Added candidate/provider metadata, instance-host docs source refs, Basic/bearer security overlay, endpoint advisory overlay, artifact registry row, and deterministic coverage. |
| Monica CRM curated | `[+]` | Added candidate/provider metadata, API docs source refs, bearer security overlay, endpoint advisory overlay, artifact registry row, and deterministic coverage. |
| Salesmate curated | `[+]` | Added candidate/provider metadata, support/API docs source refs, sessionToken/x-linkname security overlay, endpoint advisory overlay, artifact registry row, and deterministic coverage. |
| Sendy curated | `[+]` | Added candidate/provider metadata, installation-hosted API docs source refs, API-key security overlay, endpoint advisory overlay, artifact registry row, and deterministic coverage. |
| Vero curated | `[+]` | Added candidate/provider metadata, Track API docs source refs, auth-token security overlay, endpoint advisory overlay, artifact registry row, and deterministic coverage. |
| Batch quality review completed | `[+]` | Batch review and full verification completed before final commit. |

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
