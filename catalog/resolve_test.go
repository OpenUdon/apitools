package catalog

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestResolveProviderUsesBuiltInMetadata(t *testing.T) {
	resolved, err := ResolveProvider(ResolveProviderOptions{ProviderKey: "slack"})
	if err != nil {
		t.Fatalf("ResolveProvider() error = %v", err)
	}
	if resolved.OpenAPI.Source != ResolutionSourceBuiltInSpecReference {
		t.Fatalf("OpenAPI source = %q, want %q", resolved.OpenAPI.Source, ResolutionSourceBuiltInSpecReference)
	}
	if resolved.OpenAPI.SpecRefID != "slack-web-openapi-v2" {
		t.Fatalf("OpenAPI spec ref = %q", resolved.OpenAPI.SpecRefID)
	}
	if resolved.Security.Source != ResolutionSourceBuiltInSecurityOverlay {
		t.Fatalf("Security source = %q, want %q", resolved.Security.Source, ResolutionSourceBuiltInSecurityOverlay)
	}
	if resolved.Security.SourceNote == "" || resolved.Security.SourceNote != resolved.CatalogSecurityOverlays[0].SourceNote {
		t.Fatalf("Security source note = %q, want overlay note %q", resolved.Security.SourceNote, resolved.CatalogSecurityOverlays[0].SourceNote)
	}
	if resolved.SecurityStatus != AuthStatusPresentIncomplete {
		t.Fatalf("SecurityStatus = %q, want %q", resolved.SecurityStatus, AuthStatusPresentIncomplete)
	}
}

func TestResolveProviderReportsSelectedReferenceKindAndProtocol(t *testing.T) {
	tests := []struct {
		provider string
		specID   string
		kind     SpecKind
		protocol SpecProtocol
	}{
		{provider: "slack", specID: "slack-web-openapi-v2", kind: SpecKindOpenAPI, protocol: SpecProtocolSwagger},
		{provider: "aws-acm", specID: "aws-acm-smithy-model", kind: SpecKindSmithyJSON, protocol: SpecProtocolSmithy},
		{provider: "gmail", specID: "gmail-discovery-v1", kind: SpecKindGoogleDiscovery, protocol: SpecProtocolGoogleDiscovery},
		{provider: "airtable", specID: "airtable-web-api-docs", kind: SpecKindHumanDocs, protocol: SpecProtocolHumanDocs},
	}
	for _, test := range tests {
		t.Run(test.provider, func(t *testing.T) {
			resolved, err := ResolveProvider(ResolveProviderOptions{ProviderKey: test.provider})
			if err != nil {
				t.Fatalf("ResolveProvider() error = %v", err)
			}
			if resolved.OpenAPI.SpecRefID != test.specID || resolved.OpenAPI.Kind != test.kind || resolved.OpenAPI.Protocol != test.protocol {
				t.Fatalf("resolved OpenAPI source = %#v, want spec=%q kind=%q protocol=%q", resolved.OpenAPI, test.specID, test.kind, test.protocol)
			}
		})
	}
}

func TestResolvedReferenceJSONAddsKindAndProtocolCompatibly(t *testing.T) {
	const legacy = `{"source":"built-in-spec-reference","value":"https://example.test/api.json","spec_ref_id":"example"}`
	var decoded ResolvedReference
	if err := json.Unmarshal([]byte(legacy), &decoded); err != nil {
		t.Fatalf("unmarshal legacy resolved reference: %v", err)
	}
	if decoded.Source != ResolutionSourceBuiltInSpecReference || decoded.Kind != "" || decoded.Protocol != "" {
		t.Fatalf("legacy reference decoded as %#v", decoded)
	}
	encodedLegacy, err := json.Marshal(decoded)
	if err != nil {
		t.Fatalf("marshal legacy resolved reference: %v", err)
	}
	if strings.Contains(string(encodedLegacy), `"kind"`) || strings.Contains(string(encodedLegacy), `"protocol"`) {
		t.Fatalf("zero-value additive fields were not omitted: %s", encodedLegacy)
	}

	withClassification := ResolvedReference{
		Source:    ResolutionSourceBuiltInSpecReference,
		SpecRefID: "gmail-discovery-v1",
		Kind:      SpecKindGoogleDiscovery,
		Protocol:  SpecProtocolGoogleDiscovery,
	}
	encoded, err := json.Marshal(withClassification)
	if err != nil {
		t.Fatalf("marshal classified resolved reference: %v", err)
	}
	var roundTrip ResolvedReference
	if err := json.Unmarshal(encoded, &roundTrip); err != nil {
		t.Fatalf("unmarshal classified resolved reference: %v", err)
	}
	if roundTrip.Kind != withClassification.Kind || roundTrip.Protocol != withClassification.Protocol {
		t.Fatalf("classified reference round-trip = %#v", roundTrip)
	}
}

