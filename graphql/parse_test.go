package graphql

import (
	"bytes"
	"strings"
	"testing"

	"github.com/OpenUdon/apitools/internal/sourceguard"
)

func TestParseIntrospectionPreservesSchemaRootOperations(t *testing.T) {
	model, err := Parse([]byte(`{
  "data": {
    "__schema": {
      "queryType": {"name": "Query"},
      "mutationType": {"name": "Mutation"},
      "types": [
        {
          "kind": "OBJECT",
          "name": "Query",
          "fields": [
            {
              "name": "node",
              "description": "Fetch a node.",
              "args": [{"name": "id", "type": {"kind": "NON_NULL", "ofType": {"kind": "SCALAR", "name": "ID"}}}],
              "type": {"kind": "OBJECT", "name": "Node"}
            }
          ]
        },
        {
          "kind": "OBJECT",
          "name": "Mutation",
          "fields": [
            {"name": "createIssue", "type": {"kind": "OBJECT", "name": "Issue"}}
          ]
        }
      ]
    }
  }
}`))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if model.SourceKind != ArtifactKindIntrospectionJSON {
		t.Fatalf("source kind = %q, want introspection", model.SourceKind)
	}
	summaries := model.OperationSummaries()
	if len(summaries) != 2 {
		t.Fatalf("operation summaries = %d, want 2: %#v", len(summaries), summaries)
	}
	if summaries[0].ID != "mutation.createIssue" || summaries[1].ID != "query.node" {
		t.Fatalf("operation ids = %#v, want deterministic mutation/query ids", summaries)
	}
	query, ok := model.OperationByID("query.node")
	if !ok {
		t.Fatalf("missing query.node operation")
	}
	if query.Selector != "#/schema/Query/fields/node" || query.FieldType.Display != "Node" {
		t.Fatalf("query metadata = %#v", query)
	}
	target, err := model.ResolveSelector("query.node")
	if err != nil {
		t.Fatalf("ResolveSelector() error = %v", err)
	}
	if target.Kind != SelectorKindOperation || target.Summary.SourceOperationRef != "#/schema/Query/fields/node" {
		t.Fatalf("selector target = %#v", target)
	}
}

func TestParseSDLAndOperationDocument(t *testing.T) {
	model, err := Parse([]byte(`
schema {
  query: RootQuery
  mutation: RootMutation
}

type RootQuery {
  """Find one issue."""
  issue(id: ID!): Issue
}

type RootMutation {
  createIssue(input: IssueInput!): Issue
}

input IssueInput {
  title: String!
}

query GetIssue($id: ID!) {
  issue(id: $id) { id title }
}
`))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if model.SourceKind != ArtifactKindSDL {
		t.Fatalf("source kind = %q, want SDL", model.SourceKind)
	}
	for _, id := range []string{"mutation.createIssue", "query.GetIssue", "query.issue"} {
		if _, ok := model.OperationByID(id); !ok {
			t.Fatalf("missing operation %q in %#v", id, model.OperationSummaries())
		}
	}
	op, _ := model.OperationByID("query.GetIssue")
	if len(op.Variables) != 1 || op.Variables[0].Name != "id" || op.Variables[0].Type.Display != "ID!" {
		t.Fatalf("operation variables = %#v", op.Variables)
	}
	if len(op.SelectionNames) != 1 || op.SelectionNames[0] != "issue" {
		t.Fatalf("selection names = %#v", op.SelectionNames)
	}
}

func TestDetectLocalArtifact(t *testing.T) {
	detection, ok := DetectLocalArtifact([]byte(`query FindViewer { viewer { login } }`), "ops.graphql")
	if !ok {
		t.Fatalf("DetectLocalArtifact() did not detect operation document")
	}
	if detection.Kind != ArtifactKindOperationDocument || detection.Format != "graphql" {
		t.Fatalf("detection = %#v", detection)
	}
	if _, ok := DetectLocalArtifact([]byte(`{"title":"not graphql"}`), "spec.json"); ok {
		t.Fatalf("non-GraphQL JSON detected")
	}
}

func TestParseRejectsMalformedGraphQL(t *testing.T) {
	if _, err := Parse([]byte(`not graphql`)); err == nil {
		t.Fatalf("Parse() error = nil, want malformed GraphQL error")
	}
	if _, err := Parse([]byte(`query`)); err == nil {
		t.Fatalf("Parse() error = nil, want incomplete operation error")
	}
}

func TestParseRejectsResourceLimitViolations(t *testing.T) {
	t.Run("document bytes", func(t *testing.T) {
		_, err := Parse(bytes.Repeat([]byte{'x'}, sourceguard.MaxDocumentBytes+1))
		if err == nil || !strings.Contains(err.Error(), "maximum size") {
			t.Fatalf("Parse() error = %v, want maximum size", err)
		}
	})
	t.Run("text type nesting", func(t *testing.T) {
		input := `query Q($v: ` + strings.Repeat(`[`, sourceguard.MaxNestingDepth+1) + `String` + strings.Repeat(`]`, sourceguard.MaxNestingDepth+1) + `) { viewer }`
		_, err := Parse([]byte(input))
		if err == nil || !strings.Contains(err.Error(), "nesting") {
			t.Fatalf("Parse() error = %v, want nesting limit", err)
		}
	})
	t.Run("introspection nesting", func(t *testing.T) {
		ref := `{"kind":"SCALAR","name":"String"}`
		for range sourceguard.MaxNestingDepth {
			ref = `{"kind":"LIST","ofType":` + ref + `}`
		}
		input := `{"data":{"__schema":{"types":[{"kind":"OBJECT","name":"Query","fields":[{"name":"x","type":` + ref + `}]}]}}}`
		input = strings.ReplaceAll(input, `\"`, `"`)
		_, err := Parse([]byte(input))
		if err == nil || !strings.Contains(err.Error(), "JSON nesting") {
			t.Fatalf("Parse() error = %v, want JSON nesting limit", err)
		}
	})
}

