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
  [M80 review findings](../docs/history/status-M80.md) and the source-backed field-loss,
  no-reference provider, shared-link and prompt-budget regressions.

## Keep trust projections independent of prompt summaries

- **Applies when:** Producing schemas for advisory binding or verifying cached
  producer metadata against raw API sources.
- **Lesson:** Reproduce claims from exact source bytes and native selectors.
  Prompt summaries and parser-normalized subsets can lose constraints, presence,
  protocol details or numeric precision. A useful partial type is not a complete
  schema. Keep unknown security distinct from anonymous and retain OR-of-AND
  requirement grouping without credential resolution.
- **Why it matters:** Structurally valid forged/stale metadata and lossily
  summarized schemas must not become execution authority or a positive proof.
- **Evidence:** [M82](../docs/history/status-M82.md), public UWS binding checks,
  forged/stale identity, symbolic-security and no-network regressions in
  `operation_shapes_identity_test.go`.

Native protocol direction, default presence, literal-versus-member provenance
and response aliases are also proof inputs. Stage 11 intake probes at
24c36bf40102c2c1d160dc7d0e27fb161e12dbd6 showed AsyncAPI 2.x direction inversion,
defaulted GraphQL arguments still required and Smithy literal query inputs.
Lost serialization or unenforced formats cannot establish complete/known
evidence. Accepted [M83](../docs/history/status-M83.md), source
f2c5693ec39ad6981693d2ad126ff26e3fdd564c, restores these contracts after review5
and complete ordinary owner/consumer proof. Structurally valid tables alone do
not establish native semantics; changed identities require fresh consumer review.

## Bind version evidence to an exact publisher and native identity

- **Applies when:** Discovering versions across shared hosting and native source families.
- **Lesson:** Compare exact origins and repository paths, then preserve the recipe and native selector that produced source evidence. An unchanged known URL cannot establish absence of later versions; a provider aggregate cannot stand in for the selected service.
- **Why it matters:** A hostname substring can authorize the wrong origin, and an HTTP display projection can conceal native operation changes. Both undermine an otherwise bounded advisory check.
- **Evidence:** [Retired S05](../docs/history/status-S05.md), review R1/R2 and exact-host, service-family, conditional and native Discovery regressions.

String token values used for prompt text cannot prove native argument equality. M83 whole reviews showed escaped-string collisions, quoted delimiter boundary mistakes and compatible field merges lost by raw whole-selection equality. Keep native typed argument/child evidence, exact default presence and explicit parser closure/depth guards; a partial schema never grants missing-field or runtime authority. [M83 review evidence](../docs/history/status-M83.md#whole-review-iteration-4--local-code-gate-passed) and the corrected/refusal matrices preserve these lessons.

## Source-specific ignored declarations and complete archive inputs

Ignored source declarations must not fabricate input or completeness constraints.
Apply source-version/location/header-case/reference rules without weakening genuine
bindings or independent security projection. Evidence: [accepted M84](../docs/history/status-M84.md),
operation_shapes_reserved_headers_test.go and the16-case SDK/native authoring
matrix with M83 negative controls.

Immutable module seeding must preserve every archive file, including source *.lock
files; mutable-cache filtering belongs only in cache metadata. Real module directory
roots and immutable hardlinks preserve Go modverify behavior. Reuse passing broad
checks only after full source/version/sum/consumer input equality. Ordinary and local
ZIP containers may differ despite identical normalized files/module sums.
Evidence: [M84 ordinary proof](../../docs/m84-ordinary-proof.json) and retained
cache failures/restored exact lockfiles. No source workaround or user-cache mutation.

## Verify every effective-input projection before consumer acceptance

- **Applies when:** A source feeds independent shape, inventory, summary and candidate APIs.
- **Lesson:** Shape and SDK build/assessment/verification proof does not qualify candidate matching. Exercise the same source through the actual consumer selection/publication path, with reserved and genuine inputs and independent security evidence. Bound scripted confirmation so an unavailable operation fails promptly.
- **Why it matters:** A correct shape can coexist with false candidate required inputs and legitimately block publication; a host workaround would fork native semantics.
- **Evidence:** [M50-NATIVE-F01](../../../kinet/docs/m50-candidate-header-finding.md), actual local/rootless worker ced5468, published APItools0c; [M85](status-M85.md) is approved remediation, not delivered proof.
