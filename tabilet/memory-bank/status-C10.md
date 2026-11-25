# C10 — Catalog updater repairs before the first inventory

**Owner.** APItools. **Dependency.** Accepted [C06](../docs/history/status-C06.md).
**Parallel-safe: no.** C07 consumes the repaired updater for the first
publisher inventory. See the [milestone specification](milestone.md#c10--catalog-updater-repairs-before-the-first-inventory).

## Review source and baseline

The owner approved these planning changes from **C06 code review (xhigh)**
(`c06-code-review-xhigh.md`, SHA-256
`b9ee0181b8ce548eccfd21160138977412dec6583b51c5bab11f921053b6cd7a`).
The reviewed range starts at `150d08c644bd4ae4ce83a8772f317f083580d230`
and ends at `cdd16c0295b6c74f739de145eeb5c3fd4487bd84`, also the
revalidation revision. Updater source matched that revision. Counts below
use the working catalog of 334 providers, including the uncommitted C04/C05
entries and overlay registry. The committed catalog contained 316 providers.
Finding IDs follow the source order. The report labels F01–F10 Important
and F11–F15 Nit. Each task below gives the local severity.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| C10.1 — Validate reviewed artifacts and read native versions | `[+]` | **F02 (P1, partially confirmed), F03 (P2, confirmed).** Exact-pin adoption and reconstruction retain raw bytes and record validation evidence, including eight pinned parseable-invalid routes and a scoped correction. Square remains invalid for updater adoption because its pinned document contains external Square CDN schema references. The updater reads native version metadata with bounded dates, aliases and merges, while compatibility still requires complete original operation shapes and the original wire contract. Astra xhigh cleared the task after alias work accounting and strict scalar fixes. Fresh affected tests with the staged pins, offline build/vet and diff checks passed. [Validation](../../catalog_update_validation.go), [fixture evidence](../../catalog_update_validation_test.go). |
| C10.2 — Reproduce and adopt revised catalog pins | `[+]` | **F01 (P1, confirmed), F06 (P1, confirmed).** A revised pin can reuse owned cached bytes after digest, native-format and locator checks. The updater can independently reproduce an authored GitHub commit when the head has moved, redirected or failed. Baseline and head evidence stay separate. Acceptance and checkpoint use the original preference as their compare-and-swap condition. Digest history prevents replay, including equal-version artifacts. Fixtures cover migrated preferences, cache and exact locator proof, moved sources, protobuf, Discovery and competing first acceptance. Astra xhigh passed re-review. Fresh affected tests, storage race, offline build and vet, and diff checks passed. [Pin fixtures](../../catalog_update_pins_test.go). |
| C10.3 — Recognize the same resource and identical bytes | `[+]` | **F07 (P2, confirmed), F08 (P2, confirmed).** Exact URL paths and two immutable commits of the same GitHub file preserve resource identity. Generic version-like name normalization remains inventory grouping. Origin, query and distinct file boundaries still prevent compatibility acceptance. Verified identical bytes at a new commit report unchanged with observed provenance and preserve the prior preference, even without reviewed operations. Fixtures cover Confluent, Pipedrive, divergent redirect destinations and prior cache recovery. Astra xhigh passed re-review. Fresh affected offline package tests, build, vet and diff checks passed. [Resource fixtures](../../catalog_update_compatibility_test.go), [update fixtures](../../catalog_update_pins_test.go). |
| C10.4 — Retry temporary publisher failures within a bounded budget | `[+]` | **F04 (P2, fixes FND-01), F05 (P2, partially confirmed), F10 (P2).** Classifies publisher quotas, 5xx responses and timeouts across discovery, head and authored retrieval. Incomplete runs preserve separate origin deadlines and progress outside throttled origins. Resume uses exact-checkpoint CAS, three attempts and a 24-hour recovery window with explicit terminal deferral for longer hints. Fixtures cover 403/429/5xx, timeout, restart, exhausted and saturated hints, stale-checkpoint rollback, successful nil-client metadata, and resumed 334-provider CLI inventory. The live quota cutoff remains unmeasured. Astra xhigh passed re-review. Fresh affected tests, storage race, offline build, vet and diff checks passed. [Retry fixtures](../../catalog_update_retry_test.go). |
| C10.5 — Preserve failure evidence and inspect historical runs | `[+]` | **F09 (P1, confirmed), F15 (P2, confirmed).** Run lifecycle and reference attempts record UTC observations, including failed discovery, guarded redirect targets, authored checks and DNS preflight. Skips omit unobserved attempts. Schema v4 preserves historical unknowns. Completed reports read the original SQLite snapshot without opening the cache, retaining original identity and reporting unknown IDs plainly. Unfinished runs still check the catalog and config identity. Fixtures cover v2/v3 inspection and migration, deterministic export and fragment exclusion. Astra xhigh passed re-review. Fresh offline build, vet, full suite, storage and focused DNS race, catalog generation/quality, CLI help, 334-provider dry run and diff checks passed. [Run history fixtures](../../catalog_update_history_test.go), [legacy SQLite fixtures](../../internal/catalogupdate/run_history_test.go). |

## Verification and review

Execution baseline is `fae6ced44866a7370e59bf00ff3b20c124ad4f3b` with the
committed 334-provider catalog. C09 retains the original C06.1 integrated
review baseline. This launch uses ten whole review rounds and three
verification recovery attempts per required key.

Use offline publisher fixtures and isolated database and artifact directories.
Run the row regressions and the whole-unit checks in the
[specification](milestone.md#c10--catalog-updater-repairs-before-the-first-inventory).
Keep checks on guarded destinations, size deferral, parser bounds, immutable
locator ownership and concurrent updates in the regression suite.

Review gate iterations: 1

Whole review limit 10. This post-retirement intake leaves C06's
accepted six-round review intact. Follow the [test policy](test-policy.md) and
[milestone review procedure](milestone.md#milestone-review-procedure).

Whole review 1 started at `2026-11-25T14:11:44Z` with Claude Opus 5.5 at high effort.
Reviewed state `b423f3dcaaf34df6093af1c5391e1d0f8b3836d7` against execution baseline
`fae6ced44866a7370e59bf00ff3b20c124ad4f3b`. The
[review diff](../../dist/c10-goal.tZ4ZvJ/review-1.diff) has SHA-256
`d5eea440b081c1162aef570af7fbd995abc9fb8a396d70b1646857ac26e0b05f`. Fresh [offline build/vet](../../dist/c10-goal.tZ4ZvJ/static-acceptance.log),
[full tests](../../dist/c10-goal.tZ4ZvJ/tests-acceptance.log),
[storage race](../../dist/c10-goal.tZ4ZvJ/race-acceptance.log) and
[catalog checks](../../dist/c10-goal.tZ4ZvJ/catalog-acceptance.log) passed on the
final implementation. The catalog check reported zero errors and warnings.
The no-write dry run covered all 334 providers.

Whole review 1 interrupted at `2026-11-25T14:12:51Z`. Claude Code selected
`claude-opus-5-5` with an explicit high-effort launch and rejected its
configured API key before returning a review. The [attempt receipt](../../dist/c10-goal.tZ4ZvJ/review-1-stream.jsonl)
records the model and failure. Resume this same round with restored
authentication or an explicitly approved reviewer change. C10 remains
open pending the independent review, and C07 still depends on C10 acceptance
and its separate publisher-run approval.

Whole review 1 resumed at `2026-11-25T14:15:40.958Z` in Herdr APItools tab 2 at
the owner's request. The existing Claude Code session dispatched a fresh
read-only child with an explicit high effort selection. The
[dispatch receipt](../../dist/c10-goal.tZ4ZvJ/review-1-herdr-dispatch.json)
records the native selection. The implementation, review baseline and
iteration counter are the same as in the original attempt. The native child transcript confirms Claude Opus 5.5 at high effort.

Whole review 1 returned at `2026-11-25T14:22:12.747Z`. The read-only reviewer
covered all five lenses and blocked on lost adoption evidence after a retry.
[Report](../../dist/c10-goal.tZ4ZvJ/review-1-herdr-report.md).
Confirmed from current code and classified locally as P1 because an accepted
pin must remain visible in durable completed-run evidence. The same-digest
immutable locator move is a separate P2. The report also identifies
lower-priority cases and unverified observations. Revalidate those before
applying dispositions. The counter remains 1 until this review's required
fixes pass fresh verification and round 2 is started.

Round 1 dispositions. Preserve an adopted baseline and its authored artifact
when a pending retry is replaced. Keep the first verified object linkage
when the retry supplies none. An unchanged head after an adoption in the
same run reports the adoption. A terminal head failure retains its failure
status and the distinct adoption evidence. Fixtures read both completed
outcomes from SQLite. This fixes local P1 from reviewer P2-1.

The locator repair for P2-2 permits an independently reproduced same-digest
move from a held immutable GitHub locator. Unproven moves and replay of a
standing pin over a newer preference still go to review. The transport
review added a final GET-locator check and uses the artifact URL guard
inside HEAD redirects (U1). The index inventory skips single-service
preference lookup (P3-1). The GitHub seconds-based reset hint applies only
to metadata GitHub and must fit the representable range (P3-2).

A disposable interrupted-writer fixture reproduced SQLite's
`SQLITE_READONLY_ROLLBACK` (U2). Ordinary completed-run inspection stays
read-only. A resume supplied with the maintenance configuration may use
the guarded writer to recover the journal before the saved identity check
and any publisher request. The state-database-only path reports the needed
recovery without opening the writer.

P3-3 keeps bounded 5xx recovery and the existing HEAD-specific 501 size
deferral. P3-4 keeps the existing JSON-shaped key contract and explicit
validation refusal for non-response non-string keys. P3-5 retains the
separate guarded response and staged-publisher fixtures and the storage
CAS test. The additional fault combinations are optional coverage.
For U3 the nil client is normalized before the probe reads configuration.
For U4 the non-spec body flag still releases the source budget. For U5
dry-run and export branches bypass the pending-run text. The narrowed
round 2 packet includes the directly consulted URL guard, CLI dispatch
and database transaction setup (U6).

Round 1 repairs passed fresh [offline build and vet](../../dist/c10-goal.tZ4ZvJ/static-round2.log),
[full Go suite](../../dist/c10-goal.tZ4ZvJ/tests-round2.log),
[storage and focused retry/journal race checks](../../dist/c10-goal.tZ4ZvJ/race-round2.log)
and [catalog, CLI and 334-provider dry-run checks](../../dist/c10-goal.tZ4ZvJ/catalog-round2.log).
The package-local CLI exercises the affected report and resume path.
No Go callers of the maintenance API appeared in the OpenUdon or Udon
consumer checkouts.
