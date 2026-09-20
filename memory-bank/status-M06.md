# Status M06 - Catalog Data Expansion

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Expand the provider catalog with a curated, source-backed batch of popular
workflow services. M6 should increase useful coverage without weakening the
metadata-only boundary established in M2-M5.

The feature should answer these questions without live provider calls:

- Which additional popular services should be reviewed next?
- Which reviewed services have official OpenAPI or machine-readable specs?
- Which reviewed services need user-provided OpenAPI documents?
- Which reviewed services have complete, incomplete, or overlay-required
  auth/security metadata?

## Initial Provider Batch

Use this provider batch unless review discovers a blocker that should be noted
in this file before replacement:

```text
asana
box
calendly
clickup
discord
dropbox
github
gitlab
google-calendar
google-sheets
microsoft-graph
notion
salesforce
shopify
stripe
```

`../n8n/packages/nodes-base/nodes` is a priority signal only. It must not be
used as runtime compatibility evidence.

## Frozen Batch Rationale

The M6 provider batch is frozen as of May 18, 2026. Each provider has a sibling
`../n8n/packages/nodes-base/nodes` directory, which is only a priority signal
for common workflow usage. The batch intentionally spans productivity, files,
developer tooling, messaging, CRM, commerce, payments, and large platform APIs.

| Provider ID | Display Name | Category | Priority Signal |
|---|---|---|---|
| asana | Asana | project-management | `../n8n/packages/nodes-base/nodes/Asana` |
| box | Box | files | `../n8n/packages/nodes-base/nodes/Box` |
| calendly | Calendly | scheduling | `../n8n/packages/nodes-base/nodes/Calendly` |
| clickup | ClickUp | project-management | `../n8n/packages/nodes-base/nodes/ClickUp` |
| discord | Discord | messaging | `../n8n/packages/nodes-base/nodes/Discord` |
| dropbox | Dropbox | files | `../n8n/packages/nodes-base/nodes/Dropbox` |
| github | GitHub | developer-tools | `../n8n/packages/nodes-base/nodes/Github` |
| gitlab | GitLab | developer-tools | `../n8n/packages/nodes-base/nodes/Gitlab` |
| google-calendar | Google Calendar | scheduling | `../n8n/packages/nodes-base/nodes/Google/Calendar` |
| google-sheets | Google Sheets | spreadsheet | `../n8n/packages/nodes-base/nodes/Google/Sheet` |
| microsoft-graph | Microsoft Graph | productivity | `../n8n/packages/nodes-base/nodes/Microsoft` |
| notion | Notion | knowledge-management | `../n8n/packages/nodes-base/nodes/Notion` |
| salesforce | Salesforce | crm | `../n8n/packages/nodes-base/nodes/Salesforce` |
| shopify | Shopify | commerce | `../n8n/packages/nodes-base/nodes/Shopify` |
| stripe | Stripe | payments | `../n8n/packages/nodes-base/nodes/Stripe` |

## Service Tasks

Each service row is a review and commit unit once implementation begins. A row
is complete only after the full curation loop is done for that service:
official spec acquisition attempt, auth/security review, security overlay if
needed, advisory endpoint+security overlay if no official OpenAPI exists,
tracked builder/asset updates, cache manifest registration, durable catalog
metadata, tests, and notes.

