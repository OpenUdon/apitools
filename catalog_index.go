package apitools

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/artifactio"
)

const (
	CatalogOperationIndexVersion = "apitools.catalog-operation-index/v1"
	CatalogIndexMetadataVersion  = "apitools.operation-candidates/v1/catalog-index.1"
	MaxCatalogIndexBytes         = 512 << 20
	MaxCatalogIndexOperations    = 500_000
	MaxCatalogIndexArtifacts     = 10_000
	CatalogIndexDeadline         = 3 * time.Minute
)

// CatalogRegistrationReader is an installation-selected read-only snapshot
// adapter. sqlitecache.ReadCatalogSpecArtifacts provides the SQLite adapter.
// Requests cannot supply this callback or change its root configuration.
type CatalogRegistrationReader func(context.Context, catalog.RootOptions) ([]catalog.CatalogSpecArtifact, error)

type CatalogIndexOptions struct {
	Root              catalog.RootOptions
	Catalog           *catalog.Catalog
	ReadRegistrations CatalogRegistrationReader
}

type CatalogArtifactLink struct {
	ProviderID     string `json:"provider_id"`
	SpecRefID      string `json:"spec_ref_id"`
	ArtifactID     string `json:"artifact_id"`
	RegisteredKind string `json:"registered_kind"`
	Advisory       bool   `json:"advisory,omitempty"`
}

// CatalogIndexedArtifact is one raw artifact with all catalog links. Partial
// metadata remains usable positive evidence, but its scope is unexamined.
type CatalogIndexedArtifact struct {
	ArtifactID     string                `json:"artifact_id"`
	Kind           OperationSourceKind   `json:"kind"`
	SHA256         string                `json:"sha256"`
	Bytes          int64                 `json:"bytes"`
	Path           string                `json:"path"`
	Links          []CatalogArtifactLink `json:"links"`
	State          string                `json:"state"`
	Reason         string                `json:"reason,omitempty"`
	OperationCount int                   `json:"operation_count"`
	Operations     []OperationCandidate  `json:"operations,omitempty"`
	Diagnostics    []Diagnostic          `json:"diagnostics,omitempty"`
}

type CatalogIndexCoverage struct {
	ProviderID string `json:"provider_id"`
	SpecRefID  string `json:"spec_ref_id"`
	ArtifactID string `json:"artifact_id,omitempty"`
	State      string `json:"state"`
	Reason     string `json:"reason,omitempty"`
}

// CatalogOperationIndex is deterministic derived metadata, never approval or
// execution authority. It contains no absolute paths or observation times.
type CatalogOperationIndex struct {
	SchemaVersion      string                   `json:"schema_version"`
	MetadataVersion    string                   `json:"metadata_version"`
	CatalogSHA256      string                   `json:"catalog_sha256"`
	RegistrationSHA256 string                   `json:"registration_sha256"`
	Artifacts          []CatalogIndexedArtifact `json:"artifacts"`
	Coverage           []CatalogIndexCoverage   `json:"coverage"`
}

func catalogIndexConfiguration(options CatalogIndexOptions) (catalog.RootPaths, catalog.Catalog, error) {
	paths, err := catalog.ResolveRoot(options.Root)
	if err != nil {
		return paths, catalog.Catalog{}, err
	}
	if paths.Directory == "" {
		return paths, catalog.Catalog{}, fmt.Errorf("an explicit catalog root is required")
	}
	if options.ReadRegistrations == nil {
		return paths, catalog.Catalog{}, fmt.Errorf("a read-only catalog registration adapter is required")
	}
	cat := catalog.BuiltInCatalog()
	if options.Catalog != nil {
		cat = *options.Catalog
	}
	if err := cat.Validate(); err != nil {
		return paths, cat, err
	}
	return paths, cat, nil
}

