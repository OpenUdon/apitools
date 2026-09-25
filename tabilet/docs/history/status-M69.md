# Retired milestone M69 - Bounded Local Multi-Family Discovery

**Milestone.** M69
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M69.md
**Source specification.** tabilet/memory-bank/milestone.md#m69---bounded-local-multi-family-discovery
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M69 - Bounded Local Multi-Family Discovery

**Goal.** Give downstream authoring one safe local discovery API for every
first-class source family without moving workflow choices or remote lookup into
`apitools`.

**Scope.**

- Add `DiscoverLocalSources(context.Context, LocalSourceDiscoveryOptions)` and
  a versioned report with validated candidates, rejected and ambiguous files,
  truncation diagnostics, family, title, operation count, score, path,
  SHA-256, and provenance.
- Inspect only explicit regular-file or directory roots; reject symlinks,
  unsafe/special paths, oversized files, and cancellation.
- Detect and validate OpenAPI/Swagger, Google Discovery, AWS Smithy, AsyncAPI,
  GraphQL, OpenRPC, gRPC/protobuf, and OData through native parsers. Treat
  conventional directory names as hints rather than proof.
- Deduplicate identical content by digest. Require an explicit kind for
  ambiguous JSON or XML.
- Default to 10,000 visited entries, 100 accepted candidates, and 20 MiB per
  file. Reaching a count bound produces a blocking truncation diagnostic with
  narrowing guidance.
- Extend the native operation inventory entry point to the same eight source
  families while preserving their source kinds and selectors.
- Keep remote discovery, workflow selection, operation execution, credential
  behavior, and project artifact writes outside this API.

**Acceptance.** Mixed-root discovery covers every supported family;
deduplication, ambiguity, symlinks, non-regular/oversized paths, both limits,
and cancellation fail safely; existing callers remain compatible; full Go,
vet, standalone, downstream, CLI, and diff checks pass.
````

## Status record

````markdown
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
````
