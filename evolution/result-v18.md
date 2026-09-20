# Result V18 - UWS 1.4 Source Tooling Roadmap

The apitools roadmap now opens three separate source-family milestones after
M64:

- M65 GraphQL Source Support for local SDL, introspection JSON, and operation
  document metadata.
- M66 gRPC/Protobuf Source Support for local `.proto` files and descriptor
  metadata.
- M67 OData Source Support for local CSDL/service metadata.

M65 is the active milestone; M66 and M67 are planned. OpenRPC remains covered by
M63. The roadmap preserves the metadata-only boundary: no source operation
execution, live service probing, credential resolution, token fetching, or
account/project/tenant/workspace selection belongs in `apitools`.
