package apitools

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/OpenUdon/uws/binding"
)

func remediationShapeSources(t *testing.T) []ShapeSourceInput {
	t.Helper()
	sources := shapeFamilySources(t)
	for i := range sources {
		name := ""
		switch sources[i].Kind {
		case OperationSourceOpenAPI:
			name = "openapi.yaml"
		case OperationSourceGoogleDiscovery:
			name = "google-discovery.json"
		case OperationSourceAWSSmithy:
			name = "aws-smithy.json"
		case OperationSourceAsyncAPI:
			name = "asyncapi.json"
		case OperationSourceGraphQL:
			name = "graphql.graphql"
		}
		if name != "" {
			data, err := os.ReadFile(filepath.Join("testdata", "operation-shapes", "m83", "sources", name))
			if err != nil {
				t.Fatal(err)
			}
			sources[i].Content = data
		}
	}
	return sources
}

func TestOperationShapesRemediationGolden(t *testing.T) {
	sources := remediationShapeSources(t)
	table, err := BuildOperationShapeTable(context.Background(), OperationShapeOptions{Sources: sources})
	if err != nil {
		t.Fatal(err)
	}
	golden, err := os.ReadFile(filepath.Join("testdata", "operation-shapes", "m83", "eight-families.json"))
	if err != nil {
		t.Fatal(err)
	}
	data := mustShapeBytes(t, table)
	if !bytes.Equal(bytes.TrimSpace(golden), data) {
		t.Fatal("corrected table differs from reviewed fixture")
	}
	if len(table.Sources) != 8 || len(table.Operations) != 13 {
		t.Fatalf("lost family operations: %d/%d", len(table.Sources), len(table.Operations))
	}
	for _, source := range sources {
		digest := sha256.Sum256(source.Content)
		found := false
		for _, identity := range table.Sources {
			if identity.ID == source.ID {
				found = identity.SHA256 == hex.EncodeToString(digest[:])
			}
		}
		if !found {
			t.Fatal("lost raw source identity")
		}
	}
	resolver, err := binding.NewResolver(table)
	if err != nil {
		t.Fatal(err)
	}
	for _, operation := range table.Operations {
		resolution, err := resolver.Resolve(context.Background(), binding.Binding{Source: operation.Source, SelectorKind: "ref", SelectorValue: operation.Selector.Value})
		if err != nil || resolution.Status != binding.Resolved || resolution.Shape.Selector.Key != operation.Selector.Key {
			t.Fatal("lost native selector")
		}
		if operation.Complete || operation.Source.Kind != "openapi" && operation.Security.Known {
			t.Fatal("invented complete native proof")
		}
		switch operation.Source.Kind {
		case "asyncapi":
			if operation.Selector.Value == "#/operations/received" && len(operation.Outputs) != 1 || operation.Selector.Value == "#/operations/sent" && len(operation.Inputs) != 1 {
				t.Fatal("wrong AsyncAPI 2 direction")
			}
		case "graphql":
			if operation.Selector.Value == "#/operations/query.Q" && (len(operation.Outputs) != 2 || operation.Outputs[0].Name != "a" || operation.Outputs[1].Name != "b" || operation.Inputs[0].Required) {
				t.Fatal("lost GraphQL default/alias evidence")
			}
		case "aws-smithy":
			if len(operation.Inputs) != 1 || operation.Inputs[0].Name != "Id" {
				t.Fatal("fabricated literal member")
			}
		case "google-discovery":
			if operation.Path != "/api/v1/files" || operation.Inputs[0].Required || operation.Inputs[0].Schema.Known {
				t.Fatal("fabricated upload/mandatory body proof")
			}
		case "openapi":
			if operation.Method != "TRACE" || operation.Inputs[0].Schema.Known || len(operation.Security.Alternatives) != 2 {
				t.Fatal("lost unknown format/symbolic security/TRACE")
			}
		}
	}
	for i, j := 0, len(sources)-1; i < j; i, j = i+1, j-1 {
		sources[i], sources[j] = sources[j], sources[i]
	}
	other, err := BuildOperationShapeTable(context.Background(), OperationShapeOptions{Sources: sources})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, mustShapeBytes(t, other)) {
		t.Fatal("source order changed corrected bytes")
	}
	if err := VerifyOperationShapeTable(context.Background(), OperationShapeOptions{Sources: sources}, table); err != nil {
		t.Fatal(err)
	}
}

func TestOperationShapesRemediationNoFetch(t *testing.T) {
	var calls atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls.Add(1) }))
	defer server.Close()
	sources := remediationShapeSources(t)
	for i := range sources {
		sources[i].URL = server.URL + "/source?private=redacted"
	}
	table, err := BuildOperationShapeTable(context.Background(), OperationShapeOptions{Sources: sources})
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifyOperationShapeTable(context.Background(), OperationShapeOptions{Sources: sources}, table); err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 0 {
		t.Fatal("source tooling fetched provenance")
	}
}
