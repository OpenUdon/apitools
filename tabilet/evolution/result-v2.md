# Result V2 - Staged Core Catalog Plan

The catalog plan is staged across M2 through M5:

- M2 owns candidate service inventory.
- M3 owns durable catalog entries and spec references.
- M4 owns security overlays and auth/security classification.
- M5 owns catalog resolution and CLI reporting.

This keeps the first implementation step focused on evidence and prioritization
before adding catalog contracts, overlay semantics, or operator-facing commands.
The implementation boundary remains `apitools/catalog`, with no API execution,
credential resolution, signing, account selection, or n8n runtime semantics.
