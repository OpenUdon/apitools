package apitools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/OpenUdon/apitools/googlediscovery"
	"github.com/OpenUdon/apitools/internal/sourceguard"
)

type apiVersionSession struct {
	ctx        context.Context
	cancel     context.CancelFunc
	client     *Client
	req        APIVersionDiscoveryRequest
	opts       APIVersionDiscoveryOptions
	scopes     []APIVersionOfficialScope
	budget     *apiVersionBudget
	mu         sync.Mutex
	report     APIVersionDiscoveryReport
	seen       map[string]bool
	contents   map[string][]byte
	incomplete bool
	conflict   bool
	examined   bool
	reused     bool
}

// DiscoverAPIVersions performs an explicitly enabled, bounded advisory check.
// Fetch/parse/persistence failures preserve Known and return partial evidence.
func (c *Client) DiscoverAPIVersions(ctx context.Context, request APIVersionDiscoveryRequest, options APIVersionDiscoveryOptions) (APIVersionDiscoveryReport, error) {
	s, err := newAPIVersionSession(ctx, c, request, options)
	if err != nil {
		return APIVersionDiscoveryReport{}, err
	}
	defer s.cancel()
	s.discover()
	return s.finish(), nil
}

func newAPIVersionSession(ctx context.Context, c *Client, req APIVersionDiscoveryRequest, opts APIVersionDiscoveryOptions) (*apiVersionSession, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	c = c.effective()
	scopes, err := validateVersionRequest(c, req, opts)
	if err != nil {
		return nil, err
	}
	data, _ := json.Marshal(req)
	_ = json.Unmarshal(data, &req)
	if opts.BaselineInventory != nil {
		if len(opts.BaselineInventory.Operations) > 10000 {
			return nil, fmt.Errorf("baseline inventory operation budget exceeded")
		}
		data, err = json.Marshal(opts.BaselineInventory)
		if err != nil || len(data) > 32<<20 {
			return nil, fmt.Errorf("baseline inventory byte budget exceeded")
		}
		var inv OperationInventory
		_ = json.Unmarshal(data, &inv)
		opts.BaselineInventory = &inv
	}
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = DefaultAPIVersionTimeout
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	client := *c
	client.AllowedPorts = append([]int(nil), c.AllowedPorts...)
	s := &apiVersionSession{ctx: ctx, cancel: cancel, client: &client, req: req, opts: opts, scopes: scopes, budget: &apiVersionBudget{sem: make(chan struct{}, 4)}, seen: map[string]bool{}, contents: map[string][]byte{}}
	s.report = APIVersionDiscoveryReport{QueryIdentity: versionRequestIdentity(req, opts), SchemaVersion: APIVersionDiscoverySchemaVersion, Known: req.Known, Status: "unexamined", CheckedAt: time.Now().UTC(), CheckedScope: append([]APIVersionOfficialScope(nil), scopes...)}
	return s, nil
}

func versionDiagnostic(code, message string) Diagnostic {
	text, _ := sanitizePromptString(message, 2048)
	return Diagnostic{Severity: "warning", Code: code, Message: text}
}
func (s *apiVersionSession) gap(code, message string) {
	s.mu.Lock()
	s.incomplete = true
	s.report.Diagnostics = append(s.report.Diagnostics, versionDiagnostic(code, message))
	s.mu.Unlock()
}
func (s *apiVersionSession) tier(name, status string) {
	s.mu.Lock()
	s.report.Tiers = append(s.report.Tiers, APIVersionTier{Tier: name, Status: status, CheckedAt: s.report.CheckedAt})
	s.mu.Unlock()
}
func (s *apiVersionSession) add(row APIVersionRecord) {
	row.Locator.URL = strings.SplitN(strings.SplitN(row.Locator.URL, "?", 2)[0], "#", 2)[0]
	row.Locator.URLTemplate = strings.SplitN(strings.SplitN(row.Locator.URLTemplate, "?", 2)[0], "#", 2)[0]
	row.Recipe.URL = strings.SplitN(strings.SplitN(row.Recipe.URL, "?", 2)[0], "#", 2)[0]
	row.Recipe.URLTemplate = strings.SplitN(strings.SplitN(row.Recipe.URLTemplate, "?", 2)[0], "#", 2)[0]
	s.mu.Lock()
	s.report.Versions = append(s.report.Versions, row)
	s.mu.Unlock()
}

