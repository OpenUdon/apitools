# Product

`apitools` is a Go library and CLI for OpenAPI/Swagger document tooling and
provider API-source metadata. It helps downstream OpenUdon projects find
candidate API documents, import OpenAPI-compatible documents safely, summarize
operations and security requirements, rank likely operations, and classify
official server-side API description sources without crossing into runtime
execution.

## Users

- OpenUdon maintainers who need deterministic OpenAPI operation metadata for
  authoring, review, package evidence, and credential handoff.
- Runtime/compiler maintainers who need prompt-safe operation inventories and
  auth summaries without depending on product workflow behavior.
- Agents and reviewers comparing API discovery results, imported specs,
  operation indexes, and security summaries.
- Future provider-catalog maintainers curating popular service metadata,
  official spec references, and security overlays.

## Core Concepts

- **Explicit catalog root**: operator configuration identifying the local
  artifact registry and deterministic offline operation index. Metadata reads never select a
  sibling checkout implicitly and never migrate, prune or update registrations.
  No-root discovery remains pending M80 and must report insufficient evidence.

- **OpenAPI document**: an OpenAPI 3.x or Swagger 2.0 document loaded from a
  URL, public catalog, cache, or local file.
- **Spec protocol classification**: the API description protocol or model
  family represented by a catalog reference, such as OpenAPI, Swagger, Smithy,
  Google Discovery, AsyncAPI, OpenRPC, GraphQL, gRPC/protobuf, OData, Dropbox
  Stone, OpenAPI index, or human docs. It is separate from the catalog source
  reference kind and maps
  executable first-class protocols to UWS source types: OpenAPI/Swagger to
  `openapi`, Google Discovery to `google-discovery`, AWS Smithy JSON to
  `aws-smithy`, AsyncAPI to `asyncapi`, OpenRPC to `openrpc`, GraphQL to
  `graphql`, gRPC/protobuf to `grpc-protobuf`, and OData to `odata`.
- **Discovery candidate**: a scored local or remote document candidate with
  source metadata and validation evidence.
- **Remote catalog discovery**: bounded, metadata-only lookup that keeps
  APIs.guru as the primary global directory, treats LAP Registry as an
  experimental secondary directory, and checks an RFC 9727 publisher catalog
  only when the caller supplies a provider URL or hostname. Catalog results
  remain unvalidated until the original source document is downloaded through
  the normal import validator.
- **Operation inventory**: deterministic summaries of documents, operations,
  parameters, request bodies, response fields, and security requirements.
- **Step-candidate metadata contract**: a versioned, source-qualified record
  that pairs an existing operation summary with a consumer description, an
  evidence-bearing read/write/unknown effect assessment, and an advisory
  comparison against a step's purpose, inputs, outputs, and effect. The local
  producer binds each candidate to raw-content SHA-256 and a family-native
  selector, carries response-field nullability through selected schema
  ancestors, exposes source-family limitations, and never fetches provenance
  URLs or external references. An unknown leading action remains unknown even
  when later wording contains a recognized read/write token. It carries
  metadata only; it does not bind, approve, or execute an operation.
- **Operation lifecycle ranking**: conservative same-source sibling selection
  for read/update/delete roles over `OperationSummary`, with explicit source
  provenance controlling provider-specific path normalization.
- **Prompt-safety report**: sanitized and budgeted operation or ranked-context
  metadata with visible warnings for non-semantic shortening and blocking
  diagnostics whenever compaction could alter fields, selected operations, or
  security alternatives.
- **Auth summary**: prompt-safe credential and configuration metadata derived
  from OpenAPI security schemes and known dialect hints.
- **Provider catalog**: curated metadata for popular services, official or
  machine-readable spec sources, local aliases, and security overlays.
- **Security overlay**: explicit supplemental auth/security metadata kept
  separate from upstream specs when provider documents are incomplete.
- **Advisory endpoint overlay**: a service-specific OpenAPI-shaped advisory
  artifact generated from official human API docs or an official non-OpenAPI
  machine-readable source when no official OpenAPI document is available but
  endpoint instructions are sufficient for a useful reviewed subset. It
  supplements, but does not replace, auth/security overlay metadata.
- **Overlay inspection view**: a read-only, provenance-labeled view showing how
  security overlays would supplement catalog security classifications, including
  advisory conflicts for human review.
- **Catalog quality report**: a deterministic offline report for catalog
  metadata regressions such as stale verification dates, missing source notes,
  missing security status, duplicate lookup keys, and invalid overlay
  references.
