# Retired milestone M78 - Stage 1 operation-metadata remediation

**Milestone.** M78
**Outcome.** completed
**Retired.** 2026-09-27
**Source status.** tabilet/memory-bank/status-M78.md
**Source specification.** tabilet/memory-bank/milestone.md#m78---stage-1-operation-metadata-remediation
**Evidence.** b585540f33d487c0956ffab92a01c9ce9750eeae
**Worktree.** includes uncommitted changes
**Review.** passed
**Review iterations.** 1
**Verification.** Focused public regressions; `GOWORK=off go test ./... -count=1`; `GOWORK=off go test -race ./... -count=1`; `GOWORK=off go vet ./...`; OpenUdon `go test ./internal/stepauthoring ./cmd/openudon -count=1` with a temporary Go workspace pointing at this APItools revision; and `git diff --check` passed.
**Consolidated into.** [operation-candidate contract](../../../docs/operation-candidates.md); no product, architecture, tech-stack, or reusable-lesson change was needed. No evolution bump.

## Milestone specification

````markdown
## M78 - Stage 1 operation-metadata remediation

**Goal.** Correct two confirmed P2 findings from the 2026-09-27 UWS/APItools/OpenUdon stage 1 review before OpenUdon repins the APItools consumer contract. This is M77 remediation, not a reopening of its retired record.

**Scope and acceptance.** Compound operation wording that explicitly includes both read and mutation must yield `unknown`, including a read verb followed by a distant connector. OpenAPI nullable response fields must not earn compatible output type/requiredness evidence when a non-nullable step contract cannot represent null. Expose a source-grounded gap and indeterminate output match, preserve exact operation/source identity and auth metadata, and keep existing public version and legacy wire shapes. Add focused public-API regressions and conformance fixture coverage where applicable.

**Order and dependency.** [M78.1 and M78.2](status-M78.md) run in order. The existing M77 publication is the prerequisite. OpenUdon M88 consumes the completed M78 revision; Kinet W03 is a downstream consumer. No UWS wire version or provider operation is changed.

**Verification and closeout.** Run focused regressions, `GOWORK=off go test ./...`, `GOWORK=off go test -race ./...`, `GOWORK=off go vet ./...`, downstream OpenUdon tests against the candidate revision, `git diff --check`, and the persisted whole-milestone review gate. Publish the reviewed commit before OpenUdon pins it.
````

## Status record

````markdown
# Status M78 — Stage 1 operation-metadata remediation

**State:** Complete. Whole-milestone review passed in iteration 1.

**Provenance:** 2026-09-27 UWS/APItools/OpenUdon stage 1 review R1 and R4; source P2, local P2, both confirmed against clean APItools `30c6f3bd57001e521804a520e5aa21cb7dd95047`. Historical owner M77 is retired. Reproductions exercise `operation_effect.go`, `operation_candidates.go`, and the public `BuildOperationCandidates` API. The review baseline is the same commit; no uncommitted repository changes were part of the evidence.

| Item | State | Notes |
| --- | --- | --- |
| M78.1 Compound effect evidence | `[+]` | R1 fixed: scan later action verbs across the bounded documented text, retain negation and mixed-action `unknown`, and keep read contracts ineligible. Added direct and public candidate regressions; `GOWORK=off go test ./... -count=1` and `git diff --check` passed. |
| M78.2 Nullable response compatibility | `[+]` | R4 fixed: selected OpenAPI/Swagger response schemas with nullable fields emit a readiness issue, partial output capability, and an indeterminate zero-point output match. Added an OpenAPI 3.0 public-candidate regression and documented the limitation. `GOWORK=off go test ./... -count=1`, `GOWORK=off go test -race ./... -count=1`, `GOWORK=off go vet ./...`, OpenUdon authoring tests against the local APItools revision, and `git diff --check` passed. |

## Acceptance and review

Both rows and M78 verification in `milestone.md` must pass before a whole-milestone review. Persist the review iteration here before each pass (maximum ten).

**Review iteration 1:** Started after both rows and required verification passed. Reviewed the full milestone diff from `30c6f3bd57001e521804a520e5aa21cb7dd95047`, including public candidate behavior, response selection, conservative effect classification, compatibility scoring, fixtures, docs, and OpenUdon consumer checks. No P1/P2 finding remains within M78 acceptance. The output score is conservatively zero for the entire selected response when any inspected field permits null; this is explicit in the gap and avoids a false positive. No product direction, architecture boundary, or contract direction changed; no evolution bump. Current behavior is recorded in `docs/operation-candidates.md`.
````
