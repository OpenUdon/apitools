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

## Verify every effective-input projection without extending reference work

- **Applies when:** Shapes, inventories, summaries and candidates derive independent evidence from one source.
- **Lesson:** Exercise actual selection as well as SDK Build/Assess/Verify. Apply coherent ignored-input semantics at the owner. Preserve existing reference capability boundaries; an added per-parameter traversal can multiply work across operations even under a raw byte bound. Bound fixture confirmation and retain unknown unsupported references.
- **Why:** Shape-only M84 proof left W18 blocked; M85 corrected candidate semantics. Whole1 caught an unnecessary unbounded chain extension before publication, and removal preserved the original single-hop/escaped contract.
- **Evidence:** [M85](../docs/history/status-M85.md), final whole3, exact source409d/current SDK16/CLI8/ordinary proof; [M50-NATIVE-F01](../../../kinet/docs/m50-candidate-header-finding.md) preserves failed prior integration.

## Safe catalog snapshots require full-set drift verification and relative path preservation

- **Applies when:** Snapshotting, copying, or projecting read-only catalog roots for consumers.
- **Lesson:** Preserve nested relative paths rather than flattening with `filepath.Base`. Snapshot metadata and selected source alongside registry/index files under distinct bounded ceilings. Recheck the entire source set for file additions, removals, modifications, and active WAL/journal sidecars before destination normalization. Destination-only SQLite journal normalization to `DELETE` mode must not touch source databases.
- **Why it matters:** Flattening relative paths breaks consumers expecting nested registry/index layouts. Omitting metadata or selected source causes downstream validation failures. In-flight writes or concurrent updates during copying can produce corrupt or inconsistent snapshots.
- **Evidence:** [M87](../docs/history/status-M87.md), review findings F02, F03, F04; `catalog/snapshot.go` and `catalog/snapshot_test.go`.

Snapshot publication stages a complete tree and confines rollback to owned
staging; pre-existing targets are either identical/empty or refused. Use the
catalog root, rather than a nested registry's directory, as the artifact base.
Validate identities against the same explicitly selected catalog before
registration and search. Evidence: S07 snapshot preservation, root-aware cache
and custom-catalog round-trip tests.

## Catalog remote discovery fallback requires provider URL forwarding and strict URL safety

- **Applies when:** Broadening remote catalog discovery from single-source to cascaded multi-source (`SourceAuto`).
- **Lesson:** `SourceAuto` cascades across multiple catalog backends (APIs.guru, LAP Registry, RFC 9727, public-apis). For RFC 9727 discovery to succeed, an installation-side `ProviderURL` must be explicitly provided and plumbed to `SearchOptions.ProviderURL`. Validate the provider URL against SSRF, unpermitted private hosts, and non-canonical schemes before dispatching queries.
- **Why it matters:** Hardcoding a single source disables RFC 9727 provider-directed discovery; conversely, unvalidated remote hosts can lead to private network traversal or SSRF vulnerabilities.
- **Evidence:** [M88](../docs/history/status-M88.md), review finding F06; `catalog_discovery_remote.go` and `catalog_discovery_remote_test.go`.

## LLM touchpoints must use standard effect vocabularies and selectable-only candidate boundaries

- **Applies when:** Formulating operation contracts and candidate recommendations with pluggable LLM assistants.
- **Lesson:** Normalize supported model-generated effects to standard `OperationEffect` values (`read`, `write`, `unknown`); `idempotent_write` maps to `write`. Reject missing or unsupported planning effects. Avoid fabricating wildcard inputs when parameters are missing; instead, expose missing evidence via informational diagnostics (`contract.input_missing`). In candidate recommendation, restrict assistance strictly to verified `Selectable` candidates; never recommend unqualified, omitted, or hallucinated candidate indexes. Always revalidate serialized JSON and prompt bounds after attaching explanation text.
- **Why it matters:** Non-standard effect strings cause qualification failures. Wildcard inputs conceal missing parameter dependencies. Recommending non-selectable candidates violates human confirmation boundaries and can authorize incorrect operations.
- **Evidence:** [M89](../docs/history/status-M89.md), review findings F10, F11, F12, F13; `llm_plan.go`, `llm_recommend.go`, `llm_plan_test.go`, and `llm_recommend_test.go`.

Reject incomplete or mistyped completion objects and bound the entire response
before fence extraction. Sanitizing or truncating a valid JSON object's fields
does not make an invalid or oversized completion a usable plan. Deterministic
fallback preserves useful intent when model assistance fails. Evidence:
[S07 qualification](../../docs/s07-qualification.md), invalid planning and
recommendation completion regressions.
