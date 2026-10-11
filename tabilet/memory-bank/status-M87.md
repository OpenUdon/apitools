# Status M87 — Explicit catalog roots, safe snapshots, and bounded decisions

**State.** Pending; no row started. Whole-milestone review 0/10, not started.
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
| M87.R3 — Bounded decision projection | `[ ]` | Implement final JSON/prompt validation for every outcome, source-distinguishable prompts and stable selectable indexes. Owns F08-F09 (P2/P2). |
| M87.R4 — Compatibility qualification | `[ ]` | Run owner and immutable consumer/parity gates; record accepted source and publication handoff. Publication requires separately named authority and independent verification before consumer adoption. |

## Acceptance and verification

**Acceptance.** Explicit roots remain read-only; snapshots preserve all configured evidence or safely refuse it; cross-file drift/collisions cannot yield accepted output. Every returned view and prompt is bounded and distinguishes alternatives; omitted, unqualified and remote-only candidates remain unselectable.

**Verification.** GOWORK=off go test ./...; GOWORK=off go vet ./...; focused root/snapshot/projection race and parity fixtures; existing root/discovery regressions and immutable Kinet consumer fixtures; git diff --check. All data is offline/disposable.

## Review gate and scoped commit tracking

- Whole-milestone review: 0/10, not started; begin after implementation and required verification.
- Implementation commits:
  - M87.R1: catalog: add environment resolution and precedence for catalog roots
  - M87.R2: catalog: add safe bounded SnapshotRoot with drift defense and normalization
- No prior review counter or frozen status is reopened.
