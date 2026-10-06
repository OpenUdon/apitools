package apitools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/OpenUdon/apitools/catalog"
)

func TestAPIVersionTwoMissScopeAndFailures(t *testing.T) {
	var count atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count.Add(1)
		switch r.URL.Path {
		case "/v1.json":
			_, _ = w.Write([]byte(versionSpec("1")))
		case "/throttle":
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(429)
		case "/bad-docs":
			http.NotFound(w, r)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	req := versionTestRequest(server.URL)
	req.Locators = []APIVersionLocator{{Kind: "pattern", URLTemplate: server.URL + "/v{version}.json", BaselineToken: "1"}}
	c := &Client{AllowUnsafeHosts: true}
	report, err := c.DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{Network: true})
	if err != nil || count.Load() != 3 || report.Status != "none_found_in_scope" {
		t.Fatalf("two misses %#v count=%d err=%v", report, count.Load(), err)
	}
	req.Locators = nil
	req.Known.SourceURL = server.URL + "/throttle"
	before := count.Load()
	report, err = c.DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{Network: true})
	if err != nil || count.Load() != before+1 || report.Status != "unexamined" {
		t.Fatalf("429 retried %#v %v", report, err)
	}
	req.Known.SourceURL = server.URL + "/v1.json"
	req.DocsURL = server.URL + "/bad-docs"
	report, err = c.DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{Network: true})
	if err != nil || report.Status != "unexamined" {
		t.Fatalf("404 docs invented absence %#v %v", report, err)
	}
}

func TestAPIVersionNativeSourceAndUnsafeSave(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/google-v2.json" {
			_, _ = w.Write([]byte(`{"kind":"discovery#restDescription","discoveryVersion":"v1","id":"fixture:v2","name":"fixture","version":"v2","rootUrl":"https://example.com/","servicePath":"v2/","resources":{},"schemas":{}}`))
		} else {
			_, _ = w.Write([]byte(versionSpec("2")))
		}
	}))
	defer server.Close()
	req := versionTestRequest(server.URL)
	req.Known.SourceKind = OperationSourceGoogleDiscovery
	req.Known.SourceURL = server.URL + "/google-v2.json"
	dir := t.TempDir()
	report, err := (&Client{AllowUnsafeHosts: true}).DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{Network: true, SaveDir: dir})
	if err != nil || report.Status != "newer_found" {
		t.Fatalf("Discovery %#v %v", report, err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Fatal("Discovery masqueraded as Import-valid source")
	}
	req.Known.SourceKind = OperationSourceOpenAPI
	req.Known.SourceURL = server.URL + "/v2.json"
	target := t.TempDir()
	link := filepath.Join(dir, "link")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	report, err = (&Client{AllowUnsafeHosts: true}).DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{Network: true, SaveDir: link})
	if err != nil || report.Status != "newer_found" {
		t.Fatalf("save failure lost check %#v %v", report, err)
	}
	for _, row := range report.Versions {
		if row.Saved != nil {
			t.Fatal("followed unsafe save root")
		}
	}
}

func TestAPIVersionOrderAndSharedRepositoryBoundary(t *testing.T) {
	var reverse atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2.json" && reverse.Load() {
			time.Sleep(10 * time.Millisecond)
		}
		if r.URL.Path == "/v3.json" && !reverse.Load() {
			time.Sleep(10 * time.Millisecond)
		}
		v := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/v"), ".json")
		_, _ = w.Write([]byte(versionSpec(v)))
	}))
	defer server.Close()
	req := versionTestRequest(server.URL)
	req.Locators = []APIVersionLocator{{Kind: "direct", URL: server.URL + "/v2.json"}, {Kind: "direct", URL: server.URL + "/v3.json"}}
	c := &Client{AllowUnsafeHosts: true}
	first, err := c.DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{Network: true})
	if err != nil {
		t.Fatal(err)
	}
	reverse.Store(true)
	second, err := c.DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{Network: true})
	if err != nil {
		t.Fatal(err)
	}
	ids := func(r APIVersionDiscoveryReport) []string {
		var out []string
		for _, v := range r.Versions {
			out = append(out, v.ID)
		}
		return out
	}
	if !reflect.DeepEqual(ids(first), ids(second)) {
		t.Fatal("report order follows completion order")
	}
	scopes := []APIVersionOfficialScope{{Origin: "https://raw.githubusercontent.com", Repository: "owner/repo", PathPrefix: "/owner/repo/"}}
	if c.versionScopeAllows("https://raw.githubusercontent.com/other/repo/main/api.json", scopes) == nil {
		t.Fatal("cross-repository source accepted")
	}
	for _, raw := range []string{"https://example.com/spec?api-key=secret", "https://example.com/spec?token=%QQ"} {
		if _, err := c.versionURLSyntax(raw); err == nil {
			t.Fatalf("credential query accepted %s", raw)
		}
	}
}

