package catalog

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestBuiltInProviderSearchProfiles(t *testing.T) {
	profiles := BuiltInProviderSearchProfiles()
	providers := BuiltInProviders()
	if len(profiles) != len(providers) {
		t.Fatalf("got %d profiles for %d reviewed providers", len(profiles), len(providers))
	}
	identities := make(map[string]string)
	descriptions := make(map[string]string)
	seenStatus := make(map[ProviderSearchCoverageStatus]bool)
	for i, profile := range profiles {
		provider := providers[i]
		if profile.ProviderID != provider.ID || profile.DisplayName != provider.DisplayName {
			t.Fatalf("profile %d has wrong identity: %q, %q", i, profile.ProviderID, profile.DisplayName)
		}
		if prior := descriptions[profile.Description]; prior != "" {
			t.Errorf("providers %q and %q have the same description", prior, profile.ProviderID)
		}
		descriptions[profile.Description] = profile.ProviderID
		for _, alias := range provider.Aliases {
			if !slices.Contains(profile.Aliases, alias) {
				t.Errorf("profile %q omitted accepted alias %q", profile.ProviderID, alias)
			}
		}
		for _, key := range append([]string{profile.ProviderID, profile.DisplayName}, profile.Aliases...) {
			key = normalizeKey(key)
			if prior := identities[key]; prior != "" && prior != profile.ProviderID {
				t.Errorf("ambiguous identity %q between %q and %q", key, prior, profile.ProviderID)
			}
			identities[key] = profile.ProviderID
		}
		if len(profile.Coverage.Sources) == 0 {
			t.Errorf("profile %q has no source provenance", profile.ProviderID)
		}
		seenStatus[profile.Coverage.Status] = true
		for _, source := range profile.Coverage.Sources {
			if !strings.HasPrefix(source.SourceURL, "https://") || source.Authority == "" || source.SpecRefID == "" {
				t.Errorf("profile %q has incomplete source provenance: %+v", profile.ProviderID, source)
			}
			if profile.Coverage.Status == ProviderSearchSelectedSubset {
				if len(source.Operations) == 0 || !validProviderSearchDigest(source.SHA256) || source.Kind == ProviderSearchSourceReference {
					t.Errorf("profile %q has no selected operation evidence: %+v", profile.ProviderID, source)
				}
			} else if len(source.Operations) != 0 || source.Kind != ProviderSearchSourceReference ||
				(profile.Coverage.Status == ProviderSearchDocsOnly && source.Protocol != SpecProtocolHumanDocs) {
				t.Errorf("profile %q claims inconsistent source coverage: %+v", profile.ProviderID, source)
			}
		}
	}
	for _, status := range []ProviderSearchCoverageStatus{ProviderSearchSelectedSubset, ProviderSearchDocsOnly, ProviderSearchReferenceOnly} {
		if !seenStatus[status] {
			t.Errorf("missing coverage status %q", status)
		}
	}
	if identities["slack-web-api"] != "slack" || identities["seven.io"] != "sms77" || identities["square-connect"] != "square" {
		t.Fatal("missing reviewed or accepted aliases")
	}
	index := slices.IndexFunc(profiles, func(p ProviderSearchProfile) bool { return p.ProviderID == "slack" })
	if index < 0 {
		t.Fatal("missing Slack profile")
	}
	changed := &profiles[index]
	changed.Aliases[0] = "changed alias"
	changed.Keywords[0] = "changed keyword"
	changed.Coverage.Sources[0].Operations[0].Label = "changed operation"
	fresh := BuiltInProviderSearchProfiles()[index]
	if fresh.Aliases[0] == changed.Aliases[0] || fresh.Keywords[0] == changed.Keywords[0] ||
		fresh.Coverage.Sources[0].Operations[0].Label == changed.Coverage.Sources[0].Operations[0].Label {
		t.Fatal("a caller mutated the built-in search metadata")
	}
}

