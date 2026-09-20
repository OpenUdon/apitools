# Status C01 - Generated Catalog And Security Aggregation

**Acceptance.** Reviewed JSON is the catalog source of truth, deterministic
generation enforces invariants, security status is per spec with explicit mixed
and conflict outcomes, aggregation is indexed, and milestone verification
passes.

| Item | State | Notes |
|---|---|---|
| Define and populate validated catalog JSON bundles | `[+]` | Added canonical `apitools.catalog.v1` reviewed JSON for 317 candidates, 316 providers, 908 spec references, 276 classifications, and 249 overlays plus strict unknown-field, ordering, identity, scope, and cross-reference validation. Canonical JSON parity with the legacy built-ins and full/focused/static checks pass. |
| Generate deterministic catalog indexes | `[+]` | Replaced the hand-maintained candidate, provider, classification, and overlay literals with an embedded reviewed JSON bundle plus deterministic checked-in ID/provider indexes. `cmd/cataloggen` canonicalizes and validates the bundle, atomically writes a golden count/digest manifest and Go indexes, and rejects stale output in `-check` mode; package initialization verifies all bytes, counts, and index coverage before exposing cloned records. Focused tests, stale-output tests, generator check, and diff checks pass. |
| Replace overlay precedence with scoped aggregation | `[+]` | Added deterministic provider/spec `SecurityDisposition` records with separate classification and overlay statuses, explicit overlay supplementation, provider-level `mixed`, same-scope `conflict`, provider-wide fallback, and exact selected-spec resolution. Conflicting built-in dispositions and derived statuses used as evidence fail generation/validation; advisory, resolver, inspection, quality, security-report, and security-audit outputs preserve the scoped ledger. Focused tests cover ordering, mixed/conflict, generation rejection, copy safety, inspection, quality, and selection. |
| Index advisory and statistics aggregation | `[+]` | Added immutable `CatalogIndex` maps for provider aliases, exact spec references, provider overlays, and precomputed scoped security reports; advisory and resolver paths now validate/build once and return cloned lookup results instead of rescanning the full catalog per provider. Moved protocol selection/order, artifact-registry buckets, and refresh-evidence statistics from `cmd/apitools` into the reusable root `BuildCatalogStatsReport` API, including all supported source families. Determinism, copy isolation, invalid input, protocol coverage, package behavior, CLI output, vet, and staticcheck tests pass. |
| Review and verify C01 | `[+]` | Acceptance review found and fixed remaining scattered top-level count literals by exporting `BuiltInCatalogManifest`, adding bundle digest/count parity, generation-level invalid ID/reference/scope cases, and updating the curation workflow. Generator freshness, catalog quality, full and `GOWORK=off` tests, catalog race tests, vet in workspace/standalone modes, staticcheck, offline CLI JSON assertions, diff checks, OpenUdon workspace tests, and Udon workspace/standalone tests pass. OpenUdon's stale standalone apitools pin still lacks an earlier S01 field and is explicitly owned by M73 coordinated pin migration. |

## Boundary Checks

- Generated metadata remains reviewable and checked in.
- Catalog security remains advisory and never resolves credentials.
- No implicit sorted-order precedence survives generation or aggregation.
