package apitools

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/catalogupdate"
)

func TestReviewedProviderSearchExactAndCoveredOperations(t *testing.T) {
	client := &Client{}
	alias, err := client.SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{Query: " SLACK  web api "})
	if err != nil || len(alias.Results) == 0 || alias.Results[0].Profile.ProviderID != "slack" ||
		alias.Results[0].Match != ReviewedProviderExactMatch {
		t.Fatalf("exact reviewed alias: %+v, %v", alias.Results, err)
	}

	message, err := client.SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{Query: "send slack message", Limit: 1})
	if err != nil || len(message.Results) != 1 || message.Results[0].Profile.ProviderID != "slack" ||
		message.Results[0].Match != ReviewedProviderOperationMatch || len(message.Results[0].MatchedCapabilities) != 1 {
		t.Fatalf("Slack operation lookup: %+v, %v", message.Results, err)
	}
	capability := message.Results[0].MatchedCapabilities[0]
	if capability.Selector != "#/paths/~1chat.postMessage/post" || capability.ArtifactID != "slack-web-message-history-overlay" ||
		len(capability.SHA256) != 64 || message.Results[0].Profile.Coverage.Sources[0].VerifiedAt == "" {
		t.Fatalf("Slack operation provenance: %+v", capability)
	}
	if len(message.MetadataSHA256) != 64 || message.MetadataSHA256 != alias.MetadataSHA256 {
		t.Fatalf("inconsistent metadata digest: %q vs %q", message.MetadataSHA256, alias.MetadataSHA256)
	}

	for _, query := range []string{"refund a payment", "refund"} {
		refunds, err := client.SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{Query: query, Limit: 25})
		if err != nil {
			t.Fatal(err)
		}
		selected := map[string]bool{}
		for _, result := range refunds.Results {
			if result.Match == ReviewedProviderOperationMatch {
				selected[result.Profile.ProviderID] = true
			}
		}
		if !selected["stripe"] || !selected["square"] || !selected["mercado-pago"] || selected["razorpay"] {
			t.Fatalf("unexpected capability results for %q: %v", query, selected)
		}
	}

	identity, err := client.SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{Query: "intuit quickbook"})
	if err != nil || len(identity.Results) == 0 || identity.Results[0].Profile.ProviderID != "quickbooks" ||
		identity.Results[0].Match != ReviewedProviderMetadataMatch || len(identity.Results[0].MatchedCapabilities) != 0 {
		t.Fatalf("stemmed identity invented operation evidence: %+v, %v", identity, err)
	}

	for _, query := range []string{"galactic zeppelin API", "(-*~) <>", "widgets \" OR x"} {
		got, err := client.SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{Query: query})
		if err != nil || len(got.Results) != 0 {
			t.Fatalf("unexpected matches or syntax error on %q: %v, %v", query, len(got.Results), err)
		}
	}
}

func TestReviewedSearchFiltersIncidentalCapabilityWords(t *testing.T) {
	for _, tt := range []struct {
		query, relevantProvider string
	}{
		{"accounting", "quickbooks"},
		{"marketing", "mailchimp"},
		{"news", "hackernews"},
		{"maps", "google-maps-platform"},
		{"listings", ""},
		{"copper api", "copper"},
		{"get account", ""},
	} {
		t.Run(tt.query, func(t *testing.T) {
			report, err := (&Client{}).SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{Query: tt.query, Limit: 334})
			if err != nil {
				t.Fatal(err)
			}
			found := false
			for _, result := range report.Results {
				if result.Match == ReviewedProviderOperationMatch || len(result.MatchedCapabilities) > 0 {
					t.Fatalf("incidental word asserted capability %s: %+v", result.Profile.ProviderID, result.MatchedCapabilities)
				}
				found = found || result.Profile.ProviderID == tt.relevantProvider
			}
			if tt.relevantProvider != "" && !found {
				t.Fatalf("missing relevant provider %s", tt.relevantProvider)
			}
		})
	}
}

