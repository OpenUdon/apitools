# Result V19 - Adaptive Authoring Source Discovery

M69 adds `DiscoverLocalSources` as the single bounded local discovery API for
all eight executable source families. It scans only explicit roots, validates
content with native parsers, deduplicates by SHA-256, and reports prompt-safe
metadata, provenance, invalid candidates, and ambiguous JSON/XML.

The default scan bounds are 10,000 visited entries, 100 accepted candidates,
and 20 MiB per file. Reaching either count bound visibly blocks completeness
through a truncated report and narrowing diagnostic. Symlinks, special paths,
oversized files, and cancellation fail safely.

The native source operation inventory now accepts the same family set so
downstream authoring can preserve source-specific selectors. Remote discovery,
active workflow and source selection, materialization, credentials, and runtime
execution remain downstream responsibilities.
