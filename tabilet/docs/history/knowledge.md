# Superseded knowledge journal

## 2026-09-28 — Step-candidate output nullability is field-scoped

**Original source.** `docs/operation-candidates.md`, headings `Ranking and uncertainty` and `Compatibility`.

**Replacement reason and evidence.** M79.2 carries nullable state onto response fields and their nullable schema ancestors, then checks only outputs selected by the step contract. A nullable sibling no longer downgrades unrelated selected outputs; a selected nullable output remains indeterminate because the v1 step contract cannot declare that null is acceptable. The public JSON change is additive: existing field meanings remain stable while optional nullable fields are exposed. The current contract is documented in [operation-candidates.md](../../../docs/operation-candidates.md), and implementation evidence is retained in the M79 status until retirement.

**Superseded wording.**

````markdown
OpenAPI response schemas that permit null are
reported as partial output evidence. The v1 step contract has no nullability
field, so a nullable response cannot earn output compatibility points or a
compatible output status.

This is an additive API. Existing inventory, selection, ranking, and lifecycle
functions and JSON shapes remain unchanged.
````

**Retirement correction (2026-09-28).** M79 completed after OpenUdon M89 consumed
the exact published APItools revision. The complete milestone status and review
evidence are preserved in the [retired M79 record](status-M79.md).

## 2026-09-30 — Exclude the separate Ramen project from downstream verification

**Original source.** `tabilet/memory-bank/tech-stack.md`, heading `Common Commands`, downstream verification paragraph and command block.

