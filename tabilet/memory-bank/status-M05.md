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
