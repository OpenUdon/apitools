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
- **Lesson:** HTTP method alone is not enough to label a POST as create. Respect explicit update/create semantics and known action routes; on parameterized item/action paths, keep otherwise ambiguous POSTs generic rather than treating a finite action-name list as exhaustive.
- **Why it matters:** Custom verbs such as `renew` or `reprocess` can mutate an existing resource without creating it. A false create classification can distort a generated lifecycle proposal.
- **Evidence:** [S03 - Operation Lifecycle Ranking Correctness](status-S03.md), especially review iteration 1 and the `renewInvoice` regression.
