package apitools

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/uws/binding"
)

func TestOperationShapesDeterministicGolden(t *testing.T) {
	sources := shapeFamilySources(t)
	table, err := BuildOperationShapeTable(context.Background(), OperationShapeOptions{Sources: sources})
	if err != nil {
		t.Fatal(err)
	}
	data := mustShapeBytes(t, table)
	golden, err := os.ReadFile(filepath.Join("testdata", "operation-shapes", "v1", "eight-families.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, bytes.TrimSpace(golden)) {
		t.Fatal("shape wire differs from reviewed eight-family fixture")
	}
	for i, j := 0, len(sources)-1; i < j; i, j = i+1, j-1 {
		sources[i], sources[j] = sources[j], sources[i]
	}
	other, err := BuildOperationShapeTable(context.Background(), OperationShapeOptions{Sources: sources})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, mustShapeBytes(t, other)) {
		t.Fatal("source order changed shape serialization")
	}
	parsed, err := binding.ParseTable(data)
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifyOperationShapeTable(context.Background(), OperationShapeOptions{Sources: sources}, parsed); err != nil {
		t.Fatal("serialized claims did not reproduce", err)
	}
	copy, err := binding.NewResolver(table)
	if err != nil {
		t.Fatal(err)
	}
	table.Operations[0].Protocol = "mutated"
	resolved, err := copy.Resolve(context.Background(), binding.Binding{Source: parsed.Operations[0].Source, SelectorKind: parsed.Operations[0].Selector.Kind, SelectorValue: parsed.Operations[0].Selector.Value})
	if err != nil || resolved.Shape == nil || resolved.Shape.Protocol == "mutated" {
		t.Fatal("caller mutation reached reviewed snapshot", err)
	}
}

func TestOperationShapesPreserveExactSchemaNumbers(t *testing.T) {
	for _, data := range []string{
		`{"openrpc":"1.3.2","info":{"title":"Exact","version":"1"},"methods":[{"name":"test","params":[{"name":"value","schema":{"type":"number","minimum":9007199254740993,"maximum":9007199254740995,"multipleOf":0.000000000000000000123456789,"enum":[1e+03,1.2300]}}]}]}`,
		"openapi: 3.1.0\ninfo: {title: Exact, version: '1'}\npaths:\n  /test:\n    get:\n      parameters:\n        - name: value\n          in: query\n          schema:\n            type: number\n            minimum: 9007199254740993\n            maximum: 9007199254740995\n            multipleOf: 0.000000000000000000123456789\n            enum: [1e+03, 1.2300]\n      responses: {'200': {description: exact}}\n",
	} {
		kind := OperationSourceOpenRPC
		if strings.HasPrefix(data, "openapi") {
			kind = OperationSourceOpenAPI
		}
		table, err := BuildOperationShapeTable(context.Background(), OperationShapeOptions{Sources: []ShapeSourceInput{{ID: "exact", OperationSourceInput: OperationSourceInput{Kind: kind, Content: []byte(data)}}}})
		if err != nil {
			t.Fatal(err)
		}
		schema := table.Operations[0].Inputs[0].Schema.JSON
		for _, lexeme := range []string{"9007199254740993", "9007199254740995", "0.000000000000000000123456789", "1e+03", "1.2300"} {
			if !bytes.Contains(schema, []byte(lexeme)) {
				t.Fatalf("numeric lexeme %s changed in %s", lexeme, schema)
			}
		}
		if !table.Operations[0].Inputs[0].Schema.Known {
			t.Fatal("valid exact numeric schema lost knownness")
		}
	}
}

func TestOperationShapesRefuseAmbiguousSourceEncoding(t *testing.T) {
	valid, _ := json.Marshal(completeShapeFixture())
	for _, data := range [][]byte{
		append(append([]byte(nil), valid...), valid...),
		[]byte(`{"openapi":"3.1.0","openapi":"3.0.3","paths":{}}`),
		[]byte("openapi: 3.1.0\nopenapi: 3.0.3\npaths: {}\n"),
		[]byte("openapi: 3.1.0\npaths: &x {}\nother: *x\n"),
		[]byte(strings.Repeat("[", 101) + strings.Repeat("]", 101)),
	} {
		if table, err := BuildOperationShapeTable(context.Background(), OperationShapeOptions{Sources: []ShapeSourceInput{{ID: "bad", OperationSourceInput: OperationSourceInput{Kind: OperationSourceOpenAPI, Content: data}}}}); !errors.Is(err, ErrOperationShapeTable) || len(table.Sources) > 0 || len(table.Operations) > 0 {
			t.Fatalf("ambiguous input was accepted: %v", err)
		}
	}
}

