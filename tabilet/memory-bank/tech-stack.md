# Tech Stack

## M91 local qualification checkpoint — 2026-10-25

Exact implementation `0ff4ef561dcc5f784a630b1559ac688bca35802d` in the retained
clean detached qualification worktree passes full invocation1 (30.768s: tests,
vet, CLI/catalog) and affected race1 (89.189s). Go1.26.6 is selected explicitly;
`GOWORK=off GOPROXY=off GOSUMDB=off GOAUTH=off GOTOOLCHAIN=local` and read-only
existing module resolution retain 77 selected / 39 actual compiled modules
including the local owner, no replacements and exact version/sum parity.

Affected staticcheck2026.2.1 invocation1 actually exits1 with six inherited
observations, reconciled against the existing M85 policy. There are no new
helper-file diagnostics, suppression, waiver, unrelated repair or retry.
Historical raw output is absent; no byte-identical-output or clean-global claim
is made. Automatic and one explicit VCS build omit VCS fields despite successful
identical binaries; separate Git cleanliness/revision evidence is retained.

The 19 shared actual-source cases cover business and fake Cloudflare mappings;
eight-family partial shapes remain unresolved. [Local proof](../../docs/m91-local-proof.json)
and [qualification](../../docs/m91-producer-qualification.md) retain logs, counts,
source/corpus and precise limitations. [Live status](status-M91.md) owns whole
review1/evidence blocker. Required ordinary publication, consumer fixtures and
existing-consumer compatibility remain unperformed; no local replacement or
unchanged old consumer pin proves M91 adoption.

## Accepted M85 qualification — 2026-10-09

M85 is accepted/retired after four rows and final whole review3. Published
source `409d4dd54ffa3fbd8c7cddd75bd556ea6bf0e67c` resolves ordinarily as
`v0.0.0-20261009043143-409d4dd54ffa`, sum
`h1:pM+UcJViuxwN0+e2z8FL3IhHrN1rKGOmYeDtNfg/f2Y=` and GoMod sum
`h1:WmUXlfBBoaI6vtv/wzeZfyO8q/pZMhnHiBckkAthYfk=`.
Reviewed source409816b94ef377ab3c7253959169bca08cf1d125 and ordinary evidence
acf9ab9f2afb40e98c0047ae08d0033f5b2395ce are independently observed. [Handoff](../../docs/m85-release-handoff.md),
[ordinary proof](../../docs/m85-ordinary-proof.json) and [full record](../docs/history/status-M85.md)
retain exact source/consumer/review/cleanup evidence. M84 and earlier acceptance
stay frozen; Kinet adoption requires new actual workers and qualification.

Use explicit cached Go1.26.6 executable under GOTOOLCHAIN=local/GOWORK=off;
`/usr/local/go/bin/go` is1.26.1 and its auto-selected version is insufficient
provenance. Actual owner/public/private binary build information is go1.26.6.
Complete owner/current-consumer suites plus ordinary22 events and CLI8 pass;
selected77/73/171, compiled39/35/74, all649 APItools source files and204 unique
selected archives/compiled-generated inputs match. Source archives have different
raw ZIP packaging but identical normalized contents/sums. Three CLI SHA256s:
033868c6f16b2305ad83e19783843a3a25b0233c4a18f95514e56982ba25ee31,
c76757728f6d41c67a14b0f7e52586acce79e6ca02faff607b9b95f9afcf8d81,
0170c0394138ca9853e83eb124e6c4d94a14f8ca3e9c34ec6fbefe8393bdd3be.
Six compatible checker observations retain existing policy; no new waiver.
Compiler-only cleanup preserved all34 referenced generated files/hashes; source,
archive, proof, log and binary caches remain. Rejected d082/quota/old-toolchain
attempts remain recorded and supply no final acceptance.

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

