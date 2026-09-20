# Status M48 - Residual Google Discovery Coverage

| Item | State | Notes |
|---|---|---|
| Source review | `[+]` | Reviewed the official Google Discovery documents for People API (`people` v1), Perspective Comment Analyzer API (`commentanalyzer` v1alpha1), and Firebase Realtime Database Management API (`firebasedatabase` v1beta), plus the matching official REST documentation. |
| Candidate rows | `[+]` | Added Google Contacts / People API, Google Perspective, and Google Firebase Realtime Database as machine-spec candidates; added Google Ads as candidate-only review metadata with no Discovery/OpenAPI/Smithy promotion. |
| Provider rows | `[+]` | Promoted Contacts/People, Perspective, and Firebase Realtime Database Management as native Google Discovery provider entries with source refs, versions, revisions, and provider quirks. |
| Auth overlays | `[+]` | Added OAuth metadata overlays from Discovery scope declarations for the three promoted providers without resolving tokens, credentials, projects, or accounts. |
| Discovery artifacts | `[+]` | Saved ignored Discovery review artifacts under `catalog-openapi-cache/google-discovery/`, registered artifact paths in the local registry, and extended native Discovery artifact parser coverage tests. |
| Google Ads decision | `[+]` | Recorded Google Ads as a high-value candidate only. Current Google Ads REST docs describe REST/gRPC resources, but no official Google Discovery document, OpenAPI document, or Smithy model was found for M48 promotion. |
| Catalog snapshots | `[+]` | Updated deterministic provider, candidate, security, protocol-stat, alias, and classification snapshots for the new Google Discovery rows. |
| Review and verification | `[+]` | Deep review found no blocking follow-up fixes. Passed `go test ./...`, `go vet ./...`, `go run ./cmd/apitools catalog check`, focused list/spec/advisory/security/stats smoke checks for the three promoted providers, `go run ./cmd/apitools search --help`, `go run ./cmd/apitools import --help`, `git diff --check`, and `git -C ../tofu diff --check -- apitools`. |

## Boundary Notes

- M48 keeps Google Discovery native and does not reintroduce Discovery-to-OpenAPI
  conversion.
- Google Ads is not promoted until a future milestone explicitly scopes a
  protobuf/gRPC or Google API-specific source family.
- Firebase Realtime Database data paths remain documented as instance-specific
  REST behavior outside the management Discovery document.
- `apitools` must not execute Google APIs, resolve credentials, fetch tokens, or
  choose Google Cloud/Firebase projects.
