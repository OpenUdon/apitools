# Result V4 - Refresh Correction Notes

Catalog refresh validation now has an explicit source-reviewed correction layer
for selected official artifacts.

- HighLevel module specs bundle reviewed common error schemas, rewrite sibling
  common-schema references, and repair small response/schema omissions.
- Strava Swagger replaces external official model references with permissive
  local object placeholders and removes an undeclared root OAuth scope for
  validation.
- Spotify Web API OpenAPI removes an extension-only external policy reference
  and prunes schema `required` mistakes that violate OpenAPI 3.0 validation.
- SyncroMSP remains invalid because its broader schema-shape issues require a
  future, separately reviewed normalization decision.

Refresh and refresh-review results expose correction notes while preserving raw
downloaded artifacts and the no-execution, no-credential, no-runtime-ref
boundary.