func TestReviewedSearchCapabilityEvidenceMatrix(t *testing.T) {
	const slackPost = "#/paths/~1chat.postMessage/post"
	const squareRefund = "#/paths/~1v2~1refunds/post"
	const stripeRefund = "#/paths/~1v1~1refunds/post"
	const mercadoRefund = "#/paths/~1v1~1payments~1{id}~1refunds/post"
	const wordpressUsers = "#/paths/~1users/get"
	const clockifyEntries = "#/paths/~1v1~1workspaces~1{workspace_id}~1time-entries/get"
	const bambooRequests = "#/paths/~1v1~1time_off~1requests/get"
	const jotformFolders = "#/paths/~1user~1folders/get"
	cases := []struct {
		query, provider, selector string
		want                      bool
	}{
		{"send slack message", "slack", slackPost, true},
		{"send slack messages", "slack", slackPost, true},
		{"refund a Square payment", "square", squareRefund, true},
		{"refund Square payments", "square", squareRefund, true},
		{"refunds", "stripe", stripeRefund, true},
		{"refunds", "square", squareRefund, true},
		{"refunds", "mercado-pago", mercadoRefund, true},
		{"refund payments", "stripe", stripeRefund, true},
		{"refund payments", "square", squareRefund, true},
		{"refund payments", "mercado-pago", mercadoRefund, true},
		{"list user", "wordpress", wordpressUsers, true},
		{"list users", "wordpress", wordpressUsers, true},
		{"list clockify time entry", "clockify", clockifyEntries, true},
		{"list clockify time entries", "clockify", clockifyEntries, true},
		{"list user folders", "jotform", jotformFolders, true},
		{"list time off requests", "bamboohr", bambooRequests, true},
		{"get product", "gumroad", "#/paths/~1v2~1products~1{product_id}/get", true},
		{"get product", "shopify", "#/paths/~1products~1{productId}.json/get", true},
		{"list tag", "wordpress", "#/paths/~1tags/get", true},
		{"list job", "cortex", "#/paths/~1job/get", true},
		{"get npm user", "npm", "#/paths/~1-~1whoami/get", true},
		{"get reddit user", "reddit", "#/paths/~1api~1v1~1me/get", true},
		{"get current facebook user", "facebook", "#/paths/~1me/get", true},
		{"get mandrill user", "mandrill", "#/paths/~1users~1info.json/post", true},
		{"get linkedin user", "linkedin", "#/paths/~1v2~1userinfo/get", true},
		{"get jotform user", "jotform", "#/paths/~1user/get", true},
		{"get reddit user", "reddit", "#/paths/~1user~1{username}~1about/get", true},
		{"get current dropbox account", "dropbox", "#/paths/~12~1users~1get_current_account/post", true},
		{"get current digitalocean account", "digitalocean", "#/paths/~1v2~1account/get", true},
		{"fetch copper user", "copper", "#/paths/~1users~1me/get", true},
		{"telegram updates", "telegram", "#/paths/~1getUpdates/post", true},
		{"list user", "eventbrite", "#/paths/~1users~1me~1organizations~1/get", false},
		{"list current user", "eventbrite", "#/paths/~1users~1me~1organizations~1/get", false},
		{"list user", "jotform", jotformFolders, false},
		{"list user", "medium", "#/paths/~1users~1{user_id}~1publications/get", false},
		{"list user", "nvidia-dsx-air", "#/paths/~1api~1v3~1announcements~1/get", false},
		{"time", "bamboohr", bambooRequests, false},
		{"time", "clockify", clockifyEntries, false},
		{"time", "google-calendar", "#/paths/~1freeBusy/post", false},
		{"get product", "plaid", "#/paths/~1processor~1token~1permissions~1get/post", false},
		{"get one product", "plaid", "#/paths/~1processor~1token~1permissions~1get/post", false},
		{"telegram update", "telegram", "#/paths/~1getUpdates/post", false},
		{"time", "acuity-scheduling", "#/paths/~1availability~1times/get", false},
	}
	client := &Client{}
	var truePositives, expected, predicted int
	for _, tc := range cases {
		report, err := client.SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{
			Query: tc.query, Limit: len(catalog.BuiltInProviderSearchProfiles()),
		})
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for _, result := range report.Results {
			if result.Profile.ProviderID != tc.provider {
				continue
			}
			for _, capability := range result.MatchedCapabilities {
				if capability.Selector == tc.selector {
					found = true
					if result.Match != ReviewedProviderOperationMatch || capability.SpecRefID == "" || len(capability.SHA256) != 64 {
						t.Errorf("%q lost source evidence for %s %s", tc.query, tc.provider, tc.selector)
					}
				}
			}
		}
		if found {
			predicted++
		}
		if tc.want {
			expected++
			if found {
				truePositives++
			}
		}
		if found != tc.want {
			t.Errorf("%q %s %s: capability matched = %t, want %t", tc.query, tc.provider, tc.selector, found, tc.want)
		}
	}
	t.Logf("precision %d/%d, recall %d/%d across %d fixed query-selector cases", truePositives, predicted, truePositives, expected, len(cases))
}

