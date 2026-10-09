# M84 — Reserved OpenAPI header projection remediation

**Stage/owner:** STG-11 post-acceptance remediation; APItools.
**Planning approval:** User requested package-local milestones on 2026-10-08 and selected required fixes only.
**State:** Executing M84.4; M84.1/.2/.3 complete; one serial owner; whole pre-publication review1 passed (1/10); ordinary publication/final review pending.
**Current goal authority:** The user invoked memory-bank-goal, authorized task commits and needed reviewed normal fast-forward GitHub source/closure publication to the recorded existing main refs, and confirmed serial execution Udon:M52 -> APItools:M84 -> Kinet:M50. Udon:M52 is accepted/published/retired, and M84 is the next ready serial unit.
**Review source:** Stage 11 review 2 — all milestones, all five packages (`stage11-review-2.md`), AP1.
**Source priority / local severity:** P2 / P2; confirmed projection behavior and current Kinet incompatibility gate.
**Review and revalidation baseline:** `de3f16acbf12c7b632ee0ed9be02efaaf2c2be4b`; clean worktree. No relevant uncommitted changes used.
**Evidence:** operation_shapes_openapi.go parameter projection; UWS binding.input_required and Kinet W18 compatible-publication gate.
**Lineage:** Accepted M82/M83; [M83 retired record](../docs/history/status-M83.md).

## Scope, dependencies and ownership

Correct AP1 only. For supported OpenAPI 3 inputs, reserved header parameter declarations must not manufacture required workflow input bindings. Preserve the actual security-scheme/requirement metadata, genuine headers, source digests, deterministic identities and independent shape verification.

No new UWS/OpenUdon source is scheduled. Retain the accepted UWS M09 root/codec versions and current public contracts. Udon:M52 -> APItools:M84 is serial scheduling only, with no new producer dependency. Direct downstream: Kinet:M50; ordinary public SDK consumer verification must establish compatibility with the corrected projection.

## Tasks

**Accepted serial predecessor:** Udon:M52 runtime9c99eab6051f250dc49e2e52d8246526280046ca,
modulev0.0.0-20261009004332-9c99eab6051f, independently published with source/evidence
closuref2a56e71293997f73597aafb7c44b747bd0e525c and whole review4. This is not an
APItools producer dependency, but current private SDK consumer proof must use
this exact accepted runtime alongside retained UWS root/codec and SDK801.
The [owner record](../../../udon/tabilet/docs/history/status-M52.md) and
[handoff](../../../udon/docs/m52-release-handoff.md) bind both sums/CLI/graph proof.

| Item | State | Notes |
|---|---|---|
| M84.1 — Apply reserved-header parameter rules | `[+]` | AP1. Exclude reserved Accept, Content-Type and Authorization parameter declarations according to supported OpenAPI rules, including reusable local parameter references and the supported header-name semantics. A required declaration must not create binding.input_required. Preserve genuine required headers and security scheme/OR/AND information; never bind credential values as literals or resolve credentials. |
| M84.2 — Prove projection and native-authoring compatibility | `[+]` | Add deterministic bounded fixtures for supported OpenAPI versions, direct/referenced declarations, genuine headers and security requirements. Independently reproduce shape bytes and show current UWS/OpenUdon assessment admits valid native authoring without fabricated reserved inputs while incomplete real bindings still refuse. Run source/determinism/identity/forgery and metadata-only compatibility checks. |
| M84.3 — Qualify clean source and prepare reviewed publication | `[+]` | Run required offline owner and proportionate public/private consumer checks against the exact candidate graph, preserving frozen corpus/wires and all unrelated pins. Record exact source/module/build identities and a concrete fresh source-publication proposal for normal fast-forward to the existing main ref. Review the exact candidate before exercising the current publication authority; bootstrap proof remains labelled and cannot satisfy ordinary publication. |
| M84.4 — Publish reviewed source and close ordinary handoff | `[~]` | The initial missing-authority blocker is resolved by the user's current GitHub publication authorization. After M84.3 qualification and review of the exact proposal, mark this operation in progress and publish reviewed M84 source plus scoped evidence/closure by normal fast-forward to git@github.com-tabilet:OpenUdon/apitools.git refs/heads/main. Independently verify ordinary full-SHA/module/sums/source/consumer proof; complete whole review, downstream reconciliation, consolidation and retirement. No additional push approval is required within this granted scope. No force push, tags, images, deployment or live requests. |

## Acceptance and verification

AP1 is fixed at the source projection and independently verified consumer boundary. Valid workflows using reserved header declarations do not become incompatible because of fabricated required inputs; genuine header/security constraints remain intact. No schema or public API/wire change is needed, and source metadata remains bounded, deterministic, value-free and network/credential-free by default.

Run go test ./..., go vet ./..., git diff --check and existing applicable corpus/source-format/identity guards. Follow the established checker policy without inventing a new waiver. Verify the current SDK consumer and actual selected private consumer graph with retained Go 1.26.6/GOWORK=off and cached offline dependencies. Reuse identical identity-bound verification evidence; repeat checks for changed inputs or remaining failures. Publication acquisition/proof follows the exact reviewed proposal under the current scoped user grant.

