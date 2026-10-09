# Architecture

## Accepted M84 reserved-header remediation — 2026-10-09

M84 is accepted after four task gates and whole review2. Independently published
source `0c1383c6e2059e89cbf100ed8cb689be64d9ba57` resolves as `v0.0.0-20261009012335-0c1383c6e205`,
module sum `h1:GeJhKB7I42cH+bWw+8CoiPImyLmHKsqJih6E7uvUdlg=` and GoMod sum
`h1:WmUXlfBBoaI6vtv/wzeZfyO8q/pZMhnHiBckkAthYfk=`.
Reviewed source3773bd485fe31599e355719e0c5ecea0373bc553 and evidence closure
`b7425e2030489ce25a0ee523f83d694176f949f0` are independently observed on the authorized main ref.
[Handoff](../../docs/m84-release-handoff.md), [ordinary proof](../../docs/m84-ordinary-proof.json)
and [history](../docs/history/index.md) retain exact source/consumer/review evidence.
M82/M83 and frozen consumers remain at their recorded acceptance.

Bounded local-reference/name/location/duplicate checks precede the OpenAPI3
ignore rule. HTTP case-insensitive reserved header names are excluded before
requiredness/schema/serialization/completeness projection. Security scheme
OR/AND is independently retained; query/cookie/path and genuine header inputs,
Swagger2 and all public interfaces/wires are unchanged. Independent verification
reproduces corrected claims from raw bytes and rejects forged/stale metadata.

## Qualified Stage 11 source metadata — M82

The accepted M82 producer/verifier uses existing native parsers and public UWS
binding types. Explicit source IDs and local bytes/files reproduce exact SHA-256,
native selectors, schema projections and security symbols. Unknown source/dialect/
transport/presence/auth evidence stays explicit. HTTP metadata is emitted only
for declared HTTP leaves; schema compilation never loads external resources.
One invocation charges source/operation serialization and aggregate projected
schema work/bytes incrementally; over-limit or malformed inputs return no partial
table. Independent reproduction checks claimed tables before trust; hard CPU/RSS/
deadline/mount/network isolation belongs to consuming workers. Existing inventory/
candidate APIs and wires are unchanged; this package never approves or executes.
Exact public source and consumer evidence are in
[qualification](../../docs/m82-qualification.md). Kinet/OpenUdon worker adoption
and Stage 12 browser migration remain separately owned work.

## Repository Boundary

`apitools` owns a focused OpenAPI and provider API-source metadata boundary:

- safe remote catalog lookup, download, and validation of OpenAPI/Swagger
  documents and
  metadata-safe handling of catalog Google Discovery, AWS Smithy JSON, AsyncAPI,
  OpenRPC, GraphQL, gRPC/protobuf, and OData artifacts;
- local project file scanning with symlink and size protections, including a
  bounded family-aware discovery report over explicit file/directory roots;
- import into deterministic local `openapi/` files;
- operation inventories, indexes, summaries, selection ranking, and
  provenance-aware lifecycle sibling ranking;
- additive, versioned `BuildOperationCandidates` metadata that reads only
  explicit local source bytes/files, preserves SHA-256 and native selectors,
  emits consumer-readable summaries and evidence-bearing read/write/unknown
  effects, carries response nullability by field and schema ancestry, and ranks
  purpose/input/output/effect fit without binding or approving operations;
- auth/security summaries derived from OpenAPI security metadata;
- optional cache adapters;
- provider catalog, spec protocol classification, security-overlay metadata,
  and read-only overlay inspection views;
- deterministic catalog protocol to UWS source type mapping for OpenAPI,
  Google Discovery, AWS Smithy JSON, AsyncAPI, OpenRPC, GraphQL,
  gRPC/protobuf, and OData artifacts;
- Google Discovery native metadata parsing through the standalone
  `github.com/OpenUdon/googlediscovery` module;
- AWS Smithy JSON native model parsing;
- offline catalog quality reports and CLI gates for catalog metadata
  regressions.
- read-only provider advisory reports that combine catalog metadata, security
  status, overlay IDs, and existing artifact path registrations.
- offline refresh review reports that inspect existing cache registrations and
  saved artifacts without downloading documents or creating caches.
- offline catalog stats reports that summarize provider protocol buckets,
  artifact registry counts, and refresh validation buckets.
- offline catalog security audits that classify provider auth/security
  dispositions and inspect local OpenAPI/Swagger artifact security metadata;
  when root security is undeclared and scheme-bearing requirements exist on
  only a subset of operations, coverage is reported as partial-operation-
  security, with a follow-up to confirm the other operations are intentionally
  anonymous. Root declaration presence is recorded separately from its
  alternative count, so explicit anonymous root policies (`security: []` or
  `security: [{}]`) are not mistaken for absent root policy. Operation security
  declarations are counted separately from scheme-bearing operations, and
  explicit anonymous `security` arrays are declarations. More specific scheme
  diagnostics may take precedence while the partial-coverage follow-up remains
  visible.
- offline provider artifact resolution, materialization, and workflow export
  that copies existing local catalog artifacts and separate security-overlay
  JSON metadata with provenance manifests.
- pure helper packages and public helper descriptors for payload-shaping
  functions that trusted runtimes may register, without making `apitools` a
  workflow executor or credential resolver.

Related repositories:

| Path | Role |
|---|---|
| `../apitools` | OpenAPI tooling module and CLI. |
| `../openudon` | Consumes API metadata for authoring, review, package evidence, and trusted-runner handoff. |
| `../uws` | Public workflow semantics and model. |
| `../udon` | Private UWS/OpenAPI compiler, runtime execution, credential resolution, and operator OAuth token exchange. |
| `../tfconfig` | Static Terraform/OpenTofu parser only. |

## Layout

