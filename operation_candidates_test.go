package apitools

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildOperationCandidatesAdaptsEverySourceFamily(t *testing.T) {
	kinds := []OperationSourceKind{
		OperationSourceOpenAPI,
		OperationSourceGoogleDiscovery,
		OperationSourceAWSSmithy,
		OperationSourceAsyncAPI,
		OperationSourceGraphQL,
		OperationSourceOpenRPC,
		OperationSourceGRPCProtobuf,
		OperationSourceOData,
	}
	names := []string{"openapi.yaml", "google-discovery.json", "aws-smithy.json", "asyncapi.yaml", "graphql.graphql", "openrpc.json", "grpc.proto", "odata.xml"}
	sources := make([]OperationSourceInput, 0, len(kinds))
	for i, kind := range kinds {
		path := filepath.Join("testdata", "operation-candidates", "v1", "sources", names[i])
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		sources = append(sources, OperationSourceInput{Kind: kind, Path: path, Content: content})
	}
	report, err := BuildOperationCandidates(context.Background(), OperationCandidateOptions{
		Sources:  sources,
		Contract: StepContract{Purpose: "get a pet", Effect: OperationEffectRead},
	})
	if err != nil {
		t.Fatalf("BuildOperationCandidates() error = %v; diagnostics=%#v", err, report.Diagnostics)
	}
	if report.SchemaVersion != OperationCandidatesSchemaVersion || report.Truncated {
		t.Fatalf("report metadata = %#v", report)
	}
	if len(report.Sources) != len(kinds) {
		t.Fatalf("source report count = %d, want %d; reports=%#v", len(report.Sources), len(kinds), report.Sources)
	}
	reportedKinds := map[OperationSourceKind]bool{}
	for _, source := range report.Sources {
		reportedKinds[source.Kind] = true
		if source.OperationCount == 0 || len(source.SHA256) != 64 {
			t.Errorf("source report = %#v", source)
		}
	}
	for _, kind := range kinds {
		if !reportedKinds[kind] {
			t.Errorf("missing source report for %q", kind)
		}
	}
	if len(report.Candidates) < len(kinds) {
		t.Fatalf("candidate count = %d, want at least one per source family", len(report.Candidates))
	}
	seen := map[OperationSourceKind]bool{}
	for _, candidate := range report.Candidates {
		seen[candidate.Source.Kind] = true
		if candidate.Source.SHA256 == "" || !strings.HasPrefix(candidate.Source.Selector, "#/") {
			t.Errorf("candidate lost source identity: %#v", candidate.Source)
		}
		if candidate.Operation.ID == "" || candidate.Effect.Class == "" {
			t.Errorf("candidate missing native identity/effect assessment: %#v", candidate)
		}
	}
	for _, kind := range kinds {
		if !seen[kind] {
			t.Errorf("no candidate from source kind %q", kind)
		}
	}
	assertCandidateEffect(t, report.Candidates, OperationSourceAWSSmithy, "GetPet", OperationEffectRead)
	assertCandidateEffect(t, report.Candidates, OperationSourceGraphQL, "pet", OperationEffectRead)
	assertCandidateEffect(t, report.Candidates, OperationSourceOData, "action.Demo.UpdatePet", OperationEffectWrite)
	assertCandidateEffect(t, report.Candidates, OperationSourceOData, "function.Demo.GetPet", OperationEffectRead)
	assertCandidateEffect(t, report.Candidates, OperationSourceOData, "entityset.Pets", OperationEffectUnknown)
	odataEntitySet := findCandidate(report.Candidates, OperationSourceOData, "entityset.Pets")
	if odataEntitySet == nil || odataEntitySet.Capabilities[3].Status != OperationCapabilityPartial || len(odataEntitySet.Effect.Evidence) == 0 || odataEntitySet.Effect.Evidence[0].Kind != "odata.resource_kind" || !strings.Contains(strings.Join(odataEntitySet.Capabilities[3].Gaps, " "), "not a read-only HTTP operation") {
		t.Fatalf("OData resource effect limitation was not surfaced: %#v", odataEntitySet)
	}
	odataAction := findCandidate(report.Candidates, OperationSourceOData, "action.Demo.UpdatePet")
	if odataAction == nil || len(odataAction.Summary.Inputs) == 0 || odataAction.Summary.Inputs[0].Required != nil || odataAction.Capabilities[0].Status != OperationCapabilityPartial || !strings.Contains(strings.Join(odataAction.Summary.Gaps, " "), "does not establish whether a call argument is required") {
		t.Fatalf("OData parameter requiredness was not kept unknown with its limitation: %#v", odataAction)
	}
	if candidate := findCandidate(report.Candidates, OperationSourceGoogleDiscovery, "pets.pets.get"); candidate == nil || len(candidate.Operation.SecurityRequirementSets) == 0 || candidate.Capabilities[2].Status != OperationCapabilityPartial {
		t.Fatalf("Google Discovery scope/auth evidence was not preserved: %#v", candidate)
	}
	for _, source := range report.Sources {
		if source.Kind != OperationSourceGoogleDiscovery {
			continue
		}
		for _, capability := range source.Capabilities {
			for _, gap := range capability.Gaps {
				if strings.Contains(gap, "Smithy") {
					t.Fatalf("Google Discovery source capability contains an unrelated Smithy limitation: %#v", capability)
				}
			}
		}
	}
	if candidate := findCandidate(report.Candidates, OperationSourceOpenAPI, "createPet"); candidate == nil || len(candidate.Operation.SecurityRequirementSets) != 2 || candidate.Operation.SecurityRequirementSets[0].Requirements[0].Name != "bearerAuth" || candidate.Operation.SecurityRequirementSets[1].Requirements[0].Name != "apiKeyAuth" {
		t.Fatalf("OpenAPI security alternative was not preserved: %#v", candidate)
	}
	assertCandidateValues(t, report.Candidates, OperationSourceOpenAPI, "createPet", []string{"name"}, []string{"id"})
	assertCandidateValues(t, report.Candidates, OperationSourceGoogleDiscovery, "pets.pets.get", []string{"petId"}, []string{"id", "name"})
	assertCandidateValues(t, report.Candidates, OperationSourceAWSSmithy, "GetPet", []string{"petId"}, []string{"name"})
	assertCandidateValues(t, report.Candidates, OperationSourceAsyncAPI, "publishPet", []string{"petId", "name"}, nil)
	assertCandidateValues(t, report.Candidates, OperationSourceGraphQL, "pet", []string{"id"}, []string{"pet"})
	assertCandidateValues(t, report.Candidates, OperationSourceOpenRPC, "pets.get", []string{"petId"}, []string{"id", "name"})
	assertCandidateValues(t, report.Candidates, OperationSourceGRPCProtobuf, "GetPet", []string{"petId"}, []string{"id", "name"})
	assertCandidateValues(t, report.Candidates, OperationSourceOData, "entityset.Pets", nil, []string{"ID", "Name"})
}

