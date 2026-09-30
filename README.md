# apitools

`github.com/OpenUdon/apitools` is a Go library and CLI for untrusted API-source
metadata: discovery, URL-safe download, validation, local file scanning,
OpenAPI importing, caching, operation inventories, operation summaries,
auth/security summaries, deterministic operation ranking, catalog protocol
classification, and advisory endpoint overlays.

The module is intentionally narrow. It handles OpenAPI, Swagger, Google
Discovery, AWS Smithy, AsyncAPI, GraphQL, OpenRPC, gRPC/protobuf, and OData
source documents as untrusted data. It does not own workflow semantics, review
handoff contracts, runtime policy, account selection, credential resolution,
request signing, or execution. Downstream products keep those responsibilities
in their own repositories.

## Scope

Keep these responsibilities in `apitools`:

- Discover OpenAPI/Swagger candidates from URLs, APIs.guru, the experimental
  LAP Registry adapter, provider-scoped RFC 9727 catalogs, and legacy
  public-apis probes. Bounded explicit local roots additionally discover and
  validate Google Discovery, AWS Smithy JSON, AsyncAPI, GraphQL, OpenRPC,
  gRPC/protobuf, and OData sources.
- Safely download OpenAPI, Swagger, and catalog-registered API source documents over HTTP(S), with unsafe host
  rejection, redirect limits, request timeouts, and response-size limits.
- Validate OpenAPI 3.0, OpenAPI 3.1, and Swagger 2.0 roots well enough for
  import and authoring workflows.
- Import OpenAPI/Swagger or materialize registered source documents into
  source-aligned `openapi/`, `google-discovery/`, `aws-smithy/`, `asyncapi/`,
  `openrpc/`, `graphql/`, `grpc-protobuf/`, and `odata/` directories with
  deterministic file names.
- Build operation inventories, prompt-safe document summaries, auth/security
  summaries, and request-field summaries.
- Rank operations deterministically from text and structured hints.
- Optionally cache catalog search results and downloaded document bytes.

Keep these responsibilities downstream:

- UWS workflow semantics and public workflow schema live in `../uws`.

`apitools` may describe what an API source document requires. It must not decide
which production account to use, fetch secrets, sign live requests, or execute
operations.

## CLI

```bash
go run ./cmd/apitools --help
go run ./cmd/apitools search --query slack
go run ./cmd/apitools search --query slack --source lap-registry
go run ./cmd/apitools search --query slack --source rfc9727 --provider-url slack.com
go run ./cmd/apitools search --query slack --json
go run ./cmd/apitools search --query slack --cache ~/.cache/apitools/cache.sqlite
go run ./cmd/apitools import --url https://example.com/openapi.yaml --dir ./openapi --name example
go run ./cmd/apitools catalog list
go run ./cmd/apitools catalog resolve gmail openweathermap
go run ./cmd/apitools catalog advisory slack
go run ./cmd/apitools catalog inspect slack
go run ./cmd/apitools catalog materialize slack --out ./api-artifacts
go run ./cmd/apitools catalog export slack --workflow-dir ../openudon/examples/demo
go run ./cmd/apitools catalog overlay-view github
go run ./cmd/apitools catalog security-audit
go run ./cmd/apitools catalog security-report
go run ./cmd/apitools catalog check
go run ./cmd/apitools catalog stats
go run ./cmd/apitools catalog refresh-report
```

CLI help exits 0 and writes usage to stdout. Missing commands, invalid flags,
and other usage failures exit 2 on stderr; runtime failures exit 1 on stderr.
Use `--` when a provider value begins with `-`. A `--help` token consumed as a
flag value or supplied after `--` is data, not a help request.

Search uses APIs.guru first, then the experimental LAP Registry adapter. When
`--provider-url` is present, it next checks that publisher's RFC 9727
`/.well-known/api-catalog`; the legacy public-apis path probe remains the final
compatibility fallback. LAP returns its reported original `source_url`, and
RFC 9727 returns OpenAPI-like `service-desc` links. Both remain visibly
experimental, unvalidated candidates until `apitools import` downloads and
validates the selected document. RFC 9727 discovery fetches one catalog only,
does not recurse into nested catalogs, and rejects unsafe hosts and catalogs
over the service-description link bound. Swagger Catalog and Scalar Registry
are not queried as unauthenticated global catalogs.

Imported documents are treated as untrusted data.
This package never executes API operations, resolves credentials, or exchanges
tokens. OAuth consent and token refresh belong to Udon's trusted runtime; the
operator entry point is `udon oauth google login`. See
[Helper And Gmail Setup](docs/helper-gmail.md) for the pure `gmail.render_raw`
helper and downstream runtime setup.

`search` flags:

