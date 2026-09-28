# Retired milestone M79 - Correct operation effect and nullable-output metadata

**Milestone.** M79
**Outcome.** completed
**Retired.** 2026-09-28
**Source status.** tabilet/memory-bank/status-M79.md
**Source specification.** tabilet/memory-bank/milestone.md#m79--correct-operation-effect-and-nullable-output-metadata
**Evidence.** 2c7036614872f29acace56ddb7e3bec75302a672
**Worktree.** includes uncommitted changes
**Review.** passed
**Review iterations.** 1
**Verification.** APItools `go test ./... -count=1`, `go vet ./...`, CLI help checks, and `git diff --check` passed. OpenUdon tests passed against the exact published APItools module `v0.0.0-20260928033144-e3625f6ef52e`; Kinet W04.6's focused empty-package journey and checkout-backed consumer check passed against OpenUdon M89 implementation `2e2ecedb32add53e10d6d81d4b2528ca2bdfdcc8`. OpenUdon M89's retirement closure was published at `0c7c5d33da2ba7b190954b9eb402cc14b5ec1f73`.
**Consolidated into.** [milestone index](../../memory-bank/milestone.md), [catalog upgrade plan](../../../docs/catalog-upgrade-2026-09.md), and [knowledge journal](knowledge.md). No product, architecture, stack, or evolution change was required.

## Milestone specification

````markdown
## M79 — Correct operation effect and nullable-output metadata

Correct the leading-action classifier so an unrecognized leading verb cannot
be skipped in search of a later read verb; retain `unknown` when the action is
not understood. Determine nullable-output compatibility from the selected
declared/mapped response fields and their schema ancestors, not from unrelated
fields elsewhere in the response. Keep the existing conservative behavior for
unknown effects. Acceptance covers mutating summaries and operation IDs,
unrecognized actions, unrelated nullable fields, and selected nullable output
fields. Publish the accepted revision for OpenUdon M89 to consume and verify
the OpenUdon consumer behavior against that exact revision.

**Dependencies.** No technical upstream dependency. The cross-package order
requires this milestone's accepted, published revision before OpenUdon's
dependent Stage 1 milestone begins. Findings F1/F2 are from the 2026-09-28
independent Stage 1 review.
````

## Status record

````markdown
# Status M79 — Correct operation effect and nullable-output metadata

**State:** Complete. The accepted M79 revision is published and consumed by OpenUdon M89; the review gate passed in iteration 1.

**Provenance:** “Stage 1 review: Kinet, APItools, OpenUdon, UWS” (`stage1-review.md`, 2026-09-28), findings F1 and F2 (source/local P2). Review and revalidation baseline: APItools `26bb05247d6c48f8ee60b9ae178f6ef6d48bbe3d`; worktree clean. The review records sibling Kinet and OpenUdon baselines and a clean worktree for each. Earlier M78 remediation is completed and retired; M79 does not reopen it.

| Item | State | Notes |
|---|---|---|
| M79.1 — Classify the leading documented action conservatively | `[+]` | F1, source/local P2. Do not skip an unrecognized leading verb and classify a later “read” noun/verb as proof of a read operation. Return `unknown` unless the actual leading action is recognized as read or write; retain conservative unknown behavior. Added regressions for `markMessageRead`, `Mark message as read`, and `saveSearch`. Focused `go test ./... -run 'TestAssessOperationEffect' -count=1` passed. Evidence: `apitools/operation_effect.go`. |
| M79.2 — Scope nullable response checks to selected outputs | `[+]` | F2, source/local P2. Evaluate nullability only on declared/mapped response fields selected by the operation contract and their schema ancestors; an unrelated nullable field must not make every output indeterminate. Selected nullable outputs remain non-compatible unless their contract can prove otherwise. Candidate v1 now preserves nullable state on root schemas, response fields, and output metadata. Focused selected, unrelated-sibling, root-ancestor, and nested-parent regressions passed. Full APItools tests and vet passed; OpenUdon and Udon consumer suites also passed, with OpenUdon additionally tested using a temporary modfile replacing its pinned APItools revision with this worktree. Evidence: `apitools/inventory_request.go`, `apitools/operation_rank.go`. |
| M79.3 — Qualify and publish the consumer contract | `[+]` | User authorized the APItools commit and publication needed for M79 on 2026-09-28. Review iteration 1 passed with no P1/P2-or-higher findings. APItools uncached tests, vet, CLI help, and diff check pass; Udon's full suite passes against its local replacement. APItools code revision `e3625f6ef52ea54b7f78b7a4a4f1993bf8a06a46` is pushed to `origin/main` and resolves as `v0.0.0-20260928033144-e3625f6ef52e`. OpenUdon's full suite passes against that exact module version using a temporary modfile. OpenUdon M89 consumed this exact revision and completed/published its reviewed implementation at `2e2ecedb32add53e10d6d81d4b2528ca2bdfdcc8`; Kinet W04.6's exact consumer passed against that revision. OpenUdon's M89 retirement record is published at `0c7c5d33da2ba7b190954b9eb402cc14b5ec1f73`. |

**Acceptance.** Mutating actions are never silently classified as `read` merely because a later word matches; unrecognized actions remain `unknown`. Nullable fields outside the selected outputs do not block those outputs, while selected nullable fields do not become falsely compatible. Run `go test ./...`, `go vet ./...`, `git diff --check`, `go run ./cmd/apitools search --help`, and `go run ./cmd/apitools import --help`. If exported APIs change, also run the AGENTS.md dependent checks in OpenUdon and Udon. OpenUdon consumes the exact accepted, published revision and passes its consumer/boundary checks.

**Dependencies.** No technical upstream dependency. The coordinated order requires this package to accept and publish before the OpenUdon consumer milestone begins; see Kinet `docs/kinet-order.md`. M78 remains frozen in history.

**Review iteration 1 of 10 started:** 2026-09-28, before reviewing the full M79 implementation and diff. Authorization to commit and publish received from the user on 2026-09-28.

**Review iterations started:** 1 of at most 10.

**Review findings.** Iteration 1 reviewed the full M79 diff from baseline `26bb05247d6c48f8ee60b9ae178f6ef6d48bbe3d`, including conservative leading-action classification, native-family exceptions, schema-reference/ancestor nullability propagation, selected-output scoring, tests, and the public contract/memory updates. No P1, P2, or higher-severity finding remains. APItools `go test ./... -count=1`, `go vet ./...`, CLI help, `git diff --check`, OpenUdon `go test ./...` with a temporary local APItools replacement and again against exact published module version `v0.0.0-20260928033144-e3625f6ef52e`, and Udon `go test ./...` passed. The review-fix gate passes at iteration 1. No evolution bump: this is a corrective extension within the existing operation-candidate contract and Stage 1 scope; ownership and architecture boundaries are unchanged.
````
