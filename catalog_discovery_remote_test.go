package apitools_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/OpenUdon/apitools"
)

func TestCatalogDiscoveryRemoteOptInAndProvenance(t *testing.T) {
	content, err := os.ReadFile("testdata/catalog-root/openapi/notes.json")
	if err != nil {
		t.Fatal(err)
	}
	var calls, sources, cookies atomic.Int32
	var server *httptest.Server
	server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Header.Get("Cookie") != "" || r.Header.Get("Authorization") != "" {
			cookies.Add(1)
		}
		if r.URL.Path == "/list" {
			entries := map[string]any{}
			for _, id := range []string{"a", "b", "c", "d"} {
				entries[id] = map[string]any{"preferred": "1", "versions": map[string]any{"1": map[string]any{"swaggerUrl": server.URL + "/" + id + "?token=fixture-private-value", "info": map[string]any{"title": "List notes records"}}}}
			}
			_ = json.NewEncoder(w).Encode(entries)
		} else {
			sources.Add(1)
			_, _ = w.Write(content)
		}
	}))
	defer server.Close()
	options := preparedCatalogDiscovery(t)
	jar, _ := cookiejar.New(nil)
	address, _ := url.Parse(server.URL)
	jar.SetCookies(address, []*http.Cookie{{Name: "session", Value: "fixture-private-cookie"}})
	httpClient := server.Client()
	httpClient.Jar = jar
	options.RemoteClient = &apitools.Client{APIsGuruListURL: server.URL + "/list", AllowUnsafeHosts: true, HTTPClient: httpClient}
	for _, flags := range [][2]bool{{false, false}, {false, true}, {true, false}} {
		options.Request.RemoteLookup, options.RemoteEnabled = flags[0], flags[1]
		report, err := apitools.DiscoverCatalogOperations(context.Background(), options)
		if err != nil {
			t.Fatal(err)
		}
		if flags[0] && report.Outcome != apitools.CatalogDiscoveryBlocked {
			t.Fatal("unconfigured remote lookup was not refused")
		}
	}
	if calls.Load() != 0 {
		t.Fatal("remote lookup ran without both opt-ins")
	}
	options.Request.RemoteLookup, options.RemoteEnabled = true, true
	report, err := apitools.DiscoverCatalogOperations(context.Background(), options)
	if err != nil || report.Outcome != apitools.CatalogDiscoveryMatch || report.QualifiedOperations != 1 || report.ExaminedOperations != 2 || !report.Incomplete || report.Scope.Complete {
		t.Fatalf("remote dedup: %v %s count=%d examined=%d", err, report.Outcome, report.QualifiedOperations, report.ExaminedOperations)
	}
	if calls.Load() != 4 || sources.Load() != 3 || cookies.Load() != 0 {
		t.Fatalf("fetch/credential boundary: %d calls %d sources %d credential requests", calls.Load(), sources.Load(), cookies.Load())
	}
	digest := sha256.Sum256(content)
	if report.Candidates[0].Remote == nil || report.Candidates[0].Remote.SHA256 != hex.EncodeToString(digest[:]) || report.Candidates[0].Remote.Bytes != int64(len(content)) || report.Candidates[0].Remote.Authority != "public-catalog" {
		t.Fatal("fetched identity lost")
	}
	fetched := 0
	for _, lead := range report.Leads {
		if lead.Remote != nil {
			fetched++
		}
	}
	if fetched != 3 {
		t.Fatal("individual fetched-source provenance lost")
	}
	encoded, _ := json.Marshal(report)
	if bytes.Contains(encoded, []byte("fixture-private")) || bytes.Contains(encoded, []byte("?token=")) {
		t.Fatal("credential-like URL or cookies entered evidence")
	}
	options.Index.Root.Directory = ""
	remoteOnly, err := apitools.DiscoverCatalogOperations(context.Background(), options)
	if err != nil || remoteOnly.Outcome != apitools.CatalogDiscoveryMatch || remoteOnly.Candidates[0].Remote == nil || len(remoteOnly.Candidates[0].References) != 0 || remoteOnly.Scope.Complete {
		t.Fatal("ephemeral remote identity pretended to be registered local evidence")
	}
	for i := range options.Index.Catalog.Providers {
		options.Index.Catalog.Providers[i].SpecReferences[0].URL = server.URL + "/a?token=fixture-private-value"
	}
	options.Request.ProviderKeys = []string{"Example Notes", "Example Notes Shared"}
	constrained, err := apitools.DiscoverCatalogOperations(context.Background(), options)
	if err != nil || constrained.Outcome != apitools.CatalogDiscoveryMatch || len(constrained.Candidates[0].Sources) != 2 || len(constrained.Scope.ProviderIDs) != 2 {
		t.Fatalf("shared exact remote URL lost canonical provider links: err=%v outcome=%s scope=%+v candidates=%+v exclusions=%+v diagnostics=%+v", err, constrained.Outcome, constrained.Scope, constrained.Candidates, constrained.Exclusions, constrained.Diagnostics)
	}

}

