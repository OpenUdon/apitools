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