func TestOperationShapesKeepAmbiguousSelectors(t *testing.T) {
	fixture := completeShapeFixture()
	paths := mapValue(fixture["paths"])
	paths["/other"] = map[string]any{"get": map[string]any{"operationId": "getPet", "responses": mapValue(mapValue(paths["/pets/{id}"])["get"])["responses"]}}
	table, err := BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
	if err != nil {
		t.Fatal(err)
	}
	resolver, err := binding.NewResolver(table)
	if err != nil {
		t.Fatal(err)
	}
	resolution, err := resolver.Resolve(context.Background(), binding.Binding{Source: table.Sources[0], SelectorKind: "id", SelectorValue: "getPet"})
	if err != nil || resolution.Status != binding.Ambiguous || len(table.Operations) != 2 {
		t.Fatal("duplicate native IDs were resolved by sorting", resolution.Status, err)
	}
	for _, op := range table.Operations {
		resolution, err := resolver.Resolve(context.Background(), binding.Binding{Source: op.Source, SelectorKind: "ref", SelectorValue: op.Selector.Value})
		if err != nil || resolution.Status != binding.Resolved {
			t.Fatal("unique native ref was lost", err)
		}
	}
}

func TestOperationShapesRetainIncompleteReferencesAndAlternatives(t *testing.T) {
	fixture := completeShapeFixture()
	operation := mapValue(mapValue(mapValue(fixture["paths"])["/pets/{id}"])["get"])
	parameters := operation["parameters"].([]any)
	for _, schema := range []map[string]any{
		{"$ref": "#/components/schemas/Missing"},
		{"$ref": "#/components/schemas/Cycle"},
	} {
		fixture["components"] = map[string]any{"schemas": map[string]any{"Cycle": map[string]any{"$ref": "#/components/schemas/Cycle"}}}
		parameters[0].(map[string]any)["schema"] = schema
		table, err := BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
		if err != nil {
			t.Fatal(err)
		}
		if table.Operations[0].Inputs[0].Schema.Known || table.Operations[0].Complete {
			t.Fatal("missing/cyclic schema became complete")
		}
	}
	parameters[0].(map[string]any)["schema"] = map[string]any{"type": "string"}
	response := mapValue(mapValue(operation["responses"])["200"])
	mapValue(response["content"])["application/xml"] = map[string]any{"schema": map[string]any{"type": "object"}}
	table, err := BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
	if err != nil {
		t.Fatal(err)
	}
	if table.Operations[0].Complete {
		t.Fatal("multiple response media became complete")
	}
	mapValue(operation["responses"])["201"] = map[string]any{"description": "Other response", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "string"}}}}
	table, err = BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
	if err != nil || table.Operations[0].Outputs[0].Schema.Known || table.Operations[0].Complete {
		t.Fatal("multiple responses became a single proved contract", err)
	}
}

func TestOperationShapesProtocolEvidenceEdges(t *testing.T) {
	for name, source := range map[string]OperationSourceInput{
		"async-multiple": {Kind: OperationSourceAsyncAPI, Content: []byte(`{"asyncapi":"3.0.0","info":{"title":"Events","version":"1"},"operations":{"send":{"action":"send","channel":{"$ref":"#/channels/c"},"messages":[{"$ref":"#/channels/c/messages/a"},{"$ref":"#/channels/c/messages/b"}]}},"channels":{"c":{"messages":{"a":{"payload":{"type":"string"}},"b":{"payload":{"type":"integer"}}}}}}`)},
		"graphql-list":   {Kind: OperationSourceGraphQL, Content: []byte(`type Query { pets: [String!]! }`)},
		"grpc-stream":    {Kind: OperationSourceGRPCProtobuf, Content: []byte(`syntax = "proto3"; package fixture; message Request { string id = 1; } message Response { string id = 1; } service Pets { rpc Watch(stream Request) returns (stream Response); }`)},
	} {
		t.Run(name, func(t *testing.T) {
			table, err := BuildOperationShapeTable(context.Background(), OperationShapeOptions{Sources: []ShapeSourceInput{{ID: name, OperationSourceInput: source}}})
			if err != nil {
				t.Fatal(err)
			}
			if len(table.Operations) != 1 || table.Operations[0].Complete || table.Operations[0].Method != "" {
				t.Fatal("native limitation lost")
			}
			for _, in := range table.Operations[0].Inputs {
				if in.Schema.Known {
					t.Fatal("unsupported native input became complete")
				}
			}
			for _, out := range table.Operations[0].Outputs {
				if out.Schema.Known {
					t.Fatal("unsupported native output became complete")
				}
			}
		})
	}
}

