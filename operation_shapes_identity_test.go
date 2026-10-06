package apitools

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/OpenUdon/uws/binding"
)

func completeShapeFixture() map[string]any {
	return map[string]any{
		"openapi": "3.1.0", "info": map[string]any{"title": "Fixture", "version": "1"},
		"servers":  []any{map[string]any{"url": "https://fixture.invalid"}},
		"security": []any{map[string]any{}},
		"paths": map[string]any{"/pets/{id}": map[string]any{"get": map[string]any{
			"operationId": "getPet",
			"parameters":  []any{map[string]any{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "string"}}},
			"responses":   map[string]any{"200": map[string]any{"description": "Pet", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"id": map[string]any{"type": "string"}}}}}}},
		}}},
	}
}

func optionsForShapeFixture(t *testing.T, value map[string]any) OperationShapeOptions {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return OperationShapeOptions{Sources: []ShapeSourceInput{{ID: "pets", OperationSourceInput: OperationSourceInput{Kind: OperationSourceOpenAPI, Content: data}}}}
}

func TestOperationShapesProveSupportedMetadataOnly(t *testing.T) {
	table, err := BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, completeShapeFixture()))
	if err != nil {
		t.Fatal(err)
	}
	shape := table.Operations[0]
	if !shape.Complete || !shape.Security.Known || len(shape.Security.Alternatives) != 1 || len(shape.Security.Alternatives[0].Requirements) != 0 {
		t.Fatalf("explicit supported anonymous shape incomplete: %+v", shape)
	}
	resolver, err := binding.NewResolver(table)
	if err != nil {
		t.Fatal(err)
	}
	request := binding.Request{Binding: binding.Binding{Source: shape.Source, SelectorKind: "id", SelectorValue: "getPet"}, Inputs: []binding.BoundInput{{Location: "path", Name: "id", Value: "pet-1"}}}
	report, err := binding.ValidateBinding(context.Background(), resolver, request)
	if err != nil || report.Outcome != binding.Compatible {
		t.Fatalf("supported metadata did not prove binding: %+v %v", report, err)
	}
	request.Inputs[0].Value = 3
	report, err = binding.ValidateBinding(context.Background(), resolver, request)
	if err != nil || report.Outcome != binding.Incompatible {
		t.Fatalf("declared type mismatch was lost: %+v %v", report, err)
	}
}

func TestOperationShapesRetainSecurityAlternatives(t *testing.T) {
	fixture := completeShapeFixture()
	fixture["components"] = map[string]any{"securitySchemes": map[string]any{
		"key":    map[string]any{"type": "apiKey", "in": "header", "name": "X-Key"},
		"oauth":  map[string]any{"type": "oauth2", "flows": map[string]any{"clientCredentials": map[string]any{"tokenUrl": "https://fixture.invalid/token", "scopes": map[string]any{"read": "Read pets"}}}},
		"bearer": map[string]any{"type": "http", "scheme": "bearer"},
	}}
	fixture["security"] = []any{map[string]any{"key": []any{}, "oauth": []any{"read"}}, map[string]any{"bearer": []any{}}}
	table, err := BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
	if err != nil {
		t.Fatal(err)
	}
	security := table.Operations[0].Security
	if !security.Known || len(security.Alternatives) != 2 || len(security.Alternatives[0].Requirements) != 2 || len(security.Alternatives[1].Requirements) != 1 || security.Alternatives[0].Requirements[1].Scopes[0] != "read" {
		t.Fatalf("flattened/lost security: %+v", security)
	}
	resolver, err := binding.NewResolver(table)
	if err != nil {
		t.Fatal(err)
	}
	request := binding.Request{Binding: binding.Binding{Source: table.Sources[0], SelectorKind: "id", SelectorValue: "getPet"}, Inputs: []binding.BoundInput{{Location: "path", Name: "id", Value: "1"}}, Security: []binding.SecurityBinding{{Scheme: "key", CredentialSlot: "symbolic-key"}}}
	report, err := binding.ValidateBinding(context.Background(), resolver, request)
	if err != nil || report.Outcome != binding.Incompatible {
		t.Fatalf("AND requirement relaxed: %+v %v", report, err)
	}
	request.Security = []binding.SecurityBinding{{Scheme: "bearer", CredentialSlot: "symbolic-bearer"}}
	report, err = binding.ValidateBinding(context.Background(), resolver, request)
	if err != nil || report.Outcome != binding.Compatible {
		t.Fatalf("OR alternative lost: %+v %v", report, err)
	}
	delete(fixture, "components")
	table, err = BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
	if err != nil || table.Operations[0].Security.Known || table.Operations[0].Complete || len(table.Operations[0].Security.Alternatives) != 2 {
		t.Fatalf("missing definitions became anonymous/known: %+v %v", table, err)
	}
}

