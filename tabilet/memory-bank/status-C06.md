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
| C06.3 — Persist the maintained catalog and retained bytes | `[ ]` | Keep the SQLite database on native Linux storage. Deduplicate retained files outside Git by verified digest and link multiple provider references to the same bytes. Cover mutable links, offline retention, failed/truncated writes, lying size headers, disk errors and crash recovery. |
| C06.4 — Compare reviewed contracts and stage export | `[ ]` | Compare source identity, operations, required inputs, response contracts and effective auth. Adopt compatible changes transactionally in SQLite. Keep prior pins accessible. Queue ambiguous, invalid, breaking or newly authored changes for review. Verify identical export from identical state and enforce redistribution terms. |
| C06.5 — Expose the update command and verify repeatability | `[ ]` | Exercise dry-run, one-provider, full and resume modes through the CLI. Prove dry-run leaves persistent state alone. Use fixture publishers for unchanged, changed, deferred, partial and failed cases. Run the package gates and hand the accepted update contract to C07. |

## Review

Review gate iterations: 0

Whole review pending, default limit 10. Use the [test policy](test-policy.md)
and [milestone review procedure](milestone.md#milestone-review-procedure).
