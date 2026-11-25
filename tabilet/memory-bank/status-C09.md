# C09 — Monthly cron and agent maintenance skill

**Owner.** APItools. **Dependency.** Accepted [C08](../docs/history/status-C08.md). **Parallel-safe: no.**
C09 owns integrated closure for C06–C09 and the [C10 updater repairs](../docs/history/status-C10.md) against the baseline in [C06.1](../docs/history/status-C06.md#status-record). See the [milestone specification](milestone.md#c09--monthly-cron-and-agent-maintenance-skill).

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| C09.1 — Check readiness and implement the monthly job | `[+]` | Read the same cache directory configuration as a manual update. Preflight binary, database, configured directory, permissions and cron environment. Use absolute paths, local time, one run lock, bounded retries and per-provider checkpoints. Exercise an ordinary cache directory, overlapping jobs, restart and loss of an explicitly configured mount in isolated tests. |
| C09.2 — Expose schedule management and reports | `[+]` | Default the optional user schedule to `17 3 1 * *`. Provide idempotent enable, disable and status, preserving unrelated crontab entries. Produce readable and JSON reports for accepted changes, pending reviews, size deferrals, failures, disk use and embedding usage. Test with a scratch crontab. This row fixes FND-02. |
| C09.3 — Deliver the agent maintenance skill | `[+]` | Ship the repository-owned `apitools-catalog-maintenance` skill. Check its instructions against actual commands for preflight, dry-run, schedule management, report inspection, per-provider retry and generating the Git export. Document the operator-initiated first run and how to inspect it before scheduling. |
| C09.4 — Verify repeat run and integrated closure | `[+]` | Reconcile [C07 first-run evidence](../docs/history/status-C07.md#status-record) for 334 providers and all six special cases. Record the repeat run, unchanged artifact reuse and selective re-embedding. Verify repeat outcomes, bounded retry and inspection of previous runs against accepted C10 evidence. Scope any live publisher or metered embedding check separately. Exercise the operator/agent workflow, review the cumulative diff, and run the existing class O whole-unit gates and affected race and consumer checks. This row fixes FND-03 and FND-06. |

## Review

Review gate iterations: 1

Integrated whole review round 1 started under effective limit 10. Reuse accepted evidence for contributor
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

C09.4 required gate attempt c09-4-a1 started at `47621bd085e02b2be3454a6981b16d34b646f935` on diff SHA-256
`dd9c434a166c8f509bd4e720f779c6c28f422a74ef5fde38b71e652bda164aa4`. Whole offline Go build/vet/test with retained C04/C05/C07 pins
and completed C07 report, standard whole affected races, generator, catalog
quality, CLI and diff checks. [Gate receipts](../../dist/c09-goal/c09-4-gates.log).

C09.4 attempt c09-4-a1 recorded a whole ordinary test FAIL after 99.268s.
The one failing key is `apitools.TestC07HubSpotEntryListing`. Its signature
is the old listing helper treating the now-explicit HubSpot `Test At Debug`
version-label conflict as unclassified. The key used 1 of 3 attempts.
C09.4 specifically fixes FND-06, which requires separating this prerelease
label from Stable. The correction will keep all 377 original source rows,
the pinned digest, the counts of the unaffected classes, and the exact
accepted C07 TSV comparison. The current classification adds one explicit
conflict and has 306 stable candidates. Only the test's historical
projection preserves the former 307-candidate count.

C09.4 recovery attempt c09-4-a2 started at `47621bd085e02b2be3454a6981b16d34b646f935` with diff SHA-256
`cedfeebcf6c34bfe82c5f4beca4468cb829e71cbbbafaab0554384a17131bae8`. Focused test now passes with the original pinned source and frozen TSV.
Rerun all ordinary required gates and standard whole races.
[Corrective gate receipts](../../dist/c09-goal/c09-4-final-gates.log).

C09.4 passed separate Astra xhigh task review on exact diff SHA-256
`cedfeebcf6c34bfe82c5f4beca4468cb829e71cbbbafaab0554384a17131bae8`.
Whole ordinary test recovery passed in 73.185s, clearing the used gate
key. Whole root, CLI and storage races passed in 200.584s. Build, vet,
catalog generation, catalog quality, CLI help, and diff passed.
The adopted Go test state directory uses a temporary native filesystem
when available, otherwise a disposable workspace child, and accepts an
existing root selected through `APITOOLS_TEST_STATE_DIR`.

The built CLI read the C07 copy at the original run ID. Its raw run
object equals the frozen JSON for 334 providers and all 529 outcomes.
The scratch schedule was enabled, inspected and disabled using the
completed first run. The copied state SHA-256 remained unchanged,
the unrelated crontab bytes were restored, and the Git export contained
17 accepted preferred versions. [Built CLI proof](../../dist/c09-goal/c09-4-prior-run.log).
The pinned suite verified Dropbox, HubSpot, Kubernetes, Stripe,
Google Sheets and Slack against the accepted C07 qualification evidence.

The local two-month fixture ran both full scopes, returned identical
accepted version exports, and kept one content-addressed object for the
identical direct artifact. The explicit loopback search fixture sent
335 inputs on its first search and two on the second after changing one
public profile. It reopened the same separate search database. CLI
tests verified a pending saved run, bounded retry from accepted C10
behavior, distinct provider retry after releasing the monthly lock, and
read-only repeat export. FND-03 and FND-06 are fixed.
[Focused checks](../../dist/c09-goal/c09-4-focused.log).
[Historical comparison and portable-state proof](../../dist/c09-goal/c09-4-followup-focused.log).

OpenUdon passed against the exact product code. Udon's ordinary pinned
suite passed. Optional Udon composition was a setup FAIL on the
pre-existing uncached module `cloud.google.com/go/compute@v1.49.1`.
Only the pinned APItools test changed after those consumer runs.
[Consumer scope](../../dist/c09-goal/c09-4-consumers.txt).
The integrated review uses the C06 starting Git revision and exact
initial tracked diff. It recognizes the four unavailable planning-only
snapshots and verifies the remaining original inputs.
[Baseline reconstruction](../../dist/c09-goal/c09-integrated-baseline.json).

Integrated whole round 1 started at exact candidate
`cf8591db608baaa5ad963b4a54e27459b2e5311a`.
[C09 delta](../../dist/c09-goal/c09-integrated-round1-delta.diff) SHA-256
`6ba99cddd11b56d0ad731060613c818996e45663cb6722b9b43c26f5f6f528d5`.
[Canonical integrated Git diff](../../dist/c09-goal/c09-integrated-round1-full.diff)
SHA-256 `b52fc2c1d58c996b36296e1ca49e493094966c12b2698449700b3f2ab5699393`.
[Review packet](../../dist/c09-goal/c09-integrated-round1-review.txt) SHA-256
`84992ccba7378a17509df8d7935bd2fc76303decfb69f9853c9b4583dcc8d198`.
Reuse the whole test gates from C09.4 at this source, with only ledger and
findings bookkeeping changed after qualification.

Round 1 blocked on F1. [Report](../../dist/c09-goal/c09-whole-round1-report.md)
and [Opus 5.5 high dispatch](../../dist/c09-goal/c09-whole-round1-dispatch.json).
The [baseline regression](../../dist/c09-goal/c09-whole-round1-f1-red.log)
confirmed zero publisher requests after mismatch and no previously committed
checkpoints in the returned report. A diagnostic alone would leave a changed
catalog unable to service future periods. The repair gives the operator an
explicit run-ID selection, returns the original outcomes on refusal, moves
the assignment transactionally, and leaves the checkpoint and version rows
available to all existing readers. Manual writers and scoped retry now share
the lock, while inspection and export continue to use read-only SQLite paths.
The repair also covers timed retry evidence, hot journal recovery, stale
schedule proof, and the case where a deferred source loses a cache object.

Integrated corrective gate `c09-whole-r1-a1` started from
`cf8591db608baaa5ad963b4a54e27459b2e5311a`, exact implementation and ordinary
documentation diff SHA-256
`1c017e71007a2dad6ccb61af3a1d7d8ae5e01502f54662b7c6fa3bcfae6bc2ca`.
Re-run offline focused regressions, build, vet, whole pinned tests, root/CLI/
storage races, generator, quality, CLI help, skill structure and diff checks;
requalify the affected OpenUdon composition.
[Required gates](../../dist/c09-goal/c09-whole-round1-corrective-gates.log).

Round 1 observed dispatch 21:42:19Z and handback 21:50:10Z on 2026-11-25.
F2 through F6 were confirmed in the reviewed source and addressed in this
repair. F2 uses the same state lock for scheduled, manual and scoped writers.
F3 now drives real timed retries against a fixture transport through the
terminal attempt. The one-day wall limit is verified by its explicit context
deadline and cancellation tests. F4 exercises a real hot journal in disposable
state and guides read-only callers to the guarded writer. F5 shows a stale
installed provider-scope proof alongside schedule presence. F6 reports a
missing historical digest and completes the later provider without relying on
the lost bytes. Manual completed and superseded reports remain readable with
no writer; a superseded run is never restarted by the current updater.

Round 1 lower dispositions. L1, L3 and L4 are fixed with explicit blank-selector
rejection, a manual-resume hint for unfinished targeted runs, and a bounded
cache measurement after cancellation. L2 keeps the documented operator-chosen
separate scratch file and explicit config/state safeguards. L5 would broaden
version-label inference beyond pinned publisher evidence. L6 bounds carryover
selection to due periods; saved reports continue to order by recorded month.
L7 keeps the supported permanent binary location and permissive scratch
preflight. L8 uses disposable native directories as allowed by the test policy.
L9 keeps the structured cron report. L10 leaves C08's independent search-index
refusal intact. L11 is a safe refusal on a hand-edited managed block. L12 uses
the existing race budget. L13 keeps the concise text refusal and full JSON.

Gate `c09-whole-r1-a1` passed the focused tests in 32.879s, build, vet,
whole pinned tests in 91.053s and generator freshness. The catalog quality
check was invoked with the erroneous CLI subcommand `catalog quality` and
exited with usage status 2 after 1.127s. Required-gate key C09 catalog check,
working directory APItools, built-in catalog scope, runner setup CLI selector.
Signature `unknown catalog command "quality"` consumed one attempt of three,
with two remaining and one distinct signature. The CLI help in the receipt
names `catalog check` as the correct command. The first run stopped before
launching the remaining checks.

Gate `c09-whole-r1-a2` starts the exact catalog check, CLI help, skill
validation, diff and whole root/CLI/storage race checks at the same base and
candidate SHA-256 as `c09-whole-r1-a1`. Product inputs are identical, so the
passing required checks retain their evidence.
[Recovery gates](../../dist/c09-goal/c09-whole-round1-corrective-gates-resumed.log).

The corrected catalog check passed in 0.211s with no findings and cleared
the gate's key and signature counters. The CLI help, structural skill and
diff checks also passed. Exact-candidate OpenUdon and optional exact-candidate
Udon composition passed in 37.514s and 35.984s respectively. The Udon copy
required `-mod=mod` for its isolated replacement and completed offline.
[Consumer input](../../dist/c09-goal/c09-whole-round1-consumer-inputs.txt),
[OpenUdon](../../dist/c09-goal/c09-whole-round1-openudon.log), and
[Udon](../../dist/c09-goal/c09-whole-round1-udon-composed-with-mod.log).
The corrected CLI also re-read the exact C07 record and export from the
disposable database. It checked the full-run proof in scratch status and
restored the unrelated crontab bytes with byte-identical SQLite and export.
[Saved-state proof](../../dist/c09-goal/c09-whole-round1-prior-run.log).
Whole root, CLI and storage races passed in 181.596s. All required
corrective gates are clean. No review counter was charged for the harness
correction.
