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
