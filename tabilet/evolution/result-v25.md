# Result V25 - Catalog discovery for step contracts

**State:** Approved plan, 2026-09-30; not implemented or accepted. Two
milestones are active, each with all rows pending:
[M81](../memory-bank/milestone.md#m81--catalog-discovery-foundations)
([status-M81.md](../memory-bank/status-M81.md)) runs first, then
[M80](../memory-bank/milestone.md#m80--catalog-discovery-api-for-step-contracts)
([status-M80.md](../memory-bank/status-M80.md)). No closing review gate has
started. The two approved reconciliations change planning and correct
downstream verification scope; they grant no implementation or release
authority.

The material direction change is catalog-wide retrieval for a step contract,
instead of ranking only explicitly supplied documents. The additive discovery
report is to keep provider/artifact identity, authority, license and
redistribution evidence, auth, effect, rank evidence, and searched coverage
visible. Existing APIs and wire shapes remain unchanged. Only validated,
integrity-checked artifact operations can become operation candidates; absent
bytes remain reference-only leads. License notes do not imply redistribution
permission, and untyped or missing contract evidence does not become a positive
match. Separate catalog-license/expansion/typed-overlay upgrades remain
unnumbered candidates.

M81 supplies what M80 needs at catalog scale. The existing ranking path stops
at 10,000 operations and 32 sources, while the local cache holds roughly
40,000 operations; its 20 MiB parser cap also excludes the Microsoft Graph
spec that Outlook uses. M81 therefore adds an approved design record and
consumer checkpoint, a caller-supplied catalog root contract with a synthetic
fixture root, a 128 MiB limit for registered artifacts on the index path only,
a digest-bound operation index keyed by catalog-stable identity, and an
artifact-scoped export. Embedding the committed overlays as a default root
remains a candidate direction.

M80 is to distinguish match, ambiguous, no qualifying API within the checked
scope, insufficient evidence, and blocked. It must not convert missing
artifacts, unexamined scope, or indeterminate comparisons into "impossible",
and only a scoped no-match may route to a browser step. Optional provider
constraints resolve by exact catalog key, and ties break on catalog identity,
not machine paths. Ranking remains advisory; APItools neither provisions
sources nor selects browser fallbacks, accounts, credentials, approvals, or
runtime operations. The default path is offline. Explicitly configured
APIs.guru lookup reuses guarded fetching with an eight-second total deadline,
at most three documents, and 20 MiB per document, matching iCoT's remote
lookup; tests use local servers.

The approved order is M81 then M80. The external handoff is producer
qualification and separately authorized publication -> proposed OpenUdon M94
-> proposed Kinet W10. OpenUdon M94 owns exact-pin adoption, retains its
proposed M93 prerequisite, and supplies the discovery journey needed by
proposed M95's iCoT removal gate. These sibling references remain drafts and
their files are unchanged; they must absorb the changes M81's design record
lists, including an explicitly supplied catalog root. M80 qualifies its
producer/fixtures and existing consumers without a reverse dependency on
future M94 implementation. Publication stays gated on explicit authority and
observed accepted-revision evidence.

Focused existing operation-candidate/ranking and catalog checks passed during
assessment in workspace and standalone modes; generated catalog output and
diff checks passed. These checks do not establish delivery of M81 or M80.
Future qualification includes workspace/standalone tests and vet, catalog
checks, the persisted ten-iteration review gates, and affected OpenUdon/Udon
consumer checks. Ramen remains excluded as a separate project. Frozen history,
earlier direction snapshots, and current product/architecture/lessons remain
unchanged; the superseded Ramen test requirements are retained in the
[knowledge journal](../docs/history/knowledge.md#2026-09-30--exclude-the-separate-ramen-project-from-downstream-verification).
