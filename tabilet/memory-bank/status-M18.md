# Status M18 - Parseable Artifact Auth Overlays

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

For official OpenAPI/Swagger artifacts that are downloadable and parseable but
strict-invalid, continue treating the downloaded artifact as review evidence and
carry auth/security guidance through explicit catalog overlays.

## Direction

Parseable-invalid artifacts remain non-import-ready. Their auth/security
metadata is surfaced as advisory overlay metadata so downstream resolution and
inspection can see credential requirements without treating failed validation
artifacts as provider truth.

## Tasks

| Item | State | Notes |
|---|---|---|
| Auth overlays added | `[+]` | Added present-incomplete auth review overlays for parseable-invalid Asana, Bitbucket, Box, CircleCI, Cloudflare, Jira Cloud, Notion, Okta, PagerDuty, Pipedrive, Trello, and Zendesk Sunshine artifacts. |
| Existing overlays aligned | `[+]` | Updated existing ClickUp, GitHub, and Supabase overlay notes to call out parseable-invalid refresh status; existing Slack and Zendesk Support overlays remain review-only. |
| Source refs checked | `[+]` | Confirmed new official overlay source references return HTTP 200 as of 2026-05-19. |
| Security expectations updated | `[+]` | Security report now surfaces overlay IDs and present-incomplete status for parseable-invalid artifacts whose auth metadata would otherwise have been classified complete from source inspection. |
| Verification completed | `[+]` | Ran Go tests/vet, diff check, catalog quality, refresh-report, advisory/overlay-view checks, CLI smoke checks, and sibling `openudon`/`udon` tests. |

## Boundary Checks

- Do not execute provider API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not promote parseable-invalid artifacts as import-ready metadata.
- Do not treat advisory overlays as provider truth.
