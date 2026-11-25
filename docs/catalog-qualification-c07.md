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
