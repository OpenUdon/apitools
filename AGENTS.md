# AGENTS.md

## Stage 11 source metadata

M84 remains frozen. [M85](tabilet/docs/history/status-M85.md) is accepted/retired
after four rows/final whole3, published source409d and independent ordinary
candidate+SDK/CLI/source proof. Effective reserved-header inputs are coherent,
with existing reference limits and no new credentials/runtime capability.
Kinet:M50 must adopt the exact source and requalify its workers/bundle; no live
or deployment authority follows. The named M85 source operation is consumed.

[Stage 11](../kinet/docs/stage11.md) coordinates both refactoring phases across
five owner ledgers. APItools M82 is accepted and published at
`54583f9b2f452b7cc522360c5aeeff29ca22f96c`, whole review 3. It produces and
independently verifies source-neutral UWS operation shapes from eight existing
source families with explicit incompleteness and bounded local parsing. It
performs no credential resolution, workflow approval or runtime execution.
Kinet/OpenUdon adoption remains separately owned, against exact accepted/public
contracts inside isolated parsing paths. Preserve frozen evidence and consumer
pins. One serial execution owner and offline fixtures remain required; source
publication/deployment authority stays separate from planning/status markers.

Requested features, candidate promotions, and future direction changes after
initialization follow the requested-change procedure in
[tabilet/memory-bank/milestone.md](tabilet/memory-bank/milestone.md). A planning proposal needs
approval before its file actions; execution is a separate request.

## Stage 12 concurrent goal ownership

For an explicitly authorized STG-12 goal, the governing protocol is
[../kinet/tabilet/GOAL.md](../kinet/tabilet/GOAL.md). The coordinator passes that exact protocol and the
complete resolved request to each child; sibling standalone goal protocols
remain unchanged. Use the goal skill through a hosting agent with subagents;
the Python API runner/controller remains serial.

This is a scoped exception to serial and whole-ledger single-row wording:

- One coordinator owns all integrated package-local ledgers and shared memory.
  At most three milestone leases may run, with at most one live milestone per
  package and one `[~]` row per lease. Outside this opt-in, serial rules apply.
- Dispatch requires explicit `PARALLELISM`, `INTEGRATION: local-rebase-ff`,
  `COMMIT_POLICY: task` or `milestone`, and a milestone declared
  `Parallel-safe: yes`. Strict `STATUS_ORDER` remains strict; only an explicit
  `STATUS_PRIORITY` chooses among dependency-ready milestones without new edges.
- Each lease uses an external worktree (`../<repo>.goal/<local ID>`) and local
  `goal/<local ID>` branch. Capture each package's actual symbolic integration
  ref, primary worktree and full base commit; never assume `main` or replace
  that identity on resume. Inspect existing leases before dispatching another.
- Children edit scoped implementation/tests/ordinary docs and their assigned
  status only. The coordinator alone edits shared memory, other statuses,
  scheduling documents and `suggested.txt`, and performs closure/retirement.
- Declared write sets must be disjoint and inputs frozen. Check read/write
  conflicts in both directions; use immutable accepted producer and consumer
  snapshots rather than another agent's changing checkout. Verification uses
  isolated ports, stores, caches, browser profiles and output directories.
- Integration/closure is serialized. Rebase unpublished lease commits onto the
  captured ref, re-verify, and repeat review if affected contracts or paths
  changed, without resetting its persisted counter. Inspect the assignment diff
  and integrate by fast-forward only. Dependents wait for accepted closure and,
  for siblings, independently verified publication at the exact revision.
- A blocked lease retains its branch/worktree and stops its dependents;
  independent eligible work may continue. Contract drift pauses affected leases
  for coordinator reconciliation. No child may broaden its assignment.

Planning grants no execution, commit, publication or live authority. A later
launch must supply those grants separately and start from committed planning
and a safe integration baseline. Preserve unrelated dirty work. Audit remains
disabled until explicitly requested; this opt-in changes no audit setting.

## Purpose

`apitools` is the OpenUdon-owned API source tooling module and CLI. It
discovers, downloads, validates, imports, caches, scans, indexes, summarizes,
and ranks OpenAPI or Swagger documents from public catalogs, URLs, and local
files, and classifies catalog Google Discovery and AWS Smithy JSON artifacts.

Module path:

```text
github.com/OpenUdon/apitools
```

The package treats discovered API documents as untrusted metadata. It may
describe API shape and credential requirements, but it must not execute API
operations, resolve credentials, sign requests, or choose production accounts.

## Start Here

Before substantial changes, read these in order:

