package main

import (
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/OpenUdon/apitools/catalog"
)

const (
	coverageOpenAPI = "official-openapi"
	coverageMachine = "official-machine-readable"
	coverageOverlay = "docs-with-advisory-overlay"
	coverageDocs    = "docs-only"
)

var (
	slugPattern     = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	shaPattern      = regexp.MustCompile(`^[a-f0-9]{40}$`)
	hostPattern     = regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?\.)+[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]$`)
	fragmentPattern = regexp.MustCompile(`^[a-zA-Z0-9_.:/-]+$`)
)

type siteCatalog struct {
	SchemaVersion int            `json:"schema_version"`
	SourceCommit  string         `json:"source_commit"`
	BaseURL       string         `json:"base_url"`
	ContactURL    string         `json:"contact_url"`
	Providers     []siteProvider `json:"providers"`
}

type siteProvider struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
	Aliases     []string        `json:"aliases"`
	Keywords    []string        `json:"keywords"`
	Coverage    string          `json:"coverage"`
	URL         string          `json:"url"`
	Sources     []siteSource    `json:"sources"`
	Operations  []siteOperation `json:"reviewed_operations"`
}

type siteSource struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Authority   string `json:"authority"`
	URL         string `json:"url"`
	VerifiedAt  string `json:"verified_at,omitempty"`
	LicenseNote string `json:"license_note,omitempty"`
}

type siteOperation struct {
	Label    string `json:"label"`
	SourceID string `json:"source_id"`
	Evidence string `json:"evidence"`
}

func builtInSiteCatalog(commit string, base *url.URL) (siteCatalog, error) {
	providers := catalog.BuiltInProviders()
	profiles := catalog.BuiltInProviderSearchProfiles()
	data := siteCatalog{
		SchemaVersion: 1,
		SourceCommit:  commit,
		BaseURL:       base.String(),
		ContactURL:    "https://github.com/OpenUdon/apitools/issues",
		Providers:     make([]siteProvider, 0, len(providers)),
	}
	if len(providers) != len(profiles) {
		return siteCatalog{}, fmt.Errorf("provider and public profile counts differ")
	}
	hasOverlay := make(map[string]bool)
	for _, artifact := range catalog.BuiltInEndpointOverlayArtifacts() {
		hasOverlay[artifact.ProviderID] = true
	}
	for i, provider := range providers {
		profile := profiles[i]
		if provider.ID != profile.ProviderID || !slugPattern.MatchString(provider.ID) || !slugPattern.MatchString(profile.Category) ||
			provider.ReviewState != catalog.ProviderReviewedCatalogEntry ||
			(i > 0 && provider.ID <= providers[i-1].ID) {
			return siteCatalog{}, fmt.Errorf("invalid reviewed provider or profile at index %d", i)
		}
		public, err := projectProvider(provider, profile, hasOverlay[provider.ID], base)
		if err != nil {
			return siteCatalog{}, err
		}
		data.Providers = append(data.Providers, public)
	}
	return data, nil
}

func projectProvider(provider catalog.Provider, profile catalog.ProviderSearchProfile, hasOverlay bool, base *url.URL) (siteProvider, error) {
	public := siteProvider{
		ID: provider.ID, Name: profile.DisplayName, Category: profile.Category, Description: profile.Description,
		Aliases: append([]string{}, profile.Aliases...), Keywords: append([]string{}, profile.Keywords...),
		URL:     base.JoinPath("providers", provider.ID, "index.html").String(),
		Sources: make([]siteSource, 0, len(provider.SpecReferences)), Operations: []siteOperation{},
	}
	if provider.ID == "aws-elbv2" {
		public.Description = "AWS Elastic Load Balancing v2 covers Application, Network and Gateway Load Balancers, listeners and target groups."
	}
	switch {
	case provider.OfficialOpenAPIAvailability == catalog.SpecAvailabilityKnown:
		public.Coverage = coverageOpenAPI
	case provider.OfficialMachineSpecAvailability == catalog.SpecAvailabilityKnown:
		public.Coverage = coverageMachine
	case hasOverlay:
		public.Coverage = coverageOverlay
	default:
		public.Coverage = coverageDocs
	}
	for _, ref := range provider.SpecReferences {
		if !publicSourceURL(ref.URL) {
			return siteProvider{}, fmt.Errorf("provider %s source %s has an unsafe public URL", provider.ID, ref.ID)
		}
		if ref.VerifiedAt != "" {
			if _, err := time.Parse(time.DateOnly, ref.VerifiedAt); err != nil {
				return siteProvider{}, fmt.Errorf("provider %s source %s has an invalid verification date", provider.ID, ref.ID)
			}
		}
		public.Sources = append(public.Sources, siteSource{
			ID: ref.ID, Kind: string(ref.Kind), Authority: string(ref.SourceAuthority), URL: ref.URL,
			VerifiedAt: ref.VerifiedAt, LicenseNote: ref.LicenseNote,
		})
	}
	for _, source := range profile.Coverage.Sources {
		for _, operation := range source.Operations {
			if (provider.ID == "jenkins" && operation.Label == "Get CSRF crumb") ||
				(provider.ID == "filemaker" && operation.Label == "Delete session") {
				continue
			}
			public.Operations = append(public.Operations, siteOperation{
				Label: operation.Label, SourceID: source.SpecRefID, Evidence: string(source.Kind),
			})
		}
	}
	return public, nil
}

func publicHTTPS(parsed *url.URL) bool {
	host := parsed.Hostname()
	tld := strings.ToLower(host[strings.LastIndexByte(host, '.')+1:])
	return parsed.Scheme == "https" && parsed.Opaque == "" && parsed.User == nil && parsed.Host == host &&
		hostPattern.MatchString(host) && !slices.Contains([]string{"local", "localhost", "internal", "test", "invalid", "example"}, tld)
}

func publicSourceURL(raw string) bool {
	parsed, err := url.Parse(raw)
	if err != nil || !publicHTTPS(parsed) {
		return false
	}
	if parsed.Fragment != "" && !fragmentPattern.MatchString(parsed.Fragment) {
		return false
	}
	params, err := url.ParseQuery(parsed.RawQuery)
	if err != nil {
		return false
	}
	for name := range params {
		// Permit the public document selectors used by the reviewed references.
		switch strings.ToLower(name) {
		case "version", "language_content_entity", "download", "locale", "pageid", "wikiname", "from",
			"collection", "id", "language", "release", "type", "format", "vendorid":
		default:
			return false
		}
	}
	return true
}