Go1.26.6/GOWORK=off qualifies actual owner/public/private graphs77/73/171,
compiled modules39/35/74 and all15657/18160/35037 selected archive source files.
Current SDK801, UWS root/codec b099 and accepted private runtimeM52 9c remain exact,
with no directory/workspace replacements and the retained Docker-to-Moby version
mapping. Public graphs import no private Udon. Owner/public/private CLI SHA256s:
d221084d1af7bbb81c173103d713fec83d6ed7ea4366cd8163dec941a073f292,
719c1dcb6efe614ff5952b89bd00295efebcc36ea1edd3a05d4843532c303321,
c116ecafb7f910399f0c61d9d8616385a95147cf2cacdaf4bb40c4ac33c74d4a.

Full owner tests/races/vet/build/catalog/format/diff and current consumer proof
pass; all19 ordinary events and621 Git/ZIP/extracted source files match.
The same six checker observations retain the existing policy without a new waiver.
Prior environment/cache failures are retained and resolved only in owned scratch.
Real module directories and complete immutable archive bytes are required for
modverify; broad evidence is reused only after exact source/graph/consumer parity.

## M83 ordinary qualified source

Published f2c5693ec39ad6981693d2ad126ff26e3fdd564c / v0.0.0-20261008061439-f2c5693ec39a pins independently acquired UWS M09 root `v0.0.0-20261008043726-b099f6803277`, source `b099f6803277ae94c7e9f1da0904a0140b278f20`, with exact archive/mod sums recorded in the [retired M83](../docs/history/status-M83.md). Existing dependencies and consumer/browser pins stay fixed. Verification uses retained Go 1.26.6 with GOWORK/GOPROXY off and GOTOOLCHAIN local, full tests/vet/build and root/GraphQL/sourceguard races. Compatible staticcheck 2026.2.1 retains the same six byte-unchanged pre-existing observations. The successor corpus has eight sources/thirteen operations, 9,615 bytes / SHA-256 `cfe13bcd44f62b0b650a20833984231426bd0a83cee91538b96d5332d1107fad`; the original M82 corpus stays frozen. [Qualification](../../docs/m83-qualification.md), [ordinary selected/compiled proof](../../docs/m83-ordinary-proof.json) and [publication evidence](../../docs/m83-publication.md) record actual owner77/39 and public consumer64/33 selected/compiled modules, all files/sums and passing downloaded full checks. Whole review5 and all task rows pass. The preserved candidate SDK failure belongs to already-pending M99.3 declared-output owner/time qualification; it is not waived or reported as passed. No consumer migration, additional publication, source replacement or runtime serializer is implicit; reviewed retirement publication is independently observed before M99 implementation.

## Qualified Stage 11 source metadata — M82

Accepted APItools source `54583f9b2f452b7cc522360c5aeeff29ca22f96c` is public as
`v0.0.0-20261006210844-54583f9b2f45`, review 3. It directly imports accepted UWS
C09 `v0.0.0-20261006181058-6a267306032e` and the existing JSON Schema compiler
`v6.0.1`. The exact closure retains Horizon `v1.14.5` / HashiCorp HCL `v2.24.0`;
no core/HCL dependency removal is claimed. Standalone/offline tests use GOWORK=off,
GOPROXY=off; focused gate: go test -race . -run '^TestOperationShapes'.
Public UWS Marshal/ParseTable is reused without a new serializer or CLI.
The eight-family/twelve-operation fixture remains 9,829 bytes, SHA-256
`dc20d2287322a4b20d5d92b1a0d1ec8036f1bec853b642d7df8797ed096c3ba1`.
Limits: 32 sources, 20 MiB per source/64 MiB total raw, 10,000 operations,
8 MiB table/aggregate projected schemas, 256 KiB per schema, 100,000 projection
nodes total/10,000 per schema/depth 50. No hard in-process CPU/RSS claim.
Required tests/vet/races and consumer checks pass; compatible staticcheck has
six unchanged pre-existing observations, with no suppression/zero-result claim.
Consumer adoption remains with each owning pending milestone. See
[qualification](../../docs/m82-qualification.md).

## Language And Runtime

- **Language/runtime**: Go, using the version declared in `go.mod`.
- **Package type**: library plus CLI.
- **Module path**: `github.com/OpenUdon/apitools`.

## Direct Dependencies

See `go.mod` for exact versions. Key dependency roles:

| Module | Role |
|---|---|
| `github.com/OpenUdon/oas` | OpenAPI 3.x and Swagger parsing, validation, and operation inventory support. |
| `github.com/OpenUdon/awssmithy` | Native AWS Smithy JSON parser used by the deprecated compatibility wrapper and source-type classification tests. |
| `github.com/OpenUdon/googlediscovery` | Native Google Discovery parser used by the deprecated compatibility wrapper and source-type classification tests. |
| `gopkg.in/yaml.v3` | YAML parsing for local and remote specs. |
| `modernc.org/sqlite` | Pure-Go SQLite driver for `sqlitecache` and `catalog` snapshot normalization. |

## Common Commands

```bash
# Test
go test ./...
GOWORK=off go test ./...
go test ./googlediscovery
go test ./awssmithy

# Vet
go vet ./...
GOWORK=off go vet ./...

# Whitespace and patch sanity
git diff --check

# Catalog source validation and generated-output freshness
go run ./cmd/cataloggen -check

# Lint, when available
golangci-lint run ./...
staticcheck ./...

# CLI smoke checks
go run ./cmd/apitools search --help
go run ./cmd/apitools search --query slack --source lap-registry
go run ./cmd/apitools search --query slack --source rfc9727 --provider-url slack.com
go run ./cmd/apitools import --help
go run ./cmd/apitools catalog check
go run ./cmd/apitools catalog advisory slack
go run ./cmd/apitools catalog resolve slack
go run ./cmd/apitools catalog materialize --help
go run ./cmd/apitools catalog export --help
go run ./cmd/apitools catalog specs
go run ./cmd/apitools catalog stats
go run ./cmd/apitools catalog refresh-report
go run ./cmd/apitools catalog refresh --help
go run ./cmd/apitools catalog index --help
go run ./cmd/apitools catalog overlay-view github
go run ./cmd/apitools catalog security-audit
```

When exported APIs change, run dependent checks in OpenUdon and Udon when
available. Ramen is a separate project and is excluded from APItools downstream
verification, per the user's confirmed scope:

```bash
(cd ../openudon && go test ./...)
(cd ../udon && go test ./...)
(cd ../openudon && GOWORK=off go test ./...)
(cd ../udon && GOWORK=off go test ./...)
```

## Runtime Assumptions

- Unit tests should be deterministic and use `httptest` instead of live network
  services.
- Remote downloads are bounded and safe by default: URL, redirect, resolved
  address, and dial-time checks reject unsafe destinations; conventional ports
  are the default allowlist; and gzip wire bytes and decoded bytes are each
  limited to `Client.MaxBytes`.
- Remote search keeps APIs.guru primary, uses LAP Registry's `text/lap`
  targeted-search metadata as an experimental second source, and supports one
  provider-scoped RFC 9727 `application/linkset+json` request when a caller
  supplies a hostname. RFC 9727 discovery defaults to at most 100 inspected
  `service-desc` links and never follows nested catalogs.
- Local scans use bounded reads and reject unsafe path shapes.
- Public-apis list parsing uses the standard library for Markdown and JSON,
  with the shared JSON preflight: 20 MiB response, 10,000 entries, 64 KiB per
  row/JSON entry. Limits, malformed input and cancellation expose a diagnostic
  without probing a partial list; legacy empty JSON arrays remain supported.
- Refresh, materialization, export, and cache artifact paths share the internal
  `artifactio` confinement layer for regular-file reads, digest checks,
  synchronized atomic writes, collision handling, and directory transactions.
- Family-aware local discovery uses existing native parsers plus
  standard-library filesystem, SHA-256, and context primitives. Its defaults
  are 10,000 visited entries, 100 accepted candidates, and
  `DefaultMaxBytes` (20 MiB) per file.
- Step-contract candidates use `BuildOperationCandidates` with explicit local
  bytes or regular files only; it hashes the original bytes, reuses the
  existing family parsers, and never fetches provenance URLs or external refs.
  Defaults are 20 MiB per artifact, 32 sources, 10,000 operations, 100 results,
  and the shared 512 KiB context budget. Family-native selector and auth/effect
  limitations are returned with the report.
