# Milestones

Milestones are listed in priority order across long-lived domain lanes. Each
item lists scope and acceptance criteria. This file owns the roadmap, lane
meanings, active milestones, status-file index, dependencies, candidate
directions, milestone scope, and acceptance criteria.

Per-task completion state lives in one `status-<LANE><NN>.md` file per
milestone, such as [`status-M01.md`](status-M01.md). Individual status files
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

Active milestones: S02, C02, and S03 (pending review-remediation work;
S02 before C02, S03 independent). Latest completed milestone: M75 - Operation
Lifecycle Ranking Ownership. Apitools and Authoring are published, and OpenUdon/Ramen pin
both revisions with passing standalone test/vet.

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

| Milestone | Status File | State |
|---|---|---|
| M01 - Private Harness Bootstrap | [status-M01.md](status-M01.md) | Complete |
| M02 - Candidate Service Inventory | [status-M02.md](status-M02.md) | Complete |
| M03 - Catalog Entries And Spec References | [status-M03.md](status-M03.md) | Complete |
| M04 - Security Overlays And Auth Classification | [status-M04.md](status-M04.md) | Complete |
| M05 - Catalog Resolution And CLI Reporting | [status-M05.md](status-M05.md) | Complete |
| M06 - Catalog Data Expansion | [status-M06.md](status-M06.md) | Complete |
| M07 - Overlay Applied Inspection Views | [status-M07.md](status-M07.md) | Complete |
| M08 - OpenUdon Intent Advisory Integration | [status-M08.md](status-M08.md) | Complete |
| M09 - Catalog Quality Gates | [status-M09.md](status-M09.md) | Complete |
| M10 - Catalog Expansion Batch 2 | [status-M10.md](status-M10.md) | Complete |
| M11 - Catalog Spec Refresh Tooling | [status-M11.md](status-M11.md) | Complete |
| M12 - Existing Catalog Hardening | [status-M12.md](status-M12.md) | Complete |
| M13 - Infra/API Platform Expansion Batch | [status-M13.md](status-M13.md) | Complete |
| M14 - Provider Advisory Output | [status-M14.md](status-M14.md) | Complete |
| M15 - Refresh Review Reports | [status-M15.md](status-M15.md) | Complete |
| M16 - Refresh Artifact Reconciliation | [status-M16.md](status-M16.md) | Complete |
| M17 - Refresh Validation Reconciliation | [status-M17.md](status-M17.md) | Complete |
| M18 - Parseable Artifact Auth Overlays | [status-M18.md](status-M18.md) | Complete |
| M19 - Spec Protocol Classification | [status-M19.md](status-M19.md) | Complete |
| M20 - Human-Docs Overlay Completion And Catalog Stats | [status-M20.md](status-M20.md) | Complete |
| M21 - Provider Node Expansion Batch 4 | [status-M21.md](status-M21.md) | Complete |
| M22 - Provider Node Expansion Batch 5 | [status-M22.md](status-M22.md) | Complete |
| M23 - Provider Node Expansion Batch 6 | [status-M23.md](status-M23.md) | Complete |
| M24 - Provider Node Expansion Batch 7 | [status-M24.md](status-M24.md) | Complete |
| M25 - Provider Node Expansion Batch 8 | [status-M25.md](status-M25.md) | Complete |
| M26 - Provider Node Expansion Batch 9 | [status-M26.md](status-M26.md) | Complete |
| M27 - Provider Node Expansion Batch 10 | [status-M27.md](status-M27.md) | Complete |
| M28 - Provider Node Expansion Batch 11 | [status-M28.md](status-M28.md) | Complete |
| M29 - Provider Node Expansion Batch 12 | [status-M29.md](status-M29.md) | Complete |
| M30 - Provider Node Expansion Batch 13 | [status-M30.md](status-M30.md) | Complete |
| M31 - Provider Node Expansion Batch 14 | [status-M31.md](status-M31.md) | Complete |
| M32 - Provider Node Expansion Batch 15 | [status-M32.md](status-M32.md) | Complete |
| M33 - Provider Node Expansion Batch 16 | [status-M33.md](status-M33.md) | Complete |
| M34 - Provider Node Expansion Batch 17 | [status-M34.md](status-M34.md) | Complete |
| M35 - Provider Node Expansion Batch 18 And Reconciliation | [status-M35.md](status-M35.md) | Complete |
| M36 - Refresh Correction Review | [status-M36.md](status-M36.md) | Complete |
| M37 - Dropbox Stone-Derived OpenAPI Advisory Artifact | [status-M37.md](status-M37.md) | Complete |
| M38 - Google Discovery Converter Ownership | [status-M38.md](status-M38.md) | Complete |
| M39 - AWS Smithy Converter Ownership | [status-M39.md](status-M39.md) | Complete |
| M40 - AWS Smithy Native Model Ownership | [status-M40.md](status-M40.md) | Complete |
| M41 - Google Discovery Native Model Ownership | [status-M41.md](status-M41.md) | Complete |
| M42 - NVIDIA DSX Air Provider Entry | [status-M42.md](status-M42.md) | Complete |
| M43 - AWS Smithy Service Expansion | [status-M43.md](status-M43.md) | Complete |
| M44 - Google Discovery Service Expansion | [status-M44.md](status-M44.md) | Complete |
| M45 - Microsoft OpenAPI Service Expansion | [status-M45.md](status-M45.md) | Complete |
| M46 - High-Value SaaS OpenAPI Curation | [status-M46.md](status-M46.md) | Complete |
| M47 - Dropbox Stone Source Evaluation | [status-M47.md](status-M47.md) | Complete |
| M48 - Residual Google Discovery Coverage | [status-M48.md](status-M48.md) | Complete |
| M49 - n8n Public API OpenAPI Entry | [status-M49.md](status-M49.md) | Complete |
| M50 - n8n Protocol Connector Boundary Review | [status-M50.md](status-M50.md) | Complete |
| M51 - n8n Gap Report And Source Batch Freeze | [status-M51.md](status-M51.md) | Complete |
| M52 - n8n De-Emphasis And License-Risk Reduction | [status-M52.md](status-M52.md) | Complete |
| M53 - Docker API Source Coverage | [status-M53.md](status-M53.md) | Complete |
| M54 - Kubernetes Cluster API Source Coverage | [status-M54.md](status-M54.md) | Complete |
| M55 - Enterprise Applications API Source Coverage | [status-M55.md](status-M55.md) | Complete |
| M56 - OpenAPI-First Provider Batch | [status-M56.md](status-M56.md) | Complete |
| M57 - Docs-Overlay Provider Batch | [status-M57.md](status-M57.md) | Complete |
| M58 - Catalog Security Metadata Audit | [status-M58.md](status-M58.md) | Complete |
| M59 - Workato/Tray Signal Provider Batch | [status-M59.md](status-M59.md) | Complete |
| M60 - Workato-Signaled AWS Smithy Expansion | [status-M60.md](status-M60.md) | Complete |
| M61 - Advisory Overlay Linked-Docs Review | [status-M61.md](status-M61.md) | Complete |
| M62 - Portable Fnct Helper Catalog | [status-M62.md](status-M62.md) | Complete |
| M63 - OpenRPC Source Support | [status-M63.md](status-M63.md) | Complete |
| M64 - Postman/RAML/API Blueprint Import Evaluation | [status-M64.md](status-M64.md) | Complete |
| M65 - GraphQL Source Support | [status-M65.md](status-M65.md) | Complete |
| M66 - gRPC/Protobuf Source Support | [status-M66.md](status-M66.md) | Complete |
| M67 - OData Source Support | [status-M67.md](status-M67.md) | Complete |
| M68 - UWS 1.4 Source Artifact Corpus | [status-M68.md](status-M68.md) | Complete |
| M69 - Bounded Local Multi-Family Discovery | [status-M69.md](status-M69.md) | Complete |
| M70 - Bounded Remote Catalog Discovery | [status-M70.md](status-M70.md) | Complete |
| M71 - Parallel-Lane Harness Migration | [status-M71.md](status-M71.md) | Complete |
| S01 - Untrusted Source And Prompt Contract Hardening | [status-S01.md](status-S01.md) | Complete |
| M72 - Network, Storage, Cache, And OAuth Boundary Hardening | [status-M72.md](status-M72.md) | Complete |
| C01 - Generated Catalog And Security Aggregation | [status-C01.md](status-C01.md) | Complete |
| M73 - CLI And Coordinated Release Cleanup | [status-M73.md](status-M73.md) | Complete |
| M74 - Review Contract Corrections | [status-M74.md](status-M74.md) | Complete |
| M75 - Operation Lifecycle Ranking Ownership | [status-M75.md](status-M75.md) | Complete |
| S02 - Local, Offline, And Discovery Safety Remediation | [status-S02.md](status-S02.md) | Pending |
| C02 - Catalog Refresh Manifest Integrity | [status-C02.md](status-C02.md) | Pending |
| S03 - Operation Lifecycle Ranking Correctness | [status-S03.md](status-S03.md) | Pending |

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
| Next provider-catalog expansion | No provider batch or user-verifiable outcome has been approved after the completed historical catalog batches. | Approve a bounded provider-owned source batch, its source/auth review rules, and acceptance evidence. |
| Production adoption of experimental LAP/RFC 9727 discovery | The adapters are intentionally experimental, and OpenUdon still selects APIs.guru explicitly. | Collect reliability/coverage evidence and approve the downstream OpenUdon discovery policy and network budget. |
| Remove deprecated Discovery and Smithy wrappers | Sibling consumers may still rely on compatibility imports, and removal is a public breaking change. | Confirm no active sibling imports remain, approve release notes, and pass downstream compatibility checks. |
| Network and local-scan hardening batch | Review "apitools Code Review" Pass 1 (baseline `a5699e5`) findings #7 (LAP registry validates every entry by DNS before scoring/limit and fails the whole search on one bad `source_url`, `remote_discovery.go`), #9 (unsafe-IP list omits `240.0.0.0/4`, IPv4-compatible `::/96`, and `fec0::/10`, `download.go`), #11 (OpenAPI metadata paths skip the `sourceguard.CheckJSON`/`CheckYAML` preflight, `validation.go`, `local.go`, `local_source_discovery.go`), #12 (`.bin` is always treated as a protobuf descriptor, `local_source_discovery.go`), #13 (regular-file-to-FIFO swap between `Lstat` and `Open` can block, `local_read.go`, `internal/artifactio`), and #14 (`Discoverer` exposes no port/transport options, `discovery.go`). All were revalidated at `a5699e570fb5d6288fdd255c71d229fb062b50b3` as Lower severity; no supported scenario fails and no private-host or root-escape path was found. | Approve a hardening batch after S02, or promote any item earlier if it becomes reachable in a supported deployment or a downstream consumer depends on it. |
| Egress proxy support for guarded downloads | Review "apitools Code Review" Pass 1 finding #8: the guarded transport has `Proxy: nil`, and the only alternative, `AllowUnsafeHosts`, drops all host checks (`download.go`). Supporting a proxy trades dial-time IP filtering for proxy trust, which is a product/security decision. | An operator needs mandatory-proxy egress, and an approved policy defines target pre-validation, proxy trust, and the documented DNS-rebinding limitation. |

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
only when their closure evidence is available.

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

## M01 - Private Harness Bootstrap

**Goal.** Install the private planning harness for `apitools` without exposing
the harness files in the public repository.

**Scope.**

- Create canonical harness files under `../tofu/apitools`.
- Symlink `AGENTS.md`, `tabilet/memory-bank/`, and `tabilet/evolution/` from `../apitools` to
  the private harness snapshot.
- Ignore harness paths in the public `../apitools` repository.
- Preserve the existing public package behavior and verification commands.

**Acceptance.** `../apitools` exposes the harness through local symlinks,
`../tofu` tracks the canonical files, `../apitools` ignores the harness paths,
and `go test ./...`, `go vet ./...`, and `git diff --check` pass in
`../apitools`.

## M02 - Candidate Service Inventory

**Goal.** Build the first reviewed list of candidate services for a future core
OpenAPI catalog.

**Scope.**

- Define a candidate service record shape for inventory and review.
- Seed the initial candidates from `../try-n8n/reducibility/specs`, because
  those fixtures already match common OpenUdon advisory services.
- Use `../n8n/packages/nodes-base/nodes` only as a service-priority signal for
  later expansion.
- Classify each candidate by known local spec evidence, official OpenAPI status,
  official non-OpenAPI machine spec status, user-supplied-spec need, and
  auth/security review state.
- Avoid n8n runtime compatibility, node import support, or behavior coverage
  claims.

**Acceptance.** Downstream catalog work has a deterministic candidate service
inventory with source notes and review classifications, without introducing
catalog resolution, security overlays, API execution, or credential handling.

## M03 - Catalog Entries And Spec References

**Goal.** Convert reviewed candidate services into durable built-in catalog
entries and spec references.

**Scope.**

- Add a focused `catalog` package without refactoring the existing root package.
- Define provider, alias, category, source hint, and spec reference types.
- Track official OpenAPI and official non-OpenAPI machine-readable spec
  references separately.
- Preserve source authority, URL, version/date/hash when known, license notes,
  and update notes.
- Keep upstream/vendor OpenAPI documents separate from local metadata.

