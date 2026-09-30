# Lessons

Keep concise, reusable lessons that still affect decisions. Consult the topics
relevant to the current task before substantial changes; this is not a session
log or a requirement to produce one lesson per milestone.

Maintain relevant lessons when ordinary work produces reusable, evidence-backed
learning, and consolidate them during milestone closure. A still-applicable
lesson stays here even after its supporting milestone retires; closing a task
or reaching an arbitrary file size is not a reason to discard useful knowledge.

For each lesson, use a descriptive heading and record when it applies, the
lesson, why it matters, and links to supporting tasks, verification, or retired
records. Keep product facts in `product.md`, system contracts in
`architecture.md`, and commands in `tech-stack.md` rather than repeating them.

Merge duplicates. Before materially replacing or removing an obsolete lesson,
append its previous wording, source, reason, and replacement reference to
`tabilet/docs/history/knowledge.md` under the retirement rules in
[milestone.md](milestone.md#long-term-memory-and-retirement). Preserve the
evidence links when merging lessons. Revalidate historical evidence before
applying it to current work. Knowledge preservation is triggered by material
supersession or removal, even outside milestone closure; routine wording edits
need no journal entry or separate archive run.

## Keep inferred POST actions out of resource creation

- **Applies when:** Ranking lifecycle roles from API documents whose operation IDs and paths may use provider-specific naming.
- **Lesson:** HTTP method alone is not enough to label a POST as create. Respect explicit update/create semantics named on the operation itself before path shape, and known action routes; on a path whose *trailing* segment is a parameter (a POST directly against an item ID), keep an otherwise ambiguous POST generic rather than treating a finite action-name list as exhaustive. A parameter earlier in the path only scopes a parent resource (for example `{projectId}` in a nested-collection create) and must not by itself suppress a create classification.
- **Why it matters:** Custom verbs such as `renew` or `reprocess` can mutate an existing resource without creating it, and checking every path segment for a parameter (not just the trailing one) misclassifies ordinary nested-collection creates as actions. A false classification either way can distort a generated lifecycle proposal.
- **Evidence:** [S03 - Operation Lifecycle Ranking Correctness](../docs/history/status-S03.md), especially review iteration 1's `renewInvoice` regression and iteration 3's `newChild`/nested-collection regression (U4).

## Keep incomplete catalog coverage distinct from negative evidence

- **Applies when:** Building indexed discovery over heterogeneous or large
  registered API specifications.
- **Lesson:** Bind derived operations to verified raw identity and record every
  reference's coverage separately. Selected providers with no references also
  remain unexamined; curated availability is not absence proof. Preserve
  positive metadata from partially
  summarized sources while keeping their scope unexamined. Missing/stale or
  unsupported sources cannot establish that no qualifying API exists. Apply
  query link and per-operation prompt bounds before multiplicative allocation
  or publication; removed selected fields cannot become known absence.
- **Why it matters:** A valid large source can exceed schema-summary budgets
  even when it fits a raised byte limit. Counting retained candidates as full
  coverage would silently authorize downstream fallback on missing evidence.
- **Evidence:** [M81 task/review record](../docs/history/status-M81.md), including local large
  artifact measurements and index coverage, drift and failure regressions;
  [M80 review findings](status-M80.md) and the source-backed field-loss,
  no-reference provider, shared-link and prompt-budget regressions.
