# Retired milestone S06 - Public-apis source repair

**Milestone.** S06
**Outcome.** completed
**Retired.** 2026-10-06
**Source status.** tabilet/memory-bank/status-S06.md
**Source specification.** tabilet/memory-bank/milestone.md#s06--public-apis-source-repair
**Evidence.** c379ffb3d6074889e5978adbc733a300adfd558e
**Worktree.** includes uncommitted changes
**Review.** passed
**Review iterations.** 1
**Verification.** APItools full test/vet in workspace and standalone modes; focused root/CLI race checks; catalog generator/quality checks; search help; patch check; parser fuzz (13,412 executions); live default-README no-match smoke; OpenUdon/Udon candidate-bound full test/vet in workspace and standalone modes. Udon retries with reduced build concurrency passed after disk-quota and TempDir cleanup failures, retained below. Go 1.26.6; Ramen excluded. The evidence commit contains S06.1/S06.2; S06.3 documentation, closure and retirement are uncommitted at this recorded observation.
**Consolidated into.** README.md, tabilet/memory-bank/architecture.md and tabilet/memory-bank/tech-stack.md; no product-scope or lesson change.

## Milestone specification

````markdown
## S06 — Public-apis source repair

Restore the documented `public-apis` search source and the `auto` fallback
that ends in it. The default source host `api.publicapis.org` has no DNS
record (observed 2026-10-02): `apitools search --source public-apis` fails
with "no such host" and exits 1, and `--source auto` with a query that no
other source matches exits 1 with the same error instead of returning an empty
result. Both are documented behaviors (`README.md`, the CLI help), so this is
a confirmed P2 defect in a supported scenario. It was found while evaluating
the public-apis repository as a lead source, in this conversation's findings
of 2026-10-02 (source priority `not supplied`), and revalidated at
`8a52c3f602988b945a4b5c1960bce8c04170c63d` with no relevant uncommitted code.
Lineage: the source was added by retired M70
([M70 record](../docs/history/status-M70.md)); it is not reopened. The user
approved this milestone, its order before S05, and its rows on 2026-10-02.
The approved "S06/S05 plan consistency review, 2026-10-06" amends its pending
acceptance only; source priority and separate review baseline are `not supplied`.
Revalidation remains at the full commit above, including the six uncommitted
planning files and no code changes. Status S06 owns this intake's F08-F10.
The dead-endpoint finding remains the original S06 work, not a new milestone;
its live evidence is dated 2026-10-02 and was not fetched again.

**Approach.** The public-apis project is an MIT-licensed, actively maintained
list: one Markdown `README.md` of about 88 KB with roughly 1,970 entries in 52
categories, each with a name, description, docs link, and auth type. It has no
spec URLs. The default source becomes that README, read at runtime from
GitHub's raw host through the existing guarded transport. Nothing from the
list is vendored into this repository.

**Scope and acceptance.**

- **List parsing.** Parse the Markdown table (`###` category headings and
  `| [Name](link) | description | auth | https | cors |` rows) into the
  existing entry fields. Define and test explicit response-size, entry-count,
  and row-length bounds, with diagnostics naming the list URL for malformed
  or limit-exceeded input; never treat a truncated list as complete. Reject
  unrecognized or empty Markdown, but preserve a valid legacy JSON
  `{"entries": []}` as an empty result. Keep accepting the JSON
  `{"entries": [...]}` shape from a configured `PublicAPIsURL`, selecting the
  parser by content rather than by URL, so existing fixtures and mirrors work.
- **Default and wiring.** `DefaultPublicAPIsURL` points at the repository's
  README on GitHub raw. `Client`, `SearchOptions`, `Result`, and
  `SearchReport` shapes, the `PublicAPIsURL` field (which Udon passes
  through), the probe paths, the 5 s per-candidate timeout, the 30 s budget,
  the 50-candidate cap, and `auto`'s fallback order and error semantics are
  unchanged. When the final fallback cannot read its list, `auto` still
  reports that error, which is correct because nothing was checked.
- **Tests and documentation.** Synthetic Markdown and JSON fixtures served by
  local servers; a test that the default is not the dead host; no default test
  touches the network. README and architecture document the list source and
  its MIT attribution. A manual, opt-in check against live GitHub raw confirms
  `search --source public-apis` works, and is not part of default
  verification.
- **Verification.** `go test ./...`, `GOWORK=off go test ./...`, `go vet ./...`,
  `GOWORK=off go vet ./...`, `go run ./cmd/cataloggen -check`,
  `go run ./cmd/apitools catalog check`, `go run ./cmd/apitools search --help`,
  `git diff --check`, and OpenUdon/Udon compatibility checks in workspace and
  standalone modes per tech-stack.md, with Ramen excluded.

