# Status M81 — Catalog discovery foundations

**State:** Execution started, 2026-09-30, under the confirmed Kinet Stage 5
goal with `COMMIT_POLICY: task`. All six task rows are complete. Closing review passed in iteration 3; downstream reconciliation and retirement follow the task commit.

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
| M81.5 — Artifact-scoped export | `[+]` | Owns P2-4. Depends on M81.2. Add a new function or options type that exports or materializes selected artifacts by reference, verifies expected digests, includes provider- or spec-scoped security overlays and provenance, handles shared artifacts, and fails closed on a mismatch. Existing exported struct shapes, unkeyed composite literals, and provider-level export behavior stay unchanged. |
| M81.6 — Documentation and compatibility verification | `[+]` | Depends on M81.1-M81.5. Update README and architecture for the delivered behavior. Run the verification below and pass the persisted ten-iteration review gate with no open P1/P2-or-higher findings. No separate publication; M80's authorized publication carries M81. |

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

### M81.5 execution evidence — 2026-09-30

Delivered additive `CatalogArtifactReference`, `CatalogArtifactExportOptions`
and `ExportCatalogArtifacts` with no legacy exported shape changes. Selection
joins canonical provider/spec/artifact registrations and verifies catalog
identity plus actual raw digest/bytes. Shared bytes copy once while every
selected operation/provider link retains native-selector provenance. Only
provider-wide and matching spec-scoped security overlays are emitted, keeping
OR-of-AND alternatives unchanged. Advisory sources remain explicitly labeled;
registry private metadata and URL userinfo/query/fragment are excluded.
Atomic directory transactions preserve collision/reuse/Force semantics and
refuse partial publication, source-root overlap, symlink ancestors, unsafe
paths, cancellation, mismatched sources and registration drift. The downstream
OpenUdon binder still checks native selectors. No network/parser/execution or
credential handling was added. Fixed bounds are documented in README.

Focused tests prove selected-only export despite extra provider registrations,
shared physical copy, retained provider provenance, overlay scope/alternatives,
private metadata exclusion, identical reuse, collision/no partial Force,
integrity/binding failure, cancellation, registry drift and confinement.
APItools workspace/standalone tests and vet and diff checks pass. OpenUdon and
Udon full workspace consumer suites pass against the additive implementation.
Full final standalone qualification and closing review belong to M81.6.

**Review iterations started:** 3 of at most 10. Iteration 1 started on
2026-09-30 against `c79d94b30a2ca1798dcd851f083b77432b65a5e3` plus this
selection note. Review the full M81 design/root/parser/index/export change and
resume this same iteration after any interruption. Findings are not yet final;
the milestone is not accepted. An unresolved blocking finding at iteration 10
requires user direction.

Original planning assessment used read-only inspection of the code and local cache named above
and a read-only harness validation of the planning files. No implementation
verification or milestone review is claimed.

### Closing review iteration 1 — findings persisted before fixes

Full milestone review at `c79d94b30a2ca1798dcd851f083b77432b65a5e3`:

- P2 R1: read-only registry aggregate accounting omits selected `updated_at`.
  SQLite accepts oversized nonnumeric timestamp text even in an INTEGER column,
  allowing a malformed registry to exceed the documented pre-decode text bound.
  Owner: M81.6; include timestamp bytes and add a refusal regression.
- P2 R2: the index CLI compares untrimmed `--index` against `--catalog`, while
  root/artifact confinement normalizes whitespace. A whitespace-padded output
  spelling can bypass the custom-catalog input overwrite check. Owner: M81.6;
  compare the same normalized paths and test the dispatched CLI refusal.
- Lower R3: builder cancellation at an empty registry can be ignored by a
  custom installation callback; add an explicit final context check. Fix here.

The gate has not passed. Fix these findings and rerun affected verification
before starting iteration 2. No sibling work or new product scope is needed.

### Closing review iteration 2 — started 2026-09-30

R1 now counts timestamp bytes before decoding; the oversized SQLite timestamp
regression passes. R2 compares normalized CLI paths and its dispatched
whitespace-overwrite refusal passes. R3 has an explicit final context check
and the cancellation regression with an empty custom snapshot passes. Focused
registry/index/CLI checks pass after fixes. Iteration 2 reviews the entire M81
range plus these uncommitted corrections, not only the findings patch.
No acceptance is claimed until the full gate and final checks finish.

### Closing review iteration 2 — finding persisted before fix

R1–R3 remain resolved. P2 R4: an installation-selected registration callback
can supply a whitespace-padded source path. Artifact I/O trims it, while the
index destination overlap guard compares the raw spelling, permitting an
input overwrite for that malformed callback snapshot. SQLite's default adapter
normalizes paths already, but the additive library must refuse this ambiguity.
Owner: M81.6. Reject noncanonical whitespace-padded relative source/output
paths and test the library guard before starting iteration 3. No gate pass yet.

### Closing review iteration 3 — started 2026-09-30

R4 rejects whitespace-padded noncanonical source paths before indexing or
publication; the library input-preservation regression passes. Focused index,
export, registry and CLI checks pass after all fixes. Iteration 3 reviews the
full M81 range plus the persisted corrections; final verification is running.

### Closing review iteration 3 — passed; M81.6 complete

Full M81 review found no remaining P1/P2-or-higher findings or Lower findings
requiring carry-forward. R1–R4 are resolved with focused regressions. Rechecked
shared legacy helper extraction, parser defaults/exported wires, deterministic
identity and coverage, missing/partial/stale semantics, bounded work,
confinement and atomic failure, private metadata omission, selector ownership,
fixture redistribution and operator documentation.

Final APItools workspace/standalone full tests and vet, generator freshness,
catalog check, index help and diff checks pass after corrections. OpenUdon
workspace/standalone suites pass, and isolated standalone adoption explicitly
replacing APItools with the M81 worktree also passes at
`/tmp/openudon-m81-final-adoption-457mg5y_`. Udon workspace suite passes;
standalone passes with disposable module metadata at
`/tmp/udon-m81-final-standalone-ffwnbp0t`, preserving its documented preexisting
module-update distinction. No sibling file changed and Ramen is excluded.
No new evolution version: V25 already records the approved public contract,
boundaries and M81→M80 direction; implementation refines that same scope.
M80 must now reconcile to the actual M81 task commit before starting.
No separate M81 publication; authorized M80 publication carries the reviewed
sources and ordinary closure record. The approved global sequence is unchanged.