| Path | Role |
|---|---|
| `*.go` | Root `github.com/OpenUdon/apitools` package: client, discovery, validation, import, inventory, auth, and ranking APIs. |
| `operationlifecycle/` | Conservative same-source item-lifecycle ranking over root `OperationSummary` records; candidate roles require a trailing item parameter, with Discovery `{+name}` matched by method resource identity, and HEAD is treated as read. Shared purpose classification keeps POST updates/actions out of create, checking explicit create/update verbs on the operation ID before summary or tag wording; only a *trailing* path parameter marks a POST as an item-level action, so a parent-scoping parameter elsewhere in the path does not suppress a nested-collection create. A create-or-update PUT is treated as create consistently in both the seed's reported role and its sibling search. Family keys use operation IDs/paths, and update intent uses whole tokens across inflected forms. Source identity compares the highest-priority identifier both operations have populated (absolute path, then URL, then relative path/name); disagreement there is decisive, but an under-specified seed missing a higher-priority identifier still resolves against a fully described operation. A same-ID, same-source seed with its own method and path resolves to its exact match rather than reporting ambiguity. No fetch or execution. |
| `step_metadata.go`, `operation_summary.go`, `operation_effect.go`, `operation_rank.go`, `operation_candidates.go`, `operation_source_native.go` | Additive versioned operation-candidate contract, bounded local source adapters, exact digest/native-selector identity, prompt-safe consumer summaries, evidence-backed effects, and dimension-level contract ranking. External references and source URLs are never fetched; incomplete evidence remains visible or blocks ranking. |
| `cmd/apitools/` | Thin CLI wrapper over reusable package behavior. |
| `catalog/` | Metadata-only candidate inventory, durable provider entries, and official spec references. |
| `catalog/data/catalog.json` | Canonical reviewed C01 catalog bundle for candidates, providers, security classifications, overlays, and provenance. |
| `catalog/data/manifest.json` | Generated golden catalog version, SHA-256 digest, byte length, and record counts. |
| `catalog/catalog_gen.go` | Deterministic checked-in indexes generated from the reviewed catalog bundle. |
| `catalog/index.go` | Immutable reusable provider/spec/overlay/security lookup index used by resolution and advisory aggregation. |
| `cmd/cataloggen/` | Offline validator/generator with write and stale-output check modes. |
| `catalog_stats.go` | Package-owned offline catalog protocol, artifact-registry, and refresh-evidence aggregation policy. |
| `catalog-openapi-cache/` | Mixed catalog curation area: ignored downloaded specs, Discovery documents, and SQLite cache copies; tracked generated advisory overlays and service-specific overlay builders. |
| `googlediscovery/` | Deprecated compatibility wrapper over the standalone `github.com/OpenUdon/googlediscovery` native parser. |
| `awssmithy/` | Deprecated compatibility wrapper over the standalone `github.com/OpenUdon/awssmithy` native parser. |
| `openrpc/` | Metadata-only OpenRPC JSON parser and selector-summary API for JSON-RPC source documents. |
| `graphql/` | Metadata-only GraphQL parser and selector-summary API for SDL, introspection JSON, and operation documents. |
| `grpcproto/` | Metadata-only gRPC/protobuf parser and selector-summary API for `.proto`, JSON descriptor, and binary descriptor-set artifacts. |
| `odata/` | Metadata-only OData parser and selector-summary API for CSDL XML, JSON CSDL, and exported service metadata. |
| `openapidisco/` | Compatibility wrapper for local OpenAPI discovery and primary-candidate selection. |
| `sqlitecache/` | Optional SQLite implementation of the cache interface. |
| `helper/` | Pure payload-shaping helper contracts and implementations, such as Gmail raw-message rendering, intended for downstream authoring metadata and trusted runtime registration. |
| `internal/` | Private implementation helpers when needed. |
| `tabilet/memory-bank/` | Tracked lane-aware harness project memory; `milestone.md` owns lane meanings and `status-<LANE><NN>.md` files own task state. `lessons.md` keeps reusable lessons. |
| `tabilet/evolution/` | Tracked direction snapshots. |

The harness keeps completed pre-lane work in the zero-padded `M` series. Future
provider-catalog curation uses lane `C`, future source/discovery tooling uses
lane `S`, and cross-cutting work remains in `M`. Lanes classify long-lived
domains rather than phases. Independent milestones in different lanes may be
active together only when `milestone.md` records non-overlapping ownership,
resolved dependencies, and downstream impacts.

C01 introduces `apitools.catalog.v1` as the strict reviewed catalog-data
contract. Its decoder rejects unknown/trailing fields, invalid record shapes,
unsorted top-level records, duplicate provider/spec security scopes, and
provider, candidate, spec, or overlay cross-reference errors. The reviewed JSON
is embedded as the built-in catalog; `cmd/cataloggen` canonicalizes it and emits
a golden digest/count manifest plus deterministic provider, candidate,
classification, and overlay indexes. Package initialization verifies the
embedded bytes and every generated index against that manifest before exposing
independent copies through the catalog API. `go run ./cmd/cataloggen -check`
fails when any checked-in generated output is stale.

Catalog security aggregation groups evidence by provider plus optional spec
reference. A classification is baseline evidence; one or more overlays with a
single agreed status explicitly supply the effective disposition for their
scope. Differing overlay statuses inside one scope produce `conflict`, while
different effective statuses across scopes produce provider-level `mixed`.
This removes sorted last-wins behavior while retaining classification and
overlay provenance separately. Resolver selection uses an exact spec scope,
then provider-wide evidence. When a selected spec has neither, its resolved
status stays unknown; the provider aggregate remains available in the separate
report but is used as resolved evidence only when no spec was selected. The
resolver never silently selects another spec's disposition. The generator
rejects built-in same-scope conflicts, while library reports preserve explicit
conflict output for caller-supplied catalogs and quality review.

`CatalogIndex` validates and clones one catalog, then builds provider alias,
spec-reference, overlay-by-provider, and security-report maps once. Provider
advisory construction reuses that index and its precomputed security report for
every row instead of revalidating and rescanning the full catalog per provider;
all public accessors return independent copies. Catalog statistics policy lives
in the root library as `BuildCatalogStatsReport`, including primary protocol
selection, supported-family order, artifact-kind buckets, and refresh
validation buckets. `cmd/apitools` retains only input handling and rendering.

## Data Flow

Remote search is an ordered candidate lookup rather than a trust decision:

