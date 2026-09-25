# Retired milestone M49 - n8n Public API OpenAPI Entry

**Milestone.** M49
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M49.md
**Source specification.** tabilet/memory-bank/milestone.md#m49---n8n-public-api-openapi-entry
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M49 - n8n Public API OpenAPI Entry

**Goal.** Add n8n's own public REST API as a reviewed OpenAPI provider entry,
without treating the broader n8n node directory as runtime compatibility
evidence.

**Scope.**

- Promote the n8n Public API from the n8n API node priority signal into a
  durable catalog provider entry.
- Reference the official n8n repository OpenAPI 3.0 root and public REST API
  documentation.
- Save and register the ignored OpenAPI review artifact.
- Record that the repository OpenAPI root is unbundled with relative `$ref`
  entries, so the saved artifact is official source metadata but may need
  bundling before strict import.
- Classify n8n auth/security metadata from the root OpenAPI `ApiKeyAuth` and
  `BearerAuth` alternatives.
- Keep database, queue, mail, SQL, LDAP, FTP, and other protocol connector
  nodes out of the OpenAPI provider catalog unless a future native protocol
  milestone is opened.
- Do not call n8n APIs, resolve API keys or bearer tokens, choose instances,
  users, projects, credentials, workflows, executions, or data tables.

**Acceptance.** n8n appears in provider list/advisory/spec/security/stats views
as an official OpenAPI provider with a registered ignored review artifact,
complete auth metadata, and an explicit unbundled-spec quirk. Catalog quality
and diff checks pass.
````

## Status record

````markdown
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
````
