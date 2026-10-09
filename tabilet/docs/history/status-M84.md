# Retired milestone M84 — Reserved OpenAPI header projection remediation

**Milestone.** M84
**Outcome.** completed
**Retired.** 2026-10-09
**Source status.** tabilet/memory-bank/status-M84.md
**Source specification.** tabilet/memory-bank/milestone.md#m84--reserved-openapi-header-projection-remediation
**Evidence.** b7425e2030489ce25a0ee523f83d694176f949f0
**Worktree.** includes uncommitted changes
**Review.** passed
**Review iterations.** 2
**Verification.** Full owner tests/races/vet/build/catalog/format/diff,108-case header matrix,16-case current SDK and M83 negative controls, fresh ordinary Origin/full-SHA/version/sums/all621 source files, exact owner77/public73/private171 archives/compiled source/CLI parity and all19 ordinary events pass. Six unchanged checker observations retain existing policy; cache failures remain recorded.
**Consolidated into.** [product](../../memory-bank/product.md), [architecture](../../memory-bank/architecture.md), [stack](../../memory-bank/tech-stack.md), [lessons](../../memory-bank/lessons.md), [ordinary proof](../../../docs/m84-ordinary-proof.json), [publication](../../../docs/m84-publication.md) and [handoff](../../../docs/m84-release-handoff.md); exact Kinet:M50 adoption reconciled. Included coordinator closure/retirement changes postdate the observed evidence head.
**Accepted source.** 0c1383c6e2059e89cbf100ed8cb689be64d9ba57
**Module.** v0.0.0-20261009012335-0c1383c6e205

## Milestone specification

``````markdown
## M84 — Reserved OpenAPI header projection remediation

**Goal/scope:** Fix AP1's fabricated required Accept/Content-Type/Authorization inputs in supported OpenAPI 3 operation shapes. Preserve real header/security constraints, deterministic source identity, local references, bounds and independent reproduction; no credential resolution or runtime execution.
**Provenance:** Stage 11 review 2 (`stage11-review-2.md`), source P2/local P2, confirmed; report/revalidation commit `de3f16acbf12c7b632ee0ed9be02efaaf2c2be4b`, clean worktree. Lineage M82/M83; completed records stay frozen.
**Dependencies:** Existing accepted UWS M09 root/codec, no new producer milestone dependency. Udon M52 is a serial scheduling predecessor.
**Acceptance:** A reserved header declaration cannot fabricate binding.input_required and block otherwise valid native authoring. Genuine headers and security OR/AND remain correctly represented; current UWS/OpenUdon consumers independently verify corrected shapes.
**Tasks/review:** All four correction/consumer/qualification/publication units are complete. Final whole review2 passes, with exact ordinary source/consumer proof and scoped evidence publication observed at b7425e2030489ce25a0ee523f83d694176f949f0. Owner consolidation and downstream reconciliation complete; normal retirement applies.
**Verification:** Offline owner tests/vet/source/corpus/identity guards and proportionate actual consumer checks, retained toolchain/current selected graphs, existing checker policy, diff hygiene and independent ordinary module proof after named publication.
**Compatibility/downstream:** Affected shape digests require fresh downstream assessment/publication/grants. Kinet M50 owns adoption/worker/bundle qualification; old packages, public schemas/wires and unrelated pins stay frozen. No new evolution direction, live action or broader source-format hardening is included.
``````

## Status record

``````markdown
# M84 — Reserved OpenAPI header projection remediation

**Stage/owner:** STG-11 post-acceptance remediation; APItools.
**Planning approval:** User requested package-local milestones on 2026-10-08 and selected required fixes only.
**State:** Accepted 2026-10-09; all four task gates, ordinary source/current-consumer proof, whole review2 and coordinator consolidation/reconciliation pass; ready for validated retirement.
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
| M84.4 — Publish reviewed source and close ordinary handoff | `[+]` | The initial missing-authority blocker is resolved by the user's current GitHub publication authorization. After M84.3 qualification and review of the exact proposal, mark this operation in progress and publish reviewed M84 source plus scoped evidence/closure by normal fast-forward to git@github.com-tabilet:OpenUdon/apitools.git refs/heads/main. Independently verify ordinary full-SHA/module/sums/source/consumer proof; complete whole review, downstream reconciliation, consolidation and retirement. No additional push approval is required within this granted scope. No force push, tags, images, deployment or live requests. |

## Acceptance and verification

AP1 is fixed at the source projection and independently verified consumer boundary. Valid workflows using reserved header declarations do not become incompatible because of fabricated required inputs; genuine header/security constraints remain intact. No schema or public API/wire change is needed, and source metadata remains bounded, deterministic, value-free and network/credential-free by default.

