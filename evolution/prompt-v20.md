# Prompt V20 - Bounded Remote Catalog Discovery

Keep APIs.guru as the primary global OpenAPI lookup. Evaluate LAP Registry as
an experimental second source, using its reported original source URL rather
than treating LAP's compressed representation as provider truth. When a caller
has a provider hostname, try the publisher's RFC 9727
`/.well-known/api-catalog` as a separate domain-scoped mechanism.

Preserve the existing public-apis compatibility fallback. Do not treat Swagger
Catalog or Scalar Registry as unauthenticated global aggregators. Keep all
remote documents untrusted, enforce existing safe-download bounds, make partial
or malformed RFC 9727 discovery fail visibly, avoid nested catalog crawling,
and leave document validation to the normal import path.
