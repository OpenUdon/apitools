# C09 — Monthly cron and agent maintenance skill

**Owner.** APItools. **Dependency.** Accepted [C08](../docs/history/status-C08.md). **Parallel-safe: no.**
C09 owns integrated closure for C06–C09 and the [C10 updater repairs](../docs/history/status-C10.md) against the baseline in [C06.1](../docs/history/status-C06.md#status-record). See the [milestone specification](milestone.md#c09--monthly-cron-and-agent-maintenance-skill).

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| C09.1 — Check readiness and implement the monthly job | `[+]` | Read the same cache directory configuration as a manual update. Preflight binary, database, configured directory, permissions and cron environment. Use absolute paths, local time, one run lock, bounded retries and per-provider checkpoints. Exercise an ordinary cache directory, overlapping jobs, restart and loss of an explicitly configured mount in isolated tests. |
| C09.2 — Expose schedule management and reports | `[+]` | Default the optional user schedule to `17 3 1 * *`. Provide idempotent enable, disable and status, preserving unrelated crontab entries. Produce readable and JSON reports for accepted changes, pending reviews, size deferrals, failures, disk use and embedding usage. Test with a scratch crontab. This row fixes FND-02. |
| C09.3 — Deliver the agent maintenance skill | `[+]` | Ship the repository-owned `apitools-catalog-maintenance` skill. Check its instructions against actual commands for preflight, dry-run, schedule management, report inspection, per-provider retry and generating the Git export. Document the operator-initiated first run and how to inspect it before scheduling. |
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

C09.2 required gate attempt c09-2-a1 started at `f47eaa7c5d59b30f29f99a76873de59ecda2ac92` with diff SHA-256
`e654124f97673d05044093a6e4bdc506848800eae2f7f6702eccdcbb80a838e7`. Go build and vet ./..., plain and race whole affected root, CLI and
storage units, and diff validation. [Row checks](../../dist/c09-goal/c09-2-gates.log).

C09.2 initial build, vet, whole ordinary and affected race checks passed. The
review identified empty selector fallthrough, scratch input aliases and combined
publisher-head and baseline facts missing in reports. Corrections are in the CLI.
C09.2 corrective attempt c09-2-a2 started at `f47eaa7c5d59b30f29f99a76873de59ecda2ac92` with exact diff SHA-256
`abd06640f894f688a352bb1a6ae9296a212d559b45339dc8e149e00bb7fe37f3`. Run build/vet, whole CLI ordinary and race tests, and diff verification.
The previously passed root and storage race results cover their identical inputs.
[Corrective gates](../../dist/c09-goal/c09-2-final-gates.log).

C09.2 correction attempt c09-2-a3 started at `f47eaa7c5d59b30f29f99a76873de59ecda2ac92` with diff SHA-256
`8b0d5f1b18ee54899e1e0acacab439651c66c3ed7c031e10bff8271548e0d9ff`. Re-review cleared empty flags, independent head/baseline reporting and
command capture limits. It found path normalization mismatched the alias guard.
Canonical scratch paths now cover existing parent symlinks and reject outer
whitespace before file access. Re-run build/vet, whole CLI ordinary/race and
diff. [Path guard verification](../../dist/c09-goal/c09-2-path-gates.log).

C09.2 passed separate Astra xhigh re-review at SHA-256
`8b0d5f1b18ee54899e1e0acacab439651c66c3ed7c031e10bff8271548e0d9ff`.
Whole affected ordinary checks took 32.155s, and whole races took 175.126s
on the initial library/CLI candidate. Final CLI ordinary and race checks
took 9.729s and 33.304s. Build, vet and diff checks passed. Scratch tests
verified a required completed full run, idempotence, unrelated cron
content, normalization and hard-link conflicts, and both report formats.
FND-02 is fixed. Interrupted output now equals the committed checkpoint
snapshot, and read-only inspection can retrieve those unfinished outcomes.

C09.3 candidate SHA-256 `3941e4b857828f7def68ebcc26ca3fa44ac578ee2538f6b4a8da96d60cc176bf` at `7997f24568155dca30c2d704a3ce863085a2c746`. Verify skill packaging and
installed CLI options, dry-run, preflight, and scratch schedule status. Go code
and test inputs match the whole affected C09.2 pass.
[Skill verification](../../dist/c09-goal/c09-3-gates.log).

C09.3 skill packaging and installed CLI smoke passed, including preflight,
dry-run and scratch status. Earlier scratch enable/disable and persisted
report tests apply to unchanged Go inputs. The task reviewer verified
every documented command and independently traced an interrupted-run
inspection and scratch rehearsal. Astra xhigh PASS at
`3941e4b857828f7def68ebcc26ca3fa44ac578ee2538f6b4a8da96d60cc176bf`.
