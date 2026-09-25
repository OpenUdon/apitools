# Retired milestone M57 - Docs-Overlay Provider Batch

**Milestone.** M57
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M57.md
**Source specification.** tabilet/memory-bank/milestone.md#m57---docs-overlay-provider-batch
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M57 - Docs-Overlay Provider Batch

**Goal.** Add a source-first docs-overlay batch for provider APIs where the
best usable official source is documentation rather than a durable downloadable
OpenAPI artifact: Marketo, Aircall, Checkr, Airwallex, Apollo.io, AfterShip,
and Adobe Acrobat Sign.

**Scope.**

- For each provider, first confirm whether an official downloadable OpenAPI or
  Swagger artifact exists. If one exists and is suitable, record it as official
  OpenAPI metadata instead of building an advisory overlay.
- When no durable official spec is available, add catalog/candidate metadata
  plus a reviewed docs-derived endpoint overlay only from provider-owned API
  docs. Record a no-overlay decision where official docs cannot support a
  useful reviewed subset.
- Keep each advisory overlay narrow, deterministic, and clearly marked as
  non-official OpenAPI. Include source refs, operation scope, auth/security
  notes, and omission rationale in the overlay registry or no-overlay notes.
- Prefer stable workflow primitives such as contacts/leads, calls, background
  check candidates/reports, payments/transfers, enrichment/search, shipment
  tracking, and agreements only where official docs make request and response
  shapes clear enough for metadata-only import.
- Do not sign in to accounts, call APIs, create applicants/leads/shipments,
  upload documents, initiate financial movement, start background checks, fetch
  tokens, or scrape authenticated documentation.

**Acceptance.** The seven M57 providers have catalog/candidate metadata,
official source refs, auth/security classifications, and either a tracked
docs-derived advisory overlay or an explicit no-overlay decision; every
overlay is source-reviewed, deterministic, marked non-official OpenAPI, and
visible in advisory/security/reporting views; Go, catalog, consumer, and diff
checks pass.
````

## Status record

````markdown
# Status M57 - Docs-Overlay Provider Batch

| Item | State | Notes |
|---|---|---|
| Source/spec check | `[+]` | Reviewed official source surfaces for all seven providers. No stable directly importable provider-wide OpenAPI 2/3 artifact was recorded for this batch; Apollo exposes endpoint fragments, AfterShip advertises an OAS export behind an unusable direct fetch path, and Acrobat Sign SDK JSON is Swagger 1.2-style. |
| Overlay policy frozen | `[+]` | Use docs-derived overlays for REST-shaped reviewed subsets where source docs expose enough endpoint/auth metadata; use no-overlay decisions where a generic overlay would hide sensitive product, tenant, or money-movement semantics. |
| Marketo cataloged | `[+]` | Added Adobe Marketo Engage candidate/provider rows, bearer auth overlay, source refs, and `marketo-rest-api-overlay.json` for lead metadata, lead retrieval, and activity retrieval. |
| Aircall cataloged | `[+]` | Added Aircall candidate/provider rows, Basic/OAuth2 auth overlay, source refs, and `aircall-public-api-overlay.json` for users, availability, calls, and contacts. |
| Checkr cataloged | `[+]` | Added Checkr candidate/provider rows, Basic Auth overlay, source refs, and `checkr-api-overlay.json` for a conservative account, candidate, package, and report read subset. |
| Airwallex cataloged | `[+]` | Added Airwallex candidate/provider rows, bearer-token auth overlay, source refs, and an explicit no-endpoint-overlay decision because broad financial APIs need product/region/account-specific source handling. |
| Apollo.io cataloged | `[+]` | Added Apollo.io candidate/provider rows, API-key/bearer auth overlay, source refs, and `apollo-api-overlay.json` for people and organization search from official endpoint-level OpenAPI fragments. |
| AfterShip cataloged | `[+]` | Added AfterShip candidate/provider rows, current `as-api-key` auth overlay, source refs, and `aftership-tracking-api-overlay.json` for tracking and courier read endpoints. |
| Adobe Acrobat Sign cataloged | `[+]` | Added Adobe Acrobat Sign candidate/provider rows, OAuth bearer auth overlay, source refs, and `adobe-acrobat-sign-api-overlay.json` because the official SDK JSON is Swagger 1.2-style rather than importable OpenAPI 2/3. |
| Verification | `[+]` | Passed `go test ./...`, `GOWORK=off go test ./...`, `go vet ./...`, `GOWORK=off go vet ./...`, `GOWORK=off go mod tidy -diff`, `go run ./cmd/apitools catalog check`, `go run ./cmd/apitools catalog stats`, advisory/security CLI smoke checks, M57 overlay JSON metadata checks, `go test ./...` in `../openudon` and `../udon`, text-origin audit, `git diff --check`, and `git -C ../tofu diff --check -- apitools`. Final code/document review found no M57 blocker. |

## Source Notes

- Marketo REST API coverage should use Adobe Marketo Engage developer docs and
  endpoint references, with OAuth/client handling kept advisory.
- Aircall Public API coverage should start from official API references and
  keep account, rate-limit, and telephony side effects out of scope.
- Checkr coverage should be especially conservative because candidate and
  background-check workflows are regulated and side-effecting.
- Airwallex coverage should avoid financial movement semantics beyond
  metadata-only auth, source, and reviewed endpoint shape.
- Apollo.io coverage should use official API docs only and avoid copying
  unofficial specs, lead data examples beyond minimal shapes, or scraping
  authenticated pages.
- AfterShip coverage should focus on official shipping/tracking docs and avoid
  carrier account or shipment mutation behavior unless official docs support a
  narrow reviewed subset.
- Adobe Acrobat Sign has official REST docs and OpenAPI/SDK evidence; source
  review should decide whether direct OpenAPI import is durable or whether a
  narrow advisory overlay is safer.
- M57 conclusion: docs-derived endpoint overlays were built for Marketo,
  Aircall, Checkr, Apollo.io, AfterShip, and Adobe Acrobat Sign. Airwallex is
  cataloged with auth/source metadata and no endpoint overlay because broad
  financial movement surfaces need a more explicit source-family design before
  OpenAPI-shaped operation coverage is useful.
- Post-review remediation: Apollo search overlays now model documented query
  parameters instead of JSON request bodies; Marketo required lead/activity
  query parameters are marked required; and Acrobat Sign agreement read overlays
  preserve the `agreement_read` OAuth scope in both generated advisory OpenAPI
  and catalog security metadata.

## Boundary Checks

- Prefer official downloadable OpenAPI when it is durable and provider-owned;
  only use docs-derived overlays when no suitable official spec is available.
- Mark every docs-derived artifact as advisory and non-official OpenAPI.
- Do not sign in, fetch tokens, create leads/candidates/shipments/agreements,
  upload documents, initiate payments, trigger calls, or start background
  checks.
- Do not scrape authenticated docs, copy third-party connector metadata, or
  ingest unofficial community OpenAPI documents.
- Keep reviewed subsets intentionally small and document omitted operation
  families when request/response shapes or safety semantics are unclear.
````
