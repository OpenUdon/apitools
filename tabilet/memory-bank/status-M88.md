# Status M88 — Catalog management CLI and guarded remote discovery

**State.** Pending; no row started. Whole-milestone review 0/10, not started.
**Context.** Approved catalog track in the current horizon alongside Kinet STG-12, without changing the original browser scope or provisional stages.
**Specification.** [M88](milestone.md#m88--catalog-management-cli-and-guarded-remote-discovery).

**Planning authority.** The user approved the complete proposal with “Implement the plan” on 2026-10-11; these are planning records only.

**Review source.** Catalog roadmap review (2026-10-11), pasted review and `plan_cross_package_catalog_root_and_search.md` / `plan_llm_touchpoints_in_apitools.md`. Review baseline: not supplied. F01-F13 use source order; this is ordinary intake, not a closing-review iteration.

**Source fingerprints.** Cross-package plan SHA256 `6fbecdb55ceed1d089b63f14b8d064a5bc58e0d0986a9ac5d0ba26f865f464fe`; LLM plan SHA256 `811a33a621b6e4df728e8a97fb6c3fe082e9e0e8d9915a7d06484e285c50b71f`. Source plans remain read-only.

**Revalidation.** `2b74e96f0dd38aeb228d585f750d79194bdc7d06`; clean at assessment and before planning writes. Relevant uncommitted implementation evidence: none. Cross-package baselines: APItools `2b74e96f0dd38aeb228d585f750d79194bdc7d06`, Kinet `31721164eb6a4cef382b17cd3956a3a125a93f72`, Ramen `834cb8e77004c07c6014e4428c2fdb1fe2cb8345`.
**Lineage.** Retired M80/M81 and existing index/refresh/import/search commands; extend compatible contracts instead of duplicating commands.

## Review ownership

| Finding | Source priority | Local severity | Disposition | Repository evidence | Primary owner |
|---|---|---|---|---|---|
| F06 | P2 | P2 | Confirmed plan defect | client.go Search uses ProviderURL for RFC 9727; catalog_discovery_remote.go currently supplies no host and selects APIs.guru. SourceAuto alone is insufficient. | M88.R2 |

## Dependencies and ownership

**Depends on.** APItools:M87 accepted.
**Downstream impacts.** APItools:M89; Ramen:M81.
**Later implementation write set.** cmd/apitools commands and catalog/discovery helpers, tests and ordinary docs; this status. Shared planning/closure is coordinator-owned; no child edits sibling repositories or another status.

**Execution policy.** Planning authorizes no implementation, commit, publication, deployment or live-model operation. A separate execution request governs task commits. One coordinator owns shared records/integration; this catalog track is serial and **Parallel-safe: no**. No overlapping catalog or same-package Stage 12 worker/module writers. Use immutable inputs; sibling adoption requires exact accepted and independently published revisions, not workspace replacements. Publication requires separately named authority; old Stage 12 grants do not transfer. Audit remains disabled/unavailable as recorded.

## Task status

Each row is a later implementation commit unit under the governing request. Markers and bounded review/retirement rules come from [milestone.md](milestone.md). Planning completes no row.

| Item | State | Notes |
|---|---|---|
| M88.R1 — Catalog lifecycle commands | `[+]` | Delivered catalog lifecycle CLI commands (init, add, import, search, snapshot) and environment-root resolution for index/refresh in cmd/apitools; retained explicit-root no-repair/no-refresh guarantees and verified 100% offline tests. |
| M88.R2 — Provider-scoped remote discovery | `[+]` | Broadened DiscoverCatalogOperations remote lookup to SourceAuto, plumbed installation-side ProviderURL through SearchOptions.ProviderURL for RFC 9727 fallback (resolving F06), enforced URL safety and dual opt-in, and verified with offline fixtures. |
| M88.R3 — CLI and consumer qualification | `[ ]` | Run offline owner/CLI/network-fixture/consumer gates and record exact accepted source plus separately authorized publication handoff. |

## Acceptance and verification

**Acceptance.** Compatible explicit CLI operations can create, populate, index, search, refresh and snapshot catalogs. RFC 9727 receives an approved provider URL when earlier auto tiers miss. Offline paths never fetch; remote leads need local import and requalification before selection.

**Verification.** GOWORK=off go test ./...; GOWORK=off go vet ./...; CLI lifecycle/legacy fixtures; injected HTTP fixtures for provider URL, SourceAuto, refusal/deadline/partial coverage; immutable consumer checks; git diff --check.

## Review gate and scoped commit tracking

- Whole-milestone review: 0/10, not started; begin after implementation and required verification.
- Implementation commits:
  - M88.R1: feat(M88.R1): implement catalog lifecycle CLI commands and root resolution
  - M88.R2: feat(M88.R2): broaden remote catalog discovery to SourceAuto with provider URL
- No prior review counter or frozen status is reopened.