func TestAPIVersionStateCannotCrossQueryOrScope(t *testing.T) {
	var count atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { count.Add(1); http.NotFound(w, r) }))
	defer server.Close()
	req := versionTestRequest(server.URL)
	req.Locators = []APIVersionLocator{{Kind: "pattern", URLTemplate: server.URL + "/v{version}.json", BaselineToken: "1"}}
	opts := APIVersionDiscoveryOptions{Network: true, StatePath: filepath.Join(t.TempDir(), "state.json")}
	c := &Client{AllowUnsafeHosts: true}
	first, err := c.DiscoverAPIVersions(context.Background(), req, opts)
	if err != nil {
		t.Fatal(err)
	}
	before := count.Load()
	req.Known.Version = "v2"
	_, err = c.DiscoverAPIVersions(context.Background(), req, opts)
	if err != nil || count.Load() == before {
		t.Fatal("reused foreign query")
	}
	data, _ := json.Marshal(first)
	if strings.Contains(string(data), "Content") || strings.Contains(string(data), "Cookie") {
		t.Fatal("runtime source bytes or cookies persisted")
	}
}

func TestAPIVersionServiceIdentityAndURLOnlyDirectory(t *testing.T) {
	req := versionTestRequest("https://example.com")
	req.Known.SourceURL = "https://example.com/crm/v1.json"
	req.Known.ProviderKey = "provider"
	cat := catalog.Catalog{Providers: []catalog.Provider{{ID: "provider", SpecReferences: []catalog.SpecReference{{URL: "https://example.com/crm/v2.json", SourceAuthority: catalog.SourceAuthorityOfficialProvider}, {URL: "https://example.com/books/v9.json", SourceAuthority: catalog.SourceAuthorityOfficialProvider}}}}}
	leads := versionCatalogReferences(req, APIVersionDiscoveryOptions{Catalog: &cat})
	if len(leads) != 1 || leads[0].SourceURL != "https://example.com/crm/v2.json" {
		t.Fatalf("cross-service leads %#v", leads)
	}
	req.Known.ProviderKey = ""
	data := []byte(`{"example.com:crm":{"preferred":"v2","versions":{"v1":{"swaggerUrl":"https://example.com/crm/v1.json","info":{}},"v2":{"swaggerUrl":"https://example.com/crm/v2.json","info":{}}}}}`)
	rows, err := versionListRecords(context.Background(), data, req)
	if err != nil || len(rows) != 2 {
		t.Fatalf("URL-only family lost newer versions %#v %v", rows, err)
	}
	if versionTokenURL("https://www.googleapis.com/discovery/v1/apis/calendar/v3/rest") != "v3" {
		t.Fatal("Discovery protocol version confused with service version")
	}
}

func TestAPIVersionLatestPointerAndInvalidContract(t *testing.T) {
	var sourceReads atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/latest":
			http.Redirect(w, r, "/v2.json", http.StatusFound)
		case "/v2.json":
			sourceReads.Add(1)
			_, _ = w.Write([]byte(versionSpec("2")))
		default:
			_, _ = w.Write([]byte(versionSpec("1")))
		}
	}))
	defer server.Close()
	req := versionTestRequest(server.URL)
	req.Locators = []APIVersionLocator{{Kind: "pointer", URL: server.URL + "/latest"}}
	c := &Client{AllowUnsafeHosts: true}
	report, err := c.DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{Network: true})
	if err != nil || report.Status != "newer_found" || sourceReads.Load() != 1 {
		t.Fatalf("alias %#v %v count=%d", report, err, sourceReads.Load())
	}
	req.Contract = &StepContract{Purpose: "get item", Effect: OperationEffect("invalid")}
	before := sourceReads.Load()
	if _, err = c.DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{Network: true, SaveDir: t.TempDir()}); err == nil || sourceReads.Load() != before {
		t.Fatal("invalid contract reached network")
	}
}

