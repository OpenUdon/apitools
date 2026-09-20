# Status M49 - n8n Public API OpenAPI Entry

| Item | State | Notes |
|---|---|---|
| Scope freeze | `[+]` | M49 is limited to n8n's own public REST API. It does not promote database, queue, mail, SQL, LDAP, FTP, or other protocol connector nodes from the n8n tree. |
| Source review | `[+]` | Reviewed official n8n public REST API docs and the official n8n repository OpenAPI root at `packages/cli/src/public-api/v1/openapi.yml`. |
| Candidate/provider rows | `[+]` | Added `n8n` candidate and provider catalog rows with official OpenAPI, docs references, aliases, source notes, and deployment-specific-host quirks. |
| Artifact registration | `[+]` | Saved ignored `catalog-openapi-cache/openapi/n8n-public-api-openapi-v1.yml` and registered it in `catalog-openapi-cache/cache.sqlite` through the artifact registry. The artifact is marked `invalid` because the repository root is unbundled and contains relative `$ref` entries. |
| Auth/security metadata | `[+]` | Classified n8n auth as complete from the root OpenAPI `ApiKeyAuth` and `BearerAuth` alternatives. No API keys, bearer tokens, users, projects, credentials, workflows, executions, or data tables are selected or resolved. |
| Catalog snapshots | `[+]` | Updated deterministic candidate, provider, security, protocol-stat, alias, and provider classification snapshots for the new n8n row. |
| Review and verification | `[+]` | Deep review found no blocking fixes. Passed `go test ./...`, `go vet ./...`, `go run ./cmd/apitools catalog check`, focused n8n list/spec/advisory/security/stats/refresh-report smoke checks, `go run ./cmd/apitools search --help`, `go run ./cmd/apitools import --help`, `git diff --check`, and `git -C ../tofu diff --check -- apitools`. |

## Boundary Notes

- The n8n node directory remains a priority signal only, not runtime
  compatibility evidence.
- The n8n OpenAPI root is official source metadata, but strict OpenAPI import
  may require bundling the relative references first.
- M49 does not add native database, message-broker, mail, SQL, LDAP, FTP, or
  generic protocol source families.
- `apitools` must not execute n8n API operations, resolve credentials, fetch
  tokens, or choose n8n instances, users, projects, credentials, workflows,
  executions, or data tables.