```bash
go run ./cmd/apitools search \
  --query "support ticket" \
  --limit 10 \
  --source auto \
  --provider-url api.example.com \
  --public-probe 25 \
  --probe-timeout 5s \
  --probe-budget 30s \
  --cache ~/.cache/apitools/cache.sqlite \
  --cache-mode read-write \
  --json
```

`import` flags:

```bash
go run ./cmd/apitools import \
  --url https://example.com/openapi.yaml \
  --dir ./openapi \
  --name example \
  --cache ~/.cache/apitools/cache.sqlite \
  --cache-mode read-write \
  --json
```

Cache modes are `read-write`, `refresh`, `offline`, and `bypass`. The
`--offline` flag is shorthand for `--cache-mode offline`.
Offline imports validate URL syntax without DNS or network access, accept an
integrity-checked cached document regardless of its age, and report its stored
timestamp and age in the import result. Network-capable modes continue to
reject unsafe hosts before cache use or fetching.

The SQLite backend is bounded by default. Search queries/results, document
rows, and catalog artifact rows use deterministic least-recently-used
retention; cached reports, metadata, and document bytes have independent byte
budgets. Path-backed documents and catalog artifacts must use local relative
paths beneath the database directory, reject symlinks, and match their recorded
SHA-256 and byte count. `sqlitecache.OpenWithOptions` exposes lower explicit
limits, and `Cache.Prune` reports evicted row counts.

Catalog commands expose built-in provider metadata and security-overlay status
without fetching provider documents, executing operations, resolving
credentials, or claiming runtime compatibility:

```bash
go run ./cmd/apitools catalog list
go run ./cmd/apitools catalog resolve gmail openweathermap \
  --cache catalog-openapi-cache/cache.sqlite
go run ./cmd/apitools catalog advisory slack \
  --cache catalog-openapi-cache/cache.sqlite
go run ./cmd/apitools catalog inspect slack
go run ./cmd/apitools catalog inspect slack \
  --openapi ./openapi/slack.yaml \
  --security-overlay ./openapi/slack-security.json
go run ./cmd/apitools catalog overlay-view github --json
go run ./cmd/apitools catalog materialize slack \
  --out ./api-artifacts \
  --cache-dir catalog-openapi-cache \
  --cache catalog-openapi-cache/cache.sqlite
go run ./cmd/apitools catalog export slack openai \
  --workflow-dir ../openudon/examples/demo \
  --cache-dir catalog-openapi-cache \
  --cache catalog-openapi-cache/cache.sqlite
go run ./cmd/apitools catalog security-audit --json
go run ./cmd/apitools catalog security-report --json
go run ./cmd/apitools catalog check --as-of 2026-05-18 --json
go run ./cmd/apitools catalog specs --cache catalog-openapi-cache/cache.sqlite
go run ./cmd/apitools catalog stats \
  --cache-dir catalog-openapi-cache \
  --cache catalog-openapi-cache/cache.sqlite \
  --as-of 2026-05-19
go run ./cmd/apitools catalog refresh-report \
  --cache-dir catalog-openapi-cache \
  --cache catalog-openapi-cache/cache.sqlite \
  --as-of 2026-05-19
go run ./cmd/apitools catalog refresh \
  --provider slack \
  --cache-dir catalog-openapi-cache \
  --cache catalog-openapi-cache/cache.sqlite
```

Reviewed built-in catalog records live in `catalog/data/catalog.json`. After a
catalog-data review, regenerate the checked-in digest/count manifest and lookup
indexes with `go run ./cmd/cataloggen`; CI and local verification use
`go run ./cmd/cataloggen -check` to reject invalid or stale generated output.
The generator is offline and never probes provider URLs.

Catalog resolution is intentionally conservative. Explicit user OpenAPI inputs
and user security overlays take precedence over project-local documents, and
project-local documents take precedence over built-in spec references and
built-in advisory security overlays. Built-in catalog metadata is a discovery
baseline only; it does not override a team's local API contract.

Provider advisory output is a read-only summary for operators and downstream
authoring integrations. `catalog advisory [provider]` combines provider
metadata, spec references, user OpenAPI need, auth/security status, overlay IDs,
source notes, manual follow-ups, optional registered artifact paths, protocol
classifications, and UWS source type labels from an existing `cache.sqlite`.
It does not create a cache when the file is missing, fetch remote documents,
apply overlays, execute API operations, or resolve credentials.
Resolved-reference JSON includes the selected source kind and protocol; text
inspect/advisory headings use the selected non-OpenAPI protocol rather than
labeling a Discovery, Smithy, or documentation source as OpenAPI. Explicit
`--openapi` references are labeled by the caller-declared kind without reading
or validating their contents.

