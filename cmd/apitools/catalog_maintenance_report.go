package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/apitools/internal/catalogupdate"
)

type catalogMaintenanceSummary struct {
	Providers          int            `json:"providers"`
	ProvidersChecked   int            `json:"providers_checked"`
	AcceptedCompatible int            `json:"accepted_compatible"`
	ReviewedPins       int            `json:"reviewed_pins"`
	PendingReview      int            `json:"pending_review"`
	PendingRetries     int            `json:"pending_retries"`
	SizeDeferrals      int            `json:"size_deferrals"`
	DeferredRetries    int            `json:"deferred_retries"`
	Failures           int            `json:"failures"`
	ByStatus           map[string]int `json:"by_status"`
}

type catalogMaintenanceReport struct {
	Run            apitools.CatalogUpdateReport      `json:"run"`
	Summary        catalogMaintenanceSummary         `json:"summary"`
	CacheUsage     *catalogupdate.ArtifactCacheUsage `json:"cache_usage,omitempty"`
	CacheError     string                            `json:"cache_error,omitempty"`
	EmbeddingUsage apitools.ReviewedEmbeddingUsage   `json:"embedding_usage"`
	Error          string                            `json:"error,omitempty"`
}

type catalogMaintenanceReferenceFacts struct {
	reviewedPin bool
	needsReview bool
	nextRetry   time.Time
}

func loadCatalogMaintenanceReport(ctx context.Context, statePath, id string) (apitools.CatalogUpdateReport, error) {
	var report apitools.CatalogUpdateReport
	if statePath == "" {
		return report, errors.New("catalog maintenance report requires state_db")
	}
	if id == "" {
		var err error
		id, err = catalogupdate.ReadLatestMonthlyRunID(ctx, statePath)
		if err != nil {
			return report, err
		}
		if id == "" {
			return report, errors.New("no scheduled run exists; use --run to inspect a saved update")
		}
	}
	run, outcomes, err := catalogupdate.ReadExistingRun(ctx, statePath, id)
	if err != nil {
		return report, err
	}
	return apitools.CatalogUpdateReport{Schema: "apitools.catalog-update.v1", CatalogSHA256: run.CatalogSHA256,
		RunIdentity: run.Catalog, RunID: run.ID, Providers: run.Providers, Completed: run.Completed,
		StartedAt: run.StartedAt, CompletedAt: run.CompletedAt, Outcomes: outcomes}, nil
}

func writeCatalogMaintenanceReport(ctx context.Context, out, errOut io.Writer, run apitools.CatalogUpdateReport, runErr error, cfg apitools.CatalogUpdateConfig, jsonOutput bool) int {
	report := catalogMaintenanceReport{Run: run, EmbeddingUsage: apitools.ReviewedEmbeddingUsage{Status: "not-requested"}}
	if runErr != nil {
		report.Error = runErr.Error()
	}
	summary := &report.Summary
	summary.Providers = len(run.Providers)
	summary.ByStatus = make(map[string]int)
	facts := make([]catalogMaintenanceReferenceFacts, len(run.Outcomes))
	for i, outcome := range run.Outcomes {
		summary.ByStatus[outcome.Status]++
		if outcome.ReferenceID == "" {
			summary.ProvidersChecked++
			continue
		}
		var evidence struct {
			ReviewedBaseline *struct{} `json:"reviewed_baseline"`
			Diagnostic       string    `json:"diagnostic"`
			Retry            struct {
				Next time.Time `json:"next_attempt_at"`
			} `json:"publisher_retry"`
		}
		_ = json.Unmarshal(outcome.Evidence, &evidence)
		fact := &facts[i]
		fact.reviewedPin = outcome.Status == "reviewed-pin-reproduced" || evidence.ReviewedBaseline != nil
		fact.needsReview = outcome.Status == "reviewed-pin-reproduced" && evidence.Diagnostic != ""
		fact.nextRetry = evidence.Retry.Next
		switch outcome.Status {
		case "accepted-compatible":
			summary.AcceptedCompatible++
		case "review-required", "publisher-source-review", "publisher-source-unexamined", "publisher-index-unexamined", "publisher-spec-unexamined":
			fact.needsReview = true
		case catalogupdate.RetryPendingStatus:
			summary.PendingRetries++
		case "size-unknown", "over-download-limit", "over-parser-limit", "size-changed":
			summary.SizeDeferrals++
		case "publisher-retry-deferred":
			summary.DeferredRetries++
		case "publisher-check-failed", "publisher-retry-exhausted", "publisher-spec-invalid":
			summary.Failures++
		}
		if fact.reviewedPin {
			summary.ReviewedPins++
		}
		if fact.needsReview {
			summary.PendingReview++
		}
	}
	if cfg.CacheDir != "" {
		usage, err := catalogupdate.MeasureArtifactCache(ctx, cfg.CacheDir, cfg.ExpectedVolumeUUID)
		if err == nil {
			report.CacheUsage = &usage
		} else {
			report.CacheError = err.Error()
		}
	}
	if jsonOutput {
		if err := writeJSON(out, report); err != nil {
			fmt.Fprintln(errOut, err)
			return exitRuntime
		}
	} else if run.RunID != "" {
		state := "unfinished"
		if run.Completed {
			state = "complete"
		}
		fmt.Fprintf(out, "Run %s is %s\nProviders checked %d/%d\n", run.RunID, state, summary.ProvidersChecked, summary.Providers)
		fmt.Fprintf(out, "Accepted compatible %d\nReviewed pins %d\nPending review %d\nPending retries %d\nSize deferrals %d\nDeferred retries %d\nFailures %d\n",
			summary.AcceptedCompatible, summary.ReviewedPins, summary.PendingReview, summary.PendingRetries, summary.SizeDeferrals, summary.DeferredRetries, summary.Failures)
		if report.CacheUsage != nil {
			fmt.Fprintf(out, "Artifact cache %d bytes across %d objects\n", report.CacheUsage.Bytes, report.CacheUsage.Objects)
		} else if report.CacheError != "" {
			fmt.Fprintln(out, "Cache usage unavailable.", report.CacheError)
		}
		fmt.Fprintln(out, "Embeddings not requested")
		for i, outcome := range run.Outcomes {
			if outcome.ReferenceID == "" {
				fmt.Fprintf(out, "%s %s\n", outcome.ProviderID, outcome.Status)
				continue
			}
			fmt.Fprintf(out, "%s/%s %s", outcome.ProviderID, outcome.ReferenceID, outcome.Status)
			fact := facts[i]
			if fact.reviewedPin && outcome.Status != "reviewed-pin-reproduced" {
				fmt.Fprint(out, " (reviewed baseline adopted)")
			}
			if fact.needsReview {
				fmt.Fprint(out, " (publisher head needs review)")
			}
			if outcome.Status == catalogupdate.RetryPendingStatus && !fact.nextRetry.IsZero() {
				fmt.Fprint(out, " after ", fact.nextRetry.UTC().Format(time.RFC3339))
			}
			fmt.Fprintln(out)
		}
	}
	if runErr != nil {
		fmt.Fprintln(errOut, runErr)
		return exitRuntime
	}
	return exitSuccess
}