**Acceptance.** Downstream callers can list built-in providers and inspect known
spec references with provenance, while unverified candidates remain visibly
distinct from durable catalog entries.

## M04 - Security Overlays And Auth Classification

**Goal.** Add supplemental security-overlay metadata and auth/security
completeness classification for catalog providers.

**Scope.**

- Define security overlay schemas under `catalog`.
- Keep overlays separate from upstream/vendor specs and label their sources.
- Classify OpenAPI auth/security metadata as complete, present-incomplete,
  absent, overlay-required, intentionally-anonymous, or unknown.
- Support added or corrected `securitySchemes`, root-level or operation-level
  `security`, OAuth scopes, and API-key placement metadata.
- Do not resolve credentials, fetch tokens, sign requests, or choose runtime
  accounts.

**Acceptance.** Downstream callers can inspect auth/security completeness and
available overlays as metadata only, without treating overlays as upstream
provider truth or executing API operations.

## M05 - Catalog Resolution And CLI Reporting

**Goal.** Expose catalog resolution and reporting through stable library helpers
and CLI inspection commands.

**Scope.**

- Implement resolution precedence for user-provided OpenAPI paths/URLs,
  user-provided overlays, project-local documents, built-in spec references,
  and built-in overlays.
- Add deterministic reporting helpers for provider catalog and security status.
- Provide CLI inspection/reporting commands, such as `catalog list`,
  `catalog inspect`, and `catalog security-report`.
- Update the public README to explain the catalog as metadata only.
- Preserve existing exported API compatibility unless a breaking change is
  intentional and documented.

**Acceptance.** Operators and downstream callers can inspect catalog metadata,
resolve metadata with explicit local inputs taking precedence, and generate
security reports without live provider calls, operation execution, credential
resolution, or n8n compatibility claims.

## M06 - Catalog Data Expansion

**Goal.** Expand the built-in catalog with a curated batch of popular workflow
providers, preserving source-backed spec and auth/security metadata.

**Scope.**

- Freeze a 12-15 provider expansion batch before implementation begins.
- Use `../n8n/packages/nodes-base/nodes` as a priority signal only, not as
  runtime compatibility evidence.
- Add candidate records and promote reviewed providers into durable catalog
  entries only after source review.
- Verify official OpenAPI, official OpenAPI index, official non-OpenAPI
  machine-readable specs, or human-docs-only status for each provider.
- For each provider, run the same curation loop: try official
  OpenAPI/Swagger/Discovery download first; for OpenAPI/Swagger artifacts,
  inspect reusable security schemes/definitions and root/operation security
  requirements; add security overlay metadata when schemes or requirements are
  missing, ambiguous, stale, internally inconsistent, or incomplete relative to
  official auth docs and protected operations; and build a service-specific
  docs-derived endpoint advisory overlay with both endpoints and security only
  when no official OpenAPI exists and official docs expose enough endpoint
  instructions for a useful reviewed subset.
- Save generated advisory overlays and service-specific overlay builder
  programs as tracked catalog assets when an advisory overlay is generated.
- Record source authority, source notes, license notes, verified dates, quirks,
  user OpenAPI need, and auth/security completeness.
- Add source-backed security classifications or overlays where upstream specs
  are incomplete or ambiguous.
- Avoid vendoring upstream specs, executing API operations, credential
  handling, request signing, or provider account decisions.

**Acceptance.** The catalog has a larger deterministic provider set with
reviewed provenance, source-backed auth/security status, ignored local review
artifacts for downloaded official specs, tracked advisory overlays and builders
when applicable, and tests proving candidate/provider/security metadata remains
valid and ordered.

## M07 - Overlay Applied Inspection Views

**Goal.** Provide inspection-only views that show how catalog security overlays
would supplement provider metadata without mutating or exporting upstream
OpenAPI documents.

**Scope.**

- Define provenance-labeled inspection view types under `catalog`.
- Combine base catalog security classifications, spec references, and security
  overlays into a read-only advisory view.
- Preserve source labels for overlay-provided schemes, root security,
  operation security, OAuth scopes, and API-key placement.
- Report conflicts such as duplicate scheme names, missing referenced schemes,
  overlay-only additions, and unresolved operation matches.
- Expose inspection output through library helpers and catalog CLI reporting.
- Keep merged views metadata-only; do not write overlay-applied OpenAPI files.
- Do not execute API operations, parse credentials, sign requests, or treat
  overlays as upstream provider truth.

**Acceptance.** Downstream callers and operators can inspect overlay effects and
conflicts deterministically while the upstream spec and local overlay remain
separate provenance-labeled metadata.

## M08 - OpenUdon Intent Advisory Integration

**Goal.** Consume `apitools/catalog` from `../openudon` so intent authoring can
surface optional provider/spec/security advice without making catalog metadata
authoritative.

**Scope.**

- Add an OpenUdon adapter that resolves provider catalog hints from
  `apitools/catalog`.
- Improve optional `intent.hcl` authoring or review guidance with provider ID,
  known spec reference, user OpenAPI need, auth/security status, and overlay
  source notes.
- Preserve explicit project/user OpenAPI inputs as higher precedence than
  built-in catalog metadata.
- Keep catalog advice optional and fail-open when no catalog match exists.
- Avoid runtime credential binding, account selection, release gate decisions,
  workflow execution, or trusted-runner behavior changes.
- Keep `apitools` responsible for metadata and `openudon` responsible for
  workflow authoring/review behavior.

**Acceptance.** OpenUdon can surface catalog-derived intent advice when useful,
but existing explicit user inputs and non-catalog workflows continue to work
without depending on built-in catalog data.

## M09 - Catalog Quality Gates

**Goal.** Add static, offline quality checks that help maintain catalog
metadata as provider count, spec references, and overlays grow.

**Scope.**

- Define static quality rules for missing source notes, duplicate aliases,
  stale verification dates, missing auth/security status, unresolved overlay
  provider/spec references, and inconsistent user OpenAPI need.
- Add reusable catalog quality report helpers with deterministic output.
- Add CLI reporting, likely `apitools catalog check`, with text and JSON
  output.
- Keep M9 checks offline by default; do not probe remote URLs or fetch provider
  documents.
- Add tests for passing catalog data, failing records, deterministic report
  order, and CLI exit codes.
- Update README development commands if a new check command is introduced.
- Preserve existing catalog APIs unless any additive exported API is needed for
  quality reports.

**Acceptance.** Maintainers can run a deterministic offline catalog quality
check that catches common metadata regressions without network access, provider
API calls, credential handling, or runtime compatibility claims.

## M10 - Catalog Expansion Batch 2

**Goal.** Add the next frozen batch of popular workflow providers to the core
catalog with source-backed spec references, auth/security classification, and
quality-gate coverage.

**Scope.**

- Freeze Batch 2 before implementation: Zendesk, Twilio, SendGrid, Mailchimp,
  Zoom, Okta, ServiceNow, PayPal, Xero, QuickBooks, Linear, Monday.com,
  Pipedrive, Typeform, and Telegram.
- Treat each service as a separate task and commit unit once curation begins.
- Use `../n8n/packages/nodes-base/nodes` as a priority signal only, not as
  runtime compatibility evidence.
- For each service, try official OpenAPI, Swagger, OpenAPI index, or Discovery
  download first; save downloaded review artifacts under ignored
  `catalog-openapi-cache/openapi/` or `catalog-openapi-cache/google-discovery/`
  and register paths in `catalog-openapi-cache/cache.sqlite` when artifacts are
  saved.
- If no official OpenAPI source exists, record official human docs and, when
  those docs expose enough endpoint instructions, generate tracked
  service-specific docs-derived endpoint overlay assets plus builder programs in
  addition to auth/security overlay metadata.
- Add or update candidate records, durable provider entries, spec references,
  security classifications, and security overlays with source authority,
  verified dates, license notes, quirks, source refs, and source notes.
- Run `apitools catalog check`, catalog tests, and required repository checks
  after the batch.
- Avoid vendoring downloaded official specs, executing API operations,
  credential handling, request signing, account selection, or n8n runtime
  compatibility claims.

**Acceptance.** Batch 2 providers are represented in the catalog with durable
source-backed metadata, quality checks pass offline, generated advisory assets
are tracked only when appropriate, downloaded official artifacts remain ignored,
and each service can be inspected without live provider calls or credential
resolution.

## M11 - Catalog Spec Refresh Tooling

**Goal.** Make known-spec refresh and review artifact registration more
repeatable for maintainers without expanding the catalog by hand-running
one-off commands.

**Scope.**

- Add a metadata-only helper or CLI command that lists known catalog spec
  references eligible for refresh, including provider ID, spec ref ID, kind,
  URL, source authority, current verified date, and saved artifact path when
  registered.
- Add a bounded refresh/download path for selected catalog spec references that
  uses existing safe HTTP(S) download rules, unsafe-host rejection, size limits,
  redirects, and validation where applicable.
- Save refreshed official OpenAPI/Swagger/Discovery artifacts under ignored
  catalog cache directories and update cache registry paths instead of storing
  duplicate document blobs.
- Emit a deterministic refresh report with changed artifact paths, validation
  status, content hash or revision evidence when available, and required manual
  catalog metadata follow-ups.
- Keep refresh tooling opt-in and local; do not run network probes from
  `catalog check`.
- Do not automatically promote metadata, overwrite tracked advisory overlays,
  execute operations, resolve credentials, sign requests, or choose provider
  accounts.

**Acceptance.** Maintainers can refresh selected known spec references through
a safe, deterministic local workflow that records ignored review artifacts and
actionable metadata follow-ups while keeping catalog quality checks offline.

## M12 - Existing Catalog Hardening

**Goal.** Reduce unresolved status in the existing built-in catalog before
adding another provider batch.

**Scope.**

- Review providers that currently report unknown, needs-verification, or
  advisory-only status: Airtable, Calendly, Linear, Mailchimp, Monday.com,
  OpenWeatherMap, QuickBooks, Salesforce, ServiceNow, Shopify, Telegram, and
  Typeform.
- Re-run the source-backed curation loop for each targeted provider:
  official OpenAPI/Swagger/Discovery first, then official machine-readable
  alternatives, then docs-derived endpoint overlays only when no official
  OpenAPI exists and official docs expose enough endpoint instructions.
- Reconcile ignored official review artifacts, tracked docs-derived advisory
  overlays, tracked overlay builders, and SQLite artifact registrations.
- Tighten source notes, source refs, license notes, verified dates, user
  OpenAPI need, auth/security classification, and security overlays where the
  current catalog is ambiguous.
- Keep `catalog check` offline and avoid adding network probes to quality
  gates.
- Avoid API execution, credential resolution, request signing, automatic
  metadata promotion, or treating advisory overlays as provider truth.

**Acceptance.** Existing ambiguous providers have reviewed source-backed status
and artifact registration where appropriate, offline quality checks remain
clean, and unresolved cases are explicitly documented rather than silently
left as stale `unknown` or `needs-verification` metadata.

## M13 - Infra/API Platform Expansion Batch

**Goal.** Add a frozen infra/API-platform provider batch using n8n node
presence as a priority signal only.

**Scope.**

- Freeze Batch 3 before implementation: AWS S3, AWS Lambda, AWS SNS,
  Bitbucket, CircleCI, Cloudflare, Databricks, Elastic, Grafana, Jenkins,
  Netlify, Sentry, Snowflake, Splunk, and Supabase.
- Treat each service as a separate task and commit unit once curation begins.
- For each service, try official OpenAPI, Swagger, Discovery, or other
  official machine-readable sources first; save downloaded official artifacts
  under ignored catalog cache paths and register them in SQLite when saved.
- If no official OpenAPI exists, record official human docs and build tracked
  service-specific docs-derived endpoint overlay assets plus builders only when
  official docs expose enough endpoint instructions for a useful reviewed
  subset.
- Add or update candidate records, durable provider entries, spec references,
  security classifications, and security overlays with source authority,
  verified dates, source notes, source refs, license notes, quirks, and user
  OpenAPI need.
- For AWS providers, describe SigV4 as auth/security metadata only; do not
  implement signing or choose accounts.
- Avoid n8n runtime compatibility claims, provider API execution, credential
  handling, request signing, or automatic account decisions.

**Acceptance.** Batch 3 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.

## M14 - Provider Advisory Output

**Goal.** Provide a directly consumable catalog advisory report for operators
and downstream authoring integrations without making catalog metadata
authoritative.

**Scope.**

- Add reusable catalog advisory report helpers that combine provider metadata,
  spec references, user OpenAPI need, auth/security status, overlay IDs,
  registered artifact paths, and manual follow-ups.
- Add CLI reporting through `apitools catalog advisory [provider] [--cache
  <path>] [--json]`, resolving provider by ID, display name, or alias when a
  provider is supplied.
- Join SQLite artifact paths when a cache exists, but do not create a new cache
  just to render advisory output.
- Keep output deterministic and metadata-only; do not fetch remote documents,
  execute operations, resolve credentials, or apply overlays to upstream specs.
- Preserve explicit downstream user/project OpenAPI inputs as higher
  precedence conceptually; advisory output is guidance, not a release gate.

**Acceptance.** Callers can obtain deterministic provider advisory summaries in
text or JSON, including cache artifact paths when available, while offline
quality checks and existing catalog commands continue to behave as before.

