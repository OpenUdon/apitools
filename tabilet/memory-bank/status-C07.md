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
| C07.2 — Qualify Dropbox, HubSpot and Kubernetes | `[+]` | Compare Dropbox official OpenAPI with the auth and operations in the overlay, then prefer it when qualified. Preserve old Stone references as provenance; UWS already has no Stone wire type. **F14 (Nit, local P1, confirmed).** The [HubSpot inventory](../../catalog_update_sources.go) enumerates services, and the [updater](../../catalog_update.go) stops after listing them. Select the reviewed families, versions and release stages, qualify stable documents through pinned native operation and compatibility checks, and label previews and other services as candidates. Keep any service with insufficient evidence out of the qualified set. Use release-tagged upstream Kubernetes OpenAPI, keep the historical fixture labeled and preserve cluster-specific exports for actual API groups and CRDs. Pinned stable Contacts/Deals list operations and Kubernetes Core Pod list passed with scoped auth evidence. Dropbox retains the reverified advisory overlay and official Stone provenance. HubSpot version migrations remain unexamined. Astra xhigh passed re-review. |
| C07.3 — Qualify Stripe, Google Sheets and Slack | `[+]` | Separate Stripe dated stable and preview releases and check the reviewed schema and security. Compare Google Sheets Discovery revision and digest even at the same v4 URL; reconcile the advisory overlay. Retain Slack’s archived Swagger as dated evidence and check current methods and scopes against official docs before proposing gap overlays. Qualified Stripe GA with reviewed auth, compared Sheets revisions and exact pins, and added the current Slack advisory. Astra xhigh passed re-review. |
| C07.4 — Curate the remaining sources and overlay decisions | `[+]` | Seek official links for existing docs-only and overlay-primary rows. Prefer a usable provider source; otherwise review a bounded endpoint overlay and reproducible builder or record the concrete evidence gap. **F12 (Nit, local P2, confirmed).** Of 161 working-catalog providers without an eligible updater reference, 139 already have a checked-in overlay. Use the [advisory registry](../../catalog/advisory.go) in the [provider summary](../../catalog_update.go) to distinguish advisory coverage, a documented gap and unavailable evidence. Report upstream freshness only when observed. Inspect reusable schemes and effective root/operation auth. Keep each provider's supported subset and license evidence visible. Screened all 160. Qualified four official sources, retained 137 primary advisories and 19 explicit gaps, and recorded the three new local coverage outcomes. Astra xhigh passed re-review. |
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

## Qualification evidence

