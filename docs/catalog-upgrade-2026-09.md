# APItools catalog upgrade proposal: widely used services, official specs, license-clean

Discussion record: 2026-09-28.

**Status: proposal.** This is planning input for apitools' existing candidate
direction "Next provider-catalog expansion". To act on it, it enters through
`memory-bank-propose` in `apitools`. Nothing has been implemented, and no
provider was contacted. The [Kinet request-resolution note](../../kinet/docs/request-resolution.md)
references this proposal (its §9), so Kinet's stage-5 planning also plans
these APItools milestones.

**Selection basis.** Services were chosen by broad, general market use across
standard SaaS categories, not from n8n or any other integration product's node
or connector list. The M52 policy stays intact.

## 1. Summary

- The catalog already covers 316 providers:
  - 158 have an official machine-readable spec: OpenAPI or Swagger, AWS
    Smithy, Google Discovery, or Dropbox Stone.
  - 157 are covered by human docs only.
  - 1 has an OpenAPI index (HubSpot).
- The biggest remaining gaps are:
  1. widely used services that are missing entirely;
  2. big-cloud services, which the existing native source families could
     cover with no new code;
  3. entries whose best official source is now a native format that apitools
     already supports (GraphQL, protobuf), rather than human docs;
  4. license metadata: many cached specs rely on "terms apply" notes instead of
     a recorded license.
- **The proposal:**
  - **Wave 0:** a license baseline.
  - **Wave 1:** 13 new providers with official specs expected under open
    licenses.
  - **Wave 1b:** 37 cloud services in existing families.
  - **Wave 2:** 11 upgrades or explicit decisions on existing entries.
  - **Wave 3:** 30 new providers whose official sources need verifying first.
  - **Browser tier:** browsertools only as the last resort, under strict
    gates.
  - **Docs-only providers (§10):** these already produce UWS workflows through
    OpenAPI-shaped overlays built from provider docs. The work is to make
    those overlays workflow-ready: typed payloads, declared server variables,
    visible provenance, and UWS 1.12 pending steps as the signal for which
    operations to add next.

## 2. Baseline (verified at apitools `e9a555c`)

