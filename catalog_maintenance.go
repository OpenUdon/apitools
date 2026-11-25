package apitools

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/catalogupdate"
)

// RunMonthlyCatalogUpdate runs or resumes the local calendar month's reviewed
// update. Publisher retries share the ordinary persisted per-reference limit.
func (c *Client) RunMonthlyCatalogUpdate(ctx context.Context, config CatalogUpdateConfig, onStart func(string)) (CatalogUpdateReport, error) {
	return c.RunMonthlyCatalogUpdateAfter(ctx, config, "", onStart)
}

// RunMonthlyCatalogUpdateAfter additionally supports a scheduled run's
// operator-reviewed prior-run prerequisite.
// If afterRun is supplied, its completed scope is checked inside the writer's
// recovered database before contacting any publisher. An invocation is capped
// at one day.
func (c *Client) RunMonthlyCatalogUpdateAfter(ctx context.Context, config CatalogUpdateConfig, afterRun string, onStart func(string)) (CatalogUpdateReport, error) {
	return c.runMonthlyCatalogUpdate(ctx, config, time.Now().Format("2006-01"), afterRun, onStart)
}

func (c *Client) runMonthlyCatalogUpdate(ctx context.Context, config CatalogUpdateConfig, period, afterRun string, onStart func(string)) (CatalogUpdateReport, error) {
	report, err := c.UpdateCatalog(ctx, CatalogUpdateOptions{All: true, DryRun: true})
	if err != nil {
		return report, err
	}
	report.DryRun = false
	if !filepath.IsAbs(config.CacheDir) || !filepath.IsAbs(config.StateDB) {
		return report, errors.New("catalog maintenance needs absolute cache_dir and state_db paths")
	}
	if afterRun != "" {
		if _, err := os.Lstat(config.StateDB); err != nil {
			return report, err
		}
	}
	identity, err := catalogUpdateRunIdentity(config, report.CatalogSHA256)
	if err != nil {
		return report, err
	}
	ctx, cancel := context.WithTimeout(ctx, catalogUpdateRetryWindow)
	defer cancel()
	cache, err := catalogupdate.OpenArtifactCache(config.CacheDir, config.ExpectedVolumeUUID, 0)
	if err != nil {
		return report, err
	}
	defer cache.Close()
	lock, err := catalogupdate.LockMonthlyRun(ctx, config.StateDB)
	if err != nil {
		return report, err
	}
	defer lock.Close()
	store, err := catalogupdate.OpenStore(ctx, config.StateDB, cache)
	if err != nil {
		return report, err
	}
	defer store.Close()
	if afterRun != "" {
		prior, err := store.ReadRun(ctx, afterRun)
		if err != nil {
			return report, err
		}
		if !prior.Completed || !slices.Equal(prior.Providers, report.Providers) {
			return report, errors.New("schedule requires a completed run covering the full reviewed provider list")
		}
	}
	manifest, providers := report.CatalogSHA256, report.Providers
	for {
		run, err := store.MonthlyRun(ctx, period, identity, manifest, providers)
		if err != nil {
			return report, err
		}
		report.setRun(run)
		report.Outcomes = nil
		if run.Completed {
			report.Outcomes, err = store.ReadOutcomes(ctx, run.ID)
			return report, err
		}
		if run.Catalog != identity {
			report.Outcomes, err = store.ReadOutcomes(ctx, run.ID)
			return report, errors.Join(catalogRunIdentityMismatch(run), err)
		}
		if onStart != nil {
			onStart(run.ID)
		}
		report, err = c.continueMonthlyCatalogUpdate(ctx, store, report, config)
		if err != nil || run.MonthlyPeriod == period {
			return report, err
		}
		// Finish carryover under this same lock, then service the requested month.
	}
}

func (c *Client) continueMonthlyCatalogUpdate(ctx context.Context, store *catalogupdate.Store, report CatalogUpdateReport, config CatalogUpdateConfig) (CatalogUpdateReport, error) {
	passes := 1
	for _, id := range report.Providers {
		if provider, ok := catalog.FindBuiltInProvider(id); ok {
			passes += catalogUpdatePublisherAttempts * len(catalogUpdateReferenceSources(provider))
		}
	}
	for range passes {
		var err error
		report, err = c.continueCatalogUpdate(ctx, store, report, false, config)
		if err != nil || report.Completed {
			return report, err
		}
		if delay := catalogMonthlyRetryDelay(report, time.Now()); delay > 0 {
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return report, ctx.Err()
			case <-timer.C:
			}
		}
	}
	return report, errors.New("catalog maintenance reached the bounded retry limit")
}

func catalogMonthlyRetryDelay(report CatalogUpdateReport, now time.Time) time.Duration {
	checked := make(map[string]bool)
	refs := make(map[[2]string]catalogupdate.Outcome)
	blocks := make(map[string]catalogUpdateOriginBlock)
	for _, outcome := range report.Outcomes {
		mergeCatalogUpdateRetryBlocks(blocks, outcome, false)
		if outcome.ReferenceID == "" {
			checked[outcome.ProviderID] = true
		} else {
			refs[[2]string{outcome.ProviderID, outcome.ReferenceID}] = outcome
		}
	}
	var next time.Time
	for _, id := range report.Providers {
		if checked[id] {
			continue
		}
		provider, ok := catalog.FindBuiltInProvider(id)
		if !ok {
			return 0
		}
		for _, source := range catalogUpdateReferenceSources(provider) {
			outcome, recorded := refs[[2]string{id, source.ReferenceID}]
			if recorded && outcome.Status != catalogupdate.RetryPendingStatus {
				continue
			}
			_, block := catalogUpdateBlockForSource(source, blocks)
			if block.OutsideWindow {
				return 0 // The writer can record a terminal deferral immediately.
			}
			eligible := readCatalogUpdateRetry(outcome).NextAttemptAt
			if block.Until.After(eligible) {
				eligible = block.Until
			}
			if !eligible.After(now) {
				return 0
			}
			if next.IsZero() || eligible.Before(next) {
				next = eligible
			}
		}
	}
	if next.IsZero() {
		return 0
	}
	return next.Sub(now)
}
