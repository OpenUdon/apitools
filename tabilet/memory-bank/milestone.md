# Milestones

Milestones are listed in priority order across long-lived domain lanes. Each
item lists scope and acceptance criteria. This file owns the roadmap, lane
meanings, active milestones, status-file index, dependencies, candidate
directions, milestone scope, and acceptance criteria.

Per-task completion state lives in one `status-<LANE><NN>.md` file per
milestone, such as `status-S02.md`. Individual status files
are the detailed task ledgers for their milestone; do not copy their task rows
into this roadmap.

Each milestone section is a review unit: after no `[ ]`, `[~]`, or `[!]` rows
remain in its matching status file, run a deep code review and a milestone
review against the milestone acceptance criteria before closing it. Completed
`[+]`, cancelled `[X]`, and closed-historical `[-]` rows are non-actionable;
every `[-]` row must name its accepted successor. Review-driven fixes should be
verified and committed before dependent work starts. An independent milestone
in another lane may proceed when the ownership and dependency rules below
explicitly make that safe.

The review procedure below covers one milestone. To run several in order, with
dependency and downstream reconciliation between them, [../GOAL.md](../GOAL.md)
is one optional protocol for that; any equivalent works just as well.

## Status ID Pattern

Status files are named `tabilet/memory-bank/status-<LANE><NN>.md`. `<LANE>` is one
uppercase domain letter and `<NN>` is a zero-padded number from `01` through
`99` within that lane.

Lane meanings:

- `M`: default and cross-cutting work. It also retains all M01-M70 history
  allocated before domain lanes were introduced.
- `C`: provider-catalog curation, source provenance, security overlays, and
  catalog quality.
- `S`: API source families, parsing, validation, local or remote discovery,
  import safety, inventories, and prompt-safe summaries.

Rules:

- Always use the two-digit form (`S01`, not `S1`). When a lane reaches `99`,
  open a new domain letter instead of adding a third digit.
- Do not reuse an ID after its status file exists, including after retirement.
  Search both active files and the history index before allocating an ID.
  Context archive IDs use an independent namespace; a matching archive ID is
  not a status collision and does not justify renaming a status ID. Cancelled
  work keeps its record and uses `[X]`.
- Retain a consumed failed attempt or superseded row as `[-]` closed historical
  evidence, record its accepted successor, and never retry it.
- Never rename a status ID after its file exists. Resolve a lane collision by
  allocating an unused lane/ID to new work and recording lineage; preserve the
  original record and its reserved ID.
- M1-M9 were legacy identifiers from the pre-lane harness. This M71 migration
  explicitly normalizes only their filenames and headings to M01-M09; the
  historical milestone identity is unchanged, and evolution snapshots are not
  rewritten. M10 and later IDs are unchanged.
- Do not rename completed IDs merely to reclassify history. New catalog work
  starts at C01, new source-tooling work starts at S01, and cross-cutting work
  continues from M72.
- Do not create an aggregate `tabilet/memory-bank/status.md`.

## Parallel Work And Priority

The status-file index is the priority view. A lane classifies ownership; it is
not a global execution phase. More than one lane may have an active milestone,
and independent rows may be implemented in parallel only when their milestone
sections record non-overlapping file ownership, resolved prerequisites, and
downstream impacts. Work that changes a shared public contract stays in `M` or
names the affected cross-lane dependency explicitly. Within one lane, prefer
one active implementation milestone at a time.

This is a deliberate local override of the single in-progress-row default:
at most one `[~]` row per lane, and parallel `[~]` rows only across lanes whose
milestone sections record that non-overlapping ownership. Without that record,
keep zero or one general row in progress across the active ledger. Parallel
rows still have one execution owner for the ledger; they do not authorize
concurrent ledger writers. A `tabilet/GOAL.md` run keeps its own single-row rule.

## Current Dashboard

M79 completed the Stage 1 operation-metadata remediation. Its accepted API
revision `e3625f6ef52ea54b7f78b7a4a4f1993bf8a06a46` is published and consumed
by OpenUdon M89. OpenUdon M89's published implementation at
`2e2ecedb32add53e10d6d81d4b2528ca2bdfdcc8` passed Kinet W04.6's exact consumer
check and is retired in published closure commit
`0c7c5d33da2ba7b190954b9eb402cc14b5ec1f73`. Two milestones are active, both
approved on 2026-09-30. M81.1 is complete under the confirmed Stage 5 goal.
Prerequisite M81 (catalog discovery foundations) runs first, then M80 (the catalog discovery
API). Latest completed milestone: M79. M78 and all earlier milestones through M77
are retired; see the history index below.
Apitools and Authoring are published, and OpenUdon passes standalone test/vet
against its pinned APItools revision.

`apitools` is a public OpenAPI tooling module and CLI. Its planning harness
(`AGENTS.md`, `tabilet/GOAL.md`, `tabilet/memory-bank/`, and
`tabilet/evolution/`) is tracked in this repository; the earlier private
`../tofu/apitools` symlinked harness is no longer used.

## Status Files