func TestBuildOperationCandidatesMarksUnresolvedOpenAPISecurityPartial(t *testing.T) {
	content := []byte(`openapi: 3.0.3
info: {title: Missing security scheme, version: 1}
paths:
  /items:
    get:
      operationId: getItems
      summary: Get items
      security:
        - missingBearer: []
      responses: {"200": {description: ok}}
`)
	report, err := BuildOperationCandidates(context.Background(), OperationCandidateOptions{
		Sources:  []OperationSourceInput{{Kind: OperationSourceOpenAPI, Content: content}},
		Contract: StepContract{Purpose: "get items"},
	})
	if err != nil {
		t.Fatalf("BuildOperationCandidates() error = %v; diagnostics=%#v", err, report.Diagnostics)
	}
	if len(report.Candidates) != 1 {
		t.Fatalf("candidate count = %d; diagnostics=%#v", len(report.Candidates), report.Diagnostics)
	}
	candidate := report.Candidates[0]
	if len(candidate.Operation.SecurityRequirementSets) != 1 || candidate.Operation.SecurityRequirementSets[0].Requirements[0].Name != "missingBearer" {
		t.Fatalf("unresolved OpenAPI security alternative was not retained: %#v", candidate.Operation.SecurityRequirementSets)
	}
	if candidate.Capabilities[2].Status != OperationCapabilityPartial || !strings.Contains(strings.Join(candidate.Capabilities[2].Gaps, " "), "incomplete security scheme metadata") {
		t.Fatalf("unresolved OpenAPI security alternative was reported as complete: %#v", candidate.Capabilities[2])
	}
}

