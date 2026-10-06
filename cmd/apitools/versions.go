package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/apitools/internal/artifactio"
	"github.com/OpenUdon/apitools/internal/sourceguard"
)

func runVersions(args []string, out, errOut io.Writer) int {
	return runVersionsWithClient(args, out, errOut, &apitools.Client{})
}

func runVersionsWithClient(args []string, out, errOut io.Writer, client *apitools.Client) int {
	fs := newCommandFlagSet("apitools versions")
	requestFile := fs.String("request", "", "Versioned baseline/scope request JSON")
	network := fs.Bool("network", false, "Enable bounded official-source lookup")
	never := fs.Bool("no-network", false, "Use local evidence only")
	timeout := fs.Duration("timeout", apitools.DefaultAPIVersionTimeout, "Total deterministic check budget (maximum 5s)")
	state := fs.String("state", "", "Optional freshness-state JSON path outside a repository")
	cache := fs.String("list-cache", "", "Optional APIs.guru metadata cache path outside a repository")
	save := fs.String("save-dir", "", "Explicit directory for fetched newer Import-valid sources")
	inventoryFile := fs.String("baseline-inventory", "", "Optional prebuilt inventory envelope with sha256 and inventory")
	jsonOut := fs.Bool("json", false, "Write versioned JSON report")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: apitools versions --request FILE [--network|--no-network] [--timeout 5s] [--state FILE] [--list-cache FILE] [--save-dir DIRECTORY] [--baseline-inventory FILE] [--json]")
		fs.PrintDefaults()
	}
	if code, done := parseCommandFlags(fs, args, out, errOut); done {
		return code
	}
	if *requestFile == "" || (*network && *never) || fs.NArg() != 0 {
		fmt.Fprintln(errOut, "--request is required; network flags are exclusive and positional arguments are unsupported")
		return exitUsage
	}
	data, err := readVersionsCLIFile(*requestFile, 64<<10)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}
	req, err := apitools.DecodeAPIVersionDiscoveryRequest(data)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitUsage
	}
	opts := apitools.APIVersionDiscoveryOptions{Network: *network, Timeout: *timeout, StatePath: *state, ListCachePath: *cache, SaveDir: *save}
	if *inventoryFile != "" {
		data, err = readVersionsCLIFile(*inventoryFile, 32<<20)
		if err != nil {
			fmt.Fprintln(errOut, err)
			return exitRuntime
		}
		limits := sourceguard.DefaultLimits()
		limits.MaxDocumentBytes = 32 << 20
		if err = sourceguard.CheckJSONWithLimits(context.Background(), "baseline-inventory", data, limits); err != nil {
			fmt.Fprintln(errOut, err)
			return exitUsage
		}
		var envelope struct {
			SHA256    string                      `json:"sha256"`
			Inventory apitools.OperationInventory `json:"inventory"`
		}
		if err = json.Unmarshal(data, &envelope); err != nil {
			fmt.Fprintln(errOut, err)
			return exitUsage
		}
		opts.BaselineInventory = &envelope.Inventory
		opts.BaselineInventorySHA256 = envelope.SHA256
	}
	report, err := client.DiscoverAPIVersions(context.Background(), req, opts)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitUsage
	}
	if *jsonOut {
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", "  ")
		if err = encoder.Encode(report); err != nil {
			fmt.Fprintln(errOut, err)
			return exitRuntime
		}
	} else {
		fmt.Fprintf(out, "API version check: %s\nHeld baseline: %s\nChecked at: %s\n", report.Status, report.Known.Version, report.CheckedAt.Format("2006-01-02T15:04:05Z"))
		for _, row := range report.Versions {
			fmt.Fprintf(out, "  %s %s %s\n", row.Validation, row.Comparison, row.SourceURL)
		}
		for _, diagnostic := range report.Diagnostics {
			fmt.Fprintf(out, "  %s: %s\n", diagnostic.Code, diagnostic.Message)
		}
	}
	return exitSuccess
}

func readVersionsCLIFile(file string, max int64) ([]byte, error) {
	abs, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}
	data, err := artifactio.ReadFile(filepath.Dir(abs), filepath.Base(abs), artifactio.ReadOptions{MaxBytes: max})
	return data.Data, err
}
