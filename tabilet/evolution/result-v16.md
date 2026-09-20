# Result V16 - Static Catalog Materialization

`apitools` now owns a first static provider artifact handoff surface.

The catalog package exposes provider resolution, single-provider
materialization, and workflow artifact export reports. The CLI adds
`catalog resolve`, `catalog materialize`, and `catalog export`, all backed by
existing cache artifact registrations and provider catalog metadata.

Materialization remains metadata-safe: it copies existing local artifacts,
emits catalog security overlays as separate JSON files, and writes provenance
manifests. It does not fetch missing documents, apply overlays into OpenAPI
files, execute provider operations, resolve credentials, or lower native
Discovery/Smithy metadata into OpenAPI.