Run go test ./..., go vet ./..., git diff --check and existing applicable corpus/source-format/identity guards. Follow the established checker policy without inventing a new waiver. Verify the current SDK consumer and actual selected private consumer graph with retained Go 1.26.6/GOWORK=off and cached offline dependencies. Reuse identical identity-bound verification evidence; repeat checks for changed inputs or remaining failures. Publication acquisition/proof follows the exact reviewed proposal under the current scoped user grant.

**Compatibility:** Corrected shape bytes/digests may change for affected sources. Do not rewrite previously approved packages or grants; Kinet:M50 owns fresh assessment/publication and worker adoption. Existing schemas, frozen sources and history remain unchanged.
**Closing review:** whole pre-publication iteration1 and final whole iteration2 passed, 2/10; persist whole-pass starts, resume the same interrupted pass, maximum ten, no P1/P2 carry-forward.
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

## M84.4 — source and ordinary handoff complete, 2026-10-09

M84.4 was in progress before the source launcher. Under the current exact grant,
reviewed normal fast-forward advanced git@github.com-tabilet:OpenUdon/apitools.git
refs/heads/main fromde3f16acbf12c7b632ee0ed9be02efaaf2c2be4b to
3773bd485fe31599e355719e0c5ecea0373bc553. Independent ls-remote observed that exact
source/evidence head. The six outgoing commits and all scoped paths/proposal hash
were reviewed; no force, tag/image/deployment/other-owner/live operation occurred.
Fresh source-empty APItools/VCS cache resolves full SHA0c from the ordinary
Git/module origin, with exact version/time/both sums and all621 source ZIP/Git/
extracted bytes matching the qualified source. Ordinary rawZIP2333509 bytes,
SHA256012e818ca76501f84f3064dd9809ca1c4f666d26bb83f6bc2a6100a9a1593655, is
independently observed; no candidate seed/file proxy is used. Actual ordinary
owner/public/private graphs77/73/171 and complete compiled/source inventories
reproduce qualified inputs; all19 fresh events pass and all3 CLI bytes reproduce.
Affected owner/header and public/private16-case SDK races plus privateM52 unknown/
deadline controls pass. Exact source/module/consumer input parity permits reuse
of unchanged passing broad owner/SDK/private/checker evidence without repetition.
The [ordinary proof](../../docs/m84-ordinary-proof.json), [reuse proof](../../docs/m84-verification-reuse.json)
and [publication evidence](../../docs/m84-publication.md) retain exact gates.
Scoped final review/evidence publication follows the closing review; the root
coordinator owns shared-memory consolidation, Kinet:M50 reconciliation and
validated retirement. Terminal rows do not establish milestone acceptance.

## Whole review iteration2 — passed, 2026-10-09

Persisted before the final entire-milestone pass. Review both implementation
commits, all regressions/guarded invariants, exact source and complete actual
ordinary owner/public/private proof, preserved failed checks, exact proposal/
source before-and-after, metadata publication scope and downstream custody.
Resume interrupted iteration2; no P1/P2/higher carry-forward, maximum10 total.

The final entire-milestone pass and independent Tier0 read-only review found no
P1/P2/P3. Ordinary621-file source/Origin/version/time/sums, all19 event/log hashes,
exact publication range/proposal, complete normalized compiled-input parity,
unchanged consumer manifests, private/public boundaries and all3 reproducible
binaries independently revalidated. No source changed after0c; no broad suites
were repeated. All four task units and whole review2 pass. Scoped final evidence
publication is authorized; root acceptance/consolidation/reconciliation/retirement
remain separate required closure work, with no live/deployment authority.

## Coordinator closure — 2026-10-09

Root inspected the complete outgoing source diff, all four completed rows,
persisted whole review2 and ordinary source/consumer identity. No source
Go/go.mod/go.sum changed after 0c1383c6e2059e89cbf100ed8cb689be64d9ba57. Source/evidence main
publication is independently observed through b7425e2030489ce25a0ee523f83d694176f949f0; the exact
second push is retained in docs/m84-evidence-publication.json. Current product,
architecture, tooling and applicable lessons are consolidated. Kinet:M50 now
binds exact APItools source/version/sums alongside accepted M52, with real worker/
catalog/combined-graph/integration/bundle qualification still pending. Full
specification/status retention and envelope checks precede active-source removal.
No original history/corpus/wire/pin or evolution direction is changed. Coordinator
changes postdate the observed source/evidence head; no later publication is claimed.
``````
