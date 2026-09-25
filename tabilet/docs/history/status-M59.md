# Retired milestone M59 - Workato/Tray Signal Provider Batch

**Milestone.** M59
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M59.md
**Source specification.** tabilet/memory-bank/milestone.md#m59---workatotray-signal-provider-batch
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M59 - Workato/Tray Signal Provider Batch

**Goal.** Add a source-first provider batch identified from the local
`tray_workato_connectors.xlsx` market-demand signal: Canva, LaunchDarkly,
Klaviyo, Attio, Fivetran, Anthropic, Amplitude, Ashby, Braze, and Postman.

**Scope.**

- Treat Workato and Tray.io rows only as prioritization signals. Do not copy
  connector metadata, descriptions, auth behavior, operation lists, categories,
  examples, icons, or workflow behavior.
- For each provider, first look for durable provider-owned OpenAPI or Swagger
  artifacts. Use official OpenAPI where available for Canva, LaunchDarkly,
  Klaviyo, Attio, and Fivetran if source review confirms stable importable
  artifacts.
- Use docs-derived advisory overlays only when official docs support a narrow
  reviewed subset and no suitable provider-owned OpenAPI artifact is found.
  Likely docs-overlay decisions include Anthropic, Amplitude, Ashby, Braze,
  and Postman unless source review finds durable official specs.
- Add catalog/candidate/spec/source-family/auth/security/advisory metadata in
  the same slice for each provider, including M58-style security inspection
  and overlay decisions.
- Do not sign in, fetch tokens, call provider APIs, create records, send AI
  prompts, export customer data, or scrape authenticated documentation.

**Acceptance.** The ten M59 providers have durable catalog/candidate metadata,
official source refs, auth/security classification, source-family
classification, and either a provider-owned OpenAPI artifact, a tracked
docs-derived advisory overlay, or an explicit no-overlay decision; catalog
advisory, inspect, stats, security-report, security-audit, quality, Go,
consumer, and diff checks pass.
````

## Status record

````markdown
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
````
