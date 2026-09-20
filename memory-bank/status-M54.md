# Status M54 - Kubernetes Cluster API Source Coverage

| Item | State | Notes |
|---|---|---|
| Source model frozen | `[+]` | Kubernetes API metadata is recorded as cluster-published Discovery/OpenAPI artifacts, not as a static global provider spec. |
| Catalog entry scoped | `[+]` | Added Kubernetes candidate/provider/security metadata from official docs for Discovery API, OpenAPI v2, OpenAPI v3, authentication, and access control. Later added a pinned `kubernetes-v1-19-2-swagger` bootstrap artifact from the HashiCorp terraform-provider-kubernetes OpenAPI test fixture. |
| Artifact handling designed | `[+]` | Exported `/api`, `/apis`, `/openapi/v2`, and `/openapi/v3` documents remain the preferred user-provided or downstream-exported local metadata; apitools performs no live discovery. The built-in Swagger artifact is a pinned fixture snapshot for offline bootstrap and corpus use, not a canonical provider-wide replacement for exact cluster metadata. |
| Parser decision recorded | `[+]` | M54 uses catalog/documentation support only. No Kubernetes-specific parser is opened because user-provided OpenAPI documents already fit existing import paths, and live discovery would cross the module boundary. |
| Security overlay inspection | `[+]` | Reconfirmed that Kubernetes has no portable provider-wide security overlay. Real credential binding needs cluster/export-specific user metadata because authentication and authorization are cluster-configured. |
| CRD boundary documented | `[+]` | Public docs and catalog quirks state that CRDs, aggregated APIs, enabled API groups, versions, and auth posture are cluster-specific. |
| Verification | `[+]` | Go, catalog, consumer, review, and diff checks passed before completing M54. |

## Source Notes

- Kubernetes publishes API resources through Discovery API endpoints such as
  `/api` and `/apis`.
- Kubernetes serves OpenAPI v2 at `/openapi/v2`.
- Kubernetes serves OpenAPI v3 through `/openapi/v3` plus per-group/version
  documents; official docs describe v3 as the preferred richer representation.
- `catalog-openapi-cache/openapi/kubernetes-v1-19-2-swagger.json` is a pinned
  Swagger 2.0 snapshot copied from the HashiCorp terraform-provider-kubernetes
  `manifest/openapi/testdata/k8s-swagger.json` fixture at commit
  `dcdf46c9ca238b671d1159f252ec19c8fe2ed16e`; it is registered as a
  non-official bootstrap artifact.
- Kubernetes Discovery includes built-in resources and CRDs, so exported source
  artifacts must be treated as exact-cluster metadata.
- Kubernetes API authentication and authorization are cluster-configured; the
  pinned fixture declares bearer-token auth, but no portable provider-wide
  security scheme should be assumed for real cluster exports.

## Security Overlay Inspection

| Provider | Built-in overlay needed? | Notes |
|---|---|---|
| Kubernetes | No portable built-in overlay | The same OpenAPI path can represent different authn/authz setups depending on cluster configuration, API server flags, authenticating proxies, service-account issuer, OIDC integration, client certificates, and RBAC/admission policy. Downstream tooling should attach a user-provided security overlay to the exact exported cluster metadata when credential binding is needed. |

## Boundary Checks

- Do not read kubeconfig files, service-account tokens, certificates, or
  cluster environment variables.
- Do not contact API servers, perform discovery, watch resources, issue dry-run
  applies, or mutate resources.
- Do not claim compatibility with a live cluster unless downstream tooling
  separately supplies and approves the exact exported metadata artifact.
- Preserve group/version/resource provenance and CRD/extension boundaries if a
  local parser is introduced.
