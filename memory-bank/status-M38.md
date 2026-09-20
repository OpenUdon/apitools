# Status M38 - Google Discovery Converter Ownership

## Goal

Move reusable Google Discovery to OpenAPI-shaped conversion into `apitools`
while preserving the metadata-only boundary and downstream behavior.

## Status

| Item | State | Notes |
|---|---|---|
| Converter package added | `[+]` | `googlediscovery` exposes `Convert` and `ConvertMap` with Drive/Gmail fixture coverage and no `udon` import. |
| Downstream ownership split updated | `[+]` | `../udon` delegates conversion to `apitools/googlediscovery` and retains native document adaptation. |
| Verification completed | `[+]` | Ran focused/full `apitools` and `udon` test, vet, diff, CLI help, and catalog Discovery classification checks. |
| Review hardening completed | `[+]` | Follow-up review fixes preserve root-level Google Discovery parameters on operations, reject duplicate path/method collisions instead of overwriting operations, and honor non-multipart simple media upload metadata. |

## Boundary Checks

- Google Discovery conversion remains derived metadata, not official provider
  OpenAPI truth.
- Do not execute Google API operations, resolve credentials, fetch tokens, sign
  requests, or choose Google accounts.
- Keep catalog Google Discovery artifacts classified as `google-discovery` /
  structured review evidence unless an explicit caller converts them.
- Direct `../udon` adoption preserves the same derived metadata shape through
  native parsing and generator/runtime-plan lowering.
