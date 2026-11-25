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

Review gate iterations: 6

Required verification recovery limit: 3 per key, from the approved launch.
C06.5 focused task package checks use fresh Go tests for the root, CLI and
storage packages, plus build, vet, storage race and diff checks.

Whole verification attempt C06-whole-1 passed at `420fe56`. Offline build,
vet, fresh full tests with the pinned C04/C05 artifacts, fresh storage race
tests, catalog generation and quality checks, CLI help and diff checks
passed. The all-provider dry run returned 334 entries. Evidence is saved
under [closure](../../dist/c06-goal.zgalc0nf/closure).

Whole review round 1 started under the approved limit of 10. Review the
complete C06 change against the recorded baseline. The independent reviewer
is Claude Opus 5.5 at high effort.

Round 1 is interrupted. Claude Code rejected the configured API key before
returning any model output. Resume round 1 with a fresh Claude Opus 5.5 High
session after authentication is repaired. C06 closure and C07 depend on
this review. The scoped review packet and authentication failure are saved
in the closure evidence directory.

Round 1 resumed and completed in the owner's designated Herdr tab 2 with
Claude Opus 5.5 at high effort. The [review receipt](../../dist/c06-goal.zgalc0nf/closure/review-1-herdr-receipt.md)
identifies two pending acceptance fixes. Honor a newly reviewed catalog pin
when its exact bytes reproduce. Keep an unavailable cached artifact from
blocking the rest of the provider run. Review round 2 will cover these
findings, the repair diff and its direct consumers.

The owner selected GPT-6.2 Astra at high effort for C06 whole review from
round 2 onward. Round 1 remains part of the approved ten-round allowance.
Launch reference refreshed for that selection from
`bd67aef5b4798e572e107647b2d901611d501bd57234209414bad12db741ab11` to
`bbb10c571c5a6ddb5abc005718fcc52666f8e4858c78fb4eb3db908bbcc5be02`.

Whole verification attempt C06-whole-2 passed against
`420fe5645411c85ca2f5cccd5a6b6c914d0eba3d` plus the [saved repair diff](../../dist/c06-goal.zgalc0nf/closure/round2/working-tree.diff)
with SHA-256 `a32d03fba370872edfb755a2e8977b78809a531c3183a4c832ee26b0cf79f8b0`.
Fresh pinned whole tests, storage race, build, vet, catalog generation and
quality, CLI help, 334-provider dry run and diff checks passed. Independent
round 2 started with GPT-6.2 Astra High against that implementation state.

Round 2 returned three acceptance findings. Keep the accepted service key
after a reviewed locator move, distinguish a new catalog review from the
standing pin, and checkpoint a damaged cache object during immutable
baseline reconstruction. Round 3 covers these findings and their repairs.

Whole verification attempt C06-whole-3 passed at
`420fe5645411c85ca2f5cccd5a6b6c914d0eba3d` with the [saved repair diff](../../dist/c06-goal.zgalc0nf/closure/round3/working-tree.diff)
SHA-256 `cb2afe6fca34c4314712d59e4899675390a48a442d2b874e2b7a53c56f9f6343`.
The fresh pinned suite, race, build, vet, generation and quality, CLI help,
dry run and diff checks passed. Round 3 started with GPT-6.2 Astra High.

Round 3 verified locator continuity and cache recovery. Two acceptance
findings remain for upgraded state. Prove strict release order when the
previous catalog-pin identity is unknown, and keep schema v2 export readable
without migration. A newly authored pin must reproduce before a later
compatible release can advance its preference.

Whole verification attempt C06-whole-4 passed at
`420fe5645411c85ca2f5cccd5a6b6c914d0eba3d` with the [saved repair diff](../../dist/c06-goal.zgalc0nf/closure/round4/working-tree.diff)
SHA-256 `44ec92a21a43be95b30de374c509942bce028b7f9040efd4650f39c5facbf2dd`.
Fresh pinned tests, storage race, build, vet, catalog generation and quality,
CLI help, the 334-provider dry run and diff checks passed. Round 4 started
with GPT-6.2 Astra High to review the round 3 repairs and their consumers.