```text
query
  -> APIs.guru global directory
  -> experimental LAP Registry targeted search when APIs.guru has no match
  -> one RFC 9727 /.well-known/api-catalog lookup when a provider URL exists
  -> public-apis Markdown list and well-known-path probe
  -> untrusted Result candidates
  -> selected original source URL enters the existing bounded import validator
```

The LAP adapter parses `text/lap` search metadata and returns LAP's reported
original `source_url`, not the transformed LAP document. RFC 9727 discovery
requires `application/linkset+json`, reads OpenAPI-like `service-desc` links,
deduplicates URLs, rejects unsafe targets, and fails rather than returning a
partial result after the 100-link default. It fetches exactly one explicit
publisher catalog and does not follow nested `api-catalog` relations. The
result records mark both mechanisms experimental and unvalidated. A caller
context supplies any tighter total network deadline.

All default-safe remote fetches reject URL userinfo, non-HTTP(S) schemes,
non-conventional ports, localhost, non-global addresses, carrier-grade NAT,
documentation and benchmark ranges, NAT64, Teredo, and 6to4 addresses before
the request. Redirects repeat the URL checks, and a dedicated proxy-free
transport resolves and filters every dial-time address to resist DNS rebinding.
`Client.AllowedPorts` adds reviewed ports without bypassing address checks.
The client advertises gzip itself, then independently limits wire bytes and
decoded bytes to `Client.MaxBytes`; unsupported content encodings fail closed.
`AllowUnsafeHosts` is reserved for local fixtures and custom transports and
does not permit URL userinfo or unbounded bodies.

Public-apis lists accept both configured legacy JSON mirrors and Markdown
API tables with category headings. The default list is the public-apis
repository's README on GitHub raw, replacing the retired hosted JSON endpoint.
The upstream list is MIT-licensed (README attributes the public-apis
contributors); it is read at runtime and is not vendored into the module.
Parsing uses a 20-MiB response bound,
10,000-entry bound and 64-KiB row/JSON-entry bound. JSON also receives the
shared depth/token preflight. Malformed tables, cancellation and limits return
no searchable partial list; diagnostics identify the list URL. A valid empty
JSON entries array remains an empty result. Linked hosts enter the existing
guarded well-known-path probes only after list parsing succeeds.

`DiscoverLocalSources` walks only caller-provided roots. It resolves symlinked
ancestors of each selected root once, then rejects a symlink at that root or
beneath it. It counts every visited entry, rejects non-regular candidates,
reads regular files through the 20 MiB bounded reader, detects source families
from content, and validates them through their native parsers. Directory names are hints for
reporting likely invalid candidates, never proof of a family. Accepted content
is deduplicated by SHA-256 and returned with title, operation count, score,
path, and provenance. JSON or XML without exactly one family signal is reported
as ambiguous until the caller supplies an explicit kind. The 10,000-entry and
100-candidate defaults fail visibly through truncation diagnostics so a
downstream author cannot mistake a partial scan for complete evidence.
Adjacent advisory security sidecars are excluded for every supported
`source.ext.security.*`, `source.security.*`, and `source.security-overlay.*`
JSON/YAML naming form so auth metadata cannot become an API-source ambiguity.

The legacy `LocalFiles` API uses the same bounded walker in OpenAPI-only mode,
preserving draft-friendly metadata and score ordering. Invalid or unreadable
entries are skipped while valid siblings remain available; count-limit
truncation returns partial results together with an error. Project-text URL
imports process the first 16 unique URLs in source order, set a warning and
`Truncated` when more are present, and preserve successful partial results.
Same-name imports reuse byte-identical regular files, differing content uses a
collision-safe suffix, and discovery candidate reports deduplicate by
SHA-256. The historical project-URL methods retain their signatures; the
additive report method exposes truncation and diagnostics.

All direct source-parser entry points enforce the shared 20 MiB source-byte,
depth-100, and bounded structural/semantic-work contract before decoding.
Deprecated Smithy and Google Discovery `ParseMap` wrappers apply the equivalent
already-decoded value budget for nesting, retained strings, and traversal work.
Prompt-facing operations then pass through one structural sanitizer: identifiers
are capped at 256 runes, text at 2,048 runes, collections at 32 items, schemas at
60 fields, authoring work at 10,000 operations, individual operations at 32 KiB,
and ranked contexts at 512 KiB. Every sanitation or shortlist omission produces
a diagnostic. In every sanitized string, Cc controls, Unicode format
characters (Cf), tag characters (U+E0000-U+E007F), variation selectors
(U+FE00-U+FE0F and U+E0100-U+E01EF), and ANSI control sequences are replaced
with spaces before whitespace normalization; ordinary Unicode text is
preserved. Security/field compaction, unsafe selected references, work-limit
violations, and selected contexts that cannot fit fail closed rather than
silently changing the interpretation.

`BuildOperationInventory` applies the default prompt budget once.
`BuildAuthoringAPIDocuments` instead applies its caller-selected prompt budget
during inventory construction, so a reviewed larger bound is not irreversibly
truncated by the default before grouping.

OpenAPI security requirements remain ordered OR alternatives whose members are
AND requirements; an empty alternative explicitly permits anonymous access.
Downstream Authoring v2 preserves the sets, OpenUdon selects one stable
alternative before credential mapping, and Ramen refuses runnable output while
an alternative remains unresolved.

```text
service hint / project text / URL / local openapi directory
  -> Client.Search, Discoverer, LocalFiles, or Client.Import
  -> safe download or bounded local read
  -> OpenAPI/Swagger validation and metadata extraction
  -> OperationInventory / OperationIndex / AuthoringAPIDocument
  -> downstream OpenUdon authoring, review, and package evidence
```

Step-contract candidate flow:

```text
explicit local files/bytes + step purpose/inputs/outputs/effect
  -> bounded native parser for the declared source family
  -> exact raw-content digest and source-native operation selector
  -> sanitized operation summary + capability/evidence gaps
  -> advisory purpose/input/output/effect ranking
  -> report with ties, conflicts, unsupported details, or blocking truncation
```