## M15 - Refresh Review Reports

**Goal.** Add an offline review report for saved catalog artifacts and refresh
evidence so maintainers can audit local curation state without rerunning
network refreshes.

**Scope.**

- Add reusable refresh review helpers that inspect built-in refreshable spec
  references, existing SQLite artifact registrations, and saved files under the
  catalog cache directory.
- Report missing artifacts, registered paths, saved paths, bytes/hash evidence,
  parse or validation status where applicable, verified dates, and manual
  follow-ups.
- Add CLI reporting through `apitools catalog refresh-report [--cache-dir
  catalog-openapi-cache] [--cache catalog-openapi-cache/cache.sqlite] [--json]`.
- Treat OpenAPI/Swagger, Discovery/index documents, and non-OpenAPI
  machine-readable specs according to the validation categories defined in
  M11.
- Keep the report offline; do not download documents, promote metadata,
  overwrite overlays, execute operations, resolve credentials, sign requests,
  or choose provider accounts.

**Acceptance.** Maintainers can audit local refresh artifact state with a
deterministic offline report, including actionable missing/invalid artifact
findings, without turning review into a network or execution workflow.

## M16 - Refresh Artifact Reconciliation

**Goal.** Use the offline refresh review report to reconcile local catalog
artifact registrations before another provider expansion.

**Scope.**

- Run `catalog refresh-report` against the current cache and classify missing
  registrations, invalid saved artifacts, and validation failures.
- Refresh selected known spec references only when the source is already in the
  built-in catalog and the refresh command can save a validated review artifact.
- Update the tracked artifact registry source for successfully saved official
  review artifacts so `cache.sqlite` can be rebuilt as a local manifest.
- Record remaining validation failures as current importer/upstream-shape
  limitations rather than silently leaving them ambiguous.
- Keep downloaded official specs and SQLite cache files ignored, and avoid API
  execution, credential handling, request signing, or metadata promotion.

**Acceptance.** Successful refresh artifacts are registered reproducibly,
remaining invalid or missing artifacts are documented with current validation
outcomes, and offline catalog quality checks remain clean.

## M17 - Refresh Validation Reconciliation

**Goal.** Preserve official catalog refresh artifacts that are syntactically
parseable OpenAPI/Swagger but fail the current strict semantic validator, while
keeping import behavior strict and review-only.

**Scope.**

- Add curation-only validation statuses for parseable OpenAPI 3.x and Swagger
  2.0 artifacts that fail strict validation.
- Include validation diagnostics in selected refresh and offline refresh-report
  results so maintainers can distinguish strict importer limits from missing or
  malformed documents.
- Continue rejecting empty, malformed, unsupported, or external-ref documents as
  invalid refresh artifacts.
- Register newly saved parseable-invalid official artifacts only as local
  review artifacts; do not promote them as import-ready metadata or provider
  truth.
- Re-run the M16 missing/unregistered refresh set and reconcile any newly saved
  artifact paths into the tracked artifact registry source.
- Preserve boundaries: no provider API execution, credential resolution,
  request signing, account selection, automatic overlay promotion, or automatic
  metadata acceptance.

**Acceptance.** Refresh and refresh-report output classify parseable-but-strict
invalid official specs explicitly, newly saved review artifacts are
reproducibly registered where appropriate, truly invalid artifacts remain
blocked, and the full verification suite remains clean.

## M18 - Parseable Artifact Auth Overlays

**Goal.** Keep parseable-but-strict-invalid official OpenAPI/Swagger artifacts
usable as downloaded review evidence while exposing auth/security metadata
through explicit advisory overlays.

**Scope.**

- Add built-in security overlays for parseable-invalid official OpenAPI/Swagger
  artifacts whose auth metadata cannot be safely relied on through strict
  import parsing.
- Preserve existing selected refresh behavior that downloads and registers
  parseable-invalid artifacts as review-only local cache evidence.
- Keep overlay statuses review-oriented (`present-incomplete` or
  `overlay-required`) instead of promoting strict-invalid artifacts as complete
  provider truth.
- Ensure provider advisory, security report, resolve, and overlay inspection
  views surface the overlay IDs and source-backed auth/security notes.
- Preserve boundaries: no provider API execution, credential resolution,
  request signing, account selection, or automatic metadata acceptance.

**Acceptance.** Parseable-invalid official OpenAPI/Swagger artifacts remain
downloadable review artifacts, catalog auth/security guidance is available via
overlays, and catalog quality plus Go verification remain clean.

## M19 - Spec Protocol Classification

**Goal.** Represent OpenAPI-adjacent server-side API description sources as
first-class catalog metadata without treating them as OpenAPI import inputs.

**Scope.**

- Add a protocol/model classification separate from catalog source reference
  kind.
- Classify OpenAPI, Swagger, Smithy JSON, Google Discovery, Dropbox Stone,
  OpenAPI indexes, human docs, and unknown references.
- Preserve protocol version hints when known, including OpenAPI/Swagger
  versions from metadata or source notes.
- Expose classification in advisory, specs, refresh, refresh-report, and
  inspect output.
- Keep existing OpenAPI import and validation behavior strict; do not add
  provider API execution, credential resolution, request signing, or automatic
  converters for non-OpenAPI sources.

**Acceptance.** Downstream callers can distinguish source kind from API
description protocol/model family in library and CLI output, while refresh and
review workflows still save non-OpenAPI machine specs only as review artifacts
and all catalog quality checks remain clean.

## M20 - Human-Docs Overlay Completion And Catalog Stats

**Goal.** Close the known human-docs endpoint-overlay gap and make catalog
coverage statistics directly reproducible from the CLI.

**Scope.**

- Triage the remaining human-docs-primary providers without endpoint overlays:
  Databricks, Jenkins, Linear, Monday.com, Sentry, Splunk, and Telegram.
- Build tracked OpenAPI-shaped advisory endpoint overlays for REST-shaped
  services where official docs support a useful reviewed subset.
- Explicitly document no-overlay decisions where OpenAPI-shaped output would
  misrepresent the provider API model.
- Add `catalog stats` output for primary provider protocol counts, local
  artifact registry counts, refresh validation buckets, and byte totals.
- Document non-OpenAPI converter priority for Smithy JSON, Google Discovery,
  Dropbox Stone, OpenAPI indexes, human docs, and GraphQL-first services.
- Document the next catalog expansion queue and source-first acceptance rules
  before adding another provider batch.
- Preserve boundaries: no provider API execution, credential resolution,
  request signing, account selection, or treating advisory overlays as provider
  truth.

**Acceptance.** Maintainers can see catalog coverage and refresh artifact
status with one offline command, every human-docs-primary provider has either a
tracked advisory endpoint overlay or an explicit no-overlay decision, and
future non-OpenAPI conversion and provider expansion work has source-first
guidance.

## M21 - Provider Node Expansion Batch 4

**Goal.** Add another frozen batch of popular provider-node services to the
core catalog before starting non-OpenAPI converter work.

**Scope.**

- Freeze Batch 4 before implementation: Intercom, Chargebee, Webflow,
  ActiveCampaign, BambooHR, Contentful, Customer.io, Eventbrite, Freshdesk,
  Mailgun, Postmark, and Todoist.
- Treat each service as a separate task and commit unit once curation begins.
- Use `../n8n/packages/nodes-base/nodes` as a priority signal only, not as
  runtime compatibility evidence.
- For each service, try official OpenAPI, Swagger, OpenAPI index, Discovery, or
  other official machine-readable sources first; save downloaded official
  artifacts under ignored catalog cache paths and register them in SQLite when
  saved.
- If no official OpenAPI exists, record official human docs and build tracked
  docs-derived endpoint overlay assets plus builders only when official docs
  expose enough endpoint instructions for a useful reviewed subset.
- Add or update candidate records, durable provider entries, spec references,
  security classifications, and security overlays with source authority,
  verified dates, source notes, source refs, license notes, quirks, and user
  OpenAPI need.
- Keep Google Discovery, Smithy, Stone, and GraphQL converter work deferred
  until provider-node catalog coverage is substantially broader.
- Avoid n8n runtime compatibility claims, provider API execution, credential
  handling, request signing, account selection, or automatic metadata
  promotion.

**Acceptance.** Batch 4 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.

## M22 - Provider Node Expansion Batch 5

**Goal.** Add another frozen batch of popular provider-node services to broaden
REST-oriented catalog coverage before starting non-OpenAPI converter work.

**Scope.**

- Freeze Batch 5 before implementation: Acuity Scheduling, Bitly, Brevo,
  Clockify, Coda, DeepL, Figma, Ghost, Harvest, Help Scout, Iterable, and
  Mailjet.
- Treat each service as a separate task and commit unit once curation begins.
- Use `../n8n/packages/nodes-base/nodes` as a priority signal only, not as
  runtime compatibility evidence.
- For each service, try official OpenAPI, Swagger, OpenAPI index, Discovery, or
  other official machine-readable sources first; save downloaded official
  artifacts under ignored catalog cache paths and register them in SQLite when
  saved.
- If no official OpenAPI exists, record official human docs and build tracked
  docs-derived endpoint overlay assets plus builders only when official docs
  expose enough endpoint instructions for a useful reviewed subset.
- Add or update candidate records, durable provider entries, spec references,
  security classifications, and security overlays with source authority,
  verified dates, source notes, source refs, license notes, quirks, and user
  OpenAPI need.
- Keep Google Discovery, Smithy, Stone, and GraphQL converter work deferred
  until provider-node catalog coverage is substantially broader.
- Avoid n8n runtime compatibility claims, provider API execution, credential
  handling, request signing, account selection, or automatic metadata
  promotion.

**Acceptance.** Batch 5 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.

## M23 - Provider Node Expansion Batch 6

**Goal.** Add another frozen batch of REST-oriented provider-node services to
broaden catalog coverage before starting non-OpenAPI converter work.

**Scope.**

- Freeze Batch 6 before implementation: Action Network, Adalo, Affinity,
  Agile CRM, Bannerbear, Baserow, Beeminder, Brandfetch, Clearbit, ConvertKit,
  Copper, Discourse, Freshservice, Gong, and Grist.
- Treat each service as a separate task and commit unit once curation begins.
- Use `../n8n/packages/nodes-base/nodes` as a priority signal only, not as
  runtime compatibility evidence.
- For each service, try official OpenAPI, Swagger, OpenAPI index, Discovery, or
  other official machine-readable sources first; save downloaded official
  artifacts under ignored catalog cache paths and register them in SQLite when
  saved.
- If no official OpenAPI exists, record official human docs and build tracked
  docs-derived endpoint overlay assets plus builders only when official docs
  expose enough endpoint instructions for a useful reviewed subset.
- Add or update candidate records, durable provider entries, spec references,
  security classifications, and security overlays with source authority,
  verified dates, source notes, source refs, license notes, quirks, and user
  OpenAPI need.
- Keep Google Discovery, Smithy, Stone, GraphQL, and other converter work
  deferred until provider-node catalog coverage is substantially broader.
- Avoid n8n runtime compatibility claims, provider API execution, credential
  handling, request signing, account selection, or automatic metadata
  promotion.

**Acceptance.** Batch 6 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.

## M24 - Provider Node Expansion Batch 7

**Goal.** Add another frozen batch of uncovered provider-like n8n node services
to broaden catalog coverage before starting standalone converter work.

**Scope.**

- Freeze Batch 7 before implementation: UptimeRobot, PostHog, NocoDB,
  SeaTable, QuickChart, Marketstack, CoinGecko, Reddit, HackerNews, NASA,
  OpenThesaurus, and OneSimpleApi.
- Treat each service as a separate task and commit unit once curation begins.
- Use `../n8n/packages/nodes-base/nodes` as a priority signal only, not as
  runtime compatibility evidence.
- For each service, try official OpenAPI, Swagger, OpenAPI index, Discovery, or
  other official machine-readable sources first; save downloaded official
  artifacts under ignored catalog cache paths and register them in SQLite when
  saved.
- If no official OpenAPI exists, record official human docs and build tracked
  docs-derived endpoint overlay assets plus builders only when official docs
  expose enough endpoint instructions for a useful reviewed subset.
- Add or update candidate records, durable provider entries, spec references,
  security classifications, and security overlays with source authority,
  verified dates, source notes, source refs, license notes, quirks, and user
  OpenAPI need.
- Keep Google Discovery, Smithy, Stone, GraphQL, AI/vector-specific, and other
  converter work deferred unless a batch checkpoint explicitly changes that
  direction.
- Avoid n8n runtime compatibility claims, provider API execution, credential
  handling, request signing, account selection, or automatic metadata
  promotion.

**Acceptance.** Batch 7 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.

## M25 - Provider Node Expansion Batch 8

**Goal.** Add a frozen CRM, sales, and marketing provider-node batch to the
core catalog before standalone converter work.

**Scope.**

- Freeze Batch 8 before implementation: Autopilot, Drift, Freshworks CRM,
  GetResponse, HighLevel, Keap, MailerLite, Mautic, Monica CRM, Salesmate,
  Sendy, and Vero.
- Apply the same per-service curation loop, source-first evidence requirements,
  priority-only n8n evidence boundary, and no-execution/no-credential rules as
  M24.

