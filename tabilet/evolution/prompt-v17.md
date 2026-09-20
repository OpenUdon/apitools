# Prompt V17 - Portable Fnct Helper Boundary

Expose a narrow public helper catalog for pure `fnct` payload-shaping
functions, using Gmail raw-message rendering as the first case.

The helper surface should let OpenUdon author/review a stable function name and
request/output contract while letting udon import and register the actual Go
helper in its trusted runtime. Keep `apitools` out of workflow execution,
provider API calls, credential resolution, token fetching, and account
selection.
