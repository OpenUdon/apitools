# Prompt V26 - Official-source API version discovery

Plan opt-in discovery of the API versions a publisher actually offers, from
official sources beyond catalog results. Distinguish catalog preference from
the highest verified comparable version found within the checked official
scope, keep source URLs and verified full-content digests, and report conflicting
evidence instead of resolving it silently. Never claim global latest.
Choose versions by the requested capability, since a newer API may not cover
every older operation.

A version check must be cheap. The version already held is the baseline and
the product; searching for a later one is advisory and must never delay or
fail it, even when an LLM or browser is involved. The caller keeps the baseline
available independently of network completion; initialize the report from local
evidence before network work. A timeout or rate limit yields a partial report with
unexamined entries, and no check ever claims the latest version. Tiers run
cheapest first and under one shared time budget: stored results, a conditional
re-fetch of the known official URL, bounded probing of nearby versions, and the
publisher catalog or one official documentation page. Newer candidates are
preferred within the shared document-body cap; conditional checks, including
unchanged or invalid fetched source bodies, consume applicable budget. Optional freshness state lives in
a versioned file, not the shared SQLite cache.

Keep existing search compatible: add a separate library API and
`apitools versions` command. Stay metadata-only, with no account operations,
credentials, or general crawling; network access happens only on explicit
opt-in and only inside declared official origins. APItools adds no LLM client
or browser dependency. LLM or browser help comes from consumers: they supply
candidate URLs, and APItools verifies and fetches them itself. An optional
pull wrapper may call a consumer-supplied adapter under a deadline, only after
the deterministic tiers.

Treat a spec's URL, not a copy of its bytes, as the source of truth, as
APIs.guru does. APItools holds an index of locators, not spec bytes. Because
not every spec has a static URL, each record keeps a locator type and the
recipe that found it, so a broken link is repaired by replaying the recipe
before any LLM is asked. The 2026-10-02 observation found no APIs.guru spec
updates after April 2023; baselines are labeled by their recorded age rather
than a permanent assumption, and its `list.json` is cached locally
instead of cloned. Saving a fetched document is an explicit, opt-in act that
writes only validated newer versions into a directory the caller names. A
shared, published directory built from these records, seeded by APIs.guru and
the public-apis list, is a later, legally gated candidate.

Source: "APItools official-source API version discovery handoff"
(`apitools-search.md`, 2026-10-01), demonstrated with GoDaddy's Domains
specs, and review "APItools S05 — Official-source API version discovery"
(`apitools-s05-review.md`). The user approved the plan on 2026-10-01 against
`8a52c3f602988b945a4b5c1960bce8c04170c63d` with a clean worktree, and approved
the review's amendments the same day. On 2026-10-02 the user approved the
locator-and-recipe model, the cached-list and baseline-age decisions, and the
explicit-save scope, and split the dead public-apis endpoint into its own
milestone.

The user approved "S06/S05 plan consistency review, 2026-10-06" at the same
full commit, with the six pre-existing planning files uncommitted and no code
changes. It clarifies v26, not a new direction: keep shared caps (unexamined
excess hints; no integrated reset), optional supplied baseline inventory
(unexamined comparisons if absent), and a scope-qualified `preferred`.
A 304 proves only its resource unchanged. Partial/over-cap reads cannot supply
verified full digests. State/cache persistence and document saves have separate
opt-ins; only fully fetched newer Import-valid OpenAPI/Swagger is save-eligible.
Enforce confined paths, credential-free clients and official scope on every
redirect/replay; shared hosting alone is not publisher authority. Sniffing
rejects definite non-spec bodies, not valid documents with late identifying keys.

Execution, merging into the workspace checkout, publication, and consumer
adoption each need separate authority. Kinet's Stage 5 is now complete; retain
one isolated implementation worktree to protect shared builds. Compatibility
checks must explicitly resolve to that candidate worktree, not unchanged
workspace sources. S05.1 still requires separate design-record approval.
No remote facts were remeasured in the 2026-10-06 reconciliation.
