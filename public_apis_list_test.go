package apitools

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
)

const publicAPIsTestHeader = "### Communication\n\n| API | Description | Auth | HTTPS | CORS |\n|:---|:---|:---|:---|:---|\n"

func TestPublicAPIsListFormats(t *testing.T) {
	want := []publicAPIEntry{{API: "Mail", Description: "Send mail", Link: "https://example.com/docs", Category: "Communication"}}
	for _, body := range []string{
		publicAPIsTestHeader + "| [Mail](https://example.com/docs) | Send mail | apiKey | Yes | Unknown |\n",
		`{"entries":[{"API":"Mail","Description":"Send mail","Link":"https://example.com/docs","Category":"Communication","Auth":"apiKey"}]}`,
	} {
		got, err := parsePublicAPIsCatalog(context.Background(), []byte(body), "https://example.com/list")
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("entries=%#v, err=%v", got, err)
		}
	}
	empty, err := parsePublicAPIsCatalog(context.Background(), []byte(`{"entries":[]}`), "https://example.com/list")
	if err != nil || len(empty) != 0 {
		t.Fatalf("legacy empty list=%#v, err=%v", empty, err)
	}
}

func TestPublicAPIsMarkdownCategoriesAndEscapes(t *testing.T) {
	body := publicAPIsTestHeader + "| [Máil](<https://example.com/(mail)>) | Send \\| receive mail | | Yes | No |\r\n" +
		"### Weather\n\n| API | Description | Auth | HTTPS | CORS |\n|---|---|---|---|---|\n" +
		"| [Forecast](https://weather.example/docs) | Weather | OAuth | Yes | Unknown |\n"
	got, err := parsePublicAPIsCatalog(context.Background(), []byte(body), "https://example.com/README.md")
	if err != nil || len(got) != 2 {
		t.Fatalf("entries=%#v, err=%v", got, err)
	}
	if got[0].Description != "Send | receive mail" || got[0].Link != "https://example.com/(mail)" || got[1].Category != "Weather" {
		t.Fatalf("entries=%#v", got)
	}
}

func TestPublicAPIsListRefusals(t *testing.T) {
	row := "| [Mail](https://example.com/docs) | Send mail | | Yes | No |\n"
	entries := strings.Repeat(`{},`, publicAPIsMaxEntries) + `{}`
	cases := map[string]string{
		"empty":                       "",
		"HTML":                        "<html>not a list</html>",
		"header only":                 publicAPIsTestHeader,
		"malformed after valid entry": publicAPIsTestHeader + row + "| missing | columns |\n",
		"missing category":            strings.TrimPrefix(publicAPIsTestHeader, "### Communication\n\n") + row,
		"bad link":                    publicAPIsTestHeader + "| Mail | Send mail | | Yes | No |\n",
		"invalid JSON":                `{"entries":[`,
		"missing JSON entries":        `{}`,
		"null JSON entries":           `{"entries":null}`,
		"nested JSON":                 `{"entries":[],"nested":` + strings.Repeat("[", 101) + `0` + strings.Repeat("]", 101) + `}`,
		"null JSON entry":             `{"entries":[null]}`,
		"large response":              strings.Repeat("x", DefaultMaxBytes+1),
		"large Markdown row":          publicAPIsTestHeader + strings.Repeat("x", publicAPIsMaxRowBytes+1),
		"large JSON entry":            `{"entries":[{"Description":"` + strings.Repeat("x", publicAPIsMaxRowBytes) + `"}]}`,
		"too many Markdown entries":   publicAPIsTestHeader + strings.Repeat(row, publicAPIsMaxEntries+1),
		"too many JSON entries":       `{"entries":[` + entries + `]}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := parsePublicAPIsCatalog(context.Background(), []byte(body), "https://example.com/source")
			if err == nil || got != nil || !strings.Contains(err.Error(), "https://example.com/source") {
				t.Fatalf("entries=%#v, err=%v; want no partial list and source diagnostic", got, err)
			}
		})
	}
}

func TestPublicAPIsListExactBounds(t *testing.T) {
	row := "| [Mail](https://example.com/docs) | Send mail | | Yes | No |\n"
	body := publicAPIsTestHeader + strings.Repeat(row, publicAPIsMaxEntries)
	got, err := parsePublicAPIsCatalog(context.Background(), []byte(body), "fixture")
	if err != nil || len(got) != publicAPIsMaxEntries {
		t.Fatalf("entries=%d, err=%v", len(got), err)
	}
	prefix, suffix := "| [Mail](https://example.com/docs) | ", " | | Yes | No |"
	body = publicAPIsTestHeader + prefix + strings.Repeat("x", publicAPIsMaxRowBytes-len(prefix)-len(suffix)) + suffix + "\n"
	if _, err := parsePublicAPIsCatalog(context.Background(), []byte(body), "fixture"); err != nil {
		t.Fatalf("exact row limit: %v", err)
	}
}

func TestPublicAPIsMarkdownSearchAndFailure(t *testing.T) {
	var list atomic.Value
	list.Store("")
	var probes atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/list":
			w.Header().Set("Content-Type", "application/octet-stream")
			_, _ = w.Write([]byte(list.Load().(string)))
		case "/openapi.json":
			probes.Add(1)
			_, _ = w.Write([]byte(`{"openapi":"3.0.0","info":{"title":"Mail","version":"1.0"},"paths":{}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	sourceURL := server.URL + "/list"
	list.Store(publicAPIsTestHeader + "| [Mail](" + server.URL + "/docs) | Send mail | | Yes | No |\n")
	c := &Client{PublicAPIsURL: sourceURL, AllowUnsafeHosts: true}
	report, err := c.Search(context.Background(), SearchOptions{Query: "mail", Source: SourcePublicAPIs, Limit: 1})
	if err != nil || len(report.Results) != 1 || !report.Results[0].Validated || probes.Load() != 1 {
		t.Fatalf("report=%#v, err=%v, probes=%d", report, err, probes.Load())
	}
	list.Store(list.Load().(string) + "| malformed | row |\n")
	probes.Store(0)
	report, err = c.Search(context.Background(), SearchOptions{Query: "mail", Source: SourcePublicAPIs, Limit: 1})
	if err == nil || len(report.Results) != 0 || probes.Load() != 0 || report.Attempts[0].Status != "fail" {
		t.Fatalf("partial parse was searchable: report=%#v, err=%v, probes=%d", report, err, probes.Load())
	}
	list.Store(`{"entries":[]}`)
	report, err = c.Search(context.Background(), SearchOptions{Query: "mail", Source: SourcePublicAPIs})
	if err != nil || len(report.Results) != 0 || probes.Load() != 0 {
		t.Fatalf("empty JSON report=%#v, err=%v", report, err)
	}
}

func FuzzPublicAPIsList(f *testing.F) {
	f.Add([]byte(`{"entries":[]}`))
	f.Add([]byte(publicAPIsTestHeader + "| [Mail](https://example.com) | Mail | | Yes | No |\n"))
	f.Fuzz(func(t *testing.T, body []byte) {
		entries, err := parsePublicAPIsCatalog(context.Background(), body, "fixture")
		if err != nil && entries != nil {
			t.Fatal("failure returned a partial list")
		}
		if len(entries) > publicAPIsMaxEntries {
			t.Fatal("entry bound exceeded")
		}
	})
}

func TestPublicAPIsListCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	got, err := parsePublicAPIsCatalog(ctx, []byte(`{"entries":[]}`), "fixture")
	if !errors.Is(err, context.Canceled) || got != nil {
		t.Fatalf("entries=%#v, err=%v", got, err)
	}
}
