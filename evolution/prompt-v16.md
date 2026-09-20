# Prompt V16 - Provider Artifact Materialization Boundary

Move static provider catalog materialization out of downstream OpenUdon and
into `apitools`.

`apitools` should resolve provider names to catalog metadata, report whether a
provider has executable OpenAPI, Discovery-only, Smithy-only, docs-only, or
advisory-overlay-only local metadata, copy existing cache-registered artifacts
into provider- or workflow-scoped directories, emit catalog security overlays
as separate JSON metadata, and write provenance manifests.

Keep this surface offline and copy-only. Do not download missing artifacts,
execute API operations, resolve credentials, apply overlays into provider
OpenAPI documents, or lower Discovery/Smithy into OpenAPI-bound operation
metadata by default.
