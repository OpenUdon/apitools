# Status M39 - AWS Smithy Converter Ownership

## Goal

Add reusable AWS Smithy JSON to OpenAPI-shaped metadata conversion to
`apitools` while preserving the metadata-only boundary and downstream ownership
split.

## Status

| Item | State | Notes |
|---|---|---|
| Converter package added | `[+]` | `awssmithy` exposes `Convert` and `ConvertMap` for one-service Smithy JSON 2.0 documents with restJson1, restXml, awsQuery, schema, recursion, and malformed-input coverage. |
| Downstream ownership split updated | `[+]` | `../udon` adds explicit Smithy parse/generator/runtime-plan entrypoints that delegate reusable conversion to `apitools/awssmithy`. |
| Verification completed | `[+]` | Ran focused/full `apitools` and focused `udon` Smithy tests during implementation; full vet/diff checks are part of milestone closeout. |
| Review hardening completed | `[+]` | Follow-up review fixes preserve Smithy query-literal and AWS Query operation uniqueness, model `httpQueryParams` / `httpPrefixHeaders`, and preserve response HTTP bindings for payloads, headers, and status-code members. |

## Boundary Checks

- AWS Smithy conversion remains derived metadata, not official AWS OpenAPI
  truth.
- Do not execute AWS API operations, resolve credentials, fetch tokens, sign
  requests, or choose AWS accounts or regions.
- Preserve SigV4 service/signing metadata as extensions only; runtime signing
  belongs to a later downstream `udon` / `grand` / `soliton` / `httpclient`
  feature.
- Keep catalog AWS Smithy artifacts classified as `smithy` structured review
  evidence unless an explicit caller converts them.
- Direct `../udon` adoption resolves the preserved operation IDs for S3-style
  query-literal operations and SNS-style AWS Query actions through explicit
  Smithy entrypoints.
