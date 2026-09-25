# Retired milestone M48 - Residual Google Discovery Coverage

**Milestone.** M48
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M48.md
**Source specification.** tabilet/memory-bank/milestone.md#m48---residual-google-discovery-coverage
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M48 - Residual Google Discovery Coverage

**Goal.** Add the remaining n8n-visible Google services that have official
Google Discovery documents, and record why Google Ads is not promoted as a
first-class source entry in this batch.

**Scope.**

- Promote People API / Contacts, Perspective Comment Analyzer API, and
  Firebase Realtime Database Management API as native Google Discovery catalog
  providers.
- Register ignored Google Discovery review artifacts and extend parser
  coverage checks so the artifacts remain native Discovery inputs, not OpenAPI.
- Add source-backed OAuth overlay metadata from Discovery scope declarations.
- Add Google Ads as a reviewed candidate-only row with an explicit
  no-Discovery/no-OpenAPI promotion decision.
- Preserve the Firebase distinction between Discovery-covered management
  resources and instance-specific Realtime Database data REST paths documented
  outside the Discovery document.
- Do not execute Google APIs, resolve credentials, fetch tokens, choose Google
  Cloud/Firebase projects, or add protobuf/gRPC source-family support.

**Acceptance.** Contacts/People, Perspective, and Firebase Realtime Database
Management appear in provider list/advisory/spec/security views as native
Google Discovery sources with registered ignored artifacts and OAuth overlays;
Google Ads remains candidate-only with its promotion blocker documented; catalog
quality and diff checks pass.
````

## Status record

````markdown
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
````
