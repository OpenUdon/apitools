# Result V3 - Source Kind And Protocol Split

The catalog now separates source reference kind from API description protocol.

- `SpecKind` remains the curation/source reference kind.
- `SpecProtocol` records the model family: OpenAPI, Swagger, Smithy, Google
  Discovery, Dropbox Stone, OpenAPI index, human docs, or unknown.
- Advisory, specs, inspect, refresh, and refresh-report surfaces carry protocol
  classification and version hints where known.

This expands catalog metadata for official non-OpenAPI machine sources without
changing strict OpenAPI import behavior or crossing into API execution,
credential resolution, signing, account selection, or automatic conversion.
