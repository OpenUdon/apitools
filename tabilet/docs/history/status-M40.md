# Retired milestone M40 - AWS Smithy Native Model Ownership

**Milestone.** M40
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M40.md
**Source specification.** tabilet/memory-bank/milestone.md#m40---aws-smithy-native-model-ownership
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M40 - AWS Smithy Native Model Ownership

**Goal.** Make `apitools/awssmithy` own native AWS Smithy model parsing and
remove OpenAPI-shaped Smithy conversion as a runtime contract source.

**Scope.**

- Add `Parse` and `ParseMap` APIs for one-service Smithy JSON 2.0 documents.
- Preserve service, operation, shape, binding, protocol, static
  query/header/payload fields, endpoint-prefix, and signing-name metadata
  without producing OpenAPI.
- Cover `restJson1`, `restXml`, `awsQuery`, `ec2Query`, `awsJson1_0`, and
  `awsJson1_1`.
- Remove `Convert` and `ConvertMap`; downstream callers must use native Smithy
  parsing or their own explicit compatibility adapter.
- Do not execute AWS operations, resolve credentials, sign requests, fetch
  tokens, or select AWS accounts or regions.

**Acceptance.** Native parsing covers the six AWS protocol families with
focused tests, `../udon` no longer uses converted OpenAPI-shaped Smithy
metadata for runtime-plan generation, no public Smithy-to-OpenAPI conversion
API remains in `apitools`, and required tests/vet/diff checks pass.
````

## Status record

````markdown
# Status M40 - AWS Smithy Native Model Ownership

## Goal

Make `apitools/awssmithy` the native AWS Smithy model parser while removing
OpenAPI-shaped conversion from the public Smithy API.

## Status

| Item | State | Notes |
|---|---|---|
| Native parser added | `[+]` | `awssmithy.Parse` and `ParseMap` expose service, operation, protocol, binding, static field, shape, endpoint-prefix, and signing-name metadata without producing OpenAPI. |
| Protocol coverage expanded | `[+]` | Parser coverage now includes `restJson1`, `restXml`, `awsQuery`, `ec2Query`, `awsJson1_0`, and `awsJson1_1`. |
| Runtime protocol bindings preserved | `[+]` | Parser coverage now explicitly preserves execution-relevant metadata such as greedy path labels, blob payload members, response status-code members, response headers, `httpPrefixHeaders`, and `httpQueryParams` for downstream Smithy-aware runtimes. |
| Conversion removed | `[+]` | `Convert` and `ConvertMap` were removed; `awssmithy` now exposes native parsing only. |
| Downstream native adoption verified | `[+]` | `../udon` Smithy runtime-plan generation uses native Smithy metadata instead of converted OpenAPI-shaped bytes. |
| Verification completed | `[+]` | Focused/full tests passed in `apitools` and `udon`; vet and diff checks are part of closeout. |

## Boundary Checks

- `apitools` parses and classifies untrusted Smithy metadata only.
- OpenAPI-shaped Smithy conversion is not exposed by `apitools`; Smithy
  metadata is parsed natively and is not AWS provider OpenAPI truth.
- No AWS operations are executed, no credentials are resolved, no requests are
  signed, and no AWS accounts or regions are selected in `apitools`.
````
