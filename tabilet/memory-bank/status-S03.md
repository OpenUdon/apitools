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

Iteration 3 review provenance: Review "apitools uncommitted-diff review" (2026-09-26; source priority not supplied), baseline `6ee935b` plus the uncommitted implementation diff; revalidated at `6ee935b328234287cc8dfa097884a67fdebccfbd` including those uncommitted changes. Finding IDs U1-U15 follow the review's order.

**Review.** open
**Review iterations.** 3 (open; iterations 1-2 passed as recorded in the restored retirement record)
**Review findings.** Iteration 3: the external uncommitted-diff review, counted as the next full pass at the user's direction, found P2 regressions U2, U3, U4, U5, U6, U7, and U11 (fix rows below).

| Item | State | Notes |
|---|---|---|
| Item-path matching and Discovery resource identity | `[+]` | Pass 0 F01, F02, F05 (local P2) and F14 (Lower); source priority not supplied. Candidate read/update/delete roles now require a trailing item parameter and exact collection scope; nested/missing-scope and collection-level paths no longer match. Google Discovery `{+name}` paths match only when explicit source-kind provenance and the operation-ID resource prefix agree; `{+parent}` is not treated as an item ID. Removed list/collection token exclusions. Regressions cover scoped paths, collection reads, same-route wrong-resource Discovery IDs, and resources named `list`/`collection`. `go test ./operationlifecycle -count=1` and `git diff --check` pass. Lineage: M75. |
| Role classification from operation semantics | `[+]` | Pass 0 F03 (local P2) and F08 (Lower). Lifecycle purposes reuse `apitools.ClassifyOperationPurpose` and only score candidates classified as non-seed roles. POST updates remain `update`; common/colon actions remain `post`; unlisted POSTs on parameterized item/action routes now stay `post` absent explicit create semantics. `createOrUpdate` and explicit `create`/`insert`/`add` retain create. Regressions cover `updatePetWithForm`, `closeTicket`, `renewInvoice`, and explicit child creation. `go test ./operationlifecycle -count=1`, `go vet ./operationlifecycle`, and `git diff --check` pass. Lineage: M75. |
| Family tokens and goal matching | `[+]` | Pass 0 F04 (partially confirmed by code, local P2) and F07 (local P2). Family keys now tokenize only operation IDs and paths, excluding generic lifecycle/API/resource/ID words; summaries and tags cannot create family overlap. Update intent is token-matched, so "dispatch" does not trigger "patch". Regressions cover free-text exclusion, stop-word-only overlap, embedded verbs, and positive update terms. `go test ./operationlifecycle -count=1` and `git diff --check` pass. Lineage: M75. |
| Seed and source identity | `[+]` | Pass 0 F06 (Lower) and F12 (local P2). A unique operation-ID-only seed now resolves to the matching full summary; duplicate IDs produce an explicit ambiguity diagnostic instead of choosing a source. Source identity prefers absolute document paths, then URLs, before relative paths/names, so identical relative names from different roots are not paired. Regressions cover unique/ambiguous ID-only seeds, cross-root relative-path collisions, and positive absolute-path/URL matches. `go test ./operationlifecycle -count=1` and `git diff --check` pass. F13 remains unsupported. Lineage: M75. |
| Downstream verification | `[+]` | Full apitools `go test ./...`, `go vet ./...`, and `git diff --check` pass; OpenUdon `go test ./...` passes after the separately approved test-only registration fixture date correction; Ramen `go test ./...` passes (including the long-running `tests` package). All required suites were rerun after the review fix. No publishing or re-pinning performed. |

## Bounded Review-Fix Gate

- Iteration 1 of 10 started: 2026-09-25 22:27 UTC.
- Review scope: the complete S03 `operationlifecycle` implementation and tests, architecture contract, and this milestone's changes.
- Iteration 1 outcome: P2 — the finite POST action-name list classified an unlisted action on a parameterized item/action path (for example, `renewInvoice` on `/invoices/{id}`) as `create`. Fixed by preserving explicit create/update semantics and keeping other POSTs on parameterized resource paths as `post`; regression tests now cover unlisted actions and explicit child creation. Focused package tests, package vet, and diff check pass.
- Iteration 2 of 10 started: 2026-09-25 22:37 UTC, after full apitools, OpenUdon, and Ramen verification passed with the iteration-1 fix.
- Iteration 2 outcome: pass. Re-reviewed the full S03 implementation, regression tests, operation lifecycle contract, and milestone changes; no P1, P2, or higher-severity findings remain. `go test ./operationlifecycle -count=1` and `git diff --check` pass. Review-fix gate passed at iteration 2 of 10.
| Seed resolution and source identity | `[ ]` | U2 and U11, local P2. `sourceID` ranks absolute `DocumentPath` first and `normalizeSeed` matches only on operation ID: a seed with only a relative path returns `seed_source_mismatch` and no roles, and a fully specified seed whose operation ID is duplicated returns `seed_ambiguous` and no roles (probes). Match seeds by method and path when present and treat sources as the same when any shared identifier agrees. |
| Role classification regressions | `[ ]` | U3, U4, and U6, local P2. A create-or-update PUT keeps no update sibling even for an update goal (probe: Databases_*); `postOperationIsAction` treats any path parameter as an action, so nested-collection creates become `post` (probe: newChild); summary/tag wording overrides an explicit `createWidget` operation ID (probe). Skip the update role only when the seed's primary role is update; detect actions only on trailing item parameters; let explicit operation-ID verbs outrank summary/tag text. Revisit the "Keep inferred POST actions out of resource creation" lesson accordingly. |
| HEAD read siblings and inflected goal verbs | `[ ]` | U5 and U7, local P2. `bestSibling` requires `lifecyclePurpose(op) == role`, but `ClassifyOperationPurpose` returns empty for HEAD, so HEAD reads are excluded (probe); whole-token goal matching misses `patching`, `patched`, and `replaced` (probe). Map HEAD to read and match inflected verb forms. |