func catalogIndexIdentity(cat catalog.Catalog, builtin bool) (string, error) {
	if builtin {
		return catalog.BuiltInCatalogManifest().CatalogSHA256, nil
	}
	cat = catalog.Catalog{Providers: cat.ListProviders(), SecurityOverlays: cat.ListSecurityOverlays()}
	for i := range cat.Providers {
		sort.Slice(cat.Providers[i].SpecReferences, func(a, b int) bool {
			return cat.Providers[i].SpecReferences[a].ID < cat.Providers[i].SpecReferences[b].ID
		})
	}
	data, err := json.Marshal(cat)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:]), nil
}

func catalogRegistrationIdentity(rows []catalog.CatalogSpecArtifact) (string, error) {
	if len(rows) > MaxCatalogIndexArtifacts {
		return "", fmt.Errorf("too many catalog registrations")
	}
	rows = append([]catalog.CatalogSpecArtifact(nil), rows...)
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].ProviderID != rows[j].ProviderID {
			return rows[i].ProviderID < rows[j].ProviderID
		}
		return rows[i].ArtifactID < rows[j].ArtifactID
	})
	seen := map[string]bool{}
	for _, row := range rows {
		key := row.ProviderID + "\x00" + row.ArtifactID
		if seen[key] || row.ProviderID == "" || row.SpecRefID == "" || row.ArtifactID == "" || row.Bytes <= 0 || len(row.SHA256) != 64 {
			return "", fmt.Errorf("invalid or duplicate catalog registration")
		}
		if _, err := hex.DecodeString(row.SHA256); err != nil || row.SHA256 != strings.ToLower(row.SHA256) {
			return "", fmt.Errorf("invalid catalog registration digest")
		}
		if !catalogIndexRelativePath(row.Path) {
			return "", fmt.Errorf("unsafe catalog artifact path")
		}
		for _, part := range strings.Split(filepath.FromSlash(row.Path), string(filepath.Separator)) {
			if part == ".." || part == "." || part == "" {
				return "", fmt.Errorf("unsafe catalog artifact path component")
			}
		}
		seen[key] = true
	}
	data, err := json.Marshal(rows)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:]), nil
}

func catalogRegisteredSourceKind(kind string) OperationSourceKind {
	switch kind {
	case "openapi", "swagger", "advisory-overlay":
		return OperationSourceOpenAPI
	case "smithy-json", "aws-smithy":
		return OperationSourceAWSSmithy
	default:
		return OperationSourceKind(kind)
	}
}

