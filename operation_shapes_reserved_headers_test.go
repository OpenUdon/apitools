package apitools

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
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

func TestOperationShapesReservedHeadersKeepRealInputsAndSecurity(t *testing.T) {
	fixture := completeShapeFixture()
	path := mapValue(mapValue(fixture["paths"])["/pets/{id}"])
	operation := mapValue(path["get"])
	parameters := operation["parameters"].([]any)
	for _, location := range []string{"query", "cookie"} {
		for _, name := range []string{"Accept", "Content-Type", "Authorization"} {
			parameters = append(parameters, map[string]any{"name": name, "in": location, "required": true, "schema": map[string]any{"type": "string"}})
		}
	}
	parameters = append(parameters, map[string]any{"$ref": "#/components/parameters/real"})
	// Ignored definitions cannot introduce unsupported schema/serialization gaps.
	for _, name := range []string{"Accept", "Content-Type", "Authorization"} {
		parameters = append(parameters, map[string]any{"name": name, "in": "header", "required": true, "schema": map[string]any{"$ref": "https://fixture.invalid/ignored"}, "style": "simple", "explode": true})
	}
	operation["parameters"] = parameters
	fixture["components"] = map[string]any{
		"parameters": map[string]any{"real": map[string]any{"name": "X-Required", "in": "header", "required": true, "schema": map[string]any{"type": "string"}}},
		"securitySchemes": map[string]any{
			"key":    map[string]any{"type": "apiKey", "in": "header", "name": "Authorization"},
			"second": map[string]any{"type": "apiKey", "in": "header", "name": "X-Key"},
			"bearer": map[string]any{"type": "http", "scheme": "bearer"},
		},
	}
	fixture["security"] = []any{map[string]any{"key": []any{}, "second": []any{}}, map[string]any{"bearer": []any{}}}
	options := optionsForShapeFixture(t, fixture)
	table, err := BuildOperationShapeTable(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	shape := table.Operations[0]
	if !shape.Complete || len(shape.Inputs) != 8 || len(shape.Security.Alternatives) != 2 || len(shape.Security.Alternatives[0].Requirements) != 2 || len(shape.Security.Alternatives[1].Requirements) != 1 {
		t.Fatalf("real input or symbolic security lost: %+v", shape)
	}
	resolver, err := binding.NewResolver(table)
	if err != nil {
		t.Fatal(err)
	}
	request := binding.Request{Binding: binding.Binding{Source: shape.Source, SelectorKind: "id", SelectorValue: "getPet"}, Security: []binding.SecurityBinding{{Scheme: "bearer", CredentialSlot: "symbolic-bearer"}}}
	for _, input := range shape.Inputs {
		request.Inputs = append(request.Inputs, binding.BoundInput{Location: input.Location, Name: input.Name, Value: "fixture"})
	}
	report, err := binding.ValidateBinding(context.Background(), resolver, request)
	if err != nil || report.Outcome != binding.Compatible {
		t.Fatalf("genuine inputs and OR alternative refused: %+v %v", report, err)
	}
	request.Inputs = nil
	for _, input := range shape.Inputs {
		if input.Location != "header" {
			request.Inputs = append(request.Inputs, binding.BoundInput{Location: input.Location, Name: input.Name, Value: "fixture"})
		}
	}
	report, err = binding.ValidateBinding(context.Background(), resolver, request)
	if err != nil || report.Outcome != binding.Incompatible || len(report.Diagnostics) == 0 || report.Diagnostics[0].Code != "binding.input_required" {
		t.Fatalf("missing genuine header accepted: %+v %v", report, err)
	}
	request.Security = []binding.SecurityBinding{{Scheme: "key", CredentialSlot: "symbolic-key"}}
	for _, input := range shape.Inputs {
		if input.Location == "header" {
			request.Inputs = append(request.Inputs, binding.BoundInput{Location: input.Location, Name: input.Name, Value: "fixture"})
		}
	}
	report, err = binding.ValidateBinding(context.Background(), resolver, request)
	if err != nil || report.Outcome != binding.Incompatible {
		t.Fatalf("AND security requirement relaxed: %+v %v", report, err)
	}
	if err := VerifyOperationShapeTable(context.Background(), options, table); err != nil {
		t.Fatal(err)
	}
	again, err := BuildOperationShapeTable(context.Background(), options)
	if err != nil || !bytes.Equal(mustShapeBytes(t, table), mustShapeBytes(t, again)) {
		t.Fatal("shape bytes are not deterministic", err)
	}
	digest := sha256.Sum256(options.Sources[0].Content)
	if table.Sources[0].SHA256 != hex.EncodeToString(digest[:]) || shape.Source != table.Sources[0] || shape.Selector.Value != "#/paths/~1pets~1{id}/get" {
		t.Fatal("raw source/native identity changed")
	}
	forged, err := binding.ParseTable(mustShapeBytes(t, table))
	if err != nil {
		t.Fatal(err)
	}
	forged.Operations[0].Inputs = append(forged.Operations[0].Inputs, binding.Input{Location: "header", Name: "Authorization", Required: true, Schema: shape.Inputs[0].Schema})
	sortShapeInputs(forged.Operations[0].Inputs)
	if err := VerifyOperationShapeTable(context.Background(), options, forged); !errors.Is(err, ErrOperationShapeTable) {
		t.Fatal("fabricated reserved-header claim verified", err)
	}
	options.Sources[0].Content = append(options.Sources[0].Content, ' ')
	if err := VerifyOperationShapeTable(context.Background(), options, table); !errors.Is(err, ErrOperationShapeTable) {
		t.Fatal("stale reserved-header source verified", err)
	}
}

func TestOperationShapesReservedHeaderRuleKeepsPathAndSwagger(t *testing.T) {
	for _, name := range []string{"Accept", "Content-Type", "Authorization"} {
		fixture := completeShapeFixture()
		paths := mapValue(fixture["paths"])
		path := paths["/pets/{id}"]
		delete(paths, "/pets/{id}")
		paths["/pets/{"+name+"}"] = path
		operation := mapValue(mapValue(path)["get"])
		mapValue(operation["parameters"].([]any)[0])["name"] = name
		table, err := BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
		if err != nil || len(table.Operations[0].Inputs) != 1 || table.Operations[0].Inputs[0].Name != name || !table.Operations[0].Inputs[0].Required {
			t.Fatalf("real path input %s lost: %+v %v", name, table, err)
		}
	}
	fixture := completeShapeFixture()
	delete(fixture, "openapi")
	delete(fixture, "servers")
	fixture["swagger"], fixture["host"], fixture["schemes"], fixture["produces"] = "2.0", "fixture.invalid", []any{"https"}, []any{"application/json"}
	operation := mapValue(mapValue(mapValue(fixture["paths"])["/pets/{id}"])["get"])
	operation["parameters"] = []any{map[string]any{"name": "id", "in": "path", "required": true, "type": "string"}}
	for _, name := range []string{"Accept", "Content-Type", "Authorization"} {
		operation["parameters"] = append(operation["parameters"].([]any), map[string]any{"name": name, "in": "header", "required": true, "type": "string"})
	}
	operation["responses"] = map[string]any{"200": map[string]any{"description": "Pet", "schema": map[string]any{"type": "object"}}}
	table, err := BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
	if err != nil || !table.Operations[0].Complete || len(table.Operations[0].Inputs) != 4 {
		t.Fatalf("Swagger behavior changed: %+v %v", table, err)
	}
	for _, input := range table.Operations[0].Inputs {
		if !input.Required || strings.TrimSpace(input.Name) == "" {
			t.Fatal("Swagger header requiredness lost")
		}
	}
}

func TestOperationShapesReservedHeaderRefsRetainRefusals(t *testing.T) {
	fixture := completeShapeFixture()
	operation := mapValue(mapValue(mapValue(fixture["paths"])["/pets/{id}"])["get"])
	fixture["components"] = map[string]any{"parameters": map[string]any{
		"chain":            map[string]any{"$ref": "#/components/parameters/reserved~1~0header"},
		"reserved/~header": map[string]any{"name": "AUTHORIZATION", "in": "header", "required": true, "schema": map[string]any{"type": "string"}},
	}}
	operation["parameters"] = append(operation["parameters"].([]any), map[string]any{"$ref": "#/components/parameters/chain"})
	table, err := BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
	if err != nil || len(table.Operations[0].Inputs) != 1 || !table.Operations[0].Complete {
		t.Fatalf("local chain/escaped pointer not ignored: %+v %v", table, err)
	}
	for _, ref := range []string{"#/components/parameters/chain", "#/components/parameters/missing", "https://fixture.invalid/parameter"} {
		mapValue(mapValue(fixture["components"])["parameters"])["chain"] = map[string]any{"$ref": ref}
		table, err := BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
		if !errors.Is(err, ErrOperationShapeTable) || len(table.Operations) != 0 {
			t.Fatalf("unsafe/unresolved parameter ref bypassed: %+v %v", table, err)
		}
	}
}