Catalog curation follows a fixed per-service workflow: try official OpenAPI,
Swagger, Google Discovery, and AWS Smithy sources first, review auth/security
completeness, add security-overlay metadata when needed, and, when no official
OpenAPI document exists but official API docs expose usable endpoint
instructions, generate a service-specific docs-derived endpoint overlay asset
in addition to any auth/security overlay metadata. The detailed workflow lives in
[docs/catalog-curation.md](docs/catalog-curation.md). Non-OpenAPI source-family
and protocol-connector boundaries live in
[docs/non-openapi-protocols.md](docs/non-openapi-protocols.md). Downloaded specs
and local SQLite caches stay ignored under `catalog-openapi-cache/`; generated
advisory overlays and service-specific overlay builders are tracked catalog
assets. For catalog curation, `cache.sqlite` records file paths for saved
documents and overlay artifacts rather than duplicating already-saved document
bodies; the local manifest registration program lives under
`catalog-openapi-cache/artifact-registry/`.

Control-plane OpenAPI sources remain metadata-only. Catalog rows such as Docker
Engine may describe local daemon APIs, but `apitools` must not open local
sockets, contact daemon endpoints, encode registry auth headers, or execute
image, container, network, volume, registry, or cluster operations.
Kubernetes is treated the same way: catalog metadata may describe exported
cluster Discovery and OpenAPI artifacts from `/api`, `/apis`, `/openapi/v2`,
and `/openapi/v3`, but `apitools` must not read kubeconfig, service-account
tokens, certificates, or cluster environment variables, and must not contact API
servers.

Enterprise application metadata is also source-first and tenant-safe. NetSuite,
SAP S/4HANA, SAP SuccessFactors, Oracle Fusion Cloud Applications, and Workday
are cataloged from official provider documentation as high-value source
families, but account-specific OpenAPI metadata catalogs, Fusion `/describe`
responses, SAP OData/SOAP artifacts, and Workday WWS/REST metadata must be
user-provided or exported by downstream tooling. `apitools` does not sign in to
enterprise tenants, resolve OAuth clients, fetch tokens, call metadata
endpoints, infer enabled modules, or lower OData/WSDL/SOAP into OpenAPI in this
catalog step.

OpenAPI-first provider rows follow the same metadata-only rule even when a
durable provider-owned spec is available. Adyen, DocuSign, Auth0, Confluence
Cloud, BigCommerce, Cisco Meraki, and Confluent Cloud are cataloged from
official OpenAPI/Swagger documents or provider repositories; Acumatica is
cataloged as tenant-generated Swagger/OpenAPI metadata that must be supplied by
the operator. `apitools` does not create API keys, request OAuth consent, submit
payments, send envelopes, administer networks, inspect stores, choose tenants,
or call streaming/control-plane endpoints.

Overlay inspection views are metadata-only. `catalog overlay-view` reports how
built-in security overlays would supplement catalog classifications, preserving
provenance for schemes and security requirements and surfacing review conflicts
such as duplicate scheme names, missing referenced schemes, overlay-only
additions, and unresolved operation matches. It does not write overlay-applied
OpenAPI files.

Catalog quality checks are offline by default. `catalog check` validates
built-in providers, candidates, source notes, verification dates, security
classifications, and overlay references without probing URLs or downloading
provider documents. Error-level findings return a nonzero exit code; warning-only
reports remain exit code 0 for inspection.

Catalog security audits are offline. `catalog security-audit` classifies each
durable provider by effective auth/security disposition and, when local cache
artifacts are registered, inspects OpenAPI/Swagger `securitySchemes` or
`securityDefinitions` plus root and operation-level `security` requirements.
It reports missing, incomplete, or internally inconsistent artifact security
metadata without fetching provider documents or applying credentials. When
root security is undeclared, at least one operation names an authentication
scheme, and some other operations have no explicit `security` declaration,
the audit reports
`partial-operation-security` and asks maintainers to confirm that operations
without requirements are intentionally anonymous. The audit counts root
requirement alternatives, including an explicit anonymous alternative such as
`security: [{}]`, and separately records whether the root declares a security
array at all (including `security: []`). It also separates operation security
declarations, including explicit anonymous `[]` or `[{}]`, from operations
that name authentication schemes. Only a scheme-bearing operation subset with
other operations lacking an explicit security declaration is reported as
partial; a more specific scheme diagnostic can take precedence while retaining
the partial-coverage follow-up.

Security classifications and overlays are resolved by provider/spec scope.
Classifications remain baseline evidence and a scoped overlay is an explicit
reviewed supplementation for that same scope; input order never determines the
result. Reports expose each scoped disposition, derive `mixed` when distinct
scopes have different effective statuses, and derive `conflict` when overlays
disagree within one scope. Provider resolution uses the selected spec's exact
disposition, then provider-wide evidence, before falling back to the aggregate.
The reviewed built-in catalog rejects conflicting same-scope dispositions at
generation time.

