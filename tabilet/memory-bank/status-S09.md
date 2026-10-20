# S09 — Ranked search over the built-in catalog

**Stage:** none. **Owner:** APItools.
**State:** Rows complete; review gate passed at iteration 1. Review 1/10 (base `87f14076460626cda0a598bb958a90d348d0f5f4`, reviewed head `a94d440`).
**Source baseline:** `87f14076460626cda0a598bb958a90d348d0f5f4` (clean at planning).
**Specification:** [milestone.md](milestone.md#s09--ranked-search-over-the-built-in-catalog).
**Requested by:** Ramen:M110 (conditional on S09's publication).
**Authority:** task commits and local verification only. No push, tag, publication or
consumer-pin change; the user publishes this repository separately.

## Status rules

See [`milestone.md`](milestone.md) for the status ID pattern and markers. Write markers
with backticks: `` `[ ]` ``, `` `[~]` ``, `` `[+]` ``, `` `[!]` ``, `` `[X]` ``, `` `[-]` ``.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| S09.1: Add `Catalog.SearchProviders` and `SearchBuiltInProviders` | `[+]` | New `catalog/search.go` and `catalog/search_test.go`. Tests were written first and failed to compile (undefined symbols) before the implementation. Covered: `github`, `stripe`, `openweathermap`, `cloudflare` and `slack` return their own provider first at rank 1; `github` lists 93 providers, one at rank 1 and 92 at rank 3 (URL, note or overlay text only), in rank order; `cloudflare dns` needs both words; case and surrounding punctuation are ignored; display name, alias and normalized key (`agile crm`) match at rank 1; limit; empty or punctuation-only query returns none; two searches are identical; a small test catalog proves rank 1/2/3 (name, category or workflow relevance, quirk, spec reference URL, overlay note) and that a word's lowest-priority field decides the rank; an exact key outranks providers that merely contain the word; returned providers (aliases, spec references, quirks) are copies. Mutation proof: ranking other text as category fails three tests, and dropping the exact-key rule fails one. Verified: `GOWORK=off go vet ./...`, `GOWORK=off go test ./...`, gofmt and `git diff --check` pass. Note: the 297 in the planning evidence counts Ramen's per-spec-reference rows, not providers; this search finds 93 providers for `github`. |
| S09.2: Document the search | `[+]` | README Go usage (an example and a paragraph on query handling, ranks, ordering, limit, and that it is not `apitools catalog search`), CHANGELOG "Unreleased" entry, and the `catalog/search.go` row in `architecture.md`. `docs/catalog-discovery.md` was left unchanged: it describes operation discovery, and the specification allowed either it or `architecture.md`. Verified: `GOWORK=off go vet ./...`, `GOWORK=off go test ./...` and `git diff --check` pass. |

## Review

Iteration 1 (base `87f14076460626cda0a598bb958a90d348d0f5f4`, head `a94d440`; full diff of seven files, checked against the specification's acceptance lines). Result: no P1 or P2 finding; the gate passed.

- Acceptance fidelity: every acceptance line is covered by a test (`github`, `stripe`, `openweathermap` first at rank 1; `github` URL-only matches at rank 3; `cloudflare dns` needs both words; case and surrounding punctuation; limit; determinism; overlay-only word at rank 3; empty query; copies). No requirement is missing, and the matching rule is the same as Ramen's current matcher (same fields, same token trimming), so no match it finds today is lost. Nothing was added that was not asked: no CLI, no whole-word matching, no existing API changed (the diff adds two source files and edits docs and memory only).
- Lower finding L1 (fixed): a branch keeping an exact-key provider whose words did not match was unreachable, because an exact normalized key implies every word is found; its test comment claimed otherwise. Removed the branch and corrected the comment.
- Lower finding L2 (fixed): `SearchBuiltInProviders` cloned the whole catalog before `SearchProviders` cloned the matches. It now reads the built-in data directly, since the search only reads and returns copies.
- After the fixes: `GOWORK=off go vet ./...`, `GOWORK=off go test ./...`, `go test -race ./catalog -run TestSearch`, gofmt and `git diff --check` pass.
- Boundary checks: no network, file or cache access; no change to an existing exported API; no consumer pin, push or publication.
