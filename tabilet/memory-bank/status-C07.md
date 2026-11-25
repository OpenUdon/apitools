# C07 — First inventory and provider curation

**Owner.** APItools. **Dependency.** Accepted [C10 updater repairs](../docs/history/status-C10.md) following [C06](../docs/history/status-C06.md). **Parallel-safe: no.**
See the [milestone specification](milestone.md#c07--first-inventory-and-provider-curation)
for first-run scope, authority and acceptance.

F12 and F14 use the review range, revalidation revision and working-catalog
evidence in the [C06 code review (xhigh) source note](../docs/history/status-C10.md#status-record).

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| C07.1 — Prepare the first publisher inventory | `[+]` | Captured 334 unique reviewed providers and 196 machine-readable references, including 195 official references and the historical third-party Kubernetes fixture. The 181 earlier references are sorted in the [offline inventory](../../dist/c06-goal.zgalc0nf/c07-readiness/older-priority.json). Verified the launcher, mount and selected scope. Recorded the one authorized live pass with 302 provider summaries and 160 reference outcomes. C07.5 retains the remaining inventory work. Astra xhigh passed the task review. |
| C07.2 — Qualify Dropbox, HubSpot and Kubernetes | `[ ]` | Compare Dropbox official OpenAPI with the auth and operations in the overlay, then prefer it when qualified. Preserve old Stone references as provenance; UWS already has no Stone wire type. **F14 (Nit, local P1, confirmed).** The [HubSpot inventory](../../catalog_update_sources.go) enumerates services, and the [updater](../../catalog_update.go) stops after listing them. Select the reviewed families, versions and release stages, qualify stable documents through pinned native operation and compatibility checks, and label previews and other services as candidates. Keep any service with insufficient evidence out of the qualified set. Use release-tagged upstream Kubernetes OpenAPI, keep the historical fixture labeled and preserve cluster-specific exports for actual API groups and CRDs. |
| C07.3 — Qualify Stripe, Google Sheets and Slack | `[ ]` | Separate Stripe dated stable and preview releases and check the reviewed schema and security. Compare Google Sheets Discovery revision and digest even at the same v4 URL; reconcile the advisory overlay. Retain Slack’s archived Swagger as dated evidence and check current methods and scopes against official docs before proposing gap overlays. |
| C07.4 — Curate the remaining sources and overlay decisions | `[ ]` | Seek official links for existing docs-only and overlay-primary rows. Prefer a usable provider source; otherwise review a bounded endpoint overlay and reproducible builder or record the concrete evidence gap. **F12 (Nit, local P2, confirmed).** Of 161 working-catalog providers without an eligible updater reference, 139 already have a checked-in overlay. Use the [advisory registry](../../catalog/advisory.go) in the [provider summary](../../catalog_update.go) to distinguish advisory coverage, a documented gap and unavailable evidence. Report upstream freshness only when observed. Inspect reusable schemes and effective root/operation auth. Keep each provider's supported subset and license evidence visible. |
| C07.5 — Complete and reconcile the first inventory | `[ ]` | Run the authorized publisher inventory and account for each of the 334 provider IDs exactly once in the provider summary, with service/reference detail below it. Record actual versions, digests, outcomes, deferrals and review queue. Demonstrate reviewed operations for the six named providers with pinned fixtures; run catalog generation, quality and whole-unit gates. Use C10.4's accepted FND-01 failure classification and retry evidence to reconcile temporary publisher failures. |

## First publisher run

The fixed built-in catalog SHA-256 is
`fcfcd0c6d82cbde994c72e1e3d4f7df445b178adb8f4028e8b4d269ffb206cdd`.
The [host config rehearsal](../../dist/c06-goal.zgalc0nf/c07-readiness/maintenance-host.json)
uses `/run/media/peter/94A8DB9CA8DB7AE6/apitools-cache` as the
artifact directory and `/home/peter/.local/state/apitools/catalog-maintenance.sqlite`
for SQLite. The portable settings are `cache_dir` and `state_db`.
`/dev/sda3` is mounted at the expected UUID and its device number matches
the mount. The SQLite parent is on ext4. The first attempt initialized both
configured paths.

The first run was one sequential `catalog update --all` for this catalog
revision, supervised by the coordinator with a 60-minute process limit.
The authorized scope covered public publisher spec URLs and official GitHub
revision metadata through guarded HTTP(S) without credentials. The launcher
used the default 30-second request and 20 MiB parser limits and deferred
unknown sizes and files over 1 GB before their body request. Its scope
included checkpoints at the configured state path and retained artifacts at
the configured cache path. Its stop conditions were cancellation, timeout or
unavailable storage. [Class R test policy](test-policy.md#class-r-runs).

### Authorized first attempt

The operator explicitly approved one supervised C07.1 live inventory for
the 334 reviewed providers in the goal conversation. The coordinator reserved
attempt 1 of 1 at `2026-11-25T15:08Z`. The launcher uses Linux `timeout` with
the 60-minute limit and TERM on expiration, plus bounded shutdown. The
accepted source is `49862534d0bc85e1b19f0f39068dc8fecf8d4017`.
The offline Go 1.26.6 build and the exact launcher's configured dry run
passed with the expected catalog SHA-256, 334 distinct provider IDs, and
no run ID. Mount inspection confirmed `/dev/sda3` at the configured UUID
and ext4 for SQLite. [Launcher and input hashes](../../dist/c07-first-inventory.0yJaZ4/inputs.sha256).

Attempt 1 of 1 ran from `2026-11-25T15:09:12Z` to `2026-11-25T15:14:26Z`.
The CLI exited 0 and left run `2OXWJLDNCITU5LOUTAZAVOZ25V` incomplete.
The [first-pass JSON report](../../dist/c07-first-inventory.0yJaZ4/report.json)
has SHA-256 `260f33357439355da9bd20173b46d0ec8e7fdd833e02798f5bc4476c5062bf8d`.
Of 334 selected providers, 302 had provider summaries, including 161 with
no eligible updater reference. There are 160 service/reference outcomes
for 195 eligible references.
Of these, 48 require review, 35 report invalid specs, 30 report publisher
check failures, 26 have unknown sizes, 10 reproduced reviewed pins, and
11 have other recorded results or deferrals. One exceeded the parser
limit. There were no over-1-GB outcomes in this pass.

The open run lacks 32 provider summaries and 35 reference checkpoints. It
also has three explicit retry checkpoints from two Jina timeouts and
GitHub's shared HTTP 403 rate limit on Mercado Pago revision metadata. The
GitHub retry deadline is `2026-11-25T16:09:13Z`. The shared-origin deferral
skipped further GitHub references in this pass. C07.5 owns inventory
reconciliation. Any live continuation needs a new named scope for this
run ID because the one-attempt approval is consumed.

Read-only inspection verified SQLite integrity, the foreign keys, the
334-provider scope, and exact agreement of all 462 outcome rows with
the saved JSON. The cache has 20 verified digest-addressed files totaling
20,697,477 bytes. SQLite has 12 preferred versions. The offline Go suite,
vet, focused updater and storage tests, generator check, and catalog
quality check passed. [Offline verification](../../dist/c07-first-inventory.0yJaZ4/).

The independent Astra xhigh task review confirmed the scope, counts,
checkpoint evidence, cache digests and open run state. C07.2 still needs
publisher source evidence for the official Dropbox OpenAPI, selected
HubSpot releases and a release-tagged upstream Kubernetes document. This
pass returned a Dropbox source review and a HubSpot unknown-size deferral.
The historical Kubernetes fixture is third-party evidence and was outside
the eligible references.

## Review

Review gate iterations: 0

Whole review pending, default limit 10. The authorized publisher run provides
the first-run evidence for C09. Use the [test policy](test-policy.md) and
[milestone review procedure](milestone.md#milestone-review-procedure).
