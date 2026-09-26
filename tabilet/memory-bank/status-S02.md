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

**Review.** open
**Review iterations.** 4 (open; iterations 1-3 passed as recorded)
**Review findings.** Iteration 4: the external uncommitted-diff review, counted as the next full pass at the user's direction, found P2 regressions U8, U9, U10, and U13 (fix rows below) and Lower cleanups in touched code (unused `candidateDigests`, unused FileInfo from `resolveLocalScanRoot`).
Earlier iterations: Iteration 1: P2 fixed — offline URL validation now rejects empty hostnames; `TestOfflineCacheURLRejectsEmptyHostname` covers it. Lower fixed — uncached imports omit cache metadata; the report test covers it. Iteration 2: P1 fixed — restored the historical exported field shapes for `ImportedSpec`, `LocalOptions`, `LocalResult`, and `DiscoveryCandidate`; added compile-time unkeyed-literal coverage. Cache age is exposed through additive `Client.ImportWithReport`, and digest deduplication uses internal aligned metadata. Focused and full apitools tests passed. Iteration 3: whole-milestone review found no P1/P2 or higher-severity issue across root anchoring, offline cache paths, local scanner bounds/partial errors, digest-idempotent imports, URL truncation, exported compatibility, and docs. Full apitools, OpenUdon, and Udon checks passed.

| Item | State | Notes |
|---|---|---|
| Anchor symlink checks at caller-chosen roots | `[+]` | Pass 1 #1, source High, local P2. Resolve symlinked ancestors of the selected local/artifact root once, reject symlinked final roots and descendants, and retain canonical paths for scans, transactions, and SQLite cache setup. Added temp-parent regressions across local scans, discovery, artifact I/O, inventory, and cache opening. Focused root, `internal/artifactio`, and `sqlitecache` tests plus `git diff --check` pass. Lineage: M72. |
| Offline cache mode without network | `[+]` | Pass 1 #2 (Medium, P2) and #10 (Low, Lower; decision: offline ignores TTL). Offline import now performs syntax-only URL checks, skips DNS/network, loads old cached documents independent of TTL while retaining integrity validation, and reports `cache_stored_at`/`cache_age` through additive `Client.ImportWithReport` and CLI output; legacy `Import` and `ImportedSpec` shapes remain unchanged. Read-write/fetch paths retain unsafe-host rejection. Tests cover unresolvable and private cache keys, expired SQLite records, age reporting, corrupt records, and CLI output. README and architecture document the behavior; root, SQLite, and CLI suites plus `git diff --check` pass. Lineage: M72. |
| Rebuild `LocalFiles` on the bounded walker | `[+]` | Pass 1 #5 (Medium, P2) and #6 (Medium, Lower). `LocalFiles` now shares the `DiscoverLocalSources` bounded walker in OpenAPI-only mode while retaining loose/draft metadata, score ordering, and digest deduplication. It uses shared 10,000-entry/100-candidate defaults and preserves the existing configurable per-file byte limit; explicit count bounds remain available on the additive family-aware discovery API. A truncated scan returns partial results with an error. Symlink, oversized, stat, read, and walk failures become per-entry rejections so valid siblings remain available. Focused scanner tests and full `go test ./...` pass; `git diff --check` passes. Lineage: M69. |
| Idempotent project-URL import | `[+]` | Pass 1 #4 (Medium, P2). `writeUniqueFile` now safely reuses a same-name regular file only when bounded content comparison proves it identical; differing content retains collision-safe suffixes. Local and imported discovery candidate reports deduplicate using internal SHA-256 metadata without changing `LocalResult` or `DiscoveryCandidate` field shapes. Added identical-content reuse/collision, digest-dedup, and exported-shape compatibility tests. Full `go test ./...` and `git diff --check` pass. Lineage: M69, M70. |
| Project-URL limit without discovery failure | `[+]` | Pass 0 F09 and F10 (source priority not supplied; local P2). More than 16 unique URLs now sets `Truncated` and a warning, then imports the first 16 in extraction order without failing discovery; local candidates and successful URL imports remain available alongside per-URL diagnostics. Existing `ImportProjectURLs` and `ImportProjectURLsWithReport` signatures are unchanged; `ImportProjectURLsReport` exposes truncation. Tests cover bounded attempts, local-candidate preservation, warning diagnostics, and compatibility function types. Full `go test ./...` and `git diff --check` pass. Lineage: M70. |
| Docs and consumer verification | `[+]` | README and architecture now describe selected-root ancestor resolution, bounded OpenAPI-only `LocalFiles`, per-entry rejection behavior, internal digest-idempotent imports, the 16-URL truncation warning, offline cache provenance through `ImportWithReport`, and preserved legacy exported shapes. Verified apitools `go test ./...`, `go vet ./...`, and `git diff --check`; OpenUdon `go test ./...` and `git diff --check`; and Udon `go test ./...`. OpenUdon's only change is the separately approved test-fixture timestamp correction. No commits, pushes, or repinning. |
| Report truncation only when a candidate is dropped | `[ ]` | U8, local P2. `DiscoverLocalSources` (`local_source_discovery.go`) fires the candidate limit with `>=` on the 100th accepted candidate, so `LocalFiles` on exactly 100 unique specs returns 100 results plus "scan is incomplete" (probe). Fire only when an additional candidate would be dropped. |
| Keep partial local candidates and later discovery steps | `[ ]` | U9, local P2, plus Lower cleanup of unused `candidateDigests`. `discoverOpenAPIWithDigests` (`discovery.go`) returns `nil` on local truncation, so `DiscoverWithReport` yields 0 candidates and never attempts URL/APIs.guru (probe: 100-spec directory). Keep partial results, record the local attempt and diagnostic, and continue. |
| Unreadable scan root is an error | `[ ]` | U10, local P2, plus Lower cleanup of the unused FileInfo from `resolveLocalScanRoot`. A mode-000 root now returns 0 results and no error (probe); HEAD returned permission denied. Keep descendant walk errors as rejections, but fail on an unreadable root. |
| Offline mode keeps DNS-free unsafe-host checks | `[ ]` | U13, local P2 (weakens the default-on unsafe-host policy; exploitation needs a planted or previously unsafe cache entry). `validateCacheURL` (`download.go`) skips `localhost` and private-IP-literal rejection in offline mode, and the default-rejection test was deleted. Restore those checks without DNS and restore the test. Spec amended; the earlier "syntax-only" wording was the planning gap. |
