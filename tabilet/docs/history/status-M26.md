# Retired milestone M26 - Provider Node Expansion Batch 9

**Milestone.** M26
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M26.md
**Source specification.** tabilet/memory-bank/milestone.md#m26---provider-node-expansion-batch-9
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M26 - Provider Node Expansion Batch 9

**Goal.** Add a frozen forms, no-code, and app-builder provider-node batch to
the core catalog before standalone converter work.

**Scope.**

- Freeze Batch 9 before implementation: JotForm, Form.io, Formstack,
  SurveyMonkey, KoBoToolbox, Bubble, Cal, Cockpit, Stackby, Airtop, TimeSaved,
  and FileMaker.
- Apply the same per-service curation loop, source-first evidence requirements,
  priority-only n8n evidence boundary, and no-execution/no-credential rules as
  M24.

**Acceptance.** Batch 9 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.
````

## Status record

````markdown
# Status M26 - Provider Node Expansion Batch 9

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Continue broadening provider-like n8n node catalog coverage before standalone
converter work. M26 focuses on forms, no-code, and app-builder services with
source-backed metadata and auth/security classification.

## Frozen Batch

Batch 9 services:

- JotForm
- Form.io
- Formstack
- SurveyMonkey
- KoBoToolbox
- Bubble
- Cal
- Cockpit
- Stackby
- Airtop
- TimeSaved
- FileMaker

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
| Batch list frozen | `[+]` | Frozen Batch 9 services are JotForm, Form.io, Formstack, SurveyMonkey, KoBoToolbox, Bubble, Cal, Cockpit, Stackby, Airtop, TimeSaved, and FileMaker. |
| JotForm curated | `[+]` | Added candidate/provider metadata, docs source refs, API-key security overlay, endpoint advisory overlay, artifact registry row, and deterministic coverage. |
| Form.io curated | `[+]` | Added candidate/provider metadata, docs/auth source refs, token security overlay, endpoint advisory overlay, artifact registry row, and deterministic coverage. |
| Formstack curated | `[+]` | Added candidate/provider metadata, Forms API docs source refs, OAuth security overlay, endpoint advisory overlay, artifact registry row, and deterministic coverage. |
| SurveyMonkey curated | `[+]` | Added candidate/provider metadata, API v3 docs source refs, OAuth bearer security overlay, endpoint advisory overlay, artifact registry row, and deterministic coverage. |
| KoBoToolbox curated | `[+]` | Added candidate/provider metadata, official API v2 OpenAPI schema ref, ignored review artifact registration, complete auth classification, and deterministic coverage. |
| Bubble curated | `[+]` | Added candidate/provider metadata, API docs source refs, bearer security overlay, generic app-specific endpoint advisory overlay, artifact registry row, and deterministic coverage. |
| Cal curated | `[+]` | Added candidate/provider metadata, official API v2 OpenAPI ref, ignored review artifact registration, present-incomplete auth review overlay, and deterministic coverage. |
| Cockpit curated | `[+]` | Added candidate/provider metadata, self-hosted API docs source refs, token security overlay, endpoint advisory overlay, artifact registry row, and deterministic coverage. |
| Stackby curated | `[+]` | Added candidate/provider metadata, Developer API/API-key docs source refs, API-key security overlay, endpoint advisory overlay, artifact registry row, and deterministic coverage. |
| Airtop curated | `[+]` | Added candidate/provider metadata, official OpenAPI ref, ignored review artifact registration, present-incomplete auth review overlay, and deterministic coverage. |
| TimeSaved curated | `[+]` | Recorded frozen-batch review outcome as an n8n workflow metadata helper with no external OpenAPI/auth/endpoint overlay surface. |
| FileMaker curated | `[+]` | Added candidate/provider metadata, Claris Data API docs source refs, Basic/session bearer security overlay, endpoint advisory overlay, artifact registry row, and deterministic coverage. |
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
