# M84 — Reserved OpenAPI header projection remediation

**Stage/owner:** STG-11 post-acceptance remediation; APItools.
**Planning approval:** User requested package-local milestones on 2026-10-08 and selected required fixes only.
**State:** Pending; four task rows; publication authorized after qualification/review; closing review not started (0/10).
**Current goal authority:** The user invoked memory-bank-goal, authorized task commits and needed reviewed normal fast-forward GitHub source/closure publication to the recorded existing main refs, and confirmed serial execution Udon:M52 -> APItools:M84 -> Kinet:M50. This status remains pending until its serial/dependency gates pass.
**Review source:** Stage 11 review 2 — all milestones, all five packages (`stage11-review-2.md`), AP1.
**Source priority / local severity:** P2 / P2; confirmed projection behavior and current Kinet incompatibility gate.
**Review and revalidation baseline:** `de3f16acbf12c7b632ee0ed9be02efaaf2c2be4b`; clean worktree. No relevant uncommitted changes used.
**Evidence:** operation_shapes_openapi.go parameter projection; UWS binding.input_required and Kinet W18 compatible-publication gate.
**Lineage:** Accepted M82/M83; [M83 retired record](../docs/history/status-M83.md).

## Scope, dependencies and ownership

Correct AP1 only. For supported OpenAPI 3 inputs, reserved header parameter declarations must not manufacture required workflow input bindings. Preserve the actual security-scheme/requirement metadata, genuine headers, source digests, deterministic identities and independent shape verification.

No new UWS/OpenUdon source is scheduled. Retain the accepted UWS M09 root/codec versions and current public contracts. Udon:M52 -> APItools:M84 is serial scheduling only, with no new producer dependency. Direct downstream: Kinet:M50; ordinary public SDK consumer verification must establish compatibility with the corrected projection.

## Tasks

| Item | State | Notes |
|---|---|---|
| M84.1 — Apply reserved-header parameter rules | `[ ]` | AP1. Exclude reserved Accept, Content-Type and Authorization parameter declarations according to supported OpenAPI rules, including reusable local parameter references and the supported header-name semantics. A required declaration must not create binding.input_required. Preserve genuine required headers and security scheme/OR/AND information; never bind credential values as literals or resolve credentials. |
| M84.2 — Prove projection and native-authoring compatibility | `[ ]` | Add deterministic bounded fixtures for supported OpenAPI versions, direct/referenced declarations, genuine headers and security requirements. Independently reproduce shape bytes and show current UWS/OpenUdon assessment admits valid native authoring without fabricated reserved inputs while incomplete real bindings still refuse. Run source/determinism/identity/forgery and metadata-only compatibility checks. |
| M84.3 — Qualify clean source and prepare reviewed publication | `[ ]` | Run required offline owner and proportionate public/private consumer checks against the exact candidate graph, preserving frozen corpus/wires and all unrelated pins. Record exact source/module/build identities and a concrete fresh source-publication proposal for normal fast-forward to the existing main ref. Review the exact candidate before exercising the current publication authority; bootstrap proof remains labelled and cannot satisfy ordinary publication. |
| M84.4 — Publish reviewed source and close ordinary handoff | `[ ]` | The initial missing-authority blocker is resolved by the user's current GitHub publication authorization. After M84.3 qualification and review of the exact proposal, mark this operation in progress and publish reviewed M84 source plus scoped evidence/closure by normal fast-forward to git@github.com-tabilet:OpenUdon/apitools.git refs/heads/main. Independently verify ordinary full-SHA/module/sums/source/consumer proof; complete whole review, downstream reconciliation, consolidation and retirement. No additional push approval is required within this granted scope. No force push, tags, images, deployment or live requests. |

## Acceptance and verification

AP1 is fixed at the source projection and independently verified consumer boundary. Valid workflows using reserved header declarations do not become incompatible because of fabricated required inputs; genuine header/security constraints remain intact. No schema or public API/wire change is needed, and source metadata remains bounded, deterministic, value-free and network/credential-free by default.

Run go test ./..., go vet ./..., git diff --check and existing applicable corpus/source-format/identity guards. Follow the established checker policy without inventing a new waiver. Verify the current SDK consumer and actual selected private consumer graph with retained Go 1.26.6/GOWORK=off and cached offline dependencies. Reuse identical identity-bound verification evidence; repeat checks for changed inputs or remaining failures. Publication acquisition/proof follows the exact reviewed proposal under the current scoped user grant.

**Compatibility:** Corrected shape bytes/digests may change for affected sources. Do not rewrite previously approved packages or grants; Kinet:M50 owns fresh assessment/publication and worker adoption. Existing schemas, frozen sources and history remain unchanged.
**Closing review:** not started, 0/10; persist whole-pass starts, resume the same interrupted pass, maximum ten, no P1/P2 carry-forward.
**Execution/commit:** One serial owner; later confirmed goal uses task commits. No execution or commit follows from this planning write.
**External mutations:** Reviewed M84 source and scoped documentation/closure, normal fast-forward only, to git@github.com-tabilet:OpenUdon/apitools.git refs/heads/main under the current user grant. Prepare and verify the exact outgoing proposal before acting; M83 grants remain consumed.
**Downstream:** Kinet:M50 requires accepted independently published M84 identities before adoption. UWS/OpenUdon lower-priority findings are not scheduled.
