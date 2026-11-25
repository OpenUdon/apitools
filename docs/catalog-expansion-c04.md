# Catalog expansion C04

The reviewed operation records below bind the exact local source bytes to a
provider, spec reference, and native selector. Authority, licensing, versions,
and authentication sources remain in the built-in catalog. The operation
checks use a separate disposable SQLite registry per provider.

To recheck the pins from files already staged at the listed paths, run the
root package with the artifact directory supplied. The tests verify the pins,
build and write the index, scope discovery to the provider, resolve and
materialize the source, and reproduce the selector from the copied bytes.

```bash
APITOOLS_CATALOG_QUALIFICATION_DIR=catalog-openapi-cache go test -count=1 .
```

Omitting this variable leaves the provider metadata checks enabled and skips
the checks that require downloaded bytes. The command reads local files.
Use the [catalog curation workflow](catalog-curation.md) to stage provider
artifacts and verify that they match the pins. A mutable source URL may now
serve a newer document. Its bytes must match the digest before they can be
used for this qualification. Downloaded documents stay in the ignored cache.

## Cloud and observability

| Provider and surface | Reviewed operation | Native selector | Local artifact | SHA-256 |
| --- | --- | --- | --- | --- |
| DigitalOcean API v2 | Get account metadata (`account_get`) | `#/paths/~1v2~1account/get` | `openapi/digitalocean-api-v2.yaml` (4,606,025 bytes) | `120b1f5ca7ac46159c0428c318c4b1e88598d923081ba338ebe6838f8c1a08ba` |
| Vercel REST API, document 0.0.1 | Retrieve domains for a project (`getProjectDomains`) | `#/paths/~1v9~1projects~1{idOrName}~1domains/get` | `openapi/vercel-rest-openapi.json` (12,497,764 bytes) | `ae3e778a91928ef7753675d31ca52d89dbd5fb2247e4a6e404f4f986124ea027` |
| Datadog Public API v2 | Read the Synthetic Monitoring on-demand concurrency cap (`GetOnDemandConcurrencyCap`) | `#/paths/~1api~1v2~1synthetics~1settings~1on_demand_concurrency_cap/get` | `openapi/datadog-api-v2.yaml` (10,574,047 bytes) | `37702169e22c9f10aaf840d83b7e621896b17463033ae5ff06a5b91af2802189` |

