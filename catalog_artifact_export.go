package apitools

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/artifactio"
)

// CatalogArtifactReference binds a selected operation or metadata lead to one
// registered source. Selector remains native and is checked during downstream
// binding, not converted into a display method/path by this copying API.
type CatalogArtifactReference struct {
	CatalogSHA256 string              `json:"catalog_sha256"`
	ProviderID    string              `json:"provider_id"`
	SpecRefID     string              `json:"spec_ref_id"`
	ArtifactID    string              `json:"artifact_id"`
	Kind          OperationSourceKind `json:"kind"`
	SHA256        string              `json:"sha256"`
	Bytes         int64               `json:"bytes"`
	Selector      string              `json:"selector,omitempty"`
}

type CatalogArtifactExportOptions struct {
	Index       CatalogIndexOptions
	References  []CatalogArtifactReference
	WorkflowDir string
	ArtifactDir string
	Force       bool
}

type CatalogArtifactExportEntry struct {
	Reference      CatalogArtifactReference `json:"reference"`
	Path           string                   `json:"path"`
	RegisteredKind string                   `json:"registered_kind"`
	Advisory       bool                     `json:"advisory,omitempty"`
	SourceURL      string                   `json:"source_url,omitempty"`
}

type CatalogArtifactExportOverlay struct {
	ProviderID string `json:"provider_id"`
	SpecRefID  string `json:"spec_ref_id,omitempty"`
	OverlayID  string `json:"overlay_id"`
	Path       string `json:"path"`
	SHA256     string `json:"sha256"`
}

// Paths in this report are relative to ArtifactDir, never installation paths.
type CatalogArtifactExportReport struct {
	SchemaVersion    string                         `json:"schema_version"`
	CatalogSHA256    string                         `json:"catalog_sha256"`
	Artifacts        []CatalogArtifactExportEntry   `json:"artifacts"`
	SecurityOverlays []CatalogArtifactExportOverlay `json:"security_overlays,omitempty"`
	Reused           bool                           `json:"reused,omitempty"`
}

