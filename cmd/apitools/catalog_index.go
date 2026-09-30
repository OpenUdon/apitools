package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/artifactio"
	"github.com/OpenUdon/apitools/sqlitecache"
)

func runCatalogIndex(args []string, out, errOut io.Writer) int {
	fs := newCommandFlagSet("apitools catalog index")
	root := fs.String("root", "", "Explicit existing catalog root (required)")
	registry := fs.String("registry", catalog.DefaultRegistryPath, "Registration database path relative to root")
	indexPath := fs.String("index", catalog.DefaultOperationIndexPath, "Output index path relative to root")
	catalogPath := fs.String("catalog", "", "Optional synthetic/custom catalog JSON path relative to root")
	jsonOut := fs.Bool("json", false, "Write build counts as JSON")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: apitools catalog index --root DIRECTORY [--registry cache.sqlite] [--index operations.v1.json] [--catalog catalog.json] [--json]")
		fs.PrintDefaults()
	}
	if code, done := parseCommandFlags(fs, args, out, errOut); done {
		return code
	}
	if *root == "" || fs.NArg() != 0 {
		fmt.Fprintln(errOut, "an explicit --root and no positional arguments are required")
		return exitUsage
	}
	options := apitools.CatalogIndexOptions{Root: catalog.RootOptions{Directory: *root, RegistryPath: *registry, IndexPath: *indexPath}, ReadRegistrations: sqlitecache.ReadCatalogSpecArtifacts}
	if *catalogPath != "" {
		if filepath.Clean(*catalogPath) == filepath.Clean(*indexPath) {
			fmt.Fprintln(errOut, "catalog input and index output overlap")
			return exitUsage
		}
		file, err := artifactio.ReadFile(*root, *catalogPath, artifactio.ReadOptions{MaxBytes: 4 << 20})
		if err != nil {
			fmt.Fprintln(errOut, "cannot read confined catalog metadata")
			return exitRuntime
		}
		var cat catalog.Catalog
		decoder := json.NewDecoder(bytes.NewReader(file.Data))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&cat); err != nil {
			fmt.Fprintln(errOut, "invalid custom catalog metadata")
			return exitRuntime
		}
		var trailing any
		if err := decoder.Decode(&trailing); err != io.EOF {
			fmt.Fprintln(errOut, "trailing custom catalog metadata")
			return exitRuntime
		}
		options.Catalog = &cat
	}
	ctx := context.Background()
	index, err := apitools.BuildCatalogOperationIndex(ctx, options)
	if err == nil {
		err = apitools.WriteCatalogOperationIndex(ctx, options, index)
	}
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}
	operations, incomplete := 0, 0
	for _, artifact := range index.Artifacts {
		operations += len(artifact.Operations)
	}
	for _, coverage := range index.Coverage {
		if coverage.State != "indexed" {
			incomplete++
		}
	}
	summary := struct {
		SchemaVersion string `json:"schema_version"`
		Artifacts     int    `json:"artifacts"`
		Operations    int    `json:"operations"`
		Unexamined    int    `json:"unexamined"`
	}{"apitools.catalog-index-build/v1", len(index.Artifacts), operations, incomplete}
	if *jsonOut {
		if err := writeJSON(out, summary); err != nil {
			fmt.Fprintln(errOut, err)
			return exitRuntime
		}
	} else {
		fmt.Fprintf(out, "indexed %d artifacts, %d operation candidates; %d unexamined coverage entries\n", summary.Artifacts, summary.Operations, summary.Unexamined)
	}
	return exitSuccess
}