**Compatibility:** Corrected shape bytes/digests may change for affected sources. Do not rewrite previously approved packages or grants; Kinet:M50 owns fresh assessment/publication and worker adoption. Existing schemas, frozen sources and history remain unchanged.
**Closing review:** whole pre-publication iteration1 passed, 1/10; persist whole-pass starts, resume the same interrupted pass, maximum ten, no P1/P2 carry-forward.
**Execution/commit:** One serial owner; later confirmed goal uses task commits. No execution or commit follows from this planning write.
**External mutations:** Reviewed M84 source and scoped documentation/closure, normal fast-forward only, to git@github.com-tabilet:OpenUdon/apitools.git refs/heads/main under the current user grant. Prepare and verify the exact outgoing proposal before acting; M83 grants remain consumed.
**Downstream:** Kinet:M50 requires accepted independently published M84 identities before adoption. UWS/OpenUdon lower-priority findings are not scheduled.

## M84.1 — completed, 2026-10-09

The source-neutral OpenAPI projector ignores the three reserved header names
case insensitively after bounded local Parameter Object resolution and existing
identity/location/duplicate checks, before schema/requiredness/serialization
projection. Swagger is excluded from this new rule. Schemes/security alternatives
are processed independently and unchanged. The 108-case supported3.0/3.1
direct/local-ref matrix covers path-item inheritance and operation declarations
and validates without fabricated bindings. Focused Go1.26.6 check passed with
GOWORK/GOPROXY off; `/tmp/am84-proof/row1-corrected.log` retains it. The initial
build failed from an inherited absent TMPDIR; owned TMPDIR/GOTMPDIR corrected
the environment and no source workaround was added. Contract documentation links
the official OpenAPI parameter/header semantics.

## M84.2 — completed, 2026-10-09

The new owner regressions preserve genuine header, query, cookie and path inputs;
ignored schemas/serialization cannot introduce completeness gaps. Security retains
a two-member AND alternative and a one-member OR alternative, including an
API-key scheme whose parameter name is Authorization. Genuine missing inputs and
incomplete AND security refuse. Escaped local pointer chains resolve; cycles,
missing and remote Parameter Object refs refuse. Deterministic table bytes and raw
SHA/native selector equality reproduce; fabricated reserved inputs and stale raw
source identities are rejected. The Swagger projection control stays unchanged.
Focused owner shapes pass; inherited source fixtures/golden bytes remain unchanged.

Current ordinary SDK801 plus UWS M09 consumes the exact M84.1 source as the
explicitly local bootstrap version v0.0.0-20261009012011-04c158e2b41c, without
workspace/directory replacements. This is compatibility evidence, not publication.
Sixteen source-backed Build/Assess/Verify/determinism/raw-source cases pass over
3.0.3/3.1.0 and direct/ref forms; bound genuine inputs pass and missing genuine
header/query inputs retain binding.input_required. The unchanged M83 source fails
all four reserved-only controls with exactly three fabricated required-input
diagnostics, retaining the before-fix evidence. Initial fixture spelling used
unsupported request.headers; corrected to existing request.header, with no sibling
code or contract change. Inert [consumer source](../../docs/m84-consumer-fixture/public_test.go.txt)
is reproducible; logs are `/tmp/am84-proof/public-row2-corrected.log`,
`public-row2-before-fix.log` and `row2-guards.log`. Exact final clean-source and
actual public/private selected graph qualification remain M84.3/.4 work.

## Whole review iteration1 — pre-publication passed, 2026-10-09

Persisted before reviewing the full implementation, regressions, candidate-bound
owner/public/private qualification, source files, invariants and exact publication
proposal. M84.4 ordinary publication is deliberately still pending; this pass
reviews the pre-publication gate and cannot declare final acceptance. No P1/P2
carry-forward is permitted. An interrupted pass resumes iteration1.

The full pre-publication pass, including independent Tier0 read-only review, found
no P1/P2/P3. All27 event/log identities, source621-file Git/ZIP/extraction equality,
selected/compiled graphs, cache parity, preserved failures and exact outgoing
planning/proposal scope independently revalidated. Publication remains M84.4's
required gate; final whole review resumes at iteration2, never resets to1.

## M84.3 — completed, 2026-10-09

Qualified clean source0c1383c6e2059e89cbf100ed8cb689be64d9ba57, time
2026-10-09T01:23:35Z, modulev0.0.0-20261009012335-0c1383c6e205,
sumh1:GeJhKB7I42cH+bWw+8CoiPImyLmHKsqJih6E7uvUdlg= and
modsumh1:WmUXlfBBoaI6vtv/wzeZfyO8q/pZMhnHiBckkAthYfk=. The source621-file
archive/ZIP/extraction is byte-exact. Owner full tests/races/vet/build/catalog/
format/diff, current public/private consumer races/vet/build, current SDK packagev3
full tests and all3 final modverify checks pass. Actual selected graphs77/73/171
and compiled source fields6469/5586/8307 preserve SDK801/UWSb099/privateM52runtime9c
and existing version-only Docker-to-Moby mapping, with no directory/workspace
replacement or public private-Udon import. All15657/18160/35037 selected artifact
source files and actual compiled/test/embed/C/assembly ownership are qualified.

The real owned cache uses immutable hardlinks. Initial TMPDIR/missing-version/
symlink-root/native-lockfile seeding failures remain preserved; exact prior test
input byte parity proves no dependency/code changes. Compatible staticcheck has
the exact same6 accepted M83 observations and four byte-unchanged source files.
[Qualification](../../docs/m84-qualification.md), [local proof](../../docs/m84-local-proof.json),
[source inventory](../../docs/m84-source-files.json), [guards](../../docs/m84-frozen-guards.json)
and [exact publication proposal](../../docs/m84-publication-proposal.md) bind the
source/toolchain/modules/CLI identities and scoped outgoing content. Candidate
proof is labelled bootstrap, never ordinary publication. Whole prepub1 passed
without findings; M84.4 is selected under the current exact normal-FF source grant.
