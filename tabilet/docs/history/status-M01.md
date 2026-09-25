# Retired milestone M01 - Private Harness Bootstrap

**Milestone.** M01
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M01.md
**Source specification.** tabilet/memory-bank/milestone.md#m01---private-harness-bootstrap
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M01 - Private Harness Bootstrap

**Goal.** Install the private planning harness for `apitools` without exposing
the harness files in the public repository.

**Scope.**

- Create canonical harness files under `../tofu/apitools`.
- Symlink `AGENTS.md`, `tabilet/memory-bank/`, and `tabilet/evolution/` from `../apitools` to
  the private harness snapshot.
- Ignore harness paths in the public `../apitools` repository.
- Preserve the existing public package behavior and verification commands.

**Acceptance.** `../apitools` exposes the harness through local symlinks,
`../tofu` tracks the canonical files, `../apitools` ignores the harness paths,
and `go test ./...`, `go vet ./...`, and `git diff --check` pass in
`../apitools`.
````

## Status record

````markdown
# Status M01 - Private Harness Bootstrap

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

| Item | State | Notes |
|---|---|---|
| Canonical harness directory created under `../tofu/apitools` | `[+]` | Initial private harness files added. |
| `../apitools/AGENTS.md` symlinked to private harness | `[+]` | Public tracked file removed from index; local symlink points to `../tofu/apitools/AGENTS.md`. |
| `../apitools/memory-bank` symlinked to private harness | `[+]` | New local symlink points to `../tofu/apitools/memory-bank`. |
| `../apitools/evolution` symlinked to private harness | `[+]` | New local symlink points to `../tofu/apitools/evolution`. |
| Public repo ignores harness paths | `[+]` | Ignore entries cover `AGENTS.md`, `memory-bank`, and `evolution` symlink paths. |
| Verification completed | `[+]` | `go test ./...`, `go vet ./...`, `git diff --check`, `go run ./cmd/apitools search --help`, and `go run ./cmd/apitools import --help` passed in `../apitools`. |
````
