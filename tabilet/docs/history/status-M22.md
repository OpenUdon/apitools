# Retired milestone M22 - Provider Node Expansion Batch 5

**Milestone.** M22
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M22.md
**Source specification.** tabilet/memory-bank/milestone.md#m22---provider-node-expansion-batch-5
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M22 - Provider Node Expansion Batch 5

**Goal.** Add another frozen batch of popular provider-node services to broaden
REST-oriented catalog coverage before starting non-OpenAPI converter work.

**Scope.**

- Freeze Batch 5 before implementation: Acuity Scheduling, Bitly, Brevo,
  Clockify, Coda, DeepL, Figma, Ghost, Harvest, Help Scout, Iterable, and
  Mailjet.
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
- Keep Google Discovery, Smithy, Stone, and GraphQL converter work deferred
  until provider-node catalog coverage is substantially broader.
- Avoid n8n runtime compatibility claims, provider API execution, credential
  handling, request signing, account selection, or automatic metadata
  promotion.

**Acceptance.** Batch 5 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.
````

## Status record

````markdown
# Status M22 - Provider Node Expansion Batch 5

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Continue broadening provider-node catalog coverage before starting
non-OpenAPI converter work. M22 focuses on REST-oriented provider nodes with
source-backed metadata and auth/security classification.

## Frozen Batch

Batch 5 services:

- Acuity Scheduling
- Bitly
- Brevo
- Clockify
- Coda
- DeepL
- Figma
- Ghost
- Harvest
- Help Scout
- Iterable
- Mailjet

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
| Batch list frozen | `[+]` | Frozen Batch 5 services are Acuity Scheduling, Bitly, Brevo, Clockify, Coda, DeepL, Figma, Ghost, Harvest, Help Scout, Iterable, and Mailjet. |
| Acuity Scheduling curated | `[+]` | Added official-docs provider metadata, Basic-auth security overlay, and docs-derived endpoint overlay for a focused API v1 subset. |
| Bitly curated | `[+]` | Added official OpenAPI v4 metadata, saved ignored review artifact, and complete bearer-auth classification. |
| Brevo curated | `[+]` | Added official OpenAPI v3 metadata, saved ignored review artifact, and complete api-key/partner-key auth classification. |
| Clockify curated | `[+]` | Added official-docs provider metadata, X-Api-Key/X-Addon-Token security overlay, and docs-derived endpoint overlay for a focused API v1 subset. |
| Coda curated | `[+]` | Added official API v1 OpenAPI metadata, saved ignored review artifact, and complete bearer-auth classification. |
| DeepL curated | `[+]` | Added official OpenAPI repository metadata, saved ignored review artifact, and complete auth-header classification. |
| Figma curated | `[+]` | Added official rest-api-spec OpenAPI metadata, saved ignored review artifact, and complete PAT/plan-token/OAuth classification. |
| Ghost curated | `[+]` | Added official-docs provider metadata, Admin/Content API security overlay, and docs-derived Admin API endpoint overlay. |
| Harvest curated | `[+]` | Added official-docs provider metadata, bearer-auth security overlay, and docs-derived endpoint overlay with Harvest-Account-ID as request metadata. |
| Help Scout curated | `[+]` | Added official-docs provider metadata, Inbox/Docs API security overlay, and docs-derived Inbox API endpoint overlay. |
| Iterable curated | `[+]` | Added official-docs provider metadata, Api-Key security overlay, and docs-derived endpoint overlay for a focused REST subset. |
| Mailjet curated | `[+]` | Added official-docs provider metadata, Basic-auth security overlay, and docs-derived endpoint overlay for focused REST/send endpoints. |
| Human-doc endpoint overlays completed | `[+]` | Built and registered docs-derived endpoint overlays for Acuity Scheduling, Clockify, Ghost, Harvest, Help Scout, Iterable, and Mailjet; official-OpenAPI providers did not need endpoint overlays. |
| Batch quality review completed | `[+]` | Reviewed M22 diff; fixed Harvest account header so it is request metadata rather than a security credential. Verified `go test -count=1 ./...`, `go vet ./...`, `git diff --check`, and `go run ./cmd/apitools catalog check --as-of 2026-05-19`. |

## Boundary Checks

- Do not treat n8n node presence as runtime compatibility evidence.
- Do not vendor downloaded official specs, Discovery documents, or SQLite
  caches into the public repository.
- Do not execute API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not treat advisory overlays as provider truth.
- Do not start Google Discovery, Smithy, Stone, or GraphQL converter work in
  this milestone.
````
