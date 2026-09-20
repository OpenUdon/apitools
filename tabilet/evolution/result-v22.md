# Result V22 - Untrusted Metadata And Ownership Hardening

The coordinated hardening direction is implemented across apitools, Authoring,
OpenUdon, Ramen, Udon, and the private harness.

- M72 confined network and filesystem behavior, bounded SQLite growth, removed
  live OAuth and credential parsing from apitools, and moved state/PKCE-protected
  Google authorization plus secret file persistence into Udon.
- S01 applied common source byte/depth/work guards, parser fuzz regressions,
  prompt-safe identifier/text/collection/operation budgets, preserved OpenAPI
  security alternatives, and migrated OpenUdon/Ramen prompt consumers.
- C01 made the reviewed catalog JSON authoritative, generated its manifest and
  lookup indexes deterministically, and moved reusable security/statistics
  aggregation out of the CLI.
- M73 centralized CLI parsing and the 0/1/2 stdout/stderr policy, moved refresh
  registration into `sqlitecache`, published reviewed apitools and Authoring
  pseudo-versions, and pinned OpenUdon, Ramen, and Udon consumers.

The published integration revisions are apitools
`v0.0.0-20260820001033-3a9864901572` and Authoring
`v0.0.0-20260819235644-d73b65ae179c`. OpenUdon and Ramen resolve both without
the workspace; Udon pins apitools while retaining its documented private
sibling replaces.

Row-level full tests, standalone consumer tests, vet, staticcheck, generation,
and diff checks passed during implementation. The final coordinated release
audit also passed complete apitools, OpenUdon, Ramen, and Udon gates, including
sandboxed headed-browser scenarios, UI checks, variant and scorecard checks,
documentation validation, package builds, approvals, and dry runs. No stale
apitools OAuth or token-exchange path remains. M73 is complete; this result does
not broaden apitools into execution, credential, or workflow ownership.
