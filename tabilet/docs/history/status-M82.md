# Retired milestone M82 - Shape production

**Milestone.** M82
**Outcome.** completed
**Retired.** 2026-10-06
**Source status.** tabilet/memory-bank/status-M82.md
**Source specification.** tabilet/memory-bank/milestone.md#m82--shape-production
**Evidence.** 54583f9b2f452b7cc522360c5aeeff29ca22f96c
**Worktree.** includes uncommitted changes
**Review.** passed
**Review iterations.** 3
**Verification.** Standalone/offline full APItools tests/vet, affected shape races, current workspace OpenUdon/Udon regression suites, retained-pin OpenUdon and disposable-modfile Udon qualification, source-family/identity/security/dialect/number/wire/resource/privacy/refusal fixtures, catalog freshness, gofmt, patch and exact upstream/publication checks passed. Six byte-identical pre-existing staticcheck observations are recorded without suppression or a zero-result claim. Accepted source independently observed on authorized origin/main; closure publication precedes downstream execution.
**Consolidated into.** [product](../../memory-bank/product.md), [architecture](../../memory-bank/architecture.md), [stack](../../memory-bank/tech-stack.md), [lessons](../../memory-bank/lessons.md), [contract](../../../docs/operation-shapes.md), [qualification](../../../docs/m82-qualification.md) and exact pending consumer records.

## Milestone specification

``````markdown
## M82 — Shape production

**Stage/owner.** STG-11 Phase A; APItools. **Priority.** Serial position 4/18, not a review severity.
**Dependencies.** [UWS:C09](../../../uws/tabilet/docs/history/status-C09.md); exact accepted/published contract closure recorded before adoption. Serial gates and direct contract/regression dependencies are reconciled in the coordinator.
**Scope.** Project supported source metadata; Preserve identity and incomplete evidence; Serialize and test shape tables; Qualify consumers and publish.
**Acceptance.** All existing source families produce honest, deterministic metadata with explicit incompleteness. The package remains source tooling, not a workflow or credential runtime.
**Verification.** go test ./...; go vet ./...; affected race/API/source-family fixtures and OpenUdon/Udon consumer checks; git diff --check. Use bounded local artifacts, not remote discovery.
Use only disposable roots and fixtures. Preserve published schemas/wires, historic evidence, current runtime capability restrictions and the installed M44 service. Changed v3/package/worker identities require fresh approval; they do not preserve old grants.
**Downstream.** [UWS:M08](../../../uws/tabilet/memory-bank/status-M08.md), [Kinet:M46](../../../kinet/tabilet/memory-bank/status-M46.md), [OpenUdon:P09](../../../openudon/tabilet/memory-bank/status-P09.md). Reconcile exact accepted/publication revisions before advancing.
**Tasks/review.** [status-M82.md](status-M82.md), all 4 task units complete; whole review passed 3/10. Accepted source `54583f9b2f452b7cc522360c5aeeff29ca22f96c` is independently observed on authorized origin/main as `v0.0.0-20261006210844-54583f9b2f45`. Qualification, fixed findings and preserved consumer/runtime boundaries are recorded in the status.

``````

## Status record

``````markdown
# M82 — Shape production

**Stage:** Kinet STG-11, Phase A. **Owner:** APItools.
**State:** All four tasks complete; whole review passed 3/10; accepted and independently observed published source `54583f9b2f452b7cc522360c5aeeff29ca22f96c`; reconciliation/retirement in progress.
**Source baseline:** `9364afac98c97b79c9e6b9335bd9c71d4e8722d5` (clean at planning).
**Coordinator:** [Stage 11 contract](../../../kinet/docs/stage11.md); the package-local milestone/status owns acceptance.

## Dependencies and handoff

[UWS:C09](../../../uws/tabilet/docs/history/status-C09.md).
The serial predecessor is a scheduling gate; direct contract and regression impacts are also listed. Every prerequisite must pass its whole review, and required publication must be independently verified before adoption. Record exact accepted/published sources and fixture/build hashes; no Stage 11 acceptance or future pin is claimed yet.

