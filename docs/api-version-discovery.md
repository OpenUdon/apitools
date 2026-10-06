# API version discovery design

State: implemented and verified on 2026-10-06 in isolated branch
`work/s06-s05`, following the approved S05.1 design. The complete task and
two-iteration review evidence is retained in [S05](../tabilet/docs/history/status-S05.md).

## Outcome and boundaries

Check whether later API versions can be found within an explicit publisher
scope while keeping the version already held available throughout the check.
The result is advisory metadata. It neither replaces the held source nor binds,
approves or executes an operation. Existing Search, Import, catalog data,
registrations, operation indexes, SQLite schemas and consumer pins keep their
contracts. There is no LLM client, browser integration or general crawler here.

Source families for the initial check are OpenAPI/Swagger and Google Discovery.
Other families remain catalog leads with an unsupported-family diagnostic.
Discovery can be verified through its native parser; it does not become an
OpenAPI Import input. No live introspection, account lookup, tokens, cookies,
signed requests or provider operations are used.

## Additive library and wire contract

The root package adds the following interfaces; existing exported types and
JSON shapes are unchanged:

```go
func (c *Client) DiscoverAPIVersions(
    ctx context.Context,
    request APIVersionDiscoveryRequest,
    options APIVersionDiscoveryOptions,
) (APIVersionDiscoveryReport, error)

func (c *Client) VerifyHints(
    ctx context.Context,
    request APIVersionDiscoveryRequest,
    hints APIVersionHints,
    options APIVersionDiscoveryOptions,
) (APIVersionDiscoveryReport, error)

type APIVersionHintAdapter interface {
    DiscoverHints(context.Context, APIVersionHintContext) (APIVersionHints, error)
}

func (c *Client) DiscoverAPIVersionsWithAdapter(
    ctx context.Context,
    request APIVersionDiscoveryRequest,
    options APIVersionDiscoveryOptions,
    adapter APIVersionHintAdapter,
) (APIVersionDiscoveryReport, error)
```

Request/report schema is `apitools.api-version-discovery/v1`. Hint schema is
`apitools.api-version-hints/v1`; persisted state uses
`apitools.api-version-state/v1`. Strict request/hint decoders reject unknown
fields, unsupported schemas and trailing values. Input JSON is at most 64 KiB;
the hint array has at most eight entries. Existing prompt string, field and
collection budgets apply, with diagnostics for any shortening.

| Record | Required meaning |
|---|---|
| Request | Schema, known baseline, explicit official scopes, optional one docs pointer, locator recipes and optional StepContract. No transport, endpoint override or installation path in wire JSON. |
| Known baseline | Source family, version label, SHA-256, provider key or public source URL, and optional catalog updated timestamp. This metadata initializes every valid report before network work. |
| Options | Explicit Network permission (false by default), timeout, optional metadata state/list-cache paths, SaveDir, optional already-built baseline inventory/previous report, installation catalog/client settings and optional adapter deadline. These are outside request JSON. |
| Report | Schema, unchanged known baseline, status, checked_at, checked scope, ordered version records, tier outcomes, diagnostics, truncation/evidence gaps and optional preferred. No source bytes. |
| Version | Locator/recipe, evidence label, observed version tokens, source/final public URL, complete-read digest if available, observed/declared bytes, fetch time, validation result, comparison/diff, optional saved-file evidence and locator health. |
| Hint | Candidate public URL, optional version claim, bounded evidence note, producer label and optional supported locator recipe. No spec bytes, credentials, headers or asserted authority. |

The optional already-built inventory is an existing OperationInventory supplied
through Go options or a separate CLI metadata file. It is scoped by a matching
baseline SHA-256. APItools copies and bounds it; it never reads the baseline
spec bytes to reconstruct it. A mismatched, truncated or unusable inventory
makes comparison unexamined without rejecting the version check. Metadata
inventory input is limited to 32 MiB and 10,000 operations, with existing
structural/prompt budgets. Google Discovery comparisons are unexamined unless
the caller supplies source-native identity metadata; an OpenAPI-shaped display
projection alone is insufficient for a native comparison.

