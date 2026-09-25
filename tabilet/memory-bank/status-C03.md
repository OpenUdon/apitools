# Status C03 - Catalog Resolution And Security Audit Accuracy

State of each C03 milestone item. Update as items complete. See
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

Review provenance: "apitools Code Review", Passes 2-4, revalidated at
`e015231ef651bb92ee1aab9a4e8fc870d94ad020` with a clean worktree.

| Item | State | Notes |
|---|---|---|
| Report resolved reference kind and protocol | `[ ]` | Pass 2 C1, source Medium, local P2. `preferredSpecReference` (`catalog/resolve.go`) fills the `openapi` reference with human-docs (157), Smithy (29), Discovery (21), Stone (1), and index (1) sources; `ResolvedReference` has no kind. Add kind/protocol (additive JSON) and protocol-aware CLI/advisory labels. Lineage: M05, M14. |
| Flag partial operation security coverage | `[ ]` | Pass 2 C2, source Medium, local P2. `auditOpenAPISecurityArtifact` (`catalog_security_audit.go`) returns has-security-metadata when root security is empty and only some operations declare security. Add a partial-operation-security status and follow-up. Lineage: M58, C01. |
| Catalog and consumer verification | `[ ]` | Run `go test ./...`, `go vet ./...`, `go run ./cmd/apitools catalog check`, `git diff --check`, and `(cd ../openudon && go test ./...)`. |