| Service | State | Official Spec | Security | Advisory Overlay | Registry | Notes |
|---|---|---|---|---|---|---|
| asana | `[+]` | Official OpenAPI from `Asana/openapi` saved to ignored `catalog-openapi-cache/openapi/asana-openapi-v1.yaml`; repo commit `e698f9acf281833c71add9261f0f49f296da6003`, ETag recorded. | Complete: spec includes bearer PAT and OAuth2 schemes with root/operation security. | Not needed. | `register_catalog_artifacts.go` records `content_path` and artifact manifest row. | Promoted to durable provider with official spec and docs refs. |
| box | `[+]` | Official OpenAPI from `box/box-openapi` saved to ignored `catalog-openapi-cache/openapi/box-platform-openapi-v3.json`; repo commit `f5fd732661ae7a6d066c71869a59d967e712ec46`, ETag recorded. | Complete: spec includes OAuth2 authorization code scheme and root security. | Not needed. | `register_catalog_artifacts.go` records `content_path` and artifact manifest row. | Promoted to durable provider; note records Box versioned public API specs introduced in 2025. |
| calendly | `[+]` | No downloadable official OpenAPI was found from official docs, obvious provider URLs, or Stoplight export probes. Official API reference is hosted at `developer.calendly.com/api-docs` and points to `calendly.stoplight.io/docs/api-docs`. | Overlay required: Calendly docs describe personal access tokens and OAuth 2.1 bearer access tokens; no official OpenAPI security scheme is recorded. | Generated docs-derived non-official overlay at tracked `catalog-openapi-cache/advisory-overlays/calendly-public-api-overlay.json`. | `register_catalog_artifacts.go` records the advisory overlay path and builder path. | Promoted to durable provider with human docs refs; user OpenAPI need classified likely. |
| clickup | `[+]` | Official ReadMe-hosted OpenAPI saved to ignored `catalog-openapi-cache/openapi/clickup-api-v2-reference.json`; OpenAPI 3.1.0, version `2.0`, ETag and last-modified recorded. | Present-incomplete: spec has `Authorization` header API key/root security; OAuth authorization-code details are documented separately. | Not needed. | `register_catalog_artifacts.go` records `content_path` and artifact manifest row. | Promoted to durable provider with official spec, auth docs refs, and 3.1/v2-v3 quirk note. |
| discord | `[+]` | Official OpenAPI from `discord/discord-api-spec` saved to ignored `catalog-openapi-cache/openapi/discord-api-v10-openapi.json`; repo commit `ac55939ed657bbb4b69cbe02e24abffc79223026`, ETag recorded. | Complete: spec includes bot token and OAuth2 security schemes with operation-level security. | Not needed. | `register_catalog_artifacts.go` records `content_path` and artifact manifest row. | Promoted to durable provider; note records Discord public preview warning and docs-over-spec guidance. |
| dropbox | `[+]` | Official Dropbox API v2 Stone spec repository saved as ignored `catalog-openapi-cache/openapi/dropbox-api-spec-main.tar.gz`; repo commit `6c68656caafdff65a30943fe186a1800b9f35bbc`, archive ETag recorded. No official OpenAPI document is recorded. | Overlay required: official Stone spec is not OpenAPI, and OAuth scopes/security need OpenAPI mapping for OpenAPI-only consumers. | Generated docs/Stone-derived non-official overlay at tracked `catalog-openapi-cache/advisory-overlays/dropbox-core-api-overlay.json`. | `register_catalog_artifacts.go` records the Stone archive `content_path`, advisory overlay path, and builder path. | Promoted to durable provider with `dropbox-stone` machine spec kind, OAuth/security overlay, and content endpoint quirk notes. |
| github | `[+]` | Official OpenAPI from `github/rest-api-description` saved to ignored `catalog-openapi-cache/openapi/github-rest-api-openapi.json`; repo commit `f5d3342150d3748e7307c81639635706f8338a12`, ETag recorded. | Overlay required: the official OpenAPI has no security schemes, root security, or operation security; official docs describe bearer tokens and some basic-auth app endpoints. | Not needed because official OpenAPI exists. | `register_catalog_artifacts.go` records `content_path` and artifact manifest row. | Promoted to durable provider with official spec, auth docs refs, and security overlay. |
| gitlab | `[+]` | Official Swagger/OpenAPI v2 YAML from `gitlab-org/gitlab` saved to ignored `catalog-openapi-cache/openapi/gitlab-openapi-v2.yaml`; repo commit `3ef57e7c312105c20fce8b97e5ddc4c60f2d23ed`, ETag and SHA-256 recorded. | Present-incomplete: spec includes `PRIVATE-TOKEN` header/query security definitions but no root security; official docs also describe OAuth bearer and other access-token forms. | Not needed because official OpenAPI/Swagger exists. | `register_catalog_artifacts.go` records `content_path` and artifact manifest row. | Promoted to durable provider; user OpenAPI need classified possible because GitLab docs say only some endpoints are currently documented with OpenAPI. |
| google-calendar | `[+]` | Official Google Calendar API Discovery document saved to ignored `catalog-openapi-cache/google-discovery/google-calendar-discovery-v3.json`; Discovery revision `20260510` recorded. No official OpenAPI document is recorded. | Overlay required: Discovery is not OpenAPI, and OAuth scopes need OpenAPI mapping. | Generated Discovery/docs-derived non-official overlay at tracked `catalog-openapi-cache/advisory-overlays/google-calendar-v3-overlay.json`. | `register_catalog_artifacts.go` records the Discovery `content_path`, advisory overlay path, and builder path. | Promoted to durable provider with Google Discovery spec kind and Calendar OAuth/security overlay. |
| google-sheets | `[+]` | Official Google Sheets API Discovery document saved to ignored `catalog-openapi-cache/google-discovery/google-sheets-discovery-v4.json`; Discovery revision `20260511` recorded. No official OpenAPI document is recorded. | Overlay required: Discovery is not OpenAPI, and OAuth scopes need OpenAPI mapping. | Generated Discovery/docs-derived non-official overlay at tracked `catalog-openapi-cache/advisory-overlays/google-sheets-v4-overlay.json`. | `register_catalog_artifacts.go` records the Discovery `content_path`, advisory overlay path, and builder path. | Promoted to durable provider with Google Discovery spec kind and Sheets OAuth/security overlay. |
| microsoft-graph | `[+]` | Official OpenAPI v1.0 document from `microsoftgraph/msgraph-metadata` saved to ignored `catalog-openapi-cache/openapi/microsoft-graph-v1-openapi.yaml`; repo commit `ecfbb0d49b8bf7f7cb13310726a5feb3697dd42f`, ETag recorded. | Overlay required: official OpenAPI lacks security schemes; Microsoft docs describe Microsoft identity platform access tokens in the `Authorization: Bearer` header. | Not needed because official OpenAPI exists. | `register_catalog_artifacts.go` records `content_path` and artifact manifest row. | Promoted to durable provider with official spec, auth docs refs, and large-spec/caching quirk note. |
| notion | `[+]` | Official OpenAPI from Notion developer docs saved to ignored `catalog-openapi-cache/openapi/notion-api-openapi.json`; ETag, last-modified, and served-version recorded. | Complete: official OpenAPI includes bearer auth security scheme and root security requirement; docs confirm `Authorization: Bearer` and `Notion-Version` headers. | Not needed. | `register_catalog_artifacts.go` records `content_path` and artifact manifest row. | Promoted to durable provider with official spec and Notion-Version quirk note. |
| salesforce | `[+]` | No stable public downloadable official OpenAPI for core REST API is recorded; official release notes describe org/configuration-dependent OpenAPI generation beta. | Overlay required: official REST docs use OAuth 2.0 bearer access tokens, but no stable public OpenAPI security scheme is recorded. | Generated docs-derived non-official overlay at tracked `catalog-openapi-cache/advisory-overlays/salesforce-rest-core-overlay.json`. | `register_catalog_artifacts.go` records the advisory overlay path and builder path. | Promoted to durable provider with human docs refs; user OpenAPI need classified likely. |
| shopify | `[+]` | No stable public downloadable official OpenAPI is recorded from official Shopify REST Admin docs/search. | Overlay required: docs describe Admin API access tokens in `X-Shopify-Access-Token` and access scopes, but no official OpenAPI security scheme is recorded. | Generated docs-derived non-official overlay at tracked `catalog-openapi-cache/advisory-overlays/shopify-admin-rest-overlay.json`. | `register_catalog_artifacts.go` records the advisory overlay path and builder path. | Promoted to durable provider with human docs refs; user OpenAPI need classified likely. |
| stripe | `[+]` | Official Stripe OpenAPI latest GA public spec saved to ignored `catalog-openapi-cache/openapi/stripe-latest-openapi-spec3.json`; repo commit `c29a3d372d2cdd1085a6d1cc18bd968a84129c7b`, ETag recorded. | Complete: official OpenAPI includes basic and bearer HTTP auth schemes with root security requirements. | Not needed. | `register_catalog_artifacts.go` records `content_path` and artifact manifest row. | Promoted to durable provider; note records that `latest/openapi.spec3.json` is preferred over legacy `openapi/spec3.json`. |

