package apitools_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/artifactio"
	"github.com/OpenUdon/apitools/sqlitecache"
	"github.com/OpenUdon/uws/binding"
)

// Official source bytes are kept in the ignored curation cache. Setting this
// directory makes their presence and exact identity mandatory for qualification.
const catalogPromotionArtifactDir = "APITOOLS_CATALOG_QUALIFICATION_DIR"

type catalogPromotionCase struct {
	providerID string
	specID     string
	artifactID string
	kind       string
	sourceKind apitools.OperationSourceKind
	path       string
	digest     string
	bytes      int64
	selector   string
	contract   apitools.StepContract
	purpose    apitools.ContractMatchStatus
	shapeTable bool
}

var catalogPromotionCases = []catalogPromotionCase{
	{
		providerID: "amazon-selling-partner", specID: "amazon-sp-sellers-v1-swagger", artifactID: "amazon-sp-sellers-v1-swagger",
		kind: "swagger", sourceKind: apitools.OperationSourceOpenAPI, path: "openapi/amazon-sp-sellers-v1-swagger.json",
		digest: "497862ea32de8040453649986e2cd7c6fcc15b55e8022becc783e4a6d6ffcffd", bytes: 29604,
		selector: "#/paths/~1sellers~1v1~1marketplaceParticipations/get",
		contract: apitools.StepContract{Purpose: "list marketplaces seller participation", Effect: apitools.OperationEffectRead},
	},
	{
		providerID: "datadog", specID: "datadog-api-v2", artifactID: "datadog-api-v2",
		kind: "openapi", sourceKind: apitools.OperationSourceOpenAPI, path: "openapi/datadog-api-v2.yaml",
		digest: "37702169e22c9f10aaf840d83b7e621896b17463033ae5ff06a5b91af2802189", bytes: 10574047,
		selector: "#/paths/~1api~1v2~1synthetics~1settings~1on_demand_concurrency_cap/get",
		contract: apitools.StepContract{Purpose: "get on demand concurrency cap", Effect: apitools.OperationEffectRead},
	},
	{
		providerID: "digitalocean", specID: "digitalocean-api-v2", artifactID: "digitalocean-api-v2",
		kind: "openapi", sourceKind: apitools.OperationSourceOpenAPI, path: "openapi/digitalocean-api-v2.yaml",
		digest: "120b1f5ca7ac46159c0428c318c4b1e88598d923081ba338ebe6838f8c1a08ba", bytes: 4606025,
		selector: "#/paths/~1v2~1account/get",
		contract: apitools.StepContract{Purpose: "get current user account information", Effect: apitools.OperationEffectRead},
	},
	{
		providerID: "google-gemini", specID: "google-gemini-discovery-v1beta", artifactID: "google-gemini-discovery-v1beta",
		kind: "google-discovery", sourceKind: apitools.OperationSourceGoogleDiscovery, path: "google-discovery/google-gemini-discovery-v1beta.json",
		digest: "62024ad3940881c88048b9c0c04f72860c9821a89e36422c3f33401a21d4472f", bytes: 401495,
		selector: "#/methods/generativelanguage.models.list",
		contract: apitools.StepContract{Purpose: "list models available through gemini api", Effect: apitools.OperationEffectRead},
	},
	{
		providerID: "google-gemini", specID: "google-gemini-models-docs", artifactID: "google-gemini-models-overlay",
		kind: "advisory-overlay", sourceKind: apitools.OperationSourceOpenAPI, path: "advisory-overlays/google-gemini-models-overlay.json",
		digest: "8e2838356f8c3773b34553f6df1ff7c9c7701bca073239a3dd12590ca04bdf85", bytes: 2556,
		selector: "#/paths/~1v1beta~1models/get",
		contract: apitools.StepContract{Purpose: "list available gemini models", Effect: apitools.OperationEffectRead},
	},
	{
		providerID: "hugging-face-hub", specID: "hugging-face-hub-openapi", artifactID: "hugging-face-hub-openapi",
		kind: "openapi", sourceKind: apitools.OperationSourceOpenAPI, path: "openapi/hugging-face-hub-openapi.json",
		digest: "da5be4f2ee566534ee0c5fd74e1c380cc4575f19ca7a90c7fed4c51cc360c01d", bytes: 1247465,
		selector: "#/paths/~1api~1models~1{namespace}~1{repo}~1refs/get",
		contract: apitools.StepContract{Purpose: "list references", Effect: apitools.OperationEffectRead,
			Inputs: map[string]apitools.ContractValue{"namespace": {Type: "string"}, "repo": {Type: "string"}}},
	},
	{
		providerID: "miro", specID: "miro-rest-v2-openapi", artifactID: "miro-rest-v2-openapi",
		kind: "openapi", sourceKind: apitools.OperationSourceOpenAPI, path: "openapi/miro-rest-v2-openapi.json",
		digest: "f0e17b5b43c92ff2d85c7fbe9a84301e39df1d3c14e6738cd939a6d7142f44dd", bytes: 1065633,
		selector: "#/paths/~1v2~1boards~1{board_id}/get",
		contract: apitools.StepContract{Purpose: "get specific board", Effect: apitools.OperationEffectRead,
			Inputs: map[string]apitools.ContractValue{"board_id": {Type: "string"}}},
	},
	{
		providerID: "mongodb-atlas", specID: "mongodb-atlas-admin-openapi", artifactID: "mongodb-atlas-admin-openapi",
		kind: "openapi", sourceKind: apitools.OperationSourceOpenAPI, path: "openapi/mongodb-atlas-admin-openapi.yaml",
		digest: "8eaeb645ded2fb91b6c9c07d29b5315ebb9a8d5ae4397f32b512457d74df34c5", bytes: 5196444,
		selector: "#/paths/~1api~1atlas~1v2~1groups/get",
		contract: apitools.StepContract{Purpose: "return all projects", Effect: apitools.OperationEffectRead},
	},
	{
		providerID: "plaid", specID: "plaid-openapi-2020-09-14", artifactID: "plaid-openapi-2020-09-14",
		kind: "openapi", sourceKind: apitools.OperationSourceOpenAPI, path: "openapi/plaid-openapi-2020-09-14.yml",
		digest: "59246af3f2aef2e55f11a19af1baf6dd1aa19497c3c49a6fc3fe24c68634717f", bytes: 3453430,
		selector: "#/paths/~1processor~1token~1permissions~1get/post",
		contract: apitools.StepContract{Purpose: "get processor token product permissions", Effect: apitools.OperationEffectRead,
			Outputs: map[string]apitools.ContractValue{"products": {Type: "array"}}},
	},
	{
		providerID: "vercel", specID: "vercel-rest-openapi", artifactID: "vercel-rest-openapi",
		kind: "openapi", sourceKind: apitools.OperationSourceOpenAPI, path: "openapi/vercel-rest-openapi.json",
		digest: "ae3e778a91928ef7753675d31ca52d89dbd5fb2247e4a6e404f4f986124ea027", bytes: 12497764,
		selector: "#/paths/~1v9~1projects~1{idOrName}~1domains/get",
		contract: apitools.StepContract{Purpose: "retrieve domains associated with a given project", Effect: apitools.OperationEffectRead,
			Inputs: map[string]apitools.ContractValue{"idOrName": {Type: "string"}}},
	},
	{
		providerID: "google-maps-platform", specID: "google-maps-places-v1-discovery", artifactID: "google-maps-places-v1-discovery",
		kind: "google-discovery", sourceKind: apitools.OperationSourceGoogleDiscovery, path: "google-discovery/google-maps-places-v1-discovery.json",
		digest: "e6db1c21f8b1c51f9ff1b5495b1bbafdf33cac27e470f0db2c7bad3afa0f7331", bytes: 148978,
		selector: "#/methods/places.places.autocomplete",
		contract: apitools.StepContract{Purpose: "returns predictions for the given input", Effect: apitools.OperationEffectRead,
			Inputs: map[string]apitools.ContractValue{"input": {Type: "string"}}},
	},
	{
		providerID: "google-ads", specID: "google-ads-v26-customer-service", artifactID: "google-ads-v26-customer-service",
		kind: "grpc-protobuf", sourceKind: apitools.OperationSourceGRPCProtobuf, path: "grpc-protobuf/google-ads-v26-customer-service.proto",
		digest: "0cca59c738498fb0010f39f437d29eaff5dbe7577aae4bcdd011ae445795e685", bytes: 7296,
		selector: "#/services/google.ads.googleads.v26.services.CustomerService/methods/ListAccessibleCustomers",
		contract: apitools.StepContract{Purpose: "list accessible customers", Effect: apitools.OperationEffectRead},
		purpose:  apitools.ContractMatchIndeterminate,
	},
	{
		providerID: "google-maps-platform", specID: "google-maps-places-v1-docs", artifactID: "google-maps-places-autocomplete-overlay",
		kind: "advisory-overlay", sourceKind: apitools.OperationSourceOpenAPI, path: "advisory-overlays/google-maps-places-autocomplete-overlay.json",
		digest: "c5f393a1ee6135633361c7ea7be680a2c77e43f28974663e75117347f21a1794", bytes: 5326,
		selector: "#/paths/~1v1~1places:autocomplete/post",
		contract: apitools.StepContract{Purpose: "return place and query predictions for the given input", Effect: apitools.OperationEffectRead,
			Inputs: map[string]apitools.ContractValue{"input": {Type: "string"}}},
		shapeTable: true,
	},
	{
		providerID: "google-ads", specID: "google-ads-v26-customer-docs", artifactID: "google-ads-v26-accessible-customers-overlay",
		kind: "advisory-overlay", sourceKind: apitools.OperationSourceOpenAPI, path: "advisory-overlays/google-ads-v26-accessible-customers-overlay.json",
		digest: "bc432312fb516895f61ccd864cd320270003e881ee3c4c3eb60b8570235ce63a", bytes: 2639,
		selector: "#/paths/~1v26~1customers:listAccessibleCustomers/get",
		contract: apitools.StepContract{Purpose: "list accessible customers", Effect: apitools.OperationEffectRead,
			Outputs: map[string]apitools.ContractValue{"resourceNames": {Type: "array"}}},
		shapeTable: true,
	},
	{
		providerID: "ebay", specID: "ebay-sell-account-v1", artifactID: "ebay-sell-account-v1",
		kind: "openapi", sourceKind: apitools.OperationSourceOpenAPI, path: "openapi/ebay-sell-account-v1.json",
		digest: "3f04cf1fd7c80160312d68be5afd7d0de00eb57fd32d508cbb5eeec526b8bfbb", bytes: 328066,
		selector: "#/paths/~1program~1get_opted_in_programs/get",
		contract: apitools.StepContract{Purpose: "list seller programs that the seller has opted in to", Effect: apitools.OperationEffectRead,
			Outputs: map[string]apitools.ContractValue{"programs": {Type: "array"}}},
		shapeTable: true,
	},
	{
		providerID: "etsy", specID: "etsy-v3-openapi", artifactID: "etsy-v3-openapi",
		kind: "openapi", sourceKind: apitools.OperationSourceOpenAPI, path: "openapi/etsy-v3-openapi.json",
		digest: "c68d65ed002ddb6d3381d7d4497850eaa49a80587fcd237f58df7fa52e9fdd0e", bytes: 963751,
		selector: "#/paths/~1v3~1application~1shops~1{shop_id}~1holiday-preferences/get",
		contract: apitools.StepContract{Purpose: "retrieve holiday preferences available to an Etsy shop", Effect: apitools.OperationEffectRead,
			Inputs: map[string]apitools.ContractValue{"shop_id": {Type: "integer"}}},
		shapeTable: true,
	},
	{
		providerID: "square", specID: "square-connect-v2-openapi", artifactID: "square-connect-v2-openapi",
		kind: "openapi", sourceKind: apitools.OperationSourceOpenAPI, path: "openapi/square-connect-v2-openapi.json",
		digest: "76e7cc8ba3fd94c162a3339531aafe798d5556ac34ae03e8253df815b49305f1", bytes: 3279383,
		selector: "#/paths/~1v2~1merchants~1{merchant_id}/get",
		contract: apitools.StepContract{Purpose: "retrieve merchant object for given merchant id", Effect: apitools.OperationEffectRead,
			Inputs: map[string]apitools.ContractValue{"merchant_id": {Type: "string"}}},
	},
	{
		providerID: "mercado-pago", specID: "mercado-pago-rest-openapi", artifactID: "mercado-pago-rest-openapi",
		kind: "openapi", sourceKind: apitools.OperationSourceOpenAPI, path: "openapi/mercado-pago-rest-openapi.json",
		digest: "da7e3a7482982b1b6880a5f71411edcab54ce354a8eda55a970ec6c09973592b", bytes: 364763,
		selector: "#/paths/~1v1~1payments~1{id}/get",
		contract: apitools.StepContract{Purpose: "get payment by ID", Effect: apitools.OperationEffectRead,
			Inputs: map[string]apitools.ContractValue{"id": {Type: "integer"}}},
		shapeTable: true,
	},
	{
		providerID: "razorpay", specID: "razorpay-payments-v1-docs", artifactID: "razorpay-v1-payment-overlay",
		kind: "advisory-overlay", sourceKind: apitools.OperationSourceOpenAPI, path: "advisory-overlays/razorpay-v1-payment-overlay.json",
		digest: "6ed5333b09f383fe20490433377ec488cc130cc454a432b7d29b1713c9af3f28", bytes: 2537,
		selector: "#/paths/~1v1~1payments~1{id}/get",
		contract: apitools.StepContract{Purpose: "fetch a payment using its ID", Effect: apitools.OperationEffectRead,
			Inputs:  map[string]apitools.ContractValue{"id": {Type: "string"}},
			Outputs: map[string]apitools.ContractValue{"status": {Type: "string"}}},
		shapeTable: true,
	},
	{
		providerID: "wechat-pay", specID: "wechat-pay-v3-global-order-docs", artifactID: "wechat-pay-v3-global-order-overlay",
		kind: "advisory-overlay", sourceKind: apitools.OperationSourceOpenAPI, path: "advisory-overlays/wechat-pay-v3-global-order-overlay.json",
		digest: "b18df240bfbc3a7505a4515c87dbcd617e99d06ee76c4cfaefc21ec55575c969", bytes: 4218,
		selector: "#/paths/~1v3~1global~1transactions~1id~1{id}/get",
		contract: apitools.StepContract{Purpose: "query WeChat Pay order status by order ID", Effect: apitools.OperationEffectRead,
			Inputs:  map[string]apitools.ContractValue{"id": {Type: "string"}, "mchid": {Type: "string"}},
			Outputs: map[string]apitools.ContractValue{"trade_state": {Type: "string"}}},
		shapeTable: true,
	},
	{
		providerID: "twitch", specID: "twitch-helix-games-docs", artifactID: "twitch-helix-games-overlay",
		kind: "advisory-overlay", sourceKind: apitools.OperationSourceOpenAPI, path: "advisory-overlays/twitch-helix-games-overlay.json",
		digest: "f1db616feef677676dd4c1be42c9cc24f1a8a9f6030282d4835028a5ee6a0781", bytes: 4268,
		selector: "#/paths/~1helix~1games/get",
		contract: apitools.StepContract{Purpose: "get information about specified categories or games", Effect: apitools.OperationEffectRead,
			Inputs:  map[string]apitools.ContractValue{"id": {Type: "array"}},
			Outputs: map[string]apitools.ContractValue{"data": {Type: "array"}}},
		shapeTable: true,
	},
}

