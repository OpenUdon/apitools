# Changelog

## Unreleased

- Keep top-level response fields when an OpenAPI or Google Discovery response
  summary is truncated: every top-level field is kept first (the first 60 by path
  when more exist), then nested fields by path, up to 60, sorted by path. Large
  resources such as Google Storage buckets no longer lose `name`, `id` and
  `location`. Google Discovery operations now also report
  `schema.field_summary_truncated`. Request-field summaries and public types are
  unchanged.
- Add `catalog.SearchBuiltInProviders`, `Catalog.SearchProviders`,
  `ProviderMatch` and `ProviderMatchRank`: a ranked free-text search over
  provider metadata. Every query word must match; a provider is ranked by the
  lowest-priority text a word needed (ID, name and aliases, then category and
  workflow relevance, then quirks, spec references and security overlays), an
  exact provider key comes first, and ties break by provider ID. It reads only
  the catalog and returns copies. Existing behavior is unchanged.
- Add `Client.FetchDocument` and `FetchedDocument`: the guarded, bounded
  download used by `Import`, exported for consumers that fetch documents
  themselves. It applies the same host and port checks (dial time and every
  redirect), timeout, and wire and decoded byte limits, and honours
  `AllowUnsafeHosts` and `AllowedPorts`. It does not validate the body or use
  the cache. Existing behavior is unchanged.

## v0.1.0 - 2026-05-21

Initial tagged release of `github.com/OpenUdon/apitools`.

- OpenAPI/Swagger discovery, safe download, validation, import, local scanning,
  operation inventories, prompt-safe summaries, auth summaries, and operation
  ranking.
- Metadata-only provider catalog with protocol classification, advisory
  security overlays, catalog quality gates, provider advisory reports, refresh
  review reports, stats, and security audits.
- Native source-family support through standalone Google Discovery and AWS
  Smithy parser modules, with deprecated compatibility wrappers retained in
  this module.
- Source-first catalog coverage for OpenAPI/Swagger, Google Discovery, AWS
  Smithy, Dropbox Stone, OpenAPI index, tenant-exported metadata, and reviewed
  docs-derived advisory overlays.
- Safety boundary: no API operation execution, credential resolution, token
  fetching, request signing, account/region selection, or provider runtime
  behavior.
