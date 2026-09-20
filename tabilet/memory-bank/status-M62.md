# Status M62 - Portable Fnct Helper Catalog

> Historical note: M72 later moved every OAuth credential and token-exchange
> behavior to Udon. The OAuth rows below record what M62 originally delivered;
> they are no longer current apitools interfaces.

| Item | State | Notes |
|---|---|---|
| Helper descriptor contract | `[+]` | Added `helper/fnctspec.FunctionSpec` and request-body-object invocation metadata for pure runtime-importable helpers. |
| Gmail raw-message helper | `[+]` | Added `helper/gmailmsg` with `gmail.render_raw`, `RenderRaw`, `RenderRawAny`, and descriptor export. |
| Gmail OAuth2 bootstrap helper | `[+]` | Added `helper/gmailmsg.GoogleOAuth2Credential` helpers for trusted runtimes to parse local OAuth material, build consent URLs, and exchange/refresh access tokens without making OAuth a fnct helper. |
| Google OAuth login CLI | `[+]` | Added `apitools oauth google login` with localhost callback and manual `--code` exchange modes, env-first client secret handling, refresh-token env export output, and marker-based HCL snippets. |
| Shared helper catalog | `[+]` | Added a parent `helper` catalog for descriptor lookup and trusted runtime registration hooks so authoring and runtime imports do not drift. |
| AsyncAPI protocol classification side task | `[+]` | Added catalog protocol/UWS source type classification for AsyncAPI so OpenUdon can materialize source-aligned `asyncapi/` artifacts without owning protocol execution. Follow-up tracked representative AsyncAPI artifacts for Gemini WebSocket, Slack RTM, and Streetlights MQTT under `catalog-openapi-cache/asyncapi/`. |
| Helper regression tests | `[+]` | Covered base64url Gmail raw encoding, body-template rendering, missing fields, mutually exclusive body inputs, missing template keys, header injection rejection, non-ASCII subject encoding, OAuth token refresh and CLI code exchange through httptest, authorization-required errors, shared catalog lookup, registration, and descriptor metadata. |
| Memory-bank boundary update | `[+]` | Product, architecture, tech-stack, milestone, status, and evolution notes record the helper boundary. |

## Boundary Checks

- Fnct helpers are pure payload shaping utilities.
- Gmail OAuth2 bootstrap is not a fnct helper and runs only when a trusted
  runtime supplies local operator credential material.
- The OAuth login CLI is an explicit operator utility; it may exchange OAuth
  codes but must not call Gmail resource APIs or store secrets.
- No provider API calls, account selection, runtime registration, secret
  storage, or workflow execution belongs in `apitools`.
- Udon owns trusted runtime registration and execution.
- OpenUdon owns authoring, review, package evidence, and trusted handoff.
