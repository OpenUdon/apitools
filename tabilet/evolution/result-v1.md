# Result V1 - Harness Boundary

The `apitools` private harness is owned by `../tofu/apitools` and symlinked into
the local public checkout. The public repository keeps its OpenAPI tooling
surface focused on library and CLI behavior while private planning state stays
out of public history.

The provider catalog direction belongs in `../apitools`, not `../tfconfig`.
Catalog and security-overlay behavior should be implemented as a focused
subpackage so the existing root package remains stable for current downstream
consumers.
