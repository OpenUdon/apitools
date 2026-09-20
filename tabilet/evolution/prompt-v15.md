# Prompt V15 - Post-M60 Release Remediation

Capture the post-M60 / v0.1.0 state and tighten release-readiness follow-ups
without retagging the release.

Keep catalog stats CLI-only for now, document that boundary clearly, and keep
security-audit remediation metadata-only. When local OpenAPI/Swagger artifacts
lack or misstate security metadata, source-backed provider overlays may satisfy
the catalog review follow-up only when they target the exact spec reference or
cite the artifact's official spec URL.

Do not reintroduce runtime credential resolution, API execution, Smithy-to-
OpenAPI conversion, or Discovery-to-OpenAPI conversion. Deprecated
compatibility wrappers can be removed only after downstream consumers have
migrated to the standalone parser modules and the compatibility break is
documented.
