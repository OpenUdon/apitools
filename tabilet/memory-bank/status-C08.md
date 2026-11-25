# C08 — Searchable catalog and hosted embeddings

**Owner.** APItools. **Dependency.** Accepted [C07](../docs/history/status-C07.md). **Parallel-safe: no.**
See the [milestone specification](milestone.md#c08--searchable-catalog-and-hosted-embeddings).

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| C08.1 — Curate searchable metadata | `[+]` | Give each of the 334 reviewed providers a normalized category, aliases, useful keywords, a concise public description and an evidence-backed description of the supported subset. Stage description changes for review. Cover both docs-only entries and providers with overlapping services. Added 334 profiles, 50 categories, 41 additional aliases, and 1,352 selected operations for 164 providers. Astra xhigh passed re-review. |
| C08.2 — Add SQLite full-text provider retrieval | `[ ]` | Build a reproducible index from accepted metadata and support updates to earlier records. Expose matched capabilities, coverage, provenance and freshness in text and JSON. Exact provider identities and aliases have priority. |
| C08.3 — Add configurable hosted embeddings | `[ ]` | Combine keyword and semantic ranking using the configured HTTPS endpoint, model and budget. Limit payloads to approved public provider text and explicit queries. Reuse vectors for matching content/model and refresh changed inputs. Test budget exhaustion, missing configuration, endpoint failure and offline lexical behavior with a fixture service. |
| C08.4 — Integrate and verify automatic CLI search | `[ ]` | Prefer reviewed results and run existing public discovery for non-core or weakly matched queries under its existing network controls. Preserve explicit source flags, provider scoping and Go search behavior. Exercise “send slack message”, “refund payment”, an exact alias and a non-core query; verify relevant results, coverage labels, fallback and offline operation. |

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
