# C05 catalog expansion qualification

This review covers the remaining nine candidates from the
[catalog upgrade proposal](catalog-upgrade-2026-09.md). The exact reviewed
C04 artifacts and dispositions are recorded [with that batch](catalog-expansion-c04.md).
Each promoted surface has an offline test using a disposable registry and an
operation index. The test retrieves a provider-scoped operation candidate,
checks its artifact identity and native selector, materializes the pinned
bytes and reproduces the selector. The catalog and registrations contain
symbolic metadata. Qualification uses no provider operation or credential.

Use `APITOOLS_CATALOG_QUALIFICATION_DIR=catalog-openapi-cache go test -v -run
'^TestCatalogPromotion$' .` after staging the pinned upstream bytes below.
Tracked advisory documents are tested directly from the repository. An
upstream file served at a mutable URL must match its pin to qualify this
reviewed state.

## Maps, advertising and marketplaces

| Provider and surface | Reviewed operation | Native selector | Local artifact | SHA-256 |
| --- | --- | --- | --- | --- |
| Google Maps Platform, Places API (New) v1, Discovery revision `20261123` | Autocomplete place and query predictions | `#/methods/places.places.autocomplete` | `google-discovery/google-maps-places-v1-discovery.json` (148,978 bytes) | `e6db1c21f8b1c51f9ff1b5495b1bbafdf33cac27e470f0db2c7bad3afa0f7331` |
| Google Ads API v26, CustomerService | List directly accessible customers | `#/services/google.ads.googleads.v26.services.CustomerService/methods/ListAccessibleCustomers` | `grpc-protobuf/google-ads-v26-customer-service.proto` (7,296 bytes) | `0cca59c738498fb0010f39f437d29eaff5dbe7577aae4bcdd011ae445795e685` |
| eBay Sell Account API v1, document `v1.9.3` | List programs the seller has opted into (`getOptedInPrograms`) | `#/paths/~1program~1get_opted_in_programs/get` | `openapi/ebay-sell-account-v1.json` (328,066 bytes) | `3f04cf1fd7c80160312d68be5afd7d0de00eb57fd32d508cbb5eeec526b8bfbb` |

