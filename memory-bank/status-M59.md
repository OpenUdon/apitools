# Status M59 - Workato/Tray Signal Provider Batch

| Item | State | Notes |
|---|---|---|
| Source-signal freeze | `[+]` | Used `tray_workato_connectors.xlsx`, Workato, Tray, and similar integration lists only as market-demand signals. No connector metadata, code, icons, credential behavior, descriptions, or examples were copied into catalog truth. |
| Existing coverage reconciliation | `[+]` | Checked built-in providers, aliases, categories, and prior milestones before adding new rows for Canva, LaunchDarkly, Klaviyo, Attio, Fivetran, Anthropic, Amplitude, Ashby, Braze, and Postman. |
| Provider-owned spec review | `[+]` | Recorded provider-owned OpenAPI/docs evidence, source URL stability notes, license/source notes, and durable result decisions for all M59 providers. |
| Canva cataloged | `[+]` | Recorded official Canva Connect API OpenAPI YAML and docs. Upstream security metadata is complete; no endpoint or security overlay is needed. |
| LaunchDarkly cataloged | `[+]` | Recorded official LaunchDarkly REST API v2 OpenAPI JSON and docs. Upstream security metadata is complete; no endpoint or security overlay is needed. |
| Klaviyo cataloged | `[+]` | Recorded official Klaviyo stable OpenAPI JSON from `klaviyo/openapi` and docs. Upstream security metadata is complete; no endpoint or security overlay is needed. |
| Attio cataloged | `[+]` | Recorded official Attio REST API OpenAPI document and docs. Upstream OAuth security metadata is complete; no endpoint or security overlay is needed. |
| Fivetran cataloged | `[+]` | Recorded official Fivetran REST API V1 OpenAPI JSON and docs. Upstream Basic auth security metadata is complete; no endpoint or security overlay is needed. |
| Anthropic cataloged | `[+]` | Recorded official Anthropic API docs and an advisory x-api-key security overlay. No endpoint overlay is registered because prompt/model execution needs explicit policy or a reviewed source artifact. |
| Amplitude cataloged | `[+]` | Recorded official Amplitude API docs and an advisory security overlay for family-specific Basic, `api_key`, and bearer auth forms. No broad endpoint overlay is registered. |
| Ashby cataloged | `[+]` | Recorded official Ashby API docs and an advisory Basic auth overlay. No endpoint overlay is registered because recruiting PII/RPC-style operations need narrower review. |
| Braze cataloged | `[+]` | Recorded official Braze API docs and an advisory bearer REST API key overlay. No endpoint overlay is registered because messaging/profile/export operations are permission and region sensitive. |
| Postman cataloged | `[+]` | Recorded official Postman API docs and an advisory `X-API-Key` overlay. Public workspaces/collections are not treated as provider-wide source truth; no endpoint overlay is registered. |
| Security/advisory audit | `[+]` | Applied the M58 rule: OpenAPI providers were classified complete after security scheme/requirement review, and docs-only providers received advisory security overlays plus explicit no-endpoint-overlay decisions. |
| Verification | `[+]` | Full repo tests/vet, `GOWORK=off` checks, catalog quality, stats, advisory/security smoke checks, consumer tests, diff checks, and deep review completed in the implementation turn. |

## Source Notes

- M59 is sourced from provider-owned evidence. Third-party integration
  marketplaces and `tray_workato_connectors.xlsx` are prioritization signals
  only.
- Prefer stable, provider-owned OpenAPI/Swagger documents when available.
- Use docs-derived overlays only for narrow, reviewed subsets where official
  documentation exposes endpoint, parameter, response, and auth shape clearly.
- Record an explicit no-overlay decision when a provider has useful API
  metadata but no portable, safe, provider-owned operation shape for a built-in
  advisory artifact.

## Boundary Checks

- Do not sign in to provider products, fetch tokens, resolve credentials, call
  provider APIs, execute AI prompts, create customer data, send messages, or
  mutate provider resources.
- Do not scrape authenticated documentation or copy third-party connector
  metadata, generated code, credential behavior, icons, examples, descriptions,
  or workflow behavior.
- Keep all auth/security metadata advisory and credential-value-free.
- Treat provider-owned API docs and artifacts as untrusted metadata; validate
  and summarize shape only.
