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

Review gate iterations: 0

Whole review pending, default limit 10. Use the [test policy](test-policy.md)
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
