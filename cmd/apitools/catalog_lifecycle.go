package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/artifactio"
	"github.com/OpenUdon/apitools/sqlitecache"
)

func runCatalogInit(args []string, out, errOut io.Writer) int {
	fs := newCommandFlagSet("apitools catalog init")
	rootFlag := fs.String("root", "", "Explicit catalog root directory (or use APITOOLS_CATALOG_ROOT)")
	dirFlag := fs.String("dir", "", "Alias for --root")
	registry := fs.String("registry", catalog.DefaultRegistryPath, "Registration database path relative to root")
	indexPath := fs.String("index", catalog.DefaultOperationIndexPath, "Output index path relative to root")
	catalogPath := fs.String("catalog", "", "Optional custom catalog JSON path relative to root")
	jsonOut := fs.Bool("json", false, "Write JSON output")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: apitools catalog init [--root DIRECTORY] [--registry cache.sqlite] [--index operations.v1.json] [--catalog catalog.json] [--json]")
		fs.PrintDefaults()
	}
	if code, done := parseCommandFlags(fs, args, out, errOut); done {
		return code
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(errOut, "unexpected argument %q\n", fs.Arg(0))
		fs.Usage()
		return exitUsage
	}

	selectedDir := strings.TrimSpace(*rootFlag)
	if selectedDir == "" {
		selectedDir = strings.TrimSpace(*dirFlag)
	}
	rootOpts := catalog.RootOptionsWithEnvironment(catalog.RootOptions{
		Directory:    selectedDir,
		RegistryPath: *registry,
		IndexPath:    *indexPath,
	})
	if strings.TrimSpace(rootOpts.Directory) == "" {
		fmt.Fprintln(errOut, "catalog root directory is required (use --root or APITOOLS_CATALOG_ROOT)")
		return exitUsage
	}

	rootDir, err := filepath.Abs(rootOpts.Directory)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}
	if filepath.Dir(rootDir) == rootDir {
		fmt.Fprintln(errOut, "catalog root must not be a filesystem root")
		return exitRuntime
	}

	// Create root directory and specs directory
	if err := os.MkdirAll(rootDir, 0o755); err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}
	specsDir := filepath.Join(rootDir, "specs")
	if err := os.MkdirAll(specsDir, 0o755); err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}

	// Ensure registry parent directory exists
	regAbs := filepath.Join(rootDir, rootOpts.RegistryPath)
	if err := os.MkdirAll(filepath.Dir(regAbs), 0o755); err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}

	// Initialize SQLite database
	db, err := sqlitecache.Open(regAbs)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}
	_ = db.Close()

	// Initialize operations index
	indexOpts := apitools.CatalogIndexOptions{
		Root: catalog.RootOptions{
			Directory:    rootDir,
			RegistryPath: rootOpts.RegistryPath,
			IndexPath:    rootOpts.IndexPath,
		},
		ReadRegistrations: sqlitecache.ReadCatalogSpecArtifacts,
	}
	if *catalogPath != "" {
		file, err := artifactio.ReadFile(rootDir, *catalogPath, artifactio.ReadOptions{MaxBytes: 4 << 20})
		if err != nil {
			fmt.Fprintln(errOut, "cannot read custom catalog metadata")
			return exitRuntime
		}
		var cat catalog.Catalog
		dec := json.NewDecoder(bytes.NewReader(file.Data))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&cat); err != nil {
			fmt.Fprintln(errOut, "invalid custom catalog metadata")
			return exitRuntime
		}
		indexOpts.Catalog = &cat
	}

	ctx := context.Background()
	index, err := apitools.BuildCatalogOperationIndex(ctx, indexOpts)
	if err == nil {
		err = apitools.WriteCatalogOperationIndex(ctx, indexOpts, index)
	}
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}

	if *jsonOut {
		summary := struct {
			SchemaVersion string `json:"schema_version"`
			Root          string `json:"root"`
			Registry      string `json:"registry"`
			Index         string `json:"index"`
		}{
			SchemaVersion: "apitools.catalog-init/v1",
			Root:          rootDir,
			Registry:      rootOpts.RegistryPath,
			Index:         rootOpts.IndexPath,
		}
		if err := writeJSON(out, summary); err != nil {
			fmt.Fprintln(errOut, err)
			return exitRuntime
		}
		return exitSuccess
	}

	fmt.Fprintf(out, "initialized catalog root %s (registry: %s, index: %s)\n", rootDir, rootOpts.RegistryPath, rootOpts.IndexPath)
	return exitSuccess
}

