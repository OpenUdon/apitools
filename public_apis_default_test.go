package apitools

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type publicAPIsFixtureTransport struct {
	base    http.RoundTripper
	listURL string
}

func (t publicAPIsFixtureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.String() == DefaultPublicAPIsURL {
		next, err := http.NewRequestWithContext(req.Context(), req.Method, t.listURL, nil)
		if err != nil {
			return nil, err
		}
		next.Header = req.Header.Clone()
		return t.base.RoundTrip(next)
	}
	return t.base.RoundTrip(req)
}

func TestPublicAPIsDefaultAndAutoMarkdown(t *testing.T) {
	if DefaultPublicAPIsURL != "https://raw.githubusercontent.com/public-apis/public-apis/master/README.md" || strings.Contains(DefaultPublicAPIsURL, "api.publicapis.org") {
		t.Fatalf("unexpected default: %s", DefaultPublicAPIsURL)
	}
	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/list":
			_, _ = w.Write([]byte(publicAPIsTestHeader + "| [Mail](" + serverURL + "/docs) | Send mail | | Yes | No |\n"))
		case "/guru":
			_, _ = w.Write([]byte(`{}`))
		case "/lap":
			w.Header().Set("Content-Type", "text/lap")
			_, _ = w.Write([]byte("@search_query mail\n@total 0\n"))
		case "/openapi.json":
			_, _ = w.Write([]byte(`{"openapi":"3.0.0","info":{"title":"Mail","version":"1.0"},"paths":{}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	serverURL = server.URL
	client := &Client{
		HTTPClient:       &http.Client{Transport: publicAPIsFixtureTransport{base: server.Client().Transport, listURL: server.URL + "/list"}},
		APIsGuruListURL:  server.URL + "/guru",
		LAPSearchURL:     server.URL + "/lap",
		AllowUnsafeHosts: true,
	}
	for _, source := range []Source{SourcePublicAPIs, SourceAuto} {
		for _, query := range []string{"mail", "zzzznomatch"} {
			report, err := client.Search(context.Background(), SearchOptions{Query: query, Source: source, Limit: 1})
			if err != nil {
				t.Fatalf("%s %s: %v", source, query, err)
			}
			if query == "mail" {
				if len(report.Results) != 1 || report.Results[0].Source != string(SourcePublicAPIs) || !report.Results[0].Validated {
					t.Fatalf("%s: %#v", source, report)
				}
			} else if len(report.Results) != 0 {
				t.Fatalf("no-match %s: %#v", source, report)
			}
		}
	}
}