**Acceptance.** Batch 8 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.

## M26 - Provider Node Expansion Batch 9

**Goal.** Add a frozen forms, no-code, and app-builder provider-node batch to
the core catalog before standalone converter work.

**Scope.**

- Freeze Batch 9 before implementation: JotForm, Form.io, Formstack,
  SurveyMonkey, KoBoToolbox, Bubble, Cal, Cockpit, Stackby, Airtop, TimeSaved,
  and FileMaker.
- Apply the same per-service curation loop, source-first evidence requirements,
  priority-only n8n evidence boundary, and no-execution/no-credential rules as
  M24.

**Acceptance.** Batch 9 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.

## M27 - Provider Node Expansion Batch 10

**Goal.** Add a frozen communications and social provider-node batch, then run
the first embedded backlog checkpoint.

**Scope.**

- Freeze Batch 10 before implementation: Facebook, Facebook Lead Ads, LinkedIn,
  Line, Mattermost, MessageBird, Matrix, Rocket.Chat, Twake, Twist, Twitter,
  and WhatsApp.
- Apply the same per-service curation loop, source-first evidence requirements,
  priority-only n8n evidence boundary, and no-execution/no-credential rules as
  M24.
- After the 12 service rows are complete, reconcile the remaining uncovered
  provider-like backlog, stale n8n priority paths, and false-positive
  provider-like classifications before starting M28.

**Acceptance.** Batch 10 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, the checkpoint is documented, and offline quality checks
pass.

## M28 - Provider Node Expansion Batch 11

**Goal.** Add a frozen commerce, operations, and business provider-node batch
to the core catalog before standalone converter work.

**Scope.**

- Freeze Batch 11 before implementation: Magento, Paddle, Gumroad, Invoice
  Ninja, Odoo, ERPNext, WooCommerce, Wise, DHL, Onfleet, Unleashed Software,
  and Workable.
- Apply the same per-service curation loop, source-first evidence requirements,
  priority-only n8n evidence boundary, and no-execution/no-credential rules as
  M24.

**Acceptance.** Batch 11 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.

## M29 - Provider Node Expansion Batch 12

**Goal.** Add a frozen infrastructure, security, and admin provider-node batch
to the core catalog before standalone converter work.

**Scope.**

- Freeze Batch 12 before implementation: Bitwarden, Cisco, Cortex, Home
  Assistant, MISP, NetScaler, Rundeck, SecurityScorecard, UrlScan.io, Venafi,
  Wekan, and Zammad.
- Apply the same per-service curation loop, source-first evidence requirements,
  priority-only n8n evidence boundary, and no-execution/no-credential rules as
  M24.

**Acceptance.** Batch 12 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.

## M30 - Provider Node Expansion Batch 13

**Goal.** Add a frozen enrichment, AI, and data provider-node batch while
keeping AI/vector-specific converter work deferred.

**Scope.**

- Freeze Batch 13 before implementation: Humantic AI, Hunter, Jina AI,
  LingvaNex, Mistral AI, OpenAI, Perplexity, Mindee, Peekalink, Phantombuster,
  UpLead, and Dropcontact.
- Apply the same per-service curation loop, source-first evidence requirements,
  priority-only n8n evidence boundary, and no-execution/no-credential rules as
  M24.
- For AI-oriented providers, catalog only source-backed API metadata that fits
  the existing provider catalog boundary; do not add AI workflow behavior,
  model-selection policy, or runtime credential binding.

**Acceptance.** Batch 13 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.

## M31 - Provider Node Expansion Batch 14

**Goal.** Add a frozen messaging and notification provider-node batch, then run
an embedded converter/protocol blocker checkpoint.

**Scope.**

- Freeze Batch 14 before implementation: Mocean, MSG91, Plivo, Pushbullet,
  Pushcut, Pushover, SIGNL4, sms77, Vonage, Zulip, Gotify, and Emelia.
- Apply the same per-service curation loop, source-first evidence requirements,
  priority-only n8n evidence boundary, and no-execution/no-credential rules as
  M24.
- After the 12 service rows are complete, summarize blocked providers by
  protocol shape, such as GraphQL, SOAP, SDK-only, AI/vector-specific, weak
  docs, or no useful OpenAPI-shaped subset, and decide whether converter work
  would unlock a cluster before M32.

**Acceptance.** Batch 14 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, the converter/protocol checkpoint is documented, and
offline quality checks pass.

## M32 - Provider Node Expansion Batch 15

**Goal.** Add a frozen content, publishing, and community provider-node batch
to the core catalog before standalone converter work.

**Scope.**

- Freeze Batch 15 before implementation: Disqus, Medium, npm, Spotify,
  Storyblok, Strapi, WordPress, Wufoo, YOURLS, Raindrop, Toggl, and Travis CI.
- Apply the same per-service curation loop, source-first evidence requirements,
  priority-only n8n evidence boundary, and no-execution/no-credential rules as
  M24.

**Acceptance.** Batch 15 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.

## M33 - Provider Node Expansion Batch 16

**Goal.** Add a frozen operations, webinar, and productivity provider-node
batch to the core catalog before standalone converter work.

**Scope.**

- Freeze Batch 16 before implementation: Flow, GoToWebinar, HaloPSA,
  LoneScale, Lemlist, Orbit, Oura, ProfitWell, Quickbase, SyncroMSP, Taiga, and
  Tapfiliate.
- Apply the same per-service curation loop, source-first evidence requirements,
  priority-only n8n evidence boundary, and no-execution/no-credential rules as
  M24.

**Acceptance.** Batch 16 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.

## M34 - Provider Node Expansion Batch 17

**Goal.** Add a frozen support, data, and utility provider-node batch to the
core catalog before final backlog reconciliation.

**Scope.**

- Freeze Batch 17 before implementation: TheHive, TheHive Project,
  APITemplate.io, Currents, Demio, E-goi, Mailcheck, Mandrill, Metabase,
  Nextcloud, Philips Hue, and PostBin.
- Apply the same per-service curation loop, source-first evidence requirements,
  priority-only n8n evidence boundary, and no-execution/no-credential rules as
  M24.

**Acceptance.** Batch 17 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, and offline quality checks pass.

## M35 - Provider Node Expansion Batch 18 And Reconciliation

**Goal.** Finish the frozen uncovered provider-like n8n node backlog and
reconcile catalog coverage before any standalone converter milestone is opened.

**Scope.**

- Freeze Batch 18 before implementation: Segment, Strava, Zoho, and uProc.
- Apply the same per-service curation loop, source-first evidence requirements,
  priority-only n8n evidence boundary, and no-execution/no-credential rules as
  M24.
- Recompute provider-like n8n coverage, confirm the original 136-root backlog
  has no untriaged omissions, and document false positives or deferred protocol
  families before starting post-M35 work.

**Acceptance.** Batch 18 providers are represented with deterministic
source-backed catalog metadata, auth/security classifications or overlays,
tracked advisory assets only when appropriate, ignored official review
artifacts where saved, the uncovered-provider reconciliation is documented, and
offline quality checks pass.

## M36 - Refresh Correction Review

**Goal.** Reduce known invalid official OpenAPI/Swagger refresh artifacts where
source review supports narrow deterministic correction, while leaving broader
malformed artifacts invalid.

**Scope.**

- Add provider-specific refresh validation corrections for the official
  HighLevel module specs, Strava Swagger, and Spotify Web API OpenAPI artifact.
- Keep SyncroMSP invalid unless a future source-reviewed normalization can
  handle its broader schema-shape issues without broad importer relaxation.
- Preserve raw downloaded artifacts in the ignored cache and surface correction
  notes in live refresh and offline refresh-review outputs.

**Acceptance.** HighLevel, Strava, and Spotify refresh validation reports move
from invalid to valid with correction notes; SyncroMSP remains invalid; catalog
stats reflect the reduced invalid bucket; deterministic tests and offline
quality checks pass.

## M37 - Dropbox Stone-Derived OpenAPI Advisory Artifact

**Goal.** Formalize the existing Dropbox core advisory overlay as a narrow
Stone-derived OpenAPI-shaped subset while keeping Dropbox's official catalog
source classified as Dropbox Stone, not official OpenAPI.

**Scope.**

- Preserve `dropbox/dropbox-api-spec` as the official machine-readable source
  with `dropbox-stone` protocol classification and keep official OpenAPI
  availability unavailable.
- Tighten the existing `dropbox-core-api-overlay` metadata so it explicitly
  records `official_openapi: false`, Stone-derived provenance, and source refs
  to the Stone repository, HTTP docs, and OAuth guide.
- Keep the reviewed subset limited to account info, folder listing, metadata,
  upload, download, and shared-link operations.
- Update advisory artifact registry and reporting wording where needed so
  advisory overlays may be derived from official machine specs, not only human
  docs.
- Do not build a general Stone-to-OpenAPI converter, add runtime protocol
  support, execute Dropbox operations, resolve credentials, or change `../udon`.

**Acceptance.** The regenerated Dropbox advisory overlay is deterministic,
catalog advisory output still reports Dropbox's official protocol as
`dropbox-stone`, the OpenAPI-shaped overlay is exposed only as advisory
metadata, and deterministic tests plus catalog quality checks pass.

## M38 - Google Discovery Converter Ownership

**Goal.** Move reusable Google Discovery to OpenAPI-shaped conversion ownership
into `apitools` so downstream projects share one deterministic metadata-only
converter.

**Scope.**

- Add `github.com/OpenUdon/apitools/googlediscovery` with byte and decoded-map
  conversion APIs.
- Preserve the existing bounded OpenAPI 3.0.1 shaping behavior for Discovery
  resources, methods, schemas, media upload metadata, parameter names,
  deterministic ordering, and Google OAuth scope metadata.
- Port focused Drive and Gmail fixtures and converter tests into `apitools`.
- Keep catalog Google Discovery references classified as `google-discovery` and
  structured review artifacts, not official OpenAPI.
- Do not add CLI commands, execute Google API operations, resolve credentials,
  fetch tokens, sign requests, or select Google accounts.

**Acceptance.** `googlediscovery` is covered by focused tests, full `apitools`
test/vet/diff checks pass, `../udon` delegates reusable conversion to
`apitools` while keeping native document adaptation locally, and no
`apitools` package imports `../udon`.

## M39 - AWS Smithy Converter Ownership

**Goal.** Add reusable AWS Smithy JSON to OpenAPI-shaped metadata conversion to
`apitools` while preserving the metadata-only boundary and enabling direct
downstream Smithy entrypoints.

**Scope.**

- Add `github.com/OpenUdon/apitools/awssmithy` with byte and decoded-map
  conversion APIs for one-service Smithy JSON 2.0 documents.
- Support official AWS service models using `restJson1`, `restXml`, and
  `awsQuery` protocols.
- Convert service metadata, Smithy HTTP bindings, AWS Query form metadata,
  schemas, recursive refs, and AWS service/signing traits into bounded
  OpenAPI 3.0.1-shaped metadata.
- Keep catalog AWS Smithy artifacts classified as `smithy` structured review
  artifacts unless an explicit caller invokes the converter.
- Do not add CLI commands, execute AWS API operations, resolve credentials,
  fetch tokens, sign requests, or choose AWS accounts or regions.

**Acceptance.** `awssmithy` is covered by focused restJson1/restXml/awsQuery
tests and malformed-shape checks, full `apitools` test/vet/diff checks pass,
`../udon` consumes the converter through explicit Smithy entrypoints, and
SigV4 signing is documented as later runtime work.

## M40 - AWS Smithy Native Model Ownership

**Goal.** Make `apitools/awssmithy` own native AWS Smithy model parsing and
remove OpenAPI-shaped Smithy conversion as a runtime contract source.

**Scope.**

- Add `Parse` and `ParseMap` APIs for one-service Smithy JSON 2.0 documents.
- Preserve service, operation, shape, binding, protocol, static
  query/header/payload fields, endpoint-prefix, and signing-name metadata
  without producing OpenAPI.
- Cover `restJson1`, `restXml`, `awsQuery`, `ec2Query`, `awsJson1_0`, and
  `awsJson1_1`.
- Remove `Convert` and `ConvertMap`; downstream callers must use native Smithy
  parsing or their own explicit compatibility adapter.
- Do not execute AWS operations, resolve credentials, sign requests, fetch
  tokens, or select AWS accounts or regions.

**Acceptance.** Native parsing covers the six AWS protocol families with
focused tests, `../udon` no longer uses converted OpenAPI-shaped Smithy
metadata for runtime-plan generation, no public Smithy-to-OpenAPI conversion
API remains in `apitools`, and required tests/vet/diff checks pass.

## M41 - Google Discovery Native Model Ownership

**Goal.** Make `apitools/googlediscovery` own native Google Discovery model
parsing and remove OpenAPI-shaped Discovery conversion as a runtime contract
source.

**Scope.**

- Add `Parse` and `ParseMap` APIs for official Google Discovery documents.
- Preserve service metadata, root/base URL handling, schemas, inherited
  root/resource/method parameters, flattened methods, OAuth scopes, and media
  upload hints without producing OpenAPI.
- Remove `Convert` and `ConvertMap`; downstream callers must use native
  Discovery parsing or their own explicit compatibility adapter.
