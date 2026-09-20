# Result V5 - Dropbox Stone Advisory Subset

Advisory endpoint overlays may now be provenance-labeled as derived from an
official non-OpenAPI machine spec when no official OpenAPI document exists.

- Dropbox remains cataloged with official protocol `dropbox-stone` and
  official OpenAPI availability unavailable.
- `dropbox-core-api-overlay` records `official_openapi: false`,
  `derived_from_official_machine_spec: true`, and
  `source_protocol: dropbox-stone`.
- The overlay remains a reviewed subset covering account info, folder listing,
  metadata, upload, download, and shared-link operations.
- Advisory output can carry registered overlay metadata while preserving the
  no-execution, no-credential, no-runtime-protocol boundary.
