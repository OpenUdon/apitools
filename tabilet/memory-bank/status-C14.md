# C14 — Recheck publisher documentation for FND-04

**Owner.** APItools. **Dependency.** Accepted [C07](../docs/history/status-C07.md).
**Parallel-safe: no.** See the [milestone specification](milestone.md#c14--recheck-publisher-documentation-for-fnd-04).

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| C14.1 — Reconcile the 30 cases offline | `[+]` | Start with Emelia and Orbit, the two docs-only entries. Reconcile the original failure and current references, checked-in endpoint overlays, and recorded source gaps. Identify official pages and publisher-owned repositories for the supervised recheck. Check toolchain and configured cache readiness, and prepare the exact scope, supervision, bounds, and stop conditions for approval. |
| C14.2 — Supervised publisher recheck | `[ ]` | After C14.1 and explicit Class R approval, run one supervised recheck within the approved scope. Seek replacement documentation and official machine-readable specs. Prefer stable source links. Retain needed bytes in the configured cache and defer spec downloads when the trustworthy size is unknown or exceeds 1,000,000,000 bytes. Save each provider's checked sources, time, outcome, and provenance. |
| C14.3 — Qualify sources and record dispositions | `[ ]` | Fixes FND-04. Validate replacement specs and reproduce a representative operation offline before adoption. Refresh a checked-in endpoint overlay and its builder where official evidence supports a revised subset. Check applicable authentication evidence. Record a disposition for each of the 30 providers, identify any remaining gap by provider, run catalog and whole-unit checks, and reconcile FND-04 at closure. |

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

## Review

Review gate iterations: 0

Whole review pending with the default limit of 10. Use the [test policy](test-policy.md),
[verification-attempt gate](milestone.md#verification-attempt-gate), and
[milestone review procedure](milestone.md#milestone-review-procedure).