Core wire fields use the following JSON spelling; the example is synthetic:

```json
{
  "schema_version": "apitools.api-version-discovery/v1",
  "known": {
    "source_kind": "openapi",
    "provider_key": "example.com:mail",
    "source_url": "https://api.example.com/openapi/v1.json",
    "version": "v1",
    "sha256": "0000000000000000000000000000000000000000000000000000000000000000"
  },
  "official_scopes": [
    {"origin": "https://api.example.com", "path_prefix": "/openapi/"}
  ],
  "locators": [
    {"kind": "pattern", "url_template": "https://api.example.com/openapi/v{version}.json", "baseline_token": "1"}
  ]
}
```

Optional request fields are docs_url and contract. Official scopes use origin,
path_prefix and optional repository identity; at most eight scopes/locators are
accepted. Baseline catalog_updated is an optional RFC3339 timestamp. The held
digest is caller-supplied baseline identity, not a new APItools verification of
source bytes. Previous reports are supplied in Go options or versioned state,
outside the 64-KiB request, with the same 512-KiB report cap.

Version records use id, locator, recipe, evidence, version_claims, source_url,
final_url, sha256 (optional), bytes_observed, bytes_declared (optional),
fetched_at, validation, comparison, diff (optional), saved (optional),
checked_at and locator_status. Version claims retain url_token, info_version
and docs_version independently. Tier records use tier, status, request_count,
source_body_count, checked_at and diagnostics. All times are RFC3339 UTC.
Stable IDs derive from canonical locator/URL identity, never completion order;
conflicting same-version documents remain separate records rather than map
overwrites. The top-level report uses schema_version, known, status,
checked_at, checked_scope, versions, tiers, diagnostics, truncated and optional
preferred/directory_records. Locator health and validation status remain
separate from the top-level freshness status.

The Go options fields are Network, Timeout, StatePath, ListCachePath, SaveDir,
BaselineInventory, BaselineInventorySHA256, Previous, Catalog,
AdapterEnabled and AdapterTimeout. Client transport/installation endpoints
stay on the existing Client, copied under this design's safety policy.
Runtime option data is never serialized as request or report authority.

Invalid request/schema/limits/scope/credential-bearing URLs return an error
before network or writes. Network refusal, timeout, cancellation, HTTP failure,
parse failure and persistence failure return a report preserving the baseline
and visible diagnostics. They cannot establish negative evidence. The caller
can use its held baseline immediately; the synchronous method returns the
bounded check report, not an instantaneous network result.

## Status, version identity and preference

Top-level statuses are `newer_found`, `none_found_in_scope`, `unexamined` and
`conflicting`, with checked_at. Conflicting identity evidence takes precedence;
otherwise verified comparable newer evidence yields newer_found even when
other scope remains unexamined. Without positives, incomplete scope is
unexamined; none_found_in_scope requires complete examination of the declared
finite scope. Unattempted probes, unsupported sources and partial summaries
remain explicit. No outcome claims the publisher's globally latest version.

Evidence labels are catalog-preferred, catalog-highest, official-candidate and
official-verified. Catalog preferred and highest are distinct. Original URLs
from APIs.guru are leads; they gain official-verified only after publisher
scope checks and complete native validation. Invalid fetched documents retain
their locator, complete-read digest and validation diagnostics; they are not
verified versions or save candidates.

Compare integer/dotted version tokens naturally, with an optional leading v
and trailing zero components normalized for comparison. Generate sibling
URLs only from such tokens, incrementing the rightmost numeric component.
Recognizable year/month/day date tokens are not generated even when dotted.
Date, opaque or prerelease labels may be listed
from evidence but are not guessed or silently ordered against numeric labels.
Retain URL tokens, info.version and docs claims separately. Conflicts between
them remain explicit; dates or version strings are not substituted for one
another. Equal info.version across distinct URL versions is a conflict rather
than evidence that the sources are identical.

