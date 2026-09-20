# Prompt V8 - Native AWS Smithy Model Ownership

Reframe AWS Smithy support away from OpenAPI-shaped conversion. `apitools`
should expose reusable native Smithy metadata for official AWS Smithy JSON
models, while keeping conversion only as deprecated compatibility metadata.

The downstream goal is protocol-aware UWS execution in `../udon`, so Smithy
metadata must preserve AWS protocol semantics rather than treating OpenAPI as
the execution contract.
