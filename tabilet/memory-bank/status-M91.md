# M91 — Common eligibility and mapping helpers

**Stage/owner:** STG-13; APItools.
**Planning approval:** Owner-approved full rescope on 2026-10-25 (approval: “I aapaprove”), following the earlier whole-stage approval. Planning only; no execution, commit or grant activation.
**State:** M91.1–M91.3 complete; M91.4 next; review 0/10. Source publication remains denied by the human's “commit in all packages but no push” instruction.
**Specification:** [milestone.md](milestone.md#stg13-m91).
**Depends on:** None.
**Source baseline:** `fe54b3ba939427d86ebbdd4a47b0dccb519cc26c`; existing dirty planning/history was inspected and preserved. Reconcile live source and exact accepted producers before implementation.
**Ownership:** One serial coordinator; parallel-safe no. Each row is a coherent future task-commit unit.
**Authorization:** See AUTHORIZATION_REQUIREMENTS in the specification. Status grants no action; consumed producer prerequisites require independently verified publication.

**Review provenance.** Kinet Stage 13 plan review (`kinet-stage13.md`), 2026-10-25; source commit baseline not supplied. Review H4 (High/local P2) and M1 (Medium/local Lower): the parent go.work substitutes sibling modules, and U16/P06 supply no APItools contract prerequisite. M91.4 owns explicit standalone publication evidence; dispatch priority retains the product-first order. Revalidated at APItools `fe54b3ba939427d86ebbdd4a47b0dccb519cc26c` with the existing uncommitted Stage 13 planning inspected; source/history changes were not absorbed. This is review intake, not a closing review iteration.

## Tasks

| Item | State | Notes |
|---|---|---|
| M91.1 — Define common eligibility and mapping rules | `[+]` | Inventoried existing APItools ranking/shapes/lifecycle, Kinet mapping/native projection and Ramen role/save bindings against actual sources. [Frozen common behavior](../../docs/operation-mapping-v1.md) scopes explicit exact-shape eligibility and whole-field schema mappings, retaining consumer admission, nested authoring and security authority. Task review checked scope and source locations; `git diff --check` passed. No implementation or whole review claimed. |
| M91.2 — Implement additive pure helpers | `[+]` | Added `AssessOperationEligibility` and `MatchOperationFields`, reusing UWS table resolution/schema validation and existing credential-name policy. Exact source/native identity, unknown/truncated/ambiguous states, whole-schema direction, requiredness, placeholders and credential exclusions are retained; no transport/runtime/product imports. Task review corrected unresolved-pair coverage and cancellation propagation. Go1.26.6, `GOWORK=off GOPROXY=off GOSUMDB=off GOAUTH=off GOTOOLCHAIN=local GOFLAGS=-mod=readonly`: focused new mapping/eligibility tests passed (0.077s); affected existing shape/selection/candidate checks passed (1.150s). `git diff --check` passed. |
| M91.3 — Add shared positive and refusal cases | `[+]` | Added [19 versioned actual-source cases](../../testdata/operation-mapping/v1/README.md), corpus SHA256 `1f456c87fb0d4a594a2b26fdff5328f71a445157cf0feba990e17929e07fd847`, business-message and fake Cloudflare resource sources. Covers all planned field locations, exact candidate/native identities, schemas/requiredness, placeholder/duplicate/type/missing/stale/ambiguous/security refusals; existing frozen eight-family partial shapes stay unresolved. Task review added conservative session/signature/opaque-key header exclusions matching inventoried Ramen policy. Focused mapping/eligibility plus existing lifecycle tests passed (root 0.172s, lifecycle 0.071s) under the same standalone offline Go1.26.6 profile; `git diff --check` passed. |
| M91.4 — Qualify and hand off exact producer | `[ ]` | Run focused and required standalone checks with GOWORK=off; record go env GOWORK=off, exact selected/compiled dependency identities and absence of directory replacements for adoption/publication proof. Prepare independently verified ordinary publication and consumer fixtures. APItools:M90 stays reserved and unscheduled. |

## Acceptance and verification

APItools exposes small additive helpers for candidate eligibility and field mapping, with shared conformance cases and product-specific admission kept in consumers.

- Preserve source digest, native operation selector, requiredness/type/incompleteness and credential exclusions. Ambiguous or unsupported mappings remain explicit unresolved choices.
- Helpers are metadata-only: no credentials, provider transport, execution, universal planner, product imports or automatic approval. Existing API consumers continue to build.

Planned verification: Focused operation-shape/candidate/mapping/operationlifecycle tests across supported source identities; standalone tests/vet and existing affected static/race gates. Apply the [shared verification policy](../../../kinet/docs/stage13.md#verification-policy), including explicit `GOWORK=off` for adoption and closure checks. All implementation evidence is pending.

## Review

Review gate iterations: 0

Task reviews: M91.1 scope/inventory, M91.2 implementation and M91.3 conformance checks passed. Whole-milestone review: 0/10, not started. Follow the package's existing persisted review and closure procedure; planning, terminal rows and successful builds do not establish acceptance.
