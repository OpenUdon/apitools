# Retired milestone M47 - Dropbox Stone Source Evaluation

**Milestone.** M47
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M47.md
**Source specification.** tabilet/memory-bank/milestone.md#m47---dropbox-stone-source-evaluation
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M47 - Dropbox Stone Source Evaluation

**Goal.** Decide whether Dropbox Stone deserves first-class source parsing or
whether the current Stone-derived advisory overlay remains sufficient for
OpenUdon-style authoring.

**Scope.**

- Review the official `dropbox/dropbox-api-spec` Stone source and the existing
  Dropbox Stone-derived advisory overlay from M37.
- Compare the n8n Dropbox node surface with the reviewed advisory subset and
  identify missing high-value operations, auth/security metadata, and schema
  shapes.
- Prototype or specify the minimal native Stone parser shape only if source
  review shows broad reusable value beyond the current advisory overlay.
- If a parser is not justified, record a no-parser decision with the exact
  coverage and maintenance rationale.
- If a parser is justified, open a follow-up implementation milestone with
  native Stone model APIs, fixtures, catalog protocol handling, and downstream
  consumer expectations.
- Do not execute Dropbox operations, resolve credentials, fetch tokens, choose
  Dropbox teams/accounts, or treat Stone-derived advisory metadata as official
  OpenAPI.

**Acceptance.** Maintainers have a source-reviewed decision for Dropbox Stone:
either a documented no-parser/no-expansion rationale with the current advisory
overlay retained, or a scoped follow-up implementation milestone for native
Stone parsing. Catalog quality and diff checks pass.
````

## Status record

````markdown
# Status M47 - Dropbox Stone Source Evaluation

| Item | State | Notes |
|---|---|---|
| Source review | `[+]` | Re-reviewed the official `dropbox/dropbox-api-spec` Stone source, Dropbox HTTP docs, OAuth guide, and the M37 Stone-derived advisory overlay. The upstream `main` revision remains `6c68656caafdff65a30943fe186a1800b9f35bbc`, matching the cataloged review metadata. |
| n8n surface comparison | `[+]` | n8n covers Dropbox file/folder/search operations. The reviewed advisory subset covers account, list/list-continue, metadata, upload, download, and sharing links; missing n8n-visible operations are copy, delete, move, create-folder-v2, search-v2, and search-continue-v2. |
| Parser value assessment | `[+]` | Native Stone parsing is not justified for the current milestone: Stone is provider-specific, useful durable lowering would need route, alias, union, import, data-type, and auth attribute handling, and current catalog value comes from official source provenance plus the reviewed advisory overlay. |
| No-parser decision | `[+]` | Retain Dropbox Stone as official non-OpenAPI source metadata and advisory-overlay provenance. Do not add a native parser, general Stone lowering, or advisory overlay expansion in M47. |
| Follow-up milestone draft | `[+]` | No parser follow-up milestone was opened. If Dropbox workflow demand grows, the next scoped work should be targeted advisory overlay expansion for the missing n8n-visible operations before a general Stone parser is reconsidered. |
| Catalog reconciliation | `[+]` | Fixed provider resolution so Dropbox prefers the `dropbox-api-stone-spec` machine source over the authored human-doc references. Protocol classification, artifact paths, and overlay provenance remain unchanged. |
| Verification | `[+]` | Passed `go test ./catalog`, `go test ./...`, `go vet ./...`, `go run ./cmd/apitools catalog check`, Dropbox advisory/spec/stats smoke checks, `git diff --check`, and `git -C ../tofu diff --check -- apitools`. |

## Decision Notes

- Current reviewed advisory coverage remains intentionally limited to
  `users/get_current_account`, `files/list_folder`,
  `files/list_folder/continue`, `files/get_metadata`, `files/upload`,
  `files/download`, `sharing/create_shared_link_with_settings`, and
  `sharing/list_shared_links`.
- The n8n comparison identified useful future overlay candidates, but those
  gaps do not require a generic Stone parser before there is broader
  multi-provider value.
- Dropbox Stone remains a machine-readable source reference, not official
  OpenAPI. Stone-derived advisory metadata remains source-backed guidance for
  OpenAPI-only consumers.

## Boundary Notes

- Dropbox Stone remains official non-OpenAPI source metadata unless a future
  parser milestone explicitly changes the implementation scope.
- Stone-derived advisory metadata must not be treated as official OpenAPI.
- `apitools` must not execute Dropbox operations, resolve credentials, fetch
  tokens, or choose Dropbox teams/accounts.
````
