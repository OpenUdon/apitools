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