func TestReviewedSearchPinnedSelectors(t *testing.T) {
	priorDir := os.Getenv("APITOOLS_CATALOG_QUALIFICATION_DIR")
	curatedDir := os.Getenv("APITOOLS_C07_QUALIFICATION_DIR")
	if priorDir == "" || curatedDir == "" {
		t.Skip("set both retained publisher artifact directories for exact selector qualification")
	}
	files := map[string]string{
		"amazon-sp-sellers-v1-swagger":          filepath.Join(priorDir, "openapi/amazon-sp-sellers-v1-swagger.json"),
		"amplitude-developer-openapi":           filepath.Join(curatedDir, "amplitude-developer.json"),
		"datadog-api-v2":                        filepath.Join(priorDir, "openapi/datadog-api-v2.yaml"),
		"digitalocean-api-v2":                   filepath.Join(priorDir, "openapi/digitalocean-api-v2.yaml"),
		"ebay-sell-account-v1":                  filepath.Join(priorDir, "openapi/ebay-sell-account-v1.json"),
		"etsy-v3-openapi":                       filepath.Join(priorDir, "openapi/etsy-v3-openapi.json"),
		"google-ads-v26-customer-service":       filepath.Join(priorDir, "grpc-protobuf/google-ads-v26-customer-service.proto"),
		"google-gemini-discovery-v1beta":        filepath.Join(priorDir, "google-discovery/google-gemini-discovery-v1beta.json"),
		"google-maps-places-v1-discovery":       filepath.Join(priorDir, "google-discovery/google-maps-places-v1-discovery.json"),
		"sheets-discovery-v4":                   filepath.Join(curatedDir, "sheets-canonical-v4.json"),
		"sheets-sdk-discovery-v4":               filepath.Join(curatedDir, "sheets-current-sdk.json"),
		"hubspot-crm-contacts-2026-09-openapi":  filepath.Join(curatedDir, "hubspot-contacts-2026-09-424.json"),
		"hubspot-crm-deals-2026-09-openapi":     filepath.Join(curatedDir, "hubspot-deals-2026-09-424.json"),
		"hugging-face-hub-openapi":              filepath.Join(priorDir, "openapi/hugging-face-hub-openapi.json"),
		"hunter-v2-openapi":                     filepath.Join(curatedDir, "hunter-openapi.json"),
		"kubernetes-v1-37-2-core-openapi-v3":    filepath.Join(curatedDir, "kubernetes-commit-core-v3.json"),
		"linear-sdk-graphql-sdl":                filepath.Join(curatedDir, "linear-schema.graphql"),
		"mercado-pago-rest-openapi":             filepath.Join(priorDir, "openapi/mercado-pago-rest-openapi.json"),
		"miro-rest-v2-openapi":                  filepath.Join(priorDir, "openapi/miro-rest-v2-openapi.json"),
		"mongodb-atlas-admin-openapi":           filepath.Join(priorDir, "openapi/mongodb-atlas-admin-openapi.yaml"),
		"nvidia-dsx-air-openapi":                "catalog-openapi-cache/openapi/nvidia-dsx-air-openapi.yaml",
		"plaid-openapi-2020-09-14":              filepath.Join(priorDir, "openapi/plaid-openapi-2020-09-14.yml"),
		"square-connect-v2-openapi":             filepath.Join(priorDir, "openapi/square-connect-v2-openapi.json"),
		"storyblok-content-delivery-v2-openapi": filepath.Join(curatedDir, "storyblok-cdn-v2.yaml"),
		"stripe-latest-openapi-spec3":           filepath.Join(curatedDir, "stripe-ga.json"),
		"vercel-rest-openapi":                   filepath.Join(priorDir, "openapi/vercel-rest-openapi.json"),
	}
	for _, profile := range catalog.BuiltInProviderSearchProfiles() {
		for _, source := range profile.Coverage.Sources {
			if source.Kind != catalog.ProviderSearchPinnedSource {
				continue
			}
			t.Run(source.SpecRefID, func(t *testing.T) {
				path, ok := files[source.SpecRefID]
				if !ok {
					t.Fatal("missing retained publisher artifact mapping")
				}
				delete(files, source.SpecRefID)
				content, err := readLocalSpecFile(path, DefaultMaxBytes)
				if err != nil || fmt.Sprintf("%x", sha256.Sum256(content)) != source.SHA256 {
					t.Fatalf("artifact fails the catalog pin: %v", err)
				}
				kind := catalogRegisteredSourceKind(string(source.Protocol))
				selectors := make(map[string]bool)
				if kind == OperationSourceOpenAPI {
					paths := mapValue(catalogUpdateNativeRoot(t.Context(), content)["paths"])
					for route, operations := range paths {
						for method, operation := range mapValue(operations) {
							if mapValue(operation) != nil {
								selectors[openAPIOperationSelector(method, route)] = true
							}
						}
					}
				} else if kind == OperationSourceGoogleDiscovery {
					inventory, err := BuildAPISourceOperationInventory(t.Context(), APISourceInventoryOptions{Documents: []APISourceDocument{{
						Kind: APISourceKindGoogleDiscovery, Content: content,
					}}})
					if err != nil {
						t.Fatal(err)
					}
					for _, operation := range inventory.Operations {
						selectors["#/methods/"+escapeJSONPointer(operation.OperationID)] = true
					}
				} else {
					table, err := BuildOperationShapeTable(t.Context(), OperationShapeOptions{Sources: []ShapeSourceInput{{
						ID: source.SpecRefID, OperationSourceInput: OperationSourceInput{Kind: kind, Content: content},
					}}})
					if err != nil {
						t.Fatal(err)
					}
					for _, operation := range table.Operations {
						selectors[operation.Selector.Value] = true
					}
				}
				for _, operation := range source.Operations {
					if !selectors[operation.Selector] {
						t.Errorf("selected operation %q absent in exact %s source", operation.Selector, kind)
					}
				}
			})
		}
	}
	if len(files) > 0 {
		t.Errorf("unused retained publisher artifact paths: %v", files)
	}
}

