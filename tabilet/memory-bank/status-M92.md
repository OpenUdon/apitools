# M92 — Published HCL-free UWS adoption

**Stage/owner:** STG-13; APItools.
**Planning approval:** Owner-approved full rescope on 2026-10-25 (approval: “I aapaprove”), following the earlier whole-stage approval. Planning only; no execution, commit or grant activation.
**State:** Withdrawn; every untouched row is `[X]`; no work delivered. Review 0/10 remains unchanged.
**Specification:** [milestone.md](milestone.md#stg13-m92).
**Original dependencies (withdrawn; no active edges):** [UWS:C13](../../../uws/tabilet/memory-bank/status-C13.md), [APItools:M91](status-M91.md)
**Source baseline:** `fe54b3ba939427d86ebbdd4a47b0dccb519cc26c`; existing dirty planning/history was inspected and preserved. Reconcile live source and exact accepted producers before implementation.
**Ownership:** Withdrawn; no executable task owner. Original serial ownership and task boundaries are historical.
**Authorization:** None from this withdrawn record. Original requirements are retained only as historical text in the specification; no old or proposed publication grant transfers.

**Lineage:** Fresh successor to withdrawn APItools:M86. APItools:M90 catalog-root ingestion stays reserved and unscheduled.

**Withdrawal disposition:** Withdrawn by the owner-approved 2026-10-25 rescope: retain supported HCL in Ramen/Udon and the supporting UWS codecs. This obligation is undelivered, not accepted. HCL removal and legacy-wrapper retirement move to separate unnumbered Candidate Directions; no replacement is scheduled and no active milestone depends on this row.

**Review provenance.** Kinet Stage 13 plan review (`kinet-stage13.md`), 2026-10-25; source commit baseline not supplied. H2/H3 (High/local P1): current Ramen/Udon readers and writers depend on UWS HCL. The owner selected preservation. Original task wording and review counters are retained below.

## Tasks

| Item | State | Notes |
|---|---|---|
| M92.1 — Adopt the exact published UWS version | `[X]` | Pin independently verified UWS:C13 through ordinary module resolution and reconcile all UWS adapter/API compile changes. **Withdrawal:** Withdrawn by the owner-approved 2026-10-25 rescope: retain supported HCL in Ramen/Udon and the supporting UWS codecs. This obligation is undelivered, not accepted. HCL removal and legacy-wrapper retirement move to separate unnumbered Candidate Directions; no replacement is scheduled and no active milestone depends on this row. |
| M92.2 — Retire remaining HCL and wrapper paths | `[X]` | Remove obsolete UWS HCL/wrapper support from metadata/authoring adapters; preserve JSON/YAML and meaningful M91 positive/refusal coverage. **Withdrawal:** Withdrawn by the owner-approved 2026-10-25 rescope: retain supported HCL in Ramen/Udon and the supporting UWS codecs. This obligation is undelivered, not accepted. HCL removal and legacy-wrapper retirement move to separate unnumbered Candidate Directions; no replacement is scheduled and no active milestone depends on this row. |
| M92.3 — Qualify and publish the consumer-producer handoff | `[X]` | Run required standalone and helper checks, prove HCL-free selected/compiled closure, review and prepare exact publication for Udon/OpenUdon/Ramen/Kinet adoption. **Withdrawal:** Withdrawn by the owner-approved 2026-10-25 rescope: retain supported HCL in Ramen/Udon and the supporting UWS codecs. This obligation is undelivered, not accepted. HCL removal and legacy-wrapper retirement move to separate unnumbered Candidate Directions; no replacement is scheduled and no active milestone depends on this row. |

## Original acceptance and verification — withdrawn, undelivered

APItools consumes the published HCL-free UWS core and keeps metadata/authoring helpers compatible without UWS HCL input or wrappers.

- Exact published C13 is adopted ordinarily; UWS core/APItools selected and compiled closures have no HCL dependency.
- Source metadata and M91 fixtures retain their identities and behavior; no runtime or credential ownership is added.

Original planned verification (not scheduled or performed): Standalone tests/vet, operation/candidate/mapping fixtures and selected/compiled dependency proof; focused HCL/wrapper refusal and JSON/YAML preservation cases. All evidence is pending. Preserve existing package limits. Use disposable fixtures/fake providers; identify executed, mocked, compile-only and deferred proof separately.

## Review

Review gate iterations: 0

Task reviews: not started; withdrawal is not acceptance. Whole-milestone review: 0/10, not started. Follow the package's existing persisted review and closure procedure; planning, terminal rows and successful builds do not establish acceptance.
