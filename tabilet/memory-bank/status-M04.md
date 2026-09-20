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
