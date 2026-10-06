# Retired milestone S05 - Official-source API version discovery

**Milestone.** S05
**Outcome.** completed
**Retired.** 2026-10-06
**Source status.** tabilet/memory-bank/status-S05.md
**Source specification.** tabilet/memory-bank/milestone.md#s05--official-source-api-version-discovery
**Evidence.** 515663f8bf14793dcdd335dd67fadf79de7a1984
**Worktree.** includes uncommitted changes
**Review.** passed
**Review iterations.** 2
**Verification.** APItools full workspace/standalone tests and vet; full race tests; catalog generator/quality (zero findings); search/versions CLI help; patch/structural/link checks; candidate-bound OpenUdon/Udon full tests/vet in workspace and standalone modes; local-server deadline, budget, state, source-family/service, origin/repository, recipe/hint, native diff/rank and save regressions. Go 1.26.6; Ramen excluded. S05.6 implementation/closure remains uncommitted at this recorded observation; the evidence commit includes S05.1-S05.5. Live provider endpoint availability is explicitly unexamined (403); default verification uses fixtures only.
**Consolidated into.** README.md, docs/api-version-discovery.md, tabilet/memory-bank/product.md, tabilet/memory-bank/architecture.md, tabilet/memory-bank/tech-stack.md and tabilet/memory-bank/lessons.md.

## Milestone specification

````markdown
## S05 — Official-source API version discovery

Add opt-in discovery of the API versions a publisher actually offers, from
official sources beyond catalog results, without changing existing search.
This is a requested feature, not review remediation, so it carries no review
severity. It was requested in "APItools official-source API version discovery
handoff" (`apitools-search.md`, 2026-10-01) and planned against
`8a52c3f602988b945a4b5c1960bce8c04170c63d` (published source
`fb132631c9827eae5f2ec4503d03f21eabfb4113`) with a clean worktree. The user
approved the scope, external-adapter integration, lane S, the new `versions`
command, and this plan on 2026-10-01.

A second approved intake the same day applied review "APItools S05 —
Official-source API version discovery" (`apitools-s05-review.md`), revalidated
at the same commit with these uncommitted planning files as evidence. It made
a version check cheap by contract: the version already held is the baseline,
and a check can never delay or fail it. It rewrote the six pending rows in
place (same IDs and count) and added the latency and cost contract below. The
review's per-finding dispositions are in S05's status file.

