# Retired milestone M50 - n8n Protocol Connector Boundary Review

**Milestone.** M50
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M50.md
**Source specification.** tabilet/memory-bank/milestone.md#m50---n8n-protocol-connector-boundary-review
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M50 - n8n Protocol Connector Boundary Review

**Goal.** Record the post-M49 decision for n8n protocol connector nodes so
future provider expansion does not promote generic runtime connectors as
OpenAPI/Smithy/Discovery catalog sources.

**Scope.**

- Review n8n node roots for database, broker, mail, file, directory, generic
  execution, generic HTTP, webhook, and GraphQL connector signals.
- Update public documentation to classify these nodes as priority signals only,
  not provider rows, docs-derived OpenAPI overlays, native source parsers, or
  n8n runtime compatibility evidence.
- Define criteria for any future native protocol-family milestone, including
  stable schema artifacts, metadata-only parsing, and downstream-owned runtime
  binding.
- Preserve provider-specific catalog paths for true service APIs that happen to
  share protocol names, such as hosted control-plane APIs, when a future
  source-backed milestone explicitly scopes them.
- Do not add catalog providers, candidates, overlays, parser packages,
  database introspection, network protocol clients, credential resolution, host
  selection, queue/topic/database selection, or workflow execution.

**Acceptance.** Maintainers have a documented, reviewable boundary for n8n
protocol connector nodes; public docs point to the boundary; no provider or
parser behavior changes are introduced; required repository and private harness
diff checks pass.
````

## Status record

````markdown
# Status M50 - n8n Protocol Connector Boundary Review

| Item | State | Notes |
|---|---|---|
| Scope freeze | `[+]` | M50 is a boundary/documentation milestone only. It does not add catalog providers, candidates, overlays, parsers, or source-family runtime behavior. |
| n8n connector review | `[+]` | Reviewed local n8n node roots for database/key-value, broker, mail, file, directory, generic execution, generic HTTP, webhook, and GraphQL connector signals. |
| Public boundary docs | `[+]` | Updated `docs/non-openapi-protocols.md` with protocol connector families, exact n8n node-root examples, no-promotion examples, and future native protocol-family criteria. |
| Public discoverability | `[+]` | Linked the non-OpenAPI/protocol boundary notes from the README catalog curation section. |
| Boundary checks | `[+]` | Confirmed M50 introduces no provider rows, candidate rows, overlays, parser APIs, live introspection, credential handling, host selection, or workflow execution behavior. |
| Review and verification | `[+]` | Deep review found no blocking fixes. Passed `go test ./...`, `go vet ./...`, `go run ./cmd/apitools catalog check`, `go run ./cmd/apitools search --help`, `go run ./cmd/apitools import --help`, `git diff --check`, and `git -C ../tofu diff --check -- apitools`. |

## Boundary Notes

- n8n node roots remain priority signals only, not runtime compatibility
  evidence.
- Database and broker protocol nodes must not be mapped to similarly named
  service control-plane APIs unless a future source-backed milestone explicitly
  scopes those provider APIs.
- Future native protocol-family support needs stable metadata artifacts and
  must remain local/downloaded metadata parsing only.
- `apitools` must not open sockets to live databases, brokers, mail servers,
  LDAP directories, FTP servers, webhooks, or generic HTTP endpoints for
  protocol discovery.
````
