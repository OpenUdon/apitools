# Status M61 - Advisory Overlay Linked-Docs Review

| Item | State | Notes |
|---|---|---|
| Overlay inventory freeze | `[+]` | Counted 139 tracked advisory overlay JSON artifacts and 139 registered `overlayArtifacts`; file and registry counts agree. |
| Curation rule documented | `[+]` | Added the linked-docs advisory overlay rule to `tech-stack.md` and scoped M61 in `milestone.md`. |
| OpenWeatherMap overlay review | `[+]` | Expanded `openweathermap-one-call-3-overlay` from One Call 3.0 only to include linked Geocoding API helpers: `/geo/1.0/direct`, `/geo/1.0/zip`, and `/geo/1.0/reverse`; preserved `lat`/`lon` guidance and added Geocoding API source refs. |
| Registry refresh | `[+]` | Re-ran the artifact registry builder after the OpenWeatherMap overlay update so the ignored local cache reports the new overlay digest. |
| Batch A overlay reviews | `[+]` | Reviewed advisory overlays `action-network` through `contentful` for linked official API docs and related operations. Expanded Aircall, Airtable, and Clockify with webhook operations; recorded covered/no-add decisions for the remaining Batch A overlays. |
| Batch B overlay reviews | `[+]` | Reviewed advisory overlays `copper` through `google-sheets` for linked official API docs and related operations. Expanded Google Calendar with free/busy query, event watch, and channel stop operations; recorded covered/no-add decisions for the remaining Batch B overlays. |
| Batch C overlay reviews | `[+]` | Reviewed advisory overlays `gotowebinar` through `mailjet` for linked official API docs and related operations. Expanded Help Scout with a focused Docs API subset, JotForm with report/folder operations, and Mailjet with template and event callback webhook operations; recorded covered/no-add decisions for the remaining Batch C overlays. |
| Batch D overlay reviews | `[+]` | Reviewed advisory overlays `mandrill` through `postmark` for linked official API docs and related operations. Expanded Mandrill with template and webhook management operations and Postmark with message stream operations; recorded covered/no-add decisions for the remaining Batch D overlays. |
| Batch E overlay reviews | `[+]` | Reviewed advisory overlays `profitwell` through `servicenow` for linked official API docs and related operations. Expanded Salesforce with Composite resources and ServiceNow with Import Set API operations; recorded covered/no-add decisions for the remaining Batch E overlays. |
| Batch F overlay reviews | `[+]` | Reviewed advisory overlays `shopify` through `zammad` for linked official API docs and related operations. Expanded Shopify with webhook and inventory-level operations; recorded covered/no-add decisions for the remaining Batch F overlays. |
| No-add decision ledger | `[+]` | `linked-doc-review.md` now records all Batch A through Batch F decisions, including explicit no-add decisions for linked docs that are conceptual, setup/auth only, duplicated, instance/generated, unstable, or outside the reviewed workflow. |
| Full verification | `[+]` | Passed catalog check; advisory spot checks for Shopify, Salesforce, Mandrill, ServiceNow, and OpenWeatherMap; Shopify materialization smoke test; `go test ./catalog`; `go test ./...`; `go vet ./...`; `GOWORK=off go test ./...`; `GOWORK=off go vet ./...`; `git diff --check`; `git -C ../tofu diff --check -- apitools`; and sibling `go test ./...` in `../openudon` and `../udon`. |

## Source Notes

- M61 reviews existing docs-derived advisory overlays, not provider-official
  OpenAPI documents.
- Official linked API docs should be followed when they contribute operations,
  parameters, response fields, helper lookup flows, or auth/security metadata.
- Linked-doc expansion should stay workflow-relevant; unrelated product APIs
  should be recorded as no-add decisions rather than folded into a broad,
  unfocused overlay.

## Boundary Checks

- Do not call provider APIs, create records, submit requests, fetch tokens,
  resolve credentials, or sign requests.
- Do not scrape authenticated documentation or use third-party integration
  catalogs as provider truth.
- Do not label advisory overlays as official OpenAPI.
- Keep builders and generated overlay JSON in sync.
