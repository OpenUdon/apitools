# Retired milestone S03 - Operation Lifecycle Ranking Correctness

**Milestone.** S03
**Outcome.** completed
**Retired.** 2026-09-26
**Source status.** tabilet/memory-bank/status-S03.md
**Source specification.** tabilet/memory-bank/milestone.md#s03---operation-lifecycle-ranking-correctness
**Evidence.** 3bcba9a6c74e3d6b38307dfb04b3578d789ad174
**Worktree.** clean
**Review.** passed
**Review iterations.** 4
**Verification.** go build ./..., go vet ./..., go test ./... (apitools, all packages); OpenUdon full go test ./...; Ramen full go test ./... including the long-running tests package (333s). All passed at the evidence commit.
**Consolidated into.** tabilet/memory-bank/architecture.md (operationlifecycle/ row updated) and tabilet/memory-bank/lessons.md ("Keep inferred POST actions out of resource creation", revised for the trailing-vs-parent-parameter distinction).

## Milestone specification

````markdown
## S03 - Operation Lifecycle Ranking Correctness

**Goal.** Make lifecycle sibling ranking pick the true item-level
read/update/delete siblings with honest confidence, so destructive or unrelated
operations are not proposed as lifecycle siblings.

**Scope.**

- Treat a path as an item sibling only when it equals the seed's collection
  path plus a trailing parameter segment; resolve Google Discovery
  `{+name}`/`{+parent}` paths from method resource identity; remove the
  `list`/`collection` token exclusions and the dead path-match clause.
- Derive roles from operation semantics by reusing
  `apitools.ClassifyOperationPurpose`, so POST updates and actions are not
  labelled create, and score or diagnose only non-seed roles.
- Build family tokens from operation IDs and paths without free-text
  stop-words, and match goal intent on word boundaries.
- Resolve a seed identified only by operation ID, and prefer absolute document
  path or URL over relative path for source identity.
- Keep the package metadata-only; do not add workflow semantics, fetches,
  credentials, or execution.

**Review provenance.** "apitools Code Review" Pass 0 findings F01, F02, F03,
F04, F05, F06, F07, F08, F12, and F14 (source priority not supplied), stated
baseline `a5699e5`, revalidated at
`a5699e570fb5d6288fdd255c71d229fb062b50b3` with a probe of `Expand`; F04 is
partially confirmed by code evidence only. F13 was unsupported (not
reproduced). Lineage: M75 (completed; not reopened).

**Dependencies.** None; independent of S02 and C02.

**Parallel ownership.** `operationlifecycle/` only, plus its tests and docs.

**Downstream impacts.** OpenUdon and Ramen consume lifecycle roles for draft
ranking; their workspace suites must pass. Publishing and re-pinning consumers
need separate authorization.

**Acceptance.** Regression tests cover scoped collection versus item siblings,
Discovery `{+name}` resources, POST update seeds, resources named
`list`/`collection`, stop-word-only family overlap, goals containing
`dispatch` or other embedded verbs, operation-ID-only seeds, and same relative
paths from different documents. Review-iteration-3 regressions (U2-U7, U11)
also hold: a seed identified by relative path, URL, or name resolves against
inventory operations; a fully specified seed with a duplicated operation ID
resolves by method and path; a create-or-update PUT keeps its update sibling
when the goal asks for updates; nested-collection POST creates stay create;
HEAD operations can be read siblings; an explicit create/update operation ID
outranks summary or tag wording; and inflected goal verbs (`patching`,
`patched`, `replaced`) request updates. `go test ./...`, `go vet ./...`,
`git diff --check`, and the OpenUdon and Ramen workspace suites pass.
````

## Status record

````markdown
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

