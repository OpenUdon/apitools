# S09 — Ranked search over the built-in catalog

**Stage:** none. **Owner:** APItools.
**State:** In progress (S09.2). Review 0/10.
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
| S09.2: Document the search | `[ ]` | README (library use), CHANGELOG "Unreleased" entry, and `docs/catalog-discovery.md` or the `architecture.md` catalog description. Depends on S09.1. |