func TestCatalogPromotion(t *testing.T) {
	for _, input := range catalogPromotionCases {
		name := input.providerID
		if input.kind == "advisory-overlay" {
			name += "/advisory-overlay"
		}
		t.Run(name, func(t *testing.T) {
			resolved, err := catalog.ResolveProvider(catalog.ResolveProviderOptions{ProviderKey: input.providerID})
			if err != nil || ((input.kind != "advisory-overlay" || resolved.OpenAPI.Protocol == catalog.SpecProtocolHumanDocs) && resolved.OpenAPI.SpecRefID != input.specID) {
				t.Fatalf("default metadata resolution for %s: got %v, %v", input.providerID, resolved.OpenAPI, err)
			}
			var sourceURL string
			for _, ref := range resolved.CatalogSpecReferences {
				if ref.ID == input.specID {
					sourceURL = ref.URL
					if input.kind != "advisory-overlay" && ref.ProtocolClassification().UWSSourceType() != string(input.sourceKind) {
						t.Fatalf("catalog source type %q, want %q", ref.ProtocolClassification().UWSSourceType(), input.sourceKind)
					}
					break
				}
			}
			if sourceURL == "" {
				t.Fatalf("missing catalog source reference %s/%s", input.providerID, input.specID)
			}
			if resolved.SecurityStatus == catalog.AuthStatusUnknown || resolved.SecurityStatus == catalog.AuthStatusConflict {
				t.Fatalf("missing or conflicting auth review for %s: %s", input.providerID, resolved.SecurityStatus)
			}
			t.Run("pinned-artifact", func(t *testing.T) {
				artifactDir := os.Getenv(catalogPromotionArtifactDir)
				if input.kind == "advisory-overlay" {
					artifactDir = "catalog-openapi-cache"
				}
				if artifactDir == "" {
					t.Skip(catalogPromotionArtifactDir + " is unset; provider source bytes are kept outside Git")
				}
				qualifyCatalogPromotion(t, input, artifactDir, sourceURL)
			})
		})
	}
}