func runCatalogAdd(args []string, out, errOut io.Writer) int {
	var filePath string
	parseArgs := args
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		filePath = args[0]
		parseArgs = args[1:]
	}
	fs := newCommandFlagSet("apitools catalog add")
	providerFlag := fs.String("provider", "", "Provider ID (optional alias --id)")
	idFlag := fs.String("id", "", "Alias for --provider")
	specFlag := fs.String("spec", "", "Spec reference ID (optional alias --artifact)")
	artifactFlag := fs.String("artifact", "", "Alias for --spec")
	kindFlag := fs.String("kind", "openapi", "Artifact kind (openapi, swagger2, etc.)")
	sourceURL := fs.String("source-url", "", "Optional original source URL")
	rootFlag := fs.String("root", "", "Explicit catalog root directory (or use APITOOLS_CATALOG_ROOT)")
	registry := fs.String("registry", catalog.DefaultRegistryPath, "Registration database path relative to root")
	jsonOut := fs.Bool("json", false, "Write JSON output")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: apitools catalog add <file> [--provider <id>] [--spec <id>] [--kind openapi] [--source-url <url>] [--root DIRECTORY] [--registry cache.sqlite] [--json]")
		fs.PrintDefaults()
	}
	if code, done := parseCommandFlags(fs, parseArgs, out, errOut); done {
		return code
	}
	if filePath == "" {
		if fs.NArg() > 0 {
			filePath = fs.Arg(0)
		}
		if fs.NArg() > 1 {
			fmt.Fprintf(errOut, "unexpected argument %q\n", fs.Arg(1))
			fs.Usage()
			return exitUsage
		}
	} else if fs.NArg() > 0 {
		fmt.Fprintf(errOut, "unexpected argument %q\n", fs.Arg(0))
		fs.Usage()
		return exitUsage
	}
	if filePath == "" {
		fmt.Fprintln(errOut, "a spec file path is required")
		fs.Usage()
		return exitUsage
	}

	rootOpts := catalog.RootOptionsWithEnvironment(catalog.RootOptions{
		Directory:    *rootFlag,
		RegistryPath: *registry,
	})
	if strings.TrimSpace(rootOpts.Directory) == "" {
		fmt.Fprintln(errOut, "catalog root directory is required (use --root or APITOOLS_CATALOG_ROOT)")
		return exitUsage
	}
	paths, err := catalog.ResolveRoot(rootOpts)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}

	// Read and validate local file
	absFile, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}
	info, err := os.Lstat(absFile)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}
	if !info.Mode().IsRegular() {
		fmt.Fprintf(errOut, "spec file %q is not a regular file\n", filePath)
		return exitRuntime
	}
	if info.Size() > artifactio.DefaultMaxBytes {
		fmt.Fprintf(errOut, "spec file exceeds maximum allowed size (%d bytes)\n", artifactio.DefaultMaxBytes)
		return exitRuntime
	}
	data, err := os.ReadFile(absFile)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}
	if len(data) == 0 {
		fmt.Fprintln(errOut, "spec file is empty")
		return exitRuntime
	}

	kind := strings.TrimSpace(*kindFlag)
	if kind == "" {
		kind = "openapi"
	}

	// Validate content if openapi or swagger
	var specMeta apitools.SpecMetadata
	if kind == "openapi" || kind == "swagger2" {
		ctx := context.Background()
		var valErr error
		specMeta, valErr = apitools.ValidateLocalSpecContent(ctx, data)
		if valErr != nil {
			fmt.Fprintf(errOut, "invalid spec content: %v\n", valErr)
			return exitRuntime
		}
	}

	// Compute hash
	digest := sha256.Sum256(data)
	shaHex := hex.EncodeToString(digest[:])

	// Determine provider and spec IDs
	provider := strings.TrimSpace(*providerFlag)
	if provider == "" {
		provider = strings.TrimSpace(*idFlag)
	}
	specID := strings.TrimSpace(*specFlag)
	if specID == "" {
		specID = strings.TrimSpace(*artifactFlag)
	}
	baseStem := strings.TrimSuffix(filepath.Base(absFile), filepath.Ext(absFile))
	if provider == "" {
		if specMeta.Title != "" {
			provider = strings.ToLower(strings.ReplaceAll(specMeta.Title, " ", "-"))
		} else {
			provider = baseStem
		}
	}
	if specID == "" {
		specID = baseStem
	}
	provider = strings.ToLower(provider)

	// Determine destination path inside root
	var relPath string
	relToRoot, err := filepath.Rel(paths.Directory, absFile)
	if err == nil && !strings.HasPrefix(relToRoot, "..") && !filepath.IsAbs(relToRoot) {
		relPath = filepath.ToSlash(relToRoot)
	} else {
		destDir := filepath.Join(paths.Directory, "specs", provider)
		if err := os.MkdirAll(destDir, 0o755); err != nil {
			fmt.Fprintln(errOut, err)
			return exitRuntime
		}
		destPath := filepath.Join(destDir, filepath.Base(absFile))
		if err := os.WriteFile(destPath, data, 0o644); err != nil {
			fmt.Fprintln(errOut, err)
			return exitRuntime
		}
		relPath = filepath.ToSlash(filepath.Join("specs", provider, filepath.Base(absFile)))
	}

	// Open cache and store artifact
	cache, err := sqlitecache.Open(paths.RegistryPath)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}
	defer cache.Close()

	artifact := sqlitecache.CatalogArtifact{
		ProviderID: provider,
		ArtifactID: specID,
		Kind:       kind,
		Path:       relPath,
		SourceURL:  strings.TrimSpace(*sourceURL),
		SHA256:     shaHex,
		Bytes:      int64(len(data)),
	}

	ctx := context.Background()
	if err := cache.StoreCatalogArtifact(ctx, artifact); err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}

	if *jsonOut {
		if err := writeJSON(out, artifact); err != nil {
			fmt.Fprintln(errOut, err)
			return exitRuntime
		}
		return exitSuccess
	}

	fmt.Fprintf(out, "added %s/%s (%s, %d bytes, sha256: %s)\n", provider, specID, relPath, len(data), shaHex)
	return exitSuccess
}

