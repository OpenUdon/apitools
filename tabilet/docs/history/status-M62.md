# Retired milestone M62 - Portable Fnct Helper Catalog

**Milestone.** M62
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M62.md
**Source specification.** tabilet/memory-bank/milestone.md#m62---portable-fnct-helper-catalog
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M62 - Portable Fnct Helper Catalog

**Goal.** Expose a small public helper contract for Gmail runtime-adjacent
helpers, with pure raw-message rendering as the first fnct case and a narrow
OAuth2 bootstrap shape plus local login CLI for trusted runtime credential
fallback.

**Scope.**

- Add a minimal function descriptor type that names the helper, runtime type,
  invocation mode, input fields, output semantics, purity, and side effects.
- Add `helper/gmailmsg` with `gmail.render_raw`, accepting one request-body
  object and returning a Gmail API `raw` string.
- Add Gmail OAuth2 credential-bootstrap helpers that parse local operator
  credential material, produce a consent URL when authorization material is
  missing, and exchange/refresh access tokens only when a trusted runtime calls
  them.
- Add `apitools oauth google login` for operators to run a local consent flow
  or exchange a manual authorization code, then print environment export and
  marker-based data-file guidance.
- Keep helpers deterministic and side-effect free: no provider API calls,
  Gmail clients, account selection, or runtime registration in `apitools`.
  OAuth2 token acquisition is allowed only in the explicit runtime auth helper,
  not in fnct helpers or catalog workflows. The login CLI may exchange tokens
  only when the operator invokes it directly.
- Cover message encoding, template rendering, missing required fields, template
  missing-key failures, header injection rejection, and exported descriptor
  metadata.

**Acceptance.** Downstream OpenUdon and udon can import the public descriptor
and helper implementation to agree on `gmail.render_raw`, and udon can import
the Gmail OAuth2 bootstrap shape for runtime-owned credential fallback; the CLI
can produce refresh-token env guidance for marker-based `data.hcl`; `go test
./helper/...`, `go test ./...`, `go vet ./...`, dependent OpenUdon/udon
focused tests, and `git diff --check` pass.

M72 supersedes the OAuth portions of this historical milestone: credential
parsing, token refresh, consent, and authorization-code exchange now live only
in Udon's trusted runtime and CLI. `apitools` retains the pure
`gmail.render_raw` helper and descriptor.
````

## Status record

````markdown
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
````
