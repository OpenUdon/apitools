package apitools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func versionSpec(v string) string {
	return fmt.Sprintf(`{"openapi":"3.0.0","info":{"title":"Fixture","version":%q},"paths":{"/item":{"get":{"operationId":"getItem","responses":{"200":{"description":"OK"}}}}}}`, v)
}
func TestAPIVersionDiscoverySiblingsSaveAndResume(t *testing.T) {
	var requests atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		switch r.URL.Path {
		case "/v1.json":
			_, _ = w.Write([]byte(versionSpec("1")))
		case "/v2.json":
			_, _ = w.Write([]byte(versionSpec("2")))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	req := versionTestRequest(server.URL)
	req.Locators = []APIVersionLocator{{Kind: "pattern", URLTemplate: server.URL + "/v{version}.json", BaselineToken: "1"}}
	c := &Client{AllowUnsafeHosts: true}
	dir := t.TempDir()
	options := APIVersionDiscoveryOptions{Network: true, StatePath: filepath.Join(dir, "state.json"), SaveDir: filepath.Join(dir, "sources")}
	report, err := c.DiscoverAPIVersions(context.Background(), req, options)
	if err != nil || report.Status != "newer_found" || report.Preferred != "2" {
		t.Fatalf("report=%#v err=%v", report, err)
	}
	var saved *ImportedSpec
	for _, v := range report.Versions {
		if v.Comparison == "newer" {
			saved = v.Saved
		}
	}
	if saved == nil {
		t.Fatalf("no saved newer source: %#v", report)
	}
	data, e := os.ReadFile(saved.Path)
	if e != nil || string(data) != versionSpec("2") {
		t.Fatalf("save=%#v err=%v", saved, e)
	}
	before := requests.Load()
	resumed, err := c.DiscoverAPIVersions(context.Background(), req, options)
	if err != nil || resumed.Status != report.Status || requests.Load() != before {
		t.Fatalf("resume requests=%d/%d err=%v", requests.Load(), before, err)
	}
	if report.Known != req.Known {
		t.Fatal("baseline changed")
	}
}

func TestAPIVersionDiscoveryDeadlineAndDeny(t *testing.T) {
	var count atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { count.Add(1); <-r.Context().Done() }))
	defer server.Close()
	req := versionTestRequest(server.URL)
	c := &Client{AllowUnsafeHosts: true}
	report, err := c.DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{})
	if err != nil || count.Load() != 0 || report.Known != req.Known || report.Status != "unexamined" {
		t.Fatalf("deny %#v %v", report, err)
	}
	start := time.Now()
	report, err = c.DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{Network: true, Timeout: 30 * time.Millisecond})
	if err != nil || report.Known != req.Known || report.Status != "unexamined" || time.Since(start) > 280*time.Millisecond {
		t.Fatalf("deadline %#v %v", report, err)
	}
}

func TestAPIVersionDiscoveryConditionalIsNotAbsence(t *testing.T) {
	var count atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count.Add(1)
		if r.Header.Get("If-None-Match") != `"known"` {
			t.Error("missing conditional")
		}
		w.WriteHeader(304)
	}))
	defer server.Close()
	req := versionTestRequest(server.URL)
	req.Known.ETag = `"known"`
	report, err := (&Client{AllowUnsafeHosts: true}).DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{Network: true})
	if err != nil || count.Load() != 1 || report.Status != "unexamined" {
		t.Fatalf("304 %#v %v count=%d", report, err, count.Load())
	}
}

func TestAPIVersionDiscoveryPointerAndInvalidSource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/docs":
			_, _ = w.Write([]byte(`<a href="/v10.json">New</a><a href="/v2.json">Old</a>`))
		case "/v10.json":
			_, _ = w.Write([]byte(versionSpec("10")))
		case "/v2.json":
			_, _ = w.Write([]byte(`{"openapi":"3.0.0","info":{"version":"2"},"paths":{}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	req := versionTestRequest(server.URL)
	req.DocsURL = server.URL + "/docs"
	report, err := (&Client{AllowUnsafeHosts: true}).DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{Network: true, SaveDir: t.TempDir()})
	if err != nil || report.Status != "newer_found" {
		t.Fatalf("report %#v %v", report, err)
	}
	found := false
	for _, row := range report.Versions {
		if strings.HasSuffix(row.SourceURL, "v2.json") {
			found = true
			if row.Validation != "invalid" || row.SHA256 == "" || row.Saved != nil {
				t.Fatalf("invalid %#v", row)
			}
		}
	}
	if !found {
		t.Fatal("invalid source disappeared")
	}
}

func TestAPIVersionDiscoveryDirectory304(t *testing.T) {
	var requests atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		if r.Header.Get("If-None-Match") != `"list"` {
			t.Error("missing list validator")
		}
		w.WriteHeader(304)
	}))
	defer server.Close()
	req := versionTestRequest("https://example.com")
	req.Known.SourceURL = ""
	req.Known.ProviderKey = "example.com:mail"
	cache := apiVersionListCache{SchemaVersion: "apitools.api-version-list-cache/v1", URL: server.URL + "/list.json", ETag: `"list"`, CheckedAt: time.Now().Add(-25 * time.Hour), Data: json.RawMessage(`{}`)}
	data, _ := json.Marshal(cache)
	file := filepath.Join(t.TempDir(), "list.json")
	if err := os.WriteFile(file, data, 0600); err != nil {
		t.Fatal(err)
	}
	report, err := (&Client{AllowUnsafeHosts: true, APIsGuruListURL: cache.URL}).DiscoverAPIVersions(context.Background(), req, APIVersionDiscoveryOptions{Network: true, ListCachePath: file})
	if err != nil || requests.Load() != 1 || report.Known != req.Known {
		t.Fatalf("304 directory %#v %v requests=%d", report, err, requests.Load())
	}
}
