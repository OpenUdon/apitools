# Prompt V19 - Adaptive Authoring Source Discovery

Adaptive workflow authoring needs to inspect all local source evidence before
asking source and operation questions. Add one bounded, local-only apitools API
that accepts explicit file or directory roots and validates OpenAPI/Swagger,
Google Discovery, AWS Smithy, AsyncAPI, GraphQL, OpenRPC, gRPC/protobuf, and
OData.

Return validated candidates with family, title, operation count, score, path,
SHA-256, and provenance, plus rejected/ambiguous files and visible truncation
diagnostics. Reject symlinks, unsafe or non-regular paths, oversized documents,
and cancellation. Deduplicate identical content and require an explicit kind
for ambiguous JSON or XML. Use conventional directory names only as hints.

Keep the boundary narrow: `apitools` discovers and validates source facts. It
does not select an active workflow or source, perform remote authoring lookup,
materialize project files, resolve credentials, or execute operations.
