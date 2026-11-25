package apitools

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/catalogupdate"
)

func TestC07PinnedPublisherOperations(t *testing.T) {
	qualificationDir := os.Getenv("APITOOLS_C07_QUALIFICATION_DIR")
	if qualificationDir == "" {
		t.Skip("set APITOOLS_C07_QUALIFICATION_DIR to the pinned C07 source directory")
	}
	for _, fixture := range []struct {
		provider, reference, file, method, path, operationID string
	}{
		{"hubspot", "hubspot-crm-contacts-2026-09-openapi", "hubspot-contacts-2026-09-424.json", "GET", "/crm/objects/2026-09/{objectType}", "get-/crm/objects/2026-09/{objectType}_/crm/objects/2026-03/{objectType}"},
		{"hubspot", "hubspot-crm-deals-2026-09-openapi", "hubspot-deals-2026-09-424.json", "GET", "/crm/objects/2026-09/{objectType}", "get-/crm/objects/2026-09/{objectType}_/crm/objects/2026-03/{objectType}"},
		{"kubernetes", "kubernetes-v1-37-2-core-openapi-v3", "kubernetes-commit-core-v3.json", "GET", "/api/v1/namespaces/{namespace}/pods", "listCoreV1NamespacedPod"},
		{"amplitude", "amplitude-developer-openapi", "amplitude-developer.json", "GET", "/v1/projects", "listProjects"},
		{"hunter", "hunter-v2-openapi", "hunter-openapi.json", "GET", "/domain-search", "searchDomain"},
		{"storyblok", "storyblok-content-delivery-v2-openapi", "storyblok-cdn-v2.yaml", "GET", "/v2/cdn/links", "listLinks"},
	} {
		t.Run(fixture.reference, func(t *testing.T) {
			provider, ok := catalog.FindBuiltInProvider(fixture.provider)
			if !ok {
				t.Fatalf("missing provider %s", fixture.provider)
			}
			var source catalogUpdateSource
			for _, candidate := range catalogUpdateReferenceSources(provider) {
				if candidate.ReferenceID == fixture.reference {
					source = candidate
					break
				}
			}
			content, err := os.ReadFile(filepath.Join(qualificationDir, fixture.file))
			if err != nil {
				t.Fatal(err)
			}
			digest := fmt.Sprintf("%x", sha256.Sum256(content))
			if catalogUpdateReviewedDigest(source.HeldRevision) != digest {
				t.Fatalf("catalog pin did not identify the retained artifact: %s", digest)
			}
			if validation := catalogUpdateNativeOpenAPI(t.Context(), content, source.URL, source.refreshReference(), digest); !validation.accepted() {
				t.Fatalf("reviewed source failed updater validation: %+v", validation)
			}
			if expected, ok := map[string]string{"amplitude": "Proprietary", "hunter": "", "storyblok": "MIT"}[fixture.provider]; ok {
				root := catalogUpdateNativeRoot(t.Context(), content)
				license := mapValue(mapValue(root["info"])["license"])
				if got := stringValue(license["name"]); got != expected {
					t.Fatalf("source license changed to %q, want %q", got, expected)
				}
				if expected == "" {
					expected = "No document-level license declaration"
				}
				for _, ref := range provider.SpecReferences {
					if ref.ID == fixture.reference && !strings.Contains(ref.LicenseNote, expected) {
						t.Fatalf("catalog license %q omits %q for %s", ref.LicenseNote, expected, ref.ID)
					}
				}
			}
			table, err := BuildOperationShapeTable(t.Context(), OperationShapeOptions{Sources: []ShapeSourceInput{{
				ID: fixture.reference, OperationSourceInput: OperationSourceInput{Kind: OperationSourceOpenAPI, URL: source.URL, Content: content},
			}}})
			if err != nil || len(table.Sources) != 1 || table.Sources[0].SHA256 != digest {
				t.Fatalf("native source identity not established: %v", err)
			}
			for _, shape := range table.Operations {
				if shape.Method == fixture.method && shape.Path == fixture.path {
					for _, alias := range shape.Aliases {
						if alias.Value == fixture.operationID {
							return
						}
					}
					t.Fatalf("operation %s %s lacks exact alias %s", fixture.method, fixture.path, fixture.operationID)
				}
			}
			t.Fatalf("native source omitted %s %s", fixture.method, fixture.path)
		})
	}
}

