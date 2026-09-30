package apitools

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/OpenUdon/apitools/catalog"
)

func TestCatalogDiscoveryDecodeConfinesWire(t *testing.T) {
	for _, body := range []string{
		`{"schema_version":"apitools.catalog-discovery/v1","root":"/private"}`,
		`{"schema_version":"apitools.catalog-discovery/v1","remote_endpoint":"http://localhost"}`,
		`{"schema_version":"apitools.catalog-discovery/v1","contract":{"purpose":"list notes","secret":"hidden"}}`,
		`{} {}`,
		strings.Repeat("x", MaxCatalogDiscoveryRequestBytes+1),
	} {
		if _, err := DecodeCatalogDiscoveryRequest(strings.NewReader(body)); err == nil {
			t.Fatal("accepted unsafe or malformed request")
		}
	}
	if _, err := DecodeCatalogDiscoveryRequest(nil); err == nil {
		t.Fatal("accepted missing request")
	}
	request, err := DecodeCatalogDiscoveryRequest(strings.NewReader(`{"schema_version":"apitools.catalog-discovery/v1","contract":{"purpose":"list notes","effect":"read"},"provider_keys":[]}`))
	if err != nil || request.ProviderKeys == nil {
		t.Fatalf("empty explicit scope lost: %v", err)
	}
	encoded, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeCatalogDiscoveryRequest(bytes.NewReader(encoded))
	if err != nil || decoded.ProviderKeys == nil {
		t.Fatal("round trip broadened empty scope")
	}
}

func TestCatalogDiscoveryRequestExactProvidersAndBounds(t *testing.T) {
	cat := catalog.BuiltInCatalog()
	request := CatalogDiscoveryRequest{SchemaVersion: CatalogDiscoverySchemaVersion, Contract: StepContract{Purpose: "read calendar events", Effect: OperationEffectRead}, ProviderKeys: []string{"Microsoft Outlook", "microsoft-outlook"}}
	limits, ids, diagnostics := prepareCatalogDiscoveryRequest(request, cat)
	if len(diagnostics) != 0 || len(ids) != 1 || ids[0] != "microsoft-outlook" || limits.MaxOperations != DefaultCatalogDiscoveryOperations {
		t.Fatalf("exact keys: %#v %#v", ids, diagnostics)
	}
	request.ProviderKeys = []string{}
	_, ids, diagnostics = prepareCatalogDiscoveryRequest(request, cat)
	if len(ids) != 0 || len(diagnostics) != 0 {
		t.Fatal("empty selected scope broadened or invalidated")
	}
	request.ProviderKeys = nil
	_, ids, diagnostics = prepareCatalogDiscoveryRequest(request, cat)
	if len(ids) != len(cat.Providers) || len(diagnostics) != 0 {
		t.Fatal("open scope differs")
	}
	for _, mutate := range []func(*CatalogDiscoveryRequest){
		func(r *CatalogDiscoveryRequest) { r.SchemaVersion = "future" },
		func(r *CatalogDiscoveryRequest) { r.ProviderKeys = []string{"unknown provider"} },
		func(r *CatalogDiscoveryRequest) { r.Contract.Purpose = "" },
		func(r *CatalogDiscoveryRequest) { r.Contract.Effect = "purchase" },
		func(r *CatalogDiscoveryRequest) { r.Limits.MaxOperations = MaxCatalogIndexOperations + 1 },
		func(r *CatalogDiscoveryRequest) { r.Limits.MaxResults = -1 },
		func(r *CatalogDiscoveryRequest) {
			r.Limits.TimeoutMillis = int(MaxCatalogDiscoveryTimeout.Milliseconds()) + 1
		},
		func(r *CatalogDiscoveryRequest) { r.Limits.MaxContextBytes = MaxCatalogDiscoveryContextBytes + 1 },
		func(r *CatalogDiscoveryRequest) { r.Limits.MaxContextBytes = 1 },
		func(r *CatalogDiscoveryRequest) { r.ProviderKeys = make([]string, 33) },
		func(r *CatalogDiscoveryRequest) { r.Filters.Authorities = []string{"trusted-by-popularity"} },
	} {
		bad := request
		mutate(&bad)
		_, _, diagnostics := prepareCatalogDiscoveryRequest(bad, cat)
		if len(errorDiagnostics(diagnostics)) == 0 {
			t.Fatal("accepted invalid discovery request")
		}
	}
}

func TestCatalogDiscoveryFrozenWireCases(t *testing.T) {
	data, err := os.ReadFile("testdata/catalog-discovery/v1/wire-cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Name    string                  `json:"name"`
		Request CatalogDiscoveryRequest `json:"request"`
		Report  CatalogDiscoveryReport  `json:"report"`
	}
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	wanted := map[string]bool{"match": false, "ambiguous": false, "ties": false, "filtering": false, "unknowns": false, "no-root": false, "no-match": false, "insufficient": false, "blocked": false, "relocation": false}
	encodedReports := map[string][]byte{}
	for _, item := range cases {
		if _, ok := wanted[item.Name]; !ok {
			t.Fatal("unknown fixture")
		}
		wanted[item.Name] = true
		if item.Request.SchemaVersion != CatalogDiscoverySchemaVersion || item.Report.SchemaVersion != CatalogDiscoverySchemaVersion {
			t.Fatal("wire version drift")
		}
		switch item.Report.Outcome {
		case CatalogDiscoveryMatch, CatalogDiscoveryAmbiguous, CatalogDiscoveryNoQualifyingAPI, CatalogDiscoveryInsufficientEvidence, CatalogDiscoveryBlocked:
		default:
			t.Fatal("invalid outcome")
		}
		if item.Report.Outcome == CatalogDiscoveryAmbiguous && item.Report.QualifiedOperations < 2 {
			t.Fatal("fixture hides ambiguity")
		}
		encoded, err := json.Marshal(item.Report)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(encoded, []byte("/home/")) {
			t.Fatal("absolute installation path in wire")
		}
		encodedReports[item.Name] = encoded
		if item.Name == "unknowns" && (item.Report.Candidates[0].Sources[0].LicenseIdentifier != "unknown" || item.Report.Candidates[0].Sources[0].Redistribution != "unknown") {
			t.Fatal("notes inferred into permission")
		}
		var roundTrip CatalogDiscoveryReport
		if err := json.Unmarshal(encoded, &roundTrip); err != nil {
			t.Fatal(err)
		}
		again, _ := json.Marshal(roundTrip)
		if !bytes.Equal(encoded, again) {
			t.Fatal("wire round trip drift")
		}
	}
	for name, present := range wanted {
		if !present {
			t.Fatalf("missing fixture %s", name)
		}
	}
	if !bytes.Equal(encodedReports["match"], encodedReports["relocation"]) {
		t.Fatal("relocation fixture changed wire identity")
	}
	encoded, err := json.Marshal(CatalogDiscoveryOptions{Index: CatalogIndexOptions{Root: catalog.RootOptions{Directory: "/private"}}})
	if err != nil || string(encoded) != "{}" {
		t.Fatal("installation options entered request wire")
	}
}