**Dependencies and execution.** No upstream prerequisite and no technical
dependency on S05; the order is by priority. No consumer adoption is required:
OpenUdon does not use the source, and Udon only plumbs the `PublicAPIsURL`
field and uses `Import`, so that field must stay. Udon builds against the local
`../apitools` checkout through a `replace` directive. Implement S06 then S05
on one branch in one separate worktree outside `~/Workspace/go.work`, with
one ledger owner; preserve the main workspace checkout and consumer pins.
Kinet's current ledger records Stage 5 complete (rechecked 2026-10-06), so it
is not a pending dependency or merge gate. Isolation still protects shared
workspace builds and accepted evidence. Merging and publishing each need
separate approval. Probe efficiency (parallel probing, conditional requests)
is not in scope; S05's new probe primitive owns that work.

**Isolated consumer verification (S06.3; reused by S05.6).** Use disposable
workspace/modfile configuration outside tracked repositories to bind OpenUdon
and Udon explicitly to the candidate APItools worktree. Resolve relative paths
without editing their tracked manifests or the shared `go.work`; verify the
resolved APItools module directory before testing. Run candidate-bound tests
in both workspace and standalone modes. Record ordinary standalone checks
against unchanged pins/replacements separately as baseline evidence, never as
proof of candidate compatibility. Ramen is excluded. No sibling write or
consumer adoption is authorized.
````

## Status record

````markdown
# Status S06 — Public-apis source repair

**State:** Complete, 2026-10-06. All three tasks and whole-milestone review
passed in the isolated `work/s06-s05` branch. No merge or publication occurred.

**Specification:** [S06](milestone.md#s06--public-apis-source-repair).

**Provenance:** Findings from this conversation (2026-10-02), labeled F01-F12
in S05's status file; this milestone owns F01. Source priority is `not
supplied`; the local severity is P2 under milestone.md's review severity
definitions, because a documented, supported source returns an error. No
stated baseline beyond the revalidation at the full commit
`8a52c3f602988b945a4b5c1960bce8c04170c63d`, with S05's uncommitted planning
files as the only worktree changes and no uncommitted code. The user approved
the finding, this milestone, its order before S05, and the planning-file
actions on 2026-10-02. Execution, merge, and release authority remain
separate.

**Provenance (consistency intake):** "S06/S05 plan consistency review,
2026-10-06", approved by the user. Source priorities and separate review
baseline are `not supplied`; revalidation is
`8a52c3f602988b945a4b5c1960bce8c04170c63d`, including the six pre-existing
uncommitted planning files and no uncommitted code. This intake owns F08-F10
below; those IDs are distinct from the 2026-10-02 F01 dead-endpoint finding.
The findings amend existing pending acceptance, not completed M70 history.
All live endpoint/list observations below remain dated 2026-10-02; this intake
made no remote fetch and does not claim a fresh DNS or publisher measurement.

**Evidence at planning (verified 2026-10-02):**

- `client.go:15` sets `DefaultPublicAPIsURL` to
  `https://api.publicapis.org/entries`. The host has no DNS record (the
  resolver reports no such host) and the apex does not answer.
- Running the built CLI: `apitools search --query cats --source public-apis`
  exits 1 with "lookup api.publicapis.org ... no such host", and
  `apitools search --query <no match> --source auto` exits 1 with
  "public-apis: lookup api.publicapis.org ... no such host", because `auto`
  falls through to this source last (`client.go`) and returns its error.
  `README.md` and the CLI help list the source.
- `providers.go` (`searchPublicAPIs`) expects a JSON `{"entries": [...]}`
  body and then probes seven well-known paths per linked host, one at a time,
  at most 50 candidates, 5 s each, 30 s total.
- The public-apis repository is MIT-licensed and active (commits on the day of
  inspection). Its `README.md` is one Markdown file of about 88 KB (0.13 s)
  with 1,972 entries in 52 categories, each a name, description, docs link, and
  auth type. It has no spec URLs and no JSON file; the hosted JSON API was a
  separate service that is gone.
- No sibling uses the source: OpenUdon does not, and Udon only plumbs the
  `PublicAPIsURL` field into a client it uses for `Import`, so that field and
  the JSON shape must stay compatible. Udon builds against `../apitools`
  through a `replace` directive.

**Lineage:** The source was added by retired
[M70](../docs/history/status-M70.md); S06 does not reopen it.

## Consistency intake findings (2026-10-06)

