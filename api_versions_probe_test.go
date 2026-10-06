package apitools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func versionTestRequest(base string) APIVersionDiscoveryRequest {
	return APIVersionDiscoveryRequest{SchemaVersion: APIVersionDiscoverySchemaVersion, Known: APIVersionBaseline{SourceKind: OperationSourceOpenAPI, Version: "v1", SHA256: strings.Repeat("a", 64), SourceURL: base + "/v1.json"}, OfficialScopes: []APIVersionOfficialScope{{Origin: base, PathPrefix: "/"}}}
}
func TestAPIVersionRequestWireAndNaturalOrdering(t *testing.T) {
	req := versionTestRequest("https://example.com")
	data, _ := json.Marshal(req)
	if _, err := DecodeAPIVersionDiscoveryRequest(data); err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{`{"schema_version":"old"}`, string(data) + ` {}`, strings.Replace(string(data), `"known":`, `"network":true,"known":`, 1)} {
		if _, err := DecodeAPIVersionDiscoveryRequest([]byte(bad)); err == nil {
			t.Fatalf("accepted %s", bad)
		}
	}
	for _, pair := range [][2]string{{"v10", "v2"}, {"1.10", "1.2"}} {
		if n, ok := compareAPIVersions(pair[0], pair[1]); !ok || n <= 0 {
			t.Fatalf("comparison %v", pair)
		}
	}
	for _, v := range []string{"2026-10-06", "2026.10.06", "latest", "1.0-beta"} {
		if _, ok := nextAPIVersion(v); ok {
			t.Fatalf("guessed %s", v)
		}
	}
}
func TestAPIVersionProbeIsolationAndScope(t *testing.T) {
	var cookies, callbacks, requests atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		if r.Header.Get("Cookie") != "" {
			cookies.Add(1)
		}
		if r.Header.Get("User-Agent") != "apitools/api-version-discovery-v1" {
			t.Error("missing identifier")
		}
		switch r.URL.Path {
		case "/v1.json":
			http.Redirect(w, r, "/v2.json", http.StatusFound)
		case "/v2.json":
			w.Header().Set("ETag", `"next"`)
			_, _ = w.Write([]byte(`{"openapi":"3.0.0"}`))
		case "/escape":
			http.Redirect(w, r, "https://example.net/spec.json", http.StatusFound)
		}
	}))
	defer server.Close()
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(server.URL)
	jar.SetCookies(u, []*http.Cookie{{Name: "secret", Value: "not-sent"}})
	c := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Jar: jar, CheckRedirect: func(*http.Request, []*http.Request) error { callbacks.Add(1); return nil }}}
	b := &apiVersionBudget{sem: make(chan struct{}, 4)}
	scopes := []APIVersionOfficialScope{{Origin: server.URL, PathPrefix: "/"}}
	p := c.probeAPIVersion(context.Background(), b, server.URL+"/v1.json", scopes, true, "", "")
	if p.err != nil || p.digest == "" || p.etag != `"next"` || cookies.Load() != 0 || callbacks.Load() != 0 {
		t.Fatalf("probe=%#v cookies=%d callbacks=%d", p, cookies.Load(), callbacks.Load())
	}
	if count, _ := b.counts(); count != 2 || requests.Load() != 2 {
		t.Fatalf("redirect accounting %d", count)
	}
	p = c.probeAPIVersion(context.Background(), b, server.URL+"/escape", scopes, true, "", "")
	if p.err == nil {
		t.Fatal("off-origin redirect accepted")
	}
	if c.HTTPClient.Jar != jar || c.HTTPClient.CheckRedirect == nil {
		t.Fatal("caller client mutated")
	}
	for _, raw := range []string{server.URL + "/spec?token=secret", server.URL + "/%252e%252e/spec"} {
		if _, e := c.versionURLSyntax(raw); e == nil {
			t.Fatalf("accepted %s", raw)
		}
	}
	if (&Client{}).versionScopeAllows("http://127.0.0.1/spec", scopes) == nil {
		t.Fatal("unsafe host accepted")
	}
}
func TestAPIVersionProbeConditionalSniffAndBounds(t *testing.T) {
	var readBytes atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/unchanged":
			if r.Header.Get("If-None-Match") != `"same"` {
				t.Error("missing conditional")
			}
			w.WriteHeader(304)
		case "/html":
			for i := 0; i < 1024; i++ {
				n, err := w.Write([]byte("<html>" + strings.Repeat("x", 4090)))
				readBytes.Add(int64(n))
				if err != nil {
					return
				}
			}
		case "/late":
			_, _ = w.Write([]byte(`{"description":"` + strings.Repeat("x", 5000) + `","openapi":"3.0.0"}`))
		case "/big":
			w.Header().Set("Content-Length", "9000000")
		case "/slow":
			<-r.Context().Done()
		}
	}))
	defer server.Close()
	c := &Client{AllowUnsafeHosts: true}
	scopes := []APIVersionOfficialScope{{Origin: server.URL, PathPrefix: "/"}}
	b := &apiVersionBudget{sem: make(chan struct{}, 4)}
	p := c.probeAPIVersion(context.Background(), b, server.URL+"/unchanged", scopes, true, `"same"`, "")
	if p.err != nil || p.status != 304 {
		t.Fatalf("304 %#v", p)
	}
	if count, bodies := b.counts(); count != 1 || bodies != 0 {
		t.Fatal("304 accounting")
	}
	p = c.probeAPIVersion(context.Background(), b, server.URL+"/html", scopes, true, "", "")
	if !p.nonSpec || p.bytes > 4096 || p.digest != "" {
		t.Fatalf("HTML %#v", p)
	}
	p = c.probeAPIVersion(context.Background(), b, server.URL+"/late", scopes, true, "", "")
	if p.err != nil || p.digest == "" {
		t.Fatalf("late marker %#v", p)
	}
	p = c.probeAPIVersion(context.Background(), b, server.URL+"/big", scopes, true, "", "")
	if p.err == nil || p.digest != "" || p.declared == nil {
		t.Fatalf("oversize %#v", p)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()
	start := time.Now()
	p = c.probeAPIVersion(ctx, &apiVersionBudget{sem: make(chan struct{}, 4)}, server.URL+"/slow", scopes, true, "", "")
	if p.err == nil || time.Since(start) > 300*time.Millisecond {
		t.Fatalf("timeout %#v", p)
	}
}
