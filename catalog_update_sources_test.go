package apitools

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/OpenUdon/apitools/catalog"
)

func TestCatalogUpdateSourcesKeepOfficialScopedReferences(t *testing.T) {
	provider := catalog.Provider{
		ID: "example",
		SpecReferences: []catalog.SpecReference{
			{ID: "official", Kind: catalog.SpecKindOpenAPI, SourceAuthority: catalog.SourceAuthorityOfficialProvider, URL: "https://example.com/openapi.json", Version: "v2", Revision: "publisher revision"},
			{ID: "community", Kind: catalog.SpecKindOpenAPI, SourceAuthority: catalog.SourceAuthorityPublicCatalog},
			{ID: "docs", Kind: catalog.SpecKindHumanDocs, SourceAuthority: catalog.SourceAuthorityOfficialDocs},
		},
	}
	got := catalogUpdateReferenceSources(provider)
	want := catalogUpdateSource{ProviderID: "example", Service: "example.com/openapi.json", ReferenceID: "official", Kind: catalog.SpecKindOpenAPI, URL: "https://example.com/openapi.json", HeldVersion: "v2", Revision: "publisher revision"}
	if len(got) != 1 || got[0] != want {
		t.Fatalf("publisher sources = %+v; want %+v", got, want)
	}
}

type catalogUpdateTransport func(*http.Request) (*http.Response, error)

func (f catalogUpdateTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestCatalogUpdateSourceResolvesExactPublisherRevision(t *testing.T) {
	const held = "1111111111111111111111111111111111111111"
	const latest = "2222222222222222222222222222222222222222"
	transport := catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host != "api.github.com" || r.URL.Path != "/repos/example/service/commits" ||
			r.URL.Query().Get("path") != "specs/v2/openapi.json" || r.URL.Query().Get("sha") != "" {
			t.Fatalf("unexpected publisher scope %s", r.URL)
		}
		body := fmt.Sprintf(`[{"sha":"%s"}]`, latest)
		return &http.Response{
			StatusCode: http.StatusOK, Header: make(http.Header),
			ContentLength: int64(len(body)), Body: io.NopCloser(strings.NewReader(body)), Request: r,
		}, nil
	})
	client := Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: transport}}
	source := catalogUpdateSource{URL: "https://raw.githubusercontent.com/example/service/" + held + "/specs/v2/openapi.json"}
	got, err := client.resolveCatalogUpdateSource(t.Context(), source)
	if err != nil || got.Revision != latest || got.Resolution != "repository-revision" ||
		got.Scope != "default branch, specs/v2/openapi.json" ||
		got.URL != source.URL ||
		got.ResolvedURL != "https://raw.githubusercontent.com/example/service/"+latest+"/specs/v2/openapi.json" {
		t.Fatalf("resolved spec = %+v, %v", got, err)
	}
}

func TestCatalogUpdateSourcePreservesOriginalBranchAndService(t *testing.T) {
	const prefix = "https://raw.githubusercontent.com/example/service/"
	const file = "/specs/v2/openapi.json"
	revision := "1111111111111111111111111111111111111111"
	transport := catalogUpdateTransport(func(r *http.Request) (*http.Response, error) {
		if got := r.URL.Query().Get("sha"); got != "2.4" {
			t.Fatalf("publisher ref = %q, want 2.4", got)
		}
		body := fmt.Sprintf(`[{"sha":"%s"}]`, revision)
		return &http.Response{StatusCode: 200, Header: make(http.Header), ContentLength: int64(len(body)), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
	})
	client := Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: transport}}
	ref := catalog.SpecReference{ID: "example", URL: prefix + "2.4" + file}
	service, version := catalogUpdateReferenceIdentity(ref)
	source := catalogUpdateSource{URL: ref.URL, Service: service, Version: version}
	first, err := client.resolveCatalogUpdateSource(t.Context(), source)
	if err != nil || first.ResolvedURL != prefix+revision+file {
		t.Fatalf("first resolution = %+v, %v", first, err)
	}
	revision = "2222222222222222222222222222222222222222"
	second, err := client.resolveCatalogUpdateSource(t.Context(), first)
	if err != nil || second.URL != source.URL || second.ResolvedURL != prefix+revision+file {
		t.Fatalf("second resolution = %+v, %v", second, err)
	}
	ref.URL = second.ResolvedURL
	wantService, wantVersion := catalogUpdateReferenceIdentity(ref)
	if wantService != service || wantVersion != version || second.Service != service {
		t.Fatalf("service identity changed with branch lookup: %+v, %q/%q", second, wantService, wantVersion)
	}
}

func TestCatalogUpdateSourceSeparatesAPIVersionFromHeldDialect(t *testing.T) {
	provider := catalog.Provider{ID: "nocodb", SpecReferences: []catalog.SpecReference{
		{ID: "v2", Kind: catalog.SpecKindOpenAPI, URL: "https://nocodb.com/apis/v2/swagger-v2.json", Version: "3.1.0"},
		{ID: "v3", Kind: catalog.SpecKindOpenAPI, URL: "https://nocodb.com/apis/v3/swagger-v3.json", Version: "3.1.0"},
	}}
	got := catalogUpdateReferenceSources(provider)
	if got[0].Service != got[1].Service || got[0].Version != "v2" || got[1].Version != "v3" ||
		got[0].HeldVersion != "3.1.0" || got[0].ReferenceID == got[1].ReferenceID {
		t.Fatalf("API versions and held metadata = %+v", got)
	}
	result, err := (&Client{}).resolveCatalogUpdateSource(t.Context(), got[0])
	if err != nil || result.Resolution != "unexamined-publisher-version" || result.URL != got[0].URL {
		t.Fatalf("unexamined publisher result = %+v, %v", result, err)
	}
}

func TestHubSpotUpdateInventoryLabelsServiceAndReleaseStage(t *testing.T) {
	const indexURL = "https://api.hubapi.com/public/api/spec/v1/specs"
	const index = `{"results":[{"name":"Products","group":"CRM","versions":[
		{"version":"2026-09","stage":"LATEST","openApi":"https://api.hubapi.com/public/api/spec/v2/specs/release/1/version/2026-09","documentationBanner":"PUBLIC_BETA"},
		{"version":"2026-03","stage":"STABLE","openApi":"https://api.hubapi.com/public/api/spec/v2/specs/release/1/version/2026-03"}
	]}]}`
	now := time.Date(2026, 11, 25, 0, 0, 0, 0, time.UTC)
	got, err := hubSpotUpdateInventory(indexURL, now, []byte(index))
	if err != nil || !got.CheckedAt.Equal(now) || got.Scope != indexURL || len(got.Sources) != 2 {
		t.Fatalf("inventory = %+v, %v", got, err)
	}
	if got.Sources[0].Service != "CRM Products" || got.Sources[0].Version != "2026-03" ||
		!got.Sources[0].Stable || got.Sources[1].Stable {
		t.Fatalf("unexpected versions: %+v", got.Sources)
	}
}

func TestHubSpotUpdateInventoryRejectsForeignSpecLink(t *testing.T) {
	const indexURL = "https://api.hubapi.com/public/api/spec/v1/specs"
	const index = `{"results":[{"name":"Products","versions":[
		{"version":"2026-03","stage":"STABLE","openApi":"https://example.com/public/api/spec/v2/specs/release/1/version/2026-03"}
	]}]}`
	if got, err := hubSpotUpdateInventory(indexURL, time.Time{}, []byte(index)); err == nil || len(got.Sources) != 0 {
		t.Fatalf("inventory accepted foreign links: %+v, %v", got, err)
	}
}