func TestOperationShapesRefuseCollidingNativeContracts(t *testing.T) {
	for name, source := range map[string]OperationSourceInput{
		"async-duplicate-id":     {Kind: OperationSourceAsyncAPI, Content: []byte(`{"asyncapi":"2.6.0","info":{"title":"Events","version":"1"},"channels":{"a":{"publish":{"operationId":"duplicate"}},"b":{"subscribe":{"operationId":"duplicate"}}}}`)},
		"odata-overload":         {Kind: OperationSourceOData, Content: []byte(`<Schema xmlns="http://docs.oasis-open.org/odata/ns/edm" Namespace="Fixture"><Action Name="Update"><Parameter Name="id" Type="Edm.String"/></Action><Action Name="Update"><Parameter Name="id" Type="Edm.Int32"/></Action></Schema>`)},
		"discovery-duplicate-id": {Kind: OperationSourceGoogleDiscovery, Content: []byte(`{"discoveryVersion":"v1","name":"fixture","version":"v1","methods":{"a":{"id":"duplicate","path":"a","httpMethod":"GET"},"b":{"id":"duplicate","path":"b","httpMethod":"POST"}}}`)},
		"grpc-duplicate-method":  {Kind: OperationSourceGRPCProtobuf, Content: []byte(`syntax = "proto3"; package fixture; message Request { string id = 1; } service Pets { rpc Get(Request) returns (Request); rpc Get(Request) returns (Request); }`)},
	} {
		t.Run(name, func(t *testing.T) {
			table, err := BuildOperationShapeTable(context.Background(), OperationShapeOptions{Sources: []ShapeSourceInput{{ID: name, OperationSourceInput: source}}})
			if !errors.Is(err, ErrOperationShapeTable) || len(table.Operations) > 0 {
				t.Fatal("collision was silently selected", err)
			}
		})
	}
}

func TestOperationShapesBoundSerializedMetadata(t *testing.T) {
	fixture := completeShapeFixture()
	paths := map[string]any{}
	schema := map[string]any{"type": "string", "enum": []any{strings.Repeat("a", 80000)}}
	for i := 0; i < 110; i++ {
		paths[fmt.Sprintf("/item/%d", i)] = map[string]any{"get": map[string]any{"responses": map[string]any{"200": map[string]any{"description": "Large fixture", "content": map[string]any{"application/json": map[string]any{"schema": schema}}}}}}
	}
	fixture["paths"] = paths
	options := optionsForShapeFixture(t, fixture)
	if len(options.Sources[0].Content) > int(DefaultMaxBytes) {
		t.Fatal("fixture did not isolate table budget")
	}
	table, err := BuildOperationShapeTable(context.Background(), options)
	if !errors.Is(err, ErrOperationShapeTable) || len(table.Operations) > 0 {
		t.Fatal("oversized serialized table was returned", err)
	}
	options = optionsForShapeFixture(t, completeShapeFixture())
	options.Sources = append(options.Sources, ShapeSourceInput{ID: "second", OperationSourceInput: options.Sources[0].OperationSourceInput})
	options.MaxOperations = 1
	table, err = BuildOperationShapeTable(context.Background(), options)
	if !errors.Is(err, ErrOperationShapeTable) || len(table.Operations) > 0 {
		t.Fatal("aggregate operation limit returned partial positives", err)
	}
}
