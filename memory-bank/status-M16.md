# Status M16 - Refresh Artifact Reconciliation

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Use the M15 offline refresh report to reconcile local catalog artifact state
before adding more provider nodes. M16 should make successful refresh artifacts
reproducibly registered and leave remaining validation failures explicit.

## Direction

This milestone may run selected `catalog refresh` commands for built-in spec
references, but it must not promote catalog metadata automatically. Downloaded
official artifacts and SQLite cache contents remain ignored local review state;
tracked changes are limited to registry source and milestone memory.

## Tasks

| Item | State | Notes |
|---|---|---|
| Refresh report reviewed | `[+]` | Ran `catalog refresh-report --cache-dir catalog-openapi-cache --cache catalog-openapi-cache/cache.sqlite --as-of 2026-05-19` and classified missing registrations plus invalid saved artifacts. |
| Missing M13 refs refreshed | `[+]` | Selected refresh succeeded for AWS Lambda, AWS S3, AWS SNS Smithy JSON models, Grafana OpenAPI, and Snowflake OpenAPI. |
| Registry source reconciled | `[+]` | Added AWS Smithy JSON artifact rows to the tracked artifact registry source; Grafana and Snowflake rows were already present. |
| Remaining refresh failures documented | `[+]` | Bitbucket, CircleCI, Cloudflare, Elastic, and Supabase still fail current OpenAPI refresh validation and remain unregistered rather than saved as accepted artifacts. |
| Existing invalid artifacts documented | `[+]` | Existing saved artifacts for Asana, Box, ClickUp, GitHub, Jira Cloud, Notion, Okta, PagerDuty, Pipedrive, Slack, Trello, and Zendesk still fail current strict OpenAPI validation and require later importer/upstream-shape review. |
| Verification completed | `[+]` | Ran registry rebuild, refresh report, Go tests/vet, diff check, catalog quality, and CLI smoke checks. |

## Boundary Checks

- Do not execute provider API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not promote verified dates, revisions, or security classifications automatically.
- Do not vendor downloaded official specs or SQLite cache files into the public repository.
- Do not treat validation-failed artifacts as accepted provider truth.