func TestBuildOperationCandidatesOmitsCredentialShapedSchemaFields(t *testing.T) {
	content := []byte(`openapi: 3.0.3
info: {title: Credentials, version: 1}
paths:
  /tokens:
    post:
      operationId: createToken
      summary: Create token
      parameters:
        - {name: accessToken, in: query, required: true, type: string}
        - {name: clientSecret, in: query, required: true, type: string}
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [api_key, accessToken, clientSecret, name]
              properties:
                api_key: {type: string}
                accessToken: {type: string}
                clientSecret: {type: string}
                name: {type: string}
                profile:
                  type: object
                  properties:
                    refreshToken: {type: string}
      responses: {"200": {description: ok}}
`)
	report, err := BuildOperationCandidates(context.Background(), OperationCandidateOptions{
		Sources:  []OperationSourceInput{{Kind: OperationSourceOpenAPI, Content: content}},
		Contract: StepContract{Purpose: "create token"},
	})
	if err != nil {
		t.Fatalf("BuildOperationCandidates() error = %v; diagnostics=%#v", err, report.Diagnostics)
	}
	if len(report.Candidates) != 1 {
		t.Fatalf("candidate count = %d", len(report.Candidates))
	}
	candidate := report.Candidates[0]
	encoded, err := json.Marshal(candidate)
	if err != nil {
		t.Fatal(err)
	}
	for _, credentialName := range []string{"api_key", "accessToken", "clientSecret", "refreshToken"} {
		if strings.Contains(strings.ToLower(string(encoded)), strings.ToLower(credentialName)) {
			t.Fatalf("credential-shaped field %q leaked into candidate: %s", credentialName, encoded)
		}
	}
	if !strings.Contains(strings.Join(candidate.Summary.Gaps, " "), "Credential-shaped") {
		t.Fatalf("credential-shaped field leaked or was not explained: %s; gaps=%#v", encoded, candidate.Summary.Gaps)
	}
}

func TestBuildOperationCandidatesRedactsSensitiveURLProvenance(t *testing.T) {
	content := []byte(`openapi: 3.0.3
info: {title: Safe provenance, version: 1}
paths: {/items: {get: {operationId: getItems, summary: Get items, responses: {"200": {description: ok}}}}}
`)
	report, err := BuildOperationCandidates(context.Background(), OperationCandidateOptions{
		Sources: []OperationSourceInput{{
			Kind:    OperationSourceOpenAPI,
			URL:     "https://user:password@example.test/openapi.yaml?api_key=secret#accessToken=fragment-secret",
			Content: content,
		}},
		Contract: StepContract{Purpose: "get items"},
	})
	if err != nil {
		t.Fatalf("BuildOperationCandidates() error = %v; diagnostics=%#v", err, report.Diagnostics)
	}
	if len(report.Sources) != 1 || report.Sources[0].DocumentURL != "https://example.test/openapi.yaml" || !hasDiagnosticCode(report.Diagnostics, "source.url_redacted") {
		t.Fatalf("sensitive URL provenance was not safely redacted: %#v", report)
	}
	encoded, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	for _, secret := range []string{"user:password", "api_key=secret", "fragment-secret"} {
		if strings.Contains(string(encoded), secret) {
			t.Fatalf("URL credential %q leaked into report: %s", secret, encoded)
		}
	}
}

