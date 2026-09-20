# Status M64 - Postman/RAML/API Blueprint Import Evaluation

State: Complete

## Scope

Evaluate Postman Collection, RAML, and API Blueprint as source/advisory
artifact families for local detection, classification, and possible
source-aware conversion. This milestone is documentation and planning only; it
does not add parser code, remote workspace calls, UWS source types, credential
resolution, or runtime behavior.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| Local artifact detection posture | `[+]` | Recorded local-only detection signals for Postman Collection, RAML, and API Blueprint without accepting remote workspace URLs or authenticated discovery as default behavior. |
| Classification model | `[+]` | Defined documentation/planning labels that preserve provenance and distinguish provider-owned API descriptions from client/workspace artifacts. |
| Conversion decision notes | `[+]` | Recorded when source-aware OpenAPI conversion may be acceptable, when stronger native sources are preferred, and when semantic loss or weak provenance keeps artifacts advisory only. |
| Boundary checks | `[+]` | Confirmed no execution, credential lookup, account/workspace selection, remote workspace calls, parser implementation, or new UWS `sourceDescriptions[].type` values are introduced. |
| Verification | `[+]` | Ran docs/repo diff checks for `../apitools` and `../tofu`; Go tests were not needed because the milestone changed only public docs and memory-bank planning files. |
| Final review | `[+]` | Reconciled public protocol notes, memory-bank milestone state, and evolution no-bump decision before closing the milestone. |

## Completed Decision Posture

- Postman collections remain client, test, or workspace artifacts by default.
  Detect only local JSON collection files with Postman schema evidence,
  `*.postman_collection.json`, or local v3 collection source-tree signals such
  as `definition.yaml` and `*.request.yaml`. Classify them as advisory
  client/workspace evidence unless provenance proves provider ownership,
  intended API-truth status, and useful completeness.
- RAML may be provider-owned API description evidence. Detect `.raml`, `.yaml`,
  or `.yml` files whose first meaningful line is `#%RAML 1.0` or `#%RAML 0.8`.
  Classify matching local artifacts as source-aware conversion candidates, not
  first-class UWS source types.
- API Blueprint may be documentation-oriented API description evidence. Detect
  `.apib` or `.md` files with leading metadata such as `FORMAT: 1A`. Classify
  matching local artifacts as source-aware conversion candidates when review can
  preserve methods, resources, parameters, bodies, schemas/examples, auth hints,
  and provenance.
- Prefer an official OpenAPI/Swagger, Google Discovery, AWS Smithy, AsyncAPI,
  OpenRPC, or other stronger native source when one is available for the same
  provider and surface.
- Allow conversion to OpenAPI only as reviewed, source-aware output. The review
  must confirm that methods/resources, parameters, request and response bodies,
  schemas or examples, auth hints, and provenance survive without misleading
  downstream authoring tools.
- Keep artifacts advisory and do not convert when variables, environments,
  scripts/tests, includes, overlays, annotations, Markdown ambiguity, or missing
  provenance would make converted OpenAPI output misleading.
- Do not add remote Postman workspace calls, API execution, credential
  resolution, account/workspace selection, parser packages, materialization
  directories, CLI commands, catalog cache behavior, or UWS
  `sourceDescriptions[].type` values for `postman-collection`, `raml`, or
  `api-blueprint` in M64.
- The classification names `postman-collection`, `raml`, and `api-blueprint`
  are documentation/planning labels only in this milestone.

## External Source Checks

- Postman documentation reviewed on 2026-05-27 records schema 3.0.0 as a
  source-tree format with `definition.yaml`, `*.request.yaml`, and request
  resource directories, while exported 2.1 collections remain relevant for
  Newman.
- The RAML specification repository and examples use the leading `#%RAML 1.0`
  marker for RAML documents; M64 also preserves the legacy `#%RAML 0.8`
  detection marker.
- API Blueprint specification documentation uses leading metadata such as
  `FORMAT: 1A`.

## Evolution Decision

No evolution bump. M64 is an evaluation milestone that preserves the existing
metadata-only boundary and explicitly declines runtime behavior, parser
implementation, and new UWS source-type promotion.

## Verification

- Passed: `git diff --check` in `../apitools`
- Passed: `git diff --check` in `../tofu`
- Not run: `go test ./...` in `../apitools`, because the milestone changed only
  public docs and private memory-bank planning files.