func (s *apiVersionSession) discover() {
	identity := s.report.QueryIdentity
	if s.opts.Previous != nil && s.opts.Previous.QueryIdentity == identity && freshVersionReport(*s.opts.Previous) {
		data, _ := json.Marshal(s.opts.Previous)
		if len(data) <= apiVersionMaxReportBytes {
			_ = json.Unmarshal(data, &s.report)
			s.reused = true
			s.report.Tiers = []APIVersionTier{{Tier: "local", Status: "cached", CheckedAt: time.Now().UTC()}}
			return
		}
	}
	if s.opts.StatePath != "" {
		if r, ok, err := loadVersionState(s.opts.StatePath, identity); err == nil && ok {
			s.report = r
			s.reused = true
			s.report.Tiers = []APIVersionTier{{Tier: "local", Status: "cached", CheckedAt: time.Now().UTC()}}
			return
		}
	}
	if s.req.Known.CatalogUpdated != nil && s.req.Known.CatalogUpdated.Before(time.Now().AddDate(-1, 0, 0)) {
		s.report.Diagnostics = append(s.report.Diagnostics, versionDiagnostic("version.baseline_stale", "Baseline catalog evidence is older than twelve months; consider explicitly supplied hints."))
	}
	leads := versionCatalogReferences(s.req, s.opts)
	var cache apiVersionListCache
	if s.opts.ListCachePath != "" {
		var err error
		cache, err = readVersionList(s.opts.ListCachePath)
		if err == nil {
			rows, e := versionListRecords(s.ctx, cache.Data, s.req)
			if e == nil {
				leads = append(leads, rows...)
			} else {
				s.gap("version.list_invalid", e.Error())
			}
		}
	}
	s.tier("local", "examined")
	for _, lead := range leads {
		if lead.CatalogUpdated != nil && lead.CatalogUpdated.Before(time.Now().AddDate(-1, 0, 0)) {
			s.report.Diagnostics = append(s.report.Diagnostics, versionDiagnostic("version.catalog_stale", "Directory evidence is older than twelve months; consider explicitly supplied hints."))
			break
		}
	}
	if s.req.Known.SourceKind != OperationSourceOpenAPI && s.req.Known.SourceKind != OperationSourceGoogleDiscovery {
		for _, r := range leads {
			s.add(r)
		}
		s.gap("version.family_unsupported", "This check supports OpenAPI/Swagger and Google Discovery sources only.")
		return
	}
	if !s.opts.Network {
		for _, r := range leads {
			s.add(r)
		}
		s.gap("version.network_denied", "Network lookup was not enabled; remaining official sources are unexamined.")
		return
	}
	var wg sync.WaitGroup
	run := func(f func()) { wg.Add(1); go func() { defer wg.Done(); f() }() }
	if s.req.Known.SourceURL != "" {
		run(func() {
			s.fetchSource(s.req.Known.SourceURL, APIVersionLocator{Kind: "direct", URL: s.req.Known.SourceURL}, "", s.req.Known.ETag, s.req.Known.LastModified)
			s.tier("conditional", "examined")
		})
	}
	for _, l := range s.req.Locators {
		locator := l
		run(func() { s.locate(locator) })
	}
	if s.opts.Previous != nil && s.opts.Previous.QueryIdentity == identity {
		for _, old := range s.opts.Previous.Versions {
			if old.LocatorStatus == "dead" || old.LocatorStatus == "moved" {
				recipe := old.Recipe
				run(func() { s.locate(recipe) })
			}
		}
	}
	if s.req.DocsURL != "" {
		run(func() { s.pointer(APIVersionLocator{Kind: "pointer", URL: s.req.DocsURL}) })
	}
	if s.req.Known.ProviderKey != "" || s.opts.ListCachePath != "" {
		if cache.CheckedAt.IsZero() || time.Since(cache.CheckedAt) >= 24*time.Hour {
			run(func() { s.refreshList(cache) })
		}
	}
	for _, lead := range leads {
		s.add(lead)
	}
	if len(leads) > 0 {
		run(func() { s.fetchLeads(leads) })
	}
	wg.Wait()
	if err := s.ctx.Err(); err != nil {
		s.gap("version.deadline", err.Error())
	}
	if s.req.Known.SourceURL == "" && len(s.req.Locators) == 0 && s.req.DocsURL == "" && len(leads) == 0 {
		s.gap("version.source_missing", "No official source locator was available.")
	}
}