- Google lists Places API (New) as active. Its legacy Maps Platform OpenAPI
  repository redirects ongoing collection maintenance to Postman. The native
  Discovery source lists the query API key and three alternative OAuth scopes
  for Autocomplete. The docs also show the API key in the `X-Goog-Api-Key`
  header. The scoped overlay records the alternatives, including the
  Autocomplete-specific, Places-wide and Cloud Platform scopes. The endpoint
  overlay describes a text prediction request and a subset of prediction IDs
  and display fields. The Discovery document does not declare a
  redistribution license. The pinned artifact is local reference only.
  [Deprecation](https://github.com/googlemaps/openapi-specification),
  [Discovery](https://places.googleapis.com/$discovery/rest?version=v1),
  [method and scopes](https://developers.google.com/maps/documentation/places/web-service/reference/rest/v1/places/autocomplete),
  [endpoint documentation](https://developers.google.com/maps/documentation/places/web-service/place-autocomplete),
  [credential alternatives](https://developers.google.com/maps/documentation/places/web-service/get-api-key?setupProd=configure),
  [Maps Platform terms](https://cloud.google.com/maps-platform/terms),
  [site policies](https://developers.google.com/terms/site-policies),
  [overlay](../catalog-openapi-cache/advisory-overlays/google-maps-places-autocomplete-overlay.json),
  [builder](../catalog-openapi-cache/overlay-builders/build_google_maps_places_autocomplete_overlay.go).
- Google publishes the selected v26 protobuf in the Apache-2.0 Google APIs
  repository at `d7dd444ca9b0b53e46e7594bbb705391e11f89c6`. The file carries
  the same license header. API access has separate terms. The selected RPC's
  request and response definitions are in the pinned file. The parser
  does not resolve imports or project protobuf comments as purpose evidence.
  Discovery returns the correctly typed, pinned operation with an
  indeterminate purpose classification. The same selector resolves in the
  native UWS shape table. The scoped auth overlay records an OAuth access token
  with the `adwords` scope together with `developer-token`. Offline security
  inspection matches this combined requirement to
  `google.ads.googleads.v26.services.CustomerService/ListAccessibleCustomers`.
  The direct-access listing does not need `login-customer-id`. The endpoint overlay documents
  the published REST binding and the resource-name array.
  [Protobuf](https://github.com/googleapis/googleapis/blob/d7dd444ca9b0b53e46e7594bbb705391e11f89c6/google/ads/googleads/v26/services/customer_service.proto),
  [repository license](https://github.com/googleapis/googleapis/blob/d7dd444ca9b0b53e46e7594bbb705391e11f89c6/LICENSE),
  [method](https://developers.google.com/google-ads/api/reference/rpc/v26/CustomerService/ListAccessibleCustomers),
  [call structure](https://developers.google.com/google-ads/api/docs/concepts/call-structure),
  [REST binding](https://developers.google.com/google-ads/api/rest/auth),
  [API terms](https://developers.google.com/google-ads/api/docs/api-policy/terms),
  [overlay](../catalog-openapi-cache/advisory-overlays/google-ads-v26-accessible-customers-overlay.json),
  [builder](../catalog-openapi-cache/overlay-builders/build_google_ads_v26_accessible_customers_overlay.go).
- eBay publishes OpenAPI 3.0.0 for Sell Account v1. The document lists the
  eBay API License Agreement, and the pinned bytes stay in the local review
  cache. The reusable `api_auth` scheme declares authorization-code OAuth.
  The root omits security and all 37 operations declare it individually. The
  selected operation puts `sell.account` and `sell.account.readonly` in
  one requirement. The auth overlay records the alternative permissions in
  separate requirement sets under eBay's scope guidance. The earlier
  investigated Sell Inventory artifact is a different surface. Its item
  lookup exceeded field-summary limits, and the scoped Account
  operation matched the step purpose in offline discovery.
  [OpenAPI](https://developer.ebay.com/api-docs/master/sell/account/openapi/3/sell_account_v1_oas3.json),
  [seller account guide](https://edp.ebay.com/develop/guides/sell/account-management-guide),
  [operation](https://developer.ebay.com/api-docs/sell/account/resources/program/methods/getOptedInPrograms),
  [scope guidance](https://developer.ebay.com/develop/guides/sell/authorization),
  [license agreement](https://developer.ebay.com/join/api-license-agreement).

## Commerce and payments

| Provider and surface | Reviewed operation | Native selector | Local artifact | SHA-256 |
| --- | --- | --- | --- | --- |
| Etsy Open API v3, document `3.0.0` | Get holidays available for a shop to set preferences (`getHolidayPreferences`) | `#/paths/~1v3~1application~1shops~1{shop_id}~1holiday-preferences/get` | `openapi/etsy-v3-openapi.json` (963,751 bytes) | `c68d65ed002ddb6d3381d7d4497850eaa49a80587fcd237f58df7fa52e9fdd0e` |
| Square Connect v2, API version `2026-11-18` | Retrieve a merchant account profile (`RetrieveMerchant`) | `#/paths/~1v2~1merchants~1{merchant_id}/get` | `openapi/square-connect-v2-openapi.json` (3,279,383 bytes) | `76e7cc8ba3fd94c162a3339531aafe798d5556ac34ae03e8253df815b49305f1` |
| Mercado Pago REST, document `1.0.0` | Get payment by ID (`getPayment`) | `#/paths/~1v1~1payments~1{id}/get` | `openapi/mercado-pago-rest-openapi.json` (364,763 bytes) | `da7e3a7482982b1b6880a5f71411edcab54ce354a8eda55a970ec6c09973592b` |

- Etsy publishes the OpenAPI 3.0.2 download from its official reference.
  The artifact retains Etsy copyright and points to its developer terms.
  It is held for local reference. The root requires the reusable
  `x-api-key` scheme and the selected operation requires the app key
  together with the `shops_r` authorization-code OAuth scope. The app-key
  header combines the app's keystring and shared secret, and the selected
  feature applies to United States and Canada shops. The documented
  requirements agree with the source.
  [OpenAPI download](https://www.etsy.com/openapi/generated/oas/3.0.0.json),
  [operation](https://developers.etsy.com/documentation/reference/#operation/getHolidayPreferences),
  [authentication and scopes](https://developers.etsy.com/documentation/essentials/authentication/),
  [API developer terms](https://www.etsy.com/legal/api/).
- Square publishes the OpenAPI 3.0.0 artifact at
  `674f83e43690f385ef89896309fa9f5d612f312a`. The document and
  repository both declare Apache-2.0. API access has separate developer
  terms. The root omits security and the selected operation requires
  `MERCHANT_PROFILE_READ` through its reusable authorization-code OAuth
  scheme. Square also documents a personal bearer access token for an
  integration accessing its own account, which the scoped overlay describes
  alongside the OAuth alternative. The initial ListLocations operation
  exceeded the index's field-summary bounds, so qualification uses
  RetrieveMerchant. The offline refresh validator reported `invalid` for
  the original Square source. Bounded index discovery and the native
  selector inventory reproduced the selected operation from those bytes.
  [Source](https://github.com/square/connect-api-specification/blob/674f83e43690f385ef89896309fa9f5d612f312a/api.json),
  [repository license](https://github.com/square/connect-api-specification/blob/674f83e43690f385ef89896309fa9f5d612f312a/LICENSE),
  [operation](https://developer.squareup.com/reference/square/merchants-api/retrieve-merchant),
  [OAuth permissions](https://developer.squareup.com/docs/oauth-api/square-permissions#merchants),
  [access token alternatives](https://developer.squareup.com/docs/build-basics/access-tokens),
  [developer terms](https://squareup.com/us/en/legal/general/developers).
- Mercado Pago publishes the OpenAPI 3.1.0 artifact at
  `ac1d881294bf8b36869a0afe87c3acdfd9324dfd`. The document and its
  repository declare Apache-2.0. The API terms are separate. The reusable
  `bearerAuth` scheme applies at the root and is restated at the selected
  operation. The documented credential is a private server-side access
  token. The source's OAuth and public-key schemes address other
  surfaces and are outside the selected operation.
  [Source](https://github.com/mercadopago/openapi/blob/ac1d881294bf8b36869a0afe87c3acdfd9324dfd/spec3.json),
  [repository license](https://github.com/mercadopago/openapi/blob/ac1d881294bf8b36869a0afe87c3acdfd9324dfd/LICENSE),
  [operation](https://www.mercadopago.com.ar/developers/en/reference/online-payments/checkout-pro-preferences/get-payment/get),
  [credential documentation](https://www.mercadopago.com.ar/developers/en/docs/your-integrations/credentials),
  [developer terms](https://www.mercadopago.com.ar/developers/en/docs/resources/legal/terms-and-conditions).

## Payments and media with advisory endpoints

| Provider and scope | Reviewed operation | Native OpenAPI selector | Local advisory artifact | SHA-256 |
| --- | --- | --- | --- | --- |
| Razorpay Payments v1 | Fetch a payment with its ID | `#/paths/~1v1~1payments~1{id}/get` | `advisory-overlays/razorpay-v1-payment-overlay.json` (2,537 bytes) | `6ed5333b09f383fe20490433377ec488cc130cc454a432b7d29b1713c9af3f28` |
| WeChat Pay API v3 Global, common merchant mode | Query order status by WeChat Pay order number | `#/paths/~1v3~1global~1transactions~1id~1{id}/get` | `advisory-overlays/wechat-pay-v3-global-order-overlay.json` (4,218 bytes) | `b18df240bfbc3a7505a4515c87dbcd617e99d06ee76c4cfaefc21ec55575c969` |
| Twitch Helix | Get Games and category metadata | `#/paths/~1helix~1games/get` | `advisory-overlays/twitch-helix-games-overlay.json` (4,268 bytes) | `f1db616feef677676dd4c1be42c9cc24f1a8a9f6030282d4835028a5ee6a0781` |

- Razorpay's endpoint reference, last modified 2026-10-08, specifies the
  required string `id` and documents HTTP Basic credentials. The builder
  records the payment identifier, amount in currency subunits, currency,
  status, capture flag and creation timestamp. Razorpay publishes a Postman
  workspace and SDK code. The review found no provider-owned OpenAPI
  download for this surface. Razorpay's payment terms and IP provisions
  apply. Source pages remain reference-only.
  [Endpoint](https://razorpay.com/docs/api/payments/fetch-with-id),
  [auth](https://razorpay.com/docs/api/authentication),
  [payment terms](https://razorpay.com/terms/),
  [Postman publication](https://razorpay.com/blog/?p=7242),
  [overlay](../catalog-openapi-cache/advisory-overlays/razorpay-v1-payment-overlay.json),
  [builder](../catalog-openapi-cache/overlay-builders/build_razorpay_v1_payment_overlay.go).
- WeChat Pay's API v3 Global Query Order page dated 2025-02-25 documents
  the common-mode lookup with a WeChat Pay order number and required
  merchant ID query parameter. The overlay includes order identifiers,
  status and amount metadata. Its custom HTTP scheme represents the
  documented signed `Authorization` header symbolically. The merchant
  order-number route and institutional mode have different inputs beyond
  this subset. The official SDK and Postman resources did not provide a
  provider-published OpenAPI document for this Global route. Provider
  merchant access terms apply, and the source pages remain reference-only.
  [Query Order](https://pay.wechatpay.cn/doc/global/v3/en/4012356497),
  [signature scheme](https://pay.wechatpay.cn/doc/global/v3/en/4012357147),
  [interface rules](https://pay.wechatpay.cn/doc/global/v3/en/4012357105),
  [merchant access](https://pay.wechatpay.cn/doc/global/v3/en/4012357180),
  [SDK](https://github.com/wechatpay-apiv3/wechatpay-go),
  [Postman tooling](https://github.com/wechatpay-apiv3/wechatpay-postman-script),
  [overlay](../catalog-openapi-cache/advisory-overlays/wechat-pay-v3-global-order-overlay.json),
  [builder](../catalog-openapi-cache/overlay-builders/build_wechat_pay_v3_global_order_overlay.go).
- Twitch's Helix Games reference documents repeatable Twitch ID, exact
  name and IGDB ID query inputs, with at least one and at most 100 total
  values. The response includes game IDs, names, box-art URL templates and
  IGDB IDs. The security alternatives pair the app or user OAuth token
  with its matching `Client-Id`, without an additional endpoint scope.
  The review found no Twitch-published Helix OpenAPI. Twitch's developer
  services agreement covers API access and the source page remains
  reference-only. The overlay records the reviewed path and query limits.
  [Endpoint](https://dev.twitch.tv/docs/api/reference/games#get-games),
  [OAuth](https://dev.twitch.tv/docs/authentication/getting-tokens-oauth/),
  [developer services agreement](https://legal.twitch.com/en/legal/developer-agreement/),
  [overlay](../catalog-openapi-cache/advisory-overlays/twitch-helix-games-overlay.json),
  [builder](../catalog-openapi-cache/overlay-builders/build_twitch_helix_games_overlay.go).

The three overlays are original metadata descriptions of the selected
operations. Their builders regenerate the pinned artifacts without network
access. The qualification tests verify provider-scoped discovery and resolve
the same selector in the UWS shape table built from the materialized bytes.

## Combined C04 and C05 disposition

The accepted C04 baseline had 325 reviewed providers. The nine C05
promotions bring the total to 334. All 18 candidates in the two approved
batches qualified. The 120 primary OpenAPI-related providers comprise 105
OpenAPI, 14 Swagger and one OpenAPI index. The 54 other primary machine
sources comprise 29 Smithy, 23 Google Discovery, one Dropbox Stone and one
gRPC/protobuf. Of 160 human-documentation primary sources, 139 have a checked-in
advisory OpenAPI overlay and 21 have no checked-in endpoint overlay.

| Batch | Candidate | Disposition and qualifying operation |
| --- | --- | --- |
| C04 | Amazon Selling Partner | Promoted, list marketplace participation |
| C04 | Datadog | Promoted, read the on-demand concurrency cap |
| C04 | DigitalOcean | Promoted, get account metadata |
| C04 | Google Gemini API | Promoted, list available models |
| C04 | Hugging Face Hub | Promoted, list references in a model repository |
| C04 | Miro | Promoted, get a specific board |
| C04 | MongoDB Atlas | Promoted, return visible projects |
| C04 | Plaid | Promoted, get processor token product permissions |
| C04 | Vercel | Promoted, retrieve project domains |
| C05 | eBay | Promoted, list opted-in seller programs |
| C05 | Etsy | Promoted, get available shop holiday preferences |
| C05 | Google Ads | Promoted, list directly accessible customers |
| C05 | Google Maps Platform | Promoted, get Places Autocomplete predictions |
| C05 | Mercado Pago | Promoted, get payment by ID |
| C05 | Razorpay | Promoted, fetch a payment by ID |
| C05 | Square | Promoted, retrieve a merchant account profile |
| C05 | Twitch | Promoted, get game and category metadata |
| C05 | WeChat Pay | Promoted, query common-mode order status |

The aggregate registry for the combined review contains 15 pinned
provider-authored artifacts and six tracked advisory artifacts. C05 used
the UWS shape table for the primary Google Ads, Google Maps Platform,
eBay, Etsy, and Mercado Pago operations, and for all five C05 endpoint
overlays. Square reproduced its selected OpenAPI selector in the bounded
operation inventory. The [C04 record](catalog-expansion-c04.md) documents
the comparable per-source bounds for that batch.
