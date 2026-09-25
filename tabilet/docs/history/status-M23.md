# Retired milestone M23 - Provider Node Expansion Batch 6

**Milestone.** M23
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M23.md
**Source specification.** tabilet/memory-bank/milestone.md#m23---provider-node-expansion-batch-6
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M23 - Provider Node Expansion Batch 6

**Goal.** Add another frozen batch of REST-oriented provider-node services to
broaden catalog coverage before starting non-OpenAPI converter work.

**Scope.**

- Freeze Batch 6 before implementation: Action Network, Adalo, Affinity,
  Agile CRM, Bannerbear, Baserow, Beeminder, Brandfetch, Clearbit, ConvertKit,
  Copper, Discourse, Freshservice, Gong, and Grist.
- Treat each service as a separate task and commit unit once curation begins.
- Use `../n8n/packages/nodes-base/nodes` as a priority signal only, not as
  runtime compatibility evidence.
- For each service, try official OpenAPI, Swagger, OpenAPI index, Discovery, or
  other official machine-readable sources first; save downloaded official
  artifacts under ignored catalog cache paths and register them in SQLite when
  saved.
- If no official OpenAPI exists, record official human docs and build tracked
  docs-derived endpoint overlay assets plus builders only when official docs
  expose enough endpoint instructions for a useful reviewed subset.
- Add or update candidate records, durable provider entries, spec references,
  security classifications, and security overlays with source authority,
  verified dates, source notes, source refs, license notes, quirks, and user
  OpenAPI need.
- Keep Google Discovery, Smithy, Stone, GraphQL, and other converter work
  deferred until provider-node catalog coverage is substantially broader.
- Avoid n8n runtime compatibility claims, provider API execution, credential
  handling, request signing, account selection, or automatic metadata
  promotion.

**Acceptance.** Batch 6 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.
````

## Status record

````markdown
# Status M23 - Provider Node Expansion Batch 6

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Continue broadening provider-node catalog coverage before starting
non-OpenAPI converter work. M23 focuses on REST-oriented provider nodes with
source-backed metadata and auth/security classification.

## Frozen Batch

Batch 6 services:

- Action Network
- Adalo
- Affinity
- Agile CRM
- Bannerbear
- Baserow
- Beeminder
- Brandfetch
- Clearbit
- ConvertKit
- Copper
- Discourse
- Freshservice
- Gong
- Grist

`../n8n/packages/nodes-base/nodes` is a priority signal only. It is not runtime
compatibility evidence, and no n8n behavior claims should be recorded.

## Per-Service Curation Loop

For each service:

1. Try official OpenAPI, Swagger, OpenAPI index, Discovery, or other official
   machine-readable sources first.
2. Save downloaded official review artifacts under ignored
   `catalog-openapi-cache/openapi/` or
   `catalog-openapi-cache/google-discovery/` when a document is retrieved.
3. Register saved official artifact paths in `catalog-openapi-cache/cache.sqlite`
   when artifacts are saved.
4. Review auth/security completeness from the official spec or docs.
5. Add a candidate record if missing.
6. Promote the provider into durable catalog metadata only after source review.
7. Add source-backed security classification or security overlay metadata when
   upstream security is incomplete, absent, non-OpenAPI, or docs-derived.
8. If no official OpenAPI exists, build tracked docs-derived advisory overlay
   assets and service-specific builder programs only when endpoint/security
   metadata is useful for OpenAPI-only consumers.
9. Record source authority, source notes, source refs, verified dates, license
   notes, quirks, user OpenAPI need, and auth/security completeness.
10. Run focused tests and `go run ./cmd/apitools catalog check --as-of
    2026-05-19` before committing each service row.

## Tasks