func (s *apiVersionSession) refreshList(cache apiVersionListCache) {
	raw := versionCatalogURL(s.client, s.req)
	etag, modified := "", ""
	if cache.URL == raw {
		etag = cache.ETag
		modified = cache.LastModified
	}
	p := s.client.probeAPIVersion(s.ctx, s.budget, raw, nil, false, etag, modified)
	if p.err != nil {
		s.gap("version.catalog_fetch", p.err.Error())
		s.tier("directory", "unexamined")
		return
	}
	if p.status == 304 {
		if cache.URL != raw || len(cache.Data) == 0 {
			s.gap("version.catalog_304_without_cache", "Directory 304 has no matching cached body.")
			return
		}
		cache.CheckedAt = time.Now().UTC()
	} else if p.status == 404 {
		s.gap("version.catalog_missing", "Directory endpoint is unavailable.")
		return
	} else {
		cache = apiVersionListCache{SchemaVersion: "apitools.api-version-list-cache/v1", URL: raw, Data: p.content, ETag: p.etag, LastModified: p.modified, CheckedAt: time.Now().UTC()}
	}
	rows, err := versionListRecords(s.ctx, cache.Data, s.req)
	if err != nil {
		s.gap("version.catalog_parse", err.Error())
		return
	}
	for _, row := range rows {
		s.add(row)
	}
	s.fetchLeads(rows)
	if s.opts.ListCachePath != "" {
		data, _ := json.Marshal(cache)
		if err := writeVersionMetadata(s.ctx, s.opts.ListCachePath, data); err != nil {
			s.gap("version.list_save", err.Error())
		}
	}
	s.tier("directory", "examined")
}

func (s *apiVersionSession) fetchLeads(rows []APIVersionRecord) {
	copyRows := append([]APIVersionRecord(nil), rows...)
	sort.Slice(copyRows, func(i, j int) bool {
		if n, ok := compareAPIVersions(copyRows[i].VersionClaims.URLToken, copyRows[j].VersionClaims.URLToken); ok && n != 0 {
			return n > 0
		}
		return copyRows[i].SourceURL < copyRows[j].SourceURL
	})
	for _, row := range copyRows {
		if n, ok := compareAPIVersions(row.VersionClaims.URLToken, s.req.Known.Version); ok && n <= 0 {
			continue
		}
		raw := row.Locator.URL
		if err := s.client.versionScopeAllows(raw, s.scopes); err != nil {
			s.gap("version.origin_unconfirmed", "Catalog source remains an unconfirmed origin lead.")
			continue
		}
		if row.Locator.Kind == "pointer" {
			s.pointer(row.Locator)
			continue
		}
		s.fetchSource(raw, row.Locator, row.VersionClaims.URLToken, "", "")
	}
}

func (s *apiVersionSession) locate(l APIVersionLocator) {
	switch l.Kind {
	case "direct":
		s.fetchSource(l.URL, l, versionTokenURL(l.URL), "", "")
	case "pattern":
		token := l.BaselineToken
		misses := 0
		for i := 0; i < 12 && misses < 2; i++ {
			next, ok := nextAPIVersion(token)
			if !ok {
				s.gap("version.pattern_unsupported", "Only integer/dotted non-date version slots may be generated.")
				return
			}
			token = next
			status := s.fetchSource(strings.Replace(l.URLTemplate, "{version}", token, 1), l, token, "", "")
			if status == 404 {
				misses++
			} else {
				misses = 0
			}
			if status == 0 || s.ctx.Err() != nil {
				return
			}
		}
		if misses < 2 {
			s.gap("version.pattern_partial", "Sibling probing ended before the two-miss boundary.")
		} else {
			s.mu.Lock()
			s.examined = true
			s.mu.Unlock()
		}
		s.tier("siblings", "examined")
	case "pointer":
		s.pointer(l)
	case "github":
		s.github(l)
	case "template", "gated":
		s.gap("version.locator_gated", "Tenant-specific or gated locator needs caller source guidance; consider docs-derived overlays.")
		s.add(APIVersionRecord{ID: versionRecordID(l.URL + l.URLTemplate + l.Kind), Locator: l, Recipe: l, Evidence: "official-candidate", Validation: "unexamined", Comparison: "unexamined", LocatorStatus: "gated", CheckedAt: s.report.CheckedAt})
	}
}