func TestBuildOperationCandidatesReportsAsyncAPIPayloadAlternatives(t *testing.T) {
	content := []byte(`asyncapi: 3.0.0
info: {title: Multi-message event, version: 1}
operations:
  publishEvent:
    action: send
    channel: {$ref: '#/channels/events'}
    messages:
      - {$ref: '#/channels/events/messages/created'}
      - {$ref: '#/channels/events/messages/updated'}
channels:
  events:
    address: events
    messages:
      created:
        payload: {type: object, properties: {createdId: {type: string}}}
      updated:
        payload: {type: object, properties: {updatedId: {type: string}}}
`)
	report, err := BuildOperationCandidates(context.Background(), OperationCandidateOptions{
		Sources:  []OperationSourceInput{{Kind: OperationSourceAsyncAPI, Content: content}},
		Contract: StepContract{Purpose: "publish event"},
	})
	if err != nil {
		t.Fatalf("BuildOperationCandidates() error = %v; diagnostics=%#v", err, report.Diagnostics)
	}
	if len(report.Candidates) != 1 || !strings.Contains(strings.Join(report.Candidates[0].Summary.Gaps, " "), "multiple message payload alternatives") {
		t.Fatalf("AsyncAPI payload alternatives were silently collapsed: %#v", report.Candidates)
	}
}

func TestBuildOperationCandidatesReportsOpenAPIResponseBranches(t *testing.T) {
	content := []byte(`openapi: 3.0.3
info: {title: Content variants, version: 1}
paths:
  /items:
    post:
      operationId: createItem
      summary: Create item
      requestBody:
        required: true
        content:
          application/json:
            schema: {type: object, properties: {jsonName: {type: string}}}
          text/plain:
            schema: {type: string}
      responses:
        '200':
          description: JSON result
          content:
            application/json:
              schema: {type: object, required: [jsonId], properties: {jsonId: {type: string}}}
            application/xml:
              schema: {type: string}
        '201':
          description: alternate result
          content:
            application/json:
              schema: {type: object, properties: {alternateId: {type: string}}}
`)
	report, err := BuildOperationCandidates(context.Background(), OperationCandidateOptions{
		Sources:  []OperationSourceInput{{Kind: OperationSourceOpenAPI, Content: content}},
		Contract: StepContract{Purpose: "create item"},
	})
	if err != nil {
		t.Fatalf("BuildOperationCandidates() error = %v; diagnostics=%#v", err, report.Diagnostics)
	}
	if len(report.Candidates) != 1 {
		t.Fatalf("candidate count = %d", len(report.Candidates))
	}
	candidate := report.Candidates[0]
	if candidate.Capabilities[0].Status != OperationCapabilityPartial || candidate.Capabilities[1].Status != OperationCapabilityPartial {
		t.Fatalf("multiple source branches were marked complete: %#v", candidate.Capabilities)
	}
	if !containsOperationValueWithRequired(candidate.Summary.Outputs, "jsonId", true) {
		t.Fatalf("schema-declared required output was summarized as optional: %#v", candidate.Summary.Outputs)
	}
	gaps := strings.Join(candidate.Summary.Gaps, " ")
	for _, phrase := range []string{"multiple media types", "multiple successful response statuses"} {
		if !strings.Contains(gaps, phrase) {
			t.Errorf("candidate gaps do not report %q: %#v", phrase, candidate.Summary.Gaps)
		}
	}
}

