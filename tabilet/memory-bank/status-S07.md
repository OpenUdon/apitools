# S07 — Catalog lifecycle safety and compatibility remediation

Approved serial execution: human integration plan, 2026-10-15.
Source review: “Remote M87–M89 integration review,” baseline
`3157924bc65cc23e1e5fd19bdb03b294bd40a5ed`; merged revalidation `1aa23f6`.
All findings confirmed; source/local severity agree. Frozen M87–M89 and
M86 review 0/10 remain unchanged. No goal launcher or audit is activated.

Use the Ubuntu v2.13 [status and review rules](milestone.md).
Each row is a verified commit unit; zero or one row may be in progress.
Protected publication is authorized only by the human plan after qualification.

| Item | State | Notes |
|---|---|---|
| S07.1 — Transactional snapshots | `[+]` | F01/F02 P1, F05 P2: destination symlinks, destructive rollback, FIFO/source symlinks. Reuse artifactio transactions; preserve existing targets on failure, identical reuse, empty/new targets, collisions, drift and destination-only normalization. Red: preservation/symlink tests reproduced F01/F02. Green: go test ./catalog ./internal/artifactio, including bounded Linux FIFO refusal, empty/identical reuse and existing journal/collision gates; git diff --check. |
| S07.2 — Validate lifecycle mutations | `[+]` | F03/F04 P1: provider/initialization traversal. Validate paths/identities before writes, reject symlinks/special files, reuse identical sources and reject different collisions; registration failure preserves prior artifacts. Verified focused TestCatalog Init/Add/Import/Snapshot/Provider/Source gates and diff check; offline traversal, symlink, collision and registration rollback cases pass. |
| S07.3 — Nested registry and environment | `[+]` | F06/F09 P2: additive root-aware SQLite opener and explicit flag > canonical env > aliases > defaults. Verified nested init/add/index/search/snapshot, root-relative registry reads, legacy opener semantics and canonical/alias/explicit/empty/default precedence with focused tests; existing SQLite/CLI gates pass. |
| S07.4 — Explicit custom catalogs | `[+]` | F07/F08 P2: bounded --catalog plumbing for add/import/search; one selected catalog governs validation/registration/index/search. Unknown identities fail before writes. Verified offline custom add/import/index/search round trip with nested env-selected registry/index, wrong-catalog identity refusal, metadata fallback and unknown-provider/spec no-write/no-download cases; strict bounded decoder and all Catalog CLI gates pass. |
| S07.5 — Model response fallback | `[+]` | F10 P2: reject null/empty/missing/wrong-type/over-budget plans while preserving valid fences, normalized effects and missing-input clarification. Red: invalid-response tests reproduced suppressed fallback for null/empty/missing/oversized plans. Green: focused Plan/InvalidModel/ModelPlan/ClientLLM/PrepareAndValidate tests cover structural and field budgets, fences, normalized effects and clarification. |
| S07.6 — Combined qualification | `[+]` | Full/standalone tests/vet, focused race, generation/quality, CLI, OpenUdon/Udon disposable checks, docs, review and retirement. No Ramen. Required gates pass at 9fd831d; [qualification](../../docs/s07-qualification.md). Temporary modfiles cover consumer nested GOWORK=off; synthetic Udon test outputs moved to diagnostic storage. Consumer heads, worktrees and pins match baselines. Review/retirement follow below. |

## Review

Whole-milestone review: 1/10, completed with blocking findings corrected;
iteration 2 pending. Base: `1aa23f69ad0d25a1b9d8bf6c8cf5a037a62aa3ce`.
Iteration 1 reviewed the full diff from 1aa23f6 through 7fdf2d3, task log,
uncommitted qualification updates, source/test failure semantics, acceptance
fidelity, docs, consumer compatibility and frozen boundaries. Blocking findings:

- R1-model (P2): “over-budget responses” remains incomplete for recommendation
  completions: extractJSON bounds only the extracted object; missing rationale
  can also suppress fallback. Apply complete-response/required-field validation.
- R1-evidence (P2): “check drift” requires deterministic byte/add/remove drift
  evidence; the existing drift-named test only creates an active journal before
  invocation. Exercise the actual consistency verifier and failure preservation.
- R1-docs (P2): the current M89 lesson says to reject idempotent_write, contrary
  to the required supported normalization. Preserve old wording in the knowledge
  journal and correct the maintained lesson.

Safety review also requires source/destination sidecar hierarchy collision
coverage. No public signature, source family, remote policy or owner expansion.
Fixes remain under the same review counter; no evolution bump.

Iteration 1 resolution: R1-model, R1-evidence and R1-docs are corrected.
Focused red/green recommendation evidence, workspace and standalone tests/vet,
focused race gates, generation/quality and full OpenUdon/Udon compatibility
rechecks pass. Final consumer sources are immutable baseline archives with
temporary modfiles; both also pass vet. One OpenUdon archive run returned
private-storage setup refusal (attempt 1 at that stage); unchanged focused
diagnostics passed and the full suite passed with an isolated TMPDIR. No tests
were skipped or acceptance weakened. Original consumer heads/status/pins remain
unchanged. Qualification records the evidence and environment retro candidates.
