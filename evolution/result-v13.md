# Result V13 - n8n Gap Report And Source Batch

M51 adds an offline n8n gap report for catalog maintainers.

The report compares local n8n node-root directory names with built-in provider
IDs, display names, aliases, and curated category coverage. Non-covered roots
are classified under the M50 boundary as generic protocol connector
exclusions, local workflow utilities, or future protocol-family candidates.

Running the report against the local n8n tree found no uncovered
provider-shaped node roots after alias/category matching. The next
source-review batch is frozen to Chargebee, Mailgun, Mattermost, Paddle,
Plivo, PostHog, Postmark, Rocket.Chat, Vonage, and WooCommerce, where
provider-owned OpenAPI or strong official docs-derived overlay evidence
justifies review.

The command reads directory names only. It does not parse n8n node
implementation files, import workflows, probe provider APIs, resolve
credentials, choose hosts/accounts, or execute operations.
