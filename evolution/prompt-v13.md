# Prompt V13 - n8n Gap Report And Source Batch Freeze

After documenting the M50 protocol connector boundary, add a repeatable
offline report that compares the local n8n node tree with built-in catalog
coverage and freezes the next source-review batch.

Keep the work bounded:

- Compare n8n node-root directory names with provider IDs, display names,
  aliases, and curated category coverage.
- Classify non-covered roots as provider API candidates, M50 protocol
  exclusions, local workflow utilities, GraphQL/protocol-family candidates, or
  internal directory entries.
- Freeze an 8-12 service source-review batch only where provider-owned OpenAPI
  evidence or strong official docs-derived overlay evidence exists.
- Do not parse n8n node implementation files, import workflows, add provider
  rows, create overlays, probe provider APIs, resolve credentials, choose
  hosts/accounts, or execute operations.
