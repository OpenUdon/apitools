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

Iteration 3 review provenance: Review "apitools uncommitted-diff review" (2026-09-26; source priority not supplied), baseline `6ee935b` plus the uncommitted implementation diff; revalidated at `6ee935b328234287cc8dfa097884a67fdebccfbd` including those uncommitted changes. Finding IDs U1-U15 follow the review's order.

**Review.** open
**Review iterations.** 3 (open; iterations 1-2 passed as recorded)
**Review findings.** Iteration 3: the external uncommitted-diff review, counted as the next full pass at the user's direction, found P2 U12 and Lower U14 and U15 (fix rows below).
Earlier iterations: Iteration 1 found a P2: a cache-registration failure after an artifact overwrite could leave the old SQLite manifest inconsistent, and per-result registration could partially commit a batch. Fixed by prevalidating refresh registrations and committing their spec/artifact rows in one SQLite transaction; the CLI snapshots existing registered files and restores them if registration fails. Added a rollback regression. Iteration 2 whole-milestone review found no P1/P2 or higher-severity issue. Full verification passed: `go test ./...`, `go vet ./...`, standalone `GOWORK=off go test ./...`, catalog generator check, catalog quality check, and `git diff --check`.

| Item | State | Notes |
|---|---|---|
| Keep refresh artifacts and manifest consistent on partial failure | `[+]` | Pass 1 #3, source Medium, local P2. `RefreshCatalogSpecReferences` retains completed rows alongside a later provider/spec error; the CLI registers those rows before surfacing the failure. SQLite prevalidates and atomically commits each report's spec and artifact rows; the CLI snapshots and restores any overwritten registered paths if registration fails. Regressions cover a two-reference HTTP failure, partial registration/materialization, and rejected-registration rollback. Focused and full tests pass. S02 prerequisite passed: `artifactio.WriteFile` resolves selected-root ancestors once and retains relative-path confinement; C02 keeps source-aligned naming. Lineage: M11, M16. |
| Refresh reporting, docs, and catalog verification | `[+]` | Text/JSON CLI failure output retains completed results, names the failing provider/spec, and exits nonzero. README and architecture document registration-before-failure behavior and manifest consistency, including rollback when registration fails. Verified `go test ./...`, `go vet ./...`, standalone `GOWORK=off go test ./...`, `go run ./cmd/cataloggen -check`, `go run ./cmd/apitools catalog check`, and `git diff --check`. |
| Refresh protection without whole-file snapshots | `[+]` | U12, local P2, fixed. Added `artifactio.MoveAside`/`RestoreAside`/`DiscardAside`/`RemoveArtifact` (rename-based, no content read, confined and symlink-safe like the rest of the package). `snapshotCatalogRefreshArtifacts` now moves each registered artifact aside instead of reading it, so a large registered artifact is bounded only by the refresh's own `--max-bytes`, and a row whose move-aside fails (unreadable) is simply left out of the backup map, so it does not block refreshing itself or any other spec. New `finalizeCatalogRefreshArtifacts` (replacing `restoreCatalogRefreshArtifacts`) commits (discards the backup) exactly the rows `report.Results` covers when `registrationErr == nil` — so an earlier row already registered is kept even when `refreshErr` reports a later reference failed — and otherwise restores a backed-up row to its exact prior path or removes a brand-new, never-registered artifact so no orphan is left behind. Regressions: `TestCatalogRefreshRemovesNewArtifactWhenRegistrationFailsForUnregisteredSpec` (new), `TestMoveAsideRestoreDiscardAndRemoveArtifact`, `TestMoveAsideRejectsSymlinkAndConfinement` (new, in `internal/artifactio`); the existing two-reference partial-success/rejected-registration tests needed no change and now correctly commit the earlier row instead of rolling it back. `go test ./...`, `go vet ./...`, `go run ./cmd/cataloggen -check`, `go run ./cmd/apitools catalog check` (0/0) pass. |
| Refresh error prefix and shared registration SQL | `[ ]` | U14 and U15, Lower, carried here because the rows above rewrite the same code. Errors print the provider/spec prefix twice (`catalog_refresh.go`); `insertRefreshSpec` and `insertRefreshCatalogArtifact` (`sqlitecache/catalog_refresh.go`) copy the `StoreSpec`/`StoreCatalogArtifact` upserts. Prefix once; extract shared transaction-level store helpers. |
