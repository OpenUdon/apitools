package apitools

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/OpenUdon/uws/binding"
)

func shapeFamilySources(t *testing.T) []ShapeSourceInput {
	t.Helper()
	names := []string{"openapi.yaml", "google-discovery.json", "aws-smithy.json", "asyncapi.yaml", "graphql.graphql", "openrpc.json", "grpc.proto", "odata.xml"}
	kinds := []OperationSourceKind{OperationSourceOpenAPI, OperationSourceGoogleDiscovery, OperationSourceAWSSmithy, OperationSourceAsyncAPI, OperationSourceGraphQL, OperationSourceOpenRPC, OperationSourceGRPCProtobuf, OperationSourceOData}
	var sources []ShapeSourceInput
	for i, kind := range kinds {
		path := filepath.Join("testdata", "operation-candidates", "v1", "sources", names[i])
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		sources = append(sources, ShapeSourceInput{ID: string(kind), OperationSourceInput: OperationSourceInput{Kind: kind, Content: content}})
	}
	return sources
}

func TestOperationShapesEveryNativeFamily(t *testing.T) {
	sources := shapeFamilySources(t)
	table, err := BuildOperationShapeTable(context.Background(), OperationShapeOptions{Sources: sources})
	if err != nil {
		t.Fatal(err)
	}
	if len(table.Sources) != 8 || len(table.Operations) < 8 {
		t.Fatalf("missing families: sources=%d operations=%d", len(table.Sources), len(table.Operations))
	}
	seen := map[string]bool{}
	for _, operation := range table.Operations {
		seen[operation.Source.Kind] = true
		if operation.Source.Kind != "openapi" && operation.Source.Kind != "google-discovery" && (operation.Method != "" || operation.Path != "" || len(operation.Servers) > 0) {
			t.Fatalf("invented HTTP contract: %+v", operation)
		}
		if operation.Source.Kind != "openapi" && operation.Security.Known || operation.Complete {
			t.Fatal("unqualified completeness or security")
		}
		if operation.Selector.Kind != "ref" || operation.Selector.Value == "" {
			t.Fatal("missing native selector")
		}
		for _, alias := range operation.Aliases {
			resolver, err := binding.NewResolver(table)
			if err != nil {
				t.Fatal(err)
			}
			resolved, err := resolver.Resolve(context.Background(), binding.Binding{Source: operation.Source, SelectorKind: alias.Kind, SelectorValue: alias.Value})
			if err != nil || resolved.Status != binding.Resolved || resolved.Shape.Selector.Key != operation.Selector.Key {
				t.Fatalf("native alias did not resolve: %v %v", resolved.Status, err)
			}
		}
	}
	for _, source := range sources {
		if !seen[string(source.Kind)] {
			t.Fatalf("missing %s", source.Kind)
		}
		digest := sha256.Sum256(source.Content)
		found := false
		for _, identity := range table.Sources {
			if identity.ID == source.ID {
				found = true
				if identity.SHA256 != hex.EncodeToString(digest[:]) {
					t.Fatal("digest describes changed bytes")
				}
			}
		}
		if !found {
			t.Fatal("missing identity")
		}
	}
	for _, op := range table.Operations {
		if op.Source.Kind == "openapi" && op.Path == "/pets" && op.Method == "POST" {
			if len(op.Inputs) == 0 || !op.Inputs[0].Schema.Known || len(op.Outputs) == 0 || !op.Outputs[0].Schema.Known {
				t.Fatalf("lost declared OpenAPI schemas: %+v", op)
			}
		}
	}
}

func TestOperationShapesRefuseUnsafeOrPartialTables(t *testing.T) {
	good := shapeFamilySources(t)[0]
	for name, options := range map[string]OperationShapeOptions{
		"missing":           {},
		"url-only":          {Sources: []ShapeSourceInput{{ID: "api", OperationSourceInput: OperationSourceInput{Kind: OperationSourceOpenAPI, URL: "https://example.invalid/api"}}}},
		"duplicate":         {Sources: []ShapeSourceInput{good, good}},
		"byte-limit":        {Sources: []ShapeSourceInput{good}, MaxBytes: 1},
		"negative":          {Sources: []ShapeSourceInput{good}, MaxOperations: -1},
		"raised":            {Sources: []ShapeSourceInput{good}, MaxOperations: binding.MaxOperations + 1},
		"malformed-sibling": {Sources: []ShapeSourceInput{good, {ID: "bad", OperationSourceInput: OperationSourceInput{Kind: OperationSourceOpenRPC, Content: []byte(`{"openrpc":`)}}}},
	} {
		t.Run(name, func(t *testing.T) {
			table, err := BuildOperationShapeTable(context.Background(), options)
			if !errors.Is(err, ErrOperationShapeTable) || len(table.Operations) != 0 || len(table.Sources) != 0 {
				t.Fatalf("partial/positive result: %+v %v", table, err)
			}
		})
	}
	root := t.TempDir()
	path := filepath.Join(root, "api.yaml")
	if err := os.WriteFile(path, good.Content, 0600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "link.yaml")
	if err := os.Symlink(path, link); err != nil {
		t.Fatal(err)
	}
	good.Content = nil
	good.Path = link
	if _, err := BuildOperationShapeTable(context.Background(), OperationShapeOptions{Sources: []ShapeSourceInput{good}}); !errors.Is(err, ErrOperationShapeTable) {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := BuildOperationShapeTable(ctx, OperationShapeOptions{Sources: []ShapeSourceInput{good}}); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
