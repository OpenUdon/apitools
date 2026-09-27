# Status M77 - Step-contract operation metadata

**State:** Pending. Planning approved 2026-09-26. Implementation and acceptance
have not started.

**Specification:** [M77 in milestone.md](milestone.md#m77---step-contract-operation-metadata).
**Direction:** [Evolution v24](../evolution/result-v24.md).
**Priority:** Next and only active APItools milestone; begin with M77.1 under
a separate execution request.

## Evidence and approved decisions

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
This planning approval itself authorizes no implementation, commit, or push.

| Item | State | Notes |
|---|---|---|
| M77.1 Define the public metadata contract | `[ ]` | Own additive Go types, bounded purpose/typed-input/expected-output contracts, source kind/native selector/content digest identity, summary/effect evidence, capability and ranking diagnostics, and versioned conformance fixture shapes. Preserve existing exported data shapes and APIs. Document deterministic ordering, compatibility, and the source-family matrix. Reconcile the OpenUdon M87.1 handoff without requiring downstream implementation. |
| M77.2 Produce consumer-readable summaries | `[ ]` | Depends on M77.1. Generate concise descriptions of operation purpose, required inputs, and available outputs from source evidence. Distinguish facts, inferred wording, and gaps; use honest fallback wording when descriptions are inadequate. Never invent accounts, recipients, operations, or capabilities. Apply existing sanitization/budgets and test useful summaries plus missing and hostile descriptions. |
| M77.3 Classify operation effects with evidence | `[ ]` | Depends on M77.1. Return read/write/unknown with source-backed reasons and evidence provenance. Define conservative precedence for operation meaning and protocol metadata; a method alone cannot prove read. Missing, insufficient, or conflicting evidence stays unknown. Keep inferred effects separate from user-confirmed decisions and existing lifecycle purposes. Test real read/write cases, POST reads, mutating operations using misleading methods, conflicts, and ambiguous metadata. |
| M77.4 Rank operations by step contract | `[ ]` | Depends on M77.1-M77.3. Rank purpose, required input coverage, and expected output compatibility with separate explanations and compatible/incompatible/indeterminate results. Evaluate type and requiredness direction correctly, distinguish missing schema evidence from matches, expose mapping gaps without binding a workflow, and report ties/no matches. Effect incompatibility must remain visible and unknown cannot satisfy read-only constraints. Preserve deterministic source-qualified identity and legacy ranking behavior. |
| M77.5 Integrate inventories and native source metadata | `[ ]` | Depends on M77.1-M77.4. Connect the new reports to OpenAPI/Swagger inventories and adapt existing Google Discovery, AWS Smithy, AsyncAPI, GraphQL, OpenRPC, gRPC/protobuf, and OData metadata. Reuse parsers and preserve native selectors/kinds; no OpenAPI lowering or new parser ownership. Supply per-family fixtures and explicit capability gaps. Preserve OR-of-AND auth, source digests, caller-owned data, and prompt/work/byte budgets; semantic truncation fails visibly. |
| M77.6 Qualify the public contract and consumer handoff | `[ ]` | Depends on M77.1-M77.5. Complete public examples, versioned conformance fixtures, exact-source/effect/input/output/auth/ambiguity negatives, compatibility and budget regressions, and provider-free verification. Document the interface for OpenUdon M87.4 and later published-pin evidence. Update current memory only for verified behavior. Run required consumer checks and the persisted milestone review; publishing, sibling edits, and Kinet adoption remain separately owned. |

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
- Run available OpenUdon, Ramen, and Udon `go test ./...` and standalone
  `GOWORK=off go test ./...` / `GOWORK=off go vet ./...`. Record which dependency
  revision each mode checks. Compile and exercise new public API consumer
  examples against the local implementation. Default checks use local fixtures
  without provider calls, credentials, or a model service.
- Provide a release-ready handoff naming the API, fixture version, limitations,
  and evidence OpenUdon M87.4 must verify against a later published revision.
  Do not claim that preparing this handoff published or adopted the API.

## Review and closure record

**Review state:** Not started.
**Review iterations started:** 0 of at most 10.

Persist iteration 1 and its evidence before the first full-milestone review.
Resume an interrupted iteration at its recorded number. Review the full change,
fix P1/P2 or higher findings within the ten-iteration gate, and retain unresolved
blockers at the limit. Acceptance requires verified behavior, current-memory
consolidation, and downstream reconciliation before normal retirement; terminal
task markers alone are insufficient.

## Planning verification

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
