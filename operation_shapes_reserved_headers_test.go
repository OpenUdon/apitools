package apitools

import (
	"context"
	"testing"

	"github.com/OpenUdon/uws/binding"
)

func TestOperationShapesOpenAPIReservedHeaders(t *testing.T) {
	for _, version := range []string{"3.0.0", "3.0.3", "3.0.4", "3.1.0", "3.1.1", "3.1.2"} {
		for _, referenced := range []bool{false, true} {
			for _, name := range []string{"Accept", "accept", "aCcEpT", "Content-Type", "content-type", "CONTENT-TYPE", "Authorization", "authorization", "AUTHORIZATION"} {
				t.Run(version+"/"+name+map[bool]string{false: "/direct", true: "/ref"}[referenced], func(t *testing.T) {
					fixture := completeShapeFixture()
					fixture["openapi"] = version
					path := mapValue(mapValue(fixture["paths"])["/pets/{id}"])
					operation := mapValue(path["get"])
					parameter := map[string]any{"name": name, "in": "header", "required": true, "schema": map[string]any{"type": "string"}}
					var raw any = parameter
					if referenced {
						fixture["components"] = map[string]any{"parameters": map[string]any{"reserved": parameter}}
						raw = map[string]any{"$ref": "#/components/parameters/reserved"}
					}
					// Exercise both inherited and operation-level parameters.
					path["parameters"] = []any{raw}
					operation["parameters"] = append(operation["parameters"].([]any), raw)
					table, err := BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
					if err != nil {
						t.Fatal(err)
					}
					shape := table.Operations[0]
					if !shape.Complete || len(shape.Inputs) != 1 || shape.Inputs[0].Name != "id" {
						t.Fatalf("reserved header became an input: %+v", shape)
					}
					resolver, err := binding.NewResolver(table)
					if err != nil {
						t.Fatal(err)
					}
					request := binding.Request{Binding: binding.Binding{Source: shape.Source, SelectorKind: "id", SelectorValue: "getPet"}, Inputs: []binding.BoundInput{{Location: "path", Name: "id", Value: "pet-1"}}}
					report, err := binding.ValidateBinding(context.Background(), resolver, request)
					if err != nil || report.Outcome != binding.Compatible {
						t.Fatalf("reserved header manufactured a required binding: %+v %v", report, err)
					}
				})
			}
		}
	}
}
