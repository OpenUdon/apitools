# Catalog discovery design

**Status:** M81.1 design approved by the user on 2026-09-30. This document
records the approved design and delivered M81/M80 implementation below;
M81 is retired and M80 publication qualification is in progress. Design baseline: `fdc0a3f2647a0c4488f0c0a6334fb2647aa73447`.

Owners: [M81 foundations](../tabilet/docs/history/status-M81.md) and
[M80 discovery](../tabilet/memory-bank/status-M80.md). APItools implements
metadata discovery and artifact export. OpenUdon owns workflow source
confirmation; Kinet owns user interaction and browser routing. No source,
rank, license note or discovery outcome grants execution authority.

## 1. Request and report

Add a separate `apitools.catalog-discovery/v1` request/report contract using
the existing `StepContract`. Existing candidate requests, exported struct
shapes, JSON wires, parser defaults and ranking behavior remain unchanged.
The additive Go declarations and delivered behavior are documented below.

Request fields:

| Field | Meaning |
| --- | --- |
| `schema_version` | Exactly `apitools.catalog-discovery/v1`. |
| `contract` | Existing purpose, typed inputs/outputs and read/write/unknown effect. |
| `provider_keys` | Optional exact provider IDs, display names or aliases resolved by `Catalog.FindProvider`. An omitted list searches the catalog scope. |
| `filters` | Optional source-authority allowlist and license/redistribution evidence requirements; default includes explicitly labeled unknowns. |
| `limits` | Bounded retrieval work and returned candidates, with documented defaults and hard ceilings. |
| `remote_lookup` | Explicit opt-in only; absent or false means no network. |

The library accepts catalog/root configuration separately from the serializable
request. An untrusted request cannot select arbitrary host paths, transports
or endpoints. Provider keys are individually resolved, deduplicated to sorted
canonical IDs and never split on whitespace. `Microsoft Outlook` is one key.
An unknown key or unsupported request version is `blocked`, with a diagnostic;
it is not silently dropped or expanded into open retrieval.

Reports include the outcome, canonical requested scope, catalog identity,
examined and unexamined coverage, exclusions, bounded diagnostics, limits and
candidate truncation. Source-backed candidates carry native operation
identity, sanitized summaries, input/output compatibility, authentication
alternatives, effect evidence, purpose evidence and deterministic rank data.
Each states whether retrieval was provider-constrained or open. Reference-only
provider/artifact leads are separate and never masquerade as operations.

Local lookup reads an index; it neither parses provider artifacts per request
nor provisions them. Models can choose only returned operation references.
Reported explanations are source evidence and concise diagnostics, never
hidden model reasoning.

## 2. What qualifies as a match

Reuse `prepareStepContract`, source-backed summaries and existing typed
dimension comparisons. Apply a discovery-specific qualification rule without
changing `BuildOperationCandidates` or its ordering:

1. The artifact entry is current, digest-bound and has a native selector.
2. Purpose evidence comes from documented operation meaning. Identity/path
   overlap alone is a weak lead. At least two distinct normalized purpose
   terms must overlap that documentation and cover at least half the request's
   normalized purpose terms. Use the existing `purposeTokens` normalization.
   A one-term request cannot qualify and needs clarification.
3. All requested input/output comparisons are compatible, required operation
   inputs are satisfied, and no conflicting or indeterminate comparison is
   hidden by the score. Missing types, requiredness, selected nullable output
   gaps or unsupported selected semantics remain visible.
4. An explicitly requested read/write effect must be compatible. `unknown`
   imposes no read/write constraint, and an operation's unknown effect remains
   unknown; it is never an assertion of safety.
5. There is no blocking loss of identity, selected fields or security
   alternatives due to sanitation or limits. Unsupported *relevant* semantics
   prevent qualification; unrelated documented gaps do not assert absence.

Purpose overlap is a conservative lexical retrieval criterion, not proof of
semantic equivalence. The downstream user still reviews and confirms the
selected operation. A strong score cannot repair missing contract evidence.

One qualified operation is a `match` in the stated examined scope. More than
one is `ambiguous`, even when their scores differ; ranking helps the user
choose but does not resolve intent. Duplicate links to the same artifact and
native operation count once, with all provider links preserved. Unexamined
scope stays visible even alongside a positive match; no claim of global
uniqueness follows.

Sort by existing dimension score, descending, then by canonical provider ID,
spec-reference ID, artifact ID, source kind, SHA-256 and native selector,
ascending. Score ties stay reported ties; identity ordering only stabilizes
serialization. Never use an absolute path as a tie-breaker.

