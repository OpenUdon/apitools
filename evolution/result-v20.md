# Result V20 - Bounded Remote Catalog Discovery

M70 keeps APIs.guru as the primary global OpenAPI directory and adds the LAP
Registry as an experimental secondary lookup. LAP candidates retain the
registry-reported original source URL and are explicitly unvalidated until the
existing import path downloads and validates the selected document.

Callers with a provider URL or hostname can also request one RFC 9727
`/.well-known/api-catalog`. The adapter requires Linkset JSON, extracts
OpenAPI-like `service-desc` targets, resolves relative links against the final
catalog URL, rejects unsafe hosts, deduplicates targets, enforces a 100-link
default, and never follows nested catalogs. Both new sources expose provenance
and experimental trust metadata through the Go API, CLI, and cache.

Automatic search order is APIs.guru, LAP Registry, provider-scoped RFC 9727
when configured, then the legacy public-apis probe. Swagger Catalog and Scalar
Registry are not queried as unauthenticated global aggregators. No discovered
operation is executed and no credentials are resolved.
