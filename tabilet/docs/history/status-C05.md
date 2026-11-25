# Retired milestone C05 - Qualify the remaining nine catalog expansion candidates

**Milestone.** C05
**Outcome.** completed
**Retired.** 2026-11-25
**Source status.** tabilet/memory-bank/status-C05.md
**Source specification.** tabilet/memory-bank/milestone.md#c05--qualify-the-remaining-nine-catalog-expansion-candidates
**Evidence.** 150d08c644bd4ae4ce83a8772f317f083580d230
**Worktree.** includes uncommitted changes
**Review.** passed
**Review iterations.** 2
**Verification.** Offline Go 1.26.6 build, vet, 17-package fresh suite with 21 pinned artifacts, catalog generation, quality (0 errors and 0 warnings), CLI help, security report and audit, stats, refresh report, and diff check passed after the Google Ads RPC overlay fix. The original C05 gates also verified byte equality from all five new overlay builders. Closure reuses those results because their code, tests, and source inputs are unchanged. [Combined qualification evidence](../../../docs/catalog-expansion-c05.md).
**Consolidated into.** [Product](../../memory-bank/product.md), [architecture](../../memory-bank/architecture.md), and [catalog coverage lesson](../../memory-bank/lessons.md#keep-incomplete-catalog-coverage-distinct-from-negative-evidence).

## Milestone specification

````markdown
### C05 — Qualify the remaining nine catalog expansion candidates

**Owner and dependencies.** APItools, after accepted
[C04](../docs/history/status-C04.md). APItools remains the single serial
execution owner. **Parallel-safe: no.**

**Scope.** Review Google Maps Platform, Google Ads, eBay, Etsy, Square, Mercado
Pago, Razorpay, WeChat Pay, and Twitch. Check a current provider-owned source
for the selected Google Maps API surface because Google deprecated the
OpenAPI repository named in the earlier proposal. For Google Ads, qualify
the existing candidate with an official, scoped protobuf artifact and prove
the catalog-to-UWS path for a concrete RPC before promotion. C05.1 also owns
the built-in preferred-reference selection in `catalog/resolve.go` and its
regression tests. With both protobuf and human-doc references, default
`catalog inspect` must select the protobuf spec ID and report `grpc-protobuf`
regardless of declaration order. Check that `catalog resolve` lists that spec
with the expected source type and that existing reference priorities and
explicit caller overrides still pass their tests.

**Admission and acceptance for C05.**

- Record the exact API surface and selected operation, provider-owned source
  reference, applicable version or revision, artifact digest, provenance, and
  license or redistribution terms. Use the existing cache and registration
  path for provider artifacts. Keep reference-only bytes out of tracked files.
  For advisory overlays, record the official documentation passages and the
  covered operation subset.
- Review authentication and authorization against official documentation.
  For OpenAPI and Swagger, inspect both reusable security schemes and
  effective root and operation requirements. Add source-backed catalog
  security overlay metadata wherever the provider artifact leaves a gap.
- Before promoting a provider, register the pinned artifact in a disposable
  catalog root and build its operation index. Run offline catalog discovery
  with `ProviderKeys` scoped to that provider and a step contract for the
  reviewed operation. Assert that the returned candidate identifies the
  expected provider, raw source SHA-256, source type and native selector.
  Require the catalog to resolve and materialize the same artifact, then
  reproduce the selector from the materialized bytes. A skipped artifact or
  a metadata-only lead leaves the discovery requirement unmet.
  An OpenAPI index needs a qualifying child document. A docs link, snippet, or
  unqualified model alone leaves the provider outside the reviewed count.
- In every row that changes built-in catalog data, regenerate the catalog
  with `go run ./cmd/cataloggen`, include the generated manifest and lookup
  indexes in the same change, and run `go run ./cmd/cataloggen -check`
  before closing the row. The batch closing rows reconcile the README counts.
- Give each named candidate an evidence-backed disposition. Promote only
  qualifying entries. Evidence gaps retain a candidate disposition and reason. The upper
  bound of 334 reviewed providers assumes all 18 qualify.

**Write set for implementation.** Built-in catalog data and generated outputs,
catalog asset registration, any required advisory overlays and builders,
focused catalog and source verification, README coverage counts, and owning
memory-bank records. The scope uses existing source families and interfaces.
Downstream products retain their own admission, runtime and pin decisions.

**Verification.** Use offline, credential-free checks with
pinned inputs. Verify the resolved, materialized source and operation selector
for every promoted provider, and check that unqualified candidates remain
outside the reviewed provider set. Before each whole-milestone review, run
`go run ./cmd/cataloggen -check`, the `catalog check`, `security-report`,
`security-audit`, `stats`, and `refresh-report` commands with the applicable
local registry, `go test -count=1 ./...`, `go vet ./...`, and `git diff --check`.
Reconcile the README coverage counts with catalog data and tracked overlays.
Follow the [test policy](test-policy.md) and the ordinary persisted milestone
review. The planning approval leaves source acquisition, execution and
publication to their applicable separate authority. Verification does not
call provider operations or use live credentials.

**Status and integrated closure.** [C05](status-C05.md) owns its nine provider
reviews and the final combined disposition for all 18. Verify the combined
catalog and README against both approved batches. An unqualified candidate
retains its review reason outside the promoted set.
````

## Status record

````markdown
# C05 — Qualify the remaining nine catalog expansion candidates

**Owner.** APItools, one serial execution owner.
**Dependency.** Accepted [C04](../docs/history/status-C04.md).
**Planning approval.** Owner approved the two-batch proposal on 2026-11-24
with “yes”. The owner separately invoked the memory-bank Goal for C04 then C05.
**Specification.** [C05](milestone.md#c05--qualify-the-remaining-nine-catalog-expansion-candidates).
**Planning baseline.** `150d08c644bd4ae4ce83a8772f317f083580d230`.
Implementation starts from the accepted C04 state, with 325 reviewed providers
and the reusable pinned-artifact qualification test. The integrated review
uses the original C04/C05 starting state.

## Tasks

Apply the C04 per-promotion offline discovery check to each qualified operation.
Each row that changes built-in data regenerates the manifest and indexes and
runs `go run ./cmd/cataloggen -check` in that change.

| Item | State | Notes |
| --- | --- | --- |
| C05.1 — Review maps, advertising, and marketplace candidates | `[+]` | Promoted Places API (New) Autocomplete, Google Ads v26 ListAccessibleCustomers, and eBay Sell Account v1 GetOptedInPrograms. Both Google Ads reference orders reproduce the fix. Five pinned artifacts passed discovery and materialized selector reproduction. |
| C05.2 — Review commerce and payments candidates | `[+]` | Promoted Etsy v3 shop holiday preferences, Square Connect v2 RetrieveMerchant, and Mercado Pago Get payment by ID. All three provider-owned artifacts passed scoped discovery and materialized selector reproduction. |
| C05.3 — Review remaining payments and media candidates | `[+]` | Promoted Razorpay Payments v1, WeChat Pay Global v3 common-mode Query Order, and Twitch Helix Get Games with three generator-backed, source-reviewed advisory artifacts. Pin, discovery, materialization, and UWS binding passed. |
| C05.4 — Verify the combined catalog expansion | `[+]` | All 18 candidate dispositions reconciled. All 21 pinned artifacts passed. README and measured catalog agree on 334 providers, distributed 120 / 54 / 139 / 21. Required offline gates passed. |

### C05.1 evidence

- Decision, digests, auth and terms in the [qualification record](../../docs/catalog-expansion-c05.md).
- Row PASS on base `150d08c644bd4ae4ce83a8772f317f083580d230` and [saved diff](../../dist/c04-c05-goal.6LlCVA/c05/row1.diff) SHA-256 `9e5b019e59049454f0e4539ef58cadb82235e592f37aa0f14351ebf1568161b1`. Catalog SHA-256 `a9ddc4f31986450e16a75c4bb6b7d03cbe5952a7088528d4a1fa59eea84f363d`.
- `go build ./...`, `go vet ./...`, `git diff --check`, `go run ./cmd/cataloggen -check`, and `go test -v . ./catalog ./cmd/apitools ./cmd/cataloggen` with the exact artifacts. Root, CLI and generator test results reused their cache for unchanged inputs. [Row check](../../dist/c04-c05-goal.6LlCVA/c05/row1-check.log).
- [Both red reference orders](../../dist/c04-c05-goal.6LlCVA/c05/proto-priority-red.log), [actual artifact tests](../../dist/c04-c05-goal.6LlCVA/c05/row1-proof.log), [local refresh report](../../dist/c04-c05-goal.6LlCVA/c05/row1-refresh.json).

### C05.2 evidence

- [Qualification record](../../docs/catalog-expansion-c05.md#commerce-and-payments) lists the exact artifacts, terms, and source-backed auth decisions.
- Row PASS on base `150d08c644bd4ae4ce83a8772f317f083580d230` and [saved diff](../../dist/c04-c05-goal.6LlCVA/c05/row2.diff) SHA-256 `70022392959bfac1fdb368334551799ea27ac8e76c5fdf21f1105e7b225d1290`. Catalog SHA-256 `c3ede6a964843982ba7b46b8f73dd12f5d792e06542a7a820df4397bbfa8d8af`.
- `go build ./...`, `go vet ./...`, `git diff --check`, `go run ./cmd/cataloggen -check`, and `go test -v . ./catalog ./cmd/apitools ./cmd/cataloggen` with the exact artifacts. [Row check](../../dist/c04-c05-goal.6LlCVA/c05/row2-check.log).
- [Actual artifact tests](../../dist/c04-c05-goal.6LlCVA/c05/row2-proof.log), [offline refresh report](../../dist/c04-c05-goal.6LlCVA/c05/row2-refresh.json). Etsy and Mercado Pago validate strictly. Square's original source reports `invalid` under strict validation, and its selected merchant lookup passes bounded discovery and native inventory reproduction.

### C05.3 evidence

- [Qualification record](../../docs/catalog-expansion-c05.md#payments-and-media-with-advisory-endpoints) and the [endpoint decision table](../../catalog-openapi-cache/advisory-overlays/human-doc-overlay-decisions.md) record the supported subsets and source-backed terms and auth.
- Row PASS on base `150d08c644bd4ae4ce83a8772f317f083580d230` and [saved diff](../../dist/c04-c05-goal.6LlCVA/c05/row3.diff) SHA-256 `4f8641a5709ecf89c6d806330749fa7b22eddc167c66c6d06e5cac7776192344`. Catalog SHA-256 `2215f1b165acb64f777d2ec9efb841e7f05c9d34b02e47cd74d44a0e9a6d431f`.
- `go build ./...`, `go vet ./...`, `git diff --check`, `go run ./cmd/cataloggen -check`, and `go test -v . ./catalog ./cmd/apitools ./cmd/cataloggen` with the exact artifacts. The CLI and generator used cached results with unchanged inputs. [Row check](../../dist/c04-c05-goal.6LlCVA/c05/row3-check.log).
- [Actual artifact tests and typed selector reproduction](../../dist/c04-c05-goal.6LlCVA/c05/row3-proof.log). Etsy, Mercado Pago and eBay additionally reproduced their selected OpenAPI operations in the UWS shape resolver.

### C05.4 evidence

- [Combined provider dispositions](../../docs/catalog-expansion-c05.md#combined-c04-and-c05-disposition) and [README coverage](../../README.md#reviewed-provider-catalog).
- Final gate base `150d08c644bd4ae4ce83a8772f317f083580d230`, [saved diff](../../dist/c04-c05-goal.6LlCVA/c05/final-gate.diff) SHA-256 `04719a2db57eb470db4efe7cbc1fffbe25fd4960054c77861b17b1f77d36bcce`, catalog SHA-256 `2215f1b165acb64f777d2ec9efb841e7f05c9d34b02e47cd74d44a0e9a6d431f`.
- PASS for `go build ./...`, `go vet ./...`, `git diff --check`, `go run ./cmd/cataloggen -check`, and `go run ./cmd/apitools catalog check`, `search --help`, and `import --help`. [Static and CLI gate log](../../dist/c04-c05-goal.6LlCVA/c05/final-static.log).
- PASS for `go test -count=1 -v -timeout=75s ./...` with `APITOOLS_CATALOG_QUALIFICATION_DIR=dist/c04-c05-goal.6LlCVA/c05/source` and the pinned offline Go 1.26.6 toolchain. All 21 pinned-artifact cases ran and passed across 17 tested Go packages. [Full Go test log](../../dist/c04-c05-goal.6LlCVA/c05/final-tests.log).
- PASS for `go run ./cmd/apitools catalog security-report --json` and applicable-registry `catalog security-audit`, `catalog stats` and `catalog refresh-report`, using the scratch `source/cache.sqlite` and `source` artifact root. `stats` and `refresh-report` used `--as-of 2026-11-24`. There were 21 registered artifacts, 15 pinned upstream files and six tracked advisory artifacts. Legacy report commands updated SQLite access timestamps within the scratch copy. The final registry SHA-256 is `1abd50c6ce63e67a1750aef498458dcafeda1c9ffc933ceceda741f2207d9b97`. [Security report](../../dist/c04-c05-goal.6LlCVA/c05/final-security-report.json), [security audit](../../dist/c04-c05-goal.6LlCVA/c05/final-security-audit.json), [catalog statistics](../../dist/c04-c05-goal.6LlCVA/c05/final-stats.json), [refresh report](../../dist/c04-c05-goal.6LlCVA/c05/final-refresh-report.json).
- PASS for running the five C05 Go builders in an isolated output directory and comparing the resulting bytes with the advisory artifacts. [Regenerated artifact digests](../../dist/c04-c05-goal.6LlCVA/c05/final-builders.log).

## Review

Review gate iterations: 2

Use the [test policy](test-policy.md) and
[milestone review procedure](milestone.md#milestone-review-procedure).
Iteration 1 started at 2026-11-25T04:00:15Z under the default limit of 10.
The integrated scope starts from the original C04/C05 baseline at
`150d08c644bd4ae4ce83a8772f317f083580d230` plus the initial working-tree
evidence. The verified final state uses the same HEAD plus the
[final gate diff](../../dist/c04-c05-goal.6LlCVA/c05/final-gate.diff),
SHA-256 `04719a2db57eb470db4efe7cbc1fffbe25fd4960054c77861b17b1f77d36bcce`.
The closing row's fresh gates cover this exact implementation and are
reused for the review.

Iteration 1 by GPT-6.2 Astra at xhigh found one P1 in the integrated scope.
Google Ads' scoped security overlay used a dotted operation ID. Security
inspection could not match it to the slash-form ID emitted by the native
parser. The first review found no other issues across the five review lenses.
Match the native ID, verify the combined OAuth and developer-token
requirements against the indexed RPC, and reverify for iteration 2.

Iteration 2 started at 2026-11-25 04:05:23 UTC. The review covers the P1 and the fix's direct
consumers. [Saved diff](../../dist/c04-c05-goal.6LlCVA/c05/review1-fix/review2.diff)
SHA-256 `8037355abc30e209f7bbb7c5579407dc2fbea6166313deee623111e66529f690`
on the same HEAD. The catalog digest is
`fcfcd0c6d82cbde994c72e1e3d4f7df445b178adb8f4028e8b4d269ffb206cdd`.
The targeted native RPC regression failed on the dotted ID and passed
after the catalog used the parser's ID. Fresh offline build, vet, whole Go
tests with all 21 pinned artifacts, generated output, CLI help, catalog
quality (0 errors and 0 warnings), security and refresh reports, stats and
diff checks passed.
[Red](../../dist/c04-c05-goal.6LlCVA/c05/rpc-security-red.log) ·
[Green](../../dist/c04-c05-goal.6LlCVA/c05/rpc-security-green.log) ·
[Full test log](../../dist/c04-c05-goal.6LlCVA/c05/review1-fix/full-tests.log).

Fresh GPT-6.2 Astra at xhigh passed iteration 2. The Google Ads native RPC
matches one security requirement set containing both the OAuth scope and
developer-token header. There are no remaining findings. The integrated
C04/C05 review accepts all 18 promotions and 21 pinned artifacts.
No evolution bump is needed for this catalog expansion and resolver correction.
````
