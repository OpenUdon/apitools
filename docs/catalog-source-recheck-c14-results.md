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
C14.3 qualification.

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
