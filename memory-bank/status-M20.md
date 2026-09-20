# Status M20 - Human-Docs Overlay Completion And Catalog Stats

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Finish the known human-docs overlay triage, make catalog coverage statistics
repeatable, and record the next non-OpenAPI and expansion direction before
adding more providers.

## Direction

REST-shaped human-docs providers can receive small OpenAPI-shaped advisory
endpoint overlays when official docs expose enough endpoint instructions. The
overlays remain review artifacts only. GraphQL-first providers should wait for
GraphQL-aware metadata instead of receiving a misleading single-endpoint REST
wrapper.

## Tasks

| Item | State | Notes |
|---|---|---|
| Catalog stats added | `[+]` | Added `apitools catalog stats` with provider protocol counts, artifact registry kind counts, refresh validation buckets, and byte totals. |
| REST-shaped overlay gaps closed | `[+]` | Added tracked advisory overlays for Databricks, Jenkins, Sentry, Splunk, and Telegram, with registry entries and a shared M20 overlay builder. |
| No-overlay decisions recorded | `[+]` | Recorded Linear and Monday.com as GraphQL-first no-overlay cases until GraphQL-aware classification/introspection work exists. |
| Overlay structural tests added | `[+]` | Added a test that loads tracked advisory overlay artifacts into the operation index and rejects empty/invalid operation coverage. |
| Non-OpenAPI converter notes documented | `[+]` | Added public notes for Discovery, Smithy, Stone, OpenAPI index, human docs, and GraphQL-first follow-up direction. |
| Expansion queue documented | `[+]` | Added source-first criteria for the next provider expansion batch instead of adding unfrozen services ad hoc. |
| Verification completed | `[+]` | Ran Go tests/vet, diff check, catalog quality, catalog stats/specs/help smoke checks, overlay structure checks, and sibling `openudon`/`udon` tests. |

## Boundary Checks

- Do not execute provider API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not treat advisory endpoint overlays as official provider OpenAPI.
- Do not force GraphQL APIs into REST-shaped endpoint overlays.
