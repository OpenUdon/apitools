package apitools

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/OpenUdon/apitools/catalog"
)

var apiVersionNumeric = regexp.MustCompile(`(?i)^v?[0-9]+(?:\.[0-9]+)*$`)
var apiVersionInURL = regexp.MustCompile(`(?i)(?:^|[/_\-])v([0-9]+(?:\.[0-9]+)*)(?:[/_\-]|\.(?:json|ya?ml)|$)`)
var apiVersionDate = regexp.MustCompile(`^[0-9]{4}[.\-][0-9]{2}[.\-][0-9]{2}$`)

func numericAPIVersion(value string) ([]uint64, bool) {
	if !apiVersionNumeric.MatchString(value) || apiVersionDate.MatchString(strings.TrimPrefix(strings.ToLower(value), "v")) {
		return nil, false
	}
	var parts []uint64
	for _, p := range strings.Split(strings.TrimPrefix(strings.ToLower(value), "v"), ".") {
		n, err := strconv.ParseUint(p, 10, 64)
		if err != nil {
			return nil, false
		}
		parts = append(parts, n)
	}
	for len(parts) > 1 && parts[len(parts)-1] == 0 {
		parts = parts[:len(parts)-1]
	}
	return parts, true
}

func compareAPIVersions(a, b string) (int, bool) {
	x, ok := numericAPIVersion(a)
	if !ok {
		return 0, false
	}
	y, ok := numericAPIVersion(b)
	if !ok {
		return 0, false
	}
	for i := 0; i < len(x) || i < len(y); i++ {
		var p, q uint64
		if i < len(x) {
			p = x[i]
		}
		if i < len(y) {
			q = y[i]
		}
		if p > q {
			return 1, true
		}
		if p < q {
			return -1, true
		}
	}
	return 0, true
}

func apiVersionLess(a, b string) bool {
	_, an := numericAPIVersion(a)
	_, bn := numericAPIVersion(b)
	if an != bn {
		return an
	}
	if n, ok := compareAPIVersions(a, b); ok && n != 0 {
		return n > 0
	}
	return a < b
}

func nextAPIVersion(value string) (string, bool) {
	if _, ok := numericAPIVersion(value); !ok {
		return "", false
	}
	parts := strings.Split(value, ".")
	last := parts[len(parts)-1]
	prefix := ""
	if len(parts) == 1 && strings.HasPrefix(strings.ToLower(last), "v") {
		prefix = last[:1]
		last = last[1:]
	}
	n, err := strconv.ParseUint(last, 10, 64)
	if err != nil || n == ^uint64(0) {
		return "", false
	}
	parts[len(parts)-1] = prefix + strconv.FormatUint(n+1, 10)
	return strings.Join(parts, "."), true
}

func versionTokenURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	all := apiVersionInURL.FindAllStringSubmatch(u.Path, -1)
	if len(all) > 0 {
		return "v" + all[len(all)-1][1]
	}
	return ""
}

func publicVersionURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	u.User = nil
	u.RawQuery = ""
	u.Fragment = ""
	return u.String()
}

func (c *Client) versionURLSyntax(raw string) (*url.URL, error) {
	u, err := c.validateCacheURL(raw)
	if err != nil {
		return nil, err
	}
	decoded := u.Path
	for i := 0; i < 2; i++ {
		d, e := url.PathUnescape(decoded)
		if e != nil {
			return nil, e
		}
		decoded = d
	}
	for _, seg := range strings.Split(decoded, "/") {
		if seg == "." || seg == ".." {
			return nil, fmt.Errorf("URL path escape is not permitted")
		}
	}
	query, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("invalid source URL query")
	}
	for k := range query {
		key := strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(k), "-", ""), "_", "")
		if key == "key" {
			return nil, fmt.Errorf("credential-bearing source URL is not permitted")
		}
		for _, s := range []string{"token", "secret", "password", "credential", "signature", "api_key", "apikey", "authorization"} {
			if strings.Contains(key, s) {
				return nil, fmt.Errorf("credential-bearing source URL is not permitted")
			}
		}
	}
	return u, nil
}

