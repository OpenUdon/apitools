# M85 — Reserved-header candidate and inventory remediation

**Owner/stage:** APItools; STG-11 required GOAL Step4 successor.
**Approval:** User approved the complete four-row reconciliation, fix execution and source publication on 2026-10-09. Task commits; serial APItools:M85 -> Kinet:M50.
**State:** Execution started; M85.1 complete; M85.2 selected; whole review not started0/10.
**Finding:** M50-NATIVE-F01; source priority not supplied; local P2; confirmed. [Exact evidence](../../../kinet/docs/m50-candidate-header-finding.md), JSON SHA2564060c75e73ec15715971e37e99d2d510136d63118ff9d2d7308acdd4bba9738f.
**Baseline/revalidation:** APItools34f5633b941192f6719f0abfc2783cd6675434ad, clean; published source0c1383c6e2059e89cbf100ed8cb689be64d9ba57. Kinet1aff008db0a463b360502f288a08f49ae0aafaa2 preserves actual workerced5468 local/rootless mismatch and bounded diagnostic-test changes. M84 acceptance/history remain frozen.
**Dependencies:** Accepted M84; retain SDK801b45bb1aec297631b2e5d87beb1977788a1ab1, UWS root/codec b099f6803277ae94c7e9f1da0904a0140b278f20 and private consumer UdonM52 runtime9c99eab6051f250dc49e2e52d8246526280046ca. Existing module versions/sums remain until independently qualified successor source is selected downstream.

| Item | State | Notes |
|---|---|---|
| M85.1 — Align reserved-header effective input semantics | `[+]` | Correct OpenAPI3 inventory/consumer-summary/candidate paths to ignore case-insensitive Accept/Content-Type/Authorization header parameters consistently with shapes. Preserve bounded reference/inheritance handling, source identity, genuine inputs/security, Swagger2, other families, API signatures and wire schemas. No host workaround or new capability. |
| M85.2 — Regress owner projections and candidate matching | `[~]` | Cover casing, path/operation inheritance and supported local references, valid/ignored metadata, all three reserved headers, genuine required header/query/cookie/path/body controls, independent security and negative candidate matches. Prove source/selector/privacy/budget/no-fetch invariants and preserve existing shape verification. |
| M85.3 — Qualify exact source and SDK/CLI consumers | `[ ]` | Prove coherent candidate plus shape/SDK Build/Assess/Verify outcomes through current SDK801, owner/CLI and public/private selected graphs. Full owner checks/races/vet/build/catalog/format/compatible checker, exact source snapshot/CLI and selected-unused/compiled archive proof; reuse unchanged identity-bound dependencies, retain failed evidence. Prepare exact reviewed source-publication proposal. |
| M85.4 — Publish and independently qualify the successor | `[ ]` | Under named user grant, normal fast-forward reviewed source and scoped evidence/closure to the existing main ref; record independently observed exact module/version/sums/source/CLI/consumer identities. No directory bootstrap or stale gate substitution. Complete final persisted whole review, release handoff and parent-owned consolidation/reconciliation/retirement. |

## Acceptance and verification

Close the confirmed P2 without changing raw source bytes, selectors, public wire schemas, credential/privacy/runtime behavior or frozen source contracts. Reserved header declarations cannot block otherwise valid candidate binding; genuine required input/security evidence remains required. Owner and SDK/CLI candidate+shape proofs and ordinary complete source/module/CLI checks must pass. KinetM50 resumes only after accepted/published M85 gate, then rebuilds and qualifies changed workers/integration/bundle. The exact ced5468/M84 worker artifacts are preserved failed-context evidence, not final acceptance.

Use Go1.26.6, GOWORK=off, offline fake/disposable resources by default. Run full owner tests/races/vet/build/catalog/format/diff, compatible checker with unchanged existing six observations, exact public/private SDK consumer tests and complete selected/unused archives/compiled source ownership/CLI reproduction. No new waiver, dependency, SDK/UWS/Udon code or general hardening is planned.

## Review, authority and downstream

**Whole review:** not started,0/10. Persist each started pass and findings before fixes; resume interruptions at same count; maximum10, no P1/P2 carried. Completed rows do not establish acceptance.
**Commit policy:** task; required substantive clean-source checkpoints may be additional commits without premature row completion; no amend/history rewrite.
**Publication grant:** The user explicitly approved APItools:M85 fix and source publication on 2026-10-09. Authorized git@github.com-tabilet:OpenUdon/apitools.git refs/heads/main, reviewed normal fast-forward source/scoped evidence/closure after exact prepared proposal and qualification. Changed target/ref, unrelated outgoing commits or new authority needs must be resolved before that action. No force/tags/images/deployment/live provider/API/model/mail/host/account/user-ledger/registration action.
**Downstream:** M85 -> KinetM50. Root owns shared memory/suggested/downstream/consolidation/retirement; delegated executor owns this status, code/tests and source-adjacent evidence only. One serial execution owner across all ledgers. Original18 rows plus these four rows must close, along with every review/qualification/retirement gate.
**Evolution/candidates:** unchanged; this restores existing supported native-owner semantics, with no Stage12 or lower-priority promotion.

## M85.1 evidence

Inventory removes supported OpenAPI3 reserved header declarations after existing local parameter-reference lookup and before schema/readiness projection. Shapes reuse the same source-version/location/case predicate. No summary/candidate host filter is added. Focused inventory/candidate and existing108-case shape matrix pass (`/tmp/am85-proof/row1-corrected.log`); initial fixture omitted requiredness and correctly produced indeterminate matching (`row1.log`), now corrected to always-available genuine id. Public API/wires, manifests and dependencies unchanged.
