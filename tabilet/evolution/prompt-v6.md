# Prompt V6 - Google Discovery Converter Ownership

Move the reusable Google Discovery to OpenAPI-shaped metadata converter into
`apitools` as a public library package, after catalog coverage milestones made
standalone converter ownership appropriate.

Preserve the metadata-only boundary: converted Discovery output is derived
OpenAPI-shaped metadata for indexing and downstream generation, not official
provider OpenAPI truth. Do not add API execution, credential resolution, token
fetching, request signing, or account selection.
