# C14 — Recheck publisher documentation for FND-04

**Owner.** APItools. **Dependency.** Accepted [C07](../docs/history/status-C07.md).
**Parallel-safe: no.** See the [milestone specification](milestone.md#c14--recheck-publisher-documentation-for-fnd-04).

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| C14.1 — Reconcile the 30 cases offline | `[+]` | Start with Emelia and Orbit, the two docs-only entries. Reconcile the original failure and current references, checked-in endpoint overlays, and recorded source gaps. Identify official pages and publisher-owned repositories for the supervised recheck. Check toolchain and configured cache readiness, and prepare the exact scope, supervision, bounds, and stop conditions for approval. |
| C14.2 — Supervised publisher recheck | `[+]` | After C14.1 and explicit Class R approval, run one supervised recheck within the approved scope. Seek replacement documentation and official machine-readable specs. Prefer stable source links. Retain needed bytes in the configured cache and defer spec downloads when the trustworthy size is unknown or exceeds 1,000,000,000 bytes. Save each provider's checked sources, time, outcome, and provenance. |
| C14.3 — Qualify sources and record dispositions | `[+]` | Fixes FND-04. Validate replacement specs and reproduce a representative operation offline before adoption. Refresh a checked-in endpoint overlay and its builder where official evidence supports a revised subset. Check applicable authentication evidence. Record a disposition for each of the 30 providers, identify any remaining gap by provider, run catalog and whole-unit checks, and reconcile FND-04 at closure. |

## Approved scope

On 2026-11-27 the owner replied “yes” directly to recording the three-task
C14 plan. The [C07 qualification report](../../docs/catalog-qualification-c07.md#remaining-publisher-sources-and-overlays)
and retained [probe receipt](../../dist/c07-qualification.P2FWoi/docs-screen.json)
define the original population. The receipt SHA-256 is
`6882f2b9d2da808e496c810b46529905d77e3e8bc54356c58bef817ed9795e1f`.
Current catalog metadata shows 28 providers with checked-in endpoint overlays
and two with docs-only entries. This approval covers the planning updates.

| Original probe result | Provider IDs |
| --- | --- |
| Cross-origin redirect (6) | `aircall`, `home-assistant`, `mautic`, `orbit`, `tapfiliate`, `wufoo` |
| Bounded read failure (9) | `autopilot`, `clockify`, `demio`, `emelia`, `eventbrite`, `flow`, `invoice-ninja`, `todoist`, `twake` |
| HTTP 403 (10) | `aftership`, `agile-crm`, `clearbit`, `formstack`, `freshdesk`, `freshservice`, `mailjet`, `reddit`, `salesmate`, `splunk` |
| HTTP 404 (3) | `cockpit`, `erpnext`, `monica-crm` |
| HTTP 429 (1) | `contentful` |
| HTTP 502 (1) | `profitwell` |

## C14.1 offline reconciliation

The owner said “proceed” to C14.1 after the plan was recorded. The integrated
baseline is local commit `bcbb9c5b811d39f2d186ca6418f4874fb36e261b` plus the
two approved planning changes. The original
[roadmap delta](../../dist/c14-source-recheck/initial-planning.diff) has SHA-256
`f5fd319406d39580f0e34e0d131d3af7a6f6cbd3a4eaded4310fa2c9b0e05526` and the
[initial ledger](../../dist/c14-source-recheck/initial-status-C14.md) has
SHA-256 `258a0814c3a7e777fc3095506b52b13017979290632afea0b8ad73f1767f3527`.

The [recheck scope](../../docs/catalog-source-recheck-c14.md) reconciles every
sampled URL and current reference, names the 28 checked-in overlays and
their builders, and records the docs-only gaps for Emelia and Orbit. It
proposes one supervised C14.2 run with 69 starting URLs, exact GitHub
namespaces, guarded request and time limits, and provider-level evidence.
The scope document SHA-256 is
`866fe801c76c841768802c131ba56a25ae7989e687c7cd3a5da5945632fa060a`.

Offline readiness confirmed Go 1.26.6, the configured cache volume identity,
and successful scratch write and read. The checkout filesystem is ext4 for
disposable state. [Local readiness](../../dist/c14-source-recheck/c14.1-cache-toolchain.log).

Required gate `c14.1-offline-a1` starts against the exact baseline source
and the locally reconciled evidence. It runs the standard offline build,
vet, fresh whole suite with Chromium, generator and catalog quality checks,
CLI smoke checks, and whitespace validation. Gate limit 3, no charged
failure keys. Publisher access awaits separate approval.

Gate `c14.1-offline-a1` passed in 49 seconds on a clean checkout of the
baseline. Offline Go 1.26.6 `go build ./...`, `go vet ./...`,
`go test -v -count=1 ./...` with Chromium, generator freshness, catalog
quality, search/import/maintenance help, and `git diff --check` passed.
No Go or catalog data changed. [Gate log](../../dist/c14-source-recheck/c14.1-offline.log).

Independent GPT-6.2 Astra xhigh review caught an unspecified supervision
method in the proposed live scope. The corrected scope names the foreground
coordinator and how the run is observed and stopped. The re-review passed
with no open findings. The [reviewed scope revision](../../dist/c14-source-recheck/c14.1-reviewed-revision.diff)
has SHA-256
`3ae4649a2c12c0f65ad42bbbd2aef22981c3dd3f0e1a6542496d873ee9b119db`.
The unchanged build and test inputs reuse the passing offline gate, and
the corrected documentation and diff checks passed. C14.2 is pending
separate approval of the supervised run.

## C14.2 publisher recheck

On 2026-11-27 the owner replied “yes” to one supervised C14.2 live recheck
of the approved scope above. The scope binds the 30 named providers and 69
listed source URLs, starting with Emelia and Orbit, with 180 HTTP transactions,
20 minutes, and 4,000,000,000 response bytes in total. The approved run
defers specs with unknown trustworthy size or over 1,000,000,000 bytes.
Offline setup used `adef848004a0c0f7ce626a446ebd159cf60c9e3c`.
The isolated [runner](../../dist/c14-source-recheck/c14.2-jGNDPf/c14_recheck_test.go)
has source SHA-256 `c5cf6ad36a89b7e1f94f67ec54a6c84ab5accec21c8a0f2e0672e0abc52db10b`.
The foreground executable has SHA-256
`ac2542274234d1d4f4bc5bae686ae95058aa0cdcfac9cdb6871244ea65660069`.
The [offline fixtures](../../dist/c14-source-recheck/c14.2-jGNDPf/fixture.log)
passed the scope, redirect, throttling, request, and size checks. The
confirmed local source scope still matches the approved digest above.

Class R attempt `c14.2-live-a1`, 1 of 1, reserved before dispatch. Run
`c14-recheck.test -test.run '^TestC14Live$' -test.v -test.timeout=21m` in one
retained terminal with the scope, cache and fresh evidence directory for
this attempt. The runner imposes the approved 20-minute limit. Requests
are serial through the guarded transport. The coordinator observes each
provider and interrupts this foreground session with Ctrl+C on owner stop.
The evidence directory is
[c14.2-jGNDPf/evidence](../../dist/c14-source-recheck/c14.2-jGNDPf/evidence).

Attempt `c14.2-live-a1` passed in 36.45 seconds. All 30 providers have an
outcome. The terminal exited successfully after 64 GETs and 5,551,978
response-body bytes. The approved attempt is consumed. The final
[publisher recheck report](../../docs/catalog-source-recheck-c14-results.md)
records four pre-dispatch DNS failures, two Contentful deferrals after HTTP
429, 13 refused redirects, and 19 retained responses across 14 providers.
All 19 objects were independently rehashed from the configured cache. No
direct machine-readable artifact was acquired. Formstack's captured
overview contains a spec filename that remains an unverified source lead.
The [request journal](../../dist/c14-source-recheck/c14.2-jGNDPf/evidence/events.jsonl)
has SHA-256 `b205c09971795cfb7781316fc7ee12a626e4c3e58336a824d78fd0a581020b49`.
[Evidence and cache verification](../../dist/c14-source-recheck/c14.2-jGNDPf/evidence-verification.log).

The [offline gate](../../dist/c14-source-recheck/c14.2-jGNDPf/offline.log)
passed with Go 1.26.6 on the baseline in the disposable checkout. It ran
build, vet, fresh `go test -v -count=1 ./...`, generator freshness, catalog
quality (zero errors and warnings), search/import/maintenance help, and
whitespace checks. The full suite included the offline recheck fixtures and
skipped the separately authorized live entry point. No product or catalog
data changed. C14.3 owns offline qualification, any source and overlay
changes, and the final provider dispositions.

The initial full suite also skipped the optional browser fixture with its
browser variable unset. The fresh `./cmd/catalogsite` package run with
`APITOOLS_SITE_TEST_CHROME=/usr/bin/google-chrome` passed the loopback
browser fixture and the rest of the package.
[Browser gate](../../dist/c14-source-recheck/c14.2-jGNDPf/browser.log).

Independent GPT-6.2 Astra xhigh review passed with no remaining findings.
It reconciled the 69 listed URLs, all 30 outcomes, and the 19 cached objects.
Review clarified that both Todoist redirects pointed to its same-host
`/api/v1/` route. The runner deferred it based on the path filter. The
reviewed base was `adef848004a0c0f7ce626a446ebd159cf60c9e3c`, with initial
[diff](../../dist/c14-source-recheck/c14.2-jGNDPf/review.diff) SHA-256
`191602776447a7989cdca0645d0945bd01ecb8279aa034ce51ce39e236d60660`
and corrected [diff](../../dist/c14-source-recheck/c14.2-jGNDPf/review-followup.diff)
SHA-256 `e17b1cce8d880af6b9165c24ef7bb1d35f7926fd771dbb9061493bcf1489c55e`.

## C14.3 offline qualification

The owner selected C14.3 via `memory-bank-next` on 2026-11-27. The
starting checkout is `6699ba78dae56144b9cab3847d962b0f6ce435c7` with a clean tracked tree. This row uses
the repository [test policy](test-policy.md) and retained publisher evidence.
The supervised C14.2 approval is consumed.

Offline focused operation inventories reported the missing ERPNext current-user
route and the missing Invoice Ninja fixed request header before the changes.
Both checks passed after regenerating the two overlays. A baseline-built
`cataloggen` regenerated the identity files after the ordinary generator refused
the mismatched edited catalog bytes during initialization.
[Before](../../dist/c14-source-recheck/c14.3-ZwFOxp/focused-before.log).
[After](../../dist/c14-source-recheck/c14.3-ZwFOxp/focused-after.log).

Required offline gate `c14.3-offline-a1` started from the exact base above and
[implementation delta](../../dist/c14-source-recheck/c14.3-ZwFOxp/qualification-input.diff)
SHA-256 `67d90c948a26f799b23fc2c01189f4f73152f028631dc919314d74b7916876b0`.
The disposable checkout runs the Go build, vet, fresh whole suite with the
loopback Chromium fixture, builder reproducibility, generator freshness,
catalog quality, CLI smoke and whitespace checks. The default required-gate
limit is 3 per failing key, with no failing key charged at start.

Gate `c14.3-offline-a1` failed `go test -v -count=1 ./...` on
`catalog.TestProviderSearchAdvisoryOperationPins/erpnext-rest-api-overlay`
and `catalog.TestProviderSearchAdvisoryOperationPins/invoice-ninja-api-overlay`.
Both are the same signature at `catalog/search_profiles_test.go:181`, an
advisory artifact digest mismatch with the pinned search coverage metadata.
Each key has 1 charged attempt and 2 remaining. The new overlay bytes must
be pinned in the search profiles too. Both targeted pin checks passed after
updating the digests and adding the new ERPNext selector.
[Failed gate](../../dist/c14-source-recheck/c14.3-ZwFOxp/offline.log).
[Search pin checks](../../dist/c14-source-recheck/c14.3-ZwFOxp/search-pins-after.log).

Required offline gate `c14.3-offline-a2` started on the same base and command
set with [the corrected delta](../../dist/c14-source-recheck/c14.3-ZwFOxp/qualification-final-input.diff),
SHA-256 `05c04cab1dd650c69ccb74dfb07e7358eb7cba77d826f5b3b3a9107a689da6c5`.
The corrected delta also checks the Invoice Ninja header's singleton OpenAPI
enum as well as the required inventory parameter.

Gate `c14.3-offline-a2` passed on its frozen inputs, including the two pin
checks and the loopback browser fixture. Task review then identified stale
ERPNext public provenance and a missing Formstack registry lead. The ERPNext
advisory registration and search profile now point to the verified Frappe REST
reference. Provider, artifact and catalog reference identifiers are preserved.
The disposition distinguishes Formstack 2.0's `rest-api.json` registry entry
from the 2025.0 lead and retains both as unqualified.
[Passing intermediate gate](../../dist/c14-source-recheck/c14.3-ZwFOxp/offline-a2.log).

Required offline gate `c14.3-offline-a3` started with the same base and command
set on [the reviewed delta](../../dist/c14-source-recheck/c14.3-ZwFOxp/qualification-reviewed-input.diff),
SHA-256 `cad630182f39bf172936c6e0ca97ef8556f0dd57fa3c10d826a0e7ed010a0697`.

Gate `c14.3-offline-a3` passed in 64 seconds on the exact reviewed inputs.
Both overlays reproduced byte for byte from the builder. Go 1.26.6
build, vet, fresh whole-module tests with the loopback browser fixture,
generator freshness, catalog quality, CLI smoke and whitespace checks passed.
[Final gate](../../dist/c14-source-recheck/c14.3-ZwFOxp/offline-a3.log)
SHA-256 `e12c28334f7f417527de82a0bec6642963e8d6836462364f3f21331b84b5856b`.

The [30 offline dispositions](../../docs/catalog-source-recheck-c14-results.md#offline-source-qualification)
refresh ERPNext and Invoice Ninja, preserve the other 26 advisory overlays,
and preserve the docs-only decisions for Emelia and Orbit. All provider IDs
match the C07 cohort. Native adoption did not apply because no standalone
replacement qualified. Formstack and Todoist remain source leads. The
[local publisher facts](../../dist/c14-source-recheck/c14.3-ZwFOxp/publisher-facts.txt)
bind the qualifying source bodies and give the auth and operation evidence.

FND-04 originated in [C07](../docs/history/status-C07.md#status-record)
whole review 1 F2 and was deferred in [C13](../docs/history/status-C13.md#status-record).
C14.3 fixes the follow-up evidence gap through the supervised source recheck,
two source-backed corrections, and explicit dispositions for every provider.
The report preserves the provider-level gaps for future qualification.

Independent GPT-6.2 Astra xhigh task review passed the final delta above
against base `6699ba78dae56144b9cab3847d962b0f6ce435c7`. The stale search pins,
fixed-header enum coverage, ERPNext provenance and Formstack version evidence
were corrected. The reviewer matched all 30 cases to C07, rehashed all 19
retained source bodies, checked the source claims, and confirmed the final
offline gate. No findings remain for the task.

## Review

Review gate iterations: 2

Independent whole review round 1 started at 2026-11-27T10:23:20Z under the
default limit of 10. The integrated base is `bcbb9c5b811d39f2d186ca6418f4874fb36e261b`
plus the saved initial planning and ledger inputs above. The clean reviewed
HEAD is `358a60813150fde50a3e5ab639fd1363b72de023`. The full three-commit
[delta](../../dist/c14-source-recheck/c14.3-ZwFOxp/whole-review.diff) has SHA-256
`33aa22a914e308fe140439705ed37b8e63d9f2d412018c0b94da4dd47256058b`.
The independent review uses the selected Claude Opus 5.5 high role.

Whole review round 1 used a fresh Claude Code session with no tools. The
result identified model `claude-opus-5-5`; the invocation selected high effort.
It found no P0 or P1. P2-1 is resolved by the dashboard bookkeeping and
final FND-04 reconciliation. P2-2 is addressed in the report by marking the
approval scope historical and recording the exact `redirect-refused` outcome
for the unfetched Todoist target. The approved scope bytes retain their hash.
The report now explains the ERPNext search-reference association and binder
impact of the Invoice Ninja header. The overlay contains 16 operations, as
counted from its paths, correcting the review prose's count of 17. No other
review change was needed. The local gate and independent task review checked
the exact journal and cached sources beyond the packet's excerpt limits.
[Round 1 report](../../dist/c14-source-recheck/c14.3-ZwFOxp/whole-review-report.md)
SHA-256 `0c1639cb49968f708783d7a977e7f8474b52773a079e790248bf87a64d56aa09`.

Round 2 started at 2026-11-27T10:27:38Z over the documentation corrections
against the same reviewed HEAD. The [delta](../../dist/c14-source-recheck/c14.3-ZwFOxp/whole-review-r2.diff)
SHA-256 is `bc26caff8b85e1915ef859c9526083bad27c3a8773a48fdcbe53e1979ed76ed7`.
Whitespace, scope digest and operation count checks passed. The prior Class O
gate applies to the unchanged implementation and fixtures.

Whole review round 2 passed with `claude-opus-5-5` at high effort and no tools.
Both P2 findings were resolved. Declined Lower N1 because the closure ledger
retains an explicit report link to the per-provider remaining gaps and the
active FND-04 row is removed at retirement. Declined Lower N2 because the
preceding Invoice Ninja paragraph gives the exact `X-Requested-With:
XMLHttpRequest` header and identifies its fixed value.
[Round 2 report](../../dist/c14-source-recheck/c14.3-ZwFOxp/whole-review-r2-report.md)
SHA-256 `b574d8e4c82c11a302ecd2b26bd2373e46e7271752bb6418adf2f4ffae4510ff`.

FND-04 is fixed by the supervised recheck and C14.3's two overlay corrections
and 30 dispositions. Remaining provider evidence gaps are recorded in the
[C14 results](../../docs/catalog-source-recheck-c14-results.md#offline-source-qualification).
No active dependency needs a new contract. The accepted C07 baseline, C12 site
publication, consumer pins and the 334-provider count are separate. No
evolution bump after checking v32. This row refreshes catalog metadata within
the existing product and architecture boundaries.

The environment review found one bounded mechanical improvement candidate.
Overlay regeneration could verify the associated search coverage pins as part
of the same maintainer step. The catalog pin tests found and protected the
existing contract. The generator bootstrap dependency is documented in
[tech-stack.md](tech-stack.md#harnesses), and the distinction between partial
documentation and source qualification is consolidated in
[lessons.md](lessons.md#keep-incomplete-catalog-coverage-distinct-from-negative-evidence).

Use the [test policy](test-policy.md),
[verification-attempt gate](milestone.md#verification-attempt-gate), and
[milestone review procedure](milestone.md#milestone-review-procedure).
