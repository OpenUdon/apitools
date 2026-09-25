# Retired milestone M55 - Enterprise Applications API Source Coverage

**Milestone.** M55
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M55.md
**Source specification.** tabilet/memory-bank/milestone.md#m55---enterprise-applications-api-source-coverage
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M55 - Enterprise Applications API Source Coverage

**Goal.** Add source-first catalog coverage for high-value enterprise
application APIs: NetSuite, SAP S/4HANA, SAP SuccessFactors, Oracle Fusion
Cloud Applications, and Workday.

**Scope.**

- Add NetSuite as the first implementation target. Record SuiteTalk REST Web
  Services and REST API Browser metadata as official OpenAPI 3.0 source
  evidence, but treat account-specific metadata-catalog URLs as user-provided
  metadata rather than built-in downloadable artifacts.
- Add SAP S/4HANA and SAP SuccessFactors as separate provider entries or a
  clearly grouped SAP family, based on official SAP Business Accelerator Hub
  source references. Record whether each reviewed surface is OpenAPI, OData,
  SOAP, or docs-only, and defer parser work if OData semantics would be hidden
  by a generic REST overlay.
- Add Oracle Fusion Cloud Applications as the Oracle application target, not a
  generic Oracle provider. Use official Fusion REST docs and `/describe`
  metadata as source evidence, while treating tenant-specific describe URLs as
  user-provided metadata.
- Add Workday provider metadata from official Workday SOAP/WWS and REST API
  documentation. Defer deep implementation unless review finds stable public
  REST/OpenAPI artifacts or a future WSDL/SOAP source milestone is opened.
- Add auth/security classifications and advisory notes only from official
  sources. Do not sign in to tenants, resolve OAuth clients, fetch tokens,
  call metadata endpoints, or infer customer-specific modules.

**Acceptance.** NetSuite, SAP S/4HANA, SAP SuccessFactors, Oracle Fusion Cloud
Applications, and Workday have catalog/candidate metadata, official source
refs, source-family classifications, auth/security status, and clear
user-provided metadata boundaries; any OData, WSDL/SOAP, or tenant-describe
parser decisions are explicitly recorded; Go, catalog, consumer, and diff
checks pass.
````

## Status record

````markdown
# Status M55 - Enterprise Applications API Source Coverage

| Item | State | Notes |
|---|---|---|
| Source freeze | `[+]` | Froze NetSuite, SAP S/4HANA, SAP SuccessFactors, Oracle Fusion Cloud Applications, and Workday from official provider documentation and support references only. |
| NetSuite cataloged | `[+]` | Added SuiteTalk REST Web Services, REST API Browser, resource metadata, metadata catalog, and auth references; account-specific OpenAPI 3.0 metadata-catalog URLs remain user-provided. |
| SAP cataloged | `[+]` | Added SAP S/4HANA and SAP SuccessFactors source metadata from SAP Business Accelerator Hub, SAP Help, and SAP support references; OData/SOAP/OpenAPI/docs-only families are recorded without parser work. |
| Oracle Fusion cataloged | `[+]` | Added Oracle Fusion Cloud Applications from official Fusion REST docs, HCM endpoint docs, `/describe` metadata, and object metadata; no generic Oracle provider row was created. |
| Workday cataloged | `[+]` | Added Workday metadata from official SOAP/WWS and REST API directories; deep implementation is deferred to future WSDL/SOAP or stable REST/OpenAPI source work. |
| Parser decisions recorded | `[+]` | Recorded no-overlay/no-parser decisions for OData, WSDL/SOAP, and tenant-specific describe/metadata catalogs in public docs and advisory overlay decisions. |
| Security overlay inspection | `[+]` | Reconfirmed that M55 providers need tenant/account/source-family-specific auth metadata rather than portable built-in security overlays. Future user-provided or parser-generated overlays should attach to exact exported OpenAPI/OData/WSDL/describe artifacts. |
| Verification | `[+]` | Go, catalog, consumer, go mod tidy diff, text-audit, and diff checks passed before completing M55. |

## Source Notes

- NetSuite is the best first target because official docs describe SuiteTalk
  REST Web Services metadata using OpenAPI 3.0 and account-specific metadata
  catalog URLs.
- SAP source review should distinguish S/4HANA and SuccessFactors surfaces and
  avoid hiding OData semantics behind generic REST-shaped overlays.
- Oracle work should target Oracle Fusion Cloud Applications specifically,
  using official Fusion REST docs and `/describe` metadata where available.
- Workday has official SOAP/WWS and REST documentation, but public OpenAPI
  evidence is weaker; start with provider metadata and source-family notes.
- M55 records no stable public built-in downloadable provider-wide OpenAPI
  artifact for these five enterprise application rows. Useful machine-readable
  metadata is account, tenant, product-family, module, or protocol specific and
  must be user-provided until a future native parser milestone is scoped.

## Security Overlay Inspection

| Provider | Built-in overlay needed? | Notes |
|---|---|---|
| NetSuite | User/account-specific overlay needed | SuiteTalk REST auth depends on account, integration record, OAuth 2.0 or token-based auth setup, roles, and account-specific metadata-catalog URLs. A generic built-in overlay would not be enough for safe credential binding. |
| Oracle Fusion Cloud Applications | User/tenant-specific overlay needed | Fusion REST and `/describe` metadata are product-family, tenant, customization, version, and security-context specific. Credential binding should be attached to exported tenant metadata, not a provider-wide row. |
| SAP S/4HANA | Artifact-specific overlay needed | SAP Business Accelerator Hub surfaces include OpenAPI, OData, SOAP, and docs-only APIs with differing OAuth/API-key/basic patterns. Overlay work should be per selected artifact and should preserve OData/SOAP semantics. |
| SAP SuccessFactors | Tenant/module-specific overlay needed | SuccessFactors APIs are OData/SOAP oriented and tenant/module configured; useful auth metadata should attach to exported or selected API metadata rather than a generic provider-wide overlay. |
| Workday | Tenant/source-family-specific overlay needed | Workday WWS and REST directories vary by tenant, product area, and security context. Future WSDL/SOAP or REST source handling should emit or accept explicit security overlays for the selected tenant artifact. |

## Boundary Checks

- Treat tenant/account-specific metadata endpoints as user-provided source
  URLs, not built-in downloadable artifacts.
- Do not sign in to enterprise tenants, resolve OAuth clients, fetch tokens,
  call metadata endpoints, or infer customer-specific enabled modules.
- Do not implement OData, WSDL/SOAP, or Oracle/NetSuite describe parsers unless
  a milestone explicitly scopes metadata-only local parsing.
- Keep SAP, Oracle Fusion, NetSuite, and Workday auth/security notes advisory
  until source review confirms operation-level metadata.
````
