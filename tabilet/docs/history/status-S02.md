# Retired milestone S02 - Local, Offline, And Discovery Safety Remediation

**Milestone.** S02
**Outcome.** completed
**Retired.** 2026-09-26
**Source status.** tabilet/memory-bank/status-S02.md
**Source specification.** tabilet/memory-bank/milestone.md#s02---local-offline-and-discovery-safety-remediation
**Evidence.** 8edc7229480cd2e4e59134814c0f5bc8cebe4c18
**Worktree.** clean
**Review.** passed
**Review iterations.** 5
**Verification.** go build ./..., go vet ./..., go test ./... (apitools, all 17 packages); git diff --check; OpenUdon full go test ./...; Udon go build ./.... All passed at the evidence commit.
**Consolidated into.** no current-truth change; current facts already live in README and architecture.md (root anchoring, offline cache behavior, bounded `LocalFiles`, truncation reporting).

## Milestone specification

````markdown
## S02 - Local, Offline, And Discovery Safety Remediation

**Goal.** Make local scanning, offline cache use, and project-URL discovery fail
safely and predictably without weakening the untrusted-source guards.

**Scope.**

- Anchor symlink rejection at the caller-chosen root: resolve a root's own
  ancestors once, then reject symlinks only beneath it, in local reads,
  `internal/artifactio` roots, and `sqlitecache.Open`.
- Make `CacheModeOffline` perform no DNS or network access (URL syntax plus
  the DNS-free unsafe-host checks: `localhost`, scoped hosts, and private or
  reserved IP literals stay rejected unless `AllowUnsafeHosts` is set) and serve
  any integrity-checked cached copy regardless of cache TTL, reporting its
  stored age. Read-write mode keeps TTL.
- Rebuild `LocalFiles` on the bounded local walker: per-file rejections instead
  of whole-scan failure, visit/byte bounds, digest deduplication, and walk
  errors on one entry recorded as rejections in both local scanners.
- Make project-URL import idempotent by reusing identical content or a stable
  name, and deduplicate discovery candidates by digest.
- Keep the 16-URL network bound but never fail discovery on URL count: import
  the first 16 unique URLs in source order, record truncation in the additive
  `ImportProjectURLsReport`, and preserve local and partial URL candidates.
  Keep the legacy `ImportProjectURLs` and `ImportProjectURLsWithReport`
  signatures unchanged.
- Do not relax private-host rejection, add proxy support, or change public
  function signatures.