Round 4 confirmed equal-version replay protection and v2 read-only export.
Two acceptance fixes remained. A migrated preference must reproduce a new
pin before accepting later publisher bytes. A moved immutable locator must
keep one preference and queue review when the resolved head fails to
reproduce the pin. The scoped review also found that an immutable cache
hit could hide a new pin and that a concurrent migration could split
the read-only export between two snapshots. The repairs make cache hits
honor the reviewed pin and read both export queries in one transaction.

Revalidation of round 1's SQLite locking concern reproduced loss of a
write lock when a second store closed its file-inspection descriptor.
Use a Linux path-only descriptor for existing state and serialize the
initial file creation. The offline write-lock test fails before that
repair and passes after it. Catalog update and migration fixtures passed.

Whole verification attempt C06-whole-5 passed at
`420fe5645411c85ca2f5cccd5a6b6c914d0eba3d` with the [saved repair diff](../../dist/c06-goal.zgalc0nf/closure/round5/working-tree.diff)
SHA-256 `e334668d739f27bbca1190eba4f545a3ddd4f506443449c86cc84a332f7e76f8`.
Fresh pinned tests, storage race, build, vet, catalog generation and quality,
CLI help, the 334-provider dry run and diff checks passed. GPT-6.2 Astra High
started round 5 on the round 4 pin and export repairs and the verified
state-file lock repair.

Round 5 confirmed the round 4 fixes and the state-file lock repair. It
found one remaining roll-back path. A locator-only change could make the
unchanged catalog digest look newly reviewed after a compatible advance.
The standing-pin guard now compares the anchor across locator changes.
The moved-locator fixture checks that the older digest is queued for
review and exactly one preference still points to the newer bytes.

The round 5 direct-consumer checks used immutable snapshots and scratch
module overrides. OpenUdon `36ae38dfe93327a0879929c2cfe5dbbbbf9b988f`
passed its fresh suite. Udon `fdff9f37722ac6a2ab050ef112b4c3521b6b2dc5`
returned the same `invalid_plan` failures in browser fixtures in
`pkg/execute` and `pkg/runner` with the candidate and the original
APItools pin. Failure names were identical after removing timings.
The source-adapter and command packages passed.

Whole verification attempt C06-whole-6 passed at
`420fe5645411c85ca2f5cccd5a6b6c914d0eba3d` with the [saved repair diff](../../dist/c06-goal.zgalc0nf/closure/round6/working-tree.diff)
SHA-256 `6a66db134527a8fd03016a9921d094d17c68f53cf5fc0ba44009c156aa3d8e64`.
The focused fixture, fresh pinned whole tests, storage race, build, vet,
catalog generation and quality, CLI help, 334-provider dry run and diff
checks passed. Round 6 started with GPT-6.2 Astra High for the replay
guard across locator changes.

Round 6 passed. Astra High confirmed that a standing digest stays barred
from replacing a newer preference after a locator move. A genuinely new
reviewed pin can still advance the preference, and the export sees the
single preferred row. All earlier acceptance findings are closed. The
effective limit was 10 under the approved launch and 6 rounds were used.

Round 1's remaining advisory findings carry forward as FND-01 (transient
failure evidence and retry), FND-02 (partial reports), and FND-03
(portable and integrated fixtures). The write transaction contention,
mount-identity documentation, redirect parameters and descriptor locking
findings were fixed and reviewed within C06.

The review exposed two useful checks for future maintenance work. Test a
changed locator with both a newly authored digest and a standing digest.
Exercise SQLite migrations with old read-only consumers and verify that
file-identity checks leave active advisory locks intact. The fixtures in
C06 now exercise these cases. No new evolution version is needed. C06
implements the publisher-refresh direction recorded in v31.