func runCatalogImport(args []string, out, errOut io.Writer) int {
	return runCatalogImportWithClient(args, out, errOut, clientForCache)
}

func runCatalogImportWithClient(args []string, out, errOut io.Writer, newClient func(string) (*apitools.Client, func(), error)) int {
	var rawURL string
	parseArgs := args
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		rawURL = args[0]
		parseArgs = args[1:]
	}
	fs := newCommandFlagSet("apitools catalog import")
	urlFlag := fs.String("url", "", "OpenAPI document URL")
	providerFlag := fs.String("provider", "", "Provider ID (optional)")
	nameFlag := fs.String("name", "", "Suggested filename stem (optional)")
	rootFlag := fs.String("root", "", "Explicit catalog root directory (or use APITOOLS_CATALOG_ROOT)")
	registry := fs.String("registry", catalog.DefaultRegistryPath, "Registration database path relative to root")
	timeout := fs.Duration("timeout", apitools.DefaultTimeout, "Import download timeout")
	jsonOut := fs.Bool("json", false, "Write JSON output")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: apitools catalog import <url> [--provider <id>] [--name <stem>] [--root DIRECTORY] [--registry cache.sqlite] [--timeout 30s] [--json]")
		fs.PrintDefaults()
	}
	if code, done := parseCommandFlags(fs, parseArgs, out, errOut); done {
		return code
	}
	if rawURL == "" {
		if fs.NArg() > 0 {
			rawURL = fs.Arg(0)
		}
		if fs.NArg() > 1 {
			fmt.Fprintf(errOut, "unexpected argument %q\n", fs.Arg(1))
			fs.Usage()
			return exitUsage
		}
	} else if fs.NArg() > 0 {
		fmt.Fprintf(errOut, "unexpected argument %q\n", fs.Arg(0))
		fs.Usage()
		return exitUsage
	}
	if rawURL == "" {
		rawURL = *urlFlag
	}
	trimmedURL := strings.TrimSpace(rawURL)
	if trimmedURL == "" {
		fmt.Fprintln(errOut, "an OpenAPI document URL is required")
		fs.Usage()
		return exitUsage
	}
	if err := validateImportURL(trimmedURL); err != nil {
		fmt.Fprintln(errOut, err)
		fs.Usage()
		return exitUsage
	}

	rootOpts := catalog.RootOptionsWithEnvironment(catalog.RootOptions{
		Directory:    *rootFlag,
		RegistryPath: *registry,
	})
	if strings.TrimSpace(rootOpts.Directory) == "" {
		fmt.Fprintln(errOut, "catalog root directory is required (use --root or APITOOLS_CATALOG_ROOT)")
		return exitUsage
	}
	paths, err := catalog.ResolveRoot(rootOpts)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}

	// Derive provider ID
	provider := strings.TrimSpace(*providerFlag)
	if provider == "" {
		parsed, _ := url.Parse(trimmedURL)
		host := parsed.Hostname()
		parts := strings.Split(host, ".")
		if len(parts) >= 2 {
			provider = parts[len(parts)-2]
		} else {
			provider = host
		}
	}
	provider = strings.ToLower(strings.TrimSpace(provider))
	if provider == "" {
		provider = "imported"
	}

	targetDir := filepath.Join(paths.Directory, "specs", provider)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client, closeCache, err := newClient(paths.RegistryPath)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}
	defer closeCache()
	if *timeout > 0 {
		client.Timeout = *timeout
	}

	importReport, err := client.ImportWithReport(ctx, apitools.ImportOptions{
		URL:       trimmedURL,
		Dir:       targetDir,
		Name:      *nameFlag,
		CacheMode: apitools.CacheModeBypass,
	})
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}

	imported := importReport.Imported
	relPath, err := filepath.Rel(paths.Directory, imported.Path)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}
	relPath = filepath.ToSlash(relPath)

	stem := *nameFlag
	if stem == "" {
		stem = strings.TrimSuffix(filepath.Base(imported.Path), filepath.Ext(imported.Path))
	}

	// Store in SQLite cache as catalog artifact
	cache, err := sqlitecache.Open(paths.RegistryPath)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}
	defer cache.Close()

	artifact := sqlitecache.CatalogArtifact{
		ProviderID: provider,
		ArtifactID: stem,
		Kind:       "openapi",
		Path:       relPath,
		SourceURL:  trimmedURL,
		SHA256:     imported.SHA256,
		Bytes:      imported.Bytes,
	}
	if err := cache.StoreCatalogArtifact(ctx, artifact); err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}

	if *jsonOut {
		if err := writeJSON(out, artifact); err != nil {
			fmt.Fprintln(errOut, err)
			return exitRuntime
		}
		return exitSuccess
	}

	fmt.Fprintf(out, "imported %s/%s (%s, %d bytes, sha256: %s)\n", provider, stem, relPath, imported.Bytes, imported.SHA256)
	return exitSuccess
}

