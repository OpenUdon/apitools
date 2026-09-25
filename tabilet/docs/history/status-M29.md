# Retired milestone M29 - Provider Node Expansion Batch 12

**Milestone.** M29
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M29.md
**Source specification.** tabilet/memory-bank/milestone.md#m29---provider-node-expansion-batch-12
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M29 - Provider Node Expansion Batch 12

**Goal.** Add a frozen infrastructure, security, and admin provider-node batch
to the core catalog before standalone converter work.

**Scope.**

- Freeze Batch 12 before implementation: Bitwarden, Cisco, Cortex, Home
  Assistant, MISP, NetScaler, Rundeck, SecurityScorecard, UrlScan.io, Venafi,
  Wekan, and Zammad.
- Apply the same per-service curation loop, source-first evidence requirements,
  priority-only n8n evidence boundary, and no-execution/no-credential rules as
  M24.

**Acceptance.** Batch 12 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.
````

## Status record

````markdown
# Status M29 - Provider Node Expansion Batch 12

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Continue broadening provider-like n8n node catalog coverage before standalone
converter work. M29 focuses on infrastructure, security, and administration
services with source-backed metadata and auth/security classification.

## Frozen Batch

Batch 12 services:

- Bitwarden
- Cisco
- Cortex
- Home Assistant
- MISP
- NetScaler
- Rundeck
- SecurityScorecard
- UrlScan.io
- Venafi
- Wekan
- Zammad

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
| Batch list frozen | `[+]` | Frozen Batch 12 services are Bitwarden, Cisco, Cortex, Home Assistant, MISP, NetScaler, Rundeck, SecurityScorecard, UrlScan.io, Venafi, Wekan, and Zammad. |
| Bitwarden curated | `[+]` | Added candidate/provider metadata, official Public API docs refs, auth overlay, and docs-derived endpoint overlay. M29 review found no stable public downloadable official OpenAPI artifact at explicit probes. |
| Cisco curated | `[+]` | Curated as Cisco Webex because the n8n priority signal is `Cisco/Webex`; added candidate/provider metadata, official Webex docs refs, auth overlay, and docs-derived endpoint overlay. M29 review found no stable public downloadable official OpenAPI artifact at explicit probes. |
| Cortex curated | `[+]` | Added candidate/provider metadata, StrangeBee Cortex API guide refs, auth overlay, and docs-derived endpoint overlay. M29 review found no stable public official OpenAPI artifact at explicit probes. |
| Home Assistant curated | `[+]` | Added candidate/provider metadata, official REST API docs refs, auth overlay, and docs-derived endpoint overlay. M29 review found no stable public official OpenAPI artifact at explicit probes. |
| MISP curated | `[+]` | Added candidate/provider metadata with official MISP Automation API OpenAPI YAML from the MISP repository, artifact registry metadata, and complete Authorization-header API-key security classification. |
| NetScaler curated | `[+]` | Added candidate/provider metadata, official NetScaler ADC NITRO API docs refs, auth overlay, and docs-derived endpoint overlay. M29 review found no stable public official OpenAPI artifact at explicit probes. |
| Rundeck curated | `[+]` | Added candidate/provider metadata with official Rundeck / Runbook Automation OpenAPI YAML from docs.rundeck.com, artifact registry metadata, and complete API-token/session/JWT security classification. |
| SecurityScorecard curated | `[+]` | Added candidate/provider metadata with official SecurityScorecard Swagger 2.0 JSON from api.securityscorecard.io, artifact registry metadata, and complete Token header security classification. |
| UrlScan.io curated | `[+]` | Added candidate/provider metadata with official urlscan.io OpenAPI JSON bundle from docs.urlscan.io, artifact registry metadata, and complete api-key header security classification. |
| Venafi curated | `[+]` | Added candidate/provider metadata, TLS Protect Cloud and Datacenter WebSDK docs refs, auth overlay, and docs-derived endpoint overlay. M29 review found no stable public official OpenAPI artifact at explicit probes. |
| Wekan curated | `[+]` | Added candidate/provider metadata, official REST API wiki and OpenAPI generator refs, auth overlay, and docs-derived endpoint overlay. M29 review found generator tooling but no stable public generated OpenAPI artifact to catalog. |
| Zammad curated | `[+]` | Added candidate/provider metadata, official API docs refs, auth overlay, and docs-derived endpoint overlay. M29 review found no stable public official OpenAPI artifact at explicit probes. |
| Batch quality review completed | `[+]` | Final code review completed; verification passed and M29 changes are ready for the closing commits. |

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
