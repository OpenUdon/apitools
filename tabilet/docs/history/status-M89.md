# Retired milestone M89 - Optional vendor-neutral query, contract, and recommendation planning

**Milestone.** M89
**Outcome.** completed
**Retired.** 2026-10-11
**Source status.** tabilet/memory-bank/status-M89.md
**Source specification.** tabilet/memory-bank/milestone.md#m89--optional-vendor-neutral-query-contract-and-recommendation-planning
**Evidence.** 2a1c763d078f5a2285f954d72d18be67b6828827
**Worktree.** includes uncommitted changes
**Review.** passed
**Review iterations.** 1
**Verification.** GOWORK=off go test ./...; GOWORK=off go test -race ./...; GOWORK=off go vet ./...; git diff --check; model-free/mock tests for effects, missing input, invalid/oversized output, cancellation, recommendations and explanation limits.
**Consolidated into.** [Architecture](../../memory-bank/architecture.md), [stack](../../memory-bank/tech-stack.md), [lessons](../../memory-bank/lessons.md), [knowledge](knowledge.md).

## Milestone specification

````markdown
### M89 — Optional vendor-neutral query, contract, and recommendation planning

**Goal/acceptance.** Optional planners compile without package cycles/vendor SDKs, preserve deterministic qualification and effect/input semantics, expose missing evidence, and never recommend outside final selectable indexes. No model credential is needed for default use or tests.

**State/owner.** Actionable 2026-10-11; apitools; catalog track. All 3 rows pending; review 0/10, not started. Target behavior is unimplemented.

**Scope and contracts.**

Place PlanSearchQuery, PlanDiscoveryRequest and discovery-related types in root apitools, not catalog. Supply an optional vendor-neutral completion interface and deterministic prepare/validate entry points for isolated callers. Consumers own model clients, credentials, policy, cancellation, privacy and accounting. Standalone APItools CLI stays deterministic; no HTTP model adapter, provider credential environment fallback or model transport flag is scheduled.

Use the current StepContract read/write/unknown/omitted vocabulary, never idempotent_write. Return contract drafts with readiness and diagnostics for consumer review before qualification. Missing input evidence yields clarification or metadata-only leads; no wildcard input or relaxed required-input matching. Use existing PromptBudget/sanitization and strict bounded output parsing. Nil/failing/invalid completions fall back deterministically; cancellation stops rather than launching another call. Model drafts never establish user-provided facts by themselves.

RecommendCandidate takes the final bounded view. Expose only its Selectable candidates to assistance and validate returned indexes against that set. Deterministic fallback uses existing ranking order among eligible candidates with original indexes; empty eligibility yields no recommendation. Omitted/unqualified/remote-only candidates remain ineligible. Advice never selects or executes. Recheck JSON/prompt limits after attaching explanation text.

**Depends on.** APItools:M88 accepted.
**Downstream impacts.** Kinet:M57; Ramen:M82; APItools:M86? compatibility join.
**Lineage.** Retired M80/M81 and M85 deterministic qualification/candidate contracts remain authoritative.

**Rows.** M89.R1 Query and contract planners; M89.R2 Selectable-only recommendation; M89.R3 Mock and consumer qualification. [Status/provenance](status-M89.md).

**Verification.** GOWORK=off go test ./...; GOWORK=off go vet ./...; model-free/mock tests for effects, missing input, invalid/oversized output, cancellation, recommendations and explanation limits; standalone graph/consumer checks; git diff --check. No live model calls.

**Write set.** root planning/projection APIs, dependency-neutral completion contracts, tests and ordinary docs; this status. Provider transport and credentials are consumer-owned. Shared planning/closure is coordinator-owned; no child edits sibling repositories or another status.

**Execution policy.** Planning authorizes no implementation, commit, publication, deployment or live-model operation. A separate execution request governs task commits. One coordinator owns shared records/integration; this catalog track is serial and **Parallel-safe: no**. No overlapping catalog or same-package Stage 12 worker/module writers. Use immutable inputs; sibling adoption requires exact accepted and independently published revisions, not workspace replacements. Publication requires separately named authority; old Stage 12 grants do not transfer. Audit remains disabled/unavailable as recorded.
````

## Status record

````markdown
# Status M89 — Optional vendor-neutral query, contract, and recommendation planning

