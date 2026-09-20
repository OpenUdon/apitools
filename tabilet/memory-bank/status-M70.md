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