func TestBuildOperationCandidatesKeepsScalarResponseBodyRequirednessUnknown(t *testing.T) {
	content := []byte(`openapi: 3.0.3
info: {title: Scalar response, version: 1}
paths:
  /body:
    get:
      operationId: getBody
      summary: Get body
      responses:
        '200':
          description: Body value
          content:
            text/plain:
              schema: {type: string}
`)
	required := true
	report, err := BuildOperationCandidates(context.Background(), OperationCandidateOptions{
		Sources:  []OperationSourceInput{{Kind: OperationSourceOpenAPI, Content: content}},
		Contract: StepContract{Purpose: "get body", Outputs: map[string]ContractValue{"body": {Type: "string", Required: &required}}},
	})
	if err != nil {
		t.Fatalf("BuildOperationCandidates() error = %v; diagnostics=%#v", err, report.Diagnostics)
	}
	if len(report.Candidates) != 1 {
		t.Fatalf("candidate count = %d", len(report.Candidates))
	}
	candidate := report.Candidates[0]
	if len(candidate.Summary.Outputs) != 1 || candidate.Summary.Outputs[0].Required != nil {
		t.Fatalf("synthetic response-body output claims requiredness: %#v", candidate.Summary.Outputs)
	}
	if candidate.Match.Outputs.Status != ContractMatchIndeterminate || len(candidate.Match.Outputs.Conflicts) != 0 {
		t.Fatalf("unknown response-body presence was treated as a conflict: %#v", candidate.Match.Outputs)
	}
}

func TestBuildOperationCandidatesBoundsContextErrorReport(t *testing.T) {
	var paths strings.Builder
	for i := range 12 {
		paths.WriteString("  /items/")
		paths.WriteString(string(rune('a' + i)))
		paths.WriteString(": {get: {operationId: getItem")
		paths.WriteString(string(rune('a' + i)))
		paths.WriteString(", summary: Get item, responses: {\"200\": {description: ok}}}}\n")
	}
	content := []byte("openapi: 3.0.3\ninfo: {title: Many operations, version: 1}\npaths:\n" + paths.String())
	report, err := BuildOperationCandidates(context.Background(), OperationCandidateOptions{
		Sources:      []OperationSourceInput{{Kind: OperationSourceOpenAPI, Content: content}},
		Contract:     StepContract{Purpose: "get item"},
		PromptBudget: PromptBudget{MaxContextBytes: 512},
	})
	if err == nil || !report.Truncated || len(report.Candidates) != 0 || len(report.Sources) != 0 || !hasDiagnosticCode(report.Diagnostics, "candidate.context_byte_budget") {
		t.Fatalf("oversized report content was retained after context failure: report=%#v err=%v", report, err)
	}
	encoded, marshalErr := json.Marshal(report)
	if marshalErr != nil {
		t.Fatal(marshalErr)
	}
	if len(encoded) > 512 {
		t.Fatalf("bounded error report is %d bytes, exceeds context budget: %s", len(encoded), encoded)
	}
}

func TestSanitizeAdapterCandidatesContinuesAfterDiscard(t *testing.T) {
	candidates := []OperationCandidate{
		{Source: OperationSourceIdentity{Selector: "#/unsafe"}, Summary: ConsumerOperationSummary{Gaps: []string{strings.Repeat("x", DefaultPromptTextRunes+1)}}},
		{Source: OperationSourceIdentity{Selector: "#/safe"}, Operation: OperationSummary{ID: "safe"}, Summary: ConsumerOperationSummary{Gaps: []string{"safe gap"}}},
	}
	sanitized, diagnostics, err := sanitizeAdapterCandidates(candidates, resolvedPromptBudget(PromptBudget{}))
	if err == nil || len(diagnostics) != 1 || diagnostics[0].Code != "candidate.rationale_unsafe" {
		t.Fatalf("unsafe rationale was not rejected: candidates=%#v diagnostics=%#v err=%v", sanitized, diagnostics, err)
	}
	if len(sanitized) != 1 || sanitized[0].Operation.ID != "safe" {
		t.Fatalf("sanitization did not safely continue to the next candidate: %#v", sanitized)
	}
}

