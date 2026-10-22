# S10 — Keep top-level response fields when a response summary is truncated

**Stage:** none. **Owner:** APItools.
**State:** Pending. Review 0/10.
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
| S10.2: Document the change | `[ ]` | CHANGELOG "Unreleased" and `architecture.md`. Depends on S10.1. |
