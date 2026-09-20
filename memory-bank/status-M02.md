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
