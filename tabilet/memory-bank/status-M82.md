# M82 — Shape production

**Stage:** Kinet STG-11, Phase A. **Owner:** APItools.
**State:** Confirmed Stage 11 execution; M82.1/M82.2 complete, M82.3–M82.4 pending; review 0/10 not started.
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
| M82.3 — Serialize and test shape tables | `[ ]` | Emit deterministic UWS ShapeTable bytes from bounded local inputs. Exercise every source family, malformed references, collisions, security alternatives and forged/stale identity; no credential resolution, execution or implicit network access. Record the planned UWS module dependency and its retained Horizon/HashiCorp HCL transitive closure; do not claim APItools or UWS core becomes HCL-free in Stage 11. |
| M82.4 — Qualify consumers and publish | `[ ]` | Run owner and affected OpenUdon/Udon checks against exact UWS source, record the public API/fixture release, and publish only under named authority before adoption. |

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

- Review iteration: **0/10**; not started.
- Closing-review findings: none; the whole-milestone review has not started. Approved intake requirements above remain pending.
- Accepted revision: not available.
- Published revision / artifact evidence: not available.
- Verification: whole-milestone qualification remains pending; completed task verification is recorded below.

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
