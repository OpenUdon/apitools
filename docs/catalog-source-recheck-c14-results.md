# C14 publisher recheck results

One supervised recheck finished on 2026-11-27 in 36.45 seconds. It recorded
outcomes for all 30 providers in the [approved source list](catalog-source-recheck-c14.md#initial-source-urls).
It issued 64 GETs, including one followed redirect, and read 5,551,978
response-body bytes. Four listed URLs had DNS failures before HTTP dispatch.
Two Contentful URLs were deferred after the first returned HTTP 429. The run
refused 13 redirects by origin or path. Each request used the existing guarded
transport without publisher credentials.

The run retained 19 text responses from 14 providers, totaling 2,406,250 bytes.
These are inputs for offline qualification. Some responses are general support
pages or rendered reference shells. The live run found no in-scope direct spec
link to preflight, so it issued no HEAD requests and downloaded no specs.
Offline inspection of the captured Formstack overview found the filename
`api2025-openapi.yaml` in page metadata. A stable download URL and its size
remain unverified. The catalog's source and overlay decisions are pending
C14.3 qualification at the end of the live run. See the
[offline dispositions](#offline-source-qualification) for the result.

The [local event journal](../dist/c14-source-recheck/c14.2-jGNDPf/evidence/events.jsonl)
records each seed or followed hop, UTC time, status or failure, publisher
reference, and artifact digest where applicable. Its SHA-256 is
`b205c09971795cfb7781316fc7ee12a626e4c3e58336a824d78fd0a581020b49`.
The approved scope document SHA-256 is
`866fe801c76c841768802c131ba56a25ae7989e687c7cd3a5da5945632fa060a`.

## Provider outcomes

Completion times are UTC on 2026-11-27. Publisher authority for each seed is
its tracked catalog reference. An arrow gives the origin of a refused redirect.
A same-host refusal records the runner's conservative path filter decision.

| Provider | Time | Retained documentation or principal outcome | Other observed outcomes |
| --- | --- | --- | --- |
| `emelia` | 10:04:57 | DNS unavailable at `docs-old.emelia.io` | Public schema and operation documentation remain unresolved. |
| `orbit` | 10:04:58 | Shutdown notice returned HTTP 530 | Historical reference redirected to `api.orbit.love` and was refused. Lifecycle corroboration remains unresolved. |
| `aftership` | 10:04:58 | HTTP 403 on all three tracked docs | |
| `agile-crm` | 10:04:59 | [Publisher REST documentation repository](https://github.com/agilecrm/rest-api) | `www.agilecrm.com/api` returned HTTP 403. |
| `aircall` | 10:04:59 | Redirect → `developers.aircall.io` | |
| `autopilot` | 10:05:00 | [Ortto help](https://help.ortto.com/a-376-autopilot-how-to-use-autopilots-api) and [Apiary reference response](https://autopilot.docs.apiary.io/) | |
| `clearbit` | 10:05:00 | HTTP 403 on all four tracked docs | |
| `clockify` | 10:05:02 | [API overview](https://clockify.me/help/getting-started/clockify-api-overview) | `docs.clockify.me` was over the body-size limit. |
| `cockpit` | 10:05:03 | [Core documentation](https://getcockpit.com/documentation/core) | The specific API path returned HTTP 404. |
| `contentful` | 10:05:04 | API reference returned HTTP 429 | The auth and management references on that origin were deferred. |
| `demio` | 10:05:12 | [API limitations](https://help.demio.com/en/articles/4544025-api-limitations) | The public Apiary reference timed out. |
| `erpnext` | 10:05:16 | [Frappe REST API reference](https://docs.frappe.io/framework/user/en/api/rest) | The alternate ERPNext path returned HTTP 404. |
| `eventbrite` | 10:05:16 | Same-host redirect refused | |
| `flow` | 10:05:17 | DNS unavailable for all three `developer.getflow.com` URLs | |
| `formstack` | 10:05:18 | [API overview and spec filename lead](https://developers.formstack.com/v2.0/reference/api-overview) | The auth reference returned HTTP 404 and the support article HTTP 403. |
| `freshdesk` | 10:05:19 | Developer reference over the body-size limit after following `/api` to `/api/` | The support article's same-host redirect was refused. |
| `freshservice` | 10:05:21 | [API support](https://support.freshservice.com/support/solutions/articles/50000012704-working-with-apis-in-freshservice) and [upgrade guidance](https://support.freshservice.com/support/solutions/articles/50000004220-introducing-upgraded-apis-for-freshservice) | `api.freshservice.com` reference was over the body-size limit. |
| `home-assistant` | 10:05:21 | Redirect → `contributors.home-assistant.io` | |
| `invoice-ninja` | 10:05:23 | [Developer guide](https://invoiceninja.github.io/docs/developer-guide) and [API reference](https://invoiceninja.github.io/docs/api-reference/invoice-ninja-api-reference) | The separate `api-docs.invoicing.co` reference reached the body cap. |
| `mailjet` | 10:05:24 | [Publisher API guides repository](https://github.com/mailjet/api-documentation) | The three support articles returned HTTP 403. |
| `mautic` | 10:05:24 | [API knowledge base](https://kb.mautic.org/article/what-is-mautic-039%3Bs-api.html) | Developer docs redirected to `devdocs.mautic.org`. |
| `monica-crm` | 10:05:24 | HTTP 404 on the tracked API documentation | |
| `profitwell` | 10:05:26 | [API v2 Apiary reference response](https://profitwellapiv2.docs.apiary.io/) | The API v1 documentation redirected to `www.profitwell.com`. |
| `reddit` | 10:05:27 | [Developer overview](https://developers.reddit.com/docs/capabilities/server/reddit-api), [OAuth reference](https://www.reddit.com/dev/api/oauth), and [API reference](https://www.reddit.com/dev/api/) | Support wiki returned HTTP 403. |
| `salesmate` | 10:05:28 | [Public API reference response](https://apidocs.salesmate.io/) | Two support articles returned HTTP 403. |
| `splunk` | 10:05:30 | HTTP 403 at both `docs.splunk.com` references | The ACS OpenAPI page's same-host redirect was refused. |
| `tapfiliate` | 10:05:30 | Redirect → `developers.tapfiliate.com` | |
| `todoist` | 10:05:31 | Both references redirected to [`/api/v1/`](https://developer.todoist.com/api/v1/) on the listed host | The URL matches the journal's target digest. The runner's path filter deferred the target, so its content is unverified. |
| `twake` | 10:05:33 | Transport errors on all three tracked docs | |
| `wufoo` | 10:05:33 | Both references redirected to `main.wufoo.com` | |

## Cached documentation objects

Each object is stored by full SHA-256 as `sha256-<digest>` in the configured
artifact cache outside Git. The run used disposable SQLite state outside the
catalog. Each retained object passed a size and SHA-256 check after the run.
Multiple objects for the same provider appear in journal order.

| Provider | Bytes | SHA-256 |
| --- | ---: | --- |
| `agile-crm` | 227293 | `9e2686ea1b2d3bc9acaa76052cdb9d051c53fd72b6d9ee99c81e20a394a40258` |
| `autopilot` | 12362 | `452a528201b86cc949d0de4ecbbc341303dfd888ed04eb55403e06b1ce4b55a4` |
| `autopilot` | 6073 | `c903b8cdc40002dcfe0d0854d7890b4e2b4e9852358f253e1446dc14d254ec1a` |
| `clockify` | 121362 | `9d93e30580619ddeae453d7e540c55cd21463a9fb494b01c80e700aade90c816` |
| `cockpit` | 19952 | `ac47ba17363dcf650a6dd8b33f50e8294b928b631e3b0aa70bff65fedba20649` |
| `demio` | 75376 | `f5a5160dd89dffd8ed8268724b3e1978ed4b2cc542f3b75f3381bab2449ad89a` |
| `erpnext` | 235578 | `4d4489d46dbf689c520c4905f77a8258c8c51b5979fbf8a6462116c759d3be0b` |
| `formstack` | 374283 | `e59527316b8ecbf1abca7df42fcb60b80ec83075957a284709322bb81798b206` |
| `freshservice` | 117713 | `0144612caaca12c9f5744266bf5c4e46f7f76e211930d93feb67749951f891f5` |
| `freshservice` | 122312 | `a6e62213054408b0b374a92e8edde7f3639c17a4e9b01550dfb514d4a950b6ab` |
| `invoice-ninja` | 32661 | `31b83eaea42811eced5fac5deba7fc3bf8f4a305ce335b9ec9ed6523e68b6117` |
| `invoice-ninja` | 58392 | `cdf23050fc68d5d4fa2a22be07ad5dbf809405eb42373a9f0795bb4e695efa35` |
| `mailjet` | 215677 | `3f663a14e2ae4c17e6c0017b0dee5e47c293e28b05f35046ed05284802da28a6` |
| `mautic` | 47875 | `7b39f10bcd226d71ed2eae27f46fbb251a18fa5227fc99cbdb4a4b517af4276e` |
| `profitwell` | 5765 | `86f823cad1dc365fbac8c0dcb78da8e14e04fbff483e4a04451ecb8375667ac4` |
| `reddit` | 302038 | `b0c3ac8b71cce38fbeaca7c3957a3de8577d4eb224a0cbb271b7c608f927859e` |
| `reddit` | 89140 | `de1b8fb851ea8991c3cb3bd6a62c7902ea7d459c922e960b342b8cb9d05111e2` |
| `reddit` | 301524 | `3c47aadb2fc740c8b224216210eeba9fc81001f8cfbe5c89b38a11673155b6a4` |
| `salesmate` | 40874 | `88f739017c512fe6592fceb4a58169887f93205c309ce8152d7ee0bd1ea3afb8` |

## Offline source qualification

The offline review refreshed two advisory overlays, retained 26 other
advisory overlays, and retained the two docs-only entries. The retained
publisher material supported a new ERPNext operation and a request-header
correction for Invoice Ninja. It supplied no standalone machine-readable
replacement for native validation or publication. Provider and source IDs
stay stable. Dates advance only on the Frappe REST and token-auth references
and the Invoice Ninja developer guide used for these changes.

- **ERPNext.** The captured Frappe REST docs show `GET
  /api/method/frappe.auth.get_logged_user`, the `Authorization` header with a
  `token api_key:api_secret` value, and a successful JSON object containing a
  `message` string. The overlay represents that read-only operation and its
  token-auth requirement. OAuth and cookie-session paths remain outside this
  token-auth subset. The advisory registration and search profile link to
  the verified Frappe REST reference. [Operation](https://docs.frappe.io/framework/user/en/api/rest#remote-method-calls).
  [Token authentication](https://docs.frappe.io/framework/user/en/api/rest#1-token-based-authentication).
- **Invoice Ninja.** The captured developer guide requires
  `X-Requested-With: XMLHttpRequest` for API requests. The advisory overlay
  models it as a required fixed-value header parameter alongside the existing
  `X-API-TOKEN` security requirement. The optional `X-API-SECRET` applies
  to `/api/v1/login`, outside the reviewed routes. The multipart invoice-upload
  example leaves the full upload operation contract unqualified.
  [Required headers](https://invoiceninja.github.io/docs/developer-guide#api-headers).

| Provider | Disposition | Remaining evidence gap |
| --- | --- | --- |
| `emelia` | Retain docs-only. The GraphQL reference did not resolve. | Find a publisher schema or operation-level GraphQL documentation. A generic POST wrapper would hide the operation shape. |
| `orbit` | Retain historical docs-only and auth metadata. | The shutdown notice could not be read again and the historical redirect points to an unqualified API origin. Fresh lifecycle corroboration is missing. |
| `aftership` | Retain overlay. | HTTP 403 left the tracking endpoint and auth documents unverified by this run. |
| `agile-crm` | Retain overlay. The publisher repo landing page corroborates ownership of the REST docs. | The landing capture showed an unloaded README and the publisher API page returned 403. Endpoint and Basic auth details need readable source content. |
| `aircall` | Retain overlay. | The redirected publisher host needs source and path qualification. |
| `autopilot` | Retain overlay. Ortto help confirms the legacy Autopilot Journeys integration. | The Apiary shell lacks raw operation and API-key requirements. Ortto's current API is separate from the legacy surface under review. |
| `clearbit` | Retain overlay. | HubSpot-era service availability, operations and auth could not be scoped from the tracked help paths returning 403. |
| `clockify` | Retain overlay. The help article corroborates the API and API-key guidance. | The full reference exceeded 3 MiB. Its operation contracts and token alternatives remain unqualified here. |
| `cockpit` | Retain overlay. The core index links to API introduction, authentication, content and assets pages. | The specific `/documentation/core/api` returned 404. The linked child pages remain unfetched. |
| `contentful` | Retain overlay. | The origin returned 429, so this run deferred the remaining auth and management sources. |
| `demio` | Retain overlay. Publisher support confirms API rate limits and custom-field identifiers. | The Apiary reference timed out. The retained support page gives no sufficient route or API-key/secret contract for an overlay change. |
| `erpnext` | Refresh overlay. Index the authenticated read-only current-user route and `message` response. | The alternate ERPNext path returned 404. The revised operation is sourced to the Frappe REST docs. Site-specific DocType and broader auth coverage remain outside the reviewed subset. |
| `eventbrite` | Retain overlay. | The same-host redirect target remains unread and unqualified for endpoint and auth changes. |
| `flow` | Retain overlay. | The developer host did not resolve for any tracked path. Replacement docs and service lifecycle evidence are missing. |
| `formstack` | Retain overlay. The overview corroborates OAuth bearer authentication for Forms API v2. | The 2.0 branch metadata lists `rest-api.json` at the relative registry URI `/branches/2.0/apis/rest-api.json`. The 2025.0 branch separately lists `api2025-openapi.yaml`. The retained 2.0 metadata gives no sync URL or native schema body. Both leads still need stable download URLs, size checks and compatibility review. |
| `freshdesk` | Retain overlay. | The canonical developer reference exceeded the cap and the support redirect remains unqualified. |
| `freshservice` | Retain overlay. Official support corroborates the v2 migration and API-key authentication. | The full v2 contracts exceeded the cap. The articles alone leave response schemas and operation coverage unqualified. |
| `home-assistant` | Retain overlay. | The redirect to the contributor docs needs publisher path, endpoint and auth verification. |
| `invoice-ninja` | Refresh overlay. Record the required fixed-value request header alongside the API token. | The primary rendered definition reached the cap. A bounded standalone definition or additional operation contract is still needed. |
| `mailjet` | Retain overlay. The publisher repo confirms ownership of the API guides. | The landing page and blocked support sources leave endpoint and Basic auth details unqualified in this run. |
| `mautic` | Retain OAuth bearer overlay. The captured knowledge-base article documents OAuth 2 and a separately enabled Basic-auth option. | The main developer reference moved. The scope and applicability of Basic auth across current deployments and reviewed operations remain unqualified. |
| `monica-crm` | Retain overlay. | The tracked docs returned 404. A replacement publisher reference is missing. |
| `profitwell` | Retain overlay. The Apiary shell identifies the v2 API. | The shell lacks endpoint and auth definitions. The v1 redirect and current lifecycle remain unqualified. |
| `reddit` | Retain overlay. The retained OAuth operation reference confirms the existing current-user route and its `identity` scope. | The overlay still describes bearer transport. Complete operation-specific scopes and authorization flow metadata need separate qualification for the reviewed subset. |
| `salesmate` | Retain overlay. The public reference returned a rendered shell. | The capture lacks operation and dual-header details. The two supporting articles were blocked. |
| `splunk` | Retain overlay. | The Splunk Enterprise REST files were blocked. The ACS OpenAPI page's redirect is unqualified and does not establish Enterprise REST coverage. |
| `tapfiliate` | Retain overlay. | The redirected publisher origin needs a bounded source check. |
| `todoist` | Retain existing v2 overlay. | Both seeds point to `/api/v1/` on the developer docs host. The deferred target's contents and migration implications are unverified. |
| `twake` | Retain overlay. | Transport failures left the docs and auth contract unread. |
| `wufoo` | Retain overlay. | The new publisher origin needs a bounded source check before an endpoint or auth update. |
