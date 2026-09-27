# Retired milestone M77 - Step-contract operation metadata

**Milestone.** M77
**Outcome.** completed
**Retired.** 2026-09-27
**Source status.** tabilet/memory-bank/status-M77.md
**Source specification.** tabilet/memory-bank/milestone.md#m77---step-contract-operation-metadata
**Evidence.** e3b4b6ec343a18c48fa93a971a69930b993203db
**Worktree.** includes uncommitted changes
**Review.** passed
**Review iterations.** 7
**Verification.** APItools `go test -count=1 ./...`, `go vet ./...`, `GOWORK=off go test -count=1 ./...`, `GOWORK=off go vet ./...`, `go test -race -count=1 .`, and `git diff --check` passed. OpenUdon workspace and standalone test/vet passed; its boundary and doc-memory checks passed. Udon workspace test/vet and standalone test/vet passed, with a temporary tidy-equivalent module-file copy for standalone tests. The OpenUdon workspace suites were rerun sequentially because concurrent suites share a temporary path; sequential checks passed. Ramen was excluded by owner instruction.
**Consolidated into.** [product](../../memory-bank/product.md), [architecture](../../memory-bank/architecture.md), [tech stack](../../memory-bank/tech-stack.md), [evolution v24](../../evolution/result-v24.md), and the [operation-candidate contract](../../../docs/operation-candidates.md); no separate reusable lesson was warranted.

## Milestone specification

````markdown
## M77 - Step-contract operation metadata

