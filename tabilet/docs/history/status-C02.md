# Retired milestone C02 - Catalog Refresh Manifest Integrity

**Milestone.** C02
**Outcome.** completed
**Retired.** 2026-09-26
**Source status.** tabilet/memory-bank/status-C02.md
**Source specification.** tabilet/memory-bank/milestone.md#c02---catalog-refresh-manifest-integrity
**Evidence.** 02e473fb6f3bf765364699484bbfac31b90ddffc
**Worktree.** clean
**Review.** passed
**Review iterations.** 4
**Verification.** go build ./..., go vet ./..., go test ./... (apitools, all packages); go run ./cmd/cataloggen -check; go run ./cmd/apitools catalog check (0 errors/warnings); git diff --check; OpenUdon full go build ./.... All passed at the evidence commit.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## C02 - Catalog Refresh Manifest Integrity

**Goal.** Keep catalog refresh artifacts on disk consistent with the
`cache.sqlite` integrity manifest when a batch fails partway.

**Scope.**

- Stage refreshed artifacts and promote them only after the whole batch
  validates and registers, or register completed rows before reporting a
  failure; never leave an overwritten registered artifact whose SHA-256 and
  byte count disagree with the cache database.
- Surface partial results and the failing reference in the CLI report.
- Protect registered artifacts without whole-file in-memory snapshots: a
  refresh must not be blocked by, or re-read, a large or unreadable registered
  artifact, and an artifact written for an unregistered spec is also restored
  when registration fails.
- Do not change download safety, validation rules, or artifact path policy.

**Review provenance.** "apitools Code Review" Pass 1 finding #3, stated
baseline `a5699e5`, revalidated at
`a5699e570fb5d6288fdd255c71d229fb062b50b3` by code trace of
`catalog_refresh.go` and the `catalog refresh` CLI flow. Lineage: M11, M16
(completed; not reopened).

**Dependencies.** S02 row "Anchor symlink checks at caller-chosen roots",
because both touch `internal/artifactio` root handling.

**Parallel ownership.** `catalog_refresh.go`, `sqlitecache/catalog_refresh.go`,
the `catalog refresh` command in `cmd/apitools`, and their tests/docs.

**Downstream impacts.** `catalog materialize` and `catalog export` keep passing
integrity checks after a failed refresh. No consumer API change.

**Acceptance.** A two-reference refresh whose second reference fails leaves the
first artifact either unchanged or registered with its new digest, and
materialization of that provider succeeds. A registered artifact larger than
128 MiB refreshes under a matching `--max-bytes`, an unreadable registered
artifact does not block its own replacement, failure errors name the
provider/spec once, and refresh registration shares the spec/artifact upsert
code with `StoreSpec`/`StoreCatalogArtifact`. `go test ./...`,
`go vet ./...`, `go run ./cmd/cataloggen -check`,
`go run ./cmd/apitools catalog check`, and `git diff --check` pass.
````

## Status record

````markdown
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

**Review.** passed
**Review iterations.** 4 (passed; iterations 1-2 passed as recorded)
**Review findings.** Iteration 3: the external uncommitted-diff review, counted as the next full pass at the user's direction, found P2 U12 and Lower U14 and U15 (fix rows below). Iteration 4: after fixing U12 (commit `8db527a`) and U14/U15 (commit `0f004d5`), a fresh whole-milestone review of the full refresh CLI flow (`runCatalogRefreshWithClient`, `snapshotCatalogRefreshArtifacts`, `finalizeCatalogRefreshArtifacts`), `internal/artifactio`'s new move-aside helpers, `catalog_refresh.go`'s error paths, and the shared SQLite registration helpers found no P1/P2 or higher-severity issue. Confirmed the per-row commit decision correctly distinguishes "committed rows registration already covered" from "rows to roll back" across all four combinations of refresh/registration success and failure, including the zero-results case. `go test ./...`, `go vet ./...`, `go run ./cmd/cataloggen -check`, `go run ./cmd/apitools catalog check` (0/0), `git diff --check`, and OpenUdon's full `go build ./...` all pass. The gate is closed within the persisted iteration count.
Earlier iterations: Iteration 1 found a P2: a cache-registration failure after an artifact overwrite could leave the old SQLite manifest inconsistent, and per-result registration could partially commit a batch. Fixed by prevalidating refresh registrations and committing their spec/artifact rows in one SQLite transaction; the CLI snapshots existing registered files and restores them if registration fails. Added a rollback regression. Iteration 2 whole-milestone review found no P1/P2 or higher-severity issue. Full verification passed: `go test ./...`, `go vet ./...`, standalone `GOWORK=off go test ./...`, catalog generator check, catalog quality check, and `git diff --check`.