## Cross-Service Acceptance Tasks

These rows are milestone-level acceptance checks. They should be flipped only
after the relevant per-service rows above prove the batch-wide condition.

| Item | State | Notes |
|---|---|---|
| Provider batch frozen | `[+]` | Final M6 provider IDs, display names, categories, and priority-only n8n rationale recorded above. |
| Candidate records added | `[+]` | Added M6 candidate metadata with priority-only n8n evidence; local fixture evidence remains limited to real `../try-n8n/reducibility/specs` fixtures. |
| Official spec acquisition attempted | `[+]` | Every M6 provider row records an official OpenAPI/Swagger, official non-OpenAPI machine spec, or no-stable-downloadable-spec result; downloaded review artifacts were saved under ignored `catalog-openapi-cache/`. |
| Auth/security completeness reviewed | `[+]` | Each M6 provider has a complete, present-incomplete, or overlay-required security classification based on reviewed machine specs and official auth docs. |
| Security overlays added where needed | `[+]` | Added source-backed catalog security overlays for M6 providers whose official specs were absent, non-OpenAPI, or incomplete for auth/security. |
| Advisory overlays generated where needed | `[+]` | Generated non-official endpoint+security overlays for M6 providers without a stable public official OpenAPI where docs or Discovery coverage could support a useful advisory subset. |
| Overlay builders saved per service | `[+]` | Saved deterministic `//go:build ignore` builders for every generated advisory overlay under tracked `catalog-openapi-cache/overlay-builders/`. |
| Cache manifest paths recorded | `[+]` | Registered saved official documents by `content_path` and advisory overlays by overlay/builder path in `register_catalog_artifacts.go`; SQLite cache stores paths rather than duplicated document bodies. |
| Durable provider entries added | `[+]` | Promoted all frozen M6 provider IDs into built-in provider catalog entries with reviewed sources, verified dates, revisions, license notes, and quirks. |
| Auth/security metadata added | `[+]` | Added source-backed security classifications or overlays for every M6 provider. |
| User OpenAPI need classified | `[+]` | Classified each M6 provider as likely, possible, or not-expected based on official spec availability and coverage. |
| Catalog tests added | `[+]` | Expanded catalog tests for deterministic IDs, alias lookup, provider promotion, spec availability, security statuses, and copy safety across the expanded provider set. |
| Documentation updated if needed | `[+]` | No public command behavior or boundary docs changed; the M6 status file records implementation and acceptance details. |
| Verification completed | `[+]` | Ran `go test ./...`, `go vet ./...`, `git diff --check`, `go run ./cmd/apitools search --help`, `go run ./cmd/apitools import --help`, `(cd ../openudon && go test ./...)`, and `(cd ../udon && go test ./...)`. |

## Boundary Checks

- Do not execute API operations.
- Do not fetch credentials, tokens, or secrets.
- Do not sign requests.
- Do not choose accounts, tenants, workspaces, or runtime environments.
- Do not vendor upstream specs into this repository.
- Do not commit downloaded specs, Google Discovery documents, or SQLite caches
  from `catalog-openapi-cache/`.
- Do commit generated advisory overlays and service-specific builder programs
  from `catalog-openapi-cache/` as catalog assets.
- Do not treat n8n node presence as runtime compatibility evidence.
- Do not treat local overlays as upstream provider truth.