func TestC07HubSpotReleaseCompatibility(t *testing.T) {
	dir := os.Getenv("APITOOLS_C07_QUALIFICATION_DIR")
	if dir == "" {
		t.Skip("set APITOOLS_C07_QUALIFICATION_DIR to the pinned C07 source directory")
	}
	const collection = "https://raw.githubusercontent.com/HubSpot/HubSpot-public-api-spec-collection/c6f46a91db551f1f4fcbe61050ea002050e27d55/PublicApiSpecs/CRM/"
	const selectedKey = "#/paths/~1crm~1objects~12026-09~1{objectType}/get"
	for _, service := range []struct {
		name, group, legacySelector, legacyDigest, selectedDigest, rolloutDigest string
	}{
		{"contacts", "Contacts", "#/paths/~1crm~1v3~1objects~1contacts/get", "f144b194fc68e7d6ea338b7f89882d247479a80437c0f581456c10f1cc90bd5c", "f67bf5699a481ef85735c1846884f3c7c72ea85f858341de3499a90d9039f344", "0f2fb5042f2fd9c18d70e053ed8503978bb9cd641cc82c40e9d841d7f63911e9"},
		{"deals", "Deals", "#/paths/~1crm~1v3~1objects~10-3/get", "3b9f0d5a0f6b0684aa4210c661ccadb13df513238c01b811f13902a5a4cd2cf0", "cafe9ee81824c7994fdf3298ef483d5c7fddcebd869ad02ea8b23b09968f1aa1", "ef6fa7c6149224af71c811b1e86c1b6b0904fda46019525781a66d7cccb84d23"},
	} {
		t.Run(service.name, func(t *testing.T) {
			read := func(file, url, expectedDigest string) OperationSourceInput {
				t.Helper()
				content, err := os.ReadFile(filepath.Join(dir, file))
				if err != nil {
					t.Fatal(err)
				}
				if fmt.Sprintf("%x", sha256.Sum256(content)) != expectedDigest {
					t.Fatalf("qualification artifact changed: %s", file)
				}
				return OperationSourceInput{URL: collection + service.group + "/Rollouts/" + url + "/" + service.name + ".json", Content: content, Kind: OperationSourceOpenAPI}
			}
			prefix := "hubspot-" + service.name
			legacy := read(prefix+"-v3-424.json", "424/v3", service.legacyDigest)
			selected := read(prefix+"-2026-09-424.json", "424/2026-09", service.selectedDigest)
			alternate := read(prefix+"-2026-09.json", "321895/2026-09", service.rolloutDigest)
			for _, comparison := range []struct {
				baseline, candidate OperationSourceInput
				key                 string
			}{
				{legacy, selected, service.legacySelector},
				{selected, alternate, selectedKey},
			} {
				table, err := BuildOperationShapeTable(t.Context(), OperationShapeOptions{Sources: []ShapeSourceInput{{
					ID: "baseline", OperationSourceInput: comparison.baseline,
				}}})
				if _, ok := indexUpdateShapes(table.Operations)[comparison.key]; err != nil || !ok {
					t.Fatalf("baseline does not reproduce the reviewed key %s: %v", comparison.key, err)
				}
				got := compareCatalogUpdate(t.Context(), comparison.baseline, comparison.candidate, []string{comparison.key}, catalog.RefreshableSpecReference{}, "")
				if got != catalogUpdateChanged && got != catalogUpdateUnexamined {
					t.Fatalf("unproven HubSpot migration was marked %s for %s", got, comparison.key)
				}
				t.Logf("%s: %s", comparison.key, got)
			}
		})
	}
}

func TestC07HubSpotSecuritySchemesHaveSpecScope(t *testing.T) {
	view, err := catalog.BuiltInSecurityInspectionView("hubspot")
	if err != nil {
		t.Fatal(err)
	}
	for _, conflict := range view.Conflicts {
		if conflict.Type == catalog.SecurityInspectionConflictDuplicateScheme {
			t.Fatalf("HubSpot qualification introduced a conflicting scheme: %+v", conflict)
		}
	}
}