`BuildOperationCandidates` handles OpenAPI/Swagger, Google Discovery, AWS
Smithy, AsyncAPI, GraphQL, OpenRPC, gRPC/protobuf, and OData using their
existing parsers. Only explicit local bytes or paths are read; a supplied URL
is provenance only. OpenAPI security alternatives retain OR-of-AND grouping.
OpenAPI response nullability is attached to each field together with nullable
schema ancestors; only selected nullable outputs make their output match
indeterminate. An unrelated nullable response field does not downgrade other
selected outputs. Effect analysis leaves an unrecognized leading action
unknown instead of searching later wording for a read/write token.
Native `OperationSummary.Method`/`Path` fields for non-OpenAPI families are
display projections; consumers identify the source with kind, digest, and
native selector. Unsupported auth, nested shape, response, streaming, or
requiredness semantics remain partial/unsupported rather than becoming
positive compatibility evidence. Negated or conflicting effect text remains
unknown; OData action/function kinds can provide effect hints, but entity-set,
singleton, and navigation resource selectors do not assert read-only behavior.
OpenAPI security alternatives retain OR-of-AND structure; a requirement with a
missing or incomplete scheme is retained but marks candidate auth capability
partial with an explicit gap. Native schema compositions and dynamic
`additionalProperties` values are not expanded, and each omission is reported
as a summary gap. Effect analysis bounds source meaning text before
tokenization, and source adapters honor cancellation between operations.
Synthetic OpenAPI response-body requiredness remains unknown when its schema
does not establish it. An over-context result is reduced to a bounded
diagnostic-only report. Ranking is informational and never binds an operation,
resolves credentials, chooses an account, or authorizes execution.

Provider catalog flow:

```text
candidate service inventory + provider-owned source evidence
  -> candidate validation and deterministic listing
  -> reviewed service set for later catalog entries

curated provider entry + official spec reference
  -> deterministic catalog listing and provider lookup
  -> source kind and protocol/model classification
  -> downstream import, inventory, or review guidance

curated provider entry + optional security overlay
  -> catalog resolver
  -> spec metadata and auth completeness report
  -> downstream import, inventory, or review guidance

curated provider entry + security classification + optional security overlay
  -> overlay inspection view
  -> provenance-labeled effective security metadata and advisory conflicts
  -> downstream human review without writing overlay-applied OpenAPI files

curated provider entry + candidates + classifications + overlays
  -> offline catalog quality report
  -> deterministic findings for missing provenance, stale local verification
     dates, security coverage gaps, duplicate lookup keys, and invalid overlay
     references
  -> maintainer gate without network probes or provider API calls

curated provider entry + classifications + overlays + optional existing cache
artifact registrations
  -> provider advisory report
  -> deterministic metadata summary with spec references, auth status, overlay
     IDs, local artifact paths, and manual follow-ups
  -> downstream authoring/review guidance without fetching documents, applying
     overlays, or creating caches

built-in refreshable spec reference + optional existing cache artifact
registration + saved catalog-openapi-cache file
  -> optional provider-specific source-reviewed refresh correction for selected
     official artifacts, evaluated in memory without replacing raw bytes
  -> refresh review report
  -> deterministic local evidence with registered path, saved path, bytes,
     SHA-256, file status, raw validation evidence, separately labeled corrected
     validation evidence, correction notes, verified date, and manual follow-ups
  -> maintainer review without network access, cache creation, or metadata
     promotion

Selected catalog refreshes persist each successfully validated artifact and
retain those successful report rows if a later reference fails. The CLI
registers completed rows before returning a nonzero partial-failure result, so
an overwritten artifact's SHA-256 and byte count match its cache manifest.
Refresh-result metadata commits as one SQLite transaction. If registration
fails, the CLI restores previously registered artifact files where it captured
a safe backup. Text and JSON failure output include the partial report and
identify the failed provider/spec reference; catalog metadata promotion remains
a separate review action.

built-in provider catalog + optional existing cache registrations + refresh
review report
  -> catalog stats report
  -> provider protocol counts, artifact registry counts, refresh validation
     buckets, and local byte totals
  -> maintainer overview without network access or provider API calls

built-in provider catalog + security classifications + overlays + optional
existing OpenAPI/Swagger cache artifact registrations
  -> catalog security audit report
  -> provider disposition buckets plus artifact security-scheme and
     root/operation security requirement inspection, including partial
     operation coverage diagnostics when root security is absent
  -> maintainer audit without network access, provider API calls, credential
     lookup, or overlay-applied OpenAPI mutation

built-in provider catalog + optional existing cache registrations
  -> provider artifact resolution and materialization
  -> provider-scoped artifact copies under source-aligned directories, separate security-overlay JSON metadata,
     capability labels, and provenance manifests
  -> downstream workflow-directory handoff without downloads, overlay-applied
     OpenAPI mutation, Discovery/Smithy-to-OpenAPI lowering, provider API
     calls, or credential lookup

fnct helper flow:

```text
helper package implementation + FunctionSpec
  -> OpenUdon authoring/review imports names, inputs, and output semantics
  -> udon imports and registers the pure Go helper in its fnct runtime
  -> trusted runtime executes the helper during an approved workflow
```

`apitools` owns only the helper code and descriptor. It does not register the
helper into a runtime, execute workflows, call provider APIs, resolve
credentials, or choose provider accounts.

Google Discovery model flow:

```text
Google Discovery JSON bytes or decoded map
  -> github.com/OpenUdon/googlediscovery.Parse / ParseMap
  -> native Discovery model metadata with service URL, schemas, inherited
     root/resource/method parameters, flattened methods, media upload hints,
     and Google OAuth scopes
  -> downstream protocol-aware generators such as udon native Discovery
     lowering
```

Discovery-to-OpenAPI conversion has been removed. Downstream runtime planning
must consume the native Discovery model instead of derived OpenAPI-shaped
metadata. `apitools/googlediscovery` remains only as a deprecated wrapper over
the standalone parser while consumers migrate. `apitools` does not execute
Google API operations, resolve credentials, fetch tokens, sign requests, or
choose Google accounts or projects.

AWS Smithy model flow:

```text
AWS Smithy JSON bytes or decoded map
  -> github.com/OpenUdon/awssmithy.Parse / ParseMap
  -> native Smithy model metadata with service traits, AWS protocol,
     operation bindings, greedy labels, prefix/query map bindings, static
     protocol fields, shape graph, and AWS service/signing metadata
  -> downstream protocol-aware generators such as udon native Smithy lowering
