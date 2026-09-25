# Retired milestone M18 - Parseable Artifact Auth Overlays

**Milestone.** M18
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M18.md
**Source specification.** tabilet/memory-bank/milestone.md#m18---parseable-artifact-auth-overlays
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M18 - Parseable Artifact Auth Overlays

**Goal.** Keep parseable-but-strict-invalid official OpenAPI/Swagger artifacts
usable as downloaded review evidence while exposing auth/security metadata
through explicit advisory overlays.

**Scope.**

- Add built-in security overlays for parseable-invalid official OpenAPI/Swagger
  artifacts whose auth metadata cannot be safely relied on through strict
  import parsing.
- Preserve existing selected refresh behavior that downloads and registers
  parseable-invalid artifacts as review-only local cache evidence.
- Keep overlay statuses review-oriented (`present-incomplete` or
  `overlay-required`) instead of promoting strict-invalid artifacts as complete
  provider truth.
- Ensure provider advisory, security report, resolve, and overlay inspection
  views surface the overlay IDs and source-backed auth/security notes.
- Preserve boundaries: no provider API execution, credential resolution,
  request signing, account selection, or automatic metadata acceptance.

**Acceptance.** Parseable-invalid official OpenAPI/Swagger artifacts remain
downloadable review artifacts, catalog auth/security guidance is available via
overlays, and catalog quality plus Go verification remain clean.
````

## Status record

````markdown
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
````