- Do not execute Google operations, resolve credentials, fetch tokens, sign
  requests, or choose Google accounts or projects.

**Acceptance.** Native parsing is covered by focused Discovery tests,
`../udon` no longer uses converted OpenAPI-shaped Discovery metadata for
runtime-plan generation, no public Discovery-to-OpenAPI conversion API remains
in `apitools`, and required tests/vet/diff checks pass.

## M42 - NVIDIA DSX Air Provider Entry

**Goal.** Add NVIDIA DSX Air as a reviewed provider catalog entry using its
official OpenAPI schema and source-backed authentication metadata.

**Scope.**

- Add NVIDIA DSX Air candidate and provider records.
- Save the official DSX Air OpenAPI schema as a tracked review artifact under
  `catalog-openapi-cache/openapi/`.
- Register the schema in the catalog artifact registry.
- Add a security overlay for NGC API-key bearer authentication because the
  current official schema omits reusable security scheme metadata.
- Do not create API keys, call DSX Air operations, choose organizations, or
  manage DSX Air simulations.

**Acceptance.** Catalog provider, security, advisory, stats, refresh, test/vet,
and diff checks pass with DSX Air discoverable through provider lookup and
refreshable OpenAPI metadata.

## M43 - AWS Smithy Service Expansion

**Goal.** Expand first-class AWS Smithy source coverage for high-priority AWS
services represented in `../n8n`, using official AWS Smithy JSON AST models as
cataloged native API metadata.

**Scope.**

- Freeze the service batch: SQS, DynamoDB, SES or SESv2, IAM, Cognito,
  Textract, Rekognition, Comprehend, Transcribe, ACM, ELB, and ELBv2.
- Add or update durable provider entries, aliases, categories, official Smithy
  spec references, source notes, license notes, verified dates, quirks, and
  user OpenAPI need classifications for each service.
- Save official review artifacts from `aws/api-models-aws` under ignored
  catalog cache paths and register saved artifact paths when artifacts are
  checked in or intentionally tracked for review.
- Exercise `github.com/OpenUdon/awssmithy` native parsing against selected
  service fixtures that cover AWS protocol families and Smithy trait shapes
  used by the batch.
- Add or update SigV4/security metadata overlays or classifications as
  metadata only; do not resolve credentials, sign requests, choose accounts, or
  choose regions.
- Keep Smithy artifacts classified as `smithy`; do not treat them as official
  OpenAPI and do not reintroduce Smithy-to-OpenAPI conversion in `apitools`.

**Acceptance.** The AWS batch is discoverable through catalog/provider/spec
inspection with official Smithy provenance, native parser coverage protects
the batch's protocol shapes, catalog quality passes, and no API execution,
credential resolution, signing, account selection, or region selection is
introduced.

## M44 - Google Discovery Service Expansion

**Goal.** Expand first-class Google Discovery source coverage for high-priority
Google services represented in `../n8n`, using official Discovery documents as
cataloged native API metadata.

**Scope.**

- Freeze the service batch: BigQuery, Cloud Storage, Docs, Slides, Chat,
  YouTube, Tasks, Admin SDK, Analytics, Translate, Cloud Natural Language,
  Books, Business Profile, and Firestore where official Discovery coverage is
  available.
- Add or update provider entries, aliases, categories, official Discovery spec
  references, source notes, license notes, verified dates, quirks, and user
  OpenAPI need classifications for each service.
- Save official Discovery review artifacts under
  `catalog-openapi-cache/google-discovery/` when appropriate and register
  artifact paths for offline refresh/review.
- Exercise `github.com/OpenUdon/googlediscovery` native parsing against
  selected service fixtures covering nested resources, inherited parameters,
  media upload, OAuth scopes, and schema shapes used by the batch.
- Keep Discovery artifacts classified as `google-discovery`; do not treat them
  as official OpenAPI and do not reintroduce Discovery-to-OpenAPI conversion in
  `apitools`.
- Preserve metadata-only boundaries: no Google API execution, token fetching,
  credential resolution, project/account selection, or request signing.

**Acceptance.** The Google batch is discoverable through catalog/provider/spec
inspection with official Discovery provenance, native parser coverage protects
the batch's Discovery shapes, catalog quality passes, and no API execution,
credential resolution, token fetching, project selection, or account selection
is introduced.

## M45 - Microsoft OpenAPI Service Expansion

**Goal.** Make Microsoft Graph and Azure REST coverage first-class for the
Microsoft services represented in `../n8n`, using official Microsoft OpenAPI
sources where available.

**Scope.**

- Map n8n Microsoft nodes to existing or new catalog providers: Outlook,
  OneDrive, Teams, Excel, SharePoint, To Do, Entra, Graph Security,
  Microsoft Graph core, Azure Storage, and Azure Cosmos DB.
- Keep Microsoft Graph-backed services tied to official Graph OpenAPI metadata
  from `microsoftgraph/msgraph-metadata`; add service-level aliases and source
  notes so users can find each n8n-visible service without duplicating provider
  truth.
- Add Azure REST spec references from official Azure REST API specs for Azure
  Storage and Cosmos DB when source review confirms stable OpenAPI artifacts.
- Review auth/security metadata for OAuth scopes, tenant/account boundaries,
  and Azure resource scopes as metadata only; do not acquire tokens or choose
  tenants, subscriptions, accounts, or regions.
- Register saved official artifacts or explicit no-artifact decisions for
  offline refresh review.

**Acceptance.** Microsoft Graph service aliases and Azure REST providers are
discoverable with official OpenAPI provenance, auth/security classifications
are source-backed, catalog quality passes, and no Microsoft API execution,
token acquisition, tenant selection, subscription selection, or account
selection is introduced.

## M46 - High-Value SaaS OpenAPI Curation

**Goal.** Reconcile and harden first-class OpenAPI coverage for high-value
SaaS providers represented in `../n8n`, prioritizing official OpenAPI,
Swagger, or OpenAPI-index sources already present or readily source-reviewed.

**Scope.**

- Freeze the curation set: Slack, GitHub, GitLab, Jira, Trello, Stripe,
  Notion, Asana, Box, Discord, Twilio, SendGrid, Xero, Zendesk, Zoho, OpenAI,
  PagerDuty, Shopify, ServiceNow, QuickBooks, HubSpot, Airtable, Salesforce,
  and Linear.
- For providers already present in the catalog, verify provider aliases, spec
  references, registered artifact paths, auth/security classifications,
  overlays, and advisory/report output against current source-backed metadata.
- For gaps, add official OpenAPI, Swagger, or OpenAPI-index references and
  tracked/registered review artifacts where source review supports them.
- Add advisory overlays or explicit no-overlay decisions only where official
  specs are absent or incomplete and official docs support a useful reviewed
  subset.
- Preserve provider-specific quirks and manual follow-ups, especially for
  OAuth scopes, app installation models, workspace/org selection, and
  enterprise account boundaries.
- Avoid SaaS API execution, OAuth token fetching, credential resolution,
  workspace/org/account selection, or treating catalog metadata as executable
  approval.

**Acceptance.** The high-value SaaS set has reviewed catalog metadata,
source-backed official OpenAPI/Swagger/index references or explicit documented
gaps, auth/security classifications or overlays, deterministic advisory/spec
output, catalog quality passes, and no provider API execution or credential
handling is introduced.

## M47 - Dropbox Stone Source Evaluation

**Goal.** Decide whether Dropbox Stone deserves first-class source parsing or
whether the current Stone-derived advisory overlay remains sufficient for
OpenUdon-style authoring.

**Scope.**

- Review the official `dropbox/dropbox-api-spec` Stone source and the existing
  Dropbox Stone-derived advisory overlay from M37.
- Compare the n8n Dropbox node surface with the reviewed advisory subset and
  identify missing high-value operations, auth/security metadata, and schema
  shapes.
- Prototype or specify the minimal native Stone parser shape only if source
  review shows broad reusable value beyond the current advisory overlay.
- If a parser is not justified, record a no-parser decision with the exact
  coverage and maintenance rationale.
- If a parser is justified, open a follow-up implementation milestone with
  native Stone model APIs, fixtures, catalog protocol handling, and downstream
  consumer expectations.
- Do not execute Dropbox operations, resolve credentials, fetch tokens, choose
  Dropbox teams/accounts, or treat Stone-derived advisory metadata as official
  OpenAPI.

**Acceptance.** Maintainers have a source-reviewed decision for Dropbox Stone:
either a documented no-parser/no-expansion rationale with the current advisory
overlay retained, or a scoped follow-up implementation milestone for native
Stone parsing. Catalog quality and diff checks pass.

## M48 - Residual Google Discovery Coverage

**Goal.** Add the remaining n8n-visible Google services that have official
Google Discovery documents, and record why Google Ads is not promoted as a
first-class source entry in this batch.

**Scope.**

- Promote People API / Contacts, Perspective Comment Analyzer API, and
  Firebase Realtime Database Management API as native Google Discovery catalog
  providers.
- Register ignored Google Discovery review artifacts and extend parser
  coverage checks so the artifacts remain native Discovery inputs, not OpenAPI.
- Add source-backed OAuth overlay metadata from Discovery scope declarations.
- Add Google Ads as a reviewed candidate-only row with an explicit
  no-Discovery/no-OpenAPI promotion decision.
- Preserve the Firebase distinction between Discovery-covered management
  resources and instance-specific Realtime Database data REST paths documented
  outside the Discovery document.
- Do not execute Google APIs, resolve credentials, fetch tokens, choose Google
  Cloud/Firebase projects, or add protobuf/gRPC source-family support.

**Acceptance.** Contacts/People, Perspective, and Firebase Realtime Database
Management appear in provider list/advisory/spec/security views as native
Google Discovery sources with registered ignored artifacts and OAuth overlays;
Google Ads remains candidate-only with its promotion blocker documented; catalog
quality and diff checks pass.

## M49 - n8n Public API OpenAPI Entry

**Goal.** Add n8n's own public REST API as a reviewed OpenAPI provider entry,
without treating the broader n8n node directory as runtime compatibility
evidence.

**Scope.**

- Promote the n8n Public API from the n8n API node priority signal into a
  durable catalog provider entry.
- Reference the official n8n repository OpenAPI 3.0 root and public REST API
  documentation.
- Save and register the ignored OpenAPI review artifact.
- Record that the repository OpenAPI root is unbundled with relative `$ref`
  entries, so the saved artifact is official source metadata but may need
  bundling before strict import.
- Classify n8n auth/security metadata from the root OpenAPI `ApiKeyAuth` and
  `BearerAuth` alternatives.
- Keep database, queue, mail, SQL, LDAP, FTP, and other protocol connector
  nodes out of the OpenAPI provider catalog unless a future native protocol
  milestone is opened.
- Do not call n8n APIs, resolve API keys or bearer tokens, choose instances,
  users, projects, credentials, workflows, executions, or data tables.

**Acceptance.** n8n appears in provider list/advisory/spec/security/stats views
as an official OpenAPI provider with a registered ignored review artifact,
complete auth metadata, and an explicit unbundled-spec quirk. Catalog quality
and diff checks pass.

## M50 - n8n Protocol Connector Boundary Review

**Goal.** Record the post-M49 decision for n8n protocol connector nodes so
future provider expansion does not promote generic runtime connectors as
OpenAPI/Smithy/Discovery catalog sources.

**Scope.**

- Review n8n node roots for database, broker, mail, file, directory, generic
  execution, generic HTTP, webhook, and GraphQL connector signals.
- Update public documentation to classify these nodes as priority signals only,
  not provider rows, docs-derived OpenAPI overlays, native source parsers, or
  n8n runtime compatibility evidence.
- Define criteria for any future native protocol-family milestone, including
  stable schema artifacts, metadata-only parsing, and downstream-owned runtime
  binding.
- Preserve provider-specific catalog paths for true service APIs that happen to
  share protocol names, such as hosted control-plane APIs, when a future
  source-backed milestone explicitly scopes them.
- Do not add catalog providers, candidates, overlays, parser packages,
  database introspection, network protocol clients, credential resolution, host
  selection, queue/topic/database selection, or workflow execution.

**Acceptance.** Maintainers have a documented, reviewable boundary for n8n
protocol connector nodes; public docs point to the boundary; no provider or
parser behavior changes are introduced; required repository and private harness
diff checks pass.

## M51 - n8n Gap Report And Source Batch Freeze

**Goal.** Add a repeatable offline report that compares the local n8n node
tree with built-in provider catalog coverage and freezes the next
source-review batch for first-class OpenAPI or advisory-overlay work.

**Scope.**

- Add a reusable catalog report that accepts n8n node-root names and compares
  them to provider IDs, display names, aliases, and curated category coverage.
- Add CLI access through `apitools catalog n8n-gap-report`, reading directory
  names from `../n8n/packages/nodes-base/nodes` by default.
- Classify non-covered roots as provider API candidates, M50 generic protocol
  connector exclusions, local workflow utilities, GraphQL/protocol-family
  candidates, or internal directory exclusions.
- Freeze an 8-12 service source-review batch only where provider-owned OpenAPI
  evidence or strong official docs-derived overlay evidence exists.
- Update public catalog expansion docs with the command, report interpretation,
  and frozen batch.