Each linked `status-<LANE><NN>.md` file records task rows, task notes, boundary
checks, and the review state for that milestone. Completed files stay here
until they are retired under
[Long-term memory and retirement](#long-term-memory-and-retirement).

Status markers, always written in backticks in the second column:

| Symbol | Status | Interpretation |
|---|---|---|
| `[ ]` | Pending | Item is actionable when its dependencies pass. |
| `[+]` | Completed | Item finished or done. |
| `[~]` | In Progress | Item is currently being worked, under the ownership rules above. |
| `[!]` | Blocked | A current unresolved blocker still prevents progress. |
| `[X]` | Cancelled | Item is no longer needed. |
| `[-]` | Closed Historical | A consumed failed attempt or superseded row retained for audit; it is never retried and does not block its accepted successor. |

Every active row below must link its existing status file. Retired IDs and
specifications belong in `tabilet/docs/history/index.md`, not in this table.
When history exists, add one link to that index here; do not accumulate one
retired row per milestone.

Retired milestones are indexed in [the history index](../docs/history/index.md).

| ID | Status file | State |
| --- | --- | --- |
| M81 | [status-M81.md](status-M81.md) | In progress; M81.1 design approved, remaining foundation tasks pending. |
| M80 | [status-M80.md](status-M80.md) | Pending; approved plan, waits for M81, implementation not started. |

Closed milestones are recorded in the history index.

## Active Milestone Specifications

## M81 — Catalog discovery foundations

Settle the six design decisions M80 depends on, and deliver the fixes for
pre-existing limits that would otherwise stop catalog-wide discovery: parser
and ranking budgets sized for explicit sources, a 20 MiB cap that excludes
large official specs, no defined artifact root for consumers, and
provider-level-only provisioning. M81 changes no existing exported API, JSON
shape, or `BuildOperationCandidates` behavior or ordering, and does not reopen
retired M77-M79 records. It adds shared public contracts (an operation index
builder and an artifact-scoped export), so it belongs in lane M. It runs
before M80; permanent IDs do not imply order, and the status index above is
the priority view. Implementation stays sequential under one execution owner,
with no parallel ownership or sibling writes authorized.

Provenance: review "APItools M80 — Catalog discovery API for step contracts"
(`apitools-m80-review.md`, 2026-09-30), findings P2-1 through P2-7, L1, and
L2; source and local severities are recorded in M81's status file. Review
baseline `8580ff2`; revalidated at
`8580ff2485a3faff139b17bdc0b8e78d2af4ce18`, with the uncommitted M80 planning
files as part of the evidence and no uncommitted code. The user approved the
dispositions, the six decisions below, this milestone, and the M80 amendments
on 2026-09-30. This intake starts no review counter.

**Approved design decisions.**

1. **Retrieval.** An explicit refresh-time command builds a digest-bound,
   versioned operation index; discovery reads only the index.
2. **Large artifacts.** Only the index builder may read catalog-registered,
   digest-verified artifacts above 20 MiB, up to the existing 128 MiB artifact
   bound, under separately reviewed parser budgets. Every existing parser
   entry point and `BuildOperationCandidates` keep the 20 MiB default.
3. **Artifact roots.** The caller always supplies the catalog root. With no
   root, discovery returns metadata-only provider leads and insufficient
   evidence that names what would resolve it. Embedding the committed
   overlays remains a candidate direction.
4. **Provider constraints.** An optional request field resolved by exact key
   through the catalog's `FindProvider` lookup, which already matches
   multi-word names such as "Microsoft Outlook". Each result states whether
   it came from a constraint or from open retrieval.
5. **Outcomes.** Match, ambiguous, no qualifying API within the checked scope,
   insufficient evidence, and blocked. Only a scoped no-match may route a
   consumer to a browser step. Curated catalog facts, such as an unavailable
   machine spec or a docs-only provider, are evidence, never proof.
6. **Provisioning.** An additive artifact-scoped export/materialize API, added
   as a new function or options type so existing exported struct shapes and
   unkeyed composite literals stay unchanged.

**Scope and acceptance.**

- **Design record and checkpoint.** `docs/catalog-discovery.md` records the
  request with provider constraints; the outcome set; the qualification rule
  that separates a match from a weak or ambiguous result; the
  outcome-to-consumer-action table; license/redistribution defaults (unknowns
  included and labeled, `license_note` passed through verbatim, no inferred
  permission); the catalog-stable provider/artifact reference and its round
  trip into the M81 export API; the root contract; and the index approach. It
  lists the changes the proposed OpenUdon M94 and Kinet W10 drafts must absorb
  at their promotion. The user approves the record before M80's contract is
  implemented; APItools edits no sibling file.
- **Catalog root contract.** A caller-supplied root option naming the cache
  directory, the artifact registrations, and the index location; documented
  no-root behavior; operator documentation for preparing a root (refresh,
  then index); and a committed synthetic, redistributable fixture root that
  contains no third-party provider specs.
- **Large registered artifacts.** A reviewed limit up to 128 MiB for
  digest-verified registered artifacts on the index path only, with explicit
  time, memory, and structural budgets. Record measurements from an opt-in,
  provider-free local check against the two known oversized cached specs
  (`microsoft-graph-v1-openapi`, 35.4 MiB; `cloudflare-api-openapi`,
  20.9 MiB); CI uses synthetic large fixtures. Artifacts over the limit or a
  budget fail closed as unexamined scope.
- **Operation index.** A `catalog index` command and library builder that
  digest-verify each registered artifact and store sanitized, source-backed
  operation metadata: native selector, effect assessment, input/output value
  summaries, security alternatives, and capability gaps. Entries use
  catalog-stable identity (provider IDs, spec-ref ID, artifact ID, SHA-256,
  native selector) and no absolute paths; an artifact shared by several
  providers is indexed once with every provider link. The index records the
  built-in catalog identity and per-artifact coverage: indexed, missing
  registration or file, digest mismatch, oversize, parse failure, or
  unsupported family. Rebuilding from the same inputs is byte-identical, and
  relocating the root does not change the index. Readers verify the index
  version and each entry's registration digest and report stale entries as
  unexamined. The index lives under the supplied root, is never committed,
  needs no network, and must not break existing cache readers or require
  migrating existing caches.
- **Artifact-scoped provisioning.** Select one or more artifacts by reference,
  verify their expected digests, copy only them with their provider- or
  spec-scoped security overlays and provenance, fail closed on a mismatch, and
  handle shared artifacts. Existing provider-level export and materialization
  behavior is unchanged.
- **Documentation and compatibility.** README and architecture describe the
  delivered behavior, and the verification below passes: `go test ./...`,
  `GOWORK=off go test ./...`, `go vet ./...`, `GOWORK=off go vet ./...`,
  `go run ./cmd/cataloggen -check`, `go run ./cmd/apitools catalog check`,
  `go run ./cmd/apitools catalog index --help`, `git diff --check`, and
  OpenUdon/Udon compatibility checks in workspace and standalone modes per
  tech-stack.md, with Ramen excluded. Default checks stay provider-free and
  offline.

**Dependencies and downstream.** No upstream prerequisite: retired M77-M79
supply the operation metadata, and the existing catalog refresh registrations
and `internal/artifactio` reads supply digest-verified artifacts. M80 depends
on M81's completed rows and closed review gate. M81 needs no separate
publication; M80's authorized published revision carries it to the proposed
OpenUdon M94. Kinet's Stage 5 draft plan still lists only APItools M80;
updating it is outside this repository.

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
[M81](#m81--catalog-discovery-foundations) and amended the scope and rows
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
- Every result carries catalog-stable provider/artifact identity and a
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

## Requested Changes After Initialization

For a requested feature, candidate promotion, or change to future direction,
inspect the current plan and implementation first. Record user intent separately
from observed facts and unresolved assumptions. If a pending milestone already
owns a small change, propose the smallest row and acceptance update that fits.
Otherwise propose a dependency-closed milestone with an unused permanent ID.
A candidate's trigger prompts a new decision and approval; promotion is never
automatic. Schedule by approved priority and dependencies, not review finding
severity. Reconcile affected downstream acceptance and optional launch input.

Present one complete proposal with intended outcome, owner rows, acceptance,
verification, dependencies, downstream effects, and exact file actions. Write
planning changes only after approval and a fresh check of affected files,
worktree changes, and active and retired IDs. Preserve non-pending outcomes,
review counters, local policy, and frozen history. Record target behavior here
and in pending status rows until implemented; current architecture describes
only what the repository establishes as fact. Planning does not implement or
accept the requested behavior. Use `memory-bank-propose` where installed, or
follow this procedure directly with another agent.

## Candidate Directions

Candidate directions are outside the active execution horizon. They have no
lane, permanent ID, status file, or execution-order entry. A trigger starts a
fresh scope/dependency reconciliation and approval round; it does not schedule
the work automatically.

| Direction | Why Deferred | Promotion Trigger |
|---|---|---|
| Next provider-catalog expansion | No provider batch or user-verifiable outcome has been approved after the completed historical catalog batches. The [catalog upgrade proposal](../../docs/catalog-upgrade-2026-09.md) also proposes a structured license/redistribution audit and workflow-ready typed overlays; those remain separate from M80's discovery report, which must expose unknown evidence. | Approve a bounded provider-owned source batch, its source/auth/license review rules, or a bounded license/typed-overlay upgrade with acceptance evidence. Proposed catalog IDs remain unallocated until fresh reconciliation and approval. |
| Embedded baseline catalog root | M81's approved decision 3 keeps catalog roots caller-supplied. Embedding the 141 committed advisory overlays (about 2 MB) and their index would give published consumers a zero-setup offline root, but it widens distribution before the committed-overlay redistribution audit proposed in the [catalog upgrade proposal](../../docs/catalog-upgrade-2026-09.md) and adds binary size for every consumer. | The committed-overlay redistribution audit completes, and a consumer needs a zero-setup offline root. |
| Production adoption of experimental LAP/RFC 9727 discovery | The adapters are intentionally experimental, and OpenUdon still selects APIs.guru explicitly. | Collect reliability/coverage evidence and approve the downstream OpenUdon discovery policy and network budget. |
| Remove deprecated Discovery and Smithy wrappers | Sibling consumers may still rely on compatibility imports, and removal is a public breaking change. | Confirm no active sibling imports remain, approve release notes, and pass downstream compatibility checks. |
| Network and local-scan hardening batch | Review "apitools Code Review" Pass 1 (baseline `a5699e5`) findings #7 (LAP registry validates every entry by DNS before scoring/limit and fails the whole search on one bad `source_url`, `remote_discovery.go`), #9 (unsafe-IP list omits `240.0.0.0/4`, IPv4-compatible `::/96`, and `fec0::/10`, `download.go`), #11 (OpenAPI metadata paths skip the `sourceguard.CheckJSON`/`CheckYAML` preflight, `validation.go`, `local.go`, `local_source_discovery.go`), #12 (`.bin` is always treated as a protobuf descriptor, `local_source_discovery.go`), #13 (regular-file-to-FIFO swap between `Lstat` and `Open` can block, `local_read.go`, `internal/artifactio`), and #14 (`Discoverer` exposes no port/transport options, `discovery.go`). All were revalidated at `a5699e570fb5d6288fdd255c71d229fb062b50b3` as Lower severity; no supported scenario fails and no private-host or root-escape path was found. | Approve a hardening batch after S02, or promote any item earlier if it becomes reachable in a supported deployment or a downstream consumer depends on it. |
| Egress proxy support for guarded downloads | Review "apitools Code Review" Pass 1 finding #8: the guarded transport has `Proxy: nil`, and the only alternative, `AllowUnsafeHosts`, drops all host checks (`download.go`). Supporting a proxy trades dial-time IP filtering for proxy trust, which is a product/security decision. | An operator needs mandatory-proxy egress, and an approved policy defines target pre-validation, proxy trust, and the documented DNS-rebinding limitation. |
| Lifecycle and scan cleanup follow-ups | Uncommitted-diff review (2026-09-26) Lower items not required by S02, S03, or S04 acceptance: `operationlifecycle` recomputes tokens and purpose for every operation on each role (performance only); offline `search` still applies the cache TTL while offline `import` ignores it (consistency decision); `isUnicodeTagCharacter` largely repeats the Cf category check in `prompt_safety.go` (cosmetic). | Promote when lifecycle ranking becomes a measured hotspot, an operator needs offline search of expired entries, or `prompt_safety.go` is next edited. |
| Catalog and tooling hardening batch | Review "apitools Code Review" Passes 2 and 4 (revalidated at `e015231`), all Lower severity: C3 `ResolveProviders(query)` splits multi-word display names on whitespace (`catalog/materialize.go`; CLI and OpenUdon pass explicit keys); C4 `FindSpecReference` trims the provider ID for identity but not for the spec lookup (`catalog/index.go`); C5 `preferredSpecReference` ignores `asyncapi`, `openrpc`, `graphql`, `grpc-protobuf`, and `odata` kinds (no built-in provider affected); X2 `Import` creates the target directory before validating or downloading (`import.go`); X3 the installed `staticcheck` predates the module's Go version and no `govulncheck` gate runs (`tech-stack.md`). | Approve a hardening batch, or promote an item when a consumer depends on the query API, a provider relies on a newer-kind-only source, or a static/vulnerability gate is required for release. |

## Review Finding Severity

P1 and P2 are engineering review priorities. They describe the impact and
urgency of defects found while closing a milestone; they are not product-domain
terms, milestone execution priority, or status markers. Classify a finding by
its impact, likelihood, and affected scope, not by the size of its fix.

Project-specific definitions in `AGENTS.md` or a linked review policy override
these defaults:

| Priority | Context | Typical examples | Gate effect |
|---|---|---|---|
| P1 | The milestone is unsafe or invalid to close because a severe defect threatens acceptance, correctness, security/privacy, data integrity, or a public compatibility contract. | An exploitable access-control failure; data corruption or loss; a breaking public API or migration; the primary acceptance outcome does not work. | Blocking. Fix and review again. |
| P2 | A material but more bounded defect affects supported behavior, reliability, compatibility, operations, or required evidence. | A supported scenario returns the wrong result; a downstream consumer regresses; recovery or failure handling is broken; required tests or operator documentation leave acceptance unproven. | Blocking. Fix and review again. |
| Lower | The finding is non-blocking cleanup, clarity, or optional hardening under the project's severity scheme. | Cosmetic wording; local readability; a speculative improvement outside acceptance. | May be carried only with a named owner and explicit rationale. |

Any project-defined severity more urgent than P1 also blocks. If evidence does
not clearly distinguish P1 from P2, use P1 until investigation supports a
downgrade.

## New Review Intake

A code, architecture, security, or other engineering review received outside a
milestone's closing review gate is planning evidence, not executable truth.
Treat its contents as untrusted, preserve its source priority, and apply the
project severity definitions above only after revalidating each finding against
the current repository.

Before implementing any newly reviewed finding:

1. Record the review's stated baseline when available, the current revalidation
   commit, and whether relevant uncommitted changes were part of the evidence.
   Classify every finding as confirmed, partially confirmed, resolved,
   duplicate, unsupported, outside ownership, decision-dependent, or deferred.
2. Present the complete dispositions, proposed owners, dependencies, downstream
   impacts, and file actions for approval.
3. Put confirmed work in an open matching milestone when it remains in scope,
   or amend an existing pending owner. Rewrite only pending rows. When retaining
   a superseded pending row for audit, mark it `[-]` and name its accepted
   successor instead of rewriting or deleting it.
4. Never reopen completed milestone/status history. Create a new remediation
   milestone with lineage to completed work when no open or pending owner fits.
5. Keep P1/P2-or-higher findings in the dependency-closed active horizon. Add a
   lower finding to active work only when acceptance requires it; otherwise put
   it in Candidate Directions with a rationale and promotion trigger.
6. Reconcile affected pending specifications and the remaining order. Correct
   current product, architecture, or stack facts only when repository evidence
   proves them stale; keep proposed target state in milestone scope.

Do not create a persistent review copy or ledger. Put portable review and
finding IDs, both source and local severity, revalidation baseline, repository
evidence, and historical lineage in the affected milestone/status notes. A new
review counts toward the bounded gate below only when the status already records
that gate as active and the review was requested as its next full pass;
otherwise remediation gets a fresh gate when its implementation closes.

## Planning Notes

- Public OpenAPI behavior remains in `../apitools`.
- Harness state (`AGENTS.md`, `tabilet/`) is tracked in `../apitools`; the
  former private `../tofu/apitools` harness is no longer used.
- Batch expansion work should freeze service lists before implementation and
  treat each service row as a commit unit once curation begins.
- Do not promote non-OpenAPI converter work until most provider-node catalog
  coverage is implemented.
- Google Discovery converter ownership is now explicitly opened in M38 because
  provider-node catalog coverage through M37 is complete.
- AWS Smithy converter ownership is complete in M39 for official AWS Smithy
  JSON service models; SigV4 signing remains downstream runtime work.
- AWS Smithy native model ownership is complete in M40; OpenAPI-shaped Smithy
  conversion has been removed from `apitools`, and downstream runtimes consume
  native Smithy metadata.
- Google Discovery native model ownership is complete in M41; OpenAPI-shaped
  Discovery conversion has been removed from `apitools`, and downstream
  runtimes consume native Discovery metadata.
- Run `apitools catalog check` before completing each catalog metadata batch.
- M52 supersedes the prior n8n-derived catalog planning workflow. Active
  catalog expansion is provider-source-first and must not copy third-party
  integration code, generated node metadata, credential behavior, descriptions,
  icons, or runtime workflow behavior.
- First-class Smithy and Discovery work means cataloged native source metadata
  and parser coverage, not reintroducing Smithy-to-OpenAPI or
  Discovery-to-OpenAPI conversion as an `apitools` runtime contract.
- M47 retained Dropbox Stone as official source metadata and advisory-overlay
  provenance, but did not open native Stone parser work or expand the M37
  advisory overlay.
- M48 closes the residual Google Discovery batch for n8n-visible Google
  services that were not promoted in M44, and records Google Ads as a
  candidate-only/no-Discovery decision until a future protobuf/gRPC source
  family is explicitly scoped.
- M49-M51 are historical and superseded by M52 for active code and public docs.
  The n8n provider row, n8n gap-report CLI/API, and n8n-derived priority
  evidence are intentionally removed before release.
- M53 should add Docker source coverage from Docker-owned API references:
  Docker Hub as a normal public OpenAPI provider, Docker Engine as a local
  daemon/control-plane metadata source, and Docker Registry only if source
  review confirms a durable official spec path.
- M54 treats Kubernetes as cluster-published source metadata. A pinned
  v1.19.2 Swagger fixture snapshot is available as a non-official bootstrap
  artifact, but real Kubernetes OpenAPI and Discovery documents remain
  user-provided or cluster-exported artifacts; `apitools` must not read
  kubeconfig, resolve cluster credentials, contact API servers, or claim
  live-cluster compatibility.
- M55 should add enterprise application source coverage for NetSuite, SAP
  S/4HANA, SAP SuccessFactors, Oracle Fusion Cloud Applications, and Workday.
  Treat account/tenant-specific metadata endpoints as user-provided source
  URLs, and keep OData, WSDL/SOAP, and describe-metadata decisions explicit.
- M56 should add an OpenAPI-first provider batch from provider-owned official
  specs: Adyen, DocuSign, Auth0, Confluence, Acumatica, BigCommerce, Cisco
  Meraki, and Confluent Cloud. Use third-party integration catalogs only as
  market-demand signals; do not copy their connector metadata.
- M57 should add a docs-overlay provider batch for Marketo, Aircall, Checkr,
  Airwallex, Apollo.io, AfterShip, and Adobe Acrobat Sign only where official
  downloadable OpenAPI evidence is unavailable or unsuitable and official docs
  can support a reviewed subset.
- M58 is deferred until M57 is complete. It should audit all durable catalog
  providers for auth/security metadata completeness under the OpenAPI/Swagger
  security-scheme and root/operation security requirement rule, plus native
  Smithy, Discovery, Stone, human-docs, and tenant-specific source-family
  decisions.
- M59 should use `tray_workato_connectors.xlsx` only as a market-demand
  signal, then add Canva, LaunchDarkly, Klaviyo, Attio, Fivetran, Anthropic,
  Amplitude, Ashby, Braze, and Postman from provider-owned API evidence.
- M60 should expand AWS Smithy service coverage for Workato-visible AWS
  candidates using official AWS Smithy JSON models only; SigV4 remains
  advisory metadata and no AWS credential, region, or account behavior belongs
  in `apitools`.
- Deprecated `apitools/googlediscovery` and `apitools/awssmithy` wrappers can
  be removed only after active sibling consumers import
  `github.com/OpenUdon/googlediscovery` and `github.com/OpenUdon/awssmithy`
  directly, public release notes call out the compatibility break, and
  dependent checks in `../openudon` and `../udon` pass without wrapper imports.
- The next catalog-expansion milestone should start from provider-owned source
  evidence and prioritize remaining high-demand services with official
  OpenAPI/Swagger, Smithy, Discovery, Stone, official indexes, or official docs
  sufficient for reviewed advisory overlays.
- Post-M60 materialization work moves static provider artifact handoff into
  `apitools`: resolving provider names, reporting protocol capability, copying
  existing cache-registered artifacts, emitting separate security-overlay JSON,
  and exporting provenance manifests for workflow directories. It remains
  offline and must not download missing artifacts, apply overlays into OpenAPI
  documents, or reintroduce Discovery/Smithy-to-OpenAPI lowering.
- M61 reviews existing docs-derived advisory overlays under the expanded
  linked-docs curation rule: follow official API links cited by source docs,
  include related operations when they support the same provider workflow, and
  record no-add decisions where linked docs are conceptual, commercial,
  duplicated, authenticated-only, or outside the reviewed scope.
- M62 opens a narrow public helper surface for pure fnct payload shaping. The
  first helper is Gmail raw-message rendering; runtime registration remains in
  `../udon`, and OpenUdon remains the authoring/review consumer.
- M63 adds OpenRPC as metadata-only JSON-RPC source support: the `openrpc`
  package parses document metadata, method summaries, selector aliases, params,
  results, errors, and components; catalog/source materialization treats
  OpenRPC as its own `openrpc` source family; JSON-RPC execution, server
  contact, credentials, and runtime transport remain out of scope.
- M64 closes Postman Collection, RAML, and API Blueprint as local artifact
  detection/classification and conversion candidates only, following the
  Dropbox Stone decision pattern. Preserve provenance, prefer official
  OpenAPI/Swagger, Discovery, Smithy, AsyncAPI, OpenRPC, or another stronger
  native source when available, and do not add parser code, runtime behavior, or
  UWS `sourceDescriptions[].type` values for these formats.
- UWS 1.4 now defines `graphql`, `openrpc`, `grpc-protobuf`, and `odata` source
  description types. OpenRPC parser/classification work is already complete in
  M63. M65-M67 open separate metadata-only source tooling milestones for
  GraphQL, gRPC/protobuf, and OData. Each must preserve native source semantics
  and must not execute operations, probe live services, resolve credentials, or
  choose accounts/tenants/projects.
- M65 adds GraphQL as metadata-only source tooling: local SDL, introspection
  JSON, and operation documents parse into native GraphQL summaries and
  selectors; catalog refresh, resolution, and materialization identify
  GraphQL as `graphql` without live introspection or REST-shaped lowering.
- M66 adds gRPC/protobuf as metadata-only source tooling: local `.proto`, JSON
  descriptor, and binary `FileDescriptorSet` artifacts parse into native
  protobuf service/method/message summaries and selectors; catalog refresh,
  resolution, and materialization identify gRPC/protobuf as `grpc-protobuf`
  without server reflection, channel behavior, TLS negotiation, RPC execution,
  or REST-shaped lowering.
- M67 adds OData as metadata-only source tooling: local CSDL XML, JSON CSDL,
  and exported service metadata parse into native entity set, singleton,
  navigation, action, function, parameter, query-option, and selector
  summaries; catalog refresh, resolution, and materialization identify OData
  as `odata` without tenant metadata calls, login, credential resolution,
  account/tenant selection, `$batch`, request execution, or REST-shaped
  lowering.
- M68 adds a reviewed local artifact corpus for UWS 1.4 source families:
  exactly two official example artifacts each for GraphQL, OpenRPC,
  gRPC/protobuf, and OData under source-aligned cache directories. It is
  curation and validation coverage only; it does not add provider catalog rows,
  source types, OpenAPI lowering, parser behavior expansion, or runtime
  behavior.
- M70 keeps APIs.guru as the primary global lookup, evaluates LAP Registry as
  an experimental second source, and adds provider-scoped RFC 9727 discovery.
  Swagger Catalog and Scalar Registry remain outside implicit global discovery
  because they are publisher/team registries rather than unauthenticated
  global aggregators.

## Milestone Review Procedure

When the last open row in a milestone's status file closes as `[+]`, `[X]`, or
`[-]` during an agent session, perform the review before ending the turn and
before moving to the next milestone. A `[-]` row counts as closed only when its
notes identify the consumed attempt or supersession and its accepted successor:

1. Re-read the milestone scope and acceptance criteria here. Confirm the code or
   docs meet the acceptance line; do not rely on the status file alone.
2. Run the required verification before review, including proportionate checks
   for affected consumers and any applicable compatibility, migration,
   rollback, security, concurrency, or failure paths.
3. Run a bounded deep-review and fix gate using the review finding severity
   context above.
   - The initial deep-review pass is iteration 1. Read the persisted count and
     findings first. Record the iteration as started before review; resume an
     interrupted pass at that same number. Runtime round limits are separate.
     Read the `git log` range and
     review the full milestone diff for correctness, regressions, failure
     semantics, boundary drift, stale docs, and missing tests.
   - Record the iteration number and findings in the current status notes so a
     continuation cannot reset the counter. If the pass finds no P1, P2, or
     higher-severity issue, the gate passes. Otherwise, when the current
     iteration is below 10, fix every such finding in the current milestone,
     rerun affected verification, and review the whole milestone again,
     including the fixes.
   - Run at most 10 iterations; a session or reviewer change does not reset the
     counter. If iteration 10 still finds a blocking issue, add `[!]` review
     rows recording the remaining findings and iteration-limit blocker. Do not
     start another automatic fix-review cycle or move downstream; stop and ask
     for user direction.
   - Carry a lower-severity finding forward only with a named pending owner and
     explicit rationale.
4. Reconcile the memory bank. Update `product.md` if the milestone changed
   product scope, domain terminology, concept relationships, or business
   invariants. Update `architecture.md` or `tech-stack.md` if it changed
   boundaries, dependencies, commands, data flow, or runtime assumptions.
5. Check `tabilet/evolution/`. Record an explicit bump or no-bump decision. Add
   the next `prompt-vN.md` and `result-vN.md` only when product direction,
   architecture boundary, milestone target, or public/private contract direction
   materially changes.
6. Revisit candidate directions affected by the milestone. Update their reason
   or trigger; when a trigger is now true, propose a reconciled milestone and
   obtain approval before allocating its permanent ID and status file.
7. After the gate passes, reconcile affected active downstream dependencies and
   acceptance against the delivered outcome. Then consolidate and retire under
   the procedure below. During an ordered goal, the goal protocol finishes its
   downstream reconciliation before performing this retirement step.
8. Run required verification again, including the retirement links and record,
   then commit substantive closure changes under the governing commit policy.
   Do not create an extra milestone commit when the review changes nothing.
9. Report a short review summary: what was verified, the review-fix iteration
   count, what memory-bank files changed, any review commit, and whether an
   evolution bump was made. Include the retired record and durable lessons.

## Long-term Memory And Retirement

This project retires a milestone automatically after closure. “Automatically”
means the agent performs this procedure during the milestone review workflow
above, including when using `memory-bank-next` or `memory-bank-goal`. There is
no background process or separate archive invocation. Individual completed,
cancelled, or closed-historical rows remain in their active status file until
the whole milestone qualifies; unresolved work or missing closure evidence
keeps it active. Terminal rows alone do not prove acceptance.

Existing projects adopt this contract explicitly, with a compatible API runner
if used; upgrading installed skills alone does not merge project instructions
or move files. An explicit cleanup request may retire older closed milestones
only when their closure evidence is available, either a recorded review under
this contract or the legacy closure described below.

### Consolidate Before Retiring

Keep current facts in product, architecture, and stack documents. Maintain
[lessons.md](lessons.md) for applicable lessons and decision rationale with
evidence links; merge duplicates and avoid a chronological session log. Maintain
relevant lessons during ordinary work, not only at closure, and keep them active
while applicable even after their supporting milestone retires. This removes
accumulated history from active context, not a fixed number of tokens or files;
genuinely active work and useful knowledge can still grow.

Before materially superseding or removing knowledge from those documents,
append the old wording to `tabilet/docs/history/knowledge.md`. Use a unique, descriptive
dated heading, the original document and heading, retirement reason, supporting
evidence, and a link to its replacement (or an explicit reason there is none).
Preserve the old excerpt in a fenced markdown block. This applies outside
milestones too. Create this journal only when needed, link it from the history
index, and append corrections rather than rewriting old entries. Routine edits
need no journal entry; Git, when present, holds intermediate revisions.

Consolidation is ordinary project maintenance. A separate context snapshot is
optional and never a prerequisite for retirement; the archive skill's clean
baseline rule must not force commits or interrupt a no-commit goal.

### Retire A Closed Milestone

1. Require a passed review within the persisted 10-iteration limit, recorded
   verification, consolidated knowledge, and reconciled downstream work. No
   `[ ]`, `[~]`, or `[!]` row may remain. Every `[-]` row names its accepted
   successor. A cancelled or superseded milestone also needs an authorized
   disposition; it must not be described as delivered acceptance.
2. Create `tabilet/docs/history/status-<LANE><NN>.md` using the envelope below. Keep the
   full final specification and status document in separate literal markdown
   fences, including all task text, notes, acceptance, and review evidence.
   Choose fences longer than any fence in the source. Original relative paths
   inside these literal documents retain their source-document meaning.
3. Add one history index row using `Milestone | Outcome | Retired | Record |
   Summary`; the ID, outcome, date, and record link must match the envelope.
   Outcomes are `completed`, `cancelled`, or `superseded`. The record link is
   relative to the index, for example `[M01](status-M01.md)`.
4. Remove the active status file, specification section, and active index row.
   Before removal, validate every envelope field and compare both retained
   documents with their complete active sources. For Git evidence, use the
   literal full output of `git rev-parse --verify HEAD`, never an abbreviated
   hash from a log display. Stop with the active sources intact if validation
   fails; a malformed retired record is not completed closure.
   Keep one link from this file to `../docs/history/index.md`. Repair maintained
   incoming links and remaining dependencies. Never rewrite frozen snapshots;
   resolve their old paths through the record's original-path metadata.
5. Refresh an existing disposable goal suggestion to contain only remaining
   active work, or remove it when empty. Do not create one without a compatible
   goal protocol. The resolved goal in the conversation remains authoritative
   over its disposable launch input.
6. Verify the record and links before handoff. On interruption, reconcile the
   source and destination before continuing; duplicate IDs or a missing record
   are invalid, and an existing retired record must never be overwritten.

Use this envelope, with actual values. The `Evidence` field is the observed
full Git commit, not a claim that uncommitted work exists in that commit.
`Worktree` is `clean`, `includes uncommitted changes`, or `unversioned`; without
Git both provenance fields say `unversioned`. The date uses `YYYY-MM-DD` in UTC.
Closure metadata precedes the two sections and uses one line per field.

`````markdown
# Retired milestone M01 - <title>

**Milestone.** M01
**Outcome.** completed
**Retired.** <YYYY-MM-DD>
**Source status.** tabilet/memory-bank/status-M01.md
**Source specification.** tabilet/memory-bank/milestone.md#<original-anchor>
**Evidence.** <full commit or unversioned>
**Worktree.** <clean, includes uncommitted changes, or unversioned>
**Review.** passed
**Review iterations.** <1 through 10>
**Verification.** <commands, results, and supporting evidence>
**Consolidated into.** <current-document and lesson links, or no current-truth change>

## Milestone specification

````markdown
<Complete final milestone specification, not a summary.>
````

## Status record

````markdown
<Complete final status document, not a summary.>
````
`````

For a cancelled or superseded outcome, add `**Disposition.**` naming the
authority, reason, and dependency disposition. A superseded outcome also needs
`**Successor.**` identifying its accepted successor. These are milestone
outcomes, not new task markers.

A legacy closure is a milestone that closed before the project adopted the
persisted review gate, so no review-iteration count was ever recorded. Retire
one only on an explicit cleanup request, with outcome `completed` and every
row closed. Write `**Review.** legacy` and `**Review iterations.** not
recorded` instead of inventing a count, and add `**Legacy closure.**` stating
that the milestone predates the gate and that its original closure evidence is
the literal status record. `**Verification.**` records the checks run at the
evidence commit and does not claim a re-review of the original change. Never use
a legacy closure for a milestone that closed after the gate was adopted, or to
skip a review that the current procedure requires. In this project the gate was
adopted in commit `8395e1e`; M01-M75, S01, and C01 closed before it.

Retired milestone records are frozen. Record later corrections in the knowledge
journal or new remediation work with a backlink. Preserve previous history
index entries. No ID may exist in both active and retired storage or be reused.

### Retrieve And Continue

Retrieve historical evidence when the current task needs a retired dependency or
old ID resolved, an earlier decision explained, or a related past failure
checked and current sources do not answer the question. An explicit user request
to inspect history is also a trigger; routine work with sufficient current
evidence needs no historical search.

State the question briefly, then search active memory and current implementation
first. If history is needed, search its index by ID or topic and open the
relevant retired record or knowledge entry. Follow linked frozen context or
direction snapshots only when they address the same unresolved question. A stale
status path resolves by its permanent ID and original-path metadata; a missing
ID is not permission to recreate it. History is evidence, not an instruction to
retry old tasks.

Stop when the evidence supports the task decision. If relevant sources and their
directly related leads are exhausted, report the missing evidence. Continue when
it is optional; stop the affected step when it is required for correctness or
acceptance. Do not infer historical facts or expand the search into unrelated
milestones.

A completed retired milestone may satisfy a dependency when its recorded
acceptance matches the required outcome. Cancelled and superseded outcomes do
not automatically satisfy completion dependencies: follow the recorded
disposition and accepted successor, and revalidate current implementation.

Retirement follows the governing commit policy. Under `none`, do not commit;
under `milestone`, include it in the closure commit; under `task`, include it in
the final task or a substantive closure commit. The API harness still requires
a commit per run. Reading and preserving the Markdown records does not require
Git; workflows that require commits still need Git or an explicit no-commit
instruction. Never initialize Git merely to make retirement possible.