A third approved intake on 2026-10-02 applied this conversation's findings
(F01-F12 in S05's status file), revalidated at the same commit. It records
measured APIs.guru and public-apis facts, corrects the GoDaddy example, adds
the locator-and-recipe model for specs without a static URL, adds a cached
APIs.guru list as tier 0, decides what the explicit save writes, and runs
after new milestone S06. No intake starts the closing review counter.

The fourth approved intake, "S06/S05 plan consistency review, 2026-10-06",
revalidates the same full commit with all six uncommitted planning files and
no uncommitted code. Source priorities and a separate review baseline are
`not supplied`. S05 owns F01-F07; S06 owns F08-F10. These are planning gaps,
not newly demonstrated production exploits. The user selected shared caps,
an optional baseline inventory, and a scope-qualified preference, then
approved all ten dispositions and six-file actions. Existing live measurements
remain dated historical evidence; no remote fetch was made in this intake.
The nine row identities and both zero review counters remain unchanged.

**Observed problem.** For a GoDaddy query, `search` returned an older Domains
spec. APIs.guru search keeps each entry's `preferred` version and otherwise
sorts version keys lexicographically (`providers.go`); `auto` stops at the
first source with matches (`client.go`); and only the experimental RFC 9727
path, when a provider URL is supplied, consults a publisher. GoDaddy is not in
the built-in catalog. Its official Domains specs show why "latest" needs
evidence: v1 and v2 both declare `info.version` 1.0.0 while v3 declares 3.1.0;
v3 has 14 operations and none of the transfer operations that v1 (23
operations) and v2 (30) carry; and v1 and v2 fail strict parsing because their
`JsonSchema` component sets `properties.type` to the string `"object"`, a
publisher defect.

Measured on 2026-10-02: the APIs.guru baseline for that query is the entry
`ote-godaddy.com:domains`, a Swagger 2.0 file whose recorded origin is
`https://developer.godaddy.com/swagger/swagger_domains.json` (updated
2023-03-06). That URL has no version token and sits under a different path from
the newer `https://developer.godaddy.com/openapi/domains-v1.json` through
`-v3.json`, so sibling probing from the APIs.guru origin cannot reach them; a
pointer (a docs page) or the opt-in LLM tier is needed. Sibling probing helps
only publishers whose spec URLs carry a version token (57 of 178 official
machine-readable catalog references). APIs.guru itself is frozen: the newest
`updated` timestamp across its 3,992 specs is 2023-04-21 and none changed in
the last 12 months, so a baseline taken from it is old by construction. Every
APIs.guru entry carries `info.x-origin` (the publisher's original URL); its
`list.json` is 907 KB on the wire (0.43 s) and answers a conditional GET with a
304 in 0.31 s; and its repository (about 766 MB with history) does not contain
`list.json`, so a clone would not replace it.

**Boundaries.** Discovery is metadata only: no account operation, credential,
token, cookie, or signed request. Existing `Client.Search`, `apitools search`,
wire shapes, built-in catalog data, artifact registrations, catalog index
inputs, and consumer pins stay unchanged. There is no general crawling:
network access happens only when a request opts in, and only inside official
origins named by the caller or by catalog spec references with official
source authority. APItools adds no LLM client or browser dependency; LLM or
browser help arrives through an external adapter whose output is untrusted
hints. Optional freshness state is a versioned JSON file at a caller-supplied
path; it holds no secrets, is never committed, and adds no SQLite schema
change, because older readers reject a registry whose schema version is higher
than their own. APItools holds no spec bytes by default: the source of truth
for a version is its locator (where it was found and how), the digest records
what was verified, and bytes are saved only on an explicit request into a
caller-named directory that is not part of the repository. Fixtures are
synthetic, and no third-party provider spec is committed.

**Latency and cost contract.** The user's requirement is that searching for a
later version of an API already held must not take long or time out. The held
version is the baseline and the product; the check is advisory.

- The caller's baseline remains available independently of network completion;
  initialize the report from local evidence before network work. A timeout, rate
  limit, unreachable host, or bad response yields a partial report with
  `unexamined` entries, never an error that drops the baseline. Errors are
  reserved for invalid requests. A check never claims the "latest" version.
  Its statuses are `newer_found`, `none_found_in_scope`, `unexamined`, and
  `conflicting`, each with `checked_at`.
- The request carries the known baseline (provider key or spec URL, version,
  digest, and an optional already-built operation inventory; local bytes or a
  registration may identify existing evidence but do not authorize reparsing).
  The baseline is never refetched or reparsed for comparison. Without its
  inventory, capability comparison is `unexamined`, not an invalid request.
  A conditional source check is separate from the held baseline and never
  replaces it. Discovery looks upward only: order known comparable newer
  candidates newest first within the shared document budget; older versions
  are listed as unfetched candidates with URLs.
- Tiers run cheapest first. Tier 0 uses no network: a fresh prior check
  inside its time to live (24 hours), the catalog reference, and a locally
  cached APIs.guru `list.json` (at a caller-supplied path outside the
  repository) supplying the entry's `x-origin` URL, version map, and `updated`
  age. Any ETag revalidation after 24 hours is network work charged to the
  shared budget, not zero-network tier 0; a 304 costs no body. Tier 1 is one
  conditional GET (`If-None-Match` or `If-Modified-Since`) on the known
  official URL, so 304 means only that resource is unchanged, not that no
  newer version exists; a changed full-content digest means an in-place
  revision. Tier 2 is bounded same-host sibling probing from the known version
  token (for example `v1` to `v2`, `v3`), upward only, stopping after two
  consecutive misses. Tier 3 is the RFC 9727 publisher catalog and one
  supplied official documentation page. Tiers 1-3 start together under one
  shared wall-clock budget rather than receiving separate deadlines. Tier 4,
  adapter hints, runs only afterwards, only when the question is still open,
  and only when the caller opted in.
- Not every spec has a static URL, so a record stores a locator and the recipe
  that found it, not only a URL. Locator kinds are `direct` (a fixed spec
  URL), `github` (repository, path pattern, and ref), `pointer` (a docs page,
  RFC 9727 catalog, or `latest` alias that yields the current spec URL),
  `pattern` (a URL with a version slot), `template` (tenant-specific, recorded
  for guidance and never verified generally), and `gated` (login-only, or only
  human docs, which falls back to docs-derived overlays). A record also holds
  the last resolved URL, the digest, `checked_at`, and a status of `ok`,
  `moved`, `dead`, or `gated`. Repair order when a link breaks: replay the
  recipe, then re-read the pointer, then the opt-in LLM tier; `github` locators
  are verified with one tree listing per repository. The record is shaped to
  be compatible with APIs.guru's `list.json` (`preferred`, spec URLs,
  `info.x-origin`, `updated`, with APItools additions under `x-apitools-*`
  keys) so a future directory can aggregate it. `preferred` names only the
  highest verified comparable version found within the recorded checked scope,
  never the publisher's globally latest version; conflicting or incomparable
  evidence cannot establish a new preference.
- A baseline taken from APIs.guru is labeled with its entry's `updated` age
  and marked stale by source when older than 12 months. The report then
  suggests the opt-in LLM tier; it never runs it unasked.
- Starting values, which S05.1 may adjust only with measured evidence and the
  user's approval of its record: 5 s total for tiers 1-3, 3 s per request, at
  most 4 concurrent requests, 12 requests, 3 fetched document bodies, and 8 MiB
  per checked document. Conditional fetches, directory-cache revalidation,
  redirects, locator repair, and hint verification consume the applicable
  shared request/time/body budget; integrated tiers never reset it. A completed
  source-body fetch consumes a document slot even if validation fails or it
  proves unchanged. Oversized or partial bodies remain `unexamined`: distinguish
  observed bytes from declared length and omit a verified full-content digest
  unless the complete bounded document was read. The consumer may fetch a large
  document separately. The 30 s client timeout is not inherited.
- The adapter is push-first. A verify call (`VerifyHints`) accepts at most eight
  size-bounded hints, deduplicates and orders them deterministically, and fetches
  selected URLs at most once within the available shared caps. Excess hints
  remain `unexamined`, not silently dropped or verified. A standalone call uses
  the same bounds; integrated verification consumes the remaining budget, not
  a fresh allocation. The consumer runs slow LLM or browser work in the
  background and calls verify when hints arrive. An optional
  pull wrapper may call a consumer-supplied adapter, only after the
  deterministic tiers left the question open and only when the caller opted
  in, under a caller deadline (default 15 s, ceiling 60 s). Adapter timeout,
  error, or garbage never fails the check or changes a verified result; it is
  reported as a tier status.

**Scope and acceptance.**

- **Design record and checkpoint.** `docs/api-version-discovery.md` records
  the versioned request (with the known baseline) and report (proposed
  `apitools.api-version-discovery/v1`); the latency and cost contract above,
  with its starting values; the baseline-first and partial-report semantics;
  the tier order; the status vocabulary; the evidence labels
  (catalog-preferred, catalog-highest, official-candidate, official-verified);
  the official-origin policy, naming three sources of an origin (a caller
  declaration, a catalog reference with `official-*` authority, and APIs.guru's
  `x-origin`) and treating `x-origin` as a lead that becomes official-verified
  only when the caller declares the official scope or it matches a catalog
  reference's `official-*` scope. A shared hosting domain alone never establishes
  publisher ownership; repository/path identity must remain within that scope.
  Also record version identity and ordering, comparing URL version
  tokens, `info.version`, and documentation claims and reporting conflicts
  rather than resolving them silently; capability relevance; the adapter
  contract; the documentation-scanner decision (a bounded standard-library
  scanner, not a new dependency); and consumer notes (interactive consumers
  run tiers 0-3 inline within the budget and treat everything else as a
  background advisory that never gates authoring), plus the locator kinds,
  recipe, and status fields above, shaped to be `list.json`-compatible; the
  baseline-age label; and the save semantics below. It records the facts
  measured on 2026-10-02 (above) and verifies the one APIs.guru fact still
  unchecked: whether a provider-scoped endpoint exists. Its consumer notes
  also say that OpenUdon keeps its `catalog import-openapi` wrapper and its
  package installation (`step source add`), and that any relocation of the
  wrapper is OpenUdon's decision. The user approves the record before later
  rows are implemented.
- **Probe primitive, catalog evidence, and freshness state.**
  - A new unexported probe primitive reuses the existing URL, redirect, and
    dial guards and leaves `downloadBounded` and every existing caller
    unchanged. Its isolated client strips cookie jars and caller redirect
    callbacks, as M80's discovery client does, and enforces approved origins on
    every request and redirect, including recipe replay. Safe public hosts
    alone do not satisfy the official-origin policy. It supports conditional
    GET, a first-bytes sniff that abandons definite non-spec responses such as
    HTML, a per-request cap, and a stable tool identifier. Ambiguous JSON/YAML
    prefixes continue to bounded parsing; a missing early OpenAPI, Swagger, or
    Discovery marker is not grounds for rejection. Treat 404 as a miss;
    treat 429 and 5xx as `unexamined`, with no retry inside the
    budget and `Retry-After` recorded.
  - Zero-network evidence from catalog spec references and the locally cached
    APIs.guru `list.json` (network revalidation after the 24 hour time to live
    consumes the shared budget; no clone), parsing `x-origin` and the entry's
    `updated` age and, if S05.1 verifies it, using a provider-scoped fetch
    instead of the whole list. Versions are ordered by a natural comparison
    rather than lexicographically, with the catalog's
    preferred version labeled separately from the highest found. A catalog
    result never becomes official-verified.
  - Opt-in freshness state as the versioned JSON file described above, or a
    `previous` report passed back in, holding each record's locator, recipe,
    resolved URL, digest, status, and `checked_at`. A fresh state costs zero
    network. State and directory-cache persistence each need explicit caller
    paths/permission, independently of the document save directory; no paths
    means no writes. Cached evidence cannot expand the current official scope.
- **Deterministic upgrade check.** Tiers 1-3 under one shared budget with
  bounded parallelism, upward-only from the known baseline, with the gap rule,
  the request cap, and the shared three-document-body cap. Candidates are
  generated only from integer or dotted version tokens; date-style tokens are
  not guessed. Each version records its source URL, final URL, byte evidence,
  fetch time, and parse/validation status; a verified digest is present only
  for a complete bounded read. A version that fails strict parsing is reported
  with its diagnostic, not dropped. A bounded
  standard-library scanner reads at most one supplied official documentation
  page within a size cap and link cap and matches only spec-like paths. A page
  that exposes no static spec links is reported as examined with none found,
  and the check stops. Each version found records the locator kind and recipe
  of the tier that found it, and a `moved` or `dead` result replays that
  recipe, then re-reads the pointer, before any LLM tier. Report order does
  not depend on completion order.
- **Explicit save of fetched documents.** Opt-in through a caller-supplied
  save directory (library option and `--save-dir`); without it no source
  document is written, and there is no default location. Explicit metadata
  state/cache persistence is a separate permission. It writes only a newer
  version that was fetched during the check and passed the same validation `Import`
  applies, byte for byte, so the saved file is importable by `step source add`
  and matches the digest in the report. Current Import validation accepts
  OpenAPI/Swagger; Discovery metadata is not automatically eligible for saving.
  It never writes the baseline (already held), a document that failed validation
  (reported with a digest only if fully read, diagnostics, and locator), or one
  over the 8 MiB check cap. It reuses `Import`'s
  naming and digest-idempotence behavior behind confined filesystem checks:
  reject unsafe directories, symlinks, and non-regular targets, and never
  overwrite differing content. The existing helper alone is not a confinement
  guarantee. Report each saved path, digest, and size, never extend the check
  budget, and treat a failed save as a warning on that version, not a failed
  check. Saved files are a local cache:
  never committed, and not registered or indexed by S05.
- **Hint verification and optional adapter.** The APItools-owned
  `VerifyHints` call and versioned JSON hint format through which a consumer
  such as OpenUdon, Kinet, or Browsertools supplies LLM- or browser-derived
  candidate URLs and evidence notes (at most eight, size-bounded). APItools
  validates hints and fetches only the budget-selected URLs itself under the
  official-origin policy, labels hint provenance and unexamined excess, and
  never accepts spec bytes, credentials, cookies, or authority upgrades from an
  adapter. The optional pull wrapper
  follows the latency contract above. Tests use a fake adapter.
- **Capability-relevant selection.** When an already-built baseline inventory
  is supplied, diff each fetched newer version against it. Otherwise mark the
  comparison `unexamined` without fetching/reparsing the baseline or rejecting
  the version check. Optional `StepContract` ranking across fetched versions
  uses the existing operation-candidate path. Reports name which versions cover
  the requested capability and flag conflicting evidence; a newer version never
  replaces an older version's
  operations by default.
- **Command, documentation, and verification.** A new opt-in
  `apitools versions` command over the library API, with `--timeout`,
  `--no-network`, `--state`, and `--save-dir`, and README and architecture updates for the
  delivered behavior. Provider-free fixtures use local servers and simulated
  latency and reproduce sibling versions, equal `info.version` values, a
  strictly invalid version, a newer version missing operations, off-origin
  redirects, unsafe hosts, and adapter hints. A run without opt-in performs no
  network access. Timing fixtures prove the contract: a stalled host returns
  within the budget plus a small slack with the baseline intact and the
  affected tiers `unexamined`; all-404 siblings stop at no more than 12
  requests; a fresh state costs zero requests and an isolated conditional
  check returning 304 costs one (concurrent tiers may issue other requests);
  a large 200 HTML response to a probe is abandoned after a few KiB; 429 and 5xx
  responses are `unexamined` with no retry; an adapter timeout leaves the
  deterministic result unchanged; report order is stable under permuted
  completion order; cancellation leaves no goroutines; a saved file's digest
  equals the reported digest and a rerun reuses it; an invalid, over-cap, or
  baseline document is never written; no source document is written without a
  save directory, and no metadata is persisted without its separate opt-in;
  an unsafe save directory is refused; a failed save never fails the check;
  a `moved` locator is repaired by replaying its recipe; and a
  cached `list.json` revalidation with a 304 costs one request against the
  shared budget. Also test eight hints exceeding the document cap, aggregate
  budgets across redirects/replay, unchanged conditional bodies, missing
  baseline inventory, scoped/conflicting preference, partial-read digest
  omission, Discovery save ineligibility, symlinked save paths, valid documents
  with late identifying keys, cookie/callback isolation, and off-origin/shared-
  host repository escapes. Verification:
  `go test ./...`, `GOWORK=off go test ./...`, `go test -race` over the new
  concurrent paths, `go vet ./...`, `GOWORK=off go vet ./...`,
  `go run ./cmd/cataloggen -check`, `go run ./cmd/apitools catalog check`,
  `go run ./cmd/apitools search --help`, `go run ./cmd/apitools versions
  --help`, `git diff --check`, and OpenUdon/Udon compatibility checks in
  workspace and standalone modes per tech-stack.md, with Ramen excluded.

**Dependencies and execution.** No upstream prerequisite; retired M80 and the
existing remote-discovery and guarded-download code are foundations. S06 is complete and retired; S05 now follows it under one execution owner
in the same isolated branch/worktree. It has no technical dependency on S06's
new list parser; its probe primitive remains separate. Kinet's Stage 5 is complete, not a pending dependency.
Keep the branch outside `~/Workspace/go.work` so normal consumer builds retain
current workspace behavior; merging and publishing each need separate
approval. S05.6 reuses the isolated candidate-bound verification procedure in
[retired S06](../docs/history/status-S06.md),
with ordinary pinned standalone checks labeled separately. No consumer adopts
S05 in this milestone; OpenUdon or Kinet adoption is reconciled separately in
its own ledger. Curating GoDaddy into the
catalog stays under the "Next provider-catalog expansion" candidate.
````

## Status record

````markdown
# Status S05 — Official-source API version discovery

**State:** Complete, 2026-10-06, in isolated branch `work/s06-s05`.
All six task outcomes, required verification and whole-milestone review passed
in two iterations. No merge, publication, or consumer adoption occurred.

**Specification:** [S05](milestone.md#s05--official-source-api-version-discovery).

**Design record:** [api-version-discovery.md](../../docs/api-version-discovery.md).
Approved by the user on 2026-10-06; S05.1 is complete.

**Provenance (first intake):** Requested feature from "APItools official-source
API version discovery handoff" (`apitools-search.md`, 2026-10-01). It is not
review remediation and carries no review severity. Planned against
`8a52c3f602988b945a4b5c1960bce8c04170c63d` (published source
`fb132631c9827eae5f2ec4503d03f21eabfb4113`), with a clean worktree. The user
approved the scope, the external-adapter integration, lane S, the new
`versions` command, the candidate edits, the evolution pair, and the
planning-file actions on 2026-10-01. Execution, merge, and release authority
remain separate.

**Provenance (second intake):** Review "APItools S05 — Official-source API
version discovery" (`apitools-s05-review.md`, 2026-10-01). The review gave no
baseline of its own beyond its revalidation, so its stated baseline is `not
supplied`. It was revalidated at the full commit
`8a52c3f602988b945a4b5c1960bce8c04170c63d`, with these uncommitted planning
files as part of the evidence and no uncommitted code. Source priorities are
the review's stated P2 and Lower; local severities were classified
independently under milestone.md's review severity definitions. P2 here means
a gap in the planned acceptance, not a defect in existing code. The user
approved the dispositions, the five decisions as approved scope, a 5 s total
budget with 3 s per request, the push-first adapter with an optional pull
wrapper, and the planning-file actions on 2026-10-01. This intake starts no
review counter.

**Provenance (third intake):** This conversation's findings of 2026-10-02,
labeled F01-F12 below, with source priority `not supplied` and no stated
baseline beyond the revalidation at the same full commit
`8a52c3f602988b945a4b5c1960bce8c04170c63d` (this plan's uncommitted files were
part of the evidence; no uncommitted code). The public facts were measured on
2026-10-02 from the live APIs.guru list and GitHub metadata, and the repository
facts at the commit above. The user approved the dispositions, the explicit-save
decision, the S06 split and its order, the candidate directions, and the
planning-file actions on 2026-10-02. This intake starts no review counter.

**Provenance (fourth intake):** "S06/S05 plan consistency review, 2026-10-06",
approved by the user. Source priorities and separate review baseline are
`not supplied`. Revalidation is at
`8a52c3f602988b945a4b5c1960bce8c04170c63d`, including all six pre-existing
uncommitted planning files and no uncommitted code. The user selected shared
caps, an optional supplied baseline inventory, and a scope-qualified
`preferred`, then approved the full ten-finding disposition and six-file
amendment. These are planning gaps, not newly demonstrated production
exploits. F01-F07 below belong to S05; F08-F10 belong to S06. IDs are scoped
to this intake and do not rename the earlier F01-F12. No row or closing
review is started. Prior live measurements remain historical observations;
this intake made no remote fetch.

**Evidence at planning (earlier intakes):**

- `providers.go`: APIs.guru search keeps each entry's `preferred` version and
  falls back to a lexicographic sort of version keys.
- `client.go`: `auto` search returns at the first source with matches; RFC 9727
  publisher discovery runs only with an explicit provider URL and is
  experimental (`remote_discovery.go`).
- `catalog/data/catalog.json`: no GoDaddy provider or candidate.
- Local copies of GoDaddy's official Domains specs from the requesting
  session: v1 and v2 (OpenAPI 3.1.1) both declare `info.version` 1.0.0, v3
  (OpenAPI 3.1.0) declares 3.1.0; operation counts are 23, 30, and 14, and v3
  has no transfer operations. v1 and v2 fail strict parsing in the OpenAPI 3.1
  parser because `components.schemas.JsonSchema.properties.type` is the string
  `"object"`; v3 validates.
- `tabilet/memory-bank/product.md` non-goals exclude general web crawling,
  vendoring provider specs, and credential handling. APItools has no LLM
  client, and Authoring does not import APItools.
- OpenUdon pins APItools `v0.0.0-20260930205753-fb132631c982` and calls
  `Client.Search` with an explicit APIs.guru source. Kinet's Stage 5 was at
  U07.4 in progress, with W8M W28, OpenUdon M95, Kinet M20, and W8M W29
  remaining at the earlier intake. The 2026-10-06 read of Kinet's current
  milestone ledger confirms Stage 5 is now complete; isolation remains required.
- Measured 2026-10-02 (third intake): APIs.guru's baseline for the GoDaddy
  query is the entry `ote-godaddy.com:domains`, a Swagger 2.0 file whose origin
  is `https://developer.godaddy.com/swagger/swagger_domains.json` (updated
  2023-03-06), with no version token and a path unlike the official
  `/openapi/domains-vN.json` files. The APIs.guru list has 3,992 specs, the
  newest `updated` is 2023-04-21 (none in the last 12 months), and every entry
  has `info.x-origin`. Its `list.json` is 907 KB on the wire (0.43 s, 8.9 MB
  parsed) and a conditional GET returns 304 in 0.31 s; its repository is about
  766 MB with history and contains no `list.json`. 2,632 origins sit in 154
  GitHub repositories and 1,381 on 769 other hosts. Of its specs, 27%
  declare an `info.license`.

**Lineage:** Builds on retired [M80](../docs/history/status-M80.md) and the
existing remote-discovery and guarded-download code without reopening them.

## Review findings owned here

| Finding | Source priority | Local severity | Disposition | Evidence | Owner |
|---|---|---|---|---|---|
| P2-1 No latency contract; reusable primitives would break it | P2 | P2 | confirmed | `client.go:18,21,22` (30 s client timeout, 5 s probe timeout, 30 s public-apis budget); `catalog_discovery_types.go:25-28` (30 s default, 60 s ceiling); `download.go:138-191` (full GET, no early abort); no goroutines, `sync`, or `errgroup` in non-test code | S05.1 contract; S05.3 implementation and timing fixtures |
| P2-2 "We already have one" not modeled | P2 | P2 | confirmed | Earlier S05.2 and S05.5 text listed, fetched, and diffed every version | S05.1 `known` field; S05.3; S05.5 |
| P2-3 No cheap-first tiers, in-place revision check, or freshness state | P2 | P2 | confirmed | `download.go:190` returns only body, final URL, and content type; 63 of 908 catalog references record ETag text; `sqlitecache/catalog_read.go:76` rejects a registry schema above `schemaVersion` (3) | S05.1 tiers and state decision; S05.2 conditional GET and state file; S05.3 |
| P2-4 Probe mechanics unspecified; downloader cannot support them | P2 | P2 | confirmed | `readBoundedResponseBody` reads the whole body up to 20 MiB; 57 of 178 official machine-readable references have a guessable integer `v<N>` URL token, and date-style tokens are not guessable | S05.2 primitive; S05.3 probing |
| P2-5 Adapter on the critical path | P2 | P2 | confirmed | Earlier S05.4 text had no deadline, ordering, or failure isolation | S05.4; S05.1 contract |
| P2-6 Free origin lead ignored; whole directory re-downloaded | P2 | P2 | partially confirmed. The code side is confirmed. Update 2026-10-02: `x-origin` is on every live entry and the list is 907 KB on the wire (see F04); only whether a provider-scoped endpoint exists is still unverified and left to S05.1 | `providers.go:194-202` does not parse `x-origin`; `catalog_discovery_remote.go:28,51` sets `Cache = nil` and `CacheModeBypass` | S05.1 trust policy and live verification; S05.2 parsing and scoped fetch if verified |
| P2-7 Docs link extraction has no parser | P2 | P2 | confirmed | no `golang.org/x/net` entry in `go.sum`; `AGENTS.md` prefers the standard library | S05.1 decision; S05.3 scanner |
| L1 Never claim "latest" | Lower | Lower | confirmed; required by the honest-result acceptance | spec wording | S05.1 vocabulary; S05.3 |
| L2 `-race` and concurrency fixtures | Lower | Lower | confirmed; required once concurrency exists. `-race` is feasible here (cgo enabled, gcc present) | no `-race` in documented verification | S05.3 fixtures; S05.6 verification |
| L3 Identify the tool on probes | Lower | Lower | confirmed; the probe-primitive part is owned by S05.2, and changing shared downloads is deferred | no `User-Agent` set anywhere | S05.2; candidate direction "Identify APItools on guarded downloads" |
| L4 Consumer pattern documentation | Lower | Lower | confirmed; documentation | none | S05.1 consumer notes; S05.6 |

## Third-intake findings (conversation, 2026-10-02)

F01 (the dead public-apis endpoint) is owned by [S06](../docs/history/status-S06.md). The rest
retain their original approval provenance below; the fourth intake qualifies
F12's original "verified latest" wording with checked scope.

| Finding | Source priority | Local severity | Disposition | Evidence | Owner |
|---|---|---|---|---|---|
| F02 APIs.guru is frozen, so a baseline from it is old by construction | not supplied | Lower | confirmed; required by the added baseline-age acceptance | newest `updated` 2023-04-21 across 3,992 specs; none in the last 12 months | S05.2, S05.3: label the age, mark stale after 12 months, suggest (never run) the opt-in LLM tier |
| F03 Cache the APIs.guru `list.json` instead of cloning | not supplied | requested design | design change, approved | 907 KB on the wire, 304 in 0.31 s; the repository (about 766 MB with history) has no `list.json`, so a clone would need about 4,000 files parsed | S05.2 tier 0; the clone is recorded as rejected in the directory candidate |
| F04 `x-origin` is on 100% of live entries | not supplied | Lower | confirmed; settles an S05.1 verification item | 3,992 of 3,992 entries | S05.1 policy (a lead, never official on its own); S05.2 parsing; also the LLM tier's starting context |
| F05 GoDaddy correction: sibling probing from the APIs.guru origin cannot reach `/openapi/domains-vN.json` | not supplied | Lower | confirmed | origin `.../swagger/swagger_domains.json`; probing helps 57 of 178 official refs | The spec's observed-problem text (corrected) and S05.1 |
| F06 Specs without a static URL need a locator and a recipe, with self-healing | not supplied | requested design | design change, approved | 4% of APIs.guru origins unreachable by its own metrics; 66% of origins are GitHub raw URLs | S05.1 (kinds and record), S05.2 (state), S05.3 (record and replay), S05.4 (hint recipes) |
| F07 Explicit save of fetched documents: only newer versions that passed `Import` validation, byte for byte, into a caller directory; never the baseline, invalid, or over-cap documents | not supplied | requested design | design change, decision made by the assistant at the user's request and approved | `step source add` validates with APItools, so saved files must be importable; the digest is enough evidence for invalid ones | S05.1 (decision), S05.3 (implementation), S05.6 (`--save-dir`, docs) |
| F08 OpenUdon's `catalog import-openapi` wrapper and `step source add` stay in OpenUdon | not supplied | n/a | outside ownership | the wrapper calls APItools `Client.Import`; `step source add` installs into workflow packages under approved digests | Handoff note in S05.1; no sibling edit |
| F09 Shared provider version directory | not supplied | n/a | deferred | needs legal sign-off and a crawl policy first | Candidate direction "Provider version directory" |
| F10 Push data into the code repository daily | not supplied | n/a | unsupported; no work | module zip capped at 500 MiB; consumers pin exact revisions; tracked files are 13.7 MB | Rationale recorded in the directory candidate |
| F11 Identify APItools on downloads (User-Agent) | not supplied | Lower | deferred | no `User-Agent` set anywhere | Existing candidate, trigger amended |
| F12 Directory record compatible with APIs.guru's `list.json`, with `preferred` set to the verified latest | not supplied | requested design | design constraint | `Client` already accepts a configurable APIs.guru list URL, and `preferred` drives search's version choice | S05.1 shape; the directory candidate |

## Consistency intake findings (2026-10-06)

| Finding | Source priority | Local severity | Disposition | Evidence | Owner/action |
|---|---|---|---|---|---|
| F01 Concurrent tiers, one-request 304 acceptance, and eight hint fetches conflict with shared caps | not supplied | P2 | confirmed planning gap; user selected shared caps | Earlier S05 latency contract and timing fixtures require all three at once | S05.3 owns shared accounting; S05.2/S05.4 integrate it; S05.1 records the contract |
| F02 Required diffs lack guaranteed baseline inventory | not supplied | P2 | confirmed planning gap; user selected optional inventory | Earlier S05.1 request permits bytes/registration while S05.5 prohibits reparsing and requires a diff | S05.5 uses supplied inventory only, otherwise comparison is unexamined; S05.1 records this |
| F03 Verified latest contradicts bounded discovery | not supplied | P2 | confirmed planning gap; user selected scoped preference | Earlier S05 locator contract and directory candidate conflict with never claiming latest | S05.3 implements scoped preference; S05.1 records highest verified comparable version in checked scope |
| F04 Over-cap reads cannot establish a full digest | not supplied | P2 | confirmed planning gap | Earlier S05 requires length and digest over its 8 MiB cap; `download.go:readBoundedResponseBody` refuses oversized content | S05.3 omits verified digest on incomplete reads and distinguishes observed/declared byte evidence |
| F05 Save guarantees exceed helper/validator behavior and confuse persistence permissions | not supplied | P2 | confirmed planning gap | `import.go:writeUniqueFile` provides naming/idempotence, not directory confinement; `validation.go` accepts OpenAPI/Swagger, not Discovery | S05.3 owns confined eligible saves; S05.2 separates metadata persistence; S05.1 documents permissions |
| F06 Guarded-client reuse alone is not origin/credential isolation | not supplied | P1 | confirmed acceptance gap, not a demonstrated production exploit | `download.go:redirectSafeClient` retains cookies/callbacks and accepts other safe public hosts; `catalog_discovery_remote.go` strips them for M80 | S05.2 owns isolated client and every-hop official-scope guard; S05.3/S05.4 reuse it |
| F07 Prefix sniffing may reject valid late-marked documents | not supplied | P2 | confirmed planning gap | Earlier S05.2 rejects non-spec first bytes without distinguishing ambiguous JSON/YAML prefixes | S05.2 rejects only definite non-spec responses early; S05.3 fixtures cover late identifying keys |

F08-F10 are owned by [S06](../docs/history/status-S06.md).
S05.6 reuses F08's isolated consumer checks and the F09 dependency correction;
neither creates another owner row. All work fits untouched pending rows, with
no new milestone or promoted candidate.

| Item | State | Notes |
|---|---|---|
| S05.1 — Design record and checkpoint | `[+]` | Owns the original design checkpoint for P2-1/P2-2, P2-3/P2-5/P2-6/P2-7 decisions, L1/L4, and third-intake F02/F04/F05/F06/F07/F08/F12. Write `docs/api-version-discovery.md` for proposed `apitools.api-version-discovery/v1`: known baseline plus optional already-built inventory; five-second shared check, three-second requests, four concurrent requests, twelve requests, three document bodies and 8 MiB per checked document; baseline-independent availability and partial reports; tier order; `newer_found`/`none_found_in_scope`/`unexamined`/`conflicting` with `checked_at`; catalog-preferred/catalog-highest/official-candidate/official-verified evidence. Record all consistency-intake F01-F07 decisions in the specification, including origin/repository scope, no authority from shared hosting alone, scoped `preferred`, complete-read-only digest, separate metadata/document write permissions, and family eligibility. Retain locator kinds (`direct`, `github`, `pointer`, `pattern`, `template`, `gated`), recipe/replay order, resolved URL, digest when verified, `checked_at`, `ok`/`moved`/`dead`/`gated`, list-compatible `x-apitools-*` fields, baseline-age/stale-after-12-months evidence, natural ordering and version conflicts, bounded standard-library docs scanner, and adapter contract. Consumer notes keep baseline use independent of network/adapter work, and keep OpenUdon's `catalog import-openapi` and `step source add` ownership unchanged. Treat 2026-10-02 measurements as dated evidence; verify whether a provider-scoped APIs.guru endpoint exists before relying on it. Adjust starting values only with measured evidence and user approval. Complete only after separate user approval of the detailed design record; later rows cannot start first. |
| S05.2 — Probe primitive, catalog evidence, and freshness state | `[+]` | Owns P2-3/P2-4/P2-6, L3's probe part, third-intake F02/F03/F04/F06 state, and consistency-intake F06/F07. Depends on approved S05.1. Add an unexported probe using existing URL/redirect/dial guards without changing `downloadBounded` or existing callers. Isolate cookies/caller redirect callbacks and enforce declared official origin/repository scope on every request and redirect, including replay. Support conditional GET, per-request caps, tool identification, and early rejection only for definite non-spec responses; ambiguous JSON/YAML prefixes continue bounded parsing. Treat 404 as a miss, 429/5xx as unexamined without retry, recording `Retry-After`. Tier 0 reads catalog and caller-supplied cached APIs.guru list locally, parses `x-origin` and `updated`, orders versions naturally and labels catalog preferred separately. After 24 h, optional network revalidation consumes shared budget; use a provider-scoped endpoint only if verified. Add optional versioned JSON state/previous-report reuse with a 24 h TTL for none-found and locator/recipe/resolved URL/digest/status/time evidence. State and list-cache persistence require separate explicit paths/permission, not `--save-dir`; cached records cannot expand current origin policy. Preserve search output and SQLite schema. |
| S05.3 — Deterministic upgrade check | `[+]` | Owns P2-1/P2-2/P2-4/P2-7 implementation, L1/L2 fixtures, third-intake F02/F06/F07, and consistency-intake F01/F03/F04/F05. Depends on S05.2. Run tiers 1-3 concurrently under one shared deadline and request/body limits; charge revalidation, conditional requests, redirects, replay and integrated hint checks without resets. A fetched unchanged or invalid source body consumes a document slot. Probe upward from integer/dotted tokens, stop after two consecutive misses, prefer known comparable newer candidates, and preserve deterministic report order. Record source/final URL, observed versus declared bytes, time, validation status, and verified full digest only after a complete bounded read. Timeouts/rate limits produce partial reports with the baseline intact. Scan one supplied official docs page within size/link caps and report no static links honestly. Record locator/recipe and replay then re-read pointers before hints. Compute scoped highest verified comparable preference; conflicts do not create one. Save only fully fetched newer Import-valid OpenAPI/Swagger bytes into an explicit confined directory, preserving naming/idempotence without relying on the existing helper for confinement. Never save baseline, invalid, partial, oversize or merely parsed Discovery content; never overwrite differing content. No save directory means no source writes, independently of metadata persistence. Unsafe paths and failed saves remain visible warnings, never failed checks; do not extend the budget. Add all specification timing, trust and save fixtures. |
| S05.4 — Hint verification and optional adapter | `[+]` | Owns P2-5 and integrates consistency-intake F01/F06. Depends on S05.3. Add APItools-owned `VerifyHints` and versioned size-bounded JSON hints (at most eight); deduplicate and order deterministically, fetch budget-selected URLs at most once, and explicitly mark excess unexamined. A standalone call uses the same caps; an integrated call uses remaining budget, never a reset. Reuse every-hop origin/credential isolation; accept no supplied spec bytes, credentials, cookies or authority upgrades. Keep provenance and locator/recipe for verified hints (third-intake F06). The optional pull wrapper runs only after deterministic tiers leave the question open and with caller opt-in/deadline (15 s default, 60 s ceiling); adapter failure, timeout or garbage is a tier status and never changes a verified result. Test with a fake adapter; add no LLM/browser dependency. |
| S05.5 — Capability-relevant selection | `[+]` | Owns P2-2's diff and consistency-intake F02. Depends on S05.4. Diff fetched newer versions only against an optional caller-supplied already-built baseline inventory; never fetch/reparse the baseline for comparison. Without it, comparison is unexamined and the version check remains valid. Add optional `StepContract` ranking through existing operation-candidate metadata, report capability coverage and conflicting version evidence, and never replace the baseline's operations by default. |
| S05.6 — Command, documentation, and verification | `[+]` | Owns L2 verification and L4 docs. Depends on S05.1-S05.5. Add opt-in `apitools versions` with `--timeout`, `--no-network`, `--state`, and `--save-dir`; document delivered behavior in README/architecture. Run all specification fixtures and verification below, including race tests, and reuse S06.3's candidate-bound workspace/standalone consumer checks with unchanged-pin checks labeled separately (consistency-intake F08/F09). Pass the persisted ten-iteration gate with no open P1/P2-or-higher findings. Publication, merging into the workspace checkout and consumer adoption each require separate approval. |

## Dependencies and ownership

S06 completed review and retirement on 2026-10-06. The remaining active order
is S05 under the same execution owner; S05's task order is S05.1 -> S05.2 -> S05.3 -> S05.4 -> S05.5 -> S05.6.
S05 has no technical dependency on S06; both execute serially on the same
isolated branch/worktree, carrying S06's changes into S05 after its review.
No upstream prerequisite applies, and no sibling write scope is approved.

Kinet's Stage 5 is complete, not a pending prerequisite or merge gate. Retain
one separate worktree outside `~/Workspace/go.work`; the main workspace
checkout and normal consumer builds stay unchanged. Merging and publishing
each still require separate approval. S05.1 needs its own design-record
approval before S05.2 begins. Preserve the
built-in catalog data, artifact registrations, catalog index inputs, and every
consumer pin. OpenUdon or Kinet adoption is reconciled separately in its own
ledger; S05.1's consumer notes describe the background-advisory pattern for
them, and APItools edits no sibling file.

## Acceptance and verification

Every row's outcome holds. The latency and cost contract in the specification
is met with the approved starting values, or with values S05.1 adjusts on
measured evidence and the user approves. Provider-free fixtures with local
servers and simulated latency cover sibling versions, equal `info.version`
values, a strictly invalid version, a newer version missing operations,
off-origin redirects, unsafe hosts, adapter hints, and these timing cases:

- A stalled host returns within the budget plus a small slack, with the
  baseline intact and the affected tiers `unexamined`.
- All-404 siblings stop after the gap rule, with no more than 12 requests.
- A fresh state costs zero requests. An isolated conditional-check fixture
  returning 304 costs exactly one; concurrent tiers may issue other requests.
  A 304 proves only that resource unchanged, not absence of newer versions.
- A 200 HTML response larger than 1 MiB to a probe is abandoned after a few
  KiB; valid JSON/YAML with identifying keys beyond the sniff prefix is not
  rejected solely for that ordering.
- A 429 or 5xx response is `unexamined`, with no retry inside the budget.
- An adapter timeout leaves the deterministic result unchanged.
- Report order is identical across permuted completion order, and
  cancellation leaves no goroutines.
- A saved file's digest equals the reported digest and a rerun reuses it; an
  invalid, partial, over-cap, baseline, or merely parsed Discovery document is
  never written. Without a save directory no source bytes are written; metadata
  state/cache writes need their own explicit opt-in. Test confined paths,
  symlinked directories/targets, differing collisions, and failed saves as
  warnings that do not fail the check.
- A `moved` locator is repaired by replaying its recipe before any LLM tier,
  and isolated cached-list revalidation with a 304 consumes one shared request.
- Aggregate accounting includes conditional fetches, cache revalidation,
  redirects, replay and integrated hints; invalid or unchanged complete source
  bodies consume document slots. Eight hints may exceed three-document capacity;
  verify deterministic selection and explicit unexamined excess, with no reset.
- Missing baseline inventory yields an unexamined comparison without fetching
  or reparsing it; supplied inventory enables the diff.
- Scoped preference never claims global latest; incomparable/conflicting
  versions cannot establish a new preference.
- Oversized/partial reads expose observed versus declared byte evidence but no
  verified full-content digest.
- Cookies and redirect callbacks cannot supply credentials; every-hop origin
  enforcement rejects off-origin redirects and shared-host repository escapes,
  including recipes and hints.

A run without opt-in performs no network access, and existing exported APIs,
wire shapes, and `Client.Search` results are unchanged.

Required verification:

```bash
go test ./...
GOWORK=off go test ./...
go test -race ./...
go vet ./...
GOWORK=off go vet ./...
go run ./cmd/cataloggen -check
go run ./cmd/apitools catalog check
go run ./cmd/apitools search --help
go run ./cmd/apitools versions --help
git diff --check
```

Run OpenUdon and Udon compatibility checks in workspace and standalone modes
per tech-stack.md, with Ramen excluded. Follow
[S06's isolated consumer verification](../docs/history/status-S06.md):
disposable workspace/modfile configuration binds both consumers to the
candidate worktree; check the resolved APItools directory before test/vet.
Do not edit sibling manifests/pins or the shared workspace. Ordinary unchanged-
pin/replacement standalone runs are labeled baseline evidence, not candidate
qualification. Default checks stay provider-free and offline; network fixtures
use only local servers.

## Review and planning evidence

**Review iterations started:** 2 of at most 10; iteration 2 passed 2026-10-06.
No intake is a bounded-gate pass. Persist the count once implementation
reaches review and never reset it. An unresolved blocking finding at iteration
10 requires user direction.

Planning used read-only inspection of the code, catalog data, and consumer
ledgers named above, an offline parse of the local GoDaddy copies with the
cached OpenAPI 3.1 parser, a revalidation of every code-level review claim at
the commit above, and, for the third intake, read-only requests for the public
APIs.guru and public-apis data named above. No implementation verification or
milestone review is claimed.

Fourth-intake verification (2026-10-06): the four existing public-apis tests
named in [S06](../docs/history/status-S06.md) passed in workspace
and standalone modes; `git diff --check` passed. These validate the unchanged
baseline, not the planned features. All nine pending rows, both zero review
counters, retired history, code and current-truth documents are preserved.
The approved six-file amendment clarifies v26 rather than creating a new
direction version; no commit, execution loop, merge or publication is authorized
by this planning action.

## Reconciliation after S06 completion (2026-10-06)

S06's bounded Markdown/legacy JSON list and repaired default are accepted in
the isolated branch, with full APItools and candidate-bound OpenUdon/Udon
checks and review 1 passed. S05's six rows remain pending: its new conditional/
concurrent probe and official-version metadata are not implemented by S06.
Reuse S06's disposable consumer-test setup and verify candidate module identity
before testing; recorded Udon retries also explain the lower-concurrency
verification choice. The detailed S05.1 design still needs separate approval
before S05.2. No consumer pin, source registration, catalog index or SQLite
schema changes; no merge or publication performed.

## Execution checkpoint

The confirmed goal completed and retired S06 in task commits
`ae7f393`, `c379ffb`, and `828dd85`. S05.1 now has the concrete
[design record](../../docs/api-version-discovery.md) ready for separate approval;
it is not complete and no S05 feature code has started. All five later rows
are pending and the S05 review counter stays zero.

The record keeps the approved five-second/shared limits, optional prebuilt
baseline inventory, scoped preference, separate persistence permissions and
official-scope/credential isolation. It fixes wire/API names, tier scheduling,
locator replay, bounded hints, failure semantics, save eligibility and CLI
opt-in. Its source-body cap explicitly applies to API-description bodies;
lookup metadata has separate bounded reads under the same HTTP/deadline cap.
S05.2 cannot begin until the user approves the record.

S05.1 verified the documented provider-scoped APIs.guru contract through the
official repository's 2.2.0 OpenAPI definition: `GET /{provider}.json` has the
same APIs map shape as list.json. Live probes of the provider endpoints
returned HTTP 403 here; this records unexamined availability rather than a
nonexistent endpoint or empty directory. Exact provider keys may use the
documented scoped URL; local cached evidence is preserved on failure.
Earlier 2026-10-02 size/age/GoDaddy measurements remain dated observations,
not newly measured facts.

S06's final consumer artifact cleanup completed after its closure commit.
Two additional synthetic output directories were found and moved to disposable
verification storage; the [knowledge correction](../docs/history/knowledge.md#2026-10-06--s06-consumer-artifact-cleanup-timing-correction)
preserves the timing correction without rewriting the frozen S06 record.
Both consumer Git status checks are now clean; source/manifests/pins stayed
unchanged. No merge or publication has occurred.

The user approved S05's detailed design record on 2026-10-06. S05.1 is
complete after structural/link/patch verification; the separate design gate
is satisfied. All implementation rows now follow the approved record, and
no launch confirmation is required again for this continuing goal.

## Implementation notes

S05.2 completed 2026-10-06: additive request/report/hint/state records, strict
bounded decoders, natural version ordering, catalog/x-origin evidence helpers,
confined JSON state/list-cache I/O, and a private conditional/sniff probe with
shared counters, every-hop scope checks and cookie/callback isolation.
Focused local-server, state/symlink/repository and directory tests passed;
probe race fixtures passed. Existing downloads/search signatures are intact.

S05.3 completed 2026-10-06: baseline-first discovery joins local evidence,
conditional sources, sequential two-miss sibling sequences, bounded publisher
pointers/GitHub recipes, scoped directory lookup and opt-in validated saves.
Shared context/counters limit concurrent tiers; invalid/oversized sources and
fetch/persistence failures remain visible. A 304 alone is not negative version
evidence. Fresh matched state uses zero requests. Local fixtures for deny,
stall, conditional, sibling/save/resume, docs/invalid source and list 304 passed;
the combined foundations/discovery race run passed. Existing Search is unchanged.

S05.4 completed 2026-10-06: standalone VerifyHints accepts eight bounded
untrusted hints, deduplicates them and verifies only budget-selected sources.
Excess demand remains explicit unexamined records. Producer/claim text cannot
upgrade authority or change a source-verified version. The optional cooperative
pull adapter retains the original network deadline/counters and isolates its
copies; failure/late hints preserve baseline/deterministic evidence. Budget,
deduplication, dishonest-claim, timeout and callback-isolation race tests passed.

S05.5 completed 2026-10-06: digest-matched, complete caller inventory enables
method/path diffs and operationId-change evidence. Missing/mismatched/partial
metadata remains unexamined, and the baseline bytes are never reopened.
Optional StepContract ranking reuses source-native operation candidates over
newly fetched bytes with exact digest, auth/effect and truncation evidence.
Focused diff, missing-inventory, caller-preservation and ranking tests passed;
the combined version-check race gate passed. No version replaces baseline ops.

S05.6 implementation and verification completed 2026-10-06: versions CLI
uses the approved explicit network/persistence flags and 0/2/1 policy. README,
product, architecture and stack describe delivered metadata behavior.
Full APItools tests passed in workspace/standalone modes; full race tests,
both vet modes, catalog generator/quality and search/versions help passed.
OpenUdon/Udon full candidate-bound tests/vet passed in both modes, with Ramen
excluded. Synthetic Udon output directories were moved into disposable
verification storage; consumer source/manifests/pins are unchanged.

## Whole-milestone review

Iteration 1 started after verification. Review the complete S05 changes from
S06 closure through all five task commits and the S05.6 diff: request and
source identity, budgets/deadlines, metadata scope/negative evidence, shared
host/redirect/credential protections, state reuse, recipes and hints,
comparison independence, confined saves, CLI/errors and consumer compatibility.
The GitHub recipe fixture initially failed and its provenance loss was fixed;
this is retained verification evidence. No closing outcome is claimed yet.

Iteration 1 findings: R1 (P1) GitHub repository authorization used a hostname
substring rather than the exact raw-content origin; R2 (P2) a supplied native
Discovery inventory still could not produce the approved native operation diff.
Both remain owned by S05.6. Fix and verify them, then review the whole milestone
again as iteration 2; do not reset the persisted count or close at iteration 1.

R1 and R2 are fixed with exact-host and source-native Discovery selector
regressions; the focused combined race gate passed. Iteration 2 now reviews the
whole milestone, including these fixes, under the same persisted ten-pass cap.

## Closure evidence

Iteration 2 passed with no remaining P1/P2-or-higher findings and no carried
lower findings. The whole review covered source-family/service identity,
unknown/negative evidence, strict request and hint inputs, natural ordering,
exact GitHub authority and path scopes, redirect/request/body bounds, deadline
and worker ownership, conditional/cache semantics, prefix and size handling,
recipe provenance, adapter isolation, native diffs/candidate evidence,
byte-identical confined saves/collisions, CLI policies and all metadata writes.

Required final tests passed in workspace and standalone modes with Go 1.26.6;
full `go test -race -p 2 ./...`, both vet modes, catalog generator/quality
(zero findings), search/versions help and patch checks passed. OpenUdon/Udon
full test/vet against the candidate worktree passed in both modes. Consumer
source/manifests/pins remained unchanged; observed synthetic test artifacts
were moved to disposable storage, and Ramen was excluded.

Consolidation: README and docs/api-version-discovery.md describe delivered
interfaces; product/architecture/stack own current facts. Existing lessons
about incomplete coverage remain applicable. The exact-host/native-identity
lesson is recorded with this milestone. No prior product knowledge or frozen
record was removed. Evolution v26 is delivered; no new direction bump is
needed. Candidate directions remain unnumbered: S05 delivery satisfies the
directory candidate's prerequisite, but legal/operator/crawl-policy gates
still require separate approval. No candidate is promoted automatically.

Downstream reconciliation: no consumer adopts these new interfaces in this
milestone. Existing imports/pins and source registrations/catalog index/SQLite
schema remain intact; compatibility is verified. Retire the full specification
and status, refresh the history index and remove the exhausted launch input.
The confirmed S06 -> S05 implementation goal is then complete in this branch.
````
