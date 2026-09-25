# Retired milestone M70 - Bounded Remote Catalog Discovery

**Milestone.** M70
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M70.md
**Source specification.** tabilet/memory-bank/milestone.md#m70---bounded-remote-catalog-discovery
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M70 - Bounded Remote Catalog Discovery

**Goal.** Expand remote OpenAPI candidate discovery without weakening source
provenance, download safety, or the existing APIs.guru-first behavior.

**Scope.**

- Keep APIs.guru as the primary global source and the legacy public-apis probe
  as the final compatibility fallback.
- Add LAP Registry targeted search as an experimental secondary source. Return
  its reported original `source_url`, not the transformed LAP document, and
  label candidates unvalidated until normal import succeeds.
- Add provider-scoped RFC 9727 discovery when a caller supplies a provider URL
  or hostname. Require `application/linkset+json`, inspect OpenAPI-like
  `service-desc` links, deduplicate URLs, and do not recurse into nested
  catalogs.
- Preserve unsafe-host rejection, redirect and response bounds, context
  cancellation, deterministic result ordering, cache isolation by provider,
  and a visible service-description link bound.
- Expose direct CLI source selection and `--provider-url`; keep tests offline
  with local HTTP servers.
- Do not query Swagger Catalog or Scalar Registry as unauthenticated global
  aggregators, crawl the general web, fetch credentials, or execute operations.

**Acceptance.** Direct and automatic LAP/RFC 9727 lookup paths expose
provenance-labeled unvalidated candidates; APIs.guru remains first; RFC 9727
requires an explicit provider and fails safely on wrong media types, unsafe
hosts, malformed or oversized catalogs, link limits, redirects, timeouts, and
cancellation; cache and CLI contracts are covered; Go, vet, standalone,
downstream, CLI, and diff checks pass.
````

## Status record

````markdown
# Status M70 - Bounded Remote Catalog Discovery

State: Complete

## Scope

Add an APIs.guru-first remote discovery chain with an experimental LAP
Registry adapter and explicit-provider RFC 9727 lookup. Preserve untrusted
metadata posture, bounded safe HTTP behavior, and offline deterministic tests.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| M70 bounded remote catalog discovery | `[+]` | Added LAP and RFC 9727 adapters, source ordering, provenance/trust flags, CLI/cache wiring, safety tests, docs, evolution v20, downstream checks, and milestone review. |

## Boundary Checks

- LAP candidates use the registry-reported original source URL and remain
  unvalidated until normal import.
- RFC 9727 requires an explicit provider URL, fetches one Linkset catalog, and
  never follows nested catalogs.
- No provider operations, credentials, account selection, general crawling,
  Swagger Catalog, or Scalar Registry behavior enters `apitools`.

## Completion Notes

- Automatic lookup stops at the first source with matching candidates and
  preserves an ordered attempt ledger.
- LAP and RFC 9727 candidates remain visibly experimental and unvalidated;
  the selected original document must still pass the normal import validator.
- RFC 9727 discovery uses the final catalog URL as the base for relative
  references, follows normal safe redirects, caps inspected `service-desc`
  links, and does not traverse `api-catalog` relations.
- Evolution v20 records the accepted remote discovery boundary and result.

## Verification

- `go test ./...`
- `go vet ./...`
- `GOWORK=off go test ./...`
- `GOWORK=off go vet ./...`
- `go run ./cmd/apitools search --help`
- `go run ./cmd/apitools import --help`
- `go run ./cmd/apitools catalog check`
- `(cd ../openudon && go test ./...)`
- `(cd ../udon && go test ./...)`
- `git diff --check` in `../apitools`
- `git diff --check -- apitools` in `../tofu`
````
