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

| Item | State | Notes |
|---|---|---|
| Strip or escape invisible and direction-changing Unicode | `[ ]` | Pass 3 P1, source High, local P1. `sanitizePromptString` (`prompt_safety.go`) removes only Cc controls and ANSI; probe shows U+202E, U+200B, and U+E0000-block tag characters survive `SanitizeOperationSummaries`. Remove/escape Cf, tag block, and variation selectors in every sanitized field and emit a diagnostic. Lineage: S01. |
| Regression tests and consumer verification | `[ ]` | Tests for each character class through operation and inventory summaries, preserving ordinary non-ASCII text; run `go test ./...`, `go vet ./...`, `git diff --check`, and `(cd ../openudon && go test ./...)`. |
