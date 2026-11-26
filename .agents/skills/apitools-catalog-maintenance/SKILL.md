---
name: apitools-catalog-maintenance
description: Manage the APItools reviewed-provider monthly update. Use for maintenance preflight, provider-scope dry runs, saved run reports, scoped retries, cron schedule changes, accepted-version export, or catalog-site preview and release.
---

# Maintain the reviewed API catalog

Use the installed `apitools` binary and the operator's maintenance JSON file.
Resolve both to absolute paths for scheduling. The same file configures
manual and scheduled updates. It names `cache_dir` and `state_db`. Set
`expected_volume_uuid` when the cache must stay on a designated volume.
Use the [operator guide](../../../README.md#reviewed-catalog-maintenance) for
configuration and source-status details.

## Check readiness and inspect state

- `apitools catalog maintenance preflight --config FILE` checks the
  executable, storage, and cron environment. `run --config FILE --dry-run
  --json` lists the reviewed scope without contacting publishers.
- `apitools catalog maintenance status --config FILE --json` reads the
  managed schedule and latest scheduled run ID and checks the installed
  full-run proof. `status` also works without the config when only the
  schedule is needed.
- `apitools catalog maintenance report --config FILE --run RUN_ID --json`
  inspects saved outcomes, including an unfinished run's committed
  checkpoints. Omit `--run` for the latest scheduled run. Omit `--json`
  for text. Publisher head outcomes and accepted pins can coexist.

## Run, recover, and schedule

The operator starts the first full run with
`apitools catalog update --config FILE --all --json`. Capture its run ID
and inspect the completed report before enabling recurrence. A previously
completed full run in the same state database can supply that proof when it
covers the current reviewed provider list.

For a manual update with pending retries, inspect the saved `next_attempt_at`
and resume its original run with
`apitools catalog update --config FILE --resume RUN_ID --json`. For an
interrupted monthly job, resume through
`apitools catalog maintenance run --config FILE --json` with the original
configuration and catalog. The scheduled runner handles eligible retries
within its time budget and continues the due month under the same lock.
If the catalog or configuration has changed, inspect the returned run ID.
With the operator's authority to supersede that run, use
`apitools catalog maintenance supersede --config FILE --run RUN_ID`.
This releases the period assignment and keeps the original
checkpoints and accepted versions. Keep the run ID for later inspection.
Supersession is permanent. The next authorized `maintenance run` services the
due month. A manual or scheduled writer refuses while the lock is held.
If a scheduled invocation cannot start, run it after the writer finishes to
service the current month.

If inspection reports an interrupted SQLite journal, resume the original manual
run or use `maintenance run` for a scheduled run. Use the configured storage
and publisher authority. Then repeat the read-only inspection. When status
reports that the installed proof is out of date, complete and inspect a
new full run and enable with that run ID.

To start a new check of just one service, use
`apitools catalog maintenance retry --config FILE --provider PROVIDER_ID --json`.

`apitools catalog maintenance enable --config FILE --after-run RUN_ID`
installs the managed entry for 03:17 local time on the first day of the
month. `apitools catalog maintenance disable` removes the managed entry
even if the maintenance configuration or volume is unavailable. Check
status after a schedule change. For a rehearsal, pass `--crontab-file` with
a separate scratch text file to enable, disable, or status. Use the scratch
file's exact path with no leading or trailing whitespace in its components.

Use the operator's current authority for publisher access and changes to
the real user crontab. When authority covers only an inspection or rehearsal,
continue with the matching read-only or scratch commands. Scheduled updates
report embedding usage as unrequested. Hosted search embeddings have their
own endpoint and budget.

## Prepare the Git export

Run `apitools catalog update --config FILE --export` to emit deterministic
accepted-version metadata as JSON. Save it to the requested review path and
inspect the diff. Use the user's separate Git instructions for any commit or
publication.

## Preview and release the catalog site

- Follow the [catalog site operator guide](../../../docs/catalog-site.md).
  Build the approved clean commit with the pinned Go toolchain. Verify the
  artifact manifest and preview it on loopback. Compare the two reviewed
  catalog JSON files and review the monthly record before committing it.
- Before dispatch, present the exact source commit and outgoing Git commits
  to the owner. Publication uses the manual `Publish reviewed catalog` Pages
  workflow. Obtain separate approval for each push, domain verification,
  Pages configuration, DNS change, deployment and HTTPS enforcement.
- For a preview rollback, verify the known previous artifact and restart the
  loopback server with that directory. For a live rollback, obtain authority
  to dispatch the last verified published commit, then check its manifest and
  served files. For first-release rollback or retirement, restore or remove
  the routing record before disabling Pages.
