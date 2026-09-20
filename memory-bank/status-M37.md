# Status M37 - Dropbox Stone-Derived OpenAPI Advisory Artifact

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

## Intent

Formalize the existing Dropbox core advisory overlay as a reviewed
Stone-derived OpenAPI-shaped subset without treating Dropbox's official Stone
source as official OpenAPI or adding runtime protocol support.

## Tasks

| Item | State | Notes |
|---|---|---|
| Dropbox source boundary retained | `[+]` | `dropbox-api-stone-spec` remains the official machine-readable source, official OpenAPI remains unavailable, and advisory tests assert protocol `dropbox-stone`. |
| Dropbox advisory overlay metadata tightened | `[+]` | Regenerated the existing overlay with explicit non-official OpenAPI, official machine-spec derivation, `dropbox-stone` source protocol, and reviewed source refs. |
| Artifact registry provenance updated | `[+]` | Registry generation now records machine-spec-derived advisory overlay metadata for Dropbox while keeping artifact kind `advisory-overlay`. |
| Advisory wording and tests reconciled | `[+]` | Advisory endpoint overlay wording now allows non-human-doc derivation, and tests preserve registered overlay metadata. |
| Verification and review completed | `[+]` | Ran the requested test, vet, diff, catalog check, advisory, inspect, and stats commands; cache-backed Dropbox advisory output shows the overlay only as `advisory-overlay` metadata with `source_protocol=dropbox-stone` and `official_openapi=false`. |

## Boundary Checks

- Do not classify Dropbox as official OpenAPI.
- Do not build a general Stone-to-OpenAPI converter.
- Do not execute Dropbox API operations.
- Do not resolve credentials, fetch tokens, sign requests, or choose runtime
  accounts.
- Do not change `../udon` or add runtime protocol support.