**Downstream:** [UWS:M08](../../../uws/tabilet/memory-bank/status-M08.md), [Kinet:M46](../../../kinet/tabilet/memory-bank/status-M46.md), [OpenUdon:P09](../../../openudon/tabilet/memory-bank/status-P09.md). Reconcile every affected consumer against the accepted prerequisite revision before advancing.

## Tasks

| Item | State | Notes |
|---|---|---|
| M82.1 — Project supported source metadata | `[+]` | Additive BuildOperationShapeTable projects all eight existing parser families, native selectors/ID aliases and raw SHA-256 into bounded UWS shapes; RPC/event/AWS/OData projections never invent generic HTTP details. Explicit partial schemas/security remain unqualified. Eight-family/public-API fixtures, standalone full tests/vet, focused race checks and workspace OpenUdon/Udon full regression suites passed. See execution evidence below and docs/operation-shapes.md. |
| M82.2 — Preserve identity and incomplete evidence | `[+]` | VerifyOperationShapeTable independently reproduces exact local claims and refuses forged/stale fields, omitted operations and changed raw bytes. Explicit OpenAPI security preserves OR-of-AND symbols/scopes; undeclared/unresolved security remains unknown. Complete supported HTTP metadata proves advisory compatibility, while unsupported native/dialect/wire/requiredness/schema evidence remains partial. Identity, symbolic-security, type/refusal, canary and no-network fixtures passed with standalone full tests/vet, focused races and workspace OpenUdon/Udon regressions. |
| M82.3 — Serialize and test shape tables | `[+]` | Reuse public UWS Marshal/ParseTable and freeze eight-source/twelve-operation bytes at testdata/operation-shapes/v1/eight-families.json. Source-order, exact JSON/YAML numeric lexemes, resolver isolation, bounded wire, missing/cyclic refs, native collision/overload, ambiguous IDs and streaming/list/message-alternative fixtures passed. OData/AsyncAPI summary deduplication cannot select a contract silently. Legacy candidate wire fixtures are unchanged; standalone full tests/vet, focused races, dependency closure and patch checks passed. Retained Horizon/HCL remains explicit; no credential resolution/execution/network behavior is added. |
| M82.4 — Qualify consumers and publish | `[+]` | Owner and exact-source/retained-consumer checks passed, with value-free/privacy/identity/resource/dialect/wire fixtures and whole review 3/10. Accepted source 54583f9b2f452b7cc522360c5aeeff29ca22f96c was pushed normally to the exact authorized APItools origin/main and independently observed with git ls-remote. Go registry resolved v0.0.0-20261006210844-54583f9b2f45 to that full source hash. Final qualification records scope and preserved consumer pins; no live workflow or deployment authority is supplied. |

## Acceptance and verification

All existing source families produce honest, deterministic metadata with explicit incompleteness. The package remains source tooling, not a workflow or credential runtime.

go test ./...; go vet ./...; affected race/API/source-family fixtures and OpenUdon/Udon consumer checks; git diff --check. Use bounded local artifacts, not remote discovery.
Use only disposable roots and fixtures. Preserve published schemas/wires, historic evidence, current runtime capability restrictions and the installed M44 service. Changed v3/package/worker identities require fresh approval; they do not preserve old grants.

## Execution policy

One execution owner, serial execution and task commits under the later confirmed goal. Planning authorizes no code execution, commit, publication or external operation. Source publication requires separately named authority; a status marker or local build is not publication. Consumers must record exact accepted and published prerequisites before adoption. Default checks are offline, credential-free and model-free. No deployment, live ledger migration, real API/model/mail action or registration change.

## Approved review-intake amendments — 2026-10-06

Source: **Stage 11 planning review — cross-package refactoring** (2026-10-06). Review baseline and full local revalidation HEAD: `9364afac98c97b79c9e6b9335bd9c71d4e8722d5`; relevant uncommitted Stage 11 planning changes were included. The review covers the five owner baselines recorded in the coordinator. This is approved intake, not a closing review iteration; the persisted counter remains 0/10.