func runCatalogSnapshot(args []string, out, errOut io.Writer) int {
	fs := newCommandFlagSet("apitools catalog snapshot")
	target := fs.String("target", "", "Target directory for catalog snapshot (required)")
	rootFlag := fs.String("root", "", "Source catalog root directory (or use APITOOLS_CATALOG_ROOT)")
	registry := fs.String("registry", catalog.DefaultRegistryPath, "Registration database path relative to root")
	indexPath := fs.String("index", catalog.DefaultOperationIndexPath, "Index path relative to root")
	selectedSpec := fs.String("selected-spec", "", "Optional selected spec path relative to root")
	metadataSrc := fs.String("metadata-src", "", "Optional external metadata source path")
	metadataDest := fs.String("metadata-dest", "", "Optional relative metadata path in target")
	noNormalize := fs.Bool("no-normalize-journal", false, "Disable SQLite journal mode normalization to DELETE")
	maxBytes := fs.Int64("max-bytes", 0, "Maximum bytes per copied registry/index file")
	jsonOut := fs.Bool("json", false, "Write JSON output")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: apitools catalog snapshot --target DIRECTORY [--root DIRECTORY] [--registry cache.sqlite] [--index operations.v1.json] [--selected-spec PATH] [--metadata-src PATH] [--metadata-dest PATH] [--no-normalize-journal] [--max-bytes N] [--json]")
		fs.PrintDefaults()
	}
	if code, done := parseCommandFlags(fs, args, out, errOut); done {
		return code
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(errOut, "unexpected argument %q\n", fs.Arg(0))
		fs.Usage()
		return exitUsage
	}
	if strings.TrimSpace(*target) == "" {
		fmt.Fprintln(errOut, "--target directory is required")
		fs.Usage()
		return exitUsage
	}

	rootOpts := catalog.RootOptionsWithEnvironment(catalog.RootOptions{
		Directory:    *rootFlag,
		RegistryPath: *registry,
		IndexPath:    *indexPath,
	})
	if strings.TrimSpace(rootOpts.Directory) == "" {
		fmt.Fprintln(errOut, "source catalog root directory is required (use --root or APITOOLS_CATALOG_ROOT)")
		return exitUsage
	}
	paths, err := catalog.ResolveRoot(rootOpts)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}

	snapOpts := catalog.SnapshotOptions{
		MetadataSource:      strings.TrimSpace(*metadataSrc),
		MetadataDestination: strings.TrimSpace(*metadataDest),
		SelectedSourcePath:  strings.TrimSpace(*selectedSpec),
		MaxRegistryBytes:    *maxBytes,
		MaxIndexBytes:       *maxBytes,
		NormalizeJournal:    !*noNormalize,
	}

	if err := catalog.SnapshotRoot(paths, strings.TrimSpace(*target), snapOpts); err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}

	if *jsonOut {
		summary := struct {
			SchemaVersion string `json:"schema_version"`
			SourceRoot    string `json:"source_root"`
			Target        string `json:"target"`
			Normalized    bool   `json:"normalized"`
		}{
			SchemaVersion: "apitools.catalog-snapshot/v1",
			SourceRoot:    paths.Directory,
			Target:        *target,
			Normalized:    !*noNormalize,
		}
		if err := writeJSON(out, summary); err != nil {
			fmt.Fprintln(errOut, err)
			return exitRuntime
		}
		return exitSuccess
	}

	fmt.Fprintf(out, "snapshot created in %s from %s\n", *target, paths.Directory)
	return exitSuccess
}