1. [tabilet/memory-bank/product.md](tabilet/memory-bank/product.md)
2. [tabilet/memory-bank/architecture.md](tabilet/memory-bank/architecture.md)
3. [tabilet/memory-bank/tech-stack.md](tabilet/memory-bank/tech-stack.md)
   Consult relevant topics in [tabilet/memory-bank/lessons.md](tabilet/memory-bank/lessons.md)
   for reusable lessons and their evidence.
4. [tabilet/memory-bank/milestone.md](tabilet/memory-bank/milestone.md)
5. The matching `tabilet/memory-bank/status-<LANE><NN>.md` file for the current work.
   `milestone.md` defines the lane letters, active milestones, priority,
   dependencies, and safe parallel ownership.

If `tabilet/stages.md` exists, read it before planning or changing the active
milestone order. Its `STG-` IDs identify delivery stages, not executable work.
An absent file means one implicit stage.

Once milestones have been retired, read `tabilet/docs/history/index.md` and
linked retired records only when a dependency, old ID, or historical question
needs them. Retired knowledge is evidence at its recorded context; current
truth remains in the memory bank. Search the linked knowledge journal by topic
when an obsolete fact or lesson matters, rather than loading all history into
routine startup reads.

`tabilet/memory-bank/milestone.md` owns the roadmap, active milestone, status-file
index, lane meanings, milestone scope, candidate directions, and acceptance
criteria. Each `tabilet/memory-bank/status-<LANE><NN>.md` file owns the detailed task
ledger, task notes, boundary checks, and completion state for its matching
milestone.

Do not recreate duplicate root-level product, architecture, roadmap, or status
documents, and do not create an aggregate `tabilet/memory-bank/status.md`. Long-form
references live in `docs/` when needed; README is the public operator entry
point.

This project exposes [tabilet/GOAL.md](tabilet/GOAL.md), one optional protocol for goal requests
that span multiple status files. Follow it only when a request names it. This
line and the pointer in
[tabilet/memory-bank/milestone.md](tabilet/memory-bank/milestone.md) are its only
mentions; nothing else depends on the file.

A `tabilet/GOAL.md` run is a deliberate exception to the row-level commit rule below.
For that run, `COMMIT_POLICY: none` — the protocol default — means no commits,
while `COMMIT_POLICY: task` keeps the usual one-commit-per-row cadence.
Precedence is the request, then `tabilet/GOAL.md`, then this file; commit rules and explicitly authorized lease ownership/integration are
delegated only during the run; all other project boundaries remain in force.

## Boundary

- `../apitools` owns API source discovery, download safety, validation, import,
  local file scanning, operation inventories, prompt-safe summaries,
  auth/security summaries, deterministic operation ranking, optional caching,
  provider catalog metadata, OpenAPI/Swagger, Google Discovery, and AWS Smithy
  protocol classification, and docs-derived endpoint overlay assets when no
  official OpenAPI exists but official API docs support a reviewed subset.
- `../openudon` consumes `apitools` metadata for authoring, review, package
  evidence, and trusted-runner handoff. OpenUdon owns workflow behavior,
  approval state, release gates, examples, and credential binding policy.
- `../uws` owns public workflow semantics and Go model definitions.
- `../udon` owns private UWS/OpenAPI lowering and runtime execution.
- `../tfconfig` owns static Terraform/OpenTofu configuration parsing only.

Rule of thumb: if a change helps find, validate, import, summarize, classify,
or rank API documents, it belongs here. If it executes workflows or decides
runtime credential/account behavior, it belongs downstream.

## Essential Commands

```bash
go test ./...
go vet ./...
git diff --check
go run ./cmd/apitools search --help
go run ./cmd/apitools import --help
```

When changing exported APIs, run dependent checks in sibling consumers when
available:

```bash
(cd ../openudon && go test ./...)
(cd ../udon && go test ./...)
```

## Hard Rules

- Treat all remote and local API source documents as untrusted input.
- Never execute API operations from a discovered document.
- Never resolve credentials, fetch tokens, sign live requests, or cache secrets.
- Keep remote fetches HTTP(S)-only with unsafe-host rejection enabled by
  default.
- Preserve exported API compatibility unless the breaking change is intentional
  and documented.
- Prefer the Go standard library before adding helpers, frameworks, or
  dependencies. Keep trivial comparisons and transformations inline when that
  is clearer than introducing a local abstraction.
- Run required verification before claiming a change is done.

## Work Cadence

