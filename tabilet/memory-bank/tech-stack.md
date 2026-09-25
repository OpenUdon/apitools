# Tech Stack

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
| `modernc.org/sqlite` | Optional pure-Go SQLite cache backend in `sqlitecache`. |

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
go run ./cmd/apitools catalog overlay-view github
go run ./cmd/apitools catalog security-audit
```

When exported APIs change, run dependent checks in sibling consumers when
available:

```bash
(cd ../openudon && go test ./...)
(cd ../ramen && go test ./...)
(cd ../udon && go test ./...)
(cd ../openudon && GOWORK=off go test ./...)
(cd ../ramen && GOWORK=off go test ./...)
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
- Refresh, materialization, export, and cache artifact paths share the internal
  `artifactio` confinement layer for regular-file reads, digest checks,
  synchronized atomic writes, collision handling, and directory transactions.
- Family-aware local discovery uses existing native parsers plus
  standard-library filesystem, SHA-256, and context primitives. Its defaults
  are 10,000 visited entries, 100 accepted candidates, and
  `DefaultMaxBytes` (20 MiB) per file.
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

| Harness | Command | What it proves | Requirements |
|---|---|---|---|
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

## CI And Tooling

- Public repo guidance lives in README.
- Planning and milestone state lives in tracked `tabilet/memory-bank/` files in
  this repository.
- Status files use `status-<LANE><NN>.md`; task tables put a supported
  backticked state in the second column so the account-level runner can count
  actionable and blocked rows deterministically.
- Keep generated or cached files out of the public repository unless they are
  deliberate fixtures.
