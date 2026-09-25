# Retired milestone S01 - Untrusted Source And Prompt Contract Hardening

**Milestone.** S01
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-S01.md
**Source specification.** tabilet/memory-bank/milestone.md#s01---untrusted-source-and-prompt-contract-hardening
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## S01 - Untrusted Source And Prompt Contract Hardening

**Goal.** Apply one bounded untrusted-input and prompt-safety contract to all
source parsers, inventories, authoring summaries, and security alternatives.

**Scope.** Bound parser depth and work, guarantee lexer progress, bound inline
content and reference expansion, sanitize and budget prompt-facing summaries,
preserve OpenAPI security OR/AND semantics, resolve bounded local composition
references, and return visible diagnostics instead of silently truncating or
discarding selected operations.

**Dependencies.** M71. The public security and authoring contracts intentionally
break and require coordinated OpenUdon and Ramen migration before release.

**Parallel ownership.** S01 owns parser packages, root inventory/authoring and
security files, and the corresponding OpenUdon/Ramen consumers. M72 may proceed
in parallel in download, refresh, materialization, cache, and Udon OAuth files.
Shared documentation, dependency pins, and review commits are serialized.

**Downstream impacts.** OpenUdon must select among credential alternatives and
Ramen must render bounded structured metadata without flattening security sets.

**Acceptance.** Malformed and deeply nested sources fail within deterministic
limits; prompt content is sanitized and budgeted; security requirement sets
retain AND-within/OR-between semantics; duplicate ranking and bounded reference
resolution are correct; fuzz/regression, full, standalone, downstream, vet,
staticcheck, CLI, and diff checks pass.
````

## Status record

````markdown
# Status S01 - Untrusted Source And Prompt Contract Hardening

**Acceptance.** All parser and authoring inputs are bounded, prompt-facing
metadata is sanitized and budgeted with visible diagnostics, OpenAPI security
alternatives remain structurally correct, downstream authoring consumers adopt
the breaking contract, and milestone verification passes.

| Item | State | Notes |
|---|---|---|
| Bound and fuzz all untrusted parser paths | `[+]` | Commit `94e0ded` adds shared 20 MiB, depth-100, and 100,000 semantic-work limits; fixes stray-brace and EOF panics; bounds JSON, XML, text, descriptor, introspection, and binary wire parsing; and adds four fuzz targets plus the discovered protobuf crash corpus. Follow-up `80dfcd3` gives the linear JSON/YAML structural preflight its own one-million-item ceiling so reviewed Kubernetes-sized documents remain inside the byte bound. Focused fuzz smoke, full/standalone tests, vet, and diff checks pass. |
| Bound inventories and schema/reference expansion | `[+]` | Commit `60a798d` enforces inline 20 MiB and default 10,000-operation bounds, reports partial inventories, caps project URL imports at 16 before writes, preflights OpenAPI/AsyncAPI YAML and JSON, bounds schema expansion at 4,096 nodes with cancellation, resolves local path/parameter/body/response refs, merges `allOf`, exposes optional `oneOf`/`anyOf` unions with readiness blockers, and tolerates numeric Swagger 2.0 while retaining strict OpenAPI versions. Full, standalone, vet, and diff checks passed. |
| Add the prompt-safety budget layer | `[+]` | Commit `3a6df69` adds one exported prompt budget/sanitizer report, strips ANSI and control characters, normalizes whitespace, caps identifiers/text/collections/fields, enforces 32 KiB per operation and 512 KiB ranked contexts, preserves exact duplicate-operation indexes, surfaces every truncation, and returns a blocking diagnostic instead of dropping an oversized selected set. Authoring build/rank APIs now preserve diagnostics; full, standalone, vet, staticcheck-baseline, and diff checks passed. |
| Preserve security alternatives and migrate authoring consumers | `[+]` | Apitools `8ea8adc` replaces flattened security with ordered OR alternatives containing AND requirements and explicit anonymous sets. Authoring `5c529b6` publishes `authoring.prompt-context.v2`; OpenUdon `00a5a4e` forces a stable resume-safe selection and selected-field mapping; Ramen `326891a` rejects v1, blocks direct/static unresolved alternatives, and emits only incomplete non-runnable drafts. Full Apitools, Authoring, OpenUdon, Ramen, and Udon tests pass under the workspace; standalone pin verification remains in the S01 review/M73 release row. |
| Review and verify S01 | `[+]` | Deep review confirmed the six reported v2 regressions were already covered by the coordinated fixes and added Apitools `76c0bb1`, Authoring `d73b65a`, and OpenUdon `6ed56d9` for compatibility-parser/value bounds, bounded query/document/selection work, visible shortlist diagnostics, selected-only ranked context, prompt fail-closed behavior, and clean shared adapters. Four parser fuzz smokes, Apitools full/standalone/race/vet/staticcheck/CLI/diff gates, Authoring full/standalone/vet/staticcheck/compatibility, workspace OpenUdon/Ramen/Udon full tests and vet, and all repository diff checks pass. OpenUdon/Ramen/Udon standalone dependency-pin proof remains the explicit M73 release row. |

## Boundary Checks

- Parsing and summaries remain metadata-only and never execute an operation.
- Sanitization is structural and budget-based; it does not invent a semantic
  keyword denylist for untrusted descriptions.
- Remote `$ref` fetching, credential resolution, and workflow choices remain
  outside apitools.
````
