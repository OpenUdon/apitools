package apitools

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAPIVersionComparisonOptionalInventory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(versionSpec("2"))) }))
	defer server.Close()
	req := versionTestRequest(server.URL)
	req.Known.SourceURL = ""
	req.Known.ProviderKey = "fixture"
	req.Locators = []APIVersionLocator{{Kind: "direct", URL: server.URL + "/v2.json"}}
	c := &Client{AllowUnsafeHosts: true, APIsGuruListURL: server.URL + "/catalog"}
	c.APIsGuruListURL = ""
	req.Known.ProviderKey = ""
	req.Known.SourceURL = server.URL + "/v2.json"
	req.Known.Version = "1"
	baseline := OperationInventory{Operations: []OperationSummary{{OperationID: "removed", Method: "POST", Path: "/removed"}}}
	copy := baseline
	report, err := c.DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{Network: true, BaselineInventory: &baseline, BaselineInventorySHA256: req.Known.SHA256})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, row := range report.Versions {
		if row.Validation == "valid" {
			found = true
			if row.Diff == nil || !reflect.DeepEqual(row.Diff.Removed, []string{"POST /removed"}) || !reflect.DeepEqual(row.Diff.Added, []string{"GET /item"}) {
				t.Fatalf("diff %#v", row)
			}
		}
	}
	if !found || !reflect.DeepEqual(copy, baseline) {
		t.Fatal("baseline mutation or missing source")
	}
	report, err = c.DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{Network: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range report.Versions {
		if row.Diff != nil {
			t.Fatal("invented baseline comparison")
		}
	}
}

func TestAPIVersionCandidateRanking(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(versionSpec("2"))) }))
	defer server.Close()
	req := versionTestRequest(server.URL)
	req.Known.Version = "1"
	req.Known.SourceURL = server.URL + "/v2.json"
	req.Contract = &StepContract{Purpose: "get item", Effect: OperationEffectRead}
	report, err := (&Client{AllowUnsafeHosts: true}).DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{Network: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range report.Versions {
		if row.Validation == "valid" && len(row.Candidates) > 0 {
			if row.Candidates[0].Source.SHA256 != row.SHA256 {
				t.Fatal("candidate digest changed")
			}
			return
		}
	}
	t.Fatalf("no ranked source candidate: %#v", report)
}