**Review.** passed
**Review iterations.** 4 (passed; iterations 1-2 passed as recorded in the restored retirement record)
**Review findings.** Iteration 3: the external uncommitted-diff review, counted as the next full pass at the user's direction, found P2 regressions U2, U3, U4, U5, U6, U7, and U11 (fix rows below). Iteration 4: after all seven findings were fixed (commits `6ae8b0f`, `a878780`, `937bd72`), a fresh whole-milestone review of the full `operationlifecycle/operationlifecycle.go` implementation found no P1/P2 or higher-severity issue: source-identity comparison correctly falls back per-kind without reopening the same-relative-path-different-document collision (F12/`TestSameRelativePathDifferentAbsoluteDocumentsAreNotSiblings`), the create-or-update PUT/HEAD/inflected-verb fixes agree with `primaryRole` and `verbMatchesRole`, and `narrowSeedMatchesByMethodAndPath` only narrows when method+path are both present. `architecture.md`'s `operationlifecycle/` row updated to describe the corrected behavior. Full apitools `go test ./...`/`go vet ./...`, OpenUdon full `go test ./...`, and Ramen full `go test ./...` (including `tests`, 333s) all pass. The gate is closed within the persisted iteration count.

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
| Seed resolution and source identity | `[+]` | U2 and U11, local P2, fixed. Replaced the single-string `sourceID` with `sourceIdentity` (absolute path, URL, relative-path/name), compared pairwise: the highest-priority kind both sides have populated is decisive, so a disagreement there rejects the match even if a lower-priority kind agrees (preserves the same-relative-path-different-document rejection), while a kind one side lacks is skipped rather than forcing a mismatch (fixes U2's under-specified seed). `normalizeSeed` now also narrows same-ID, same-source candidates by the seed's own method+path when it carries them, resolving a duplicated operation ID to its exact match instead of reporting ambiguity (fixes U11); a bare operation-ID-only seed with no method/path is still reported ambiguous. Regressions: `TestUnderSpecifiedSeedResolvesAgainstFullyDescribedOperation`, `TestFullySpecifiedSeedResolvesDuplicateOperationIDByMethodAndPath` (covers both the resolved and the still-ambiguous case). All prior seed/source tests (`TestOperationIDOnlyAmbiguousSeedDoesNotChooseDocument`, `TestSameRelativePathDifferentAbsoluteDocumentsAreNotSiblings`, `TestAbsoluteDocumentPathAndURLOutrankRelativePath`) still pass. `go test ./operationlifecycle -v`, `go vet ./operationlifecycle`, full apitools `go test ./...`/`go vet ./...`, and `go build`+`go test ./internal/icot/...` in OpenUdon (the consumer package) pass. |
| Role classification regressions | `[+]` | U3, U4, and U6, local P2, fixed. **U3:** `Expand`'s sibling-skip decision and `primaryRole`'s final label now share one `seedPrimaryPurpose` helper, so a create-or-update PUT is treated as "create" in both places; the update role is no longer skipped in the sibling search merely because `ClassifyOperationPurpose`'s raw, method-based answer for any PUT is "update". Regression: `TestCreateOrUpdatePUTKeepsUpdateSiblingWhenGoalAsksForUpdate` (Databases_CreateOrUpdate/Get/Update, goal "update databases"). **U4:** `postOperationIsAction`'s path-parameter check now looks only at the trailing path segment, not every segment, so a parent-scoping parameter earlier in the path (nested-collection create) no longer suppresses the create role; a POST directly against a parameterized item is still treated as an action. Regression: `TestNestedCollectionPOSTStaysCreateDespiteParentPathParameter` (covers both the nested-create and the still-generic item-action case). **U6:** added `operationIDHasAny`/`operationIDTokens` (operation-ID-only tokenization, no summary/tags) and reordered `lifecyclePurpose` so an explicit create verb in the operation ID is checked before the update-name check, which still considers summary/tag text. Regression: `TestExplicitCreateOperationIDOutranksUpdateWordingInSummary`. Revisited the "Keep inferred POST actions out of resource creation" lesson to state the trailing-vs-parent-parameter distinction. All prior classification tests (`TestExpandClassifiesPostUpdatesAndActions`, `TestExpandHyphenAndCreateOrUpdateFamilies`, `TestExpandCollectionItemLifecycle`) still pass unchanged. `go test ./operationlifecycle -v`, `go vet ./operationlifecycle`, full apitools `go test ./...`/`go vet ./...`, and OpenUdon's `go build ./...` plus `go test ./internal/icot/...` pass. |
| HEAD read siblings and inflected goal verbs | `[+]` | U5 and U7, local P2, fixed. **U5:** `lifecyclePurpose` now maps HEAD to `"read"` directly (matching how `verbMatchesRole`/`methodRole` already treat it) instead of falling through to `ClassifyOperationPurpose`, which does not classify HEAD at all. Regression: `TestHeadOperationCanBeReadSibling`. **U7:** `goalWantsUpdate`'s word list now includes the missing inflected forms for every verb family (`updating`; `patched`, `patching`; `modified`; `replaced`) instead of just the three the review probed, since the underlying gap (word-boundary tokenization does not stem "-ing"/"-ed") applies uniformly across all four verbs. The existing `TestGoalWantsUpdateUsesWordBoundaries` asserted "no patching required" should *not* match update; that assertion was the exact gap this fix corrects, so it was rewritten to move "patching"/"patched"/"replaced" into the positive cases (matching the review's stated baseline: substring matching at commit `6ee935b` returned true for these before word-boundary tokenization dropped them). `dispatch a thing` and `the updater flow` remain negative (real word-boundary false positives, unrelated to inflection). `go test ./operationlifecycle -v`, `go vet ./operationlifecycle`, full apitools `go test ./...`/`go vet ./...` pass; OpenUdon full `go test ./...` passes; Ramen full `go test ./...` passes, including the long-running `tests` package (333s). |
````