M80 records named-provider `CatalogPlan` parity fixtures and a small labeled
relevance set derived from OpenUdon `examples/eval` step shapes. Record
precision@k as a baseline, not a new acceptance threshold. All examples are
synthetic or repository-owned, without user data or model calls.

## 3. Outcomes and consumer actions

| Outcome | Required evidence | OpenUdon/Kinet action |
| --- | --- | --- |
| `match` | One operation qualifies in the examined scope. Coverage gaps remain explicit. | Present the operation and evidence for intent/source confirmation. |
| `ambiguous` | Multiple qualified operations need a choice. | Ask an intent/provider question; do not silently pick the first. |
| `no_qualifying_api` | Checked scope is complete under the applied constraints and filters; every relevant operation is definitively ruled out. | May offer a browser step or pending step, clearly naming the checked scope. |
| `insufficient_evidence` | No qualified operation, and scope or comparison remains indeterminate: missing root/index/artifact, unsupported metadata, weak purpose, stale entries, limits, cancellation or partial remote search. | Request the missing intent, configuration or document, or let the user explicitly choose another route. |
| `blocked` | Invalid request/configuration, unsafe path/host, corrupt index, identity violation or invalid provider key prevents supported lookup. | Explain refusal and the action needed to correct it. |

Outcome evaluation order: validate request/configuration; refuse blockers;
report qualified ambiguity or match; otherwise return insufficient evidence
if any relevant scope/comparison is unresolved; only then report scoped
no-match. A hard caller cancellation returns promptly with no positive result
and preserves the incomplete-search diagnostic where a report can be returned.
Never present partial evaluation as a completed no-match.

A result limit applies after qualification counts are known. Displaying only
one of several matches does not change ambiguity. If a work limit prevents
those counts being known, report incomplete coverage. Intentional filter
exclusions are named and narrow the scope; they do not prove global absence.
An empty user-specified scope or a metadata-only catalog with no examined
operations is insufficient evidence, not proof that a service has no API.

Curated `unavailable`, docs-only and similar provider facts are leads. They
never prove no API exists. Only `no_qualifying_api` permits **automatic**
browser fallback; a user can explicitly choose browser authoring after any
outcome without turning it into an absence claim. APItools performs no routing.

## 4. Catalog root and operator preparation

Add a root configuration naming:

- The explicit catalog/cache directory.
- Its existing artifact registration database, normally `cache.sqlite`.
- The confined index location, proposed default `operations.v1.json`.
- The catalog metadata, normally the embedded reviewed catalog; tests may use
  a synthetic catalog.

Registry and index paths are non-empty relative paths confined under the
selected root. Resolve symlinked ancestors of that root consistently with
`artifactio`; reject a symlink at the selected root and beneath it, path
escapes, special files and conflicting identities. Request/report wires do
not contain absolute host paths. The path is operator configuration, not
catalog evidence.

Without a root, return metadata-only provider leads and
`insufficient_evidence`, naming explicit root preparation as the remedy.
Do not search a sibling checkout, the working directory or a user's home.
An explicitly configured missing registry/index is also insufficient evidence;
discovery never creates, migrates, prunes or refreshes a cache.

Operator sequence after M81 lands:

1. Select an explicit root and prepare artifact registrations through the
   existing refresh workflow or existing local registration APIs. Network
   refresh is a separate explicit operation; indexing does not refresh.
2. Run the new offline `catalog index` command with that root and registry.
3. Configure the consumer with that exact root/index. Refresh registrations
   and rebuild the index explicitly when artifacts change.

The existing refresh CLI supports `--cache-dir` and `--cache`; M81.2 documents
the exact new index syntax when implemented. Commit a small synthetic fixture
root/catalog and fixture preparation recipe. Include no downloaded provider
specifications, SQLite runtime files or generated production indexes.

## 5. Index identity, coverage and resource bounds

Add a separate versioned index, proposed schema
`apitools.catalog-operation-index/v1`. It is derived metadata, with no cache
schema migration or changes to existing readers.

The header binds catalog SHA-256, schema version and builder/metadata contract
version. Each artifact binds registration ID, source family, exact byte count,
raw SHA-256 and sorted provider/spec-reference links. Operations bind native
selectors and sanitized source-backed metadata. Shared artifact content is
parsed once; provider-specific provenance and overlays remain separately
scoped. Different registration assertions for the same identity cannot be
silently merged.

Record all catalog references, including these coverage states:

| State | Meaning |
| --- | --- |
| `indexed` | Registered content verified and operations extracted within bounds. |
| `missing` | No registration or regular artifact file; include the reason. |
| `digest_mismatch` | File identity differs from registration. |
| `oversize` | Byte or structural/work budget exceeded; include the limit. |
| `parse_failure` | Registered content could not produce usable metadata. |
| `unsupported` | Source family or relevant semantics cannot be indexed. |