- **Provider advisory report**: a deterministic metadata-only provider summary
  combining catalog references, provider/spec-scoped auth/security
  dispositions, aggregate mixed/conflict status, overlay IDs, optional
  registered artifact paths, and manual follow-ups for operator review.
- **Provider artifact materialization**: an offline copy/export workflow that
  resolves provider names, reports protocol capability, copies existing
  cache-registered artifacts into provider-scoped target directories, emits
  security overlays as separate JSON metadata, and writes provenance manifests.
- **Fnct helper catalog**: public, pure Go helper contracts for shaping API
  payloads before trusted runtime execution. Helpers may describe names,
  request fields, output semantics, and side-effect posture, but do not execute
  provider APIs.
- **Catalog refresh review report**: an offline, read-only report that joins
  built-in refreshable spec references with existing SQLite registrations and
  saved cache files, including byte/SHA-256 evidence, distinct raw and corrected
  validation evidence, and manual follow-ups.
- **Catalog statistics report**: reusable offline library aggregation for
  provider protocol families, registered artifact kinds, and refresh evidence;
  CLI commands render this package-owned policy rather than reimplementing it.
- **Catalog refresh correction**: a source-reviewed, provider-specific
  validation aid for official OpenAPI/Swagger artifacts whose current upstream
  shape needs narrow local normalization, such as bundling known sibling
  schemas or removing invalid extension-only refs. Corrections are reported as
  review metadata and do not execute API operations or resolve external refs at
  runtime.
- **Advisory overlay artifact**: a tracked OpenAPI-shaped review artifact for
  services that have official documentation or official non-OpenAPI machine
  specs but no recorded official OpenAPI document. It may include endpoint and
  security metadata, but it is not provider truth.
- **Google Discovery parsing**: a bounded, deterministic parser for official
  Google Discovery documents that preserves service metadata, schemas,
  inherited parameters, methods, media upload hints, and OAuth scopes for
  downstream protocol-aware generators without producing OpenAPI-shaped
  metadata. The parser source of truth is the standalone
  `github.com/OpenUdon/googlediscovery` module; `apitools` keeps only a
  temporary compatibility wrapper.
- **AWS Smithy model parsing**: bounded, deterministic parsing of official AWS
  Smithy JSON service models through the standalone
  `github.com/OpenUdon/awssmithy` parser, preserving native AWS protocol
  metadata for downstream generators without producing OpenAPI-shaped
  metadata.
- **OpenRPC parsing**: bounded, deterministic parsing of provider-owned
  OpenRPC JSON documents, preserving JSON-RPC method contracts, params,
  results, errors, components, and local selectors without contacting servers
  or producing OpenAPI-shaped metadata.
- **GraphQL parsing**: bounded, deterministic parsing of local/provider-owned
  GraphQL SDL, introspection JSON, and operation documents, preserving type,
  field, operation, variable, selector, and source-ref metadata without live
  introspection, endpoint calls, or OpenAPI-shaped lowering.
- **gRPC/protobuf parsing**: bounded, deterministic parsing of
  local/provider-owned `.proto` files, JSON descriptor artifacts, and binary
  `FileDescriptorSet` bytes, preserving package, service, method, message,
  field, streaming, selector, and source-ref metadata without server
  reflection, channel creation, RPC execution, or OpenAPI-shaped lowering.
- **OData parsing**: bounded, deterministic parsing of local/provider-owned
  CSDL XML, JSON CSDL, and exported service metadata, preserving entity set,
  singleton, navigation, action, function, parameter, query-option, selector,
  and source-ref metadata without tenant metadata calls, login, credential
  lookup, `$batch`, request execution, or OpenAPI-shaped lowering.

## Primary Workflows

1. Search public sources for an API document matching a service or intent.
   APIs.guru remains primary; LAP Registry is an experimental second source;
   an explicit provider hostname enables one non-recursive RFC 9727 lookup;
   public-apis probing remains a compatibility fallback.
2. Safely download, validate, and import a selected OpenAPI/Swagger document.
3. Scan local project files for API documents without network access.
   Family-aware scans require explicit file or directory roots, validate all
   executable source families, deduplicate identical content by digest, and
   expose ambiguity, rejection, and truncation evidence.
4. Build operation inventories and authoring summaries for downstream prompts.
5. Rank operations from structured hints and text while preserving ambiguity.
   Downstream authoring may also expand one selected operation into a
   conservative same-source lifecycle set without executing any operation.
