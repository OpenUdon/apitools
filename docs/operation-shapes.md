# Local UWS operation shapes

`BuildOperationShapeTable` is an additive source-tooling API over the public
`github.com/OpenUdon/uws/binding` contract. Supply `OperationShapeOptions` with
explicit `ShapeSourceInput` IDs and local `OperationSourceInput` bytes or files.
It returns a UWS ShapeTable, without changing the existing inventory/candidate
APIs or their wire bytes. This is M82 work in progress, not milestone acceptance
or a published consumer release.

The producer preserves raw SHA-256, family-native operation selectors and
symbolic IDs. It uses the existing native parsers. JSON numeric constraints
retain their lexemes; YAML uses guarded nodes with exact JSON-compatible numeric
scalars. Non-JSON YAML scalar tags and numeric spellings refuse. Annotations,
examples and defaults are omitted from schema projections. Missing, recursive,
external or unsupported schema evidence has `known: false`; absence is distinct
from an explicitly unconstrained schema. No external reference is loaded.

| Family | Projection and current limits |
|---|---|
| OpenAPI/Swagger | Native path/method pointer plus operation-ID alias; inherited/overridden parameters and local request/response schemas. Multiple media types or successful response alternatives remain unknown. |
| Google Discovery | Method ID/pointer, inherited native parameters and local referenced body schemas. Scope/media/server completeness remains pending. |
| AWS Smithy | Absolute operation shape ID/pointer, native AWS protocol and input-member locations. Compound/member schemas and authorization remain partial. No synthetic generic HTTP server/method is emitted. |
| AsyncAPI | Operation ID/pointer and send/receive payload direction. A single local message schema can be projected; multiple messages, bindings and auth remain partial. |
| GraphQL | Native operation/root-field pointer, variables/arguments and selected outputs. Scalar types are partial; schema selection, custom scalars, lists and auth remain unproved. |
| OpenRPC | Native method ID/pointer, ordered parameter descriptors and local result schemas. Positional parameters preserve declaration order; transport/auth remains unproved. |
| gRPC/protobuf | Native service/method pointer and request/response message identities. Message projections remain partial; streams are unknown. No RPC connection/reflection occurs. |
| OData | Native resource/action/function/import/navigation identities and declared parameters/results. Facets, response presence, resource HTTP selection and auth remain partial. |

At M82.1, operations retain `complete: false` and `security.known: false`.
M82.2 owns exact identity verification, symbolic security alternatives and
incomplete-evidence qualification. M82.3 owns deterministic wire/limit fixtures;
M82.4 owns consumer qualification, whole review and publication. A caller must
not interpret partial native type hints as a complete schema or authority.
Effects are not a ShapeTable field; this producer never infers execution effects
from HTTP methods. Existing effect/candidate evidence remains separately advisory.

Defaults and maximums are 32 sources, 20 MiB per source, 64 MiB aggregate raw
bytes, 10,000 operations, UWS 8 MiB serialized table and 256 KiB per schema.
Optional byte/operation limits only tighten these bounds. Structural guards
limit depth and semantic work. Malformed input, duplicate source IDs, unsafe
local leaves, cancellation and limits return no partial table. Parser and file
errors use the value-free `ErrOperationShapeTable`; context cancellation remains
identifiable. Cooperative parsing is not a hard CPU/RSS/deadline boundary.
Downstream hosted/untrusted parsing still requires the approved isolated worker.

Source URLs are provenance only and never fetched; userinfo/query/fragment are
removed. No credential resolution, account selection, provider call or operation
execution occurs. Source bytes and private paths are absent from the returned
table. Consumers must independently reproduce producer claims against reviewed
source bytes before deriving authority. UWS structure validation alone is not
proof of producer correctness.

Focused gate: `go test . -run '^TestOperationShapes'` and its `-race` variant.
Standalone qualification uses `GOWORK=off` and the exact published UWS C09 pin.
The UWS module retains its Horizon/HashiCorp HCL transitive dependency closure;
this API does not make UWS or APItools HCL-free.
