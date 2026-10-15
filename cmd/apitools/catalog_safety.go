package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/artifactio"
	"github.com/OpenUdon/apitools/sqlitecache"

	"github.com/OpenUdon/apitools/internal/sourceguard"
)

func safeCatalogRelative(value string) bool {
	if !filepath.IsLocal(value) || value == "." || strings.ContainsAny(value, "\\\x00\r\n\t") {
		return false
	}
	for _, part := range strings.Split(filepath.ToSlash(value), "/") {
		if part == "" || part == "." || part == ".." {
			return false
		}
	}
	return true
}

func safeCatalogID(value string) bool {
	if value == "" || value == "." || value == ".." {
		return false
	}
	for i, r := range value {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			continue
		}
		if i > 0 && (r == '-' || r == '_' || r == '.') {
			continue
		}
		return false
	}
	return true
}

// Validate every path before the first mkdir/database open/copy. Missing paths
// are permitted; existing symlinks and special files at/below root are refused.
func validateCatalogMutationPaths(opts catalog.RootOptions, files, dirs []string) error {
	reg, idx := opts.RegistryPath, opts.IndexPath
	if reg == "" {
		reg = catalog.DefaultRegistryPath
	}
	if idx == "" {
		idx = catalog.DefaultOperationIndexPath
	}
	files = append([]string{reg, idx, reg + "-wal", reg + "-shm", reg + "-journal"}, files...)
	for i, file := range files {
		if !safeCatalogRelative(file) {
			return fmt.Errorf("invalid catalog file path %q", file)
		}
		for _, prev := range files[:i] {
			if file == prev || strings.HasPrefix(file, prev+string(filepath.Separator)) || strings.HasPrefix(prev, file+string(filepath.Separator)) {
				return fmt.Errorf("catalog file paths overlap: %q and %q", prev, file)
			}
		}
	}
	root, err := filepath.Abs(opts.Directory)
	if err != nil {
		return err
	}
	if filepath.Dir(root) == root {
		return fmt.Errorf("catalog root must not be a filesystem root")
	}
	if info, err := os.Lstat(root); err == nil {
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("catalog root must be a regular directory")
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	for _, dir := range dirs {
		for _, file := range files {
			if dir == file || strings.HasPrefix(dir, file+string(filepath.Separator)) {
				return fmt.Errorf("catalog directory %q overlaps file %q", dir, file)
			}
		}
	}
	for _, dir := range dirs {
		if !safeCatalogRelative(dir) {
			return fmt.Errorf("invalid catalog directory %q", dir)
		}
	}
	for i, rel := range append(files, dirs...) {
		parts := strings.Split(filepath.ToSlash(rel), "/")
		for j := range parts {
			path := filepath.Join(root, filepath.Join(parts[:j+1]...))
			info, err := os.Lstat(path)
			if os.IsNotExist(err) {
				break
			}
			if err != nil {
				return err
			}
			wantDir := j < len(parts)-1 || i >= len(files)
			if info.Mode()&os.ModeSymlink != 0 || wantDir && !info.IsDir() || !wantDir && !info.Mode().IsRegular() {
				return fmt.Errorf("unsafe catalog path %q", path)
			}
		}
	}
	return nil
}

func storeCatalogSource(ctx context.Context, paths catalog.RootPaths, data []byte, artifact sqlitecache.CatalogArtifact) error {
	if err := validateCatalogMutationPaths(catalog.RootOptions{Directory: paths.Directory, RegistryPath: relativeCatalogPath(paths.Directory, paths.RegistryPath), IndexPath: relativeCatalogPath(paths.Directory, paths.IndexPath)}, []string{artifact.Path}, nil); err != nil {
		return err
	}
	cache, err := sqlitecache.OpenCatalog(catalog.RootOptions{Directory: paths.Directory, RegistryPath: relativeCatalogPath(paths.Directory, paths.RegistryPath), IndexPath: relativeCatalogPath(paths.Directory, paths.IndexPath)})
	if err != nil {
		return err
	}
	defer cache.Close()
	result, err := artifactio.WriteFile(paths.Directory, artifact.Path, data, artifactio.WriteOptions{Mode: 0644})
	if err != nil {
		return err
	}
	if err := cache.StoreCatalogArtifact(ctx, artifact); err != nil {
		if !result.Reused {
			if removeErr := artifactio.RemoveArtifact(paths.Directory, artifact.Path); removeErr != nil {
				return fmt.Errorf("%w; rollback: %v", err, removeErr)
			}
		}
		return err
	}
	return nil
}

func relativeCatalogPath(root, path string) string { rel, _ := filepath.Rel(root, path); return rel }

// Flag defaults must not mask environment configuration. Explicit empty path
// flags select the built-in default; an explicit empty root stays unconfigured.
func catalogCLIOptions(fs *flag.FlagSet, options catalog.RootOptions) catalog.RootOptions {
	resolved := catalog.RootOptionsWithEnvironment(options)
	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "root", "dir":
			resolved.Directory = strings.TrimSpace(options.Directory)
		case "registry":
			resolved.RegistryPath = strings.TrimSpace(options.RegistryPath)
		case "index":
			resolved.IndexPath = strings.TrimSpace(options.IndexPath)
		}
	})
	if resolved.RegistryPath == "" {
		resolved.RegistryPath = catalog.DefaultRegistryPath
	}
	if resolved.IndexPath == "" {
		resolved.IndexPath = catalog.DefaultOperationIndexPath
	}
	return resolved
}