- **P3-8** — source P3; local Lower; partially confirmed. Evidence: go.mod; ../uws/go.mod; status-M82.md M82.3. Ownership/lineage: M82.3 records the planned transitive closure; Udon:M48.2 verifies it before adoption.

## Accepted binding prerequisite — 2026-10-06

UWS:C09 is accepted and independently observed on origin/main at `6a267306032edc687a298cefc8bba7019d3ad059`, whole review 3/10. [Public contract](../../../uws/docs/binding-reference.md) and [qualification](../../../uws/docs/c09-qualification.md) pin source-neutral ShapeTable/Resolver APIs, metadata-only binding checks and deterministic flow observations. Known/unknown evidence remains explicit; nested/typed templates use exact projections and unproved constraints stay indeterminate. Metadata producer claims are independently reproduced/verified before authority. No APItools/private-runtime import, source parser, credential/provider I/O, ordinary-validator/schema change or execution permission is supplied by C09. Retirement closure `8e5be730aa68aa4cb4f6c591a3a9425c6b8bd55c` was independently observed on authorized origin/main, satisfying the prerequisite publication gate. [Publication evidence](../../../uws/docs/c09-publication.md) records accepted-source ancestry.

## Persisted review

- Review iteration: **3/10**; whole milestone review passed on 2026-10-06 after all fixes, required owner/consumer verification, terminal task/publication evidence and final closure audit. No remaining P1/P2.
- Closing-review findings: review 1 M82-R1-1/R1-2 and review 2 M82-R2-1 are resolved below; review 3 found no remaining P1/P2. Six unchanged pre-existing staticcheck observations do not change the required gate or producer behavior.
- Accepted revision: `54583f9b2f452b7cc522360c5aeeff29ca22f96c`.
- Published revision / artifact evidence: exact accepted source independently observed on authorized origin/main; `v0.0.0-20261006210844-54583f9b2f45` registry origin hash matches. Fixture dc20d2287322a4b20d5d92b1a0d1ec8036f1bec853b642d7df8797ed096c3ba1 remains unchanged; closure publication follows retirement.
- Verification: required standalone full tests/vet, affected shape races, workspace consumer suites, retained-pin/disposable-modfile qualification, catalog freshness, format/patch and exact upstream/publication checks passed; six byte-identical baseline staticcheck observations are recorded without a zero-diagnostic claim.

After all tasks finish, perform the whole-milestone review with persisted iteration/finding state and fix every P1/P2 before acceptance. Resume an interrupted pass at the same counter. Consolidate current facts, reconcile downstream work and retire under this package’s normal procedure.

## Execution evidence — M82.1, 2026-10-06

Resumed the previously selected row after the old execution session's token
revocation; its saved dependency/status changes were retained. No initial
planning commit or completed prerequisite was replayed. Exact accepted UWS C09
`v0.0.0-20261006181058-6a267306032e` is now a direct dependency; standalone
`GOWORK=off go mod tidy` recorded its existing Horizon/HCL transitive closure.
The producer reads only explicit local bytes/files and returns no partial table
on source, identity, resource or cancellation errors. Numeric/schema projections
remain independent of prompt-sanitized summaries. No old API/wire change,
credential/provider operation or workflow execution is added.

Passed with Go1.26.6: `GOWORK=off GOPROXY=off go test ./...`,
`GOWORK=off GOPROXY=off go vet ./...`,
`GOWORK=off GOPROXY=off go test -race . -run '^TestOperationShapes'`,
and `GOPROXY=off go test ./...` in both OpenUdon and Udon using the existing
operator-owned workspace. `git diff --check` and focused gofmt checks passed.
M82.2–M82.4, the closing review, consumer acceptance and publication remain
pending; no M82 acceptance or new consumer pin is claimed.

## Execution evidence — M82.2, 2026-10-06

