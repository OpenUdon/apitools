# Result V7 - AWS Smithy Converter Package

`apitools` now owns reusable AWS Smithy JSON conversion through
`github.com/OpenUdon/apitools/awssmithy`.

- `Convert` and `ConvertMap` produce bounded OpenAPI 3.0.1-shaped metadata from
  one-service Smithy JSON 2.0 documents.
- The converter covers AWS `restJson1`, `restXml`, and `awsQuery` protocols,
  including Smithy HTTP bindings, AWS Query form defaults, schemas, recursive
  refs, and AWS service/signing extensions.
- Follow-up review hardening preserves query-literal and AWS Query operation
  uniqueness, models Smithy prefix/query parameter bindings, and maps response
  payload/header/status-code bindings into derived OpenAPI-shaped responses.
- Catalog AWS Smithy artifacts remain classified as `smithy` structured review
  evidence unless an explicit caller converts them.
- `../udon` uses explicit Smithy parsing and generator/runtime-plan entrypoints
  while leaving SigV4 signing and credential binding to later runtime work.
