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

Iteration 4 review provenance: Review "apitools uncommitted-diff review" (2026-09-26; source priority not supplied), baseline `6ee935b` plus the uncommitted implementation diff; revalidated at `6ee935b328234287cc8dfa097884a67fdebccfbd` including those uncommitted changes. Finding IDs U1-U15 follow the review's order.

**Review.** open
**Review iterations.** 4 (open; iterations 1-3 passed as recorded in the restored retirement record)
**Review findings.** Iteration 4: the external uncommitted-diff review, counted as the next full pass at the user's direction, found P2 regression U1 (fix row below) and a Lower duplicate of `isPartialOperationSecurity` inlined in `cmd/apitools/main.go`.

| Item | State | Notes |
|---|---|---|
| Report resolved reference kind and protocol | `[+]` | Pass 2 C1, source Medium, local P2. `ResolvedReference` now additively reports the selected `SpecKind` and protocol for built-in references, declaration-based OpenAPI metadata for explicit inputs, and matching spec evidence for security classifications. Inspect/advisory headings identify non-OpenAPI protocols; details include kind/protocol. Tests cover OpenAPI/Swagger, human docs, Smithy, Discovery, backward-compatible JSON, and inspect/advisory text+JSON. `go test ./catalog ./cmd/apitools`, targeted vet, and `git diff --check` pass. Lineage: M05, M14. |
| Flag partial operation security coverage | `[+]` | Pass 2 C2, source Medium, local P2. When root security is undeclared, at least one operation names a scheme, and other operations lack explicit declarations, the audit reports `partial-operation-security` unless a more specific missing/undeclared-scheme status takes precedence; the artifact and provider report retain the follow-up. Root declaration presence and alternative count are separate; explicit anonymous root/operation arrays are preserved, and operation declarations are separate from scheme-bearing operation counts. Text output surfaces partial follow-ups; JSON retains structured status and counts. Tests cover partial coverage, anonymous root/operation security, follow-up text, and summary reporting. Targeted package/CLI tests, vet, and `git diff --check` pass. Lineage: M58, C01. |
| Catalog and consumer verification | `[+]` | Re-run after review fixes on 2026-09-26: apitools `go test ./...`, `go vet ./...`, `go run ./cmd/apitools catalog check` (zero errors/warnings), and `git diff --check` pass; OpenUdon `go test ./...` passes with the additive `ResolvedReference` fields and security-audit report additions. No catalog content, pins, or external state changed. |

## Bounded Review-Fix Gate

- Iteration 1 of 10 started: 2026-09-25 22:55 UTC.
- Review scope: the complete C03 resolver/advisory/security-audit implementation and tests, CLI output, architecture/README updates, and this milestone's changes.
- Finding P2: root `security: [{}]` explicitly permits anonymous access, but `RootSecurityRequirementCount` counted named schemes, so this valid nonempty root requirement could be mistaken for absent root security and spuriously report partial operation coverage. Fixed by counting root requirement alternatives while retaining named-scheme inspection separately. Added a regression proving the anonymous root alternative prevents the false partial status/follow-up. README, architecture, and milestone acceptance now state this distinction. Focused tests, CLI follow-up test, targeted vet, and `git diff --check` pass.
- Iteration 2 of 10 started: 2026-09-26 00:00 UTC. Review scope: full C03 implementation and diff, including the iteration 1 correction and its tests/docs.
- Finding P2: anonymous security declarations were still handled inconsistently. An operation-level `security: [{}]` had no named scheme and was counted as an operation without a declared requirement, which could spuriously report partial coverage; scheme-less anonymous root/operation declarations could also be mislabeled as missing metadata or requirements. Fixed by separately tracking root declaration presence, root alternatives, explicit operation security arrays (including `[]` and `[{}]`), and scheme-bearing operations; scheme metadata is required only when a declared requirement names a scheme. Added regression coverage for explicit empty and anonymous declarations at both levels. Focused package/CLI tests, vet, and diff check pass.
- Iteration 3 of 10 started: 2026-09-26 00:10 UTC. Review scope: complete C03 implementation and diff, including the iteration 2 anonymous-policy/count correction, report JSON, CLI rendering, tests, and docs.
- Iteration 3 review completed: no P1/P2-or-higher findings. The bounded review-fix gate passes at iteration 3; full acceptance verification passed afterward on 2026-09-26 (see the verification row above).
| Anonymous declarations do not satisfy requirement checks | `[ ]` | U1, local P2, plus Lower removal of the CLI's inline copy of `isPartialOperationSecurity`. The missing-metadata and missing-requirements cases in `auditOpenAPISecurityArtifact` (`catalog_security_audit.go`) test `OperationSecurityDeclarationCount == 0`, which counts `security: []`; one anonymous `/health` turns missing-requirements into has-security-metadata, hiding the gap the AGENTS.md overlay rule targets. Count only scheme-bearing requirements for those checks. |