- Do not add provider rows, candidate rows, overlays, parser packages, n8n
  workflow import, node runtime inspection, provider API probes, credential
  resolution, host/account selection, or execution behavior.

**Acceptance.** Maintainers can run the report offline against the local n8n
node tree, see deterministic classifications and the frozen source-review
batch, and use the public queue docs to start curation without treating n8n
node presence as runtime compatibility evidence.

## M52 - n8n De-Emphasis And License-Risk Reduction

**Goal.** Remove active n8n-specific catalog surfaces and make provider
curation source-first before release.

**Scope.**

- Remove the n8n Public API provider/candidate row, TimeSaved workflow helper
  row, artifact registration, auth classification, and related tests.
- Remove the `catalog n8n-gap-report` CLI/API and exported report types.
- Remove n8n node-directory evidence from built-in candidates and avoid
  active references to n8n node roots, `nodes-base`, or `try-n8n` fixtures.
- Rewrite public docs and active memory-bank guidance around provider-owned
  OpenAPI, Smithy, Discovery, Stone, official indexes, and docs-derived
  overlays.
- Keep historical milestone/status/evolution records as chronology, but mark
  M49-M51 as superseded for active implementation guidance.
- Do not copy third-party integration code, generated node metadata,
  credential behavior, descriptions, icons, or runtime workflow behavior.

**Acceptance.** Code, tests, public docs, catalog registry, and active
memory-bank guidance no longer position n8n as a provider source or priority
driver; remaining n8n mentions are historical records or explicit hygiene
tests; Go, catalog, consumer, and diff checks pass.

## M53 - Docker API Source Coverage

**Goal.** Add Docker-owned API source coverage for high-value container
workflows while keeping daemon and registry behavior metadata-only.

**Scope.**

- Add Docker Hub as a durable provider catalog entry using Docker's official
  Hub API reference and downloadable Swagger/OpenAPI specification.
- Add Docker Engine API metadata as a local daemon/control-plane API source
  using Docker's official Engine API OpenAPI references by version.
- Review Docker Registry HTTP API source availability and either add a scoped
  provider/protocol row or record a no-promotion decision if no durable
  official OpenAPI artifact is available.
- Add auth/security classifications and overlays as metadata only; do not log
  in to Docker Hub, contact a registry, connect to a Docker socket, pull
  images, inspect containers, or execute daemon operations.
- Update public docs and catalog queue notes where Docker control-plane APIs
  need a local-daemon safety distinction from SaaS provider APIs.

**Acceptance.** Docker Hub and Docker Engine appear in catalog inspection,
spec, stats, advisory, and security views with official source refs and clear
metadata-only boundaries; Docker Registry is either represented from official
source evidence or explicitly deferred; Go, catalog, consumer, and diff checks
pass.

## M54 - Kubernetes Cluster API Source Coverage

**Goal.** Add Kubernetes API source coverage based on cluster-published
Discovery and OpenAPI documents without contacting live clusters.

**Scope.**

- Add Kubernetes as a cataloged source family for cluster API metadata, with
  official docs for Discovery API, OpenAPI v2, and OpenAPI v3.
- Treat Kubernetes OpenAPI as user-provided or exported local metadata because
  each cluster publishes its own API surface, including enabled groups,
  versions, and CRDs.
- Add metadata-only handling guidance for `/api`, `/apis`, `/openapi/v2`, and
  `/openapi/v3` artifacts. Prefer OpenAPI v3 for richer schema fidelity while
  documenting v2 compatibility.
- Decide whether this milestone needs a lightweight Kubernetes source parser or
  only catalog/documentation support; if parser work is opened, it must parse
  local bytes/maps only and preserve group/version/resource provenance.
- Do not read kubeconfig, resolve cluster credentials, contact API servers,
  perform discovery calls, watch resources, dry-run applies, mutate resources,
  or claim runtime compatibility with any cluster.

**Acceptance.** Maintainers have a clear Kubernetes source model and catalog
entry or native-source plan that supports exported OpenAPI/Discovery artifacts
without live-cluster access; docs explain CRD/cluster-specific behavior and
runtime boundaries; Go, catalog, consumer, and diff checks pass.

## M55 - Enterprise Applications API Source Coverage

**Goal.** Add source-first catalog coverage for high-value enterprise
application APIs: NetSuite, SAP S/4HANA, SAP SuccessFactors, Oracle Fusion
Cloud Applications, and Workday.

**Scope.**

- Add NetSuite as the first implementation target. Record SuiteTalk REST Web
  Services and REST API Browser metadata as official OpenAPI 3.0 source
  evidence, but treat account-specific metadata-catalog URLs as user-provided
  metadata rather than built-in downloadable artifacts.
- Add SAP S/4HANA and SAP SuccessFactors as separate provider entries or a
  clearly grouped SAP family, based on official SAP Business Accelerator Hub
  source references. Record whether each reviewed surface is OpenAPI, OData,
  SOAP, or docs-only, and defer parser work if OData semantics would be hidden
  by a generic REST overlay.
- Add Oracle Fusion Cloud Applications as the Oracle application target, not a
  generic Oracle provider. Use official Fusion REST docs and `/describe`
  metadata as source evidence, while treating tenant-specific describe URLs as
  user-provided metadata.
- Add Workday provider metadata from official Workday SOAP/WWS and REST API
  documentation. Defer deep implementation unless review finds stable public
  REST/OpenAPI artifacts or a future WSDL/SOAP source milestone is opened.
- Add auth/security classifications and advisory notes only from official
  sources. Do not sign in to tenants, resolve OAuth clients, fetch tokens,
  call metadata endpoints, or infer customer-specific modules.

**Acceptance.** NetSuite, SAP S/4HANA, SAP SuccessFactors, Oracle Fusion Cloud
Applications, and Workday have catalog/candidate metadata, official source
refs, source-family classifications, auth/security status, and clear
user-provided metadata boundaries; any OData, WSDL/SOAP, or tenant-describe
parser decisions are explicitly recorded; Go, catalog, consumer, and diff
checks pass.

## M56 - OpenAPI-First Provider Batch

**Goal.** Add a source-first OpenAPI catalog batch for high-value provider APIs
with provider-owned OpenAPI evidence: Adyen, DocuSign, Auth0, Confluence,
Acumatica, BigCommerce, Cisco Meraki, and Confluent Cloud.

**Scope.**

- Freeze the batch from provider-owned official OpenAPI references only.
  Third-party integration catalogs may explain priority but must not provide
  provider names, descriptions, auth behavior, operation lists, categories, or
  ranking metadata.
- Add catalog, candidate, spec-reference, source-family, category, alias,
  auth/security, refresh, and advisory metadata for Adyen, DocuSign, Auth0,
  Confluence, Acumatica, BigCommerce, Cisco Meraki, and Confluent Cloud.
- Prefer durable downloadable OpenAPI/Swagger artifacts, official API
  reference download links, or official provider GitHub repositories. Treat
  tenant-specific Acumatica endpoint Swagger/OpenAPI URLs as user-provided
  metadata rather than built-in downloadable artifacts.
- Keep Confluence scoped to official Atlassian Confluence Cloud REST surfaces
  unless source review explicitly records a separate Data Center/server
  boundary.
- Do not sign in to provider tenants, create API keys, fetch OAuth tokens,
  access customer data, infer enabled tenant modules, or execute any operation
  from the discovered documents.

**Acceptance.** The eight M56 providers appear in catalog inspection, spec,
stats, advisory, and security views with official source refs, source-family
classification, aliases/categories, and auth/security classifications; any
tenant-specific or versioned artifact boundaries are recorded; Go, catalog,
consumer, and diff checks pass.

## M57 - Docs-Overlay Provider Batch

**Goal.** Add a source-first docs-overlay batch for provider APIs where the
best usable official source is documentation rather than a durable downloadable
OpenAPI artifact: Marketo, Aircall, Checkr, Airwallex, Apollo.io, AfterShip,
and Adobe Acrobat Sign.

**Scope.**

- For each provider, first confirm whether an official downloadable OpenAPI or
  Swagger artifact exists. If one exists and is suitable, record it as official
  OpenAPI metadata instead of building an advisory overlay.
- When no durable official spec is available, add catalog/candidate metadata
  plus a reviewed docs-derived endpoint overlay only from provider-owned API
  docs. Record a no-overlay decision where official docs cannot support a
  useful reviewed subset.
- Keep each advisory overlay narrow, deterministic, and clearly marked as
  non-official OpenAPI. Include source refs, operation scope, auth/security
  notes, and omission rationale in the overlay registry or no-overlay notes.
- Prefer stable workflow primitives such as contacts/leads, calls, background
  check candidates/reports, payments/transfers, enrichment/search, shipment
  tracking, and agreements only where official docs make request and response
  shapes clear enough for metadata-only import.
- Do not sign in to accounts, call APIs, create applicants/leads/shipments,
  upload documents, initiate financial movement, start background checks, fetch
  tokens, or scrape authenticated documentation.

**Acceptance.** The seven M57 providers have catalog/candidate metadata,
official source refs, auth/security classifications, and either a tracked
docs-derived advisory overlay or an explicit no-overlay decision; every
overlay is source-reviewed, deterministic, marked non-official OpenAPI, and
visible in advisory/security/reporting views; Go, catalog, consumer, and diff
checks pass.

## M58 - Catalog Security Metadata Audit

**Goal.** After M57 is complete, audit all durable catalog provider rows for
auth/security metadata completeness and remediate gaps with source-backed
security overlays, classifications, or explicit no-portable-overlay decisions.

**Scope.**

- For OpenAPI and Swagger rows, inspect `components.securitySchemes` or
  `securityDefinitions` plus root and operation-level `security`
  requirements. Add overlays when schemes or requirements are missing,
  ambiguous, stale, internally inconsistent, or incomplete relative to
  official auth docs and protected operations.
- For Smithy, Google Discovery, Dropbox Stone, and other native
  machine-readable source families, verify that native auth metadata is
  represented in catalog classifications or overlays without reintroducing
  source-to-OpenAPI conversion behavior.
- For human-docs-only and docs-derived advisory rows, verify there is either a
  source-backed security overlay, a tracked advisory overlay carrying reviewed
  auth/security metadata, or an explicit no-overlay/no-useful-subset decision.
- For account, tenant, cluster, instance, or user-specific rows, preserve
  `present-incomplete` unless exact exported user metadata can carry safe
  security metadata; record no-portable-overlay decisions where appropriate.
- Do not use n8n or other integration marketplaces as provider truth, and do
  not log in, fetch tokens, call provider APIs, or resolve credentials.

**Acceptance.** Every durable provider row is classified as complete from
upstream metadata, complete via overlay, remediated with an added overlay,
explicitly user/tenant-specific with no portable overlay, or queued with a
source re-review finding; advisory, security-report, inspect, catalog check,
Go, consumer, and diff checks pass.

## M59 - Workato/Tray Signal Provider Batch

**Goal.** Add a source-first provider batch identified from the local
`tray_workato_connectors.xlsx` market-demand signal: Canva, LaunchDarkly,
Klaviyo, Attio, Fivetran, Anthropic, Amplitude, Ashby, Braze, and Postman.

**Scope.**

- Treat Workato and Tray.io rows only as prioritization signals. Do not copy
  connector metadata, descriptions, auth behavior, operation lists, categories,
  examples, icons, or workflow behavior.
- For each provider, first look for durable provider-owned OpenAPI or Swagger
  artifacts. Use official OpenAPI where available for Canva, LaunchDarkly,
  Klaviyo, Attio, and Fivetran if source review confirms stable importable
  artifacts.
- Use docs-derived advisory overlays only when official docs support a narrow
  reviewed subset and no suitable provider-owned OpenAPI artifact is found.
  Likely docs-overlay decisions include Anthropic, Amplitude, Ashby, Braze,
  and Postman unless source review finds durable official specs.
- Add catalog/candidate/spec/source-family/auth/security/advisory metadata in
  the same slice for each provider, including M58-style security inspection
  and overlay decisions.
- Do not sign in, fetch tokens, call provider APIs, create records, send AI
  prompts, export customer data, or scrape authenticated documentation.

**Acceptance.** The ten M59 providers have durable catalog/candidate metadata,
official source refs, auth/security classification, source-family
classification, and either a provider-owned OpenAPI artifact, a tracked
docs-derived advisory overlay, or an explicit no-overlay decision; catalog
advisory, inspect, stats, security-report, security-audit, quality, Go,
consumer, and diff checks pass.

## M60 - Workato-Signaled AWS Smithy Expansion

**Goal.** Expand AWS Smithy-native source coverage for Workato-visible AWS
services that are not yet first-class catalog rows.

**Scope.**

- Freeze the batch from Workato-visible AWS service demand, then verify every
  promoted service against official `aws/api-models-aws` Smithy JSON models.
- Initial candidates are EC2, CloudWatch, Kinesis, RDS, SageMaker, API
  Gateway, Athena, Bedrock, Glue, GuardDuty, KMS, Secrets Manager, Security
  Hub, and Inspector2.
- Add catalog/candidate/spec/source-family/security metadata using native
  Smithy references only. Do not create Smithy-to-OpenAPI conversion behavior
  or advertise import-ready OpenAPI for AWS Smithy services.
- Add review artifact registration and native parser fixture coverage where
  local official Smithy artifacts are saved for maintainer review.
