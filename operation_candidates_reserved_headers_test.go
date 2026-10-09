package apitools

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestOperationCandidatesReservedHeaders(t *testing.T) {
	fixture := completeShapeFixture()
	operation := mapValue(mapValue(mapValue(fixture["paths"])["/pets/{id}"])["get"])
	operation["summary"] = "Get pet"
	operation["parameters"] = append(operation["parameters"].([]any),
		map[string]any{"name": "Accept", "in": "header", "required": true, "schema": map[string]any{"type": "string"}},
		map[string]any{"name": "Content-Type", "in": "header", "required": true, "schema": map[string]any{"type": "string"}},
		map[string]any{"name": "Authorization", "in": "header", "required": true, "schema": map[string]any{"type": "string"}})
	source, err := json.Marshal(fixture)
	if err != nil {
		t.Fatal(err)
	}
	inventory, err := BuildOperationInventory(context.Background(), InventoryOptions{Documents: []InventoryDocument{{Content: source}}})
	if err != nil || len(inventory.Operations) != 1 || len(inventory.Operations[0].Parameters) != 1 {
		t.Fatalf("effective inventory: %+v %v", inventory, err)
	}
	report, err := BuildOperationCandidates(context.Background(), OperationCandidateOptions{
		Sources:  []OperationSourceInput{{Kind: OperationSourceOpenAPI, Content: source}},
		Contract: StepContract{Purpose: "Get pet", Inputs: map[string]ContractValue{"id": {Type: "string", Required: boolPointer(true)}}, Effect: OperationEffectRead},
	})
	if err != nil || len(report.Candidates) != 1 {
		t.Fatalf("candidate report: %+v %v", report, err)
	}
	candidate := report.Candidates[0]
	if len(candidate.Summary.Inputs) != 1 || len(candidate.Summary.Gaps) != 0 || candidate.Match.Inputs.Status != ContractMatchCompatible {
		t.Fatalf("reserved declaration constrained candidate: %+v", candidate)
	}
}

func TestOperationCandidatesReservedHeaderMatrix(t *testing.T) {
	for _, version := range []string{"3.0.0", "3.0.3", "3.0.4", "3.1.0", "3.1.1", "3.1.2"} {
		for _, name := range []string{"Accept", "accept", "aCcEpT", "Content-Type", "content-type", "CONTENT-TYPE", "Authorization", "authorization", "AUTHORIZATION"} {
			for _, mode := range []string{"direct", "ref", "chain"} {
				t.Run(version+"/"+name+"/"+mode, func(t *testing.T) {
					fixture := completeShapeFixture()
					fixture["openapi"] = version
					path := mapValue(mapValue(fixture["paths"])["/pets/{id}"])
					operation := mapValue(path["get"])
					operation["summary"] = "Get pet"
					// Invalid/unsupported metadata on an ignored declaration cannot create
					// a schema, requiredness, serialization, or credential-review gap.
					parameter := map[string]any{"name": name, "in": "header", "required": "ignored", "schema": map[string]any{"$ref": "https://fixture.invalid/never-fetched"}, "style": "unsupported", "explode": true}
					var raw any = parameter
					if mode != "direct" {
						fixture["components"] = map[string]any{"parameters": map[string]any{"reserved/~header": parameter, "chain": map[string]any{"$ref": "#/components/parameters/reserved~1~0header"}}}
						target := "reserved~1~0header"
						if mode == "chain" {
							target = "chain"
						}
						raw = map[string]any{"$ref": "#/components/parameters/" + target}
					}
					path["parameters"] = []any{raw}
					operation["parameters"] = append(operation["parameters"].([]any), raw)
					source, err := json.Marshal(fixture)
					if err != nil {
						t.Fatal(err)
					}
					original := append([]byte(nil), source...)
					inventory, err := BuildOperationInventory(context.Background(), InventoryOptions{Documents: []InventoryDocument{{Content: source}}})
					if err != nil || len(inventory.Operations) != 1 || len(inventory.Operations[0].Parameters) != 1 || len(inventory.ReadinessIssues) != 0 {
						t.Fatalf("effective inventory: %+v %v", inventory, err)
					}
					summary, _, err := summarizeOperationForConsumer(inventory.Operations[0], "#/paths/~1pets~1{id}/get", PromptBudget{})
					if err != nil || len(summary.Inputs) != 1 || len(summary.Gaps) != 0 {
						t.Fatalf("consumer summary: %+v %v", summary, err)
					}
					contract := StepContract{Purpose: "Get pet", Inputs: map[string]ContractValue{"id": {Type: "string", Required: boolPointer(true)}}, Effect: OperationEffectRead}
					options := OperationCandidateOptions{Sources: []OperationSourceInput{{Kind: OperationSourceOpenAPI, Content: source, URL: "https://fixture.invalid/provenance"}}, Contract: contract}
					report, err := BuildOperationCandidates(context.Background(), options)
					if err != nil || len(report.Candidates) != 1 || report.Truncated {
						t.Fatalf("candidates: %+v %v", report, err)
					}
					candidate := report.Candidates[0]
					digest := sha256.Sum256(source)
					if len(candidate.Operation.Parameters) != 1 || len(candidate.Summary.Inputs) != 1 || len(candidate.Summary.Gaps) != 0 || candidate.Match.Inputs.Status != ContractMatchCompatible || candidate.Source.SHA256 != hex.EncodeToString(digest[:]) || candidate.Source.Selector != "#/paths/~1pets~1{id}/get" {
						t.Fatalf("candidate constraint/identity drift: %+v", candidate)
					}
					table, err := BuildOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture))
					if err != nil || !table.Operations[0].Complete || len(table.Operations[0].Inputs) != 1 {
						t.Fatalf("shape coherence: %+v %v", table, err)
					}
					if err := VerifyOperationShapeTable(context.Background(), optionsForShapeFixture(t, fixture), table); err != nil {
						t.Fatal(err)
					}
					again, err := BuildOperationCandidates(context.Background(), options)
					if err != nil {
						t.Fatal(err)
					}
					left, _ := json.Marshal(report)
					right, _ := json.Marshal(again)
					if !bytes.Equal(left, right) || !bytes.Equal(source, original) {
						t.Fatal("nondeterministic metadata or mutated raw source")
					}
				})
			}
		}
	}
}

