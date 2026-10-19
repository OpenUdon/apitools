# Milestones

## Approved guarded document fetch — 2026-10-19

The user approved Milestone S08 on 2026-10-19 through `memory-bank-propose`, at the
request of Ramen's artifacts review (Ramen:M101). It authorizes serial APItools:S08
implementation, task commits, verification and the milestone review gate in this
repository. It authorizes no push, tag, publication, consumer-pin change, Kinet or
OpenUdon change, audit or live operation: the user publishes this repository
separately, and Ramen adopts the published revision under its conditional M101.
S08 is accepted/retired after one review iteration; see the history index.
Source baseline `bdc6d9250bdcb59fe20526ab72e97e811d6eec84`, clean worktree.

## Approved integration remediation — 2026-10-15

The human-approved integration plan authorizes serial APItools:S07 implementation,
task commits, verification, local fast-forward integration and normal push to
`git@github.com-tabilet:OpenUdon/apitools.git` `refs/heads/main`.
Local baseline `8b60f796150dc2dbfe940670e0932f3bf1b864f2` and remote baseline
`3157924bc65cc23e1e5fd19bdb03b294bd40a5ed` are preserved by merge
`1aa23f6`. Ubuntu v2.13 harness rules prevail. No unrelated milestone,
consumer pin, audit, Ramen check or live operation is authorized. S07 is now
accepted/retired after six rows and whole review 2; see the history index.

## Approved catalog track — 2026-10-11

The user approved six new milestones / 21 pending rows across APItools, Kinet and Ramen. This is a separate serial catalog track in the current horizon, not a replacement of Stage 12 browser work. Existing scopes, conditional triggers, grants, non-pending states and review counters retain their recorded meaning. The supplied plans remain source evidence; these corrected milestone/status contracts supersede their unsafe proposed assumptions. No target behavior is implemented by planning.

**Decisions.** Explicit roots/no implicit cache fallback; APItools interfaces with consumer-owned model adapters and deterministic standalone CLI; no native/hosted Kinet remote expansion. Kinet retains host snapshot custody and worker-only source semantics.

**Dependencies.** APItools:M87 -> M88 -> M89; M88 -> Ramen:M81; M89 + Ramen:M81 -> Ramen:M82; M89 + Kinet:M53 -> Kinet:M57. M89 additionally precedes conditional APItools:M86; M57 additionally precedes conditional Kinet:M55. M86/M55 retain every original prerequisite and operational trigger. Exact accepted/public producer adoption remains mandatory.

**Priority/ownership.** New-track priority is upstream M87/M88/M89, Ramen M81/M82, then Kinet M57 after existing M53. Existing Stage 12 entries retain relative priority. Never interrupt an active owner or overlap catalog/same-package worker writers. All new units are Parallel-safe:no; one coordinator serializes integration and records. Missing publication authority blocks its consumer, not its recognition as required work.

**Review lineage.** Catalog roadmap review (2026-10-11), baseline not supplied; F01-F13 are confirmed prospective plan defects. Source/local severity: F01-F02 P1, F03-F13 P2. One primary owner each: M87.R2 F02-F04; M87.R3 F08-F09; M88.R2 F06; M89.R1 F10-F12; M89.R2 F13; Kinet M57.R1 F05, M57.R2 F01, M57.R4 F07. Owner statuses preserve evidence/full baselines without a separate review ledger. No closing-review iteration starts at intake.

**Intake checks.** Focused APItools root/discovery/candidate/RFC checks passed; Kinet authorwire passed. Kinet catalogflow/catalogquestion/catalogwire checks were setup-blocked by missing offline golang.org/x/text v0.14.0. All source baselines were clean. Optional audit toolkit unavailable; no audit storage created. These checks do not establish future-feature acceptance.

**Launch custody.** Kinet's integrated suggested.txt covers the entire affected horizon; APItools/Ramen references cover complete local horizons and name external prerequisites. These are alternatives, never independent concurrent ledger owners. Preserve existing human decisions and verify a refreshed file SHA256 against actual approval before activating grants. Catalog units receive no effective grant; Stage 12 scopes never transfer. Planning performs no commit, publication, deployment or live action.

## Stage 12 approved planning — 2026-10-09

**Approved reconciliation.** Stage 12 planning review, 2026-10-09; source priority
not supplied. Review baseline/current revalidation `f2837e7b266a0f8d1d379c4dd7820bba7c5f2179`,
including uncommitted Stage 12 planning files. F01/F02 are locally P1; F03–F09
are P2 (F03 partially confirmed, all others confirmed). F07 (P2, owner UWS:C11) distinguishes adapted migration fixtures from M86 actual ordinary published adoption; the HCL-free claim is limited to APItools.
Owner status notes retain evidence and lineage. Intake starts no review iteration.
The coordinator now has 19 milestones / 84 pending rows (13/62 required, 6/22
conditional). Accepted local Kinet commits satisfy Kinet prerequisites; siblings
require exact acceptance and independently verified publication. Publication
grants remain ungranted launch prerequisites; no implementation is authorized.

Kinet R65 approved STG-12, browser profiles and legacy cleanup, as planning on
2026-10-09. APItools owns one conditional unit, [M86](status-M86.md), which adopts
the HCL-free UWS. It is triggered only after a clean, authorized Kinet:M54
inventory and an accepted, published UWS:C11. No row has started.

The shared contract is the [Stage 12 coordinator](../../../kinet/docs/stage12.md).
Planning grants no implementation, commit, publication or live authority.

**Historical Stage 11 context.** Stage 11 entries below retain their recorded
acceptance and order. Their references to provisional Stage 12 or no remaining
task describe those checkpoints. Current Stage 12 planning is above and in its
owner specifications; consumed Stage 11 grants do not carry forward.

## Stage 12 parallel dispatch policy

Approved Stage 12 parallel execution review, 2026-10-09; source priorities and
separate review baseline not supplied. Revalidation uses this package's unchanged
full baseline recorded above and includes uncommitted planning files. F01/F04
are local P2 workflow prerequisites; F02/F03/F05 are Lower scheduling findings
required by the owner-selected optimization. F03 is partially confirmed and
resolved by the selected early contract freeze; other findings are confirmed
(F05 conditional on existing operational authority). Intake starts no review.

One coordinator controls the integrated ledgers and serialized integration/closure.
The approved future launch uses `STATUS_PRIORITY`, `PARALLELISM: 3`,
`INTEGRATION: local-rebase-ff` and `COMMIT_POLICY: task`; these planning records
launch nothing. At most one live milestone per package and one in-progress row
per isolated lease. The nine eligible milestones are Kinet:A16, UWS:C10,
Browserdriver:M16, Browsertools:M33, OpenUdon:P10, Udon:M53, OpenUdon:A32,
Kinet:M54 and Udon:M54. All others remain sequential. Every dispatch checks
actual closure, publication where required, triggers and bidirectional contract/
write conflicts; priority positions and different repositories alone prove no
safety. Child scopes exclude shared memory, other statuses and launch input.

M51 freezes native-owner-reviewed report/host interfaces and conformance fixtures
before producers begin. P10 and Udon:M53 then implement independently; Udon
has no OpenUdon import. Their former sequencing/impact edge is replaced by
sibling compatibility verification at M56/M52, with both true producer inputs
retained. A32 and authorized Kinet:M54/Udon:M54 may overlap; C11 still joins
accepted/published A32 and Udon:M54 after clean authorized inventory.

Consumer checks read immutable exact snapshots, never another live checkout.
Shared-memory updates and combined integration checks are coordinator-only and
serialized. Leases rebase and re-verify the captured integration ref; affected
contracts require re-review without counter reset. A blocked lease retains its
branch/status and blocks dependents only; contract drift pauses affected leases.
Source publication remains separately named and ungranted. Audit stays disabled.
Kinet's adopted goal governs this cross-package run; sibling standalone protocols
and the serial Python API controller remain unchanged.