Any omitted operations or incomplete artifact traversal make the affected
scope unexamined. `indexed` cannot hide dropped operations. Catalog references
without machine-readable operations remain leads. Query readers validate the
index schema/catalog identity and current registration digests; stale entries
become unexamined. Neither registration nor an index signature proves runtime
authorization. Export rechecks actual source bytes before publication.

Canonical JSON, sorted arrays/maps and stable relative identifiers produce
byte-identical rebuilds from the same catalog, registrations and artifact
bytes. Exclude timestamps, elapsed measurements and absolute root paths from
the deterministic index; put measurements in a separate build report. Rebuild
through synchronized staging and atomic replacement so readers observe one
complete generation. Registry changes during building must prevent claiming
that generation is current.

Only the index builder may use the larger registered-artifact path, capped at
128 MiB. Existing direct parsers and `BuildOperationCandidates` retain their
20 MiB defaults. M81.3 must explicitly review and enforce structural, operation,
output and time budgets for that larger path and record peak-memory/time
measurements for the cached Graph and Cloudflare artifacts. Raising byte
limits alone does not qualify the parser. Deadline/cancellation checks must
cover extraction, not merely file reads. Limits cannot silently turn omissions
into `no_qualifying_api`.

M81.3 implementation uses the private OpenAPI index path with 4,000,000
structural items, depth 100, 100,000 operations, 256 MiB metadata and a
45-second cooperative context deadline. YAML nodes are checked once and reused
for decoding. Other families retain their 20-MiB parser limits. Decode and
normalization are bounded structurally but cannot be preempted mid-call by a
Go context; this is not a hard process CPU/RSS cap. Read/extraction cancellation
discards the result. The measured Cloudflare and Graph files fit these budgets
but retain explicit operation-summary gaps. M81.4 still owns the index builder
and coverage publication; no direct parser limit was raised.

M81.4 freezes index read/build ceilings sized for the recorded approximately
40,000-operation catalog, with synthetic boundary fixtures; M80 local queries
bound work and report size separately from whole-catalog indexing. Index
queries reuse sanitized metadata and never reparse each artifact. Avoid
claiming a hard process-RSS limit unless an implementation actually enforces
it; measured memory and enforced structural bounds are distinct evidence.

## 6. Artifact references and export round trip

Define an additive catalog artifact reference containing:

- Catalog identity and canonical provider ID.
- Spec-reference ID and artifact registration ID.
- Source kind, exact byte count and raw SHA-256.
- Native operation selector for operation selections; absent for leads.

No absolute path is identity. A single source artifact may have multiple
provider/spec links. Selecting one link does not authorize exporting every
artifact of that provider.

M81.5 adds a new artifact-scoped export/materialize entry point and options
type; it does not add fields to existing exported option structs. The selected
references round-trip as follows:

1. Validate catalog/provider/spec/registration linkage and expected digest.
2. Resolve only the registered confined source path under the explicit root.
3. Verify regular-file safety, byte count and digest from the actual bytes.
4. Copy only selected artifacts into source-family directories, plus the
   selected spec's applicable security overlays and provenance manifest.
5. Publish through the existing directory transaction after every selected
   artifact and overlay passes checks. Any failure publishes no partial tree.

Deduplicate repeated references without dropping provider provenance. Preserve
provider-wide overlays plus matching spec-scoped overlays and OR-of-AND auth
alternatives; never borrow another spec's auth disposition or apply overlays
silently into source documents. Source native selectors are preserved, not
rewritten into display method/path projections. OpenUdon checks the selected
selector against the exported source during binding. A digest/selector drift
requires fresh discovery and confirmation; no silent replacement or download.

Existing provider-level export, collision/reuse/force policy and all exported
struct shapes remain unchanged. This API copies metadata; OpenUdon retains
package confirmation, credentials, accounts and runtime authority.

## 7. Authority, licensing and optional remote evidence

Report source authority from the selected source, with explicit unknowns.
Keep provider-official sources, official-docs-derived advisory artifacts and
public-catalog provenance distinguishable. A provider's popularity or reviewed
catalog entry does not upgrade an artifact's authority.

Pass `license_note` through verbatim as evidence, under the metadata size
bounds, never as an instruction. Missing notes and `terms apply` do not
establish permission. License identification and redistribution permission
are separate unknown fields unless supported by reviewed evidence; do not
parse arbitrary notes into granted rights. Default filters include and label
unknowns. Explicit evidence requirements can exclude unknowns, but reports
must name those exclusions. Sanitized prompt projections are separate from
verbatim bounded source evidence, so raw notes are not forwarded blindly.