func TestReviewedProviderPublisherObservationIdentity(t *testing.T) {
	source := catalog.ProviderSearchSource{SpecRefID: "example-spec", SourceURL: "https://example.test/spec", SHA256: "digest-a"}
	observed := catalogupdate.Version{
		ReferenceID: source.SpecRefID, Locator: source.SourceURL,
		Object: catalogupdate.Object{SHA256: source.SHA256},
	}
	if !reviewedPublisherObservationMatches(source, observed) {
		t.Fatal("an exact legacy object did not match the reviewed source")
	}
	observed.Object.SHA256 = "digest-new"
	if reviewedPublisherObservationMatches(source, observed) {
		t.Fatal("unrelated object matched")
	}
	observed.CatalogPin = source.SHA256
	if !reviewedPublisherObservationMatches(source, observed) {
		t.Fatal("accepted successor did not match its anchor")
	}
	observed.Locator += "/other"
	if reviewedPublisherObservationMatches(source, observed) {
		t.Fatal("unrelated locator matched")
	}
}

func TestReviewedProviderIndexRefreshAndAtomicRollback(t *testing.T) {
	path := filepath.Join(t.TempDir(), "reviewed.sqlite")
	db, err := openReviewedProviderIndex(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	profiles := []catalog.ProviderSearchProfile{{
		ProviderID: "example", DisplayName: "Example", Description: "opal warehousing", Aliases: []string{"example api"},
	}}
	first, digest, err := rankReviewedProviders(t.Context(), db, profiles, "opal warehousing")
	if err != nil || len(first) != 1 {
		t.Fatalf("first snapshot: %d, %q, %v", len(first), digest, err)
	}
	profiles[0].Description = "birch carpentry"
	updated, current, err := rankReviewedProviders(t.Context(), db, profiles, "birch carpentry")
	if err != nil || len(updated) != 1 || current == digest {
		t.Fatalf("updated snapshot: %d, %q, %v", len(updated), current, err)
	}
	if stale, _, err := rankReviewedProviders(t.Context(), db, profiles, "opal warehousing"); err != nil || len(stale) != 0 {
		t.Fatalf("stale FTS records survived: %d, %v", len(stale), err)
	}
	rejected := append(append([]catalog.ProviderSearchProfile{}, profiles...), profiles[0])
	if _, _, err := rankReviewedProviders(t.Context(), db, rejected, "birch carpentry"); err == nil {
		t.Fatal("duplicate IDs unexpectedly replaced the good snapshot")
	}
	var digestAfterFailure string
	if err := db.QueryRow("SELECT digest FROM search_metadata").Scan(&digestAfterFailure); err != nil || digestAfterFailure != current {
		t.Fatalf("failed rebuild altered the snapshot: %q, %v", digestAfterFailure, err)
	}
	if got, _, err := rankReviewedProviders(t.Context(), db, profiles, "birch carpentry"); err != nil || len(got) != 1 {
		t.Fatalf("good snapshot lost after failure: %d, %v", len(got), err)
	}
}

func TestReviewedProviderIndexRejectsForeignDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "foreign.sqlite")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("CREATE TABLE example (note TEXT); INSERT INTO example VALUES ('unchanged');"); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = (&Client{}).SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{Query: "slack", IndexDB: path})
	if err == nil || !strings.Contains(err.Error(), "another application") {
		t.Fatalf("foreign database check: %v", err)
	}
	after, err := os.ReadFile(path)
	if err != nil || string(before) != string(after) {
		t.Fatalf("foreign database changed: %v", err)
	}
}