func qualifyCatalogPromotion(t *testing.T, input catalogPromotionCase, artifactDir, sourceURL string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	root := catalog.RootOptions{Directory: t.TempDir()}
	raw, err := artifactio.ReadFile(artifactDir, input.path, artifactio.ReadOptions{
		MaxBytes: apitools.DefaultMaxBytes, SHA256: input.digest, Bytes: input.bytes,
	})
	if err != nil {
		t.Fatal(err)
	}
	staged := filepath.Join(root.Directory, filepath.FromSlash(input.path))
	if err := os.MkdirAll(filepath.Dir(staged), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(staged, raw.Data, 0o600); err != nil {
		t.Fatal(err)
	}
	registry, err := sqlitecache.OpenCatalog(root)
	if err != nil {
		t.Fatal(err)
	}
	registration := sqlitecache.CatalogArtifact{
		ProviderID: input.providerID, ArtifactID: input.artifactID, Kind: input.kind, Path: input.path,
		SourceURL: sourceURL, SHA256: input.digest, Bytes: input.bytes,
		Metadata: map[string]string{"spec_ref_id": input.specID},
	}
	if input.kind == "advisory-overlay" {
		registration.OverlayPath = input.path
	}
	registrationErr := registry.StoreCatalogArtifact(ctx, registration)
	closeErr := registry.Close()
	if registrationErr != nil || closeErr != nil {
		t.Fatalf("register pinned source: %v; close registry: %v", registrationErr, closeErr)
	}

	options := apitools.CatalogIndexOptions{Root: root, ReadRegistrations: sqlitecache.ReadCatalogSpecArtifacts}
	index, err := apitools.BuildCatalogOperationIndex(ctx, options)
	if err != nil {
		t.Fatal(err)
	}
	if len(index.Artifacts) != 1 || index.Artifacts[0].OperationCount == 0 {
		t.Fatalf("no indexed operation from the pinned artifact: %+v", index.Artifacts)
	}
	var operationID string
	var indexed bool
	for _, operation := range index.Artifacts[0].Operations {
		if operation.Source.Selector == input.selector {
			operationID = operation.Operation.OperationID
			indexed = true
			break
		}
	}
	if !indexed {
		for _, diagnostic := range index.Artifacts[0].Diagnostics {
			if diagnostic.Path == input.selector {
				t.Logf("index diagnostic for the selected operation: %s: %s", diagnostic.Code, diagnostic.Message)
			}
		}
		t.Fatalf("the pinned source did not index selector %s", input.selector)
	}
	if err := apitools.WriteCatalogOperationIndex(ctx, options, index); err != nil {
		t.Fatal(err)
	}

	report, err := apitools.DiscoverCatalogOperations(ctx, apitools.CatalogDiscoveryOptions{
		Index: options,
		Request: apitools.CatalogDiscoveryRequest{
			SchemaVersion: apitools.CatalogDiscoverySchemaVersion, ProviderKeys: []string{input.providerID},
			Contract: input.contract,
			Limits:   apitools.CatalogDiscoveryLimits{MaxResults: apitools.MaxCatalogDiscoveryResults, MaxContextBytes: apitools.MaxCatalogDiscoveryContextBytes},
		},
	})
	if err != nil || !report.Scope.ProviderConstrained || len(report.Scope.ProviderIDs) != 1 || report.Scope.ProviderIDs[0] != input.providerID {
		t.Fatalf("provider-scoped offline discovery: scope=%+v, error=%v", report.Scope, err)
	}
	found := false
	purpose := input.purpose
	if purpose == "" {
		purpose = apitools.ContractMatchCompatible
	}
	for _, candidate := range report.Candidates {
		source := candidate.Candidate.Source
		if source.Selector != input.selector || source.SHA256 != input.digest || source.Kind != input.sourceKind {
			continue
		}
		for _, ref := range candidate.References {
			if ref.ProviderID == input.providerID && ref.SpecRefID == input.specID && ref.ArtifactID == input.artifactID &&
				ref.SHA256 == input.digest && ref.Bytes == input.bytes && ref.Kind == input.sourceKind && ref.Selector == input.selector {
				if candidate.Candidate.Match.Purpose.Status != purpose {
					t.Fatalf("discovery purpose = %+v, want %s", candidate.Candidate.Match.Purpose, purpose)
				}
				found = true
				break
			}
		}
	}
	if !found {
		t.Fatalf("offline discovery did not return %s from %s/%s at %s (outcome=%s, examined=%d)",
			input.selector, input.providerID, input.specID, input.digest, report.Outcome, report.ExaminedOperations)
	}
	if input.providerID == "google-ads" && input.sourceKind == apitools.OperationSourceGRPCProtobuf {
		assertGoogleAdsRPCSecurity(t, input.specID, operationID)
	}

	artifacts, err := sqlitecache.ReadCatalogSpecArtifacts(ctx, root)
	if err != nil {
		t.Fatal(err)
	}
	materialized, err := catalog.MaterializeProvider(ctx, catalog.MaterializeOptions{
		ProviderKey: input.providerID, TargetDir: t.TempDir(), CacheDir: root.Directory, Artifacts: artifacts, WriteManifest: true,
	})
	if err != nil || materialized.Resolution.ProviderID != input.providerID || len(materialized.Resolution.Artifacts) != 1 || len(materialized.Artifacts) != 1 {
		t.Fatalf("resolve and materialize registered source: artifacts=%+v, err=%v", materialized.Artifacts, err)
	}
	resolvedArtifact := materialized.Resolution.Artifacts[0]
	if !resolvedArtifact.Materializable || resolvedArtifact.SpecRefID != input.specID || resolvedArtifact.ArtifactID != input.artifactID || resolvedArtifact.SHA256 != input.digest || resolvedArtifact.Bytes != input.bytes {
		t.Fatalf("resolved a different registration: %+v", resolvedArtifact)
	}
	artifact := materialized.Artifacts[0]
	if artifact.SpecRefID != input.specID || artifact.ArtifactID != input.artifactID || artifact.Kind != input.kind || artifact.SHA256 != input.digest || artifact.Bytes != input.bytes {
		t.Fatalf("materialized a different registration: %+v", artifact)
	}
	// Advisory assets retain the human-docs protocol of their catalog reference.
	// Their OpenAPI source type is established by indexing and native reproduction.
	if input.kind != "advisory-overlay" && artifact.UWSSourceType != string(input.sourceKind) {
		t.Fatalf("materialized source type %q, want %q", artifact.UWSSourceType, input.sourceKind)
	}
	reproduceCatalogPromotionSelector(t, ctx, input, artifact.TargetPath, sourceURL, operationID)
	t.Logf("qualified indexed artifact %s/%s %s %s (%d bytes; selector %s; %d candidates from %d visited operations)",
		input.providerID, input.specID, input.sourceKind, input.digest, input.bytes, input.selector,
		len(index.Artifacts[0].Operations), index.Artifacts[0].OperationCount)
}

func assertGoogleAdsRPCSecurity(t *testing.T, specID, operationID string) {
	t.Helper()
	view, err := catalog.BuildSecurityInspectionView(catalog.SecurityInspectionOptions{
		ProviderKey:     "google-ads",
		KnownOperations: []catalog.OperationMatch{{OperationID: operationID}},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, conflict := range view.Conflicts {
		if conflict.SpecRefID == specID && conflict.Type == catalog.SecurityInspectionConflictUnresolvedOperation {
			t.Fatalf("Google Ads security overlay did not match indexed RPC %q: %+v", operationID, conflict.Match)
		}
	}
	for _, security := range view.OperationSecurity {
		if security.SpecRefID != specID || security.Match.OperationID != operationID {
			continue
		}
		if len(security.SecuritySets) != 1 || len(security.SecuritySets[0].Requirements) != 2 {
			t.Fatalf("Google Ads RPC must have one combined OAuth and developer-token requirement: %+v", security.SecuritySets)
		}
		var oauth, developerToken bool
		for _, effective := range security.SecuritySets[0].Requirements {
			requirement := effective.Requirement
			switch requirement.Scheme {
			case "googleAdsOAuth":
				oauth = len(requirement.Scopes) == 1 && requirement.Scopes[0] == "https://www.googleapis.com/auth/adwords"
			case "googleAdsDeveloperToken":
				developerToken = len(requirement.Scopes) == 0
			}
		}
		if !oauth || !developerToken {
			t.Fatalf("incorrect Google Ads RPC security requirements: %+v", security.SecuritySets[0].Requirements)
		}
		return
	}
	t.Fatalf("missing source-scoped security requirements for indexed RPC %q", operationID)
}

func reproduceCatalogPromotionSelector(t *testing.T, ctx context.Context, input catalogPromotionCase, targetPath, sourceURL, operationID string) {
	t.Helper()
	artifact, err := artifactio.ReadFile(filepath.Dir(targetPath), filepath.Base(targetPath), artifactio.ReadOptions{
		MaxBytes: apitools.DefaultMaxBytes, SHA256: input.digest, Bytes: input.bytes,
	})
	if err != nil {
		t.Fatal(err)
	}
	if input.sourceKind == apitools.OperationSourceOpenAPI && !input.shapeTable {
		// Whole-source UWS shape tables have their own bounds. Reproduce the
		// selected OpenAPI or Swagger path and verb in the materialized source.
		inventory, err := apitools.BuildOperationInventory(ctx, apitools.InventoryOptions{
			Documents: []apitools.InventoryDocument{{Name: input.artifactID, Content: artifact.Data}},
		})
		if err != nil {
			t.Fatal(err)
		}
		escaped := strings.NewReplacer("~", "~0", "/", "~1")
		for _, operation := range inventory.Operations {
			selector := "#/paths/" + escaped.Replace(operation.Path) + "/" + strings.ToLower(operation.Method)
			if selector == input.selector && operation.OperationID == operationID {
				return
			}
		}
		t.Fatalf("materialized source did not reproduce operation %s, selector %s", operationID, input.selector)
	}
	table, err := apitools.BuildOperationShapeTable(ctx, apitools.OperationShapeOptions{Sources: []apitools.ShapeSourceInput{{
		ID:                   input.providerID,
		OperationSourceInput: apitools.OperationSourceInput{Kind: input.sourceKind, Content: artifact.Data, URL: sourceURL},
	}}})
	if err != nil {
		t.Fatalf("reproduce materialized operation: %v", err)
	}
	if len(table.Sources) != 1 || table.Sources[0].ID != input.providerID || table.Sources[0].Kind != string(input.sourceKind) || table.Sources[0].SHA256 != input.digest {
		t.Fatalf("materialized UWS source identity: %+v", table.Sources)
	}
	resolver, err := binding.NewResolver(table)
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := resolver.Resolve(ctx, binding.Binding{
		Source:       table.Sources[0],
		SelectorKind: "ref", SelectorValue: input.selector,
	})
	if err != nil || resolved.Status != binding.Resolved || resolved.Shape.Selector.Value != input.selector {
		t.Fatalf("reproduce materialized selector %s: status=%s, err=%v", input.selector, resolved.Status, err)
	}
}
