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
