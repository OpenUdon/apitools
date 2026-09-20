# Status M75 - Operation Lifecycle Ranking Ownership

State: Complete

| Item | State | Notes |
|---|---|---|
| Add provenance-aware lifecycle ranking | `[+]` | Added `operationlifecycle` over `apitools.OperationSummary`, documented scoring constants, retained same-source/ambiguity behavior, and restricted `/upload` normalization to explicit Google Discovery provenance. |
| Migrate Authoring consumers | `[+]` | OpenUdon and Ramen now translate downstream context to source-bearing operation summaries; Authoring removes its misplaced package and public-surface entries. |
| Verify workspace integration | `[+]` | Apitools full test/vet, Authoring standalone/race, OpenUdon/Ramen full workspace suites, OpenUdon scorecard, compatibility, and diff checks pass. |
| Publish and pin coordinated revisions | `[+]` | Published Apitools `v0.0.0-20260820042238-d51b61ead067` before Authoring `v0.0.0-20260820042256-2f73e3526583`; OpenUdon and Ramen pin both without replacements and pass standalone `GOWORK=off` test/vet. |

## Boundary Checks

- Ranking consumes metadata only and never fetches or executes operations.
- Product goals influence optional lifecycle role selection but no workflow,
  credential, account, or desired-state semantics move into apitools.
- Provider-specific path handling requires explicit source provenance.

Commit tracking: Apitools implementation `d51b61e`; coordinated Authoring
`243121d`, `eee0c0d`, and `2f73e35`; OpenUdon adapter/pins `43294bf` and
`065bb84`; Ramen adapter/pins `330fd89` and `721b540`.
