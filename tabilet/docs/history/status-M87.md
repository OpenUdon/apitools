# Retired milestone M87 - Explicit catalog roots, safe snapshots, and bounded decisions

**Milestone.** M87
**Outcome.** completed
**Retired.** 2026-10-11
**Source status.** tabilet/memory-bank/status-M87.md
**Source specification.** tabilet/memory-bank/milestone.md#m87--explicit-catalog-roots-safe-snapshots-and-bounded-decisions
**Evidence.** 76057a7deba5fb3c7393ad0693c757deeb90e108
**Worktree.** includes uncommitted changes
**Review.** passed
**Review iterations.** 1
**Verification.** GOWORK=off go test ./...; GOWORK=off go vet ./...; git diff --check; race tests on root, snapshot, decision; immutable Kinet consumer parity on internal/catalogquestion.
**Consolidated into.** [Architecture](../../memory-bank/architecture.md), [stack](../../memory-bank/tech-stack.md), [lessons](../../memory-bank/lessons.md), [knowledge](knowledge.md).

## Milestone specification

````markdown
### M87 — Explicit catalog roots, safe snapshots, and bounded decisions

**Goal/acceptance.** Explicit roots remain read-only; snapshots preserve all configured evidence or safely refuse it; cross-file drift/collisions cannot yield accepted output. Every returned view and prompt is bounded and distinguishes alternatives; omitted, unqualified and remote-only candidates remain unselectable.

**State/owner.** Approved planning 2026-10-11; apitools; catalog track. All 4 rows pending; review 0/10, not started. Target behavior is unimplemented.

**Scope and contracts.**

Add environment-to-RootOptions resolution. Precedence is explicit CLI configuration, APITOOLS_CATALOG_ROOT, CATALOG_ROOT, KINET_CATALOG_ROOT. Canonical registry/index aliases retain legacy fallbacks. ResolveRoot stays read-only: absent configuration produces metadata-only leads without implicit cache lookup or root creation.

Add catalog.SnapshotRoot for APItools consumers and CLI. Preserve registry/index/selected-source paths relative to the source root, creating destination parents. Metadata has an explicit source and destination mapping, not a Kinet-specific library filename. Default ceilings: 64 MiB per registry/index, 2 MiB metadata, 8 MiB selected source. Check file/prefix/sidecar collisions, symlinks/special files, root overlap and active journals. Missing optional registry/index evidence stays missing; selected source and configured metadata must copy successfully. After all copies, recheck the complete source set, missing-file transitions and journal state before destination-only SQLite normalization. Drift returns an error; failures leave no accepted partial snapshot or producer mutation. Kinet retains its host Prepare custody and independent 32 MiB native aggregate limit.

Add ProjectCatalogDecision and view/candidate types in root apitools using existing outcome names. Preserve original indexes/order, scope/contract evidence, omitted counts, provider/source identity, references, operation metadata and qualification gaps. Only qualified displayed candidates with local references are selectable. Enforce independent 64 KiB bounds on serialized view and rendered prompt on every successful return, including negative outcomes, oversized base scope, rollback and explanations. Marshal/validation failures return no accepted view. Selection and execution remain consumer-owned.

**Depends on.** Accepted APItools M80/M81 root/discovery contracts; no new pending prerequisite.
**Downstream impacts.** APItools:M88.
**Lineage.** Retired APItools M80/M81; Kinet W18/M50 are consumer/isolation evidence, not reopened work.

**Rows.** M87.R1 Explicit environment resolution; M87.R2 Bounded snapshot and consistency; M87.R3 Bounded decision projection; M87.R4 Compatibility qualification. [Status/provenance](status-M87.md).

**Verification.** GOWORK=off go test ./...; GOWORK=off go vet ./...; focused root/snapshot/projection race and parity fixtures; existing root/discovery regressions and immutable Kinet consumer fixtures; git diff --check. All data is offline/disposable.

**Write set.** catalog root/snapshot code, root catalog_decision API, tests and ordinary docs; this status. No Kinet host imports or sibling edits. Shared records/closure remain coordinator-owned.

**Execution policy.** Planning authorizes no implementation, commit, publication, deployment or live-model operation. A separate execution request governs task commits. One coordinator owns shared records/integration; this catalog track is serial and **Parallel-safe: no**. No overlapping catalog or same-package Stage 12 worker/module writers. Use immutable inputs; sibling adoption requires exact accepted and independently published revisions, not workspace replacements. Publication requires separately named authority; old Stage 12 grants do not transfer. Audit remains disabled/unavailable as recorded.
````

## Status record

````markdown
# Status M87 — Explicit catalog roots, safe snapshots, and bounded decisions

