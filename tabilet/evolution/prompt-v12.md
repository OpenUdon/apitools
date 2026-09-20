# Prompt V12 - n8n Protocol Connector Boundary Review

After adding n8n's public REST API as a provider entry, review the broader n8n
node tree for connector families that should not be promoted as first-class
OpenAPI/Smithy/Discovery sources.

Keep the milestone narrow:

- Document database, broker, mail, file, directory, generic execution, generic
  HTTP, webhook, and GraphQL nodes as priority signals only.
- Do not add catalog providers, candidates, docs-derived OpenAPI overlays, or
  native parser work for generic protocol connectors.
- Preserve a future path for native protocol-family support only when backed by
  stable metadata artifacts and metadata-only parsing.
- Do not execute provider or protocol operations, introspect live services,
  resolve credentials, choose hosts, queues, topics, databases, accounts, or
  workflow runtime state.
