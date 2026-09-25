# Retired milestone M43 - AWS Smithy Service Expansion

**Milestone.** M43
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M43.md
**Source specification.** tabilet/memory-bank/milestone.md#m43---aws-smithy-service-expansion
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M43 - AWS Smithy Service Expansion

**Goal.** Expand first-class AWS Smithy source coverage for high-priority AWS
services represented in `../n8n`, using official AWS Smithy JSON AST models as
cataloged native API metadata.

**Scope.**

- Freeze the service batch: SQS, DynamoDB, SES or SESv2, IAM, Cognito,
  Textract, Rekognition, Comprehend, Transcribe, ACM, ELB, and ELBv2.
- Add or update durable provider entries, aliases, categories, official Smithy
  spec references, source notes, license notes, verified dates, quirks, and
  user OpenAPI need classifications for each service.
- Save official review artifacts from `aws/api-models-aws` under ignored
  catalog cache paths and register saved artifact paths when artifacts are
  checked in or intentionally tracked for review.
- Exercise `github.com/OpenUdon/awssmithy` native parsing against selected
  service fixtures that cover AWS protocol families and Smithy trait shapes
  used by the batch.
- Add or update SigV4/security metadata overlays or classifications as
  metadata only; do not resolve credentials, sign requests, choose accounts, or
  choose regions.
- Keep Smithy artifacts classified as `smithy`; do not treat them as official
  OpenAPI and do not reintroduce Smithy-to-OpenAPI conversion in `apitools`.

**Acceptance.** The AWS batch is discoverable through catalog/provider/spec
inspection with official Smithy provenance, native parser coverage protects
the batch's protocol shapes, catalog quality passes, and no API execution,
credential resolution, signing, account selection, or region selection is
introduced.
````

## Status record

````markdown
# Status M43 - AWS Smithy Service Expansion

| Item | State | Notes |
|---|---|---|
| Batch source review | `[+]` | Confirmed official `aws/api-models-aws` Smithy JSON model paths, service dates, and GitHub blob revisions for SQS, DynamoDB, SES classic, IAM, Cognito Identity Provider, Textract, Rekognition, Comprehend, Transcribe, ACM, ELB, and ELBv2. |
| Provider/catalog rows | `[+]` | Added durable provider entries, aliases, categories, official Smithy spec references, source notes, license notes, verified dates, quirks, and user OpenAPI need classifications for the full batch. |
| Review artifacts | `[+]` | Saved selected official Smithy JSON review artifacts under ignored `catalog-openapi-cache/openapi/aws-*-smithy-model.json` paths and registered artifact paths in the artifact registry workflow. |
| Native parser coverage | `[+]` | Extended `awssmithy` catalog artifact coverage so the downloaded batch artifacts remain Smithy-shaped and parse through the native wrapper when present. |
| Security metadata | `[+]` | Added source-backed SigV4 overlays for the batch as metadata only, without credential resolution, request signing, account selection, or region selection. |
| Catalog/reporting checks | `[+]` | Verified provider lookup, `catalog specs`, `catalog advisory`, `catalog inspect`, `catalog stats`, `catalog security-report`, `catalog refresh-report`, and catalog quality output for the batch. |
| Documentation and memory bank | `[+]` | Product, architecture, and tech-stack boundaries already cover native AWS Smithy parsing; this status file records the M43 service expansion without changing the broader contract. |
| Verification | `[+]` | `go test ./catalog`, `go test ./awssmithy`, `go test ./...`, `go vet ./...`, `go run ./cmd/apitools catalog check`, catalog CLI smoke checks, and diff checks pass. |

## Boundary Notes

- AWS Smithy models remain native `smithy` source metadata, not official
  OpenAPI and not an OpenAPI import contract.
- `apitools` must not execute AWS operations, resolve credentials, fetch
  tokens, sign requests, or choose AWS accounts or regions.
- `../n8n` is only a priority signal for service selection.
````
