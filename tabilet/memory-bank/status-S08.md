# S08 — Exported guarded document fetch

**Stage:** none. **Owner:** APItools.
**State:** In progress. Review 0/10.
**Source baseline:** `bdc6d9250bdcb59fe20526ab72e97e811d6eec84` (clean at planning).
**Specification:** [milestone.md](milestone.md#s08--exported-guarded-document-fetch).
**Authority:** task commits and local verification only. No push, tag, publication or
consumer-pin change; the user publishes this repository separately.

## Status rules

See [`milestone.md`](milestone.md) for the status ID pattern and markers. Write markers
with backticks: `` `[ ]` ``, `` `[~]` ``, `` `[+]` ``, `` `[!]` ``, `` `[X]` ``, `` `[-]` ``.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| S08.1: Add `Client.FetchDocument` | `[+]` | Tests first, beside the `TestDownloadBounded*` set: refuses localhost, loopback, private, link-local and scoped hosts; refuses a public-to-private redirect; enforces the byte limit and the timeout; returns the final URL and Content-Type; returns `HTTPStatusError` on non-2xx; reaches loopback only with `AllowUnsafeHosts`. Then the exported method and `FetchedDocument` type over `downloadBoundedWithAccept`. Done in `fetch_document.go` / `fetch_document_test.go`. Red first against a stub for the positive cases; the refusal cases were proved by swapping in an unguarded `http.DefaultClient` fetch, which failed them (one test was tightened to require the guard's own refusal message, since an unreachable address also errors). Redirect-to-private refusal rides the shared `redirectSafeClient`, proved by the existing `TestRedirectToPrivateHostRejected`. |
| S08.2: Document the method | `[ ]` | README (library use), CHANGELOG entry, and the `architecture.md` description of the guarded downloader. |

## Verification Plan

```bash
GOWORK=off go vet ./...
GOWORK=off go test ./...
go test -race . -run 'FetchDocument'
git diff --check
```