func TestC07StripeAndSheetsPinnedVersions(t *testing.T) {
	dir := os.Getenv("APITOOLS_C07_QUALIFICATION_DIR")
	if dir == "" {
		t.Skip("set APITOOLS_C07_QUALIFICATION_DIR to the pinned C07 source directory")
	}
	read := func(name, digest string) []byte {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		if fmt.Sprintf("%x", sha256.Sum256(data)) != digest {
			t.Fatalf("qualification artifact changed: %s", name)
		}
		return data
	}
	ga := read("stripe-ga.json", "e1fd0c9ed163b6ec53594ac30d60b251d7cef5a0ee911760bc89cdceafbe227b")
	preview := read("stripe-preview.json", "572ffda97268638f84807dbbc1bcdbf643bc75688d89065895f686c19cbc032c")
	var gaDoc map[string]any
	if err := json.Unmarshal(ga, &gaDoc); err != nil {
		t.Fatal(err)
	}
	var previewDoc map[string]any
	if err := json.Unmarshal(preview, &previewDoc); err != nil {
		t.Fatal(err)
	}
	if mapValue(gaDoc["info"])["version"] != "2026-11-18.endive" || mapValue(previewDoc["info"])["version"] != "2026-11-18.preview" {
		t.Fatal("Stripe GA and preview artifacts must remain separate")
	}
	stripe, _ := catalog.FindBuiltInProvider("stripe")
	source := catalogUpdateReferenceSources(stripe)[0]
	if pin := catalogUpdateReviewedDigest(source.HeldRevision); pin != fmt.Sprintf("%x", sha256.Sum256(ga)) || !strings.Contains(source.URL, "/latest/") ||
		!catalogUpdateNativeOpenAPI(t.Context(), ga, source.URL, source.refreshReference(), pin).accepted() {
		t.Fatal("Stripe GA pin failed updater validation")
	}
	get := mapValue(mapValue(mapValue(gaDoc["paths"])["/v1/payment_intents/{intent}"])["get"])
	response := mapValue(mapValue(get["responses"])["200"])
	schema := mapValue(mapValue(mapValue(response["content"])["application/json"])["schema"])
	if get["operationId"] != "GetPaymentIntentsIntent" || schema["$ref"] != "#/components/schemas/payment_intent" || get["security"] != nil {
		t.Fatal("Stripe PaymentIntent retrieval contract changed")
	}
	params, _ := get["parameters"].([]any)
	intentPath := false
	for _, param := range params {
		info := mapValue(param)
		intentPath = intentPath || info["name"] == "intent" && info["in"] == "path" && info["required"] == true
	}
	if !intentPath {
		t.Fatal("Stripe PaymentIntent retrieval lacks a required intent path parameter")
	}
	schemes := mapValue(mapValue(gaDoc["components"])["securitySchemes"])
	security, ok := gaDoc["security"].([]any)
	if !ok || len(security) != 2 || mapValue(schemes["basicAuth"])["scheme"] != "basic" || mapValue(schemes["bearerAuth"])["scheme"] != "bearer" ||
		mapValue(security[0])["basicAuth"] == nil || mapValue(security[1])["bearerAuth"] == nil {
		t.Fatal("Stripe GA does not provide the expected Basic and Bearer alternatives")
	}

	sheets, _ := catalog.FindBuiltInProvider("google-sheets")
	for _, version := range []struct{ file, revision, digest, reference string }{
		{"sheets-prior-sdk.json", "20260213", "60b9c02762346133a6ff3ee24b3c3f9e1f469666fb772048c0ea0a231030cae1", ""},
		{"sheets-current-sdk.json", "20261109", "1ec492cc174692abaf513b32286b6e455aa68b08e17b60b8fba21afeaec7925d", "sheets-sdk-discovery-v4"},
		{"sheets-canonical-v4.json", "20261118", "82c209f3b283637103a1332ec34b350cb617716a3c6d956e89de44784a9e1573", "sheets-discovery-v4"},
	} {
		data := read(version.file, version.digest)
		doc, err := decodeShapeDocument(t.Context(), data)
		if err != nil || doc["version"] != "v4" || doc["revision"] != version.revision {
			t.Fatalf("unexpected Sheets Discovery version in %s: %v", version.file, err)
		}
		inventory, err := BuildAPISourceOperationInventory(t.Context(), APISourceInventoryOptions{Documents: []APISourceDocument{{
			Kind: APISourceKindGoogleDiscovery, Content: data,
		}}})
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for _, operation := range inventory.Operations {
			found = found || operation.OperationID == "sheets.spreadsheets.values.get" && operation.Method == "GET" &&
				strings.TrimPrefix(operation.Path, "/") == "v4/spreadsheets/{spreadsheetId}/values/{range}"
		}
		if !found {
			t.Fatalf("missing Sheets values.get in revision %s", version.revision)
		}
		var reviewed catalogUpdateSource
		for _, candidate := range catalogUpdateReferenceSources(sheets) {
			if candidate.ReferenceID == version.reference {
				reviewed = candidate
			}
		}
		if version.reference == "" {
			repin := catalogUpdateSource{
				ProviderID: "google-sheets", ReferenceID: "sheets-sdk-discovery-v4",
				URL:          "https://raw.githubusercontent.com/googleapis/google-api-go-client/e9015cc2681ce9d6300f4fae2ae6b075c6d75ead/sheets/v4/sheets-api.json",
				HeldRevision: "sha256:" + version.digest,
			}
			if _, accepted := qualifyCatalogUpdateNativeSource(t.Context(), repin, OperationSourceGoogleDiscovery, repin.URL, version.digest, data); accepted {
				t.Fatal("a catalog repin alone authorized the reduced Sheets proof")
			}
			continue
		}
		proof, accepted := qualifyCatalogUpdateNativeSource(t.Context(), reviewed, OperationSourceGoogleDiscovery, reviewed.URL, version.digest, data)
		if !accepted || proof == nil || proof.ReviewedOperation != "sheets.spreadsheets.values.get" {
			t.Fatalf("missing exact-pin Sheets proof for %s", version.reference)
		}
		reviewed.HeldRevision = ""
		if _, accepted := qualifyCatalogUpdateNativeSource(t.Context(), reviewed, OperationSourceGoogleDiscovery, reviewed.URL, version.digest, data); accepted {
			t.Fatal("unreviewed Sheets bytes passed the selected-method proof")
		}
	}
}

