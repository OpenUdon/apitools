# Prompt V23 - Operation Lifecycle Ranking Ownership

Move conservative API lifecycle sibling ranking out of the generic Authoring
engine and into apitools, where `OperationSummary` and source provenance are
owned. Preserve reviewed scoring and ambiguity behavior, name the scoring
constants, and normalize Google Discovery `/upload` paths only when explicit
provenance is present. Migrate OpenUdon and Ramen without moving workflow,
desired-state, credential, account, or execution semantics into apitools.
