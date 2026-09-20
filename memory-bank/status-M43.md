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
