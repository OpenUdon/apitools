# Retired milestone M51 - n8n Gap Report And Source Batch Freeze

**Milestone.** M51
**Outcome.** completed
**Retired.** 2026-09-25
**Source status.** tabilet/memory-bank/status-M51.md
**Source specification.** tabilet/memory-bank/milestone.md#m51---n8n-gap-report-and-source-batch-freeze
**Evidence.** dfb78c04118301667bed6a200c111e4489a4c437
**Worktree.** clean
**Review.** legacy
**Review iterations.** not recorded
**Legacy closure.** Closed before the persisted ten-iteration review gate was adopted; its original closure evidence is the literal status record below.
**Verification.** go test ./..., go vet ./..., go run ./cmd/cataloggen -check, go run ./cmd/apitools catalog check, and git diff --check passed at the evidence commit; this confirms current repository health, not a re-review of the original milestone diff.
**Consolidated into.** no current-truth change; current facts already live in product.md, architecture.md, and tech-stack.md.

## Milestone specification

````markdown
## M51 - n8n Gap Report And Source Batch Freeze

**Goal.** Add a repeatable offline report that compares the local n8n node
tree with built-in provider catalog coverage and freezes the next
source-review batch for first-class OpenAPI or advisory-overlay work.

**Scope.**

- Add a reusable catalog report that accepts n8n node-root names and compares
  them to provider IDs, display names, aliases, and curated category coverage.
- Add CLI access through `apitools catalog n8n-gap-report`, reading directory
  names from `../n8n/packages/nodes-base/nodes` by default.
- Classify non-covered roots as provider API candidates, M50 generic protocol
  connector exclusions, local workflow utilities, GraphQL/protocol-family
  candidates, or internal directory exclusions.
- Freeze an 8-12 service source-review batch only where provider-owned OpenAPI
  evidence or strong official docs-derived overlay evidence exists.
- Update public catalog expansion docs with the command, report interpretation,
  and frozen batch.
- Do not add provider rows, candidate rows, overlays, parser packages, n8n
  workflow import, node runtime inspection, provider API probes, credential
  resolution, host/account selection, or execution behavior.

**Acceptance.** Maintainers can run the report offline against the local n8n
node tree, see deterministic classifications and the frozen source-review
batch, and use the public queue docs to start curation without treating n8n
node presence as runtime compatibility evidence.
````

## Status record

````markdown
# Status M51 - n8n Gap Report And Source Batch Freeze

| Item | State | Notes |
|---|---|---|
| Scope freeze | `[+]` | M51 is limited to an offline report and source-review batch freeze. It does not add provider rows, candidates, overlays, parsers, or runtime behavior. |
| Report API | `[+]` | Added catalog report types and `BuildN8nNodeGapReport` for deterministic node-root classification against provider aliases/category coverage and M50 boundary maps. |
| CLI report | `[+]` | Added `apitools catalog n8n-gap-report` with `--nodes-dir` and `--json`; the command reads directory names only and skips non-directory/symlink entries. |
| Classification review | `[+]` | Running the report against `../n8n/packages/nodes-base/nodes` found 307 node roots: 225 already covered, 0 uncovered provider API candidates, 18 M50 generic protocol exclusions, 61 local workflow utilities, and 3 GraphQL/protocol-family candidates. |
| Frozen batch | `[+]` | Froze the next 10-service source-review batch: Chargebee, Mailgun, Mattermost, Paddle, Plivo, PostHog, Postmark, Rocket.Chat, Vonage, and WooCommerce. |
| Public docs | `[+]` | Updated README and `docs/catalog-expansion-queue.md` with the new command, report interpretation, and frozen batch. |
| Review and verification | `[+]` | Deep review found no blocking fixes. Passed `go test ./...`, `go test ./catalog ./cmd/apitools`, `go vet ./...`, `go run ./cmd/apitools catalog check`, `go run ./cmd/apitools search --help`, `go run ./cmd/apitools import --help`, `go run ./cmd/apitools catalog n8n-gap-report --help`, real-tree `n8n-gap-report`, `git diff --check`, and `git -C ../tofu diff --check -- apitools`. |

## Boundary Notes

- n8n node presence remains priority-only evidence.
- The report does not parse n8n node implementation files or import n8n runtime
  behavior.
- The frozen batch is a source-review queue, not a provider promotion. Each
  service still needs normal curation, artifact registration, auth/security
  review, and catalog quality checks before any durable metadata change.
- The report currently finds no uncovered provider-shaped n8n node roots after
  provider alias/category matching.
````