func TestProviderSearchCoverageExamples(t *testing.T) {
	profiles := make(map[string]ProviderSearchProfile)
	for _, profile := range BuiltInProviderSearchProfiles() {
		profiles[profile.ProviderID] = profile
	}
	for _, fixture := range []struct {
		id, category, label, selector, refID string
		kind                                 ProviderSearchSourceKind
	}{
		{"slack", "team-chat", "Send Slack message as text in a participating conversation", "#/paths/~1chat.postMessage/post", "slack-web-api-docs", ProviderSearchAdvisoryOverlay},
		{"stripe", "payments", "Refund payment by charge or payment intent", "#/paths/~1v1~1refunds/post", "stripe-latest-openapi-spec3", ProviderSearchPinnedSource},
		{"square", "payments", "Refund a Square payment", "#/paths/~1v2~1refunds/post", "square-connect-v2-openapi", ProviderSearchPinnedSource},
		{"mercado-pago", "payments", "Refund payment by ID", "#/paths/~1v1~1payments~1{id}~1refunds/post", "mercado-pago-rest-openapi", ProviderSearchPinnedSource},
		{"razorpay", "payments", "Fetch a payment using its ID", "#/paths/~1v1~1payments~1{id}/get", "razorpay-payments-v1-docs", ProviderSearchAdvisoryOverlay},
		{"nvidia-dsx-air", "cloud-infrastructure", "List announcements visible to an NVIDIA Air user", "#/paths/~1api~1v3~1announcements~1/get", "nvidia-dsx-air-openapi", ProviderSearchPinnedSource},
		{"plaid", "financial-data", "Get product permissions for one Plaid processor token", "#/paths/~1processor~1token~1permissions~1get/post", "plaid-openapi-2020-09-14", ProviderSearchPinnedSource},
		{"etsy", "commerce", "List available holidays for a US or Canadian shop to set preferences", "#/paths/~1v3~1application~1shops~1{shop_id}~1holiday-preferences/get", "etsy-v3-openapi", ProviderSearchPinnedSource},
	} {
		t.Run(fixture.id, func(t *testing.T) {
			profile := profiles[fixture.id]
			if profile.Category != fixture.category || profile.Coverage.Status != ProviderSearchSelectedSubset {
				t.Fatalf("category/coverage: got %q/%q", profile.Category, profile.Coverage.Status)
			}
			for _, source := range profile.Coverage.Sources {
				if source.Kind != fixture.kind || source.SpecRefID != fixture.refID {
					continue
				}
				for _, op := range source.Operations {
					if op.Selector == fixture.selector && op.Label == fixture.label {
						return
					}
				}
			}
			t.Fatalf("missing search selection %q in %q", fixture.label, fixture.refID)
		})
	}
	if profiles["acumatica"].Coverage.Status != ProviderSearchDocsOnly ||
		profiles["adyen"].Coverage.Status != ProviderSearchReferenceOnly {
		t.Fatal("a documentation or machine source lead acquired unsupported operation coverage")
	}
}

func TestProviderSearchAdvisoryRejectsUnrelatedReference(t *testing.T) {
	provider, ok := FindBuiltInProvider("action-network")
	if !ok {
		t.Fatal("missing Action Network")
	}
	artifacts := make(map[string]CatalogSpecArtifact)
	for _, artifact := range BuiltInEndpointOverlayArtifacts() {
		artifacts[artifact.ArtifactID] = artifact
	}
	var rows []providerSearchCoverageData
	if err := decodeProviderSearchData(providerSearchCoverageJSON, &rows); err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		if row.ProviderID == provider.ID {
			row.SpecRefID = "action-network-webhooks-docs"
			_, err := loadProviderSearchCoverage(provider, []providerSearchCoverageData{row}, artifacts)
			if err == nil || !strings.Contains(err.Error(), "does not match the provider and reference") {
				t.Fatalf("unrelated Action Network reference produced %v", err)
			}
			return
		}
	}
	t.Fatal("missing Action Network evidence")
}