Each provider/service record uses the APIs.guru-compatible versions map and
optional preferred, plus x-apitools-* fields for scope, locator, recipe,
verification and age. Preferred identifies the highest verified comparable
version found in that recorded scope; conflicting/incomparable evidence cannot
create a new preference. A future published directory remains a separate
legally gated candidate, and this milestone publishes no directory or specs.
Use swaggerUrl/swaggerYamlUrl only for real OpenAPI/Swagger documents. Native
Discovery locators belong in x-apitools-* extensions and x-origin evidence;
they never masquerade as OpenAPI or trigger conversion. The full report keeps
native records even when a legacy directory reader cannot consume them.

## Shared cost and deadline

Approved defaults remain five seconds for a deterministic invocation, three
seconds per request, four concurrent requests, twelve HTTP requests and three
API-source body slots, with 8 MiB per checked API document. The caller may
choose a smaller timeout; increasing a default needs the design approval
specified by S05.1. The old 30-second Client timeout is not inherited.

The request limit includes redirects and conditional/revalidation/replay/hint
requests. Reserve a source-body slot before scheduling an API-source GET;
304/404 or definite non-spec sniff rejection releases the slot, but a fully
read unchanged or invalid source consumes it. Oversized/partial reads are
unexamined and cannot carry a verified full-content digest. Observed bytes and
Content-Length are separate evidence. Count source slots globally across
integrated tiers; never grant each tier another three documents.

Lookup metadata is separately size-bounded under the same request/deadline
limits: the APIs.guru list/provider map at 20 MiB decoded, one docs page at
1 MiB, one RFC 9727 catalog at 1 MiB/100 links, and a GitHub tree response at
1 MiB/10,000 entries. The three source slots refer to API-description bodies,
not these lookup maps/pages. Wire and decoded limits apply independently.
Documentation yields at most 32 spec-like links; excess is unexamined rather
than silently complete evidence. Report projection remains within the shared
512-KiB prompt context bound; omissions retain named diagnostics/coverage gaps.

All workers share context, counters and a deterministic task order. Cancel and
join owned workers before returning. Parsing and regular-file I/O are bounded
and cooperatively check the deadline; this is not a hard kernel/process CPU
or filesystem latency guarantee. Do not begin a save after the deadline.
No background APItools worker may outlive an invocation.

## Tiers and scheduling

| Tier | Work and stop behavior |
|---|---|
| 0 | Local baseline, fresh previous/state evidence (24-hour TTL), catalog refs and a caller-supplied cached APIs.guru list. No network. Cached evidence is reusable only for matching baseline, official scope and comparison inputs. |
| 1 | Conditional GET of the known official source using ETag/Last-Modified. A 304 proves only that resource unchanged. A changed complete digest records an in-place revision without replacing the held baseline. |
| 2 | Upward integer/dotted sibling probes within the approved origin/path scope. Two consecutive 404 misses end each generated sequence; 429/5xx are unexamined and are not retried. |
| 3 | One RFC 9727 pointer, one supplied official docs page, or supported locator replay. One tree listing per declared GitHub repository; truncated trees remain incomplete. |
| 4 | Explicit hints after deterministic work remains open. Integrated verification uses remaining budget; a later standalone VerifyHints call is a new explicitly requested bounded invocation. |

After tier 0, tiers 1-3 run together rather than receiving additive deadlines.
Each generated sibling sequence is probed in order so two consecutive misses
stop it before more siblings are scheduled; concurrency spans independent
tiers, locators and already-known candidates.
Any stale directory-cache revalidation is network work charged to that same
budget. Sort known comparable newer candidates newest-first before scheduling;
when new links arrive, maintain the same canonical priority among ready tasks.
Unknown future links cannot retroactively change an already-started request.
Sort final reports by comparable version and canonical locator/URL/digest, not
completion time. Preserve findings when a sibling tier fails.

