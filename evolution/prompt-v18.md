# Prompt V18 - UWS 1.4 Source Tooling Roadmap

After UWS 1.4 adds `graphql`, `openrpc`, `grpc-protobuf`, and `odata` source
description types, align `apitools` source metadata work with that contract.

OpenRPC metadata support is already complete. Open separate metadata-only
milestones for GraphQL, gRPC/protobuf, and OData so each source family can
preserve native selector and schema semantics without being flattened into
OpenAPI.

Keep the boundary strict: parse or inspect local/provider-owned source
artifacts, expose prompt-safe summaries and selector metadata, and classify
catalog/materialization output. Do not execute operations, probe live services,
perform live introspection or server reflection, resolve credentials, fetch
tokens, or select accounts, projects, tenants, endpoints, workspaces, or
regions.