// ExportCatalogArtifacts copies only selected digest-bound registrations with
// provider-wide and matching spec-scoped security overlays. It never downloads,
// executes, applies auth overlays to source bytes, or verifies native selectors
// on behalf of the downstream binder. All sources pass before atomic publication.
func ExportCatalogArtifacts(ctx context.Context, options CatalogArtifactExportOptions) (CatalogArtifactExportReport, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithTimeout(ctx, CatalogIndexDeadline)
	defer cancel()
	if len(options.References) == 0 || len(options.References) > 64 {
		return CatalogArtifactExportReport{}, fmt.Errorf("select between 1 and 64 catalog artifact references")
	}
	paths, cat, err := catalogIndexConfiguration(options.Index)
	if err != nil {
		return CatalogArtifactExportReport{}, err
	}
	digest, err := catalogIndexIdentity(cat, options.Index.Catalog == nil)
	if err != nil {
		return CatalogArtifactExportReport{}, err
	}
	rows, err := options.Index.ReadRegistrations(ctx, options.Index.Root)
	if err != nil {
		return CatalogArtifactExportReport{}, err
	}
	registrationDigest, err := catalogRegistrationIdentity(rows)
	if err != nil {
		return CatalogArtifactExportReport{}, err
	}
	registered := map[string]catalog.CatalogSpecArtifact{}
	for _, row := range rows {
		registered[row.ProviderID+"\x00"+row.ArtifactID] = row
	}
	if strings.TrimSpace(options.WorkflowDir) == "" {
		return CatalogArtifactExportReport{}, fmt.Errorf("an existing workflow directory is required")
	}
	workflow, err := filepath.Abs(options.WorkflowDir)
	if err != nil {
		return CatalogArtifactExportReport{}, err
	}
	// Resolve ancestors consistently with artifactio; the chosen directory itself
	// may not be a symlink. Refuse overlap with the source root even with Force.
	info, err := os.Lstat(workflow)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return CatalogArtifactExportReport{}, fmt.Errorf("workflow root must be an existing regular directory")
	}
	workflow, err = filepath.EvalSymlinks(workflow)
	if err != nil {
		return CatalogArtifactExportReport{}, err
	}
	directory := options.ArtifactDir
	if directory == "" {
		directory = "api-artifacts"
	}
	if !catalogIndexRelativePath(directory) {
		return CatalogArtifactExportReport{}, fmt.Errorf("artifact directory must be a confined relative path")
	}
	if err := catalogExportAncestors(workflow, directory); err != nil {
		return CatalogArtifactExportReport{}, err
	}
	target := filepath.Join(workflow, filepath.FromSlash(directory))
	if target == paths.Directory || strings.HasPrefix(target, paths.Directory+string(filepath.Separator)) || strings.HasPrefix(paths.Directory, target+string(filepath.Separator)) {
		return CatalogArtifactExportReport{}, fmt.Errorf("export destination overlaps catalog source root")
	}
	tx, err := artifactio.BeginDir(target, options.Force)
	if err != nil {
		return CatalogArtifactExportReport{}, err
	}
	defer tx.Rollback()
	report := CatalogArtifactExportReport{SchemaVersion: "apitools.catalog-artifact-export/v1", CatalogSHA256: digest, Artifacts: []CatalogArtifactExportEntry{}}
	refs := append([]CatalogArtifactReference(nil), options.References...)
	sort.Slice(refs, func(i, j int) bool {
		return catalogArtifactReferenceKey(refs[i]) < catalogArtifactReferenceKey(refs[j])
	})
	seenRefs := map[string]bool{}
	savedSources := map[string]string{}
	selected := map[string]map[string]bool{}
	var total int64
	for _, ref := range refs {
		if err := ctx.Err(); err != nil {
			return CatalogArtifactExportReport{}, err
		}
		key := catalogArtifactReferenceKey(ref)
		if seenRefs[key] {
			continue
		}
		seenRefs[key] = true
		if ref.CatalogSHA256 != digest || !catalogIndexDigest(ref.SHA256) || ref.Bytes <= 0 || ref.Bytes > maxCatalogArtifactBytes || !validOperationSourceKind(ref.Kind) {
			return CatalogArtifactExportReport{}, fmt.Errorf("invalid or stale catalog artifact reference")
		}
		if value, changed := sanitizePromptString(ref.Selector, DefaultPromptTextRunes); changed || value != ref.Selector {
			return CatalogArtifactExportReport{}, fmt.Errorf("unsafe catalog operation selector")
		}
		provider, ok := cat.FindProvider(ref.ProviderID)
		row, registeredOK := registered[ref.ProviderID+"\x00"+ref.ArtifactID]
		if !ok || provider.ID != ref.ProviderID || !registeredOK || row.SpecRefID != ref.SpecRefID || row.SHA256 != ref.SHA256 || row.Bytes != ref.Bytes || catalogRegisteredSourceKind(row.Kind) != ref.Kind {
			return CatalogArtifactExportReport{}, fmt.Errorf("catalog artifact registration binding mismatch")
		}
		refFound := false
		for _, spec := range provider.SpecReferences {
			if spec.ID == ref.SpecRefID {
				refFound = true
				if row.Kind != "advisory-overlay" && catalogRegisteredSourceKind(string(spec.Kind)) != ref.Kind {
					return CatalogArtifactExportReport{}, fmt.Errorf("catalog source family mismatch")
				}
			}
		}
		if !refFound && row.Kind != "advisory-overlay" {
			return CatalogArtifactExportReport{}, fmt.Errorf("catalog artifact has no spec reference")
		}
		// Every selected registration is verified, even shared links with different
		// source paths. Physical output copies deduplicate by raw family+digest.
		file, err := artifactio.ReadFile(paths.Directory, row.Path, artifactio.ReadOptions{MaxBytes: maxCatalogArtifactBytes, SHA256: ref.SHA256, Bytes: ref.Bytes})
		if err != nil {
			return CatalogArtifactExportReport{}, fmt.Errorf("selected artifact failed integrity verification")
		}
		sourceKey := string(ref.Kind) + "\x00" + ref.SHA256
		relative, exists := savedSources[sourceKey]
		if !exists {
			total += file.Bytes
			if total > MaxCatalogIndexBytes {
				return CatalogArtifactExportReport{}, fmt.Errorf("selected export exceeds aggregate byte budget")
			}
			// Family retains native source shape. The suffix only assists local
			// loaders; content and selector remain unmodified.
			ext := strings.ToLower(filepath.Ext(row.Path))
			switch ext {
			case ".json", ".yaml", ".yml", ".graphql", ".gql", ".proto", ".xml":
			default:
				ext = ".source"
			}
			relative = filepath.ToSlash(filepath.Join(string(ref.Kind), ref.SHA256+ext))
			if _, err := tx.Stage(relative, file.Data, 0o600); err != nil {
				return CatalogArtifactExportReport{}, err
			}
			savedSources[sourceKey] = relative
		}
		provenanceURL, _, err := sanitizeOperationSourceURL(row.SourceURL)
		if err != nil {
			return CatalogArtifactExportReport{}, fmt.Errorf("invalid registered provenance URL")
		}
		report.Artifacts = append(report.Artifacts, CatalogArtifactExportEntry{Reference: ref, Path: relative, RegisteredKind: row.Kind, Advisory: row.Kind == "advisory-overlay", SourceURL: provenanceURL})
		if selected[ref.ProviderID] == nil {
			selected[ref.ProviderID] = map[string]bool{}
		}
		selected[ref.ProviderID][ref.SpecRefID] = true
	}
	for _, overlay := range cat.ListSecurityOverlays() {
		if err := ctx.Err(); err != nil {
			return CatalogArtifactExportReport{}, err
		}
		if specs := selected[overlay.ProviderID]; specs == nil || (overlay.SpecRefID != "" && !specs[overlay.SpecRefID]) {
			continue
		}
		data, err := json.MarshalIndent(overlay, "", "  ")
		if err != nil {
			return CatalogArtifactExportReport{}, err
		}
		data = append(data, '\n')
		total += int64(len(data))
		if total > MaxCatalogIndexBytes {
			return CatalogArtifactExportReport{}, fmt.Errorf("selected export exceeds aggregate byte budget")
		}
		relative := filepath.ToSlash(filepath.Join("security-overlays", overlay.ProviderID, overlay.ID+".json"))
		if _, err := tx.Stage(relative, data, 0o600); err != nil {
			return CatalogArtifactExportReport{}, err
		}
		hash := sha256.Sum256(data)
		report.SecurityOverlays = append(report.SecurityOverlays, CatalogArtifactExportOverlay{ProviderID: overlay.ProviderID, SpecRefID: overlay.SpecRefID, OverlayID: overlay.ID, Path: relative, SHA256: hex.EncodeToString(hash[:])})
	}
	manifest, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return CatalogArtifactExportReport{}, err
	}
	manifest = append(manifest, '\n')
	if total+int64(len(manifest)) > MaxCatalogIndexBytes {
		return CatalogArtifactExportReport{}, fmt.Errorf("selected export exceeds aggregate byte budget")
	}
	if _, err := tx.Stage("provenance.json", manifest, 0o600); err != nil {
		return CatalogArtifactExportReport{}, err
	}
	current, err := options.Index.ReadRegistrations(ctx, options.Index.Root)
	if err != nil {
		return CatalogArtifactExportReport{}, err
	}
	currentDigest, err := catalogRegistrationIdentity(current)
	if err != nil || currentDigest != registrationDigest {
		return CatalogArtifactExportReport{}, fmt.Errorf("catalog registrations changed during export")
	}
	if err := ctx.Err(); err != nil {
		return CatalogArtifactExportReport{}, err
	}
	reused, err := tx.Commit()
	if err != nil {
		return CatalogArtifactExportReport{}, err
	}
	report.Reused = reused
	return report, nil
}

func catalogExportAncestors(workflow, directory string) error {
	root, err := os.OpenRoot(workflow)
	if err != nil {
		return err
	}
	defer root.Close()
	relative := ""
	for _, part := range strings.Split(filepath.FromSlash(directory), string(filepath.Separator)) {
		relative = filepath.Join(relative, part)
		info, err := root.Lstat(relative)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("export directory ancestor must be a regular directory")
		}
	}
	return nil
}

func catalogArtifactReferenceKey(ref CatalogArtifactReference) string {
	return ref.ProviderID + "\x00" + ref.SpecRefID + "\x00" + ref.ArtifactID + "\x00" + string(ref.Kind) + "\x00" + ref.SHA256 + "\x00" + ref.Selector + "\x00" + ref.CatalogSHA256 + fmt.Sprintf("/%d", ref.Bytes)
}