func TestNativeSchemaFieldsReportsUnsupportedComposition(t *testing.T) {
	_, gaps, err := nativeSchemaFields(map[string]any{
		"oneOf": []any{
			map[string]any{"type": "object", "properties": map[string]any{"first": map[string]any{"type": "string"}}},
			map[string]any{"type": "object", "properties": map[string]any{"second": map[string]any{"type": "string"}}},
		},
		"additionalProperties": map[string]any{"type": "string"},
	}, DefaultPromptBudget(), nil)
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(gaps, " ")
	for _, phrase := range []string{"oneOf composition is not expanded", "additionalProperties dictionary values are not summarized"} {
		if !strings.Contains(joined, phrase) {
			t.Fatalf("unsupported JSON Schema shape %q was not surfaced: %#v", phrase, gaps)
		}
	}
}

func containsOperationValueWithRequired(values []OperationValueSummary, name string, required bool) bool {
	for _, value := range values {
		if strings.EqualFold(value.Name, name) && value.Required != nil && *value.Required == required {
			return true
		}
	}
	return false
}

func TestBuildOperationCandidatesChecksNativeOperationLimitBeforeAdaptation(t *testing.T) {
	content := []byte(`type Query { first: String second: String }
`)
	report, err := BuildOperationCandidates(context.Background(), OperationCandidateOptions{
		Sources:  []OperationSourceInput{{Kind: OperationSourceGraphQL, Content: content}},
		Contract: StepContract{Purpose: "query data"}, MaxOperations: 1,
	})
	if err == nil || !report.Truncated || len(report.Candidates) != 0 || !hasDiagnosticCode(report.Diagnostics, "candidate.operation_work_budget") {
		t.Fatalf("native operation budget was not enforced before ranking: report=%#v err=%v", report, err)
	}
	if len(report.Sources) != 1 || report.Sources[0].OperationCount != 2 {
		t.Fatalf("source count diagnostic lost parsed operation count: %#v", report.Sources)
	}
}

func TestNativeOperationAdapterHonorsMidwayCancellation(t *testing.T) {
	ctx := &cancelAfterErrContext{Context: context.Background(), cancelAt: 2}
	result, err := adaptGraphQLSource(ctx, loadedOperationSource{
		input:   OperationSourceInput{Kind: OperationSourceGraphQL},
		content: []byte("type Query { first: String second: String }"),
	}, DefaultPromptBudget())
	if !errors.Is(err, context.Canceled) || ctx.checks != 2 || len(result.candidates) != 1 {
		t.Fatalf("native adapter did not stop between operations: candidates=%d checks=%d err=%v", len(result.candidates), ctx.checks, err)
	}
}