// BuildCatalogOperationIndex reads only registered local artifacts. It does
// not publish a file, fetch sources or change cache registrations. A changed
// registry at completion refuses the whole generation.
func BuildCatalogOperationIndex(ctx context.Context, options CatalogIndexOptions) (CatalogOperationIndex, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithTimeout(ctx, CatalogIndexDeadline)
	defer cancel()
	paths, cat, err := catalogIndexConfiguration(options)
	if err != nil {
		return CatalogOperationIndex{}, err
	}
	rows, err := options.ReadRegistrations(ctx, options.Root)
	if err != nil {
		return CatalogOperationIndex{}, err
	}
	registryDigest, err := catalogRegistrationIdentity(rows)
	if err != nil {
		return CatalogOperationIndex{}, err
	}
	if err := catalogIndexDestination(paths, rows); err != nil {
		return CatalogOperationIndex{}, err
	}
	catalogDigest, err := catalogIndexIdentity(cat, options.Catalog == nil)
	if err != nil {
		return CatalogOperationIndex{}, err
	}
	index := CatalogOperationIndex{SchemaVersion: CatalogOperationIndexVersion, MetadataVersion: CatalogIndexMetadataVersion, CatalogSHA256: catalogDigest, RegistrationSHA256: registryDigest, Artifacts: []CatalogIndexedArtifact{}, Coverage: []CatalogIndexCoverage{}}
	groups := map[string][]catalog.CatalogSpecArtifact{}
	covered := map[string]bool{}
	for _, row := range rows {
		provider, ok := cat.FindProvider(row.ProviderID)
		if !ok || provider.ID != row.ProviderID {
			return CatalogOperationIndex{}, fmt.Errorf("registration names an unknown catalog provider")
		}
		refFound := false
		for _, ref := range provider.SpecReferences {
			if ref.ID == row.SpecRefID {
				refFound = true
				if row.Kind != "advisory-overlay" && catalogRegisteredSourceKind(string(ref.Kind)) != catalogRegisteredSourceKind(row.Kind) {
					return CatalogOperationIndex{}, fmt.Errorf("registration source family conflicts with catalog reference")
				}
				break
			}
		}
		if !refFound && row.Kind != "advisory-overlay" {
			return CatalogOperationIndex{}, fmt.Errorf("registration has no catalog spec reference")
		}
		key := row.ArtifactID + "\x00" + string(catalogRegisteredSourceKind(row.Kind)) + "\x00" + row.SHA256 + fmt.Sprintf("/%d", row.Bytes)
		groups[key] = append(groups[key], row)
		covered[row.ProviderID+"\x00"+row.SpecRefID] = true
	}
	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	totalOperations, totalBytes := 0, 0
	for _, key := range keys {
		if err := ctx.Err(); err != nil {
			return CatalogOperationIndex{}, err
		}
		links := groups[key]
		sort.Slice(links, func(i, j int) bool {
			if links[i].Path != links[j].Path {
				return links[i].Path < links[j].Path
			}
			return links[i].ProviderID < links[j].ProviderID
		})
		first := links[0]
		artifact := CatalogIndexedArtifact{ArtifactID: first.ArtifactID, Kind: catalogRegisteredSourceKind(first.Kind), SHA256: first.SHA256, Bytes: first.Bytes, Path: first.Path, State: "indexed"}
		for _, row := range links {
			artifact.Links = append(artifact.Links, CatalogArtifactLink{ProviderID: row.ProviderID, SpecRefID: row.SpecRefID, ArtifactID: row.ArtifactID, RegisteredKind: row.Kind, Advisory: row.Kind == "advisory-overlay"})
		}
		sort.Slice(artifact.Links, func(i, j int) bool {
			if artifact.Links[i].ProviderID != artifact.Links[j].ProviderID {
				return artifact.Links[i].ProviderID < artifact.Links[j].ProviderID
			}
			return artifact.Links[i].SpecRefID < artifact.Links[j].SpecRefID
		})
		var result operationSourceAdapterResult
		var parseErr error
		switch {
		case first.Bytes > maxCatalogArtifactBytes:
			artifact.State = "oversize"
		case !validOperationSourceKind(artifact.Kind):
			artifact.State = "unsupported"
		case artifact.Kind != OperationSourceOpenAPI && first.Bytes > DefaultMaxBytes:
			artifact.State = "oversize"
		default:
			// Every link's registered file must have the expected raw identity.
			// Content is parsed once even when providers share it.
			for _, row := range links[1:] {
				if row.Path == first.Path {
					continue
				}
				_, parseErr = artifactio.ReadFile(paths.Directory, row.Path, artifactio.ReadOptions{MaxBytes: maxCatalogArtifactBytes, SHA256: row.SHA256, Bytes: row.Bytes})
				if parseErr != nil {
					break
				}
			}
			if parseErr == nil {
				first.Kind = string(artifact.Kind)
				result, parseErr = parseRegisteredCatalogArtifact(ctx, paths.Directory, first)
			}
			if parseErr != nil {
				switch {
				case errors.Is(parseErr, os.ErrNotExist):
					artifact.State = "missing"
				case strings.Contains(parseErr.Error(), "SHA256") || strings.Contains(parseErr.Error(), "bytes, want"):
					artifact.State = "digest_mismatch"
				case strings.Contains(parseErr.Error(), "budget") || strings.Contains(parseErr.Error(), "over limit") || strings.Contains(parseErr.Error(), "exceeds maximum"):
					artifact.State = "oversize"
				default:
					artifact.State = "parse_failure"
				}
			}
		}
		if err := ctx.Err(); err != nil {
			return CatalogOperationIndex{}, err
		}
		artifact.OperationCount = result.operationCount
		artifact.Operations = result.candidates
		artifact.Diagnostics = result.diagnostics
		if len(artifact.Operations) != artifact.OperationCount && artifact.State == "indexed" {
			artifact.State = "unsupported"
		}
		if artifact.State != "indexed" {
			artifact.Reason = "source or summaries were not completely indexed; scope remains unexamined"
		}
		sort.Slice(artifact.Operations, func(i, j int) bool {
			return artifact.Operations[i].Source.Selector < artifact.Operations[j].Source.Selector
		})
		sortDiagnostics(artifact.Diagnostics)
		for _, link := range artifact.Links {
			index.Coverage = append(index.Coverage, CatalogIndexCoverage{ProviderID: link.ProviderID, SpecRefID: link.SpecRefID, ArtifactID: link.ArtifactID, State: artifact.State, Reason: artifact.Reason})
		}
		totalOperations += len(artifact.Operations)
		data, err := json.Marshal(artifact)
		if err != nil {
			return CatalogOperationIndex{}, err
		}
		totalBytes += len(data)
		if totalOperations > MaxCatalogIndexOperations || totalBytes > MaxCatalogIndexBytes {
			return CatalogOperationIndex{}, fmt.Errorf("catalog index exceeds generation bounds")
		}
		index.Artifacts = append(index.Artifacts, artifact)
	}
	for _, provider := range cat.ListProviders() {
		for _, ref := range provider.SpecReferences {
			if !covered[provider.ID+"\x00"+ref.ID] {
				state := "missing"
				if ref.Kind == catalog.SpecKindHumanDocs {
					state = "unsupported"
				}
				index.Coverage = append(index.Coverage, CatalogIndexCoverage{ProviderID: provider.ID, SpecRefID: ref.ID, State: state, Reason: "no registered operation artifact"})
			}
		}
	}
	sortCatalogCoverage(index.Coverage)
	current, err := options.ReadRegistrations(ctx, options.Root)
	if err != nil {
		return CatalogOperationIndex{}, err
	}
	currentDigest, err := catalogRegistrationIdentity(current)
	if err != nil || currentDigest != registryDigest {
		return CatalogOperationIndex{}, fmt.Errorf("catalog registrations changed during indexing")
	}
	if err := ctx.Err(); err != nil {
		return CatalogOperationIndex{}, err
	}
	return index, nil
}