func TestProviderSearchAdvisoryOperationPins(t *testing.T) {
	curatedLabels := map[string]map[string]string{
		"slack": {
			"#/paths/~1chat.postMessage/post":     "Send Slack message as text in a participating conversation",
			"#/paths/~1conversations.history/get": "Read public Slack channel message history",
		},
		"google-calendar": {
			"#/paths/~1freeBusy/post": "Check calendar availability and free or busy time",
		},
		"woocommerce": {
			"#/paths/~1wp-json~1wc~1v3~1orders~1{order_id}~1refunds/post": "Create an order refund",
		},
	}
	artifacts := make(map[string]CatalogSpecArtifact)
	for _, artifact := range BuiltInEndpointOverlayArtifacts() {
		artifacts[artifact.ArtifactID] = artifact
	}
	for _, profile := range BuiltInProviderSearchProfiles() {
		for _, source := range profile.Coverage.Sources {
			if source.Kind != ProviderSearchAdvisoryOverlay {
				continue
			}
			t.Run(source.ArtifactID, func(t *testing.T) {
				artifact, ok := artifacts[source.ArtifactID]
				if !ok || artifact.SourceURL != source.SourceURL {
					t.Fatalf("search evidence has no matching advisory registration: %+v", source)
				}
				delete(artifacts, source.ArtifactID)
				paths := readProviderSearchOpenAPIPaths(t, filepath.Join("..", "catalog-openapi-cache", artifact.Path), source.SHA256)
				for _, selection := range source.Operations {
					op := providerSearchOpenAPIOperation(t, paths, selection.Selector)
					var summary string
					if err := json.Unmarshal(op["summary"], &summary); err != nil || summary == "" {
						t.Fatalf("selection %q lacks a source operation summary", selection.Selector)
					}
					if label := curatedLabels[profile.ProviderID][selection.Selector]; label != "" {
						summary = label
					}
					if selection.Label != summary {
						t.Errorf("selection %q label %q diverged from %q", selection.Selector, selection.Label, summary)
					}
				}
			})
		}
	}
	if len(artifacts) != 0 {
		t.Errorf("missing operation search metadata for the advisory registrations %v", artifacts)
	}
}

func TestProviderSearchPinnedRefundOperations(t *testing.T) {
	for _, fixture := range []struct {
		provider, fixtureDirEnv, file, selector, operationID string
	}{
		{"stripe", "APITOOLS_C07_QUALIFICATION_DIR", "stripe-ga.json", "#/paths/~1v1~1refunds/post", "PostRefunds"},
		{"square", "APITOOLS_CATALOG_QUALIFICATION_DIR", "openapi/square-connect-v2-openapi.json", "#/paths/~1v2~1refunds/post", "RefundPayment"},
		{"mercado-pago", "APITOOLS_CATALOG_QUALIFICATION_DIR", "openapi/mercado-pago-rest-openapi.json", "#/paths/~1v1~1payments~1{id}~1refunds/post", "createRefund"},
	} {
		t.Run(fixture.provider, func(t *testing.T) {
			dir := os.Getenv(fixture.fixtureDirEnv)
			if dir == "" {
				t.Skipf("set %s to the retained publisher artifact directory", fixture.fixtureDirEnv)
			}
			for _, profile := range BuiltInProviderSearchProfiles() {
				if profile.ProviderID != fixture.provider {
					continue
				}
				for _, source := range profile.Coverage.Sources {
					for _, selection := range source.Operations {
						if selection.Selector != fixture.selector || source.Kind != ProviderSearchPinnedSource {
							continue
						}
						paths := readProviderSearchOpenAPIPaths(t, filepath.Join(dir, fixture.file), source.SHA256)
						op := providerSearchOpenAPIOperation(t, paths, fixture.selector)
						var id string
						if err := json.Unmarshal(op["operationId"], &id); err != nil || id != fixture.operationID {
							t.Fatalf("refund operation ID: got %q, want %q (%v)", id, fixture.operationID, err)
						}
						return
					}
				}
			}
			t.Fatalf("no pinned selection for %s", fixture.selector)
		})
	}
}

func readProviderSearchOpenAPIPaths(t *testing.T, file, digest string) map[string]map[string]json.RawMessage {
	t.Helper()
	content, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if got := fmt.Sprintf("%x", sha256.Sum256(content)); got != digest {
		t.Fatalf("artifact %q has digest %s, want %s", file, got, digest)
	}
	var document struct {
		Paths map[string]map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(content, &document); err != nil {
		t.Fatal(err)
	}
	return document.Paths
}

func providerSearchOpenAPIOperation(t *testing.T, paths map[string]map[string]json.RawMessage, selector string) map[string]json.RawMessage {
	t.Helper()
	value, ok := strings.CutPrefix(selector, "#/paths/")
	encodedPath, method, split := strings.Cut(value, "/")
	if !ok || !split || !slices.Contains([]string{"get", "put", "post", "delete", "options", "head", "patch", "trace"}, method) {
		t.Fatalf("invalid OpenAPI operation selector %q", selector)
	}
	path := strings.NewReplacer("~1", "/", "~0", "~").Replace(encodedPath)
	var operation map[string]json.RawMessage
	if err := json.Unmarshal(paths[path][method], &operation); err != nil || len(operation) == 0 {
		t.Fatalf("selector %q did not resolve to an operation: %v", selector, err)
	}
	return operation
}