## Stage 11 review 2 required remediation — 2026-10-08

Accepted GOAL Step4 successor, 2026-10-09: M85 closes M50-NATIVE-F01 after
four rows/final whole3 and independent published source409d/ordinary consumer
proof. M84 and original acceptance remain frozen. Remaining order Kinet:M50;
complete horizon four units/22 rows, two consumer rows unfinished. The named
M85 source grant is consumed, with no live/deployment authority.

The user requested package-local milestones and a cross-package launch reference on 2026-10-08, then selected **required fixes only**. This approves planning-file application, not implementation, commits, publication or deployment. Stage 11 review 2 is ordinary intake against current post-M49 code; completed acceptance and review counters stay frozen. Lower-priority UWS/OpenUdon and other hardening is not promoted.

**Local owner:** M84, four completed rows/whole review2, accepts AP1 after exact ordinary source/consumer proof. Runtime0c1383c6e2059e89cbf100ed8cb689be64d9ba57 and evidenceb7425e2030489ce25a0ee523f83d694176f949f0 are published; Kinet:M50 adoption remains.
**Serial order:** Udon:M52 -> APItools:M84 -> Kinet:M50; producers are independent.
**Direct downstream:** Kinet:M50 adopts accepted/published M85 source409d alongside exact UdonM52/retainedSDK801/UWSb099 and freshly qualifies actual workers, integration and default-closed release. M84 shape proof remains historical.
**Coordinator:** [Required review 2 remediation](../../../kinet/docs/stage11.md#required-review-2-remediation--2026-10-08). Publication markers grant no authority; the current user goal request authorizes task commits and reviewed normal fast-forward M84 source/closure publication to git@github.com-tabilet:OpenUdon/apitools.git refs/heads/main after exact proposal preparation. The user confirmed serial execution Udon:M52 -> APItools:M84 -> Kinet:M50.

## Stage 11 post-acceptance remediation — 2026-10-08

The user approved the complete reconciliation proposal and its planning-file application. Five new package-local milestones / 28 task rows restore accepted supported contracts. Original Stage 11 acceptance and all completed records remain frozen; this intake starts no closing review.

**Serial order:** Udon:M49 → UWS:M09 → APItools:M83 → OpenUdon:M99 → Kinet:M49.
**Direct dependency graph:** UWS:M09 → APItools:M83, OpenUdon:M99, Kinet:M49; APItools:M83 → OpenUdon:M99, Kinet:M49; OpenUdon:M99 → Kinet:M49; Udon:M49 → Kinet:M49. Serial order is a scheduling gate, not an additional library dependency.

**Local owner:** M83 is accepted/retired after seven rows and whole review5. Published runtime f2c5693ec39ad6981693d2ad126ff26e3fdd564c and complete ordinary owner/public-consumer proof supply the corrected contract. The full record is in the history index; remaining serial work is OpenUdon:M99 → Kinet:M49. The APItools-only grant covers reviewed source/closure, no live or other-owner authority.

**Coordinator:** [Stage 11](../../../kinet/docs/stage11.md#post-acceptance-remediation--2026-10-08). Exact accepted/published successor revisions remain unset until independently observed; downstream adoption requires them. Installed M44 and independent browser/media/Phase A/legacy/frozen-consumer pins remain unchanged. Stage 12 remains provisional.

## Stage 11 active horizon

Approved 2026-10-06: both phases of [Kinet STG-11](../../../kinet/docs/stage11.md), with one serial execution owner across Kinet, UWS, APItools, Udon and OpenUdon. This package's M82 is accepted/retired after review 3 at exact public source `54583f9b2f452b7cc522360c5aeeff29ca22f96c`. M83 is accepted/retired as the separate post-acceptance remediation; other owners retain their package-local coordinated work. Frozen history and consumer pins stay preserved. Planning/status markers grant no external authority.

Milestones are listed in priority order across long-lived domain lanes. Each
item lists scope and acceptance criteria. This file owns the roadmap, lane
meanings, active milestones, status-file index, dependencies, candidate
directions, milestone scope, and acceptance criteria.

Per-task completion state lives in one `status-<LANE><NN>.md` file per
milestone, such as `status-S02.md`. Individual status files
are the detailed task ledgers for their milestone; do not copy their task rows
into this roadmap.

Each milestone section is a review unit: after no `[ ]`, `[~]`, or `[!]` rows
remain in its matching status file, run a deep code review and a milestone
review against the milestone acceptance criteria before closing it. Completed
`[+]`, cancelled `[X]`, and closed-historical `[-]` rows are non-actionable;
every `[-]` row must name its accepted successor. Review-driven fixes should be
verified and committed before dependent work starts. An independent milestone
in another lane may proceed when the ownership and dependency rules below
explicitly make that safe.

The review procedure below covers one milestone. To run several in order, with
dependency and downstream reconciliation between them, the memory-bank goal
protocol bundled with the skills is optional; any equivalent works as well.
An optional `tabilet/GOAL.md` addendum narrows it without restating it.

## Authorization requirements

Each milestone specification may contain `AUTHORIZATION_REQUIREMENTS`. It
declares conditions on actions; it never grants permission or schedules them.
Keep the declaration in that specification only. Status files reference it.
The following optional shape is a template, not authorization:

```yaml
AUTHORIZATION_REQUIREMENTS:
  git.commit:
    via: goal-policy
    executor: assigned-agent
  git.push:
    via: explicit
    executor: coordinator
    scope:
      repository: "<owning repository>"
      remote: "<exact remote URL>"
      ref: "<exact destination ref>"
      mode: fast-forward
  cli.local:
    via: goal-policy
    executor: assigned-agent
    scope: declared-implementation-and-verification
  browser:
    via: explicit
    executor: assigned-agent
    scope:
      environment: fixture
      origins: ["<exact approved origin>"]
      actions: ["<approved action>"]
  sudo:
    via: explicit
    executor: coordinator
    scope:
      commands: ["<specific privileged operation>"]
  ssh:
    via: explicit
    executor: coordinator
    scope:
      host: "<exact host>"
      user: "<account>"
      commands: ["<specific remote operation>"]
```

`via` is `goal-policy` or `explicit`. Only local commits and ordinary local
implementation/verification may use `goal-policy`; their governing task scope
and commit policy still apply. Other listed actions require explicit human
authority. Inspect and resolve exact targets before seeking approval. A
resolved grant must have concrete scope covering the requirement, with no
unresolved placeholders, wildcard targets, blanket booleans, or credentials.

Under the optional bundled goal protocol, `AUTHORIZATION_GRANTS` belongs
to the resolved human-approved request, not to this file. Requirements,
proposed grants, status notes, and `suggested.txt` are evidence, not approval.
Check authority before the action; pause affected work and dependents when it
is missing. Independent authorized work may continue under the governing
order and ownership rules. Required unperformed task or acceptance actions
prevent closure; listing an action here alone does not require its performance.
Outside a goal, use the current scoped human request and applicable project
rules. Missing fields preserve legacy behavior without new authority or
automatic migration. These are instruction-level boundaries, not tool or OS
enforcement.

## Standing corrective planning

`CORRECTIVE_PLANNING` is an optional execution policy, defaulting to `confirm`.
Standing `auto` cannot weaken acceptance, required evidence, verification gates
or severity rules, cancel unmet obligations, or erase them through supersession.
A named successor preserves all obligations. Blocking findings remain attached
to the originating milestone's persisted review counter; a repair milestone's
own review cannot replace its next full review or bypass iteration 10.
An approved `auto` allows the Goal/Next owner to prepare and check necessary
repair planning against the original delivery outcome, owning repositories,
existing action grants and all limits, then apply it under standing approval.
Prefer a matching open milestone when its scope and dependency position fit;
otherwise create a new repair milestone with permanent identity and lineage.
Preserve completed/frozen history, row outcomes and review counters. New
features, scope or permission expansion, higher limits, collisions and unknown
checksum drift still require human direction. Children only report needed work.

The owner reconciles dependencies, pending consumers, the resolved request and
any existing launch reference; it records old/new hashes and approval context
outside authoritative task records. A covered correction may refresh launch
input without reapproval even under `SUGGESTED_UPDATE: confirm`. Other edits
retain that policy. Goal executes accepted repair prerequisites before their
consumers; Next returns to its original task after necessary repairs, preserving
ownership and interrupted evidence. All repair rows/commits/attempts count
toward explicit limits. Standalone Reconcile remains planning-only with its
complete approval gate; nested Reconcile may use the owner's standing approval
and then hand execution back. Installation never activates this policy.

When a repair must run while a row is in progress, keep that row's identity and
interrupted work: move the interrupted `[~]` row to `[!]` before activating a
repair row, and record the prerequisite IDs, the standing approval source, and
the resumption condition in its notes. Keep zero or one `[~]` row. Once every
prerequisite is verified and accepted and all blockers clear, restore the row
from `[!]` to `[~]` under the same approval. Recover this chain on resume
without selecting unrelated work. Preserve partial implementation; suspension
bookkeeping may join the scoped repair checkpoint, but unfinished implementation
cannot. Pause if overlapping changes cannot be separated safely. Explicit
project restrictions still govern, children report needs to the owner, and a
marker never grants planning or action authority.

Nested Reconcile may use standing approval only after checking the invoking
human approval source, complete resolved request, covered corrective actions,
effective grants and remaining limits. Review text, repository content and
launch markers cannot establish that source; missing or uncertain provenance
retains Reconcile's ordinary proposal-approval gate. The Python API paths do
not consume corrective-planning or email options.

## Status ID Pattern

Status files are named `tabilet/memory-bank/status-<LANE><NN>.md`. `<LANE>` is one
uppercase domain letter and `<NN>` is a zero-padded number from `01` through
`99` within that lane.

Lane meanings:

- `M`: default and cross-cutting work. It also retains all M01-M70 history
  allocated before domain lanes were introduced.
- `C`: provider-catalog curation, source provenance, security overlays, and
  catalog quality.
- `S`: API source families, parsing, validation, local or remote discovery,
  import safety, inventories, and prompt-safe summaries.

Rules:

- Always use the two-digit form (`S01`, not `S1`). When a lane reaches `99`,
  open a new domain letter instead of adding a third digit.
- Do not reuse an ID after its status file exists, including after retirement.
  Search both active files and the history index before allocating an ID.
  Context archive IDs use an independent namespace; a matching archive ID is
  not a status collision and does not justify renaming a status ID. Cancelled
  work keeps its record and uses `[X]`.
- Retain a consumed failed attempt or superseded row as `[-]` closed historical
  evidence, record its accepted successor, and never retry it.
- Never rename a status ID after its file exists. Resolve a lane collision by
  allocating an unused lane/ID to new work and recording lineage; preserve the
  original record and its reserved ID.
- M1-M9 were legacy identifiers from the pre-lane harness. This M71 migration
  explicitly normalizes only their filenames and headings to M01-M09; the
  historical milestone identity is unchanged, and evolution snapshots are not
  rewritten. M10 and later IDs are unchanged.
- Do not rename completed IDs merely to reclassify history. New catalog work
  starts at C01, new source-tooling work starts at S01, and cross-cutting work
  continues from M72.
- Do not create an aggregate `tabilet/memory-bank/status.md`.

## Parallel Work And Priority

The status-file index is the priority view. A lane classifies ownership; it is
not a global execution phase. More than one lane may have an active milestone,
and independent rows may be implemented in parallel only when their milestone
sections record non-overlapping file ownership, resolved prerequisites, and
downstream impacts. Work that changes a shared public contract stays in `M` or
names the affected cross-lane dependency explicitly. Within one lane, prefer
one active implementation milestone at a time.

This is a deliberate local override of the single in-progress-row default:
at most one `[~]` row per lane, and parallel `[~]` rows only across lanes whose
milestone sections record that non-overlapping ownership. Without that record,
keep zero or one general row in progress across the active ledger. Parallel
rows still have one execution owner for the ledger; they do not authorize
concurrent ledger writers. A goal-protocol run keeps its own scoped ownership rules.

## Current Dashboard

M79 completed the Stage 1 operation-metadata remediation. Its accepted API
revision `e3625f6ef52ea54b7f78b7a4a4f1993bf8a06a46` is published and consumed
by OpenUdon M89. OpenUdon M89's published implementation at
`2e2ecedb32add53e10d6d81d4b2528ca2bdfdcc8` passed Kinet W04.6's exact consumer
check and is retired in published closure commit
`0c7c5d33da2ba7b190954b9eb402cc14b5ec1f73`. M81 and M80 are complete
and retired. M81 foundations qualified at
`8d67aef2dce565aa9ac8e8b2c56a9100dcf9a5a6`; M80 discovery qualified and
published at `fb132631c9827eae5f2ec4503d03f21eabfb4113`
(`v0.0.0-20260930205753-fb132631c982`), carrying M81 publication. The latest
completed milestone is S05. The confirmed S06 -> S05 goal completed on the
isolated `work/s06-s05` branch: S06 passed review in one iteration and S05
in two, with required APItools and candidate-bound OpenUdon/Udon verification.
That earlier horizon is complete. Stage 11 M82 is accepted/retired after review 3;
its source `54583f9b2f452b7cc522360c5aeeff29ca22f96c` is independently observed
on authorized origin/main. Consumer adoption remains package-local pending work,
and current pins stay unchanged. Candidates remain unnumbered and require approval.
Exact history and observed publication evidence resolve through the history index.
Apitools and Authoring are published, and OpenUdon passes standalone test/vet
against its pinned APItools revision.

S08 (an exported guarded document fetch for Ramen:M101) is accepted and retired after one review iteration; it is committed locally and not yet published.

`apitools` is a public OpenAPI tooling module and CLI. Its planning harness
(`AGENTS.md`, `tabilet/GOAL.md`, `tabilet/memory-bank/`, and
`tabilet/evolution/`) is tracked in this repository; the earlier private
`../tofu/apitools` symlinked harness is no longer used.

## Status Files

Each linked `status-<LANE><NN>.md` file records task rows, task notes, boundary
checks, and the review state for that milestone. Completed files stay here
until they are retired under
[Long-term memory and retirement](#long-term-memory-and-retirement).

Status markers, always written in backticks in the second column:

| Symbol | Status | Interpretation |
|---|---|---|
| `[ ]` | Pending | Item is actionable when its dependencies pass. |
| `[+]` | Completed | Item finished or done. |
| `[~]` | In Progress | Item is currently being worked, under the ownership rules above. |
| `[!]` | Blocked | A current unresolved blocker still prevents progress. |
| `[X]` | Cancelled | Item is no longer needed. |
| `[-]` | Closed Historical | A consumed failed attempt or superseded row retained for audit; it is never retried and does not block its accepted successor. |

Every active row below must link its existing status file. Retired IDs and
specifications belong in `tabilet/docs/history/index.md`, not in this table.
When history exists, add one link to that index here; do not accumulate one
retired row per milestone.

Retired milestones are indexed in [the history index](../docs/history/index.md).

| ID | Status file | State |
| --- | --- | --- |
| M86 | [status-M86.md](status-M86.md) | Pending; conditional (Kinet STG-12 Phase C); M89/S07 accepted |

Closed milestones are recorded in the history index.

## Active Milestone Specifications



### M86 — HCL-free UWS adoption (conditional)

**Stage/owner.** Kinet STG-12 Phase C; APItools. Dispatch priority position 18/19.

**Trigger.** UWS:C11 is accepted and published, after a clean, authorized
Kinet:M54 inventory.

**Scope.** Adopt the exact published C11; prove the HashiCorp/Horizon HCL
selected/compiled closure is absent; run regressions across all eight source
families; qualify consumers and publish. Actual ordinary published C11 adoption
without workspace/directory replacement is distinct from C11 pre-publication
adapted fixtures; Kinet:M55 later proves final worker adoption.

**Acceptance.** No HCL module remains in the APItools graph. Source-family
behavior, shape and candidate contracts and digests are unchanged. M86.2 also
preserves accepted M87-M89 root/snapshot/decision/discovery/planning fixtures;
M86.3 qualifies the existing Kinet consumer APIs (Ramen checks excluded). M89 acceptance cannot activate the original
C11/M54 conditional trigger.

**Downstream.** Kinet:M55.

**Tasks/review.** [status-M86.md](status-M86.md), 3 rows, review 0/10. The
[coordinator](../../../kinet/docs/stage12.md) holds the shared order and triggers.

M83/M84/M85 and M87–M89 remain accepted/retired. S07 is accepted/retired
after whole review 2; M86 remains conditional with accepted M89/S07
compatibility prerequisites.
Original Stage 11 is accepted across the five owner ledgers.
Completed source repair and version discovery remain in the history index.
Other candidates still require fresh reconciliation and approval.

**Depends on.** UWS:C11, APItools:M89 and APItools:S07. All required prerequisites must have
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

## Requested Changes After Initialization

For stage work, `memory-bank-propose --stages` requests creation or revision of
the optional `tabilet/stages.md` overview. `--stage STG-01` selects an existing
stage; the request must say when milestones should be planned. Equivalent
plain-English references are valid. Later-stage ideas may have only an ID,
tentative name, and brief intent. Additional context can be proposed later,
without assigning milestone IDs or task rows now. Stage IDs stay stable across
renames and reorderings and are never reused. Do not advance stages or schedule
later work merely because the current milestone closes.

For execution environment and authorization policy, `--environment dangerous`,
`--environment sandbox` (default when omitted), or `--environment least` sets the
hands-free stance in launch references and milestone policies, optionally modified
by expressions like `plus <action>`, `minus <action>`, or `but <action> in milestone <ID>`.

For a requested feature, candidate promotion, or change to future direction,
inspect the current plan and implementation first. Record user intent separately
from observed facts and unresolved assumptions. If a pending milestone already
owns a small change, propose the smallest row and acceptance update that fits.
Otherwise propose a dependency-closed milestone with an unused permanent ID.
A candidate's trigger prompts a new decision and approval; promotion is never
automatic. Schedule by approved priority and dependencies, not review finding
severity. Reconcile affected downstream acceptance and optional launch input.

Present one complete proposal with intended outcome, owner rows, acceptance,
verification, dependencies, downstream effects, and exact file actions. Write
planning changes only after approval and a fresh check of affected files,
worktree changes, and active and retired IDs. Preserve non-pending outcomes,
review counters, local policy, and frozen history. Record target behavior here
and in pending status rows until implemented; current architecture describes
only what the repository establishes as fact. Planning does not implement or
accept the requested behavior. Use `memory-bank-propose` where installed, or
follow this procedure directly with another agent.

**Expensive checks.** Some acceptance depends on a check that is expensive or
protected; see the [verification-attempt gate](#verification-attempt-gate).
For such a milestone:

- Limit acceptance to behavior the milestone owns, proved with the smallest
  workload. Put actual pairing with real producers or consumers in the
  milestone that adopts them.
- When the check needs host capabilities (privileges, operating-system or
  kernel features, devices, browsers, or external accounts), make a cheap
  readiness probe its first row.
- Name the uncounted local diagnostic loop, and order cases so that an
  unrelated failure cannot hide the milestone's own cases.

To defer already approved work through a stage rescope, present the exact
untouched pending rows and affected downstream work for approval. After
approval, retain their identity and original task text, mark them `[X]`, and
record the authority, reason, destination stage, and dependency disposition.
Cancellation proves no delivered acceptance. A future replacement receives a
new approved identity and names its cancelled lineage. Preserve `[~]`, `[!]`,
`[+]`, `[X]`, and `[-]` rows and all review evidence. Withdraw affected pending
dependents or keep their prerequisites in the active horizon; an unresolved
dependency cannot be deferred by changing a stage label. An entirely cancelled
milestone remains active until its normal closure and retirement disposition.

## Candidate Directions

Candidate directions are outside the active execution horizon. They have no
lane, permanent ID, status file, or execution-order entry. A trigger starts a
fresh scope/dependency reconciliation and approval round; it does not schedule
the work automatically.

| Direction | Why Deferred | Promotion Trigger |
|---|---|---|
| Automatic catalog-root discovery | M87 keeps explicit roots and metadata-only no-root outcomes; no implicit cache lookup/creation. | A named consumer and approved compatibility/privacy review justify implicit discovery. |
| APItools-owned HTTP model adapters | M89 uses consumer-owned model clients; standalone APItools remains deterministic. | A standalone adopter approves provider transport, opt-in, credential/privacy and offline-test contracts. |
| Next provider-catalog expansion | No provider batch or user-verifiable outcome has been approved after the completed historical catalog batches. The [catalog upgrade proposal](../../docs/catalog-upgrade-2026-09.md) also proposes a structured license/redistribution audit and workflow-ready typed overlays; those remain separate from M80's discovery report, which must expose unknown evidence. GoDaddy is absent from the catalog; its official Domains specs (`https://developer.godaddy.com/openapi/domains-v1.json` through `domains-v3.json`) were observed on 2026-10-01, and curating v1 or v2 would need a source-reviewed refresh correction for a publisher `properties.type` string defect. | Approve a bounded provider-owned source batch, its source/auth/license review rules, or a bounded license/typed-overlay upgrade with acceptance evidence. Proposed catalog IDs remain unallocated until fresh reconciliation and approval. |
| Embedded baseline catalog root | M81's approved decision 3 keeps catalog roots caller-supplied. Embedding the 141 committed advisory overlays (about 2 MB) and their index would give published consumers a zero-setup offline root, but it widens distribution before the committed-overlay redistribution audit proposed in the [catalog upgrade proposal](../../docs/catalog-upgrade-2026-09.md) and adds binary size for every consumer. | The committed-overlay redistribution audit completes, and a consumer needs a zero-setup offline root. |
| Search default-behavior refinements | `Client.Search` keeps APIs.guru's `preferred` version and falls back to a lexicographic sort of version keys when it is missing (`providers.go`), and `auto` stops at the first source with matches (`client.go`). Changing either default alters results that OpenUdon and other callers rely on, so S05 reports versions through a separate opt-in command instead. Separately, the retired M80.4 remote tier sets `client.Cache = nil` with `CacheModeBypass` (`catalog_discovery_remote.go`), so each remote lookup re-downloads the whole APIs.guru directory list; S05 avoids that cost for its own checks but does not change the retired behavior. | Approve a default-behavior change with consumer reconciliation, observe APIs.guru entries whose `preferred` version is missing or unusable, or measure the directory download as a latency problem for a consumer. |
| Identify APItools on guarded downloads | No `User-Agent` is set on any guarded request, so hosts see Go's default agent. S05's new probe primitive sends a stable tool identifier, but adding one to the shared download path changes the requests every existing caller makes, including OpenUdon's pinned use. | A publisher rejects the default agent, an operator requires identified egress, S05's probe evidence shows the identifier is worth sharing, or the provider version directory's crawl policy requires it; then approve it with consumer reconciliation. |
| Provider version directory | A shared, published directory that grows with demand: one record per provider and service with its locator, recipe, resolved URL, digest, `checked_at`, evidence label, and the age of any APIs.guru baseline, in an APIs.guru-`list.json`-compatible shape with `x-apitools-*` extensions and `preferred` naming the highest verified comparable version within the recorded checked scope (no global latest claim or new preference from conflicting evidence). Seeds are APIs.guru's list (CC0 for its contributors' work; frozen since April 2023, with `x-origin` on every entry) and the public-apis README (MIT; about 1,970 docs links, 20% overlapping APIs.guru), plus the curated catalog; a monthly deterministic job (conditional GETs and one GitHub tree listing per repository; 2,632 of APIs.guru's origins sit in 154 GitHub repositories, the other 1,381 spread over 769 hosts) re-verifies records, and LLM discovery runs only on demand per provider. It holds metadata only, never spec bytes, and is published as a separate artifact, not inside this module or on the code repository's `main`. User submissions are untrusted leads until the job verifies them. Considered and rejected: cloning APIs.guru (the clone has no `list.json`, and the data is frozen), mirroring spec bytes (only 27% of APIs.guru specs declare a license, and the module zip is capped at 500 MiB), and daily pushes of data into the code repository (consumers pin exact revisions). S05 delivered the per-provider record this would aggregate; its legal, crawl-policy and operator gates remain separate. | S05 is delivered and its record shape is stable, legal sign-off covers publication and a crawl policy (identifying User-Agent, rate limits, robots, takedown handling, GitHub API terms), and an operator wants a shared directory; then approve a bounded milestone with its own gating row for that sign-off. |
| Production adoption of experimental LAP/RFC 9727 discovery | The adapters are intentionally experimental, and OpenUdon still selects APIs.guru explicitly. | Collect reliability/coverage evidence and approve the downstream OpenUdon discovery policy and network budget. |
| Remove deprecated Discovery and Smithy wrappers | Sibling consumers may still rely on compatibility imports, and removal is a public breaking change. | Confirm no active sibling imports remain, approve release notes, and pass downstream compatibility checks. |
| Network and local-scan hardening batch | Review "apitools Code Review" Pass 1 (baseline `a5699e5`) findings #7 (LAP registry validates every entry by DNS before scoring/limit and fails the whole search on one bad `source_url`, `remote_discovery.go`), #9 (unsafe-IP list omits `240.0.0.0/4`, IPv4-compatible `::/96`, and `fec0::/10`, `download.go`), #11 (OpenAPI metadata paths skip the `sourceguard.CheckJSON`/`CheckYAML` preflight, `validation.go`, `local.go`, `local_source_discovery.go`), #12 (`.bin` is always treated as a protobuf descriptor, `local_source_discovery.go`), #13 (regular-file-to-FIFO swap between `Lstat` and `Open` can block, `local_read.go`, `internal/artifactio`), and #14 (`Discoverer` exposes no port/transport options, `discovery.go`). All were revalidated at `a5699e570fb5d6288fdd255c71d229fb062b50b3` as Lower severity; no supported scenario fails and no private-host or root-escape path was found. | Approve a hardening batch after S02, or promote any item earlier if it becomes reachable in a supported deployment or a downstream consumer depends on it. |
| Egress proxy support for guarded downloads | Review "apitools Code Review" Pass 1 finding #8: the guarded transport has `Proxy: nil`, and the only alternative, `AllowUnsafeHosts`, drops all host checks (`download.go`). Supporting a proxy trades dial-time IP filtering for proxy trust, which is a product/security decision. | An operator needs mandatory-proxy egress, and an approved policy defines target pre-validation, proxy trust, and the documented DNS-rebinding limitation. |
| Lifecycle and scan cleanup follow-ups | Uncommitted-diff review (2026-09-26) Lower items not required by S02, S03, or S04 acceptance: `operationlifecycle` recomputes tokens and purpose for every operation on each role (performance only); offline `search` still applies the cache TTL while offline `import` ignores it (consistency decision); `isUnicodeTagCharacter` largely repeats the Cf category check in `prompt_safety.go` (cosmetic). | Promote when lifecycle ranking becomes a measured hotspot, an operator needs offline search of expired entries, or `prompt_safety.go` is next edited. |
| Catalog and tooling hardening batch | Review "apitools Code Review" Passes 2 and 4 (revalidated at `e015231`), all Lower severity: C3 `ResolveProviders(query)` splits multi-word display names on whitespace (`catalog/materialize.go`; CLI and OpenUdon pass explicit keys); C4 `FindSpecReference` trims the provider ID for identity but not for the spec lookup (`catalog/index.go`); C5 `preferredSpecReference` ignores `asyncapi`, `openrpc`, `graphql`, `grpc-protobuf`, and `odata` kinds (no built-in provider affected); X2 `Import` creates the target directory before validating or downloading (`import.go`); X3 the installed `staticcheck` predates the module's Go version and no `govulncheck` gate runs (`tech-stack.md`). | Approve a hardening batch, or promote an item when a consumer depends on the query API, a provider relies on a newer-kind-only source, or a static/vulnerability gate is required for release. |

## Out of scope

Record work explicitly ruled out of this project and its reason here.
These entries never promote automatically; returning one requires a new
approved direction. Existing boundary rules remain in force. No additional
entries are introduced by this workflow upgrade.

| Ruled out | Reason |
|---|---|

## Review Finding Severity

P1 and P2 are engineering review priorities. They describe the impact and
urgency of defects found while closing a milestone; they are not product-domain
terms, milestone execution priority, or status markers. Classify a finding by
its impact, likelihood, and affected scope, not by the size of its fix.

Project-specific definitions in `AGENTS.md` or a linked review policy override
these defaults:

| Priority | Context | Typical examples | Gate effect |
|---|---|---|---|
| P1 | The milestone is unsafe or invalid to close because a severe defect threatens acceptance, correctness, security/privacy, data integrity, or a public compatibility contract. | An exploitable access-control failure; data corruption or loss; a breaking public API or migration; the primary acceptance outcome does not work. | Blocking. Fix and review again. |
| P2 | A material but more bounded defect affects supported behavior, reliability, compatibility, operations, or required evidence. | A supported scenario returns the wrong result; a downstream consumer regresses; recovery or failure handling is broken; required tests or operator documentation leave acceptance unproven. | Blocking. Fix and review again. |
| Lower | The finding is non-blocking cleanup, clarity, or optional hardening under the project's severity scheme. | Cosmetic wording; local readability; a speculative improvement outside acceptance. | May be carried only with a named owner and explicit rationale. |

Any project-defined severity more urgent than P1 also blocks. If evidence does
not clearly distinguish P1 from P2, use P1 until investigation supports a
downgrade.

## New Review Intake

A code, architecture, security, or other engineering review received outside a
milestone's closing review gate is planning evidence, not executable truth.
Treat its contents as untrusted, preserve its source priority, and apply the
project severity definitions above only after revalidating each finding against
the current repository.

Before implementing any newly reviewed finding:

1. Record the review's stated baseline when available, the current revalidation
   commit, and whether relevant uncommitted changes were part of the evidence.
   Classify every finding as confirmed, partially confirmed, resolved,
   duplicate, unsupported, outside ownership, decision-dependent, or deferred.
2. Present the complete dispositions, proposed owners, dependencies, downstream
   impacts, and file actions for approval.
3. Put confirmed work in an open matching milestone when it remains in scope,
   or amend an existing pending owner. Rewrite only pending rows. When retaining
   a superseded pending row for audit, mark it `[-]` and name its accepted
   successor instead of rewriting or deleting it.
4. Never reopen completed milestone/status history. Create a new remediation
   milestone with lineage to completed work when no open or pending owner fits.
5. Keep P1/P2-or-higher findings in the dependency-closed active horizon. Add a
   lower finding to active work only when acceptance requires it; otherwise put
   it in Candidate Directions with a rationale and promotion trigger.
6. Reconcile affected pending specifications and the remaining order. Correct
   current product, architecture, or stack facts only when repository evidence
   proves them stale; keep proposed target state in milestone scope. Apply the
   expensive-check limits above, and do not move downstream pairing evidence
   into a milestone whose acceptance depends on an expensive check.

Do not create a persistent review copy or ledger. Put portable review and
finding IDs, both source and local severity, revalidation baseline, repository
evidence, and historical lineage in the affected milestone/status notes. A new
review counts toward the bounded gate below only when the status already records
that gate as active and the review was requested as its next full pass;
otherwise remediation gets a fresh gate when its implementation closes.

## Verification-attempt gate

An expensive or protected check needs privileged setup, scarce hardware or
quota, a physical or live environment, explicitly granted actions, or minutes
per run. For a row whose acceptance depends on one:

- Reproduce only when safely possible within current authority and budget.
  Prefer the cheapest environment that can show the failure, such as an
  unprivileged local run on disposable fixtures. For unavailable environments,
  protected actions, exhausted attempt limits or non-reproducible incidents,
  record existing evidence and the reproduction limitation instead of forcing
  retries or inventing a failure. Required acceptance evidence and verification
  gates still apply; if missing, pause. Diagnostic runs are ordinary verification
  and never satisfy acceptance.
- Make each expensive run execute every case whose safety does not depend on an
  earlier result, and report all failures. Stop early only when continuing
  would be unsafe.
- Redaction protects real data. Raw diagnostics from synthetic, disposable
  fixtures may be kept in private, ignored artifact paths. Never commit them,
  and never treat them as evidence of effects, non-dispatch, retry eligibility,
  or authority.
- Each expensive run is one attempt. Persist the row's attempt count in its
  status notes; a session change does not reset it. Stop after three
  consecutive failed attempts at the same stage (or the limit the milestone
  sets), or when the failing stage lies outside the behavior the milestone owns.
  Mark the row `[!]` with a short diagnosis and options, such as a correction, a
  rescope or an environment change, and ask the user. Approval to continue sets
  a new explicit limit. Never repeat an unchanged failed attempt.

## Planning Notes

- Public OpenAPI behavior remains in `../apitools`.
- Harness state (`AGENTS.md`, `tabilet/`) is tracked in `../apitools`; the
  former private `../tofu/apitools` harness is no longer used.
- Batch expansion work should freeze service lists before implementation and
  treat each service row as a commit unit once curation begins.
- Do not promote non-OpenAPI converter work until most provider-node catalog
  coverage is implemented.
- Google Discovery converter ownership is now explicitly opened in M38 because
  provider-node catalog coverage through M37 is complete.
- AWS Smithy converter ownership is complete in M39 for official AWS Smithy
  JSON service models; SigV4 signing remains downstream runtime work.
- AWS Smithy native model ownership is complete in M40; OpenAPI-shaped Smithy
  conversion has been removed from `apitools`, and downstream runtimes consume
  native Smithy metadata.
- Google Discovery native model ownership is complete in M41; OpenAPI-shaped
  Discovery conversion has been removed from `apitools`, and downstream
  runtimes consume native Discovery metadata.
- Run `apitools catalog check` before completing each catalog metadata batch.
- M52 supersedes the prior n8n-derived catalog planning workflow. Active
  catalog expansion is provider-source-first and must not copy third-party
  integration code, generated node metadata, credential behavior, descriptions,
  icons, or runtime workflow behavior.
- First-class Smithy and Discovery work means cataloged native source metadata
  and parser coverage, not reintroducing Smithy-to-OpenAPI or
  Discovery-to-OpenAPI conversion as an `apitools` runtime contract.
- M47 retained Dropbox Stone as official source metadata and advisory-overlay
  provenance, but did not open native Stone parser work or expand the M37
  advisory overlay.
- M48 closes the residual Google Discovery batch for n8n-visible Google
  services that were not promoted in M44, and records Google Ads as a
  candidate-only/no-Discovery decision until a future protobuf/gRPC source
  family is explicitly scoped.
- M49-M51 are historical and superseded by M52 for active code and public docs.
  The n8n provider row, n8n gap-report CLI/API, and n8n-derived priority
  evidence are intentionally removed before release.
- M53 should add Docker source coverage from Docker-owned API references:
  Docker Hub as a normal public OpenAPI provider, Docker Engine as a local
  daemon/control-plane metadata source, and Docker Registry only if source
  review confirms a durable official spec path.
- M54 treats Kubernetes as cluster-published source metadata. A pinned
  v1.19.2 Swagger fixture snapshot is available as a non-official bootstrap
  artifact, but real Kubernetes OpenAPI and Discovery documents remain
  user-provided or cluster-exported artifacts; `apitools` must not read
  kubeconfig, resolve cluster credentials, contact API servers, or claim
  live-cluster compatibility.
- M55 should add enterprise application source coverage for NetSuite, SAP
  S/4HANA, SAP SuccessFactors, Oracle Fusion Cloud Applications, and Workday.
  Treat account/tenant-specific metadata endpoints as user-provided source
  URLs, and keep OData, WSDL/SOAP, and describe-metadata decisions explicit.
- M56 should add an OpenAPI-first provider batch from provider-owned official
  specs: Adyen, DocuSign, Auth0, Confluence, Acumatica, BigCommerce, Cisco
  Meraki, and Confluent Cloud. Use third-party integration catalogs only as
  market-demand signals; do not copy their connector metadata.
- M57 should add a docs-overlay provider batch for Marketo, Aircall, Checkr,
  Airwallex, Apollo.io, AfterShip, and Adobe Acrobat Sign only where official
  downloadable OpenAPI evidence is unavailable or unsuitable and official docs
  can support a reviewed subset.
- M58 is deferred until M57 is complete. It should audit all durable catalog
  providers for auth/security metadata completeness under the OpenAPI/Swagger
  security-scheme and root/operation security requirement rule, plus native
  Smithy, Discovery, Stone, human-docs, and tenant-specific source-family
  decisions.
- M59 should use `tray_workato_connectors.xlsx` only as a market-demand
  signal, then add Canva, LaunchDarkly, Klaviyo, Attio, Fivetran, Anthropic,
  Amplitude, Ashby, Braze, and Postman from provider-owned API evidence.
- M60 should expand AWS Smithy service coverage for Workato-visible AWS
  candidates using official AWS Smithy JSON models only; SigV4 remains
  advisory metadata and no AWS credential, region, or account behavior belongs
  in `apitools`.
- Deprecated `apitools/googlediscovery` and `apitools/awssmithy` wrappers can
  be removed only after active sibling consumers import
  `github.com/OpenUdon/googlediscovery` and `github.com/OpenUdon/awssmithy`
  directly, public release notes call out the compatibility break, and
  dependent checks in `../openudon` and `../udon` pass without wrapper imports.
- The next catalog-expansion milestone should start from provider-owned source
  evidence and prioritize remaining high-demand services with official
  OpenAPI/Swagger, Smithy, Discovery, Stone, official indexes, or official docs
  sufficient for reviewed advisory overlays.
- Post-M60 materialization work moves static provider artifact handoff into
  `apitools`: resolving provider names, reporting protocol capability, copying
  existing cache-registered artifacts, emitting separate security-overlay JSON,
  and exporting provenance manifests for workflow directories. It remains
  offline and must not download missing artifacts, apply overlays into OpenAPI
  documents, or reintroduce Discovery/Smithy-to-OpenAPI lowering.
- M61 reviews existing docs-derived advisory overlays under the expanded
  linked-docs curation rule: follow official API links cited by source docs,
  include related operations when they support the same provider workflow, and
  record no-add decisions where linked docs are conceptual, commercial,
  duplicated, authenticated-only, or outside the reviewed scope.
- M62 opens a narrow public helper surface for pure fnct payload shaping. The
  first helper is Gmail raw-message rendering; runtime registration remains in
  `../udon`, and OpenUdon remains the authoring/review consumer.
- M63 adds OpenRPC as metadata-only JSON-RPC source support: the `openrpc`
  package parses document metadata, method summaries, selector aliases, params,
  results, errors, and components; catalog/source materialization treats
  OpenRPC as its own `openrpc` source family; JSON-RPC execution, server
  contact, credentials, and runtime transport remain out of scope.
- M64 closes Postman Collection, RAML, and API Blueprint as local artifact
  detection/classification and conversion candidates only, following the
  Dropbox Stone decision pattern. Preserve provenance, prefer official
  OpenAPI/Swagger, Discovery, Smithy, AsyncAPI, OpenRPC, or another stronger
  native source when available, and do not add parser code, runtime behavior, or
  UWS `sourceDescriptions[].type` values for these formats.
- UWS 1.4 now defines `graphql`, `openrpc`, `grpc-protobuf`, and `odata` source
  description types. OpenRPC parser/classification work is already complete in
  M63. M65-M67 open separate metadata-only source tooling milestones for
  GraphQL, gRPC/protobuf, and OData. Each must preserve native source semantics
  and must not execute operations, probe live services, resolve credentials, or
  choose accounts/tenants/projects.
- M65 adds GraphQL as metadata-only source tooling: local SDL, introspection
  JSON, and operation documents parse into native GraphQL summaries and
  selectors; catalog refresh, resolution, and materialization identify
  GraphQL as `graphql` without live introspection or REST-shaped lowering.
- M66 adds gRPC/protobuf as metadata-only source tooling: local `.proto`, JSON
  descriptor, and binary `FileDescriptorSet` artifacts parse into native
  protobuf service/method/message summaries and selectors; catalog refresh,
  resolution, and materialization identify gRPC/protobuf as `grpc-protobuf`
  without server reflection, channel behavior, TLS negotiation, RPC execution,
  or REST-shaped lowering.
- M67 adds OData as metadata-only source tooling: local CSDL XML, JSON CSDL,
  and exported service metadata parse into native entity set, singleton,
  navigation, action, function, parameter, query-option, and selector
  summaries; catalog refresh, resolution, and materialization identify OData
  as `odata` without tenant metadata calls, login, credential resolution,
  account/tenant selection, `$batch`, request execution, or REST-shaped
  lowering.
- M68 adds a reviewed local artifact corpus for UWS 1.4 source families:
  exactly two official example artifacts each for GraphQL, OpenRPC,
  gRPC/protobuf, and OData under source-aligned cache directories. It is
  curation and validation coverage only; it does not add provider catalog rows,
  source types, OpenAPI lowering, parser behavior expansion, or runtime
  behavior.
- M70 keeps APIs.guru as the primary global lookup, evaluates LAP Registry as
  an experimental second source, and adds provider-scoped RFC 9727 discovery.
  Swagger Catalog and Scalar Registry remain outside implicit global discovery
  because they are publisher/team registries rather than unauthenticated
  global aggregators.

## Milestone Review Procedure

When the last open row in a milestone's status file closes as `[+]`, `[X]`, or
`[-]` during an agent session, perform the review before ending the turn and
before moving to the next milestone. A `[-]` row counts as closed only when its
notes identify the consumed attempt or supersession and its accepted successor:

1. Re-read the milestone scope and acceptance criteria here. Confirm the code or
   docs meet the acceptance line; do not rely on the status file alone.
2. Run the required verification before review, including proportionate checks
   for affected consumers and any applicable compatibility, migration,
   rollback, security, concurrency, or failure paths.
3. Run a bounded deep-review and fix gate using the review finding severity
   context above.
   - The initial deep-review pass is iteration 1. Read the persisted count and
     findings first. Record the iteration as started before review; resume an
     interrupted pass at that same number. Runtime round limits are separate.
     Pin the milestone's base commit, inspect the full diff including staged,
     unstaged and relevant untracked work under the commit policy, read the
     `git log` range, and review the full milestone diff for correctness,
     regressions, failure semantics, boundary drift, stale docs, and missing
     tests. Review acceptance fidelity separately: requirements missing or
     partial, behavior nobody asked for, and requirements implemented wrongly,
     each quoting its acceptance line. Record findings grouped by concern.
     An empty implementation diff still needs acceptance review; do not
     manufacture changes. Lenses are coverage obligations: combine them or
     run batches within host capacity. When an interactive request selects
     `MODEL_POLICY`, follow the Goal/Next bundled model-policy reference for
     fresh read-only whole-milestone reviewers; the owner alone records results
     and fixes under `DEFAULT`. Absence preserves host model settings.
   - Record the iteration number and findings in the current status notes so a
     continuation cannot reset the counter. If the pass finds no P1, P2, or
     higher-severity issue, the gate passes. Otherwise, when the current
     iteration is below 10, fix every such finding in the current milestone,
     rerun affected verification, and review the whole milestone again,
     including the fixes.
   - Run at most 10 iterations; a session or reviewer change does not reset the
     counter. If iteration 10 still finds a blocking issue, add `[!]` review
     rows recording the remaining findings and iteration-limit blocker. Do not
     start another automatic fix-review cycle or move downstream; stop and ask
     for user direction.
   - Carry a lower-severity finding forward only with a named pending owner and
     explicit rationale.
4. Reconcile the memory bank. Update `product.md` if the milestone changed
   product scope, domain terminology, concept relationships, or business
   invariants. Update `architecture.md` or `tech-stack.md` if it changed
   boundaries, dependencies, commands, data flow, or runtime assumptions.
   Then run an environment retro on mistakes found during the milestone:
   mechanical issues become proposed deterministic checks; judgement issues
   become lessons or review rules. Include missing navigation or information and
   instructions that change nothing. Report these candidates by severity for a
   later approved proposal; create no unapproved task rows.
5. Check `tabilet/evolution/`. Record an explicit bump or no-bump decision. Add
   the next `prompt-vN.md` and `result-vN.md` only when product direction,
   architecture boundary, milestone target, or public/private contract direction
   materially changes.
6. Revisit candidate directions affected by the milestone. Update their reason
   or trigger; when a trigger is now true, propose a reconciled milestone and
   obtain approval before allocating its permanent ID and status file.
7. After the gate passes, reconcile affected active downstream dependencies and
   acceptance against the delivered outcome. Then consolidate and retire under
   the procedure below. During an ordered goal, the goal protocol finishes its
   downstream reconciliation before performing this retirement step.
8. Run required verification again, including the retirement links and record,
   then commit substantive closure changes under the governing commit policy.
   Do not create an extra milestone commit when the review changes nothing.
9. Report a short review summary: what was verified, the review-fix iteration
   count, what memory-bank files changed, any review commit, and whether an
   evolution bump was made. Include the retired record and durable lessons.

## Goal lease cleanup

When a goal uses external lease worktrees, follow the selected goal protocol's
cleanup gate. Remove one only after recorded successful closure, matching
ownership, no active writer, verified integration and a clean worktree. Inspect
ignored data as well as tracked changes. Before freezing closure records,
preserve required qualification/recovery artifacts in an approved durable
location with their hashes and citations. Never rewrite frozen citations or
delete their required evidence; preserve and report uncertain leftovers.
This rule does not authorize cleanup of existing leases during an upgrade.

## Long-term Memory And Retirement

This project retires a milestone automatically after closure. “Automatically”
means the agent performs this procedure during the milestone review workflow
above, including when using `memory-bank-next` or `memory-bank-goal`. There is
no background process or separate archive invocation. Individual completed,
cancelled, or closed-historical rows remain in their active status file until
the whole milestone qualifies; unresolved work or missing closure evidence
keeps it active. Terminal rows alone do not prove acceptance.

Existing projects adopt this contract explicitly, with a compatible API runner
if used; upgrading installed skills alone does not merge project instructions
or move files. An explicit cleanup request may retire older closed milestones
only when their closure evidence is available, either a recorded review under
this contract or the legacy closure described below.

### Consolidate Before Retiring

Keep current facts in product, architecture, and stack documents. Maintain
[lessons.md](lessons.md) for applicable lessons and decision rationale with
evidence links; merge duplicates and avoid a chronological session log. Maintain
relevant lessons during ordinary work, not only at closure, and keep them active
while applicable even after their supporting milestone retires. This removes
accumulated history from active context, not a fixed number of tokens or files;
genuinely active work and useful knowledge can still grow.

Before materially superseding or removing knowledge from those documents,
append the old wording to `tabilet/docs/history/knowledge.md`. Use a unique, descriptive
dated heading, the original document and heading, retirement reason, supporting
evidence, and a link to its replacement (or an explicit reason there is none).
Preserve the old excerpt in a fenced markdown block. This applies outside
milestones too. Create this journal only when needed, link it from the history
index, and append corrections rather than rewriting old entries. Routine edits
need no journal entry; Git, when present, holds intermediate revisions.

Consolidation is ordinary project maintenance. A separate context snapshot is
optional and never a prerequisite for retirement; the archive skill's clean
baseline rule must not force commits or interrupt a no-commit goal.

### Retire A Closed Milestone

1. Require a passed review within the persisted 10-iteration limit, recorded
   verification, consolidated knowledge, and reconciled downstream work. No
   `[ ]`, `[~]`, or `[!]` row may remain. Every `[-]` row names its accepted
   successor. A cancelled or superseded milestone also needs an authorized
   disposition; it must not be described as delivered acceptance.
2. Create `tabilet/docs/history/status-<LANE><NN>.md` using the envelope below. Keep the
   full final specification and status document in separate literal markdown
   fences, including all task text, notes, acceptance, and review evidence.
   Choose fences longer than any fence in the source. Original relative paths
   inside these literal documents retain their source-document meaning.
3. Add one history index row using `Milestone | Outcome | Retired | Record |
   Summary`; the ID, outcome, date, and record link must match the envelope.
   Outcomes are `completed`, `cancelled`, or `superseded`. The record link is
   relative to the index, for example `[M01](status-M01.md)`.
4. Remove the active status file, specification section, and active index row.
   Before removal, validate every envelope field and compare both retained
   documents with their complete active sources. For Git evidence, use the
   literal full output of `git rev-parse --verify HEAD`, never an abbreviated
   hash from a log display. Stop with the active sources intact if validation
   fails; a malformed retired record is not completed closure.
   Keep one link from this file to `../docs/history/index.md`. Repair maintained
   incoming links and remaining dependencies. Never rewrite frozen snapshots;
   resolve their old paths through the record's original-path metadata.
5. Refresh an existing disposable goal suggestion to contain only remaining
   active work, or remove it when empty. Do not create one without a compatible
   goal protocol. The resolved goal in the conversation remains authoritative
   over its disposable launch input.
6. Verify the record and links before handoff. On interruption, reconcile the
   source and destination before continuing; duplicate IDs or a missing record
   are invalid, and an existing retired record must never be overwritten.

Use this envelope, with actual values. The `Evidence` field is the observed
full Git commit, not a claim that uncommitted work exists in that commit.
`Worktree` is `clean`, `includes uncommitted changes`, or `unversioned`; without
Git both provenance fields say `unversioned`. The date uses `YYYY-MM-DD` in UTC.
Closure metadata precedes the two sections and uses one line per field.

`````markdown
# Retired milestone M01 - <title>

**Milestone.** M01
**Outcome.** completed
**Retired.** <YYYY-MM-DD>
**Source status.** tabilet/memory-bank/status-M01.md
**Source specification.** tabilet/memory-bank/milestone.md#<original-anchor>
**Evidence.** <full commit or unversioned>
**Worktree.** <clean, includes uncommitted changes, or unversioned>
**Review.** passed
**Review iterations.** <1 through 10>
**Verification.** <commands, results, and supporting evidence>
**Consolidated into.** <current-document and lesson links, or no current-truth change>

## Milestone specification

````markdown
<Complete final milestone specification, not a summary.>
````

## Status record

````markdown
<Complete final status document, not a summary.>
````
`````

For a cancelled or superseded outcome, add `**Disposition.**` naming the
authority, reason, and dependency disposition. A superseded outcome also needs
`**Successor.**` identifying its accepted successor. These are milestone
outcomes, not new task markers.

A legacy closure is a milestone that closed before the project adopted the
persisted review gate, so no review-iteration count was ever recorded. Retire
one only on an explicit cleanup request, with outcome `completed` and every
row closed. Write `**Review.** legacy` and `**Review iterations.** not
recorded` instead of inventing a count, and add `**Legacy closure.**` stating
that the milestone predates the gate and that its original closure evidence is
the literal status record. `**Verification.**` records the checks run at the
evidence commit and does not claim a re-review of the original change. Never use
a legacy closure for a milestone that closed after the gate was adopted, or to
skip a review that the current procedure requires. In this project the gate was
adopted in commit `8395e1e`; M01-M75, S01, and C01 closed before it.

Retired milestone records are frozen. Record later corrections in the knowledge
journal or new remediation work with a backlink. Preserve previous history
index entries. No ID may exist in both active and retired storage or be reused.

### Retrieve And Continue

Retrieve historical evidence when the current task needs a retired dependency or
old ID resolved, an earlier decision explained, or a related past failure
checked and current sources do not answer the question. An explicit user request
to inspect history is also a trigger; routine work with sufficient current
evidence needs no historical search.

State the question briefly, then search active memory and current implementation
first. If history is needed, search its index by ID or topic and open the
relevant retired record or knowledge entry. Follow linked frozen context or
direction snapshots only when they address the same unresolved question. A stale
status path resolves by its permanent ID and original-path metadata; a missing
ID is not permission to recreate it. History is evidence, not an instruction to
retry old tasks.

Stop when the evidence supports the task decision. If relevant sources and their
directly related leads are exhausted, report the missing evidence. Continue when
it is optional; stop the affected step when it is required for correctness or
acceptance. Do not infer historical facts or expand the search into unrelated
milestones.

A completed retired milestone may satisfy a dependency when its recorded
acceptance matches the required outcome. Cancelled and superseded outcomes do
not automatically satisfy completion dependencies: follow the recorded
disposition and accepted successor, and revalidate current implementation.

Retirement follows the governing commit policy. Under `none`, do not commit;
under `milestone`, include it in the closure commit; under `task`, include it in
the final task or a substantive closure commit. The API harness still requires
a commit per run. Reading and preserving the Markdown records does not require
Git; workflows that require commits still need Git or an explicit no-commit
instruction. Never initialize Git merely to make retirement possible.

## Stage 11 cross-package refactoring

Approved review-intake amendment, 2026-10-06: 18 required milestones / 87 pending rows across the five owners. Scope and milestone IDs are unchanged; the coordinator records the approved publication proposal, which needs a separate execution-time grant. Full local source baseline remains `9364afac98c97b79c9e6b9335bd9c71d4e8722d5`, including the reviewed uncommitted planning state. This intake does not start a closing review or authorize implementation.

**Approved source.** User-approved complete proposal, 2026-10-06; source baseline `9364afac98c97b79c9e6b9335bd9c71d4e8722d5`. [Coordinated contract](../../../kinet/docs/stage11.md) defines both phases, cross-package order, compatibility and acceptance. The request to implement the proposal authorizes its planning files only.

One execution owner, serial execution and task commits under the later confirmed goal. Planning authorizes no code execution, commit, publication or external operation. Source publication requires separately named authority; a status marker or local build is not publication. Consumers must record exact accepted and published prerequisites before adoption. Default checks are offline, credential-free and model-free. No deployment, live ledger migration, real API/model/mail action or registration change.

## Stage 11 candidate dispositions

M82 is new source-tooling work with a shared public UWS contract, so it uses the cross-cutting M lane. It promotes no catalog/provider expansion, remote discovery, credential/runtime behavior or unrelated hardening candidate.
