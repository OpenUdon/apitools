# Status M34 - Provider Node Expansion Batch 17

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Continue broadening provider-like n8n node catalog coverage before final backlog
reconciliation. M34 focuses on support, data, and utility services with
source-backed metadata and auth/security classification.

## Frozen Batch

Batch 17 services:

- TheHive
- TheHive Project
- APITemplate.io
- Currents
- Demio
- E-goi
- Mailcheck
- Mandrill
- Metabase
- Nextcloud
- Philips Hue
- PostBin

`../n8n/packages/nodes-base/nodes` is a priority signal only. It is not runtime
compatibility evidence, and no n8n behavior claims should be recorded.

## Per-Service Curation Loop

For each service, try official machine-readable sources first, save and
register ignored official review artifacts when retrieved, review
auth/security completeness, add candidate/provider metadata, add security
classification or overlays when needed, generate docs-derived advisory overlays
only for useful REST-shaped subsets, record provenance fields, and run focused
tests plus `go run ./cmd/apitools catalog check --as-of 2026-05-20` before the
single end-of-batch commit.

## Tasks

| Item | State | Notes |
|---|---|---|
| Batch list frozen | `[+]` | Frozen Batch 17 services are TheHive, TheHive Project, APITemplate.io, Currents, Demio, E-goi, Mailcheck, Mandrill, Metabase, Nextcloud, Philips Hue, and PostBin. |
| TheHive curated | `[+]` | Added candidate/provider rows and registered the official TheHive 4 OpenAPI artifact from the TheHive Project API documentation repository. |
| TheHive Project curated | `[+]` | Added candidate/provider rows and registered the official TheHive 5 OpenAPI artifact from StrangeBee documentation. |
| APITemplate.io curated | `[+]` | Added candidate/provider rows, X-API-KEY auth overlay, and docs-derived overlay `apitemplate-io-api-overlay.json` for template listing and PDF/image creation. |
| Currents curated | `[+]` | Added candidate/provider rows, bearer API-key auth overlay, and docs-derived overlay `currents-api-overlay.json` for projects, runs, tests, actions, signatures, instances, spec files, and webhooks. |
| Demio curated | `[+]` | Added candidate/provider rows, paired Api-Key/Api-Secret auth overlay, and docs-derived overlay `demio-api-overlay.json` for events, registrations, and reports. |
| E-goi curated | `[+]` | Added candidate/provider rows and registered the official E-goi Marketing API v3 OpenAPI artifact. |
| Mailcheck curated | `[+]` | Added candidate/provider rows and registered the official Mailcheck OpenAPI artifact from app.mailcheck.co. |
| Mandrill curated | `[+]` | Added candidate/provider rows, advisory key-placement auth overlay, and docs-derived overlay `mandrill-transactional-api-overlay.json` for messages, templates, users, and webhooks. |
| Metabase curated | `[+]` | Added candidate/provider rows, X-Metabase-Session auth overlay, and docs-derived overlay `metabase-api-overlay.json` for session, current user, databases, questions, alerts, and metrics. |
| Nextcloud curated | `[+]` | Added candidate/provider rows and registered the official Nextcloud server OpenAPI artifact; WebDAV remains protocol-specific metadata. |
| Philips Hue curated | `[+]` | Added candidate/provider rows and OAuth2 auth overlay; no endpoint overlay because the full official API reference is login-gated and public pages do not expose enough reviewed endpoint detail. |
| PostBin curated | `[+]` | Added candidate/provider rows and intentionally anonymous security classification; no endpoint overlay because stable official endpoint documentation was not found. |
| Batch quality review completed | `[+]` | Final review completed; verification rerun before the end-of-batch commit. |

## Boundary Checks

- Do not treat n8n node presence as runtime compatibility evidence.
- Do not vendor downloaded official specs, Discovery documents, or SQLite
  caches into the public repository.
- Do not execute API operations.
- Do not resolve credentials, tokens, or secrets.
- Do not sign requests or choose provider accounts.
- Do not treat advisory overlays as provider truth.
- Do not start standalone converter work in this milestone.
