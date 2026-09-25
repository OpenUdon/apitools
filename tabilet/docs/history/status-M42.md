# Retired milestone M42 - NVIDIA DSX Air Provider Entry

**Milestone.** M42
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M42.md
**Source specification.** tabilet/memory-bank/milestone.md#m42---nvidia-dsx-air-provider-entry
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M42 - NVIDIA DSX Air Provider Entry

**Goal.** Add NVIDIA DSX Air as a reviewed provider catalog entry using its
official OpenAPI schema and source-backed authentication metadata.

**Scope.**

- Add NVIDIA DSX Air candidate and provider records.
- Save the official DSX Air OpenAPI schema as a tracked review artifact under
  `catalog-openapi-cache/openapi/`.
- Register the schema in the catalog artifact registry.
- Add a security overlay for NGC API-key bearer authentication because the
  current official schema omits reusable security scheme metadata.
- Do not create API keys, call DSX Air operations, choose organizations, or
  manage DSX Air simulations.

**Acceptance.** Catalog provider, security, advisory, stats, refresh, test/vet,
and diff checks pass with DSX Air discoverable through provider lookup and
refreshable OpenAPI metadata.
````

## Status record

````markdown
# Status M42 - NVIDIA DSX Air Provider Entry

| Item | State | Notes |
|---|---|---|
| Candidate and provider metadata | `[+]` | Added NVIDIA DSX Air as a reviewed provider with aliases, official docs, and official OpenAPI schema reference. |
| OpenAPI review artifact | `[+]` | Saved the official DSX Air OpenAPI 3.0.3 schema from `https://dsx-air.nvidia.com/api/schema/` under `catalog-openapi-cache/openapi/`. |
| Auth overlay | `[+]` | Added NGC API-key bearer authentication overlay because the current official schema omits reusable security scheme metadata. |
| Artifact registry | `[+]` | Registered the DSX Air schema path in the catalog artifact registration helper. |
| Verification | `[+]` | `go test ./catalog`, `go test ./...`, `go vet ./...`, `go run ./cmd/apitools catalog check`, DSX Air advisory/overlay/spec/refresh checks, `git diff --check`, and tofu harness diff checks passed. |

## Boundary Notes

- DSX Air API-key creation, NGC organization selection, role/scope assignment,
  simulation management, worker registration, and API operation execution are
  outside `apitools`.
- The catalog records metadata and review artifacts only; downstream runtime
  code remains responsible for credential binding and execution policy.
````
