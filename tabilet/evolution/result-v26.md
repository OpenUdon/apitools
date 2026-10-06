# Result V26 - Official-source API version discovery

**State:** Approved plan, 2026-10-01, amended on 2026-10-01, 2026-10-02,
and 2026-10-06;
not implemented or accepted.
[S05](../memory-bank/milestone.md#s05--official-source-api-version-discovery)
has six pending rows in [status-S05.md](../memory-bank/status-S05.md) and runs
after [S06](../memory-bank/milestone.md#s06--public-apis-source-repair), a
separate P2 defect fix for the dead public-apis endpoint that is not part of
this direction. No closing review gate has started, and the plan grants no
implementation, merge, or release authority.

The material direction change is that APItools will discover versions from
official publisher sources, not only from catalogs, and will accept LLM or
browser help only as untrusted hints. Today, search keeps APIs.guru's
preferred version (falling back to a lexicographic sort), `auto` stops at the
first source with matches, and only the experimental RFC 9727 path consults a
publisher. GoDaddy illustrates the gap: it is absent from the catalog, its v1
and v2 Domains specs both declare `info.version` 1.0.0, its v3 spec lacks v1's
transfer operations, and v1 and v2 fail strict parsing because of a publisher
schema defect.

The amended plan makes the check cheap by contract. The held version is the
baseline, available independently of network completion; a timeout, rate
limit, or bad response gives a
partial report and never an error or a "latest" claim. Starting values are 5 s
total for the deterministic tiers, 3 s per request, at most 4 concurrent
requests, 12 requests, and 3 fetched document bodies, with 8 MiB per checked
document. Existing primitives would not
meet this on their own: the client timeout is 30 s, downloads read whole
bodies, no code reads response headers, and nothing uses concurrency. S05 adds
an unexported probe primitive, zero-network evidence including APIs.guru's
`x-origin`, an opt-in freshness state file, and a push-first adapter whose
verify step is cheap while the consumer runs slow LLM or browser work in the
background.

S05 first records an approved design, then adds the probe primitive and
catalog evidence, the deterministic upgrade check, hint verification, a
capability-relevant diff against the local baseline, and an opt-in
`apitools versions` command. Existing search, catalog data, registrations,
index inputs, the shared SQLite schema, and consumer pins stay unchanged.
Fixtures are synthetic with simulated latency, and no third-party spec is
committed.

The 2026-10-02 amendments make a spec's location, not its bytes, the source of
truth. Each record holds a locator kind (direct, github, pointer, pattern,
template, or gated), the recipe that found it, the resolved URL, digest, and
status, so a moved or dead link is repaired by replaying the recipe, then
re-reading the pointer, before any LLM tier. Facts measured on 2026-10-02 shaped this:
APIs.guru's newest spec update is 2023-04-21, every entry carries `x-origin`,
and its 907 KB `list.json` answers a conditional request with a 304, so S05
caches the list instead of cloning a 766 MB repository. For GoDaddy the
APIs.guru origin is `swagger_domains.json` with no version token, so sibling
probing cannot reach `/openapi/domains-v3.json`; the pointer or opt-in LLM
tier is needed. A baseline from APIs.guru is labeled by age. The explicit save
writes only newer versions that passed validation, never the baseline, an
invalid or oversize document, or source bytes without a caller directory.
Metadata state/cache persistence is separately opt-in.

Implementation remains on one branch in one separate worktree outside the
shared Go workspace, preserving normal consumer builds. Kinet's Stage 5 is
complete, not a pending dependency or merge gate. GoDaddy catalog curation, search default changes, identifying
APItools on shared downloads, and a shared, published provider version
directory (metadata only, legally gated) stay as candidate directions. Current product,
architecture, tech-stack, lessons, and frozen history are unchanged.

The approved "S06/S05 plan consistency review, 2026-10-06" revalidated
`8a52c3f602988b945a4b5c1960bce8c04170c63d` with the six uncommitted planning
files and no code changes; source priorities and separate review baseline were
not supplied. It amends the existing nine pending rows and preserves both zero
review counters. Findings and owners are recorded in
[S05](../memory-bank/status-S05.md#consistency-intake-findings-2026-10-06)
and [S06](../memory-bank/status-S06.md#consistency-intake-findings-2026-10-06).

Shared limits now explicitly cover conditional checks, revalidation, redirects,
recipe replay and integrated hint verification, without resets. Up to eight
hints may be accepted, but excess document demand remains unexamined; standalone
verification uses the same caps. A 304 is one request in an isolated fixture
and says nothing about other versions. Comparison uses optional supplied
baseline inventory; missing inventory is unexamined, not an invalid request.
`preferred` is scoped highest verified comparable evidence, never global
latest; conflicts cannot establish a new preference.

Only complete bounded reads have verified full digests. Confined document saves
are distinct from metadata persistence, and parsed Discovery is not automatically
eligible under OpenAPI/Swagger Import validation. Credential-free clients enforce
official origin/repository scope on every hop. Ambiguous sniff prefixes continue
bounded parsing. S06 also preserves valid empty legacy JSON while rejecting
unrecognized/empty Markdown and refusing to present truncated lists as complete.

S06.3 establishes candidate-bound OpenUdon/Udon tests through disposable
workspace/modfile configuration; S05.6 reuses them. Unchanged-pin/replacement
standalone checks remain separate baseline evidence. Ramen stays excluded;
sibling files, consumer pins, merge and publication authority remain unchanged.
The four existing public-apis baseline tests passed in workspace and standalone
modes, and the patch check passed; no new feature acceptance is claimed.
Live measurements were not repeated. S05.1's detailed design approval remains
a checkpoint. No evolution bump is needed: this clarifies v26 rather than
changing direction.
