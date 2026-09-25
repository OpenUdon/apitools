# Retired milestone M73 - CLI And Coordinated Release Cleanup

**Milestone.** M73
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M73.md
**Source specification.** tabilet/memory-bank/milestone.md#m73---cli-and-coordinated-release-cleanup
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M73 - CLI And Coordinated Release Cleanup

**Goal.** Finish the breaking-contract rollout with a thin, consistent CLI and
verified sibling dependency updates.

**Scope.** Centralize command parsing and exit/output policy, keep catalog stats
in the root package and refresh registration in `sqlitecache`, remove
duplicated flag scaffolding, update public/private documentation, migrate all
sibling consumers, and pin the reviewed apitools version in OpenUdon, Ramen,
and Udon.

**Dependencies.** S01, M72, and C01 reviews must be complete.

**Parallel ownership.** M73 is the serialized integration and release lane.

**Downstream impacts.** Help exits 0 on stdout, usage failures exit 2 on
stderr, runtime failures exit 1 on stderr, positional values are not mistaken
for help flags, and the former OAuth entry point is documented under Udon.

**Acceptance.** All migrated sibling consumers compile and pass their complete
test suites; apitools has no stale compatibility surface or staticcheck
findings; CLI contract tests pass; memory-bank and evolution records match the
released behavior; full cross-repository verification and diff checks pass.
````

## Status record

````markdown
# Status M73 - CLI And Coordinated Release Cleanup

**Acceptance.** The CLI is thin and consistent, all breaking contracts are
migrated and pinned across sibling consumers, documentation matches ownership,
and full cross-repository verification passes.

| Item | State | Notes |
|---|---|---|
| Centralize CLI parsing and exit/output policy | `[+]` | Shared non-terminating parser routes genuine help to stdout, usage diagnostics to stderr, and uses named 0/1/2 exit policy; regressions cover every command, missing values, and `--help` as data. Verified with full tests, vet, staticcheck, and diff checks. |
| Move remaining business policy out of cmd | `[+]` | Stats aggregation is root-package policy via `BuildCatalogStatsReport`; refresh-result registration and cache-relative containment now live in `sqlitecache.RegisterCatalogRefreshResults`. The CLI retains argument adaptation and rendering only. Verified with full tests, vet, staticcheck, and diff checks. |
| Migrate and pin all sibling consumers | `[+]` | Published apitools `v0.0.0-20260820001033-3a9864901572` and Authoring `v0.0.0-20260819235644-d73b65ae179c`; OpenUdon and Ramen pin both reviewed contracts, and Udon pins apitools while retaining its documented private-workspace replace. All three pass standalone full tests and vet; resolved module graphs show the intended versions. |
| Reconcile public docs and private harness | `[+]` | README, Gmail helper docs, product, architecture, tech stack, milestone/status, and evolution v22 now describe the multi-family source boundary, Udon OAuth ownership, CLI exit/stream contract, package-owned cache registration, current sibling pins, and remaining final release gate. |
| Run coordinated release verification | `[+]` | Apitools workspace/standalone tests and vet, staticcheck, generated-catalog freshness, offline catalog quality, CLI help, and diff checks pass; OpenUdon's complete release SaaS gate passes the 11-component integration matrix, 21 headed loopback scenarios, 8 journeys, UI browser checks, variants, scorecard, docs, fixture builds, approvals, and dry runs; Ramen workspace/standalone tests and vet and Udon's quality gate pass. Resolved modules use the published apitools/Authoring revisions, and the stale-surface audit finds no apitools OAuth command, OAuth dependency, token exchange, callback server, or refresh-token path. |

## Boundary Checks

- Command names remain stable except the explicitly relocated OAuth command.
- CLI code performs presentation and argument adaptation, not catalog or cache
  policy.
- No tfconfig changes are required.
````