func sortCatalogCoverage(coverage []CatalogIndexCoverage) {
	sort.Slice(coverage, func(i, j int) bool {
		if coverage[i].ProviderID != coverage[j].ProviderID {
			return coverage[i].ProviderID < coverage[j].ProviderID
		}
		if coverage[i].SpecRefID != coverage[j].SpecRefID {
			return coverage[i].SpecRefID < coverage[j].SpecRefID
		}
		return coverage[i].ArtifactID < coverage[j].ArtifactID
	})
}

// WriteCatalogOperationIndex atomically replaces the configured index only
// after checking its generation against the current registration snapshot.
func WriteCatalogOperationIndex(ctx context.Context, options CatalogIndexOptions, index CatalogOperationIndex) error {
	if ctx == nil {
		ctx = context.Background()
	}
	paths, cat, err := catalogIndexConfiguration(options)
	if err != nil {
		return err
	}
	if err := validateCatalogIndexHeader(index, cat, options.Catalog == nil); err != nil {
		return err
	}
	rows, err := options.ReadRegistrations(ctx, options.Root)
	if err != nil {
		return err
	}
	if err := catalogIndexDestination(paths, rows); err != nil {
		return err
	}
	if err := validateCatalogIndexBody(index, cat, rows, true); err != nil {
		return err
	}
	digest, err := catalogRegistrationIdentity(rows)
	if err != nil || digest != index.RegistrationSHA256 {
		return fmt.Errorf("catalog registrations changed before publication")
	}
	data, err := json.Marshal(index)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if len(data) > MaxCatalogIndexBytes {
		return fmt.Errorf("catalog index exceeds byte bound")
	}
	if err := checkCatalogIndexJSON(ctx, data); err != nil {
		return err
	}
	relative, err := filepath.Rel(paths.Directory, paths.IndexPath)
	if err != nil {
		return err
	}
	_, err = artifactio.WriteFile(paths.Directory, relative, data, artifactio.WriteOptions{Mode: 0o600, Force: true})
	return err
}