func versionOrigin(u *url.URL) string {
	port := u.Port()
	if port == "" {
		if u.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	return strings.ToLower(u.Scheme + "://" + u.Hostname() + ":" + port)
}

func (c *Client) versionScopeAllows(raw string, scopes []APIVersionOfficialScope) error {
	u, err := c.versionURLSyntax(raw)
	if err != nil {
		return err
	}
	for _, s := range scopes {
		base, e := url.Parse(s.Origin)
		if e != nil || versionOrigin(u) != versionOrigin(base) {
			continue
		}
		prefix := strings.TrimSuffix(s.PathPrefix, "/")
		if prefix == "" {
			prefix = "/"
		}
		if prefix == "/" || u.Path == prefix || strings.HasPrefix(u.Path, prefix+"/") {
			return nil
		}
	}
	return fmt.Errorf("source URL is outside the approved official scope")
}

func scopeForVersionURL(raw string) APIVersionOfficialScope {
	u, _ := url.Parse(raw)
	s := APIVersionOfficialScope{Origin: u.Scheme + "://" + u.Host, PathPrefix: path.Dir(u.Path) + "/"}
	if strings.EqualFold(u.Hostname(), "raw.githubusercontent.com") || strings.EqualFold(u.Hostname(), "github.com") {
		p := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(p) >= 2 {
			s.Repository = p[0] + "/" + p[1]
			s.PathPrefix = "/" + s.Repository + "/"
		}
	}
	return s
}

func versionReferenceMatchesBaseline(raw, known string) bool {
	a, e := url.Parse(raw)
	if e != nil {
		return false
	}
	b, e := url.Parse(known)
	if e != nil || versionOrigin(a) != versionOrigin(b) {
		return false
	}
	normalize := func(p string) string {
		all := apiVersionInURL.FindAllStringSubmatchIndex(p, -1)
		if len(all) > 0 {
			m := all[len(all)-1]
			return p[:m[2]] + "{version}" + p[m[3]:]
		}
		return p
	}
	return normalize(a.Path) == normalize(b.Path)
}

func validateVersionRequest(c *Client, req APIVersionDiscoveryRequest, opts APIVersionDiscoveryOptions) ([]APIVersionOfficialScope, error) {
	data, err := json.Marshal(req)
	if err != nil || len(data) > 64<<10 {
		return nil, fmt.Errorf("API version request exceeds its metadata budget")
	}
	if req.SchemaVersion != APIVersionDiscoverySchemaVersion {
		return nil, fmt.Errorf("unsupported API version request schema %q", req.SchemaVersion)
	}
	digest, err := hex.DecodeString(req.Known.SHA256)
	if err != nil || len(digest) != sha256.Size || req.Known.Version == "" || (req.Known.ProviderKey == "" && req.Known.SourceURL == "") {
		return nil, fmt.Errorf("known version, SHA256 and provider key or source URL are required")
	}
	if req.Known.SourceKind == "" {
		return nil, fmt.Errorf("known source_kind is required")
	}
	switch req.Known.SourceKind {
	case OperationSourceOpenAPI, OperationSourceGoogleDiscovery, OperationSourceAWSSmithy, OperationSourceAsyncAPI, OperationSourceGraphQL, OperationSourceOpenRPC, OperationSourceGRPCProtobuf, OperationSourceOData:
	default:
		return nil, fmt.Errorf("unknown source_kind")
	}
	if req.Contract != nil {
		_, _, _, diagnostics := prepareStepContract(*req.Contract, DefaultPromptBudget())
		if errors := errorDiagnostics(diagnostics); len(errors) > 0 {
			return nil, DiagnosticError{Diagnostics: errors}
		}
	}
	if opts.Timeout < 0 || opts.Timeout > DefaultAPIVersionTimeout || opts.AdapterTimeout < 0 || opts.AdapterTimeout > 60*1e9 {
		return nil, fmt.Errorf("invalid API version timeout")
	}
	if len(req.OfficialScopes) > 8 || len(req.Locators) > 8 {
		return nil, fmt.Errorf("at most eight official scopes and locators are supported")
	}
	for _, v := range []string{req.Known.Version, req.Known.ETag, req.Known.LastModified} {
		_, changed := sanitizePromptString(v, 256)
		if len(v) > 256 || changed {
			return nil, fmt.Errorf("invalid baseline text or validator")
		}
	}
	for _, raw := range []string{req.Known.SourceURL, req.DocsURL} {
		if raw != "" {
			if _, err := c.versionURLSyntax(raw); err != nil {
				return nil, err
			}
		}
	}
	for _, l := range req.Locators {
		switch l.Kind {
		case "direct", "pointer", "pattern", "github", "template", "gated":
		default:
			return nil, fmt.Errorf("unsupported locator kind %q", l.Kind)
		}
		if l.URL != "" {
			if _, err := c.versionURLSyntax(l.URL); err != nil {
				return nil, err
			}
		}
		if l.URLTemplate != "" {
			if _, err := c.versionURLSyntax(strings.ReplaceAll(l.URLTemplate, "{version}", "1")); err != nil {
				return nil, err
			}
		}
		if l.Kind == "pattern" && (strings.Count(l.URLTemplate, "{version}") != 1 || l.BaselineToken == "") {
			return nil, fmt.Errorf("pattern locator requires one version slot and baseline token")
		}
		if l.Kind == "github" {
			p := strings.Split(l.Repository, "/")
			if len(p) != 2 || p[0] == "" || p[1] == "" || strings.ContainsAny(l.Repository+l.Ref, "?#\\") || strings.Contains(l.Repository, "..") {
				return nil, fmt.Errorf("invalid GitHub repository")
			}
			if l.Ref == "" || l.PathPattern == "" || strings.Contains(l.PathPattern, "..") {
				return nil, fmt.Errorf("GitHub ref and local path pattern are required")
			}
		}
	}
	scopes := append([]APIVersionOfficialScope(nil), req.OfficialScopes...)
	for _, s := range scopes {
		u, err := c.versionURLSyntax(s.Origin)
		if err != nil {
			return nil, err
		}
		if (u.Path != "" && u.Path != "/") || u.RawQuery != "" || u.Fragment != "" || !strings.HasPrefix(s.PathPrefix, "/") {
			return nil, fmt.Errorf("invalid official origin/path scope")
		}
		candidate := *u
		candidate.Path = s.PathPrefix
		if _, err := c.versionURLSyntax(candidate.String()); err != nil {
			return nil, err
		}
		if (u.Hostname() == "github.com" || u.Hostname() == "raw.githubusercontent.com") && (s.Repository == "" || !strings.HasPrefix(s.PathPrefix, "/"+s.Repository+"/")) {
			return nil, fmt.Errorf("shared GitHub origin requires repository/path scope")
		}
	}
	cat := catalog.BuiltInCatalog()
	if opts.Catalog != nil {
		cat = *opts.Catalog
	}
	if req.Known.ProviderKey != "" {
		if provider, ok := cat.FindProvider(req.Known.ProviderKey); ok {
			for _, ref := range provider.SpecReferences {
				if (ref.SourceAuthority == catalog.SourceAuthorityOfficialProvider || ref.SourceAuthority == catalog.SourceAuthorityOfficialGitHub || ref.SourceAuthority == catalog.SourceAuthorityOfficialDocs) && (req.Known.SourceURL == "" || versionReferenceMatchesBaseline(ref.URL, req.Known.SourceURL)) {
					if _, e := c.versionURLSyntax(ref.URL); e == nil {
						scopes = append(scopes, scopeForVersionURL(ref.URL))
					}
				}
			}
		}
	}
	return scopes, nil
}
