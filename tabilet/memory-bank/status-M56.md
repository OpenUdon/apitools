# Status M56 - OpenAPI-First Provider Batch

| Item | State | Notes |
|---|---|---|
| Source freeze | `[+]` | Frozen on 2026-05-21 from provider-owned official OpenAPI/Swagger documents, provider-owned repositories, and official docs only; no Workato or third-party integration metadata was used as provider truth. |
| Adyen cataloged | `[+]` | Added Adyen Checkout Service v72 from the official Adyen OpenAPI repository, plus API Explorer and credential docs, with payment-operation safety notes. |
| DocuSign cataloged | `[+]` | Added DocuSign eSignature REST v2.1 Swagger from the official DocuSign OpenAPI-Specifications repository, plus eSignature docs and auth boundaries. |
| Auth0 cataloged | `[+]` | Added Auth0 Management API OpenAPI schema from official Auth0 docs and recorded tenant-domain, scope, client, and token boundaries. |
| Confluence cataloged | `[+]` | Added Atlassian Confluence Cloud REST v2 OpenAPI from official Atlassian docs; Data Center/server APIs remain out of scope. |
| Acumatica cataloged | `[+]` | Added Acumatica REST metadata as tenant-generated Swagger/OpenAPI requiring user-provided exported metadata; no built-in provider-wide spec or endpoint overlay. |
| BigCommerce cataloged | `[+]` | Added BigCommerce Catalog Products, Orders, and Carts OpenAPI references from the official BigCommerce API specs repository and official docs. |
| Cisco Meraki cataloged | `[+]` | Added Cisco Meraki Dashboard API v1 OpenAPI from the official Meraki repository and Cisco developer docs. |
| Confluent Cloud cataloged | `[+]` | Added Confluent Cloud organization, Kafka REST, and Schema Registry OpenAPI references from the official Confluent Cloud SDK repository and docs. |
| Security overlay inspection | `[+]` | Inspected durable M56 OpenAPI/Swagger security sections and implemented the needed built-in overlays. Adyen, DocuSign, and Confluence Cloud now have reviewed overlays and report `complete`; Auth0, BigCommerce, Cisco Meraki, and Confluent Cloud report `complete` from official OpenAPI security metadata without built-in overlays; Acumatica remains user-generated metadata only. |
| Verification | `[+]` | Passed `go test ./...`, `GOWORK=off go test ./...`, `go vet ./...`, `GOWORK=off go vet ./...`, `GOWORK=off go mod tidy -diff`, catalog check/specs/stats/advisory/inspect/security-report/refresh-report smoke checks, `go run ./cmd/apitools search --help`, `go run ./cmd/apitools import --help`, consumer tests in `../openudon` and `../udon`, source-hygiene text audit, and public/private diff checks. |

## Source Notes

- Adyen publishes provider-owned OpenAPI definitions and an API Explorer for
  its API families.
- DocuSign publishes official REST API Swagger/OpenAPI specifications, with
  eSignature as the likely first surface.
- Auth0 Management API docs identify an OpenAPI-backed schema; tenant domains
  and credentials remain downstream concerns.
- Confluence should start with Atlassian Cloud REST APIs; Data Center/server
  APIs are a separate compatibility surface.
- Acumatica exposes Swagger/OpenAPI per ERP instance endpoint, so built-in
  catalog rows should describe source shape and require user-provided metadata
  URLs for concrete imports.
- BigCommerce, Cisco Meraki, and Confluent Cloud have strong official
  OpenAPI/API-reference evidence suitable for source-first cataloging.
- M56 catalog rows use `AuthStatusPresentIncomplete` unless existing official
  OpenAPI security metadata or a reviewed built-in overlay covers reusable
  credential placement. Credential binding, tenant, account, scope, and
  permission decisions still stay downstream.

## Security Overlay Inspection

| Provider | Built-in overlay needed? | Notes |
|---|---|---|
| Adyen | Complete | Implemented `adyen-checkout-v72-auth-overlay`. Checkout Service v72 defines `ApiKeyAuth` and `BasicAuth`, and most operations have security requirements; the overlay records reviewed operation security for `post-validateShopperId`, the one operation lacking operation-level security. |
| DocuSign | Complete | Implemented `docusign-esignature-rest-v2-1-auth-overlay`. eSignature REST v2.1 Swagger has no `securityDefinitions`, root security, or operation security, so the overlay records OAuth bearer placement from official auth docs while preserving environment, account, consent, and token issuance boundaries downstream. |
| Auth0 | Complete, no overlay needed | Management API OpenAPI defines bearer and OAuth2 client-credentials schemes, root security, and operation security for all inspected operations. Tenant domains, scopes, API clients, and access-token issuance remain downstream. |
| Confluence Cloud | Complete | Implemented `confluence-cloud-rest-v2-auth-overlay`. The OpenAPI defines basic and OAuth2 schemes and most operations carry security requirements; the overlay records operation-level security for the user-email and Atlassian Connect dynamic-module operations that lacked operation security. |
| Acumatica | No portable built-in overlay | Acumatica Swagger/OpenAPI is generated per ERP instance, endpoint, tenant, version, and customization. Credential binding needs user-provided exported metadata and likely a user-provided overlay. |
| BigCommerce | Complete, no overlay needed | Catalog Products, Orders, and Carts OpenAPI specs define `X-Auth-Token` as an API-key header and root security. Store hashes, OAuth/private-app credentials, scopes, and channel availability remain downstream. |
| Cisco Meraki | Complete, no overlay needed | Dashboard API OpenAPI defines `X-Cisco-Meraki-API-Key` and bearer auth schemes with root security. Missing operation-level security inherits root requirements, so no supplemental overlay is needed for basic credential placement. |
| Confluent Cloud | Complete, no overlay needed | Organization, Kafka REST, and Schema Registry specs define Basic API-key and OAuth2 schemes with operation-level security across inspected operations. Tenant environment, cluster, and endpoint selection remain downstream concerns. |

## Boundary Checks

- Do not use Workato or any other integration marketplace as a metadata source;
  at most, use such lists as uncopied market-demand signals.
- Do not log in to provider tenants, create API keys, fetch OAuth tokens, call
  APIs, inspect customer resources, or infer account-specific modules.
- Treat payment, e-signature, network administration, identity management,
  ecommerce, ERP, and streaming-control APIs as metadata-only surfaces.
- Preserve tenant-specific/versioned artifact boundaries in spec refs and
  advisory notes instead of normalizing them away.
