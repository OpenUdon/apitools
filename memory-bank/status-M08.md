# Status M08 - OpenUdon Intent Advisory Integration

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Use `apitools/catalog` from `../openudon` to improve optional `intent.hcl`
authoring and review guidance. Catalog metadata should help users discover
known provider/spec/security context, but it must not become authoritative
workflow behavior.

The feature should answer these questions in OpenUdon:

- Does a workflow intent mention a provider known to the apitools catalog?
- Is there a known official spec reference or a likely need for user OpenAPI?
- What auth/security status or overlay note should authoring show?
- Did an explicit project/user OpenAPI input override catalog metadata?

## Integration Direction

Start with OpenUdon intent advisory. Reporting-only or apitools-only expansion
is deferred unless implementation discovers a blocker.

## Tasks

| Item | State | Notes |
|---|---|---|
| OpenUdon integration point identified | `[+]` | Review evidence now receives optional catalog advice after OpenAPI candidates/discovery and before inferred data flow. |
| Catalog adapter added in OpenUdon | `[+]` | `internal/workflowintent` exposes a small adapter over `apitools/catalog` without duplicating catalog records downstream. |
| Provider matching implemented | `[+]` | Intent provider and OpenAPI hints match catalog IDs, display names, aliases, and common spec filename stems deterministically. |
| Advisory output added | `[+]` | Review evidence can surface provider ID, spec reference, user OpenAPI need, auth/security status, overlay IDs, and source notes. |
| Explicit input precedence preserved | `[+]` | Intent/project OpenAPI paths are recorded as explicit overrides so catalog specs remain advisory. |
| No-match fallback implemented | `[+]` | Unknown providers and disabled advice render no catalog section and do not affect review generation. |
| OpenUdon tests added | `[+]` | Tests cover advisory presence, no-match fallback, explicit-input precedence, disabled behavior, filename matching, and review rendering. |
| Cross-repo verification completed | `[+]` | Ran OpenUdon, apitools, and udon checks required by AGENTS.md, including public-module OpenUdon tests. |

## Boundary Checks

- Do not move workflow semantics into `apitools`.
- Do not make catalog metadata mandatory for OpenUdon workflows.
- Do not bind credentials, choose accounts, sign requests, or execute workflows.
- Do not change trusted-runner handoff semantics.
- Do not treat catalog advice as a release gate unless a later milestone
  explicitly defines that policy.
