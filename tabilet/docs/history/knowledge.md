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
