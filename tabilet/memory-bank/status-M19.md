# Status M19 - Spec Protocol Classification

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Represent official server-side API description sources beyond OpenAPI as
catalog metadata that downstream tools can understand without broadening strict
OpenAPI import behavior.

## Direction

Catalog source kind identifies where a reference sits in the curation model;
spec protocol identifies the API description model family. OpenAPI/Swagger
remain the only strict import targets today, while Smithy, Google Discovery,
Dropbox Stone, OpenAPI indexes, and human docs remain review or advisory
metadata until an explicit importer/converter exists.

## Tasks

| Item | State | Notes |
|---|---|---|
| Protocol model added | `[+]` | Added `SpecProtocol` and protocol classification helpers for OpenAPI, Swagger, Smithy, Google Discovery, Dropbox Stone, OpenAPI indexes, human docs, and unknown references. |
| Reports expose protocol | `[+]` | Advisory, specs, inspect, refresh, and refresh-report result structs include protocol and version fields where available. |
| Version hints preserved | `[+]` | Static catalog rows infer conservative OpenAPI/Swagger version hints; refresh and refresh-report use parsed OpenAPI/Swagger metadata when available. |
| Tests updated | `[+]` | Added protocol classification coverage for catalog rows, advisory output, refresh results, refresh-review results, and CLI specs/inspect output. |
| Verification completed | `[+]` | Ran Go tests/vet, diff check, catalog quality, specs/advisory/inspect/refresh-report CLI smoke checks, and sibling `openudon`/`udon` tests. |

## Boundary Checks

- Do not execute provider API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not treat Smithy, Discovery, Stone, OpenAPI indexes, or human docs as
  import-ready OpenAPI documents without an explicit importer/converter.
- Do not treat advisory overlays as provider truth.
