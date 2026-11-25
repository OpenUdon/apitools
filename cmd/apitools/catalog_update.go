package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/OpenUdon/apitools"
)

func runCatalogUpdate(args []string, out, errOut io.Writer) int {
	return runCatalogUpdateWithClient(args, out, errOut, &apitools.Client{})
}

func runCatalogUpdateWithClient(args []string, out, errOut io.Writer, client *apitools.Client) int {
	fs := newCommandFlagSet("apitools catalog update")
	configPath := fs.String("config", "", "Catalog maintenance JSON configuration file")
	var options apitools.CatalogUpdateOptions
	fs.BoolVar(&options.All, "all", false, "Check all reviewed providers")
	fs.StringVar(&options.Provider, "provider", "", "Check one reviewed provider by ID or alias")
	fs.StringVar(&options.Resume, "resume", "", "Resume an unfinished run or inspect a completed report")
	fs.BoolVar(&options.Export, "export", false, "Export deterministic accepted-version metadata")
	fs.BoolVar(&options.DryRun, "dry-run", false, "Show selected providers without network or persistent writes")
	jsonOutput := fs.Bool("json", false, "Write the update report as JSON")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: apitools catalog update [--config FILE] (--all | --provider KEY | --resume ID | --export) [--dry-run] [--json]")
		fs.PrintDefaults()
	}
	if code, done := parseCommandFlags(fs, args, out, errOut); done {
		return code
	}
	if fs.NArg() != 0 || (*configPath == "" && !options.DryRun) {
		fs.Usage()
		return exitUsage
	}
	if *configPath != "" {
		config, err := readCatalogUpdateConfig(*configPath)
		if err != nil {
			fmt.Fprintln(errOut, err)
			return exitUsage
		}
		options.Config = config
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	options.OnStart = func(id string) { fmt.Fprintln(errOut, "Catalog update run", id) }
	report, err := client.UpdateCatalog(ctx, options)
	return writeCatalogUpdateResult(out, errOut, report, err, *jsonOutput, options.Export)
}

func writeCatalogUpdateResult(out, errOut io.Writer, report apitools.CatalogUpdateReport, err error, jsonOutput, export bool) int {
	if jsonOutput {
		if writeErr := writeJSON(out, report); writeErr != nil {
			fmt.Fprintln(errOut, writeErr)
			return exitRuntime
		}
	} else if export && err == nil {
		if _, err := out.Write(report.Export); err != nil {
			fmt.Fprintln(errOut, err)
			return exitRuntime
		}
	} else if report.DryRun {
		fmt.Fprintf(out, "Dry run covers %d reviewed providers\n", len(report.Providers))
		for _, provider := range report.Providers {
			fmt.Fprintln(out, provider)
		}
	} else if report.RunID != "" {
		fmt.Fprintf(out, "Run %s covers %d providers\n", report.RunID, len(report.Providers))
		for _, row := range report.Outcomes {
			fmt.Fprintf(out, "%s %s %s\n", row.ProviderID, row.ReferenceID, row.Status)
		}
		if !report.Completed && err == nil {
			fmt.Fprintln(out, "Publisher retries are pending. Resume this run after the next attempt times in its JSON report.")
		}
	}
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntime
	}
	return 0
}

func readCatalogUpdateConfig(path string) (apitools.CatalogUpdateConfig, error) {
	info, err := os.Stat(path)
	if err != nil {
		return apitools.CatalogUpdateConfig{}, err
	}
	const max = 64 << 10
	if !info.Mode().IsRegular() || info.Size() > max {
		return apitools.CatalogUpdateConfig{}, errors.New("catalog maintenance configuration must be a regular file under 64 KiB")
	}
	f, err := os.Open(path)
	if err != nil {
		return apitools.CatalogUpdateConfig{}, err
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, max+1))
	if err != nil {
		return apitools.CatalogUpdateConfig{}, err
	}
	if len(data) > max {
		return apitools.CatalogUpdateConfig{}, errors.New("catalog maintenance configuration exceeds 64 KiB")
	}
	var cfg apitools.CatalogUpdateConfig
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("decode catalog maintenance configuration: %w", err)
	}
	if err := dec.Decode(new(any)); err != io.EOF {
		return cfg, errors.New("catalog maintenance configuration contains trailing data")
	}
	dir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return cfg, err
	}
	for _, target := range []*string{&cfg.CacheDir, &cfg.StateDB} {
		if *target != "" && !filepath.IsAbs(*target) {
			*target = filepath.Join(dir, *target)
		}
	}
	return cfg, nil
}
