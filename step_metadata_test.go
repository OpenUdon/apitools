package apitools

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestStepMetadataV1WireFixtures(t *testing.T) {
	requestBytes, err := os.ReadFile("testdata/operation-candidates/v1/request.json")
	if err != nil {
		t.Fatal(err)
	}
	var request OperationCandidateRequest
	if err := json.Unmarshal(requestBytes, &request); err != nil {
		t.Fatal(err)
	}
	contract := request.Contract
	if request.SchemaVersion != OperationCandidatesSchemaVersion ||
		contract.Purpose != "send a message" ||
		contract.Effect != OperationEffectWrite ||
		contract.Inputs["message"].Type != "string" || contract.Inputs["message"].Required == nil || !*contract.Inputs["message"].Required ||
		contract.Outputs["message_id"].Type != "string" || contract.Outputs["message_id"].Required == nil || !*contract.Outputs["message_id"].Required {
		t.Fatalf("decoded step contract request = %#v", request)
	}
	encoded, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	var roundTrip OperationCandidateRequest
	if err := json.Unmarshal(encoded, &roundTrip); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(request, roundTrip) {
		t.Fatalf("round trip = %#v, want %#v", roundTrip, request)
	}

	reportBytes, err := os.ReadFile("testdata/operation-candidates/v1/empty-report.json")
	if err != nil {
		t.Fatal(err)
	}
	var report OperationCandidateReport
	if err := json.Unmarshal(reportBytes, &report); err != nil {
		t.Fatal(err)
	}
	if report.SchemaVersion != OperationCandidatesSchemaVersion || len(report.Candidates) != 0 || report.Truncated {
		t.Fatalf("decoded candidate report = %#v", report)
	}

	candidateBytes, err := os.ReadFile("testdata/operation-candidates/v1/candidate.json")
	if err != nil {
		t.Fatal(err)
	}
	var candidate OperationCandidate
	if err := json.Unmarshal(candidateBytes, &candidate); err != nil {
		t.Fatal(err)
	}
	if candidate.Source.Kind != OperationSourceOpenAPI ||
		candidate.Source.SHA256 != "7c216f3755d84f6ad9adc7a3ca6bd80aa8c89c636e4e0ae642db0f2098f2f459" ||
		candidate.Source.Selector != "#/paths/~1messages/post" ||
		candidate.Summary.Description == "" ||
		candidate.Effect.Class != OperationEffectWrite ||
		candidate.Match.Inputs.Status != ContractMatchCompatible ||
		candidate.Match.Outputs.Status != ContractMatchCompatible {
		t.Fatalf("decoded operation candidate = %#v", candidate)
	}
	if candidate.Match.Score != candidate.Match.Purpose.Score+candidate.Match.Inputs.Score+candidate.Match.Outputs.Score+candidate.Match.Effect.Score || candidate.Match.Effect.Score != 20 {
		t.Fatalf("candidate fixture score does not follow v1 dimension weights: %#v", candidate.Match)
	}
	candidateRoundTrip, err := json.Marshal(candidate)
	if err != nil {
		t.Fatal(err)
	}
	var decodedCandidate OperationCandidate
	if err := json.Unmarshal(candidateRoundTrip, &decodedCandidate); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(candidate, decodedCandidate) {
		t.Fatalf("candidate round trip = %#v, want %#v", decodedCandidate, candidate)
	}
}

func TestOperationCandidateV1FixtureMatchesProducerRanking(t *testing.T) {
	requestBytes, err := os.ReadFile("testdata/operation-candidates/v1/request.json")
	if err != nil {
		t.Fatal(err)
	}
	var request OperationCandidateRequest
	if err := json.Unmarshal(requestBytes, &request); err != nil {
		t.Fatal(err)
	}
	sourceBytes, err := os.ReadFile("testdata/operation-candidates/v1/sources/send-message.yaml")
	if err != nil {
		t.Fatal(err)
	}
	fixtureBytes, err := os.ReadFile("testdata/operation-candidates/v1/candidate.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture OperationCandidate
	if err := json.Unmarshal(fixtureBytes, &fixture); err != nil {
		t.Fatal(err)
	}
	report, err := BuildOperationCandidates(context.Background(), OperationCandidateOptions{
		Sources:  []OperationSourceInput{{Kind: OperationSourceOpenAPI, Path: "openapi/mail.yaml", Content: sourceBytes}},
		Contract: request.Contract,
	})
	if err != nil {
		t.Fatalf("BuildOperationCandidates() error = %v; diagnostics=%#v", err, report.Diagnostics)
	}
	if len(report.Candidates) != 1 || !reflect.DeepEqual(report.Candidates[0].Match, fixture.Match) {
		t.Fatalf("producer ranking differs from v1 fixture:\nproducer: %#v\nfixture: %#v", report.Candidates, fixture)
	}
	if report.Candidates[0].Source.SHA256 != fixture.Source.SHA256 || report.Candidates[0].Source.Selector != fixture.Source.Selector {
		t.Fatalf("producer source identity differs from v1 fixture: producer=%#v fixture=%#v", report.Candidates[0].Source, fixture.Source)
	}
}

func TestEffectAndContractStatusWireEnums(t *testing.T) {
	for _, value := range []OperationEffect{OperationEffectRead, OperationEffectWrite, OperationEffectUnknown} {
		encoded, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		var decoded OperationEffect
		if err := json.Unmarshal(encoded, &decoded); err != nil || decoded != value {
			t.Fatalf("effect %q decoded as %q, error %v", value, decoded, err)
		}
	}
	for _, value := range []ContractMatchStatus{ContractMatchCompatible, ContractMatchIncompatible, ContractMatchIndeterminate} {
		encoded, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		var decoded ContractMatchStatus
		if err := json.Unmarshal(encoded, &decoded); err != nil || decoded != value {
			t.Fatalf("match status %q decoded as %q, error %v", value, decoded, err)
		}
	}
}

func TestOperationSourceReportFixtureAndInputContentPrivacy(t *testing.T) {
	fixture, err := os.ReadFile("testdata/operation-candidates/v1/source-report.json")
	if err != nil {
		t.Fatal(err)
	}
	var report OperationCandidateReport
	if err := json.Unmarshal(fixture, &report); err != nil {
		t.Fatal(err)
	}
	if report.SchemaVersion != OperationCandidatesSchemaVersion || len(report.Sources) != 1 {
		t.Fatalf("source report fixture = %#v", report)
	}
	source := report.Sources[0]
	if source.Kind != OperationSourceGraphQL || source.OperationCount != 0 || len(source.Capabilities) != 1 ||
		source.Capabilities[0].Status != OperationCapabilityUnsupported || source.Capabilities[0].Dimension != "auth" {
		t.Fatalf("decoded source capability report = %#v", source)
	}

	input := OperationSourceInput{Kind: OperationSourceOpenAPI, Path: "openapi/service.json", Content: []byte("private source bytes")}
	encoded, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if _, exists := decoded["content"]; exists {
		t.Fatalf("source bytes leaked into JSON input metadata: %s", encoded)
	}
}

func TestContractRequirednessPreservesUnknownAndOptional(t *testing.T) {
	optional := false
	values := map[string]ContractValue{
		"unknown":  {},
		"optional": {Required: &optional},
	}
	encoded, err := json.Marshal(values)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]ContractValue
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["unknown"].Required != nil || decoded["optional"].Required == nil || *decoded["optional"].Required {
		t.Fatalf("requiredness tri-state lost in %s: %#v", encoded, decoded)
	}
}
