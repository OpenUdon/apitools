# Prompt V1 - Private Harness And Catalog Direction

Bootstrap private project memory for `apitools` using the same harness shape as
the official skills package and existing OpenUdon/tofu harnesses.

The harness consists of:

- `AGENTS.md`
- `memory-bank/`
- `evolution/`

Canonical files should live under `../tofu/apitools`, while the public
`../apitools` checkout should expose those files through local symlinks and
ignore them so the harness does not become public repository content.

The next product direction is to add an `apitools` provider catalog for popular
workflow services, official OpenAPI or machine-readable spec references, and
security overlays that supplement incomplete provider specs without executing
API operations or resolving credentials.