func TestC07SlackAdvisoryOperations(t *testing.T) {
	p := "catalog-openapi-cache/advisory-overlays/slack-web-message-history-overlay.json"
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	table, err := BuildOperationShapeTable(t.Context(), OperationShapeOptions{Sources: []ShapeSourceInput{{
		ID: "slack", OperationSourceInput: OperationSourceInput{Kind: OperationSourceOpenAPI, Content: data},
	}}})
	if err != nil || len(table.Operations) != 2 {
		t.Fatalf("Slack advisory native operations: %d, %v", len(table.Operations), err)
	}
	doc, err := decodeShapeDocument(t.Context(), data)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []struct{ method, path, scope string }{
		{"POST", "/chat.postMessage", "chat:write"},
		{"GET", "/conversations.history", "channels:history"},
	} {
		operation := mapValue(mapValue(mapValue(doc["paths"])[expected.path])[strings.ToLower(expected.method)])
		security, ok := operation["security"].([]any)
		if !ok || len(security) != 1 {
			t.Fatalf("%s %s lacks one scoped requirement", expected.method, expected.path)
		}
		scopes, ok := mapValue(security[0])["slackOAuth2"].([]any)
		if !ok || len(scopes) != 1 || scopes[0] != expected.scope {
			t.Fatalf("%s %s lacks scoped OAuth requirement %s", expected.method, expected.path, expected.scope)
		}
	}
}