**Goal.** Give OpenUdon's non-interactive `step candidates` consumer readable
operation summaries, evidence-backed `read`/`write`/`unknown` effects, and
deterministic ranking by step purpose, available inputs, and expected outputs.
This is APItools' part of sibling item S2a in
[Kinet's design, section 7](../../../kinet/docs/icot.md#7-proposed-changes-by-package).
It is new work building on existing inventories and ranking; no completed
milestone is reopened and no existing candidate direction is promoted.

**Direction.** [Evolution v24](../evolution/result-v24.md) records the approved
future public metadata contract. Implementation and acceptance remain pending.

**Scope and contract.**

- Add public Go types and versioned conformance fixtures for operation/source
  identity, source kind, native selector, content digest, consumer summary,
  effect evidence, and a bounded step contract containing purpose, typed
  available inputs, and expected outputs. Preserve existing exported APIs,
  data shapes, and legacy scoring behavior through an additive interface.
- Produce concise, source-grounded descriptions of what an operation does,
  requires, and returns. Distinguish extracted facts, inferred wording, and
  missing evidence. Do not invent accounts, recipients, capabilities, fields,
  or operation identifiers. The first version is deterministic and offline;
  no model service or provider lookup is part of summary generation.
- Classify effects from operation meaning and available protocol/source
  evidence, returning the class, evidence provenance, and concise reasons.
  HTTP method alone cannot establish `read`. Missing, insufficient, or
  conflicting evidence remains `unknown`; it cannot satisfy a confirmed
  read-only requirement. Existing lifecycle purposes are not execution-safety
  evidence. Inferred classification never becomes user confirmation or approval.
- Rank against purpose, required request inputs, and expected response outputs
  with separate compatibility evidence and score explanations. Distinguish
  compatible, incompatible, and indeterminate relationships; missing schema
  evidence is not a positive match. Expose absent inputs/outputs, type conflicts,
  ties, no matches, and truncation. Natural-language similarity is advisory and
  cannot prove that a workflow's intended outcome will be achieved.
- Integrate OpenAPI/Swagger inventories and adapt existing native metadata for
  Google Discovery, AWS Smithy, AsyncAPI, GraphQL, OpenRPC, gRPC/protobuf, and
  OData. Reuse existing parsers, preserve native kinds and selectors, and
  document per-family evidence and capability limits. Unsupported details
  produce explicit diagnostics rather than silently disappearing or becoming
  successful compatibility claims. This does not add protocol parsers or
  lower native source documents into OpenAPI.
- Preserve exact source identity and authentication OR-of-AND alternatives.
  Apply existing prompt-safety, collection, byte, and work budgets to the new
  requests, summaries, evidence, and ranked reports. Report compaction and
  fail visibly if it would change identity, field, security, or effect meaning.
  Public helpers return independent data and do not mutate caller metadata.

**Ownership and order.** Lane M owns this shared public consumer contract.
M77.1 defines it and its fixture shapes; M77.2 and M77.3 provide summaries and
effects; M77.4 ranks contracts; M77.5 integrates the source paths; M77.6
qualifies the full surface and consumer handoff. One execution owner and at
most one general row may be in progress. Ownership covers APItools metadata,
summary/classification/ranking adapters, tests, fixtures, and documentation;
catalog curation, remote discovery policy, runtime behavior, and sibling
implementation are outside this milestone.

**Dependencies and downstream effects.** Existing inventory/prompt hardening
(S01, M74, S04) and lifecycle ranking (M75, S03) provide completed baseline
contracts, resolved through the history index. There is no active upstream
implementation prerequisite. OpenUdon's pending M87.1 drafts its own command
contract independently; reconcile the metadata handoff against that work
without making either contract-design row wait for the other's implementation.
OpenUdon M87.4 production integration and M87.6 acceptance require a compatible
published APItools revision. APItools qualifies its public API and fixtures
without depending on Kinet W03 delivery. OpenUdon and Udon remain compatibility
consumers; Ramen is a separate project and is outside this milestone's
downstream verification. Publication, dependency re-pinning, and changes to
sibling ledgers or code require their own authority and are not performed by
this planning approval.

Kinet retains the planning loop and confirmed step decisions. OpenUdon retains
command JSON, step binding/checking, account/destination constraints, approval,
and package behavior. UWS and Udon retain workflow effect semantics and runtime
enforcement. Effect override authority, browser actions, simulation, and iCoT
retirement are not APItools deliverables or prerequisites for M77.

**Acceptance and verification.** Public fixtures demonstrate useful
source-grounded summaries and actual read/write/unknown classifications, with
reasons, rather than a permanently unknown substitute for classification.
Rankings change appropriately when required inputs or expected outputs change
while purpose text stays constant. Negative cases cover misleading methods,
conflicting evidence, missing schemas, type/requiredness conflicts, unsupported
source details, duplicate IDs across sources, ambiguity, no match, hostile
descriptions, budget exhaustion, and caller-data preservation. Results retain
source identity and auth alternatives with deterministic ordering. Preserve
existing inventory, selection, authoring-context, and lifecycle regressions.

Run focused metadata/adapter/conformance tests, affected race checks,
`go test ./...`, `go vet ./...`, `GOWORK=off go test ./...`,
`GOWORK=off go vet ./...`, and `git diff --check`. Run available OpenUdon and
Udon consumer suites in workspace and standalone modes, distinguishing tests
against the local API from standalone tests of existing pinned versions.
New public API examples and fixtures must compile and run against the local
implementation. Keep default verification credential-free and provider-free.
Document the adoption contract and evidence required for later published-pin
verification. Complete the persisted maximum-ten-iteration review, current-memory
consolidation, and downstream reconciliation before normal milestone retirement.
````

## Status record

````markdown
# Status M77 - Step-contract operation metadata

**State:** Complete. Planning approved 2026-09-26; the combined M77 -> M87
run was confirmed 2026-09-27 with `COMMIT_POLICY: milestone` and publication
authority scoped to the APItools and OpenUdon `origin/main` branches.

**Specification:** [M77 in milestone.md](milestone.md#m77---step-contract-operation-metadata).
**Direction:** [Evolution v24](../evolution/result-v24.md).
**Priority:** M77 was the only active APItools milestone; it is complete under
the confirmed goal run.

## Evidence and approved decisions

- Combined release scope confirmed 2026-09-27: finish M77 before OpenUdon M87.4
  and M87.6; create one closure commit per milestone; push only the configured
  `origin/main` branches of APItools and OpenUdon; publish the compatible
  APItools Go-module revision through its pushed commit. No tag release is
  planned unless the documented module flow requires it. Ramen and other
  sibling repositories remain out of scope. Each branch currently has one
  unpushed planning commit, which will be included in its authorized push.
- The user requested consumer-readable operation summaries, read/write/unknown
  effect classification, and ranking by purpose/inputs/outputs for OpenUdon's
  `step candidates`, then explicitly approved the four planning-file actions.
- [Kinet's iCoT replacement design](../../../kinet/docs/icot.md#7-proposed-changes-by-package)
  assigns this work to APItools in sibling S2a. The label S2a is a cross-package
  sequencing item, not APItools' permanent status ID S02.
- Inspection baseline: clean APItools
  `096cf5b695f8d88b4040c0f4288df2f3337ea9a9`. `inventory_types.go`,
  `inventory.go`, and `authoring_api.go` already expose request/response summaries,
  auth alternatives, and bounded authoring context. `operation_selection.go`
  provides hint/text selection and method-based purpose heuristics, including
  AWS Query exceptions; it does not implement the requested effect contract.
  Existing tests in `authoring_api_test.go`, `prompt_safety_test.go`, and
  `operationlifecycle/operationlifecycle_test.go` establish compatibility,
  budget, caller-copy, source identity, and ambiguity requirements.
- Existing native parser summary types preserve source-specific information,
  such as GraphQL operation kind and OData action/function identity. Adapters
  must preserve those distinctions and visibly report unavailable semantics.
- OpenUdon was inspected at
  `55b24d29279c8efe67ae931f9d73f094f929efab`, with uncommitted M87 planning;
  [M87.4](../../../openudon/tabilet/memory-bank/status-M87.md)
  already names the APItools metadata dependency. Kinet was inspected at
  `9719a364904743c0d2cb8b3f47cb4e50ed6322e8`, with uncommitted design/planning
  changes. These are planning evidence, not delivered consumer behavior.
- The approved first version is deterministic, offline, additive, and based
  on source evidence. Preserve legacy exported contracts and scoring. No LLM,
  live operation, credential resolution, account selection, or approval policy
  moves into APItools.
- This is new work. Existing candidate directions remain deferred under their
  recorded triggers. Historical M75/S03 ranking and S01/M74/S04 hardening
  remain completed; consult the [history index](../docs/history/index.md)
  for their evidence without reopening them.

## Task ledger

Markers: `[ ]` pending, `[~]` in progress, `[+]` complete, `[!]` blocked,
`[X]` cancelled, `[-]` closed historical with an accepted successor.
Each row is an implementation commit unit under the governing execution
policy. Keep one execution owner and at most one general in-progress row.
The original feature approval authorized implementation only. The confirmed
combined run overrides that default for this run with
`COMMIT_POLICY: milestone` and narrowly authorizes OpenUdon M87 work plus pushes
to the existing APItools and OpenUdon `origin/main` branches and APItools Go
module publication through the pushed revision. No other sibling edits or
external mutations are authorized.

| Item | State | Notes |
|---|---|---|
| M77.1 Define the public metadata contract | `[+]` | Added additive request/report/source/candidate DTOs, tri-state typed values, exact source identity, evidence, capability status, and per-dimension match results with versioned JSON fixtures. Documented deterministic ordering, safety, compatibility, and the source-family matrix in `docs/operation-candidates.md`; no existing APIs or JSON shapes changed. The interoperable purpose, inputs, outputs, and effect names match the pending OpenUdon/UWS contract; APItools keeps richer value metadata local to its v1 interface, with no cross-repository Go dependency or sibling edit. Focused and full tests pass. |
| M77.2 Produce consumer-readable summaries | `[+]` | Added bounded source-backed purpose text and structured request/response values with stable local evidence references. Missing purpose/schema/output evidence is explicit; unsafe selectors and semantic field/parameter budget overflow block rather than silently truncate. Hostile text is sanitized with diagnostics, credential-shaped parameters are omitted as workflow data with an auth-review gap, caller-owned summaries remain unchanged, and no provider/account/recipient claims are invented. Focused tests and `go test ./...` pass. |
| M77.3 Classify operation effects with evidence | `[+]` | Added conservative operation-ID/summary/description effect signals with evidence refs and public rationale. HTTP method is never evidence; method-only and neutral operations remain unknown, POST searches can classify read, misleading-method mutations can classify write, and compound/conflicting meaning stays unknown. Added a native-protocol evidence hook; actual family wiring and fixtures belong to M77.5. Focused effect tests and `go test ./...` pass. |
| M77.4 Rank operations by step contract | `[+]` | Ranks purpose, available request inputs, expected outputs, and explicit effect constraints with separate scores/explanations and compatible/incompatible/indeterminate statuses. Primitive type direction and requiredness conflicts are tested; unresolved named types remain indeterminate even when names match. Missing evidence is not a positive match, ties/no matches are reported, effect conflicts remain visible, and unknown cannot satisfy read/write constraints. Source-qualified order is deterministic and caller contracts/candidates remain unchanged. Focused ranking tests and `go test ./...` plus `git diff --check` pass. |
| M77.5 Integrate inventories and native source metadata | `[+]` | Added public `BuildOperationCandidates` over explicitly supplied local files/bytes, raw-content SHA-256, native selectors, OpenAPI security alternatives, and adapters for all eight planned families using existing parsers. Added a versioned local source fixture per family. Request/response capabilities and known gaps are explicit; credential-shaped fields are omitted from the returned candidate operation and summary with an auth-review gap. File/symlink, cancellation, exact digest/selector, duplicate IDs across files, family fields/effects/auth, semantic OpenAPI field truncation, native JSON Schema composition gaps, and stable diagnostics under source reordering have focused coverage. Corrected the Google Discovery source-level capability wording; incomplete OpenAPI security schemes are explicitly partial. Focused and full `go test ./...`, `go vet ./...`, standalone `GOWORK=off go test ./...` / `go vet ./...`, affected race tests, and `git diff --check` pass. |
| M77.6 Qualify the public contract and consumer handoff | `[+]` | Fixed the false-positive credential filtering for camel-case `projectKey` and `idempotencyKey`; standalone `key`, delimiter-separated key paths, and explicit credential contexts remain filtered. Reconciled OpenUdon's M87.4 consumer gate and handed off the local API contract; M87.4 will verify the exact published revision. Uncached APItools workspace/standalone tests and vet, focused race tests, `git diff --check`, OpenUdon/Udon workspace and standalone consumer checks, OpenUdon boundary/doc-memory checks, and Udon workspace tests/vet pass. Ramen remains excluded by owner instruction. The full-milestone review-fix gate passed after iteration 7. The confirmed run authorizes the APItools `origin/main` push and Go-module publication; this milestone's code acceptance does not claim OpenUdon adoption. |

## Dependencies and ownership

There is no pending upstream APItools milestone. Completed inventory, prompt
safety, and lifecycle contracts are the baseline to preserve. M77 owns APItools
metadata and its public summary/classification/ranking interface. Scope includes
the existing source families listed in M77.5, with visible evidence limitations;
it does not claim arbitrary schema compatibility or infer safe execution from
textual similarity. A source with unavailable details must remain identifiable
with explicit diagnostics rather than silently being omitted.

OpenUdon M87.1 can define command contracts independently; M77.1 can likewise
define reusable metadata independently. Reconcile their requirements before
freezing the handoff. OpenUdon M87.4 and M87.6 require the compatible published
APItools API. APItools acceptance is its tested library/fixture delivery and
compatibility evidence, not downstream command completion. Existing standalone
consumer suites using older pins cannot prove adoption of the new interface;
record that distinction and the later revision gate explicitly.

Kinet owns authoring decisions and repair orchestration; OpenUdon owns step
binding, command JSON, accounts/destinations, and package approval; UWS/Udon own
workflow semantics and execution enforcement. Unknown effects are conservative
metadata and must never be promoted into a read-only guarantee. User overrides
and their authority/audit remain downstream. Browser steps, simulation, runtime
probing, parser replacement, catalog expansion, and iCoT retirement are outside
M77. No sibling file change or remote publication is part of this plan.

## Acceptance and planned verification

- Summaries explain useful operation behavior, required inputs, and available
  outputs using source evidence, with explicit missing-information fallbacks.
- Effect fixtures include meaningful read and write classifications as well
  as unknown; returning unknown for every operation does not satisfy acceptance.
  Method-only and contradictory evidence cannot establish read safety.
- With the same purpose text, changing available input types/requiredness or
  expected outputs changes ranking compatibility evidence appropriately.
  Similar names alone cannot hide incompatibility or missing evidence.
- All listed source families have metadata adapter fixtures. Preserve exact
  source identity, digest, native selector, and auth alternatives; document
  unsupported structural or semantic details explicitly. Duplicate operation
  IDs across different sources remain distinct, and ties stay ambiguous.
- Cover no matches, malformed contracts, hostile/invisible/control text,
  budget exhaustion, semantically unsafe truncation, deterministic ordering,
  independent return values, and unchanged caller graphs. Preserve existing
  authoring, selection, inventory, and lifecycle tests.
- Run focused metadata, source-adapter, conformance, and affected race tests;
  `go test ./...`; `go vet ./...`; `GOWORK=off go test ./...`;
  `GOWORK=off go vet ./...`; and `git diff --check` in APItools.
- Run available OpenUdon and Udon `go test ./...` and standalone
  `GOWORK=off go test ./...` / `GOWORK=off go vet ./...`. Record which dependency
  revision each mode checks. Compile and exercise new public API consumer
  examples against the local implementation. Default checks use local fixtures
  without provider calls, credentials, or a model service.
- Provide a release-ready handoff naming the API, fixture version, limitations,
  and evidence OpenUdon M87.4 must verify against a later published revision.
  Do not claim that preparing this handoff published or adopted the API.

## Review and closure record

**Review state:** Passed after iteration 7.
**Review iterations started:** 7 of at most 10.

Iteration 1 findings were fixed and verified before iteration 2:

P1 URL provenance is now stripped of userinfo/query/fragment and reported with a
warning; camelCase credential fields are omitted from candidate operations,
summaries, and nested request/response fields. P2 fixes report AsyncAPI payload
alternatives, mark OpenAPI multi-media/success-response branches partial, carry
schema-required output fields, enforce non-OpenAPI operation limits before
candidate adaptation, and align the v1 candidate/source fixtures with producer
semantics. Added end-to-end URL, credential, branch, work-budget, required-output,
and producer-fixture tests. Focused tests, full and standalone tests/vet,
affected race tests, and `git diff --check` pass.

Iteration 2 review findings, persisted before fixes, were resolved and verified:

- **P2:** OData CSDL parameter `Nullable` metadata does not establish whether
  a call argument is required, but the adapter inherited the legacy bool's
  zero value and emitted `required:false` in the new tri-state summary. Keep
  OData argument requiredness unknown and attach the limitation to the
  per-operation capability/gap evidence.
- **P2:** input alias resolution accepts the first of multiple step-contract
  names that map to one source input, silently ignoring a conflicting alias
  type; output aliases can count one source output more than once. Treat such
  duplicate alias mappings as indeterminate and test both directions.

OData requiredness is now unknown in operation summaries, with explicit input
and output capability gaps. Input alias collisions are indeterminate; duplicate
output aliases are indeterminate and do not double-count ranking score. Added
regression tests for both. Focused tests, `go test ./...`, `go vet ./...`,
standalone `GOWORK=off` test/vet, affected race tests, and `git diff --check`
pass.

Iteration 3 findings, persisted before fixes, were resolved and verified:

- **P1:** OData entity-set, singleton, and navigation resource selectors were
  classified as reads even though they do not identify a read-only operation
  and can support mutations. Keep those effects unknown with a partial effect
  capability; only source operation kinds with protocol-defined action/function
  semantics may provide a read/write hint.
- **P1:** candidate rationale sanitization deleted from a slice while iterating
  it by index. An unsafe first candidate followed by another candidate can
  panic from an out-of-range index, rather than safely omitting the affected
  item. Rebuild the sanitized candidate slice and test a rejected candidate
  followed by a valid one.
- **P2:** operation meaning text containing a negated action (for example,
  “does not delete”) is classified from the action verb alone. Negation can
  invert a read/write hint and distort downstream compatibility. Conservatively
  keep negated action phrases unknown and cover common negation forms.
- **P2:** OpenAPI scalar/opaque successful response schemas are projected as a
  synthetic `body` field with `required:false`, although response-body
  requiredness is not established by that Boolean. The summary can then report
  an unsupported incompatibility for a required step output. Preserve unknown
  body requiredness and test candidate ranking.
- **P2:** an over-context candidate report can still return all of its large
  candidates alongside a budget error, violating the result bound if a caller
  inspects or serializes the error report. Replace an oversized result with a
  bounded diagnostic-only report and test the public producer path.

OData action/function kinds now provide write/read hints respectively; resource
selectors remain unknown with partial capability and source-kind evidence.
Negated actions remain unknown with evidence; synthetic OpenAPI response-body
requiredness remains unknown. Candidate rationale filtering rebuilds the slice
without panicking, and over-context results shed candidate/source payloads in
favor of a bounded diagnostic. Regression tests cover each finding. Focused
tests, `go test ./...`, `go vet ./...`, standalone `GOWORK=off` test/vet,
affected race tests, and `git diff --check` pass.

Iteration 4 review finding, persisted before fix:

- **P1:** effect classification tokenizes raw native-family IDs and
  descriptions before the candidate sanitizer applies prompt-text bounds. A
  single bounded-size local document can therefore trigger millions of token
  allocations or exhaust memory. Cap effect evidence work before tokenization
  and keep oversized text unknown with source evidence.
- **P2:** the native JSON Schema walker expands object properties and local
  references but silently ignores `allOf`, `oneOf`, and `anyOf`. Important
  fields can disappear without a specific summary gap. Surface the unexpanded
  composition limitation and cover it with a source-schema test.
- **P2:** non-OpenAPI adapters do not check the caller context while enumerating
  parsed operations. A cancellation during a bounded but large source therefore
  waits until all operation candidates have been built. Check cancellation
  between native operations and exercise a mid-adaptation cancellation.

Iteration 4 findings were resolved and verified before iteration 5:

- Effect classification now scans only through the prompt-text rune budget and
  leaves oversized meaning text unknown while preserving bounded source
  evidence. A large untrusted-description regression test covers the bound.
- Native schema summaries report unexpanded `allOf`, `oneOf`, and `anyOf`, and
  dictionary-valued `additionalProperties`, without treating the unsupported
  branches as complete fields. Composition and dictionary gaps have regression
  coverage; `additionalProperties: false` does not produce a false gap.
- Every non-OpenAPI native adapter checks the caller context between parsed
  operations; the OpenAPI operation loop does as well. Mid-adaptation
  cancellation is tested.

Focused tests, `go test ./...`, `go vet ./...`, standalone `GOWORK=off` test and
vet, affected race tests, and `git diff --check` pass. Iteration 5 was recorded
as started before beginning the next full-milestone review.

Iteration 5 review finding, persisted before fix:

- **P2:** OpenAPI operations whose security requirements name an undeclared or
  incomplete security scheme retain a requirement record with an empty scheme
  type, but the candidate still reports auth capability as supported with no
  gap. Mark this evidence partial and explain the unresolved scheme so a
  consumer cannot mistake an incomplete alternative for a complete auth
  summary.

Iteration 5 finding was resolved and verified before iteration 6. OpenAPI
candidate auth capability now becomes partial and carries a bounded explicit
gap if any retained security requirement lacks a name or scheme type; the
original security alternative remains available for review. The regression
test confirms both the retained requirement and partial capability. The
candidate documentation and architecture notes describe this behavior.
Focused tests, full and standalone tests/vet, affected race tests, and
`git diff --check` pass. Iteration 6 was recorded as started before beginning
the next full-milestone review.

Iteration 6 reviewed the full M77 implementation and its accumulated fixes;
no P1, P2, or higher-severity findings remained at that point.

Iteration 7 was started after OpenUdon handed off its shared-workspace consumer
failure. The **P2** compatibility finding was that splitting camel-case names
into tokens made ordinary `projectKey` and `idempotencyKey` fields match the
generic `key` credential token and disappear from candidate metadata. The
classifier now treats a camel-boundary `key` as credential-shaped only with an
explicit credential context; standalone and delimiter-separated `key` names
remain conservatively filtered. Added regressions for the two consumer fields,
standalone and path `key`, and credential-context keys. Uncached APItools and
OpenUdon workspace test suites pass after the fix. Iteration 7 re-reviewed the
M77 contract, adapters, summaries, effect evidence, ranking, budgets, fixtures,
and shared credential-filter call sites; no further P1, P2, or higher-severity
findings remain. The review-fix gate passes after seven iterations.

M77.6 verification recorded on 2026-09-27: APItools `go test -count=1 ./...`,
`go vet ./...`, `GOWORK=off go test -count=1 ./...`, `GOWORK=off go vet ./...`,
`go test -race -count=1 .`, and `git diff --check` passed. OpenUdon full
workspace tests/vet passed using `/home/peter/Workspace/go.work`, and standalone
test/vet passed with `GOWORK=off`. Udon workspace tests/vet passed; standalone
test/vet passed with `GOWORK=off` and a temporary tidy-equivalent module-file
copy because the tracked Udon `go.mod` otherwise requests tidy changes. Its
worktree was unchanged. Ramen checks are excluded by the owner's 2026-09-27
instruction that it is a separate project. The OpenUdon M87.4 status-only
update reconciles the formerly stale workspace-failure record; M87.4 and M87.6
remain the downstream owners for verification against the exact published
module revision. The compatible Go-module revision will be published by the
authorized APItools `origin/main` push in this combined run.

Persist iteration 1 and its evidence before the first full-milestone review.
Resume an interrupted iteration at its recorded number. Review the full change,
fix P1/P2 or higher findings within the ten-iteration gate, and retain unresolved
blockers at the limit. Acceptance requires verified behavior, current-memory
consolidation, and downstream reconciliation before normal retirement; terminal
task markers alone are insufficient.

## Planning verification (historical, before implementation)

The approved file scope is `milestone.md`, this status file, and
`prompt-v24.md` / `result-v24.md`. Planning checks cover ledger parsing, six
pending task rows, permanent-ID uniqueness across active/history records,
local links/anchors, exact changed-file scope, and whitespace. Implementation
verification and consumer acceptance remain pending.

Verified 2026-09-26: the installed runner's read-only parser recognizes exactly
six pending M77 rows and no in-progress row. Active and retired IDs do not
overlap; all 15 rendered local links/anchors across the four planning documents
resolve. Candidate Directions are unchanged, only the four approved files
are modified or created, and `git diff --check` plus explicit whitespace checks
for new files pass. No implementation tests, commit, or publication were run.
````
