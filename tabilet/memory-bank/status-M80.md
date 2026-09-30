# Status M80 — Catalog discovery API for step contracts

**State:** Execution started, 2026-09-30, under the confirmed Stage 5 goal.
M81 is retired at closure `d32995f628b9bdf91b6963d566fbac9c16fe6afe`
(qualified source `8d67aef2dce565aa9ac8e8b2c56a9100dcf9a5a6`). M80.1 is
complete; M80.2 is complete, M80.3–M80.5 are pending. No publication or closing review started.

**Specification:** [M80](milestone.md#m80--catalog-discovery-api-for-step-contracts).

**Provenance (first intake):** "APItools stage 5 draft — M80" and its draft
status, from the Kinet stage 5 draft plan (revision 2, 2026-09-30; staging
material outside Kinet's tracked history), with that plan's deep review
(finding 2) and Kinet `docs/request-resolution.md` G1/S2d. The source draft
has no finding IDs of its own; F01-F06 were assigned in source order during
reconciliation. Review and revalidation baseline:
`8580ff2485a3faff139b17bdc0b8e78d2af4ce18`; worktree clean during assessment
and immediately before planning writes. F01's source priority is G1's
"highest priority" (no P-level supplied); source priority for F02-F06 is
`not supplied`. All six local classifications are P2 acceptance-plan gaps,
not defects in the existing explicit-source API. The user approved the
complete dispositions, ownership, dependencies, and seven planning-file
actions on 2026-09-30. Execution and release authority remain separate.

**Provenance (second intake):** Review "APItools M80 — Catalog discovery API
for step contracts" (`apitools-m80-review.md`, 2026-09-30), at the same
baseline, with these uncommitted planning files as evidence. Its P2 findings,
L1, and L2 are owned by prerequisite [M81](../docs/history/status-M81.md) and by the amended
rows below. L3 (links to untracked drafts), L4 (finding labels without finding
text), and L5 (missing provenance) were fixed directly in this planning text.
The user approved the dispositions, M81, and these amendments on 2026-09-30.

**Lineage:** The [M77](../docs/history/status-M77.md),
[M78](../docs/history/status-M78.md), and
[M79](../docs/history/status-M79.md) operation-metadata foundations are
completed and retired. M80 adds catalog discovery; it does not reopen their
acceptance or rewrite frozen history. The strict catalog and safe artifact
helpers remain existing foundations.

## First-intake findings (F01-F06)

The claims below restate the source draft in source order. They were
reconstructed at the second intake from the draft text and each row's
ownership, because the first intake recorded labels without claims. Parts not
recorded at the first intake are marked as such.

| Label | Draft claim | Disposition | Unconfirmed part | Owner |
|---|---|---|---|---|
| F01 | A library API takes a step contract and returns ranked providers and artifacts from the whole local catalog. | confirmed | none | M80.2 |
| F02 | Each result carries authority, license and redistribution fields, authentication needs, effect class, and rank evidence. | partially confirmed | Not recorded at the first intake; the row limits license and redistribution to recorded evidence or explicit unknowns, since the catalog has only a free-text `LicenseNote`. | M80.1, M80.2 |
| F03 | When nothing qualifies, return an explicit "no API for this task" decision instead of a weak match. | decision-dependent; resolved by the approved scoped-outcome policy and M81 decision 5 | not applicable | M80.3 |
| F04 | Ranking is deterministic and offline over the local catalog; an optional remote tier such as APIs.guru runs only when explicitly configured. | partially confirmed | Not recorded at the first intake. | M80.2, M80.4 |
| F05 | Acceptance includes consumer compatibility with OpenUdon's `step discover` (M94), and the accepted revision is published for M94. | decision-dependent; resolved as separate producer and adoption gates | not applicable | M80.5 |
| F06 | There is no technical upstream dependency. | partially confirmed at the first intake; superseded at the second intake, where M81 became the upstream prerequisite | The first intake held it only if unknown evidence stayed explicit. | M80.1 |

| Item | State | Notes |
|---|---|---|
| M80.1 — Contract and fixtures | `[+]` | Owns F06 and F02's contract part, plus the M80 parts of review findings P2-4, P2-6, and L2. Depends on M81. Implement the additive versioned request/report approved at M81.1 over the existing StepContract: optional provider constraints resolved by exact catalog key; catalog-stable provider/artifact identity with a reference that round-trips into M81.5's artifact-scoped export; authority and license filters with unknowns included by default and labeled; auth, effect, and rank metadata; coverage; and the five outcomes. Freeze fixtures for matches, ambiguity, ties, filtering, unknowns, no root, no-match, insufficient evidence, blockers, and root relocation. Evidence: `step_metadata.go`, `catalog/provider.go` (free-text LicenseNote), `docs/catalog-discovery.md` (from M81.1). Existing exported APIs and wires stay unchanged. |
| M80.2 — Local catalog retrieval | `[+]` | Owns F01, F02's retrieval part, F04's local part, and the M80 part of P2-1. Depends on M80.1. Read the M81 operation index under an explicitly supplied root, with no per-call parsing of catalog artifacts. Verify the index version and registration digests; report stale, oversized, missing, invalid, or unsupported artifacts as unexamined scope. Without a root, return metadata-only provider leads and insufficient evidence naming what would resolve it. Return authority, license/redistribution evidence or explicit unknowns, authentication needs, effect, and source-qualified ranking inputs. Keep artifacts without indexed operations as reference-only leads; record coverage, exclusions, missing evidence, and limits; bound local time and memory. Evidence: `operation_candidates.go` requires explicit sources; M81.4 supplies the index. |
| M80.3 — Contract ranking and outcomes | `[ ]` | Owns F03 and the M80 parts of P2-5, P2-6, and L1. Depends on M80.2. Reuse `prepareStepContract` and dimension matching over index entries with exact digests and native selectors. Apply the approved qualification rule so one shared purpose term cannot make a match; report ambiguity; break ties on catalog-stable identity, never machine paths. Distinguish match, ambiguous, no qualifying API within the checked scope, insufficient evidence, and blocked; only a scoped no-match is browser-routable, and curated catalog facts are evidence only. Add iCoT `CatalogPlan`-parity fixtures for named providers and a small labeled relevance set derived from OpenUdon `examples/eval` step shapes (no user data), recording precision@k as a baseline. Evidence: `operation_rank.go:276-297,803-826,845-846`. Browser routing, provisioning, account choices, approval, and execution remain downstream. |
| M80.4 — Opt-in remote lookup | `[ ]` | Owns F04's remote part. Depends on M80.3. Default to offline with no lookup or implicit provisioning; explicitly configured lookup initially uses APIs.guru through existing guarded search/download. Bound total network time to eight seconds, fetched source documents to three, and each document to 20 MiB, matching iCoT's remote lookup (OpenUdon `internal/icot/elicitor/remote_sources.go`). Report each remote result's fetched digest and final URL. Preserve safe URL/redirect/dial-time checks and parser/prompt limits; report cancellation, timeout, empty results, and partial evidence. Keep remote/public-catalog provenance distinct from reviewed official or docs-derived sources. Evidence: `client.go`, `download.go`, `remote_discovery.go`, and `operation_candidates.go`. Use local-server fixtures only; no general crawling, credentials, or provider operations. |
| M80.5 — Qualification and publication handoff | `[ ]` | Owns F05. Depends on M80.1-M80.4. Run the verification below, including an opt-in, provider-free scale benchmark over a populated local root; document the additive contract and evidence limitations; and pass the persisted ten-iteration review gate with no open P1/P2-or-higher findings. Qualify existing OpenUdon/Udon consumers without depending on future M94 implementation; Ramen is excluded. The handoff states that consumers supply the catalog root explicitly and must not rely on iCoT's sibling-checkout default. Publication remains a required handoff but requires separate explicit commit/push/publication authority; it carries M81's changes. Record the actual accepted full revision and module version once authorized and published. OpenUdon's proposed M94 owns exact-pin adoption; APItools does not edit sibling ledgers or code. Evidence: the Kinet stage 5 draft plan's OpenUdon M94 status waits for accepted, published M80. Do not complete the publication portion without authority and observed evidence. |

## Dependencies and ownership

The complete APItools active order is M81 then M80, with M80 task order
M80.1 -> M80.2 -> M80.3 -> M80.4 -> M80.5. One execution owner controls the
ledger; no parallel implementation or sibling write scope is approved.
M81 is the technical upstream prerequisite: every M81 row must be complete and
its review gate closed before M80.1 starts. Completed M77-M79 satisfy the
existing operation-metadata foundation.

The external handoff is qualified, authorized, published M80 -> proposed
OpenUdon M94 -> proposed Kinet W10. M94 retains its proposed M93 prerequisite;
proposed OpenUdon M95 must retain the replacement discovery journey as an iCoT
removal gate. These are draft references, not active sibling allocations made
by this reconciliation. M80 producer qualification must not wait for M94 to
implement an API that M94 itself waits for M80 to publish. Existing consumer
compatibility checks do not assert that the future discovery journey works.
The M94 and W10 drafts must absorb the changes listed in M81.1's design record
at their promotion; APItools does not edit them.

Catalog-wide license audits, provider expansion, and typed-overlay upgrades
remain unnumbered candidates in milestone.md. Unknown or missing evidence is
explicit in M80; no catalog batch or proposed catalog ID is promoted here.

## Acceptance and verification

Implement the specification's additive discovery contract and all five rows
on top of M81. Fixtures must prove deterministic ordering under input
permutations and root relocation; matches, ambiguity, and ties; provider
constraints and `CatalogPlan` parity; authority/license filters; unknown
license/redistribution, auth, and effect evidence; no root, empty, and stale
indexes; scoped no-match versus insufficient evidence; unexamined scope;
bounded work/results; cancellation; unsafe paths/hosts and digest failures;
and offline defaults. Operations remain source-backed with exact
digests/native selectors. No absence or unexamined scope is proof of global
impossibility. The library neither provisions sources nor chooses browser
fallbacks, accounts, credentials, approval, or runtime execution.

Required producer verification:

```bash
go test ./...
GOWORK=off go test ./...
go vet ./...
GOWORK=off go vet ./...
go run ./cmd/cataloggen -check
go run ./cmd/apitools catalog check
go run ./cmd/apitools search --help
go run ./cmd/apitools import --help
git diff --check
```

Run affected OpenUdon and Udon compatibility checks, including workspace and
standalone modes where available, per tech-stack.md. Default checks remain
provider-free/offline; all new network tests use local servers, and the scale
benchmark is opt-in. Before claiming M80 acceptance, require the bounded
full-milestone review, documented qualified fixtures, and the authorized
publication handoff. M94 owns subsequent adoption and its exact-pin
discovery/provisioning journey checks. Obtain authority before release
actions; planning approval and a status marker grant none.

## Review and planning evidence

**Review iterations started:** 0 of at most 10; closing gate not started.
Both intakes are ordinary approved intake, not bounded-gate passes. Persist
and resume the count once implementation reaches review, including any
qualification, publication, and closure review passes; never reset it. An
unresolved blocking finding at iteration 10 requires user direction.

Assessment-only checks at the revalidation baseline passed: focused existing
operation-candidate/ranking and catalog tests in workspace and standalone
modes, `go run ./cmd/cataloggen -check`, and `git diff --check`. They establish
existing building-block behavior, not acceptance of the unimplemented
discovery API. The second intake added read-only inspection of the ranking
budgets, parser limits, local cache, and iCoT's catalog path, plus a read-only
harness validation. No full implementation verification or milestone review
is claimed. Evolution V25 records the catalog-discovery target, including the
M81 prerequisite, as planned.

### Upstream reconciliation — M81 accepted, 2026-09-30

M81 completed and passed closing review in iteration 3. Exact qualified source:
`8d67aef2dce565aa9ac8e8b2c56a9100dcf9a5a6`; ordinary retirement follows in a
separate closure commit. No M80 row is implemented by this reconciliation.
The technical prerequisite is satisfied; M80 remains sequential approved work.

Use `CatalogIndexOptions` with an explicit `catalog.RootOptions`, selected
catalog and read-only `sqlitecache.ReadCatalogSpecArtifacts` callback.
`ReadCatalogOperationIndex` consumes validated `apitools.catalog-operation-index/v1`
metadata without source parsing; retain complete coverage and stale-generation
candidate removal. `CatalogIndexedArtifact` groups shared raw identity with
all canonical provider/spec links and explicit advisory flags. Partial metadata
is positive evidence with unexamined coverage, never a definitive no-match.
`CatalogArtifactReference` plus `ExportCatalogArtifacts` now provide the exact
provisioning round trip; keep selectors native and leave binding/approvals to
OpenUdon. Existing APIs and exported option shapes are unchanged.

Discovery must respect the index's coverage and fixed budgets, then apply its
own bounded query/result/context work. No root stays metadata-only/insufficient;
missing/invalid/stale evidence must not route automatically to a browser.
The default SQLite reader is migration/pruning/access-time free. Final Udon
standalone compatibility needs a disposable modfile for its preexisting
module metadata updates; record that distinction, never edit Udon dependencies
from APItools. OpenUdon can qualify actual APItools adoption with a disposable
replacement before publication. M80's publication carries M81 and waits for
review of the exact outgoing diff under the confirmed launch policy. No producer
qualification depends on future M94 implementation or changes sibling ledgers.

### M80.1 execution evidence — 2026-09-30

Implemented the additive `apitools.catalog-discovery/v1` request/report types,
strict 64-KiB decoder, fixed lookup/result/context/time limits and exact-key
validation over the existing StepContract. Installation root/adapter/client
options serialize as `{}` and cannot enter request JSON. Multiword names remain
one exact key; invalid keys/version/limits are diagnosed rather than broadened.
Nil/omitted/null provider keys mean open scope; explicit empty arrays preserve
empty scope through JSON serialization. Structured unknown license and
redistribution evidence remain distinct from verbatim notes. Candidates retain
M81-native artifact references, metadata, source evidence, qualification counts,
coverage/exclusions and tied rank data under all five outcomes.

Frozen synthetic wire examples cover match, ambiguity despite a display limit,
ties, filtering, unknowns, missing root, scoped no-match, stale/insufficient
scope, blockers and relocation. They are serialization fixtures, not observed
lookup runs or claims that retrieval/ranking is implemented. The synthetic tie
case declaration is explicitly documented. Behavioral source-backed tests and
lookup/ranking are M80.2/M80.3; remote lookup remains M80.4. No old exported API
or wire changed. README/design/current architecture describe the exact scope.

Focused decoder/exact-provider/empty-scope/bounds/unknown/wire/relocation tests,
APItools workspace/standalone full tests and vet, OpenUdon/Udon full workspace
consumer tests and diff checks pass. No model, provider operation, credentials,
sibling edits or publication. Closing review has not started (count 0/10).

### M80.2 execution evidence — 2026-09-30

Delivered local metadata retrieval through `DiscoverCatalogOperations` and the
validated M81 index/read-only registry adapter, with installation paths/client
options outside request JSON. It selects exact canonical providers, preserves
shared native references, authority or explicit unknowns, bounded verbatim
license notes, unchanged auth alternatives/effect metadata and named evidence
filter exclusions. No root/index yields metadata-only leads and insufficient
evidence. Invalid configuration/index identity blocks; stale entries contain
no candidates. Missing/unsupported scope and per-query `work_limit` coverage
remain explicit. Report projections and local work/time are bounded. Notes
over 64 KiB are omitted with a source evidence gap, never altered into a
verbatim claim or inferred permission. No raw source parsing, fetching,
provisioning, migration or access-time writes occur.

This row provides unqualified metadata candidates; M80.3 owns actual contract
ranking and five-outcome qualification, and M80.4 owns explicit remote support.
The intermediate function therefore never claims a match or scoped no-match
from metadata alone. Existing public APIs and wires remain unchanged.

Focused shared/source-deletion-without-reparse, relative identity/relocation,
exact multiword scope, no-root/missing/stale/corrupt/unsafe/empty/cancel cases,
license/authority exclusions, work/context bounds and oversized raw-note gap
regressions pass. APItools workspace/standalone full tests and vet, OpenUdon/Udon
full workspace consumer suites and diff checks pass. No sibling edits,
provider/model calls or publication. Closing review remains 0/10.
