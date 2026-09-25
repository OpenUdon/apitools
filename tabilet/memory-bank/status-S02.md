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

| Item | State | Notes |
|---|---|---|
| Anchor symlink checks at caller-chosen roots | `[ ]` | Pass 1 #1, source High, local P2. Root-anchored checks reject any symlinked ancestor (`local_read.go` `lstatLocalPathNoSymlinks`; `internal/artifactio` `ensureDirectory`/`verifyDirectoryChain`; `sqlitecache.Open`). Revalidated by probe through a symlinked parent. Resolve the chosen root's ancestors once; keep rejecting symlinks beneath it. Test with a symlinked temp parent. Lineage: M72. |
| Offline cache mode without network | `[ ]` | Pass 1 #2 (Medium, P2) and #10 (Low, Lower; decision: offline ignores TTL). `downloadSpecWithCache` → `validateCacheURL` → `rejectHost` does DNS; probe: cached import fails with "no such host". Offline: syntax-only URL checks, serve integrity-checked copy regardless of TTL, report stored age. Lineage: M72. |
| Rebuild `LocalFiles` on the bounded walker | `[ ]` | Pass 1 #5 (Medium, P2) and #6 (Medium, Lower). `local.go` aborts on one oversized file/symlink/unreadable entry (probe: `results=0` + error) and has no visit/byte bounds. Reuse the `DiscoverLocalSources` walker filtered to OpenAPI; walk errors become rejections in both scanners. Lineage: M69. |
| Idempotent project-URL import | `[ ]` | Pass 1 #4 (Medium, P2). `DiscoverWithReport` → `ImportURL` → `writeUniqueFile` (`O_EXCL`) adds a copy per run; OpenUdon `synthesize.go` calls it every non-local run. Reuse identical content or a stable name; dedupe candidates by digest. Lineage: M69, M70. |
| Project-URL limit without discovery failure | `[ ]` | Pass 0 F09 and F10 (source priority not supplied; local P2). `discovery.go` returns `DiagnosticError` for >16 URLs even with local candidates, and `ImportProjectURLsWithReport` drops the error. Decision: import the first 16 unique URLs in order with a truncation diagnostic; no public signature change. Lineage: M70. |
| Docs and consumer verification | `[ ]` | Update README and `architecture.md` for delivered behavior; run `go test ./...`, `go vet ./...`, `git diff --check`, and `(cd ../openudon && go test ./...)`. Publishing/re-pinning needs separate authorization. |