| Item | State | Notes |
|---|---|---|
| Batch list frozen | `[+]` | Frozen Batch 6 services are Action Network, Adalo, Affinity, Agile CRM, Bannerbear, Baserow, Beeminder, Brandfetch, Clearbit, ConvertKit, Copper, Discourse, Freshservice, Gong, and Grist. |
| Action Network curated | `[+]` | Added Action Network candidate/provider metadata, recorded no stable public official OpenAPI from explicit probes, added OSDI-API-Token auth overlay, and generated a tracked API v2 docs-derived endpoint advisory overlay plus builder coverage. |
| Adalo curated | `[+]` | Added Adalo candidate/provider metadata, recorded no stable public official OpenAPI from explicit probes, added bearer API-key auth overlay, and generated a tracked API docs-derived endpoint advisory overlay plus builder coverage. |
| Affinity curated | `[+]` | Added Affinity candidate/provider metadata, recorded no stable public official OpenAPI from explicit probes, added basic/bearer API-key auth overlay, and generated a tracked V1 docs-derived endpoint advisory overlay plus builder coverage. |
| Agile CRM curated | `[+]` | Added Agile CRM candidate/provider metadata, recorded no stable public official OpenAPI from explicit probes, added HTTP Basic email/API-key auth overlay, and generated a tracked REST API docs-derived endpoint advisory overlay plus builder coverage. |
| Bannerbear curated | `[+]` | Added Bannerbear candidate/provider metadata, recorded no stable public official OpenAPI from explicit probes, added bearer API-key auth overlay, and generated a tracked API v2 docs-derived endpoint advisory overlay plus builder. |
| Baserow curated | `[+]` | Added Baserow candidate/provider metadata, official API OpenAPI reference, saved ignored review artifact registration, and complete bearer-token/JWT/database-token auth classification; saved artifact is parseable OpenAPI but strict-invalid for review. |
| Beeminder curated | `[+]` | Added Beeminder candidate/provider metadata, recorded no stable public official OpenAPI from explicit probes, added auth_token/OAuth bearer auth overlay, and generated a tracked API v1 docs-derived endpoint advisory overlay plus builder coverage. |
| Brandfetch curated | `[+]` | Added Brandfetch candidate/provider metadata, official Brand API OpenAPI reference, saved ignored review artifact registration, and complete bearer-token auth classification. |
| Clearbit curated | `[+]` | Added Clearbit candidate/provider metadata, recorded no stable public official OpenAPI from explicit probes, added secret API-key Basic auth overlay, and generated a tracked docs-derived Prospector/Enrichment endpoint advisory overlay plus builder coverage. |
| ConvertKit curated | `[+]` | Added Kit/ConvertKit candidate/provider metadata, official API V4 OpenAPI reference, saved ignored review artifact registration, and complete API-key/OAuth auth classification; saved artifact is parseable OpenAPI but strict-invalid for review. |
| Copper curated | `[+]` | Added Copper candidate/provider metadata, recorded no stable public official OpenAPI from explicit probes, added Developer API token/user-email/application header auth overlay, and generated a tracked docs-derived endpoint advisory overlay plus builder coverage. |
| Discourse curated | `[+]` | Added Discourse candidate/provider metadata, official OpenAPI 3.1 reference, saved ignored review artifact registration, and advisory Api-Key/Api-Username auth overlay because the official spec currently lacks OpenAPI securitySchemes. |
| Freshservice curated | `[+]` | Added Freshservice candidate/provider metadata, recorded no stable public official OpenAPI from explicit probes, added API-key Basic auth overlay, and generated a tracked API v2 docs-derived endpoint advisory overlay plus builder coverage. |
| Gong curated | `[+]` | Added Gong candidate/provider metadata, recorded no stable public official OpenAPI from explicit probes, added access-key/secret Basic plus OAuth bearer auth overlay, and generated a narrow tracked docs-derived users/call-upload advisory overlay plus builder coverage. |
| Grist curated | `[+]` | Added Grist candidate/provider metadata, recorded no stable public official OpenAPI from explicit probes, added bearer API-key auth overlay, and generated a tracked REST API docs-derived endpoint advisory overlay plus builder. |
| Human-doc endpoint overlays completed | `[+]` | Every M23 human-docs-only REST provider has either a tracked docs-derived endpoint advisory overlay plus builder coverage or an official OpenAPI reference that makes an endpoint overlay unnecessary. |
| Batch quality review completed | `[+]` | Final M23 verification completed with catalog tests, full tests, vet, diff check, catalog quality check, and CLI smoke checks on 2026-05-19. |
| Post-review Copper auth semantics fixed | `[+]` | Follow-up review fix committed in `2acc84c` corrected the Copper advisory endpoint overlay to require the access-token, application, and user-email headers together on each generated operation. |

## Boundary Checks

- Do not treat n8n node presence as runtime compatibility evidence.
- Do not vendor downloaded official specs, Discovery documents, or SQLite
  caches into the public repository.
- Do not execute API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not treat advisory overlays as provider truth.
- Do not start Google Discovery, Smithy, Stone, GraphQL, or other converter
  work in this milestone.
````