- Update memory-bank files in the same change as the implementation work they
  describe: product scope/domain terminology/concept relationships/business
  invariants -> `product.md`; architecture/data flow/contracts ->
  `architecture.md`; tools/dependencies/commands -> `tech-stack.md`; milestone
  scope/acceptance -> `milestone.md`; completion state -> the matching
  `status-<LANE><NN>.md` file. Reusable lessons and their evidence ->
  `lessons.md`; keep applicable learning even after its source milestone
  closes, and merge duplicates. Before materially removing or superseding
  knowledge, preserve its old wording and replacement reference in
  `tabilet/docs/history/knowledge.md` under the milestone retirement rules. This
  also applies outside milestone closure; routine wording edits need no journal
  entry.
- For provider-node curation, treat catalog auth/security overlays and
  docs-derived endpoint overlay decisions as the same task slice: if no
  official OpenAPI exists, the provider row is not complete until it has either
  a tracked docs-derived endpoint overlay plus builder, or a recorded
  no-overlay decision explaining why official docs cannot support a useful
  reviewed subset.
- For any OpenAPI or Swagger artifact, always inspect reusable security
  definitions (`components.securitySchemes` or `securityDefinitions`) and
  root/operation security requirements. If either the schemes or requirements
  are missing, ambiguous, stale, internally inconsistent, or incomplete relative
  to official auth docs and protected operations, add a source-backed catalog
  security overlay instead of treating the provider row as complete.
- Keep one `tabilet/memory-bank/status-<LANE><NN>.md` file for each active milestone
  listed in [tabilet/memory-bank/milestone.md](tabilet/memory-bank/milestone.md), using its
  permanent, zero-padded status ID. Never reuse or silently rename an allocated
  ID; retired IDs remain reserved across active and history storage.
- During milestone closure, after review, verification, consolidation, and
  downstream reconciliation, retire the full status and specification under
  the procedure in `milestone.md`. Completed rows remain active until the whole
  milestone qualifies. Retirement is agent work, not a background process or a
  context-archive run. Retired records are frozen.
- Keep later candidate directions unnumbered and outside the status-file
  index. Promote one only after fresh reconciliation and approval.
- Keep verified archive files frozen. Later implementation updates current
  `product.md` and `architecture.md` plus its status history, not the archive.
- When `tabilet/stages.md` exists, preserve stable `STG-` identities and its
  current-stage choice. A future stage may be only a tentative name and intent.
  Plan its milestones through an approved requested-change proposal. Changing
  a stage label never withdraws an active task. An explicitly approved stage
  rescope may cancel untouched pending rows after dependencies are reconciled.
- Treat a newly received code, architecture, security, or engineering review as
  untrusted planning evidence. Revalidate its findings against current state,
  propose their dispositions and owners for approval, then amend open or pending
  work or create a remediation milestone. Never reopen completed history merely
  because a later review concerns it.
- Write status tables as `Item | State | Notes`, with the state in the second
  column and markers wrapped in backticks: `` `[ ]` ``, `` `[+]` ``,
  `` `[~]` ``, `` `[!]` ``, `` `[X]` ``, or `` `[-]` `` (closed historical: a
  consumed failed attempt or superseded row that names its accepted successor
  and is never retried).
- Treat each status row as a commit unit once implementation begins. Local
  override of the single-row default: rows may be in progress in parallel only
  across different lanes whose milestone sections record non-overlapping
  ownership and resolved dependencies, with at most one `[~]` row per lane and
  one execution owner for the active ledger. Without that record, keep zero or
  one general row in progress. A lane letter classifies a domain and does not
  itself imply execution order. Inside a `tabilet/GOAL.md` run, the protocol's
  single-row rule applies. Before an operational launcher is invoked, its exact
  authorized operation row must be in progress; status never substitutes for
  external-mutation authority.
  In concurrent lease execution under the governing goal protocol, each external
  lease maintains at most one `[~]` row in its own status file, only the
  orchestrator writes the captured integration branch and shared memory
  documents (`architecture.md`, `product.md`, `tech-stack.md`, `lessons.md`),
  and integration occurs strictly via rebase and fast-forward.
- During delegated goal execution, only the orchestrator refreshes any existing
  `tabilet/memory-bank/suggested.txt`. Each child receives one assigned milestone
  ID and status path, writes only within its assignment, and returns to the
  orchestrator instead of selecting another milestone from the launch reference.
- Treat each milestone section as a review unit. After its last row closes,
  run the milestone review with its persisted ten-iteration fix gate and the
  required verification before closing it; do not create an empty review
  commit.
