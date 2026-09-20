# Status M41 - Google Discovery Native Model Ownership

## Goal

Make `apitools/googlediscovery` the native Google Discovery parser while
removing OpenAPI-shaped conversion from the public Discovery API.

## Status

| Item | State | Notes |
|---|---|---|
| Native parser added | `[+]` | `googlediscovery.Parse` and `ParseMap` expose service, schema, inherited parameter, method, OAuth scope, and media upload metadata without producing OpenAPI. The implementation now lives in standalone `github.com/OpenUdon/googlediscovery`; `apitools/googlediscovery` is a deprecated wrapper. |
| Runtime-relevant metadata preserved | `[+]` | Parser coverage preserves flattened operation IDs, service URL/path-prefix behavior, request/response refs, all declared media upload protocols, selected simple/multipart execution hints, and per-method scopes for downstream runtimes. |
| Conversion removed | `[+]` | `Convert` and `ConvertMap` were removed; the standalone `googlediscovery` module exposes native parsing only. |
| Downstream native adoption verified | `[+]` | `../udon` Discovery runtime-plan generation uses native Discovery metadata instead of converted OpenAPI-shaped bytes. |
| Verification completed | `[+]` | Focused tests passed in `apitools`, `udon`, and `grand`; full test/vet/diff checks are part of closeout. |

## Boundary Checks

- `apitools` parses and classifies untrusted Google Discovery metadata only by
  delegating to `github.com/OpenUdon/googlediscovery`.
- OpenAPI-shaped Discovery conversion is not exposed by `apitools`; Discovery
  metadata is parsed natively and is not Google provider OpenAPI truth.
- No Google operations are executed, no credentials are resolved, no tokens are
  fetched, and no Google accounts or projects are selected in `apitools`.