```

Smithy-to-OpenAPI conversion has been removed. Downstream runtime planning
must consume the native Smithy model instead of derived OpenAPI-shaped
metadata. `apitools/awssmithy` remains only as a deprecated compatibility
wrapper over `github.com/OpenUdon/awssmithy`. `apitools` does not execute AWS
API operations, resolve credentials, fetch tokens, sign requests, or choose
AWS accounts or regions.

OpenRPC model flow:

```text
OpenRPC JSON bytes
  -> openrpc.Parse
  -> native JSON-RPC metadata with document info, servers, methods, params,
     results, errors, components, and stable method selectors
  -> downstream source-aware authoring and review without OpenAPI lowering
```

OpenRPC support is metadata-only. `apitools` does not execute JSON-RPC
methods, contact RPC servers, resolve credentials, create channels, or
implement runtime transport behavior.

GraphQL model flow:

```text
GraphQL SDL, introspection JSON, or operation document bytes
  -> graphql.Parse
  -> native GraphQL metadata with schema roots, types, fields, operations,
     variables, selection names, and stable source operation refs
  -> downstream source-aware authoring and review without OpenAPI lowering
```

GraphQL support is metadata-only. `apitools` does not execute GraphQL
operations, perform live introspection, contact GraphQL endpoints, resolve
variables, choose endpoints, resolve credentials, fetch tokens, or select
accounts/workspaces.

gRPC/protobuf model flow:

```text
.proto files, JSON descriptor artifacts, or binary FileDescriptorSet bytes
  -> grpcproto.Parse
  -> native protobuf metadata with packages, services, methods,
     request/response messages, streaming shape, fields, and stable source refs
  -> downstream source-aware authoring and review without OpenAPI lowering
```

gRPC/protobuf support is metadata-only. `apitools` does not execute RPCs, use
server reflection, open channels, negotiate TLS, contact gRPC endpoints,
resolve credentials, fetch tokens, or select accounts/projects.

OData model flow:

```text
CSDL XML, JSON CSDL, or exported service metadata bytes
  -> odata.Parse
  -> native OData metadata with schemas, entity sets, singletons, navigation
     properties, actions, functions, parameters, query options, and stable
     source refs
  -> downstream source-aware authoring and review without OpenAPI lowering
```

OData support is metadata-only. `apitools` does not call tenant `$metadata`
endpoints, log in to ERP or SaaS tenants, resolve credentials, choose
accounts/tenants, execute `$batch`, execute OData requests, or implement
runtime transport behavior.

Catalog curation artifact flow:

```text
service candidate
  -> try official OpenAPI, Swagger, OpenAPI index, Smithy JSON, Google
     Discovery, AsyncAPI, OpenRPC, GraphQL, gRPC/protobuf, OData, Dropbox
     Stone, or other official machine-readable source
  -> save downloaded review artifact under ignored catalog-openapi-cache/openapi/,
     google-discovery/, aws-smithy/, asyncapi/, openrpc/, graphql/,
     grpc-protobuf/, or odata/
  -> register saved artifact path in catalog-openapi-cache/cache.sqlite
  -> review existing registrations and saved files offline with catalog
     refresh-report before rerunning downloads or promoting metadata
  -> classify parseable OpenAPI/Swagger artifacts that fail strict semantic
     validation as review artifacts, not import-ready provider truth
  -> apply source-reviewed provider-specific validation corrections only for
     explicitly supported official artifacts, while preserving raw downloaded
     artifacts and surfacing correction notes for review
  -> review OpenAPI/Swagger securitySchemes/securityDefinitions plus root and
     operation security requirements for auth/security completeness
  -> expose auth/security for strict-invalid but parseable official artifacts
     through advisory overlays rather than relying on strict import parsing
  -> add catalog security overlay metadata when upstream schemes or
     requirements are missing, ambiguous, stale, internally inconsistent, or
     incomplete for documented protected operations
  -> if no official OpenAPI exists but official API docs or an official
     non-OpenAPI machine spec expose usable endpoint instructions, generate a
     tracked endpoint advisory overlay plus service-specific tracked builder
     under catalog-openapi-cache/
  -> keep endpoint advisory overlays in addition to catalog auth/security
     overlay metadata, not as provider truth
  -> register advisory overlay path in catalog-openapi-cache/cache.sqlite
  -> promote only durable metadata into catalog/
