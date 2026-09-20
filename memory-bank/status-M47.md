# Status M47 - Dropbox Stone Source Evaluation

| Item | State | Notes |
|---|---|---|
| Source review | `[+]` | Re-reviewed the official `dropbox/dropbox-api-spec` Stone source, Dropbox HTTP docs, OAuth guide, and the M37 Stone-derived advisory overlay. The upstream `main` revision remains `6c68656caafdff65a30943fe186a1800b9f35bbc`, matching the cataloged review metadata. |
| n8n surface comparison | `[+]` | n8n covers Dropbox file/folder/search operations. The reviewed advisory subset covers account, list/list-continue, metadata, upload, download, and sharing links; missing n8n-visible operations are copy, delete, move, create-folder-v2, search-v2, and search-continue-v2. |
| Parser value assessment | `[+]` | Native Stone parsing is not justified for the current milestone: Stone is provider-specific, useful durable lowering would need route, alias, union, import, data-type, and auth attribute handling, and current catalog value comes from official source provenance plus the reviewed advisory overlay. |
| No-parser decision | `[+]` | Retain Dropbox Stone as official non-OpenAPI source metadata and advisory-overlay provenance. Do not add a native parser, general Stone lowering, or advisory overlay expansion in M47. |
| Follow-up milestone draft | `[+]` | No parser follow-up milestone was opened. If Dropbox workflow demand grows, the next scoped work should be targeted advisory overlay expansion for the missing n8n-visible operations before a general Stone parser is reconsidered. |
| Catalog reconciliation | `[+]` | Fixed provider resolution so Dropbox prefers the `dropbox-api-stone-spec` machine source over the authored human-doc references. Protocol classification, artifact paths, and overlay provenance remain unchanged. |
| Verification | `[+]` | Passed `go test ./catalog`, `go test ./...`, `go vet ./...`, `go run ./cmd/apitools catalog check`, Dropbox advisory/spec/stats smoke checks, `git diff --check`, and `git -C ../tofu diff --check -- apitools`. |

## Decision Notes

- Current reviewed advisory coverage remains intentionally limited to
  `users/get_current_account`, `files/list_folder`,
  `files/list_folder/continue`, `files/get_metadata`, `files/upload`,
  `files/download`, `sharing/create_shared_link_with_settings`, and
  `sharing/list_shared_links`.
- The n8n comparison identified useful future overlay candidates, but those
  gaps do not require a generic Stone parser before there is broader
  multi-provider value.
- Dropbox Stone remains a machine-readable source reference, not official
  OpenAPI. Stone-derived advisory metadata remains source-backed guidance for
  OpenAPI-only consumers.

## Boundary Notes

- Dropbox Stone remains official non-OpenAPI source metadata unless a future
  parser milestone explicitly changes the implementation scope.
- Stone-derived advisory metadata must not be treated as official OpenAPI.
- `apitools` must not execute Dropbox operations, resolve credentials, fetch
  tokens, or choose Dropbox teams/accounts.
