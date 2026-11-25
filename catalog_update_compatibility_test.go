package apitools

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/OpenUdon/uws/binding"
)

func TestCatalogUpdateCompatibilityRequiresReviewedCompleteContracts(t *testing.T) {
	schema := binding.Schema{Known: true, JSON: json.RawMessage(`{"type":"string"}`)}
	shape := binding.OperationShape{
		Source:   binding.Source{SHA256: "held"},
		Selector: binding.Selector{Key: "getRecord"},
		Method:   "GET", Path: "/records/{id}", Complete: true,
		Security: binding.Security{Known: true, Alternatives: []binding.SecurityAlternative{{}}},
		Inputs:   []binding.Input{{Location: "path", Name: "id", Required: true, Schema: schema}},
		Outputs:  []binding.Output{{Location: "body", Name: "value", Schema: schema}},
	}
	before := []binding.OperationShape{shape}
	after := []binding.OperationShape{shape, {Selector: binding.Selector{Key: "newOperation"}}}
	after[0].Source.SHA256 = "new"
	keys := []string{"getRecord"}
	if got := compareReviewedCatalogShapes(before, after, keys); got != catalogUpdateCompatible {
		t.Fatalf("unchanged reviewed shape with another operation = %s", got)
	}
	if got := compareReviewedCatalogShapes(before, after[1:], keys); got != catalogUpdateChanged {
		t.Fatalf("removed reviewed operation = %s", got)
	}
	after[0].Security.Alternatives = []binding.SecurityAlternative{{Requirements: []binding.SecurityRequirement{{Scheme: "token", Type: "http"}}}}
	if got := compareReviewedCatalogShapes(before, after, keys); got != catalogUpdateChanged {
		t.Fatalf("changed auth = %s", got)
	}
	after[0].Security.Known = false
	if got := compareReviewedCatalogShapes(before, after, keys); got != catalogUpdateUnexamined {
		t.Fatalf("unknown auth = %s", got)
	}
}

func TestCatalogUpdateRouteOverlapComparisonIsBounded(t *testing.T) {
	before, after := map[string]any{}, map[string]any{}
	for i := range 250 {
		before[fmt.Sprintf("/existing/%d", i)] = nil
		after[fmt.Sprintf("/additional/%d", i)] = nil
	}
	if removeUnrelatedCatalogUpdatePaths(t.Context(), before, after) {
		t.Fatal("unbounded additional path overlap check")
	}
}

func TestCatalogUpdateCompatibilityChecksNativeHTTPContract(t *testing.T) {
	const original = `{
	  "openapi":"3.0.3",
	  "info":{"title":"Records","version":"1.0"},
	  "servers":[{"url":"https://example.com"}],
	  "security":[{"account":[]}],
	  "components":{"securitySchemes":{"account":{"type":"http","scheme":"basic"}}},
	  "paths":{"/records":{"get":{"operationId":"list","responses":{
	    "200":{"description":"ok","content":{"application/json":{"schema":{"type":"string"}}}}
	  }}}}
	}`
	source := func(content string) OperationSourceInput {
		return OperationSourceInput{Kind: OperationSourceOpenAPI, URL: "https://example.com/spec.json", Content: []byte(content)}
	}
	keys := []string{"#/paths/~1records/get"}
	t.Run("informational version bump", func(t *testing.T) {
		next := strings.Replace(original, `"1.0"`, `"2.0"`, 1)
		if got := compareCatalogUpdate(t.Context(), source(original), source(next), keys); got != catalogUpdateCompatible {
			t.Fatalf("informational version bump = %s", got)
		}
	})
	t.Run("unrelated literal route", func(t *testing.T) {
		next := strings.Replace(original, `"paths":{`, `"paths":{"/unrelated":{"get":{"operationId":"unrelated","responses":{"204":{"description":"empty"}}}},`, 1)
		if got := compareCatalogUpdate(t.Context(), source(original), source(next), keys); got != catalogUpdateCompatible {
			t.Fatalf("unrelated path = %s", got)
		}
	})
	t.Run("authorization with native info reference", func(t *testing.T) {
		old := strings.Replace(original, `"title":"Records","version":"1.0"`, `"title":"Records","version":"1.0","contact":{"x-auth":{"type":"http","scheme":"basic"}}`, 1)
		old = strings.Replace(old, `"account":{"type":"http","scheme":"basic"}`, `"account":{"$ref":"#/info/contact/x-auth"}`, 1)
		next := strings.Replace(old, `"basic"`, `"bearer"`, 1)
		if got := compareCatalogUpdate(t.Context(), source(old), source(next), keys); got == catalogUpdateCompatible {
			t.Fatal("adopted an auth change in a referenced extension")
		}
	})
	for _, tc := range []struct{ name, from, to string }{
		{"HTTP authorization scheme", `"basic"`, `"bearer"`},
		{"response status", `"200"`, `"202"`},
		{"response media", `"application/json"`, `"application/vnd.records+json"`},
		{"unused external ref", `"securitySchemes":`, `"schemas":{"External":{"$ref":"https://example.com/unreviewed.json"}},"securitySchemes":`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			next := strings.Replace(original, tc.from, tc.to, 1)
			if got := compareCatalogUpdate(t.Context(), source(original), source(next), keys); got == catalogUpdateCompatible {
				t.Fatalf("accepted changed native contract: %s", next)
			}
		})
	}
}

func TestCatalogUpdateResourcePreservesQueryAndPublisherFamily(t *testing.T) {
	const index = "https://api.hubapi.com/public/api/spec/v1/specs"
	const release = "https://api.hubapi.com/public/api/spec/v2/specs/release/1/version/"
	for _, tc := range []struct {
		locator, left, right string
		same                 bool
	}{
		{"https://example.com/openapi.json", "https://example.com/openapi.json?service=orders", "https://example.com/openapi.json?service=billing", false},
		{index, release + "2026-03", release + "2026-09", true},
		{index, release + "2026-03", release + "2026-09/other", false},
		{index, release + "2026-03", release + "2026-09?scope=other", false},
		{index, release + "2026-03", strings.Replace(release, "/1/", "/2/", 1) + "2026-09", false},
	} {
		if same := sameCatalogUpdateResource(tc.locator, tc.left, tc.right); same != tc.same {
			t.Fatalf("resource identity %s / %s = %t, want %t", tc.left, tc.right, same, tc.same)
		}
	}
}
