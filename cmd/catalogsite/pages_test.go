package main

import (
	"bytes"
	"html"
	"net/url"
	"path"
	"regexp"
	"strings"
	"testing"
)

func TestSiteNavigationAndAccessibility(t *testing.T) {
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
	home := files["index.html"]
	privacy := files["privacy/index.html"]
	for _, want := range []string{
		`<h1>Privacy &amp; Site Notes</h1>`,
		`<link rel="canonical" href="https://apis.kinet.com/privacy/index.html">`,
		`href="` + data.ContactURL + `">APItools issue tracker</a>`,
	} {
		if !bytes.Contains(privacy, []byte(want)) {
			t.Errorf("privacy page lacks %s", want)
		}
	}
	for _, provider := range data.Providers {
		name := "providers/" + provider.ID + "/index.html"
		category := "categories/" + provider.Category + "/index.html"
		if !bytes.Contains(home, []byte(`href="`+name+`"`)) || !bytes.Contains(home, []byte(`href="`+category+`"`)) || len(files[name]) == 0 || len(files[category]) == 0 {
			t.Fatalf("provider %s or category %s missing static navigation", provider.ID, provider.Category)
		}
		for _, source := range provider.Sources {
			if source.VerifiedAt == "" {
				continue
			}
			categoryHTML := files[category]
			_, entry, found := bytes.Cut(categoryHTML, []byte(`href="../../`+name+`"`))
			entry, _, closed := bytes.Cut(entry, []byte("</li>"))
			if !found || !closed || !bytes.Contains(entry, []byte(`href="`+html.EscapeString(source.URL)+`"`)) ||
				!bytes.Contains(entry, []byte(`datetime="`+source.VerifiedAt+`"`)) ||
				!bytes.Contains(entry, []byte(strings.ReplaceAll(source.Authority, "-", " "))) {
				t.Errorf("%s category entry lacks verification evidence for %s", provider.ID, source.ID)
			}
			break
		}
	}
	for name, contents := range files {
		if !strings.HasSuffix(name, ".html") {
			continue
		}
		for _, check := range []struct {
			name string
			want string
		}{
			{"document language", `<html lang="en">`},
			{"main landmark", `<main class="wrap" id="content"`},
			{"labeled navigation landmark", `<nav class="wrap" aria-label="Main">`},
			{"footer landmark", `<footer class="site-footer">`},
			{"content security policy", `content="default-src 'none'; style-src 'self'; script-src 'self'; base-uri 'none'; form-action 'none'"`},
			{"canonical URL", `<link rel="canonical" href="https://apis.kinet.com/`},
		} {
			if !bytes.Contains(contents, []byte(check.want)) {
				t.Errorf("%s lacks %s", name, check.name)
			}
		}
		if !titleTag.Match(contents) {
			t.Errorf("%s lacks a nonempty title", name)
		}
		_, footer, _ := bytes.Cut(contents, []byte(`<footer class="site-footer">`))
		privacyHref := strings.Repeat("../", strings.Count(name, "/")) + "privacy/index.html"
		if !bytes.Contains(footer, []byte(`href="`+privacyHref+`"`)) || !bytes.Contains(footer, []byte(`>Privacy &amp; Site Notes</a>`)) {
			t.Errorf("%s lacks a privacy page link in the footer", name)
		}
		if strings.HasPrefix(name, "categories/") || name == "index.html" {
			if !bytes.Contains(contents, []byte(`<label for="search-query">Search providers</label>`)) || !bytes.Contains(contents, []byte(`<input type="search" id="search-query"`)) {
				t.Errorf("%s lacks a visible search label associated with its input", name)
			}
		}
		for _, match := range anchorTag.FindAllSubmatch(contents, -1) {
			if strings.TrimSpace(html.UnescapeString(stripTags.ReplaceAllString(string(match[1]), ""))) == "" {
				t.Errorf("%s has a link without readable text", name)
			}
		}
		assertLinksResolve(t, files, name, contents)
	}
}

var (
	titleTag  = regexp.MustCompile(`<title>[^<\s][^<]*</title>`)
	anchorTag = regexp.MustCompile(`(?s)<a\s+[^>]*>(.*?)</a>`)
	stripTags = regexp.MustCompile(`<[^>]*>`)
	linkAttr  = regexp.MustCompile(`(?:href|src)="([^"]+)"`)
)

func assertLinksResolve(t *testing.T, files map[string][]byte, name string, contents []byte) {
	t.Helper()
	for _, match := range linkAttr.FindAllSubmatch(contents, -1) {
		link, err := url.Parse(html.UnescapeString(string(match[1])))
		if err != nil || (link.IsAbs() && !publicHTTPS(link)) || link.Host != "" && !link.IsAbs() {
			t.Errorf("%s has invalid link %q", name, match[1])
			continue
		}
		if link.IsAbs() {
			continue
		}
		target := name
		if link.Path != "" {
			target = path.Join(path.Dir(name), link.Path)
		}
		linked, exists := files[target]
		if !exists {
			t.Errorf("%s link %q does not resolve to an output file", name, match[1])
			continue
		}
		if link.Fragment != "" && !bytes.Contains(linked, []byte(`id="`+link.Fragment+`"`)) {
			t.Errorf("%s link %q does not resolve to an anchor", name, match[1])
		}
	}
}

func TestPublicProfileSearchTextAndEscaping(t *testing.T) {
	data := siteCatalog{BaseURL: "https://apis.kinet.com/", ContactURL: "https://github.com/OpenUdon/apitools/issues", Providers: []siteProvider{{
		ID: "fixture", Name: "<script>alert(1)</script>", Description: `"><img src=x>`, Category: "data-platform",
		Aliases: []string{"Lake Engine"}, Keywords: []string{"warehouse"}, Coverage: coverageOverlay,
		URL:        "https://apis.kinet.com/providers/fixture/index.html",
		Sources:    []siteSource{{ID: "fixture-docs", URL: "https://docs.example.com/", LicenseNote: "<iframe>"}},
		Operations: []siteOperation{{Label: "List unique widget", SourceID: "fixture-docs", Evidence: "advisory-overlay"}},
	}}}
	files, err := buildFiles(data)
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range []string{"index.html", "categories/data-platform/index.html", "providers/fixture/index.html"} {
		contents := files[file]
		for _, unsafe := range []string{"<script>alert(1)</script>", "<img src=x>", "<iframe>"} {
			if bytes.Contains(contents, []byte(unsafe)) {
				t.Errorf("%s rendered unescaped markup %q", file, unsafe)
			}
		}
	}
	directory := files["index.html"]
	if !bytes.Contains(directory, []byte("Lake Engine warehouse")) || bytes.Contains(directory, []byte("List unique widget")) {
		t.Error("directory search text does not isolate public profiles from operation labels")
	}
	if !bytes.Contains(files["providers/fixture/index.html"], []byte("List unique widget")) {
		t.Error("provider page omits selected operation evidence")
	}
}
