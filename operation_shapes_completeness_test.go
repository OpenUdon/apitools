package apitools

import (
	"context"
	"encoding/json"
	"testing"
)

func TestOperationShapesCannotProveMalformedHTTPMetadata(t *testing.T) {
	for name, change := range map[string]func(map[string]any){
		"missing-path-parameter": func(f map[string]any) {
			delete(mapValue(mapValue(mapValue(f["paths"])["/pets/{id}"])["get"]), "parameters")
		},
		"unsupported-location": func(f map[string]any) {
			p := mapValue(mapValue(mapValue(f["paths"])["/pets/{id}"])["get"])["parameters"].([]any)[0].(map[string]any)
			p["in"] = "custom"
		},
		"wrong-path-parameter": func(f map[string]any) {
			p := mapValue(mapValue(mapValue(f["paths"])["/pets/{id}"])["get"])["parameters"].([]any)[0].(map[string]any)
			p["name"] = "other"
		},
		"malformed-template": func(f map[string]any) {
			paths := mapValue(f["paths"])
			paths["/pets/{id"] = paths["/pets/{id}"]
			delete(paths, "/pets/{id}")
		},
	} {
		t.Run(name, func(t *testing.T) {
			fixture := completeShapeFixture()
			change(fixture)
			table, err := BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
			if err == nil && table.Operations[0].Complete {
				t.Fatal("malformed HTTP metadata became complete")
			}
		})
	}
}

func TestOperationShapesUnknownDialectAndReferenceSiblings(t *testing.T) {
	for name, change := range map[string]func(map[string]any){
		"dialect": func(f map[string]any) { f["jsonSchemaDialect"] = "https://fixture.invalid/custom-schema" },
		"reference-sibling": func(f map[string]any) {
			f["components"] = map[string]any{"schemas": map[string]any{"ID": map[string]any{"type": "string"}}}
			p := mapValue(mapValue(mapValue(f["paths"])["/pets/{id}"])["get"])["parameters"].([]any)[0].(map[string]any)
			p["schema"] = map[string]any{"$ref": "#/components/schemas/ID", "type": "integer"}
		},
		"legacy-exclusive-bound": func(f map[string]any) {
			p := mapValue(mapValue(mapValue(f["paths"])["/pets/{id}"])["get"])["parameters"].([]any)[0].(map[string]any)
			p["schema"] = map[string]any{"type": "number", "minimum": 1, "exclusiveMinimum": true}
		},
	} {
		t.Run(name, func(t *testing.T) {
			fixture := completeShapeFixture()
			change(fixture)
			table, err := BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
			if err != nil {
				t.Fatal(err)
			}
			if table.Operations[0].Complete || table.Operations[0].Inputs[0].Schema.Known {
				t.Fatal("unproved dialect/ref semantics became known")
			}
		})
	}
}

func TestOperationShapesRespectSourceSchemaDialect(t *testing.T) {
	for _, value := range []any{true, false, map[string]any{"type": []any{"string", "null"}}, map[string]any{"type": "null"}} {
		for _, version := range []string{"3.0.3", "3.1.0"} {
			fixture := completeShapeFixture()
			fixture["openapi"] = version
			parameter := mapValue(mapValue(mapValue(fixture["paths"])["/pets/{id}"])["get"])["parameters"].([]any)[0].(map[string]any)
			parameter["schema"] = value
			table, err := BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
			if err != nil {
				t.Fatal(err)
			}
			if table.Operations[0].Inputs[0].Schema.Known != (version == "3.1.0") {
				t.Fatal("boolean/union/null schema dialect was ignored", version, value)
			}
		}
	}
	for _, version := range []string{"3.0.3", "3.1.0"} {
		fixture := completeShapeFixture()
		fixture["openapi"] = version
		parameter := mapValue(mapValue(mapValue(fixture["paths"])["/pets/{id}"])["get"])["parameters"].([]any)[0].(map[string]any)
		parameter["schema"] = map[string]any{"type": "string", "const": "declared"}
		table, err := BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
		if err != nil {
			t.Fatal(err)
		}
		if table.Operations[0].Inputs[0].Schema.Known != (version == "3.1.0") {
			t.Fatal("OpenAPI schema dialect was ignored", version)
		}
	}
	for _, explicit := range []bool{false, true} {
		schema := map[string]any{"type": "array", "prefixItems": []any{map[string]any{"type": "string"}}}
		if explicit {
			schema["$schema"] = "https://json-schema.org/draft/2020-12/schema"
		}
		value := map[string]any{"openrpc": "1.3.2", "info": map[string]any{"title": "Dialect", "version": "1"}, "methods": []any{map[string]any{"name": "test", "params": []any{map[string]any{"name": "value", "schema": schema}}}}}
		data, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		table, err := BuildOperationShapeTable(context.Background(), OperationShapeOptions{Sources: []ShapeSourceInput{{ID: "rpc", OperationSourceInput: OperationSourceInput{Kind: OperationSourceOpenRPC, Content: data}}}})
		if err != nil {
			t.Fatal(err)
		}
		if table.Operations[0].Inputs[0].Schema.Known != explicit {
			t.Fatal("default/explicit OpenRPC dialect evidence was ignored")
		}
	}
}
