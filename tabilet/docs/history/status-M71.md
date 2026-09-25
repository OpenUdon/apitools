# Retired milestone M71 - Parallel-Lane Harness Migration

**Milestone.** M71
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M71.md
**Source specification.** tabilet/memory-bank/milestone.md#m71---parallel-lane-harness-migration
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M71 - Parallel-Lane Harness Migration

**Goal.** Adopt the current lane-aware memory-bank contract without losing or
reclassifying the completed M01-M70 history.

**Scope.**

- Define permanent `M`, `C`, and `S` lane meanings, cross-lane priority,
  dependency, and safe-parallelism rules.
- Normalize the legacy M1-M9 filenames/headings to M01-M09 and record that
  explicit migration while leaving later historical IDs and evolution
  snapshots intact.
- Make every status file discoverable by the two-digit lane glob and every task
  row parseable as `Item | State | Notes` with a backticked second-column state.
- Keep later work as unnumbered candidate directions with deferral reasons and
  promotion triggers.
- Update bootstrap instructions and evolution v21; do not add a mandatory goal
  protocol or alter public `apitools` behavior.

**Dependencies.** None. This is private harness metadata only.

**Parallel ownership.** No implementation milestone is active concurrently.
Future `C` and `S` milestones may run in parallel when their sections name
non-overlapping files and no shared public contract dependency.

**Downstream impacts.** The account-level memory-bank API runner can discover
all canonical status files and choose actionable work across lanes. Public
`apitools`, OpenUdon, Udon, and UWS contracts are unchanged.

**Acceptance.** Every status filename matches
`status-[A-Z][0-9][0-9].md`; the index links exactly one file per milestone;
all state-bearing rows use a supported backticked marker in the second column;
the updated unattended runner reports no actionable rows without needing an
API call; `git diff --check` passes in both harness-visible repositories; and
the public Go test suite remains green.
````

## Status record

````markdown
# Status M71 - Parallel-Lane Harness Migration

**Acceptance.** All status files satisfy the canonical two-digit lane pattern,
their task rows are parseable by the updated API runner, milestone links and
lane rules are internally consistent, the runner finds no actionable rows
after completion, `git diff --check` passes in both repositories, and
`go test ./...` passes in `../apitools`.

| Item | State | Notes |
|---|---|---|
| Migrate the private harness to parallel status lanes | `[+]` | Preserved M01-M70 history; defined future catalog and source-tooling lanes; normalized filenames, tables, instructions, candidates, and evolution records; validated all 71 indexed files and 671 state rows; ran the no-action runner smoke test; and left public package behavior unchanged. |

## Boundary Checks

- `M` remains the permanent home of all legacy IDs; completed milestones are
  not reclassified into the new domain lanes.
- `C` and `S` classify future work by product domain, not by team, priority, or
  calendar phase.
- No public Go API, CLI, source artifact, provider metadata, credential policy,
  or runtime behavior changes in this milestone.

## Verification

- A structural audit validated 71 unique index links, canonical filenames, and
  671 parser-recognized state rows with no raw or unsupported table markers.
- `../skills/harness/tackle-memory-bank-api-loop --model lane-audit .`
  discovered all lane files and reported no actionable rows without an API
  request.
- `go test ./...` passed in `../apitools`.
- `git diff --check` passed in `../apitools` and for `apitools/` in `../tofu`.
````