func TestC07SheetsSDKReusesVerifiedObject(t *testing.T) {
	dir := os.Getenv("APITOOLS_C07_QUALIFICATION_DIR")
	if dir == "" {
		t.Skip("set APITOOLS_C07_QUALIFICATION_DIR to the pinned C07 source directory")
	}
	content, err := os.ReadFile(filepath.Join(dir, "sheets-current-sdk.json"))
	if err != nil {
		t.Fatal(err)
	}
	provider, _ := catalog.FindBuiltInProvider("google-sheets")
	var source catalogUpdateSource
	for _, candidate := range catalogUpdateReferenceSources(provider) {
		if candidate.ReferenceID == "sheets-sdk-discovery-v4" {
			source = candidate
		}
	}
	commit, _ := catalogUpdateGitFile(source.URL)
	if commit == "" {
		t.Fatal("missing the pinned Google SDK file locator")
	}
	ctx := t.Context()
	config := newCatalogRepairFixture(t)
	cache, err := catalogupdate.OpenArtifactCache(config.CacheDir, "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	store, err := catalogupdate.OpenStore(ctx, config.StateDB, cache)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	var artifactRequests int
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "api.github.com" && req.URL.Path == "/repos/googleapis/google-api-go-client/commits" {
			return catalogRepairResponse(req, fmt.Sprintf(`[{"sha":%q}]`, commit)), nil
		}
		if req.URL.String() != source.URL {
			t.Fatalf("unexpected request %s", req.URL)
		}
		artifactRequests++
		return catalogRepairResponse(req, string(content)), nil
	})}}
	for i, status := range []string{"reviewed-pin-reproduced", "unchanged"} {
		run, err := store.StartRun(ctx, catalog.BuiltInCatalogManifest().CatalogSHA256, []string{provider.ID})
		if err != nil {
			t.Fatal(err)
		}
		outcome, err := client.checkCatalogUpdateSource(ctx, store, run.ID, provider, source, config, nil)
		if err != nil || outcome.Status != status || artifactRequests != 2 {
			t.Fatalf("SDK attempt %d status %s used %d artifact requests: %v", i+1, outcome.Status, artifactRequests, err)
		}
		if i == 0 {
			var evidence catalogUpdateEvidence
			if err := json.Unmarshal(outcome.Evidence, &evidence); err != nil || evidence.Source.NativeRevision != "20261109" || evidence.Source.Revision != commit || evidence.Native == nil {
				t.Fatalf("missing separate native and Git revision evidence: %v", err)
			}
		}
		preferred, err := store.PreferredVersion(ctx, provider.ID, source.ReferenceID, source.Service)
		if err != nil || preferred.Revision != commit {
			t.Fatalf("SDK preference lost the immutable revision: %v", err)
		}
	}
}

func TestC07CatalogUpdateReportsLocalAdvisoryCoverage(t *testing.T) {
	for _, fixture := range []struct{ provider, status string }{
		{"airtable", "advisory-coverage"},
		{"odoo", "documented-source-gap"},
	} {
		t.Run(fixture.provider, func(t *testing.T) {
			client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(req *http.Request) (*http.Response, error) {
				t.Fatalf("local coverage made a publisher request to %s", req.URL)
				return nil, nil
			})}}
			options := CatalogUpdateOptions{Provider: fixture.provider, Config: newCatalogRepairFixture(t)}
			report, err := client.UpdateCatalog(t.Context(), options)
			if err != nil || !report.Completed || len(report.Outcomes) != 1 || report.Outcomes[0].Status != fixture.status {
				t.Fatalf("provider coverage status: %+v, %v", report.Outcomes, err)
			}
			var data catalog.AdvisoryEndpointCoverage
			if err := json.Unmarshal(report.Outcomes[0].Evidence, &data); err != nil || data.Status != fixture.status {
				t.Fatalf("provider coverage evidence: %+v, %v", data, err)
			}
			if fixture.provider == "airtable" && (len(data.EndpointOverlays) != 1 || data.EndpointOverlays[0].BuilderPath == "") {
				t.Fatalf("checked-in overlay evidence missing: %+v", data)
			}
			if fixture.provider == "odoo" && (data.Gap == "" || len(data.EndpointOverlays) != 0) {
				t.Fatalf("documented gap evidence missing: %+v", data)
			}
			resumed, err := client.UpdateCatalog(t.Context(), CatalogUpdateOptions{Resume: report.RunID, Config: CatalogUpdateConfig{StateDB: options.Config.StateDB}})
			if err != nil || len(resumed.Outcomes) != 1 || string(resumed.Outcomes[0].Evidence) != string(report.Outcomes[0].Evidence) {
				t.Fatalf("saved provider evidence changed: %+v, %v", resumed.Outcomes, err)
			}
		})
	}
	if got := catalog.BuiltInAdvisoryEndpointCoverage("unmeasured"); got.Status != "unexamined-source-coverage" || got.Gap != "" {
		t.Fatalf("unmeasured coverage was treated as a documented gap: %+v", got)
	}
}

