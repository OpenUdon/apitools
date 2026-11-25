# Retired milestone C07 - First inventory and provider curation

**Milestone.** C07
**Outcome.** completed
**Retired.** 2026-11-25
**Source status.** tabilet/memory-bank/status-C07.md
**Source specification.** tabilet/memory-bank/milestone.md#c07--first-inventory-and-provider-curation
**Evidence.** 5aa3bb78122916c57e527d1fe3b48f8bebb5a048
**Worktree.** includes uncommitted changes
**Review.** passed
**Review iterations.** 2
**Verification.** Fresh offline Go 1.26.6 build, vet, full suite with C04/C05 and C07 pinned artifacts, focused and storage race checks, generator and catalog quality, CLI help, identical 334-provider scope, and byte-identical read-only inspection of the closed 529-outcome SQLite report passed. Offline reconciliation verified the 25 retained artifacts. OpenUdon passed against the corrected candidate. Udon's ordinary pinned suite passed. Its optional composed test retains the recorded baseline offline module limit. Closure reuses the required code gates because their inputs are unchanged. [Verification receipts](../../../dist/c07-qualification.P2FWoi/).
**Consolidated into.** [README](../../../README.md#reviewed-provider-catalog), [qualification and reconciled inventory](../../../docs/catalog-qualification-c07.md), [catalog maintenance architecture](../../memory-bank/architecture.md#catalog-maintenance-storage-c06), and [source review lessons](../../memory-bank/lessons.md#bind-version-evidence-to-an-exact-publisher-and-native-identity).

## Milestone specification

````markdown
### C07 — First inventory and provider curation

**Dependency and owner.** Uses the accepted [C10 updater](../docs/history/status-C10.md)
against the retired [C06 baseline](../docs/history/status-C06.md).
APItools owns the provider recipes, curated sources, security and endpoint
overlays, and first inventory.
[C07 ledger](status-C07.md).

Inventory the 334 reviewed providers, starting with the 181 older
machine-readable references. Seek an official source for the 139 providers
whose primary evidence is a docs-derived overlay and the 21 docs-only
providers. For each unresolved service, use official docs to build or refresh
a bounded, reproducible overlay and builder, or explain the specific missing
evidence. Review effective security schemes and requirements for protected
operations. Stage authored changes for review. Acceptance requires a
provider-level outcome for every entry and operation-level evidence for the
six named providers in the ledger. Qualify the reviewed HubSpot services and
their operations against pinned evidence. Show which other index entries are
candidates or previews. Where the updater has no eligible machine-readable
reference, report existing advisory coverage and any remaining evidence gap.
Distinguish redirects, failures, deferrals and partial coverage from a verified update.

````

## Status record

````markdown
# C07 — First inventory and provider curation

**Owner.** APItools. **Dependency.** Accepted [C10 updater repairs](../docs/history/status-C10.md) following [C06](../docs/history/status-C06.md). **Parallel-safe: no.**
See the [milestone specification](milestone.md#c07--first-inventory-and-provider-curation)
for first-run scope, authority and acceptance.

F12 and F14 use the review range, revalidation revision and working-catalog
evidence in the [C06 code review (xhigh) source note](../docs/history/status-C10.md#status-record).

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| C07.1 — Prepare the first publisher inventory | `[+]` | Captured 334 unique reviewed providers and 196 machine-readable references, including 195 official references and the historical third-party Kubernetes fixture. Prepared the 181 earlier references first in the [offline inventory](../../dist/c06-goal.zgalc0nf/c07-readiness/older-priority.json). Verified the launcher, mount and selected scope. Recorded the one authorized live pass with 302 provider summaries and 160 reference outcomes. C07.5 reconciled the deterministic provider-order inventory and the curated catalog. Astra xhigh passed the task review. |
| C07.2 — Qualify Dropbox, HubSpot and Kubernetes | `[+]` | Compare Dropbox official OpenAPI with the auth and operations in the overlay, then prefer it when qualified. Preserve old Stone references as provenance; UWS already has no Stone wire type. **F14 (Nit, local P1, confirmed).** The [HubSpot inventory](../../catalog_update_sources.go) enumerates services, and the [updater](../../catalog_update.go) stops after listing them. Select the reviewed families, versions and release stages, qualify stable documents through pinned native operation and compatibility checks, and label previews and other services as candidates. Keep any service with insufficient evidence out of the qualified set. Use release-tagged upstream Kubernetes OpenAPI, keep the historical fixture labeled and preserve cluster-specific exports for actual API groups and CRDs. Pinned Contacts/Deals list operations in indexed stable family versions and Kubernetes Core Pod list passed with scoped auth evidence. Dropbox retains the reverified advisory overlay and official Stone provenance. HubSpot version migrations remain unexamined. Astra xhigh passed re-review. |
| C07.3 — Qualify Stripe, Google Sheets and Slack | `[+]` | Separate Stripe dated stable and preview releases and check the reviewed schema and security. Compare Google Sheets Discovery revision and digest even at the same v4 URL; reconcile the advisory overlay. Retain Slack’s archived Swagger as dated evidence and check current methods and scopes against official docs before proposing gap overlays. Qualified Stripe GA with reviewed auth, compared Sheets revisions and exact pins, and added the current Slack advisory. Astra xhigh passed re-review. |
| C07.4 — Curate the remaining sources and overlay decisions | `[+]` | Seek official links for existing docs-only and overlay-primary rows. Prefer a usable provider source; otherwise review a bounded endpoint overlay and reproducible builder or record the concrete evidence gap. **F12 (Nit, local P2, confirmed).** Of 161 working-catalog providers without an eligible updater reference, 139 already have a checked-in overlay. Use the [advisory registry](../../catalog/advisory.go) in the [provider summary](../../catalog_update.go) to distinguish advisory coverage, a documented gap and unavailable evidence. Report upstream freshness only when observed. Inspect reusable schemes and effective root/operation auth. Keep each provider's supported subset and license evidence visible. Attempted all 160 scoped pages, read 130, and deferred 30. Qualified four official sources, retained 137 primary advisories and 19 explicit gaps, and recorded the three new local coverage outcomes. Astra xhigh passed re-review. |
| C07.5 — Complete and reconcile the first inventory | `[+]` | Run the authorized publisher inventory and account for each of the 334 provider IDs exactly once in the provider summary, with service/reference detail below it. Record actual versions, digests, outcomes, deferrals and review queue. Demonstrate reviewed operations for the six named providers with pinned fixtures; run catalog generation, quality and whole-unit gates. Use C10.4's accepted FND-01 failure classification and retry evidence to reconcile temporary publisher failures. The live run is complete with 334 provider summaries and 195 reference outcomes. Both Jina retries exhausted. C07.2–C07.4 passed. The closed SQLite report reproduced byte for byte, all 25 cache objects matched, and the six source proofs and whole-unit gates passed. Astra xhigh passed the task review. |

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

### Authorized resume

The operator explicitly approved run `2OXWJLDNCITU5LOUTAZAVOZ25V` in the
goal conversation. The coordinator reserved one supervised `catalog update
--resume` pass at `2026-11-25T17:03Z`. The original binary, configuration,
catalog fingerprint, approval scope, cache, SQLite location and 60-minute
supervision limit carry forward. Read-only inspection found the same
unfinished run with 302 provider summaries and 160 reference outcomes.
The three saved retry deadlines had elapsed. This is inventory attempt 2
and the only invocation under the new approval.

The resume ran from `2026-11-25T17:03:58Z` to `2026-11-25T17:05:44Z` and
exited 0. The [resume JSON report](../../dist/c07-first-inventory.0yJaZ4/resume-1/report.json)
has SHA-256 `f7bd4ea213d466539a0ca70f56d5af54c13652f2a0f488b03232b92cf0d06287`.
It saved 31 additional provider summaries and 35 additional reference
checkpoints. The GitHub retry for Mercado Pago reproduced its reviewed
pin. All 195 eligible references now have outcomes, including 16 reviewed
pins reproduced, 56 requiring review, 43 invalid specs and 35 publisher
check failures. Ten documents exceeded the 20 MiB parser limit and 26
had unknown sizes.

The run had 333 provider summaries after this pass. Only Jina AI was
missing its summary. Its `jina-reader-openapi` and `jina-search-openapi`
checks each timed out on attempt 2. The later of their two retry times
was `2026-11-25T17:05:18.388394271Z`.

Read-only verification confirmed SQLite integrity and that its 528
outcome rows match the resume report. The original 459 non-pending
outcomes are unchanged. Checked all 25 stored file digests and byte
counts, totaling 30,760,406 bytes. The database now has 17 preferred
versions. The one supervised resume approval has been consumed.

### Authorized final resume

The operator approved one final supervised resume for the two pending
Jina references in run `2OXWJLDNCITU5LOUTAZAVOZ25V`. The coordinator
reserved inventory attempt 3 at `2026-11-25T17:11Z`. Read-only inspection
verified the original binary, config, scope, mounts and SQLite run. Only
the Jina Reader and Search checkpoints are eligible, both on their third
attempt. The same 60-minute process limit and storage paths apply.

Attempt 3 ran from `2026-11-25T17:11:38Z` to `2026-11-25T17:12:38Z` and
exited 0. The [completed JSON report](../../dist/c07-first-inventory.0yJaZ4/resume-2/report.json)
has SHA-256 `c712ed2d84c5fdd7e201e090ddd5cc691e4f97e866d8555d6293d55757c2e687`.
Jina Reader and Search both timed out and reached `publisher-retry-exhausted`
on attempt 3. This terminal disposition allowed the original run to record
the Jina provider summary and finish. Its third reference retains the earlier
unknown-size observation. The exhaustion establishes a fetch failure, with
freshness and document validity for those sources unresolved.

Read-only checks verified one summary for each of the 334 selected provider
IDs, exactly the 195 eligible reference outcomes, and zero pending retries.
All 529 outcomes and the completion timestamp matched SQLite. The resume
preserved 526 earlier outcomes and changed only the two pending references
plus the new Jina summary. All 25 cached artifacts passed size and digest
checks. SQLite retained 17 preferred versions. The final live approval is
consumed. C07.5 subsequently reconciled this snapshot with the C07.2–C07.4
results and recorded the acceptance gates.

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

## Inventory acceptance proof

C07.5 reconciled the completed inventory and the curated catalog in a
separate, read-only report. It validated the exact original reference and
provider set against the first inventory's catalog fingerprint, including
334 provider summaries and 195 eligible reference outcomes, with zero
pending retries. A stateless dry run on the curated catalog selects the
same 334-provider set. Its new catalog SHA-256 is
`0118c98b3f07cce24e3e23f898e460211704f3743c267d1f260d6e8463c991a8`.
The current catalog has 203 eligible references. The read-only state
inspection reproduced the archived report byte for byte. The persisted
database SHA-256 was unchanged. The 25 retained files matched all
recorded byte counts and digests, totaling 30,760,406 bytes.
[Per-provider reconciliation](../../dist/c07-qualification.P2FWoi/c07-5-provider-reconciliation.json),
SHA-256 `342d094654d8e3355ee64ab2458a8196ffb883f866c91a40856e9d7e7c020556`.
[Digest and scope check](../../dist/c07-qualification.P2FWoi/c07-5-reconciliation.log).

The [qualification report](../../docs/catalog-qualification-c07.md#first-inventory-reconciliation)
summarizes the 195 reference dispositions, pin/review outcomes and deferrals.
The per-provider receipt retains each reference checkpoint.
Jina Reader and Search exhausted their bounded retries. Their publisher
versions remain unresolved. Postmark's candidate Swagger files require
further qualification. Existing advisory and documented gaps stay distinct
from demonstrated publisher freshness.

Go 1.26.6 build, vet, fresh whole Go suite with both pinned source
directories, catalog storage and targeted updater races, generator,
catalog quality (zero findings), search/import/update help, the
334-provider dry run, historical SQLite inspection and diff checks all
passed at `a85c93bcac83ff648cecdea3f507fc0d28459d36` with C07.5
documentation and status edits in the worktree.
[Acceptance checks](../../dist/c07-qualification.P2FWoi/c07-5-gates.log).

The independent Astra xhigh task review passed the C07.5 reconciliation on
`a85c93bcac83ff648cecdea3f507fc0d28459d36` and the
[reviewed diff](../../dist/c07-qualification.P2FWoi/c07-5-reviewed.diff),
SHA-256 `7d55bd8817343bcefbe94b7d9b5134aced9e725cce57a094ff9f39e115b3aee9`.
The [verbose operation gate](../../dist/c07-qualification.P2FWoi/c07-5-pinned-operations.log)
shows the active pin directories and the selected tests. Supplemental
OpenUdon snapshot tests passed against this APItools source. The Udon
ordinary pinned CLI test passed. Its supplemental workspace test with
local APItools is limited by the missing offline
`cloud.google.com/go/compute@v1.49.1` module. The same workspace test has
that limit at the pre-C07 APItools baseline, and C07 leaves the Go module
inputs unchanged.
[Consumer evidence](../../dist/c07-qualification.P2FWoi/c07-5-consumers/).

## Review

The operator approved public publisher documentation reads and bounded,
unauthenticated spec downloads for C07.2–C07.4 in this goal conversation.
The scope is the six named providers and the 160 existing overlay/docs-primary
entries, limited to official sites and publisher-owned repositories, with
source copies in ignored storage and provenance recorded. The completed
first inventory stays closed. C07 execution uses the approved Sol high
implementation, Astra xhigh task reviews and Opus 5.5 high whole review.

Review gate iterations: 2

Whole review round 1 of 10 started on candidate
`f64f1a64c5652272756316cccffeafeebbe094e0`, pre-C07 base
`49862534d0bc85e1b19f0f39068dc8fecf8d4017`. The round 1 dispatch covered that
candidate. The [canonical full diff](../../dist/c07-qualification.P2FWoi/c07-whole-round1-full.diff)
has SHA-256 `5c4b4bdea3b48d0ffe9ceef1e18f857298ef0c5d55cab1935a69e95b402e22df`.
The [Opus review brief](../../dist/c07-qualification.P2FWoi/c07-whole-round1-review.txt)
has SHA-256 `abff8d3a6af9daae6c6b97b43cb421b176122319c663b9fc97896125fec2cc23`.

The whole offline gate passed afresh on that revision with Go 1.26.6:
build, vet, full Go test suite with pinned C04/C05 and C07 sources, catalog
generator, catalog quality, search/import help and diff check.
[Whole verification](../../dist/c07-qualification.P2FWoi/c07-whole-gates.log).
Effective required gate recovery limit 3 from the approved launch. The
initial whole build launch on that revision failed setup because PATH selected
Go 1.25.6 while GOROOT selected the Go 1.26.6 compiler. The explicit Go 1.26.6
binary recovered the gate on retry. One setup charge, two remaining for
that key. [Initial setup receipt](../../dist/c07-qualification.P2FWoi/c07-whole-toolchain-selection.log).
OpenUdon's isolated workspace full suite and Udon's ordinary pinned full
suite passed. The optional Udon workspace against the local APItools
candidate retained the pre-C07 offline module limitation above.

The approved whole review limit is 10. The authorized publisher run provides
the first-run evidence for C09. Use the [test policy](test-policy.md) and
[milestone review procedure](milestone.md#milestone-review-procedure).

Round 1 BLOCK from Claude Opus 5.5 high. F1 (P1) confirmed that the active
evidence gave HubSpot release categories without the index entry listing
or stage/banner for the selected families. The retained publisher index
contains the missing evidence. The round also reported scope, pin,
security, consumer compatibility and evidence findings F2–F14.
[Full review](../../dist/c07-qualification.P2FWoi/c07-whole-round1-report.md),
SHA-256 `85ecce2138d25557349f48504f5b8bfcbe3bcae236ac24f953a0386883643ecd`.
[Model and dispatch receipt](../../dist/c07-qualification.P2FWoi/c07-whole-round1-dispatch.json).

The Round 1 correction records the 377 HubSpot entries with per-entry stage,
banner and candidate/preview disposition. The pinned index shows Contacts
release `165147` and Deals release `165150`, each `LATEST` with the `NONE`
banner. The selected native proofs use independent repository rollout `424`
at the same family and version. Compatibility stays unexamined. The index
listing is tracked. [HubSpot evidence](../../docs/catalog-qualification-c07.md#dropbox-hubspot-and-kubernetes).

F3 and F4 were fixed with digest-bound Dropbox Stone and Slack Swagger
proofs and the matching earlier Slack Git blob. F5 was fixed by disclosing
the three spec-selection security audit items separately from the historical
56 reference-level review outcomes. F7 and F8 were addressed by documenting
and testing provider defaults, checking HubSpot's per-reference auth and
encoding Linear's alternatives as separate single-requirement sets. F9–F12
were fixed by retiring stale gap metadata, clarifying the README, checking
the raw auth facts and requiring the actual HubSpot compatibility
disposition. F13 was addressed with a tracked HubSpot listing and explicit
local per-reference evidence and distinct result counts. F14 was clarified
with the earlier-priority preparation and the actual updater order. The
Slack artifact registry uses the primary message docs URL and the artifact
retains both method sources.

F2 carries as FND-04 for the 30 deferred publisher pages. The current
scope records existing overlay or gap evidence and the exact observation
for each. F6 carries as FND-05 on pre-existing Google Slides auth.
Neither carried finding prevents the documented C07 operations or the
C08 search metadata task.

The correction produced catalog SHA-256
`d68d18e9b301f3f5a3fde091f8011f37564ec65a4c4c33aac48aefe5be382592`.
The [updated per-provider reconciliation](../../dist/c07-qualification.P2FWoi/c07-5-reconciled-after-review.json)
has SHA-256 `80f42e3e2348e80157e36289cdd9c283c0a57f602c26e58cb140003ebcd4f65f`.
It reproduces the original 334 providers and 195 reference outcomes against
the corrected catalog and checks the same 25 cache objects.

Started the round 1 corrective gates on
`5aa3bb78122916c57e527d1fe3b48f8bebb5a048` with the explicit Go 1.26.6
binary. The gates cover the full offline suite, pinned C07 operations,
race checks on the updater and state store, catalog checks, CLI help,
dry-run scope and read-only historical SQLite replay.

The corrective gates passed, including all new pinned tests, race gates,
catalog generation and zero quality findings. The curated dry run selects
the same 334 provider IDs. The historical CLI output and the SQLite SHA-256
remain byte-identical. [Round 2 gates](../../dist/c07-qualification.P2FWoi/c07-whole-round2-gates.log).

Whole review round 2 of 10 started on
`5aa3bb78122916c57e527d1fe3b48f8bebb5a048`. It covers F1, all corrections
since the first reviewed candidate, and their direct consumers.
[Corrective diff](../../dist/c07-qualification.P2FWoi/c07-whole-round2.diff)
SHA-256 `a054bcd419b58902f874386806e3354c06d2c43711422b4e979c364e327c9158`.
[Independent review brief](../../dist/c07-qualification.P2FWoi/c07-whole-round2-review.txt)
SHA-256 `bf625b5c93abd8bd8a906d4885b4a8086dc9ce8519549d410fd023ce8b203bc1`.

Round 2 PASS from Claude Opus 5.5 high. It confirmed F1 against the full
tracked listing, exact publisher stages and separate repository artifact
identity. No blocking findings remain. [Full review](../../dist/c07-qualification.P2FWoi/c07-whole-round2-report.md),
SHA-256 `55489dd54b7426c016355948328f101136b1c87c65362bd7e3aeaf6832ece9f3`.
[Model and dispatch receipt](../../dist/c07-qualification.P2FWoi/c07-whole-round2-dispatch.json).

Round 2 N1 carries as FND-06 for C09.4. The HubSpot index puts Test At Debug
`2026-09-beta` at `LATEST` with banner `NONE`. The entry is an unqualified
candidate. The listing preserves the version label and its conflicting
publisher labels. The C07 candidate identification and validated
operations remain scoped to Contacts and Deals.

N3 and N4 are lower priority additions to the existing exact-pin
regressions. The raw listing proof runs with the retained publisher
fixture, and the reviewed operations retain the source and scoped catalog
checks. N5 joins FND-04 as clarification of the inherited gap reasons.
N6 is a lower priority portability suggestion. The reviewed raw receipts
and reconciliation remain local, with hashes in the records and the
HubSpot listing checked in. N7 is a lower priority wording suggestion.
The release document explains stage/banner versus the exact repository
identity, and the review chronology now records the candidate for each
dispatch. N8 follows the open findings ownership rules. FND-04 and FND-05
await later source-specific work, with Google Slides providing the
confirmed example and the broader root-scope question still open.

An evolution version bump is unnecessary. This milestone curates publisher
metadata and first-run evidence within the existing catalog maintenance
boundary. C08 consumes the accepted provider identities, source and
coverage evidence, and freshness limits. C09 consumes the original
supervised inventory evidence and its separate repeat-run requirements.

The round 2 N2 check passed on the corrected candidate. The full OpenUdon
suite ran with immutable OpenUdon and APItools sources. No flat security
reader was found in the tracked OpenUdon or Udon source scan. Udon's
ordinary pinned suite passed earlier. Its optional composed check retains
the same documented offline dependency limit.
[Post-review consumer check](../../dist/c07-qualification.P2FWoi/c07-round2-consumers/).
````
