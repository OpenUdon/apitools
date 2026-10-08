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
