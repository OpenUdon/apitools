# Status S03 - Operation Lifecycle Ranking Correctness

State of each S03 milestone item. Update as items complete. See
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
| Item-path matching and Discovery resource identity | `[ ]` | Pass 0 F01, F02, F05 (local P2) and F14 (Lower); source priority not supplied. Probe: collection `GET getInstances` chosen over item `fetchInstance`; Discovery `instances.delete` paired with `operations.get`; Qdrant `*_collection` update/delete lost. `lifecyclePathsMatch` accepts any parameterized path. Lineage: M75. |
| Role classification from operation semantics | `[ ]` | Pass 0 F03 (local P2) and F08 (Lower). Probe: `create:updatePetWithForm(high)`. Reuse `apitools.ClassifyOperationPurpose`; score and diagnose only non-seed roles. Lineage: M75. |
| Family tokens and goal matching | `[ ]` | Pass 0 F04 (partially confirmed by code, local P2) and F07 (local P2). `operationTokens` uses free-text summary/tag words; probe: goal "dispatch a thing" adds `update:updateThing(high)`. Use IDs/paths without stop-words; word-boundary goal matching. Lineage: M75. |
| Seed and source identity | `[ ]` | Pass 0 F06 (Lower) and F12 (local P2). Probe: operation-ID-only seed returns only itself; `/a` and `/b` documents with the same relative path are paired. Resolve by operation ID; prefer absolute path/URL identity. F13 not reproduced (unsupported). Lineage: M75. |
| Downstream verification | `[ ]` | Run `go test ./...`, `go vet ./...`, `git diff --check`, and the OpenUdon and Ramen workspace suites. Publishing/re-pinning needs separate authorization. |
