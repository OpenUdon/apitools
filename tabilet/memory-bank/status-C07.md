# C07 — First inventory and provider curation

**Owner.** APItools. **Dependency.** Accepted [C06](../docs/history/status-C06.md). **Parallel-safe: no.**
See the [milestone specification](milestone.md#c07--first-inventory-and-provider-curation)
for first-run scope, authority and acceptance.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| C07.1 — Prepare the first publisher inventory | `[ ]` | Fix the baseline list and reviewed service scopes. Verify the configured cache directory, any host-specific mount requirement, limits, existing pins and first-run configuration. Rehearse locally; record the concrete operator authority and budget before the class R inventory. Prioritize the 181 older machine-readable references. |
| C07.2 — Qualify Dropbox, HubSpot and Kubernetes | `[ ]` | Compare Dropbox official OpenAPI with the auth and operations in the overlay, then prefer it when qualified. Preserve old Stone references as provenance; UWS already has no Stone wire type. Enumerate HubSpot by family/version/stage, link stable documents and distinguish previews and unreviewed services. Use release-tagged upstream Kubernetes OpenAPI, keep the historical fixture labeled and preserve cluster-specific exports for actual API groups and CRDs. |
| C07.3 — Qualify Stripe, Google Sheets and Slack | `[ ]` | Separate Stripe dated stable and preview releases and check the reviewed schema and security. Compare Google Sheets Discovery revision and digest even at the same v4 URL; reconcile the advisory overlay. Retain Slack’s archived Swagger as dated evidence and check current methods and scopes against official docs before proposing gap overlays. |
| C07.4 — Curate the remaining sources and overlay decisions | `[ ]` | Seek official links for existing docs-only and overlay-primary rows. Prefer a usable provider source; otherwise review a bounded endpoint overlay and reproducible builder or record the concrete evidence gap. Inspect reusable schemes and effective root/operation auth. Keep each provider's supported subset and license evidence visible. |
| C07.5 — Complete and reconcile the first inventory | `[ ]` | Run the authorized publisher inventory and account for each of the 334 provider IDs exactly once in the provider summary, with service/reference detail below it. Record actual versions, digests, outcomes, deferrals and review queue. Demonstrate reviewed operations for the six named providers with pinned fixtures; run catalog generation, quality and whole-unit gates. This row fixes FND-01. |

## Review

Review gate iterations: 0

Whole review pending, default limit 10. The authorized publisher run provides
the first-run evidence for C09. Use the [test policy](test-policy.md) and
[milestone review procedure](milestone.md#milestone-review-procedure).
