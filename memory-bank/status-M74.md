# Status M74 - Review Contract Corrections

**Acceptance.** All seven focused review findings are corrected without a new
compatibility break, their regressions are covered, and coordinated verification
passes.

| Item | State | Notes |
|---|---|---|
| Correct catalog security scope and restore discovery compatibility | `[+]` | Commit `8c40d36` keeps selected-spec auth unknown when only another spec has evidence, preserves the aggregate inspection report, restores both historical project-URL import signatures, and adds `ImportProjectURLsReport`; focused catalog/root tests pass. |
| Preserve valid protobuf and authoring inputs | `[+]` | Commit `03e4dc1` falls back to binary descriptors after a failed text heuristic, deep-copies nested ranking inputs, applies caller prompt budgets once during inventory construction, and marks selected-context overflow as truncated; focused parser/root tests pass. |
| Reject empty catalog artifacts | `[+]` | Commit `4557def` rejects zero-byte regular files before catalog-artifact persistence and verifies no poison row remains; focused SQLite cache tests pass. |
| Run milestone review and coordinated verification | `[+]` | Final review found no remaining cross-spec, compatibility, ownership, budgeting, truncation, parser-fallback, or empty-artifact gap. Uncached focused tests, full workspace/standalone apitools tests and vet, staticcheck, catalog generator/quality, CLI help, OpenUdon/Ramen/Udon workspace and standalone tests/vet, and apitools/tofu diff checks pass. |

## Boundary Checks

- Discovered metadata remains non-executable and credential-free.
- Remote and local source safety defaults are unchanged.
- Existing exported discovery method signatures remain source-compatible.
- No tfconfig changes are required.
