# C07 reviewed catalog source qualification

The first publisher inventory recorded 334 provider summaries and 195 eligible
reference outcomes for the catalog fingerprint in the C07 ledger. This review
qualifies additional source metadata against pinned evidence. The inventory
report records the separate observations made during that run.
[C07 ledger](../tabilet/memory-bank/status-C07.md).

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

The HubSpot discovery index advertised 130 service families. It marks some
`LATEST` and `STABLE` releases with a `NONE` banner and marks other entries with
`PUBLIC_BETA` or `DEVELOPER_PREVIEW` evidence. The updater preserves both the
stage and the banner. Only the two pinned rollout `424` artifacts above have
qualified operations here. Other stable index entries and previews remain
candidates. The index's own spec URLs returned an unusable zero-length HEAD
response, so qualification used HubSpot's separately published repository.
The larger rollout `321895` variants produced an operation inventory but failed
the native shape-table gate. The two selected rollout `424` documents passed it.
[Publisher index](https://api.hubapi.com/public/api/spec/v1/specs),
[HubSpot's public spec collection](https://github.com/HubSpot/HubSpot-public-api-spec-collection),
[HubSpot's release notes](https://developers.hubspot.com/changelog/fall-2026-spotlight).

For both families, the updater comparator returned
`compatibility-unexamined` for the pinned v3-to-2026-09 comparison and the
rollout `424`-to-`321895` comparison. These checks use the current review
selectors and the pinned publisher bytes. Neither transition is qualified
for automatic migration.

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
`dfea73e06d146c368d7f94b52ac90796dc4e27e1`, remains dated evidence. It lists
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
