# Retired milestone M02 - Candidate Service Inventory

**Milestone.** M02
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M02.md
**Source specification.** tabilet/memory-bank/milestone.md#m02---candidate-service-inventory
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M02 - Candidate Service Inventory

**Goal.** Build the first reviewed list of candidate services for a future core
OpenAPI catalog.

**Scope.**

- Define a candidate service record shape for inventory and review.
- Seed the initial candidates from `../try-n8n/reducibility/specs`, because
  those fixtures already match common OpenUdon advisory services.
- Use `../n8n/packages/nodes-base/nodes` only as a service-priority signal for
  later expansion.
- Classify each candidate by known local spec evidence, official OpenAPI status,
  official non-OpenAPI machine spec status, user-supplied-spec need, and
  auth/security review state.
- Avoid n8n runtime compatibility, node import support, or behavior coverage
  claims.

**Acceptance.** Downstream catalog work has a deterministic candidate service
inventory with source notes and review classifications, without introducing
catalog resolution, security overlays, API execution, or credential handling.
````

## Status record

````markdown
# Status M02 - Candidate Service Inventory

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Create a reviewed candidate service inventory before committing to durable
catalog entries, spec references, security overlays, or CLI behavior.

The inventory should answer these questions without live provider calls:

- Which popular workflow services should be considered first?
- Which candidates already have local OpenAPI fixture evidence?
- Which candidates have known official OpenAPI sources to verify later?
- Which candidates have official machine-readable specs that are not OpenAPI?
- Which candidates likely need user-supplied OpenAPI documents?
- Which candidates need auth/security review before becoming catalog entries?

Use n8n's built-in service list only as a practical prioritization signal for
common workflow services. Do not claim n8n runtime compatibility, node import
support, or complete n8n behavior coverage.

## Candidate Source Direction

Start with `../try-n8n/reducibility/specs`, because it already contains local
fixtures for the nine OpenUdon advisory services:

| Provider | Initial inventory note |
|---|---|
| Airtable | Local fixture exists; official OpenAPI status needs verification. |
| Gmail | Local fixture exists; official Google Discovery status needs verification. |
| Google Drive | Local fixture exists; official Google Discovery status needs verification. |
| HubSpot | Local fixture exists; official OpenAPI status needs verification. |
| Jira Cloud | Local fixture exists; official OpenAPI status needs verification. |
| OpenWeatherMap | Local fixture exists; official OpenAPI status needs verification. |
| PagerDuty | Local fixture exists; official OpenAPI status needs verification. |
| Slack | Local fixture exists; official OpenAPI source status needs verification. |
| Trello | Local fixture exists; official OpenAPI status needs verification. |

After this seed set, inspect `../n8n/packages/nodes-base/nodes` for broader
candidate prioritization only. Do not import n8n runtime semantics into
`apitools`.

## Candidate Record Direction

Candidate records should cover at least:

- stable candidate ID, display name, aliases, and source notes;
- category and workflow relevance notes;
- evidence source, such as local fixture, n8n priority signal, public catalog,
  official docs, or user report;
- local OpenAPI fixture availability;
- official OpenAPI status: known, needs-verification, unavailable, or unknown;
- official non-OpenAPI machine spec status, such as Google Discovery;
- likely user-supplied OpenAPI need;
- auth/security review state: not-reviewed, likely-incomplete, likely-complete,
  intentionally-anonymous, or unknown.

Candidate records are review inputs. They are not yet durable built-in catalog
entries.

## Tasks

| Item | State | Notes |
|---|---|---|
| Candidate record shape designed | `[+]` | Added `catalog.Candidate`, evidence, status, user OpenAPI need, and auth/security review types. |
| Nine-service seed inventory added | `[+]` | Airtable, Gmail, Drive, HubSpot, Jira, OpenWeatherMap, PagerDuty, Slack, Trello. |
| Local fixture evidence captured | `[+]` | References `../try-n8n/reducibility/specs` as local fixture evidence only. |
| Broader n8n candidate signal scoped | `[+]` | References n8n node paths as `priority-only` evidence, not runtime compatibility. |
| Candidate classification implemented | `[+]` | Captures OpenAPI/non-OpenAPI/user-supplied/auth review states. |
| Inventory tests added | `[+]` | Covers deterministic loading, IDs, aliases, copy behavior, and classification values. |
| Verification completed | `[+]` | Ran `go test ./...`, `go vet ./...`, `git diff --check`, CLI help smoke checks, and sibling `../openudon`/`../udon` tests. |

## Boundary Checks

- Do not execute API operations.
- Do not fetch credentials, tokens, or secrets.
- Do not sign requests.
- Do not choose accounts, tenants, workspaces, or runtime environments.
- Do not import n8n node runtime semantics.
- Do not treat local fixtures as official provider sources.
````