func TestBuildOperationCandidatesUsesExactDigestAndNativeSelector(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("testdata", "operation-candidates", "v1", "sources", "openapi.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	expected := sha256.Sum256(content)
	report, err := BuildOperationCandidates(context.Background(), OperationCandidateOptions{
		Sources:  []OperationSourceInput{{Kind: OperationSourceOpenAPI, Path: "first.yaml", Content: content}},
		Contract: StepContract{Purpose: "create a pet"},
	})
	if err != nil {
		t.Fatalf("BuildOperationCandidates() error = %v", err)
	}
	wantDigest := hex.EncodeToString(expected[:])
	for _, candidate := range report.Candidates {
		if candidate.Source.SHA256 != wantDigest {
			t.Fatalf("source digest = %q, want %q", candidate.Source.SHA256, wantDigest)
		}
	}
	if candidate := findCandidate(report.Candidates, OperationSourceOpenAPI, "createPet"); candidate == nil || candidate.Source.Selector != "#/paths/~1pets/post" {
		t.Fatalf("native OpenAPI selector was not preserved: %#v", candidate)
	}
}

func TestBuildOperationCandidatesRetainsDuplicateIDsAcrossSources(t *testing.T) {
	content := []byte(`openapi: 3.0.3
info: {title: Duplicate, version: 1}
paths:
  /items:
    get:
      operationId: getItem
      summary: Get item
      responses: {"200": {description: ok}}
`)
	report, err := BuildOperationCandidates(context.Background(), OperationCandidateOptions{
		Sources: []OperationSourceInput{
			{Kind: OperationSourceOpenAPI, Path: "a.yaml", Content: content},
			{Kind: OperationSourceOpenAPI, Path: "b.yaml", Content: content},
		},
		Contract: StepContract{Purpose: "get item"},
	})
	if err != nil {
		t.Fatalf("BuildOperationCandidates() error = %v", err)
	}
	if len(report.Candidates) != 2 || report.Candidates[0].Source.DocumentPath == report.Candidates[1].Source.DocumentPath {
		t.Fatalf("duplicate native IDs across source documents collapsed: %#v", report.Candidates)
	}
}

func TestBuildOperationCandidatesDiagnosticAndSourceOrderingIsStable(t *testing.T) {
	first := OperationSourceInput{Kind: OperationSourceOpenAPI, Path: "z-invalid.yaml", Content: []byte("openapi: [")}
	second := OperationSourceInput{Kind: OperationSourceOpenAPI, Path: "a-invalid.yaml", Content: []byte("openapi: [")}
	build := func(sources []OperationSourceInput) []byte {
		t.Helper()
		report, err := BuildOperationCandidates(context.Background(), OperationCandidateOptions{
			Sources: sources, Contract: StepContract{Purpose: "list items"},
		})
		if err == nil || !report.Truncated {
			t.Fatalf("malformed sources should produce a blocking diagnostic; report=%#v err=%v", report, err)
		}
		encoded, marshalErr := json.Marshal(report)
		if marshalErr != nil {
			t.Fatal(marshalErr)
		}
		return encoded
	}
	forward := build([]OperationSourceInput{first, second})
	reverse := build([]OperationSourceInput{second, first})
	if string(forward) != string(reverse) {
		t.Fatalf("source input order changed normalized report:\nforward: %s\nreverse: %s", forward, reverse)
	}
}

func TestBuildOperationCandidatesBlocksReadFailuresAndSemanticTruncation(t *testing.T) {
	valid := []byte(`openapi: 3.0.3
info: {title: Good, version: 1}
paths: {/items: {get: {operationId: listItems, summary: List items, responses: {"200": {description: ok}}}}}
`)
	report, err := BuildOperationCandidates(context.Background(), OperationCandidateOptions{
		Sources: []OperationSourceInput{
			{Kind: OperationSourceOpenAPI, Name: "good", Content: valid},
			{Kind: OperationSourceOpenAPI, Path: filepath.Join(t.TempDir(), "missing.yaml")},
		},
		Contract: StepContract{Purpose: "list items"},
	})
	if err == nil || !report.Truncated || len(report.Candidates) != 1 || !hasDiagnosticCode(report.Diagnostics, "source.read") {
		t.Fatalf("partial local-source failure was not blocking: report=%#v err=%v", report, err)
	}

	fields := make([]string, 61)
	for i := range fields {
		fields[i] = "field" + strings.Repeat("x", 0)
	}
	var properties strings.Builder
	for i := range fields {
		properties.WriteString("                  field")
		properties.WriteString(string(rune('a' + i%26)))
		properties.WriteString("x")
		properties.WriteString(string(rune('A' + i/26)))
		properties.WriteString(": {type: string}\n")
	}
	largeSchema := []byte("openapi: 3.0.3\ninfo: {title: Wide, version: 1}\npaths:\n  /items:\n    post:\n      operationId: createItem\n      summary: Create item\n      requestBody:\n        required: true\n        content:\n          application/json:\n            schema:\n              type: object\n              properties:\n" + properties.String() + "      responses: {\"200\": {description: ok}}\n")
	report, err = BuildOperationCandidates(context.Background(), OperationCandidateOptions{
		Sources:  []OperationSourceInput{{Kind: OperationSourceOpenAPI, Content: largeSchema}},
		Contract: StepContract{Purpose: "create item"},
	})
	if err == nil || !report.Truncated || !hasDiagnosticCode(report.Diagnostics, "candidate.schema_field_truncated") {
		t.Fatalf("semantically truncated OpenAPI fields were not blocked: report=%#v err=%v", report, err)
	}
}

func TestBuildOperationCandidatesRejectsUnsafeLocalInputsAndCancellation(t *testing.T) {
	temp := t.TempDir()
	target := filepath.Join(temp, "source.yaml")
	if err := os.WriteFile(target, []byte("not used"), 0o600); err != nil {
		t.Fatal(err)
	}
	symlink := filepath.Join(temp, "link.yaml")
	if err := os.Symlink(target, symlink); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	report, err := BuildOperationCandidates(context.Background(), OperationCandidateOptions{
		Sources:  []OperationSourceInput{{Kind: OperationSourceOpenAPI, Path: symlink}},
		Contract: StepContract{Purpose: "read data"},
	})
	if err == nil || !report.Truncated || !hasDiagnosticCode(report.Diagnostics, "source.read") {
		t.Fatalf("symlinked source was not blocked: report=%#v err=%v", report, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = BuildOperationCandidates(ctx, OperationCandidateOptions{Sources: []OperationSourceInput{{Kind: OperationSourceOpenAPI, Content: []byte("bad")}}})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation error = %v, want context.Canceled", err)
	}
}

func findCandidate(candidates []OperationCandidate, kind OperationSourceKind, id string) *OperationCandidate {
	for i := range candidates {
		candidate := &candidates[i]
		if candidate.Source.Kind == kind && (strings.EqualFold(candidate.Operation.ID, id) || strings.EqualFold(candidate.Operation.OperationID, id) || strings.Contains(strings.ToLower(candidate.Operation.ID), strings.ToLower(id))) {
			return candidate
		}
	}
	return nil
}

func assertCandidateEffect(t *testing.T, candidates []OperationCandidate, kind OperationSourceKind, id string, effect OperationEffect) {
	t.Helper()
	for _, candidate := range candidates {
		if candidate.Source.Kind == kind && (strings.Contains(strings.ToLower(candidate.Operation.ID), strings.ToLower(id)) || strings.Contains(strings.ToLower(candidate.Operation.OperationID), strings.ToLower(id))) {
			if candidate.Effect.Class != effect {
				t.Fatalf("%s candidate %q effect = %q, want %s", kind, id, candidate.Effect.Class, effect)
			}
			return
		}
	}
	t.Fatalf("%s candidate %q not found", kind, id)
}

func assertCandidateValues(t *testing.T, candidates []OperationCandidate, kind OperationSourceKind, id string, inputs, outputs []string) {
	t.Helper()
	candidate := findCandidate(candidates, kind, id)
	if candidate == nil {
		t.Fatalf("%s candidate %q not found", kind, id)
	}
	for _, name := range inputs {
		if !containsOperationValue(candidate.Summary.Inputs, name) {
			t.Errorf("%s candidate %q is missing input %q: %#v", kind, id, name, candidate.Summary.Inputs)
		}
	}
	for _, name := range outputs {
		if !containsOperationValue(candidate.Summary.Outputs, name) {
			t.Errorf("%s candidate %q is missing output %q: %#v", kind, id, name, candidate.Summary.Outputs)
		}
	}
}

func containsOperationValue(values []OperationValueSummary, name string) bool {
	for _, value := range values {
		if strings.EqualFold(value.Name, name) {
			return true
		}
	}
	return false
}

type cancelAfterErrContext struct {
	context.Context
	cancelAt int
	checks   int
}

func (ctx *cancelAfterErrContext) Err() error {
	ctx.checks++
	if ctx.checks >= ctx.cancelAt {
		return context.Canceled
	}
	return nil
}
