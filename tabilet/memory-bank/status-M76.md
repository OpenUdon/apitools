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
**Review iterations.** 2 (passed; iteration 1 passed as recorded in the restored retirement record)
**Review findings.** Iteration 2: the external uncommitted-diff review, counted as a confirming full pass at the user's direction, found no P1/P2 or higher-severity issue in M76's scope; no finding concerns the search/import exit-code handlers. Closure-ready: retire after final verification together with the shared commit plan.

| Item | State | Notes |
|---|---|---|
| Return exit 2 for search/import usage errors | `[+]` | Pass 4 X1, source Medium, local P2. `search` now rejects missing/blank and too-short queries before client creation; `import` validates required URL/directory plus absolute HTTP(S) syntax/host before client/library work. Host safety, network, and content failures remain runtime exit 1. Tests cover stderr usage, no client creation for invalid arguments, and a private-host runtime error. `go test ./...`, `go vet ./...`, `git diff --check`, and executable smoke checks pass: help exits 0; missing/short/scheme/URL usage exits 2; private-host policy exits 1. C03 is accepted; this row owns distinct `search`/`import` handlers in the same CLI file ([C03 history](status-C03.md)). Lineage: M73. |

## Bounded Review-Fix Gate

- Iteration 1 of 10 started: 2026-09-26 00:28 UTC.
- Review scope: the M76 search/import validation and usage output, CLI tests, and the relevant README/milestone/status contract.
- Iteration 1 review completed: no P1/P2-or-higher findings. The bounded review-fix gate passes at iteration 1; full acceptance verification is recorded in the task row above.