Remote lookup is off by default. M80's explicit opt-in uses guarded APIs.guru
search/download only: eight seconds total, at most three source documents,
20 MiB per document, existing redirect/dial-time/unsafe-host protection and
bounded parsing. Preserve digest and final URL provenance, omitting secret
userinfo/query/fragment from ordinary reports. Search limits, timeouts and
empty remote results do not prove global absence. Remote source bytes remain
ephemeral unless the operator separately authorizes provisioning; the offline
index is not implicitly changed. A remote candidate without a registered
artifact needs explicit provisioning and a subsequent digest-bound reference
before it can use the local artifact export. No network access, provider
operation or credential use is inferred from a `match`.

## 8. Downstream changes and acceptance

| Owner | Required consumer behavior |
| --- | --- |
| OpenUdon M94 | Adopt the accepted M80 revision including M81; configure an explicit root/index, preserve five outcomes and coverage, provision exact artifact references through the new export API, and retain source confirmation/native-selector checks. |
| Kinet W10 | Ask candidate-backed intent questions, keep missing evidence distinct from no-match, offer automatic browser fallback only after scoped no-match, and stage only confirmed sources. Hosted workers remain offline with explicit read-only metadata mounts. |
| OpenUdon M95 | Preserve the qualified discovery/provisioning journey as an iCoT-removal prerequisite. |
| Existing OpenUdon/Udon consumers | Preserve current exported APIs and verify workspace and standalone compatibility; no dependency on future M94 implementation. Ramen is excluded. |

The coordinated Stage 5 ledgers have already promoted M94/W10 and included
M81 before M80. This design lists their required contract reconciliation;
APItools M81.1 changes no sibling files. Reconcile consumers to the exact
accepted published APItools revision before their implementation starts.

M81/M80 verification must cover complete/incomplete scope, all five outcomes,
named multiword provider constraints, weak purpose and indeterminate types,
ties and relocation, stale/digest/unsafe-path failures, selected overlays,
shared artifacts, no partial export, cancellation and bounded work. Default
tests are offline and model/credential-free; remote tests use local servers.
Large-artifact and scale measurements are explicitly invoked local checks.
Record relevant rank precision separately from correctness acceptance.

Run APItools tests/vet in workspace and standalone modes, catalog generation
and quality checks, new command help and affected consumer compatibility at
their owner milestones. M81 publishes through M80, which records actual source
revision/module identity after the required publication review. No release or
implementation is claimed by approval of this document alone.

### Delivered index implementation (M81.4)

`catalog index --root DIRECTORY` (optionally `--registry`, `--index` and a
confined `--catalog catalog.json` for synthetic/custom fixtures) uses
`BuildCatalogOperationIndex` and atomic `WriteCatalogOperationIndex` with the
read-only `sqlitecache.ReadCatalogSpecArtifacts` adapter. No network is used.
The builder groups shared raw identity once and records every provider link
and catalog-reference coverage. It refuses registration drift before returning.
Partial metadata remains positive evidence with unexamined source coverage.
`ReadCatalogOperationIndex` checks identities, bindings and complete coverage
without re-reading source bytes. Registration drift marks the whole saved
generation stale and strips candidates. Operator-owned indexes are derived
metadata, not signatures or approval records. The byte/depth/work limits and
cooperative deadlines are documented in README; no hard RSS claim is made.

### Delivered selected-artifact export (M81.5)

`CatalogArtifactReference`, `CatalogArtifactExportOptions` and
`ExportCatalogArtifacts` implement the selected-reference round trip. Export
uses the configured catalog/root/read-only registration adapter, preserves
native selectors and raw bytes, verifies raw integrity and registration links,
and copies one physical file per shared family/digest. Every selected provider
link retains provenance and advisory labeling. Only provider-wide and matching
spec-scoped security overlays are emitted without transforming auth
alternatives or source documents. The index is not parsed during export.
The confined workflow output uses existing transaction/collision/reuse/force
semantics. Invalid sources, cancellation or registration drift leave no
partial output. The README documents fixed bounds and source-root exclusion.

### Delivered wire foundations (M80.1)

The additive Go types and strict bounded decoder now live in
`catalog_discovery_types.go` and `catalog_discovery_request.go`. Installation
options serialize as `{}` and cannot enter request JSON. Nil/omitted/null
provider keys mean open scope; an explicit empty array preserves empty scope
through serialization and needs insufficient-evidence handling. Structured
license identification and redistribution stay explicitly unknown unless
reviewed source fields support them; free-text notes never create permission.
The README records defaults/ceilings and the synthetic v1 wire fixture set.
Local retrieval/ranking and explicit remote opt-in are delivered.

