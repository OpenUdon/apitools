package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// The fixture supplies UI input to the production search without judging it.
// Chrome renders the response and the Go test checks the resulting DOM.
const browserSearchInput = `addEventListener("DOMContentLoaded", () => {
  const params = new URLSearchParams(location.search);
  const input = document.querySelector("#search-query");
  input.value = params.get("q") || "";
  input.dispatchEvent(new Event("input", {bubbles: true}));
  if (params.has("reset")) document.querySelector("#provider-search button").click();
});`

func TestBrowserProviderSearch(t *testing.T) {
	chrome := os.Getenv("APITOOLS_SITE_TEST_CHROME")
	if chrome == "" {
		t.Skip("set APITOOLS_SITE_TEST_CHROME for the catalog site browser gate")
	}
	if !filepath.IsAbs(chrome) {
		t.Fatal("APITOOLS_SITE_TEST_CHROME must name the installed Chromium binary by absolute path")
	}
	base, err := parseBaseURL("https://apis.kinet.com/")
	if err != nil {
		t.Fatal(err)
	}
	data, err := builtInSiteCatalog(testCommit, base)
	if err != nil {
		t.Fatal(err)
	}
	files, err := buildFiles(data)
	if err != nil {
		t.Fatal(err)
	}
	testPage := strings.Replace(string(files["index.html"]), "</body>", `<script src="/search-input.js" defer></script></body>`, 1)
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/":
			response.Header().Set("Content-Type", "text/html; charset=utf-8")
			if request.URL.Query().Has("noscript") {
				response.Header().Set("Content-Security-Policy", "script-src 'none'")
				fmt.Fprint(response, string(files["index.html"]))
				return
			}
			fmt.Fprint(response, testPage)
		case "/assets/search.js":
			response.Header().Set("Content-Type", "text/javascript")
			response.Write(files["assets/search.js"])
		case "/search-input.js":
			response.Header().Set("Content-Type", "text/javascript")
			fmt.Fprint(response, browserSearchInput)
		case "/assets/site.css":
			response.Header().Set("Content-Type", "text/css")
			response.Write(files["assets/site.css"])
		default:
			http.NotFound(response, request)
		}
	}))
	defer server.Close()

	for _, tc := range []struct {
		name, query, visible, hidden, status string
		reset, noscript                      bool
	}{
		{name: "name", query: "Action Network", visible: "action-network"},
		{name: "alias", query: "osdi", visible: "action-network"},
		{name: "category", query: "community", visible: "action-network"},
		{name: "keyword", query: "activist outreach", visible: "action-network"},
		{name: "no operation inference", query: "list user", hidden: "google-calendar"},
		{name: "empty", query: "qqq-unmatched-zz", hidden: "action-network", status: "No matching providers"},
		{name: "clear", query: "qqq-unmatched-zz", visible: "action-network", reset: true, status: fmt.Sprintf("%d providers", len(data.Providers))},
		{name: "static navigation", query: "qqq-unmatched-zz", visible: "action-network", noscript: true, status: fmt.Sprintf("%d providers", len(data.Providers))},
	} {
		t.Run(tc.name, func(t *testing.T) {
			params := url.Values{"q": {tc.query}}
			if tc.reset {
				params.Set("reset", "1")
			}
			if tc.noscript {
				params.Set("noscript", "1")
			}
			ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
			defer cancel()
			home := t.TempDir()
			command := exec.CommandContext(ctx, chrome,
				"--headless", "--disable-background-networking", "--disable-component-update", "--no-first-run",
				"--host-resolver-rules=MAP * ~NOTFOUND, EXCLUDE 127.0.0.1", "--no-proxy-server",
				"--user-data-dir="+filepath.Join(home, "profile"), "--dump-dom", server.URL+"/?"+params.Encode())
			command.Env = []string{"HOME=" + home, "XDG_CONFIG_HOME=" + home, "PATH=/usr/bin:/bin"}
			command.WaitDelay = 2 * time.Second
			stdout, err := command.Output()
			if err != nil {
				t.Fatalf("Chromium did not render the loopback catalog: %v", err)
			}
			visible := map[string]bool{}
			for _, row := range browserProviderRow.FindAllStringSubmatch(string(stdout), -1) {
				visible[row[2]] = !strings.Contains(row[1], `hidden=""`)
			}
			if len(visible) != len(data.Providers) {
				t.Fatalf("Chromium rendered %d provider rows, want %d", len(visible), len(data.Providers))
			}
			if tc.visible != "" && !visible[tc.visible] {
				t.Errorf("%q should show %s", tc.query, tc.visible)
			}
			if tc.hidden != "" && visible[tc.hidden] {
				t.Errorf("%q should hide %s", tc.query, tc.hidden)
			}
			if tc.status != "" && !strings.Contains(string(stdout), `role="status">`+tc.status+`</p>`) {
				t.Errorf("%q should show status %q", tc.query, tc.status)
			}
		})
	}
}

var browserProviderRow = regexp.MustCompile(`(?s)<li (data-search-text="[^"]*"[^>]*)>\s*<div class="provider-title">\s*<h3><a href="providers/([a-z0-9-]+)/index.html">`)