The private probe primitive reuses host, port, redirect, DNS/dial and gzip
guards. It adds conditional headers, response headers, a stable
`apitools/api-version-discovery-v1` User-Agent, and a 4-KiB sniff prefix.
Reject definite HTML/non-spec bodies early. Ambiguous JSON/YAML continues
bounded parsing; source markers may occur later. 404 is a miss; 401/403 marks
gated/unexamined evidence; 429/5xx records Retry-After and remains unexamined
without retry. Existing downloadBounded and its callers are unchanged.

## Official scope and locator replay

Allowed evidence is a caller-declared official scope, a catalog official-*
reference, or an APIs.guru x-origin lead. A scope includes scheme/host/port
and, for shared hosts, repository/path identity. Declaring a shared hostname
does not authorize every publisher on it. Normalize paths and reject path
escape, userinfo, credential-like query parameters and HTTPS downgrade.
Every initial/redirect/replay/hint URL must pass both transport safety and
official scope policy. Reports redact query/userinfo/fragment; a URL requiring
secret-bearing parameters is gated and not fetched.

Owned clients drop cookie jars, source caches and caller redirect callbacks.
Use the existing default-safe proxy-free transport and DNS rebinding guards;
unsafe-host overrides are test-only. GitHub locators explicitly identify owner,
repository, ref and allowed path pattern. The fixed GitHub tree API and raw
content origins are paired only for that declared repository; no repository
search, authenticated calls or unrestricted cross-origin expansion occurs.

| Locator kind | Retained recipe |
|---|---|
| direct | Fixed public source URL and validators. |
| github | Owner/repository/ref, bounded path pattern and resolved tree/ref identity. |
| pointer | One supplied docs page, publisher catalog or latest alias plus selected link. |
| pattern | Known URL with one integer/dotted version slot and bounded upward sequence. |
| template | Tenant-specific guidance retained without interpolation or generalized verification. |
| gated | Auth-required or human-only source, with a docs-overlay follow-up suggestion. |

Each version stores last resolved URL, checked_at, verified digest when
available and locator health ok/moved/dead/gated. A broken locator first
replays its recipe, then rereads its explicit pointer; only remaining open
questions may reach the opt-in hint tier. Both actions consume remaining
budget. A gated/human-only source does not cause automatic overlay generation.

## Hints and adapter failure isolation

VerifyHints accepts at most eight bounded hints, validates and deduplicates
them, and prioritizes comparable newer candidates followed by canonical URLs.
Fetch each selected URL at most once. Hints beyond remaining capacity receive
unexamined/budget-exhausted outcomes. Adapter provenance and notes do not
upgrade authority. Invalid individual hints produce diagnostics; invalid
envelopes are request errors. No supplied bytes/headers/auth are accepted.

Push-first integration is the default: consumers run their LLM/browser work
independently and submit hints when ready. The optional pull wrapper requires
explicit AdapterEnabled and a context-cooperative adapter, with a 15-second
default/60-second ceiling on its caller deadline. It starts only after the
deterministic result remains open. Its wait cannot replace the five-second
network budget or reset counters: late hints remain unexamined, and consumers
can later request standalone verification. Adapter timeout/error/garbage is a
tier outcome, preserves deterministic evidence and never fails the baseline.
APItools adds no adapter implementation or goroutine that forcibly abandons
an uncooperative foreign adapter.

## Metadata persistence, saves and comparison

State and directory-cache writes each require separate explicit paths outside
the repository. SaveDir permits source-document writes only. With none of
these paths, no files are written. Use versioned bounded JSON, atomic confined
regular-file writes and no symlink following; do not migrate the catalog DB.
Never persist cookies, credentials, request headers or source bytes in state.
The state records baseline/scope/contract identity, check age, locator evidence
and coverage so a fresh but different query cannot inherit negative evidence.

Only fully fetched newer OpenAPI/Swagger content accepted by existing Import
validation is eligible for SaveDir, byte-for-byte with matching SHA-256.
Reuse Import naming and digest-idempotence behind filesystem confinement.
Identical targets are reused, differing content receives collision-safe names,
and symlink/non-regular targets or unsafe directories are refused. Do not
refetch via Import to save bytes already held. Baselines, invalid/partial/
oversized sources and parsed Discovery documents are not saved. Failed or
deadline-skipped saves are per-version warnings; no registration/index update
is made. Saved path/digest/bytes are reported, and no default directory exists.