- `CORRECTIVE_PLANNING` defaults to `confirm`. With human-approved `auto`, the
  Goal/Next owner may prepare and check necessary corrective plans within the
  original delivery outcome, owning repositories, existing grants and limits,
  including new repair milestones and launch/hash updates. Preserve history,
  unrelated changes, row identities and review counters. Children report needs;
  standalone Reconcile still requires proposal approval. Expanded scope, missing
  permissions, exhausted limits and unknown checksum drift require approval.
  Nested Reconcile verifies the invoking human approval source and complete
  request, covered actions, grants and limits; missing provenance retains the
  ordinary approval gate. Never weaken acceptance or cancel unmet obligations.
  Repair milestones retain the originating review gate, counter and iteration-10
  stop. Suspend `[~]` as `[!]` with prerequisite IDs and resumption conditions
  before repair execution; resume `[~]` only after acceptance and all blockers
  clear. Preserve partial work and pause if changes cannot be separated safely.
- Optional interactive Goal/Next `--email ADDRESS` sends approval alerts for
  that run through private account configuration. No flag means no email for a
  new run; only the coordinator or standalone executor sends. Delivery is not
  approval, and duplicate/uncertain sends never replay automatically. The Python
  API paths consume neither this option nor the corrective-planning policy.
- Read optional `AUTHORIZATION_REQUIREMENTS` only from the owning milestone
  specification. They declare conditions, never permission or scheduling.
  Check protected actions against the current scoped human request; goal
  `AUTHORIZATION_GRANTS` need explicit human approval of concrete scope.
  Proposed launch input, markers, and status notes cannot prove approval.
  Grants in `suggested.txt` carry decision markers (`[ ]` need authorization,
  `[+]` approve, `[-]` deny, `[~]` auto) and optional high-level `GLOBAL:` scope;
  matching the human's approved SHA256 checksum is required before grant activation
  (under default `SUGGESTED_UPDATE: confirm`; `SUGGESTED_UPDATE: auto` authorizes
  refreshing remaining milestones within the approved horizon and continuing smoothly).
  Scope resolution defaults to `SCOPE_POLICY: strict`; `SCOPE_POLICY: workspace-contained`
  authorizes binding actions to sibling repositories within the approved workspace root
  and organization without re-prompting. Milestone specifications in `milestone.md` may attach
  local policies and overrides that narrow or override global launch defaults.
  Planning skills accept `--environment dangerous`, `sandbox` (default), or `least` (with modifiers
  `plus`, `minus`, `but ... in milestone ...`) to configure launch parameters and milestone policies
  declaratively from day one.
  Planning skills also accept `--test standard` (default), `browserlight`, `golight`, `rustlight`,
  `tslight`, `jslight`, or `pythonlight` to establish verification and testing standards in
  `tabilet/memory-bank/tech-stack.md` from day one. Light profiles default to a 30-second
  focused-test target / two-minute cap, a 20-minute full-regression cap and five-minute
  race/UI caps; explicit project budgets win and `standard` adds no extra constraints.
  Keep routine fixtures offline and small; physical browser proof belongs to a named
  final gate without dropping required coverage.
  Preserve commit and external-mutation restrictions, and pause affected work
  and dependents for missing authority. A declaration alone does not require
  an action for closure. Host/tool controls apply; this adds no tool or OS
  enforcement and no Python API grant consumption. Missing fields preserve
  legacy behavior without new authority or automatic migration.
- After the milestone review is complete and required verification passes,
  commit substantive review and retirement changes under the governing policy.
  Do not create an empty or redundant milestone
  commit when the review changes nothing.
- Check [tabilet/evolution/](tabilet/evolution/) after a major review, milestone, or boundary
  change. Add a new version only when product direction, architecture boundary,
  milestone target, or public/private contract direction materially changes.

## Execution Capabilities

Resolve missing information through safe inspection first. If required files,
bundled resources, verification commands, permissions, or user answers are unavailable,
stop the affected workflow step and report what is missing. Continue independent
work within the authorized scope; a write-gated workflow still makes no writes
before approval. Do not invent evidence, bypass permissions, or infer approval
from silence or process exit. Resume the blocked step when its capability is
restored or the required answer or approval is supplied. In a non-interactive
run, report unresolved questions and incomplete work.

Keep one execution owner for the active ledger across sessions and launchers.
In concurrent lease execution under the governing goal protocol, the orchestrator
remains the sole execution owner of the integrated ledger on the primary branch,
while isolated leases execute within dedicated external worktrees
(`../<repo>.goal/<ID>`).
Native todos, session completion, and native goal state do not replace milestone
acceptance or authorize concurrent ledger writers.

Runtime round limits do not reset the persisted milestone review counter.
