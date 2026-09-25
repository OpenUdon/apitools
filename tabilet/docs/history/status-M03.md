# Retired milestone M03 - Catalog Entries And Spec References

**Milestone.** M03
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M03.md
**Source specification.** tabilet/memory-bank/milestone.md#m03---catalog-entries-and-spec-references
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M03 - Catalog Entries And Spec References

**Goal.** Convert reviewed candidate services into durable built-in catalog
entries and spec references.

**Scope.**

- Add a focused `catalog` package without refactoring the existing root package.
- Define provider, alias, category, source hint, and spec reference types.
- Track official OpenAPI and official non-OpenAPI machine-readable spec
  references separately.
- Preserve source authority, URL, version/date/hash when known, license notes,
  and update notes.
- Keep upstream/vendor OpenAPI documents separate from local metadata.

**Acceptance.** Downstream callers can list built-in providers and inspect known
spec references with provenance, while unverified candidates remain visibly
distinct from durable catalog entries.
````

## Status record

````markdown
# Status M03 - Catalog Entries And Spec References

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Convert reviewed M2 candidates into durable provider catalog entries and spec
references under `apitools/catalog`.

The catalog should answer these questions without live provider calls:

- Which providers are built into the core catalog?
- Which providers have official OpenAPI documents?
- Which providers have official machine-readable specs that are not OpenAPI?
- Which providers need user-supplied OpenAPI documents?
- What source authority, URL, version/date/hash, license, and update notes are
  known for each spec reference?

## Proposed Package Shape

Prefer a focused subpackage over a root-package refactor:

```text
catalog/
  catalog.go
  provider.go
  load.go
  testdata/
```

Keep existing root-package APIs stable. Add root re-exports only after real
downstream usage proves they are needed.

## Data Model Direction

Catalog entries should cover at least:

- stable provider ID, display name, aliases, and source hints;
- service category and workflow relevance notes;
- official OpenAPI availability;
- official non-OpenAPI machine spec availability, such as Google Discovery;
- spec references with URL, kind, source authority, version/date/hash when
  known, and license/source notes;
- notes for provider-specific quirks that OpenAPI alone cannot model.

## Tasks

| Item | State | Notes |
|---|---|---|
| Catalog package skeleton designed | `[+]` | Added `catalog.Catalog` around the existing focused `catalog` package without root package refactors. |
| Provider entry types implemented | `[+]` | Added durable provider entries with IDs, display names, aliases, categories, source hints, and review state. |
| Spec reference types implemented | `[+]` | Added OpenAPI, OpenAPI index, Google Discovery, and human-docs reference kinds with provenance fields. |
| Deterministic loading implemented | `[+]` | Built-in providers load as copy-safe sorted slices and validate duplicate IDs/aliases. |
| Seed catalog entries added | `[+]` | Promoted the nine reviewed M2 candidates into durable provider entries. |
| Catalog tests added | `[+]` | Covers loading, validation, duplicate detection, deterministic ordering, copy behavior, and spec classifications. |
| Verification completed | `[+]` | Ran `go test ./...`, `go vet ./...`, `git diff --check`, CLI help smoke checks, and sibling `../openudon`/`../udon` tests. |

## Boundary Checks

- Do not execute API operations.
- Do not fetch credentials, tokens, or secrets.
- Do not sign requests.
- Do not choose accounts, tenants, workspaces, or runtime environments.
- Keep upstream/vendor OpenAPI documents separate from local metadata.
- Do not vendor large third-party specs until source, license, size, and update
  strategy are explicit.
````
