# Status S06 — Public-apis source repair

**State:** Execution started, 2026-10-06, after confirmation of the S06 -> S05
goal with task commits and no external mutations. S06.1-S06.2 are complete on
branch `work/s06-s05` in the isolated worktree; S06.3 remains pending.
The approved planning amendments and the main workspace checkout are preserved.
No merge, publication, or closing review has started.

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
| S06.3 — Documentation, compatibility, and verification | `[ ]` | Depends on S06.1-S06.2. Owns consistency-intake F08-F09. Document the list source and MIT attribution in README and architecture. Run the verification below against the candidate worktree, recording ordinary unchanged-pin/replacement checks separately. Keep the live GitHub-raw smoke check separately opt-in, not a default test. Preserve worktree isolation without treating completed Kinet Stage 5 as a pending gate. Pass the persisted ten-iteration review gate with no open P1/P2-or-higher findings. Merging and publication each need separate approval. |

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

**Review iterations started:** 0 of at most 10; closing gate not started.
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