Catalog stats are offline. `catalog stats` summarizes primary provider
protocol classifications, local catalog artifact registry counts by kind, and
refresh artifact validation buckets without probing URLs or executing provider
operations.
The same aggregation is available to library callers through
`apitools.BuildCatalogStatsReport`; the CLI owns only flag parsing and text/JSON
rendering. Repeated advisory and resolver lookups use `catalog.CatalogIndex`,
which validates provider/spec/overlay/security metadata once and returns cloned
records from immutable lookup maps.

Catalog materialization is offline and copy-only. `catalog resolve` reports
provider IDs, protocol capability, registered local artifacts, and security
overlay IDs. JSON reports include protocol classifications and
`uws_source_type` for first-class UWS API sources. `catalog materialize` copies
existing cache-registered artifacts for one provider into source-aligned
directories and emits catalog security overlays as separate JSON metadata with
provenance. Every registered source must be a relative regular file under the
selected cache directory with a matching SHA-256 digest. Materialization and
workflow export stage a complete tree and publish it atomically; identical
destinations are reused, while differing destinations are rejected unless the
operator supplies `--force`, which uses a rollback-capable replacement.
`catalog export` requires `--artifact-dir` to stay beneath the workflow root.
These commands do not download missing files, apply
overlays into OpenAPI documents, lower Discovery or Smithy into OpenAPI,
execute operations, or resolve credentials.

Catalog spec refresh is opt-in and selected. `catalog specs` lists known
machine-readable built-in spec references without network access, optionally
joining registered local artifact paths from `cache.sqlite`. `catalog
refresh-report` reviews built-in refreshable spec references against existing
SQLite registrations and saved files under `catalog-openapi-cache/` without
creating a cache, downloading documents, or promoting metadata. It reports
missing registrations, missing files, SHA-256 and byte evidence, raw validation
status, separately labeled corrected-validation evidence, stale verification
dates, and deterministic manual follow-ups. Provider-specific corrections are
applied only to an in-memory validation copy: the saved file, registered digest,
byte count, and cached metadata remain evidence about the raw downloaded bytes.
It reads saved artifacts with a bounded local file limit and rejects symlinked
artifact paths. `catalog refresh`
downloads only the selected provider/spec reference using the same safe HTTP(S)
download limits as imports, saves the artifact under ignored
`catalog-openapi-cache/openapi/`, `catalog-openapi-cache/google-discovery/`, or
`catalog-openapi-cache/aws-smithy/`, and registers file paths in SQLite instead
of storing duplicate document blobs. Registration and cache-relative path
containment are owned by `sqlitecache.RegisterCatalogRefreshResults`; the CLI
registers completed results before surfacing a later refresh failure. A partial
failure exits nonzero while retaining successful rows in text or JSON output
and naming the failing provider/spec reference, keeping each overwritten
artifact's registered digest and byte count aligned with its saved file. The
SQLite result registration is one transaction; if registration fails, the CLI
restores previously registered artifact files where a safe backup was captured.
Legacy Google Discovery or AWS Smithy cache rows that still point under
`openapi/` are normalized to `google-discovery/` or `aws-smithy/` in catalog
outputs; rerun the artifact registry when accepting those source-aligned paths.
Refresh reports are review inputs only: they do not edit provider metadata,
verified dates, tracked advisory overlays, or security classifications.

### Explicit catalog roots

`catalog.RootOptions` names an operator-supplied directory, a confined relative
registry path (default `cache.sqlite`) and index path (default
`operations.v1.json`). `catalog.ResolveRoot` checks these paths without creating
anything. An empty root means metadata-only leads for the future discovery API;
no sibling checkout or home directory is searched automatically.

`sqlitecache.ReadCatalogArtifacts(ctx, options)` reads existing registrations
without migration, pruning or access-time writes. Missing files remain missing
evidence; invalid/newer registry schemas fail closed. Reads are bounded to five
seconds, 10,000 rows and 32 MiB aggregate registration text. Existing read/write
cache APIs remain unchanged.

For a local root, explicitly refresh selected artifacts with `catalog refresh
--cache-dir ROOT --cache ROOT/cache.sqlite`, then build its index with the
`catalog index --root ROOT` command. Refresh has network effects and is not
part of a discovery call. A redistributable, provider-free root preparation
recipe is in [testdata/catalog-root](testdata/catalog-root/README.md).