- Add SigV4 metadata overlays as advisory credential-placement metadata only;
  do not resolve AWS credentials, choose regions/accounts, sign requests, or
  execute AWS operations.

**Acceptance.** Each promoted AWS service has official Smithy source refs,
native protocol/source classification, advisory SigV4 metadata, catalog
advisory/inspect/specs/stats/security-report/security-audit visibility, native
parser coverage for saved artifacts where present, and Go, catalog, consumer,
and diff checks pass.

## M61 - Advisory Overlay Linked-Docs Review

**Goal.** Re-review existing docs-derived advisory OpenAPI artifacts so linked
official API docs are followed and related operations are captured when they
materially improve provider workflow metadata.

**Scope.**

- Inventory all registered advisory overlay artifacts and their builder files.
- For each overlay, inspect `x-apitools-overlay.source_refs`, builder source
  refs, and linked official API documentation reachable from those docs.
- Add linked operations to the same advisory OpenAPI when they are official,
  API-shaped, unauthenticated to document, and directly support the provider
  workflow represented by the overlay.
- Preserve original provider parameter descriptions where they help LLMs plan
  data pipelines, including constraints, accepted formats, alternate lookup
  flows, and related helper APIs.
- Add source refs for linked API docs that contribute operations, schemas,
  parameters, or auth/security metadata.
- Record explicit no-add notes when linked docs are conceptual, pricing or auth
  only, duplicated by existing operations, authenticated-only, unstable, or
  outside the overlay's reviewed workflow.
- Regenerate overlay JSON from builders; do not hand-edit generated artifacts
  without updating the builder.
- Do not scrape authenticated docs, call provider APIs, fetch tokens, resolve
  credentials, or claim provider-official OpenAPI status for advisory overlays.

**Acceptance.** Every registered advisory overlay has a reviewed linked-docs
decision, expanded operation/schema coverage where official linked docs support
it, updated source refs and no-add notes where needed, regenerated builders and
artifacts, refreshed artifact registry metadata, and catalog, Go, consumer, and
diff checks pass.

## M62 - Portable Fnct Helper Catalog

**Goal.** Expose a small public helper contract for Gmail runtime-adjacent
helpers, with pure raw-message rendering as the first fnct case and a narrow
OAuth2 bootstrap shape plus local login CLI for trusted runtime credential
fallback.

**Scope.**

- Add a minimal function descriptor type that names the helper, runtime type,
  invocation mode, input fields, output semantics, purity, and side effects.
- Add `helper/gmailmsg` with `gmail.render_raw`, accepting one request-body
  object and returning a Gmail API `raw` string.
- Add Gmail OAuth2 credential-bootstrap helpers that parse local operator
  credential material, produce a consent URL when authorization material is
  missing, and exchange/refresh access tokens only when a trusted runtime calls
  them.
- Add `apitools oauth google login` for operators to run a local consent flow
  or exchange a manual authorization code, then print environment export and
  marker-based data-file guidance.
- Keep helpers deterministic and side-effect free: no provider API calls,
  Gmail clients, account selection, or runtime registration in `apitools`.
  OAuth2 token acquisition is allowed only in the explicit runtime auth helper,
  not in fnct helpers or catalog workflows. The login CLI may exchange tokens
  only when the operator invokes it directly.
- Cover message encoding, template rendering, missing required fields, template
  missing-key failures, header injection rejection, and exported descriptor
  metadata.

**Acceptance.** Downstream OpenUdon and udon can import the public descriptor
and helper implementation to agree on `gmail.render_raw`, and udon can import
the Gmail OAuth2 bootstrap shape for runtime-owned credential fallback; the CLI
can produce refresh-token env guidance for marker-based `data.hcl`; `go test
./helper/...`, `go test ./...`, `go vet ./...`, dependent OpenUdon/udon
focused tests, and `git diff --check` pass.

M72 supersedes the OAuth portions of this historical milestone: credential
parsing, token refresh, consent, and authorization-code exchange now live only
in Udon's trusted runtime and CLI. `apitools` retains the pure
`gmail.render_raw` helper and descriptor.

## M63 - OpenRPC Source Support

**Goal.** Add metadata-only OpenRPC support for JSON-RPC source documents so
provider-owned method contracts can participate in source-aware authoring and
review without being flattened into OpenAPI.

**Scope.**

- Keep the `openrpc` package as the public boundary for OpenRPC parsing and
  summary APIs.
- Parse OpenRPC documents offline, preserving document info, servers, methods,
  params, results, errors, components, and source refs.
- Expose method summary and selector metadata for downstream OpenUdon/UWS
  binding and review evidence.
- Classify OpenRPC artifacts as a distinct source family in catalog and
  materialization output.
- Copy reviewed OpenRPC artifacts and provenance into output packages when
  materialized.
- Do not execute JSON-RPC methods, contact servers, resolve credentials, or
  implement runtime transport behavior in `apitools`.

**Acceptance.** OpenRPC artifacts can be parsed and summarized offline with
stable method selectors, catalog/materialization identifies OpenRPC as its own
source family, downstream authoring tools have enough metadata to bind UWS
operations, and parser, catalog, Go, and diff checks pass.

## M64 - Postman/RAML/API Blueprint Import Evaluation

**Goal.** Decide how `apitools` should treat Postman Collection, RAML, and API
Blueprint artifacts without promoting them to normative UWS source types or
adding runtime behavior.

**Scope.**

- Evaluate local artifact detection and classification posture for Postman
  Collection, RAML, and API Blueprint files.
- Treat these formats like Dropbox Stone-style source/advisory evidence until
  source review proves a stronger posture: preserve provenance, classify source
  family, and avoid overstating import readiness.
- For Postman collections, assume they are client/workspace artifacts rather
  than provider-wide API truth unless source evidence proves otherwise.
- For RAML and API Blueprint, evaluate source-aware import or conversion paths
  first, preferring OpenAPI or a stronger native source when semantics survive.
- Record conversion and no-conversion decision notes, including semantic loss
  risks, provenance requirements, and downstream review posture.
- Do not add UWS `sourceDescriptions[].type` values, execution behavior,
  credential resolution, remote workspace calls, or provider API calls in this
  milestone.

**Acceptance.** Maintainers have documented local detection, classification,
conversion/no-conversion, boundary, and verification decisions for Postman
Collection, RAML, and API Blueprint. Public notes describe their advisory
status, M64 is marked complete, and docs/repo diff checks pass.

## M65 - GraphQL Source Support

**Goal.** Add metadata-only GraphQL source support for UWS 1.4 `graphql`
source descriptions so provider-owned GraphQL schemas and operation documents
can participate in source-aware authoring and review without being flattened
into OpenAPI.

**Scope.**

- Create a focused GraphQL metadata package/API boundary.
- Detect local SDL, introspection JSON, and operation documents without probing
  live GraphQL endpoints.
- Parse enough schema, type, field, operation, variable, and source-ref
  metadata to build prompt-safe summaries.
- Define stable selector metadata for `sourceOperationId` and
  `sourceOperationRef`.
- Classify and materialize GraphQL as its own source family and UWS source type.
- Do not execute GraphQL operations, perform live introspection, resolve
  credentials, fetch tokens, choose endpoints, or select accounts/workspaces.

**Acceptance.** GraphQL artifacts can be parsed or inspected offline into
deterministic metadata, selector summaries are available to downstream
authoring/review tools, catalog/materialization identifies GraphQL as its own
source family, and focused parser/catalog, Go, vet, and diff checks pass.

## M66 - gRPC/Protobuf Source Support

**Goal.** Add metadata-only gRPC/protobuf source support for UWS 1.4
`grpc-protobuf` source descriptions so provider-owned protobuf descriptors can
participate in source-aware authoring and review without REST-shaped lowering.

**Scope.**

- Create a focused gRPC/protobuf metadata package/API boundary.
- Detect local `.proto`, descriptor set, and reviewed descriptor artifacts
  without using server reflection.
- Parse packages, services, methods, request/response messages, streaming
  shape, fields, and source refs.
- Define stable selector metadata for service/method `sourceOperationId` and
  descriptor `sourceOperationRef`.
- Classify and materialize gRPC/protobuf as its own source family and UWS source
  type.
- Do not execute RPCs, use server reflection, open network channels, negotiate
  TLS, resolve credentials, or select accounts/projects.

**Acceptance.** gRPC/protobuf artifacts can be parsed or inspected offline into
deterministic service/method metadata, selector summaries are available to
downstream authoring/review tools, catalog/materialization identifies
gRPC/protobuf as its own source family, and focused parser/catalog, Go, vet, and
diff checks pass.

## M67 - OData Source Support

**Goal.** Add metadata-only OData source support for UWS 1.4 `odata` source
descriptions so provider-owned or user-exported CSDL/service metadata can
participate in source-aware authoring and review without generic REST-shaped
overlays.

**Scope.**

- Create a focused OData metadata package/API boundary.
- Detect local CSDL XML/JSON and exported service metadata without calling live
  tenant `$metadata` endpoints.
- Parse entity types, entity sets, singletons, navigation properties, functions,
  actions, parameters, and source refs.
- Define stable selector metadata for entity operations, functions, and actions.
- Classify and materialize OData as its own source family and UWS source type.
- Do not execute OData operations, call tenant metadata endpoints, resolve
  credentials, log in to ERP tenants, execute `$batch`, or select
  accounts/tenants.

**Acceptance.** OData artifacts can be parsed or inspected offline into
deterministic metadata, selector summaries are available to downstream
authoring/review tools, catalog/materialization identifies OData as its own
source family, and focused parser/catalog, Go, vet, and diff checks pass.

## M68 - UWS 1.4 Source Artifact Corpus

**Goal.** Add a small reviewed local artifact corpus for the four UWS 1.4 source
families so refresh-review validation covers real official examples without
expanding parser or runtime behavior.

**Scope.**

- Add two official example artifacts each for GraphQL, OpenRPC,
  gRPC/protobuf, and OData under source-aligned `catalog-openapi-cache/`
  directories.
- Register the examples as curation-only artifact rows with provider IDs
  `graphql-examples`, `openrpc-examples`, `grpc-protobuf-examples`, and
  `odata-examples`.
- Add root refresh-review coverage that validates each committed artifact as
  `valid-structured`, asserts the expected protocol, records byte/SHA-256
  evidence, and confirms nonzero native operation summaries.
- Keep the examples out of the durable built-in provider catalog.
- Do not add new source types, OpenAPI lowering, parser expansion, runtime
  transport behavior, credentials, endpoint calls, or provider account
  behavior.

**Acceptance.** All eight selected official artifacts are tracked, registered
as curation-only local artifacts, validated through the refresh-review path with
nonzero operation summaries, documented as a local corpus, and focused
parser/catalog/root, full Go, vet, CLI smoke, and diff checks pass.

## M69 - Bounded Local Multi-Family Discovery

**Goal.** Give downstream authoring one safe local discovery API for every
first-class source family without moving workflow choices or remote lookup into
`apitools`.

**Scope.**

- Add `DiscoverLocalSources(context.Context, LocalSourceDiscoveryOptions)` and
  a versioned report with validated candidates, rejected and ambiguous files,
  truncation diagnostics, family, title, operation count, score, path,
  SHA-256, and provenance.
- Inspect only explicit regular-file or directory roots; reject symlinks,
  unsafe/special paths, oversized files, and cancellation.
- Detect and validate OpenAPI/Swagger, Google Discovery, AWS Smithy, AsyncAPI,
  GraphQL, OpenRPC, gRPC/protobuf, and OData through native parsers. Treat
  conventional directory names as hints rather than proof.
- Deduplicate identical content by digest. Require an explicit kind for
  ambiguous JSON or XML.
- Default to 10,000 visited entries, 100 accepted candidates, and 20 MiB per
  file. Reaching a count bound produces a blocking truncation diagnostic with
  narrowing guidance.
- Extend the native operation inventory entry point to the same eight source
  families while preserving their source kinds and selectors.
- Keep remote discovery, workflow selection, operation execution, credential
  behavior, and project artifact writes outside this API.

**Acceptance.** Mixed-root discovery covers every supported family;
deduplication, ambiguity, symlinks, non-regular/oversized paths, both limits,
and cancellation fail safely; existing callers remain compatible; full Go,
vet, standalone, downstream, CLI, and diff checks pass.

## M70 - Bounded Remote Catalog Discovery

**Goal.** Expand remote OpenAPI candidate discovery without weakening source
provenance, download safety, or the existing APIs.guru-first behavior.

**Scope.**

- Keep APIs.guru as the primary global source and the legacy public-apis probe
  as the final compatibility fallback.
- Add LAP Registry targeted search as an experimental secondary source. Return
  its reported original `source_url`, not the transformed LAP document, and
  label candidates unvalidated until normal import succeeds.
- Add provider-scoped RFC 9727 discovery when a caller supplies a provider URL
  or hostname. Require `application/linkset+json`, inspect OpenAPI-like
  `service-desc` links, deduplicate URLs, and do not recurse into nested
  catalogs.
- Preserve unsafe-host rejection, redirect and response bounds, context
  cancellation, deterministic result ordering, cache isolation by provider,
  and a visible service-description link bound.
- Expose direct CLI source selection and `--provider-url`; keep tests offline
  with local HTTP servers.
