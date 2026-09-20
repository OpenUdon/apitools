# AGENTS.md

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
4. [tabilet/memory-bank/milestone.md](tabilet/memory-bank/milestone.md)
5. The matching `tabilet/memory-bank/status-<LANE><NN>.md` file for the current work.
   `milestone.md` defines the lane letters, active milestones, priority,
   dependencies, and safe parallel ownership.

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
that span multiple status files. Follow it only when a request names it.

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
- Run required verification before claiming a change is done.

## Work Cadence

- Update memory-bank files in the same change as the implementation work they
  describe: product scope -> `product.md`; architecture/data flow/contracts ->
  `architecture.md`; tools/dependencies/commands -> `tech-stack.md`; milestone
  scope/acceptance -> `milestone.md`; completion state -> the matching
  `status-<LANE><NN>.md` file.
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
- Keep one `tabilet/memory-bank/status-<LANE><NN>.md` file for each milestone listed in
  [tabilet/memory-bank/milestone.md](tabilet/memory-bank/milestone.md), using its permanent,
  zero-padded status ID. Never reuse or silently rename an allocated ID.
- Keep later candidate directions unnumbered and outside the status-file
  index. Promote one only after fresh reconciliation and approval.
- Write status tables as `Item | State | Notes`, with the state in the second
  column and markers wrapped in backticks: `` `[ ]` ``, `` `[+]` ``,
  `` `[~]` ``, `` `[!]` ``, or `` `[X]` ``.
- Treat each status row as a commit unit once implementation begins. When
  independent rows in different lanes are worked in parallel, record
  non-overlapping ownership and dependencies in `milestone.md`; a lane letter
  classifies a domain and does not itself imply execution order.
- Treat each milestone section as a review unit. After its last row completes,
  run the milestone review and required verification before closing it; do not
  create an empty review commit.
- Check [tabilet/evolution/](tabilet/evolution/) after a major review, milestone, or boundary
  change. Add a new version only when product direction, architecture boundary,
  milestone target, or public/private contract direction materially changes.
