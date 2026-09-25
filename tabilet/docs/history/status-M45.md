# Retired milestone M45 - Microsoft OpenAPI Service Expansion

**Milestone.** M45
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M45.md
**Source specification.** tabilet/memory-bank/milestone.md#m45---microsoft-openapi-service-expansion
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M45 - Microsoft OpenAPI Service Expansion

**Goal.** Make Microsoft Graph and Azure REST coverage first-class for the
Microsoft services represented in `../n8n`, using official Microsoft OpenAPI
sources where available.

**Scope.**

- Map n8n Microsoft nodes to existing or new catalog providers: Outlook,
  OneDrive, Teams, Excel, SharePoint, To Do, Entra, Graph Security,
  Microsoft Graph core, Azure Storage, and Azure Cosmos DB.
- Keep Microsoft Graph-backed services tied to official Graph OpenAPI metadata
  from `microsoftgraph/msgraph-metadata`; add service-level aliases and source
  notes so users can find each n8n-visible service without duplicating provider
  truth.
- Add Azure REST spec references from official Azure REST API specs for Azure
  Storage and Cosmos DB when source review confirms stable OpenAPI artifacts.
- Review auth/security metadata for OAuth scopes, tenant/account boundaries,
  and Azure resource scopes as metadata only; do not acquire tokens or choose
  tenants, subscriptions, accounts, or regions.
- Register saved official artifacts or explicit no-artifact decisions for
  offline refresh review.

**Acceptance.** Microsoft Graph service aliases and Azure REST providers are
discoverable with official OpenAPI provenance, auth/security classifications
are source-backed, catalog quality passes, and no Microsoft API execution,
token acquisition, tenant selection, subscription selection, or account
selection is introduced.
````

## Status record

````markdown
# Status M45 - Microsoft OpenAPI Service Expansion

| Item | State | Notes |
|---|---|---|
| Service mapping review | `[+]` | Mapped n8n Microsoft nodes to Graph-backed Outlook, OneDrive, Teams, Excel, SharePoint, To Do, Entra, and Graph Security coverage, plus Azure Storage Blob and Azure Cosmos DB REST coverage. |
| Microsoft Graph metadata | `[+]` | Confirmed official `microsoftgraph/msgraph-metadata` v1.0 OpenAPI source and added service-level catalog aliases/source notes that point back to the shared Graph artifact. |
| Azure REST specs | `[+]` | Added official Azure REST API Specs references for Azure Blob Storage data plane and Cosmos DB Resource Manager Swagger; Cosmos DB SQL data-plane exact item/query shape remains documented as a possible generated/user-spec need. |
| Provider/catalog rows | `[+]` | Added Microsoft service-level providers, Azure providers, aliases, categories, source notes, verified dates, quirks, and user OpenAPI need classifications. |
| Review artifacts | `[+]` | Registered Azure artifacts and the shared Microsoft Graph artifact for service-level rows; added refresh-review reuse for duplicate registered artifact paths. |
| Security metadata | `[+]` | Added Microsoft Graph OAuth scope overlays and Azure Storage/Cosmos auth overlays as metadata only, with no token, signing, tenant, subscription, account, database, or endpoint resolution. |
| Catalog/reporting checks | `[+]` | Verified provider lookup, `catalog specs`, `catalog inspect`, `catalog stats`, refresh review, and catalog quality output for the batch. |
| Verification | `[+]` | Ran `go test ./catalog`, `go test ./...`, `go vet ./...`, `go run ./cmd/apitools catalog check`, catalog CLI smoke checks, and `git diff --check`. |

## Review Notes

- Microsoft 365 service rows intentionally share `microsoft-graph-v1-openapi`
  as provider truth instead of inventing service-specific OpenAPI artifacts.
- Azure Storage and Cosmos DB source artifacts are Swagger 2.0 documents from
  Azure REST API Specs, so protocol stats count them under Swagger.
- The Cosmos DB Resource Manager artifact requires an offline refresh
  correction for official external Swagger refs; this correction is review-only
  and does not alter provider truth.

## Boundary Notes

- Microsoft Graph and Azure REST sources remain OpenAPI metadata for review and
  downstream planning only.
- `apitools` must not execute Microsoft APIs, acquire tokens, resolve
  credentials, or choose tenants, subscriptions, accounts, resources, or
  regions.
- `../n8n` is only a priority signal for service selection.
````