Independent reproduction compares exact deterministic UWS tables, not only raw
digests or structurally valid producer claims. Security symbols retain explicit
anonymous alternatives and unknown missing definitions. JSON Schema compilation
uses the existing exact `jsonschema/v6 v6.0.1` pin with an external loader that
always refuses; native parser losses cannot claim known schemas. Full native
RPC/GraphQL/OData aliases and protobuf field hints remain metadata only. CSDL
nullability does not establish call-argument requiredness; incomplete operation
evidence cannot establish absence or supported transport/auth.

Passed: standalone offline full APItools tests and vet; focused shape races;
workspace offline full OpenUdon and Udon regressions; final focused tests/races
and vet after stricter requiredness/Swagger-server checks; `git diff --check`.
Supported compatibility and mismatches were checked through the public UWS
binding consumer, alongside structurally valid forgeries/stale raw identity,
OR/AND requirements, unknown definitions/schema/dialect and no-network canaries.
The remaining wire/source-family edge fixtures and whole consumer/publication
gate stay with M82.3/M82.4. No new acceptance or publication is claimed.

## Execution evidence — M82.3, 2026-10-06

The independently inspected synthetic golden contains eight exact raw source
identities and twelve native operations, 9,829 bytes with SHA-256
`dc20d2287322a4b20d5d92b1a0d1ec8036f1bec853b642d7df8797ed096c3ba1`.
All prior candidate v1 wire files and source artifacts are unchanged. Exact
large-integer, decimal and exponent lexemes survive JSON/YAML schema projection;
no float64/JCS equality oracle is used. Missing/cyclic schema references remain
unknown; duplicate/trailing/alias encodings and colliding native contracts
refuse. Unique OpenAPI refs retain ambiguous duplicate ID aliases. Message,
stream/list and multiple response/media evidence stays incomplete. The 8 MiB
serialized-table limit and aggregate operation bound return no partial table.

Passed: `GOWORK=off GOPROXY=off go test ./...`, `go vet ./...` under the same
standalone/offline environment, focused shape races, `go list -m all`, gofmt
and `git diff --check`. UWS remains pinned to accepted/published C09; Horizon
`v1.14.5` / HCL `v2.24.0` remain in the module closure. M82.4 owns the final
consumer qualification, publication and whole review; review remains 0/10.

## Execution evidence — M82.4 continuation, 2026-10-06

M82.4 is the sole in-progress row. [Draft qualification](../../docs/m82-qualification.md)
records exact source/fixture/upstream reachability and consumer evidence.
Prepublication checks exposed overly strong completeness for malformed HTTP
path parameters and unproved schema dialect/reference-sibling/exclusive-bound
semantics. The new negative fixtures reproduced these cases; M82.4 guards now
refuse malformed HTTP metadata or keep schema evidence unknown. Final standalone
full tests/vet, focused races, catalog freshness and patch/format checks passed.
OpenUdon's standalone retained-pin suite passed. Udon's original readonly
standalone check needs the UWS metadata update; its disposable-modfile suite
passed with only the exact accepted C09 require changed, without tracked consumer
changes. The source publication destination and fast-forward baseline are
verified; no M82 push occurred. Closing review remains 0/10, not started.
Compatible retained staticcheck reported six diagnostics, all in four files
byte-identical to pre-M82; none concerns the new shape implementation. No zero
staticcheck result is claimed, and no unrelated baseline cleanup or suppression
was introduced. The original Udon modfile/sum were preserved and the temporary
qualification copies removed.
Resume this row through final review, publication, reconciliation and retirement;
do not restart M82.1–M82.3 or accepted prerequisites.

## Closing review 1 — findings resolved, 2026-10-06

Full scope: producer/verifier boundaries, every native adapter, schema/security
projection, source identity/refusal, determinism/resource/cancellation behavior,
legacy API/wire preservation, consumer closure and documentation/publication.

- **M82-R1-1 (P2):** Serialized-table and aggregate-operation limits are checked
  after a whole source's shapes have been accumulated. Compact documents can
  reuse large local schemas across many operations/parameters, multiplying
  projection/compilation work and resident metadata before the final 8 MiB
  refusal. Per-schema work/byte limits do not bound their aggregate expansion.
  Own in M82 before acceptance: charge schema/projection work and serialized
  operations incrementally against one invocation budget; stop at the first
  exceeded bound with no partial table. Add a compact shared-ref amplification
  regression. Worker hard CPU/RSS/deadline isolation remains separately owned.
