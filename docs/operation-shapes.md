# Local UWS operation shapes

`BuildOperationShapeTable` is an additive source-tooling API over the public
`github.com/OpenUdon/uws/binding` contract. Supply `OperationShapeOptions` with
explicit `ShapeSourceInput` IDs and local `OperationSourceInput` bytes or files.
It returns a UWS ShapeTable, without changing the existing inventory/candidate
APIs or their wire bytes. M82 is accepted after whole review 3 at public source
`54583f9b2f452b7cc522360c5aeeff29ca22f96c`;
`v0.0.0-20261006210844-54583f9b2f45` resolves to that exact revision.

The producer preserves raw SHA-256, family-native operation selectors and
symbolic IDs. It uses the existing native parsers. Direct JSON Schema numeric
constraints retain their lexemes; YAML uses guarded nodes with exact JSON-compatible numeric
scalars. Non-JSON YAML scalar tags and numeric spellings refuse. Annotations,
examples and defaults are omitted from schema projections. Missing, recursive,
external or unsupported schema evidence has `known: false`; absence is distinct
from an explicitly unconstrained schema. No external reference is loaded.
Native-parser normalized subsets retain explicit unknownness rather than
claiming complete schema or lossless numeric semantics.

| Family | Projection and current limits |
|---|---|
| OpenAPI/Swagger | Native path/method pointer plus operation-ID alias; inherited/overridden parameters, local request/response/header schemas and symbolic security. Supported explicit HTTP contracts can be complete. Multiple media types/successful response alternatives, undeclared security, unresolved schemes, unsupported dialects or transport details remain incomplete. |
| Google Discovery | Method ID/pointer, inherited native parameters, declared server and local referenced body schemas. Parser-normalized schemas and OAuth scope alternatives remain partial; no default endpoint or complete account/auth policy is inferred. |
| AWS Smithy | Absolute operation shape ID/pointer, native AWS protocol and input-member locations. Compound/member schemas and authorization remain partial. No synthetic generic HTTP server/method is emitted. |
| AsyncAPI | Operation ID/pointer and send/receive payload direction. A single local message schema can be projected; multiple messages, bindings and auth remain partial. |
| GraphQL | Native operation/root-field pointer, variables/arguments and selected outputs. Scalar types are partial; schema selection, custom scalars, lists and auth remain unproved. |
| OpenRPC | Native method ID/pointer, ordered parameter descriptors and local result schemas. Positional parameters preserve declaration order; transport/auth remains unproved. |
| gRPC/protobuf | Native service/method pointer/full-method aliases and request/response message identities with available top-level field shapes. Wire presence, numeric encoding, nested messages, enums and streams remain partial. No RPC connection/reflection occurs. |
| OData | Native resource/action/function/import/navigation identities and declared parameters/results. Facets, response presence, resource HTTP selection and auth remain partial. |

M82.2 adds `VerifyOperationShapeTable`, which reproduces claims from the exact
local source set and compares deterministic UWS bytes. It rejects structurally
valid forged fields, omitted operations and stale raw identities. This is
independent reproduction, not a signature or execution approval.

OpenAPI security preserves declaration-order OR alternatives and sorted AND
requirements inside each alternative. Explicit empty security is anonymous;
missing declarations/definitions remain unknown, with symbolic alternatives
retained. Operation security overrides root security. Scope symbols never load
tokens or contact token URLs. Complete schemas are checked with the existing
JSON Schema compiler and a loader that always refuses external resources.
Native-family projections with unsupported dialect/presence/transport/auth
evidence remain incomplete, even when a shallow type is available.
M82 qualified deterministic wire/limit fixtures, consumer compatibility, whole
review and publication. Each consuming worker/package owns its adoption. A caller must
not interpret partial native type hints as a complete schema or authority.
Effects are not a ShapeTable field; this producer never infers execution effects
from HTTP methods. Existing effect/candidate evidence remains separately advisory.

Defaults and maximums are 32 sources, 20 MiB per source, 64 MiB aggregate raw
bytes, 10,000 operations, UWS 8 MiB serialized table and 256 KiB per schema.
Optional byte/operation limits only tighten these bounds. Structural guards
limit depth and semantic work. One invocation also charges schema projection
work (100,000 nodes total, 10,000 per schema, depth 50) and projected schema
bytes (8 MiB total) before further compilation/accumulation. Source and operation
serialization are checked incrementally, including decoded claimed tables;
compact repeated references cannot grow a whole oversized table before refusal.
Malformed input, duplicate source IDs, unsafe
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
The public UWS `ShapeTable.Marshal` method supplies deterministic serialization;
`ParseTable` checks the bounded closed wire before independent reproduction.
The [eight-family fixture](../testdata/operation-shapes/v1/README.md) freezes
exact bytes with independent semantic/refusal tests, source-order independence
and resolver isolation. JSON/YAML lexeme tests do not use float64 or JCS equality
as their oracle. Native IDs that the legacy parsers deduplicate ambiguously
(including OData overloads and repeated AsyncAPI IDs) refuse rather than silently
choosing one contract; duplicate OpenAPI IDs keep distinct native pointers and
resolve as ambiguous by ID.
Schema keyword knownness follows the source dialect: OpenAPI 3.0's supported
subset, OpenAPI 3.1's default 2020-12 context, or the common draft-07 subset for
OpenRPC/AsyncAPI. Unsupported modern keywords remain unknown; a supported
explicit 2020-12 schema can establish its modern semantics. Reference siblings,
legacy exclusive bounds and unresolved dialect overrides remain conservative.
Standalone qualification uses `GOWORK=off` and the exact published UWS C09 pin.
The UWS module retains its Horizon/HashiCorp HCL transitive dependency closure;
this API does not make UWS or APItools HCL-free.
