package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/apitools/internal/catalogupdate"
)

type catalogMaintenanceReadiness struct {
	Binary   string `json:"binary"`
	Config   string `json:"config"`
	CacheDir string `json:"cache_dir"`
	StateDB  string `json:"state_db"`
	CronHome string `json:"cron_home"`
	TimeZone string `json:"time_zone"`
}

func runCatalogMaintenance(args []string, out, errOut io.Writer) int {
	return runCatalogMaintenanceWithClient(args, out, errOut, &apitools.Client{})
}

func runCatalogMaintenanceWithClient(args []string, out, errOut io.Writer, client *apitools.Client) int {
	const usage = "Usage: apitools catalog maintenance preflight|run|enable|disable|status|report|retry [--config FILE] [--json]"
	if len(args) == 0 {
		fmt.Fprintln(errOut, usage)
		return exitUsage
	}
	action := args[0]
	if action == "--help" || action == "-h" || action == "help" {
		fmt.Fprintln(out, usage)
		return exitSuccess
	}
	switch action {
	case "preflight", "run", "enable", "disable", "status", "report", "retry":
	default:
		fmt.Fprintf(errOut, "unknown catalog maintenance command %q\n", action)
		fmt.Fprintln(errOut, usage)
		return exitUsage
	}
	fs := newCommandFlagSet("apitools catalog maintenance " + action)
	configFlag := fs.String("config", "", "Catalog maintenance JSON configuration file")
	jsonOutput := fs.Bool("json", false, "Write JSON output")
	var dryRun bool
	var afterRun, cronFile, inspectRun, provider string
	switch action {
	case "run":
		fs.BoolVar(&dryRun, "dry-run", false, "List the provider scope without opening state or contacting publishers")
		fs.StringVar(&afterRun, "after-run", "", "Require this completed full reviewed run before publisher work")
	case "enable":
		fs.StringVar(&afterRun, "after-run", "", "Completed full reviewed run the operator has inspected")
	case "report":
		fs.StringVar(&inspectRun, "run", "", "Inspect this saved run instead of the latest scheduled run")
	case "retry":
		fs.StringVar(&provider, "provider", "", "Check one reviewed provider in a new maintenance run")
	}
	if action == "enable" || action == "disable" || action == "status" {
		fs.StringVar(&cronFile, "crontab-file", "", "Use a scratch crontab file instead of the user crontab")
	}
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: apitools catalog maintenance %s [options]\n", action)
		fs.PrintDefaults()
	}
	if code, done := parseCommandFlags(fs, args[1:], out, errOut); done {
		return code
	}
	emptyFlag := ""
	fs.Visit(func(option *flag.Flag) {
		if (option.Name == "crontab-file" || option.Name == "after-run") && strings.TrimSpace(option.Value.String()) == "" {
			emptyFlag = option.Name
		}
	})
	if emptyFlag != "" {
		fmt.Fprintf(errOut, "--%s requires a nonempty value\n", emptyFlag)
		return exitUsage
	}
	if fs.NArg() != 0 || (*configFlag == "" && action != "disable" && action != "status") ||
		(action == "enable" && afterRun == "") || (action == "retry" && provider == "") {
		fs.Usage()
		return exitUsage
	}
	var config apitools.CatalogUpdateConfig
	var configPath string
	if *configFlag != "" {
		var err error
		configPath, err = filepath.Abs(*configFlag)
		if err == nil {
			config, err = readCatalogUpdateConfig(configPath)
		}
		if err != nil {
			fmt.Fprintln(errOut, err)
			return exitUsage
		}
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if action == "enable" || action == "disable" || action == "status" {
		if err := manageCatalogCron(ctx, action, cronFile, configPath, afterRun, config, out, *jsonOutput); err != nil {
			fmt.Fprintln(errOut, err)
			return exitRuntime
		}
		return exitSuccess
	}
	if action == "report" {
		report, err := loadCatalogMaintenanceReport(ctx, config.StateDB, inspectRun)
		return writeCatalogMaintenanceReport(ctx, out, errOut, report, err, config, *jsonOutput)
	}
	if action == "preflight" {
		readiness, err := catalogMaintenancePreflight(ctx, configPath, config)
		if err != nil {
			fmt.Fprintln(errOut, err)
			return exitRuntime
		}
		if *jsonOutput {
			if err := writeJSON(out, readiness); err != nil {
				fmt.Fprintln(errOut, err)
				return exitRuntime
			}
		} else {
			fmt.Fprintf(out, "Ready for catalog maintenance in %s\nBinary %s\nCache %s\nState %s\n",
				readiness.TimeZone, readiness.Binary, readiness.CacheDir, readiness.StateDB)
		}
		return exitSuccess
	}
	if dryRun {
		report, err := client.UpdateCatalog(ctx, apitools.CatalogUpdateOptions{All: true, DryRun: true})
		return writeCatalogUpdateResult(out, errOut, report, err, *jsonOutput, false)
	}
	onStart := func(id string) {
		fmt.Fprintln(errOut, "Catalog maintenance run", id)
	}
	if action == "retry" {
		lock, err := catalogupdate.LockMonthlyRun(ctx, config.StateDB)
		if err != nil {
			fmt.Fprintln(errOut, err)
			return exitRuntime
		}
		defer lock.Close()
		report, err := client.UpdateCatalog(ctx, apitools.CatalogUpdateOptions{Config: config, Provider: provider, OnStart: onStart})
		return writeCatalogMaintenanceReport(ctx, out, errOut, report, err, config, *jsonOutput)
	}
	report, err := client.RunMonthlyCatalogUpdateAfter(ctx, config, afterRun, onStart)
	return writeCatalogMaintenanceReport(ctx, out, errOut, report, err, config, *jsonOutput)
}

func catalogMaintenancePreflight(ctx context.Context, path string, cfg apitools.CatalogUpdateConfig) (catalogMaintenanceReadiness, error) {
	binary, err := os.Executable()
	if err != nil {
		return catalogMaintenanceReadiness{}, err
	}
	for _, executable := range []string{binary, "/bin/sh"} {
		info, err := os.Stat(executable)
		if err != nil {
			return catalogMaintenanceReadiness{}, err
		}
		if !filepath.IsAbs(executable) || !info.Mode().IsRegular() || info.Mode().Perm()&0o111 == 0 {
			return catalogMaintenanceReadiness{}, fmt.Errorf("cron executable is not an absolute executable file: %q", executable)
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return catalogMaintenanceReadiness{}, err
	}
	info, err := os.Stat(home)
	if err != nil {
		return catalogMaintenanceReadiness{}, err
	}
	if !filepath.IsAbs(home) || !info.IsDir() {
		return catalogMaintenanceReadiness{}, errors.New("cron HOME must be an existing absolute directory")
	}
	if err := catalogupdate.CheckMaintenanceStorage(ctx, cfg.CacheDir, cfg.StateDB, cfg.ExpectedVolumeUUID); err != nil {
		return catalogMaintenanceReadiness{}, err
	}
	zone, _ := time.Now().Zone()
	return catalogMaintenanceReadiness{Binary: binary, Config: path, CacheDir: cfg.CacheDir,
		StateDB: cfg.StateDB, CronHome: home, TimeZone: zone}, nil
}
