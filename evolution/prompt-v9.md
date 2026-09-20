# Prompt V9 - Source Coverage Expansion After n8n Review

Use `../n8n/packages/nodes-base/nodes` as a priority signal to plan the next
first-class API source coverage sequence, without adding n8n runtime
compatibility or UWS source-type changes.

Group the agreed source families into logical milestones:

- AWS Smithy service expansion for additional AWS nodes.
- Google Discovery service expansion for additional Google nodes.
- Microsoft Graph and Azure OpenAPI service expansion.
- High-value SaaS OpenAPI curation and reconciliation.
- Dropbox Stone parser/no-parser decision.

Keep Smithy and Discovery native. Do not reintroduce OpenAPI-shaped conversion
as an `apitools` runtime contract.
