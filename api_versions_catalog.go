package apitools

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/sourceguard"
)

func versionRecordID(raw string) string {
	digest := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(digest[:])[:24]
}

func versionCatalogURL(c *Client, req APIVersionDiscoveryRequest) string {
	base := c.APIsGuruListURL
	if base != "" && base != DefaultAPIsGuruListURL {
		return base
	}
	if req.Known.ProviderKey != "" {
		provider := strings.SplitN(req.Known.ProviderKey, ":", 2)[0]
		if !strings.ContainsAny(provider, "/?#%\\") {
			return "https://api.apis.guru/v2/" + url.PathEscape(provider) + ".json"
		}
	}
	return DefaultAPIsGuruListURL
}

func versionCatalogReferences(req APIVersionDiscoveryRequest, opts APIVersionDiscoveryOptions) []APIVersionRecord {
	cat := catalog.BuiltInCatalog()
	if opts.Catalog != nil {
		cat = *opts.Catalog
	}
	provider, ok := cat.FindProvider(req.Known.ProviderKey)
	if !ok {
		return nil
	}
	var out []APIVersionRecord
	for _, ref := range provider.SpecReferences {
		kind := "direct"
		if ref.Kind == catalog.SpecKindHumanDocs || ref.Kind == catalog.SpecKindOpenAPIIndex {
			kind = "pointer"
		}
		l := APIVersionLocator{Kind: kind, URL: ref.URL}
		out = append(out, APIVersionRecord{ID: versionRecordID(ref.URL), Locator: l, Recipe: l, SourceURL: publicVersionURL(ref.URL), VersionClaims: APIVersionClaims{URLToken: ref.Version}, Evidence: "catalog-highest", Validation: "unexamined", Comparison: "unexamined", LocatorStatus: "ok"})
	}
	return out
}

func versionListRecords(ctx context.Context, data []byte, req APIVersionDiscoveryRequest) ([]APIVersionRecord, error) {
	if err := sourceguard.CheckJSONWithLimits(ctx, "api-version-list", data, sourceguard.DefaultLimits()); err != nil {
		return nil, err
	}
	var list map[string]struct {
		Preferred string `json:"preferred"`
		Versions  map[string]struct {
			SwaggerURL     string          `json:"swaggerUrl"`
			SwaggerYAMLURL string          `json:"swaggerYamlUrl"`
			Updated        time.Time       `json:"updated"`
			Info           json.RawMessage `json:"info"`
		} `json:"versions"`
	}
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	if len(list) > 10000 {
		return nil, fmt.Errorf("directory entry budget exceeded")
	}
	var keys []string
	for k := range list {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var out []APIVersionRecord
	for _, key := range keys {
		entry := list[key]
		if req.Known.ProviderKey != "" && key != req.Known.ProviderKey {
			continue
		}
		var versions []string
		for k := range entry.Versions {
			versions = append(versions, k)
		}
		sort.Slice(versions, func(i, j int) bool {
			if n, ok := compareAPIVersions(versions[i], versions[j]); ok && n != 0 {
				return n > 0
			}
			return versions[i] < versions[j]
		})
		if len(versions) > 100 {
			return nil, fmt.Errorf("provider version budget exceeded")
		}
		for _, v := range versions {
			e := entry.Versions[v]
			var info struct {
				Origin json.RawMessage `json:"x-origin"`
			}
			if err := json.Unmarshal(e.Info, &info); err != nil {
				continue
			}
			var origins []struct {
				URL string `json:"url"`
			}
			if len(info.Origin) > 0 && info.Origin[0] == '[' {
				_ = json.Unmarshal(info.Origin, &origins)
			} else {
				var origin struct {
					URL string `json:"url"`
				}
				_ = json.Unmarshal(info.Origin, &origin)
				origins = append(origins, origin)
			}
			raw := firstNonEmpty(e.SwaggerURL, e.SwaggerYAMLURL)
			if len(origins) > 0 && origins[0].URL != "" {
				raw = origins[0].URL
			}
			if req.Known.ProviderKey == "" && raw != req.Known.SourceURL {
				continue
			}
			if len(raw) > 2048 {
				continue
			}
			l := APIVersionLocator{Kind: "direct", URL: raw}
			evidence := "catalog-highest"
			if v == entry.Preferred {
				evidence = "catalog-preferred"
			}
			var updated *time.Time
			if !e.Updated.IsZero() {
				t := e.Updated
				updated = &t
			}
			out = append(out, APIVersionRecord{ID: versionRecordID(raw + "|" + v), Locator: l, Recipe: l, Evidence: evidence, SourceURL: publicVersionURL(raw), CatalogUpdated: updated, VersionClaims: APIVersionClaims{URLToken: v}, Validation: "unexamined", Comparison: "unexamined", LocatorStatus: "ok"})
		}
	}
	return out, nil
}

func readVersionList(file string) (apiVersionListCache, error) {
	data, err := readVersionMetadata(file, DefaultMaxBytes)
	if err != nil {
		return apiVersionListCache{}, err
	}
	var envelope map[string]json.RawMessage
	if err = json.Unmarshal(data, &envelope); err != nil {
		return apiVersionListCache{}, err
	}
	if _, ok := envelope["schema_version"]; ok {
		var cache apiVersionListCache
		if err = decodeVersionJSON(data, &cache, DefaultMaxBytes); err != nil {
			return cache, err
		}
		if cache.SchemaVersion != "apitools.api-version-list-cache/v1" {
			return cache, fmt.Errorf("unsupported directory cache schema")
		}
		return cache, nil
	}
	info, err := os.Lstat(file)
	if err != nil {
		return apiVersionListCache{}, err
	}
	return apiVersionListCache{SchemaVersion: "apitools.api-version-list-cache/v1", CheckedAt: info.ModTime().UTC(), Data: data}, nil
}