func (s *apiVersionSession) fetchSource(raw string, l APIVersionLocator, claim, etag, modified string) int {
	if raw == "" {
		s.gap("version.locator_missing", "Source locator has no URL.")
		return 0
	}
	if err := s.client.versionScopeAllows(raw, s.scopes); err != nil {
		s.gap("version.origin_unconfirmed", err.Error())
		return 0
	}
	s.mu.Lock()
	if s.seen[raw] {
		s.mu.Unlock()
		return 200
	}
	s.seen[raw] = true
	s.mu.Unlock()
	if claim == "" {
		claim = versionTokenURL(raw)
	}
	p := s.client.probeAPIVersion(s.ctx, s.budget, raw, s.scopes, true, etag, modified)
	row := APIVersionRecord{ID: versionRecordID(raw), SourceKind: s.req.Known.SourceKind, Locator: l, Recipe: l, Evidence: "official-candidate", SourceURL: publicVersionURL(raw), FinalURL: publicVersionURL(p.finalURL), SHA256: p.digest, BytesObserved: p.bytes, BytesDeclared: p.declared, CheckedAt: s.report.CheckedAt, VersionClaims: APIVersionClaims{URLToken: claim}, Validation: "unexamined", Comparison: "unexamined", LocatorStatus: "ok", ETag: p.etag, LastModified: p.modified}
	if p.status == 404 {
		row.LocatorStatus = "dead"
		s.add(row)
		return 404
	}
	if p.status == 304 {
		row.Validation = "unchanged"
		row.SHA256 = s.req.Known.SHA256
		s.add(row)
		return 304
	}
	if p.status == 401 || p.status == 403 {
		row.LocatorStatus = "gated"
	}
	if p.err != nil || p.nonSpec {
		message := "Definite non-spec response was refused after the sniff prefix."
		if p.err != nil {
			message = p.err.Error()
		}
		if p.retryAfter != "" {
			message += "; Retry-After: " + p.retryAfter
		}
		row.Diagnostics = append(row.Diagnostics, versionDiagnostic("version.fetch", message))
		s.gap("version.fetch", message)
		s.add(row)
		return 0
	}
	now := time.Now().UTC()
	row.FetchedAt = &now
	root, err := parseInventoryDocument(p.content)
	if err == nil {
		if info, ok := root["info"].(map[string]any); ok {
			row.VersionClaims.InfoVersion = stringValue(info["version"])
		}
		if _, ok := downloadedSpecMetadata(s.ctx, p.content, p.finalURL); !ok {
			err = fmt.Errorf("source failed strict Import validation")
		}
		row.SourceKind = OperationSourceOpenAPI
	} else {
		var google map[string]any
		if e := json.Unmarshal(p.content, &google); e == nil && stringValue(google["discoveryVersion"]) != "" {
			if e = sourceguard.CheckJSONWithLimits(s.ctx, "version-discovery", p.content, sourceguard.DefaultLimits()); e == nil {
				if model, e := googlediscovery.Parse(p.content); e == nil {
					row.SourceKind = OperationSourceGoogleDiscovery
					row.VersionClaims.InfoVersion = model.Version
					err = nil
				} else {
					err = e
				}
			}
		}
	}
	if err != nil {
		row.Validation = "invalid"
		row.Diagnostics = append(row.Diagnostics, versionDiagnostic("version.validation", err.Error()))
		s.gap("version.validation", err.Error())
		s.add(row)
		return p.status
	}
	info, changed := sanitizePromptString(row.VersionClaims.InfoVersion, 256)
	row.VersionClaims.InfoVersion = info
	if changed {
		row.Diagnostics = append(row.Diagnostics, versionDiagnostic("version.metadata_loss", "Version metadata needs source review after sanitation."))
		s.gap("version.metadata_loss", "Version identity metadata was shortened or sanitized.")
		s.add(row)
		return p.status
	}
	row.Validation = "valid"
	row.Evidence = "official-verified"
	version := firstNonEmpty(claim, info)
	n, comparable := compareAPIVersions(version, s.req.Known.Version)
	if comparable && n > 0 {
		row.Comparison = "newer"
	} else if comparable {
		row.Comparison = "not_newer"
	}
	if claim != "" && info != "" {
		if n, ok := compareAPIVersions(claim, info); ok && n != 0 {
			row.Comparison = "conflicting"
			s.mu.Lock()
			s.conflict = true
			s.mu.Unlock()
			row.Diagnostics = append(row.Diagnostics, versionDiagnostic("version.identity_conflict", "URL and info.version disagree."))
		}
	}
	if row.Comparison == "unexamined" {
		s.gap("version.order_unknown", "Version labels cannot be compared honestly.")
	}
	if row.SourceKind != s.req.Known.SourceKind {
		row.Comparison = "unexamined"
		s.gap("version.family_changed", "Different source families cannot establish a comparable upgrade.")
	}
	if s.opts.SaveDir != "" && row.Comparison == "newer" && row.SourceKind == OperationSourceOpenAPI {
		saved, e := saveVersionSource(s.ctx, s.opts.SaveDir, row, p.content)
		if e != nil {
			row.Diagnostics = append(row.Diagnostics, versionDiagnostic("version.save", e.Error()))
		} else {
			row.Saved = saved
		}
	}
	s.mu.Lock()
	s.contents[row.ID] = p.content
	s.mu.Unlock()
	s.add(row)
	return p.status
}

