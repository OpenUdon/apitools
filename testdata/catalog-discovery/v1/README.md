# Catalog discovery v1 wire fixtures

Synthetic shape examples for the additive request/report contract. All operation
metadata is repository-written test data; no model, provider or user data is used.
These are wire fixtures, not observed discovery runs or acceptance evidence for
the still-pending retrieval/ranking implementation. The tie case declares a
synthetic second operation to exercise serialization; it is not an assertion
that the prepared M81 root contains that operation. Behavioral source fixtures
and outcome checks belong to M80.2/M80.3.

Cases cover match, ambiguity after a display limit, explicit tied rank data,
intentional evidence filtering, unknown license/redistribution, no root, scoped
no-match on incompatible output typing, stale evidence, blocked exact provider
selection, and relocation (identical relative-path wire). The unknowns case
keeps the raw license note and explicitly unknown structured permission. Root
or transport installation configuration is never included in the request wire.