func runCatalogSearch(args []string, out, errOut io.Writer) int {
	return runCatalogSearchWithClient(args, out, errOut, clientForCache)
}

func runCatalogSearchWithClient(args []string, out, errOut io.Writer, newClient func(string) (*apitools.Client, func(), error)) int {
	var query string
	parseArgs := args
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		query = args[0]
		parseArgs = args[1:]
	}
	fs := newCommandFlagSet("apitools catalog search")
	queryFlag := fs.String("query", "", "Search query (or first positional argument)")
	rootFlag := fs.String("root", "", "Catalog root directory (or use APITOOLS_CATALOG_ROOT)")
	registry := fs.String("registry", catalog.DefaultRegistryPath, "Registration database path relative to root")
	indexPath := fs.String("index", catalog.DefaultOperationIndexPath, "Index path relative to root")
	remote := fs.Bool("remote", false, "Enable remote catalog lookup")
	providerURL := fs.String("provider-url", "", "Provider URL or hostname for RFC 9727 discovery")
	limit := fs.Int("limit", 10, "Maximum result count")
	jsonOut := fs.Bool("json", false, "Write JSON output")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: apitools catalog search <query> [--root DIRECTORY] [--registry cache.sqlite] [--index operations.v1.json] [--remote] [--provider-url URL] [--limit 10] [--json]")
		fs.PrintDefaults()
	}
	if code, done := parseCommandFlags(fs, parseArgs, out, errOut); done {
		return code
	}
	if query == "" {
		if fs.NArg() > 0 {
			query = fs.Arg(0)
		}
		if fs.NArg() > 1 {
			fmt.Fprintf(errOut, "unexpected argument %q\n", fs.Arg(1))
			fs.Usage()
			return exitUsage
		}
	} else if fs.NArg() > 0 {
		fmt.Fprintf(errOut, "unexpected argument %q\n", fs.Arg(0))
		fs.Usage()
		return exitUsage
	}
	if query == "" {
		query = *queryFlag
	}
	trimmedQuery := strings.TrimSpace(query)
	if trimmedQuery == "" {
		fmt.Fprintln(errOut, "a search query is required")
		fs.Usage()
		return exitUsage
	}

	rootOpts := catalog.RootOptionsWithEnvironment(catalog.RootOptions{
		Directory:    *rootFlag,
		RegistryPath: *registry,
		IndexPath:    *indexPath,
	})

	// If root is configured and exists, search the index via DiscoverCatalogOperations
	if rootOpts.Directory != "" {
		paths, err := catalog.ResolveRoot(rootOpts)
		if err == nil {
			indexFile := paths.IndexPath
			if fi, err := os.Lstat(indexFile); err == nil && fi.Mode().IsRegular() {
				// Read commands never implicitly repair, refresh or index
				discoveryOpts := apitools.CatalogDiscoveryOptions{
					Index: apitools.CatalogIndexOptions{
						Root:              rootOpts,
						ReadRegistrations: sqlitecache.ReadCatalogSpecArtifacts,
					},
					Request: apitools.CatalogDiscoveryRequest{
						SchemaVersion: apitools.CatalogDiscoverySchemaVersion,
						Contract: apitools.StepContract{
							Purpose: trimmedQuery,
							Effect:  apitools.OperationEffectRead,
						},
						Limits: apitools.CatalogDiscoveryLimits{
							MaxResults: *limit,
						},
						RemoteLookup: *remote,
					},
					RemoteEnabled: *remote,
					ProviderURL:   strings.TrimSpace(*providerURL),
				}
				if *remote {
					remoteClient, closeCache, clientErr := newClient(paths.RegistryPath)
					if clientErr == nil {
						defer closeCache()
						discoveryOpts.RemoteClient = remoteClient
					}
				}

				ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
				defer stop()

				report, err := apitools.DiscoverCatalogOperations(ctx, discoveryOpts)
				if err != nil {
					fmt.Fprintln(errOut, err)
					return exitRuntime
				}
				if *jsonOut {
					if err := writeJSON(out, report); err != nil {
						fmt.Fprintln(errOut, err)
						return exitRuntime
					}
					return exitSuccess
				}

				if len(report.Candidates) == 0 {
					fmt.Fprintf(out, "No matching operations found for %q (outcome: %s).\n", trimmedQuery, report.Outcome)
					return exitSuccess
				}
				fmt.Fprintf(out, "Found %d candidates (outcome: %s):\n", len(report.Candidates), report.Outcome)
				for i, cand := range report.Candidates {
					fmt.Fprintf(out, "%d. %s %s [qualified: %t]\n", i+1, cand.Candidate.Operation.Method, cand.Candidate.Operation.Path, cand.Qualified)
					if cand.Candidate.Operation.OperationID != "" {
						fmt.Fprintf(out, "   operation: %s\n", cand.Candidate.Operation.OperationID)
					}
					if len(cand.References) > 0 {
						fmt.Fprintf(out, "   provider:  %s (%s)\n", cand.References[0].ProviderID, cand.References[0].SpecRefID)
					}
					if cand.Candidate.Operation.Summary != "" {
						fmt.Fprintf(out, "   summary:   %s\n", singleLine(cand.Candidate.Operation.Summary))
					}
				}
				return exitSuccess
			}
		}
	}

	// Fallback when no root/index is available: search built-in providers
	cat := catalog.BuiltInCatalog()
	type providerMatch struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
		Category    string `json:"category"`
	}
	var matches []providerMatch
	lowerQuery := strings.ToLower(trimmedQuery)
	for _, p := range cat.ListProviders() {
		if strings.Contains(strings.ToLower(p.ID), lowerQuery) ||
			strings.Contains(strings.ToLower(p.DisplayName), lowerQuery) ||
			strings.Contains(strings.ToLower(p.Category), lowerQuery) {
			matches = append(matches, providerMatch{
				ID:          p.ID,
				DisplayName: p.DisplayName,
				Category:    p.Category,
			})
			if len(matches) >= *limit {
				break
			}
		}
	}

	if *jsonOut {
		if err := writeJSON(out, matches); err != nil {
			fmt.Fprintln(errOut, err)
			return exitRuntime
		}
		return exitSuccess
	}

	if len(matches) == 0 {
		fmt.Fprintf(out, "No catalog providers or operations found for %q.\n", trimmedQuery)
		return exitSuccess
	}
	fmt.Fprintf(out, "Matching catalog providers (%d):\n", len(matches))
	for i, m := range matches {
		fmt.Fprintf(out, "%d. %s (%s) [%s]\n", i+1, m.DisplayName, m.ID, m.Category)
	}
	return exitSuccess
}
