# Retired milestone M91 - Common eligibility and mapping helpers

**Milestone.** M91
**Outcome.** completed
**Retired.** 2026-10-25
**Source status.** tabilet/memory-bank/status-M91.md
**Source specification.** tabilet/memory-bank/milestone.md#stg13-m91
**Evidence.** 3954824c6bfdd036f2f8c5f9b57a17b385b51d59
**Worktree.** includes uncommitted changes
**Review.** passed
**Review iterations.** 2
**Verification.** Qualified0ff local full1 PASS30.768s/race1 PASS89.189s; inherited staticexit1sixdiagnostics/absentVCS retained. Exact d8creation-only publication1 PASS1.720709s/independentGitHub identities; ordinaryacquisition1 PASS9.797954s/exact sums/origin/707files. Reusabletypedconsumer49RUN49PASS/vet/build/correctedprobe selected92/compiled41/no replacements. ExistingOpenUdon full1 PASS27.075342s; approvedexactUdoncopy readiness2 PASS0.958568s/346frozen identities/892packages, full1 PASS81.822281s/test-vet-build-CLI/compiled76/exactM91/no directoryreplace/original692files unchanged. All earlier failures and unqualified scopes retained.
**Consolidated into.** [product](../../memory-bank/product.md), [architecture](../../memory-bank/architecture.md), [tech-stack](../../memory-bank/tech-stack.md), [lessons](../../memory-bank/lessons.md); docs/m91-published-proof.json, reusable fixture and exact W22/M136 handoffs. Actual pre-closure HEAD; uncommitted review/P3/retirement metadata is not claimed in that commit.

## Milestone specification

````markdown
<a id="stg13-m91"></a>

### M91 — Common eligibility and mapping helpers

**Stage/owner.** Kinet STG-13; APItools. **State.** Completed; all four task rows and fresh whole2 pass.
**Depends on.** None.
**Downstream.** [Kinet:W22](../../../kinet/tabilet/memory-bank/status-W22.md), [Ramen:M136](../../../ramen/tabilet/memory-bank/status-M136.md)
**Outcome.** APItools exposes small additive helpers for candidate eligibility and field mapping, with shared conformance cases and product-specific admission kept in consumers.

**Write set.** APItools public operation metadata/eligibility/mapping helpers, operationlifecycle integration where needed, versioned fixtures and SDK documentation. Package-local documentation/status follows the owning change; sibling changes have their own owners. **Parallel-safe: no.**

**Acceptance.**

- Preserve source digest, native operation selector, requiredness/type/incompleteness and credential exclusions. Ambiguous or unsupported mappings remain explicit unresolved choices.
- Helpers are metadata-only: no credentials, provider transport, execution, universal planner, product imports or automatic approval. Existing API consumers continue to build.

