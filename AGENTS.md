# AGENTS.md

## Approved Stage 11 planning

[Stage 11](../kinet/docs/stage11.md) coordinates both refactoring phases across five package-local ledgers. Read the local [milestones](tabilet/memory-bank/milestone.md) before selecting work. APItools will produce source-neutral UWS operation shapes from its existing eight source families. It still performs no credential resolution, workflow approval or runtime execution. Kinet and OpenUdon consume exact published metadata contracts through isolated parsing paths.
Planning is approved; implementation and named publication/deployment authority are separate. One serial execution owner, offline fixtures, exact upstream reconciliation and persisted milestone reviews apply. Completed records and frozen consumer pins stay preserved.

Requested features, candidate promotions, and future direction changes after
initialization follow the requested-change procedure in
[tabilet/memory-bank/milestone.md](tabilet/memory-bank/milestone.md). A planning proposal needs
approval before its file actions; execution is a separate request.

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
Precedence is the request, then `tabilet/GOAL.md`, then this file; only commits are
delegated, and only during the run.

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
- Treat each milestone section as a review unit. After its last row closes,
  run the milestone review with its persisted ten-iteration fix gate and the
  required verification before closing it; do not create an empty review
  commit.
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
Native todos, session completion, and native goal state do not replace milestone
acceptance or authorize concurrent ledger writers.

Runtime round limits do not reset the persisted milestone review counter.
