# Retired milestone M80 - Catalog discovery API for step contracts

**Milestone.** M80
**Outcome.** completed
**Retired.** 2026-09-30
**Source status.** tabilet/memory-bank/status-M80.md
**Source specification.** tabilet/memory-bank/milestone.md#m80--catalog-discovery-api-for-step-contracts
**Evidence.** fb132631c9827eae5f2ec4503d03f21eabfb4113
**Worktree.** includes uncommitted changes
**Review.** passed
**Review iterations.** 3
**Verification.** APItools workspace/standalone tests and vet, catalog generation/quality (zero errors/warnings), existing help and diff checks pass; OpenUdon/Udon workspace and actual-source standalone adoption pass through disposable modfiles without changing consumer modules. Source-backed conformance and 12,000-operation benchmark evidence are retained below. Exact qualified source is published and its module version observed.
**Consolidated into.** [product](../../memory-bank/product.md), [architecture](../../memory-bank/architecture.md), [tech stack](../../memory-bank/tech-stack.md), [coverage lesson](../../memory-bank/lessons.md#keep-incomplete-catalog-coverage-distinct-from-negative-evidence), [design](../../../docs/catalog-discovery.md), operator README and reconciled downstream package-local status files.

## Milestone specification

````markdown
## M80 — Catalog discovery API for step contracts

Add an additive, versioned library discovery contract that takes a step's
purpose, typed inputs/outputs, effect, and optional provider constraints, and
returns ranked providers and artifact references across the local catalog by
reading the M81 operation index. Reuse the operation-candidate metadata and
conservative ranking delivered by M77-M79; do not reopen their retired
records. This shared public contract belongs in lane M. The active order is
M81 then M80; implementation stays sequential under one execution owner, with
no parallel ownership or sibling writes authorized.

The user approved this reconciliation of "APItools stage 5 draft — M80"
(Kinet stage 5 draft plan, revision 2, 2026-09-30; staging material outside
Kinet's tracked history) on 2026-09-30. Source context is that draft plan's
deep review (finding 2) and Kinet's
[request-resolution G1/S2d](../../../kinet/docs/request-resolution.md). Review
and revalidation baseline are both
`8580ff2485a3faff139b17bdc0b8e78d2af4ce18`; the worktree was clean during that
assessment and immediately before its planning writes. Assigned source finding
labels F01-F06, their evidence, and their owners are retained in
[status-M80.md](status-M80.md). Local P2 classifications describe gaps in the
proposed acceptance, not regressions in the existing API. The draft was
promoted before Kinet's stage 4 retired, consistent with the draft plan's note
that APItools M80 does not depend on stage 4.

A second approved intake on 2026-09-30 applied review "APItools M80 — Catalog
discovery API for step contracts" (`apitools-m80-review.md`) at the same
baseline, with the uncommitted planning files as evidence. It moved the six
design decisions and the pre-existing limits into prerequisite
[M81](../docs/history/status-M81.md) and amended the scope and rows
below. Neither intake starts the closing review counter.

**Scope and acceptance.**

- Implement the contract approved at M81's design checkpoint. Optional
  provider constraints resolve by exact catalog key, and each result records
  whether it came from a constraint or from open retrieval.
- Retrieve deterministically from the M81 operation index under an explicitly
  supplied catalog root, with caller-visible authority and
  license/redistribution filters and rank evidence. Do not parse catalog
  artifacts per call. Verify the index version and registration digests;
  stale, oversized, missing, invalid, or unsupported artifacts are reported as
  unexamined scope. Without a root, return metadata-only provider leads and
  insufficient evidence that names what would resolve it. Report catalog
  identity, searched scope, exclusions, missing evidence, limits, and
  cancellation, and bound local discovery time and memory. A catalog reference
  without indexed operations remains a reference-only lead, not an invented
  operation or positive contract match.
- Every indexed local operation result carries catalog-stable provider/artifact identity and a
  reference that round-trips into M81's artifact-scoped export, source
  authority, license and redistribution evidence or explicit unknowns
  (unknowns included by default and labeled, `license_note` passed through
  verbatim, no inferred permission), authentication needs,
  read/write/unknown effect, and concise rank evidence. Keep official sources,
  docs-derived advisory artifacts, and public-catalog provenance distinct;
  auth evidence stays scoped to the selected spec, with alternatives intact.
  Existing license notes such as "terms apply" do not establish permission.
- Reuse step-contract preparation and dimension matching only where source
  evidence supports them. Apply the approved qualification rule, so a single
  shared purpose term cannot make a match. Break ties on catalog-stable
  identity, never on machine paths. Preserve exact digests and native
  selectors, ties, typed input/output gaps, and conservative unknown effects.
  Model consumers may choose only returned source-backed operations; scores
  never approve or bind an operation.
- Distinguish match, ambiguous, no qualifying API within the checked scope,
  insufficient evidence, and blocked. Missing documents, unsupported
  semantics, indeterminate comparisons, unexamined scope, and failed reads
  must not become a definitive no-API result. Only a scoped no-match is
  browser-routable, and curated catalog facts are evidence, never proof.
  Never claim that a task is impossible. Browser routing and source
  provisioning belong to OpenUdon/Kinet; this library performs neither.
- Default to offline operation. An explicitly configured remote tier initially
  uses APIs.guru through the existing guarded search/download path, with an
  eight-second total network deadline, at most three source documents, and
  the 20 MiB per-document bound, matching iCoT's remote lookup (OpenUdon
  `internal/icot/elicitor/remote_sources.go`). Report each remote result's
  fetched digest and final URL so later provisioning can detect drift.
  Preserve URL/redirect/dial-time unsafe-host rejection, cancellation, prompt
  budgets, and visible timeout/empty/partial outcomes. No general crawling,
  provider operation, credential resolution, or implicit artifact
  provisioning is permitted.
- Add fixtures for ordering under input permutations and root relocation,
  matches, ambiguity, ties, provider constraints including iCoT `CatalogPlan`
  parity, authority/license filters, unknown metadata, no root, empty and
  stale indexes, scoped no-match, unexamined scope, cancellation, unsafe
  paths/hosts, digest failures, and offline defaults. Add a small labeled
  relevance set of step contracts derived from OpenUdon's `examples/eval` step
  shapes, with no user data, and record precision@k as a baseline rather than
  a pass/fail gate. Add an opt-in, provider-free scale benchmark over a
  populated local root. Keep existing exported APIs and wire shapes unchanged;
  introduce the discovery report as a separate additive contract. Run
  workspace and standalone tests/vet, catalog generation/quality checks, CLI
  help, diff checks, and affected OpenUdon/Udon compatibility checks, with
  Ramen excluded. Network tests use local servers and default verification
  remains provider-free.

**Dependencies and downstream acceptance.** M81 is the technical upstream
prerequisite: its approved design record, catalog root contract,
large-artifact limits, operation index, and artifact-scoped export must be
complete, and its review gate closed, before M80.1 starts. M77-M79 supply the
existing operation-metadata foundation and are completed, retired
dependencies. Catalog-wide license audits, provider expansion, and
typed-overlay upgrades remain unnumbered candidates; M80 represents their
missing evidence explicitly rather than requiring or claiming their delivery.
The coordinated handoff is M80 qualification and authorized publication ->
OpenUdon's proposed M94 -> Kinet's proposed W10. M94 also retains its proposed
M93 prerequisite; OpenUdon's proposed M95 iCoT removal gate requires the
replacement discovery journey. These sibling IDs are draft references, not
allocations or execution authority in this ledger. The handoff states that
consumers must supply the catalog root explicitly and must not rely on iCoT's
sibling-checkout default.

Qualify the producer and existing consumers independently of future M94
adoption; M94 owns re-pinning and verification against the exact accepted,
published revision. M80 must not wait for M94 implementation to qualify,
because M94 itself waits for M80. Publication remains a required handoff gate,
but planning approval authorizes no commit, push, tag, release, or sibling
edit. Obtain separate authority before publishing and record the actual
accepted revision; leave the publication portion incomplete while authority
or evidence is missing. Preserve the ten-iteration review count across
qualification/publication and normal milestone closure; no counter reset.
````

## Status record

````markdown
# Status M80 — Catalog discovery API for step contracts

**State:** All tasks complete, 2026-09-30, under the confirmed Stage 5 goal; final publication/closure review iteration 3 passed.
M81 is retired at closure `d32995f628b9bdf91b6963d566fbac9c16fe6afe`
(qualified source `8d67aef2dce565aa9ac8e8b2c56a9100dcf9a5a6`). M80.1 is
complete; M80.2–M80.3 are complete; M80.4 is complete; M80.5 is complete. Qualified source is published at `fb132631c9827eae5f2ec4503d03f21eabfb4113`; all acceptance and review gates passed; closure publication is the remaining operational handoff.

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
| M80.3 — Contract ranking and outcomes | `[+]` | Owns F03 and the M80 parts of P2-5, P2-6, and L1. Depends on M80.2. Reuse `prepareStepContract` and dimension matching over index entries with exact digests and native selectors. Apply the approved qualification rule so one shared purpose term cannot make a match; report ambiguity; break ties on catalog-stable identity, never machine paths. Distinguish match, ambiguous, no qualifying API within the checked scope, insufficient evidence, and blocked; only a scoped no-match is browser-routable, and curated catalog facts are evidence only. Add iCoT `CatalogPlan`-parity fixtures for named providers and a small labeled relevance set derived from OpenUdon `examples/eval` step shapes (no user data), recording precision@k as a baseline. Evidence: `operation_rank.go:276-297,803-826,845-846`. Browser routing, provisioning, account choices, approval, and execution remain downstream. |
| M80.4 — Opt-in remote lookup | `[+]` | Owns F04's remote part. Depends on M80.3. Default to offline with no lookup or implicit provisioning; explicitly configured lookup initially uses APIs.guru through existing guarded search/download. Bound total network time to eight seconds, fetched source documents to three, and each document to 20 MiB, matching iCoT's remote lookup (OpenUdon `internal/icot/elicitor/remote_sources.go`). Report each remote result's fetched digest and final URL. Preserve safe URL/redirect/dial-time checks and parser/prompt limits; report cancellation, timeout, empty results, and partial evidence. Keep remote/public-catalog provenance distinct from reviewed official or docs-derived sources. Evidence: `client.go`, `download.go`, `remote_discovery.go`, and `operation_candidates.go`. Use local-server fixtures only; no general crawling, credentials, or provider operations. |
| M80.5 — Qualification and publication handoff | `[+]` | Owns F05. Depends on M80.1-M80.4. Run the verification below, including an opt-in, provider-free scale benchmark over a populated local root; document the additive contract and evidence limitations; and pass the persisted ten-iteration review gate with no open P1/P2-or-higher findings. Qualify existing OpenUdon/Udon consumers without depending on future M94 implementation; Ramen is excluded. The handoff states that consumers supply the catalog root explicitly and must not rely on iCoT's sibling-checkout default. Publication remains a required handoff but requires separate explicit commit/push/publication authority; it carries M81's changes. Record the actual accepted full revision and module version once authorized and published. OpenUdon's proposed M94 owns exact-pin adoption; APItools does not edit sibling ledgers or code. Evidence: the Kinet stage 5 draft plan's OpenUdon M94 status waits for accepted, published M80. Do not complete the publication portion without authority and observed evidence. |

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

**Review iterations started:** 3 of at most 10; iteration 2 passed with no open P1/P2 findings; iteration 3 publication/closure review passed with no open P1/P2 findings.
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

### M80.3 execution evidence — 2026-09-30

Implemented discovery-specific strong purpose qualification over the existing
`prepareStepContract` / typed dimension comparisons, without changing the older
candidate API. Scores cannot hide missing required inputs, unused requested
inputs, nullable/unknown selected outputs, constrained unknown effects,
incomplete auth alternatives or sanitation/compaction loss. Shared artifact
links count once per native operation. Ranking and ties use canonical catalog
identity, never host paths. All qualification counts precede display limits;
multiple matches remain ambiguous even with one displayed result. Projection
that omits every qualified result becomes insufficient. Positive matches retain
unexamined coverage. No qualifying API requires complete scope and definitive
conflicts for every examined operation; weak wording/empty/unexamined scope do
not prove absence. The cooperative query deadline includes retrieval, ranking
and projection; caller cancellation clears positives.

Source-backed fixtures verify strong versus one-term/identity-only wording,
ambiguity/ties/display limits, typed conflicts versus unknown evidence, missing
auth, omitted fields, selected nullable outputs, constrained unknown effects,
unused inputs, positive partial scope and selected-reference M81 export.
Named-provider CatalogPlan parity preserves exact canonical selections and
multiword keys. Three repository-owned OpenUdon eval step shapes yield observed
precision@1 3/3 (1.000); this is a recorded baseline, not an acceptance threshold.
Slack qualifies, Jira remains insufficient after credential-shaped `key`
response removal, and Drive remains insufficient because the existing effect
classifier does not recognize upload as write. Fixtures preserve these source
limitations rather than treating them as no API. Provenance and labels live in
`testdata/catalog-discovery/relevance/README.md`.

Focused behavioral checks, APItools full workspace/standalone tests and vet,
OpenUdon/Udon full workspace tests and diff checks pass. No models, provider
operations, credentials, sibling edits or publication. Closing review remains
0/10. Remote lookup is still the separately pending M80.4 task.

### M80.4 execution evidence — 2026-09-30

Delivered explicitly enabled APIs.guru lookup through existing guarded search/
download, with both request and installation opt-in. Installation client and
endpoint stay outside JSON. Copied clients drop caches, cookie jars and redirect
callbacks; default URL/host/port, redirect and dial-time protection remains.
Bounds are eight cooperative seconds, three fetched documents, 20 MiB each
(stricter configured caps stay effective), existing inventory/prompt limits and
16 MiB candidate metadata per source. Default calls make no remote requests.

Every fetched source retains public-catalog authority, sanitized final URL,
exact digest/bytes and ephemeral/provisioning status in a lead; operations retain
that remote evidence too. Local/remote digest/native duplicates count once.
Remote-only candidates have no fabricated registered reference. Constrained
remote sources require exact catalog URL linkage and preserve all matching
provider/spec links. Filters exclude named evidence before source fetching.
Lookup empty/failure/timeout/partial parsing never establishes absence; positive
local/remote matches keep coverage gaps. Unsafe URL/transport/redirect refusals
block; hard caller cancellation clears positives. Nothing is registered,
persisted, provisioned or executed.

Local HTTP/TLS fixtures verify opt-ins and zero default calls, the three-document
cap, duplicate operations, each fetched digest/final URL, query-secret omission,
no cookie/authorization propagation, ephemeral remote-only matches, shared exact
URL provider links, empty/invalid/oversized evidence, timeout/cancellation,
unsafe lists/sources/userinfo/redirects, provider constraints and filtering.
APItools full workspace/standalone tests/vet and OpenUdon/Udon full workspace
consumer tests pass; diff checks pass. No external network, models, provider
operations, credentials or sibling changes. Closing review remains 0/10.

### M80.5 qualification and review — iteration 1 started, 2026-09-30

All prior task rows are committed. Qualification includes the new opt-in scale
benchmark, visible-tie projection remapping and associated evidence below.
Required APItools workspace/standalone tests and vet, catalog generator/quality,
search/import help and diff checks pass. Existing OpenUdon/Udon workspace suites
pass. OpenUdon literal standalone tests and standalone actual-current-APItools
adoption tests pass (`/tmp/openudon-m80-final-adoption-p222qx85/test.log`). Udon
standalone actual adoption passes via a disposable modfile, preserving its
original modules (`/tmp/udon-m80-final-standalone-8zsa5dqr/test.log`); its preexisting
literal module-update requirement is unchanged, as recorded at M81. This does
not claim future M94 discovery/adoption qualification. No sibling module or
ledger was changed by the producer checks.

The explicitly invoked provider-free populated-root benchmark prepared a new
12,000-operation synthetic artifact with exact registrations and a complete M81
index (one shared artifact, zero unexamined entries). It measured 2.724309939
seconds/query, 12,000 examined/qualified operations, 95,668 JSON report bytes,
787,978,664 total allocated bytes/query and 237,676 KiB maximum process RSS.
Source, registry and index digests were byte-identical before/after querying.
Evidence: `/tmp/apitools-m80-scale-rqhtf1mj/{index.json,query.txt,query-memory.txt,before.json}`.
These are observations, not throughput or RSS guarantees. Default tests skip
this benchmark unless the operator selects an existing indexed root. The fixture
preparation recipe optionally generates 1–100,000 operations in a new directory.

Iteration 1 now reviews all M80 implementation and uncommitted qualification
changes against the full scope, safety/compatibility/failure semantics, tests,
documentation and precise upstream contract. Publication and retirement are
still incomplete; M80.5 remains the sole in-progress row.

#### Review iteration 1 findings — recorded before fixes

- P2 R80-1: Selected field removal can become a false scoped no-match. Existing
  summaries omit credential-shaped fields; the output comparator can then call
  the requested removed field definitively absent. Qualification rejects the
  lossy candidate, but no-match evaluation still trusts that incompatible
  dimension. Owner: M80.5. Require unresolved metadata loss to remain insufficient
  rather than proving absence, with a raw-source-backed selected-field fixture.
- P2 R80-2: Wrapped remote timeout/cancellation errors use equality rather than
  `errors.Is`, losing the explicit timeout classification despite incomplete
  scope. Owner: M80.5. Preserve wrapped cancellation/timeout diagnostics and
  verify a local-server transport-timeout path.
- Review of stable ordering, work-limit traversal and scope documentation is
  still in progress at iteration 1; do not advance or reset the counter.

- P2 R80-3: The M81 reader accepts equivalent index arrays in any order, but
  discovery uses their stored operation/link order for work-limit traversal
  and shared references. A valid permuted index changes report ordering and
  potentially outcomes under a work cap. Owner: M80.5. Canonicalize discovery's
  traversal/reference order without altering M81 storage and verify equivalent
  arrays with and without a work cap.

- P2 R80-4: A valid catalog provider may have no source references (unknown or
  unavailable availability). M81 correctly covers catalog references, but
  discovery omits that selected provider's unexamined scope. Alongside complete
  conflicting operations from another provider, it can incorrectly report a
  complete no-match. Owner: M80.5. Retain a provider-level lead/unsupported
  coverage or an explicit filter exclusion; curated unavailability never proves
  API absence. Verify an indexed mixed-provider scope and no-root leads.

- P2 R80-5: Discovery's global context projection does not reapply the existing
  32-KiB per-operation prompt bound after typed comparison adds rank evidence.
  A candidate can fit the global report but exceed the existing operation
  budget. Owner: M80.5. Keep oversized ranked operations as reference-only
  source leads with incomplete evidence; never expose an oversized operation
  or turn its omission into a no-match. Verify a contract-expanded candidate.
- P2 R80-6: Retrieval copies all shared provider references for every operation
  before the 32-link display guard runs. A valid large shared group can multiply
  reference memory far beyond the query's intended bounds. Owner: M80.5. Apply
  the existing 32-link limit before per-operation allocation, retain all allowed
  source links as leads with named unexamined coverage, and allow an explicit
  narrower provider constraint to recover supported operation retrieval.

#### Review iteration 1 fixes applied — required verification passed

R80-1 now treats critical metadata loss as incomplete evidence even when retained
comparisons look incompatible; the raw required `api_key` output regression
reproduced false no-match before the fix and now passes. R80-2 uses wrapped error
classification and a dedicated remote-timeout diagnostic, exercised with a local
transport timeout. R80-3 canonicalizes index link, operation and artifact order
before query traversal; permuted arrays reproduced report differences before
fixing, and now serialize identically with and without a work cap. R80-4 retains
provider-level unsupported coverage/leads (or named exclusions), including
no-root lookup; the mixed-scope no-reference regression reproduced the issue.
R80-5 reapplies the 32-KiB ranked-operation bound, replacing oversized operations
with reference-only leads rather than publishing oversized prompt metadata or
negative evidence. R80-6 applies the 32-link limit before per-operation reference
copies; a 33-provider source stays unexamined with all links as leads, and a
single-provider constraint recovers the supported query. Current product,
architecture, technical docs, README/design and durable coverage lesson match
these contracts. Focused verification passes; the required full checks and
iteration-2 whole-milestone review remain before qualification.

#### Review iteration 1 verification complete; iteration 2 started

All six recorded P2 findings have fixes and focused regressions. The operation
budget fixture is source-backed and measures 27,609 bytes before comparison;
rank evidence crosses the existing 32-KiB limit, so no oversized operation is
published. Full workspace/standalone APItools tests/vet, catalog generator/quality
and diff checks pass after these fixes. The repeated populated-root benchmark
still examines/qualifies 12,000 operations and preserves source/registry/index
digests. After fixes: 2.717633410 seconds/query, 95,668 report bytes, 844,918,144
cumulative allocation bytes/query and 237,792 KiB maximum process RSS.
Evidence: `/tmp/apitools-m80-scale-rqhtf1mj/query-review-fixes{,-memory}.txt`.
The earlier measurements are retained, not overwritten or relabeled. Iteration
2 now reviews the complete M80 implementation, source-backed coverage, prompt
and work bounds, and all iteration-1 fixes. Publication remains unstarted.

#### Review iteration 2 passed — publication remains pending

Reviewed the full additive request/decoder, local retrieval and registration
identity checks, dimension ranking and scoped outcomes, guarded opt-in remote
lookup, bounded projection/ties, selected-artifact round trip, offline/default
safety, all source-backed fixtures, consumer compatibility, resource observations
and documentation. All six iteration-1 P2 findings remain fixed with regressions.
No open P1/P2-or-higher finding remains. Corrected stale introductory descriptions
that still called the delivered APIs pending; historical design baseline remains.
This is qualified implementation review, not publication or milestone acceptance.
M80.5 remains in progress until the exact outgoing diff checkpoint, authorized
fast-forward publication, observed full revision/module version, and closure.
Evolution v25 already records this approved direction; no boundary or target
change requires another snapshot.

#### Final qualified-source verification — before publication checkpoint

APItools full tests and vet pass in workspace and standalone modes after all
review fixes. Catalog generation is current; catalog quality reports zero
errors/warnings; unchanged search/import help and `git diff --check` pass.
OpenUdon full workspace, literal standalone published-pin, and actual APItools
worktree adoption through a disposable replacement modfile pass. Udon full
workspace and standalone current-source adoption through its disposable modfile
pass. Its existing standalone module-metadata gap is confined to that temporary
modfile; no consumer go.mod/go.sum or tracked source/ledger changed. Udon's
untracked `cmd/udon/tmp/` test output is preserved and excluded from all commits.
Final logs: `/tmp/apitools-m80-review2-{workspace,standalone}.log`,
`/tmp/openudon-m80-review2-{workspace,published-pin,adoption}.log`, and
`/tmp/udon-m80-review2-{workspace,adoption}.log`. All checks are provider/model-free.
No live remote lookup, sibling implementation, publication or launch occurred.

The qualification source commit is the prerequisite artifact for the required
publication task: it does not close M80.5. Its observed full revision will be
shown with the exact outgoing diff at the human checkpoint. Publication, module
identity, downstream handoff and retirement remain incomplete.

#### Publication observed; final closure review iteration 3 started

The user approved the exact outgoing APItools publication at the human
checkpoint. The reviewed patch SHA-256 was
`cfb073a380b7362b83259e785c0a613d3a0d1387bb09c5bfe3cbf96c0642cee0`,
covering 13 approved commits from remote baseline
`8580ff2485a3faff139b17bdc0b8e78d2af4ce18` to qualified source
`fb132631c9827eae5f2ec4503d03f21eabfb4113`. Before pushing, verified the
exact clean local source, unchanged remote baseline, patch digest and
fast-forward ancestry. Normal push succeeded; subsequent `git ls-remote`
confirmed `origin/main` equals the full qualified source revision.

Go's published-module query resolves
`github.com/OpenUdon/apitools v0.0.0-20260930205753-fb132631c982`
with `Origin.Hash` equal to that full source revision. Observed JSON is at
`/tmp/apitools-m80-published-module.json`; exact reviewed outgoing diff and
commit list are at `/tmp/apitools-stage5-publication-fb132631c982/`.
This publication includes M81; no tags, history rewrite or deployment occurred.
M80.5 is complete. Final review iteration 3 now checks the entire milestone
with publication evidence, consolidation, downstream reconciliation and
retirement; it does not reset the earlier implementation reviews.

#### Final review iteration 3 passed; consolidation and retirement

Rechecked the complete M80 scope and source at the published qualified revision,
all five task outcomes, prior P2 fixes/regressions, required verification and
benchmark observations, exact outgoing-diff approval, remote publication and
Go module resolution. No open P1/P2-or-higher finding remains. No source changed
after iteration 2's successful full producer/consumer verification.

The coordinator reconciled OpenUdon `status-M94.md` and `status-M95.md`, Kinet
`status-W10.md` and its disposable launch reference, and Udon `status-M45.md`
to the exact published APItools source/module contract. Those ledgers keep their
pending tasks and independently owned acceptance; no sibling code/dependency
pin or milestone ID changed. Udon's preexisting standalone module-metadata
qualification remains within M45.3's frozen-build scope, without adding an
upstream blocker. OpenUdon M94 retains M93 and qualifies its own new command.

Current product, architecture, technical docs and the durable coverage lesson
retain the delivered boundaries, evidence unknowns and query safety limits.
Evolution v25 remains the approved direction; no new snapshot is needed.
Retirement retains the exact final specification and full status in the portable
history envelope, anchored to the actual full published source and explicitly
identifying uncommitted closure changes. M81's frozen record is unchanged.
Final closure publication has its own required exact-diff checkpoint before
advancing to Udon M45; it is not authority to execute any later milestone.
````
