# Status M32 - Provider Node Expansion Batch 15

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Continue broadening provider-like n8n node catalog coverage before standalone
converter work. M32 focuses on content, publishing, and community services with
source-backed metadata and auth/security classification.

## Frozen Batch

Batch 15 services:

- Disqus
- Medium
- npm
- Spotify
- Storyblok
- Strapi
- WordPress
- Wufoo
- YOURLS
- Raindrop
- Toggl
- Travis CI

`../n8n/packages/nodes-base/nodes` is a priority signal only. It is not runtime
compatibility evidence, and no n8n behavior claims should be recorded.

## Per-Service Curation Loop

For each service, try official machine-readable sources first, save and
register ignored official review artifacts when retrieved, review
auth/security completeness, add candidate/provider metadata, add security
classification or overlays when needed, generate docs-derived advisory overlays
only for useful REST-shaped subsets, record provenance fields, and run focused
tests plus `go run ./cmd/apitools catalog check --as-of 2026-05-20` before the
single end-of-batch commit.

## Tasks

| Item | State | Notes |
|---|---|---|
| Batch list frozen | `[+]` | Frozen Batch 15 services are Disqus, Medium, npm, Spotify, Storyblok, Strapi, WordPress, Wufoo, YOURLS, Raindrop, Toggl, and Travis CI. |
| Disqus curated | `[+]` | Added candidate/provider rows, access_token/api_key auth overlay, and docs-derived overlay `disqus-api-overlay.json` for forums, threads, and posts. |
| Medium curated | `[+]` | Added candidate/provider rows, bearer auth overlay, and docs-derived overlay `medium-api-overlay.json` for current user, publications, and publishing posts. |
| npm curated | `[+]` | Added candidate/provider rows, bearer auth overlay, and docs-derived overlay `npm-registry-api-overlay.json` for package metadata, search, dist-tags, and whoami. |
| Spotify curated | `[+]` | Added candidate/provider rows and registered the official Spotify Web API OpenAPI schema; no docs-derived endpoint overlay needed. |
| Storyblok curated | `[+]` | Added candidate/provider rows, content-token and management bearer auth overlay, and docs-derived overlay `storyblok-api-overlay.json` for Content Delivery and Management story endpoints. |
| Strapi curated | `[+]` | Added candidate/provider rows, bearer JWT/API-token auth overlay, and docs-derived overlay `strapi-rest-api-overlay.json` for generic collection entry routes and local auth; precise schemas remain instance-generated. |
| WordPress curated | `[+]` | Added candidate/provider rows, Basic and WordPress.com bearer auth overlay, and docs-derived overlay `wordpress-rest-api-overlay.json` for core wp/v2 resources; site/plugin schemas remain variable. |
| Wufoo curated | `[+]` | Added candidate/provider rows, API-key Basic auth overlay, and docs-derived overlay `wufoo-api-v3-overlay.json` for forms, entries, reports, users, and webhooks. |
| YOURLS curated | `[+]` | Added candidate/provider rows, signature query auth overlay, and docs-derived overlay `yourls-api-overlay.json` for shorturl, expand, stats, and URL stats actions. |
| Raindrop curated | `[+]` | Added candidate/provider rows, OAuth bearer auth overlay, and docs-derived overlay `raindrop-api-overlay.json` for raindrops, collections, tags, filters, highlights, and user metadata. |
| Toggl curated | `[+]` | Added candidate/provider rows and registered the official Toggl Track Swagger artifact; auth classified present-incomplete because the reviewed Swagger defines BasicAuth but references undeclared OAuth2 on some operations. |
| Travis CI curated | `[+]` | Added candidate/provider rows, token/API-version header auth overlay, and docs-derived overlay `travis-ci-api-overlay.json` for repositories, builds, jobs, requests, branches, and environment variables. |
| Batch quality review completed | `[+]` | Final code review completed; provider inspect/advisory checks, tests, vet, diff checks, and catalog quality checks passed before commit. |

## Boundary Checks

- Do not treat n8n node presence as runtime compatibility evidence.
- Do not vendor downloaded official specs, Discovery documents, or SQLite
  caches into the public repository.
- Do not execute API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not treat advisory overlays as provider truth.
- Do not start standalone converter work in this milestone.
