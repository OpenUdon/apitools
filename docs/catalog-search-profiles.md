# Reviewed provider search profiles

`catalog.BuiltInProviderSearchProfiles()` reads embedded profiles keyed by
the 334 built-in provider IDs. Each profile has a normalized category, a
public description, search keywords, and exact aliases. The profiles reuse
catalog identities and add researched spellings such as `slack web api`.
The source catalog and its generated manifest keep their own lifecycle.

The public search fields are the display name, category, description,
aliases, keywords, and selected operation labels. Descriptions and keywords
describe provider services. Operation labels require local source evidence.
Each label carries the selector and digest of its exact artifact, with
source identity and reference verification date available alongside it.
Selectors, URLs, revision strings, and authentication metadata stay outside
the public text used for matching or embedding.

| Coverage | Evidence |
| --- | --- |
| `selected-subset` | Selected operations in checked-in documentation-derived advisories or byte-pinned publisher sources. |
| `docs-only` | Documentation references with no selected operation evidence. |
| `reference-only` | A catalog reference to a machine-readable source or index with no selected operation evidence. |

The current profiles record selected operation evidence for 164 providers,
covering 146 advisories and 26 pinned publisher sources. The other 170
profiles expose source references, including 19 with documentation only.
Advisory capability labels come from selected checked-in operation summaries.
Some cover session setup or administrative operations documented by the
publisher. Their inclusion records source coverage. The source reference
date records the catalog's check of that reference. An advisory's digest
identifies its reviewed local bytes.

Slack message posting and public-channel history come from the scoped
advisory verified in [C07](catalog-qualification-c07.md#stripe-google-sheets-and-slack).
The catalog distinguishes these routes from Slack's archived Swagger. The
C08 refund selections use local publisher bytes at the existing pins. Stripe
has `PostRefunds`, Square has `RefundPayment`, and Mercado Pago has
`createRefund`. Razorpay's advisory covers payment lookup by ID. WooCommerce's
advisory covers creating an order refund.

To maintain a profile, edit its row in
[provider_search_profiles.json](../catalog/data/provider_search_profiles.json).
Keep profile rows in provider ID order and use exact provider names or
recognized variants for aliases. Use keywords for unqualified capabilities
that can apply to several providers. Add operation evidence in
[provider_search_coverage.json](../catalog/data/provider_search_coverage.json)
with the provider ID, catalog reference ID, exact artifact digest, and
selector. Checked-in advisories also require the registered artifact ID.
Verify a new or changed selection against the identified source before updating
its annotation. Pinned selections come from the reviewed [C04](catalog-expansion-c04.md),
[C05](catalog-expansion-c05.md), and [C07](catalog-qualification-c07.md) sources
and the local C08 refund checks.

`go test ./catalog` checks coverage, alias ambiguity, and the checked-in
advisory pins and selectors. To reproduce the refund checks, set
`APITOOLS_C07_QUALIFICATION_DIR` to the retained C07 publisher directory
containing `stripe-ga.json` and `APITOOLS_CATALOG_QUALIFICATION_DIR` to the
retained catalog artifact root. Run `go test -count=1 ./catalog -run
'^TestProviderSearchPinnedRefundOperations$'` with those local inputs.
With both directories set, run `go test -count=1 -run
'^TestReviewedSearchPinnedSelectors$' .` to resolve every selected operation
against all 26 pinned publisher sources through the local parsers.

## SQLite provider search

Use `apitools search --query "send slack message"` for a reviewed catalog
lookup. Add `--source reviewed` to include general metadata matches even
when the query has no strong reviewed result. Add `--json` for the versioned
`apitools.reviewed-provider-search.v1` report. Exact provider IDs, names,
and aliases rank first under the catalog's identity normalization. Full-text
search proposes operations. A covered match requires substantive evidence in
the selected label, allowing simple singular and plural forms such as
“message” and “messages.” Provider identity words keep their exact matching
rules. A query for “list user” can match “list users.” A label that mentions
a user or product while describing another resource needs additional evidence
from the query. A broad word such as “time” alone can still find providers by
metadata. The report identifies the match type
and gives selectors and digests for up to 10 matched operations per provider.

With no source or provider URL specified, the CLI returns reviewed results
when the top match is an exact provider or a covered operation. Metadata and
semantic suggestions allow public discovery. When that chain cannot supply a
result, available reviewed suggestions remain visible. If the public lookup
fails, the CLI reports it on stderr in both text and JSON modes. Offline mode
uses local reviewed data and any configured public cache. Explicit `--source auto`,
individual public sources and `--provider-url` use the public path directly
with its original output, cache and network controls. `--config` applies only
to the reviewed path. Automation that requires a stable output shape can
select `--source reviewed` or `--source auto`.

Search builds a private in-memory SQLite index by default. To reuse an
index on disk, pass `--config FILE` with `index_db` set to a path in an
existing local directory. The index refreshes transactionally when the
profile metadata digest changes. It uses a separate SQLite database from
catalog maintenance and the public-directory search cache. The optional
`state_db` setting points to the existing maintenance database. Relative
paths resolve from the configuration file's directory.

The report includes the catalog review date for each source reference.
An accepted publisher check is shown separately when its reference,
locator, and available catalog pin or exact held digest agree with the
reviewed source.
Accepted checks come from the maintenance database in read-only mode.

## Optional hosted embeddings

The optional `embeddings` object configures an OpenAI-compatible HTTPS
embeddings endpoint and per-search budgets. A full cold lookup can need
up to 334 public profile inputs and the explicit query. The budgets count
requests, inputs and UTF-8 text bytes. For example, a reusable configuration
can contain these fields.

```json
{
  "index_db": "reviewed.sqlite",
  "state_db": "catalog-maintenance.sqlite",
  "embeddings": {
    "endpoint": "https://embeddings.example.net/v1/embeddings",
    "model": "operator-selected-model",
    "max_requests": 16,
    "max_inputs": 350,
    "max_bytes": 1048576,
    "min_similarity": 0.8
  }
}
```

The adapter sends the documented public profile fields and the exact
submitted query. Source URLs, digests, selectors and auth evidence stay local.
Use an endpoint URL without a query string and an endpoint that accepts
requests without credentials from APItools.
Provider vectors are reused by endpoint, model and content digest. The
explicit query vector is computed for each semantic lookup. Completed
provider batches whose dimensions agree with the query are reusable after
an interruption. The adapter refreshes missing providers and vectors in an
outdated dimension within the configured budgets. The optional `min_similarity`
sets the cosine threshold for semantic suggestions and defaults to 0.8.

The reviewed text and JSON reports include counted usage when embedding
configuration is supplied. Default search can attempt embeddings before
falling through to a public report, which keeps its original fields.
Exact identity matches, offline mode, missing configuration,
budget exhaustion and endpoint errors use lexical results. A semantic
provider match indicates relevance to a profile with reviewed operations.
Operation matches carry their own selector and artifact provenance.
