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
