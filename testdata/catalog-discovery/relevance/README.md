# Repository-owned relevance baseline

These three synthetic OpenAPI documents are copied byte-for-byte from OpenUdon
`examples/eval/m28-incident-drive-slack/openapi/{jira,slack,drive}.yaml` at
`d2c4f157db47c81eb80c31bcd9ae0b3770932d23`. They describe throwaway `.test`
services, contain no user data or credential values, and are never executed.

| Step shape / purpose | Required string inputs | Relevant operation | Observed discovery |
| --- | --- | --- | --- |
| Create an incident issue | summary, description, severity | createIssue | insufficient evidence: `key` response omitted by existing credential filtering |
| Post an incident alert to Slack | channel, text | postMessage | match |
| Upload an incident report file | name, parentId, content | uploadFile | insufficient evidence: existing effect classifier reports unknown |

All three contracts request write effects. Open retrieval evaluates all three
artifacts together, using temporary digest-bound registrations and a real M81
index. `TestCatalogDiscoveryEvalRelevanceBaseline` records precision@1: the
relevant operation ranked first for 3/3 queries (1.000). This small baseline is
not an acceptance threshold or a claim of semantic equivalence. Outcome checks
are independent of rank and preserve the two unsupported cases.

Named-provider parity fixtures reuse CatalogPlan's explicit provider-selection
shapes (weather/Gmail and Jira/Slack/Google Drive). Exact catalog keys select the
same canonical provider identities without splitting names or borrowing model
planning, source migration, fallback, step sequencing or execution authority.
M81's source-backed notes fixture separately verifies native-reference export.
