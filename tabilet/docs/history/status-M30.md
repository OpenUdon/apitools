# Retired milestone M30 - Provider Node Expansion Batch 13

**Milestone.** M30
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M30.md
**Source specification.** tabilet/memory-bank/milestone.md#m30---provider-node-expansion-batch-13
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M30 - Provider Node Expansion Batch 13

**Goal.** Add a frozen enrichment, AI, and data provider-node batch while
keeping AI/vector-specific converter work deferred.

**Scope.**

- Freeze Batch 13 before implementation: Humantic AI, Hunter, Jina AI,
  LingvaNex, Mistral AI, OpenAI, Perplexity, Mindee, Peekalink, Phantombuster,
  UpLead, and Dropcontact.
- Apply the same per-service curation loop, source-first evidence requirements,
  priority-only n8n evidence boundary, and no-execution/no-credential rules as
  M24.
- For AI-oriented providers, catalog only source-backed API metadata that fits
  the existing provider catalog boundary; do not add AI workflow behavior,
  model-selection policy, or runtime credential binding.

**Acceptance.** Batch 13 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.
````

## Status record

````markdown
# Status M30 - Provider Node Expansion Batch 13

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Continue broadening provider-like n8n node catalog coverage before standalone
converter work. M30 focuses on enrichment, AI, and data services while keeping
AI/vector-specific converter work outside this milestone.

## Frozen Batch

Batch 13 services:

- Humantic AI
- Hunter
- Jina AI
- LingvaNex
- Mistral AI
- OpenAI
- Perplexity
- Mindee
- Peekalink
- Phantombuster
- UpLead
- Dropcontact

`../n8n/packages/nodes-base/nodes` is a priority signal only. It is not runtime
compatibility evidence, and no n8n behavior claims should be recorded.

## Per-Service Curation Loop

For each service, try official machine-readable sources first, save and
register ignored official review artifacts when retrieved, review
auth/security completeness, add candidate/provider metadata, add security
classification or overlays when needed, generate docs-derived advisory overlays
only for useful REST-shaped subsets, record provenance fields, and run focused
tests plus `go run ./cmd/apitools catalog check --as-of 2026-05-19` before
committing each service row. For AI-oriented providers, catalog only
source-backed API metadata that fits the existing provider catalog boundary.

## Tasks

| Item | State | Notes |
|---|---|---|
| Batch list frozen | `[+]` | Frozen Batch 13 services are Humantic AI, Hunter, Jina AI, LingvaNex, Mistral AI, OpenAI, Perplexity, Mindee, Peekalink, Phantombuster, UpLead, and Dropcontact. |
| Humantic AI curated | `[+]` | Added candidate/provider metadata from the n8n Humantic AI node and official API root; no stable public official OpenAPI was found, so auth metadata and a small profile advisory overlay are docs-derived. |
| Hunter curated | `[+]` | Added candidate/provider metadata from the n8n Hunter node and official API docs; no stable public official OpenAPI was found, so api_key query auth and email/domain/leads/campaigns advisory operations are docs-derived. |
| Jina AI curated | `[+]` | Added candidate/provider metadata from the n8n Jina AI node, registered Search Foundation, Reader, and Search official OpenAPI artifacts, and added a bearer auth overlay for docs-backed services whose public specs omit security metadata. |
| LingvaNex curated | `[+]` | Added candidate/provider metadata from the n8n LingvaNex node and official API docs/GitHub SDK references; no stable public official OpenAPI was found, so bearer auth and translate/detect/languages/dictionary advisory operations are docs-derived. |
| Mistral AI curated | `[+]` | Added candidate/provider metadata from the n8n Mistral AI node and official API/API-key docs; no stable public official OpenAPI endpoint was found, so bearer auth and common v1 advisory operations are docs-derived. |
| OpenAI curated | `[+]` | Added candidate/provider metadata from the n8n OpenAI node and registered the official documented OpenAPI artifact linked from the official OpenAI OpenAPI repository; auth is complete from the spec. |
| Perplexity curated | `[+]` | Added candidate/provider metadata from the n8n Perplexity node and registered the official docs OpenAPI artifact; auth is complete from operation-level bearer security in the spec. |
| Mindee curated | `[+]` | Added candidate/provider metadata from the n8n Mindee node and official API overview/API-key docs; no stable public official OpenAPI was found, so token header auth and common document prediction advisory operations are docs-derived. |
| Peekalink curated | `[+]` | Added candidate/provider metadata from the n8n Peekalink node and registered the official Mintlify OpenAPI artifact advertised by official docs; auth is complete from the spec. |
| Phantombuster curated | `[+]` | Added candidate/provider metadata from the n8n Phantombuster node and official API docs; the probed OpenAPI endpoint requires authorization, so X-Phantombuster-Key auth and agent/container/script/org advisory operations are docs-derived. |
| UpLead curated | `[+]` | Added candidate/provider metadata from the n8n UpLead node and official docs; no stable public official OpenAPI was found, so Authorization header auth and person/company/credits advisory operations are docs-derived. |
| Dropcontact curated | `[+]` | Added candidate/provider metadata from the n8n Dropcontact node and official developer docs; no stable public official OpenAPI was found, so X-Access-Token auth and batch status advisory operations are docs-derived. |
| Batch quality review completed | `[+]` | Final review found no blocking issues. Verified with focused catalog/CLI tests, full `go test -count=1 ./...`, `go vet ./...`, whitespace checks, catalog check/stats, and per-provider inspect/advisory sanity checks on 2026-05-19. |

## Boundary Checks

- Do not treat n8n node presence as runtime compatibility evidence.
- Do not vendor downloaded official specs, Discovery documents, or SQLite
  caches into the public repository.
- Do not execute API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not treat advisory overlays as provider truth.
- Do not add AI workflow behavior, model-selection policy, runtime credential
  binding, or standalone converter work in this milestone.
````
