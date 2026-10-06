package apitools_test

import (
	"context"
	"testing"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/uws/binding"
)

func TestOperationShapesPublicContract(t *testing.T) {
	table, err := apitools.BuildOperationShapeTable(context.Background(), apitools.OperationShapeOptions{Sources: []apitools.ShapeSourceInput{{ID: "pets", OperationSourceInput: apitools.OperationSourceInput{Kind: apitools.OperationSourceOpenRPC, Content: []byte(`{"openrpc":"1.3.2","info":{"title":"pets","version":"1"},"methods":[{"name":"pets.get","params":[{"name":"id","required":true,"schema":{"type":"string"}}],"result":{"name":"pet","schema":{"type":"object"}}}]}`)}}}})
	if err != nil {
		t.Fatal(err)
	}
	resolver, err := binding.NewResolver(table)
	if err != nil {
		t.Fatal(err)
	}
	result, err := resolver.Resolve(context.Background(), binding.Binding{Source: table.Sources[0], SelectorKind: "id", SelectorValue: "pets.get"})
	if err != nil || result.Status != binding.Resolved || result.Shape.Protocol != "json-rpc" || result.Shape.Inputs[0].Name != "id" || !result.Shape.Inputs[0].Schema.Known {
		t.Fatalf("public shape projection failed: %+v %v", result, err)
	}
}