`BuildCatalogOperationIndex`, `WriteCatalogOperationIndex` and
`ReadCatalogOperationIndex` take `CatalogIndexOptions` with the explicit root
and installation-selected `sqlitecache.ReadCatalogSpecArtifacts` adapter.
Indexing is offline; requests cannot raise its limits. It verifies raw digests,
indexes shared artifacts once, and records every catalog reference as indexed
or unexamined (missing, digest mismatch, oversize, parse failure, unsupported).
Partial candidate metadata stays positive evidence with incomplete coverage.
The atomic index is deterministic and contains no absolute paths or timestamps.
Readers validate version, catalog identity, coverage, and native source binding
without parsing source files. A changed registration snapshot makes the saved
generation stale and removes its candidates; explicitly rebuild it. A saved
index is installation-owned derived metadata, not a cryptographic signature
or evidence that an operation was approved.

Generation has a three-minute cooperative deadline, 10,000 registrations,
500,000 candidates and a 512-MiB index bound. Reading/publishing applies depth
100 and 32-million JSON-token bounds before typed decoding. These are work and
metadata limits, not a hard process RSS ceiling. Malformed or incomplete
indexes fail closed; they cannot prove that no qualifying API exists.

### Selected artifact export

`ExportCatalogArtifacts(ctx, CatalogArtifactExportOptions{...})` accepts
`CatalogArtifactReference` values from the discovery/index contract and an
installation-selected `CatalogIndexOptions`. It verifies catalog/provider/spec
linkage and exact raw SHA-256/byte count. It exports only selected sources,
provider-wide and matching spec-scoped security overlays, and a relative-path
provenance manifest. Shared bytes are copied once without dropping any selected
provider/operation links. Source bytes and native selectors remain unchanged;
OpenUdon checks selectors during binding. Advisory artifacts remain labeled.
No fetching, credentials, operation execution or implicit overlay application
occurs. Registry private metadata and URL userinfo/query/fragment are excluded.

The workflow directory must exist. `ArtifactDir` defaults to `api-artifacts`,
is confined beneath it, and cannot overlap the catalog source root or use
symlink ancestors. All selected sources and overlays pass before the directory
transaction publishes. Identical output is reused; different output requires
`Force`. Failure preserves any old output without a partial replacement.
Requests are limited to 64 references, 128 MiB per raw source, 512 MiB aggregate
export and a three-minute cooperative deadline. Existing provider-level export
functions and option shapes keep their behavior.

### Catalog discovery contract

M80.1 adds `CatalogDiscoveryRequest` / `CatalogDiscoveryReport` as
`apitools.catalog-discovery/v1`, separate from existing candidate wires.
`DecodeCatalogDiscoveryRequest` bounds JSON to 64 KiB, rejects unknown fields
and trailing data, and never accepts installation paths/endpoints. Local lookup, ranking and explicitly enabled remote lookup are implemented.

Optional provider keys resolve individually by exact catalog key, including
multiword display names. Omitted/null keys mean open scope; an explicit empty
array means empty constrained scope and must not broaden silently. Reports
carry one of match, ambiguous, no_qualifying_api, insufficient_evidence or
blocked, plus coverage, exclusions, unknown license/redistribution evidence,
native references and advisory auth/effect/rank metadata. Raw license notes
are separate bounded source evidence, never inferred permission or instructions.

`DiscoverCatalogOperations(ctx, CatalogDiscoveryOptions{Request: request,
Index: indexOptions})` now retrieves source-backed local metadata with no
source parsing, provider calls or provisioning. It compares typed inputs/outputs, required inputs, documented purpose
and explicitly requested effects. Qualification needs at least two documented
purpose terms covering half the request terms; identity wording alone cannot
qualify. Unknown effects never assert read safety. Missing root/index/registrations yield
metadata-only leads with insufficient evidence; invalid configuration/index
identity yields blocked. Stale entries have no candidates. Exact provider
constraints and evidence filters narrow scope with named exclusions.

Reports preserve relative native references, auth alternatives and effect
metadata. A displayed candidate permits at most 32 source/provider links;
larger shared groups become reference-only leads and `link_limit` coverage
before per-operation reference allocation. Narrow provider constraints recover
supported lookup. Ranked operation metadata retains the existing 32-KiB prompt
bound; larger operations become source leads with incomplete evidence. Work-limited coverage is labeled `work_limit`; report omissions
remain incomplete evidence. Raw license notes above 64 KiB are omitted with
an explicit evidence gap rather than changed and called verbatim. Query
metadata limits do not imply source/parser failure or negative API evidence.