- `operationlifecycle` uses only the standard library and root
  `OperationSummary`. Google Discovery `/upload` normalization requires the
  explicit `x-uws-source-kind=google-discovery` extension; generic paths are
  never rewritten by provider-name heuristics.
- `internal/sourceguard` supplies direct-parser defaults of 20 MiB, depth 100,
  100,000 semantic work items, and 1,000,000 linear JSON/YAML structural items.
  `PromptBudget` defaults to 256 identifier runes, 2,048 text runes, 32
  collection items, 60 schema fields, 10,000 operations, 32 KiB per operation,
  and 512 KiB per ranked context.
- Catalog specs, resolution, and materialization expose protocol metadata plus
  `uws_source_type` for first-class UWS API sources. Materialization prefers
  source-aligned output directories: `openapi/` for OpenAPI, Swagger, and
  advisory overlays; `google-discovery/` for Google Discovery; and
  `aws-smithy/` for AWS Smithy JSON; `asyncapi/` for AsyncAPI; `openrpc/` for
  OpenRPC; `graphql/` for GraphQL; `grpc-protobuf/` for gRPC/protobuf; and
  `odata/` for OData.
- OpenRPC parsing uses the Go standard library JSON decoder only. It is
  metadata-only and must not add JSON-RPC transport clients, credential
  resolvers, server probes, or network behavior.
- GraphQL parsing uses only local bytes and the Go standard library. It is
  metadata-only and must not add GraphQL clients, live introspection, endpoint
  probes, variable resolution, credential resolvers, token fetching, or network
  behavior.
- gRPC/protobuf parsing uses only local bytes and the Go standard library. It
  is metadata-only and must not add gRPC clients, server reflection, channel
  creation, TLS negotiation, endpoint probes, credential resolvers, token
  fetching, or network behavior.
- OData parsing uses only local bytes and the Go standard library. It is
  metadata-only and must not add OData clients, tenant `$metadata` calls,
  login flows, endpoint probes, account/tenant selection, credential
  resolvers, token fetching, `$batch`, request execution, or network behavior.
- Helper packages under `helper/` expose pure payload-shaping helpers and public
  function descriptors. They must not import runtime executors, provider
  clients, credential resolvers, token exchange, or network callers.
- The CLI must not require credentials for discovery/import operations.
- SQLite support uses `modernc.org/sqlite`, a pure-Go SQLite driver, so the
  cache package does not require cgo. Schema v3 adds access timestamps for
  deterministic LRU pruning. `sqlitecache.Options` applies bounded defaults for
  row counts plus report, metadata, and artifact bytes; path-backed rows use
  `artifactio` confinement and mandatory digest/byte checks.
  `sqlitecache.RegisterCatalogRefreshResults` owns refresh-result metadata
  mapping and rejects files outside the database directory.
- Every CLI subcommand uses the shared non-terminating flag parser. Help is
  stdout/exit 0, usage diagnostics are stderr/exit 2, and runtime failures are
  stderr/exit 1.

## Advisory Overlay Curation

- Advisory OpenAPI overlays are generated review artifacts, not provider-
  official OpenAPI documents. Update the overlay builder and regenerate the
  JSON artifact together.
- When official docs referenced by an overlay link to other official API docs,
  inspect those linked docs before closing the overlay review.
- Add linked operations to the same advisory overlay when they are
  provider-owned, API-shaped, documented without login, and directly support
  the same workflow or required lookup path.
- Preserve useful provider parameter descriptions and constraints so downstream
  LLM planning can understand data-pipeline inputs, accepted formats, lookup
  alternatives, and related helper APIs.
- Add every linked doc that contributes operations, schemas, parameters, or
  auth/security metadata to `x-apitools-overlay.source_refs`.
- Record a no-add decision when linked docs are conceptual, pricing-only,
  duplicated, authenticated-only, unstable, or outside the reviewed workflow.
- Do not call provider APIs, fetch tokens, resolve credentials, scrape
  authenticated docs, or label advisory overlays as official OpenAPI.

