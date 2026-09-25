# Status C02 - Catalog Refresh Manifest Integrity

State of each C02 milestone item. Update as items complete. See
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

Review provenance: "apitools Code Review", stated baseline `a5699e5`,
revalidated at `a5699e570fb5d6288fdd255c71d229fb062b50b3`; relevant uncommitted
changes at revalidation were harness documentation only.

| Item | State | Notes |
|---|---|---|
| Keep refresh artifacts and manifest consistent on partial failure | `[ ]` | Pass 1 #3, source Medium, local P2. `RefreshCatalogSpecReferences` writes with `Force: true` to registered paths, then returns an empty report on a later error; the CLI skips `RegisterCatalogRefreshResults`, so materialize/export fail SHA checks. Stage-then-promote or register completed rows; add a two-reference test with one failure. Depends on S02 root anchoring. Lineage: M11, M16. |
| Refresh reporting, docs, and catalog verification | `[ ]` | Surface partial results and the failing reference in `catalog refresh` output; run `go test ./...`, `go vet ./...`, `go run ./cmd/cataloggen -check`, `go run ./cmd/apitools catalog check`, and `git diff --check`. |
