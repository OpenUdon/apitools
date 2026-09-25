# Retired milestone M13 - Infra/API Platform Expansion Batch

**Milestone.** M13
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M13.md
**Source specification.** tabilet/memory-bank/milestone.md#m13---infraapi-platform-expansion-batch
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M13 - Infra/API Platform Expansion Batch

**Goal.** Add a frozen infra/API-platform provider batch using n8n node
presence as a priority signal only.

**Scope.**

- Freeze Batch 3 before implementation: AWS S3, AWS Lambda, AWS SNS,
  Bitbucket, CircleCI, Cloudflare, Databricks, Elastic, Grafana, Jenkins,
  Netlify, Sentry, Snowflake, Splunk, and Supabase.
- Treat each service as a separate task and commit unit once curation begins.
- For each service, try official OpenAPI, Swagger, Discovery, or other
  official machine-readable sources first; save downloaded official artifacts
  under ignored catalog cache paths and register them in SQLite when saved.
- If no official OpenAPI exists, record official human docs and build tracked
  service-specific docs-derived endpoint overlay assets plus builders only when
  official docs expose enough endpoint instructions for a useful reviewed
  subset.
- Add or update candidate records, durable provider entries, spec references,
  security classifications, and security overlays with source authority,
  verified dates, source notes, source refs, license notes, quirks, and user
  OpenAPI need.
- For AWS providers, describe SigV4 as auth/security metadata only; do not
  implement signing or choose accounts.
- Avoid n8n runtime compatibility claims, provider API execution, credential
  handling, request signing, or automatic account decisions.

**Acceptance.** Batch 3 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.
````

## Status record

````markdown
# Status M13 - Infra/API Platform Expansion Batch

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Add a frozen infra/API-platform catalog batch with source-backed provider
metadata, auth/security classification, and quality-gate coverage.

## Frozen Batch

Batch 3 services:

- AWS S3
- AWS Lambda
- AWS SNS
- Bitbucket
- CircleCI
- Cloudflare
- Databricks
- Elastic
- Grafana
- Jenkins
- Netlify
- Sentry
- Snowflake
- Splunk
- Supabase

`../n8n/packages/nodes-base/nodes` is a priority signal only. It is not runtime
compatibility evidence, and no n8n behavior claims should be recorded.

## Per-Service Curation Loop

For each service:

1. Try official OpenAPI, Swagger, Discovery, or other official
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
| Batch list frozen | `[+]` | Frozen Batch 3 services are AWS S3, AWS Lambda, AWS SNS, Bitbucket, CircleCI, Cloudflare, Databricks, Elastic, Grafana, Jenkins, Netlify, Sentry, Snowflake, Splunk, and Supabase. |
| AWS S3 curated | `[+]` | Added AWS S3 candidate/provider metadata, official AWS API Models Smithy JSON service model reference, S3 API and SigV4 docs refs, SigV4 advisory overlay metadata, and focused catalog checks. |
| AWS Lambda curated | `[+]` | Added AWS Lambda candidate/provider metadata, official AWS API Models Smithy JSON service model reference, Lambda API and SigV4 docs refs, SigV4 advisory overlay metadata, and focused catalog checks. |
| AWS SNS curated | `[+]` | Added AWS SNS candidate/provider metadata, official AWS API Models Smithy JSON service model reference, SNS API and SigV4 docs refs, SigV4 advisory overlay metadata, and focused catalog checks. |
| Bitbucket curated | `[+]` | Added Bitbucket Cloud candidate/provider metadata, official canonical Swagger/OpenAPI reference, REST/auth docs refs, complete auth classification, and focused checks. Current refresh validation rejects the official Swagger document, so no saved artifact path was registered. |
| CircleCI curated | `[+]` | Added CircleCI candidate/provider metadata, official API v2 OpenAPI reference, API/token docs refs, complete auth classification, and focused checks. Current refresh validation rejects the official OpenAPI document, so no saved artifact path was registered. |
| Cloudflare curated | `[+]` | Added Cloudflare candidate/provider metadata, official API Schemas OpenAPI reference, API/token docs refs, complete auth classification, and focused checks. The official OpenAPI YAML is larger than the default refresh limit and current refresh validation still rejects it, so no saved artifact path was registered. |
| Databricks curated | `[+]` | Added Databricks candidate/provider metadata, official REST/auth/PAT docs refs, no-public-OpenAPI classification, bearer-token advisory auth overlay, and focused checks. |
| Elastic curated | `[+]` | Added Elastic/Elasticsearch candidate/provider metadata, official downloadable Elasticsearch OpenAPI reference, API/auth/spec repo docs refs, complete auth classification, and focused checks. Current refresh validation rejects the official OpenAPI document, so no saved artifact path was registered. |
| Grafana curated | `[+]` | Added Grafana candidate/provider metadata, official HTTP API OpenAPI/docs/auth/new-API refs, complete auth classification from the official OpenAPI, and focused checks. |
| Jenkins curated | `[+]` | Added Jenkins candidate/provider metadata, official Remote Access API and API-token docs refs, no-public-OpenAPI classification, basic/API-token plus crumb advisory auth overlay, and focused checks. |
| Netlify curated | `[+]` | Added Netlify candidate/provider metadata, official OpenAPI/API docs refs, complete auth classification, refreshed and registered ignored OpenAPI artifact path, and focused checks. |
| Sentry curated | `[+]` | Added Sentry candidate/provider metadata, official API/auth/permissions docs refs, no-public-OpenAPI classification, bearer/legacy-key advisory auth overlay, and focused checks. |
| Snowflake curated | `[+]` | Added Snowflake candidate/provider metadata, official SQL API OpenAPI plus REST/auth docs refs, complete auth classification from the official OpenAPI, and focused checks. |
| Splunk curated | `[+]` | Added Splunk candidate/provider metadata, official REST API/usage/ACS OpenAPI docs refs, no-general-public-OpenAPI classification, token/basic advisory auth overlay, and focused checks. |
| Supabase curated | `[+]` | Added Supabase candidate/provider metadata, official Management API OpenAPI and docs refs, present-incomplete bearer auth review overlay, and focused checks. Current refresh validation rejects the official OpenAPI document, so no saved artifact path was registered. |
| Batch quality review completed | `[+]` | Ran full verification and milestone review after the final provider task: `go test ./...`, `go vet ./...`, `git diff --check`, catalog quality check, and CLI help checks passed. Follow-up review fixes promoted Grafana and Snowflake to official OpenAPI-backed entries and corrected Smithy JSON resolution before human-doc fallbacks. |

## Boundary Checks

- Do not treat n8n node presence as runtime compatibility evidence.
- Do not vendor downloaded official specs or Discovery documents.
- Do not execute API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not implement AWS SigV4 signing; record it as metadata only.
- Do not treat advisory overlays as provider truth.
````
