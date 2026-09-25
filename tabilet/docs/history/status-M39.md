# Retired milestone M39 - AWS Smithy Converter Ownership

**Milestone.** M39
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M39.md
**Source specification.** tabilet/memory-bank/milestone.md#m39---aws-smithy-converter-ownership
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M39 - AWS Smithy Converter Ownership

**Goal.** Add reusable AWS Smithy JSON to OpenAPI-shaped metadata conversion to
`apitools` while preserving the metadata-only boundary and enabling direct
downstream Smithy entrypoints.

**Scope.**

- Add `github.com/OpenUdon/apitools/awssmithy` with byte and decoded-map
  conversion APIs for one-service Smithy JSON 2.0 documents.
- Support official AWS service models using `restJson1`, `restXml`, and
  `awsQuery` protocols.
- Convert service metadata, Smithy HTTP bindings, AWS Query form metadata,
  schemas, recursive refs, and AWS service/signing traits into bounded
  OpenAPI 3.0.1-shaped metadata.
- Keep catalog AWS Smithy artifacts classified as `smithy` structured review
  artifacts unless an explicit caller invokes the converter.
- Do not add CLI commands, execute AWS API operations, resolve credentials,
  fetch tokens, sign requests, or choose AWS accounts or regions.

**Acceptance.** `awssmithy` is covered by focused restJson1/restXml/awsQuery
tests and malformed-shape checks, full `apitools` test/vet/diff checks pass,
`../udon` consumes the converter through explicit Smithy entrypoints, and
SigV4 signing is documented as later runtime work.
````

## Status record

````markdown
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
````