func (s *apiVersionSession) pointer(l APIVersionLocator) {
	if err := s.client.versionScopeAllows(l.URL, s.scopes); err != nil {
		s.gap("version.pointer_scope", err.Error())
		return
	}
	p := s.client.probeAPIVersion(s.ctx, s.budget, l.URL, s.scopes, false, "", "")
	if p.err != nil || p.status < 200 || p.status >= 300 || len(p.content) > 1<<20 {
		s.gap("version.pointer", "Pointer could not be examined within its bounds.")
		return
	}
	var links []string
	if strings.Contains(p.contentType, "application/linkset+json") {
		doc, err := parseRFC9727Document(p.content)
		if err != nil {
			s.gap("version.pointer_parse", err.Error())
			return
		}
		count := 0
		for _, set := range doc.Linkset {
			targets, err := parseRFC9727Targets(set.ServiceDesc)
			if err != nil {
				s.gap("version.pointer_parse", err.Error())
				return
			}
			for _, t := range targets {
				count++
				if count > 100 {
					s.gap("version.link_limit", "Publisher catalog link limit reached.")
					break
				}
				if openAPIDescriptionTarget(t) {
					links = append(links, t.Href)
				}
			}
		}
	} else {
		links = versionDocsLinks(string(p.content))
		if len(links) > 32 {
			links = links[:32]
			s.gap("version.link_limit", "Documentation link limit reached; narrow the pointer.")
		}
	}
	base, _ := url.Parse(p.finalURL)
	seen := map[string]bool{}
	var rows []APIVersionRecord
	for _, link := range links {
		u, err := base.Parse(link)
		if err != nil {
			continue
		}
		raw := u.String()
		if seen[raw] {
			continue
		}
		seen[raw] = true
		if err := s.client.versionScopeAllows(raw, s.scopes); err != nil {
			s.gap("version.pointer_scope", err.Error())
			continue
		}
		locator := l
		rows = append(rows, APIVersionRecord{SourceURL: raw, Locator: APIVersionLocator{Kind: "direct", URL: raw}, Recipe: locator, VersionClaims: APIVersionClaims{URLToken: versionTokenURL(raw)}})
	}
	sort.Slice(rows, func(i, j int) bool {
		if n, ok := compareAPIVersions(rows[i].VersionClaims.URLToken, rows[j].VersionClaims.URLToken); ok && n != 0 {
			return n > 0
		}
		return rows[i].SourceURL < rows[j].SourceURL
	})
	for _, r := range rows {
		if n, ok := compareAPIVersions(r.VersionClaims.URLToken, s.req.Known.Version); ok && n <= 0 {
			s.add(APIVersionRecord{ID: versionRecordID(r.SourceURL), Locator: r.Locator, Recipe: l, SourceURL: publicVersionURL(r.SourceURL), Evidence: "official-candidate", VersionClaims: r.VersionClaims, Validation: "unfetched", Comparison: "not_newer", LocatorStatus: "ok"})
			continue
		}
		s.fetchSource(r.SourceURL, l, r.VersionClaims.URLToken, "", "")
	}
	s.mu.Lock()
	s.examined = true
	s.mu.Unlock()
	s.tier("publisher", "examined")
}