- Do not query Swagger Catalog or Scalar Registry as unauthenticated global
  aggregators, crawl the general web, fetch credentials, or execute operations.

**Acceptance.** Direct and automatic LAP/RFC 9727 lookup paths expose
provenance-labeled unvalidated candidates; APIs.guru remains first; RFC 9727
requires an explicit provider and fails safely on wrong media types, unsafe
hosts, malformed or oversized catalogs, link limits, redirects, timeouts, and
cancellation; cache and CLI contracts are covered; Go, vet, standalone,
downstream, CLI, and diff checks pass.

## M71 - Parallel-Lane Harness Migration

**Goal.** Adopt the current lane-aware memory-bank contract without losing or
reclassifying the completed M01-M70 history.

**Scope.**

- Define permanent `M`, `C`, and `S` lane meanings, cross-lane priority,
  dependency, and safe-parallelism rules.
- Normalize the legacy M1-M9 filenames/headings to M01-M09 and record that
  explicit migration while leaving later historical IDs and evolution
  snapshots intact.
- Make every status file discoverable by the two-digit lane glob and every task
  row parseable as `Item | State | Notes` with a backticked second-column state.
- Keep later work as unnumbered candidate directions with deferral reasons and
  promotion triggers.
- Update bootstrap instructions and evolution v21; do not add a mandatory goal
  protocol or alter public `apitools` behavior.

**Dependencies.** None. This is private harness metadata only.

**Parallel ownership.** No implementation milestone is active concurrently.
Future `C` and `S` milestones may run in parallel when their sections name
non-overlapping files and no shared public contract dependency.

**Downstream impacts.** The account-level memory-bank API runner can discover
all canonical status files and choose actionable work across lanes. Public
`apitools`, OpenUdon, Udon, and UWS contracts are unchanged.

**Acceptance.** Every status filename matches
`status-[A-Z][0-9][0-9].md`; the index links exactly one file per milestone;
all state-bearing rows use a supported backticked marker in the second column;
the updated unattended runner reports no actionable rows without needing an
API call; `git diff --check` passes in both harness-visible repositories; and
the public Go test suite remains green.

## S01 - Untrusted Source And Prompt Contract Hardening

**Goal.** Apply one bounded untrusted-input and prompt-safety contract to all
source parsers, inventories, authoring summaries, and security alternatives.

**Scope.** Bound parser depth and work, guarantee lexer progress, bound inline
content and reference expansion, sanitize and budget prompt-facing summaries,
preserve OpenAPI security OR/AND semantics, resolve bounded local composition
references, and return visible diagnostics instead of silently truncating or
discarding selected operations.

**Dependencies.** M71. The public security and authoring contracts intentionally
break and require coordinated OpenUdon and Ramen migration before release.

**Parallel ownership.** S01 owns parser packages, root inventory/authoring and
security files, and the corresponding OpenUdon/Ramen consumers. M72 may proceed
in parallel in download, refresh, materialization, cache, and Udon OAuth files.
Shared documentation, dependency pins, and review commits are serialized.

**Downstream impacts.** OpenUdon must select among credential alternatives and
Ramen must render bounded structured metadata without flattening security sets.

**Acceptance.** Malformed and deeply nested sources fail within deterministic
limits; prompt content is sanitized and budgeted; security requirement sets
retain AND-within/OR-between semantics; duplicate ranking and bounded reference
resolution are correct; fuzz/regression, full, standalone, downstream, vet,
staticcheck, CLI, and diff checks pass.

## M72 - Network, Storage, Cache, And OAuth Boundary Hardening

**Goal.** Close network and filesystem integrity gaps and move live OAuth token
exchange into the runtime-owned Udon credential boundary.

**Scope.** Reject unsafe URL credentials, ports, and special address ranges;
enforce wire and decoded download bounds; use one confined atomic write/read
layer for refresh, materialization, export, and cache artifacts; separate raw
and corrected refresh validation; require cache digests and bounded pruning;
and replace the apitools OAuth command with a state/PKCE-protected Udon command
that writes secrets only to an exclusive mode-0600 file.

**Dependencies.** M71. Udon migration and apitools OAuth removal are one
coordinated row. M72 file ownership is disjoint from active S01 work except for
serialized documentation and dependency updates.

**Parallel ownership.** M72 owns `download.go`, refresh, catalog materialize,
SQLite cache, apitools OAuth removal, and Udon credentials/CLI. S01 owns parser,
inventory, authoring, and security summary contracts.

**Downstream impacts.** Udon becomes the sole owner of authorization-code
exchange and refresh-token handling. Materialization and refresh callers gain
explicit collision and raw/corrected validation outcomes.

**Acceptance.** Default-safe downloads reject unsafe hosts, rebinding, ports,
userinfo, redirects, and decompression abuse; all artifact I/O is confined,
digest-checked, atomic, and rollback-safe; cache growth is bounded; apitools has
no token exchange path; Udon never emits secrets to stdout/stderr; focused,
full, standalone, downstream, vet, staticcheck, CLI, and diff checks pass.

## C01 - Generated Catalog And Security Aggregation

**Goal.** Replace hand-maintained catalog literals and implicit overlay
precedence with validated source data and deterministic generated indexes.

**Scope.** Store reviewed provider bundles as JSON, generate checked-in Go
indexes deterministically, validate catalog and overlay invariants at generation
time, build advisory indexes once, and represent security disposition per spec
with explicit mixed and conflict outcomes.

**Dependencies.** M72 safe-file and refresh contracts. C01 starts only after
M72 review closes.

**Parallel ownership.** C01 owns catalog source data, generator, generated
provider/candidate/overlay indexes, advisory aggregation, and catalog tests.

**Downstream impacts.** Catalog callers consume per-spec security dispositions
and explicit provider-level mixed status rather than sorted last-wins overlay
behavior.

**Acceptance.** Generation is deterministic and checkable; invalid IDs,
references, overlay scopes, or conflicting dispositions fail generation;
catalog golden count/digest replaces scattered count assertions; advisory
aggregation is indexed; full catalog, root, standalone, downstream, vet,
staticcheck, CLI, and diff checks pass.

## M73 - CLI And Coordinated Release Cleanup

**Goal.** Finish the breaking-contract rollout with a thin, consistent CLI and
verified sibling dependency updates.

**Scope.** Centralize command parsing and exit/output policy, keep catalog stats
in the root package and refresh registration in `sqlitecache`, remove
duplicated flag scaffolding, update public/private documentation, migrate all
sibling consumers, and pin the reviewed apitools version in OpenUdon, Ramen,
and Udon.

**Dependencies.** S01, M72, and C01 reviews must be complete.

**Parallel ownership.** M73 is the serialized integration and release lane.

**Downstream impacts.** Help exits 0 on stdout, usage failures exit 2 on
stderr, runtime failures exit 1 on stderr, positional values are not mistaken
for help flags, and the former OAuth entry point is documented under Udon.

**Acceptance.** All migrated sibling consumers compile and pass their complete
test suites; apitools has no stale compatibility surface or staticcheck
findings; CLI contract tests pass; memory-bank and evolution records match the
released behavior; full cross-repository verification and diff checks pass.

## M74 - Review Contract Corrections

**Goal.** Correct the focused post-release regressions in catalog security
scope, discovery API compatibility, protobuf format fallback, authoring prompt
ownership and budgeting, truncation reporting, and cache artifact integrity.

**Scope.** Keep catalog security evidence on the selected spec or provider-wide
scope; restore both historical project-URL import methods while exposing the
bounded report separately; fall back from heuristic proto-text parsing to
binary descriptors; deep-copy caller-owned authoring metadata; apply custom
prompt budgets once; mark selected-context overflow as truncated; and reject
zero-byte catalog artifacts before persistence.

**Dependencies.** M73. These changes restore documented public and safety
contracts without introducing a new product direction or compatibility break.

**Parallel ownership.** M74 is a serialized cross-cutting correction milestone
covering catalog resolution, discovery, gRPC/protobuf parsing, root authoring
and inventory, SQLite cache integrity, and their focused tests.

**Downstream impacts.** Existing discovery callers compile unchanged; OpenUdon,
Ramen, and Udon retain their current API integrations and gain corrected local
behavior when updated to the reviewed revision.

**Acceptance.** Each reported case has a deterministic regression test; full
workspace and standalone tests, vet, staticcheck, generated-catalog checks, CLI
help, sibling consumer tests, and diff checks pass; the milestone receives a
final review before closure.

## M75 - Operation Lifecycle Ranking Ownership

**Goal.** Own conservative API lifecycle sibling ranking beside the operation
summaries and source provenance it evaluates, removing that API-specific policy
from the generic Authoring engine.

**Scope.** Add `operationlifecycle` over `apitools.OperationSummary`, retain the
reviewed API-first scoring behavior with named constants, reject ambiguous
ties, and restrict `/upload` normalization to explicit Google Discovery
provenance. Migrate OpenUdon and Ramen consumers and remove the old Authoring
package. Do not add workflow semantics, source fetches, credential behavior,
account selection, or operation execution.

**Dependencies.** M74. Coordinated with Authoring I01/D01/M28 and the
OpenUdon/Ramen adapter changes.

**Parallel ownership.** M75 owns only `operationlifecycle/` in apitools plus
the explicitly coordinated downstream import bridges and memory-bank records.
It does not overlap parser, catalog, network, cache, or runtime behavior.

**Downstream impacts.** OpenUdon and Ramen must first consume a published
apitools revision containing the new package, then consume the coordinated
Authoring revision and pass standalone checks.

**Acceptance.** Full apitools test/vet and workspace consumer suites pass;
provenance regressions prove generic `/upload` paths are untouched; Authoring
contains no lifecycle-ranking package or apitools dependency; after publication
both consumers pin the new revisions and pass `GOWORK=off` test/vet.

## S02 - Local, Offline, And Discovery Safety Remediation

**Goal.** Make local scanning, offline cache use, and project-URL discovery fail
safely and predictably without weakening the untrusted-source guards.

**Scope.**

- Anchor symlink rejection at the caller-chosen root: resolve a root's own
  ancestors once, then reject symlinks only beneath it, in local reads,
  `internal/artifactio` roots, and `sqlitecache.Open`.
- Make `CacheModeOffline` perform no DNS or network access (syntax-only URL
  validation of cache keys) and serve any integrity-checked cached copy
  regardless of cache TTL, reporting its stored age. Read-write mode keeps TTL.
- Rebuild `LocalFiles` on the bounded local walker: per-file rejections instead
  of whole-scan failure, visit/byte bounds, digest deduplication, and walk
  errors on one entry recorded as rejections in both local scanners.
- Make project-URL import idempotent by reusing identical content or a stable
  name, and deduplicate discovery candidates by digest.
- Keep the 16-URL network bound but never fail discovery on URL count: import
  the first 16 unique URLs in order and record a truncation diagnostic;
  `ImportProjectURLsWithReport` exposes truncation through its attempts.
- Do not relax private-host rejection, add proxy support, or change public
  function signatures.

**Review provenance.** "apitools Code Review" Pass 1 (#1, #2, #4, #5, #6,
#10) and Pass 0 (F09, F10), stated baseline `a5699e5`, revalidated at
`a5699e570fb5d6288fdd255c71d229fb062b50b3`; relevant uncommitted changes at
revalidation were harness documentation only. Lineage: M69, M70, M72, S01
(completed; not reopened).

**Dependencies.** None.

**Parallel ownership.** `local.go`, `local_read.go`,
`local_source_discovery.go`, `download.go`, `import.go`, `discovery.go`,
`internal/artifactio/`, and the `sqlitecache` open path, plus their tests and
docs. Does not own `operationlifecycle/` or catalog refresh behavior.

**Downstream impacts.** C02 builds on the root-anchoring semantics in
`internal/artifactio`. OpenUdon `internal/synthesize` stops accumulating
duplicate imports and no longer fails discovery on long briefs; its suite must
pass against the workspace checkout. Publishing and re-pinning consumers need
separate authorization.

**Acceptance.** Regression tests cover a symlinked temp-parent root, offline
import of a non-resolvable cached hostname and of an entry older than the TTL,
an oversized file and a symlink beside a valid spec, repeated discovery
producing one file and one candidate per document, and more than 16 project
URLs with local candidates present. `go test ./...`, `go vet ./...`,
`git diff --check`, and `(cd ../openudon && go test ./...)` pass, and README and
`architecture.md` describe the delivered behavior.

## C02 - Catalog Refresh Manifest Integrity

**Goal.** Keep catalog refresh artifacts on disk consistent with the
`cache.sqlite` integrity manifest when a batch fails partway.

**Scope.**

- Stage refreshed artifacts and promote them only after the whole batch
  validates and registers, or register completed rows before reporting a
  failure; never leave an overwritten registered artifact whose SHA-256 and
  byte count disagree with the cache database.
- Surface partial results and the failing reference in the CLI report.
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
materialization of that provider succeeds. `go test ./...`, `go vet ./...`,
`go run ./cmd/cataloggen -check`, `go run ./cmd/apitools catalog check`, and
`git diff --check` pass.

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
paths from different documents. `go test ./...`, `go vet ./...`,
`git diff --check`, and the OpenUdon and Ramen workspace suites pass.