**State.** Accepted; all rows complete. Whole-milestone review 1/10 passed.
**Context.** Approved catalog track in the current horizon alongside Kinet STG-12, without changing the original browser scope or provisional stages.
**Specification.** [M89](milestone.md#m89--optional-vendor-neutral-query-contract-and-recommendation-planning).

**Planning authority.** The user approved the complete proposal with “Implement the plan” on 2026-10-11; these are planning records only.

**Review source.** Catalog roadmap review (2026-10-11), pasted review and `plan_cross_package_catalog_root_and_search.md` / `plan_llm_touchpoints_in_apitools.md`. Review baseline: not supplied. F01-F13 use source order; this is ordinary intake, not a closing-review iteration.

**Source fingerprints.** Cross-package plan SHA256 `6fbecdb55ceed1d089b63f14b8d064a5bc58e0d0986a9ac5d0ba26f865f464fe`; LLM plan SHA256 `811a33a621b6e4df728e8a97fb6c3fe082e9e0e8d9915a7d06484e285c50b71f`. Source plans remain read-only.

**Revalidation.** `2b74e96f0dd38aeb228d585f750d79194bdc7d06`; clean at assessment and before planning writes. Relevant uncommitted implementation evidence: none. Cross-package baselines: APItools `2b74e96f0dd38aeb228d585f750d79194bdc7d06`, Kinet `31721164eb6a4cef382b17cd3956a3a125a93f72`, Ramen `834cb8e77004c07c6014e4428c2fdb1fe2cb8345`.
**Lineage.** Retired M80/M81 and M85 deterministic qualification/candidate contracts remain authoritative.

## Review ownership

| Finding | Source priority | Local severity | Disposition | Repository evidence | Primary owner |
|---|---|---|---|---|---|
| F10 | P2 | P2 | Confirmed plan defect | client.go and catalog_discovery_types.go define root types; root imports catalog, so catalog cannot define methods on root Client or import it without a cycle. | M89.R1 |
| F11 | P2 | P2 | Confirmed plan defect | operation_rank.go prepareStepContract rejects idempotent_write with contract.effect_invalid. | M89.R1 |
| F12 | P2 | P2 | Confirmed plan defect | operation_rank.go required-input matching has no wildcard semantics; empty inputs or literal * do not satisfy zone_id. | M89.R1 |
| F13 | P2 | P2 | Confirmed plan defect | catalog_discovery_types.go separates Qualified from References; Kinet catalogquestion further restricts eligibility to displayed candidates. Ranking is not selection authority. | M89.R2 |

## Dependencies and ownership

**Depends on.** APItools:M88 accepted.
**Downstream impacts.** Kinet:M57; Ramen:M82; APItools:M86? compatibility join.
**Later implementation write set.** root planning/projection APIs, dependency-neutral completion contracts, tests and ordinary docs; this status. Provider transport and credentials are consumer-owned. Shared planning/closure is coordinator-owned; no child edits sibling repositories or another status.

**Execution policy.** Planning authorizes no implementation, commit, publication, deployment or live-model operation. A separate execution request governs task commits. One coordinator owns shared records/integration; this catalog track is serial and **Parallel-safe: no**. No overlapping catalog or same-package Stage 12 worker/module writers. Use immutable inputs; sibling adoption requires exact accepted and independently published revisions, not workspace replacements. Publication requires separately named authority; old Stage 12 grants do not transfer. Audit remains disabled/unavailable as recorded.

## Task status

Each row is a later implementation commit unit under the governing request. Markers and bounded review/retirement rules come from [milestone.md](milestone.md). Planning completes no row.

| Item | State | Notes |
|---|---|---|
| M89.R1 — Query and contract planners | `[+]` | Delivered vendor-neutral LLMClient interface, PlanSearchQuery, PlanDiscoveryRequest, and deterministic draft preparers/validators in root apitools (resolving F10); enforced standard StepContract effects normalizing idempotent_write to write (resolving F11); rejected wildcard inputs and added contract.input_missing clarification diagnostics (resolving F12); verified offline mock and deterministic fallback tests. |
| M89.R2 — Selectable-only recommendation | `[+]` | Delivered RecommendCandidate with strict selectable-only validation (resolving F13), deterministic ranking fallback, empty eligibility abstention, and decision view bounds/validity rechecks. |
| M89.R3 — Mock and consumer qualification | `[+]` | Ran offline mock and consumer compatibility gates with race detectors; whole-milestone review 1/10 passed with 0 findings; recorded exact accepted source. |

## Acceptance and verification

**Acceptance.** Optional planners compile without package cycles/vendor SDKs, preserve deterministic qualification and effect/input semantics, expose missing evidence, and never recommend outside final selectable indexes. No model credential is needed for default use or tests.

**Verification.** GOWORK=off go test ./...; GOWORK=off go vet ./...; model-free/mock tests for effects, missing input, invalid/oversized output, cancellation, recommendations and explanation limits; standalone graph/consumer checks; git diff --check. No live model calls.

## Review gate and scoped commit tracking

- Whole-milestone review: 1/10, passed; no P1/P2/Lower findings.
- Implementation commits:
  - M89.R1: feat(M89.R1): add vendor-neutral query and contract planners with deterministic fallbacks
  - M89.R2: feat(M89.R2): implement selectable-only candidate recommendation and bounds recheck
- No prior review counter or frozen status is reopened.
````