| Fact | Value |
|---|---|
| Providers | 316: OpenAPI 94, Swagger 13, Smithy 29, Google Discovery 21, Dropbox Stone 1, OpenAPI index 1, human docs 157 |
| Human-docs providers with a docs-derived overlay | 139 of 157 |
| Human-docs providers with no machine-readable artifact | 21: acumatica, airwallex, amplitude, anthropic, ashby, braze, emelia, linear, monday-com, netsuite, odoo, oracle-fusion-cloud-applications, orbit, philips-hue, postbin, postman, sap-s4hana, sap-successfactors, segment, uproc, workday |
| License notes on machine-readable references | 57 name an explicit open license (27 OpenAPI, 29 Smithy, 1 Stone). 124 say only "…terms/repository license apply" (101 OpenAPI, 22 Discovery, 1 index). |
| Redistribution today | Specs download into a local refresh cache that git ignores. Only fixtures are committed, such as `kubernetes-v1-19-2-swagger.json` (taken from HashiCorp's terraform-provider-kubernetes, not Kubernetes itself) and `nvidia-dsx-air-openapi.yaml`. The 141 advisory overlays are committed. |
| Refresh health | 181 refreshable specs, 19 missing registrations, 9 invalid artifacts |
| Unpromoted candidate | `google-ads`: M48 parked it until a protobuf source family existed. That family now exists. |
| Next free catalog-lane milestone ID | C04 |

## 3. Principles

### 3.1 How services are selected (no n8n)

1. **Category coverage** of a typical SMB/SaaS automation stack: payments,
   commerce, CRM, marketing, support, communications, productivity, dev and
   cloud, data, AI, HR, finance, maps and social.
2. **Public, license-neutral usage facts**, measured in the first row of each
   batch:
   - download counts of the *provider's official SDK* from public npm or PyPI
     registry statistics;
   - scale figures the provider publishes itself;
   - OpenUdon user demand.

   Record these as numbers only. No third-party integration code, connector
   metadata or node lists is copied or used as evidence.
3. The ordering below is a judgment of broad market usage. Row 1 of each batch
   should confirm or reorder it with the measured facts.

### 3.2 Source ladder

This is applied per provider, and per operation where needed.

| Level | Source | Catalog treatment |
|---|---|---|
| **L1** | Official machine-readable spec under an explicit open license (MIT, Apache-2.0, BSD, CC0, MPL-2.0) | Catalog reference plus digest. Refreshed into the local cache. May be committed as a fixture, keeping the license notice. |
| **L2** | Official machine-readable spec under provider terms only | Catalog reference plus digest. Local cache only; **never committed**. `materialize` and `export` must show the terms in the provenance manifest. |
| **L3** | Official human docs only | The existing docs-derived advisory-overlay process (M61 linked-docs rule). It carries **facts only**: method, path, parameters, auth. No copied prose. |
| **L4** | No usable official API for the task | A **browsertools** browser-capability profile, as the last resort (see §8). |

**Always excluded:**
- undocumented or private endpoints;
- third-party integration code or metadata (n8n, Zapier, Make and the like);
- unofficial community specs as catalog sources;
- automating any service whose terms prohibit automated access.

### 3.3 Precedence

An official source always outranks derived material. A docs overlay never
replaces an official spec. It may only fill an operation gap that is recorded,
with a provider-doc citation.

## 4. Wave 0: license-clean baseline (do this first)

| # | Item |
|---|---|
| 0.1 | Add structured license fields to `spec_references`: `license_spdx` (or `provider-terms`), `license_basis` (repository license, document metadata, or provider terms), and `redistribution` (permitted, reference-only, or unknown). Keep `license_note` for human-readable detail. |
| 0.2 | Reclassify the 124 "terms only" references. Several are probably open once checked; for example, the Microsoft Graph metadata and Azure REST API specs repositories. Others must stay reference-only; for example, the Grafana repository is AGPL-3.0. |
| 0.3 | Replace the committed Kubernetes fixture (a third-party HashiCorp snapshot, MPL-2.0) with the official `kubernetes/kubernetes` `api/openapi-spec` (Apache-2.0). Record the license and notice of `nvidia-dsx-air-openapi.yaml`, the other committed fixture. |
| 0.4 | Audit the 141 committed advisory overlays: they should hold facts and original wording, with no provider prose copied verbatim. |
| 0.5 | Fix refresh health: 19 missing registrations and 9 invalid artifacts. |
| 0.6 | Make `catalog materialize` and `catalog export` carry the license and redistribution fields into the workflow-directory provenance manifest, because OpenUdon packages redistribute those artifacts. |

## 5. Wave 1: widely used services missing from the catalog, with official specs expected under open licenses (L1)

All 13 are new, except `google-ads`, which is promoted from its existing
candidate record.

| # | Service | Category | Expected official source | Format | Expected license | Confidence |
|---|---|---|---|---|---|---|
| 1 | DigitalOcean | cloud | `digitalocean/openapi` | OpenAPI 3 | Apache-2.0 | high |
| 2 | Google Maps Platform | maps/geo | `googlemaps/openapi-specification` | OpenAPI 3 | Apache-2.0 | high |
| 3 | Plaid | banking data | `plaid/plaid-openapi` | OpenAPI 3 | MIT | high |
| 4 | Amazon Selling Partner API | marketplace commerce | `amzn/selling-partner-api-models` | Swagger/OpenAPI per API | Apache-2.0 | high |
| 5 | Algolia | search | `algolia/api-clients-automation` (its `specs/`) | OpenAPI 3 | MIT | high |
| 6 | Datadog | observability | the OpenAPI schemas shipped in the `DataDog/datadog-api-client-*` repositories (v1, v2) | OpenAPI 3 | Apache-2.0 | high |
| 7 | Google Gemini API | AI | Google Discovery document for `generativelanguage` | Discovery | Google Developers Site Policies (same basis as the 22 existing Discovery docs) | high |
| 8 | Google Ads API (promote the `google-ads` candidate) | advertising | protobuf definitions in `googleapis/googleapis` | gRPC/protobuf | Apache-2.0 | high (large surface: pick a scope) |
| 9 | MongoDB Atlas Administration API | database | `mongodb/openapi` | OpenAPI 3 | Apache-2.0 | medium: verify |
| 10 | Square | payments/commerce | Square's official API-spec repository on GitHub (verify its current location) | OpenAPI/Swagger | Apache-2.0 expected | medium: verify |
| 11 | Pinterest | social/ads | `pinterest/api-description` | OpenAPI 3 | verify | medium: verify |
| 12 | Dropbox Sign (HelloSign) | e-signature | `hellosign/hellosign-openapi` | OpenAPI 3 | MIT expected | medium: verify |
| 13 | AssemblyAI | speech AI | `AssemblyAI/assemblyai-api-spec` | OpenAPI | MIT expected | medium: verify |

## 6. Wave 1b: big-cloud services in existing native families (L1 or L2, no new code)

These follow the M60 precedent: official models only, and auth metadata is
advisory.

| Family | Official source | Services to add |
|---|---|---|
| AWS Smithy (15) | `aws/api-models-aws` (Apache-2.0) | ECS, EKS, ECR, Route 53, CloudFormation, Step Functions, EventBridge, CloudFront, Systems Manager, CloudWatch Logs, STS, Organizations, Bedrock Runtime, Cost Explorer, SES v2 |
| Azure OpenAPI (10) | `Azure/azure-rest-api-specs` (MIT expected; see 0.2) | Compute, Network, Resources (resource groups and deployments), Key Vault (data plane), Monitor, App Service, AKS (ContainerService), Azure OpenAI inference (data plane), Service Bus, Event Grid |
| Google Discovery (12) | Google Discovery service | Compute Engine, Cloud Run, Pub/Sub, Cloud Functions, Secret Manager, Cloud Logging, Cloud Monitoring, IAM, Resource Manager, GKE (container), Cloud SQL Admin, Google Forms |

## 7. Wave 2: upgrade or settle existing entries (verify first)

Earlier catalog decisions (M12, M46, M59) are respected. An entry is re-opened
only where a *native* format that apitools now supports offers a better
official source.

| Entry | Now | Proposed target | Expected result |
|---|---|---|---|
| Linear | human docs (GraphQL, per its catalog quirk) | official GraphQL SDL published with Linear's SDK (MIT expected) | native `graphql`, L1 |
| monday.com | human docs (GraphQL) | official GraphQL schema, if downloadable | native `graphql`, L1/L2, else stays L3 |
| Shopify | human docs (REST Admin) | Admin **GraphQL** schema, if an official public schema exists (REST Admin is the legacy surface) | native `graphql`, L2, else stays L3 |
| Okta | minimal Management subset | a broader subset of Okta's official OpenAPI tree, chosen per workflow | L1, larger |
| Twilio | core v2010 only | add official product specs from the same repository: Messaging, Verify, Conversations | L1, more specs |
| Zendesk | Sunshine Conversations only | Support/Ticketing API: an official downloadable spec if present, else a docs overlay | L2 or L3 |
| Slack | archived official repository | freshness check; add a docs overlay for methods newer than the archive | L1 plus L3 |
| Salesforce | human docs | explicit decision: org-generated OpenAPI supplied by the user, plus a docs overlay | L3 (record the decision) |
| Anthropic | human docs (M59: none found) | re-check only whether an official spec is now published | L2 or stays L3 |
| Mailchimp | human docs (M12: none found) | re-check the Marketing API spec status | L2 or stays L3 |
| Amplitude, Ashby, Braze, Segment | human docs, no artifact | build overlays from provider docs, finishing the M59 intent | L3 |

## 8. Wave 3: widely used services missing from the catalog, official sources to verify

The target is L1 or L2. The fallback is an L3 overlay; the browser is never
the first choice.

| Category | Services |
|---|---|
| AI and ML | Cohere, ElevenLabs, Replicate, Hugging Face (Inference), Pinecone, Deepgram, Groq |
| Productivity and collaboration | Miro, Smartsheet, Wrike, Basecamp |
| Commerce and finance | Etsy, eBay, Brex, Ramp |
| Support and communications | Front, RingCentral |
| Dev, cloud and security | Vercel, Snyk, 1Password (Connect), Heroku (its platform schema is JSON hyper-schema, a format apitools doesn't parse, so expect L3) |
| HR and recruiting | Gusto, Deel, Rippling, Personio, Greenhouse, Lever |
| Analytics, email and media | Mixpanel, Resend, Vimeo |
| Upgrades of existing overlays | Mistral AI, Sentry: check for official specs |

For each service, record:
- the official source found, or "none found";
- the level: L1, L2, L3 or L4;
- the structured license fields;
- the auth/security review (M58 rule);
- the digest and refresh registration;
- `apitools catalog check`.

## 9. Browser tier: browsertools as the last resort (L4)

Use a browsertools browser-capability profile only when **all** of these hold:

1. No L1–L3 source covers the task. That is recorded as a catalog decision,
   for example "UI-only: no public API for <task>".
2. The service's terms allow automated operation of the user's own account.
   For example, LinkedIn's User Agreement prohibits bots and scraping, so
   LinkedIn member actions never qualify.
3. The action is narrow and reviewed. It follows browsertools' existing
   safety model: exact origins, confirmation for side effects, no
   credentials or captured values in the profile, and human-only CAPTCHA
   handling.
4. The profile is labeled as weaker evidence than an API source. It is
   re-verified when the UI changes, and it is retired as soon as an official
   API appears.

**Typical fits:**
- UI-only administrative tasks in tenant consoles;
- small SaaS tools with no API;
- supplier or public-sector portals the user operates for their own account,
  where the terms allow it.

**Ownership:** these profiles live in browsertools and OpenUdon, not in
apitools. apitools only records the "no API for this task" decision, so that
the workflow author can route to browsertools. Today that author is OpenUdon's
iCoT; after the kinet migration it is kinet's W loop (see [Kinet request-resolution note](../../kinet/docs/request-resolution.md),
gap G5). That routing would be a separate cross-repo proposal.

## 10. Building UWS workflows from docs-only providers

### 10.1 No new UWS feature is needed

UWS 1.5 added the `browser-profile` source type for website UIs. It did not
add a human-docs source type, and none should exist: a UWS operation binds to
a machine-readable source (UWS §4.4).

The bridge is apitools' docs-derived advisory overlays. Each one is a real
OpenAPI 3.0.3 document, with `paths`, `servers` and `securitySchemes`, and it
records its provenance in `x-apitools-overlay`: `derived_from_docs: true`,
`official_openapi: false`, and `source_refs` pointing to the provider docs. A
workflow binds to it as an ordinary `openapi` source.

**The pipeline already works end to end:**
1. **apitools:** 136 of the 157 docs-only providers have an overlay.
   `catalog materialize` copies the overlay and its security-overlay sidecar
   into the workflow directory.
2. **OpenUdon iCoT:** treats `advisory-overlay` as an OpenAPI source candidate
   (`internal/icot/elicitor/catalog.go:292`, `catalog_plan.go:471`).
3. **The UWS document:** binds operations to overlay `operationId`s.
4. **udon:** executes over HTTP. It fills server placeholders such as
   `{domain}` from its runtime provider settings
   (`soliton.Provider.ResolveServerURLVariables`), not from the UWS document.

**Note: the ongoing kinet migration.** Step 2 is the legacy iCoT path.
Under [Kinet `docs/icot.md`](../../kinet/docs/icot.md), iCoT retires and kinet's W loop takes over
authoring through OpenUdon's non-interactive step commands (`step
candidates|bind|check|source add`). Those commands currently accept only
local source files the user names. So the overlays reach the new path only
once catalog discovery exists there: gap G1 in [Kinet request-resolution note](../../kinet/docs/request-resolution.md), which that
note recommends as a precondition for retiring iCoT.

```yaml
uws: "1.12.0"
info: {title: Freshdesk escalation, version: "1.0.0"}
sourceDescriptions:
  - name: freshdesk
    type: openapi
    url: ./sources/freshdesk-api-v2-overlay.json   # docs-derived, not official
operations:
  - operationId: create_ticket
    sourceDescription: freshdesk
    sourceOperationId: createFreshdeskTicket      # overlay operationId
    effect: write                                  # UWS 1.12 effect classification
    request:
      body: {subject: $inputs.subject, email: $inputs.email, priority: 1, status: 2}
```

### 10.2 Where it falls short

Measured across the 139 overlay files at `e9a555c`:

- **Payloads are almost never typed.**
  - Only 11 of 443 request bodies and 20 of 1,309 success responses have
    typed schemas.
  - In the example above, `subject`, `email`, `priority` and `status` come
    from Freshdesk's docs, not from the overlay, whose request body is an
    untyped object.
  - Nothing can check `$response.body.<path>` outputs statically, and
    content-trust resolvers have no field-level shape to work with.
- **Coverage is thin.** There are 1,328 operations in total, a median of 8 per
  provider (from 1 to 42). Many workflows need an operation the overlay
  doesn't have.
- **Server placeholders are often undeclared.**
  - 46 overlays use a templated server URL, and 27 of them do not declare its
    `variables`: for example `https://{site}`, `https://{wordpress_host}/…`
    and `https://{shop}.myshopify.com/…`.
  - Execution still works, because udon fills placeholders from its runtime
    settings. But authoring cannot tell which deployment values to ask for.
  - Telegram (`bot{token}`) and Plivo (`/Account/{auth_id}`) put a
    credential in the server URL. That value must come from credential
    binding, never from UWS request values.
- **Provenance is invisible in the UWS document.** The `sourceDescription` only
  points to a file. The "derived, not official" status lives inside the
  overlay.
- **21 docs-only providers have no artifact at all** (listed in §2).

### 10.3 Proposed improvements

These are apitools changes unless marked otherwise. UWS core stays unchanged.

| # | Improvement | Owner |
|---|---|---|
| D1 | **Typed workflow slices.** For the operations that OpenUdon examples and users actually call, add request and response schemas taken from the provider docs: field names, types, required flags and enums, as facts only. Prioritize by demand (D2), not by covering all 1,328 operations. | apitools |
| D2 | **Pending steps as the demand loop.** When an operation is missing, the author writes a UWS 1.12 `pending` step with `purpose`, typed `inputs` and `outputs` schemas, and `effect`. The draft passes ordinary validation, while executable validation blocks it. apitools adds the operation to the overlay from the provider docs (M61 linked-docs rule), and the step then becomes an ordinary `operationRef` step. | OpenUdon (authoring), apitools (overlay rows) |
| D3 | **Declare every server variable.** Add a description and default for each of the 27 undeclared templates, and mark credential-bearing variables (for example a proposed `x-apitools-credential-variable: true`) so that OpenUdon routes them to credential binding and asks for the other deployment values. | apitools |
| D4 | **Carry provenance into UWS.** A proposed convention: when materializing, copy `x-apitools-source: {authority: derived-docs, overlay_id, digest}` onto `sourceDescriptions[]`. Extensions are already allowed there. Set `effect` on operations backed by overlay writes (apitools `operation_effect` can suggest a value), so review flags derived writes. | apitools (materialize), OpenUdon (review display) |
| D5 | **Evidence of working operations.** Record sandbox proof-run results, under OpenUdon's existing proof-run policy, as evidence that a derived operation works. The results store status, operation and digest, but no response values. | OpenUdon |
| D6 | **The 21 providers without artifacts, by kind.** Linear and monday.com: their GraphQL schemas through UWS's `graphql` source type (Wave 2). Tenant ERPs (SAP S/4HANA, SAP SuccessFactors, Oracle Fusion, NetSuite, Workday, Acumatica, Odoo): specs exported from the customer's own tenant, following the Kubernetes pattern of user- or cluster-exported artifacts from M54, or else an overlay. Everything else: a new overlay. browsertools only for tasks with no API (§9). | apitools |
| D7 | **Replace overlays when official specs appear.** Per §3.3, an official spec supersedes the overlay (Wave 2), and the overlay is kept only for operations the official spec omits. | apitools |

**Acceptance for a "workflow-ready" overlay:**
- every operation used by a published OpenUdon example has typed request and
  response schemas with doc citations;
- every server template variable is declared, and credential-bearing ones are
  marked;
- `x-apitools-source` provenance appears in the materialized UWS package;
- at least one sandbox proof-run record exists (D5), or the absence of one is
  recorded explicitly;
- `apitools catalog check` passes.

## 11. Suggested milestone shape (apitools, catalog lane)

| ID | Milestone | Scope | Depends on |
|---|---|---|---|
| C04 | License-clean catalog baseline | Wave 0 (0.1–0.6) | none |
| C05 | Official-spec provider batch | Wave 1 (13 providers) | C04 |
| C06 | Native cloud family expansion | Wave 1b (37 services) | C04 |
| C07 | Native-format upgrades and settled decisions | Wave 2 (11 entries) | C04 |
| C08 | Verify-and-add batch | Wave 3 (30 services), each with a recorded L1–L4 decision | C05 |
| C09 | Workflow-ready docs overlays | §10.3 D1, D3, D4 (materialize side), D6 and D7 for the overlays OpenUdon uses first | C04 |

Keep two items as unnumbered candidate directions until OpenUdon agrees on the
contracts:
- "UI-only routing to browsertools";
- the pending-step demand loop and proof-run evidence (§10.3 D2 and D5, plus
  the OpenUdon side of D4).

**Reconcile with [catalog-expansion-queue.md](catalog-expansion-queue.md) before
promotion:**
- The queue prefers batches of 8–12 services. Split Wave 1 (13 providers) into
  two batches, and split Wave 1b (37 services) into batches per source family.
- The queue says GraphQL-first services should wait for GraphQL support.
  APItools now has a GraphQL source family, so the Wave 2 GraphQL moves
  (Linear, monday.com, Shopify) no longer need to wait. The queue text should
  be updated when this proposal is promoted.

**Acceptance for each provider:**
- the official source is cited with `source_authority` `official-*`;
- the structured license fields are set;
- the auth and security metadata is reviewed;
- the digest and refresh registration are in place;
- `apitools catalog check` and the existing tests pass;
- where the result is L3 or L4, a no-spec or UI-only decision is recorded.

**Batch-level acceptance:**
- the source-ladder distribution and the license distribution are reported;
- no committed artifact lacks a redistribution basis.

**Metrics to track:**
- share of providers with an official machine-readable source (now 158 of 316,
  about 50%);
- share of machine-readable references with an explicit license (now 57 of
  181);
- count of categories covered;
- count of refresh failures;
- share of typed request bodies and success responses in overlays (now 11 of
  443 and 20 of 1,309);
- count of undeclared server-template variables (now 27 overlays).

### 11.1 Mapping to kinet's stages

kinet plans six stages. Stage 1 is in progress, and no milestones are planned
for later stages. kinet recommends stage 5 next, then 2, 3, 4 and 6
([Kinet `docs/ideas.md`](../../kinet/docs/ideas.md) §10). See [Kinet request-resolution note](../../kinet/docs/request-resolution.md) §7.1 for the full mapping
of the resolution work.

| Item | Stage | Reason |
|---|---|---|
| C04–C08 (license baseline and catalog waves) | **No kinet stage: apitools' own catalog lane** | kinet-order allows sibling preparation during stage 1. M79 completed and retired after OpenUdon M89 consumed its exact published contract; see the [M79 history record](../tabilet/docs/history/status-M79.md). |
| C09 (workflow-ready docs overlays) | **No kinet stage, but prioritize it before kinet stage 5** | It raises candidate quality for stage 5's step-contract ranking and catalog discovery. |
| Catalog discovery API for step contracts ([Kinet request-resolution note](../../kinet/docs/request-resolution.md) G1) | **Stage 5, first slice**, with the contract designed during stage 1 | It must exist before iCoT retires (S3). This is apitools' part of a proposed kinet sibling item S2d. |
| UI-only routing to browsertools (§9) | **Stage 5** | Needs S2b's supervised browser acquisition. |
| Pending-step demand loop (§10.3 D2) | **Stage 5** | Placeholders inside the package are a stage-5 outcome. |
| Proof-run evidence (§10.3 D5) | **Tier 1 in stage 5; live runs in stage 4** | Real calls need kinet's stage-4 real execution. |

## 12. Caveats

- Repository names marked "high" confidence are well-known official
  locations. Everything else is marked "verify", and nothing here was fetched
  or checked online. Batch row 1 must confirm location, license, freshness and
  format before a provider is added.
- License classifications are engineering judgments, not legal advice. L2
  (provider-terms) artifacts stay reference-only unless their terms clearly
  permit redistribution.
- The ordering reflects broad market usage by category. Measured SDK download
  counts and OpenUdon demand should confirm or reorder it.
