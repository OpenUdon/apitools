package apitools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// nestedObjectSchema returns an object schema with n string properties.
func nestedObjectSchema(n int) map[string]any {
	properties := map[string]any{}
	for i := 0; i < n; i++ {
		properties[fmt.Sprintf("p%03d", i)] = map[string]any{"type": "string"}
	}
	return map[string]any{"type": "object", "properties": properties}
}

func openAPIFieldsInventory(t *testing.T, request, response map[string]any) OperationSummary {
	t.Helper()
	operation := map[string]any{
		"operationId": "probe",
		"responses": map[string]any{"200": map[string]any{
			"description": "ok",
			"content":     map[string]any{"application/json": map[string]any{"schema": response}},
		}},
	}
	if request != nil {
		operation["requestBody"] = map[string]any{"content": map[string]any{"application/json": map[string]any{"schema": request}}}
	}
	document, err := json.Marshal(map[string]any{
		"openapi": "3.0.0",
		"info":    map[string]any{"title": "Fields", "version": "1"},
		"paths":   map[string]any{"/probe": map[string]any{"post": operation}},
	})
	if err != nil {
		t.Fatal(err)
	}
	inventory, err := BuildOperationInventory(context.Background(), InventoryOptions{
		Documents: []InventoryDocument{{Name: "fields", Content: document}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(inventory.Operations) != 1 {
		t.Fatalf("operations = %#v", inventory.Operations)
	}
	return inventory.Operations[0]
}

func hasIssue(op OperationSummary, code string) bool {
	for _, issue := range op.ReadinessIssues {
		if issue.Code == code {
			return true
		}
	}
	return false
}

func requireTopLevelKept(t *testing.T, fields []RequestFieldSummary, names ...string) {
	t.Helper()
	if len(fields) > 60 {
		t.Fatalf("fields = %d, want at most 60", len(fields))
	}
	got := responseFieldNames(fields)
	if !sort.StringsAreSorted(got) {
		t.Fatalf("fields are not sorted by path: %v", got)
	}
	have := map[string]bool{}
	for _, name := range got {
		have[name] = true
	}
	for _, name := range names {
		if !have[name] {
			t.Fatalf("top-level field %q was dropped; fields = %v", name, got)
		}
	}
}

func TestOpenAPIResponseTruncationKeepsTopLevelFields(t *testing.T) {
	// "a" sorts first and holds 70 nested fields, so a depth-first cut at 60
	// would drop every later top-level field.
	properties := map[string]any{"a": nestedObjectSchema(70)}
	for _, name := range []string{"id", "location", "name", "storageClass"} {
		properties[name] = map[string]any{"type": "string"}
	}
	op := openAPIFieldsInventory(t, nil, map[string]any{"type": "object", "properties": properties})
	if op.ResponseBody == nil {
		t.Fatal("missing response body")
	}
	requireTopLevelKept(t, op.ResponseBody.Fields, "a", "id", "location", "name", "storageClass")
	if len(op.ResponseBody.Fields) != 60 {
		t.Fatalf("fields = %d, want 60", len(op.ResponseBody.Fields))
	}
	if !hasIssue(op, "schema.field_summary_truncated") {
		t.Fatalf("truncation issue missing: %#v", op.ReadinessIssues)
	}
	// The kept nested fields are the first by path.
	if names := responseFieldNames(op.ResponseBody.Fields); names[1] != "a.p000" || names[55] != "a.p054" {
		t.Fatalf("nested fields not selected by path: %v", names)
	}
}

func TestOpenAPIResponseOverSixtyTopLevelFieldsKeepsFirstSixtyByPath(t *testing.T) {
	op := openAPIFieldsInventory(t, nil, nestedObjectSchema(75))
	fields := op.ResponseBody.Fields
	if len(fields) != 60 {
		t.Fatalf("fields = %d, want 60", len(fields))
	}
	if fields[0].Path != "p000" || fields[59].Path != "p059" {
		t.Fatalf("first/last = %q/%q", fields[0].Path, fields[59].Path)
	}
	if !hasIssue(op, "schema.field_summary_truncated") {
		t.Fatalf("truncation issue missing: %#v", op.ReadinessIssues)
	}
}

func TestOpenAPIResponseUnderBudgetIsUnchanged(t *testing.T) {
	op := openAPIFieldsInventory(t, nil, map[string]any{"type": "object", "properties": map[string]any{
		"name": map[string]any{"type": "string"},
		"meta": map[string]any{"type": "object", "properties": map[string]any{"x": map[string]any{"type": "string"}}},
	}})
	if got := strings.Join(responseFieldNames(op.ResponseBody.Fields), ","); got != "meta,meta.x,name" {
		t.Fatalf("fields = %s", got)
	}
	if hasIssue(op, "schema.field_summary_truncated") {
		t.Fatalf("unexpected truncation issue: %#v", op.ReadinessIssues)
	}
}

func TestOpenAPIRequestFieldsKeepDepthFirstCut(t *testing.T) {
	properties := map[string]any{"a": nestedObjectSchema(70), "zeta": map[string]any{"type": "string"}}
	op := openAPIFieldsInventory(t, map[string]any{"type": "object", "properties": properties}, nestedObjectSchema(2))
	fields := op.RequestBody.Fields
	if len(fields) != 60 {
		t.Fatalf("request fields = %d, want 60", len(fields))
	}
	for _, field := range fields {
		if field.Path == "zeta" {
			t.Fatalf("request-field summary must keep the depth-first cut; fields = %v", responseFieldNames(fields))
		}
	}
}

func TestGoogleDiscoveryStorageResponseKeepsTopLevelFields(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("testdata", "inventory", "google-storage-discovery.json"))
	if err != nil {
		t.Fatal(err)
	}
	inventory, err := BuildAPISourceOperationInventory(context.Background(), APISourceInventoryOptions{
		Documents: []APISourceDocument{{Kind: APISourceKindGoogleDiscovery, Name: "storage", Content: content}},
	})
	if err != nil {
		t.Fatal(err)
	}
	var op *OperationSummary
	for i := range inventory.Operations {
		if inventory.Operations[i].OperationID == "storage.buckets.get" {
			op = &inventory.Operations[i]
		}
	}
	if op == nil || op.ResponseBody == nil {
		t.Fatalf("storage.buckets.get response missing: %#v", op)
	}
	requireTopLevelKept(t, op.ResponseBody.Fields, "name", "id", "location", "storageClass")
	if len(op.ResponseBody.Fields) != 60 {
		t.Fatalf("fields = %d, want 60", len(op.ResponseBody.Fields))
	}
	if !hasIssue(*op, "schema.field_summary_truncated") {
		t.Fatalf("truncation issue missing: %#v", op.ReadinessIssues)
	}
}

func TestGoogleDiscoveryResponseUnderBudgetIsUnchanged(t *testing.T) {
	content := `{
  "discoveryVersion": "v1", "name": "svc", "version": "v1",
  "rootUrl": "https://svc.example/", "servicePath": "v1/",
  "schemas": {"Thing": {"id": "Thing", "type": "object", "properties": {"name": {"type": "string"}, "id": {"type": "string"}}}},
  "resources": {"things": {"methods": {"get": {"id": "svc.things.get", "path": "things", "httpMethod": "GET", "response": {"$ref": "Thing"}}}}}
}`
	inventory, err := BuildAPISourceOperationInventory(context.Background(), APISourceInventoryOptions{
		Documents: []APISourceDocument{{Kind: APISourceKindGoogleDiscovery, Name: "svc", Content: []byte(content)}},
	})
	if err != nil {
		t.Fatal(err)
	}
	op := inventory.Operations[0]
	if got := strings.Join(responseFieldNames(op.ResponseBody.Fields), ","); got != "id,name" {
		t.Fatalf("fields = %s", got)
	}
	if hasIssue(op, "schema.field_summary_truncated") {
		t.Fatalf("unexpected truncation issue: %#v", op.ReadinessIssues)
	}
}