## Harnesses

M81.3 keeps all existing direct parser and `BuildOperationCandidates` limits at
20 MiB. Focused registered-index checks are `go test . ./internal/sourceguard
-run 'Test(RegisteredCatalogArtifactLargeAndIdentity|CatalogOpenAPIStructuralRefusal|ExplicitIndexLimitsRemainBounded)'`.
The explicitly invoked local-only measurement is
`APITOOLS_CATALOG_MEASURE_ROOT=/absolute/registered/root go test . -run
TestCatalogRegisteredArtifactMeasurement -v -count=1`; it reads existing
Cloudflare/Graph registrations and source bytes without network or provider
credentials. Default tests skip it. The 2026-09-30 runs measured 5.727 seconds /
430,364 KiB peak process RSS for Cloudflare (21,930,908 bytes), and 40.111
seconds / 870,608 KiB for Graph (37,110,274 bytes). GNU time measured the
largest process in each Go test invocation, not a guaranteed RSS ceiling.
Both exposed incomplete operation summaries rather than suppressing gaps.

M81.2's focused root/registration gate is `go test ./catalog ./sqlitecache`.
`go run testdata/catalog-root/prepare.go /absolute/new/directory` creates an
explicit new disposable synthetic root without fetching documents. No new
module dependency is used. The offline operation index command is delivered by M81.4.

| Harness | Command | What it proves | Requirements |
|---|---|---|---|
| Step-candidate contract and source-adapter tests | `go test . -run 'Test(StepMetadata|EffectAndContract|OperationSourceReport|ContractRequiredness|OperationCandidateV1|SummarizeOperationForConsumer|AssessOperationEffect|RankOperationCandidates|BuildOperationCandidates|SanitizeAdapterCandidates|NativeSchemaFields|NativeOperationAdapter)'` | Versioned wire fixtures validated against producer ranking, source-backed summaries and effect rationale, typed contract matching, all eight local source adapters, exact digest/native selector preservation, auth alternatives, truncation blockers, cancellation, and source safety. | Go toolchain. |
| Unit tests | `go test ./...` | Core discovery, validation, import, inventory, auth, ranking, and cache behavior. | Go toolchain. |
| Lifecycle ranking tests | `go test ./operationlifecycle` | Same-source role scoring, ambiguity, API-first compatibility, and explicit Discovery provenance. | Go toolchain. |
| Helper tests | `go test ./helper/...` | Public helper descriptors and pure Gmail raw-message encoding without credential or network behavior. | Go toolchain. |
| Google Discovery wrapper tests | `go test ./googlediscovery` | The deprecated `apitools/googlediscovery` wrapper still delegates to the standalone parser and preserves compatibility fixtures. | Go toolchain and the versioned standalone module from `go.mod`. |
| AWS Smithy wrapper tests | `go test ./awssmithy` | The deprecated `apitools/awssmithy` wrapper still delegates to the standalone parser and preserves compatibility fixtures. | Go toolchain and the versioned standalone module from `go.mod`. |
| OpenRPC parser tests | `go test ./openrpc` | OpenRPC documents parse into deterministic metadata, selectors, summaries, and malformed-input failures without runtime behavior. | Go toolchain. |
| GraphQL parser tests | `go test ./graphql` | GraphQL SDL, introspection JSON, and operation documents parse into deterministic metadata, selectors, summaries, and malformed-input failures without runtime behavior. | Go toolchain. |
| gRPC/protobuf parser tests | `go test ./grpcproto` | `.proto`, JSON descriptor, and binary descriptor-set artifacts parse into deterministic metadata, selectors, summaries, and malformed-input failures without runtime behavior. | Go toolchain. |
| OData parser tests | `go test ./odata` | CSDL XML, JSON CSDL, and exported service metadata parse into deterministic metadata, selectors, summaries, and malformed-input failures without runtime behavior. | Go toolchain. |
| Vet | `go vet ./...` | Basic static correctness. | Go toolchain. |
| Lint | `golangci-lint run ./...` | Extended static linting when the tool is installed. | `golangci-lint`. |
| Staticcheck | `staticcheck ./...` | Whole-module static analysis, including unused code and simplification diagnostics. | `staticcheck`. |
| Catalog generation | `go run ./cmd/cataloggen -check` | The reviewed catalog JSON, golden digest/count manifest, and checked-in indexes are valid, deterministic, and synchronized. | Go toolchain. |
| Catalog quality | `go run ./cmd/apitools catalog check` | Built-in catalog metadata passes offline quality gates with no error-level findings. | Go toolchain. |
| CLI smoke | `go run ./cmd/apitools search --help`, `go run ./cmd/apitools import --help`, `go run ./cmd/apitools catalog advisory slack`, `go run ./cmd/apitools catalog resolve slack`, `go run ./cmd/apitools catalog materialize --help`, `go run ./cmd/apitools catalog export --help`, `go run ./cmd/apitools catalog specs`, `go run ./cmd/apitools catalog stats`, `go run ./cmd/apitools catalog refresh-report`, `go run ./cmd/apitools catalog overlay-view github`, and `go run ./cmd/apitools catalog security-audit` | CLI wiring and the stdout/stderr/0/1/2 policy remain available, including documented LAP/RFC 9727 source flags. | Go toolchain. |
| Memory-bank API runner | `../skills/harness/tackle-memory-bank-api-loop .` | Discovers every canonical lane file and either reports no actionable rows without an API call or selects work across lanes using the milestone priority. | Python 3 and `LLM_MODEL`; API credentials and explicit host-shell acknowledgement only when actionable work is executed. |
| Patch check | `git diff --check` | No whitespace errors. | Git. |

