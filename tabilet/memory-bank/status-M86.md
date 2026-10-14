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

## Status rules

See [`milestone.md`](milestone.md) for milestone definitions and the status ID
pattern that names this file.

Status markers:

| Symbol | Suggested Status | Interpretation |
|---|---|---|
| `[ ]` | Pending | Item is actionable when its dependencies pass. |
| `[+]` | Completed | Item finished or done. |
| `[~]` | In Progress | Item is the one general row currently being worked. |
| `[!]` | Blocked | A current unresolved blocker still prevents progress. |
| `[X]` | Cancelled | The approved item was withdrawn; any later replacement needs fresh planning. Cancellation does not prove delivery. |
| `[-]` | Closed Historical | A consumed failed attempt or superseded row retained for audit; it is never retried and does not block its accepted successor. |

Write markers with backticks, exactly as in the table above: `` `[ ]` ``, not
`[ ]`.

Across the active ledger, zero or one general row may be `[~]`, except under the local
per-lane override in [`milestone.md`](milestone.md#parallel-work-and-priority). Resume it
before selecting another row. An operational launcher additionally requires its exact
authorized operation row to be `[~]` before invocation; the marker records the
selected operation but does not grant missing external-mutation authority.

Read this milestone's optional `AUTHORIZATION_REQUIREMENTS` in
[`milestone.md`](milestone.md#authorization-requirements); do not duplicate its
declaration here. Check protected actions against the current human-approved
request and effective grants before acting. Notes may record a grant ID,
performed-action evidence, or a missing-authority blocker; they never prove
approval. Required unperformed actions prevent closure. No new task marker is
introduced, and absent authorization metadata grants no additional authority.

With human-approved `CORRECTIVE_PLANNING: auto`, the owner may prepare and
check necessary repair plans within the original delivery outcome, existing
grants and explicit limits. Preserve this row's identity and interrupted work;
move the interrupted `[~]` row to `[!]` before activating a repair row. Record
prerequisite IDs, standing approval source and resumption condition. Keep zero
or one `[~]` row. Once every prerequisite is verified and accepted and all
blockers clear, restore this row from `[!]` to `[~]` under the same approval.
Recover this chain on resume; do not select unrelated work. Preserve partial
implementation; suspension bookkeeping may join the scoped repair checkpoint,
but unfinished implementation cannot. Pause if overlapping changes cannot be
separated safely. Explicit project restrictions still govern.
Preserve review counts and frozen history. Children report needs
to the owner; a marker never grants planning or action authority.

Each table row is a commit unit: after changing its state, verify the change,
update the memory bank/docs, and make a scoped `git commit` before starting the
next row, except the scoped repair suspension bookkeeping above. A `[-]` row's
notes must record the consumed attempt or supersession
and identify its accepted successor. If two rows are inseparable, redefine them
as one row before starting rather than closing several rows in one commit.
During a goal run, the resolved `COMMIT_POLICY` governs instead; `none` means
no commits even when a requirement or grant mentions `git.commit`.

Keep this file a current-state record, not a log. Replace a superseded checkpoint
instead of appending dated narrative, and link per-attempt detail in evidence
files. A row whose acceptance depends on an expensive check records its attempt
count here under the [verification-attempt gate](milestone.md#verification-attempt-gate).

After the milestone's review, consolidation, and downstream reconciliation
pass, follow [the retirement procedure](milestone.md#long-term-memory-and-retirement).
Preserve this complete document and the full milestone specification in its
retired record; earlier rows, notes, and original path context must survive.
Retirement follows the governing commit policy and never changes the task
markers or permits reuse of this milestone ID.

## Dispatch and lease boundaries

**Depends on.** UWS:C11. All required prerequisites must have
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

**Upstream.** UWS:C11 (accepted and published).

**Downstream.** [Kinet:M55](../../../kinet/tabilet/memory-bank/status-M55.md).

Accepted M85 candidate and shape contracts stay unchanged.

## Tasks

| Item | State | Notes |
|---|---|---|
| M86.1 — Adopt HCL-free UWS core | `[ ]` | Adopt exact accepted/published C11 ordinarily with `GOWORK=off` and no workspace/directory substitution. Prove HashiCorp/Horizon HCL modules are absent from APItools selected and compiled closures, with tidy/go.sum evidence. C11 adapted migration fixtures are pre-publication compatibility evidence, not this actual adoption proof. |
| M86.2 — Source-family regression | `[ ]` | Regression tests for all eight source families: OpenAPI/Swagger, Discovery, Smithy, AsyncAPI, GraphQL, OpenRPC, gRPC/protobuf and OData. Covers shape tables, candidates and diagnostics, with exact digests unchanged. |
| M86.3 — Qualify and publish | `[ ]` | Consumer builds (OpenUdon, Kinet workers), tidy check, and a publication handoff under named authority. |

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
