# Status M58 - Catalog Security Metadata Audit

| Item | State | Notes |
|---|---|---|
| Scope freeze after M57 | `[+]` | M57 was complete and committed before M58 audit implementation started. |
| Catalog-wide security inventory | `[+]` | Added offline `CatalogSecurityAuditReport` and `apitools catalog security-audit` for durable provider disposition, source-family, auth status, spec refs, overlay IDs, registered artifacts, and source-note coverage. |
| OpenAPI/Swagger security audit | `[+]` | The audit inspects local OpenAPI/Swagger `components.securitySchemes` or `securityDefinitions`, root `security`, operation-level `security`, operation counts, and undeclared requirement scheme names. |
| Native source-family audit | `[+]` | Native Smithy, Google Discovery, and Dropbox Stone rows remain classified through catalog security classifications/overlays; M58 adds no conversion behavior. |
| Human-docs/advisory overlay audit | `[+]` | Docs-only/docs-derived rows resolve through existing source-backed security overlays or explicit incomplete/anonymous decisions in the audit disposition buckets. |
| User/tenant-specific decisions | `[+]` | Account, tenant, cluster, instance, or user-specific rows remain `present-incomplete-reviewed` where no portable built-in credential binding is safe. |
| Overlay/classification remediation | `[+]` | Added source-backed overlays for Matrix signedRequest, SyncroMSP apiKey, TheHive 5 ApiKey/Basic/Session, and Zulip basicAuth scheme mismatches found by the artifact audit. |
| Advisory/security-report verification | `[+]` | `catalog security-audit`, `catalog security-report`, and tests confirm 292 providers classify into complete-upstream, complete-via-overlay, present-incomplete-reviewed, or intentionally-anonymous; queued re-review is currently zero. |
| Deep milestone review | `[+]` | Reviewed M58 diff against the acceptance criteria; no runtime credential behavior, provider API execution, conversion behavior, or third-party connector dependency was introduced. |
| Verification | `[+]` | Passed `go test ./...`, `GOWORK=off go test ./...`, `go vet ./...`, `GOWORK=off go vet ./...`, catalog check/security-audit/security-report/advisory/inspect/stats/overlay-view/help smoke checks, `../openudon` tests, `../udon` tests, and public/private `git diff --check`. |

## Source Notes

- M58 must not begin until M57 is complete and committed.
- The audit target is the durable built-in provider catalog after M57, not
  third-party integration marketplaces or historical n8n-derived priority data.
- OpenAPI/Swagger rows must be checked for both reusable security definitions
  and root/operation security requirements before they can be considered
  complete.
- Native Smithy, Discovery, Stone, and other source-family rows should preserve
  native auth metadata and must not reintroduce source-to-OpenAPI conversion as
  an `apitools` runtime contract.
- Rows whose auth posture is inherently account, tenant, cluster, instance, or
  user specific should keep explicit no-portable-overlay or user-provided
  metadata decisions instead of generic built-in overlays.

## Audit Snapshot

- `catalog security-audit` covers 292 durable providers.
- Provider dispositions after remediation: 49 complete-upstream, 202
  complete-via-overlay, 37 present-incomplete-reviewed, and 4
  intentionally-anonymous.
- Security report auth statuses remain 56 complete, 195 overlay-required, 37
  present-incomplete, and 4 intentionally-anonymous.
- Local artifact inspection remains advisory evidence: some cached artifacts
  still omit or mismatch native OpenAPI security metadata, but every durable
  provider row now has an explicit catalog disposition and source-backed
  classification or overlay.

## Boundary Checks

- Do not start M58 audit or remediation while M57 is still open.
- Do not call provider APIs, fetch OAuth tokens, resolve credentials, read
  local credential stores, sign requests, or choose accounts, tenants, clusters,
  regions, or environments.
- Do not treat n8n, Workato, or other integration catalogs as provider truth.
- Do not copy third-party connector code, generated metadata, credential
  behavior, descriptions, icons, examples, or workflow behavior.
- Keep all security metadata advisory and credential-value-free.
