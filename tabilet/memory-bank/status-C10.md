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
| C10.2 — Reproduce and adopt revised catalog pins | `[ ]` | **F01 (P1, confirmed), F06 (P1, confirmed).** F01 extends C06 whole review 1's P2-1. If a revised catalog pin matches the held artifact, verify the owned cached bytes and locator before using it as the comparison baseline. Reproduce a reviewed immutable file at its authored commit even when the publisher has moved on. Test the D1 → D2 catalog repin followed by compatible D3, migrated v2 preferences, absent and existing preferences, and references with and without reviewed operations. Keep older pins from replacing a newer preference and condition preference/checkpoint writes on the expected prior state. [Adoption and baseline reconstruction](../../catalog_update.go), [version storage](../../internal/catalogupdate/versions.go). |
| C10.3 — Recognize the same resource and identical bytes | `[ ]` | **F07 (P2, confirmed), F08 (P2, confirmed).** The locator comparison can reject a path with multiple version-like components even when compared with itself. Test exact locators and the same GitHub file at different commits, including version-like owner/repo names, Confluent paths and Pipedrive paths. Distinct origins, services, file paths or meaningful queries must remain distinct. For the same reviewed resource and verified digest at a new commit, report unchanged content or a provenance-only update, including when reviewed operations are empty. Preserve validation, security checks and the newer preference. [Resource comparison](../../catalog_update_compatibility.go), [reference identity](../../catalog_update_sources.go), [unchanged decision](../../catalog_update.go). |
| C10.4 — Retry temporary publisher failures within a bounded budget | `[ ]` | **F04 (P2, confirmed, fixes FND-01), F05 (P2, partially confirmed), F10 (P2, confirmed).** Preserve safe evidence for rate limits, retryable server errors and timeouts. Support bounded targeted retries and per-origin backoff that lets unrelated providers finish and deferred references progress across attempts. Use offline 403/429, 5xx, timeout, restart and quota fixtures. Verify atomic checkpoints, retry hints and exhaustion reporting. The catalog exposes 96 eligible raw GitHub references (88 at the revalidation commit). The live cutoff was unmeasured. Use the effective client on GitHub metadata probes and verify nil-receiver behavior with isolated transport. [Checkpoint and resume](../../catalog_update.go), [GitHub resolution](../../catalog_update_sources.go), [client and retry evidence](../../api_versions_probe.go). |
| C10.5 — Preserve failure evidence and inspect historical runs | `[ ]` | **F09 (P1, confirmed), F15 (P2, confirmed).** Early failures can record zero artifact check times and an empty artifact URL even though the source carries the publisher URL. Persist UTC attempt and lifecycle times as observed, the safe original locator and any known resolved locator. Distinguish skipped and unobserved events in SQLite and export. Exercise failed discovery and preflight, additive migration of v2/v3 state and deterministic export. Leave unknown historical values unknown. Read completed outcomes using their original run identity when current inputs change, diagnose an unknown ID, and keep identity checks on incomplete resumes. Verify historical inspection is read-only. [Run identity and evidence](../../catalog_update.go), [download preflight](../../catalog_update_download.go), [SQLite run records](../../internal/catalogupdate/store.go). |

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

Review gate iterations: 0

Whole review pending, limit 10. This post-retirement intake leaves C06's
accepted six-round review intact. Follow the [test policy](test-policy.md) and
[milestone review procedure](milestone.md#milestone-review-procedure).
