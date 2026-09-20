# Result V8 - Native AWS Smithy Model Parser

`apitools/awssmithy` now exposes native AWS Smithy parsing alongside deprecated
OpenAPI-shaped compatibility conversion.

- `Parse` and `ParseMap` return protocol-aware Smithy metadata with service,
  operation, binding, static request field, shape, endpoint-prefix, and
  signing-name information.
- Runtime-relevant bindings such as greedy path labels, blob payload members,
  prefix headers, and query-param maps are preserved for downstream
  Smithy-aware execution planning.
- Native parser coverage includes `restJson1`, `restXml`, `awsQuery`,
  `ec2Query`, `awsJson1_0`, and `awsJson1_1`.
- `Convert` and `ConvertMap` remain available only as deprecated derived
  metadata compatibility wrappers.
- `../udon` no longer uses converted OpenAPI-shaped Smithy metadata for
  Smithy runtime-plan generation.
