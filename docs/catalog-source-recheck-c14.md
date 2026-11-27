# C14 publisher documentation recheck

## Offline reconciliation

C14 follows up on the 30 unreadable C07.4 documentation checks. The retained
receipt matches the SHA-256 in the C07 qualification report. Every sampled URL
still appears in its provider's tracked references. At this checkout the same
30 providers have 72 official documentation references, yielding 69 distinct
initial GET URLs on 45 hosts after removing fragments. Each of the 28 advisory
providers has its checked-in endpoint overlay and builder present. There are
28 overlay files and 14 distinct builder files. [C07 qualification](catalog-qualification-c07.md#remaining-publisher-sources-and-overlays).

- **Emelia** has an official documentation reference for GraphQL and a recorded
  no-overlay decision. Search for public schema or operation documentation
  within the approved publisher scope. A generic GraphQL POST overlay would
  lose operation semantics. The live GraphQL endpoint is outside the scope.
- **Orbit** has an API reference and a separate publisher shutdown notice.
  The next pass seeks lifecycle evidence for this historical entry first.
  Its historical auth metadata is retained for catalog inspection.

The corresponding overlay decisions are recorded in the catalog.
[Overlay and gap decisions](../catalog-openapi-cache/advisory-overlays/human-doc-overlay-decisions.md).

## Proposed supervised run

The run is scoped to the 30 entries below. Start with Emelia and Orbit.
For Orbit, read the shutdown notice before its historical API reference.
Try an alternate listed reference before repeating a path that previously
failed. Fetch each listed URL at most once during this run. The table
identifies the initial URLs, including two tracked publisher repositories
in the `agilecrm/rest-api` and `mailjet/api-documentation` GitHub namespaces.

- **Scope.** Use unauthenticated GETs for docs and guarded HEAD/GET acquisition
  for verified public specs. Follow a redirect or direct documentation/spec
  link only within an origin listed for that provider, using the publisher
  area for the same service. For the two listed GitHub repositories, restrict
  paths to that exact owner and repository on `github.com` and its matching
  `raw.githubusercontent.com` path. Record and defer out-of-scope destinations.
  Reuse the existing guarded transport, including unsafe-host checks at each
  hop. The run includes no API operation execution, sign-in or credential use.
- **Bounds.** One supervised invocation, at most 180 HTTP transactions
  including HEAD/GET calls and redirect hops, within 20 minutes. Run at most
  two requests concurrently and one per origin. Documentation retains C07's
  eight-second timeout and 3 MiB response limit. Consider at most one
  publisher spec per provider in this invocation. Defer a spec before body
  download if its trustworthy size is unknown or exceeds 1,000,000,000 bytes.
  Cap total response bodies at 4,000,000,000 bytes. Record throttling and
  continue with independent sources. Stop on scope, safety, capacity, or
  budget limits. A further supervised invocation needs a separate decision.
- **Supervision.** The coordinator runs the approved command in one retained
  foreground terminal session and monitors the per-provider progress and
  request count. On an owner stop request or a stop condition, interrupt the
  foreground process group with Ctrl+C, wait for it to exit, and retain the
  completed outcomes. Inspect that state before considering a new invocation.
- **Evidence.** Keep one result per provider with the attempted URL, final
  URL or refused redirect, time, status, error or deferral, publisher authority
  evidence, and artifact digest where applicable. Retain qualifying bytes only
  as needed in the configured `cache_dir` outside Git. Use a disposable state
  directory on a native Linux filesystem for this recheck. Check an existing
  `expected_volume_uuid` before cache writes. The local configuration and
  credentials remain outside the tracked report.
- **Adoption.** C14.3 verifies any replacement or overlay change offline,
  including the operation, source provenance, and authentication evidence.
  Unresolved cases retain a named evidence gap. The accepted catalog state,
  recurring maintenance and website release operate under their own scopes.

## Initial source URLs

All rows below are derived from the checked-in provider references. A listed
URL is a candidate for a live check. Its current reachability is unverified.

| Provider | C07 result | Checked-in endpoint overlay | Initial URLs |
| --- | --- | --- | --- |
| `emelia` | `bounded-docs-read-failed` | Docs only | <https://docs-old.emelia.io/> |
| `orbit` | `http-301` | Docs only | <https://orbit.love/blog/orbit-is-joining-postman><br><https://docs.orbit.love/reference> |
| `aftership` | `http-403` | [`aftership-tracking-api-overlay`](../catalog-openapi-cache/advisory-overlays/aftership-tracking-api-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m57_human_docs_overlays.go) | <https://www.aftership.com/docs/tracking/fcd9acb5f448a-api-overview><br><https://www.aftership.com/docs/tracking/2024-07/quickstart/authentication><br><https://www.aftership.com/docs/tracking> |
| `agile-crm` | `http-403` | [`agile-crm-rest-api-overlay`](../catalog-openapi-cache/advisory-overlays/agile-crm-rest-api-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m23_human_docs_overlays.go) | <https://github.com/agilecrm/rest-api><br><https://www.agilecrm.com/api> |
| `aircall` | `http-301` | [`aircall-public-api-overlay`](../catalog-openapi-cache/advisory-overlays/aircall-public-api-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m57_human_docs_overlays.go) | <https://developer.aircall.io/api-references/> |
| `autopilot` | `bounded-docs-read-failed` | [`autopilot-api-overlay`](../catalog-openapi-cache/advisory-overlays/autopilot-api-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m25_human_docs_overlays.go) | <https://help.ortto.com/a-376-autopilot-how-to-use-autopilots-api><br><https://autopilot.docs.apiary.io/> |
| `clearbit` | `http-403` | [`clearbit-api-overlay`](../catalog-openapi-cache/advisory-overlays/clearbit-api-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m23_human_docs_overlays.go) | <https://help.clearbit.com/hc/en-us/articles/6480449602967-Integrate-Clearbit-Prospector-with-Google-Sheets-Using-Zapier><br><https://help.clearbit.com/hc/en-us/articles/6045527495191-How-Do-I-Access-My-Clearbit-API-Key><br><https://help.clearbit.com/hc/en-us/articles/8502992633111-Autocomplete-Name-to-Domain-and-Risk-API-FAQ><br><https://help.clearbit.com/hc/en-us/categories/360000913214-APIs> |
| `clockify` | `bounded-docs-read-failed` | [`clockify-api-v1-overlay`](../catalog-openapi-cache/advisory-overlays/clockify-api-v1-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m22_human_docs_overlays.go) | <https://clockify.me/help/getting-started/clockify-api-overview><br><https://docs.clockify.me/> |
| `cockpit` | `http-404` | [`cockpit-api-overlay`](../catalog-openapi-cache/advisory-overlays/cockpit-api-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m26_human_docs_overlays.go) | <https://getcockpit.com/documentation/core><br><https://getcockpit.com/documentation/core/api> |
| `contentful` | `http-429` | [`contentful-management-api-overlay`](../catalog-openapi-cache/advisory-overlays/contentful-management-api-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m21_human_docs_overlays.go) | <https://www.contentful.com/developers/docs/references/content-delivery-api/><br><https://www.contentful.com/developers/docs/references/authentication/><br><https://www.contentful.com/developers/docs/references/content-management-api/> |
| `demio` | `bounded-docs-read-failed` | [`demio-api-overlay`](../catalog-openapi-cache/advisory-overlays/demio-api-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m34_human_docs_overlays.go) | <https://help.demio.com/en/articles/4544025-api-limitations><br><https://publicdemioapi.docs.apiary.io/> |
| `erpnext` | `http-404` | [`erpnext-rest-api-overlay`](../catalog-openapi-cache/advisory-overlays/erpnext-rest-api-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m28_human_docs_overlays.go) | <https://docs.frappe.io/framework/user/en/api/rest><br><https://docs.frappe.io/erpnext/rest-api> |
| `eventbrite` | `bounded-docs-read-failed` | [`eventbrite-platform-api-v3-overlay`](../catalog-openapi-cache/advisory-overlays/eventbrite-platform-api-v3-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m21_human_docs_overlays.go) | <https://www.eventbrite.com/platform/api> |
| `flow` | `bounded-docs-read-failed` | [`flow-api-overlay`](../catalog-openapi-cache/advisory-overlays/flow-api-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m33_human_docs_overlays.go) | <https://developer.getflow.com/><br><https://developer.getflow.com/errors><br><https://developer.getflow.com/api> |
| `formstack` | `http-403` | [`formstack-api-v2-overlay`](../catalog-openapi-cache/advisory-overlays/formstack-api-v2-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m26_human_docs_overlays.go) | <https://developers.formstack.com/v2.0/reference/api-overview><br><https://developers.formstack.com/reference/authorization><br><https://help.formstack.com/hc/en-us/articles/44593115536403> |
| `freshdesk` | `http-403` | [`freshdesk-api-v2-overlay`](../catalog-openapi-cache/advisory-overlays/freshdesk-api-v2-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m21_human_docs_overlays.go) | <https://support.freshdesk.com/en/support/solutions/articles/225441-is-there-any-documentation-for-the-apis-on-freshdesk-><br><https://developers.freshdesk.com/api> |
| `freshservice` | `http-403` | [`freshservice-api-v2-overlay`](../catalog-openapi-cache/advisory-overlays/freshservice-api-v2-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m23_human_docs_overlays.go) | <https://support.freshservice.com/support/solutions/articles/50000012704-working-with-apis-in-freshservice><br><https://support.freshservice.com/support/solutions/articles/50000004220-introducing-upgraded-apis-for-freshservice><br><https://api.freshservice.com/> |
| `home-assistant` | `http-301` | [`home-assistant-rest-api-overlay`](../catalog-openapi-cache/advisory-overlays/home-assistant-rest-api-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m29_human_docs_overlays.go) | <https://developers.home-assistant.io/docs/api/rest/> |
| `invoice-ninja` | `bounded-docs-read-failed` | [`invoice-ninja-api-overlay`](../catalog-openapi-cache/advisory-overlays/invoice-ninja-api-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m28_human_docs_overlays.go) | <https://invoiceninja.github.io/docs/developer-guide><br><https://invoiceninja.github.io/docs/api-reference/invoice-ninja-api-reference><br><https://api-docs.invoicing.co/> |
| `mailjet` | `http-403` | [`mailjet-rest-api-overlay`](../catalog-openapi-cache/advisory-overlays/mailjet-rest-api-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m22_human_docs_overlays.go) | <https://documentation.mailjet.com/hc/en-us/articles/360043225693-What-is-an-API-key><br><https://documentation.mailjet.com/hc/en-us/articles/360043230093-What-are-the-endpoints-available-for-the-API><br><https://github.com/mailjet/api-documentation><br><https://documentation.mailjet.com/hc/en-us/articles/360044088173-REST-API> |
| `mautic` | `http-301` | [`mautic-api-overlay`](../catalog-openapi-cache/advisory-overlays/mautic-api-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m25_human_docs_overlays.go) | <https://kb.mautic.org/article/what-is-mautic-039%3Bs-api.html><br><https://developer.mautic.org/> |
| `monica-crm` | `http-404` | [`monica-crm-api-overlay`](../catalog-openapi-cache/advisory-overlays/monica-crm-api-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m25_human_docs_overlays.go) | <https://www.monicahq.com/api> |
| `profitwell` | `http-502` | [`profitwell-api-v2-overlay`](../catalog-openapi-cache/advisory-overlays/profitwell-api-v2-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m33_human_docs_overlays.go) | <https://profitwell.com/docs/api/v1/><br><https://profitwellapiv2.docs.apiary.io/> |
| `reddit` | `http-403` | [`reddit-api-overlay`](../catalog-openapi-cache/advisory-overlays/reddit-api-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m24_human_docs_overlays.go) | <https://www.reddit.com/dev/api/oauth><br><https://developers.reddit.com/docs/capabilities/server/reddit-api><br><https://support.reddithelp.com/hc/en-us/articles/14945211791892-Reddit-Data-API-Wiki><br><https://www.reddit.com/dev/api/> |
| `salesmate` | `http-403` | [`salesmate-api-overlay`](../catalog-openapi-cache/advisory-overlays/salesmate-api-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m25_human_docs_overlays.go) | <https://support.salesmate.io/hc/en-us/articles/360043653751-Auto-enroll-Contacts-to-sequence-whenever-a-Deal-stage-is-updated-via-Webhooks><br><https://apidocs.salesmate.io/><br><https://support.salesmate.io/hc/en-us/articles/12864176787609-What-Are-API-s-Usage-and-Working> |
| `splunk` | `http-403` | [`splunk-enterprise-rest-overlay`](../catalog-openapi-cache/advisory-overlays/splunk-enterprise-rest-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m20_human_docs_overlays.go) | <https://docs.splunk.com/Documentation/Splunk/latest/RESTUM/RESTusing><br><https://dev.splunk.com/enterprise/docs/devtools/acs/acsopenapi/><br><https://docs.splunk.com/Documentation/Splunk/latest/RESTREF> |
| `tapfiliate` | `http-301` | [`tapfiliate-api-overlay`](../catalog-openapi-cache/advisory-overlays/tapfiliate-api-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m33_human_docs_overlays.go) | <https://tapfiliate.com/docs/rest/> |
| `todoist` | `bounded-docs-read-failed` | [`todoist-rest-api-v2-overlay`](../catalog-openapi-cache/advisory-overlays/todoist-rest-api-v2-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m21_human_docs_overlays.go) | <https://developer.todoist.com/guides/><br><https://developer.todoist.com/rest/v2/> |
| `twake` | `bounded-docs-read-failed` | [`twake-developers-api-overlay`](../catalog-openapi-cache/advisory-overlays/twake-developers-api-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m27_human_docs_overlays.go) | <https://doc.twake.app/developers-api/api-reference/auth><br><https://doc.twake.app/developers-api/api-reference/message/post-request><br><https://doc.twake.app/developers-api/get-started> |
| `wufoo` | `http-301` | [`wufoo-api-v3-overlay`](../catalog-openapi-cache/advisory-overlays/wufoo-api-v3-overlay.json), [builder](../catalog-openapi-cache/overlay-builders/build_m32_human_docs_overlays.go) | <https://www.wufoo.com/docs/api/v3/forms/><br><https://www.wufoo.com/docs/api/v3/> |