- **M82-R1-2 (P2):** The schema projector recognizes modern JSON Schema
  keywords without tying them to the source's schema dialect. An OpenAPI 3.0
  or default draft-07 OpenRPC/AsyncAPI schema could incorrectly claim known
  semantics for an unsupported modern keyword. Own in M82: retain useful known
  common constraints, but mark unsupported dialect keywords unknown, including
  within resolved schemas; only the supported explicit/default 2020-12 context
  can prove modern semantics. Add dialect regressions before acceptance.

Resolution: one invocation now charges projected schema bytes/work and validates/
accounts source/operation serialization incrementally; the claim verifier checks
amplification before marshaling. Compact shared-ref and decoded-claim amplification
vectors refuse without a partial table. Source-specific dialect knowledge keeps
modern constraints unknown outside supported contexts and retains known common/
explicit-2020 semantics. The unchanged golden and all focused tests passed;
standalone full APItools tests/vet, shape races and full workspace OpenUdon/Udon
regressions passed after both fixes. Both findings are resolved in the current
uncommitted M82.4 implementation. No existing API/wire/source fixture changed.

## Closing review 2 — finding resolved, 2026-10-06

Review the full implementation, source/parser/worker boundaries, projection/
verification/cancellation/resource failures, fixtures, metadata/privacy, legacy
contracts, exact dependency closure, consumer results and public documentation.
Source publication, acceptance and retirement remain pending this gate.

- **M82-R2-1 (P2):** The review-1 dialect guard covers unsupported keyword
  names but not boolean schema values or union/null `type` values in OpenAPI
  3.0. The neutral JSON Schema compiler accepts them, so this unsupported source
  evidence can still claim Known. Own in M82: keep these projections unknown
  under the OpenAPI 3.0 profile, including nested/ref-resolved schemas, while
  retaining valid draft-07/2020 schemas. Add explicit positive/negative vectors.

Resolution: boolean and union/null schemas cannot claim known OpenAPI 3.0
semantics; nested and resolved projections inherit the same profile. Supported
OpenAPI 3.1/draft-07/explicit-2020 projections retain their positive proofs.
Focused shape tests/races and full vet passed; the frozen wire is unchanged.

## Closing review 3 — whole code gate passed, 2026-10-06

Re-review the complete M82 change and tests, source/native/parser boundaries,
all schema dialect/ref/number/security cases, identity/serialization/claims,
resource/cancellation/privacy failures, unchanged APIs/wires/dependencies,
consumer qualification, documentation and the named publication/retirement gates.

Result: no remaining P1/P2. The whole producer/verifier implementation, all native
families and negative/positive fixtures, source-specific schema evidence,
incremental bounds, exact numeric/wire identity, value-free failures, symbolic
security, privacy/cancellation and unchanged legacy behavior satisfy M82's code
acceptance. Required full APItools tests/vet and affected races passed; current
workspace OpenUdon/Udon suites passed, with retained-pin/disposable-module
qualification recorded. Evolution decision: no bump; this implements the already
approved Stage 11 target. Publication and the final acceptance/retirement evidence
remain the selected M82.4 operation, not authority supplied by these markers.
The reviewed source must be committed before the named fast-forward publication;
the source commit is the required artifact for that operation. Its exact observed
hash/publication and substantive closure evidence follow in this same task.

Final closure audit: all four required task rows are terminal and every M82
acceptance/verification/publication requirement is satisfied. Normal publication
of accepted source `54583f9b2f452b7cc522360c5aeeff29ca22f96c` to the unchanged
authorized origin/main succeeded; an independent ls-remote returned that exact
hash and the configured Go registry resolved the same origin revision. Downstream
UWS:M08/Kinet:M46/OpenUdon:P09 consume this exact producer contract, its limitations
and independently reproduced claims; consumer pin changes still belong to their
own task units. Closure records will be published before any downstream execution.
``````