func versionDocsLinks(text string) []string {
	var out []string
	for len(text) > 0 {
		index := strings.Index(strings.ToLower(text), "href")
		if index < 0 {
			break
		}
		text = text[index+4:]
		text = strings.TrimLeft(text, " \t\r\n")
		if !strings.HasPrefix(text, "=") {
			continue
		}
		text = strings.TrimLeft(text[1:], " \t\r\n")
		if len(text) == 0 {
			break
		}
		quote := text[0]
		if quote != '\'' && quote != '"' {
			continue
		}
		text = text[1:]
		end := strings.IndexByte(text, quote)
		if end < 0 {
			break
		}
		link := text[:end]
		text = text[end+1:]
		lower := strings.ToLower(strings.SplitN(link, "?", 2)[0])
		if strings.HasSuffix(lower, ".json") || strings.HasSuffix(lower, ".yaml") || strings.HasSuffix(lower, ".yml") {
			out = append(out, link)
		}
		if len(out) > 32 {
			break
		}
	}
	return out
}

func (s *apiVersionSession) github(l APIVersionLocator) {
	prefix := "https://raw.githubusercontent.com/" + l.Repository + "/"
	authorized := false
	for _, scope := range s.scopes {
		if strings.Contains(scope.Origin, "raw.githubusercontent.com") && scope.Repository == l.Repository {
			authorized = true
		}
	}
	if !authorized {
		s.gap("version.github_scope", "GitHub repository was not declared official.")
		return
	}
	raw := "https://api.github.com/repos/" + l.Repository + "/git/trees/" + url.PathEscape(l.Ref) + "?recursive=1"
	p := s.client.probeAPIVersion(s.ctx, s.budget, raw, []APIVersionOfficialScope{{Origin: "https://api.github.com", PathPrefix: "/repos/" + l.Repository + "/"}}, false, "", "")
	if p.err != nil {
		s.gap("version.github_tree", p.err.Error())
		return
	}
	var tree struct {
		SHA       string `json:"sha"`
		Truncated bool   `json:"truncated"`
		Tree      []struct {
			Path string `json:"path"`
			Type string `json:"type"`
		} `json:"tree"`
	}
	if err := json.Unmarshal(p.content, &tree); err != nil || tree.Truncated || len(tree.Tree) > 10000 {
		s.gap("version.github_tree", "GitHub tree is invalid or incomplete.")
		return
	}
	var rows []APIVersionRecord
	for _, item := range tree.Tree {
		match, _ := path.Match(l.PathPattern, item.Path)
		if !match || item.Type != "blob" {
			continue
		}
		u := prefix + url.PathEscape(firstNonEmpty(tree.SHA, l.Ref)) + "/" + item.Path
		rows = append(rows, APIVersionRecord{SourceURL: u, Locator: APIVersionLocator{Kind: "direct", URL: u}, Recipe: l, VersionClaims: APIVersionClaims{URLToken: versionTokenURL(u)}})
	}
	s.fetchLeads(rows)
	s.tier("github", "examined")
}