func TestAPIVersionGitHubTreeAndRecipe(t *testing.T) {
	var trees atomic.Int64
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		body := ""
		status := 200
		switch req.URL.Hostname() {
		case "api.github.com":
			trees.Add(1)
			body = `{"sha":"pinned","truncated":false,"tree":[{"path":"api/v2.json","type":"blob"}]}`
		case "raw.githubusercontent.com":
			body = versionSpec("2")
		default:
			status = 404
		}
		return &http.Response{StatusCode: status, Status: fmt.Sprint(status), Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: req, ContentLength: int64(len(body))}, nil
	})}}
	req := versionTestRequest("https://raw.githubusercontent.com")
	req.Known.SourceURL = "https://raw.githubusercontent.com/owner/repo/main/api/v1.json"
	req.OfficialScopes = []APIVersionOfficialScope{{Origin: "https://raw.githubusercontent.com", Repository: "owner/repo", PathPrefix: "/owner/repo/"}}
	req.Locators = []APIVersionLocator{{Kind: "github", Repository: "owner/repo", Ref: "main", PathPattern: "api/*.json"}}
	report, err := client.DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{Network: true})
	if err != nil || trees.Load() != 1 {
		t.Fatalf("tree %#v %v", report, err)
	}
	found := false
	for _, r := range report.Versions {
		if strings.Contains(r.SourceURL, "/pinned/") {
			found = true
			if r.Recipe.Kind != "github" {
				t.Fatal("recipe lost")
			}
		}
	}
	if !found {
		t.Fatal("tree source lost")
	}
}

func TestAPIVersionSavedDigestReuseAndCollision(t *testing.T) {
	data := []byte(versionSpec("2"))
	row := APIVersionRecord{FinalURL: "https://example.com/v2.json"}
	dir := t.TempDir()
	first, err := saveVersionSource(context.Background(), dir, row, data)
	if err != nil {
		t.Fatal(err)
	}
	second, err := saveVersionSource(context.Background(), dir, row, data)
	if err != nil || first.Path != second.Path || first.SHA256 != second.SHA256 {
		t.Fatal("identical source not reused")
	}
	third, err := saveVersionSource(context.Background(), dir, row, []byte(strings.Replace(versionSpec("2"), "Fixture", "Changed", 1)))
	if err != nil || third.Path == first.Path {
		t.Fatal("differing source overwritten")
	}
	content, err := os.ReadFile(first.Path)
	if err != nil || string(content) != string(data) {
		t.Fatal("old source damaged")
	}
}

func TestAPIVersionGitHubOriginRequiresExactHost(t *testing.T) {
	var trees atomic.Int64
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Hostname() == "api.github.com" {
			trees.Add(1)
		}
		return &http.Response{StatusCode: 404, Status: "404", Header: make(http.Header), Body: io.NopCloser(strings.NewReader("")), Request: req}, nil
	})}}
	req := versionTestRequest("https://raw.githubusercontent.com.example.com")
	req.Known.SourceURL = "https://raw.githubusercontent.com.example.com/owner/repo/v1.json"
	req.OfficialScopes = []APIVersionOfficialScope{{Origin: "https://raw.githubusercontent.com.example.com", Repository: "owner/repo", PathPrefix: "/owner/repo/"}}
	req.Locators = []APIVersionLocator{{Kind: "github", Repository: "owner/repo", Ref: "main", PathPattern: "*.json"}}
	report, err := client.DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{Network: true})
	if err != nil || trees.Load() != 0 || report.Status != "unexamined" {
		t.Fatalf("hostname substring granted GitHub authority: %#v %v", report, err)
	}
}

func TestAPIVersionNativeDiscoveryInventoryDiff(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"kind":"discovery#restDescription","discoveryVersion":"v1","id":"fixture:v2","name":"fixture","version":"v2","rootUrl":"https://example.com/","servicePath":"v2/","methods":{"getItem":{"id":"fixture.get","path":"items","httpMethod":"GET"}},"schemas":{}}`))
	}))
	defer server.Close()
	req := versionTestRequest(server.URL)
	req.Known.SourceKind = OperationSourceGoogleDiscovery
	req.Known.SourceURL = server.URL + "/v2.json"
	baseline := OperationInventory{Operations: []OperationSummary{{OperationID: "fixture.old", Method: "GET", Path: "old", Extensions: map[string]string{"x-uws-source-kind": "google-discovery"}}}}
	report, err := (&Client{AllowUnsafeHosts: true}).DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{Network: true, BaselineInventory: &baseline, BaselineInventorySHA256: req.Known.SHA256})
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range report.Versions {
		if row.Validation == "valid" {
			if row.Diff == nil || !reflect.DeepEqual(row.Diff.Added, []string{"fixture.get"}) || !reflect.DeepEqual(row.Diff.Removed, []string{"fixture.old"}) {
				t.Fatalf("native diff lost: %#v", row.Diff)
			}
			return
		}
	}
	t.Fatal("no native source")
}