### Delivered local retrieval (M80.2)

`DiscoverCatalogOperations` now returns scoped index metadata and references
under the M81 installation options. M80.3 applies the qualification/outcome rules above; scores remain advisory. Missing root/index yields metadata-only leads, stale index
entries have no candidates, and invalid configuration/identity blocks. Exact
provider and evidence filters report their narrowed scope and exclusions.
Work limits are explicit `work_limit` query coverage; display/context omissions
remain incomplete. Verbatim notes exceeding 64 KiB are omitted with an explicit
source evidence gap, never truncated into altered evidence or permission.
No source parsing, network activity, migration, access-time write, provisioning
or operation execution occurs. Explicit remote lookup is delivered by M80.4 as described below.

### Delivered ranking (M80.3)

`catalog_discovery_rank.go` reuses the existing typed comparisons without
changing the older candidate API. It requires strong documented-purpose overlap,
compatible input/output dimensions and requested effects, complete auth
alternatives, and no blocking metadata loss. Unused contract inputs remain
indeterminate. Each shared artifact/native operation counts once, preserving all
provider links. Counts precede display limits; positive scope gaps remain visible.
If projection omits every qualified candidate, the report becomes insufficient
rather than claiming a reviewable match. Sorting/ties use catalog identity.

The query's cooperative deadline covers index reading, comparison and report
projection. Caller cancellation clears positive evidence; budget expiration is
reported as incomplete evidence. Empty/filtered/unexamined scope never proves
absence. Explicit remote lookup is now delivered by M80.4.

The small repository-owned relevance baseline is described in
`testdata/catalog-discovery/relevance/README.md`. Observed precision@1 is 3/3;
this is a ranking baseline, not a semantic correctness threshold. Slack qualifies;
Jira remains insufficient due to credential-shaped field removal, and Drive
remains insufficient because the existing effect classifier does not recognize
upload as a documented write. Neither limitation is concealed or treated as
no API. No source/model/provider operations are performed by these fixtures.

### Delivered remote lookup (M80.4)

The request and installation must both enable remote lookup. The optional
installation client selects APIs.guru and cannot enter request JSON. A copied
client drops caches, cookie jars and redirect callbacks; existing owned URL,
redirect and dial-time guards remain in force. Only guarded APIs.guru search and
at most three guarded source downloads occur. The eight-second cooperative
budget covers search/download/parsing, within the overall query deadline. Each
source and catalog response is capped at 20 MiB, with lower configured bounds
preserved. Parsing uses existing inventory/prompt guards and a 16-MiB candidate
metadata cap per source; omitted operations remain partial evidence.

Every fetched source has a lead with sanitized final URL, exact digest/bytes and
public-catalog authority; operation candidates retain that same remote evidence.
Matching local/remote digest/native identities count once. Sources stay ephemeral;
registration/provisioning and credentials are downstream decisions. Exact
provider-constrained remote selection requires a matching catalog source URL;
public-catalog labels cannot establish that identity. Evidence filters apply
before downloading source documents. Empty, failed, timed-out or partially parsed
remote results remain insufficient unless other operations qualify. Unsafe
URLs/transports/redirects block; caller cancellation clears positive results.
Local HTTP fixtures verify these paths without network/provider/model operations.

### Scale and projection qualification (M80.5)

The opt-in root benchmark and synthetic preparation recipe are in README.
Qualification used a complete 12,000-operation index (one shared artifact, no
unexamined entries), exceeding the older whole-source ranking count limit.
The query reported all qualification counts before projecting 20 displayed
candidates; ambiguity stayed explicit. Visible tied indexes are remapped to the
returned candidate array, even when a larger candidate cannot fit. Projection
loss never fabricates uniqueness or turns partial comparisons into absence.
The measured evidence is a baseline, not a general throughput or RSS guarantee.

Iteration 1 qualification tightened metadata-loss handling: removed selected
fields never establish absence; providers without references retain explicit
unexamined coverage rather than borrowing another provider's completeness.
Discovery canonicalizes validated index arrays before bounded traversal, so
input order cannot change checked operations or shared reference order. Wrapped
remote timeouts retain their explicit classification. Shared groups beyond 32
allowed links become reference-only leads with `link_limit` coverage before
per-operation allocation; a narrower provider constraint can recover lookup.
The existing 32-KiB operation prompt budget is reapplied after rank evidence is
added, with oversized operations replaced by source leads/incomplete evidence.
