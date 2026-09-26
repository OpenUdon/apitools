# Retired milestone S04 - Prompt Sanitizer Invisible-Unicode Hardening

**Milestone.** S04
**Outcome.** completed
**Retired.** 2026-09-26
**Source status.** tabilet/memory-bank/status-S04.md
**Source specification.** tabilet/memory-bank/milestone.md#s04---prompt-sanitizer-invisible-unicode-hardening
**Evidence.** 07d7ac4d2171600dc793ec6478808e26b0c9025f
**Worktree.** clean
**Review.** passed
**Review iterations.** 3
**Verification.** go build ./..., go vet ./..., go test ./... (apitools, all 17 packages, including focused sanitizer tests); git diff --check; OpenUdon full go test ./.... All passed at the evidence commit.
**Consolidated into.** no current-truth change; current facts already live in architecture.md (prompt sanitizer's Unicode guarantees).

## Milestone specification

````markdown
## S04 - Prompt Sanitizer Invisible-Unicode Hardening

**Goal.** Keep prompt-safe summaries free of invisible or direction-changing
Unicode that could hide instructions from human reviewers while remaining
visible to authoring models.

**Scope.**

- Extend the prompt sanitizer to remove or visibly escape Unicode format
  characters (general category Cf, including bidi embeddings/overrides/isolates
  and zero-width characters), the Unicode tag block (U+E0000-U+E007F), and
  variation selectors, and report a diagnostic when it does so.
- Apply the same rule to every sanitized field (identifiers, text,
  collections, maps, request fields, security summaries).
- Do not change budgets, truncation, or JSON shapes beyond the added
  diagnostic.

**Review provenance.** "apitools Code Review" Pass 3 finding P1, source High,
local P1 (exploitable hidden-instruction path into authoring prompts);
revalidated at `e015231ef651bb92ee1aab9a4e8fc870d94ad020` with a clean
worktree by probing `SanitizeOperationSummaries` (bidi, zero-width, and tag
characters survive unchanged; ANSI/BEL are stripped). Lineage: S01
(completed; not reopened).

**Dependencies.** None. Ordered first as the only P1 finding.

**Parallel ownership.** `prompt_safety.go` and its tests.

**Downstream impacts.** OpenUdon drafting prompts receive sanitized summaries;
its suite must pass against the workspace checkout. Publishing and re-pinning
consumers need separate authorization.

**Acceptance.** Regression tests prove bidi overrides, zero-width characters,
tag characters, and variation selectors are removed or escaped with a
diagnostic in operation and inventory summaries, while ordinary non-ASCII
text is preserved. `go test ./...`, `go vet ./...`, `git diff --check`, and
`(cd ../openudon && go test ./...)` pass.
````

## Status record

````markdown
# Status S04 - Prompt Sanitizer Invisible-Unicode Hardening

State of each S04 milestone item. Update as items complete. See
[`milestone.md`](milestone.md) for milestone definitions and for the status ID
pattern that names this file.

Status markers:

| Symbol | Suggested Status | Interpretation |
|---|---|---|
| `[ ]` | Pending | Item is actionable when its dependencies pass. |
| `[+]` | Completed | Item finished or done. |
| `[~]` | In Progress | Item is currently being worked, under the ownership rules in `milestone.md`. |
| `[!]` | Blocked | A current unresolved blocker still prevents progress. |
| `[X]` | Cancelled | Item is no longer needed. |
| `[-]` | Closed Historical | A consumed failed attempt or superseded row retained for audit; it is never retried and does not block its accepted successor. |

Write markers with backticks, exactly as in the table above. Each table row is
a commit unit: after changing its state, verify the change, update the memory
bank/docs, and make a scoped `git commit` before starting the next row. A
`[-]` row's notes must record the consumed attempt or supersession and
identify its accepted successor.

After the milestone's review, consolidation, and downstream reconciliation
pass, follow [the retirement procedure](milestone.md#long-term-memory-and-retirement).
Preserve this complete document and the full milestone specification in its
retired record.

Review provenance: "apitools Code Review", Passes 2-4, revalidated at
`e015231ef651bb92ee1aab9a4e8fc870d94ad020` with a clean worktree.

Iteration 2 review provenance: Review "apitools uncommitted-diff review" (2026-09-26; source priority not supplied), baseline `6ee935b` plus the uncommitted implementation diff; revalidated at `6ee935b328234287cc8dfa097884a67fdebccfbd` including those uncommitted changes. Finding IDs U1-U15 follow the review's order.

**Review.** passed
**Review iterations.** 3 (passed; iteration 1 passed as recorded in the restored retirement record)
**Review findings.** Iteration 2: the external uncommitted-diff review, counted as a confirming full pass at the user's direction, found no P1/P2 or higher-severity issue in S04's scope; the only related note is a Lower cosmetic item routed to Candidate Directions. Iteration 3: closing whole-milestone review re-traced `sanitizeOperationSummary`/`sanitizeOperationFields`/`sanitizePromptString` end to end against the acceptance criteria (bidi embeddings/overrides/isolates and zero-width characters via the Cf category, both variation-selector ranges, and the U+E0000-U+E007F tag block are stripped with a diagnostic; ordinary non-ASCII text such as "café"/"漢字" is preserved) and found no P1/P2 or higher-severity issue. Focused sanitizer tests and full `go test ./...`, `go vet ./...` pass. Gate closes at iteration 3; retiring alongside this pass.

| Item | State | Notes |
|---|---|---|
| Strip or escape invisible and direction-changing Unicode | `[+]` | Pass 3 P1, source High, local P1. `sanitizePromptString` now replaces Unicode format controls, tag characters, and variation selectors with spaces. Operation and inventory summary tests confirm U+202E, U+200B, tag characters, and both variation-selector ranges are removed, natural non-ASCII text is preserved, and sanitization emits diagnostics. Lineage: S01. |
| Regression tests and consumer verification | `[+]` | Focused Unicode operation/inventory tests, apitools `go test ./...`, `go vet ./...`, and `git diff --check` pass. With user approval, OpenUdon's test-only registration fixture now anchors its synthetic observation at current UTC minus two minutes; its focused test, full `go test ./...`, and `git diff --check` pass. No OpenUdon product code changed. This goal run has `COMMIT_POLICY: none`. |

## Bounded Review Gate

Iteration 1 passed after review of the complete S04 diff, sanitizer call paths,
diagnostics, test behavior, downstream fixture correction, and acceptance.
Resolved the original P1; no new P1, P2, or higher-severity findings. No lower
findings carried.

## Reconciliation

Updated `architecture.md` to record the prompt sanitizer's Unicode guarantees.
Product scope and tools did not change; no lesson addition or evolution bump is
warranted. OpenUdon's required consumer suite passes after its approved test-only
fixture correction; no additional pending downstream work was found.
````