C07.2 verified the Dropbox advisory builder and all eight Stone route scope
mappings. The publisher's official model remains Stone. The new HubSpot
Contacts and Deals references are publisher-owned `2026-09` rollout `424`
artifacts. Latest/stable index entries with a `NONE` banner are stable
candidates. Other releases, including previews, need their own qualification.
The selected HubSpot and Kubernetes references have operation-scoped auth
overlays. Kubernetes Core v1 OpenAPI v3 passed against the pinned release
tag and exact commit. The historical third-party fixture remains labeled.
[Source and operation evidence](../../docs/catalog-qualification-c07.md#dropbox-hubspot-and-kubernetes),
[offline qualification](../../dist/c07-qualification.P2FWoi/c07-2-compatibility-final.log).

The first C07.2 review found a missing HubSpot compatibility check and
scheme-name collisions in the provider inspection. Pinned Contacts and
Deals checks now record `compatibility-unexamined` for both transitions.
The auth scheme identifiers are service-specific, and inspection verifies
that they produce no duplicate-scheme finding. Verification ran `go build
./...`, `go vet ./...`, `go test -count=1 . ./catalog ./cmd/apitools`,
`go run ./cmd/cataloggen -check`, `go run ./cmd/apitools catalog check`,
and `git diff --check` on `5ebadc0039d0ca222fda35964c7e3fe4c9509528`
with the C04/C05 and C07 pinned source directories. All passed.
[Affected-unit checks](../../dist/c07-qualification.P2FWoi/c07-2-final-gates.log).
The [reviewed diff](../../dist/c07-qualification.P2FWoi/c07-2-reviewed.diff)
has SHA-256 `1f29b1f7618e55c990d0fcda9fd73f4b214d22bb713cf7bd470f9f1d93defd27`.

C07.3 qualified the Stripe GA PaymentIntent retrieval schema and Basic or
Bearer alternatives from the Endive pin. Preview bytes remain a candidate.
The same Sheets v4 URL returned a newer Discovery revision under a bounded
read. The prior catalog revision had no retained digest. The SDK snapshots
and canonical pin reproduce the selected native method, and the updater
has a bounded exact-pin proof for the two reviewed references. Slack
retains the archived Swagger provenance and adds a reproducible current
message and public-channel history overlay.
[Source and operation evidence](../../docs/catalog-qualification-c07.md#stripe-google-sheets-and-slack).

Verification on `e469981476c62ec52d3ad9d49a1a1a95d5a70ecf` passed. It ran
`go build ./...`, `go vet ./...`, `go test -count=1 . ./catalog ./cmd/apitools`
with the C04/C05 and C07 source pins, the generator and catalog quality
checks, and `git diff --check`. The first review required a fixed per-ref
digest for the reduced Sheets proof and separate Git and Discovery revision
identities. A legacy catalog repin now fails the selected-method proof, and
the SDK second pass reuses its verified object without fetching the body.
[Affected-unit checks](../../dist/c07-qualification.P2FWoi/c07-3-repeat-gates.log),
[focused checks](../../dist/c07-qualification.P2FWoi/c07-3-repeat-final.log).
The [re-reviewed diff](../../dist/c07-qualification.P2FWoi/c07-3-repeat-reviewed.diff)
has SHA-256 `c3b53c3f9b028df19685a0a1b055fb96666d9ed5ea95df90c8e4046aafe1ac96`.

C07.4 screened the 160 scoped publisher docs pages and obtained four direct
candidate files plus the separately identified Amplitude and Linear
repository artifacts. The selected Amplitude Developer, Hunter v2 and
Storyblok CDN operations passed guarded OpenAPI validation and native
shape proof. Linear passed native GraphQL qualification for `Query.viewer`.
Both Postmark Swagger candidates failed guarded parsing. The Postmark
advisory remains reproducible, with separate account and server auth
requirements on its reviewed operations. Storyblok records CDN and
Management fallback auth separately. The 137 remaining primary overlay
providers have a tracked artifact and builder, and 19 have a documented
source gap. New provider checkpoints distinguish advisory coverage,
documented gap and unexamined coverage.
[Curation evidence](../../docs/catalog-qualification-c07.md#remaining-publisher-sources-and-overlays).

The C07.4 first review found the missing document license declarations,
fallback auth scope, and mismatched source reference precedence in the
overlay registry accessor. The revised catalog records the document
licenses separately from repository terms, uses exact fallback operation
auth matches, and retains the same reference precedence as the SQLite
registry. All 146 registration literals match after field relocation.
Go 1.26.6 build, vet, affected whole Go packages with both approved pin
directories, catalog generator, catalog quality, artifact launcher build,
three byte-identical overlay rebuilds and diff checks passed on
`1568dc639aa9ecb41646c559cf34d231cddbaf12`.
[Verification](../../dist/c07-qualification.P2FWoi/c07-4-final-gates.log).
Astra xhigh passed the [re-reviewed diff](../../dist/c07-qualification.P2FWoi/c07-4-final-reviewed.diff)
with SHA-256 `d4795cb48b0009f41dc2ad38c928f29703462cbd4b218aa44f17112430f3e25a`.

## Review

The operator approved public publisher documentation reads and bounded,
unauthenticated spec downloads for C07.2–C07.4 in this goal conversation.
The scope is the six named providers and the 160 existing overlay/docs-primary
entries, limited to official sites and publisher-owned repositories, with
source copies in ignored storage and provenance recorded. The completed
first inventory stays closed. C07 execution uses the approved Sol high
implementation, Astra xhigh task reviews and Opus 5.5 high whole review.

Review gate iterations: 0

Whole review pending, default limit 10. The authorized publisher run provides
the first-run evidence for C09. Use the [test policy](test-policy.md) and
[milestone review procedure](milestone.md#milestone-review-procedure).