func TestOperationCandidatesReservedHeadersKeepGenuineInputs(t *testing.T) {
	for _, location := range []string{"header", "query", "cookie", "path", "body"} {
		t.Run(location, func(t *testing.T) {
			fixture := completeShapeFixture()
			path := mapValue(mapValue(fixture["paths"])["/pets/{id}"])
			operation := mapValue(path["get"])
			operation["summary"] = "Get pet"
			name := "X-Required"
			if location == "query" || location == "cookie" {
				name = "Accept"
			}
			if location == "path" {
				name = "id"
			}
			operation["parameters"] = []any{map[string]any{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "string"}}}
			if location == "body" {
				name = "body"
				operation["requestBody"] = map[string]any{"required": true, "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "required": []any{"value"}, "properties": map[string]any{"value": map[string]any{"type": "string"}}}}}}
			} else if location != "path" {
				operation["parameters"] = append(operation["parameters"].([]any), map[string]any{"name": name, "in": location, "required": true, "schema": map[string]any{"type": "string"}})
			}
			for _, reserved := range []string{"Accept", "Content-Type", "Authorization"} {
				operation["parameters"] = append(operation["parameters"].([]any), map[string]any{"name": reserved, "in": "header", "required": true, "schema": map[string]any{"type": "string"}})
			}
			source, _ := json.Marshal(fixture)
			inputs := map[string]ContractValue{"id": {Type: "string", Required: boolPointer(true)}, name: {Type: "string", Required: boolPointer(true)}}
			if location == "body" {
				inputs[name] = ContractValue{Type: "object", Required: boolPointer(true), Properties: map[string]ContractValue{"value": {Type: "string", Required: boolPointer(true)}}}
			}
			options := OperationCandidateOptions{Sources: []OperationSourceInput{{Kind: OperationSourceOpenAPI, Content: source}}, Contract: StepContract{Purpose: "Get pet", Inputs: inputs, Effect: OperationEffectRead}}
			good, err := BuildOperationCandidates(context.Background(), options)
			if err != nil || len(good.Candidates) != 1 || good.Candidates[0].Match.Inputs.Status != ContractMatchCompatible {
				t.Fatalf("bound genuine input: %+v %v", good, err)
			}
			inputs[name] = ContractValue{Type: "boolean", Required: boolPointer(true)}
			wrong, err := BuildOperationCandidates(context.Background(), options)
			if err != nil || len(wrong.Candidates) != 1 || wrong.Candidates[0].Match.Inputs.Status == ContractMatchCompatible {
				t.Fatalf("wrong genuine input type accepted: %+v %v", wrong, err)
			}
			delete(inputs, name)
			bad, err := BuildOperationCandidates(context.Background(), options)
			if err != nil || len(bad.Candidates) != 1 || bad.Candidates[0].Match.Inputs.Status == ContractMatchCompatible {
				t.Fatalf("missing genuine input accepted: %+v %v", bad, err)
			}
		})
	}
}