**Review provenance.** "apitools Code Review" Pass 1 (#1, #2, #4, #5, #6,
#10) and Pass 0 (F09, F10), stated baseline `a5699e5`, revalidated at
`a5699e570fb5d6288fdd255c71d229fb062b50b3`; relevant uncommitted changes at
revalidation were harness documentation only. Lineage: M69, M70, M72, S01
(completed; not reopened).

**Dependencies.** None.

**Parallel ownership.** `local.go`, `local_read.go`,
`local_source_discovery.go`, `download.go`, `import.go`, `discovery.go`,
`internal/artifactio/`, and the `sqlitecache` open path, plus their tests and
docs. Does not own `operationlifecycle/` or catalog refresh behavior.

**Downstream impacts.** C02 (completed; retired, see the [retired C02
record](status-C02.md)) built on the root-anchoring semantics in
`internal/artifactio`. OpenUdon `internal/synthesize` stops accumulating
duplicate imports and no longer fails discovery on long briefs; its suite must
pass against the workspace checkout. Publishing and re-pinning consumers need
separate authorization.

**Acceptance.** Regression tests cover a symlinked temp-parent root, offline
import of a non-resolvable cached hostname and of an entry older than the TTL,
an oversized file and a symlink beside a valid spec, repeated discovery
producing one file and one candidate per document, and more than 16 project
URLs with local candidates present. Review-iteration-4 regressions (U8-U10,
U13) also hold: exactly 100 unique specs are not reported as truncated,
discovery keeps partial local candidates and still attempts URL/APIs.guru
imports when the local limit is reached, an unreadable scan root is an error
rather than an empty result, and offline imports of `localhost` or private IP
literals are rejected by default. `go test ./...`, `go vet ./...`,
`git diff --check`, and `(cd ../openudon && go test ./...)` pass, and README and
`architecture.md` describe the delivered behavior.
````

## Status record

````markdown
# Status S02 - Local, Offline, And Discovery Safety Remediation

State of each S02 milestone item. Update as items complete. See
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

Review provenance: "apitools Code Review", stated baseline `a5699e5`,
revalidated at `a5699e570fb5d6288fdd255c71d229fb062b50b3`; relevant uncommitted
changes at revalidation were harness documentation only.

Iteration 4 review provenance: Review "apitools uncommitted-diff review" (2026-09-26; source priority not supplied), baseline `6ee935b` plus the uncommitted implementation diff; revalidated at `6ee935b328234287cc8dfa097884a67fdebccfbd` including those uncommitted changes. Finding IDs U1-U15 follow the review's order.

**Review.** passed
**Review iterations.** 5 (passed)
**Review findings.** Iteration 4: the external uncommitted-diff review, counted as the next full pass at the user's direction, found P2 regressions U8, U9, U10, and U13 (fix rows below) and Lower cleanups in touched code (unused `candidateDigests`, unused FileInfo from `resolveLocalScanRoot`). Iteration 5: whole-milestone review closing the gate. The two Lower cleanups noted at iteration 4 are now done — `DiscoverWithReport` no longer tracks per-candidate digests past the point they stop being read (`discovery.go`), and `resolveLocalScanRoot` no longer returns an unused `os.FileInfo` (`local_read.go`, `local.go`). Re-reviewed root anchoring, offline cache paths, local scanner bounds/partial-truncation errors (U8-U10), the DNS-free offline host-literal check (U13), digest-idempotent imports, URL truncation, and exported compatibility across the whole milestone: no P1/P2 or higher-severity issue found. `go build ./...` and `go vet ./...` are clean and `go test ./...` passes across all 17 apitools packages; OpenUdon and Udon build cleanly against the workspace module. Gate passes at iteration 5.
Earlier iterations: Iteration 1: P2 fixed — offline URL validation now rejects empty hostnames; `TestOfflineCacheURLRejectsEmptyHostname` covers it. Lower fixed — uncached imports omit cache metadata; the report test covers it. Iteration 2: P1 fixed — restored the historical exported field shapes for `ImportedSpec`, `LocalOptions`, `LocalResult`, and `DiscoveryCandidate`; added compile-time unkeyed-literal coverage. Cache age is exposed through additive `Client.ImportWithReport`, and digest deduplication uses internal aligned metadata. Focused and full apitools tests passed. Iteration 3: whole-milestone review found no P1/P2 or higher-severity issue across root anchoring, offline cache paths, local scanner bounds/partial errors, digest-idempotent imports, URL truncation, exported compatibility, and docs. Full apitools, OpenUdon, and Udon checks passed.

| Item | State | Notes |
|---|---|---|
| Anchor symlink checks at caller-chosen roots | `[+]` | Pass 1 #1, source High, local P2. Resolve symlinked ancestors of the selected local/artifact root once, reject symlinked final roots and descendants, and retain canonical paths for scans, transactions, and SQLite cache setup. Added temp-parent regressions across local scans, discovery, artifact I/O, inventory, and cache opening. Focused root, `internal/artifactio`, and `sqlitecache` tests plus `git diff --check` pass. Lineage: M72. |
| Offline cache mode without network | `[+]` | Pass 1 #2 (Medium, P2) and #10 (Low, Lower; decision: offline ignores TTL). Offline import now performs syntax-only URL checks, skips DNS/network, loads old cached documents independent of TTL while retaining integrity validation, and reports `cache_stored_at`/`cache_age` through additive `Client.ImportWithReport` and CLI output; legacy `Import` and `ImportedSpec` shapes remain unchanged. Read-write/fetch paths retain unsafe-host rejection. Tests cover unresolvable and private cache keys, expired SQLite records, age reporting, corrupt records, and CLI output. README and architecture document the behavior; root, SQLite, and CLI suites plus `git diff --check` pass. Lineage: M72. |
| Rebuild `LocalFiles` on the bounded walker | `[+]` | Pass 1 #5 (Medium, P2) and #6 (Medium, Lower). `LocalFiles` now shares the `DiscoverLocalSources` bounded walker in OpenAPI-only mode while retaining loose/draft metadata, score ordering, and digest deduplication. It uses shared 10,000-entry/100-candidate defaults and preserves the existing configurable per-file byte limit; explicit count bounds remain available on the additive family-aware discovery API. A truncated scan returns partial results with an error. Symlink, oversized, stat, read, and walk failures become per-entry rejections so valid siblings remain available. Focused scanner tests and full `go test ./...` pass; `git diff --check` passes. Lineage: M69. |
| Idempotent project-URL import | `[+]` | Pass 1 #4 (Medium, P2). `writeUniqueFile` now safely reuses a same-name regular file only when bounded content comparison proves it identical; differing content retains collision-safe suffixes. Local and imported discovery candidate reports deduplicate using internal SHA-256 metadata without changing `LocalResult` or `DiscoveryCandidate` field shapes. Added identical-content reuse/collision, digest-dedup, and exported-shape compatibility tests. Full `go test ./...` and `git diff --check` pass. Lineage: M69, M70. |
| Project-URL limit without discovery failure | `[+]` | Pass 0 F09 and F10 (source priority not supplied; local P2). More than 16 unique URLs now sets `Truncated` and a warning, then imports the first 16 in extraction order without failing discovery; local candidates and successful URL imports remain available alongside per-URL diagnostics. Existing `ImportProjectURLs` and `ImportProjectURLsWithReport` signatures are unchanged; `ImportProjectURLsReport` exposes truncation. Tests cover bounded attempts, local-candidate preservation, warning diagnostics, and compatibility function types. Full `go test ./...` and `git diff --check` pass. Lineage: M70. |
| Docs and consumer verification | `[+]` | README and architecture now describe selected-root ancestor resolution, bounded OpenAPI-only `LocalFiles`, per-entry rejection behavior, internal digest-idempotent imports, the 16-URL truncation warning, offline cache provenance through `ImportWithReport`, and preserved legacy exported shapes. Verified apitools `go test ./...`, `go vet ./...`, and `git diff --check`; OpenUdon `go test ./...` and `git diff --check`; and Udon `go test ./...`. OpenUdon's only change is the separately approved test-fixture timestamp correction. No commits, pushes, or repinning. |
| Report truncation only when a candidate is dropped | `[+]` | U8, local P2, fixed. `discoverLocalSources` now checks the candidate limit before accepting a new unique candidate (`local_source_discovery.go`), so truncation fires only when a candidate is actually excluded. Regression `TestDiscoverLocalSourcesExactCandidateLimitIsNotTruncated` proves exactly N candidates is not truncated and N+1 still is; existing candidate-limit tests still pass. Full `go test ./...`, `go vet ./...` pass. |
| Keep partial local candidates and later discovery steps | `[+]` | U9, local P2, fixed. Added `ErrLocalScanTruncated` so `DiscoverOpenAPI`/`discoverOpenAPIWithDigests` return partial candidates alongside a truncation error instead of discarding them; `DiscoverWithReport` treats that error as recoverable, keeps the partial local candidates, records a `discovery.local.truncated` diagnostic, and still attempts URL/APIs.guru. Regression `TestDiscoverWithReportKeepsPartialLocalCandidatesOnTruncation` covers it. The Lower `candidateDigests`-unused cleanup noted alongside this finding was done at iteration 5 (see review header). Full `go test ./...`, `go vet ./...` pass. |
| Unreadable scan root is an error | `[+]` | U10, local P2, fixed. `localDiscoveryState` now tracks the root path it is walking (`currentRoot`); a walk error reported for that exact path fails the scan, while a descendant walk error still becomes a rejection. Regression `TestLocalFilesFailsOnUnreadableRoot` (skips itself if the process can read a mode-000 directory, e.g. running as root). The Lower unused-`FileInfo`-from-`resolveLocalScanRoot` cleanup noted alongside this finding was done at iteration 5 (see review header). Full `go test ./...`, `go vet ./...` pass. |
| Offline mode keeps DNS-free unsafe-host checks | `[+]` | U13, local P2, fixed. Split `rejectHost` into DNS-free `rejectHostLiteral` (empty host, `localhost`, IPv6 zone scope, unsafe IP literal) and the DNS-dependent lookup; `validateCacheURL` now calls `rejectHostLiteral` so offline mode refuses localhost/private literals by default without any network access. Replaced the incorrect `TestOfflineCachedPrivateURLUsesCacheWithoutHostResolution` (which asserted the vulnerable behavior) with `TestOfflineCachedPrivateURLRejectedByDefault` and added `TestOfflineCachedLocalhostURLRejectedByDefault`; the DNS-skip and `AllowUnsafeHosts` tests still pass unchanged. Full `go test ./...`, `go vet ./...` pass. |
````
