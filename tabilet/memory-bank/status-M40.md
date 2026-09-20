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