6. Summarize auth/security requirements without resolving secrets.
7. Curate provider catalog entries and security overlays for common services.
8. Inspect overlay effects and conflicts without writing overlay-applied
   OpenAPI documents.
9. Generate tracked service-specific endpoint overlay artifacts and
   per-service builders, in addition to auth/security overlay metadata, when a
   provider has no official OpenAPI document but official docs or an official
   non-OpenAPI machine spec support a useful reviewed endpoint subset.
10. Run offline catalog quality gates before promoting catalog metadata.
11. Render provider advisory summaries for downstream authoring integrations
    without making catalog metadata authoritative.
12. Review saved refresh artifacts offline before rerunning downloads or
    promoting catalog metadata.
13. Classify catalog spec references by protocol/model family so official
    non-OpenAPI server-side API descriptions can be reviewed without pretending
    they are import-ready OpenAPI documents.
14. Apply narrow source-reviewed refresh corrections to selected official
    OpenAPI/Swagger artifacts during catalog refresh validation, with
    correction notes surfaced for human review.
15. Parse Google Discovery documents into native Discovery metadata for
    downstream protocol-aware planning without using Discovery-to-OpenAPI
    conversion as an execution contract.
16. Parse AWS Smithy JSON service models into native Smithy metadata for
    downstream protocol-aware planning without using Smithy-to-OpenAPI
    conversion as an execution contract.
17. Resolve and materialize provider artifacts from existing local catalog
    registrations into provider- or workflow-scoped directories with
    provenance. OpenAPI/Swagger artifacts materialize under `openapi/`, Google
    Discovery under `google-discovery/`, AWS Smithy JSON under `aws-smithy/`,
    AsyncAPI under `asyncapi/`, OpenRPC under `openrpc/`, GraphQL under
    `graphql/`, gRPC/protobuf under `grpc-protobuf/`, OData under `odata/`,
    and advisory overlays remain OpenAPI-shaped review artifacts under
    `openapi/`.
18. Import pure helper contracts, such as Gmail raw-message rendering, from
    `helper/...` packages so downstream authoring tools and runtimes can agree
    on function names and request/output shapes without moving execution into
    `apitools`.

## Scope

- OpenAPI/Swagger discovery, download, validation, import, local scanning, and
  optional caching.
- Bounded remote OpenAPI-candidate lookup through APIs.guru, LAP Registry's
  targeted search metadata, and provider-scoped RFC 9727 Linkset catalogs.
  LAP and RFC 9727 candidates expose provenance and remain experimental and
  unvalidated until import; RFC 9727 discovery does not recursively crawl
  nested catalogs.
- Bounded local multi-family discovery for OpenAPI/Swagger, Google Discovery,
  AWS Smithy, AsyncAPI, GraphQL, OpenRPC, gRPC/protobuf, and OData, with
  explicit roots, validated metadata, SHA-256 provenance, ambiguity reporting,
  and no remote lookup.
- Operation indexing, request/response/security summaries, and deterministic
  operation selection.
- Provenance-aware operation lifecycle sibling ranking over prompt-safe
  summaries, without workflow semantics or runtime execution.
- Prompt-safe metadata for downstream authoring and review workflows.
- Provider catalog metadata and supplemental security overlays for common
  services, as metadata only.
- Provenance-labeled inspection views for catalog classifications and overlay
  effects, including deterministic conflict reporting.
- Static offline catalog quality checks for source/provenance notes, stale
  verification dates, security status coverage, overlay references, duplicate
  lookup keys, and user OpenAPI need consistency.
- Ignored local curation artifacts used to review official specs, Discovery
  documents, and cache entries before durable metadata is promoted.
- Catalog protocol classification for OpenAPI, Swagger, Smithy JSON, Google
  Discovery, AsyncAPI, OpenRPC, GraphQL, gRPC/protobuf, OData, Dropbox Stone,
  OpenAPI indexes, and human-docs references.
- Deterministic catalog-to-UWS source type mapping for first-class executable
  API source protocols.
- Tracked endpoint advisory overlays and builders for services without official
  OpenAPI documents when official docs or official non-OpenAPI machine specs
  expose enough endpoint instructions to build a useful reviewed subset.
- Bounded SQLite cache metadata with deterministic LRU retention, record byte
  budgets, mandatory digest/byte evidence, and confined relative paths for
  saved catalog artifacts that already exist on disk. The cache package, not
  the CLI, owns refresh-result registration and cache-relative path mapping.
- Read-only provider advisory reports that can join existing local artifact
  registrations without creating caches or fetching remote documents.
