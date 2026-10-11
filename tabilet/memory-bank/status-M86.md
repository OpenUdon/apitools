# M86 — HCL-free UWS adoption (conditional)

**Stage:** Kinet STG-12, Phase C. **Owner:** APItools.
**State:** Pending and conditional. No row has started. Review 0/10.
**Trigger:** [UWS:C11](../../../uws/tabilet/memory-bank/status-C11.md) is accepted
and published. That in turn requires a clean, authorized
[Kinet:M54](../../../kinet/tabilet/memory-bank/status-M54.md) inventory. If the
trigger is absent, the goal skips this milestone without completing or
cancelling it.
**Source baseline:** `f2837e7b266a0f8d1d379c4dd7820bba7c5f2179` (clean at planning).
**Coordinator:** [Stage 12 contract](../../../kinet/docs/stage12.md#conditional-hcl-phase).
Lineage: Stage 11 recorded that APItools' UWS adoption retains the
HashiCorp/Horizon HCL transitive closure until a separately approved removal.

**Planning reconciliation.** Stage 12 planning review, 2026-10-09; source priority not supplied. Review baseline and current revalidation: `f2837e7b266a0f8d1d379c4dd7820bba7c5f2179`; includes the uncommitted Stage 12 planning files. F07 (confirmed P2; pre-publication owner UWS:C11) places actual APItools adoption in M86.1 after publication. F04/F06 closure distinctions preserve the APItools-only HCL-free requirement. Evidence: `go.mod`, the original C11.3 consumer gate and `../uws/uws1/hcl.go`; M85 accepted contracts remain unchanged. This is approved review intake, not a closing-review iteration; all rows and review counters remain pending/0.

**Parallel planning provenance.** Stage 12 parallel execution review,
2026-10-09; source priority and separate review baseline not supplied. Current
revalidation `f2837e7b266a0f8d1d379c4dd7820bba7c5f2179`, including uncommitted planning files. F01 (confirmed local P2) adopts scoped workflow ownership; F02 (confirmed Lower) replaces strict ordering with readiness; F04 (confirmed local P2) requires frozen consumer checks.
Evidence: the old Kinet goal/launch rules, owning agent rules, existing pending
dependencies and sibling consumer checks; Udon `go.mod`/`pkg/execute` import no
OpenUdon code. This approved intake changes no row state or review counter.

## Catalog compatibility reconciliation — 2026-10-11

The approved proposal adds accepted M89 (including M87/M88) before M86. Original C11 publication and authorized M54 inventory triggers remain unchanged. Preserve accepted root/snapshot/decision/discovery/planning contracts in M86.2 and qualify consumers in M86.3. This is downstream compatibility, not duplicate finding ownership. Baseline `2b74e96f0dd38aeb228d585f750d79194bdc7d06`, clean before planning; rows stay pending and review 0/10.

## Dispatch and lease boundaries

**Depends on.** UWS:C11 and APItools:M89. All required prerequisites must have
accepted closure at exact revisions; sibling producer adoption also needs
independently verified publication. A priority position never supplies authority.

**Downstream impacts.** Kinet:M55?.

**Write set.** The owning `apitools/` package's implementation, tests, ordinary
documentation, manifests and qualification outputs only as required by this
milestone's existing scope, plus `tabilet/memory-bank/status-M86.md` in its
assigned worktree. Excludes `AGENTS.md`, `tabilet/GOAL.md`, shared memory-bank
files, other statuses, evolution, stages, history/knowledge, the package audit
database/sidecars, coordination docs and launch input. The coordinator alone applies shared-memory and closure
changes serially; no child writes a sibling repository or user ledger.

**Contracts read.** Immutable exact prerequisite artifacts listed above, the
M51 native-owner-reviewed contract/fixtures when applicable, the assigned
package baseline and frozen shared-memory/consumer snapshots captured at
dispatch. Cross-package checks use read-only exact snapshots or approved
published module inputs, never changing sibling checkouts. Record full source,
artifact and fixture hashes in the later execution brief; contract drift pauses
affected leases for coordinator reconciliation. Existing no-workspace/no-directory
substitution requirements for ordinary published adoption remain in force.

**Parallel-safe.** no. Execute sequentially; this milestone is on a required integration/contract chain.
At most one live milestone per package. All tests use private lease ports,
disposable stores/caches/browser profiles and unique output directories.

## Dependencies and handoff

**Upstream.** UWS:C11 (accepted and published), plus [M89](status-M89.md) accepted with M87/M88 contracts.

**Downstream.** [Kinet:M55](../../../kinet/tabilet/memory-bank/status-M55.md).

Accepted M85 candidate and shape contracts stay unchanged.

## Tasks

| Item | State | Notes |
|---|---|---|
| M86.1 — Adopt HCL-free UWS core | `[ ]` | Adopt exact accepted/published C11 ordinarily with `GOWORK=off` and no workspace/directory substitution. Prove HashiCorp/Horizon HCL modules are absent from APItools selected and compiled closures, with tidy/go.sum evidence. C11 adapted migration fixtures are pre-publication compatibility evidence, not this actual adoption proof. |
| M86.2 — Source-family regression | `[ ]` | Regression tests for all eight source families: OpenAPI/Swagger, Discovery, Smithy, AsyncAPI, GraphQL, OpenRPC, gRPC/protobuf and OData. Covers shape tables, candidates and diagnostics, with exact digests unchanged. Preserve accepted M87-M89 explicit-root, snapshot/projection, remote-policy and optional-planning regression contracts. |
| M86.3 — Qualify and publish | `[ ]` | Consumer builds (OpenUdon, Kinet workers and Ramen catalog APIs), tidy check, and publication handoff under named authority; verify accepted M89 contracts survive adoption. |

## Acceptance and verification

**Acceptance.** No HCL module remains in the APItools graph, and every source
family and published shape contract behaves identically.

**Verification.** `GOWORK=off GOPROXY=off go test ./...`, `go vet ./...`, the
tidy check, consumer builds and `git diff --check`.

## Execution policy

One coordinator owns the integrated ledgers, shared memory and serialized
integration/closure. Serial execution remains the default. Concurrent leases
require this milestone's declared safety, frozen inputs, a complete explicit
Kinet goal request and the Stage 12 agent-rule opt-in. Each lease has one
in-progress row and one assigned milestone; its persisted review count survives
resume/rebase. Commit policy comes from that later request. Source publication,
deployment and live operations retain separate named authority. Planning and
status markers grant none; audit stays disabled.

## Review

Whole-milestone review: 0/10, not started.
