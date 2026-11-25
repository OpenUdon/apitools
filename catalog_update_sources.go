package apitools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/OpenUdon/apitools/catalog"
)

type catalogUpdateSource struct {
	ProviderID   string
	Service      string
	ReferenceID  string
	Kind         catalog.SpecKind
	URL          string
	ResolvedURL  string
	Version      string
	HeldVersion  string
	HeldRevision string
	Revision     string
	Stage        string
	Stable       bool
	Resolution   string
	Scope        string
}

var catalogUpdateGitCommit = regexp.MustCompile(`^[0-9a-f]{40}$`)

// resolveCatalogUpdateSource resolves GitHub-backed paths to the last commit
// that touched that exact file. A held commit establishes the file's identity.
// GitHub's default branch supplies the check scope when the reference is pinned.
// Other direct URLs remain publisher locators with a separately checked digest.
func (c *Client) resolveCatalogUpdateSource(ctx context.Context, source catalogUpdateSource) (catalogUpdateSource, error) {
	u, err := c.versionURLSyntax(source.URL)
	if err != nil {
		return source, err
	}
	if u.Hostname() != "raw.githubusercontent.com" {
		source.Scope = u.Scheme + "://" + u.Host + u.EscapedPath()
		source.Resolution = "unexamined-publisher-version"
		return source, nil
	}
	parts := strings.SplitN(strings.TrimPrefix(u.Path, "/"), "/", 4)
	if u.Scheme != "https" || u.Port() != "" || len(parts) != 4 ||
		parts[0] == "" || parts[1] == "" || parts[2] == "" || parts[3] == "" ||
		parts[2] == "refs" {
		source.Resolution = "unexamined-locator"
		return source, nil
	}
	repository := url.PathEscape(parts[0]) + "/" + url.PathEscape(parts[1])
	origin := "https://api.github.com"
	prefix := "/repos/" + repository
	q := url.Values{"path": {parts[3]}, "per_page": {"1"}}
	source.Scope = "default branch, " + parts[3]
	if !catalogUpdateGitCommit.MatchString(parts[2]) {
		q.Set("sha", parts[2])
		source.Scope = "ref " + parts[2] + ", " + parts[3]
	}
	budget := &apiVersionBudget{sem: make(chan struct{}, 1)}
	scope := []APIVersionOfficialScope{{Origin: origin, PathPrefix: prefix + "/"}}
	probe := c.probeAPIVersion(ctx, budget, origin+prefix+"/commits?"+q.Encode(), scope, false, "", "")
	if probe.err != nil {
		source.Resolution = "unexamined-locator"
		return source, probe.err
	}
	var commits []struct{ SHA string }
	if err := json.Unmarshal(probe.content, &commits); err != nil || len(commits) != 1 ||
		!catalogUpdateGitCommit.MatchString(commits[0].SHA) {
		source.Resolution = "unexamined-locator"
		return source, fmt.Errorf("publisher returned no exact commit for the reviewed file")
	}
	revision := commits[0].SHA
	u.Path = "/" + parts[0] + "/" + parts[1] + "/" + revision + "/" + parts[3]
	u.RawPath = ""
	source.ResolvedURL = u.String()
	source.Revision = revision
	source.Resolution = "repository-revision"
	return source, nil
}

type catalogUpdateInventory struct {
	CheckedAt time.Time
	Scope     string
	Sources   []catalogUpdateSource
}

// catalogUpdateReferenceSources retains existing publisher locators, including
// direct artifacts and source indexes. Docs are evidence for later curation.
func catalogUpdateReferenceSources(provider catalog.Provider) []catalogUpdateSource {
	var sources []catalogUpdateSource
	for _, ref := range provider.SpecReferences {
		if ref.Kind == catalog.SpecKindHumanDocs || ref.SourceAuthority == catalog.SourceAuthorityPublicCatalog {
			continue
		}
		service, version := catalogUpdateReferenceIdentity(ref)
		sources = append(sources, catalogUpdateSource{
			ProviderID: provider.ID, Service: service, ReferenceID: ref.ID,
			Kind: ref.Kind, URL: ref.URL, Version: version, HeldVersion: ref.Version, HeldRevision: ref.Revision,
		})
	}
	return sources
}

