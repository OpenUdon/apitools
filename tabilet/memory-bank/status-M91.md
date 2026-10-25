# M91 — Common eligibility and mapping helpers

**Stage/owner:** STG-13; APItools.
**Planning approval:** Owner-approved full rescope on 2026-10-25 (approval: “I aapaprove”), following the earlier whole-stage approval. Planning only; no execution, commit or grant activation.
**State:** M91.1–M91.3 complete; M91.4 locally qualified with publication/ordinary evidence pending. Authoritative whole review iteration 1 completed with one P2 evidence blocker and no production-code P1/P2 defect. The full gate has not passed. The human's no-push instruction remains binding.
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
| M91.4 — Qualify and hand off exact producer | `[!]` | Clean source `0ff4ef561dcc5f784a630b1559ac688bca35802d` in coordinated detached input `/home/peter/Workspace/apitools.goal/m91-qualification`. [Local qualification](../../docs/m91-producer-qualification.md) / [exact module proof](../../docs/m91-local-proof.json): full invocation 1 PASS 30.768s (tests/vet/CLI/catalog), race invocation 1 PASS 89.189s; Go1.26.6, selected 77 / compiled 39 including main, zero replacements. Static invocation 1 exit 1 at 17.660s retains six pre-existing observations under coordinator-reconciled unchanged M85 policy; actual log SHA/codes/locations retained, historical raw log absent, no byte-identical-output claim/new-file diagnostic/suppression/repair/retry. One explicit VCS build passed 2.541s with identical CLI bytes, but both BuildInfo outputs omit VCS fields; Git clean-revision evidence is separate. Narrow evidence/readiness task review and format/patch checks passed; whole review 1/10 completed with evidence blocker M91-W1. Required publication denied and ordinary new-consumer acquisition/adoption unperformed; resume only after scoped authority/evidence reconciliation. APItools:M90 stays reserved and unscheduled. |

## Acceptance and verification

APItools exposes small additive helpers for candidate eligibility and field mapping, with shared conformance cases and product-specific admission kept in consumers.

- Preserve source digest, native operation selector, requiredness/type/incompleteness and credential exclusions. Ambiguous or unsupported mappings remain explicit unresolved choices.
- Helpers are metadata-only: no credentials, provider transport, execution, universal planner, product imports or automatic approval. Existing API consumers continue to build.

Required verification: Focused operation-shape/candidate/mapping/operationlifecycle tests across supported source identities; standalone tests/vet and existing affected static/race gates. Apply the [shared verification policy](../../../kinet/docs/stage13.md#verification-policy), including explicit `GOWORK=off` for adoption and closure checks. Local implementation/qualification evidence is recorded above and in the linked proof; publication, ordinary consumer evidence and authoritative whole review remain pending.

## Review

Review gate iterations: 1

Task reviews: M91.1 scope/inventory, M91.2 implementation and M91.3 conformance checks passed. M91.4 narrow DEFAULT evidence/readiness review passed for local proof only; required ordinary publication/consumer evidence is missing. Authoritative whole-milestone review: 1/10 completed with evidence blocker M91-W1; no pass or final acceptance is claimed. Follow the package's existing persisted review and closure procedure; planning, terminal rows and successful builds do not establish acceptance.

## Authoritative whole review iteration 1 — blocked on evidence

Fresh read-only WHOLE_MILESTONE_REVIEW gpt-6.1-sol/xhigh reviewed the complete
candidate from integration baseline `1aa42bbbd6cce01b902e1ce612243ec0b61b928a`
through `d7ef01bd6195ba12e44df4e8ac4933adf2570c74`, plus these review-state updates.
Qualified production/test/corpus bytes are `0ff4ef561dcc5f784a630b1559ac688bca35802d`;
later commits record qualification and documentation only. All five lenses are
required, including precise allocation of publication and consumer evidence
under original M91 acceptance/task text. Existing M85 static policy, actual
nonzero diagnostics and absent VCS stamps remain explicit. No requirement,
push denial or downstream published-consumption gate is waived. The review
completed without a production-code P1/P2 defect or any P3 finding. Its sole
P2 finding, M91-W1, is the unperformed independently verified ordinary producer
publication, reusable consumer proof and available existing-consumer checks.
The owner-local `(devel)` build cannot satisfy those obligations. The human's
no-push denial and unapproved module acquisition remain binding.

The review covered all five lenses: correctness/reliability;
security/privacy; compatibility/contracts; tests/docs; acceptance/scope.
The reviewer inspected both new test files, all 19 corpus cases and four
sources, eight-family preservation, retained log/corpus hashes, existing
consumer adapters, and the unchanged six checker diagnostics. No tests,
builds, cache writes or mutations were performed by the reviewer. Existing
nonzero static output and absent VCS stamps remain disclosed.

M91 owns producer qualification, independently verified publication, reusable
consumer fixtures and existing-consumer compatibility evidence. New Kinet
author-worker/business-flow adoption belongs to W22; Ramen resource-lifecycle
adoption belongs to M136. Their completion is not an M91 prerequisite. The
whole gate remains blocked on M91-W1; after its evidence is obtained, the next
complete review is iteration 2 of the same ten-iteration gate. No counter resets.
