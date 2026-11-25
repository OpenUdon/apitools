# C08 — Searchable catalog and hosted embeddings

**Owner.** APItools. **Dependency.** Accepted [C07](../docs/history/status-C07.md). **Parallel-safe: no.**
See the [milestone specification](milestone.md#c08--searchable-catalog-and-hosted-embeddings).

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| C08.1 — Curate searchable metadata | `[ ]` | Give each of the 334 reviewed providers a normalized category, aliases, useful keywords, a concise public description and an evidence-backed description of the supported subset. Stage description changes for review. Cover both docs-only entries and providers with overlapping services. |
| C08.2 — Add SQLite full-text provider retrieval | `[ ]` | Build a reproducible index from accepted metadata and support updates to earlier records. Expose matched capabilities, coverage, provenance and freshness in text and JSON. Exact provider identities and aliases have priority. |
| C08.3 — Add configurable hosted embeddings | `[ ]` | Combine keyword and semantic ranking using the configured HTTPS endpoint, model and budget. Limit payloads to approved public provider text and explicit queries. Reuse vectors for matching content/model and refresh changed inputs. Test budget exhaustion, missing configuration, endpoint failure and offline lexical behavior with a fixture service. |
| C08.4 — Integrate and verify automatic CLI search | `[ ]` | Prefer reviewed results and run existing public discovery for non-core or weakly matched queries under its existing network controls. Preserve explicit source flags, provider scoping and Go search behavior. Exercise “send slack message”, “refund payment”, an exact alias and a non-core query; verify relevant results, coverage labels, fallback and offline operation. |

## Review

Review gate iterations: 0

Whole review pending, default limit 10. Use the [test policy](test-policy.md)
and [milestone review procedure](milestone.md#milestone-review-procedure).