```

All local artifact reads and writes share `internal/artifactio`. Reads require
local relative paths beneath an explicit root, resolve symlinked ancestors of
the selected root once, then reject a symlink at that root and beneath it,
bound bytes, and verify declared size and SHA-256. Single-file
refresh writes use synchronized sibling staging and atomic rename. Provider
materialization and workflow export build a complete sibling directory tree,
reuse byte-identical targets, reject differing collisions by default, and use
a backup/rename/rollback transaction only when force is explicit. Workflow
artifact directories must be non-empty relative paths confined beneath the
workflow root. A failed source, digest, overlay, provider, or manifest stage
publishes no partial tree.

`sqlitecache` schema v3 stores a nanosecond `accessed_at` value for deterministic
LRU retention while preserving second-based `updated_at` timestamps for TTL
behavior. `OpenWithOptions` applies nonzero defaults for query, result,
document, catalog-artifact, report, metadata, and artifact-byte bounds; opening
an existing cache migrates access columns and prunes it immediately. Stores
prune in their transaction, and `Cache.Prune` exposes removal counts. Cached
document and catalog paths are local relative paths under the database
directory. Reads reject path escapes and symlinks and require valid SHA-256 and
positive exact byte counts; inline document rows receive the same checks.

Offline document imports validate URL syntax without DNS or network access,
read integrity-checked cached documents without applying the freshness TTL, and
report the stored timestamp and computed age through the additive
`Client.ImportWithReport` API and CLI output. The existing `Client.Import`
signature and `ImportedSpec` JSON/data shape remain unchanged. Read-write and
refresh modes keep unsafe-host checks before cache use or any network fetch.

No API operation execution is part of either flow.

## Public Contracts

M81.3 supplies a private registered-artifact index parser, not a new direct
parser entry point. It confines and verifies a registration's exact bytes and
SHA-256 before adapting metadata. The OpenAPI index path accepts up to 128 MiB,
4,000,000 JSON tokens/YAML nodes, depth 100, 100,000 operations and 256 MiB
inventory/candidate metadata. A 45-second context deadline is checked during
structural traversal and extraction; decode/normalization phases are
cooperative and this is not a hard process CPU or RSS guarantee. Other source
families retain their 20-MiB direct-parser limit. Shared inventory/candidate
conversion and explicit internal sourceguard limits avoid a second parser or
changed direct-parser defaults. Summary omissions remain diagnostics and
unexamined coverage; a large artifact is not automatically complete evidence.
M81.4 consumes these helpers and preserves that distinction in its index.

M81.2 adds `catalog.RootOptions`/`ResolveRoot` for explicit root, registry and
index path configuration. No root selects no implicit directory; M80
discovery returns metadata-only leads and insufficient evidence. The
selected root rejects symlinks at or below it while resolving symlinked
ancestors. Relative registry/index paths cannot overlap or escape the root.
Resolution creates no files. `sqlitecache.ReadCatalogArtifacts` uses a
read-only SQLite transaction and shared row validation without migration,
pruning or access-time updates, including legacy registry schemas 1–3; newer
schemas refuse. Snapshot bounds are five seconds, 10,000 rows and 32 MiB
aggregate registration text before decoding. Missing registry/index remains
missing evidence. M80 discovery and M81 index generation/read/publication
and selected-artifact export are delivered and published at M80 source
`fb132631c9827eae5f2ec4503d03f21eabfb4113`. The synthetic shared-artifact fixture
is under `testdata/catalog-root` with an explicit disposable preparation recipe.

- Go module path: `github.com/OpenUdon/apitools`.
- The versioned step-candidate wire contract and source-family matrix are
  documented in [`docs/operation-candidates.md`](../../docs/operation-candidates.md).
- Public root package APIs include `Client`, `Search`, `Import`,
  additive `ImportWithReport`, `LocalFiles`, `BuildOperationInventory`, `LoadOperationIndex`,
  `BuildAuthoringAPIDocuments`, auth summaries, and operation selection.
  Step-candidate contracts are additive types in `step_metadata.go`; the
  versioned request carries purpose, typed inputs/outputs, and effect, while
  the result envelope pairs existing operation data with exact source identity,
  evidence, consumer summary, effect, and match diagnostics. These types do
  not change existing operation-summary or selection JSON shapes.
- The historical exported struct field shapes for `ImportedSpec`, `LocalOptions`,
  `LocalResult`, and `DiscoveryCandidate` are preserved. `LocalFiles` uses the
  shared 10,000-entry and 100-candidate defaults, while the existing
  `LocalOptions.MaxBytes` remains configurable.
- `Discoverer.ImportProjectURLs` and `ImportProjectURLsWithReport` retain their
  historical signatures. `ImportProjectURLsReport` exposes bounded attempts,
  diagnostics, and truncation state without breaking existing callers.
- `Search` sources include `apis-guru`, experimental `lap-registry`,
  provider-scoped experimental `rfc9727`, and legacy `public-apis`. `auto`
  uses that order and consults RFC 9727 only when `SearchOptions.ProviderURL`
  is non-empty. Search cache keys include the provider URL.
- CLI entry point: `go run ./cmd/apitools`.
- CLI parsing uses one non-terminating flag adapter. Genuine help exits 0 on
  stdout, usage errors exit 2 on stderr, and runtime failures exit 1 on stderr;
  the parser decides whether `--help` is a flag or a value instead of scanning
  raw arguments.
- Catalog metadata package:
  `github.com/OpenUdon/apitools/catalog`.
- `catalog.ResolvedReference` additively reports selected catalog `SpecKind`
  and protocol. Built-in classifications come from the selected reference;
  explicit `--openapi` inputs are classified by that caller-declared contract,
  without parsing or fetching the document. Text inspect/advisory labels use
  the selected non-OpenAPI protocol rather than implying every source is
  OpenAPI.
- Catalog inspection helpers expose provenance-labeled overlay views without
  mutating or exporting OpenAPI documents.
- Catalog security overlays preserve OpenAPI-style auth semantics for combined
  requirements: schemes in the same requirement set are required together,
  while multiple sets represent alternatives. Legacy singleton requirement
  lists remain for compatibility.
- Catalog quality helpers expose deterministic offline findings with severity,
  provider/candidate/overlay/spec references, fields, and messages.
- Catalog advisory helpers expose deterministic provider summaries with
  spec references, user OpenAPI need, auth/security status, overlay IDs,
  existing artifact paths, protocol/model classification, and manual
  follow-ups.
- Catalog materialization helpers expose provider name resolution, protocol
  capability labels, `uws_source_type` labels, local artifact copy/export reports into
  source-aligned directories, emitted security overlay JSON paths, and provenance manifests for
  existing artifact registry rows. They normalize legacy Google Discovery and AWS Smithy artifact
  paths out of `openapi/`. They are copy-only and do not download missing artifacts, lower
  Discovery or Smithy to OpenAPI, or apply overlays into provider documents.
- Catalog refresh review helpers expose offline saved-artifact reports with
  registered paths, saved paths, byte/SHA-256 evidence, file status, raw
  validation evidence, separately labeled corrected-validation evidence,
  correction notes, protocol/model classification, and deterministic manual
  follow-ups. Corrections validate an in-memory copy; saved bytes, digests, and
  cached metadata remain raw provenance. Parseable-but-strict invalid
  OpenAPI/Swagger artifacts are reported separately from malformed or
  unsupported documents and remain review-only.
- Catalog stats aggregation is reusable through `BuildCatalogStatsReport`.
  The root package owns provider protocol classification order, artifact
  registry kinds, and refresh validation buckets; `cmd/apitools` only parses
  flags and renders the report without network probes.
- Optional bounded cache package: `github.com/OpenUdon/apitools/sqlitecache`,
  with schema-v3 migration, explicit `Options`, reported `Prune` behavior, and
  refresh-result registration through `RegisterCatalogRefreshResults`.
  Cache-relative path containment and refresh metadata mapping remain package
  policy; `cmd/apitools` only selects the report to register.
- Compatibility discovery package:
  `github.com/OpenUdon/apitools/openapidisco`.
- Google Discovery parser package:
  `github.com/OpenUdon/googlediscovery`.
- Deprecated Google Discovery compatibility wrapper:
  `github.com/OpenUdon/apitools/googlediscovery`.
- AWS Smithy parser package:
  `github.com/OpenUdon/awssmithy`.
- Deprecated AWS Smithy compatibility wrapper:
  `github.com/OpenUdon/apitools/awssmithy`.

## Safety Boundary

- Reject non-HTTP(S) remote fetches.
- Reject localhost, private, link-local, multicast, and unspecified hosts by
  default.
- Keep redirect limits, response-size limits, and request timeouts in place.
- Require RFC 9727 catalogs to use `application/linkset+json`, cap inspected
  `service-desc` links, reject unsafe description URLs, and do not recursively
  fetch nested catalogs.
- Treat LAP and RFC 9727 results as unvalidated metadata until the selected
  original document passes normal download and OpenAPI/Swagger validation.
- Resolve symlinked ancestors of an explicitly selected local scan root, then
  reject a symlink at the selected root or beneath it, along with directories,
  special files, and files over the configured size limit.
- Never cache secrets or workflow execution data.
- Never execute API operations, resolve credentials, sign requests, or choose
  runtime accounts.
- Google Discovery parsing is metadata-only and must not perform Google API
  execution, token fetching, request signing, or account selection.
- AWS Smithy parsing is metadata-only and must not perform AWS API
  execution, credential resolution, token fetching, request signing, or account
  and region selection.
- Catalog refresh corrections are explicit per-provider validation aids; they
  must not perform network dereferencing, credential use, provider account
  selection, or broad importer relaxation.

## Risky Change Workflow

For changes to download safety, local file safety, validation, auth summaries,
operation ranking, or provider catalog security overlays:

1. Add or update deterministic tests before broadening behavior.
2. Confirm boundary text in README and memory bank remains accurate.
3. Run `go test ./...`, `go vet ./...`, and `git diff --check`.
4. For exported API changes, run available sibling consumer tests.

### Offline catalog operation index (M81.4)

The root library builds one sanitized operation entry per shared raw identity,
with every provider/spec link and coverage for every catalog reference. It
rechecks the read-only registry snapshot before returning a generation. The
writer validates bindings and complete coverage against the current snapshot,
then atomically replaces only the confined index. Readers enforce schema,
metadata/catalog identities, structural bounds and complete coverage without
reading source artifacts. Snapshot drift removes candidates and marks all
coverage stale/unexamined. This is installation-owned derived metadata, not
signed provider evidence or execution authority. See README for operator
commands and budgets. Missing, unsupported or partially summarized sources
never become definitive negative API evidence.

### Selected catalog artifact provisioning (M81.5)

`CatalogArtifactReference` records catalog/provider/spec/artifact identity,
source family, raw digest/bytes and optional unchanged native selector.
`ExportCatalogArtifacts` joins only selected read-only registrations, verifies
actual source bytes, copies each shared raw family/digest once, and preserves
all selected links in relative-path provenance. It includes only provider-wide
and matching spec-scoped security overlays, preserving auth alternatives
without modifying sources. A read-only registration recheck precedes atomic
directory publication. Invalid binding/digest, unsafe paths, cancellation,
collision or drift publishes no partial tree. Private registration metadata
and URL secrets are omitted; advisory source kind remains labeled. Native
selector validation and approvals belong to the downstream OpenUdon binder.
Legacy provider-level materialization/export APIs and struct shapes are intact.

### Catalog discovery wire foundations (M80.1)

The additive `apitools.catalog-discovery/v1` types reuse StepContract and
OperationCandidate without changing their shapes. Installation-selected root,
registration adapter and remote client remain outside the request wire.
The bounded strict decoder refuses hidden path/endpoint options. Exact-key
validation preserves multiword names and distinguishes nil/open provider scope
from explicit empty scope. Unknown licenses/redistribution and verbatim bounded
notes are source evidence; filters and exclusions stay explicit. Coverage,
qualification counts, native references and tied ranking data preserve the
five-outcome consumer contract. Local retrieval/ranking and explicitly enabled remote lookup are implemented.
Frozen synthetic wire cases are serialization fixtures, not observed lookup
results or source-binding acceptance. See README for supported request limits.

### Local catalog retrieval (M80.2)

`DiscoverCatalogOperations` reads the validated M81 index and current read-only
registration snapshot under explicit installation root/catalog options. It
joins canonical selected provider/spec/artifact links, preserves all shared
references, source authority or unknowns, bounded verbatim license evidence,
relative native identity, auth alternatives and effects. Provider popularity
never upgrades advisory artifacts into official source truth. Authority/license
filters produce named exclusions. Missing configuration/index yields separate
metadata-only leads; invalid identity/configuration blocks, stale coverage
contains no operation candidates, and work limits mark relevant coverage
`work_limit`. All local scope and output work is bounded without parsing raw
sources or fetching URLs. Oversized notes produce evidence gaps; report
projection losses remain incomplete. Discovery qualification now applies the approved strong documented-purpose
criterion and existing typed dimension comparisons. Scores order all candidates
by canonical catalog identity; qualified counts precede display limits. Multiple
qualified native operations remain ambiguous. Complete scope with definitive
conflicts permits a scoped no-match; unexamined scope or indeterminate
comparisons remain insufficient. Required inputs, selected nullable outputs,
unknown requested effects and incomplete auth cannot be hidden by scores.
Sanitation/compaction loss requires source review. The query deadline covers
retrieval, ranking and projection; caller cancellation clears positive results.

### Explicit remote discovery (M80.4)

`DiscoverCatalogOperations` requires both request `RemoteLookup` and installation
`RemoteEnabled` before guarded APIs.guru search/download. Installation endpoint/
HTTP options remain outside request JSON. Copied clients drop source caches,
cookie jars and redirect callbacks; owned host/port, redirect and dial-time
checks remain intact. Eight seconds, three documents and 20 MiB per document
bound lookup; existing parser/prompt limits and a per-source 16-MiB candidate
metadata budget bound projections. Stricter installation byte caps are retained.
Every fetched source retains public-catalog authority, final URL without query/
userinfo/fragment, exact SHA/bytes and ephemeral/provisioning status. Digest/native
duplicates count once. Constrained providers need an exact catalog URL link;
filters exclude named evidence before source fetching. Partial/empty/failing
remote lookup never proves absence. Unsafe configuration/URLs block; caller
cancellation clears positives. No remote bytes are registered or persisted.

Discovery canonicalizes equivalent validated index arrays before bounded
traversal and source-reference selection. Query coverage includes selected
providers without catalog references as unsupported provider-level evidence,
unless an explicit evidence filter excludes them. Critical metadata loss cannot
turn a sanitized absent field into complete negative evidence. Shared groups
above 32 allowed links become leads/`link_limit` coverage before multiplying
references across operations. Ranking reapplies the existing 32-KiB operation
prompt budget after comparison, replacing oversized operations with exact source
leads and incomplete evidence. Wrapped remote timeout diagnostics remain explicit.

### API version check foundations (S05.2)

Additive v1 request/report/hint types and private probe/state helpers support
the approved version-check design. The probe reuses guarded DNS/dial/ports,
adds conditional GET and early HTML refusal, and shares twelve requests,
three API-source body slots and four concurrent fetches. Owned clients discard
cookie jars, caches and redirect callbacks; every source hop must match explicit
official path/repository scope. State/list metadata writes are confined,
atomic and opt-in outside source checkouts. The public discovery orchestration
is supplied by S05.3; these helpers do not execute provider operations.

### Official API version discovery (S05)

DiscoverAPIVersions, VerifyHints and the optional cooperative adapter wrapper
emit additive v1 advisory reports preserving the held baseline. Local evidence
precedes opt-in guarded conditional/sibling/pointer/tree checks; all workers
join before return. Numeric versions, provider scope, source identity and
strict validation remain separate evidence. A lone 304 cannot establish
absence. Request/source/report bounds and every-hop origin checks fail visibly.
Shared-host scopes include repository paths. Catalog x-origin remains a lead.

State and list-cache paths permit bounded metadata persistence outside a
checkout; SaveDir separately permits confined newer Import-valid OpenAPI/
Swagger writes. Discovery remains native metadata and is not saved as OpenAPI.
Digest-matched prebuilt inventory enables operation diffs; absent/partial
inventory remains unexamined without reading the baseline source. Candidate
ranking reuses existing native summaries/effects/auth alternatives. No source
adoption, credential resolution, account selection or operation execution occurs.

### M83 AsyncAPI direction correction

Shape projection interprets AsyncAPI 2.x publish as consumed/output and subscribe as produced/input, corresponding to 3.x receive/send. Native operation IDs/selectors and partial schema/security evidence are retained; this adds no transport or execution support. Equivalent-version regressions live in `operation_shapes_remediation_test.go`.

### M83 GraphQL defaults and response keys

The owned native GraphQL model retains explicit default presence separately from non-null type metadata and each selected field's response key separately from its underlying name. Non-null inputs with defaults are optional. Shape outputs use response keys without changing field/operation selectors; genuine duplicate keys refuse rather than being lost by underlying-name deduplication. Existing summary SelectionNames retain their historical underlying-name meaning. No GraphQL transport is added.

### M83 Smithy member provenance

Smithy URI query literals do not create shape inputs. A same-wire-name native binding is emitted only when the raw modeled input member declares the matching httpQuery trait; a synthetic nonempty MemberName or a same-name body member does not prove query membership. Native selectors and protocol stay unchanged.

### M83 Discovery projection

The bounded raw method/resource tree supplies normal endpoint and request declaration proof independently of the native parser's preferred media-upload path. Upload-only methods cannot supply a normal method/path contract. Requests retain partial schema evidence without asserting mandatory body presence; absent requests emit no body. Only nonempty declared root/base URL evidence may emit a server.

### M83 serialization and format proof

OpenAPI serialization details that the shape contract does not reproduce (including explode, collectionFormat and content/encoding) prevent completeness. Body schemas are unknown for non-JSON or encoded representations; Swagger body proof needs one explicit inherited JSON media declaration. Format knownness permits only exact built-in assertions of the retained JSON Schema compiler, leaving native numeric/binary and custom formats unknown. No runtime serializer or new constraint translation is added.

### M83 bounded OpenAPI compatibility

YAML numeric-key normalization is restricted to integral 100–599 response codes at native operation response maps, including reusable path items; normalized collisions and other non-string keys refuse. Legal path-map x- annotations are skipped regardless of their value type. Swagger detection uses one string/number interpretation throughout shapes, servers and security. TRACE is retained as source shape metadata without changing generic HTTP/runtime methods.

### M83 whole-review GraphQL/YAML proof refinements

GraphQL native selections preserve response keys, underlying fields, deterministic typed argument identity and child selections. Repeated compatible fields/arguments merge across self-aliases, argument/object order and compatible child unions; known conflicts at child response paths refuse. String argument proof uses native lexemes and proper quoted/block-string values, never lossy prompt text. Lexer/balanced delimiter handling charges punctuation only. Child merge validation charges the existing 100,000-node projection budget and retains at most 8 MiB of response-path/argument identity. Full GraphQL validation, fragments, conditional presence and return-type proof remain outside the complete schema claim. Numeric YAML response recognition includes native callback/webhook contexts without adding projection/execution, and bounds lexeme/magnitude before arithmetic.
