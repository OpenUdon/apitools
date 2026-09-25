# Retired milestone C01 - Generated Catalog And Security Aggregation

**Milestone.** C01
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-C01.md
**Source specification.** tabilet/memory-bank/milestone.md#c01---generated-catalog-and-security-aggregation
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## C01 - Generated Catalog And Security Aggregation

**Goal.** Replace hand-maintained catalog literals and implicit overlay
precedence with validated source data and deterministic generated indexes.

**Scope.** Store reviewed provider bundles as JSON, generate checked-in Go
indexes deterministically, validate catalog and overlay invariants at generation
time, build advisory indexes once, and represent security disposition per spec
with explicit mixed and conflict outcomes.

**Dependencies.** M72 safe-file and refresh contracts. C01 starts only after
M72 review closes.

**Parallel ownership.** C01 owns catalog source data, generator, generated
provider/candidate/overlay indexes, advisory aggregation, and catalog tests.

**Downstream impacts.** Catalog callers consume per-spec security dispositions
and explicit provider-level mixed status rather than sorted last-wins overlay
behavior.

**Acceptance.** Generation is deterministic and checkable; invalid IDs,
references, overlay scopes, or conflicting dispositions fail generation;
catalog golden count/digest replaces scattered count assertions; advisory
aggregation is indexed; full catalog, root, standalone, downstream, vet,
staticcheck, CLI, and diff checks pass.
````

## Status record

````markdown
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
````
