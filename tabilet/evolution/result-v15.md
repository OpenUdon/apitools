# Result V15 - Post-M60 Release Remediation

M60 completed the Workato-signaled AWS Smithy expansion and left `apitools` in
a post-v0.1.0 remediation phase.

The active direction remains provider-source-first catalog curation, native
Smithy and Google Discovery parsing through standalone modules, and
metadata-only OpenAPI/Swagger audit/reporting. Catalog stats remains a CLI
maintainer overview rather than a public Go report API.

Security audit follow-ups now distinguish upstream artifact defects from
reviewed overlay coverage: a missing or inconsistent artifact row can point to
an existing source-backed overlay when that overlay targets the same spec ref
or cites the artifact's official spec URL. NocoDB v2 and Toggl Track gained
review overlays while preserving present-incomplete status because their
upstream auth metadata still needs normalization or source review.