| Finding | Source priority | Local severity | Disposition | Evidence | Owner/action |
|---|---|---|---|---|---|
| F08 Isolated implementation may be checked against unchanged APItools | not supplied | P2 | confirmed planning gap | Shared `go.work`; Udon's `go.mod` replaces APItools with `../apitools`, while OpenUdon standalone uses its published pin | S06.3 owns the candidate-binding verification procedure; S05.6 reuses it |
| F09 Stage 5 protection rationale is stale | not supplied | Lower | confirmed; necessary planning correction | Kinet's current `tabilet/memory-bank/milestone.md` records Stage 5 complete; rechecked 2026-10-06 | S06.3 dependency/verification notes, reflected in S05; preserve isolation, remove the stale pending gate |
| F10 Rejecting all zero-entry responses changes legacy JSON behavior | not supplied | P2 | confirmed planning gap | `providers.go:searchPublicAPIs` accepts `{"entries":[]}`; the earlier S06.1 acceptance rejected every zero-entry body | S06.1 preserves valid empty JSON while rejecting unrecognized/empty Markdown |

F01-F07 of this intake belong to [S05](status-S05.md#consistency-intake-findings-2026-10-06).
No duplicate task or milestone is allocated.

| Item | State | Notes |
|---|---|---|
| S06.1 — Markdown list parser and bounds | `[+]` | Owns the 2026-10-02 F01 parsing and consistency-intake F10. Parse the public-apis README table (`###` headings and `\| [Name](link) \| description \| auth \| https \| cors \|` rows) into existing entry fields. Define and test explicit response-size, entry-count, and row-length bounds; malformed/limit-exceeded diagnostics name the list URL and never present a partial list as complete. Reject unrecognized/empty Markdown, but preserve valid legacy JSON, including `{"entries":[]}` as an empty result. Select the parser by content, not URL. Use synthetic Markdown/JSON fixtures only; commit no copy of the list. |
| S06.2 — Default source switch and wiring | `[+]` | Owns F01's default. Depends on S06.1. Point `DefaultPublicAPIsURL` at the repository's README on GitHub raw, read through the existing guarded transport. Leave `Client`, `SearchOptions`, `Result`, and `SearchReport` shapes, the `PublicAPIsURL` field, the probe paths, the 5 s and 30 s budgets, the 50-candidate cap, and `auto`'s fallback order and error semantics unchanged. Add tests that the default is not the dead host, that `auto` and `--source public-apis` work against local fixtures, and that no default test touches the network. |
| S06.3 — Documentation, compatibility, and verification | `[+]` | Depends on S06.1-S06.2. Owns consistency-intake F08-F09. Document the list source and MIT attribution in README and architecture. Run the verification below against the candidate worktree, recording ordinary unchanged-pin/replacement checks separately. Keep the live GitHub-raw smoke check separately opt-in, not a default test. Preserve worktree isolation without treating completed Kinet Stage 5 as a pending gate. Pass the persisted ten-iteration review gate with no open P1/P2-or-higher findings. Merging and publication each need separate approval. |

## Dependencies and ownership

The APItools active order is S06 then S05, by priority. S06 has no upstream
prerequisite and no technical dependency on S05. Both use the same isolated
branch/worktree, carrying S06's changes into S05 after S06's review. One execution owner controls the ledger, and no
sibling write scope is approved.

Implement on a branch in a separate git worktree outside
`~/Workspace/go.work`; the workspace checkout stays on its current `main`, so
normal workspace builds of OpenUdon and Udon keep their current APItools.
Kinet's Stage 5 is complete, not a pending dependency or merge gate. Isolation
still protects shared builds and accepted evidence. Merging into that checkout
and publishing each need separate approval. S06 does not change probe efficiency:
parallel probing and conditional requests belong to S05's probe primitive.

## Acceptance and verification

The documented `public-apis` source works with its default configuration, and
`auto` no longer fails a no-match query because of a dead host. Existing
exported APIs, wire shapes, probe behavior, and the legacy JSON shape are
unchanged, and the default test suite uses no network.

Required verification:

```bash
go test ./...
GOWORK=off go test ./...
go vet ./...
GOWORK=off go vet ./...
go run ./cmd/cataloggen -check
go run ./cmd/apitools catalog check
go run ./cmd/apitools search --help
git diff --check
```

Run OpenUdon and Udon compatibility checks in workspace and standalone modes
per tech-stack.md, with Ramen excluded. S06.3's candidate verification uses
disposable workspace/modfile configuration outside tracked repositories,
explicitly binding both consumers to the candidate APItools worktree. Resolve
relative paths without changing tracked manifests or the shared workspace, and
check the resolved APItools module directory before test/vet. Candidate-bound
workspace and standalone checks prove compatibility; ordinary standalone runs
against unchanged pins/replacements are separate baseline evidence only.
S05.6 reuses this procedure. No sibling write or adoption is authorized.

Include regressions for a valid empty JSON list, empty/unrecognized Markdown,
malformed rows, and every parser bound, alongside successful Markdown/JSON
search and auto fallback. None may fetch live provider data.

## Review and planning evidence

**Review iterations started:** 1 of at most 10; iteration 1 passed 2026-10-06.
This intake is not a bounded-gate pass. Persist the count once implementation
reaches review and never reset it. An unresolved blocking finding at
iteration 10 requires user direction.

Planning used read-only inspection of the code and sibling usage, DNS and CLI
runs against the dead host, and reads of the public repository's metadata and
README. No implementation verification or milestone review is claimed.

The approved 2026-10-06 reconciliation amends planning only. Four focused
existing tests passed in workspace and standalone modes:
`go test . -run 'Test(AutoFallsBackToPublicAPIs|SearchPublicAPIsFindsValidatedSpec|PublicAPIsProbeBudgetStopsSlowProbes|PartialPublicAPIsProbeBudgetDoesNotCacheSearch)$' -count=1`
and the same command with `GOWORK=off`. `git diff --check` also passed.
These checks establish the unchanged baseline, not acceptance of unimplemented
S06/S05 behavior. The previous live observations were not repeated.

## Execution notes

- S06.1 completed 2026-10-06: added bounded Markdown/JSON list parsing,
  category/link extraction, escaped-pipe handling and fail-closed parse
  diagnostics. Bounds are 20 MiB response, 10,000 entries and 64 KiB per row
  or JSON entry, with shared JSON depth/token checks and cancellation.
  Valid empty JSON arrays remain supported; no partial list reaches probes.
  Focused public-apis/search tests and their race run passed with Go 1.26.6;
  `git diff --check` passed. The first task commit also carries the six
  explicitly approved planning files into the isolated branch.
- S06.2 completed 2026-10-06: the default now selects
  `https://raw.githubusercontent.com/public-apis/public-apis/master/README.md`.
  Local transport/server fixtures prove the unconfigured source URL and CLI
  `public-apis`/`auto` paths, including no-match success. Legacy JSON and probe
  budgets/order remain intact. Focused root/CLI tests and the default/CLI race
  fixtures passed with Go 1.26.6; patch check passed.
- S06.3 implementation/verification completed 2026-10-06: README now
  attributes the runtime-read public-apis list and documents bounds, mirrors,
  empty results and unchanged probe behavior. APItools full tests and vet
  passed with the disposable workspace and with `GOWORK=off`;
  `cataloggen -check`, `catalog check` (zero warnings/errors), search help,
  focused race checks, patch check and a three-second parser fuzz run
  (13,412 executions) passed.
- OpenUdon and Udon full tests/vet passed in workspace and standalone modes
  with explicit candidate module resolution. Disposable modfiles/workspace
  selected this branch; no consumer pin or manifest changed, and Ramen was
  excluded. Udon's first workspace attempt hit a temporary link disk-quota
  limit; its first standalone attempt hit an asynchronous triggerhost TempDir
  cleanup failure in a package that does not import APItools. Both full suites
  passed on retry with `-p 2` and a task-owned build temp directory. These
  failed attempts are retained as verification evidence, not hidden.
- The separately invoked live smoke
  `apitools search --source public-apis --query zzapitoolsnomatchs06 --json`
  fetched and parsed the default README and returned an empty result with one
  passing list attempt, without probing provider operations. Default tests
  use only local fixtures; no upstream list bytes were committed.
- Review iteration 1 started after verification: inspect the entire
  `8a52c3f..HEAD` milestone changes plus pending S06.3 documentation for
  parsing safety, compatibility, bounds, cancellation, fallback/error
  behavior, source attribution and downstream evidence before recording the
  gate outcome.

## Whole-milestone review and closure

- Iteration 1 passed: no P1/P2-or-higher or carried lower findings.
  Reviewed the full milestone implementation and task commits, including
  JSON/Markdown selection, empty legacy arrays, size/count/row/depth bounds,
  fail-closed errors before probing, cancellation, deterministic ordering,
  shared transport safety, unchanged exported shapes/probe limits, source
  attribution, and candidate-bound consumer tests.
- Evolution decision: no bump. This repairs the existing source contract;
  v26 continues to describe planned S05 behavior.
- S05 reconciliation: retain all six pending rows and its separate design
  approval checkpoint. S06 changes do not implement S05's conditional/concurrent
  probe; its new primitive still leaves legacy search behavior unchanged.
  Consumer adoption and publication remain separate.
- Consolidation: current list/default/bounds contracts are in architecture
  and tech-stack; README has public attribution. Product scope and reusable
  lessons did not materially change. No current knowledge was removed.
  Retire the full specification/status under the project's standard procedure.
- Udon tests generated an untracked directory under its CLI package; the
  known test artifacts were moved into disposable verification storage.
  Both sibling Git worktrees are clean and their tracked files/pins unchanged.
````
