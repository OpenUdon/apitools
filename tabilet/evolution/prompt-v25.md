# Prompt V25 - Catalog discovery for step contracts

Plan APItools' contribution to Kinet's G1/S2d: a reusable library API that
retrieves ranked providers and artifact references across the local catalog
for a step's purpose, typed inputs/outputs, and effect. Build on the completed
M77-M79 operation metadata without changing their exported contracts or history.
Expose authority, license/redistribution evidence, auth needs, effects, rank
evidence, coverage, and uncertainty. Keep missing metadata explicit and scope
negative results to the evidence actually searched. Offline is the default;
remote lookup is a separately configured bounded tier.

Before building discovery, settle its design decisions and fix the
pre-existing limits that would stop it at catalog scale. Build a digest-bound
operation index at refresh time and have discovery read only that index. Let
only the index builder read registered artifacts above 20 MiB, up to the
existing 128 MiB artifact bound. Keep catalog roots caller-supplied, accept
optional provider constraints by exact key, and add artifact-scoped export so
a discovered reference can be provisioned explicitly. Separate a match from an
ambiguous or weak result, and let only a scoped no-match route a consumer to a
browser step.

Source: "APItools stage 5 draft — M80" in the Kinet stage 5 draft plan
(revision 2, 2026-09-30; staging material outside Kinet's tracked history),
that plan's deep review (finding 2), and
[Kinet request-resolution G1/S2d](../../../kinet/docs/request-resolution.md).
The user approved the seven-file planning reconciliation on 2026-09-30 against
`8580ff2485a3faff139b17bdc0b8e78d2af4ce18` with a clean worktree. The same day,
the user approved a second reconciliation from review "APItools M80 — Catalog
discovery API for step contracts" (`apitools-m80-review.md`), which added
prerequisite M81 and its six design decisions.
Implementation, commits, publication, and sibling changes require separate
execution authority. Browser routing, provisioning, workflow approval, and
runtime execution remain downstream responsibilities.
