# Result V11 - M49 n8n Public API Entry

M49 adds n8n's own public REST API as a reviewed OpenAPI provider entry.

The catalog now records the official n8n repository OpenAPI 3.0 root, public
REST API docs, ignored local review artifact registration, deployment-specific
host quirks, and complete auth metadata from the root `ApiKeyAuth` and
`BearerAuth` alternatives. The saved artifact remains review-only for strict
import until the relative `$ref` tree is bundled.

The boundary remains metadata-only. The broader n8n node directory is still
only a priority signal, and native database, broker, mail, SQL, LDAP, FTP, or
other protocol source families are not added.
