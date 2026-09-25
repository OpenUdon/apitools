# Retired milestone M53 - Docker API Source Coverage

**Milestone.** M53
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M53.md
**Source specification.** tabilet/memory-bank/milestone.md#m53---docker-api-source-coverage
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M53 - Docker API Source Coverage

**Goal.** Add Docker-owned API source coverage for high-value container
workflows while keeping daemon and registry behavior metadata-only.

**Scope.**

- Add Docker Hub as a durable provider catalog entry using Docker's official
  Hub API reference and downloadable Swagger/OpenAPI specification.
- Add Docker Engine API metadata as a local daemon/control-plane API source
  using Docker's official Engine API OpenAPI references by version.
- Review Docker Registry HTTP API source availability and either add a scoped
  provider/protocol row or record a no-promotion decision if no durable
  official OpenAPI artifact is available.
- Add auth/security classifications and overlays as metadata only; do not log
  in to Docker Hub, contact a registry, connect to a Docker socket, pull
  images, inspect containers, or execute daemon operations.
- Update public docs and catalog queue notes where Docker control-plane APIs
  need a local-daemon safety distinction from SaaS provider APIs.

**Acceptance.** Docker Hub and Docker Engine appear in catalog inspection,
spec, stats, advisory, and security views with official source refs and clear
metadata-only boundaries; Docker Registry is either represented from official
source evidence or explicitly deferred; Go, catalog, consumer, and diff checks
pass.
````

## Status record

````markdown
# Status M53 - Docker API Source Coverage

| Item | State | Notes |
|---|---|---|
| Source freeze | `[+]` | Frozen from Docker-owned docs/specs only: Docker Hub latest OpenAPI, Docker Engine API v1.54 Swagger/OpenAPI, and Docker Hub-supported Registry API OpenAPI. |
| Docker Hub cataloged | `[+]` | Added provider/candidate/spec/auth metadata from Docker's official Hub API reference and downloadable Swagger/OpenAPI artifact. |
| Docker Engine cataloged | `[+]` | Added Docker Engine API v1.54 metadata from official Docker Engine OpenAPI references, with local-daemon/control-plane safety notes. |
| Docker Registry reviewed | `[+]` | Added scoped Docker Registry metadata from Docker's official Docker Hub-supported Registry API OpenAPI; recorded that it is not the full OCI Distribution Specification or every registry implementation. |
| Security overlay inspection | `[+]` | Inspected official Docker Hub, Engine, and Registry OpenAPI security metadata. Docker Hub does not need a built-in security overlay; Docker Engine needs deployment/user-specific auth metadata; Docker Registry now has a reviewed bearer-placement overlay while token-service discovery remains downstream. |
| Safety boundary documented | `[+]` | Public docs and catalog quirks state that `apitools` never logs in to Docker Hub, connects to registries, opens Docker sockets, pulls images, inspects containers, or executes daemon operations. |
| Verification | `[+]` | Targeted catalog tests passed; full Go, catalog, consumer, text-audit, and diff checks are run before the M53 commit. |

## Source Notes

- Docker Hub API: official Docker docs advertise reference documentation and a
  downloadable Swagger/OpenAPI specification.
- Docker Engine API: official Docker docs advertise versioned Engine API
  reference documentation and downloadable Swagger/OpenAPI specifications.
- Docker Registry API: Docker Docs publishes a narrower OpenAPI document for
  the Docker Hub-supported subset of Registry HTTP API V2; generic OCI
  Distribution compatibility remains out of scope for this milestone.

## Security Overlay Inspection

| Provider | Built-in overlay needed? | Notes |
|---|---|---|
| Docker Hub | No | Official Docker Hub OpenAPI defines `bearerAuth` and `bearerSCIMAuth`; every operation carries a security requirement or an explicit anonymous override. Classification can stay `complete` without a supplemental overlay. |
| Docker Engine | No portable built-in overlay | Official Engine Swagger has no security definitions. Access is deployment-specific through local socket permissions, SSH/TCP transport, TLS client certificates, or daemon policy, so downstream/user security metadata is required for real credential binding. |
| Docker Registry | Complete, present-incomplete status retained | Implemented `docker-registry-bearer-auth-overlay` for the official Docker Hub-supported Registry OpenAPI. The overlay records `Authorization: Bearer` placement for protected endpoints, but the provider remains `present-incomplete` because WWW-Authenticate challenge handling, token-service discovery, repository scopes, and credential resolution are runtime/downstream concerns. |

## Boundary Checks

- Treat Docker Engine as local control-plane metadata, not a SaaS provider.
- Do not connect to `/var/run/docker.sock`, TCP daemon sockets, Docker Hub,
  registries, or credential helpers.
- Do not pull, push, tag, build, inspect, start, stop, or delete images,
  containers, networks, volumes, or registry artifacts.
````