**Verification.** Focused operation-shape/candidate/mapping/operationlifecycle tests across supported source identities; standalone tests/vet and existing affected static/race gates. Apply the [shared verification policy](../../../kinet/docs/stage13.md#verification-policy), including `GOWORK=off` for every adoption/closure check. Clean implementation source `0ff4ef561dcc5f784a630b1559ac688bca35802d` passes full local qualification (30.768s) and affected races (89.189s), with 77 selected/39 compiled modules and no replacements. [Qualification](../../docs/m91-producer-qualification.md) preserves six current static observations under the existing M85 policy, actual absent VCS stamps, and all earlier failures. [Published ordinary proof](../../docs/m91-published-proof.json) records exact source/acquisition/reusable typed consumer and both available existing-consumer checks. Fresh whole2 passes; no zero-checker/global-clean/live claim follows.

**Review provenance.** Kinet Stage 13 plan review (`kinet-stage13.md`), 2026-10-25; source commit baseline not supplied. Review H4 (High/local P2) and M1 (Medium/local Lower): the parent go.work substitutes sibling modules, and U16/P06 supply no APItools contract prerequisite. M91.4 owns explicit standalone publication evidence; dispatch priority retains the product-first order. Revalidated at APItools `fe54b3ba939427d86ebbdd4a47b0dccb519cc26c` with the existing uncommitted Stage 13 planning inspected; source/history changes were not absorbed. This is review intake, not a closing review iteration.

**Execution policy.** Future execution follows the [coordinator](../../../kinet/docs/stage13.md#verification-policy) and effective goal policy. Planning activates no action, grant or old budget; required publication is independently verified before consumption.

```yaml
AUTHORIZATION_REQUIREMENTS:
  git.commit:
    via: goal-policy
    executor: coordinator
  cli.local:
    via: goal-policy
    executor: coordinator
    scope:
      repository: /home/peter/Workspace/apitools
      task: declared-implementation-and-verification
  git.push:
    via: explicit
    executor: coordinator
    scope:
      repository: /home/peter/Workspace/apitools
      remote: git@github.com-tabilet:OpenUdon/apitools.git
      ref: refs/heads/stage13/m91-operation-mapping-d8c3555d
      revision: d8c3555d365bd0414b4d89fcab7074542983e5c6
      mode: fast-forward
      creation_only: true
      expected_ref: absent; empty expected-ref lease
```

**Task/review owner.** [status-M91.md](status-M91.md) owns the 4 pending rows and persisted review; [the coordinator](../../../kinet/docs/stage13.md) owns cross-package acceptance and priority.

**Scoped execution approval — 2026-10-25.** The owner said “I approve” after both exact linked proposals. The [approval record](../../../kinet/docs/stage13-scoped-approval-20261025.json) binds the reviewed hashes. Only this exact new review branch, disclosed source-test CI, two named module acquisitions and exact disposable compatibility inputs are excepted from the no-push/external-mutation policy. Main denials remain; publication, consumption and acceptance are not yet claimed. Existing review counter1 and every earlier outcome remain.

**Final acceptance.** Both original acceptance lines and all four tasks fulfilled. Fresh whole2 passes all five lenses/no P1/P2; sole Lower stale-state summaries reconciled. Original source full/race, exact independently verified publication/acquisition, reusable49test ordinary consumer and existing OpenUdon/Udon checks qualify at exact recorded scopes. Inherited checker diagnostics, absent VCS stamps and failures remain; W22/M136 actual product adoption is downstream.

````

## Status record

````markdown
# M91 — Common eligibility and mapping helpers

**Stage/owner:** STG-13; APItools.
**Planning approval:** Owner-approved full rescope on 2026-10-25 (approval: “I aapaprove”), following the earlier whole-stage approval. Planning only; no execution, commit or grant activation.
**State:** Completed and accepted at fresh whole2/10. All four original task rows, exact producer publication/ordinary acquisition, reusable typed consumer and available OpenUdon/Udon compatibility pass. Inherited static/VCS and original failures remain.
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
| M91.4 — Qualify and hand off exact producer | `[+]` | Local full/race and inherited six-diagnostic/VCS limits retained. Exact d8 publication/ordinary acquisition and707Git/ZIP/extracted files verified. Reusable consumer49RUN49PASS/vet/build/corrected probe, selected92/compiled41/no replacements. Existing OpenUdon full1 PASS27.075342s. Approved exact da05 Udon-copy metadata correction: readiness2 PASS0.958568s/346frozen identities/892packages; full1 PASS81.822281s tests/vet/build/CLI/actual new M91 BuildInfo. Original Udon untouched, existing version-only Moby replacement retained; grants consumed. [Exact proof](../../docs/m91-published-proof.json). Fresh whole2 started, original failures retained. |

## Acceptance and verification

APItools exposes small additive helpers for candidate eligibility and field mapping, with shared conformance cases and product-specific admission kept in consumers.

- Preserve source digest, native operation selector, requiredness/type/incompleteness and credential exclusions. Ambiguous or unsupported mappings remain explicit unresolved choices.
- Helpers are metadata-only: no credentials, provider transport, execution, universal planner, product imports or automatic approval. Existing API consumers continue to build.

Required verification: Focused operation-shape/candidate/mapping/operationlifecycle tests across supported source identities; standalone tests/vet and existing affected static/race gates. Apply the [shared verification policy](../../../kinet/docs/stage13.md#verification-policy), including explicit `GOWORK=off` for adoption and closure checks. Local implementation/qualification evidence is recorded above and in the linked proof; publication, ordinary consumer evidence and authoritative whole review remain pending.

## Review

Review gate iterations: 2

All four task reviews and authoritative fresh whole2 pass. WHOLE_MILESTONE_REVIEW
gpt-6.1-sol/xhigh reviewed the full baseline1aa42bbbd6cce01b902e1ce612243ec0b61b928a
through3954824c6bfdd036f2f8c5f9b57a17b385b51d59. All five lenses pass; no P1/P2
remain, prior evidence finding M91-W1 resolved. One Lower/P3 M91-W2 identified
stale current status/stack/lesson/spec summaries; coordinator reconciles them
during closure, preserving phase evidence below. This metadata-only correction
changes no helper/schema/corpus/test/dependency input and needs no new run or
whole iteration. Existing six static diagnostics/actualexit1, absent VCS stamps
and all failed preflights/probe-oracle corrections remain explicit, not waived.

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

## Scoped continuation approval — 2026-10-25

The invoking human “I approve” approved both exact reviewed scopes identified
in the [approval record](../../../kinet/docs/stage13-scoped-approval-20261025.json).
The coordinator may create only the named frozen review branch, acquire only
the two exact modules after independent publication verification, and qualify
the named ordinary consumer fixtures. Main/other pushes, deployment and all
other acquisitions remain denied. Publication/consumption are still unperformed;
the row stays suspended while Kinet:M59 owns serial execution. Review counter1,
prior findings and all consumed checks remain; no acceptance is claimed.

## Coordinator publication1 — selected

The actually approved exact d8c3555d new review-branch creation is selected
with independent GitHub absent-ref verification, exact origin/tree/no hook
preflight, empty expected-ref lease, hard300s and existing batch Git SSH/current
host-key checks. ONE attempt1 reserved; no success is inferred. All other pushes
including main/tags/future revisions remain denied. Root is sole writer and
other general rows are complete/suspended. Whole counter1 remains; nextwhole2
follows the unperformed ordinary consumer/existing-consumer evidence.

## Actual publication/consumer checkpoint — 2026-10-25

Approved d8 creation-only push1 PASS1.720709s and exact ordinary acquisition1 PASS9.797954s; independent GitHub ref/tree/time, both sums/origin and all707Git/ZIP/extracted files agree. Actions consumed/ended, no replay. [Exact proof](../../docs/m91-published-proof.json) retains every receipt, selected92/actual41compiled includingmain/no replacements, reusable inputs and all failed checks. Combined consumer conformance tests 49RUN/49PASS/zeroFAILskip, vet/build pass; failed new probe oracle corrected at source and changed probe build/run/BuildInfo PASS1.648532s. Unchanged tests/vet reused;13operations/eightfamilies and typed card execute. Producer code/corpus/acceptance unchanged.

Exact e486 OpenUdon disposable consumer full1 PASS27.075342s (public boundary/alltests/vet/build/publicprobe/BuildInfo). Exact da05 Udon copy readonly readiness1 FAIL1.033451s before compilation: missing indirect requirements. Nonmutating offline tidy-diff diagnostic1 exit1/diffproduced/stderr0 identifies seven root declarations/additional cached sums; original/copy module files unchanged, full envelope unspent. Root promotion may expand/prune-selected graph: freeze precise proposed result before asking approval, no unchanged337 claim. Current version-only Docker/Moby replacement remains, no directory replacement. Kinet docs/stage13-udon-compatibility-correction.md is proposed only. Row suspended, zero general rows in progress; whole1 partly unresolved, nextwhole2 not started. All original scope, limits, failures, UNKNOWNs and audit-off preserved.

## Exact Udon-copy correction actually approved — 2026-10-25

Invoking human I approve binds frozen Kinet proposal89852c and graphc5b729.
Root resumes M91.4: apply only exact two prepared metadata files to disposable
da05 copy, ONE corrected readonly readiness2<=120s and, only after graph matches
all346frozen identities and packages are ready, the existing unspent full1<=1200s.
No new full allowance, source/owner Udon mutation, download/push/private action,
acceptance waiver or reviewcounter change. Earlier failures and graph337/whole1
remain. Selection is no success; only this row is general in progress.

## Fresh authoritative whole2 — started

Actual human I approve binds exact89852c correction. Approved metadata bytes
are applied only to disposable da05 copy; readonly readiness2 PASS0.958568s,
all346frozen graph identities/892package records/errorfree. Conditional original
full1 PASS81.822281s: tests68.498508s, vet6.839445s, build3.632932s, CLIbuild
2.725535s, help0.115339s and actualBuildInfo0.008447s, all exit0/originalWait/no
timeout. Output uses regular-file captures; no pipe-EOF claim. CLI stderr is
normal help usage, not a failure/zero-stderr claim. Exact new M91 dependency
version/sum agrees selected versus compiled; all original692tracked files and
all non-module disposable sources match da05. No source repair or original
Udon mutation. Individual RUN/skip counts are not retained by ordinary package
output and no live browser/provider/private-runtime qualification is inferred.

Root inspected the complete task4 evidence/reusable inputs/current module
proof and original-source preservation; task review passes. M91-W1 evidence
is now resolved at original scope. Fresh read-only whole2/10 reviews the entire
milestone; whole1, all inherited static observations, missing VCS stamps and
earlier failed preflights/oracle remain. Grants for exact pushes/acquisitions/
Udon correction/readiness/full are consumed/ended with no replay. Current
source/helper/corpus behavior is unchanged; no final acceptance until whole2.

## Final authoritative whole2 — passed and reconciled

All five lenses pass: correctness/reliability, security/privacy,
compatibility/contracts, tests/docs and acceptance/scope. No P1/P2; sole
Lower/P3 M91-W2 stale current summaries is resolved in final status/stack/
lesson/spec. Reviewer inspected exact resolution/schema/presence/duplicates/
placeholders/cancellation, metadata-only credential exclusions and no-fetch
bounds, both test files,19cases/foursources/eightfamily preservation and full
acceptance fidelity. Helpers remain exactly qualified0ff, no missing/partial/
incorrect acceptance or unrequested implementation.

Read-only review checked all79retained output hashes/sizes, all707published
Git/ZIP/extracted files, ordinary49RUN49PASS/reusabletypedprobe/selected92/
compiled41/no replacements, existingOpenUdon full/source1215files, existing
Udon approvedmetadata/frozen346graph/892errorfreepackages/compiled76 and all
692originaltracked/non-modulecopy source files. Prior evidenceP2M91-W1 is
resolved. Reviewer ran no test/build/network/cache/mutation. All exact action
grants are consumed/ended; no main/futurepush/download/oldgrant/UNKNOWN reset.
Current facts/lesson and W22/M136 exact producer handoffs are consolidated;
actual product adoption remains downstream. Stage13 stays open.
````
