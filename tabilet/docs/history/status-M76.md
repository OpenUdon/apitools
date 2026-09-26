# Retired milestone M76 - CLI Usage Exit Contract

**Milestone.** M76
**Outcome.** completed
**Retired.** 2026-09-26
**Source status.** tabilet/memory-bank/status-M76.md
**Source specification.** tabilet/memory-bank/milestone.md#m76---cli-usage-exit-contract
**Evidence.** 0b19929f08b607a0d759a24ccb7759d076afd39a
**Worktree.** clean
**Review.** passed
**Review iterations.** 3
**Verification.** go build ./..., go vet ./..., go test ./... (apitools, all 17 packages, including focused search/import CLI tests); git diff --check; documented CLI smoke checks (help exits 0; missing/short-query, missing/malformed-URL/scheme, and missing-dir usage exit 2; private-host runtime failure exits 1) against the built binary. All passed at the evidence commit.
**Consolidated into.** no current-truth change; current facts already live in the M73 0/1/2 exit-contract documentation this milestone restored.

## Milestone specification

````markdown
## M76 - CLI Usage Exit Contract

**Goal.** Restore the documented CLI contract that usage failures exit 2 on
stderr for `search` and `import`.

**Scope.**

- Validate required flags (`--query`, `--url`, `--dir`), minimum query
  length, and URL syntax/scheme in the CLI layer before calling the library,
  returning exit 2 with usage on stderr.
- Keep runtime failures (network, validation of downloaded content, cache
  errors) at exit 1, and leave library error behavior unchanged.

**Review provenance.** "apitools Code Review" Pass 4 finding X1, source Medium,
local P2 (regression of the M73 0/1/2 contract); revalidated at
`e015231ef651bb92ee1aab9a4e8fc870d94ad020` with a clean worktree by running
the built CLI (all listed usage errors exit 1). Lineage: M73 (completed; not
reopened).

**Dependencies.** C03 is completed (retired; see the [retired C03 record](status-C03.md)). M76 edits distinct handlers in `cmd/apitools/main.go`; the dependency only kept shared-file ownership sequential and is now satisfied.

**Parallel ownership.** The `search` and `import` command handlers in
`cmd/apitools/main.go` and CLI contract tests.

**Downstream impacts.** Scripts that treated exit 1 as "any failure" still see
a non-zero exit; no consumer API changes.

**Acceptance.** CLI contract tests cover every listed usage error returning 2
with usage on stderr and runtime failures returning 1. `go test ./...`,
`go vet ./...`, `git diff --check`, and the documented CLI smoke checks pass.
````

## Status record

````markdown
# Status M76 - CLI Usage Exit Contract

State of each M76 milestone item. Update as items complete. See
[`milestone.md`](milestone.md) for milestone definitions and for the status ID
pattern that names this file.

Status markers:

| Symbol | Suggested Status | Interpretation |
|---|---|---|
| `[ ]` | Pending | Item is actionable when its dependencies pass. |
| `[+]` | Completed | Item finished or done. |
| `[~]` | In Progress | Item is currently being worked, under the ownership rules in `milestone.md`. |
| `[!]` | Blocked | A current unresolved blocker still prevents progress. |
| `[X]` | Cancelled | Item is no longer needed. |
| `[-]` | Closed Historical | A consumed failed attempt or superseded row retained for audit; it is never retried and does not block its accepted successor. |

Write markers with backticks, exactly as in the table above. Each table row is
a commit unit: after changing its state, verify the change, update the memory
bank/docs, and make a scoped `git commit` before starting the next row. A
`[-]` row's notes must record the consumed attempt or supersession and
identify its accepted successor.

After the milestone's review, consolidation, and downstream reconciliation
pass, follow [the retirement procedure](milestone.md#long-term-memory-and-retirement).
Preserve this complete document and the full milestone specification in its
retired record.

Review provenance: "apitools Code Review", Passes 2-4, revalidated at
`e015231ef651bb92ee1aab9a4e8fc870d94ad020` with a clean worktree.

Iteration 2 review provenance: Review "apitools uncommitted-diff review" (2026-09-26; source priority not supplied), baseline `6ee935b` plus the uncommitted implementation diff; revalidated at `6ee935b328234287cc8dfa097884a67fdebccfbd` including those uncommitted changes. Finding IDs U1-U15 follow the review's order.

**Review.** passed
**Review iterations.** 3 (passed; iteration 1 passed as recorded in the restored retirement record)
**Review findings.** Iteration 2: the external uncommitted-diff review, counted as a confirming full pass at the user's direction, found no P1/P2 or higher-severity issue in M76's scope; no finding concerns the search/import exit-code handlers. Iteration 3: closing whole-milestone review re-traced `runSearchWithClient`/`runImportWithClient`/`validateImportURL`/`parseCommandFlags` against the acceptance criteria: query/URL/dir/scheme/hostname validation runs and returns `exitUsage` (2) with usage on stderr before any client is created, `-h`/`-help` returns `exitSuccess` (0), and only runtime failures (client creation, private-host policy, unknown source) return `exitRuntime` (1). No P1/P2 or higher-severity issue found. `go build`, `go vet`, and `go test ./...` pass across all 17 packages, and the documented CLI smoke checks (help exits 0; missing/short-query, missing/malformed-URL/scheme, and missing-dir usage exit 2; private-host runtime failure exits 1) pass against the built binary. Gate closes at iteration 3; retiring alongside this pass.

| Item | State | Notes |
|---|---|---|
| Return exit 2 for search/import usage errors | `[+]` | Pass 4 X1, source Medium, local P2. `search` now rejects missing/blank and too-short queries before client creation; `import` validates required URL/directory plus absolute HTTP(S) syntax/host before client/library work. Host safety, network, and content failures remain runtime exit 1. Tests cover stderr usage, no client creation for invalid arguments, and a private-host runtime error. `go test ./...`, `go vet ./...`, `git diff --check`, and executable smoke checks pass: help exits 0; missing/short/scheme/URL usage exits 2; private-host policy exits 1. C03 is accepted; this row owns distinct `search`/`import` handlers in the same CLI file ([C03 history](status-C03.md)). Lineage: M73. |

## Bounded Review-Fix Gate

- Iteration 1 of 10 started: 2026-09-26 00:28 UTC.
- Review scope: the M76 search/import validation and usage output, CLI tests, and the relevant README/milestone/status contract.
- Iteration 1 review completed: no P1/P2-or-higher findings. The bounded review-fix gate passes at iteration 1; full acceptance verification is recorded in the task row above.
````
