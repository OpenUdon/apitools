# Retired milestone M74 - Review Contract Corrections

**Milestone.** M74
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M74.md
**Source specification.** tabilet/memory-bank/milestone.md#m74---review-contract-corrections
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M74 - Review Contract Corrections

**Goal.** Correct the focused post-release regressions in catalog security
scope, discovery API compatibility, protobuf format fallback, authoring prompt
ownership and budgeting, truncation reporting, and cache artifact integrity.

**Scope.** Keep catalog security evidence on the selected spec or provider-wide
scope; restore both historical project-URL import methods while exposing the
bounded report separately; fall back from heuristic proto-text parsing to
binary descriptors; deep-copy caller-owned authoring metadata; apply custom
prompt budgets once; mark selected-context overflow as truncated; and reject
zero-byte catalog artifacts before persistence.

**Dependencies.** M73. These changes restore documented public and safety
contracts without introducing a new product direction or compatibility break.

**Parallel ownership.** M74 is a serialized cross-cutting correction milestone
covering catalog resolution, discovery, gRPC/protobuf parsing, root authoring
and inventory, SQLite cache integrity, and their focused tests.

**Downstream impacts.** Existing discovery callers compile unchanged; OpenUdon,
Ramen, and Udon retain their current API integrations and gain corrected local
behavior when updated to the reviewed revision.

**Acceptance.** Each reported case has a deterministic regression test; full
workspace and standalone tests, vet, staticcheck, generated-catalog checks, CLI
help, sibling consumer tests, and diff checks pass; the milestone receives a
final review before closure.
````

## Status record

````markdown
# Status M74 - Review Contract Corrections

**Acceptance.** All seven focused review findings are corrected without a new
compatibility break, their regressions are covered, and coordinated verification
passes.

| Item | State | Notes |
|---|---|---|
| Correct catalog security scope and restore discovery compatibility | `[+]` | Commit `8c40d36` keeps selected-spec auth unknown when only another spec has evidence, preserves the aggregate inspection report, restores both historical project-URL import signatures, and adds `ImportProjectURLsReport`; focused catalog/root tests pass. |
| Preserve valid protobuf and authoring inputs | `[+]` | Commit `03e4dc1` falls back to binary descriptors after a failed text heuristic, deep-copies nested ranking inputs, applies caller prompt budgets once during inventory construction, and marks selected-context overflow as truncated; focused parser/root tests pass. |
| Reject empty catalog artifacts | `[+]` | Commit `4557def` rejects zero-byte regular files before catalog-artifact persistence and verifies no poison row remains; focused SQLite cache tests pass. |
| Run milestone review and coordinated verification | `[+]` | Final review found no remaining cross-spec, compatibility, ownership, budgeting, truncation, parser-fallback, or empty-artifact gap. Uncached focused tests, full workspace/standalone apitools tests and vet, staticcheck, catalog generator/quality, CLI help, OpenUdon/Ramen/Udon workspace and standalone tests/vet, and apitools/tofu diff checks pass. |

## Boundary Checks

- Discovered metadata remains non-executable and credential-free.
- Remote and local source safety defaults are unchanged.
- Existing exported discovery method signatures remain source-compatible.
- No tfconfig changes are required.
````
