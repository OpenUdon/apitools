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
| S07.4 — Explicit custom catalogs | `[ ]` | F07/F08 P2: bounded --catalog plumbing for add/import/search; one selected catalog governs validation/registration/index/search. Unknown identities fail before writes. |
| S07.5 — Model response fallback | `[ ]` | F10 P2: reject null/empty/missing/wrong-type/over-budget plans while preserving valid fences, normalized effects and missing-input clarification. |
| S07.6 — Combined qualification | `[ ]` | Full/standalone tests/vet, focused race, generation/quality, CLI, OpenUdon/Udon disposable checks, docs, review and retirement. No Ramen. |

## Review

Whole-milestone review: 0/10, not started. Base: `1aa23f6`.
