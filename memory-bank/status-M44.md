# Status M44 - Google Discovery Service Expansion

| Item | State | Notes |
|---|---|---|
| Batch source review | `[+]` | Confirmed official Discovery document URLs, versions, and revisions for BigQuery, Cloud Storage, Docs, Slides, Chat, YouTube, Tasks, Admin SDK Directory, Analytics Data, Translate v2, Cloud Natural Language, Books, Business Profile Business Information and Account Management, and Firestore. |
| Provider/catalog rows | `[+]` | Added durable provider entries, aliases, categories, official Discovery spec references, source notes, license notes, verified dates, quirks, and user OpenAPI need classifications for the batch. |
| Review artifacts | `[+]` | Saved selected official Discovery documents under ignored `catalog-openapi-cache/google-discovery/` paths and registered them for offline review; Business Profile has two registered current Discovery artifacts. |
| Native parser coverage | `[+]` | Extended `googlediscovery` catalog artifact coverage so downloaded batch artifacts remain native Discovery-shaped and parse through the wrapper when present. |
| Security metadata | `[+]` | Added source-backed OAuth scope overlays as metadata only, without token fetching, credential resolution, account selection, or project selection. |
| Catalog/reporting checks | `[+]` | Verified provider lookup, `catalog specs`, `catalog advisory`, `catalog inspect`, `catalog stats`, `catalog security-report`, `catalog refresh-report`, and catalog quality output for the batch. |
| Documentation and memory bank | `[+]` | Product, architecture, and tech-stack boundaries already cover native Google Discovery parsing; this status file records the M44 service expansion without changing the broader contract. |
| Verification | `[+]` | `go test ./catalog`, `go test ./googlediscovery`, `go test ./...`, `go vet ./...`, `go run ./cmd/apitools catalog check`, catalog CLI smoke checks, and diff checks pass. |

## Boundary Notes

- Google Discovery documents remain native `google-discovery` source metadata,
  not official OpenAPI and not an OpenAPI import contract.
- `apitools` must not execute Google API operations, resolve credentials, fetch
  tokens, sign requests, or choose Google accounts or projects.
- `../n8n` is only a priority signal for service selection.
