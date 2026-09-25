# Retired milestone M75 - Operation Lifecycle Ranking Ownership

**Milestone.** M75
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M75.md
**Source specification.** tabilet/memory-bank/milestone.md#m75---operation-lifecycle-ranking-ownership
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M75 - Operation Lifecycle Ranking Ownership

**Goal.** Own conservative API lifecycle sibling ranking beside the operation
summaries and source provenance it evaluates, removing that API-specific policy
from the generic Authoring engine.

**Scope.** Add `operationlifecycle` over `apitools.OperationSummary`, retain the
reviewed API-first scoring behavior with named constants, reject ambiguous
ties, and restrict `/upload` normalization to explicit Google Discovery
provenance. Migrate OpenUdon and Ramen consumers and remove the old Authoring
package. Do not add workflow semantics, source fetches, credential behavior,
account selection, or operation execution.

**Dependencies.** M74. Coordinated with Authoring I01/D01/M28 and the
OpenUdon/Ramen adapter changes.

**Parallel ownership.** M75 owns only `operationlifecycle/` in apitools plus
the explicitly coordinated downstream import bridges and memory-bank records.
It does not overlap parser, catalog, network, cache, or runtime behavior.

**Downstream impacts.** OpenUdon and Ramen must first consume a published
apitools revision containing the new package, then consume the coordinated
Authoring revision and pass standalone checks.

**Acceptance.** Full apitools test/vet and workspace consumer suites pass;
provenance regressions prove generic `/upload` paths are untouched; Authoring
contains no lifecycle-ranking package or apitools dependency; after publication
both consumers pin the new revisions and pass `GOWORK=off` test/vet.
````

## Status record

````markdown
# Status M75 - Operation Lifecycle Ranking Ownership

State: Complete

| Item | State | Notes |
|---|---|---|
| Add provenance-aware lifecycle ranking | `[+]` | Added `operationlifecycle` over `apitools.OperationSummary`, documented scoring constants, retained same-source/ambiguity behavior, and restricted `/upload` normalization to explicit Google Discovery provenance. |
| Migrate Authoring consumers | `[+]` | OpenUdon and Ramen now translate downstream context to source-bearing operation summaries; Authoring removes its misplaced package and public-surface entries. |
| Verify workspace integration | `[+]` | Apitools full test/vet, Authoring standalone/race, OpenUdon/Ramen full workspace suites, OpenUdon scorecard, compatibility, and diff checks pass. |
| Publish and pin coordinated revisions | `[+]` | Published Apitools `v0.0.0-20260820042238-d51b61ead067` before Authoring `v0.0.0-20260820042256-2f73e3526583`; OpenUdon and Ramen pin both without replacements and pass standalone `GOWORK=off` test/vet. |

## Boundary Checks

- Ranking consumes metadata only and never fetches or executes operations.
- Product goals influence optional lifecycle role selection but no workflow,
  credential, account, or desired-state semantics move into apitools.
- Provider-specific path handling requires explicit source provenance.

Commit tracking: Apitools implementation `d51b61e`; coordinated Authoring
`243121d`, `eee0c0d`, and `2f73e35`; OpenUdon adapter/pins `43294bf` and
`065bb84`; Ramen adapter/pins `330fd89` and `721b540`.
````