Defaults are 50,000 evaluated operations, 20 returned candidates, 512 KiB
context and a 30-second cooperative local deadline. Ceilings are the index's
500,000 operations, 100 results, 2 MiB context and 60 seconds. Nonzero context
limits must allow at least 4 KiB for a bounded report. At most 32 provider keys
and five authority filters are accepted. Prompt safety and nested typed-field
limits stay unchanged. JSON bounds/deadlines do not claim hard process RSS.
Multiple qualified operations remain ambiguous even with a one-result display
limit. A no-match requires complete scope and definitive conflicts for every
examined operation. Weak wording, omitted fields, incomplete authentication,
nullable selected outputs or unknown constrained effects remain insufficient. Selected-field sanitation
loss never establishes an absent field. A selected provider without catalog
references remains a provider-level lead and unexamined coverage.
Scores rank evidence without selecting or approving operations. Sorting uses
canonical provider/spec/artifact/native identity, never host paths.
Remote lookup needs `request.RemoteLookup` and installation `RemoteEnabled`.
The optional `RemoteClient` and APIs.guru endpoint stay outside request JSON.
Lookup uses guarded APIs.guru search/download only, capped at eight seconds,
three source documents and 20 MiB per document (stricter configured caps remain).
HTTP cookies, redirect callbacks and caches are not inherited as credentials or
persistence. Default host/redirect/dial protection remains enabled; unsafe-host
opt-in is installation-only for reviewed tests. Fetched digests/final URLs remain
public-catalog evidence, with URL userinfo/query/fragment excluded from reports.
Remote bytes are ephemeral and need separately approved registration before local
export. Provider-constrained remote sources require an exact catalog URL link;
public catalog labels cannot broaden the selection. Empty/partial remote results
never establish API absence. No discovery result grants routing or action authority.

### Opt-in local discovery benchmark

Default tests skip the scale benchmark. Select an existing prepared/indexed root
explicitly; no raw source parsing, network refresh or writes occur during queries.
`APITOOLS_CATALOG_BENCHMARK_CATALOG=catalog.json` selects a confined custom catalog
when the index uses one (omit it for the built-in catalog).

```bash
go run testdata/catalog-root/prepare.go /absolute/new/root --operations 12000
go run ./cmd/apitools catalog index --root /absolute/new/root --catalog catalog.json
APITOOLS_CATALOG_BENCHMARK_ROOT=/absolute/new/root \
APITOOLS_CATALOG_BENCHMARK_CATALOG=catalog.json \
go test . -run '^$' -bench '^BenchmarkCatalogDiscoveryLocalRoot$' -benchtime=1x
```

A disposable 12,000-operation synthetic root measured 2.718 seconds/query,
12,000 examined/qualified operations and 95,668 returned JSON bytes on the
qualification host. Peak process RSS was 237,792 KiB; Go reported 844,918,144
cumulative allocated bytes/query. These are observations, not memory or latency
guarantees. Source, registry and index digests remained unchanged. Large real
sources' incomplete summaries remain explicit as documented in the M81 evidence.

## Go Usage

```go
import apitools "github.com/OpenUdon/apitools"

ctx := context.Background()
client := &apitools.Client{}

report, err := client.Search(ctx, apitools.SearchOptions{
	Query:  "slack",
	Source: apitools.SourceAuto,
	Limit:  10,
	// ProviderURL: "slack.com", // enables RFC 9727 after global sources miss
})
_, _ = report, err
```

Local project directories can be scanned without network access:

```go
results, err := apitools.LocalFiles(ctx, apitools.LocalOptions{
	Dir:     "./openapi",
	BaseDir: ".",
	Query:   "slack messages",
	// MaxBytes: 4 * 1024 * 1024, // optional; defaults to apitools.DefaultMaxBytes
})
_, _ = results, err
```

For family-aware authoring discovery, callers provide every file or directory
root explicitly. The report validates OpenAPI/Swagger, Google Discovery, AWS
Smithy, AsyncAPI, GraphQL, OpenRPC, gRPC/protobuf, and OData, and includes
content digests and visible rejections instead of guessing from directory
names:

```go
sources, err := apitools.DiscoverLocalSources(ctx, apitools.LocalSourceDiscoveryOptions{
	Roots: []string{"./api-sources", "./openapi/pets.yaml"},
	Sources: []apitools.LocalSource{
		{Kind: apitools.APISourceKindOpenRPC, ID: "billing", Path: "./metadata/service.json"},
	},
	Query: "fetch pet and publish billing event",
})
_, _ = sources, err
```

Discovery defaults to 10,000 visited entries, 100 accepted candidates, and 20
MiB per file. Reaching either count bound marks the report truncated with a
blocking diagnostic and narrowing guidance. Identical content is deduplicated
by SHA-256. Ambiguous JSON or XML requires an explicit `LocalSource.Kind`;
conventional source directory names affect neither detection nor validation.

The legacy `LocalFiles` helper uses the same bounded walker in OpenAPI-only
mode, while retaining draft-friendly metadata and score ordering. Invalid,
oversized, symlinked, or unreadable entries are skipped so healthy siblings
remain available; reaching its visit or candidate limit returns partial results
with an error. It uses the shared 10,000-entry and 100-candidate defaults;
`LocalOptions.MaxBytes` remains configurable. Use `DiscoverLocalSources` when
structured rejection details or explicit count bounds are needed.

