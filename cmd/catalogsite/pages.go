package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"maps"
	"net/url"
	"slices"
	"strings"
)

//go:embed layout.html assets/*
var pageAssets embed.FS

type categoryView struct {
	ID    string
	Label string
	Count int
}

type coverageView struct {
	ID    string
	Label string
	Count int
}

var coverageLabels = []coverageView{
	{ID: coverageOpenAPI, Label: "Official OpenAPI"},
	{ID: coverageMachine, Label: "Other official spec"},
	{ID: coverageOverlay, Label: "Advisory OpenAPI"},
	{ID: coverageDocs, Label: "Documentation only"},
}

type pageView struct {
	Title, Kind, Canonical, Root string
	ContactURL                   string
	Provider                     *siteProvider
	Category                     categoryView
	Categories                   []categoryView
	Providers                    []siteProvider
	Coverage                     []coverageView
	History                      []monthlyRecord
}

var siteTemplate = template.Must(template.New("layout.html").Funcs(template.FuncMap{
	"categoryLabel": categoryLabel,
	"coverageLabel": func(id string) string {
		for _, coverage := range coverageLabels {
			if coverage.ID == id {
				return coverage.Label
			}
		}
		return id
	},
	"readable": func(id string) string { return strings.ReplaceAll(id, "-", " ") },
	"categoryReference": func(sources []siteSource) *siteSource {
		for _, source := range sources {
			if source.VerifiedAt != "" {
				return &source
			}
		}
		if len(sources) != 0 {
			return &sources[0]
		}
		return nil
	},
	"searchText": func(provider siteProvider) string {
		return strings.Join([]string{
			provider.Name, provider.ID, provider.Category, provider.Description,
			strings.Join(provider.Aliases, " "), strings.Join(provider.Keywords, " "),
		}, " ")
	},
}).ParseFS(pageAssets, "layout.html"))

func categoryLabel(id string) string {
	words := strings.Split(id, "-")
	for i, word := range words {
		switch word {
		case "ai", "api", "cms", "crm", "erp", "hr", "iot", "sms":
			words[i] = strings.ToUpper(word)
		default:
			if word != "" {
				words[i] = strings.ToUpper(word[:1]) + word[1:]
			}
		}
	}
	return strings.Join(words, " ")
}

func addPages(files map[string][]byte, data siteCatalog) error {
	base, err := url.Parse(data.BaseURL)
	if err != nil {
		return err
	}
	byCategory := make(map[string][]siteProvider)
	counts := make(map[string]int)
	for _, provider := range data.Providers {
		byCategory[provider.Category] = append(byCategory[provider.Category], provider)
		counts[provider.Coverage]++
	}
	categories := make([]categoryView, 0, len(byCategory))
	for _, id := range slices.Sorted(maps.Keys(byCategory)) {
		categories = append(categories, categoryView{ID: id, Label: categoryLabel(id), Count: len(byCategory[id])})
	}
	coverage := slices.Clone(coverageLabels)
	for i := range coverage {
		coverage[i].Count = counts[coverage[i].ID]
	}
	home := pageView{
		Title: "Reviewed API catalog", Kind: "home", Canonical: base.String(), ContactURL: data.ContactURL,
		Providers: data.Providers, Categories: categories, Coverage: coverage,
	}
	if err := addPage(files, "index.html", home); err != nil {
		return err
	}
	for _, item := range categories {
		name := "categories/" + item.ID + "/index.html"
		page := pageView{
			Title: item.Label, Kind: "category", Canonical: base.JoinPath(name).String(), Root: "../../",
			ContactURL: data.ContactURL, Category: item, Providers: byCategory[item.ID],
		}
		if err := addPage(files, name, page); err != nil {
			return err
		}
	}
	for _, item := range data.Providers {
		page := pageView{
			Title: item.Name, Kind: "provider", Canonical: item.URL, Root: "../../",
			ContactURL: data.ContactURL, Provider: &item,
		}
		if err := addPage(files, "providers/"+item.ID+"/index.html", page); err != nil {
			return err
		}
	}
	history, err := loadChangeHistory()
	if err != nil {
		return err
	}
	changes := pageView{
		Title: "Monthly changes", Kind: "changes", Canonical: base.JoinPath("changes", "index.html").String(),
		Root: "../", ContactURL: data.ContactURL, History: history,
	}
	if err := addPage(files, "changes/index.html", changes); err != nil {
		return err
	}
	for _, name := range []string{"assets/site.css", "assets/search.js"} {
		contents, err := pageAssets.ReadFile(name)
		if err != nil {
			return err
		}
		files[name] = contents
	}
	return nil
}

func addPage(files map[string][]byte, name string, page pageView) error {
	var out bytes.Buffer
	if err := siteTemplate.Execute(&out, page); err != nil {
		return fmt.Errorf("render %s: %w", name, err)
	}
	files[name] = out.Bytes()
	return nil
}
