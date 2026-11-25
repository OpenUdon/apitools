# C06 — Publisher refresh and storage

**Owner.** APItools. **Dependency.** Accepted C05. **Parallel-safe: no.**
The [milestone specification](milestone.md#c06--publisher-refresh-and-storage)
owns the storage, acquisition and adoption contract.

For this machine's local configuration, set `cache_dir` to
`/run/media/peter/94A8DB9CA8DB7AE6/apitools-cache`. Its expected volume
is UUID `94A8DB9CA8DB7AE6`, currently mounted from `/dev/sda3`. Shared
configuration examples and tests use portable directory paths.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| C06.1 — Capture the baseline and check storage readiness | `[+]` | Initial revision `150d08c644bd4ae4ce83a8772f317f083580d230`, 334 providers and 196 machine-readable references. Starting tracked diff and untracked hashes are in [baseline evidence](../../dist/c06-goal.zgalc0nf/baseline.json). Added cache preparation with a directory path, optional Linux volume identity, writability and capacity checks. Astra XHigh reviewed the per-parent mount and allocation-unit fixes. Fresh affected-package tests, offline build and vet, and diff checks passed. |
| C06.2 — Discover versions and preflight artifact size | `[+]` | Added publisher-scoped HEAD preflight, decimal GB and parser deferrals, checked redirects, and size and digest validation during staging. GitHub lookup retains the original branch or pin separately from the immutable file commit. The HubSpot index separates service, version and release stage. Astra XHigh accepted the scope and provenance fixes. Fresh affected-package tests, offline build and vet, and diff checks passed. |
| C06.3 — Persist the maintained catalog and retained bytes | `[+]` | Added a dedicated SQLite run registry on durable Linux storage and verified digest-addressed files on the configured artifact filesystem. Covered interrupted staging, orphan reuse, substituted paths, low space and transaction rollback. Astra XHigh accepted the storage and concurrency fixes. Fresh package and race tests, offline build and vet, and diff checks passed. |
| C06.4 — Compare reviewed contracts and stage export | `[+]` | Added bounded shape and native OpenAPI comparison for reviewed operations, including wire security and responses. SQLite preferences use provider, reference and service. Updates check the full previous identity and preserve prior pins. Export is deterministic and metadata-only. Astra XHigh accepted the native comparison, service identity, provenance and work-budget fixes. Fresh tests, storage race checks, offline build and vet, and diff checks passed. |
| C06.5 — Expose the update command and verify repeatability | `[+]` | Added configurable catalog update, dry-run, one-provider, all-providers, resume and metadata export. Offline CLI and updater fixtures covered no-write dry-run, all 334 checkpoints, mixed outcomes, resumed references, pinned baseline reconstruction and compatible adoption. The preference and reference checkpoint commit together, and first committed results survive a competing writer. Astra XHigh cleared re-review. Fresh package and race tests, offline build and vet, and diff checks passed. Whole verification and closure follow below. |

## Review

Review gate iterations: 0

Required verification recovery limit: 3 per key, from the approved launch.
C06.5 focused task package checks use fresh Go tests for the root, CLI and
storage packages, plus build, vet, storage race and diff checks.

Whole review pending, default limit 10. Use the [test policy](test-policy.md)
and [milestone review procedure](milestone.md#milestone-review-procedure).
