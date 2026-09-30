# Status M81 — Catalog discovery foundations

**State:** Execution started, 2026-09-30, under the confirmed Kinet Stage 5
goal with `COMMIT_POLICY: task`. M81.1–M81.4 are complete; M81.5–M81.6 pending. Closing review has not started.

**Specification:** [M81](milestone.md#m81--catalog-discovery-foundations).

**Provenance:** Review "APItools M80 — Catalog discovery API for step
contracts" (`apitools-m80-review.md`, 2026-09-30). The review's own finding
IDs are used below. Review baseline `8580ff2`; revalidated at
`8580ff2485a3faff139b17bdc0b8e78d2af4ce18`, with the uncommitted M80 planning
files (the reviewed plan) as part of the evidence and no uncommitted code.
Source priorities are the review's stated severities; local severities were
classified independently under milestone.md's review severity definitions.
The user approved every disposition, the six design decisions in the
specification, this milestone, the M80 amendments, and the planning-file
actions on 2026-09-30. Execution and release authority remain separate.

**Lineage:** Builds on the retired
[M77](../docs/history/status-M77.md), [M78](../docs/history/status-M78.md),
and [M79](../docs/history/status-M79.md) operation-metadata foundations
without reopening them. M80 depends on this milestone.

## Findings owned here

| Finding | Source priority | Local severity | Disposition | Evidence | Owner |
|---|---|---|---|---|---|
| P2-1 Catalog scale exceeds the reused ranking budgets | P2 | P2 | confirmed | `inventory.go:11` (10,000 operations); `operation_rank.go:51-60`; `operation_candidates.go:58-61,182-192`; local cache of about 222 MB and roughly 40,000 operations | M81.4; M80.2 reads the index |
| P2-2 Oversized official specs cannot become candidates | P2 | P2 | confirmed | `client.go:19` (20 MiB); `operation_candidates.go:49-51`; `microsoft-graph-v1-openapi` 35.4 MiB, used by `microsoft-graph` and `microsoft-outlook`; `internal/artifactio` bound is 128 MiB | M81.3 |
| P2-3 Consumers have no defined artifact root | P2 | P2 | confirmed; default settled by decision 3 | `.gitignore` excludes cached specs and `cache.sqlite`; `catalog/builtin_data.go:10` embeds only catalog metadata; iCoT defaults to a sibling checkout (OpenUdon `internal/icot/elicitor/catalog.go`) | M81.1, M81.2 |
| P2-4 Returned artifact references have no provisioning path | P2 | P2 | confirmed | `catalog/materialize.go:204-345` exports per provider only; one artifact can serve two providers | M81.5; M80.1 reference |
| P2-5 "Match" is undefined and lexical matching over-matches | P2 | P2 | confirmed | `operation_rank.go:276-297` (one shared term makes purpose compatible); `operation_rank.go:803-826` | Rule in M81.1; implemented in M80.3 |
| P2-6 No provider constraint for iCoT parity | P2 | P2 | confirmed | OpenUdon `internal/icot/elicitor/catalog.go` matches named providers; `catalog_plan.go` limits the model to 16 artifacts; `FindProvider` already resolves multi-word keys | Decision in M81.1; M80.1, M80.3 |
| P2-7 Consumer drafts expect a single "no API" result | P2 | P2 | confirmed | Kinet stage 5 draft plan (OpenUdon M94 and Kinet W10 drafts); `docs/catalog-upgrade-2026-09.md` §9 | M81.1 checkpoint |
| L1 Ties broken by machine path | Lower | Lower | confirmed; required by M80's determinism acceptance | `operation_rank.go:845-846` | M81.4 identity; M80.3 relocation fixture |
| L2 License-filter default unstated | Lower | Lower | confirmed; required by M80's filter acceptance | `catalog/provider.go:68` (free-text `LicenseNote`); 124 of 181 machine-readable references say only "terms apply" | M81.1; M80.1 |

Findings L3-L5 (planning citations and provenance) were fixed directly in the
M80 planning text. Existing candidates C3 and C4 stay in Candidate Directions:
the approved design resolves provider constraints and references by exact key
and ID, so neither is required.

| Item | State | Notes |
|---|---|---|
| M81.1 — Design record and contract checkpoint | `[+]` | Owns P2-7, the decision parts of P2-3/P2-5/P2-6, and L2. Write `docs/catalog-discovery.md`: request with optional provider constraints; the five outcomes; the qualification rule separating a match from a weak or ambiguous result; the outcome-to-consumer-action table (only a scoped no-match is browser-routable; curated catalog facts are evidence only); license defaults (unknowns included and labeled, `license_note` verbatim, no inferred permission); the catalog-stable reference and its round trip into M81.5; the root contract; the index approach; and the changes the proposed OpenUdon M94 and Kinet W10 drafts must absorb. Complete only after the user approves the record. No sibling edits. |
| M81.2 — Catalog root contract and fixture root | `[+]` | Owns the implementation part of P2-3. Depends on M81.1. Add a caller-supplied root option (cache directory, artifact registrations, index location) and the documented no-root behavior. Document how an operator prepares a root (refresh, then index). Commit a synthetic, redistributable fixture root with no third-party provider specs, usable by APItools tests and later consumer conformance fixtures. |
| M81.3 — Large registered artifact limits | `[+]` | Owns P2-2. Depends on M81.2. Allow the index path only to read and parse catalog-registered, digest-verified artifacts up to 128 MiB, with separately reviewed time, memory, and structural budgets. Keep the 20 MiB default for every existing parser entry point and `BuildOperationCandidates`. Record measurements from an opt-in, provider-free local check against `microsoft-graph-v1-openapi` and `cloudflare-api-openapi`; CI uses synthetic large fixtures. Over-limit artifacts fail closed as unexamined scope. |
| M81.4 — Digest-bound operation index | `[+]` | Owns P2-1 and L1. Depends on M81.3. Add a `catalog index` command and library builder that digest-verify registered artifacts and store sanitized, source-backed operation metadata under catalog-stable identity with no absolute paths. Index shared artifacts once with every provider link. Record the catalog identity and per-artifact coverage (indexed, missing, digest mismatch, oversize, parse failure, unsupported). Rebuilds are byte-identical and root relocation leaves the index unchanged. Readers verify the index version and registration digests and report stale entries as unexamined. No network; existing cache readers keep working without migration. |
| M81.5 — Artifact-scoped export | `[ ]` | Owns P2-4. Depends on M81.2. Add a new function or options type that exports or materializes selected artifacts by reference, verifies expected digests, includes provider- or spec-scoped security overlays and provenance, handles shared artifacts, and fails closed on a mismatch. Existing exported struct shapes, unkeyed composite literals, and provider-level export behavior stay unchanged. |
| M81.6 — Documentation and compatibility verification | `[ ]` | Depends on M81.1-M81.5. Update README and architecture for the delivered behavior. Run the verification below and pass the persisted ten-iteration review gate with no open P1/P2-or-higher findings. No separate publication; M80's authorized publication carries M81. |

## Dependencies and ownership

The APItools active order is M81 then M80. Task order is
M81.1 -> M81.2 -> M81.3 -> M81.4 -> M81.5 -> M81.6; M81.5 depends only on
M81.2 but stays sequential under one execution owner. No parallel
implementation or sibling write scope is approved. M81 has no upstream
prerequisite. M80 starts after every M81 row is complete and its review gate
has closed.

Kinet's stage 5 draft plan lists only APItools M80 in its cross-package order;
adding M81 before M80 there is outside this repository and was reported, not
applied.

## Acceptance and verification

Every row's outcome holds, and existing exported APIs, JSON shapes, and
`BuildOperationCandidates` behavior and ordering are unchanged.

Required verification:

```bash
go test ./...
GOWORK=off go test ./...
go vet ./...
GOWORK=off go vet ./...
go run ./cmd/cataloggen -check
go run ./cmd/apitools catalog check
go run ./cmd/apitools catalog index --help
git diff --check
```

Run OpenUdon and Udon compatibility checks in workspace and standalone modes
per tech-stack.md, with Ramen excluded. Default checks stay provider-free and
offline; the large-artifact measurements are opt-in local checks.

## Review and planning evidence

### M81.1 execution checkpoint — 2026-09-30

The user confirmed Kinet's complete Stage 5 launch request, including this
repository's M81 → M80 order and `COMMIT_POLICY: task`. One execution owner
selected M81.1; M81.2–M81.6 and all M80 rows remain pending.

Prepared [the discovery design record](../../docs/catalog-discovery.md) against
clean baseline `fdc0a3f2647a0c4488f0c0a6334fb2647aa73447`. Read the current
instructions, current truth, M81/M80 specifications/statuses and relevant
operation-candidate, ranking, registry, artifact safety and export contracts;
inspected OpenUdon's catalog hints and CatalogPlan as read-only consumer
evidence. The record includes all six approved decisions, the stronger purpose
qualification rule, five-outcome precedence, explicit root/read-only index
policy, native reference/export round trip, scoped security/license evidence,
remote bounds and required M94/W10 reconciliation.

Focused verification passed: `git diff --check`; all design-record Markdown
file links resolve; a coverage check found the required decision/outcome/root/
identity/consumer topics; the APItools ledger has exactly one in-progress task,
M81.1. These are document checks, not implementation acceptance, and no new
code or runtime verification is claimed. No sibling file was changed. The
design record and ledger updates remain uncommitted until the user's required
approval completes M81.1; then make its scoped task commit before M81.2.

**Approved:** The user approved the design record on 2026-09-30. M81.1 is
complete; its scoped task commit precedes M81.2. This is not a closing-review
iteration. No implementation or publication acceptance is implied.

### M81.2 — explicit root and read-only registrations

Delivered additive `catalog.RootOptions`/`RootPaths`/`ResolveRoot` with no
implicit root, confined non-overlapping registry/index paths, symlink and
special-file rejection, resolved ancestors and no creation. Added
`sqlitecache.ReadCatalogArtifacts` with read-only snapshot/schema validation,
no migration/pruning/access-time writes, cancellation and five-second,
10,000-row/32-MiB bounds. Shared row validation preserves the existing cache
reader's behavior. The redistributable synthetic fixture links two providers
to one exact OpenAPI artifact and includes an explicit new-root preparation
recipe; no SQLite database or provider specification is committed.

Passed `go test ./...`, `GOWORK=off go test ./...`, both workspace/standalone
`go vet ./...`, focused root/cache tests after the final validation tightening,
`git diff --check`, and the fixture preparation/identity/existing-root-refusal
smoke at `/tmp/apitools-m81-2-9pg353kh/root`. Tests cover no writes, legacy
schema without migration, newer schema refusal, unsafe paths, missing evidence,
registry bounds and cancellation. Initial fixture validation caught missing
source notes and an incorrect machine-spec availability label; both were
corrected before acceptance. No discovery/index behavior is claimed yet.

### M81.3 — registered large-artifact index path

Delivered private digest/size-verified registered artifact parsing, reusing the
existing inventory and native candidate adapters. Only the new private OpenAPI
index path raises bytes to 128 MiB; ordinary direct parsers and
BuildOperationCandidates retain their exact 20-MiB defaults and wires. Internal
explicit structural limits share sourceguard's checks; YAML nodes are reused
without parsing twice. Reviewed bounds: depth 100, four million structural
items, 100,000 operations, 256 MiB inventory/candidate metadata and a 45-second
cooperative context deadline. Retained metadata is bounded during construction;
structural limits and measured memory do not claim a hard process RSS ceiling.
Decode/normalization cannot be interrupted mid-call; results after expiry are
refused. Unsupported large non-OpenAPI families retain their old limits.

Passed workspace/standalone tests and vet, focused synthetic >20-MiB accepted
registered-source versus ordinary-parser refusal, digest/registration mismatch,
>128-MiB sparse-file refusal, aliases, trailing JSON, version/depth/structural
limits and cancellation, plus diff checks. Opt-in local measurements read only
existing cache registrations and bytes (no credentials/network): Cloudflare
21,930,908 bytes, 3,008 visited operations / 2,745 candidates, 666 diagnostics,
5.727 seconds, 430,364 KiB maximum process RSS; Graph 37,110,274 bytes, 16,422
operations / 10,897 candidates, 7,801 diagnostics, 40.111 seconds, 870,608 KiB
maximum process RSS. Evidence: `/tmp/apitools-m81-3.Zcbpg8mQ` result/time files.
Existing schema-summary limits leave incomplete evidence; M81.4 must preserve
omitted operations and diagnostics as unexamined coverage, never fully indexed
scope or definitive no-match. No official source bytes were committed.

### M81.4 execution evidence — 2026-09-30

Delivered the additive index types, offline builder, validated atomic writer,
source-free reader and `catalog index` CLI. Shared raw identity parses once
with all provider links; complete coverage preserves missing, digest mismatch,
oversize, parse failure and unsupported scope. Snapshot drift removes saved
candidates and marks coverage stale. Unknown versions, malformed identities,
incomplete/duplicate coverage, false authority, wrong registration bindings,
unsafe paths, secret provenance and input/sidecar overwrite are refused.
Limits: three cooperative minutes, 10,000 registrations, 500,000 candidates,
512 MiB index; preflight depth 100/32 million JSON tokens before typed decode.
No network, schema migration or legacy exported shape change.

APItools workspace/standalone tests and vet, generator freshness, catalog
check, index help and diff checks pass. Focused regressions cover identity,
coverage, drift/cancellation, source deletion without reader reparsing,
invalid OpenAPI coverage and byte-identical input permutation/root relocation.
CLI fixture preparation and two deterministic rebuilds passed at
`/tmp/apitools-m81-4-a362a70b/root` (1 artifact, 2 candidates, no gaps).
OpenUdon workspace/standalone full tests pass; Udon workspace full tests pass.
Udon's literal standalone command requests module metadata updates unrelated
to this additive API. The standalone suite passed offline with an isolated
modfile at `/tmp/udon-m81-4-standalone-i0or7vo4`; its recorded delta resolves
local replacements' testify 1.12.1 and YAML 3.0.5 requirements. Udon go.mod,
go.sum and all sibling worktrees were left unchanged. Repeat this qualification
at final review and retain the distinction from literal readonly success.

**Review iterations started:** 0 of at most 10; closing gate not started.
This is ordinary approved intake, not a bounded-gate pass. Persist the count
once implementation reaches review and never reset it. An unresolved blocking
finding at iteration 10 requires user direction.

Assessment used read-only inspection of the code and local cache named above
and a read-only harness validation of the planning files. No implementation
verification or milestone review is claimed.
