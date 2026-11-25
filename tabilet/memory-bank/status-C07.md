# C07 — First inventory and provider curation

**Owner.** APItools. **Dependency.** Accepted [C06](../docs/history/status-C06.md). **Parallel-safe: no.**
See the [milestone specification](milestone.md#c07--first-inventory-and-provider-curation)
for first-run scope, authority and acceptance.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| C07.1 — Prepare the first publisher inventory | `[!]` | Captured 334 unique reviewed providers and 196 machine-readable references, including 195 official references and the historical third-party Kubernetes fixture. The 181 earlier references are sorted in the [offline inventory](../../dist/c06-goal.zgalc0nf/c07-readiness/older-priority.json). The configured CLI dry run covered all 334 without writes. Verified the host paths and volume identity. Waiting for the named publisher-run authority described below. |
| C07.2 — Qualify Dropbox, HubSpot and Kubernetes | `[ ]` | Compare Dropbox official OpenAPI with the auth and operations in the overlay, then prefer it when qualified. Preserve old Stone references as provenance; UWS already has no Stone wire type. Enumerate HubSpot by family/version/stage, link stable documents and distinguish previews and unreviewed services. Use release-tagged upstream Kubernetes OpenAPI, keep the historical fixture labeled and preserve cluster-specific exports for actual API groups and CRDs. |
| C07.3 — Qualify Stripe, Google Sheets and Slack | `[ ]` | Separate Stripe dated stable and preview releases and check the reviewed schema and security. Compare Google Sheets Discovery revision and digest even at the same v4 URL; reconcile the advisory overlay. Retain Slack’s archived Swagger as dated evidence and check current methods and scopes against official docs before proposing gap overlays. |
| C07.4 — Curate the remaining sources and overlay decisions | `[ ]` | Seek official links for existing docs-only and overlay-primary rows. Prefer a usable provider source; otherwise review a bounded endpoint overlay and reproducible builder or record the concrete evidence gap. Inspect reusable schemes and effective root/operation auth. Keep each provider's supported subset and license evidence visible. |
| C07.5 — Complete and reconcile the first inventory | `[ ]` | Run the authorized publisher inventory and account for each of the 334 provider IDs exactly once in the provider summary, with service/reference detail below it. Record actual versions, digests, outcomes, deferrals and review queue. Demonstrate reviewed operations for the six named providers with pinned fixtures; run catalog generation, quality and whole-unit gates. This row fixes FND-01. |

## First publisher run

The fixed built-in catalog SHA-256 is
`fcfcd0c6d82cbde994c72e1e3d4f7df445b178adb8f4028e8b4d269ffb206cdd`.
The [host config rehearsal](../../dist/c06-goal.zgalc0nf/c07-readiness/maintenance-host.json)
uses `/run/media/peter/94A8DB9CA8DB7AE6/apitools-cache` as the
artifact directory and `/home/peter/.local/state/apitools/catalog-maintenance.sqlite`
for SQLite. The portable settings are `cache_dir` and `state_db`.
`/dev/sda3` is mounted at the expected UUID and its device number matches
the mount. The proposed SQLite parent is on ext4. Both configured paths are
pending first use.

The proposed first run is one sequential `catalog update --all` for this
catalog revision, supervised by the coordinator with the standard process
timeout of 60 minutes. It reads public publisher spec URLs and official
GitHub revision metadata within the resolved source scope. It uses guarded
HTTP(S), no credentials, and the default 30-second request and 20 MiB
spec parser limits. It defers unknown sizes and artifacts over 1 GB before
their body request. It writes checkpoints at the configured state path and
retained artifacts at the configured cache path. End the run on cancellation,
timeout or unavailable storage and retain completed checkpoints. This run
requires separate human authority under the [Class R test policy](test-policy.md#class-r-runs).

## Review

Review gate iterations: 0

Whole review pending, default limit 10. The authorized publisher run provides
the first-run evidence for C09. Use the [test policy](test-policy.md) and
[milestone review procedure](milestone.md#milestone-review-procedure).