| Item | State | Notes |
|---|---|---|
| Keep refresh artifacts and manifest consistent on partial failure | `[+]` | Pass 1 #3, source Medium, local P2. `RefreshCatalogSpecReferences` retains completed rows alongside a later provider/spec error; the CLI registers those rows before surfacing the failure. SQLite prevalidates and atomically commits each report's spec and artifact rows; the CLI snapshots and restores any overwritten registered paths if registration fails. Regressions cover a two-reference HTTP failure, partial registration/materialization, and rejected-registration rollback. Focused and full tests pass. S02 prerequisite passed: `artifactio.WriteFile` resolves selected-root ancestors once and retains relative-path confinement; C02 keeps source-aligned naming. Lineage: M11, M16. |
| Refresh reporting, docs, and catalog verification | `[+]` | Text/JSON CLI failure output retains completed results, names the failing provider/spec, and exits nonzero. README and architecture document registration-before-failure behavior and manifest consistency, including rollback when registration fails. Verified `go test ./...`, `go vet ./...`, standalone `GOWORK=off go test ./...`, `go run ./cmd/cataloggen -check`, `go run ./cmd/apitools catalog check`, and `git diff --check`. |
| Refresh protection without whole-file snapshots | `[+]` | U12, local P2, fixed. Added `artifactio.MoveAside`/`RestoreAside`/`DiscardAside`/`RemoveArtifact` (rename-based, no content read, confined and symlink-safe like the rest of the package). `snapshotCatalogRefreshArtifacts` now moves each registered artifact aside instead of reading it, so a large registered artifact is bounded only by the refresh's own `--max-bytes`, and a row whose move-aside fails (unreadable) is simply left out of the backup map, so it does not block refreshing itself or any other spec. New `finalizeCatalogRefreshArtifacts` (replacing `restoreCatalogRefreshArtifacts`) commits (discards the backup) exactly the rows `report.Results` covers when `registrationErr == nil` — so an earlier row already registered is kept even when `refreshErr` reports a later reference failed — and otherwise restores a backed-up row to its exact prior path or removes a brand-new, never-registered artifact so no orphan is left behind. Regressions: `TestCatalogRefreshRemovesNewArtifactWhenRegistrationFailsForUnregisteredSpec` (new), `TestMoveAsideRestoreDiscardAndRemoveArtifact`, `TestMoveAsideRejectsSymlinkAndConfinement` (new, in `internal/artifactio`); the existing two-reference partial-success/rejected-registration tests needed no change and now correctly commit the earlier row instead of rolling it back. `go test ./...`, `go vet ./...`, `go run ./cmd/cataloggen -check`, `go run ./cmd/apitools catalog check` (0/0) pass. |
| Refresh error prefix and shared registration SQL | `[+]` | U14 and U15, Lower, fixed. **U14:** `refreshCatalogSpecReference` and `validateCatalogRefreshContentRaw` (`catalog_refresh.go`) no longer prefix their own errors with `provider/spec`; the outer `RefreshCatalogSpecReferences` loop is now the single place that adds it, so a failure no longer prints `slack/slack-web: slack/slack-web: ...`. `catalog_refresh_review.go` benefits too, since its report already carries `ProviderID`/`SpecRefID` as separate fields. Regression: `TestRefreshCatalogSpecReferencesRejectsUnparseableOpenAPI` now asserts the prefix appears exactly once. **U15:** extracted `storeSpecDocumentTx`/`storeCatalogArtifactTx` (`sqlitecache/cache.go`) so `StoreSpec`/`StoreCatalogArtifact` and `insertRefreshSpec`/`insertRefreshCatalogArtifact` (`sqlitecache/catalog_refresh.go`) share one copy of the `spec_documents`/`catalog_artifacts` upsert SQL instead of two. `go test ./...`, `go vet ./...`, `go run ./cmd/cataloggen -check`, `go run ./cmd/apitools catalog check` (0/0) pass. This closes C02's last pending fix row. |
````
