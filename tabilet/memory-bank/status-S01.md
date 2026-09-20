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
