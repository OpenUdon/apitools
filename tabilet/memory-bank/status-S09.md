# S09 — Ranked search over the built-in catalog

**Stage:** none. **Owner:** APItools.
**State:** Pending. No row has started. Review 0/10.
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
| S09.1: Add `Catalog.SearchProviders` and `SearchBuiltInProviders` | `[ ]` | Tests first in `catalog`: `github` gives GitHub first and URL-only matches at rank 3; `stripe` and `openweathermap` give their own provider at rank 1; `cloudflare dns` needs both words; case and punctuation ignored; limit; deterministic; overlay-only word gives rank 3 (small test catalog); empty query gives none; returned providers are copies. Then the implementation (types, normalization, rank, order). |
| S09.2: Document the search | `[ ]` | README (library use), CHANGELOG "Unreleased" entry, and `docs/catalog-discovery.md` or the `architecture.md` catalog description. Depends on S09.1. |