Project-text URL discovery processes at most 16 unique URLs in first-seen order.
If more are present, `ImportProjectURLsReport` sets `Truncated` and adds a
warning, but discovery still returns local and successfully imported URL
candidates. Existing `ImportProjectURLs` and `ImportProjectURLsWithReport`
signatures remain available. Repeated imports reuse same-name files when their
bytes match; candidate reports deduplicate identical documents by SHA-256.
For offline imports, `Client.ImportWithReport` exposes cache timestamp and age
without changing the historical `ImportedSpec` fields or `Client.Import`
signature. The CLI JSON output keeps the imported fields flattened and adds
`cache_stored_at` and `cache_age` when a cached document was used.

Prompt-safe OpenAPI operation context is available through inventories and
document summaries:

```go
inventory, err := apitools.BuildOperationInventory(ctx, apitools.InventoryOptions{
	Documents: []apitools.InventoryDocument{{Path: "openapi/support.yaml"}},
	Query:     "create support ticket",
	// MaxBytes: 4 * 1024 * 1024, // optional for path-backed documents
})
docs, err := apitools.BuildAuthoringAPIDocuments(ctx, apitools.AuthoringAPIDocumentOptions{
	Documents: []apitools.InventoryDocument{{Path: "openapi/support.yaml"}},
	Query:     "create support ticket",
})
_, _ = inventory, docs
```

Operation ranking helpers keep API selection deterministic while leaving
runtime policy downstream:

```go
selection := apitools.SelectOperationByHints(apitools.OperationSelectionHints{
	Provider: "aws",
	Action:   "CreateQueue",
	Purpose:  "create",
}, inventory.Operations)
_ = selection
```

Auth helpers summarize credential and configuration requirements without
resolving secrets or signing requests:

```go
sets := apitools.AuthRequirementSetsForOperation("stripe", inventory.Operations[0])
requestFields := apitools.RequiredRequestFields(inventory.Operations[0])
credentialFieldSets := apitools.CredentialFieldSets(inventory.Operations[0])
_, _, _ = sets, requestFields, credentialFieldSets
```

Security sets preserve OpenAPI semantics: outer entries are OR alternatives,
requirements inside one entry are AND, and an empty entry explicitly permits
anonymous access. Callers must select one alternative before treating its
symbolic credential fields as required; the helpers never flatten alternatives
or resolve credential values.

Provider catalog helpers live in `github.com/OpenUdon/apitools/catalog`:

```go
import "github.com/OpenUdon/apitools/catalog"

resolved, err := catalog.ResolveProvider(catalog.ResolveProviderOptions{
	ProviderKey: "slack",
	UserOpenAPI: "./openapi/slack.yaml",
})
_, _ = resolved, err

securityReport, err := catalog.BuiltInSecurityReport()
_, _ = securityReport, err

view, err := catalog.BuiltInSecurityInspectionView("github")
_, _ = view, err

advisory, err := catalog.BuiltInProviderAdvisoryReport(catalog.ProviderAdvisoryOptions{
	ProviderKey: "slack",
})
_, _ = advisory, err

quality := catalog.BuiltInCatalogQualityReport(catalog.CatalogQualityOptions{})
_ = quality

resolutions, err := catalog.ResolveProvidersWithOptions(catalog.ProviderResolutionOptions{
	ProviderKeys: []string{"slack"},
})
_, _ = resolutions, err

materialized, err := catalog.MaterializeProvider(ctx, catalog.MaterializeOptions{
	ProviderKey:             "slack",
	TargetDir:               "./api-artifacts",
	CacheDir:                "catalog-openapi-cache",
	// Artifacts: existing CatalogSpecArtifact rows with relative Path, SHA256,
	// and Bytes, usually loaded from sqlitecache.
	IncludeSecurityOverlays: true,
})
_, _ = materialized, err
```

Catalog security audits are available from the root package:

```go
audit, err := apitools.BuiltInCatalogSecurityAuditReport(apitools.CatalogSecurityAuditOptions{
	CacheDir: "catalog-openapi-cache",
})
_, _ = audit, err
```

Non-built-in catalog consumers can use the lower-level `Build*` report APIs,
including `catalog.BuildSecurityReport`, `catalog.BuildProviderAdvisoryReport`,
`catalog.BuildCatalogQualityReport`, and
`apitools.BuildCatalogSecurityAuditReport`.

Caching is optional through `github.com/OpenUdon/apitools/sqlitecache` or the
CLI `--cache` flag. Cache modes include `read-write`, `refresh`, `offline`, and
`bypass`.

