# Retired milestone M72 - Network, Storage, Cache, And OAuth Boundary Hardening

**Milestone.** M72
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M72.md
**Source specification.** tabilet/memory-bank/milestone.md#m72---network-storage-cache-and-oauth-boundary-hardening
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M72 - Network, Storage, Cache, And OAuth Boundary Hardening

**Goal.** Close network and filesystem integrity gaps and move live OAuth token
exchange into the runtime-owned Udon credential boundary.

**Scope.** Reject unsafe URL credentials, ports, and special address ranges;
enforce wire and decoded download bounds; use one confined atomic write/read
layer for refresh, materialization, export, and cache artifacts; separate raw
and corrected refresh validation; require cache digests and bounded pruning;
and replace the apitools OAuth command with a state/PKCE-protected Udon command
that writes secrets only to an exclusive mode-0600 file.

**Dependencies.** M71. Udon migration and apitools OAuth removal are one
coordinated row. M72 file ownership is disjoint from active S01 work except for
serialized documentation and dependency updates.

**Parallel ownership.** M72 owns `download.go`, refresh, catalog materialize,
SQLite cache, apitools OAuth removal, and Udon credentials/CLI. S01 owns parser,
inventory, authoring, and security summary contracts.

**Downstream impacts.** Udon becomes the sole owner of authorization-code
exchange and refresh-token handling. Materialization and refresh callers gain
explicit collision and raw/corrected validation outcomes.

**Acceptance.** Default-safe downloads reject unsafe hosts, rebinding, ports,
userinfo, redirects, and decompression abuse; all artifact I/O is confined,
digest-checked, atomic, and rollback-safe; cache growth is bounded; apitools has
no token exchange path; Udon never emits secrets to stdout/stderr; focused,
full, standalone, downstream, vet, staticcheck, CLI, and diff checks pass.
````

## Status record

````markdown
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
````
