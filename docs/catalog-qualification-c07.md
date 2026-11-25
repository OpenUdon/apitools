# C07 reviewed catalog source qualification

The first publisher inventory recorded 334 provider summaries and 195 eligible
reference outcomes for the catalog fingerprint in the C07 ledger. This review
qualifies additional source metadata against pinned evidence. The inventory
report records the separate observations made during that run.
[C07 retired record](../tabilet/docs/history/status-C07.md).

## Dropbox, HubSpot and Kubernetes

| Provider and selected scope | Pinned source | Native operation |
|---|---|---|
| Dropbox folder listing, metadata, upload, download, sharing and current account | [Advisory OpenAPI](../catalog-openapi-cache/advisory-overlays/dropbox-core-api-overlay.json) rebuilt from its [builder](../catalog-openapi-cache/overlay-builders/build_dropbox_overlay.go) and checked against Dropbox Stone revision `4876d1a552f34e10a46e293aaf7e344c452af13d` | `POST /2/files/list_folder` and the other seven operations in the overlay |
| HubSpot CRM Contacts `2026-09`, rollout `424` | [Contacts OpenAPI](https://raw.githubusercontent.com/HubSpot/HubSpot-public-api-spec-collection/c6f46a91db551f1f4fcbe61050ea002050e27d55/PublicApiSpecs/CRM/Contacts/Rollouts/424/2026-09/contacts.json) | `GET /crm/objects/2026-09/{objectType}` with `contacts` |
| HubSpot CRM Deals `2026-09`, rollout `424` | [Deals OpenAPI](https://raw.githubusercontent.com/HubSpot/HubSpot-public-api-spec-collection/c6f46a91db551f1f4fcbe61050ea002050e27d55/PublicApiSpecs/CRM/Deals/Rollouts/424/2026-09/deals.json) | `GET /crm/objects/2026-09/{objectType}` with `deals` |
| Kubernetes release `v1.37.2`, Core v1 group | [Tagged OpenAPI v3](https://raw.githubusercontent.com/kubernetes/kubernetes/v1.37.2/api/openapi-spec/v3/api__v1_openapi.json) | `GET /api/v1/namespaces/{namespace}/pods` |

Dropbox's official SDK model still uses Stone. No standalone official OpenAPI
was found through the publisher's API model and HTTP documentation. All eight
advisory routes have matching route definitions and OAuth scopes in the
reviewed Stone files. Upload and download use the content host. The checked-in
overlay remains the bounded OpenAPI option, with Stone retained as provenance.
[Dropbox API model](https://github.com/dropbox/dropbox-api-spec/tree/4876d1a552f34e10a46e293aaf7e344c452af13d),
[OAuth guide](https://developers.dropbox.com/oauth-guide).

The retained [HubSpot index snapshot](../dist/c07-qualification.P2FWoi/hubspot-index.json)
was acquired on 2026-11-25 at 17:37:52 UTC. Its 157,619 bytes have SHA-256
`f87ebe0f3832e411ae631bbb0c03c97c7592562cc3ca94de17d931fd0690d836`.
The [entry listing](catalog-hubspot-index-c07.tsv) records all 377 versions
across 130 service families. There are 307 stable candidates, 36 public
betas, 32 developer previews and two family/versions with a qualified
repository artifact. The listing has SHA-256
`826513e3cada75828f0b8d60fd4f6dedd008a9ba69c873c6ee50aae46b17301a`.
The stage and banner are preserved for each row.
Developer preview takes precedence when the publisher also shows a beta
banner. The exact index spec URL for a row is
`https://api.hubapi.com/public/api/spec/v2/specs/release/{index_release}/version/{version}`.

The selected family/versions are CRM Contacts `2026-09`, index release
`165147`, and CRM Deals `2026-09`, index release `165150`. Both were
indexed as `LATEST` with the `NONE` banner. Their spec URLs returned a
zero-length HEAD result. Qualification used the exact rollout `424`
bytes from HubSpot's public repository at the same family and version.
The stage describes the indexed family/version. The Git rollout ID is a
separate locator, and the index response bytes remain unverified. The
larger rollout `321895` variants produced operation inventories and failed
native shape projection. The rollout `424` documents passed that gate.
[Publisher index](https://api.hubapi.com/public/api/spec/v1/specs),
[HubSpot's public spec collection](https://github.com/HubSpot/HubSpot-public-api-spec-collection),
[HubSpot's release notes](https://developers.hubspot.com/changelog/fall-2026-spotlight).

For both families, the updater comparator returned
`compatibility-unexamined` for the pinned v3-to-2026-09 comparison and the
rollout `424`-to-`321895` comparison. These checks use the current review
selectors and the pinned publisher bytes. The independent proof covers
the selected list operations. Automatic migration requires a conclusive
compatibility result.

Provider-level resolution selects the Contacts document for HubSpot, Core v1
for Kubernetes, and the Developer API for Amplitude. Select the exact
service reference for other scopes. HubSpot's broader publisher index and
Kubernetes' historical fixture remain separate catalog references.

The two selected HubSpot list operations use separate Contacts and Deals
scopes. The upstream files describe private apps with a synthetic header.
The catalog overlays record the documented Authorization bearer credential and
the OAuth read scope for each exact operation. These are selected operation
proofs. Existing API version paths and other operations require their own
review for migration and auth. [Scopes](https://developers.hubspot.com/docs/apps/developer-platform/build-apps/authentication/scopes),
[HubSpot app authentication](https://developers.hubspot.com/docs/apps/legacy-apps/authentication/intro-to-auth).

Kubernetes publishes the selected Core v1 document under its release tag. The
tag and its resolved commit `c50e5395f1ae2d4a2c05f28fdcc38c0916cc1a32` produced
identical bytes. The Swagger v2 aggregate did not pass the native shape-table
contract. The smaller OpenAPI v3 Core document reproduced the selected Pod
operation. It declares a bearer header but omits the selected operation's
effective requirement. The catalog overlay records the bearer option for that
operation. Cluster authentication and authorization, enabled groups and CRDs
require cluster-specific evidence. The historical HashiCorp fixture is retained
with its third-party classification. [Kubernetes release](https://github.com/kubernetes/kubernetes/releases/tag/v1.37.2),
[API metadata](https://kubernetes.io/docs/concepts/overview/kubernetes-api/),
[Kubernetes authentication](https://kubernetes.io/docs/reference/access-authn-authz/authentication/).

The HubSpot collection describes its specifications as proprietary. The
Kubernetes release uses Apache-2.0, and Dropbox's Stone repository uses MIT.
Upstream bytes are retained in ignored qualification storage. The catalog
tracks URLs and digests, and the Dropbox advisory retains its reproducible
builder. [Kubernetes license](https://github.com/kubernetes/kubernetes/blob/v1.37.2/LICENSE).

## Stripe, Google Sheets and Slack

The Stripe public GA artifact at `0f6d4dfafe157c10df393ce8a27d00b186ab3955`
declares `2026-11-18.endive`. The separate public preview SDK artifact at
that commit declares `2026-11-18.preview`, and remains a candidate. The
reviewed GA `GET /v1/payment_intents/{intent}` reproduces
`GetPaymentIntentsIntent`, a required path ID, and a PaymentIntent response
schema. It inherits the source's Basic or Bearer security alternatives.
The previous pin was `2026-04-22.dahlia`. A dated major-version change still
requires a compatibility review for any existing workflow. [Stripe GA files](https://github.com/stripe/openapi/tree/0f6d4dfafe157c10df393ce8a27d00b186ab3955/latest),
[preview files](https://github.com/stripe/openapi/tree/0f6d4dfafe157c10df393ce8a27d00b186ab3955/preview),
[versioning](https://docs.stripe.com/api/versioning),
[selected operation](https://docs.stripe.com/api/payment_intents/retrieve).

Google's canonical Sheets v4 Discovery URL is stable across revisions. The
catalog held revision `20260511` without a legacy byte digest. The canonical
HEAD returned HTTP 400. A one-byte range preflight returned a 394,610-byte
Content-Length and no complete body. A GET bounded to that declared size
returned revision `20261118` and SHA-256
`82c209f3b283637103a1332ec34b350cb617716a3c6d956e89de44784a9e1573` at the
same canonical URL. The updater's routine HEAD guard still applies. Google's
SDK repository also supplies independent, downloadable v4 snapshots. The pinned
SDK copy at `558e74f3dbcd94a38e2f3e9960c9c032f3a9ee8a` reports revision
`20261109` and SHA-256
`1ec492cc174692abaf513b32286b6e455aa68b08e17b60b8fba21afeaec7925d`.
The earlier SDK copy at `e9015cc2681ce9d6300f4fae2ae6b075c6d75ead` reports
revision `20260213` and SHA-256
`60b9c02762346133a6ff3ee24b3c3f9e1f469666fb772048c0ea0a231030cae1`.
All three reproduce `sheets.spreadsheets.values.get` in the native inventory.
The full native shape proof remains bounded. The updater can retain the two
exact reviewed pins using the selected operation and its read scope. A new
digest requires independent qualification. The seven routes in the
OpenAPI advisory remain a bounded subset with per-operation read or write
scopes. The catalog's auth evidence now preserves that distinction for the
canonical and SDK references. A GitHub snapshot keeps its commit as the
artifact identity for cache reuse and reports the Discovery revision
separately after reading the document. [Canonical URL](https://developers.google.com/workspace/sheets/api/quickstart/js),
[pinned SDK copy](https://raw.githubusercontent.com/googleapis/google-api-go-client/558e74f3dbcd94a38e2f3e9960c9c032f3a9ee8a/sheets/v4/sheets-api.json),
[selected operation](https://developers.google.com/workspace/sheets/api/reference/rest/v4/spreadsheets.values/get).

Slack's archived Swagger `1.7.0`, pinned at
`dfea73e06d146c368d7f94b52ac90796dc4e27e1`, remains dated evidence.
The explicit commit replaces the moving branch URL. The retained bytes
reproduce the earlier Git blob `f7b1affd1fb34f9473cd87980428883019b29c07`
and SHA-256 `8b92da26a3c5b11d20042a9f36d81f1fa6fc9382c5ddc471babb68b91936bc3a`.
The first inventory's validator rejected the raw Swagger. The retained
copy provides historical auth evidence. It lists
classic `chat:write:user` and `chat:write:bot` together for message posting.
It groups all four history scopes into a single requirement. Current docs
specify `chat:write` for text messages and history scopes relevant to the
conversation type. The reproducible [message and history overlay](../catalog-openapi-cache/advisory-overlays/slack-web-message-history-overlay.json)
covers text posting in participating conversations and public-channel
history with `channels:history`, using Slack's current OAuth v2 metadata.
[Archived spec](https://github.com/slackapi/slack-api-specs),
[text message docs](https://docs.slack.dev/reference/methods/chat.postMessage/),
[history docs](https://docs.slack.dev/reference/methods/conversations.history/),
[Slack OAuth v2](https://docs.slack.dev/authentication/installing-with-oauth/).

## Remaining publisher sources and overlays

C07.4 attempted one official documentation page for each of the 160
overlay-first and docs-first providers in the original inventory scope.
It read 130 pages and deferred the 30 listed below. The guarded,
unauthenticated reads used an eight-second timeout
and 3 MiB body limit. They stopped at cross-origin redirects. The
[screen receipt](../dist/c07-qualification.P2FWoi/docs-screen.json) retains
each checked URL and time, successful page digest, observed spec link or
deferral. Its SHA-256 is
`6882f2b9d2da808e496c810b46529905d77e3e8bc54356c58bef817ed9795e1f`.

| Observation | Providers |
| --- | ---: |
| Same-origin spec file links found | 3 |
| Page read with no explicit same-origin spec file link | 127 |
| Cross-origin redirect | 6 |
| Bounded page read failed | 9 |
| Publisher HTTP 403 / 404 / 429 / 502 | 10 / 3 / 1 / 1 |

These observations apply to the sampled page and checked links. They leave
provider-wide availability, other locations, and unchecked versions open.
Additional publisher repositories supplied the Amplitude and Linear sources.

| Official source | Reviewed operation | SHA-256 |
| --- | --- | --- |
| [Amplitude Developer OpenAPI 0.1.0](https://github.com/amplitude/developer-cli/blob/d870b9b1ea5b5a8c29b414eb4db0633a35fd77ba/openapi/bundled/openapi.bundled.json) | `GET /v1/projects` (`listProjects`) | `a11be514b525b0098a2d77d75910907a98abd4e38a0a563313f5d08e30296a16` |
| [Hunter OpenAPI 2.0.0](https://hunter.io/openapi.json) | `GET /domain-search` (`searchDomain`) with the `/v2` server base | `c5ae61475d11a9ea2b6ffcff10990095fc432e099cbe92cad3051e3575d85620` |
| [Linear SDK GraphQL SDL](https://github.com/linear/linear/blob/2440c73cfd297aec89c85d5e46ea2c7068d4b0ba/packages/sdk/src/schema.graphql) | `Query.viewer` (`query.viewer`) | `ee3d23e0c68f31f5297f752e86003592e6c7c7c440559a979a785516298fbc2a` |
| [Storyblok Content Delivery v2 OpenAPI](https://www.storyblok.com/docs/openapi-spec/cdn-v2.openapi.yaml) | `GET /v2/cdn/links` (`listLinks`) | `2e0e46fc67e06525ee67da8e756b90a39e492ba7833a730205259ab6a8a8ec96` |

The official OpenAPI inputs passed guarded updater validation and reproduced
the selected native shapes. The GraphQL SDL passed native qualification and
inventory. Amplitude's Developer API bearer scheme has an operation-specific
`projects:read` requirement in vendor metadata, recorded in the scoped
auth advisory. Other Amplitude API families retain separate docs and auth
guidance. Hunter's selected operation inherits three publisher-modeled
alternatives. The Storyblok CDN operation inherits the `token` query scheme.
The Storyblok fallback records CDN and Management authentication under
their distinct source references and reviewed routes. Linear's
transport advisory distinguishes the raw personal API key in
`Authorization` from OAuth bearer auth. A single generic GraphQL POST
overlay would lose the operation identity.
[Amplitude Developer CLI](https://github.com/amplitude/developer-cli),
[Hunter API reference](https://hunter.io/api-documentation),
[Storyblok Content Delivery](https://www.storyblok.com/docs/api/content-delivery/v2),
[Linear GraphQL auth](https://linear.app/developers/graphql).

The [Postmark API Explorer](https://postmarkapp.com/api-explorer) links
[Server Swagger 1.0.0](https://postmarkapp.com/swagger/server.yml)
and [Account Swagger 0.9.0](https://postmarkapp.com/swagger/account.yml).
The tested raw copies, SHA-256
`fdae25f55e2ef5fb71997f1128cc4148d52e0e395aed04b4f8f9e1c86fdf44a2`
and `bbbcda9432851b8d7d6781d9a090b9d8d70f4396a9288393f428bfa3d5c2102b`,
failed guarded projection. Both use unquoted numeric keys in the reusable
top-level `responses` map. They also omit reusable security schemes and
root requirements. Retain the checked-in Postmark advisory with the
distinct official server/account token headers applied per operation.
The linked documents
remain qualification candidates.
[Postmark authentication](https://postmarkapp.com/developer/api/overview),
[advisory and builder decisions](../catalog-openapi-cache/advisory-overlays/human-doc-overlay-decisions.md).

The Amplitude document declares Proprietary and links the provider terms.
The Developer CLI repository separately carries an MIT license. The
Storyblok document declares MIT. Linear's SDK repository carries MIT.
The inspected Hunter and Postmark files carry no document-level license
declaration. The catalog tracks source links, digests and local advisory
metadata. The downloaded publisher copies are in the ignored qualification
directory.
[Amplitude document terms](https://amplitude.com/terms),
[Amplitude repository license](https://github.com/amplitude/developer-cli/blob/d870b9b1ea5b5a8c29b414eb4db0633a35fd77ba/LICENSE),
[Storyblok document license](https://opensource.org/licenses/MIT),
[Linear repository license](https://github.com/linear/linear/blob/2440c73cfd297aec89c85d5e46ea2c7068d4b0ba/LICENSE).

Of the original 160, four now have a qualified publisher reference, 137
retain checked-in overlay coverage as primary, and 19 retain an explicit
source gap. C07 retained the existing 137 overlays and verified their
checked-in paths, builders and loads. The Postmark, Hunter and Storyblok
builders were rerun with byte-identical output. The updater reports the
overlay and builder, documented gap,
or unexamined coverage in new runs when a provider has no eligible
publisher references. The completed first inventory keeps its original
checkpoints for all 334 providers. Local overlay/gap evidence does not
establish upstream freshness.

## First inventory reconciliation

Run `2OXWJLDNCITU5LOUTAZAVOZ25V` finished on its third supervised
invocation. The original report holds exactly one outcome for each of
334 provider IDs and exactly one outcome for each of 195 eligible
publisher references. The 25 retained objects match their recorded
digests and total 30,760,406 bytes. The state-only inspection from the
current CLI reproduced the original report.
[Original report](../dist/c07-first-inventory.0yJaZ4/resume-2/report.json).

| Original reference outcome | Count |
| --- | ---: |
| Reviewed pin reproduced | 16 |
| Review required | 56 |
| Publisher spec invalid | 43 |
| Publisher check failed | 35 |
| Size unknown | 26 |
| Over the parser limit | 10 |
| Size changed | 4 |
| Publisher retries exhausted | 2 |
| Publisher source review | 1 |
| Publisher source unexamined | 1 |
| Publisher spec unexamined | 1 |

The two exhausted references were Jina Reader and Search. Their version
and validity remain unresolved. The original run recorded no artifact
over 1 GB in the tested scope. The per-reference receipt retains the
actual observed revisions, versions, digests, and reasons for review or
deferral, where the publisher supplied evidence.

The curated catalog contains the same 334 IDs and 203 eligible publisher
references. The current coverage totals are 123 OpenAPI/Swagger/index,
55 other official machine-readable sources, 137 primary advisory overlays
and 19 documented source gaps. The provider-level security audit queues
Amplitude, Hunter and Storyblok for source selection because their
reviewed auth dispositions differ by reference. This queue is separate
from the 56 reference outcomes marked review-required in the historical
publisher inventory. Current curation and supplemental docs
evidence appear alongside the historical outcomes in the
[334-provider reconciliation](../dist/c07-qualification.P2FWoi/c07-5-reconciled-after-review.json),
SHA-256 `80f42e3e2348e80157e36289cdd9c283c0a57f602c26e58cb140003ebcd4f65f`.
The [reconciliation check](../dist/c07-qualification.P2FWoi/c07-5-reconciled-after-review.log)
compares exact provider and reference identities and reads the cached
objects for their lengths and digests.