## Safety Boundary

`apitools` may discover, validate, import, index, and summarize OpenAPI
documents.

It must not:

- Execute API operations or workflows.
- Select production accounts or runtime environments.
- Resolve credentials, acquire tokens, or sign requests.
- Store secrets, workflow execution data, approval state, or runtime state.
- Treat discovered documents as trusted input.

Remote fetches are limited to HTTP(S). Unsafe hosts are rejected by default,
including localhost, private, link-local, multicast, unspecified, carrier-grade
NAT, documentation, benchmarking, NAT64, Teredo, and 6to4 addresses. URL
userinfo is always rejected. Safe clients permit HTTP port 80 and HTTPS port
443; callers may add reviewed ports with `Client.AllowedPorts` without
relaxing address checks. Redirect targets and dial-time DNS results pass the
same checks. Remote bodies are bounded both before and after gzip decoding.
Callers that intentionally need local fixtures can opt out of host and port
checks with `Client.AllowUnsafeHosts`; userinfo and response byte limits remain
enforced.

Local file reads are also fail-closed. Scanners resolve symlinked ancestors of
the caller-selected root once, then reject a symlink at that root or beneath it.
Path-backed `BuildOperationInventory` and `LoadOperationIndex` inputs similarly
resolve ancestor symlinks while rejecting a symlink at the selected file.
Directories, special files, and files larger than the resolved byte limit are
never parsed. Path-backed local reads use bounded I/O;
`LocalOptions.MaxBytes` and `InventoryOptions.MaxBytes` can lower or raise the
limit, and `0` uses `DefaultMaxBytes` (`20 MiB`), matching remote downloads.
In-memory `InventoryDocument.Content` is unchanged. Regular `.json`, `.yaml`,
and `.yml` files that parse but are not OpenAPI or Swagger are still ignored by
`LocalFiles`.

Direct source parsers apply the same untrusted-input contract: 20 MiB source
bytes, depth 100, and bounded structural or semantic work. This includes the
deprecated Smithy and Google Discovery compatibility wrappers; their
already-decoded `ParseMap` entry points bound retained string bytes, nesting,
and traversal work before delegating upstream.

Prompt-facing operation metadata is sanitized before use. The default
`PromptBudget` caps identifiers at 256 runes, text at 2,048 runes, collections
at 32 entries, schema fields at 60 entries, authoring work at 10,000 operations,
each operation at 32 KiB, and a ranked authoring context at 512 KiB.
ANSI/control removal and non-semantic shortening produce visible
warnings. Any compaction that could remove fields or security alternatives,
an oversized selected set, or an aggregate/work-budget violation returns a
blocking `DiagnosticError`; selected operations are never silently discarded.

## Packages

- `github.com/OpenUdon/apitools`: core client, discovery, import, validation,
  inventories, authoring summaries, auth summaries, and operation ranking.
- `github.com/OpenUdon/apitools/catalog`: metadata-only candidate inventory,
  durable provider catalog entries, official spec references, security
  overlays, auth/security reports, overlay inspection views, and provider
  resolution helpers for catalog curation. Catalog metadata is not provider
  runtime compatibility.
- `github.com/OpenUdon/apitools/sqlitecache`: optional bounded SQLite cache
  implementation for the core `Cache` interface, including confined catalog
  refresh-result registration through `RegisterCatalogRefreshResults`.
- `github.com/OpenUdon/apitools/openapidisco`: compatibility wrapper around
  local OpenAPI discovery and primary-candidate selection.
- `github.com/OpenUdon/apitools/googlediscovery`: deprecated compatibility
  wrapper around the standalone `github.com/OpenUdon/googlediscovery` native
  Google Discovery parser.
- `github.com/OpenUdon/apitools/awssmithy`: deprecated compatibility wrapper
  around the standalone `github.com/OpenUdon/awssmithy` native AWS Smithy JSON
  parser.

## Development

```bash
go test ./...
go vet ./...
GOWORK=off go test ./...
GOWORK=off go vet ./...
git diff --check
go run ./cmd/apitools search --help
go run ./cmd/apitools import --help
go run ./cmd/apitools catalog check
go run ./cmd/apitools catalog list
go run ./cmd/apitools catalog resolve slack
go run ./cmd/apitools catalog advisory slack
go run ./cmd/apitools catalog inspect slack
go run ./cmd/apitools catalog materialize --help
go run ./cmd/apitools catalog export --help
go run ./cmd/apitools catalog security-audit
go run ./cmd/apitools catalog security-report
go run ./cmd/apitools catalog stats
```

When changing exported APIs, run dependent checks in sibling consumers when
available:

```bash
(cd ../openudon && go test ./...)
(cd ../udon && go test ./...)
```