func TestC07CuratedPublisherSecurityAndSelection(t *testing.T) {
	for _, fixture := range []struct {
		provider, reference string
		protocol            catalog.SpecProtocol
		status              catalog.AuthCompletenessStatus
	}{
		{"amplitude", "amplitude-developer-openapi", catalog.SpecProtocolOpenAPI, catalog.AuthStatusPresentIncomplete},
		{"hunter", "hunter-v2-openapi", catalog.SpecProtocolOpenAPI, catalog.AuthStatusComplete},
		{"linear", "linear-sdk-graphql-sdl", catalog.SpecProtocolGraphQL, catalog.AuthStatusOverlayRequired},
		{"storyblok", "storyblok-content-delivery-v2-openapi", catalog.SpecProtocolOpenAPI, catalog.AuthStatusComplete},
	} {
		t.Run(fixture.provider, func(t *testing.T) {
			got, err := catalog.ResolveProvider(catalog.ResolveProviderOptions{ProviderKey: fixture.provider})
			if err != nil || got.OpenAPI.SpecRefID != fixture.reference || got.OpenAPI.Protocol != fixture.protocol ||
				got.SecurityStatus != fixture.status {
				t.Fatalf("resolved source or security: %+v, %v", got, err)
			}
		})
	}
}

func TestC07PinnedLinearSchemaAndRejectedPostmarkCandidates(t *testing.T) {
	dir := os.Getenv("APITOOLS_C07_QUALIFICATION_DIR")
	if dir == "" {
		t.Skip("set APITOOLS_C07_QUALIFICATION_DIR to the pinned C07 source directory")
	}
	linear, _ := catalog.FindBuiltInProvider("linear")
	source := catalogUpdateReferenceSources(linear)[0]
	read := func(file, expected string) []byte {
		t.Helper()
		content, err := os.ReadFile(filepath.Join(dir, file))
		if err != nil || fmt.Sprintf("%x", sha256.Sum256(content)) != expected {
			t.Fatalf("artifact %s did not reproduce %s: %v", file, expected, err)
		}
		return content
	}
	content := read("linear-schema.graphql", catalogUpdateReviewedDigest(source.HeldRevision))
	_, accepted := qualifyCatalogUpdateNativeSource(t.Context(), source, OperationSourceGraphQL, source.URL, catalogUpdateReviewedDigest(source.HeldRevision), content)
	if !accepted {
		t.Fatal("Linear pin failed native updater qualification")
	}
	inventory, err := BuildAPISourceOperationInventory(t.Context(), APISourceInventoryOptions{Documents: []APISourceDocument{{
		Kind: APISourceKindGraphQL, Content: content,
	}}})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, op := range inventory.Operations {
		if op.OperationID == "query.viewer" && op.Method == "query" && op.Path == "#/schema/Query/fields/viewer" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("Linear GraphQL schema omitted Query.viewer")
	}
	for _, fixture := range []struct{ file, digest string }{
		{"postmark-server.yml", "fdae25f55e2ef5fb71997f1128cc4148d52e0e395aed04b4f8f9e1c86fdf44a2"},
		{"postmark-account.yml", "bbbcda9432851b8d7d6781d9a090b9d8d70f4396a9288393f428bfa3d5c2102b"},
	} {
		content := read(fixture.file, fixture.digest)
		url := "https://postmarkapp.com/swagger/" + strings.TrimPrefix(fixture.file, "postmark-")
		got := catalogUpdateNativeOpenAPI(t.Context(), content, url, catalog.RefreshableSpecReference{}, fixture.digest)
		if got.accepted() || got.RawStatus != CatalogRefreshInvalid {
			t.Fatalf("the recorded Postmark qualification needs to be reviewed: %+v", got)
		}
	}
}

