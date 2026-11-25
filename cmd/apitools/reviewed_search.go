package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/OpenUdon/apitools"
)

func checkedReviewedSearchMode(value string, offline bool) (apitools.CacheMode, error) {
	if offline {
		return apitools.CacheModeOffline, nil
	}
	mode := apitools.CacheMode(value)
	switch mode {
	case "", apitools.CacheModeReadWrite, apitools.CacheModeRefresh, apitools.CacheModeOffline, apitools.CacheModeBypass:
		return mode, nil
	default:
		return "", fmt.Errorf("unknown cache mode %q", value)
	}
}

func readReviewedSearchConfig(path string) (apitools.ReviewedProviderSearchOptions, error) {
	var options apitools.ReviewedProviderSearchOptions
	if path == "" {
		return options, nil
	}
	info, err := os.Stat(path)
	if err != nil {
		return options, err
	}
	const max = 64 << 10
	if !info.Mode().IsRegular() || info.Size() > max {
		return options, errors.New("reviewed search configuration must be a regular file under 64 KiB")
	}
	file, err := os.Open(path)
	if err != nil {
		return options, err
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, max+1))
	if err != nil {
		return options, err
	}
	if len(data) > max {
		return options, errors.New("reviewed search configuration exceeds 64 KiB")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var config *struct {
		IndexDB    string                            `json:"index_db"`
		StateDB    string                            `json:"state_db"`
		Embeddings *apitools.ReviewedEmbeddingConfig `json:"embeddings"`
	}
	if err := decoder.Decode(&config); err != nil {
		return options, fmt.Errorf("decode reviewed search configuration: %w", err)
	}
	if config == nil {
		return options, errors.New("reviewed search configuration must be a JSON object")
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return options, errors.New("reviewed search configuration contains trailing data")
	}
	dir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return options, err
	}
	for _, value := range []*string{&config.IndexDB, &config.StateDB} {
		if *value != "" && !filepath.IsAbs(*value) {
			*value = filepath.Join(dir, *value)
		}
	}
	options.IndexDB, options.StateDB = config.IndexDB, config.StateDB
	options.Embeddings = config.Embeddings
	return options, nil
}

func printReviewedSearchReport(out, errOut io.Writer, report apitools.ReviewedProviderSearchReport, jsonOut bool) int {
	if jsonOut {
		if err := writeJSON(out, report); err != nil {
			fmt.Fprintln(errOut, err)
			return exitRuntime
		}
		return exitSuccess
	}
	if report.Embeddings != nil {
		usage := report.Embeddings
		fmt.Fprintf(out, "Embeddings %s (%d requests, %d inputs, %d text bytes)\n",
			usage.Status, usage.Requests, usage.Inputs, usage.Bytes)
	}
	if len(report.Results) == 0 {
		fmt.Fprintln(out, "No reviewed providers matched.")
		return exitSuccess
	}
	if report.PublisherState != "unconfigured" {
		fmt.Fprintf(out, "Publisher state %s\n", report.PublisherState)
	}
	for i, result := range report.Results {
		profile := result.Profile
		fmt.Fprintf(out, "%d. %s (%s)\n", i+1, profile.DisplayName, profile.ProviderID)
		fmt.Fprintf(out, "   Category %s\n   Match %s\n   Coverage %s\n", profile.Category, result.Match, profile.Coverage.Status)
		fmt.Fprintf(out, "   %s\n", profile.Description)
		for _, capability := range result.MatchedCapabilities {
			fmt.Fprintf(out, "   Operation %s (%s)\n", capability.Label, capability.Selector)
			fmt.Fprintf(out, "   Operation source %s", capability.SpecRefID)
			if capability.ArtifactID != "" {
				fmt.Fprintf(out, " (%s)", capability.ArtifactID)
			}
			fmt.Fprintf(out, " sha256:%s\n", capability.SHA256)
		}
		for _, source := range profile.Coverage.Sources {
			fmt.Fprintf(out, "   Source %s %s %s\n", source.Kind, source.SpecRefID, source.SourceURL)
			reviewedAt := source.VerifiedAt
			if reviewedAt == "" {
				reviewedAt = "unknown"
			}
			fmt.Fprintf(out, "   Catalog review date %s\n", reviewedAt)
		}
		for _, check := range result.PublisherChecks {
			fmt.Fprintf(out, "   Accepted publisher check %s %s %s\n",
				check.SpecRefID, check.Version, check.CheckedAt.UTC().Format(time.RFC3339))
		}
	}
	return exitSuccess
}