func TestResolveProviderDoesNotBorrowSecurityFromAnotherSpec(t *testing.T) {
	resolved, err := ResolveProvider(ResolveProviderOptions{ProviderKey: "zoho"})
	if err != nil {
		t.Fatalf("ResolveProvider() error = %v", err)
	}
	if resolved.OpenAPI.SpecRefID != "zoho-crm-v8-apis-openapi" {
		t.Fatalf("OpenAPI spec ref = %q", resolved.OpenAPI.SpecRefID)
	}
	if resolved.SecurityStatus != AuthStatusUnknown {
		t.Fatalf("SecurityStatus = %q, want %q", resolved.SecurityStatus, AuthStatusUnknown)
	}
	if resolved.Security.Source != ResolutionSourceNone {
		t.Fatalf("Security source = %#v, want no cross-spec evidence", resolved.Security)
	}
	if resolved.SecurityReport.Status != AuthStatusComplete {
		t.Fatalf("provider aggregate status = %q, want retained report evidence", resolved.SecurityReport.Status)
	}
}

func TestResolveProviderPrefersSmithyBeforeHumanDocs(t *testing.T) {
	resolved, err := ResolveProvider(ResolveProviderOptions{ProviderKey: "aws-s3"})
	if err != nil {
		t.Fatalf("ResolveProvider() error = %v", err)
	}
	if resolved.OpenAPI.Source != ResolutionSourceBuiltInSpecReference {
		t.Fatalf("OpenAPI source = %q, want %q", resolved.OpenAPI.Source, ResolutionSourceBuiltInSpecReference)
	}
	if resolved.OpenAPI.SpecRefID != "aws-s3-smithy-model" {
		t.Fatalf("OpenAPI spec ref = %q, want aws-s3-smithy-model", resolved.OpenAPI.SpecRefID)
	}
}

func TestResolveProviderPrefersDropboxStoneBeforeHumanDocs(t *testing.T) {
	resolved, err := ResolveProvider(ResolveProviderOptions{ProviderKey: "dropbox"})
	if err != nil {
		t.Fatalf("ResolveProvider() error = %v", err)
	}
	if resolved.OpenAPI.Source != ResolutionSourceBuiltInSpecReference {
		t.Fatalf("OpenAPI source = %q, want %q", resolved.OpenAPI.Source, ResolutionSourceBuiltInSpecReference)
	}
	if resolved.OpenAPI.SpecRefID != "dropbox-api-stone-spec" {
		t.Fatalf("OpenAPI spec ref = %q, want dropbox-api-stone-spec", resolved.OpenAPI.SpecRefID)
	}
	if resolved.OpenAPI.Value != "https://github.com/dropbox/dropbox-api-spec" {
		t.Fatalf("OpenAPI value = %q, want official Dropbox Stone repository", resolved.OpenAPI.Value)
	}
}

func TestResolveProviderPreservesAuthoredHumanDocsPriority(t *testing.T) {
	tests := []struct {
		provider  string
		specRefID string
	}{
		{provider: "airtable", specRefID: "airtable-web-api-docs"},
		{provider: "linear", specRefID: "linear-graphql-docs"},
		{provider: "quickbooks", specRefID: "quickbooks-online-api-docs"},
		{provider: "salesforce", specRefID: "salesforce-rest-docs"},
		{provider: "servicenow", specRefID: "servicenow-rest-api-docs"},
		{provider: "shopify", specRefID: "shopify-admin-rest-docs"},
	}
	for _, test := range tests {
		resolved, err := ResolveProvider(ResolveProviderOptions{ProviderKey: test.provider})
		if err != nil {
			t.Fatalf("%s: ResolveProvider() error = %v", test.provider, err)
		}
		if resolved.OpenAPI.SpecRefID != test.specRefID {
			t.Fatalf("%s: OpenAPI spec ref = %q, want %q", test.provider, resolved.OpenAPI.SpecRefID, test.specRefID)
		}
	}
}

