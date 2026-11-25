package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/artifactio"
	"github.com/OpenUdon/apitools/internal/catalogupdate"
)

const (
	catalogCronStart = "# BEGIN apitools catalog maintenance"
	catalogCronEnd   = "# END apitools catalog maintenance"
	catalogCronWhen  = "17 3 1 * *"
	catalogCronLimit = 1 << 20
)

type catalogCronStatus struct {
	Enabled     bool   `json:"enabled"`
	Entry       string `json:"entry,omitempty"`
	LatestRunID string `json:"latest_run_id,omitempty"`
}

type catalogCronOutput struct {
	buffer bytes.Buffer
	limit  int
}

func (output *catalogCronOutput) Write(data []byte) (int, error) {
	if len(data) > output.limit-output.buffer.Len() {
		return 0, errors.New("crontab command output exceeds the maintenance limit")
	}
	return output.buffer.Write(data)
}

func canonicalCatalogCronPath(path string) (string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	parent, err := filepath.EvalSymlinks(filepath.Dir(absolute))
	if errors.Is(err, os.ErrNotExist) {
		return absolute, nil
	}
	if err != nil {
		return "", err
	}
	return filepath.Join(parent, filepath.Base(absolute)), nil
}

func resolveCatalogCronScratchPath(file string, protected ...string) (string, error) {
	if file == "" {
		return "", nil
	}
	path, err := canonicalCatalogCronPath(file)
	if err != nil {
		return "", err
	}
	// The confined file helpers trim path arguments. Refuse ambiguous
	// names before comparing identities or passing the path to them.
	for _, component := range strings.Split(path, string(filepath.Separator)) {
		if component != strings.TrimSpace(component) {
			return "", errors.New("scratch crontab path components must exclude leading and trailing whitespace")
		}
	}
	info, err := os.Stat(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	for _, input := range protected {
		if input == "" {
			continue
		}
		input, err := canonicalCatalogCronPath(input)
		if err != nil {
			return "", err
		}
		if path == input {
			return "", errors.New("scratch crontab must be separate from maintenance configuration and state files")
		}
		if info != nil {
			other, err := os.Stat(input)
			if err == nil && os.SameFile(info, other) {
				return "", errors.New("scratch crontab aliases a maintenance configuration or state file")
			}
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return "", err
			}
		}
	}
	return path, nil
}

func readCatalogCrontab(ctx context.Context, file string) (string, error) {
	if file != "" {
		data, err := artifactio.ReadFile(filepath.Dir(file), filepath.Base(file), artifactio.ReadOptions{MaxBytes: catalogCronLimit})
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		if err == nil && (!utf8.Valid(data.Data) || bytes.IndexByte(data.Data, 0) >= 0) {
			return "", errors.New("scratch crontab must contain text")
		}
		return string(data.Data), err
	}
	cmd := exec.CommandContext(ctx, "crontab", "-l")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	stdout := catalogCronOutput{limit: catalogCronLimit}
	stderr := catalogCronOutput{limit: 8 << 10}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.WaitDelay = time.Second
	err := cmd.Run()
	if err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) && exit.ExitCode() == 1 && stdout.buffer.Len() == 0 &&
			strings.HasPrefix(strings.TrimSpace(stderr.buffer.String()), "no crontab for ") {
			return "", nil
		}
		return "", fmt.Errorf("read user crontab: %w (%s)", err, strings.TrimSpace(stderr.buffer.String()))
	}
	return stdout.buffer.String(), nil
}

func saveCatalogCrontab(ctx context.Context, file, original, updated string) error {
	if updated == original {
		return nil
	}
	if len(updated) > catalogCronLimit {
		return errors.New("updated crontab exceeds the maintenance size limit")
	}
	current, err := readCatalogCrontab(ctx, file)
	if err != nil {
		return err
	}
	if current != original {
		return errors.New("crontab changed during the maintenance edit")
	}
	if file != "" {
		mode := os.FileMode(0o600)
		if info, err := os.Lstat(file); err == nil {
			mode = info.Mode().Perm()
		}
		_, err := artifactio.WriteFile(filepath.Dir(file), filepath.Base(file), []byte(updated), artifactio.WriteOptions{Mode: mode, Force: true})
		return err
	}
	cmd := exec.CommandContext(ctx, "crontab", "-")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	cmd.Stdin = strings.NewReader(updated)
	diagnostics := catalogCronOutput{limit: 8 << 10}
	cmd.Stdout, cmd.Stderr = &diagnostics, &diagnostics
	cmd.WaitDelay = time.Second
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("install user crontab: %w (%s)", err, strings.TrimSpace(diagnostics.buffer.String()))
	}
	return nil
}