func validateCatalogIndexHeader(index CatalogOperationIndex, cat catalog.Catalog, builtin bool) error {
	digest, err := catalogIndexIdentity(cat, builtin)
	if err != nil {
		return err
	}
	if index.SchemaVersion != CatalogOperationIndexVersion || index.MetadataVersion != CatalogIndexMetadataVersion || index.CatalogSHA256 != digest || !catalogIndexDigest(index.RegistrationSHA256) {
		return fmt.Errorf("unsupported or mismatched catalog index identity")
	}
	return nil
}

// ReadCatalogOperationIndex never parses provider artifacts. A changed registry
// marks all saved coverage stale and therefore unexamined; it does not rewrite
// the saved index or manufacture a no-match decision.
func ReadCatalogOperationIndex(ctx context.Context, options CatalogIndexOptions) (CatalogOperationIndex, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	paths, cat, err := catalogIndexConfiguration(options)
	if err != nil {
		return CatalogOperationIndex{}, err
	}
	relative, err := filepath.Rel(paths.Directory, paths.IndexPath)
	if err != nil {
		return CatalogOperationIndex{}, err
	}
	file, err := artifactio.ReadFile(paths.Directory, relative, artifactio.ReadOptions{MaxBytes: MaxCatalogIndexBytes})
	if err != nil {
		return CatalogOperationIndex{}, err
	}
	if err := checkCatalogIndexJSON(ctx, file.Data); err != nil {
		return CatalogOperationIndex{}, err
	}
	var index CatalogOperationIndex
	decoder := json.NewDecoder(bytes.NewReader(file.Data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&index); err != nil {
		return CatalogOperationIndex{}, fmt.Errorf("invalid catalog operation index")
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return CatalogOperationIndex{}, fmt.Errorf("trailing catalog operation index data")
	}
	if err := validateCatalogIndexHeader(index, cat, options.Catalog == nil); err != nil {
		return CatalogOperationIndex{}, err
	}
	rows, err := options.ReadRegistrations(ctx, options.Root)
	if err != nil {
		return CatalogOperationIndex{}, err
	}
	digest, err := catalogRegistrationIdentity(rows)
	if err != nil {
		return CatalogOperationIndex{}, err
	}
	if err := validateCatalogIndexBody(index, cat, rows, digest == index.RegistrationSHA256); err != nil {
		return CatalogOperationIndex{}, err
	}
	if digest != index.RegistrationSHA256 {
		for i := range index.Coverage {
			index.Coverage[i].State = "stale"
			index.Coverage[i].Reason = "registration snapshot changed; rebuild the index"
		}
		for i := range index.Artifacts {
			index.Artifacts[i].State = "stale"
			index.Artifacts[i].Operations = nil
		}
	}
	if err := ctx.Err(); err != nil {
		return CatalogOperationIndex{}, err
	}
	return index, nil
}

func catalogIndexDestination(paths catalog.RootPaths, rows []catalog.CatalogSpecArtifact) error {
	for _, row := range rows {
		for _, relative := range []string{row.Path, row.OverlayPath, row.BuilderPath} {
			if relative == "" {
				continue
			}
			if !catalogIndexRelativePath(relative) {
				return fmt.Errorf("unsafe registered source path")
			}
			source := filepath.Join(paths.Directory, filepath.FromSlash(relative))
			if source == paths.RegistryPath || source == paths.IndexPath || strings.HasPrefix(source, paths.IndexPath+string(filepath.Separator)) || strings.HasPrefix(paths.IndexPath, source+string(filepath.Separator)) {
				return fmt.Errorf("index or registry overlaps registered source content")
			}
		}
	}
	return nil
}