## Verification and testing standards

- Preserve the existing verification commands, milestone-specific gates and
  explicit limits. Standard testing adds no light-profile time caps.
- Run focused checks during task work and the required full regression once
  per unchanged closure candidate. Repeat affected checks after fixes, failures
  or unresolved coverage; separate compilation time from test execution.
  Required race and other acceptance gates remain required.
- Preserve meaningful success and failure coverage through a replacement
  implementation before removing the old implementation's tests. Remove tests
  only when their behavior has moved to its successor or has been removed.
- Acceptance readiness: identify missing final-acceptance inputs, capabilities
  and authority early. Before each counted expensive attempt, run a fast
  preflight of the applicable fixtures, schemas, commands and launchers.
  Missing real acceptance evidence remains missing; mocks do not replace it.
  External testing-tool repair is optional unless explicitly promoted by the
  owner under the milestone-family rules.
- A human-approved standing verification budget follows the selected goal
  protocol and the [authorization rules](milestone.md#authorization-requirements).
  Persist invocation usage and respect remaining attempts, stop conditions and
  all smaller tool or authorization caps. These standards activate no grant.

## Memory-bank API runner compatibility

The separately installed `../skills/harness/tackle-memory-bank-api-loop`
recognizes root and child filenames, the six task markers, and literal
retirement envelopes. Its current specification parser assumes level-two
milestone headings. It rejects the frozen M87, M88, M89, S07, S08, S09 and S10
records because their retained specifications use level-three headings, and it
does not extract M86's active specification for the same reason.

Automated API-runner execution and retirement compatibility remain blocked
until a separately authorized runner repair/check supports this harness's
actual heading structure. Preserve the existing project retirement rules,
frozen records and M86 specification; do not normalize history to satisfy the
tool. This rule upgrade neither changes the external runner nor launches it.

## CI And Tooling

- Public repo guidance lives in README.
- Planning and milestone state lives in tracked `tabilet/memory-bank/` files in
  this repository.
- Status files use `status-<LANE><NN>.md`; task tables put a supported
  backticked state in the second column so the account-level runner can count
  actionable and blocked rows deterministically.
- Keep generated or cached files out of the public repository unless they are
  deliberate fixtures.

- The offline `catalog index --root DIRECTORY` command uses the read-only SQLite
  registration adapter, registered-artifact parsers and atomic artifact I/O.
  Generation is bounded to three cooperative minutes, 10,000 registrations,
  500,000 candidates and 512 MiB serialized metadata. Reader/writer JSON
  preflight uses depth 100 and 32 million tokens; these limits do not claim
  a hard process RSS ceiling. Unknown coverage remains explicit.

- Selected catalog export reuses `artifactio` confined reads and directory
  transactions. Bounds: 64 references, 128 MiB per raw source, 512 MiB total,
  three-minute cooperative deadline. No index parser, network call or new
  dependency is used by export.

- Catalog discovery's additive strict request decoder accepts at most 64 KiB
  and no installation paths/endpoints. Default local limits: 50,000 operations,
  20 results, 512 KiB report, 30 cooperative seconds; ceilings: 500,000 operations,
  100 results, 2 MiB, 60 seconds; nonzero context floor 4 KiB. Provider and
  authority request arrays are bounded to 32 and 5; existing typed field and
  prompt-safety limits are reused. The contract is delivered; query behavior
  now supports local qualification and opt-in guarded APIs.guru lookup.
  Remote lookup is capped at eight seconds, three documents and 20 MiB each;
  installation configuration stays outside request JSON, and source bytes remain
  ephemeral. Default checks use only local HTTP fixtures.

- Local discovery reads only the saved index and read-only SQLite registrations;
  no raw-source parsing or URL fetching occurs. Report projection is bounded
  independently of evaluated-operation work. Raw license notes are retained up
  to 64 KiB, otherwise emitted as a named evidence gap. Query work limits use
  `work_limit` report coverage; stale registrations remove indexed candidates.

- The explicitly selected `BenchmarkCatalogDiscoveryLocalRoot` skips unless
  `APITOOLS_CATALOG_BENCHMARK_ROOT` names a populated/indexed root. Optional
  `APITOOLS_CATALOG_BENCHMARK_CATALOG` is one confined root-relative filename.
  The synthetic preparation recipe supports `--operations 12000` (1–100,000
  generated operations) in a new directory; it updates exact registrations.
  The benchmark reads metadata only. The 12,000-operation qualification run
  measured 2.718 seconds/query and 237,792 KiB peak process RSS; observations
  are not hard performance/memory bounds. Source/registry/index digests were
  unchanged. README has the exact command and allocation/report-size evidence.

- Discovery applies its 32-link shared-source limit before per-operation
  reference allocation, emits `link_limit` coverage and supports narrowing via
  exact provider constraints. The existing 32-KiB operation prompt budget is
  reapplied after typed rank evidence; oversized operations stay source leads.
  Registered source/index parser ceilings are unchanged by these query limits.

- API version probe defaults are 5 s total, 3 s per request, four concurrent
  requests, twelve HTTP requests and three source bodies at 8 MiB each.
  JSON state/list helpers use artifactio and the shared structural checks;
  there is no new dependency or SQLite migration. Focused foundations gate:
  `go test . -run '^TestAPIVersion'` and its `-race` variant.

- Delivered API version CLI: `go run ./cmd/apitools versions --help`.
  `versions --request FILE` uses local evidence by default; `--network`
  enables bounded lookup. State/list-cache/save/inventory flags are independent.
  Library network tests use local servers; focused gate is
  `go test -race . ./cmd/apitools -run 'TestAPIVersion|TestVersionsCLI'`.

S07 focused lifecycle gate: `go test ./catalog ./sqlitecache ./cmd/apitools
./internal/artifactio`. Root-aware writable catalog registration uses the
additive `sqlitecache.OpenCatalog`; existing cache openers and schemas remain.
Omitted CLI registry/index flags honor APITOOLS_CATALOG_*, then CATALOG_*,
then KINET_CATALOG_*; explicit flags win and empty path flags select defaults.

Consumer compatibility qualification can use immutable source copies plus
temporary modfiles pointing APItools at the integration candidate. Propagate an
absolute `GOFLAGS=-modfile=...` to nested `GOWORK=off` subprocesses; never edit
consumer pins for a producer check. Compare original consumer heads, status and
manifest hashes afterward. S07's commands and outcomes are in
[qualification](../../docs/s07-qualification.md).