func catalogUpdateReferenceIdentity(ref catalog.SpecReference) (service, version string) {
	u, err := url.Parse(ref.URL)
	if err != nil || u.Host == "" {
		return ref.ID, ""
	}
	normalized := u.EscapedPath()
	if u.Hostname() == "raw.githubusercontent.com" {
		parts := strings.SplitN(strings.TrimPrefix(normalized, "/"), "/", 4)
		if len(parts) == 4 && parts[2] != "refs" {
			normalized = "/" + parts[0] + "/" + parts[1] + "/" + parts[3]
		}
	}
	parts := apiVersionInURL.FindAllStringSubmatchIndex(normalized, -1)
	if len(parts) == 0 {
		return u.Host + normalized, ""
	}
	release := normalized[parts[0][2]:parts[0][3]]
	consistent := true
	for i := len(parts) - 1; i >= 0; i-- {
		begin, end := parts[i][2], parts[i][3]
		consistent = consistent && normalized[begin:end] == release
		normalized = normalized[:begin] + "{version}" + normalized[end:]
	}
	if !consistent {
		return ref.ID, ""
	}
	return u.Host + normalized, "v" + release
}

// hubSpotUpdateInventory enumerates the versioned links in the public index.
// It does not select a preferred version or fetch linked specifications.
func hubSpotUpdateInventory(sourceURL string, checkedAt time.Time, index []byte) (catalogUpdateInventory, error) {
	const origin = "https://api.hubapi.com"
	if sourceURL != origin+"/public/api/spec/v1/specs" {
		return catalogUpdateInventory{}, fmt.Errorf("unexpected HubSpot index URL")
	}
	var doc struct {
		Results []struct {
			Name     string
			Group    string
			Versions []struct {
				Version             string
				Stage               string
				OpenAPI             string
				DocumentationBanner string
			}
		}
	}
	if err := json.Unmarshal(index, &doc); err != nil {
		return catalogUpdateInventory{}, fmt.Errorf("read HubSpot version index: %w", err)
	}
	if len(doc.Results) == 0 {
		return catalogUpdateInventory{}, fmt.Errorf("HubSpot version index has no services")
	}
	report := catalogUpdateInventory{Scope: sourceURL, CheckedAt: checkedAt}
	seen := make(map[catalogUpdateSource]bool)
	for _, family := range doc.Results {
		service := strings.TrimSpace(family.Group + " " + family.Name)
		if strings.TrimSpace(family.Name) == "" {
			return catalogUpdateInventory{}, fmt.Errorf("HubSpot index contains a service without a name")
		}
		for _, entry := range family.Versions {
			if entry.OpenAPI == "" {
				continue
			}
			u, err := (&Client{}).versionURLSyntax(entry.OpenAPI)
			if err != nil || versionOrigin(u) != versionOrigin(&url.URL{Scheme: "https", Host: "api.hubapi.com"}) ||
				!strings.HasPrefix(u.Path, "/public/api/spec/") {
				return catalogUpdateInventory{}, fmt.Errorf("HubSpot index spec link is outside publisher scope")
			}
			if strings.TrimSpace(entry.Version) == "" {
				return catalogUpdateInventory{}, fmt.Errorf("HubSpot index contains an unnamed API version")
			}
			record := catalogUpdateSource{
				ProviderID: "hubspot", Service: service, ReferenceID: "hubspot-public-api-spec-index",
				Kind: catalog.SpecKindOpenAPI, URL: u.String(), Version: entry.Version, Stage: entry.Stage,
				Stable: (entry.Stage == "LATEST" || entry.Stage == "STABLE") && entry.DocumentationBanner == "",
			}
			if !seen[record] {
				report.Sources = append(report.Sources, record)
				seen[record] = true
			}
		}
	}
	sort.SliceStable(report.Sources, func(i, j int) bool {
		a, b := report.Sources[i], report.Sources[j]
		if a.Service != b.Service {
			return a.Service < b.Service
		}
		if a.Version != b.Version {
			return a.Version < b.Version
		}
		return a.URL < b.URL
	})
	return report, nil
}
