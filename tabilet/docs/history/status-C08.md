# Retired milestone C08 - Searchable catalog and hosted embeddings

**Milestone.** C08
**Outcome.** completed
**Retired.** 2026-11-25
**Source status.** tabilet/memory-bank/status-C08.md
**Source specification.** tabilet/memory-bank/milestone.md#c08--searchable-catalog-and-hosted-embeddings
**Evidence.** e5e5f53fab5e3e1ba904c3753e6c3011f4f635d0
**Worktree.** includes uncommitted changes
**Review.** passed
**Review iterations.** 2
**Verification.** Go 1.26.6 build, vet, ordinary full suite with both retained publisher directories, whole root/CLI races with standard fixtures, and all 26 exact publisher selector cases under a separate race run passed. The combined pinned race timed out in the inherited DigitalOcean promotion fixture. Catalog generation and quality, CLI help, built search, local SQLite and TLS embedding fixtures, and the exact OpenUdon consumer suite passed. The ordinary Udon pinned suite passed, and its optional composed workspace retained the offline module setup failure. Closure reuses the exact code gate evidence because its inputs are unchanged. [Verification receipts](../../../dist/c08-goal/).
**Consolidated into.** [Operator documentation](../../../docs/catalog-search-profiles.md), [product workflows](../../memory-bank/product.md#primary-workflows), [search architecture](../../memory-bank/architecture.md#reviewed-provider-search-c08), and [search evidence lesson](../../memory-bank/lessons.md#separate-operation-evidence-from-retrieval-relevance).

## Milestone specification

````markdown
### C08 — Searchable catalog and hosted embeddings

**Dependency and owner.** After accepted [C07](../docs/history/status-C07.md). APItools owns public provider
metadata, the SQLite search index, embedding adapter and CLI behavior.
[C08 ledger](status-C08.md).

Populate descriptions, normalized categories, keywords and aliases for all
reviewed providers. Separately describe the operations that our sources
cover. Return coverage, freshness, provenance and matched capabilities in
text and JSON. Combine SQLite full-text ranking with optional hosted
embeddings at a configured OpenAI-compatible HTTPS endpoint. Send reviewed
public metadata and explicitly submitted query text only. Cache vectors by
model and content, enforce the configured budget and fall back to lexical
search when offline, unconfigured, or over budget. Default automatic CLI
search resolves exact reviewed identities first, then capability matches.
Insufficient local relevance permits the existing public-directory chain
for non-core queries. Preserve explicit source selection, offline behavior
and the existing Go search contract.

````

## Status record

````markdown
# C08 — Searchable catalog and hosted embeddings

**Owner.** APItools. **Dependency.** Accepted [C07](../docs/history/status-C07.md). **Parallel-safe: no.**
See the [milestone specification](milestone.md#c08--searchable-catalog-and-hosted-embeddings).

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| C08.1 — Curate searchable metadata | `[+]` | Give each of the 334 reviewed providers a normalized category, aliases, useful keywords, a concise public description and an evidence-backed description of the supported subset. Stage description changes for review. Cover both docs-only entries and providers with overlapping services. Added 334 profiles, 50 categories, 41 additional aliases, and 1,352 selected operations for 164 providers. Astra xhigh passed re-review. |
| C08.2 — Add SQLite full-text provider retrieval | `[+]` | Build a reproducible index from accepted metadata and support updates to earlier records. Expose matched capabilities, coverage, provenance and freshness in text and JSON. Exact provider identities and aliases have priority. Reproducible FTS5, reviewed CLI output and read-only accepted publisher observations passed. Astra xhigh passed re-review. |
| C08.3 — Add configurable hosted embeddings | `[+]` | Combine keyword and semantic ranking using the configured HTTPS endpoint, model and budget. Limit payloads to approved public provider text and explicit queries. Reuse vectors for matching content/model and refresh changed inputs. Test budget exhaustion, missing configuration, endpoint failure and offline lexical behavior with a fixture service. Dedicated HTTPS transport and budget/caching/fallback behavior passed Astra xhigh re-review. |
| C08.4 — Integrate and verify automatic CLI search | `[+]` | Prefer reviewed results and run existing public discovery for non-core or weakly matched queries under its existing network controls. Preserve explicit source flags, provider scoping and Go search behavior. Exercise “send slack message”, “refund payment”, an exact alias and a non-core query; verify relevant results, coverage labels, fallback and offline operation. Default and explicit sources passed built CLI and loopback checks. Astra xhigh passed re-review. |

## Review

Review gate iterations: 2

Whole review passed after 2 iterations of the effective limit of 10. Use the [test policy](test-policy.md)
and [milestone review procedure](milestone.md#milestone-review-procedure).

C08.1 verification started against base `a4caf6f56c10e10f2fd43cee1c9810fa3bef0e86` and candidate diff SHA-256 `4203cdf099c63d1a78d9a5eba5ae43418a82464b86b54c816ba90d5d13ab6cb2`. The row gates are build, vet, affected catalog tests with retained source pins, generator and quality checks, and diff validation.

C08.1 corrective gates started after independent review. Candidate diff SHA-256 `cf41bd37e2279945c73583bb466dac905b589c00a49b8a4929cd4518704fd9a1`.

C08.1 retains the accepted C07 catalog bytes. The additive profile dataset
distinguishes selected coverage, 151 reference-only profiles, and 19
docs-only profiles. Local exact pins reproduce Slack message posting and
Stripe, Square, and Mercado Pago payment refunds. The first task review
found NVIDIA's omitted pinned route, overbroad Plaid/Etsy labels, and an
ambiguous legacy advisory reference association. All three were corrected.
Fresh Go 1.26.6 build, vet, catalog tests with both retained publisher
directories, generator and quality checks, and diff checks passed.
[Row checks](../../dist/c08-goal/c08-1-final-gates.log).
The [re-reviewed diff](../../dist/c08-goal/c08-1-final.diff) has SHA-256
`cf41bd37e2279945c73583bb466dac905b589c00a49b8a4929cd4518704fd9a1`.

C08.2 gates started against base `8e76fbac58bd6451ff3d9772ef2c914371ba265e` and candidate diff SHA-256 `02537387d9bc673e8d0f0cada7a1db365659b5bd750ea4b15a70e1cdc0e4ffe3`. The gates cover build, vet, affected packages, concurrent SQLite initialization, catalog checks and CLI evidence.

C08.2 corrective verification started. The re-review candidate diff SHA-256 is `32a92e74eeb24f23e86a4f1c6dfae922d0f9deb7ac982c737e39acc68c053703`.

C08.2 adds the reviewed-only CLI source and reusable SQLite retrieval.
Atomic rebuild and concurrent initialization passed. The first review found
a tokenizer mismatch that invented operation matches, omission of one-term
operations, ambiguous text provenance, and missing legacy exact digest
acceptance. The re-reviewed candidate uses SQLite stemming for identity
and operations, supports the single-word refund query, prints the source
and digest for each operation, and preserves both exact and anchored
accepted publisher checks. Fresh whole affected Go packages, targeted
races, generator/quality, built CLI proofs, build, vet and diff checks
passed. [Row checks](../../dist/c08-goal/c08-2-final-gates.log).
[Exact reviewed diff](../../dist/c08-goal/c08-2-final.diff) SHA-256
`32a92e74eeb24f23e86a4f1c6dfae922d0f9deb7ac982c737e39acc68c053703`.

C08.3 gates started at base `468fef327154cf248fb53ab5e0b765898ce46067` with candidate diff SHA-256 `f941a98f015a390fd96c0c44b4bd4475f357207088f01d30abd11ddbd7e1f0f8`. Hosted verification uses only disposable loopback TLS.

C08.3 corrective gates started. Candidate diff SHA-256
`cde667405298caa61bd9ec3a58d1a985c17ecb0c9eac1d21aba6e11df68c5cdb`.

C08.3 uses a dedicated HTTPS transport, rejects redirect and endpoint query
inputs, and preflights UTF-8 count budgets. The first review identified
ambient transport inheritance, mistyped cache modes, query-url credentials,
and malformed UTF-8 accounting. All four corrections passed independent
Astra xhigh re-review. Fresh build, vet, complete affected packages, focused
race checks, generator, catalog quality, and diff validation passed.
[Row checks](../../dist/c08-goal/c08-3-final-gates.log) and
[exact re-reviewed diff](../../dist/c08-goal/c08-3-final.diff).

C08.4 gates started at base `53860fb134c53a36933a9c3a8c37034a3e983234` with candidate diff SHA-256 `c6245c770be80e3a63b36c95b7fcae89af0b5ca9be1fef35d9e57cfad4b55adb`. Public search uses loopback fixtures and a local cache.

C08.4 corrective gates started for root/catalog help and usage documentation. Candidate diff SHA-256 `db2293d6e5bb97e9f43961888823d8b7285ad4d2c54c4b3cfc279d4240f132ea`.

C08.4 exercise covered default Slack message posting, payment refunds,
the exact Slack Web API alias, metadata-only and non-core public fallback,
offline cache reuse, explicit public sources, and scoped RFC 9727 discovery.
The task review found incorrect help placement and imprecise wording about
embedding usage on public fallback. Both corrections passed re-review.
Build, vet, CLI suite and targeted race checks, catalog generator/quality,
built CLI and help checks, and diff validation passed. Root Go test inputs
were unchanged by the help correction, so the
[original root suite](../../dist/c08-goal/c08-4-gates.log) remains applicable.
[Corrective gates](../../dist/c08-goal/c08-4-final-gates.log) and
[exact re-reviewed diff](../../dist/c08-goal/c08-4-final.diff).

C08 whole gates passed at `001e747d67ccf881d01ee22f74cadfed7ca373c8`.
Fresh build, vet, full tests with retained pins, reviewed races, catalog
generation and quality, search/import help, built-default CLI proofs and
diff validation passed. The first attempt encountered an incomplete
nested consumer snapshot during Go discovery. The complete snapshot
was present before the passing rerun.
[Whole verification](../../dist/c08-goal/c08-whole-gates.log).
OpenUdon passed against that exact candidate and Udon passed its ordinary
pinned suite. The optional Udon composed workspace hit the pre-existing
uncached `cloud.google.com/go/compute@v1.49.1` setup limit.
[Consumer receipts](../../dist/c08-goal/consumers/).

Whole round 1 reviews the exact candidate above against
`a4caf6f56c10e10f2fd43cee1c9810fa3bef0e86`.
[Review packet](../../dist/c08-goal/c08-whole-round1-review.txt) SHA-256
`8f4afc4b4383c46a03ba90562755641915aff72951c7aef60c727a5685420744`.
[Canonical diff](../../dist/c08-goal/c08-whole-round1-full.diff) SHA-256
`f578c17985e62366a1da8895113c79755352f2b108f209b5b65d4a7caeb14d44`.

Round 1 returned BLOCK on P1-1. The review carrier dispatched a fresh
Claude Opus 5.5 child at high effort.
[Dispatch receipt](../../dist/c08-goal/c08-whole-round1-dispatch.json).
[Report](../../dist/c08-goal/c08-whole-round1-report.md) SHA-256
`c73014aec56707491b58bf722aad1060d067c34f6ad48d048eb98039209379a2`.

Corrective verification started on diff SHA-256 `38ca9a26fd121b33c8fd187fc071c750bee2ba205c374891b806640cf43022da`.

- P1-1 and L-9: FTS proposes operation candidates. Surface terms distinguish
  operation evidence from provider identity and incidental or stem-only
  words. Accounting, news, maps and listings reproduced false operation
  evidence in round 1. Marketing was already metadata-only, and Home
  Assistant API was already an exact alias. Both remain covered by
  regressions. Slack posting and all three payment refund pins still match.
- P2-1 and P2-2: Weak metadata and semantic matches allow the public chain.
  If public yields no result or fails, return available reviewed suggestions.
  The optional tested cosine threshold supports different models. Offline
  retrieval with no public cache retains available reviewed suggestions.
- P2-3: Exact aliases use the catalog's existing whitespace and hyphen
  normalization. The catalog does not equate underscores with hyphens.
- P2-4 and P2-5: The query establishes the dimension for each lookup.
  Validate and store complete provider batches individually, and recheck
  the per-search budget before refreshing outdated dimensions.
- P2-6: Verify every selected selector against all 26 retained exact sources
  with both existing qualification directories set. Full native inventory
  identifies the reviewed Sheets methods within the larger Discovery specs.
- P2-7: The optional Udon workspace composition was a setup FAIL at
  `cloud.google.com/go/compute@v1.49.1` in the offline module cache.
  It is not a C08 acceptance gate. The legacy `Client.Search` implementation
  and Go dependency files are unchanged. Udon's ordinary pinned suite
  passed. OpenUdon was tested against the exact candidate.
- L-1: Run new whole-package races for the affected root and CLI packages.
- L-2 through L-4 and L-7: Clarify advisory labels, text publisher evidence,
  the operation display limit, and config scope.
- L-5, L-6, L-8: Preserve explicit bad index path errors, reproducible stored
  metadata and reusable vectors, and the original public output shape.
  The reviewed report retains embedding usage.

Required offline recovery limit is 3 under the approved C08 goal and
stock protocol. The first whole build setup and the retained-artifact race
timeout have distinct keys and each consumed one attempt. The first build
setup key recovered after the archive completed and the ordinary whole
gates passed. [First setup failure](../../dist/c08-goal/c08-whole-setup-failed.log).

The round 2 whole gate on the saved diff passed build, vet and
`go test -count=1 -timeout=2m ./...` with both retained source directories.
The combined `go test -count=1 -timeout=4m -race . ./cmd/apitools`
attempt timed out in the existing
`TestCatalogPromotion/digitalocean/pinned-artifact` path. It was
actively processing the catalog ranker under instrumentation and recorded
no race. The CLI package passed in that invocation.
[Combined attempt](../../dist/c08-goal/c08-whole-round2-gates.log).

Round 2 race recovery uses the same timeout and the whole root/CLI package
runner with the existing standard fixture configuration. Tests that require
opt-in retained artifacts report their own skips. The ordinary whole
suite still qualified those exact artifacts. The separate C08 selector
check uses both retained directories and passes all 26 sources under race.
Catalog quality, generation, help, and built CLI proofs passed.
[CLI and exact selector race receipt](../../dist/c08-goal/c08-whole-round2-cli-and-pins.log).
[OpenUdon corrective-source suite](../../dist/c08-goal/consumers/openudon-round2.log) passed.

The standard whole-package offline races passed for the APItools root and
CLI after the timed-out combined run. The recovery completed without a
timeout increase or test change. C08's exact-source selector race passed
separately. [Standard race results](../../dist/c08-goal/c08-whole-round2-standard-race.log).

Whole round 2 started after verification against full HEAD
`001e747d67ccf881d01ee22f74cadfed7ca373c8` and the
[corrective diff](../../dist/c08-goal/c08-whole-round2.diff) SHA-256
`38ca9a26fd121b33c8fd187fc071c750bee2ba205c374891b806640cf43022da`.
[Round 2 packet](../../dist/c08-goal/c08-whole-round2-review.txt) SHA-256
`f8273c4c54a9a1b28ba6fcfaa6f534dc88608ec30d39c455145c706dc3b80796`.

Round 2 returned PASS on the exact corrective delta. The independent
review found no blocker. The prior P1 and the seven prior P2 findings
are closed with the above evidence and dispositions.
[Review report](../../dist/c08-goal/c08-whole-round2-report.md) SHA-256
`702d9cc19a22bd7605978f84cc1a6c89bd4d7c4f41c951c001ac0ebc433808c4`.
[Model and dispatch receipt](../../dist/c08-goal/c08-whole-round2-dispatch.json).

- N-1 carried as FND-07. The built CLI confirms that plural queries can
  lose covered matches and that incidental label terms can still count.
  Recovery needs relevance tests for broader inflections and common
  multi-provider words together. C09's maintenance acceptance uses the
  embedding and evidence contracts without changing capability ranking.
- N-2 checked with the reviewer query `send slack messages`. The built
  CLI returns Slack as provider metadata. A rejected operation proposal
  adds no coverage evidence. Independent provider text determines whether
  that provider is relevant as metadata. The demonstrated plural recall
  defect is recorded in FND-07.
- N-3 carried as FND-08. The default CLI returns reviewed suggestions
  on a public-directory lookup error without a distinct failure signal.
  Explicit public searches retain their error contract. A later reporting
  change should preserve the accepted result and exit behavior.
- The remaining Lower details are recorded in the review. The NVIDIA
  pinned input is tracked in Git and verified by the selector test. The
  default volatile index, JSON-only extra provenance, configured cache
  error behavior and separately budgeted optional semantic lookup retain
  their documented contracts. These details do not require new task rows.
[Checked query evidence](../../dist/c08-goal/c08-whole-round2-followup-query-evidence.log).

No evolution bump. C08 delivers the approved v31 catalog search direction.
C09 remains pending with the accepted search and cache contract available
for its repeat-run evidence. The protected launch reference stays
byte-for-byte unchanged per the owner's instruction.
````
