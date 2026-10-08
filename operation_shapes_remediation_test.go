package apitools

import (
	"context"
	"fmt"
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