**State.** Accepted; all rows complete. Whole-milestone review 1/10 passed.
**Context.** Approved catalog track in the current horizon alongside Kinet STG-12, without changing the original browser scope or provisional stages.
**Specification.** [M87](milestone.md#m87--explicit-catalog-roots-safe-snapshots-and-bounded-decisions).

**Planning authority.** The user approved the complete proposal with “Implement the plan” on 2026-10-11; these are planning records only.

**Review source.** Catalog roadmap review (2026-10-11), pasted review and `plan_cross_package_catalog_root_and_search.md` / `plan_llm_touchpoints_in_apitools.md`. Review baseline: not supplied. F01-F13 use source order; this is ordinary intake, not a closing-review iteration.

**Source fingerprints.** Cross-package plan SHA256 `6fbecdb55ceed1d089b63f14b8d064a5bc58e0d0986a9ac5d0ba26f865f464fe`; LLM plan SHA256 `811a33a621b6e4df728e8a97fb6c3fe082e9e0e8d9915a7d06484e285c50b71f`. Source plans remain read-only.

**Revalidation.** `2b74e96f0dd38aeb228d585f750d79194bdc7d06`; clean at assessment and before planning writes. Relevant uncommitted implementation evidence: none. Cross-package baselines: APItools `2b74e96f0dd38aeb228d585f750d79194bdc7d06`, Kinet `31721164eb6a4cef382b17cd3956a3a125a93f72`, Ramen `834cb8e77004c07c6014e4428c2fdb1fe2cb8345`.
**Lineage.** Retired APItools M80/M81; Kinet W18/M50 are consumer/isolation evidence, not reopened work.

## Review ownership

| Finding | Source priority | Local severity | Disposition | Repository evidence | Primary owner |
|---|---|---|---|---|---|
| F02 | P1 | P1 | Confirmed plan defect | Kinet internal/catalogflow/catalog.go Prepare copies metadata and selected source; internal/catalogflow/native.go NativeProjection consumes both. The source plan omits both copies. | M87.R2 |
| F03 | P2 | P2 | Confirmed plan defect | catalog/root.go supports nested registry/index paths; Kinet Prepare/NativeProjection preserve them. filepath.Base in the plan loses them. | M87.R2 |
| F04 | P2 | P2 | Confirmed plan defect | Kinet internal/catalogflow/catalog.go rereads registry/index after all copying and before normalization; the proposed replacement lacks this check. | M87.R2 |
| F08 | P2 | P2 | Confirmed plan defect | Kinet internal/catalogquestion/question.go bounded/Valid validate final JSON and prompt on all successful outcomes; the plan bypasses this for negative/base-scope paths. | M87.R3 |
| F09 | P2 | P2 | Confirmed plan defect | Kinet internal/catalogquestion/question.go Prompt includes source/reference/operation/gap evidence, unlike the proposed method/path-only prompt. | M87.R3 |

## Dependencies and ownership

**Depends on.** Accepted APItools M80/M81 root/discovery contracts; no new pending prerequisite.
**Downstream impacts.** APItools:M88.
**Later implementation write set.** catalog root/snapshot code, root catalog_decision API, tests and ordinary docs; this status. No Kinet host imports or sibling edits. Shared planning/closure is coordinator-owned; no child edits sibling repositories or another status.

**Execution policy.** Planning authorizes no implementation, commit, publication, deployment or live-model operation. A separate execution request governs task commits. One coordinator owns shared records/integration; this catalog track is serial and **Parallel-safe: no**. No overlapping catalog or same-package Stage 12 worker/module writers. Use immutable inputs; sibling adoption requires exact accepted and independently published revisions, not workspace replacements. Publication requires separately named authority; old Stage 12 grants do not transfer. Audit remains disabled/unavailable as recorded.

## Task status

Each row is a later implementation commit unit under the governing request. Markers and bounded review/retirement rules come from [milestone.md](milestone.md). Planning completes no row.

| Item | State | Notes |
|---|---|---|
| M87.R1 — Explicit environment resolution | `[+]` | Delivered RootOptionsWithEnvironment, RootOptionsFromEnvironment, ResolveDefaultRoot with aliases and precedence (CLI/APITOOLS_CATALOG_*/CATALOG_*/KINET_CATALOG_*); verified missing-root and read-only behavior with 100% offline tests. |
| M87.R2 — Bounded snapshot and consistency | `[+]` | Delivered catalog.SnapshotRoot with nested relative path preservation (resolving F03), metadata and selected source copying with distinct ceilings (resolving F02), active journal detection, collision/overlap defenses, full-set drift verification (resolving F04), destination-only SQLite journal normalization to DELETE, and rollback on error. |
| M87.R3 — Bounded decision projection | `[+]` | Delivered ProjectCatalogDecision, CatalogDecisionView, and CatalogDecisionCandidate in root apitools with strict 64-KiB view/prompt bounds for every outcome (resolving F08), full source/reference/operation/gap prompt evidence (resolving F09), local-reference selectable indexing, and answer validation. |
| M87.R4 — Compatibility qualification | `[+]` | Ran owner tests, race verification on root/snapshot/decision, and immutable Kinet consumer parity checks on internal/catalogquestion; whole-milestone review 1/10 passed with 0 findings; prepared retirement record. |

## Acceptance and verification

**Acceptance.** Explicit roots remain read-only; snapshots preserve all configured evidence or safely refuse it; cross-file drift/collisions cannot yield accepted output. Every returned view and prompt is bounded and distinguishes alternatives; omitted, unqualified and remote-only candidates remain unselectable.

**Verification.** GOWORK=off go test ./...; GOWORK=off go vet ./...; focused root/snapshot/projection race and parity fixtures; existing root/discovery regressions and immutable Kinet consumer fixtures; git diff --check. All data is offline/disposable.

## Review gate and scoped commit tracking

- Whole-milestone review: 1/10, passed; no P1/P2/Lower findings.
- Implementation commits:
  - M87.R1: catalog: add environment resolution and precedence for catalog roots
  - M87.R2: catalog: add safe bounded SnapshotRoot with drift defense and normalization
  - M87.R3: root: add bounded ProjectCatalogDecision and evidence-bearing decision cards
  - M87.R4: docs(M87): qualify compatibility and retire milestone M87
- No prior review counter or frozen status is reopened.
````
