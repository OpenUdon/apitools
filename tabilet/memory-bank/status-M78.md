# Status M78 — Stage 1 operation-metadata remediation

**State:** Active. Review gate has not started.

**Provenance:** 2026-09-27 UWS/APItools/OpenUdon stage 1 review R1 and R4; source P2, local P2, both confirmed against clean APItools `30c6f3bd57001e521804a520e5aa21cb7dd95047`. Historical owner M77 is retired. Reproductions exercise `operation_effect.go`, `operation_candidates.go`, and the public `BuildOperationCandidates` API. The review baseline is the same commit; no uncommitted repository changes were part of the evidence.

| Item | State | Notes |
| --- | --- | --- |
| M78.1 Compound effect evidence | `[+]` | R1 fixed: scan later action verbs across the bounded documented text, retain negation and mixed-action `unknown`, and keep read contracts ineligible. Added direct and public candidate regressions; `GOWORK=off go test ./... -count=1` and `git diff --check` passed. |
| M78.2 Nullable response compatibility | `[ ]` | Correct R4; preserve or diagnose nullable response evidence without changing legacy wire shapes, add OpenAPI 3.0 regression plus candidate fixture coverage and documentation/current truth in the same commit. |

## Acceptance and review

Both rows and M78 verification in `milestone.md` must pass before a whole-milestone review. Persist the review iteration here before each pass (maximum ten). No gate iteration has started during review intake.
