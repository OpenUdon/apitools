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
