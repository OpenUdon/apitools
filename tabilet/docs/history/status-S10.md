# Retired milestone S10 - Keep top-level response fields when a response summary is truncated

**Milestone.** S10
**Outcome.** completed
**Retired.** 2026-10-22
**Source status.** tabilet/memory-bank/status-S10.md
**Source specification.** tabilet/memory-bank/milestone.md#s10--keep-top-level-response-fields-when-a-response-summary-is-truncated
**Evidence.** db9e9aad9717c103cd309b71f307e9a44d779e2e
**Worktree.** includes uncommitted changes
**Review.** passed
**Review iterations.** 1
**Verification.** Tests written first and red (Storage and synthetic OpenAPI dropped `name`/`id`); then GOWORK=off go vet ./... passed; GOWORK=off go test ./... passed (all packages); gofmt -l clean; git diff --check passed; read-only OpenUdon go test ./... from a scratch copy with a replace directive passed unchanged.
**Consolidated into.** [architecture.md](../../memory-bank/architecture.md) (response-field summary note), CHANGELOG.md; no new lesson.

## Milestone specification

`````markdown
### S10 — Keep top-level response fields when a response summary is truncated

**Stage/owner.** APItools (lane S: inventories and prompt-safe summaries). Requested by
Ramen:M128, which adopts the published revision.

**Scope.** For OpenAPI and Google Discovery response bodies, collect the response fields within the
existing depth and work bounds, and when more than 60 remain, keep every top-level field first (by
path), then nested fields by path, up to 60; then sort the kept fields by path. This replaces the
cut-then-append of `name`/`id` for these response paths, so the later prompt-safety trim removes
nothing. When more than 60 top-level fields exist, the first 60 by path are kept. The existing
`schema.field_summary_truncated` issue is still reported. Request-field summaries, the native
schema walkers (AsyncAPI, OpenRPC and the rest, which refuse over budget) and the public types are
unchanged. CHANGELOG and `architecture.md` record the behavior.

**Rows.**
- S10.1: the selection for OpenAPI and Google Discovery response fields, with tests written first:
  Google Storage `storage.buckets.get` keeps `name`, `id`, `location` and `storageClass`, has at most
  60 fields, is sorted by path and still reports the truncation issue (the Storage Discovery document
  copied into `testdata/`); a synthetic OpenAPI response with more than 60 fields behaves the same;
  a response under the budget is unchanged; request-field summaries are unchanged; more than 60
  top-level fields keep the first 60 by path.
- S10.2: CHANGELOG "Unreleased" entry and the `architecture.md` note.

**Acceptance.** The tests above pass; `GOWORK=off go vet ./...`, `GOWORK=off go test ./...` and
`git diff --check` pass offline. Because the kept fields of truncated responses change, a read-only
consumer check runs OpenUdon's tests against the S10 candidate from a scratch copy with a `replace`
directive (no OpenUdon edit); a difference is reported for OpenUdon's own planning, not fixed here.

**Depends on.** Nothing open. Ramen:M128 depends on S10's publication.

```yaml
AUTHORIZATION_REQUIREMENTS:
  git.commit:
    via: goal-policy
    executor: assigned-agent
  cli.local:
    via: goal-policy
    executor: assigned-agent
    scope: declared-implementation-and-verification, including a read-only OpenUdon test run from a scratch copy
```
`````

## Status record

`````markdown
# S10 — Keep top-level response fields when a response summary is truncated

**Stage:** none. **Owner:** APItools.
**State:** Complete. Review 1/10, passed (base `7f1a22869252dc4fa9d379a757651718ad3b7195`, reviewed head `db9e9aad9717c103cd309b71f307e9a44d779e2e`; no P1/P2 finding: request-field path equivalent to the former walk, `looksLikeCredentialName` filtering and nullable propagation kept, Google candidates use the native walker and are unaffected, OpenUdon consumer check unchanged).
**Source baseline:** `8389a15f2bf811072c09fff1e9a10a2946315b57` (clean at planning).
**Specification:** [milestone.md](milestone.md#s10--keep-top-level-response-fields-when-a-response-summary-is-truncated).
**Requested by:** Ramen:M128 (conditional on S10's publication).
**Authority:** task commits, local verification and a read-only OpenUdon consumer check only. No push,
tag, publication or consumer-pin change; the user publishes this repository separately.

## Status rules

See [`milestone.md`](milestone.md) for the status ID pattern and markers. Write markers
with backticks: `` `[ ]` ``, `` `[~]` ``, `` `[+]` ``, `` `[!]` ``, `` `[X]` ``, `` `[-]` ``.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| S10.1: Keep top-level fields when truncating OpenAPI and Google Discovery response summaries | `[X]` | `responseFieldSummaries` (`inventory_request.go`) keeps every top-level field first (first 60 by path), then nested fields by path, sorted by path; replaces cut-then-append of `name`/`id`. Google inventory operations now also report `schema.field_summary_truncated` (candidates use the native walker and are unaffected). Tests in `inventory_response_fields_test.go`, fixture `testdata/inventory/google-storage-discovery.json`. `GOWORK=off go vet`, `go test ./...`, `git diff --check` pass; OpenUdon `go test ./...` from a scratch copy with `replace` passes unchanged. |
| S10.2: Document the change | `[X]` | CHANGELOG "Unreleased" entry and `architecture.md` note (Data Flow, prompt budget section). |
`````
