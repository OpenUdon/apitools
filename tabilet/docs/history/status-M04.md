# Retired milestone M04 - Security Overlays And Auth Classification

**Milestone.** M04
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M04.md
**Source specification.** tabilet/memory-bank/milestone.md#m04---security-overlays-and-auth-classification
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M04 - Security Overlays And Auth Classification

**Goal.** Add supplemental security-overlay metadata and auth/security
completeness classification for catalog providers.

**Scope.**

- Define security overlay schemas under `catalog`.
- Keep overlays separate from upstream/vendor specs and label their sources.
- Classify OpenAPI auth/security metadata as complete, present-incomplete,
  absent, overlay-required, intentionally-anonymous, or unknown.
- Support added or corrected `securitySchemes`, root-level or operation-level
  `security`, OAuth scopes, and API-key placement metadata.
- Do not resolve credentials, fetch tokens, sign requests, or choose runtime
  accounts.

**Acceptance.** Downstream callers can inspect auth/security completeness and
available overlays as metadata only, without treating overlays as upstream
provider truth or executing API operations.
````

## Status record

````markdown
# Status M04 - Security Overlays And Auth Classification

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Add metadata-only auth/security classification and supplemental security
overlays for catalog providers whose upstream specs are incomplete or ambiguous.

The feature should answer these questions without live provider calls:

- Which accepted specs include complete OpenAPI auth/security metadata?
- Which specs have auth/security metadata that is present but incomplete?
- Which specs need supplemental security overlays?
- What source notes explain why an overlay exists?

## Package Direction

Keep auth/security overlay behavior inside `catalog`:

```text
catalog/
  overlay.go
  security.go
  report.go
```

Security overlays should stay separate from upstream/vendor specs. A merged view
may be produced for inspection, but the overlay must remain provenance-labeled
metadata.

## Data Model Direction

Security overlays should cover at least:

- provider ID and optional spec/operation matching rules;
- added or corrected `securitySchemes`;
- root-level or operation-level `security` requirements;
- OAuth scopes or API-key/header/query/cookie placement;
- source notes explaining why the overlay exists.

Auth/security status values should include:

- complete;
- present-incomplete;
- absent;
- overlay-required;
- intentionally-anonymous;
- unknown.

## Tasks

| Item | State | Notes |
|---|---|---|
| Security overlay types implemented | `[+]` | Added metadata-only overlay, security scheme, OAuth flow, requirement, and operation match types under `catalog/`. |
| Auth completeness classifier implemented | `[+]` | Added deterministic OpenAPI-style security metadata classifier without credential or account handling. |
| Overlay validation implemented | `[+]` | Validates provider/spec references, scheme names, security requirements, operation matches, source refs, and source notes. |
| Seed overlay metadata added | `[+]` | Added source-backed overlays for Airtable, Gmail, Google Drive, HubSpot, OpenWeatherMap, and Slack. |
| Security report helpers added | `[+]` | Added deterministic built-in and caller-supplied security report helpers for downstream use. |
| Overlay tests added | `[+]` | Covers classifier states, overlay validation failures, copy safety, and deterministic built-in reports. |
| Verification completed | `[+]` | Ran required checks, sibling consumer tests, non-cached repo tests, and catalog race tests after review fixes. |

## Boundary Checks

- Do not execute API operations.
- Do not fetch credentials, tokens, or secrets.
- Do not sign requests.
- Do not choose accounts, tenants, workspaces, or runtime environments.
- Do not treat local overlays as upstream provider truth.
````