func TestOperationShapesVerifyReproducedClaims(t *testing.T) {
	options := optionsForShapeFixture(t, completeShapeFixture())
	table, err := BuildOperationShapeTable(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifyOperationShapeTable(context.Background(), options, table); err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*binding.ShapeTable){
		func(t *binding.ShapeTable) { t.Operations[0].Method = "POST" },
		func(t *binding.ShapeTable) { t.Operations[0].Inputs[0].Required = false },
		func(t *binding.ShapeTable) {
			t.Operations[0].Inputs[0].Schema.JSON = json.RawMessage(`{"type":"integer"}`)
		},
		func(t *binding.ShapeTable) {
			t.Sources[0].SHA256 = strings.Repeat("a", 64)
			t.Operations[0].Source = t.Sources[0]
		},
		func(t *binding.ShapeTable) { t.Operations[0].Security.Known = false },
		func(t *binding.ShapeTable) { t.Operations = nil },
	} {
		copy, err := binding.ParseTable(mustShapeBytes(t, table))
		if err != nil {
			t.Fatal(err)
		}
		mutate(&copy)
		if copy.Validate() != nil {
			t.Fatal("test did not exercise a structurally valid forged claim")
		}
		if err := VerifyOperationShapeTable(context.Background(), options, copy); !errors.Is(err, ErrOperationShapeTable) {
			t.Fatal("forged claims accepted", err)
		}
	}
	options.Sources[0].Content = append(options.Sources[0].Content, ' ')
	if err := VerifyOperationShapeTable(context.Background(), options, table); !errors.Is(err, ErrOperationShapeTable) {
		t.Fatal("stale raw identity accepted", err)
	}
}

func mustShapeBytes(t *testing.T, table binding.ShapeTable) []byte {
	t.Helper()
	data, err := table.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestOperationShapesUnknownEvidenceAndNoNetwork(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls.Add(1) }))
	defer server.Close()
	fixture := completeShapeFixture()
	operation := mapValue(mapValue(mapValue(fixture["paths"])["/pets/{id}"])["get"])
	parameter := operation["parameters"].([]any)[0].(map[string]any)
	for _, schema := range []any{
		nil, map[string]any{"$ref": server.URL + "/schema"},
		map[string]any{"type": "string", "x-opaque-constraint": true},
		map[string]any{"type": "string", "nullable": true},
		map[string]any{"type": 17},
	} {
		parameter["schema"] = schema
		options := optionsForShapeFixture(t, fixture)
		options.Sources[0].URL = server.URL + "/source?secret=canary-value#private"
		table, err := BuildOperationShapeTable(context.Background(), options)
		if err != nil {
			t.Fatal(err)
		}
		if table.Operations[0].Inputs[0].Schema.Known || table.Operations[0].Complete {
			t.Fatal("incomplete schema evidence became complete")
		}
		if bytes.Contains(mustShapeBytes(t, table), []byte("canary-value")) {
			t.Fatal("private URL data retained")
		}
	}
	if calls.Load() != 0 {
		t.Fatal("source provenance/reference triggered network")
	}
	delete(fixture, "security")
	parameter["schema"] = map[string]any{"type": "string"}
	table, err := BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
	if err != nil || table.Operations[0].Security.Known || table.Operations[0].Complete {
		t.Fatal("undeclared auth became anonymous", err)
	}
}