func TestReviewedProviderSearchConcurrentInitialization(t *testing.T) {
	path := filepath.Join(t.TempDir(), "reviewed.sqlite")
	missingState := filepath.Join(t.TempDir(), "missing.sqlite")
	var wg sync.WaitGroup
	for _, query := range []string{"slack web api", "refund payment"} {
		wg.Go(func() {
			report, err := (&Client{}).SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{
				Query: query, IndexDB: path, StateDB: missingState, Limit: 1,
			})
			if err != nil || len(report.Results) != 1 || report.PublisherState != "unavailable" {
				t.Errorf("concurrent lookup %q: %d, %q, %v", query, len(report.Results), report.PublisherState, err)
			}
		})
	}
	wg.Wait()
	if _, err := os.Stat(missingState); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("publisher state was unexpectedly created: %v", err)
	}
}

func TestReviewedProviderSearchReadsAcceptedStateWithoutChangingIt(t *testing.T) {
	config := newCatalogRepairFixture(t)
	provider, ok := catalog.FindBuiltInProvider("adyen")
	if !ok {
		t.Fatal("missing Adyen fixture provider")
	}
	source := catalogUpdateReferenceSources(provider)[0]
	seedCatalogUpdateFixture(t, config, source, catalogRepairFixture)
	before, err := os.ReadFile(config.StateDB)
	if err != nil {
		t.Fatal(err)
	}
	report, err := (&Client{}).SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{
		Query: "adyen", StateDB: config.StateDB, Limit: 1,
	})
	if err != nil || report.PublisherState != "available" || len(report.Results) != 1 ||
		len(report.Results[0].PublisherChecks) != 1 {
		t.Fatalf("accepted-state lookup did not match the exact reference: %+v, %v", report, err)
	}
	check := report.Results[0].PublisherChecks[0]
	if check.SpecRefID != source.ReferenceID || check.CheckedAt.IsZero() || len(check.SHA256) != 64 {
		t.Fatalf("incorrect accepted-state evidence: %+v", check)
	}
	after, err := os.ReadFile(config.StateDB)
	if err != nil || string(before) != string(after) {
		t.Fatalf("state file changed during search: %v", err)
	}
}
