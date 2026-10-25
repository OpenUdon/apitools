# M91 — Common eligibility and mapping helpers

**Stage/owner:** STG-13; APItools.
**Planning approval:** Owner-approved full rescope on 2026-10-25 (approval: “I aapaprove”), following the earlier whole-stage approval. Planning only; no execution, commit or grant activation.
**State:** Pending; no task started. Review 0/10.
**Specification:** [milestone.md](milestone.md#stg13-m91).
**Depends on:** None.
**Source baseline:** `fe54b3ba939427d86ebbdd4a47b0dccb519cc26c`; existing dirty planning/history was inspected and preserved. Reconcile live source and exact accepted producers before implementation.
**Ownership:** One serial coordinator; parallel-safe no. Each row is a coherent future task-commit unit.
**Authorization:** See AUTHORIZATION_REQUIREMENTS in the specification. Status grants no action; consumed producer prerequisites require independently verified publication.

**Review provenance.** Kinet Stage 13 plan review (`kinet-stage13.md`), 2026-10-25; source commit baseline not supplied. Review H4 (High/local P2) and M1 (Medium/local Lower): the parent go.work substitutes sibling modules, and U16/P06 supply no APItools contract prerequisite. M91.4 owns explicit standalone publication evidence; dispatch priority retains the product-first order. Revalidated at APItools `fe54b3ba939427d86ebbdd4a47b0dccb519cc26c` with the existing uncommitted Stage 13 planning inspected; source/history changes were not absorbed. This is review intake, not a closing review iteration.

## Tasks

| Item | State | Notes |
|---|---|---|
| M91.1 — Define common eligibility and mapping rules | `[ ]` | Inventory existing Kinet/Ramen rules; freeze the small shared behavior over published operation shapes, leaving source selection/admission to each product. Keep unknown/truncated shapes and ambiguity explicit. |
| M91.2 — Implement additive pure helpers | `[ ]` | Provide deterministic eligibility and typed field/binding matching with requiredness and credential exclusions; keep native source/operation identity intact and avoid provider transport or runtime calls. |
| M91.3 — Add shared positive and refusal cases | `[ ]` | Cover supported shapes, path/query/header/body/response bindings, placeholders, incompatible types, missing fields, stale identities and ambiguous matches. Exercise both product use cases without broad new protocol features. |
| M91.4 — Qualify and hand off exact producer | `[ ]` | Run focused and required standalone checks with GOWORK=off; record go env GOWORK=off, exact selected/compiled dependency identities and absence of directory replacements for adoption/publication proof. Prepare independently verified ordinary publication and consumer fixtures. APItools:M90 stays reserved and unscheduled. |

## Acceptance and verification

APItools exposes small additive helpers for candidate eligibility and field mapping, with shared conformance cases and product-specific admission kept in consumers.

- Preserve source digest, native operation selector, requiredness/type/incompleteness and credential exclusions. Ambiguous or unsupported mappings remain explicit unresolved choices.
- Helpers are metadata-only: no credentials, provider transport, execution, universal planner, product imports or automatic approval. Existing API consumers continue to build.

Planned verification: Focused operation-shape/candidate/mapping/operationlifecycle tests across supported source identities; standalone tests/vet and existing affected static/race gates. Apply the [shared verification policy](../../../kinet/docs/stage13.md#verification-policy), including explicit `GOWORK=off` for adoption and closure checks. All implementation evidence is pending.

## Review

Review gate iterations: 0

Task reviews: not started. Whole-milestone review: 0/10, not started. Follow the package's existing persisted review and closure procedure; planning, terminal rows and successful builds do not establish acceptance.
