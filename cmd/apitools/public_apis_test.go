package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OpenUdon/apitools"
)

func TestSearchPublicAPIsMarkdownCLI(t *testing.T) {
	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/list":
			_, _ = w.Write([]byte("### Communication\n\n| API | Description | Auth | HTTPS | CORS |\n|---|---|---|---|---|\n| [Mail](" + serverURL + "/docs) | Send mail | | Yes | No |\n"))
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
	for _, source := range []string{"public-apis", "auto"} {
		for _, query := range []string{"mail", "zzzznomatch"} {
			var out, errOut bytes.Buffer
			code := runSearchWithClient([]string{"--query", query, "--source", source, "--limit", "1", "--json"}, &out, &errOut, func(string) (*apitools.Client, func(), error) {
				return &apitools.Client{PublicAPIsURL: server.URL + "/list", APIsGuruListURL: server.URL + "/guru", LAPSearchURL: server.URL + "/lap", AllowUnsafeHosts: true}, func() {}, nil
			})
			if code != 0 || errOut.Len() != 0 {
				t.Fatalf("%s %s: code=%d stderr=%s", source, query, code, errOut.String())
			}
			var report apitools.SearchReport
			if err := json.Unmarshal(out.Bytes(), &report); err != nil {
				t.Fatal(err)
			}
			if query == "mail" {
				if len(report.Results) != 1 || !report.Results[0].Validated || report.Results[0].SpecURL != server.URL+"/openapi.json" {
					t.Fatalf("%s: %#v", source, report)
				}
			} else if len(report.Results) != 0 {
				t.Fatalf("no-match %s: %#v", source, report)
			}
		}
	}
}