func TestC07EndpointRegistrationsKeepSourceIdentity(t *testing.T) {
	artifacts := catalog.BuiltInEndpointOverlayArtifacts()
	for _, expected := range []struct{ artifact, reference string }{
		{"slack-web-message-history-overlay", "slack-web-api-docs"},
		{"airtable-web-api-overlay", "airtable-web-api-overlay"},
	} {
		found := false
		for _, artifact := range artifacts {
			if artifact.ArtifactID == expected.artifact {
				found = true
				if artifact.SpecRefID != expected.reference {
					t.Fatalf("%s source identity = %s, want %s", artifact.ArtifactID, artifact.SpecRefID, expected.reference)
				}
			}
		}
		if !found {
			t.Fatalf("missing registry artifact %s", expected.artifact)
		}
	}
	report, err := catalog.BuiltInProviderAdvisoryReport(catalog.ProviderAdvisoryOptions{ProviderKey: "slack", Artifacts: artifacts})
	if err != nil || len(report.Providers) != 1 {
		t.Fatalf("Slack advisory: %+v, %v", report, err)
	}
	for _, ref := range report.Providers[0].SpecReferences {
		if ref.ID == "slack-web-api-docs" && ref.RegisteredArtifactPath == "advisory-overlays/slack-web-message-history-overlay.json" {
			return
		}
	}
	t.Fatalf("Slack docs reference missed the checked-in advisory path: %+v", report.Providers[0].SpecReferences)
}

func TestC07FallbackSecurityMatchesReviewedOperations(t *testing.T) {
	for _, provider := range []string{"postmark", "storyblok"} {
		t.Run(provider, func(t *testing.T) {
			path := filepath.Join("catalog-openapi-cache", "advisory-overlays", provider+"-api-overlay.json")
			body, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			var artifact struct {
				Paths map[string]map[string]struct {
					OperationID string                `json:"operationId"`
					Security    []map[string][]string `json:"security"`
				} `json:"paths"`
			}
			if err := json.Unmarshal(body, &artifact); err != nil {
				t.Fatal(err)
			}
			overlays := catalog.SecurityOverlaysForProvider(provider)
			for route, item := range artifact.Paths {
				for method, operation := range item {
					if operation.OperationID == "" {
						continue
					}
					reference := "postmark-api-docs"
					if provider == "storyblok" {
						reference = "storyblok-content-delivery-api-docs"
						if strings.HasPrefix(route, "/v1/spaces/") {
							reference = "storyblok-management-api-docs"
						}
					}
					matches := 0
					for _, overlay := range overlays {
						if len(overlay.RootSecurity) != 0 || len(overlay.RootSecuritySets) != 0 {
							t.Fatalf("%s supplied fallback security without a route", overlay.ID)
						}
						for _, rule := range overlay.OperationSecurity {
							if rule.Match.OperationID != operation.OperationID || rule.Match.Method != strings.ToUpper(method) || rule.Match.Path != route {
								continue
							}
							matches++
							if overlay.SpecRefID != reference || len(operation.Security) != 1 ||
								len(operation.Security[0]) != 1 || len(rule.Security) != 1 || len(rule.SecuritySets) != 0 {
								t.Fatalf("unexpected scoped security for %s %s: %s, %+v", method, route, overlay.SpecRefID, rule)
							}
							if _, ok := operation.Security[0][rule.Security[0].Scheme]; !ok {
								t.Fatalf("security scheme mismatch on %s %s: %+v", method, route, rule)
							}
						}
					}
					if matches != 1 {
						t.Fatalf("%s %s has %d matching catalog security rules", method, route, matches)
					}
				}
			}
		})
	}
}
