# Operation-shape wire fixture

`eight-families.json` is the exact compact UWS `uws.shape-table.v1` output
(plus one final newline) of `BuildOperationShapeTable` over the eight existing
local source fixtures in `../../operation-candidates/v1/sources/`. IDs equal
the source kinds; inputs use inline raw bytes without paths or URLs. No source
artifact or earlier candidate-wire fixture was modified.

The reviewed table contains eight raw source identities and twelve native
operations: AsyncAPI 1, Smithy 1, Discovery 1, GraphQL 2, gRPC 1, OData 3,
OpenAPI 2 and OpenRPC 1. The expected kind/protocol/input/output projections
are reviewed alongside the independent semantic/refusal tests. Unknown native
schema/transport/auth evidence remains explicit, and no operation in this
corpus claims full completeness. Both OpenAPI and OpenRPC preserve complete
supported local JSON Schema projections where available; only OpenAPI has
complete declared symbolic security evidence.

Fixture size: **9,829 bytes**. SHA-256:
`dc20d2287322a4b20d5d92b1a0d1ec8036f1bec853b642d7df8797ed096c3ba1`.
UWS dependency: `v0.0.0-20261006181058-6a267306032e` (accepted C09 source).
The standalone graph retains Horizon `v1.14.5` / HashiCorp HCL `v2.24.0`;
the new table API does not change the existing core import closure.

`TestOperationShapesDeterministicGolden` checks exact bytes, reversed source
order, UWS parse/reproduction and resolver snapshot isolation. Separate source
fixtures in `operation_shapes_wire_test.go` check exact JSON/YAML numeric
lexemes, duplicate/trailing/alias encodings, cyclic/missing references, selector
ambiguity, native collisions/overloads, stream/list/message alternatives, and
serialized/aggregate limits. These checks are offline and use synthetic data.

Regenerate only after reviewing a deliberate contract change: read the eight
local fixtures, supply the same explicit kind/ID pairs to the public producer,
call UWS `ShapeTable.Marshal`, append a newline, and independently review the
new output and semantic tests. The default test run never rewrites this file.
