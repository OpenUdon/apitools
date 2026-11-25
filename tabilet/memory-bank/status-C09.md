# C09 — Monthly cron and agent maintenance skill

**Owner.** APItools. **Dependency.** Accepted [C08](../docs/history/status-C08.md). **Parallel-safe: no.**
C09 owns integrated closure for C06–C09 and the [C10 updater repairs](../docs/history/status-C10.md) against the baseline in [C06.1](../docs/history/status-C06.md#status-record). See the [milestone specification](milestone.md#c09--monthly-cron-and-agent-maintenance-skill).

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| C09.1 — Check readiness and implement the monthly job | `[+]` | Read the same cache directory configuration as a manual update. Preflight binary, database, configured directory, permissions and cron environment. Use absolute paths, local time, one run lock, bounded retries and per-provider checkpoints. Exercise an ordinary cache directory, overlapping jobs, restart and loss of an explicitly configured mount in isolated tests. |
| C09.2 — Expose schedule management and reports | `[ ]` | Default the optional user schedule to `17 3 1 * *`. Provide idempotent enable, disable and status, preserving unrelated crontab entries. Produce readable and JSON reports for accepted changes, pending reviews, size deferrals, failures, disk use and embedding usage. Test with a scratch crontab. This row fixes FND-02. |
| C09.3 — Deliver the agent maintenance skill | `[ ]` | Ship the repository-owned `apitools-catalog-maintenance` skill. Check its instructions against actual commands for preflight, dry-run, schedule management, report inspection, per-provider retry and generating the Git export. Document the operator-initiated first run and how to inspect it before scheduling. |
| C09.4 — Verify repeat run and integrated closure | `[ ]` | Reconcile [C07 first-run evidence](../docs/history/status-C07.md#status-record) for 334 providers and all six special cases. Record the repeat run, unchanged artifact reuse and selective re-embedding. Verify repeat outcomes, bounded retry and inspection of previous runs against accepted C10 evidence. Scope any live publisher or metered embedding check separately. Exercise the operator/agent workflow, review the cumulative diff, and run the existing class O whole-unit gates and affected race and consumer checks. This row fixes FND-03 and FND-06. |

## Review

Review gate iterations: 0

Whole review pending, default limit 10. Reuse accepted evidence for contributor
milestones and identify the exact final state. C08 supplies the separate SQLite
search index, optional hosted embedding budgets, and local evidence that
completed vector batches are reused. C09 still owns the repeat-run proof.
Use the [test policy](test-policy.md)
and [milestone review procedure](milestone.md#milestone-review-procedure).
Required-gate records include elapsed wall-clock time. Persist
[required gate recovery](milestone.md#required-gate-recovery) state here before
each interactive Goal/Next required-gate dispatch.

C09 starts at `3fffb94d1f6d1207e3a2e11dc855b9fbfcc7dbae`; the integrated
baseline remains C06.1 at `150d08c644bd4ae4ce83a8772f317f083580d230`.
The owner selected C09 directly with the existing approved local goal scope.
Stock protocol SHA-256 is
`4facf0102dd672fc43ff6e4a79384b3a5d0bbfe77a06ed2e55c401a8650d05ae`.
There is no package addendum. Preserve the existing disposable launch
reference (SHA-256 `7bc0ef7af40917aa0309f73d81a4561d115e7be03fad36aab8b14a188d10c361`)
without depending on its changed checksum as launch authority. Scope is
serial implementation, local task commits and offline verification.
The review limit is 10 and the required-gate recovery limit is 3.

C09.1 gate attempt `c09-1-a1` started from `3fffb94d1f6d1207e3a2e11dc855b9fbfcc7dbae` and diff SHA-256
`54ced31ad7a1baf3722f7aa0510ad345574c476040a491c68a246e627f18eb82`. Build `./...`, vet `./...`, ordinary whole root,
CLI and internal catalogupdate tests, affected maintenance and storage races,
and diff validation. Existing recovery limit 3. [Row evidence](../../dist/c09-goal/c09-1-gates.log).

C09.1 initial build, vet, ordinary affected units and focused races passed.
The task reviewer found carryover skipped the due month, preflight could
disagree with the writer, and a future unrelated retry could delay ready
work. It also required the whole affected race units. Corrections passed
the focused cases before the whole-unit rerun.

C09.1 corrective gate attempt c09-1-a2 started from
`3fffb94d1f6d1207e3a2e11dc855b9fbfcc7dbae` and exact diff SHA-256
`911b7543b8b11bd887fa8e52d73eab9a360a9f6287ec4042d70f024c82a73cf9`.
Run build and vet ./..., whole ordinary and race tests in root, CLI and
internal/catalogupdate, and git diff --check.
[Corrective verification](../../dist/c09-goal/c09-1-final-gates.log).

C09.1 passed separate Astra xhigh corrective review at diff SHA-256
`911b7543b8b11bd887fa8e52d73eab9a360a9f6287ec4042d70f024c82a73cf9`.
Build, vet, ordinary affected units, whole root/CLI/storage races and diff
checks passed. Ordinary whole tests took 23.270s and races took 171.144s.
The fixture verifies that one invocation finishes the interrupted month
and the due month under one lock, with no repeat requests in the same
completed month. Preflight and writer share schema validation, and
mixed-origin retries select the earliest eligible pending reference.
