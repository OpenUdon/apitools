# Retired milestone M05 - Catalog Resolution And CLI Reporting

**Milestone.** M05
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M05.md
**Source specification.** tabilet/memory-bank/milestone.md#m05---catalog-resolution-and-cli-reporting
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M05 - Catalog Resolution And CLI Reporting

**Goal.** Expose catalog resolution and reporting through stable library helpers
and CLI inspection commands.

**Scope.**

- Implement resolution precedence for user-provided OpenAPI paths/URLs,
  user-provided overlays, project-local documents, built-in spec references,
  and built-in overlays.
- Add deterministic reporting helpers for provider catalog and security status.
- Provide CLI inspection/reporting commands, such as `catalog list`,
  `catalog inspect`, and `catalog security-report`.
- Update the public README to explain the catalog as metadata only.
- Preserve existing exported API compatibility unless a breaking change is
  intentional and documented.

**Acceptance.** Operators and downstream callers can inspect catalog metadata,
resolve metadata with explicit local inputs taking precedence, and generate
security reports without live provider calls, operation execution, credential
resolution, or n8n compatibility claims.
````

## Status record

````markdown
# Status M05 - Catalog Resolution And CLI Reporting

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Expose catalog metadata through stable resolution helpers and CLI reporting so
operators and downstream projects can inspect provider/spec/security status.

The feature should answer these questions without live provider calls:

- What provider metadata is available from built-in catalog data?
- What local user-provided spec or overlay overrides built-in metadata?
- What auth/security status should be reported for a provider?
- Which providers still need user-supplied OpenAPI documents or security review?

## Resolution Precedence

When downstream callers resolve provider metadata, use this precedence:

1. user-provided OpenAPI path or URL;
2. user-provided security overlay;
3. project-local OpenAPI document;
4. built-in catalog spec reference;
5. built-in catalog security overlay.

The catalog is a convenience and discovery baseline. It must not override a
team's explicit local API contract.

## CLI Direction

Suggested commands:

```text
apitools catalog list
apitools catalog inspect <provider>
apitools catalog security-report
```

The CLI should describe metadata and source notes only. It must not execute API
operations, resolve credentials, or claim provider runtime compatibility.

## Tasks

| Item | State | Notes |
|---|---|---|
| Resolution precedence implemented | `[+]` | Added `catalog.ResolveProvider` with user OpenAPI, user security overlay, project-local OpenAPI, built-in spec, and built-in overlay precedence. |
| Catalog list command added | `[+]` | Added deterministic `apitools catalog list` text and JSON output. |
| Catalog inspect command added | `[+]` | Added provider inspection with spec references, source notes, resolution source, and auth/security status. |
| Catalog security-report command added | `[+]` | Added auth completeness/security overlay reporting across built-in providers. |
| Public README updated | `[+]` | Documented catalog commands, metadata-only boundary, and resolution precedence. |
| CLI and resolver tests added | `[+]` | Covers precedence, deterministic output, JSON report output, and command behavior. |
| Verification completed | `[+]` | Ran required checks, catalog CLI smoke checks, sibling consumer tests, non-cached repo tests, and catalog race tests after review fixes. |

## Boundary Checks

- Do not execute API operations.
- Do not fetch credentials, tokens, or secrets.
- Do not sign requests.
- Do not choose accounts, tenants, workspaces, or runtime environments.
- Do not import n8n node runtime semantics.
- Do not let built-in catalog metadata override explicit user inputs.
````