For OpenAPI/Swagger, diff operation identity (method/path plus operationId
evidence) against the supplied inventory, never the baseline source. Missing
inventory and incomplete summaries yield unexamined comparisons; absence
claims require complete inventories. Optional StepContract ranking reuses
BuildOperationCandidates over fetched bytes and retains its capability,
effect, auth-alternative and truncation diagnostics. A newer version lacking
the requested capability does not displace an older matching version.

## CLI and consumer integration

`apitools versions --request FILE` reads the versioned request JSON. Network
is off by default; `--network` explicitly enables the bounded check and
`--no-network` explicitly denies it. Combining both is a usage error.
Additional flags are `--timeout` (default 5s), `--state`, `--list-cache`,
`--save-dir`, `--baseline-inventory` and `--json`. Output uses the existing
help/usage/runtime 0/2/1 policy; a valid partial/unexamined check exits 0 with
its report. The CLI implements no pull adapter or LLM invocation.

OpenUdon retains catalog import-openapi and step source add; wrapper relocation
is its decision. Consumers keep using the held baseline, may run tiers 0-3
within the bounded check, and treat slow hint production as background advisory.
The result never gates authoring, source adoption, review or execution approval.
Adoption in OpenUdon/Kinet is reconciled separately in their ledgers.

## Evidence and verification checkpoint

The 2026-10-02 GoDaddy/list measurements in S05 remain dated evidence: distinct
URL versions share info.version, newer v3 drops transfer operations, v1/v2
strict parsing fails, and the old APIs.guru origin has no sibling version slot.
The measured frozen-list timestamps are not a permanent claim about future
directory updates; report each observed updated age and mark older-than-12-
month baselines stale by source, suggesting the opt-in hints without invoking
them. No third-party spec/list bytes become repository fixtures.

The [APIs.guru 2.2.0 definition](https://raw.githubusercontent.com/APIs-guru/openapi-directory/main/APIs/apis.guru/2.2.0/openapi.yaml)
documents GET /{provider}.json as provider-scoped metadata using the same APIs
map as list.json. Its schema establishes the endpoint contract. Read-only
probes of ote-godaddy.com.json and apis.guru.json on 2026-10-06 returned HTTP
403 here; live availability remains unexamined. Use the documented scoped
endpoint when an exact provider key is supplied; failure remains partial,
with existing local cached-list evidence preserved. Otherwise consult the
caller-supplied global list cache; never clone the repository or infer an
endpoint from a provider display name. If global list revalidation is explicitly
enabled, its conditional request consumes the same network budget.

The scoped URL is derived only for the default APIs.guru endpoint and an
exact provider key, with the service suffix kept as a result filter. A custom
Client.APIsGuruListURL remains a list/mirror endpoint; it does not authorize
guessing adjacent provider URLs. Its cached/conditional list follows the
existing configured endpoint under this design's budget.

Provider-free httptest fixtures must cover the complete S05 acceptance matrix:
shared accounting and cancellation; stalled hosts; zero-request fresh state;
isolated 304s; two-miss sibling stop; eight hints exceeding capacity; no retry
on 429/5xx; HTML sniff/late markers; off-origin/shared-repository escape;
cookie/callback isolation; missing inventory; version conflicts and capability
loss; over-cap digest omission; recipe repair; independent metadata/save opt-ins;
identical reuse/collision/symlink refusal; and late/failing adapters. Timing
fixtures allow 250ms scheduling slack for a stalled host. Test concurrent
paths with -race and use no live provider calls in default verification.

Run full test/vet in workspace and standalone modes, catalog generation/quality,
search and versions help, patch checks, and candidate-bound OpenUdon/Udon
compatibility per retired S06. Its disposable workspace/modfile procedure and
module-directory assertions remain applicable. Ramen stays excluded. No
budget or authority change is approved merely by writing this record.