**Replacement reason and evidence.** The user explicitly directed APItools to
ignore Ramen when testing downstreams because it is now a separate project,
and approved this planning correction on 2026-09-30. AGENTS.md names OpenUdon
and Udon as the dependent checks; the accepted V24 direction already records
the same Ramen exclusion. This changes verification scope, not a claim about
Ramen's API compatibility or runtime behavior. Replacement:
[current downstream commands](../../memory-bank/tech-stack.md#common-commands).

**Superseded wording.**

````markdown
When exported APIs change, run dependent checks in sibling consumers when
available:

```bash
(cd ../openudon && go test ./...)
(cd ../ramen && go test ./...)
(cd ../udon && go test ./...)
(cd ../openudon && GOWORK=off go test ./...)
(cd ../ramen && GOWORK=off go test ./...)
(cd ../udon && GOWORK=off go test ./...)
```
````

## 2026-10-06 — S06 consumer artifact cleanup timing correction

**Original source.** [Retired S06](status-S06.md), retained status section
`Whole-milestone review and closure`.

**Correction and evidence.** S06's candidate-bound consumer checks passed.
After the closure commit, a final Udon status check found two further untracked
test-output directories under `pkg/runner/home` and `spider/home`, in addition
to the already moved CLI test artifacts. They contained only the observed
synthetic test outputs from the task-owned build-temp path and were moved into
disposable verification storage. Subsequent Git status is clean in both
consumers. Thus the retained clean-worktree statement was premature in timing;
tracked manifests, pins and source files were unchanged throughout. The frozen
record is not rewritten. Replacement evidence is this correction and
[S05's execution notes](../../memory-bank/status-S05.md#execution-checkpoint).

**Original wording.**

````markdown
- Udon tests generated an untracked directory under its CLI package; the
  known test artifacts were moved into disposable verification storage.
  Both sibling Git worktrees are clean and their tracked files/pins unchanged.
````

## Stage 11 approved target provenance — 2026-10-06

Source: `AGENTS.md` at `9364afac98c97b79c9e6b9335bd9c71d4e8722d5` (clean). The following literal prior boundary remains evidence for existing behavior. The approved Stage 11 proposal adds a scoped target/transition; it does not claim that code has already moved. Replacement target: [Stage 11](../../../../kinet/docs/stage11.md) and the active milestone specifications. No retired record is changed.

````markdown
## Boundary

- `../apitools` owns API source discovery, download safety, validation, import,
  local file scanning, operation inventories, prompt-safe summaries,
  auth/security summaries, deterministic operation ranking, optional caching,
  provider catalog metadata, OpenAPI/Swagger, Google Discovery, and AWS Smithy
  protocol classification, and docs-derived endpoint overlay assets when no
  official OpenAPI exists but official API docs support a reviewed subset.
- `../openudon` consumes `apitools` metadata for authoring, review, package
  evidence, and trusted-runner handoff. OpenUdon owns workflow behavior,
  approval state, release gates, examples, and credential binding policy.
- `../uws` owns public workflow semantics and Go model definitions.
- `../udon` owns private UWS/OpenAPI lowering and runtime execution.
- `../tfconfig` owns static Terraform/OpenTofu configuration parsing only.

Rule of thumb: if a change helps find, validate, import, summarize, classify,
or rank API documents, it belongs here. If it executes workflows or decides
runtime credential/account behavior, it belongs downstream.
````

## M82 replaces planning/progress wording in product.md — 2026-10-06

Source: `tabilet/memory-bank/product.md` at observed baseline `54583f9b2f452b7cc522360c5aeeff29ca22f96c`; prior Stage 11 planning/progress wording is superseded by qualified M82 source tooling. Replacement: [product.md](../../memory-bank/product.md#qualified-stage-11-source-metadata--m82) and [qualification](../../../docs/m82-qualification.md). Frozen records remain unchanged.

````markdown
## Stage 11 implementation progress — M82.1/M82.2

The additive local `BuildOperationShapeTable` producer projects all eight
existing API source families into public UWS metadata. Native selectors and
explicit incomplete schemas/protocol evidence remain visible. Supported explicit
OpenAPI contracts can prove metadata compatibility; exact independent reproduction
refuses forged/stale claims, while unsupported native details stay partial.
Symbolic security preserves OR-of-AND structure without resolving credentials.
This implementation is under M82,
whose acceptance, exact verification and publication remain pending. It performs
no fetching, credential resolution or execution; see [contract](../../docs/operation-shapes.md).

## Approved Stage 11 product direction — not implemented

[Stage 11](../../../kinet/docs/stage11.md) and [local milestones](milestone.md#stage-11-cross-package-refactoring) define the approved target. APItools will produce source-neutral UWS operation shapes from its existing eight source families. It still performs no credential resolution, workflow approval or runtime execution. Kinet and OpenUdon consume exact published metadata contracts through isolated parsing paths.
Both phases belong to one stage. Current facts below remain the observed implementation; no new acceptance, publication or installed behavior is claimed. The installed Kinet M44 service remains unchanged, and Stage 12 owns the browser-dependent removal gates.

````

## M82 replaces planning/progress wording in architecture.md — 2026-10-06

Source: `tabilet/memory-bank/architecture.md` at observed baseline `54583f9b2f452b7cc522360c5aeeff29ca22f96c`; prior Stage 11 planning/progress wording is superseded by qualified M82 source tooling. Replacement: [architecture.md](../../memory-bank/architecture.md#qualified-stage-11-source-metadata--m82) and [qualification](../../../docs/m82-qualification.md). Frozen records remain unchanged.

````markdown
## Stage 11 implementation progress — M82.1/M82.2

`operation_shapes*.go` adds an independent metadata projection over the existing
native parsers and public UWS binding types. Explicit source IDs plus bounded
local bytes/files produce raw SHA-256 identities, native selector aliases,
schema projections and visible incompleteness. The producer returns no partial
table on malformed input or limits; URLs are sanitized provenance only. Existing
inventory/candidate contracts are preserved. `VerifyOperationShapeTable`
independently reproduces exact claims against the source bytes. OpenAPI symbolic
security preserves OR-of-AND grouping; absent/unresolved security never becomes
anonymous. Complete projected schemas compile with an always-refusing resource
loader; unsupported native dialect/wire/requiredness evidence remains partial.
Wire qualification and consumer publication remain M82.3–M82.4 work.

## Approved Stage 11 architecture target — not implemented

[Stage 11](../../../kinet/docs/stage11.md) and [local milestones](milestone.md#stage-11-cross-package-refactoring) define the approved target. APItools will produce source-neutral UWS operation shapes from its existing eight source families. It still performs no credential resolution, workflow approval or runtime execution. Kinet and OpenUdon consume exact published metadata contracts through isolated parsing paths.
Both phases belong to one stage. Current facts below remain the observed implementation; no new acceptance, publication or installed behavior is claimed. The installed Kinet M44 service remains unchanged, and Stage 12 owns the browser-dependent removal gates.

````

## M82 replaces planning/progress wording in tech-stack.md — 2026-10-06

Source: `tabilet/memory-bank/tech-stack.md` at observed baseline `54583f9b2f452b7cc522360c5aeeff29ca22f96c`; prior Stage 11 planning/progress wording is superseded by qualified M82 source tooling. Replacement: [tech-stack.md](../../memory-bank/tech-stack.md#qualified-stage-11-source-metadata--m82) and [qualification](../../../docs/m82-qualification.md). Frozen records remain unchanged.

````markdown
## Stage 11 implementation progress — M82.1

The local shape producer now directly imports UWS C09 at
`v0.0.0-20261006181058-6a267306032e`, accepted source
`6a267306032edc687a298cefc8bba7019d3ad059`. Its module closure includes Horizon
`v1.14.5` and HashiCorp HCL `v2.24.0`; no HCL-free claim is made. Exact dependency
acquisition/tidy used `GOWORK=off`; offline tests use `GOPROXY=off`. Focused shape
checks are `go test . -run '^TestOperationShapes'` and the `-race` variant.
The same pinned `jsonschema/v6 v6.0.1` dependency is directly used to check
projected schemas with external loading disabled; no dependency version changed.
M82 whole qualification and consumer publication remain pending.

M82.3 freezes the UWS eight-family table at 9,829 bytes, SHA-256
`dc20d2287322a4b20d5d92b1a0d1ec8036f1bec853b642d7df8797ed096c3ba1`.
The public UWS Marshal/ParseTable pair is reused; no duplicate serializer or CLI
is introduced. Standalone `go list -m all` confirms the exact UWS, schema-compiler
and retained Horizon/HCL closure; legacy candidate v1 wire fixtures are unchanged.

## Approved Stage 11 tooling target — not implemented

[Stage 11](../../../kinet/docs/stage11.md) and [local milestones](milestone.md#stage-11-cross-package-refactoring) define the approved target. Exact published module revisions, frozen build closures and standalone verification are acceptance gates. Planned APItools:M82 will depend on UWS; while UWS core retains legacy HCL support, the module closure includes Horizon and HashiCorp HCL. This is a future adoption consequence, not a claim that current go.mod already imports UWS. go test ./...; go vet ./...; affected race/API/source-family fixtures and OpenUdon/Udon consumer checks; git diff --check. Use bounded local artifacts, not remote discovery.
Both phases belong to one stage. Current facts below remain the observed implementation; no new acceptance, publication or installed behavior is claimed. The installed Kinet M44 service remains unchanged, and Stage 12 owns the browser-dependent removal gates.

````

## M82 source-tooling bootstrap replacement — 2026-10-06

Source: AGENTS.md at observed baseline `54583f9b2f452b7cc522360c5aeeff29ca22f96c`. The planned producer is now qualified; consuming workers and deployment remain separately owned. Replacement: [AGENTS.md](../../../AGENTS.md#stage-11-source-metadata) and [qualification](../../../docs/m82-qualification.md).

````markdown
## Approved Stage 11 planning

[Stage 11](../kinet/docs/stage11.md) coordinates both refactoring phases across five package-local ledgers. Read the local [milestones](tabilet/memory-bank/milestone.md) before selecting work. APItools will produce source-neutral UWS operation shapes from its existing eight source families. It still performs no credential resolution, workflow approval or runtime execution. Kinet and OpenUdon consume exact published metadata contracts through isolated parsing paths.
Planning is approved; implementation and named publication/deployment authority are separate. One serial execution owner, offline fixtures, exact upstream reconciliation and persisted milestone reviews apply. Completed records and frozen consumer pins stay preserved.

````

## M82 replaces pending dashboard/horizon wording — 2026-10-06

Source: milestone.md at observed `54583f9b2f452b7cc522360c5aeeff29ca22f96c`, with completed-task closure metadata. Replacement: [current milestone dashboard](../../memory-bank/milestone.md) and [M82 record](status-M82.md).

````markdown
## Stage 11 active horizon

Approved 2026-10-06: both phases of [Kinet STG-11](../../../kinet/docs/stage11.md), with one serial execution owner across Kinet, UWS, APItools, Udon and OpenUdon. This package owns M82; all rows are pending and each review is 0/10. [Specifications](#stage-11-cross-package-refactoring) below are the current horizon. Earlier completed horizons and records remain historical; planning grants no implementation or external authority.


That earlier horizon is complete; Stage 11 M82 is now planned and unimplemented. Merge/publication and consumer
adoption are separately authorized; the main workspace and consumer pins are
unchanged. Candidate directions remain unnumbered and require fresh approval.
````
