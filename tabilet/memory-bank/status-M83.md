# M83 — Stage 11 native shape projection remediation

**Stage:** STG-11 post-acceptance remediation. **Owner:** APItools.
**State:** Authorized serial implementation started, 2026-10-08. M83.1 is complete; six rows remain pending and publication remains a separate gate.
**Authority:** The complete review reconciliation was approved, followed by the confirmed serial GOAL request on 2026-10-08: Udon:M49 → UWS:M09 → APItools:M83 → OpenUdon:M99 → Kinet:M49, COMMIT_POLICY: task and EXTERNAL_MUTATIONS: none. This authorizes scoped implementation after accepted/published prerequisites; source publication still requires this owner's separate fresh named grant. The Udon/UWS publication exceptions do not extend to this owner or live operations.
**Review source:** stage11-siblings-review.md — Stage 11 code review — sibling packages; APItools section.
**Review baseline/range:** `9364afac98c97b79c9e6b9335bd9c71d4e8722d5` → `24c36bf40102c2c1d160dc7d0e27fb161e12dbd6`.
**Revalidation HEAD:** `24c36bf40102c2c1d160dc7d0e27fb161e12dbd6`; clean worktree, no relevant uncommitted code in the evidence. Approved planning changes are not implementation evidence.
**Lineage:** [M82](../docs/history/status-M82.md). Existing acceptance, review counters, statuses and frozen evidence stay preserved.
**Coordinator:** [Stage 11](../../../kinet/docs/stage11.md#post-acceptance-remediation--2026-10-08); package-local specification and status own acceptance.

## Dependencies and handoff

Accepted M82 metadata contract plus accepted and independently published UWS:M09 root contract before exact adoption. Serial scheduling follows UWS:M09. No sibling parser changes or runtime-capability expansion is authorized.

**Exact successor acceptance/publication/build identities:** unset; record full independently observed revisions and hashes during the later execution. Never substitute local HEAD, directory replacements or prior consumed publication authority.

**Downstream:** OpenUdon:M99 independent source/shape reproduction and Kinet:M49 author worker/package qualification. Corrected table digests require fresh source-backed assessments/publication decisions; old tables/packages are historical evidence, not silently rewritten.

## Scope and acceptance

Correct native-family direction/default/member/request/response projections and preserve honest OpenAPI completeness/schema proof. Restore supported OpenAPI YAML/extension/Swagger/TRACE metadata handling while retaining source identity, selector and symbolic security contracts.

AsyncAPI 2.x publish consumes/output and subscribe produces/input consistently with v3 receive/send. GraphQL defaulted non-null arguments/variables are optional and response aliases remain distinct response keys with stable underlying field/selectors. Smithy static URI query literals are not fabricated required members. Discovery emits only independently declared request/endpoint evidence, leaving unproved requiredness/projection unknown. Unsupported serialization/formats cannot be marked complete/known. Valid numeric YAML response keys and path extensions do not discard a whole source; Swagger version detection is consistent and declared TRACE metadata is retained without widening runtime support. Exact bounded raw-source identities, duplicates, native selectors, security OR/AND and no-network guarantees survive all eight-family qualification and review.

## Tasks

| Item | State | Notes |
|---|---|---|
| M83.1 — Correct AsyncAPI version-dependent message direction | `[+]` | Map 2.x publish to consumed/output and subscribe to produced/input, corresponding to v3 receive/send. Retain partial native security/schema evidence and stable selectors; test equivalent 2/3 declarations. Source A1. |
| M83.2 — Preserve GraphQL defaults and response keys | `[+]` | Use default-value presence for non-null argument/variable requiredness. Add response-key metadata where needed in the owned GraphQL parser/model without changing underlying field names/selectors; preserve aliases and distinguish genuine duplicate response keys. No GraphQL transport capability is added. Sources A2/P3.7. |
| M83.3 — Separate Smithy static query literals from modeled inputs | `[+]` | Use raw/native member provenance to omit static URI query literals with no modeled input member. Do not filter only on an empty MemberName, since synthetic names are populated. Preserve real query members, including same-name cases, and source-native protocols/selectors. Source A3. |
| M83.4 — Correct Discovery request and endpoint projection | `[ ]` | Preserve declared method/request/normal versus upload endpoint evidence without inventing HTTP contracts or mandatory bodies. Use indeterminate/incomplete evidence when native/raw source cannot prove requiredness or distinguish a supported variant. Keep declared server/source metadata and native selectors. Source P3.6. |
| M83.5 — Keep OpenAPI completeness and format claims honest | `[ ]` | Account for explode-only, Swagger collectionFormat and non-JSON content/encoding loss. Mark unsupported serialization incomplete and unenforced native formats unknown unless constraints are actually proved; do not silently widen schemas or add runtime serialization promises. Sources P3.1/P3.2. |
| M83.6 — Restore bounded OpenAPI and Swagger metadata compatibility | `[ ]` | Normalize only numeric YAML response-code keys with duplicate checks; ignore legal path-level vendor annotations of any type without globally coercing mapping keys. Use consistent Swagger 2.0 detection for schemas/security/servers, retain declared TRACE metadata and remove dead method case logic without enabling TRACE runtime execution. Sources P3.3/P3.4/P3.5. |
| M83.7 — Qualify all native families and publish the exact handoff | `[ ]` | Adopt exact accepted/published UWS:M09 root, preserving conservative schema/dialect/path evidence, then run eight-family source-backed shape/identity/security/no-network compatibility and standalone owner/consumer checks, preserve old qualification evidence and add corrected successor fixtures. Persist review counters and record accepted/published source/module sums plus OpenUdon/Kinet impacts; publication needs fresh named authority. |

## Active finding provenance

Only approved active findings are recorded here. Source priorities and local severity are separate; task references identify one owner for each required outcome. Unsupported and unscheduled findings remain in the conversational handoff; optional directions belong only in milestone.md.

| Source finding | Source priority | Local severity | Disposition | Current repository evidence | Task owner |
|---|---|---|---|---|---|
| A1 | P2 | P2 | confirmed | operation_shapes_native.go:241–244; offline AsyncAPI 2.6 publish projection yields required input | M83.1 |
| A2 | P2 | P2 | confirmed | operation_shapes_native.go:298,313; first:Int!=10 probe remains required | M83.2 |
| A3 | P2 | P2 | confirmed | operation_shapes_native.go:166; static URI ?x-id literal probe becomes required input | M83.3 |
| P3.1 | P3 | P2 | confirmed | operation_shapes_openapi.go serialization/completeness; operation_shapes_schema.go content projection | M83.5 |
| P3.2 | P3 | P2 | confirmed | operation_shapes_schema.go:42 marks unenforced OpenAPI formats known; out-of-range int32 probe | M83.5 |
| P3.3 | P3 | P2 | confirmed | operation_shapes_openapi.go:31–36 omits trace; uppercase check is unreachable | M83.6 |
| P3.4 | P3 | P2 | partially confirmed | shape YAML scalar/map and path-object guards reject numeric response keys/non-object x- extensions | M83.6 |
| P3.5 | P3 | Lower | confirmed; OpenAPI compatibility acceptance | operation_shapes_security.go:27,76 and schema server projection compare Swagger numeric version to string | M83.6 |
| P3.6 | P3 | P2 | partially confirmed | operation_shapes_native.go:113,129 uses native upload path and always-required request body | M83.4 |
| P3.7 | P3 | P2 | confirmed | operation_shapes_native.go:318 and graphql/parse.go selection projection lose response aliases | M83.2 |

## Verification and compatibility

go test ./...; go vet ./...; focused OperationShapes/GraphQL/sourceguard tests and races; existing compatible checker policy with no new waiver; GOWORK=off GOPROXY=off standalone public source/UWS module checks; defaulted arguments/variables, aliased selections, AsyncAPI 2/3 equivalent direction, Smithy literals versus genuine same-name members, Discovery request/upload variants, serialization/format unknowns, numeric response keys/path extensions/Swagger/TRACE fixtures; eight-family exact source/selector/security/digest and no-network cases; git diff --check.

Use retained Go 1.26.6 and exact ordinary modules without ambient workspace substitution. Default verification is offline, credential-free and model-free. No live user ledger, host, provider, account, mail, Cloudflare, registration or consumer adoption operation is included.

Public schemas/wires, published grammar/version bytes, accepted historical qualification and independently retained browser/media/legacy/frozen-consumer pins stay preserved. Corrected derived metadata and new worker/package identities require fresh consumer assessment and explicit authority; historical approvals are never upgraded automatically. Source publication requires a new separately named request and independent resolution before downstream adoption. Planning rows may remain pending on this external prerequisite; none is started here.

## Closing review

**Review iterations:** 0/10.
**Review state:** not started; this is review intake, not a pass of an existing gate.
**Findings/fixes:** no implementation or fix verification claimed.
**Execution owner:** one serial owner across the five ledgers; no row is in progress.
**Commit policy:** The user separately authorized a planning commit on 2026-10-08 with “git commit and then report the index refresh issue in ~/skill-index.md”. This authorizes one commit of the approved planning changes in this owner repository; implementation, publication and deployment remain outside this request. Future task commits follow the separately invoked GOAL/request policy.
**Closure:** persist each started review iteration before reviewing; resume an interrupted pass at the same number. No open P1/P2 may remain at acceptance. Required verification, exact downstream reconciliation and owner-specific consolidation/retirement follow implementation; never reopen completed Stage 11 history.

## Accepted UWS:M09 prerequisite — 2026-10-08

All five UWS rows and whole review 7/10 pass; accepted root/codec runtime is
b099f6803277ae94c7e9f1da0904a0140b278f20, independently resolved as
v0.0.0-20261008043726-b099f6803277. Published evidence
51a74545b016b8ab75e30d454338c52f7945a836 is independently observed. Root sum is
h1:4xy+/HBNh1CSJDO+qOzWdV/0zC/yFCKAz2kOBWufA7g=; codec sum is
h1:OJsmDK/RcFpyMAMGy84DjcX6kNalYH/E0Q/C3XwD5Uo=. GoMod sums are
h1:DlqFOnO9lbmYWLLIh5WicNX6NTWIuytU6mIHmxj9BVw= (root) and
h1:0cR/xLzEP8vJ9FAUhsLbaKVkU7UarXPc51nJjEaMP5Q= (codec).
[Ordinary publication proof](../../../uws/docs/m09-publication.md) records
full provenance, exact 360/21 files and complete 52/50/52 selected closures.
These counts describe UWS proofs, not a prescribed downstream closure size.
Normal M09 retirement closure 989e3f2c88cac5c0f5a2911dfe04c36a61e43126
is independently observed on the approved origin/main before adoption. Runtime
b099f6803277ae94c7e9f1da0904a0140b278f20 is its ancestor; UWS worktree is clean.

Corrected binding proofs retain containing constraints, original dialect and
indeterminate outcomes. Flow and strict portability use actual root goto,
trigger/dependency iteration contexts and separate step/operation output owners.
Untrusted HCL/shape parsing refuses depth above 100 before recursive decoding;
HCL views require deterministic canonical bytes plus independent value/numeric
proof. Non-NFC presentation remains fail-closed. Ordinary validation/execution,
wire/schema/digest algorithms, published versions/corpora and frozen consumer
pins stay preserved. Changed diagnostics/derived assessments/package or worker
identities require fresh assessments/confirmations/grants. This prerequisite
note changes no implementation row to complete and supplies no live authority.

## M83.1 implementation evidence

AsyncAPI 2.x publish emits an output and subscribe a required payload input, matching v3 receive/send. Equivalent 2/3 raw fixtures independently reproduce selector and partial security/completeness evidence. Offline retained Go 1.26.6 `GOWORK=off GOPROXY=off GOTOOLCHAIN=local go test . -run '^TestOperationShapes'` passes. The accepted v3 eight-family corpus bytes stay unchanged.

## M83.2 implementation evidence

Owned GraphQL Argument/Variable metadata adds explicit HasDefault presence; nullable introspection defaultValue is absence, while declared SDL/operation defaults make non-null inputs optional without changing type nullability. Selection metadata retains underlying FieldName and ResponseKey per occurrence, alongside the existing deduplicated underlying SelectionNames. Shapes emit distinct response keys and refuse real duplicate keys; native field/operation selectors remain stable. Default collection stops before the next variable marker. Offline focused OperationShapes and all GraphQL parser/default tests pass; the frozen candidate/shape corpus remains unchanged. Exported additions will receive module-only sibling checks in M83.7.

## M83.3 implementation evidence

Smithy shape projection checks raw input-member httpQuery provenance before emitting a native binding that shares a URI query-literal wire name. Synthetic names are not evidence; genuine same-name query members remain, and same-name body members cannot manufacture query inputs. Literal-only, body/query same-name and differently named wire-member regressions pass with all offline OperationShapes tests. Source selectors/protocol and partial-schema policy remain intact.
