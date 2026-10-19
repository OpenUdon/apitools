# Retired milestone S08 - Exported guarded document fetch

**Milestone.** S08
**Outcome.** completed
**Retired.** 2026-10-19
**Source status.** tabilet/memory-bank/status-S08.md
**Source specification.** tabilet/memory-bank/milestone.md#s08--exported-guarded-document-fetch
**Evidence.** 6f3d18545d9c7a37cc571388c9a0938883c07532
**Worktree.** includes uncommitted changes
**Review.** passed
**Review iterations.** 1
**Verification.** GOWORK=off go vet ./... passed; GOWORK=off go test ./... passed (all packages); go test -race . -run FetchDocument passed; git diff --check passed. The refusal tests were proved by swapping in an unguarded http.DefaultClient fetch, which failed them. Not published: no push, tag or consumer-pin change.
**Consolidated into.** [lessons](../../memory-bank/lessons.md), [architecture](../../memory-bank/architecture.md), [README](../../../README.md), [CHANGELOG](../../../CHANGELOG.md). No product or stack change.

## Milestone specification

````markdown
### S08 — Exported guarded document fetch

**Stage/owner.** APItools (lane S: import safety). Requested by Ramen:M101; Ramen's
2026-10-19 artifacts review found that its GUI downloader used `http.DefaultClient`
and so bypassed the host guard that `Import` applies.

**Scope.** Add `(*Client).FetchDocument(ctx, rawURL) (FetchedDocument, error)`, a thin
exported wrapper over `downloadBoundedWithAccept`. `FetchedDocument` holds `Content`
([]byte), `FinalURL` (*url.URL) and `ContentType` (string). It applies the client's
existing host and port checks at dial time and on every redirect, its timeout, its byte
limit (wire and decoded gzip) and its `AllowUnsafeHosts`/`AllowedPorts` settings. A
non-2xx response returns the existing `HTTPStatusError`. It validates no document
content, performs no caching and reads no credentials. Existing behavior of `Import`,
`ValidateURL` and `Search` is unchanged. README, CHANGELOG and `architecture.md` record
the method.

**Acceptance.** `FetchDocument` refuses `localhost`, loopback, private, link-local and
scoped hosts; uses the shared redirect-safe client, so a redirect to a private host is
refused (proved by the existing `TestRedirectToPrivateHostRejected`, since a public
origin is not reachable offline and `AllowUnsafeHosts` also lifts the redirect check); enforces the byte limit
and the timeout; returns the final URL and Content-Type of an allowed fetch; returns
`HTTPStatusError` for a non-2xx response; and reaches a loopback server only when
`AllowUnsafeHosts` is set. `GOWORK=off go vet ./...`, `GOWORK=off go test ./...`,
`go test -race` on the new tests and `git diff --check` pass offline.

**Depends on.** Nothing open.

**Downstream impacts.** Ramen:M101 (conditional on the user publishing this repository
and Ramen re-pinning). No sibling pin changes here.

**Write set.** `download.go` (or a new `fetch.go`), its tests, `README.md`,
`CHANGELOG.md`, `tabilet/memory-bank/` records. No other source file.

**Contracts read.** The existing `downloadBoundedWithAccept` behavior and its tests.

**Parallel-safe.** no. Execute sequentially.

**Outcome.** S08 completed and is retired (1 review iteration). Unpublished: Ramen adopts it after the user publishes this repository.
````

## Status record

````markdown
# S08 — Exported guarded document fetch

**Stage:** none. **Owner:** APItools.
**State:** Complete. Review 1/10, passed.
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
| S08.2: Document the method | `[+]` | README (library use), CHANGELOG entry, and the `architecture.md` description of the guarded downloader. Done: README feature bullet, CHANGELOG "Unreleased" entry, and the architecture data-flow note. |
| S08.3: Review iteration 1 test | `[+]` | `TestFetchDocumentRefusesAnUnlistedPort`: the README and spec say `FetchDocument` honours the port policy, which the earlier tests did not exercise. A public literal address on an unlisted port is refused with "refusing URL port", with no lookup or connection. |

## Review

Base for the diff: `2341e9f` (planning commit); code and docs commits `6aa071a`, `96c9a44`. Verification before review: `GOWORK=off go vet ./...`, `GOWORK=off go test ./...`, `go test -race . -run FetchDocument` and `git diff --check` passed.

Iteration 1 findings:

| # | Finding | Severity | Disposition |
| --- | --- | --- | --- |
| 1 | The README and the specification say `FetchDocument` honours the port policy (`AllowedPorts`), but no test exercised it; only the host guard, limits and timeout were covered. | Lower | Fixed by S08.3 (test only). |

No P1 or P2 finding: the method adds no code path of its own; it returns what `downloadBoundedWithAccept` returns, which `Import` already uses and its existing tests cover (host, redirect, size, gzip, status, timeout). The accepted gap is the end-to-end redirect-to-private case, which cannot be shown offline and rides the shared `redirectSafeClient` (`TestRedirectToPrivateHostRejected`).

Iteration 1 passed.

## Verification Plan

```bash
GOWORK=off go vet ./...
GOWORK=off go test ./...
go test -race . -run 'FetchDocument'
git diff --check
```
````