func TestCatalogDiscoveryRemoteFailureScopeAndSafety(t *testing.T) {
	for _, name := range []string{"empty", "invalid-spec", "oversize", "timeout", "cancel", "unsafe-list", "unsafe-source", "userinfo", "redirect", "constrained", "filtered", "transport-timeout"} {
		t.Run(name, func(t *testing.T) {
			var sourceCalls atomic.Int32
			var server *httptest.Server
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/list" {
					if name == "timeout" || name == "cancel" {
						<-r.Context().Done()
						return
					}
					if name == "empty" {
						_, _ = w.Write([]byte("{}"))
						return
					}
					source := server.URL + "/source"
					if name == "unsafe-source" {
						source = "file:///private/spec.json"
					}
					if name == "userinfo" {
						source = "http://user:pass@localhost/source"
					}
					_ = json.NewEncoder(w).Encode(map[string]any{"notes.example": map[string]any{"preferred": "1", "versions": map[string]any{"1": map[string]any{"swaggerUrl": source, "info": map[string]any{"title": "List notes records"}}}}})
					return
				}
				sourceCalls.Add(1)
				if name == "transport-timeout" {
					<-r.Context().Done()
					return
				}
				if name == "redirect" {
					http.Redirect(w, r, "file:///private/spec.json", 302)
					return
				}
				if name == "oversize" {
					_, _ = w.Write(bytes.Repeat([]byte("x"), 2048))
					return
				}
				_, _ = w.Write([]byte("not an API document"))
			}))
			defer server.Close()
			options := preparedCatalogDiscovery(t)
			options.Index.Root.Directory = "" // no local operations can mask failure evidence
			options.Request.RemoteLookup, options.RemoteEnabled = true, true
			options.RemoteClient = &apitools.Client{APIsGuruListURL: server.URL + "/list", AllowUnsafeHosts: true}
			want := apitools.CatalogDiscoveryInsufficientEvidence
			switch name {
			case "unsafe-list":
				options.RemoteClient.AllowUnsafeHosts = false
				want = apitools.CatalogDiscoveryBlocked
			case "unsafe-source", "userinfo", "redirect":
				want = apitools.CatalogDiscoveryBlocked
			case "oversize":
				options.RemoteClient.MaxBytes = 1024
			case "transport-timeout":
				options.RemoteClient.HTTPClient = &http.Client{Timeout: 20 * time.Millisecond}
			case "timeout":
				options.Request.Limits.TimeoutMillis = 20
			case "constrained":
				options.Request.ProviderKeys = []string{"Example Notes"}
			case "filtered":
				options.Request.Filters.RequireLicenseEvidence = true
			}
			ctx := context.Background()
			if name == "cancel" {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, 20*time.Millisecond)
				defer cancel()
			}
			report, err := apitools.DiscoverCatalogOperations(ctx, options)
			if name == "cancel" {
				if !errors.Is(err, context.DeadlineExceeded) {
					t.Fatalf("caller cancellation: %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			if report.Outcome != want || report.Scope.Complete || report.QualifiedOperations != 0 {
				t.Fatalf("failure/negative safety: %s %+v", report.Outcome, report)
			}
			if (name == "constrained" || name == "filtered") && (sourceCalls.Load() != 0 || len(report.Exclusions) == 0) {
				t.Fatal("remote lookup broadened constrained scope or bypassed filters")
			}
			if name == "transport-timeout" {
				found := false
				for _, diagnostic := range report.Diagnostics {
					found = found || diagnostic.Code == "discovery.remote_timeout"
				}
				if !found {
					t.Fatal("wrapped transport timeout lost explicit timeout evidence")
				}
			}
			if name == "invalid-spec" && (len(report.Leads) == 0 || report.Leads[0].Remote == nil) {
				t.Fatal("unparsed source identity lost")
			}
		})
	}
}