func TestResolveProviderPrecedence(t *testing.T) {
	tests := []struct {
		name         string
		options      ResolveProviderOptions
		wantOpenAPI  ResolutionSource
		wantSecurity ResolutionSource
		wantStatus   AuthCompletenessStatus
	}{
		{
			name: "user openapi wins over built in spec and overlay",
			options: ResolveProviderOptions{
				ProviderKey: "slack",
				UserOpenAPI: "./openapi/slack.yaml",
			},
			wantOpenAPI:  ResolutionSourceUserOpenAPI,
			wantSecurity: ResolutionSourceNone,
			wantStatus:   AuthStatusUnknown,
		},
		{
			name: "user security overlay wins over built in overlay",
			options: ResolveProviderOptions{
				ProviderKey:         "slack",
				UserSecurityOverlay: "./security/slack-overlay.json",
			},
			wantOpenAPI:  ResolutionSourceBuiltInSpecReference,
			wantSecurity: ResolutionSourceUserSecurityOverlay,
			wantStatus:   AuthStatusUnknown,
		},
		{
			name: "project local openapi wins over built in spec and overlay",
			options: ResolveProviderOptions{
				ProviderKey:         "slack",
				ProjectLocalOpenAPI: "./openapi/slack.yaml",
			},
			wantOpenAPI:  ResolutionSourceProjectLocalOpenAPI,
			wantSecurity: ResolutionSourceNone,
			wantStatus:   AuthStatusUnknown,
		},
		{
			name: "user openapi wins over project local openapi",
			options: ResolveProviderOptions{
				ProviderKey:         "slack",
				UserOpenAPI:         "https://example.com/slack.yaml",
				ProjectLocalOpenAPI: "./openapi/slack.yaml",
			},
			wantOpenAPI:  ResolutionSourceUserOpenAPI,
			wantSecurity: ResolutionSourceNone,
			wantStatus:   AuthStatusUnknown,
		},
	}
	for _, test := range tests {
		resolved, err := ResolveProvider(test.options)
		if err != nil {
			t.Fatalf("%s: ResolveProvider() error = %v", test.name, err)
		}
		if resolved.OpenAPI.Source != test.wantOpenAPI {
			t.Fatalf("%s: OpenAPI source = %q, want %q", test.name, resolved.OpenAPI.Source, test.wantOpenAPI)
		}
		if resolved.Security.Source != test.wantSecurity {
			t.Fatalf("%s: Security source = %q, want %q", test.name, resolved.Security.Source, test.wantSecurity)
		}
		if resolved.SecurityStatus != test.wantStatus {
			t.Fatalf("%s: SecurityStatus = %q, want %q", test.name, resolved.SecurityStatus, test.wantStatus)
		}
	}
}

func TestResolveProviderClassifiesExplicitOpenAPISourcesByDeclaration(t *testing.T) {
	tests := []struct {
		name    string
		options ResolveProviderOptions
		source  ResolutionSource
	}{
		{
			name:    "user openapi",
			options: ResolveProviderOptions{ProviderKey: "slack", UserOpenAPI: "./openapi/slack.yaml"},
			source:  ResolutionSourceUserOpenAPI,
		},
		{
			name:    "project local openapi",
			options: ResolveProviderOptions{ProviderKey: "slack", ProjectLocalOpenAPI: "./openapi/slack.yaml"},
			source:  ResolutionSourceProjectLocalOpenAPI,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resolved, err := ResolveProvider(test.options)
			if err != nil {
				t.Fatalf("ResolveProvider() error = %v", err)
			}
			if resolved.OpenAPI.Source != test.source || resolved.OpenAPI.Kind != SpecKindOpenAPI || resolved.OpenAPI.Protocol != SpecProtocolOpenAPI {
				t.Fatalf("resolved explicit source = %#v, want %q/openapi/openapi", resolved.OpenAPI, test.source)
			}
		})
	}
}

func TestResolveProviderRejectsUnknownProvider(t *testing.T) {
	if _, err := ResolveProvider(ResolveProviderOptions{ProviderKey: "missing"}); err == nil {
		t.Fatalf("ResolveProvider() expected unknown provider error")
	}
}

func TestResolveProviderRejectsUnsupportedURLScheme(t *testing.T) {
	_, err := ResolveProvider(ResolveProviderOptions{
		ProviderKey: "slack",
		UserOpenAPI: "ftp://example.com/slack.yaml",
	})
	if err == nil {
		t.Fatalf("ResolveProvider() expected unsupported URL scheme error")
	}
}

func TestSecurityReportRowsReturnCopies(t *testing.T) {
	report, err := BuiltInSecurityReport()
	if err != nil {
		t.Fatal(err)
	}
	rows := SecurityReportRows(report)
	rows[0].OverlayIDs = append(rows[0].OverlayIDs, "mutated")
	if len(rows[0].Dispositions) > 0 {
		rows[0].Dispositions[0].SourceRefs = append(rows[0].Dispositions[0].SourceRefs, "mutated")
	}
	fresh := SecurityReportRows(report)
	for _, id := range fresh[0].OverlayIDs {
		if id == "mutated" {
			t.Fatalf("SecurityReportRows leaked overlay ids slice")
		}
	}
	if len(fresh[0].Dispositions) > 0 {
		for _, ref := range fresh[0].Dispositions[0].SourceRefs {
			if ref == "mutated" {
				t.Fatal("SecurityReportRows leaked disposition source refs slice")
			}
		}
	}
}
