package apitools

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
)

func TestCatalogUpdateArtifactRefusesMismatchedGETResponseLocator(t *testing.T) {
	const target = "https://example.com/artifact.json"
	client := &Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogUpdateTransport(func(req *http.Request) (*http.Response, error) {
		resp := catalogRepairResponse(req, "fixture")
		if req.Method == http.MethodGet {
			redirected := req.Clone(req.Context())
			redirected.URL, _ = url.Parse(target + "#access_token=fixture-secret")
			resp.Request = redirected
		}
		return resp, nil
	})}}
	var body bytes.Buffer
	got, err := client.acquireCatalogUpdateArtifact(t.Context(), target, DefaultMaxBytes, &body)
	if err == nil || got.URL != target || body.Len() != 0 {
		t.Fatalf("unchecked GET locator = %+v, bytes=%d, err=%v", got, body.Len(), err)
	}
}

func TestCatalogUpdateArtifactDefersBeforeBodyRequest(t *testing.T) {
	for _, tc := range []struct {
		name, length, status string
		code                 int
		limit                int64
	}{
		{name: "unknown", code: 200, status: catalogUpdateSizeUnknown},
		{name: "no HEAD", code: 405, status: catalogUpdateSizeUnknown},
		{name: "too large", length: "1000000001", code: 200, status: catalogUpdateOversized},
		{name: "at download limit", length: "1000000000", code: 200, status: catalogUpdateParserLimit},
		{name: "parser limit", length: "8", code: 200, limit: 7, status: catalogUpdateParserLimit},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var gets atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodHead {
					gets.Add(1)
				}
				if tc.length != "" {
					w.Header().Set("Content-Length", tc.length)
				}
				w.WriteHeader(tc.code)
			}))
			defer server.Close()
			var dst bytes.Buffer
			client := Client{AllowUnsafeHosts: true}
			got, err := client.acquireCatalogUpdateArtifact(t.Context(), server.URL, tc.limit, &dst)
			if err != nil || got.Status != tc.status || gets.Load() != 0 || dst.Len() != 0 {
				t.Fatalf("outcome = %+v, %v, got %d body requests and %d bytes", got, err, gets.Load(), dst.Len())
			}
		})
	}
}

func TestCatalogUpdateArtifactChecksFinalRedirectAndDigest(t *testing.T) {
	const digest = "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/alias" {
			if r.Method != http.MethodHead {
				t.Error("body request reached redirect alias")
			}
			http.Redirect(w, r, "/artifact", http.StatusFound)
			return
		}
		w.Header().Set("Content-Length", "3")
		w.Header().Set("ETag", `"rev1"`)
		if r.Method == http.MethodGet {
			if r.Header.Get("If-Match") != `"rev1"` {
				t.Error("body request missing probed strong ETag")
			}
			fmt.Fprint(w, "abc")
		}
	}))
	defer server.Close()
	var dst bytes.Buffer
	client := Client{AllowUnsafeHosts: true}
	got, err := client.acquireCatalogUpdateArtifact(t.Context(), server.URL+"/alias", 3, &dst)
	if err != nil || got.Status != catalogUpdateDownloaded || got.SHA256 != digest ||
		got.URL != server.URL+"/artifact" || dst.String() != "abc" {
		t.Fatalf("outcome = %+v, %v; body = %q", got, err, dst.String())
	}
}

func TestCatalogUpdateArtifactRejectsChangedBodyResponse(t *testing.T) {
	for _, scenario := range []string{"changed length", "body redirect", "truncated"} {
		t.Run(scenario, func(t *testing.T) {
			var redirected atomic.Bool
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/unprobed" {
					redirected.Store(true)
				}
				w.Header().Set("Content-Length", "3")
				if r.Method == http.MethodHead {
					return
				}
				switch scenario {
				case "changed length":
					w.Header().Set("Content-Length", "4")
					fmt.Fprint(w, "four")
				case "body redirect":
					w.Header().Del("Content-Length")
					w.Header().Set("Location", "/unprobed")
					w.WriteHeader(http.StatusFound)
				case "truncated":
					fmt.Fprint(w, "a")
				}
			}))
			defer server.Close()
			var dst bytes.Buffer
			client := Client{AllowUnsafeHosts: true}
			got, err := client.acquireCatalogUpdateArtifact(t.Context(), server.URL, 3, &dst)
			if got.SHA256 != "" || got.Status == catalogUpdateDownloaded || redirected.Load() {
				t.Fatalf("unverified body admitted: %+v, %v", got, err)
			}
			if scenario == "changed length" && (got.Status != catalogUpdateSizeChanged || dst.Len() != 0) {
				t.Fatalf("changed length was not deferred before staging: %+v, %v", got, err)
			}
			if scenario != "changed length" && err == nil {
				t.Fatalf("expected refusal of %s: %+v", scenario, got)
			}
		})
	}
}

func TestCatalogUpdateArtifactRejectsOutOfScopeHeadRedirect(t *testing.T) {
	var requests atomic.Int32
	foreign := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
	}))
	defer foreign.Close()
	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, foreign.URL+"/artifact.json", http.StatusFound)
	}))
	defer source.Close()
	var dst bytes.Buffer
	client := Client{AllowUnsafeHosts: true}
	if _, err := client.acquireCatalogUpdateArtifact(t.Context(), source.URL, 100, &dst); err == nil ||
		requests.Load() != 0 {
		t.Fatalf("out-of-scope publisher received %d requests, error = %v", requests.Load(), err)
	}
}

func TestCatalogUpdateArtifactRejectsAddedRedirectParameters(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		http.Redirect(w, r, "/artifact?private=fixture", http.StatusFound)
	}))
	defer server.Close()
	var dst bytes.Buffer
	client := Client{AllowUnsafeHosts: true}
	if _, err := client.acquireCatalogUpdateArtifact(t.Context(), server.URL, 100, &dst); err == nil || requests.Load() != 1 {
		t.Fatalf("redirect with added parameters received %d requests, error = %v", requests.Load(), err)
	}
}
