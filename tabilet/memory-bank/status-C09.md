# C09 — Monthly cron and agent maintenance skill

**Owner.** APItools. **Dependency.** Accepted C08. **Parallel-safe: no.**
C09 owns integrated closure for C06–C09 and the [C10 updater repairs](../docs/history/status-C10.md) against the baseline in [C06.1](../docs/history/status-C06.md#status-record). See the [milestone specification](milestone.md#c09--monthly-cron-and-agent-maintenance-skill).

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| C09.1 — Check readiness and implement the monthly job | `[ ]` | Read the same cache directory configuration as a manual update. Preflight binary, database, configured directory, permissions and cron environment. Use absolute paths, local time, one run lock, bounded retries and per-provider checkpoints. Exercise an ordinary cache directory, overlapping jobs, restart and loss of an explicitly configured mount in isolated tests. |
| C09.2 — Expose schedule management and reports | `[ ]` | Default the optional user schedule to `17 3 1 * *`. Provide idempotent enable, disable and status, preserving unrelated crontab entries. Produce readable and JSON reports for accepted changes, pending reviews, size deferrals, failures, disk use and embedding usage. Test with a scratch crontab. This row fixes FND-02. |
| C09.3 — Deliver the agent maintenance skill | `[ ]` | Ship the repository-owned `apitools-catalog-maintenance` skill. Check its instructions against actual commands for preflight, dry-run, schedule management, report inspection, per-provider retry and generating the Git export. Document the operator-initiated first run and how to inspect it before scheduling. |
| C09.4 — Verify repeat run and integrated closure | `[ ]` | Reconcile [C07 first-run evidence](../docs/history/status-C07.md#status-record) for 334 providers and all six special cases. Record the repeat run, unchanged artifact reuse and selective re-embedding. Verify repeat outcomes, bounded retry and inspection of previous runs against accepted C10 evidence. Scope any live publisher or metered embedding check separately. Exercise the operator/agent workflow, review the cumulative diff, and run the existing class O whole-unit gates and affected race and consumer checks. This row fixes FND-03 and FND-06. |

## Review

Review gate iterations: 0

Whole review pending, default limit 10. Reuse accepted evidence for contributor
milestones and identify the exact final state. Use the [test policy](test-policy.md)
and [milestone review procedure](milestone.md#milestone-review-procedure).
