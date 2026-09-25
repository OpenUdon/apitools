# Retired milestone M52 - n8n De-Emphasis And License-Risk Reduction

**Milestone.** M52
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M52.md
**Source specification.** tabilet/memory-bank/milestone.md#m52---n8n-de-emphasis-and-license-risk-reduction
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M52 - n8n De-Emphasis And License-Risk Reduction

**Goal.** Remove active n8n-specific catalog surfaces and make provider
curation source-first before release.

**Scope.**

- Remove the n8n Public API provider/candidate row, TimeSaved workflow helper
  row, artifact registration, auth classification, and related tests.
- Remove the `catalog n8n-gap-report` CLI/API and exported report types.
- Remove n8n node-directory evidence from built-in candidates and avoid
  active references to n8n node roots, `nodes-base`, or `try-n8n` fixtures.
- Rewrite public docs and active memory-bank guidance around provider-owned
  OpenAPI, Smithy, Discovery, Stone, official indexes, and docs-derived
  overlays.
- Keep historical milestone/status/evolution records as chronology, but mark
  M49-M51 as superseded for active implementation guidance.
- Do not copy third-party integration code, generated node metadata,
  credential behavior, descriptions, icons, or runtime workflow behavior.

**Acceptance.** Code, tests, public docs, catalog registry, and active
memory-bank guidance no longer position n8n as a provider source or priority
driver; remaining n8n mentions are historical records or explicit hygiene
tests; Go, catalog, consumer, and diff checks pass.
````

## Status record

````markdown
# Status M52 - n8n De-Emphasis And License-Risk Reduction

| Item | State | Notes |
|---|---|---|
| n8n report API removed | `[+]` | Removed the exported catalog n8n gap-report API, tests, CLI command, README usage, and tech-stack smoke command. |
| n8n provider rows removed | `[+]` | Removed n8n Public API and TimeSaved candidate/provider/auth metadata, alias expectations, and artifact registration. |
| Candidate evidence cleaned | `[+]` | Removed built-in n8n node-directory priority evidence and local `try-n8n` fixture seeding from candidates. Built-in candidate evidence is now source-first. |
| Public docs updated | `[+]` | Rewrote catalog queue and non-OpenAPI protocol notes to avoid n8n-specific provider prioritization and to require provider-owned or protocol-owned source artifacts. |
| Active memory superseded | `[+]` | Added this milestone and updated active product, architecture, tech-stack, and milestone notes so M49-M51 are historical only. |
| Verification | `[+]` | Passed full Go, catalog, consumer, diff, and text-audit checks. |

## Boundary Notes

- Do not copy third-party integration code, generated node metadata,
  credential behavior, descriptions, icons, or runtime workflow behavior.
- Runtime connector catalogs may suggest user demand, but they are not provider
  truth and must not drive catalog rows without provider-owned source evidence.
- Historical memory-bank and evolution records may mention n8n as chronology;
  active public docs, code, and release guidance should not depend on it.
````