func splitCatalogCronTable(raw string) (unmanaged, entry string, err error) {
	var rest strings.Builder
	inside, found := false, false
	for _, line := range strings.SplitAfter(raw, "\n") {
		content := strings.TrimRight(line, "\r\n")
		switch content {
		case catalogCronStart:
			if inside || found {
				return "", "", errors.New("crontab contains duplicate apitools maintenance markers")
			}
			inside = true
		case catalogCronEnd:
			if !inside || entry == "" {
				return "", "", errors.New("crontab contains an invalid apitools maintenance block")
			}
			inside, found = false, true
		default:
			if inside {
				if entry != "" || !strings.HasPrefix(content, catalogCronWhen+" ") ||
					!strings.Contains(content, " catalog maintenance run --config ") {
					return "", "", errors.New("crontab contains an invalid apitools maintenance entry")
				}
				entry = content
			} else {
				rest.WriteString(line)
			}
		}
	}
	if inside {
		return "", "", errors.New("crontab contains an unterminated apitools maintenance block")
	}
	return rest.String(), entry, nil
}

func makeCatalogCronEntry(binary, config, runID string) (string, error) {
	if !filepath.IsAbs(binary) || !filepath.IsAbs(config) {
		return "", errors.New("cron binary and configuration paths must be absolute")
	}
	var quoted []string
	for _, arg := range []string{binary, config, runID} {
		if arg == "" || strings.ContainsAny(arg, "\x00\r\n%") {
			return "", errors.New("cron paths and run ID must exclude percent signs and line breaks")
		}
		quoted = append(quoted, "'"+strings.ReplaceAll(arg, "'", "'\\''")+"'")
	}
	return catalogCronWhen + " " + quoted[0] + " catalog maintenance run --config " + quoted[1] + " --after-run " + quoted[2] + " --json", nil
}

func requireCatalogMaintenanceFullRun(ctx context.Context, statePath, runID string) error {
	if runID == "" || statePath == "" {
		return errors.New("select a completed full reviewed run with --after-run and a configured state_db")
	}
	run, _, err := catalogupdate.ReadExistingRun(ctx, statePath, runID)
	if err != nil {
		return err
	}
	var providers []string
	for _, provider := range catalog.BuiltInProviders() {
		providers = append(providers, provider.ID)
	}
	if !run.Completed || !slices.Equal(run.Providers, providers) {
		return errors.New("schedule requires a completed run covering the full reviewed provider list")
	}
	return nil
}

func manageCatalogCron(ctx context.Context, action, file, configPath, runID string, cfg apitools.CatalogUpdateConfig, out io.Writer, jsonOutput bool) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	protected := []string{configPath, cfg.StateDB}
	if cfg.StateDB != "" {
		for _, suffix := range []string{"-journal", "-wal", "-shm", ".monthly.lock"} {
			protected = append(protected, cfg.StateDB+suffix)
		}
	}
	file, err := resolveCatalogCronScratchPath(file, protected...)
	if err != nil {
		return err
	}
	original, err := readCatalogCrontab(ctx, file)
	if err != nil {
		return err
	}
	unmanaged, entry, err := splitCatalogCronTable(original)
	if err != nil {
		return err
	}
	var updated string
	switch action {
	case "enable":
		readiness, err := catalogMaintenancePreflight(ctx, configPath, cfg)
		if err != nil {
			return err
		}
		if err := requireCatalogMaintenanceFullRun(ctx, cfg.StateDB, runID); err != nil {
			return err
		}
		entry, err = makeCatalogCronEntry(readiness.Binary, readiness.Config, runID)
		if err != nil {
			return err
		}
		// Schedule the managed entry before any unrelated cron environment assignments.
		updated = catalogCronStart + "\n" + entry + "\n" + catalogCronEnd + "\n" + unmanaged
	case "disable":
		entry = ""
		updated = unmanaged
	}
	if action != "status" {
		if err := saveCatalogCrontab(ctx, file, original, updated); err != nil {
			return err
		}
	}
	status := catalogCronStatus{Enabled: entry != "", Entry: entry}
	if action == "status" && cfg.StateDB != "" {
		status.LatestRunID, err = catalogupdate.ReadLatestMonthlyRunID(ctx, cfg.StateDB)
		if err != nil {
			return err
		}
	}
	if jsonOutput {
		return writeJSON(out, status)
	}
	if status.Enabled {
		fmt.Fprintf(out, "Catalog maintenance schedule enabled\n%s\n", entry)
	} else {
		fmt.Fprintln(out, "Catalog maintenance schedule disabled")
	}
	if status.LatestRunID != "" {
		fmt.Fprintln(out, "Latest scheduled run", status.LatestRunID)
	}
	return nil
}
