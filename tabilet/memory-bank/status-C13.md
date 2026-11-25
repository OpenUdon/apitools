# C13 — Reviewed catalog findings

**Owner.** APItools. **Dependency.** Accepted [C09](../docs/history/status-C09.md).
**Parallel-safe: no.** Accepted C13 gates [C12](status-C12.md).
See the [milestone specification](milestone.md#c13--reviewed-catalog-findings)
for acceptance, verification and authorization requirements.

## Planning approval and deferral

On 2026-11-25 the owner replied “I approve” directly in tab 1 to the C13
planning proposal, including the FND-04 deferral and the listed file actions.
FND-04 retains the [C07 evidence gaps](../docs/history/status-C07.md#status-record)
for the 30 unreadable publisher documentation probes. Reconsider the
deferral when a named provider's documentation gap matters to approved work.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| C13.1 — Reviewed search relevance | `[+]` | Fixes FND-07. Record before-and-after precision and recall on fixed reviewed labels and exact selectors. Require genuine singular and plural matches, including `send slack messages`. Reject named incidental false positives for queries such as `list user` and `time`, and retain rejected stem cases including `accounting` and `intuit quickbook`. Preserve exact-provider ranking, metadata suggestions and source provenance. |
| C13.2 — Public fallback reporting | `[+]` | Fixes FND-08. If an implicit public lookup fails and reviewed suggestions are returned, emit a clear stderr diagnostic in both text and JSON modes. Preserve reviewed stdout and the successful exit status. Use offline fixtures for empty successful lookup, failed lookup, explicit-source behavior and publisher-scoped routing. |
| C13.3 — Google Slides scope evidence | `[ ]` | Fixes FND-05. Replace blanket read-only advisory coverage with method- and path-scoped evidence for `presentations.get`, `presentations.pages.get`, `presentations.pages.getThumbnail`, `presentations.create` and `presentations.batchUpdate`. Use the retained Google Discovery snapshots to verify the reviewed presentation scopes and record alternatives, ambiguities and source revision. Other methods and the wider root-scope class remain unreviewed. If evidence is insufficient, pause this row for separate Class R authority. |

## Review

Review gate iterations: 0

The owner launched C13 by approving the raw launch SHA-256
`f2d3b9aaa0b756a5c992d196ec3cf7ae4f4d8e930c948d5197916bbaae9d087e`.
The stock protocol SHA-256 is
`4facf0102dd672fc43ff6e4a79384b3a5d0bbfe77a06ed2e55c401a8650d05ae`.
APItools has no project addendum. C13 starts from published
`c5072d43a83e9601a0965ca770669fa8e29c19bb` and the approved planning edits.
The first task commit carries that planning context. Execution is serial,
with task commits, local verification and no external effects. Retain
verification attempts and review counters separately.

Whole review pending, limit 10. Follow the
[model roles](milestone.md#monthly-reviewed-catalog),
[test policy](test-policy.md) and
[milestone review procedure](milestone.md#milestone-review-procedure).
Required-gate recovery limit is 3. Record elapsed wall-clock time and persist
[required-gate recovery](milestone.md#required-gate-recovery) state before each
interactive Goal/Next required-gate dispatch.

At retirement, retarget the C12 dependency and the register's C13 links,
including the FND-04 approval source, to the frozen C13 record.

C13.1 required gate `c13-1-a1` started from the launch baseline with
search code, test and guide diff SHA-256
`6863dd1dfb50c93a275cefac70579493cfe658a8c32f011980d817b4f1f96f95`.
Run offline build, vet, whole root and CLI tests, and diff validation.
[Gate receipt](../../dist/c13-goal/c13-1-gates.log).

The fixed query-selector matrix went from precision 7/12 and recall 7/10
to precision 10/10 and recall 10/10 across 15 named cases. An intermediate
candidate matched a pluralized brand name in a label. The existing
`intuit quickbook` test caught the false promotion. The corrected candidate
requires exact provider identity and passes that test.
[Baseline](../../dist/c13-goal/c13-1-relevance-before.log) and
[corrected matrix](../../dist/c13-goal/c13-1-relevance-verified.log).

Gate c13-1-a1 passed build, vet, whole root and CLI tests (57.526s)
and diff validation. The first Astra xhigh task review blocked C13.1
because the matrix omitted named FND-07 false positives for NVIDIA
announcements and Plaid permissions. The reviewer also found ordinary
short nouns missing from the inflection match. The correction treats
user and product as evidence when they are direct operation targets or
another query word supplies the actual resource. It accepts short pairs
such as tag/tags and job/jobs and preserves the rejected new/news stem.

C13.1 corrective gate `c13-1-a2` started at the same launch baseline
with applied code, test and guide diff SHA-256
`aa9e7b77c024b4332e2c9694763f592e9cfa0895344b8f5deedc4ef28cdc0e66`.
Run offline build, vet, whole root and CLI tests, and diff validation.
[Gate receipt](../../dist/c13-goal/c13-1-corrected-gates.log).
The fixed 28-case matrix on clean pre-change source measured precision
11/19 and recall 11/20. With the correction it measured precision 20/20
and recall 20/20. The named negative selectors and previous stem tests
passed. [Baseline](../../dist/c13-goal/c13-1-expanded-final-before.log)
and [corrected matrix](../../dist/c13-goal/c13-1-expanded-after.log).

Gate c13-1-a2 passed build, vet, whole root and CLI tests (41.981s)
and diff validation. The second Astra xhigh task review confirmed the
named fixes and found three related P2 edges. Preserve verified npm and
Reddit authenticated-user matches. Keep broad words generic on both sides
of inflection, and prevent a quantifier from independently promoting a
qualified product match. Add exact-selector cases for the three edges.

C13.1 corrective gate `c13-1-a3` started from the same baseline with
applied code, test and guide diff SHA-256
`4c66c18081d755737bed4f0cc6fce28a776c4c6bb55761e03587d8bf8ce7fe17`.
Run offline build, vet, whole root and CLI tests, and diff validation.
[Gate receipt](../../dist/c13-goal/c13-1-final-gates.log).

Gate c13-1-a3 passed build, vet, whole root and CLI tests (39.099s)
and diff validation. The Astra xhigh re-review confirmed the three P2
fixes and identified the compound Facebook user-or-page target as a
remaining recall regression. Treat an explicitly stated alternative
target as direct evidence and retain the Facebook selector in the matrix.

C13.1 corrective gate `c13-1-a4` started from the same baseline with
applied code, test and guide diff SHA-256
`ecf7a7c6bad149edf02e12c3e592ecdd24bf88404514415c5413e3f09d45e294`.
Run the final fixed matrix, offline build, vet, whole root and CLI tests,
and diff validation. [Gate receipt](../../dist/c13-goal/c13-1-ready-gates.log).

Gate c13-1-a4 passed build, vet, whole root and CLI tests (76.902s)
and diff validation. The final 35-case matrix improved from 15/24
precision and recall to 24/24 for both. Existing stem refusals and
exact-identity regressions passed. [Baseline matrix](../../dist/c13-goal/c13-1-ready-before.log).
The fresh Astra xhigh task re-review passed the exact c13-1-a4 candidate
with no findings. FND-07 is fixed for the accepted scope.

C13.2 required gate `c13-2-a1` started from
`964110dad8448d1d63c393a00570ec15b3491a79` and scoped diff SHA-256
`b69be095e4c71df8bea59eb5eabbed4a32a2733940f324e3e853d1483ca97ae1`.
Run focused routing checks, offline build, vet, whole CLI tests, and
diff validation. FND-08's red test shows empty stderr in both output modes
when the implicit local fixture public lookup fails. The corrected check
compares successful-empty and failed lookup stdout byte for byte.
[Baseline](../../dist/c13-goal/c13-2-before.log) and
[gate receipt](../../dist/c13-goal/c13-2-gates.log).

C13.2 gate c13-2-a1 passed build, vet, the focused routing checks,
whole CLI tests (23.229s) and diff validation. Astra xhigh passed
the behavior and found a Lower test assertion gap. Check for the
warning even when it is appended to an explicit-source error.

C13.2 corrective gate `c13-2-a2` started at the same base with
scoped diff SHA-256
`40e778d80aa2f40b209a15a254fe34cea85507dceeeeb6d13822ba2813d0b5b3`.
Run offline focused routing checks, vet, whole CLI tests and diff
validation. The product build from c13-2-a1 used the same Go source.
[Gate receipt](../../dist/c13-goal/c13-2-verified-gates.log).

Gate c13-2-a2 passed the focused routing tests, vet, whole CLI
tests (15.431s) and diff validation. Fresh Astra xhigh re-review
passed the exact corrected test and behavior. FND-08 is fixed.