func TestOperationCandidatesReservedHeadersSecurityAndConservativeControls(t *testing.T) {
	fixture := completeShapeFixture()
	operation := mapValue(mapValue(mapValue(fixture["paths"])["/pets/{id}"])["get"])
	operation["summary"] = "Get pet"
	fixture["components"] = map[string]any{"securitySchemes": map[string]any{"key": map[string]any{"type": "apiKey", "in": "header", "name": "Authorization"}, "bearer": map[string]any{"type": "http", "scheme": "bearer"}}}
	fixture["security"] = []any{map[string]any{"key": []any{}, "bearer": []any{}}, map[string]any{}}
	operation["parameters"] = append(operation["parameters"].([]any), map[string]any{"name": "Authorization", "in": "header", "required": true, "schema": map[string]any{"type": "string"}})
	source, _ := json.Marshal(fixture)
	options := OperationCandidateOptions{Sources: []OperationSourceInput{{Kind: OperationSourceOpenAPI, Content: source}}, Contract: StepContract{Purpose: "Get pet", Inputs: map[string]ContractValue{"id": {Type: "string", Required: boolPointer(true)}}, Effect: OperationEffectRead}}
	report, err := BuildOperationCandidates(context.Background(), options)
	if err != nil || len(report.Candidates) != 1 {
		t.Fatal(err)
	}
	candidate := report.Candidates[0]
	if len(candidate.Operation.SecurityRequirementSets) != 2 || len(candidate.Operation.SecurityRequirementSets[0].Requirements) != 2 || len(candidate.Operation.SecurityRequirementSets[1].Requirements) != 0 || candidateCapability(candidate, "auth").Status != OperationCapabilitySupported || len(candidate.Summary.Gaps) != 0 {
		t.Fatalf("independent security changed: %+v", candidate)
	}
	delete(mapValue(fixture["components"]), "securitySchemes")
	source, _ = json.Marshal(fixture)
	options.Sources[0].Content = source
	report, err = BuildOperationCandidates(context.Background(), options)
	if err != nil || len(report.Candidates) != 1 || candidateCapability(report.Candidates[0], "auth").Status == OperationCapabilitySupported {
		t.Fatalf("unknown security upgraded: %+v %v", report, err)
	}
	for _, kind := range []string{"external", "missing", "cycle"} {
		t.Run(kind, func(t *testing.T) {
			fixture := completeShapeFixture()
			operation := mapValue(mapValue(mapValue(fixture["paths"])["/pets/{id}"])["get"])
			ref := "https://fixture.invalid/never-fetched"
			if kind == "missing" {
				ref = "#/components/parameters/missing"
			}
			if kind == "cycle" {
				ref = "#/components/parameters/cycle"
				fixture["components"] = map[string]any{"parameters": map[string]any{"cycle": map[string]any{"$ref": ref}}}
			}
			operation["parameters"] = append(operation["parameters"].([]any), map[string]any{"$ref": ref})
			source, _ := json.Marshal(fixture)
			inventory, err := BuildOperationInventory(context.Background(), InventoryOptions{Documents: []InventoryDocument{{Content: source}}})
			if err != nil || len(inventory.Operations) != 1 || len(inventory.Operations[0].Parameters) != 2 {
				t.Fatalf("unknown declaration silently ignored: %+v %v", inventory, err)
			}
			options.Sources[0].Content = source
			report, err := BuildOperationCandidates(context.Background(), options)
			if err != nil || len(report.Candidates) != 1 || len(report.Candidates[0].Summary.Gaps) == 0 {
				t.Fatalf("unknown declaration lost gaps: %+v %v", report, err)
			}
		})
	}
	swagger := []byte(`{"swagger":"2.0","info":{"title":"Fixture","version":"1"},"paths":{"/pets":{"get":{"operationId":"getPet","parameters":[{"name":"Accept","in":"header","required":true,"type":"string"},{"name":"Content-Type","in":"header","required":true,"type":"string"}],"responses":{"200":{"description":"ok"}}}}}}`)
	inventory, err := BuildOperationInventory(context.Background(), InventoryOptions{Documents: []InventoryDocument{{Content: swagger}}})
	if err != nil || len(inventory.Operations[0].Parameters) != 2 {
		t.Fatalf("Swagger semantics changed: %+v %v", inventory, err)
	}
	options.Sources[0].Content = swagger
	report, err = BuildOperationCandidates(context.Background(), options)
	if err != nil || len(report.Candidates) != 1 || len(report.Candidates[0].Summary.Inputs) != 2 {
		t.Fatalf("Swagger inputs changed: %+v %v", report, err)
	}
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := BuildOperationCandidates(cancelled, options); err == nil {
		t.Fatal("cancelled build accepted")
	}
	options.MaxBytes = 1
	if _, err := BuildOperationCandidates(context.Background(), options); err == nil {
		t.Fatal("source budget ignored")
	}
	if strings.Contains(fmt.Sprint(report), "Credential-shaped parameter Authorization") {
		t.Fatal("OpenAPI reserved header manufactured credential gap")
	}
}

func TestOperationInventoryReservedNamesInGenuineLocations(t *testing.T) {
	for _, location := range []string{"query", "cookie", "path"} {
		for _, name := range []string{"Accept", "Content-Type", "Authorization"} {
			t.Run(location+"/"+name, func(t *testing.T) {
				root := map[string]any{"openapi": "3.0.3"}
				parameters := []any{map[string]any{"name": name, "in": location, "required": true, "schema": map[string]any{"type": "string"}}}
				got := parameterSummaries(root, parameters, &OperationSummary{})
				if len(got) != 1 || got[0].Name != name || got[0].In != location || !got[0].Required {
					t.Fatalf("genuine input lost: %+v", got)
				}
			})
		}
	}
}
