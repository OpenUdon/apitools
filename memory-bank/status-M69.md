# Status M69 - Bounded Local Multi-Family Discovery

State: Complete

## Scope

Add a bounded, local-only multi-family discovery API for adaptive authoring.
`apitools` owns source facts and validation; downstream products still own the
active workflow, source choice, network policy, materialization, and writes.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| M69 local multi-family discovery | `[+]` | Added the versioned discovery contract, native validation for all eight families, explicit-root safety and bounds, SHA-256 evidence/deduplication, ambiguity handling, native inventory parity, tests, docs, boundary review, and required verification. |
| M69 review remediation | `[+]` | Aligned local scanning with every supported adjacent security-sidecar suffix so advisory auth metadata cannot become an ambiguous API-source blocker. Implementation: `ef32163`. |

## Completion Notes

- The API does not infer roots from the current directory and does not contact
  catalogs or remote services.
- A partial bounded scan is never silently successful: the report is marked
  truncated and carries an error-severity diagnostic.
- Evolution v19 records the public discovery boundary and its downstream use.

## Verification

- `go test ./...`
- `go vet ./...`
- `GOWORK=off go test ./...`
- `GOWORK=off go vet ./...`
- `go run ./cmd/apitools search --help`
- `(cd ../openudon && go test ./...)`
- `git diff --check` in `../apitools`
- `git diff --check -- apitools` in `../tofu`
