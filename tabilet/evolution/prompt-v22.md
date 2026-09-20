# Prompt V22 - Untrusted Metadata And Ownership Hardening

Resolve the reviewed parser, prompt-safety, network, filesystem, cache,
catalog, CLI, and OAuth-boundary weaknesses as one coordinated breaking
release. Apply uniform budgets to untrusted source parsing and prompt-facing
summaries, preserve OpenAPI security alternatives, make artifact provenance
and writes confined and atomic, generate the provider catalog from validated
review data, and move live Google OAuth authorization-code exchange into Udon.

Use the lane-aware harness to run source and cross-cutting hardening with
explicit ownership, then serialize catalog generation and final sibling
migration. Break the public v0 contracts now instead of retaining flattening or
ambiguous compatibility fields. Keep apitools metadata-only: it may describe
auth requirements but must not resolve credentials, fetch tokens, execute API
operations, or choose runtime accounts.