- Offline provider resolution, materialization, and workflow export helpers for
  existing catalog artifact registrations and separate security-overlay JSON
  metadata, with provenance manifests.
- Offline refresh review reports that validate saved artifacts and report
  missing registrations or files without creating caches or downloading specs.
- Provider-specific refresh corrections for selected official artifacts where
  source review establishes a small deterministic normalization needed for
  validation, with raw and corrected validation kept separate and correction
  notes visible in refresh outputs. Saved bytes and digests always describe the
  raw download.
- Google Discovery parsing for official Discovery documents via
  `github.com/OpenUdon/googlediscovery`, including
  Discovery paths, resources, methods, inherited root/resource/method
  parameters, schemas, media upload metadata, and OAuth scope metadata.
- AWS Smithy JSON model parsing for official one-service AWS models through
  `github.com/OpenUdon/awssmithy`, including service metadata, operation
  bindings, protocol static fields, greedy path labels, prefix
  header/query-params map bindings, shape metadata, and AWS service/signing
  metadata.
- OpenRPC JSON document parsing for JSON-RPC method metadata, including
  document info, servers, methods, params, results, errors, components, local
  selector aliases, and operation summaries without JSON-RPC invocation.
- GraphQL source parsing for SDL, introspection JSON, and operation documents,
  including schema root types, type and field metadata, operation kind,
  variables, selection names, stable source operation IDs, source refs, and
  operation summaries without operation execution or live introspection.
- gRPC/protobuf source parsing for `.proto`, JSON descriptor, and binary
  `FileDescriptorSet` artifacts, including packages, services, methods,
  request/response messages, streaming flags, message fields, stable source
  operation IDs, source refs, and operation summaries without RPC execution or
  server reflection.
- OData source parsing for CSDL XML, JSON CSDL, and exported service metadata,
  including entity types, entity sets, singletons, navigation properties,
  actions, functions, parameters, query options, stable source operation IDs,
  source refs, and operation summaries without tenant metadata calls, login,
  credential lookup, `$batch`, request execution, or account/tenant selection.
- `helper/...` packages for runtime-importable payload-shaping helpers and
  metadata descriptors. Helpers must remain deterministic, side-effect free,
  provider-client free, credential free, and token-exchange free.

## Non-Goals

- UWS workflow semantics, package review, approval state, or trusted-runner
  handoff; those belong in `../openudon` and `../uws`.
- Runtime execution, request signing, credential parsing or resolution, token
  acquisition or refresh, provider account selection, or production side
  effects; those belong outside `apitools`. Operator OAuth consent and token
  persistence are owned by Udon's trusted `udon oauth google login` command.
- Terraform/OpenTofu parsing; that belongs in `../tfconfig`.
- Treating catalog data or discovered specs as trusted executable input.
- General web crawling, treating authenticated organization registries such as
  Swagger Catalog or Scalar Registry as public global aggregators, or treating
  LAP-transformed metadata as a substitute for its reported original source.
- Vendoring downloaded provider specs, Discovery documents, or local SQLite
  caches into the public repository.
- Treating Smithy, Dropbox Stone, or advisory overlay metadata as
  OpenAPI-compatible import inputs without an explicit importer or converter.
- Treating Google Discovery as official provider OpenAPI truth, executing
  Google API operations, resolving Google credentials, or selecting Google
  accounts.
- Treating AWS Smithy as official AWS OpenAPI truth, executing AWS API
  operations, resolving credentials, signing requests, fetching tokens, or
  selecting AWS accounts or regions.
- Treating OpenRPC, GraphQL, gRPC/protobuf, or OData as official provider
  OpenAPI truth, executing JSON-RPC methods, GraphQL operations, RPCs, OData
  requests, or `$batch`, contacting RPC/GraphQL/gRPC/OData endpoints,
  resolving credentials, resolving variables, live-introspecting services,
  calling tenant metadata endpoints, using server reflection, negotiating TLS,
  selecting accounts/tenants, or implementing runtime transport behavior.
- Lowering Google Discovery, AWS Smithy, AsyncAPI, OpenRPC, GraphQL, or
  gRPC/protobuf, or OData into OpenAPI-bound operation metadata as a default
  materialization/export behavior.
- Copying third-party integration implementation code, generated node
  metadata, credential behavior, descriptions, icons, or runtime workflow
  behavior into catalog evidence.
- Runtime registration, execution, provider API calls, credential lookup,
  secret storage, or account selection for helper functions.
