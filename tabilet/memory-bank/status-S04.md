# Status S04 - Prompt Sanitizer Invisible-Unicode Hardening

State of each S04 milestone item. Update as items complete. See
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
**Review findings.** Iteration 2: the external uncommitted-diff review, counted as a confirming full pass at the user's direction, found no P1/P2 or higher-severity issue in S04's scope; the only related note is a Lower cosmetic item routed to Candidate Directions. Closure-ready: retire after final verification together with the shared commit plan.

| Item | State | Notes |
|---|---|---|
| Strip or escape invisible and direction-changing Unicode | `[+]` | Pass 3 P1, source High, local P1. `sanitizePromptString` now replaces Unicode format controls, tag characters, and variation selectors with spaces. Operation and inventory summary tests confirm U+202E, U+200B, tag characters, and both variation-selector ranges are removed, natural non-ASCII text is preserved, and sanitization emits diagnostics. Lineage: S01. |
| Regression tests and consumer verification | `[+]` | Focused Unicode operation/inventory tests, apitools `go test ./...`, `go vet ./...`, and `git diff --check` pass. With user approval, OpenUdon's test-only registration fixture now anchors its synthetic observation at current UTC minus two minutes; its focused test, full `go test ./...`, and `git diff --check` pass. No OpenUdon product code changed. This goal run has `COMMIT_POLICY: none`. |

## Bounded Review Gate

Iteration 1 passed after review of the complete S04 diff, sanitizer call paths,
diagnostics, test behavior, downstream fixture correction, and acceptance.
Resolved the original P1; no new P1, P2, or higher-severity findings. No lower
findings carried.

## Reconciliation

Updated `architecture.md` to record the prompt sanitizer's Unicode guarantees.
Product scope and tools did not change; no lesson addition or evolution bump is
warranted. OpenUdon's required consumer suite passes after its approved test-only
fixture correction; no additional pending downstream work was found.