// All lifecycle commands share the same confined, bounded catalog decoder.
func readCatalogMetadata(root, relative string) (*catalog.Catalog, error) {
	if relative == "" {
		return nil, nil
	}
	if strings.TrimSpace(root) == "" {
		return nil, fmt.Errorf("custom catalog requires an explicit root")
	}
	if !safeCatalogRelative(relative) {
		return nil, fmt.Errorf("invalid custom catalog path")
	}
	file, err := artifactio.ReadFile(root, relative, artifactio.ReadOptions{MaxBytes: 4 << 20})
	if err != nil {
		return nil, fmt.Errorf("cannot read custom catalog metadata: %w", err)
	}
	if err := sourceguard.CheckJSON("custom catalog", file.Data); err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(file.Data))
	dec.DisallowUnknownFields()
	var cat *catalog.Catalog
	if err := dec.Decode(&cat); err != nil {
		return nil, fmt.Errorf("invalid custom catalog metadata: %w", err)
	}
	if cat == nil {
		return nil, fmt.Errorf("custom catalog must be an object")
	}
	var trailing any
	if err := dec.Decode(&trailing); err != io.EOF {
		return nil, fmt.Errorf("trailing custom catalog metadata")
	}
	if err := cat.Validate(); err != nil {
		return nil, err
	}
	return cat, nil
}

func catalogIdentity(cat *catalog.Catalog, provider, specID, kind string) error {
	if !safeCatalogID(provider) || !safeCatalogID(specID) {
		return fmt.Errorf("provider and spec IDs must be safe catalog identifiers")
	}
	selected := catalog.BuiltInCatalog()
	if cat != nil {
		selected = *cat
	}
	p, ok := selected.FindProvider(provider)
	if !ok || p.ID != provider {
		return fmt.Errorf("unknown provider %q in selected catalog; declare custom identities with --catalog", provider)
	}
	for _, ref := range p.SpecReferences {
		if ref.ID != specID {
			continue
		}
		normalized := kind
		switch kind {
		case "swagger2":
			normalized = "openapi"
		case "aws-smithy":
			normalized = "smithy-json"
		}
		if normalized == string(ref.Kind) || kind == "advisory-overlay" && (ref.Kind == catalog.SpecKindHumanDocs || ref.Kind == catalog.SpecKindOpenAPI) {
			return nil
		}
		return fmt.Errorf("artifact kind %q disagrees with selected catalog reference %q", kind, specID)
	}
	return fmt.Errorf("unknown spec %q for provider %q in selected catalog; declare custom identities with --catalog", specID, provider)
}

func nonemptyCatalogFiles(paths ...string) []string {
	var files []string
	for _, path := range paths {
		if path != "" {
			files = append(files, path)
		}
	}
	return files
}
