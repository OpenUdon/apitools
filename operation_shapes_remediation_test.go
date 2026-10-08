package apitools

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestOperationShapesAsyncAPIVersionDirections(t *testing.T) {
	for _, direction := range []struct {
		action2, action3 string
		input            bool
	}{
		{"publish", "receive", false},
		{"subscribe", "send", true},
	} {
		t.Run(direction.action2, func(t *testing.T) {
			sources := []string{
				fmt.Sprintf(`{"asyncapi":"2.6.0","info":{"title":"Events","version":"1"},"channels":{"events":{"%s":{"operationId":"event","message":{"payload":{"type":"string"}}}}}}`, direction.action2),
				fmt.Sprintf(`{"asyncapi":"3.0.0","info":{"title":"Events","version":"1"},"operations":{"event":{"action":"%s","channel":{"$ref":"#/channels/events"},"messages":[{"$ref":"#/channels/events/messages/event"}]}},"channels":{"events":{"messages":{"event":{"payload":{"type":"string"}}}}}}`, direction.action3),
			}
			for _, source := range sources {
				options := OperationShapeOptions{Sources: []ShapeSourceInput{{ID: "events", OperationSourceInput: OperationSourceInput{Kind: OperationSourceAsyncAPI, Content: []byte(source)}}}}
				table, err := BuildOperationShapeTable(context.Background(), options)
				if err != nil {
					t.Fatal(err)
				}
				if len(table.Operations) != 1 {
					t.Fatal("missing native operation")
				}
				operation := table.Operations[0]
				if (len(operation.Inputs) == 1) != direction.input || (len(operation.Outputs) == 1) == direction.input {
					t.Fatalf("wrong %s/%s direction: %+v", direction.action2, direction.action3, operation)
				}
				if operation.Complete || operation.Security.Known || operation.Selector.Value != "#/operations/event" {
					t.Fatal("lost native selector or partial evidence")
				}
				if err := VerifyOperationShapeTable(context.Background(), options, table); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestOperationShapesGraphQLDefaultsAndAliases(t *testing.T) {
	options := OperationShapeOptions{Sources: []ShapeSourceInput{{ID: "graphql", OperationSourceInput: OperationSourceInput{Kind: OperationSourceGraphQL, Content: []byte(`type Query { field(first: Int! = 10, required: Int!): String } query Q($first: Int! = 10, $required: Int!) { a: field b: field }`)}}}}
	table, err := BuildOperationShapeTable(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	for _, operation := range table.Operations {
		if len(operation.Inputs) != 2 || operation.Inputs[0].Name != "first" || operation.Inputs[0].Required || !operation.Inputs[1].Required {
			t.Fatalf("wrong default requiredness: %+v", operation)
		}
		if operation.Selector.Value == "#/operations/query.Q" {
			if len(operation.Outputs) != 2 || operation.Outputs[0].Name != "a" || operation.Outputs[1].Name != "b" {
				t.Fatalf("aliases collapsed: %+v", operation)
			}
		}
	}
	options.Sources[0].Content = []byte(`query Q { same: first same: second }`)
	if _, err := BuildOperationShapeTable(context.Background(), options); err == nil {
		t.Fatal("duplicate response key silently accepted")
	}
}

func TestOperationShapesSmithyLiteralMemberProvenance(t *testing.T) {
	for name, members := range map[string]string{
		"literal":              `{}`,
		"same-name-body":       `{"x-id":{"target":"smithy.api#String"}}`,
		"same-name-query":      `{"x-id":{"target":"smithy.api#String","traits":{"smithy.api#httpQuery":"x-id","smithy.api#required":{}}}}`,
		"different-name-query": `{"Id":{"target":"smithy.api#String","traits":{"smithy.api#httpQuery":"x-id"}}}`,
	} {
		t.Run(name, func(t *testing.T) {
			source := fmt.Sprintf(`{"smithy":"2.0","shapes":{"example#Service":{"type":"service","version":"1","operations":[{"target":"example#Get"}],"traits":{"aws.protocols#restJson1":{}}},"example#Get":{"type":"operation","input":{"target":"example#Input"},"traits":{"smithy.api#http":{"method":"GET","uri":"/things?x-id=Get","code":200}}},"example#Input":{"type":"structure","members":%s}}}`, members)
			table, err := BuildOperationShapeTable(context.Background(), OperationShapeOptions{Sources: []ShapeSourceInput{{ID: "smithy", OperationSourceInput: OperationSourceInput{Kind: OperationSourceAWSSmithy, Content: []byte(source)}}}})
			if err != nil {
				t.Fatal(err)
			}
			operation := table.Operations[0]
			if operation.Complete || operation.Method != "" || operation.Selector.Value != "#/shapes/example#Get" {
				t.Fatal("native protocol/identity lost")
			}
			if name == "literal" && len(operation.Inputs) != 0 {
				t.Fatalf("fabricated literal input: %+v", operation.Inputs)
			}
			if name != "literal" && len(operation.Inputs) != 1 {
				t.Fatalf("lost real member or retained synthetic member: %+v", operation.Inputs)
			}
			if name == "same-name-body" && operation.Inputs[0].Location != "body" {
				t.Fatal("synthetic literal hid body member")
			}
			if name == "same-name-query" && (operation.Inputs[0].Location != "query" || !operation.Inputs[0].Required) {
				t.Fatal("lost real required query member")
			}
			if name == "different-name-query" && (operation.Inputs[0].Location != "query" || operation.Inputs[0].Name != "Id") {
				t.Fatal("lost native member name")
			}
		})
	}
}

func TestOperationShapesDiscoveryRequestsAndEndpoints(t *testing.T) {
	for name, method := range map[string]string{
		"normal":         `"path":"files","request":{"$ref":"File"}`,
		"normal-upload":  `"path":"files","request":{"$ref":"File"},"mediaUpload":{"protocols":{"simple":{"path":"/upload/files","multipart":true},"resumable":{"path":"/resumable/files"}}}`,
		"upload-only":    `"mediaUpload":{"protocols":{"simple":{"path":"/upload/files"}}}`,
		"inline":         `"path":"files","request":{"type":"object","properties":{"id":{"type":"string"}}}`,
		"missing-schema": `"path":"files","request":{"$ref":"Missing"}`,
		"no-request":     `"path":"files"`,
	} {
		t.Run(name, func(t *testing.T) {
			source := fmt.Sprintf(`{"discoveryVersion":"v1","name":"files","version":"v1","rootUrl":"https://example.invalid/","servicePath":"api/v1/","schemas":{"File":{"type":"object"}},"resources":{"files":{"methods":{"insert":{"id":"files.insert","httpMethod":"POST",%s}}}}}`, method)
			table, err := BuildOperationShapeTable(context.Background(), OperationShapeOptions{Sources: []ShapeSourceInput{{ID: "discovery", OperationSourceInput: OperationSourceInput{Kind: OperationSourceGoogleDiscovery, Content: []byte(source)}}}})
			if err != nil {
				t.Fatal(err)
			}
			operation := table.Operations[0]
			if operation.Complete || operation.Selector.Value != "#/methods/files.insert" || len(operation.Servers) != 1 || operation.Servers[0] != "https://example.invalid" {
				t.Fatalf("lost native identity/server or incomplete evidence: %+v", operation)
			}
			if name == "upload-only" {
				if operation.Path != "" || operation.Method != "" {
					t.Fatal("upload endpoint became normal HTTP contract")
				}
			} else if operation.Path != "/api/v1/files" || operation.Method != "POST" {
				t.Fatalf("lost declared normal endpoint: %+v", operation)
			}
			hasRequest := name != "upload-only" && name != "no-request"
			if (len(operation.Inputs) == 1) != hasRequest {
				t.Fatalf("lost/fabricated request: %+v", operation.Inputs)
			}
			for _, input := range operation.Inputs {
				if input.Required || input.Schema.Known {
					t.Fatal("unproved mandatory body/schema")
				}
			}
		})
	}
	source := `{"discoveryVersion":"v1","name":"files","version":"v1","rootUrl":"","baseUrl":"","methods":{"get":{"id":"get","httpMethod":"GET","path":"files"}}}`
	table, err := BuildOperationShapeTable(context.Background(), OperationShapeOptions{Sources: []ShapeSourceInput{{ID: "empty-server", OperationSourceInput: OperationSourceInput{Kind: OperationSourceGoogleDiscovery, Content: []byte(source)}}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(table.Operations[0].Servers) != 0 {
		t.Fatal("empty declaration enabled parser default server")
	}
}

func TestOperationShapesUnsupportedSerializationAndFormats(t *testing.T) {
	for name, change := range map[string]func(map[string]any){
		"explode-only": func(operation map[string]any) { operation["parameters"].([]any)[0].(map[string]any)["explode"] = true },
		"response-xml": func(operation map[string]any) {
			response := mapValue(mapValue(operation["responses"])["200"])
			response["content"] = map[string]any{"application/xml": map[string]any{"schema": map[string]any{"type": "object"}}}
		},
		"response-encoding": func(operation map[string]any) {
			media := mapValue(mapValue(mapValue(mapValue(operation["responses"])["200"])["content"])["application/json"])
			media["encoding"] = map[string]any{}
		},
		"request-form": func(operation map[string]any) {
			operation["requestBody"] = map[string]any{"content": map[string]any{"application/x-www-form-urlencoded": map[string]any{"schema": map[string]any{"type": "object"}}}}
		},
		"header-explode": func(operation map[string]any) {
			response := mapValue(mapValue(operation["responses"])["200"])
			response["headers"] = map[string]any{"X-Test": map[string]any{"schema": map[string]any{"type": "array", "items": map[string]any{"type": "string"}}, "explode": true}}
		},
	} {
		t.Run(name, func(t *testing.T) {
			fixture := completeShapeFixture()
			operation := mapValue(mapValue(mapValue(fixture["paths"])["/pets/{id}"])["get"])
			change(operation)
			table, err := BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
			if err != nil {
				t.Fatal(err)
			}
			if table.Operations[0].Complete {
				t.Fatalf("unsupported serialization proved complete: %+v", table.Operations[0])
			}
			if (name == "response-xml" || name == "response-encoding") && table.Operations[0].Outputs[0].Schema.Known {
				t.Fatal("unsupported representation proved schema")
			}
		})
	}
	for _, format := range []string{"int32", "int64", "float", "double", "byte", "binary", "password", "custom-format"} {
		t.Run(format, func(t *testing.T) {
			fixture := completeShapeFixture()
			operation := mapValue(mapValue(mapValue(fixture["paths"])["/pets/{id}"])["get"])
			operation["parameters"].([]any)[0].(map[string]any)["schema"] = map[string]any{"type": "integer", "format": format}
			table, err := BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
			if err != nil {
				t.Fatal(err)
			}
			if table.Operations[0].Inputs[0].Schema.Known || table.Operations[0].Complete {
				t.Fatal("unenforced format proved known")
			}
		})
	}
}

func TestOperationShapesSwaggerSerialization(t *testing.T) {
	for _, change := range []string{"none", "collectionFormat", "xml", "missing-media"} {
		t.Run(change, func(t *testing.T) {
			fixture := map[string]any{"swagger": "2.0", "info": map[string]any{"title": "Swagger", "version": "1"}, "host": "example.invalid", "schemes": []any{"https"}, "security": []any{}, "produces": []any{"application/json"}, "paths": map[string]any{"/things": map[string]any{"get": map[string]any{"parameters": []any{map[string]any{"name": "ids", "in": "query", "type": "array", "items": map[string]any{"type": "string"}}}, "responses": map[string]any{"200": map[string]any{"description": "Things", "schema": map[string]any{"type": "object"}}}}}}}
			operation := mapValue(mapValue(mapValue(fixture["paths"])["/things"])["get"])
			if change == "collectionFormat" {
				operation["parameters"].([]any)[0].(map[string]any)["collectionFormat"] = "pipes"
			}
			if change == "xml" {
				fixture["produces"] = []any{"application/xml"}
			}
			if change == "missing-media" {
				delete(fixture, "produces")
			}
			table, err := BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
			if err != nil {
				t.Fatal(err)
			}
			if table.Operations[0].Complete != (change == "none") {
				t.Fatalf("wrong Swagger completeness %s: %+v", change, table.Operations[0])
			}
		})
	}
}

func TestOperationShapesOpenAPIMetadataCompatibility(t *testing.T) {
	source := `openapi: 3.1.0
info: {title: YAML, version: '1'}
servers: [{url: https://example.invalid}]
security: []
paths:
  x-list: [a, b]
  x-scalar: text
  x-null: null
  /things:
    x-number: 42
    trace:
      responses:
        200: {description: Things, content: {application/json: {schema: {type: object}}}}
`
	options := OperationShapeOptions{Sources: []ShapeSourceInput{{ID: "yaml", OperationSourceInput: OperationSourceInput{Kind: OperationSourceOpenAPI, Content: []byte(source)}}}}
	table, err := BuildOperationShapeTable(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	if len(table.Operations) != 1 || table.Operations[0].Method != "TRACE" || table.Operations[0].Path != "/things" || !table.Operations[0].Outputs[0].Schema.Known {
		t.Fatalf("lost valid metadata: %+v", table)
	}
	for _, bad := range []string{
		source + "        '200': {description: duplicate}\n",
		source + "        200.0: {description: normalized-duplicate}\n",
		source + "        true: {description: boolean}\n",
		source + "        200.5: {description: fraction}\n",
		"openapi: 3.1.0\ninfo: {title: Numeric, version: '1'}\npaths:\n  200: {}\n",
		"openapi: 3.1.0\ninfo: {title: Numeric, version: '1'}\npaths: {}\ncomponents:\n  schemas:\n    Test:\n      type: object\n      properties:\n        200: {type: string}\n",
	} {
		options.Sources[0].Content = []byte(bad)
		if _, err := BuildOperationShapeTable(context.Background(), options); err == nil {
			t.Fatal("coerced arbitrary or colliding YAML key")
		}
	}
}

func TestOperationShapesNumericSwaggerVersion(t *testing.T) {
	for _, version := range []string{`"2.0"`, `2.0`} {
		source := fmt.Sprintf(`{"swagger":%s,"info":{"title":"Swagger","version":"1"},"host":"example.invalid","schemes":["https"],"basePath":"/v1","produces":["application/json"],"securityDefinitions":{"basic":{"type":"basic"},"oauth":{"type":"oauth2","flow":"implicit","authorizationUrl":"https://example.invalid/auth","scopes":{"read":"read"}}},"security":[{"basic":[],"oauth":["read"]}],"paths":{"/things":{"get":{"responses":{"200":{"description":"Things","schema":{"type":"object"}}}}}}}`, version)
		table, err := BuildOperationShapeTable(context.Background(), OperationShapeOptions{Sources: []ShapeSourceInput{{ID: "swagger", OperationSourceInput: OperationSourceInput{Kind: OperationSourceOpenAPI, Content: []byte(source)}}}})
		if err != nil {
			t.Fatal(err)
		}
		operation := table.Operations[0]
		if !operation.Complete || !operation.Security.Known || len(operation.Servers) != 1 || operation.Servers[0] != "https://example.invalid/v1" || len(operation.Security.Alternatives[0].Requirements) != 2 {
			t.Fatalf("inconsistent Swagger version: %+v", operation)
		}
	}
}

func TestOperationShapesGraphQLCompatibleRepeats(t *testing.T) {
	for _, selection := range []string{`field field`, `same: field same: field`, `field(id: 1) field(id: 1)`, `field(id: "1") field(id: "1")`, `field { id } field { id }`, `field { id } field { name }`, `field(id: 1) field: field(id: 1)`, `field(a: 1, b: 2) field(b: 2, a: 1)`, `field(input: {a: 1, b: [2, 3]}) field(input: {b: [2, 3], a: 1})`, `field(text: "a") field(text: "\u0061")`, `field(text: "n") field(text: """n""")`} {
		options := OperationShapeOptions{Sources: []ShapeSourceInput{{ID: "graphql", OperationSourceInput: OperationSourceInput{Kind: OperationSourceGraphQL, Content: []byte("query Q { " + selection + " }")}}}}
		table, err := BuildOperationShapeTable(context.Background(), options)
		if err != nil {
			t.Fatalf("compatible repeated selection %s refused: %v", selection, err)
		}
		if len(table.Operations[0].Outputs) != 1 {
			t.Fatal("compatible repeated output not merged")
		}
	}
	for _, selection := range []string{`same: first same: second`, `field(id: 1) field(id: 2)`, `field(id: "1") field(id: 1)`, `field { same: first } field { same: second }`, `field(text: "\n") field(text: "n")`, `field(text: "\t") field(text: "t")`, `field(text: "\u0061") field(text: "u0061")`, `field(input: [1,2]) field(input: [2,1])`} {
		options := OperationShapeOptions{Sources: []ShapeSourceInput{{ID: "graphql", OperationSourceInput: OperationSourceInput{Kind: OperationSourceGraphQL, Content: []byte("query Q { " + selection + " }")}}}}
		if _, err := BuildOperationShapeTable(context.Background(), options); err == nil {
			t.Fatalf("unproved conflicting response selection %s accepted", selection)
		}
	}
}

func TestOperationShapesYAMLNestedResponseContexts(t *testing.T) {
	main := `openapi: 3.1.0
info: {title: Nested, version: '1'}
paths:
  /things:
    get:
      responses: {200: {description: Things}}
`
	for _, extra := range []string{
		`      callbacks:
        event:
          '{$request.body#/url}':
            post:
              responses: {200: {description: Event}}
`,
		`webhooks:
  event:
    post:
      responses: {200: {description: Event}}
`,
		`components:
  callbacks:
    event:
      '{$request.body#/url}':
        post:
          responses: {200: {description: Event}}
`,
		`components:
  pathItems:
    event:
      post:
        responses: {200.0: {description: Event}}
`,
	} {
		options := OperationShapeOptions{Sources: []ShapeSourceInput{{ID: "yaml", OperationSourceInput: OperationSourceInput{Kind: OperationSourceOpenAPI, Content: []byte(main + extra)}}}}
		table, err := BuildOperationShapeTable(context.Background(), options)
		if err != nil {
			t.Fatal("legal nested response keys discarded source", err)
		}
		if len(table.Operations) != 1 || table.Operations[0].Path != "/things" || table.Operations[0].Complete {
			t.Fatal("nested response recognition expanded projection")
		}
	}
	for _, key := range []string{`!!float 200e1000000`, `!!int ` + strings.Repeat("2", 10000)} {
		options := OperationShapeOptions{Sources: []ShapeSourceInput{{ID: "yaml", OperationSourceInput: OperationSourceInput{Kind: OperationSourceOpenAPI, Content: []byte(strings.Replace(main, "200:", key+":", 1))}}}}
		if _, err := BuildOperationShapeTable(context.Background(), options); err == nil {
			t.Fatal("unbounded malformed response code accepted")
		}
	}
}

func TestOperationShapesGraphQLQuotedSelectionDelimiters(t *testing.T) {
	for _, selection := range []string{`a: field(text: "(") b: field`, `a: field @skip(if: "(") b: field`, `a: field { child(text: "}") } b: field`, `a: field(text: "$") b: field`, `a: field(text: """a\"""b""") b: field`} {
		options := OperationShapeOptions{Sources: []ShapeSourceInput{{ID: "graphql", OperationSourceInput: OperationSourceInput{Kind: OperationSourceGraphQL, Content: []byte("query Q { " + selection + " }")}}}}
		table, err := BuildOperationShapeTable(context.Background(), options)
		if err != nil {
			t.Fatalf("quoted selection delimiter corrupted parser: %s %v", selection, err)
		}
		if len(table.Operations[0].Outputs) != 2 || table.Operations[0].Outputs[0].Name != "a" || table.Operations[0].Outputs[1].Name != "b" {
			t.Fatal("quoted delimiter fabricated selection")
		}
	}
}

func TestOperationShapesGraphQLNativeStringAndFragmentProof(t *testing.T) {
	block := "\"\"\"\n  a\n \n  b\n\"\"\""
	for _, selection := range []string{
		`field(text: "\u{61}") field(text: "a")`,
		`field(text: "\uD83D\uDCA9") field(text: "💩")`,
		`field(input: ["]", {child: ["}","("]}])`,
		`field { ... on User { id } } field { User: name }`,
		`field(text: ` + block + `) field(text: "a\n\nb")`,
	} {
		options := OperationShapeOptions{Sources: []ShapeSourceInput{{ID: "graphql", OperationSourceInput: OperationSourceInput{Kind: OperationSourceGraphQL, Content: []byte("query Q { " + selection + " }")}}}}
		table, err := BuildOperationShapeTable(context.Background(), options)
		if err != nil {
			t.Fatal("compatible native value/fragment metadata refused", err)
		}
		if len(table.Operations[0].Outputs) != 1 || table.Operations[0].Complete || table.Operations[0].Outputs[0].Schema.Known {
			t.Fatal("fragment/string repair invented complete schema")
		}
	}
	for _, selection := range []string{
		`field(text: "\uD800") field(text: "�")`,
		`field(text: ` + block + `) field(text: "a\n \nb")`,
	} {
		options := OperationShapeOptions{Sources: []ShapeSourceInput{{ID: "graphql", OperationSourceInput: OperationSourceInput{Kind: OperationSourceGraphQL, Content: []byte("query Q { " + selection + " }")}}}}
		if _, err := BuildOperationShapeTable(context.Background(), options); err == nil {
			t.Fatal("invalid/conflicting native string values merged")
		}
	}
}
