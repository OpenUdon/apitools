# M82 — Shape production

**Stage:** Kinet STG-11, Phase A. **Owner:** APItools.
**State:** Approved planning on 2026-10-06; 4 pending rows, no implementation or acceptance.
**Source baseline:** `9364afac98c97b79c9e6b9335bd9c71d4e8722d5` (clean at planning).
**Coordinator:** [Stage 11 contract](../../../kinet/docs/stage11.md); the package-local milestone/status owns acceptance.

## Dependencies and handoff

[UWS:C09](../../../uws/tabilet/memory-bank/status-C09.md).
The serial predecessor is a scheduling gate; direct contract and regression impacts are also listed. Every prerequisite must pass its whole review, and required publication must be independently verified before adoption. Record exact accepted/published sources and fixture/build hashes; no Stage 11 acceptance or future pin is claimed yet.

**Downstream:** [UWS:M08](../../../uws/tabilet/memory-bank/status-M08.md), [Kinet:M46](../../../kinet/tabilet/memory-bank/status-M46.md), [OpenUdon:P09](../../../openudon/tabilet/memory-bank/status-P09.md). Reconcile every affected consumer against the accepted prerequisite revision before advancing.

## Tasks

| Item | State | Notes |
|---|---|---|
| M82.1 — Project supported source metadata | `[ ]` | Use existing parsers for OpenAPI/Swagger, Google Discovery, AWS Smithy, AsyncAPI, GraphQL, OpenRPC, gRPC/protobuf and OData. Preserve source-native selectors and protocol-specific limits rather than forcing every source into HTTP. |
| M82.2 — Preserve identity and incomplete evidence | `[ ]` | Bind raw source digests and symbolic security alternatives without flattening OR/AND requirements. Preserve unknown schemas, effects and capability limitations; URLs are provenance, not permission to fetch. |
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

## Persisted review

- Review iteration: **0/10**; not started.
- Closing-review findings: none; the whole-milestone review has not started. Approved intake requirements above remain pending.
- Accepted revision: not available.
- Published revision / artifact evidence: not available.
- Verification: pending implementation; no test result is claimed by this planning record.

After all tasks finish, perform the whole-milestone review with persisted iteration/finding state and fix every P1/P2 before acceptance. Resume an interrupted pass at the same counter. Consolidate current facts, reconcile downstream work and retire under this package’s normal procedure.
