# Status M46 - High-Value SaaS OpenAPI Curation

| Item | State | Notes |
|---|---|---|
| Curation set freeze | `[+]` | Frozen service list: Slack, GitHub, GitLab, Jira, Trello, Stripe, Notion, Asana, Box, Discord, Twilio, SendGrid, Xero, Zendesk, Zoho, OpenAI, PagerDuty, Shopify, ServiceNow, QuickBooks, HubSpot, Airtable, Salesforce, and Linear. |
| Existing coverage review | `[+]` | Reviewed catalog rows, aliases, spec refs, artifact paths, overlays, quirks, and manual follow-ups for the full frozen set. Added focused catalog tests so this high-value set remains covered. |
| Gap fill | `[+]` | Reconfirmed official OpenAPI/Swagger/OpenAPI-index coverage for the 18 machine-readable providers. No new stable public provider-wide OpenAPI artifact was found for Airtable, Linear, QuickBooks, Salesforce, ServiceNow, or Shopify. |
| Auth/security reconciliation | `[+]` | Reconciled OAuth/API-key/header-token/GraphQL bearer metadata through existing security classifications and overlays; refreshed docs references and review notes for the official-doc-only providers. |
| Advisory overlay decisions | `[+]` | Retained docs-derived endpoint overlays for Airtable, QuickBooks, Salesforce, ServiceNow, and Shopify. Retained Linear's explicit no-endpoint-overlay decision because a generic GraphQL POST overlay would hide native operation semantics. |
| Catalog/reporting checks | `[+]` | Verified provider lookup, `catalog specs`, `catalog advisory`, `catalog stats`, refresh review, and catalog quality output for the curation set. Fixed built-in resolution so human-doc-only providers preserve curated primary-doc order instead of alphabetic spec-ref order. |
| Downstream notes | `[+]` | Downstream OpenUdon/udon follow-up remains metadata review only: official-doc-only providers need user/generated OpenAPI or reviewed advisory overlays, and no runtime credential/account/workspace choice moves into `apitools`. |
| Verification | `[+]` | `go test ./catalog`, full tests, vet, `catalog check`, catalog CLI smoke checks, review passes, and diff checks were run for the final M46 batch. |

## Review Notes

- M46 did not change the public product boundary: catalog metadata remains
  discovery/review guidance only.
- Resolution now preserves authored source-reference priority within each
  source kind. This keeps primary official docs ahead of auth, versioning, or
  GraphQL-adjacent references when no machine-readable OpenAPI exists.
- Official docs reviewed during this milestone include Shopify REST Admin API,
  Linear GraphQL, ServiceNow REST API Explorer OpenAPI export, QuickBooks
  Online Accounting API, Airtable Web API/PAT docs, Salesforce REST/OpenAPI
  generation notes, and HubSpot's public OpenAPI spec index.

## Boundary Notes

- The curation set is OpenAPI/Swagger/OpenAPI-index metadata work, not provider
  execution support.
- `apitools` must not execute SaaS API operations, fetch OAuth tokens, resolve
  credentials, or choose workspaces, organizations, tenants, accounts, or
  enterprise scopes.
- `../n8n` is only a priority signal for service selection.
