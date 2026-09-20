# Status M72 - Network, Storage, Cache, And OAuth Boundary Hardening

**Acceptance.** Safe network policy is enforced through connection time,
artifact I/O is confined and integrity checked, cache growth is bounded,
refresh truthfully distinguishes raw from corrected validation, OAuth token
exchange lives only in Udon, and milestone verification passes.

| Item | State | Notes |
|---|---|---|
| Harden network destinations and response budgets | `[+]` | URL userinfo, unsafe ranges, ports, redirects, and dial-time rebinding fail closed; reviewed added ports retain address checks; wire and decoded bodies are independently bounded. |
| Introduce confined atomic artifact I/O | `[+]` | Shared artifact I/O confines and verifies reads, atomically replaces refresh files, and commits materialize/export trees with identical reuse, collision errors, explicit-force backup/rollback, and workflow-root containment. |
| Correct refresh and SQLite cache integrity contracts | `[+]` | Refresh and review reports separate immutable raw provenance from in-memory corrected validation; SQLite schema v3 confines path-backed rows, requires digest/byte evidence, bounds row and byte growth, and prunes deterministically by LRU. |
| Move Google OAuth authorization into Udon | `[+]` | Udon now owns refresh and a state/PKCE loopback login with exclusive synchronized 0600 credential output; apitools no longer imports oauth2, exchanges tokens, or exposes an OAuth command. Workspace, standalone, race, downstream, vet, quality, CLI, and diff gates pass; only documented pre-existing staticcheck findings remain. |
| Review and verify M72 | `[+]` | Independent review removed the last two staticcheck findings and stale OAuth-helper documentation. Focused race tests plus full workspace/standalone tests, Udon quality/vet/CLI checks, OpenUdon and Ramen suites, apitools vet/staticcheck/CLI checks, and diff checks pass. The coordinated evolution v22 direction remains current. |

## Boundary Checks

- Additional allowed ports do not bypass address classification.
- Raw downloaded bytes remain immutable provenance; corrected interpretations
  are never mislabeled as the stored artifact.
- apitools never resolves credentials, fetches tokens, or prints secrets.

## Review Notes

- Default-safe requests are revalidated on redirects and at dial time; custom
  transports still require the explicit unsafe-host test escape hatch.
- Refresh writes preserve downloaded bytes and digest as raw provenance while
  corrected validation remains separately labeled in memory and reports.
- Materialization verifies registered digests and publishes one provider or
  workflow tree transactionally; cache rows are confined and deterministically
  pruned under finite defaults.
- Udon owns credential parsing, PKCE authorization, refresh exchange, and
  private token-file output. No OAuth implementation or dependency remains in
  apitools.
