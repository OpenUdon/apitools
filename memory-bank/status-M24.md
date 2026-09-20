# Status M24 - Provider Node Expansion Batch 7

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Continue broadening provider-like n8n node catalog coverage before standalone
converter work. M24 focuses on a mixed utility, data, and public-API batch with
source-backed metadata and auth/security classification.

## Frozen Batch

Batch 7 services:

- UptimeRobot
- PostHog
- NocoDB
- SeaTable
- QuickChart
- Marketstack
- CoinGecko
- Reddit
- HackerNews
- NASA
- OpenThesaurus
- OneSimpleApi

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
| Batch list frozen | `[+]` | Frozen Batch 7 services are UptimeRobot, PostHog, NocoDB, SeaTable, QuickChart, Marketstack, CoinGecko, Reddit, HackerNews, NASA, OpenThesaurus, and OneSimpleApi. |
| UptimeRobot curated | `[+]` | Added UptimeRobot candidate/provider metadata, official API v3 OpenAPI reference, saved ignored review artifact registration, and complete bearer API-key auth classification; saved artifact is parseable OpenAPI but strict-invalid for review. |
| PostHog curated | `[+]` | Added PostHog candidate/provider metadata, official Cloud API OpenAPI schema reference, saved ignored review artifact registration, and present-incomplete personal API key/project-token auth classification; saved artifact is parseable OpenAPI 3.1 but strict-invalid for review. |
| NocoDB curated | `[+]` | Added NocoDB candidate/provider metadata, official public v2 and v3 OpenAPI references, saved ignored review artifact registrations, and present-incomplete token auth classification; saved artifacts are parseable OpenAPI 3.1 but strict-invalid for review. |
| SeaTable curated | `[+]` | Added SeaTable candidate/provider metadata, official v6.2 Authentication/User/Base/File OpenAPI references from the public OpenAPI index, saved ignored review artifact registrations, and complete bearer token auth classification; User/Base artifacts are parseable OpenAPI but strict-invalid for review. |
| QuickChart curated | `[+]` | Added QuickChart candidate/provider metadata, official human docs refs, intentionally anonymous auth classification, and a docs-derived advisory overlay for chart, QR, and graph rendering endpoints; M24 review found no stable public official OpenAPI at common probes. |
| Marketstack curated | `[+]` | Added Marketstack candidate/provider metadata, official APILayer-hosted Marketstack API v2 OpenAPI reference from the official docs payload, saved ignored review artifact registration, and complete access_key query auth classification. |
| CoinGecko curated | `[+]` | Added CoinGecko candidate/provider metadata, official human docs refs, overlay-required Demo/Pro API-key auth metadata, and a docs-derived advisory overlay for selected v3 market-data endpoints; M24 review found no stable public official OpenAPI at common probes. |
| Reddit curated | `[+]` | Added Reddit candidate/provider metadata, official generated API and Developer Platform docs refs, overlay-required OAuth bearer metadata, and a docs-derived advisory overlay for selected oauth.reddit.com endpoints; M24 review found no stable public official OpenAPI at common probes. |
| HackerNews curated | `[+]` | Added Hacker News candidate/provider metadata, official Firebase API docs refs, intentionally anonymous auth classification, and a docs-derived advisory overlay for item, user, story-list, and updates endpoints. |
| NASA curated | `[+]` | Added NASA candidate/provider metadata, official Open APIs docs refs, overlay-required api_key query metadata, and a docs-derived advisory overlay for APOD, Mars Rover Photos, NeoWs, EPIC, and DONKI endpoints; M24 review found no stable public official OpenAPI at common probes. |
| OpenThesaurus curated | `[+]` | Added OpenThesaurus candidate/provider metadata, official API docs refs, intentionally anonymous auth classification, and a docs-derived advisory overlay for synonym search with JSON/XML options. |
| OneSimpleApi curated | `[+]` | Added OneSimpleApi candidate/provider metadata, official docs/API-token refs, overlay-required token query metadata, and a docs-derived advisory overlay for selected utility endpoints; M24 review found no stable public official OpenAPI at common probes. |
| Batch quality review completed | `[+]` | Full M24 verification passed for catalog validation, deterministic tests, full repository tests, vet, diff whitespace checks, catalog stats, and inspect output for all Batch 7 providers. |

## Boundary Checks

- Do not treat n8n node presence as runtime compatibility evidence.
- Do not vendor downloaded official specs, Discovery documents, or SQLite
  caches into the public repository.
- Do not execute API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not treat advisory overlays as provider truth.
- Do not start Google Discovery, Smithy, Stone, GraphQL, AI/vector-specific, or
  other converter work in this milestone.
