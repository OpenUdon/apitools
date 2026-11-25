# Retired milestone C04 - Qualify the first nine catalog expansion candidates

**Milestone.** C04
**Outcome.** completed
**Retired.** 2026-11-25
**Source status.** tabilet/memory-bank/status-C04.md
**Source specification.** tabilet/memory-bank/milestone.md#c04--qualify-the-first-nine-catalog-expansion-candidates
**Evidence.** 150d08c644bd4ae4ce83a8772f317f083580d230
**Worktree.** includes uncommitted changes
**Review.** passed
**Review iterations.** 2
**Verification.** Offline Go 1.26.6 build, vet, fresh whole tests with ten pinned proofs, generator, catalog quality (0 errors and 0 warnings), CLI help, security reports and audit, stats and refresh report, and diff check passed on the review 2 snapshot. Artifact identities and durable procedure are in [the C04 evidence](../../../docs/catalog-expansion-c04.md). The closure reuses those gates because their code, tests and source inputs are unchanged.
**Consolidated into.** [Product](../../memory-bank/product.md), [verification standards](../../memory-bank/tech-stack.md), and [catalog coverage lesson](../../memory-bank/lessons.md#keep-incomplete-catalog-coverage-distinct-from-negative-evidence).

## Milestone specification

````markdown
### C04 — Qualify the first nine catalog expansion candidates

**Approval and owner.** The owner approved the bounded two-batch catalog plan
on 2026-11-24 (“yes”). APItools owns both milestones, with one serial execution
owner and C04 first. **Parallel-safe: no.** Planning baseline
`150d08c644bd4ae4ce83a8772f317f083580d230` includes a pre-existing README
working-tree edit recording the 316-provider baseline. Capture the live
baseline and preserve that edit when starting implementation.

**Scope.** Review DigitalOcean, Vercel, Datadog, MongoDB Atlas Administration,
Google Gemini API, Amazon Selling Partner API, Plaid, Miro, and Hugging Face
Hub. Prefer a provider-published source document supported by APItools and UWS.
When official docs support the selected operation, a checked-in docs-derived
advisory OpenAPI overlay and its tracked builder can supply missing coverage.
[Catalog curation](../../docs/catalog-curation.md) governs asset and security
review. The earlier [expansion proposal](../../docs/catalog-upgrade-2026-09.md)
remains evidence for the source leads and the separate unscheduled work.

**Planning review.** On 2026-11-24 the owner approved the three C04/C05
planning review findings. C04.1 owns the first offline discovery proof
(High, local P1), and every promotion row applies it. C05.1 owns the Google
Ads preferred-reference fix (Medium, local P2). Each catalog-changing row
owns its generated outputs (Medium, local P2).

**Admission and acceptance for C04 and C05.**

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

**Verification for C04 and C05.** Use offline, credential-free checks with
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

**Status and dependencies.** [C04](status-C04.md) owns the nine provider
reviews and its closing evidence. No prerequisite milestone. C05 follows
accepted C04 and owns the final combined review and count reconciliation.
````

## Status record

````markdown
# C04 — Qualify the first nine catalog expansion candidates

**Owner.** APItools, one serial execution owner. No prerequisite milestone.
**Planning approval.** Owner approved the two-batch proposal on 2026-11-24
with “yes”. The owner subsequently launched the ordered C04 and C05 goal.
**Specification.** [C04](milestone.md#c04--qualify-the-first-nine-catalog-expansion-candidates).
**Planning baseline.** `150d08c644bd4ae4ce83a8772f317f083580d230`.
The pre-existing README coverage table is part of the starting context.

## Tasks

Each promotion row must discover its expected operation by provider in a
disposable offline catalog root and match its digest and selector to the
materialized artifact. Each row that changes built-in data regenerates the
manifest and indexes and runs `go run ./cmd/cataloggen -check` in that change.

| Item | State | Notes |
| --- | --- | --- |
| C04.1 — Review cloud and observability candidates | `[+]` | Promoted DigitalOcean account metadata, Vercel project domains and Datadog on-demand concurrency cap. [Artifact and auth evidence](../../docs/catalog-expansion-c04.md#cloud-and-observability). The reusable offline test verifies the SQLite registration, scoped discovery and exact materialized selector. |
| C04.2 — Review database, AI, and marketplace candidates | `[+]` | Promoted MongoDB Atlas project list, Gemini Developer API model list, and Amazon Sellers v1 marketplace participation. The Gemini selector reproduced through the local UWS resolver from Discovery bytes, and the separate docs-derived endpoint overlay passed the same offline discovery and materialization checks. [Sources and auth review](../../docs/catalog-expansion-c04.md#database-ai-and-marketplace). |
| C04.3 — Review banking, collaboration, and model hub candidates | `[+]` | Promoted Plaid processor token permission read, Miro specific board, and Hugging Face Hub model repository refs. [Sources, selectors and auth evidence](../../docs/catalog-expansion-c04.md#banking-collaboration-and-model-hub). |
| C04.4 — Verify the batch and reconcile catalog coverage | `[+]` | All nine promoted providers and the separate Gemini advisory overlay passed the pinned offline checks. Reconciled the README to 325 providers, with primary coverage of 116 official OpenAPI/Swagger/index, 52 other official machine-readable, 136 human docs with endpoint overlays, and 21 human docs without endpoint overlays. C05 can extend the same pinned-artifact test table. |

## Review

Review gate iterations: 2

Iteration 1 started at 2026-11-25T03:33:36Z under the default limit of 10.
Review input is `150d08c644bd4ae4ce83a8772f317f083580d230` plus the
[saved working-tree diff](../../dist/c04-c05-goal.6LlCVA/c04-review1.diff)
with SHA-256 `94e7d9dde034b4a02aa726b45609470461dc5e3f9c262b0ae2de476c9ca43ae2`.

Use the [test policy](test-policy.md) and
[milestone review procedure](milestone.md#milestone-review-procedure).
Fresh whole verification passed with offline Go 1.26.6, including build,
vet, `go test -count=1 -v ./...` with all ten pinned operation proofs,
`cataloggen -check`, both CLI help checks, catalog quality (0 errors and
0 warnings), security report and audit, stats and refresh report, and diff
check. [Full test log](../../dist/c04-c05-goal.6LlCVA/c04-full-tests.log).

Iteration 1 by GPT-6.2 Astra at xhigh reported two findings. P1 concerns
provenance and acceptance. The pinned MongoDB OpenAPI declares CC BY-NC-SA
3.0 US and service terms, which need to appear alongside the separate SDK
license. P2 concerns the advisory schema. Gemini accepts page sizes above
1,000 and caps the response, and the model version is required. Record the
corrected evidence and reverify affected gates before iteration 2.

Iteration 2 started at 2026-11-25 03:38:10 UTC. The review covers the two fixes and their
direct consumers. Input is the same HEAD plus [saved diff](../../dist/c04-c05-goal.6LlCVA/c04-review2.diff),
SHA-256 `4a4a41760d4263b3e67ff312086212d7c9495ac1007b766700b0d957568b5ac0`.
The fresh whole check set passed after separating the MongoDB document license
from the SDK license and rebuilding the Gemini overlay with the correct
pagination behavior and required model version. Catalog quality reported
0 errors and 0 warnings. [Full test log with ten pinned proofs](../../dist/c04-c05-goal.6LlCVA/review1-fixes/full-tests.log).

Fresh GPT-6.2 Astra at xhigh passed iteration 2 across correctness,
security and provenance, compatibility, verification, and acceptance.
Both findings are resolved. C04 accepted nine promotions and the
separate Gemini advisory overlay. C05 inherits the exact admission gate,
reusable promotion test, and pinned C04 artifacts for integrated review.
No evolution bump is needed for additive catalog curation.

## Row verification

- C04.1 PASS on `150d08c644bd4ae4ce83a8772f317f083580d230` plus [saved diff](../../dist/c04.1tC087/row1.diff), SHA-256 `7ed3ee6d7263f0389f7e418339ddcc67541c47cb0d7e2e62ae79761f49594a4b`. Offline Go 1.26.6 `go build ./...`, `go vet ./...`, `APITOOLS_CATALOG_QUALIFICATION_DIR=dist/c04.1tC087/source go test -timeout 90s -v . ./catalog ./cmd/apitools ./cmd/cataloggen`, `go run ./cmd/cataloggen -check`, and `git diff --check`. [Log](../../dist/c04.1tC087/row1-check.log).
- C04.2 PASS on the same `HEAD` plus [saved diff](../../dist/c04.1tC087/row2.diff), SHA-256 `fa6ad0d165a20ea65a7cff7cb8808a9ee5f4a2d478406d0e90a19a5b9542b064`. The same offline build, vet, whole affected packages, generator check and diff checks passed with the new pinned artifacts and generated overlay. [Log](../../dist/c04.1tC087/row2-check.log).
- C04.3 PASS on the same `HEAD` plus [saved diff](../../dist/c04.1tC087/row3.diff), SHA-256 `bfdc69d2e40d3716ccf3b0b418e3482f55e9ba8764aad20a9d89541981ccd756`. The same offline build, vet, whole affected packages, generator check and diff checks passed with nine pinned sources and the Gemini overlay. [Log](../../dist/c04.1tC087/row3-check.log).
- C04.4 PASS on the same `HEAD` plus [saved diff](../../dist/c04.1tC087/row4.diff), SHA-256 `4658381310e93dd6f3b7e6b60f1c1809c6726a91522af83f87579da14af6f7ac`. Offline Go 1.26.6 `go build ./...`, `go vet ./...`, the saved artifact registration builder compile, `APITOOLS_CATALOG_QUALIFICATION_DIR=dist/c04.1tC087/source go test -timeout 60s -v . ./catalog ./cmd/apitools ./cmd/cataloggen`, `go run ./cmd/cataloggen -check`, `go run ./cmd/apitools catalog check`, and `git diff --check` passed. The generated catalog digest is `6dc84c29e01ac28766da85c35e38fc0e939f45fde62cd0d4a0b0a23e8169f477`. [Log](../../dist/c04.1tC087/row4-check.log).
- The offline `catalog security-report`, `catalog security-audit`, `catalog stats`, and `catalog refresh-report` commands passed against the applicable local inputs. The audit, stats and refresh commands used `--cache-dir dist/c04.1tC087/source --cache dist/c04.1tC087/source/cache.sqlite`. [Security report](../../dist/c04.1tC087/row4-security-report.json), [security audit](../../dist/c04.1tC087/row4-security-audit.json), [stats](../../dist/c04.1tC087/row4-counts.json), [refresh report](../../dist/c04.1tC087/row4-refresh-report.json), [README coverage counts](../../dist/c04.1tC087/row4-coverage.json).
````
