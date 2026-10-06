# Milestones

## Stage 11 active horizon

Approved 2026-10-06: both phases of [Kinet STG-11](../../../kinet/docs/stage11.md), with one serial execution owner across Kinet, UWS, APItools, Udon and OpenUdon. This package owns M82; all rows are pending and each review is 0/10. [Specifications](#stage-11-cross-package-refactoring) below are the current horizon. Earlier completed horizons and records remain historical; planning grants no implementation or external authority.

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
`0c7c5d33da2ba7b190954b9eb402cc14b5ec1f73`. M81 and M80 are complete
and retired. M81 foundations qualified at
`8d67aef2dce565aa9ac8e8b2c56a9100dcf9a5a6`; M80 discovery qualified and
published at `fb132631c9827eae5f2ec4503d03f21eabfb4113`
(`v0.0.0-20260930205753-fb132631c982`), carrying M81 publication. The latest
completed milestone is S05. The confirmed S06 -> S05 goal completed on the
isolated `work/s06-s05` branch: S06 passed review in one iteration and S05
in two, with required APItools and candidate-bound OpenUdon/Udon verification.
That earlier horizon is complete; Stage 11 M82 is now planned and unimplemented. Merge/publication and consumer
adoption are separately authorized; the main workspace and consumer pins are
unchanged. Candidate directions remain unnumbered and require fresh approval.
Exact history and observed publication evidence resolve through the history index.
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
| M82 | [status-M82.md](status-M82.md) | in progress; review 0/10 |

Closed milestones are recorded in the history index.

## Active Milestone Specifications

Stage 11 M82 is the active planned milestone; its specification is below.
Completed source repair and version discovery remain in the history index.
Other candidates still require fresh reconciliation and approval.

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
| Next provider-catalog expansion | No provider batch or user-verifiable outcome has been approved after the completed historical catalog batches. The [catalog upgrade proposal](../../docs/catalog-upgrade-2026-09.md) also proposes a structured license/redistribution audit and workflow-ready typed overlays; those remain separate from M80's discovery report, which must expose unknown evidence. GoDaddy is absent from the catalog; its official Domains specs (`https://developer.godaddy.com/openapi/domains-v1.json` through `domains-v3.json`) were observed on 2026-10-01, and curating v1 or v2 would need a source-reviewed refresh correction for a publisher `properties.type` string defect. | Approve a bounded provider-owned source batch, its source/auth/license review rules, or a bounded license/typed-overlay upgrade with acceptance evidence. Proposed catalog IDs remain unallocated until fresh reconciliation and approval. |
| Embedded baseline catalog root | M81's approved decision 3 keeps catalog roots caller-supplied. Embedding the 141 committed advisory overlays (about 2 MB) and their index would give published consumers a zero-setup offline root, but it widens distribution before the committed-overlay redistribution audit proposed in the [catalog upgrade proposal](../../docs/catalog-upgrade-2026-09.md) and adds binary size for every consumer. | The committed-overlay redistribution audit completes, and a consumer needs a zero-setup offline root. |
| Search default-behavior refinements | `Client.Search` keeps APIs.guru's `preferred` version and falls back to a lexicographic sort of version keys when it is missing (`providers.go`), and `auto` stops at the first source with matches (`client.go`). Changing either default alters results that OpenUdon and other callers rely on, so S05 reports versions through a separate opt-in command instead. Separately, the retired M80.4 remote tier sets `client.Cache = nil` with `CacheModeBypass` (`catalog_discovery_remote.go`), so each remote lookup re-downloads the whole APIs.guru directory list; S05 avoids that cost for its own checks but does not change the retired behavior. | Approve a default-behavior change with consumer reconciliation, observe APIs.guru entries whose `preferred` version is missing or unusable, or measure the directory download as a latency problem for a consumer. |
| Identify APItools on guarded downloads | No `User-Agent` is set on any guarded request, so hosts see Go's default agent. S05's new probe primitive sends a stable tool identifier, but adding one to the shared download path changes the requests every existing caller makes, including OpenUdon's pinned use. | A publisher rejects the default agent, an operator requires identified egress, S05's probe evidence shows the identifier is worth sharing, or the provider version directory's crawl policy requires it; then approve it with consumer reconciliation. |
| Provider version directory | A shared, published directory that grows with demand: one record per provider and service with its locator, recipe, resolved URL, digest, `checked_at`, evidence label, and the age of any APIs.guru baseline, in an APIs.guru-`list.json`-compatible shape with `x-apitools-*` extensions and `preferred` naming the highest verified comparable version within the recorded checked scope (no global latest claim or new preference from conflicting evidence). Seeds are APIs.guru's list (CC0 for its contributors' work; frozen since April 2023, with `x-origin` on every entry) and the public-apis README (MIT; about 1,970 docs links, 20% overlapping APIs.guru), plus the curated catalog; a monthly deterministic job (conditional GETs and one GitHub tree listing per repository; 2,632 of APIs.guru's origins sit in 154 GitHub repositories, the other 1,381 spread over 769 hosts) re-verifies records, and LLM discovery runs only on demand per provider. It holds metadata only, never spec bytes, and is published as a separate artifact, not inside this module or on the code repository's `main`. User submissions are untrusted leads until the job verifies them. Considered and rejected: cloning APIs.guru (the clone has no `list.json`, and the data is frozen), mirroring spec bytes (only 27% of APIs.guru specs declare a license, and the module zip is capped at 500 MiB), and daily pushes of data into the code repository (consumers pin exact revisions). S05 delivered the per-provider record this would aggregate; its legal, crawl-policy and operator gates remain separate. | S05 is delivered and its record shape is stable, legal sign-off covers publication and a crawl policy (identifying User-Agent, rate limits, robots, takedown handling, GitHub API terms), and an operator wants a shared directory; then approve a bounded milestone with its own gating row for that sign-off. |
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

## Stage 11 cross-package refactoring

Approved review-intake amendment, 2026-10-06: 18 required milestones / 87 pending rows across the five owners. Scope and milestone IDs are unchanged; the coordinator records the approved publication proposal, which needs a separate execution-time grant. Full local source baseline remains `9364afac98c97b79c9e6b9335bd9c71d4e8722d5`, including the reviewed uncommitted planning state. This intake does not start a closing review or authorize implementation.

**Approved source.** User-approved complete proposal, 2026-10-06; source baseline `9364afac98c97b79c9e6b9335bd9c71d4e8722d5`. [Coordinated contract](../../../kinet/docs/stage11.md) defines both phases, cross-package order, compatibility and acceptance. The request to implement the proposal authorizes its planning files only.

One execution owner, serial execution and task commits under the later confirmed goal. Planning authorizes no code execution, commit, publication or external operation. Source publication requires separately named authority; a status marker or local build is not publication. Consumers must record exact accepted and published prerequisites before adoption. Default checks are offline, credential-free and model-free. No deployment, live ledger migration, real API/model/mail action or registration change.

## M82 — Shape production

**Stage/owner.** STG-11 Phase A; APItools. **Priority.** Serial position 4/18, not a review severity.
**Dependencies.** [UWS:C09](../../../uws/tabilet/docs/history/status-C09.md); exact accepted/published contract closure recorded before adoption. Serial gates and direct contract/regression dependencies are reconciled in the coordinator.
**Scope.** Project supported source metadata; Preserve identity and incomplete evidence; Serialize and test shape tables; Qualify consumers and publish.
**Acceptance.** All existing source families produce honest, deterministic metadata with explicit incompleteness. The package remains source tooling, not a workflow or credential runtime.
**Verification.** go test ./...; go vet ./...; affected race/API/source-family fixtures and OpenUdon/Udon consumer checks; git diff --check. Use bounded local artifacts, not remote discovery.
Use only disposable roots and fixtures. Preserve published schemas/wires, historic evidence, current runtime capability restrictions and the installed M44 service. Changed v3/package/worker identities require fresh approval; they do not preserve old grants.
**Downstream.** [UWS:M08](../../../uws/tabilet/memory-bank/status-M08.md), [Kinet:M46](../../../kinet/tabilet/memory-bank/status-M46.md), [OpenUdon:P09](../../../openudon/tabilet/memory-bank/status-P09.md). Reconcile exact accepted/publication revisions before advancing.
**Tasks/review.** [status-M82.md](status-M82.md), M82.1 complete and 3 pending task commit units; review 0/10, not started. The local eight-family producer is implemented; exact identity/security, wire qualification, consumer acceptance and publication remain pending.

## Stage 11 candidate dispositions

M82 is new source-tooling work with a shared public UWS contract, so it uses the cross-cutting M lane. It promotes no catalog/provider expansion, remote discovery, credential/runtime behavior or unrelated hardening candidate.