func FuzzParse(f *testing.F) {
	for _, seed := range [][]byte{
		[]byte(`type Query { viewer: User }`),
		[]byte(`query Viewer { viewer { id } }`),
		[]byte(`{"data":{"__schema":{"types":[]}}}`),
		[]byte(strings.Repeat(`[`, sourceguard.MaxNestingDepth+1)),
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = Parse(data)
	})
}

func TestDefaultsAndResponseKeys(t *testing.T) {
	model, err := Parse([]byte(`type Query { field(first: Int! = 10, text: String! = "", required: Int!): String } input Filter { limit: Int! = 10 } query Q($first: Int! = 10, $text: String! = "", $required: Int!) { a: field b: field field a: other }`))
	if err != nil {
		t.Fatal(err)
	}
	operation, _ := model.OperationByID("query.Q")
	if len(operation.Variables) != 3 {
		t.Fatal(operation.Variables)
	}
	for i, variable := range operation.Variables {
		if variable.Required != (i == 2) || variable.HasDefault != (i != 2) || !variable.Type.Required {
			t.Fatalf("wrong variable default: %+v", variable)
		}
	}
	if operation.Variables[0].DefaultValue != "10" {
		t.Fatal("default swallowed next variable")
	}
	var field *Field
	var input *Argument
	for _, typ := range model.Types {
		if typ.Name == "Query" {
			field = typ.Fields[0]
		}
		if typ.Name == "Filter" {
			input = typ.InputFields[0]
		}
	}
	if field == nil || input == nil {
		t.Fatal("missing schema types")
	}
	for i, argument := range field.Args {
		if argument.Required != (i == 2) || argument.HasDefault != (i != 2) || !argument.Type.Required {
			t.Fatalf("wrong argument default: %+v", argument)
		}
	}
	if input.Required {
		t.Fatal("defaulted input-object field required")
	}
	if len(operation.Selections) != 4 || len(operation.SelectionNames) != 2 || operation.Selections[0].FieldName != "field" || operation.Selections[0].ResponseKey != "a" || operation.Selections[1].ResponseKey != "b" || operation.Selections[3].ResponseKey != "a" {
		t.Fatalf("lost response keys or underlying fields: %+v", operation)
	}
	if field.Name != "field" || field.Selector != "#/schema/Query/fields/field" || operation.Selector != "#/operations/query.Q" {
		t.Fatal("changed native fields/selectors")
	}
}

func TestIntrospectionDefaultPresence(t *testing.T) {
	ref := map[string]any{"kind": "NON_NULL", "ofType": map[string]any{"kind": "SCALAR", "name": "Int"}}
	for _, value := range []any{nil, "10", "null", ""} {
		argument := parseIntrospectionArgument(map[string]any{"name": "first", "type": ref, "defaultValue": value})
		if argument.Required != (value == nil) || argument.HasDefault != (value != nil) || !argument.Type.Required {
			t.Fatalf("wrong default presence: %+v", argument)
		}
	}
}

func TestDefaultQuotedDelimiters(t *testing.T) {
	for _, value := range []string{"$", "@", "}", ")", "{", "["} {
		model, err := Parse([]byte(`query Q($value: String! = "` + value + `", $required: Int!) { field }`))
		if err != nil {
			t.Fatal(err)
		}
		operation, _ := model.OperationByID("query.Q")
		if len(operation.Variables) != 2 || operation.Variables[0].DefaultValue != value || !operation.Variables[0].HasDefault || operation.Variables[0].Required || !operation.Variables[1].Required {
			t.Fatalf("quoted default changed: %+v", operation.Variables)
		}
	}
}

func TestSelectionBoundariesAndNativeValues(t *testing.T) {
	for _, source := range []string{
		`query Q { field(arg:1`,
		`query Q { field { child`,
		`query Q { field(input: [1) }`,
		`query Q { ` + strings.Repeat(`f @d(]) { `, 5000) + `id` + strings.Repeat(`}`, 5001),
		`query Q { field(text: "\uD800") }`,
		`query Q { field(text: "\uDC00") }`,
		`query Q { field(text: "\u{110000}") }`,
		`query Q { field(text: "\u{D800}") }`,
		`query Q { field(text: "\x61") }`,
	} {
		if _, err := Parse([]byte(source)); err == nil {
			t.Fatal("malformed/unbounded native value accepted")
		}
	}
	for _, source := range []string{
		`query Q { field(input: ["]", {child: ["}","("]}]) }`,
		`query Q { field(text: "\u{61}") }`,
		`query Q { field(text: "\uD83D\uDCA9") }`,
	} {
		if _, err := Parse([]byte(source)); err != nil {
			t.Fatal("valid native structured value refused", err)
		}
	}
}