func (s *apiVersionSession) finish() APIVersionDiscoveryReport {
	if s.reused {
		return s.report
	}
	s.compareVersions()
	sort.Slice(s.report.Versions, func(i, j int) bool {
		a, b := s.report.Versions[i], s.report.Versions[j]
		av := firstNonEmpty(a.VersionClaims.URLToken, a.VersionClaims.InfoVersion)
		bv := firstNonEmpty(b.VersionClaims.URLToken, b.VersionClaims.InfoVersion)
		if n, ok := compareAPIVersions(av, bv); ok && n != 0 {
			return n > 0
		}
		if a.ID != b.ID {
			return a.ID < b.ID
		}
		return a.Evidence < b.Evidence
	})
	var preferred string
	for _, r := range s.report.Versions {
		if r.Validation == "valid" && r.Comparison == "newer" {
			v := firstNonEmpty(r.VersionClaims.URLToken, r.VersionClaims.InfoVersion)
			if preferred == "" {
				preferred = v
			} else if n, ok := compareAPIVersions(v, preferred); ok && n > 0 {
				preferred = v
			}
		}
	}
	if s.conflict {
		s.report.Status = "conflicting"
	} else if preferred != "" {
		s.report.Status = "newer_found"
		s.report.Preferred = preferred
	} else if s.examined && !s.incomplete {
		s.report.Status = "none_found_in_scope"
	} else {
		s.report.Status = "unexamined"
	}
	key := s.req.Known.ProviderKey
	if key != "" {
		versions := map[string]any{}
		for _, r := range s.report.Versions {
			if r.Validation != "valid" {
				continue
			}
			v := firstNonEmpty(r.VersionClaims.URLToken, r.VersionClaims.InfoVersion)
			if v == "" {
				continue
			}
			item := map[string]any{"info": map[string]any{"version": r.VersionClaims.InfoVersion, "x-origin": []any{map[string]any{"url": r.FinalURL}}}, "updated": r.CheckedAt, "x-apitools-locator": r.Locator, "x-apitools-recipe": r.Recipe, "x-apitools-sha256": r.SHA256, "x-apitools-evidence": r.Evidence}
			if r.SourceKind == OperationSourceOpenAPI {
				item["swaggerUrl"] = r.FinalURL
			} else {
				item["x-apitools-native-url"] = r.FinalURL
				item["x-apitools-source-kind"] = r.SourceKind
			}
			if _, exists := versions[v]; exists {
				s.report.Status = "conflicting"
				s.report.Preferred = ""
			} else {
				versions[v] = item
			}
		}
		entry := map[string]any{"versions": versions, "x-apitools-checked-scope": s.scopes, "x-apitools-checked-at": s.report.CheckedAt}
		if s.report.Preferred != "" {
			entry["preferred"] = s.report.Preferred
		}
		s.report.DirectoryRecords = map[string]any{key: entry}
	}
	count, bodies := s.budget.counts()
	kept := s.report.Tiers[:0]
	for _, t := range s.report.Tiers {
		if t.Tier != "network" {
			kept = append(kept, t)
		}
	}
	s.report.Tiers = kept
	s.report.Tiers = append(s.report.Tiers, APIVersionTier{Tier: "network", Status: s.report.Status, RequestCount: count, SourceBodyCount: bodies, CheckedAt: s.report.CheckedAt})
	sort.Slice(s.report.Tiers, func(i, j int) bool { return s.report.Tiers[i].Tier < s.report.Tiers[j].Tier })
	sort.Slice(s.report.Diagnostics, func(i, j int) bool {
		a, b := s.report.Diagnostics[i], s.report.Diagnostics[j]
		if a.Code != b.Code {
			return a.Code < b.Code
		}
		return a.Message < b.Message
	})
	data, _ := json.Marshal(s.report)
	if len(data) > apiVersionMaxReportBytes {
		s.report.Versions = nil
		s.report.DirectoryRecords = nil
		s.report.Preferred = ""
		s.report.Status = "unexamined"
		s.report.Truncated = true
		s.report.Diagnostics = []Diagnostic{versionDiagnostic("version.report_limit", "Version report exceeded the prompt context budget; narrow the query.")}
	}
	if s.opts.StatePath != "" {
		data, _ := json.Marshal(apiVersionState{SchemaVersion: APIVersionStateSchemaVersion, Identity: s.report.QueryIdentity, Report: s.report})
		if err := writeVersionMetadata(s.ctx, s.opts.StatePath, data); err != nil {
			s.report.Diagnostics = append(s.report.Diagnostics, versionDiagnostic("version.state_save", err.Error()))
		}
	}
	return s.report
}
