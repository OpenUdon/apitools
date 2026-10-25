# Operation mapping v1 conformance

`cases.json` is a reusable metadata-only corpus with version
`apitools.operation-mapping/v1`. It contains 19 reviewed positive/refusal
cases and a raw SHA256 manifest for four small local OpenAPI source files.
No source is fetched and no described provider operation runs. All names and
data are synthetic; the Cloudflare resource shape is a fake test use case,
not an assertion about a live Cloudflare API or account.

The business-message and resource-role cases cover exact native selectors,
path/query/header/body inputs and whole response schemas. Refusals cover
placeholders, optional/unknown presence, missing fields, mismatched primitive
and nested types, missing required body properties, duplicate mappings,
credential-shaped schema fields, declared security query parameters,
truncation, stale source digests, candidate selector drift, ambiguous native
IDs and an unknown external schema. The external schema URL is a no-fetch
fixture. Product source admission and authority remain outside the corpus.

Each case names a source, native selector kind/value, explicit field mappings,
expected overall outcome/resolution status, field codes and input coverage.
`want_native_selector` records the reviewed canonical native reference even
when selection uses an ID alias. Source identity uses lowercase, unprefixed
raw SHA256. Consumer fixtures should pin the producer revision and these bytes,
reproduce the table from source bytes, retain its complete schema/security
evidence and compare the expected outcomes. A copied fixture is conformance
evidence; it does not prove ordinary published module selection or compilation.

`TestOperationMappingV1Conformance` uses public APIs to build and independently
reproduce shapes, generate actual candidates, and exercise the mappings.
Additional tests cover cancellation, bounded/malformed metadata, value-free
errors, reserved and credential headers, candidate source drift and every
existing native source family. The older eight-family fixture remains frozen;
its incomplete operations remain unresolved.

The corpus is deliberately maintained data, not a generated snapshot of helper
output. Review behavior and expected results before changing it. Tests never
rewrite fixtures.

`cases.json`: 44,469 bytes; SHA256
`1f456c87fb0d4a594a2b26fdff5328f71a445157cf0feba990e17929e07fd847`.

| Source | Raw SHA256 |
|---|---|
| business.json | `17f3693b85299444e0d8d5f648bca3b5f28348baa3afb9e3bba16f8cd016b23c` |
| resource.json | `462b9206828557af364ffaf993cf1a139b8008c2cf3a79552ba3856ac88b624f` |
| ambiguous.json | `17cfe84a3ab3b0dc651d925baa49d65b8478f30602ea391da4e1619a2d18b59a` |
| unknown.json | `d9081abd88644e6142e67920d804bb0b3d3c21b297c89d3ddc7b9026d0dd8634` |
