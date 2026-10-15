# S07 catalog lifecycle qualification

The integration preserves local `8b60f796150dc2dbfe940670e0932f3bf1b864f2`
and remote `3157924bc65cc23e1e5fd19bdb03b294bd40a5ed` through merge
`1aa23f6`. Ubuntu v2.13 harness rules and GOAL.md deletion remain.
M87–M89 records and review counts are frozen; M86 remains pending, review 0/10,
with its C11/M54 triggers and an S07 compatibility prerequisite.

## Regression evidence

| Findings | Offline evidence |
|---|---|
| F01/F02/F05 | Snapshot existing-target preservation, symlink refusal, bounded Linux FIFO probe, empty/new target, identical reuse, collisions, journals and normalization tests in catalog; artifactio transaction checks. |
| F03/F04 | Catalog init traversal/symlink no-write tests; add/import provider traversal before client construction; collision and registration rollback tests. |
| F06/F09 | Nested init/add/index/search/snapshot and custom import; root-relative registry read; unchanged legacy cache opener; canonical/alias/explicit/empty/default flag precedence. |
| F07/F08 | Explicit custom-catalog add/import/index/search round trip; mismatched catalog identity refusal; selected metadata fallback; unknown identity no-write/no-download tests. |
| F10 | Invalid model completion red/green tests; null/empty/missing/wrong-type responses, total byte and field/collection/rune budgets, fences, idempotent_write normalization and missing-input clarification. |

Snapshot preservation tests reproduced overwritten targets and destructive
rollback before repair. Invalid-model tests reproduced suppressed fallback
before repair. Other cases were revalidated against source and offline tests.

## Verification

Initial implementation qualification at `9fd831d182b9ceed63900545c0dd2886276738bd`
(review fixes are requalified below):

- `go test ./...` and `go vet ./...`: pass in the existing workspace.
- `GOWORK=off GOPROXY=off go test ./...` and standalone vet: pass.
- `go test -race ./catalog ./sqlitecache ./internal/artifactio ./cmd/apitools -run 'Test(Snapshot|Catalog|CustomCatalog|OpenCatalog|DirTransaction)'`: pass.
- `go test -race . -run 'Test(Plan|InvalidModel|ModelPlan|Recommend|CatalogDecision|CatalogDiscovery)'`: pass.
- `go run ./cmd/cataloggen -check`: generated outputs current.
- `go run ./cmd/apitools catalog check`: zero errors and warnings.
- Fifteen offline CLI smoke commands: search/import help, six catalog lifecycle
  help commands, specs/security-audit/advisory/resolve/overlay-view/stats/
  refresh-report. Cache-report smoke also passes with an empty disposable
  cache fixture.
- `git diff --check`: pass.

## Consumer configuration

Compatibility checks use disposable configuration with APItools pointing at this
integration checkout. No sibling pins or worktree content are edited; Ramen is
excluded. Initial standalone attempts lacked cached sibling revisions; an
alternate workspace then exposed OpenUdon's child `GOWORK=off` check. Temporary
modfiles with local dependency replacements and inherited absolute `GOFLAGS`
cover that nested command, without changing its code or skipping its check.

Captured consumer baselines:

- OpenUdon: `53191bc0d9fb6bfdae130e3af41f28beaa3e39b2`, clean.
- Udon: `d2df23793da9c667838ebdb5b74947dda5c42862`, clean.

Disposable logs/configuration are under
`/tmp/apitools-s07-verification.r6Xkug`; this is diagnostic storage, not a
permanent source contract or consumer adoption/pin change.

Both complete consumer suites pass with GOWORK=off, GOPROXY=off and
GOFLAGS='-mod=mod -modfile=/tmp/apitools-s07-verification.r6Xkug/<consumer>/candidate.mod'.
The temporary modfiles replace APItools and available workspace dependencies
with their inspected local checkouts. Both consumers also pass vet under the
disposable workspace. Consumer heads, clean status and go.mod/go.sum SHA-256
hashes match the captured baselines. Udon tests generated untracked synthetic
HCL fixtures under three relative srv/ trees; their shape/timestamps were checked
and the invocation-owned trees were moved into disposable diagnostic storage.
This proves local candidate compatibility, not ordinary published adoption.

Required qualification passes. Whole review passed at iteration 2 after the
iteration 1 fixes. Final implementation: `fb94a77db9e17c9750bdf39e0ae7fe4df318109f`.

## Review iteration 1 fixes

Complete-response and required-rationale validation now also protects candidate
recommendations. Its regression reproduced acceptance of a missing rationale
before the fix. Deterministic local consistency checks cover content changes,
source removal/addition and symlink replacement. A corrupt SQLite normalization
failure preserves existing destination files and removes staging. Sidecar
hierarchy and symlink-aliased overlap tests reject unsafe destinations.
The outdated maintained effect lesson is corrected; its original wording is
retained in the knowledge journal. Focused tests pass.

Final consumer runs use immutable archive copies of the captured baseline
commits in disposable source directories, so generated fixture outputs stay
outside the original sibling worktrees. Their temporary modfiles still point
at the final APItools integration checkout.

## Environment retro

P2 planning candidate: a deterministic consumer-check wrapper should propagate
temporary modfiles to nested Go commands, use immutable source copies, and
compare original heads/status/pin hashes afterward. This run demonstrates the
recipe; it creates no new milestone or workflow requirement. Lower candidate:
cache-report smoke should use a disposable empty cache rather than parse a
developer's populated review corpus. No audit or new capability was enabled.

No evolution bump: S07 restores existing contracts and preserves package
ownership, source families, remote opt-in, consumer pins and frozen evidence.

Final workspace/standalone APItools tests and vet, both focused race commands,
generation/quality and both complete consumer test/vet suites pass after the
review fixes. OpenUdon's first archive-based run reported private-store setup
refusal; the unchanged focused package passed, and a complete run with explicit
isolated TMPDIR passed. No consumer test was skipped or modified.

S07 is accepted and retired after two review iterations. The
[full retired specification/status](../tabilet/docs/history/status-S07.md)
retains all six rows and review evidence. M86 remains pending at review 0/10
with its original operational triggers and accepted S07 compatibility input.
M87–M89 records remain frozen. Publication follows the authorized normal
fast-forward main operation; no unrelated milestone is launched.