- DigitalOcean publishes the bundled OpenAPI 3.0.0 artifact under Apache-2.0.
  The root and selected operation use the reusable HTTP bearer scheme. The
  operation attaches `account:read` to that scheme using a scope array, outside
  the OpenAPI 3.0 scope representation for HTTP bearer. The catalog overlay
  records the bearer requirement and the documented permission. [Bundle
  source and license](https://github.com/digitalocean/openapi), [account
  endpoint](https://docs.digitalocean.com/reference/api/reference/account/),
  [token scope](https://docs.digitalocean.com/reference/api/scopes/account/read/).
- Vercel publishes OpenAPI 3.0.3 at its REST schema endpoint. The selected
  operation explicitly uses its reusable bearer scheme. The root and some
  operation requirements are empty, while official docs require access tokens
  for protected resources. The overlay records the verified project-domain
  scope. The artifact lacks an explicit redistribution license and its bytes
  are reference-only. [OpenAPI](https://openapi.vercel.sh/), [endpoint
  documentation](https://vercel.com/docs/rest-api/projects/retrieve-project-domains-by-project-by-id-or-name),
  [REST authentication](https://vercel.com/docs/rest-api),
  [API terms](https://vercel.com/legal/api-terms).
- Datadog's official Python client publishes this OpenAPI 3.0.0 API v2 source
  at commit `455958fd542397eed7b7b33a1a0aa2e1b0985932` under Apache-2.0. The
  global security alternatives separate the key pair from `AuthZ`. The selected
  operation inherits those alternatives and declares `billing_read` in an
  extension. The scoped catalog overlay carries the documented personal or
  service access token alternative and the permission without inferring an
  OAuth grant. [Source](https://github.com/DataDog/datadog-api-client-python/blob/455958fd542397eed7b7b33a1a0aa2e1b0985932/.generator/schemas/v2/openapi.yaml),
  [license](https://github.com/DataDog/datadog-api-client-python/blob/455958fd542397eed7b7b33a1a0aa2e1b0985932/LICENSE),
  [endpoint and permission](https://docs.datadoghq.com/api/latest/synthetics/get-the-on-demand-concurrency-cap/),
  [API authentication](https://docs.datadoghq.com/api/latest/authentication/).

Strict refresh validation classifies all three upstream artifacts as
`parseable-openapi-invalid`. The existing catalog index extracts reviewed
operations from these bytes and can leave other operations out. Initial list
operations (Droplets, projects, and roles) exceeded the index's field summary
bounds and were excluded. The selectors above reproduced from materialized
bytes using the native OpenAPI operation inventory. Full-source UWS shape table
construction separately refused all three originals.

## Database, AI and marketplace

| Provider and surface | Reviewed operation | Native selector | Local artifact | SHA-256 |
| --- | --- | --- | --- | --- |
| MongoDB Atlas Administration API v2, `2023-01-01` media version | Return all visible projects (`listGroups`) | `#/paths/~1api~1atlas~1v2~1groups/get` | `openapi/mongodb-atlas-admin-openapi.yaml` (5,196,444 bytes) | `8eaeb645ded2fb91b6c9c07d29b5315ebb9a8d5ae4397f32b512457d74df34c5` |
| Gemini Developer API v1beta, Discovery revision `20261124` | List available models | `#/methods/generativelanguage.models.list` | `google-discovery/google-gemini-discovery-v1beta.json` (401,495 bytes) | `62024ad3940881c88048b9c0c04f72860c9821a89e36422c3f33401a21d4472f` |
| Amazon Selling Partner Sellers API v1 | List marketplace participation (`getMarketplaceParticipations`) | `#/paths/~1sellers~1v1~1marketplaceParticipations/get` | `openapi/amazon-sp-sellers-v1-swagger.json` (29,604 bytes) | `497862ea32de8040453649986e2cd7c6fcc15b55e8022becc783e4a6d6ffcffd` |

- MongoDB includes the official, untransformed OpenAPI 3.0.1 document in its
  Atlas Go SDK at commit `b5034144f37e30894e7dbc863b71165ff685b4a7`. The
  OpenAPI document declares CC BY-NC-SA 3.0 US and links MongoDB's service
  terms. The SDK repository separately declares Apache-2.0. The reviewed
  document is retained for local reference. Project listing inherits the
  reusable OAuth2 client-credentials and legacy HTTP Digest scheme
  alternatives from the root. The authorization
  depends on the principal's Atlas roles and any configured API access list.
  Strict refresh validation records `parseable-openapi-invalid`, and the
  selected operation reproduced from the materialized OpenAPI inventory.
  [Source](https://github.com/mongodb/atlas-sdk-go/blob/b5034144f37e30894e7dbc863b71165ff685b4a7/openapi/atlas-api.yaml),
  [SDK repository license](https://github.com/mongodb/atlas-sdk-go/blob/b5034144f37e30894e7dbc863b71165ff685b4a7/LICENSE),
  [document license](https://creativecommons.org/licenses/by-nc-sa/3.0/us/),
  [MongoDB service terms](https://www.mongodb.com/mongodb-management-service-terms-and-conditions),
  [endpoint](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-listgroups),
  [authentication](https://www.mongodb.com/docs/atlas/api/api-authentication/).
- Gemini's native Discovery document is reference-only. The selected method
  lacks method-level OAuth scopes, and the artifact's Google Cloud Storage
  OAuth scope does not establish the Gemini API-key requirement. The separate
  auth overlay describes the documented `x-goog-api-key` header and the
  provider's project or service-account key restrictions. The selected method
  reproduced in the UWS shape resolver from materialized Discovery bytes.
  The checked-in advisory OpenAPI carries the documented URL, pagination
  parameters and selected output fields. [Discovery](https://generativelanguage.googleapis.com/$discovery/rest?version=v1beta),
  [models list](https://aistudio.google.com/api-reference/models/list-models),
  [model fields](https://aistudio.google.com/api-reference/models),
  [API keys](https://aistudio.google.com/docs/api-key),
  [site policies](https://developers.google.com/terms/site-policies),
  [Gemini terms](https://aistudio.google.com/terms),
  [advisory overlay](../catalog-openapi-cache/advisory-overlays/google-gemini-models-overlay.json),
  [builder](../catalog-openapi-cache/overlay-builders/build_google_gemini_models_overlay.go).
- Amazon's valid Swagger 2.0 Sellers v1 model is published under Apache-2.0
  at commit `9047916fede443a856e4802a68b322d7eada89ef`. It omits both reusable
  security definitions and root and operation requirements. The scoped review
  overlay records the documented Login with Amazon token transport in
  `x-amz-access-token`. Amazon's role matrix admits Product Listing or Selling
  Partner Insights in all listed regions and additional roles in the Far East.
  [Model](https://github.com/amzn/selling-partner-api-models/blob/9047916fede443a856e4802a68b322d7eada89ef/models/sellers-api-model/sellers.json),
  [license](https://github.com/amzn/selling-partner-api-models/blob/9047916fede443a856e4802a68b322d7eada89ef/LICENSE),
  [endpoint](https://developer-docs.amazon.com/sp-api/reference/getmarketplaceparticipations),
  [auth transport](https://developer-docs.amazon.com/sp-api/docs/using-postman-for-selling-partner-api-models),
  [role matrix](https://developer-docs.amazon.com/sp-api/docs/role-mappings#sellers-api-v1).

## Banking, collaboration and model hub

| Provider and surface | Reviewed operation | Native selector | Local artifact | SHA-256 |
| --- | --- | --- | --- | --- |
| Plaid API `2020-09-14`, build `1.822.2` | Get product permissions of a processor token (`processorTokenPermissionsGet`) | `#/paths/~1processor~1token~1permissions~1get/post` | `openapi/plaid-openapi-2020-09-14.yml` (3,453,430 bytes) | `59246af3f2aef2e55f11a19af1baf6dd1aa19497c3c49a6fc3fe24c68634717f` |
| Miro REST API v2 | Get a specific board (`get-specific-board`) | `#/paths/~1v2~1boards~1{board_id}/get` | `openapi/miro-rest-v2-openapi.json` (1,065,633 bytes) | `f0e17b5b43c92ff2d85c7fbe9a84301e39df1d3c14e6738cd939a6d7142f44dd` |
| Hugging Face Hub API, document `0.0.1` | List references in a model repository | `#/paths/~1api~1models~1{namespace}~1{repo}~1refs/get` | `openapi/hugging-face-hub-openapi.json` (1,247,465 bytes) | `da5be4f2ee566534ee0c5fd74e1c380cc4575f19ca7a90c7fed4c51cc360c01d` |

- Plaid publishes its official OpenAPI 3.0.0 at commit
  `6e75c6b8a7440a3db2d392b636d8cdf024c32e57`. The schema and repository do not
  declare a redistribution license. Provider terms govern access, and the
  document is reference-only. The selected operation inherits the root
  `PLAID-CLIENT-ID`, `PLAID-SECRET` and `Plaid-Version` header set. Plaid's
  official docs also allow the client credentials in the request body and
  make `Plaid-Version` an optional protocol selector. The operation independently
  requires a processor token. The scoped overlay records the reusable header
  credential pair and describes the other placement and version rule.
  [Model](https://github.com/plaid/plaid-openapi/blob/6e75c6b8a7440a3db2d392b636d8cdf024c32e57/2020-09-14.yml),
  [operation](https://plaid.com/docs/api/processors/#processortokenpermissionsget),
  [versioning](https://plaid.com/docs/api/versioning/),
  [terms](https://plaid.com/legal/terms-of-use/),
  [developer policy](https://plaid.com/developer-policy/).
- Miro publishes OpenAPI 3.0.1 in its MIT-licensed official API clients
  repository at commit `1c61f315e70ce8af6326d33c857fd0229541230f`. Its source
  defines the OAuth2 authorization code scheme and names `boards:read`, while
  its root and selected operation omit security requirements. The scoped
  catalog overlay identifies the existing scheme and required `boards:read`
  permission. [Specification](https://github.com/miroapp/api-clients/blob/1c61f315e70ce8af6326d33c857fd0229541230f/packages/generator/spec.json),
  [license](https://github.com/miroapp/api-clients/blob/1c61f315e70ce8af6326d33c857fd0229541230f/LICENSE),
  [endpoint](https://developers.miro.com/reference/get-specific-board-1),
  [OAuth documentation](https://developers.miro.com/docs/getting-started-with-oauth).
- Hugging Face publishes OpenAPI 3.1.0 from its Hub API documentation. The
  artifact lacks an explicit redistribution license and is reference-only.
  The selected operation has no operation ID and inherits the global bearer
  requirement. Hugging Face's client documents anonymous access for its public
  repositories and an optional token for access allowed by the user's
  permissions. The scoped security overlay records the optional bearer
  scheme and describes the access-dependent anonymous alternative in prose,
  since the catalog's requirement sets cannot represent an empty alternative.
  [Specification](https://huggingface.co/.well-known/openapi.json),
  [Hub API docs](https://huggingface.co/docs/hub/api),
  [client reference](https://huggingface.co/docs/huggingface_hub/main/package_reference/hf_api#huggingface_hub.HfApi.list_repo_refs),
  [token scopes](https://huggingface.co/docs/hub/security-tokens),
  [terms](https://huggingface.co/terms-of-service).

Strict refresh validation returned `parseable-openapi-invalid` for all three
originals. Miro's board list and Plaid's institution lookup exceeded the
index's field-summary limits. The narrower operations in the table passed
provider-scoped discovery and reproduced from the same materialized OpenAPI
bytes. Plaid's endpoint documents a read purpose. The source-based effect
classifier leaves that operation's effect unknown.

## Reviewed catalog coverage

All nine candidates passed the artifact qualification. The catalog contains
325 reviewed providers. Its primary coverage is 101 OpenAPI, 14 Swagger, one
OpenAPI index, 29 Smithy, 22 Google Discovery, one Dropbox Stone and 157
human-documentation references. Of the human-documentation providers, 136 have
an advisory endpoint overlay and 21 have no checked-in endpoint overlay.
Dropbox, Google Calendar, Google Gemini API and Google Sheets have advisory
endpoint overlays alongside their primary machine-readable sources.
