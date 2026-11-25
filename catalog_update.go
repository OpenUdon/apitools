package apitools

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"time"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/catalogupdate"
)

// CatalogUpdateConfig stores operator-selected paths and existing reviewed
// operation keys. Provider sources come from the built-in reviewed catalog.
type CatalogUpdateConfig struct {
	CacheDir           string              `json:"cache_dir"`
	StateDB            string              `json:"state_db"`
	ExpectedVolumeUUID string              `json:"expected_volume_uuid,omitempty"`
	RetainImmutable    bool                `json:"retain_immutable,omitempty"`
	ReviewedOperations map[string][]string `json:"reviewed_operations,omitempty"`
}

type CatalogUpdateOptions struct {
	Config   CatalogUpdateConfig
	Provider string
	All      bool
	Resume   string
	DryRun   bool
	Export   bool
	OnStart  func(runID string)
}

const catalogUpdateSchema = "apitools.catalog-update.v1"

type CatalogUpdateReport struct {
	Schema        string                  `json:"schema"`
	CatalogSHA256 string                  `json:"catalog_sha256,omitempty"`
	RunIdentity   string                  `json:"run_identity,omitempty"`
	RunID         string                  `json:"run_id,omitempty"`
	Providers     []string                `json:"providers,omitempty"`
	DryRun        bool                    `json:"dry_run,omitempty"`
	Completed     bool                    `json:"completed,omitempty"`
	StartedAt     time.Time               `json:"started_at,omitzero"`
	CompletedAt   time.Time               `json:"completed_at,omitzero"`
	MonthlyPeriod string                  `json:"monthly_period,omitempty"`
	SupersededAt  time.Time               `json:"superseded_at,omitzero"`
	Outcomes      []catalogupdate.Outcome `json:"outcomes,omitempty"`
	Export        json.RawMessage         `json:"export,omitempty"`
}

// UpdateCatalog checks the explicitly selected reviewed provider scope. A
// dry run resolves that scope without opening state or consulting publishers.
func (c *Client) UpdateCatalog(ctx context.Context, options CatalogUpdateOptions) (CatalogUpdateReport, error) {
	report := CatalogUpdateReport{Schema: catalogUpdateSchema, CatalogSHA256: catalog.BuiltInCatalogManifest().CatalogSHA256}
	selected := 0
	for _, setting := range []bool{options.All, options.Provider != "", options.Resume != "", options.Export} {
		if setting {
			selected++
		}
	}
	if selected != 1 || (options.DryRun && (options.Resume != "" || options.Export)) {
		return report, errors.New("select exactly one of --all, --provider, --resume or --export; dry-run supports --all and --provider")
	}
	if options.Provider != "" {
		provider, ok := catalog.FindBuiltInProvider(options.Provider)
		if !ok {
			return report, fmt.Errorf("unknown reviewed provider %q", options.Provider)
		}
		report.Providers = []string{provider.ID}
	} else if options.All {
		for _, provider := range catalog.BuiltInProviders() {
			report.Providers = append(report.Providers, provider.ID)
		}
	}
	if options.DryRun {
		report.DryRun = true
		return report, nil
	}
	config := options.Config
	if options.Export {
		if strings.TrimSpace(config.StateDB) == "" {
			return report, errors.New("catalog update export requires state_db")
		}
		export, err := catalogupdate.ReadVersionExport(ctx, config.StateDB)
		report.Export = export
		return report, err
	}
	if strings.TrimSpace(config.StateDB) == "" {
		return report, errors.New("catalog update requires state_db")
	}
	identity, err := catalogUpdateRunIdentity(config, report.CatalogSHA256)
	if err != nil {
		return report, err
	}
	catalogSHA256 := report.CatalogSHA256
	if options.Resume != "" {
		saved, outcomes, err := catalogupdate.ReadExistingRun(ctx, config.StateDB, options.Resume)
		switch {
		case errors.Is(err, catalogupdate.ErrRunJournalRecovery):
			if strings.TrimSpace(config.CacheDir) == "" {
				return report, fmt.Errorf("%w: resume with the maintenance configuration to recover it", err)
			}
			// The normal guarded writer open can roll back the journal. The
			// saved run identity is checked below before any publisher checks.
		case err != nil:
			return report, err
		default:
			report.setRun(saved)
			report.Outcomes = outcomes
			if saved.Completed || !saved.MonthlySupersededAt.IsZero() {
				return report, nil
			}
			if saved.Catalog != identity {
				return report, catalogRunIdentityMismatch(saved)
			}
		}
	}
	if strings.TrimSpace(config.CacheDir) == "" {
		return report, errors.New("catalog update requires cache_dir for new or unfinished runs")
	}
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
	var run catalogupdate.Run
	if options.Resume != "" {
		run, err = store.ReadRun(ctx, options.Resume)
	} else {
		run, err = store.StartRunSnapshot(ctx, identity, catalogSHA256, report.Providers)
	}
	if err != nil {
		return report, err
	}
	report.setRun(run)
	if !run.MonthlySupersededAt.IsZero() {
		report.Outcomes, err = store.ReadOutcomes(ctx, run.ID)
		return report, err
	}
	if !run.Completed && run.Catalog != identity {
		report.Outcomes, err = store.ReadOutcomes(ctx, run.ID)
		return report, errors.Join(catalogRunIdentityMismatch(run), err)
	}
	if options.OnStart != nil && !run.Completed {
		options.OnStart(run.ID)
	}
	return c.continueCatalogUpdate(ctx, store, report, run.Completed, config)
}

func catalogUpdateRunIdentity(config CatalogUpdateConfig, catalogSHA256 string) (string, error) {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha256.Sum256(append([]byte(catalogSHA256+"\n"), configBytes...))), nil
}

func (report *CatalogUpdateReport) setRun(run catalogupdate.Run) {
	report.RunID, report.Providers, report.RunIdentity = run.ID, run.Providers, run.Catalog
	report.Completed, report.StartedAt, report.CompletedAt = run.Completed, run.StartedAt, run.CompletedAt
	report.CatalogSHA256 = run.CatalogSHA256
	report.MonthlyPeriod, report.SupersededAt = run.MonthlyPeriod, run.MonthlySupersededAt
}

func catalogRunIdentityMismatch(run catalogupdate.Run) error {
	if run.MonthlyPeriod == "" {
		return errors.New("catalog or maintenance configuration changed since the update run started")
	}
	return fmt.Errorf("scheduled run %s for %s uses a different catalog or maintenance configuration; restore the original inputs or run `apitools catalog maintenance supersede --config FILE --run %s`, then retry maintenance run", run.ID, run.MonthlyPeriod, run.ID)
}

func (c *Client) continueCatalogUpdate(ctx context.Context, store *catalogupdate.Store, report CatalogUpdateReport, completed bool, config CatalogUpdateConfig) (result CatalogUpdateReport, resultErr error) {
	defer func() {
		if resultErr == nil {
			return
		}
		// A canceled publisher call does not undo already committed checkpoints.
		readCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
		defer cancel()
		outcomes, err := store.ReadOutcomes(readCtx, report.RunID)
		if err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("read committed checkpoints: %w", err))
			return
		}
		result.Outcomes = outcomes
	}()
	report.Completed = completed
	outcomes, err := store.ReadOutcomes(ctx, report.RunID)
	if err != nil {
		return report, err
	}
	report.Outcomes = outcomes
	if completed {
		return report, nil
	}
	checked := make(map[string]bool)
	type checkpoint struct{ provider, reference string }
	refsChecked := make(map[checkpoint]catalogupdate.Outcome)
	blockedUntil := make(map[string]catalogUpdateOriginBlock)
	for _, row := range outcomes {
		mergeCatalogUpdateRetryBlocks(blockedUntil, row, false)
		if row.ReferenceID == "" {
			checked[row.ProviderID] = true
		} else {
			refsChecked[checkpoint{row.ProviderID, row.ReferenceID}] = row
		}
	}
	pending := false
	for _, id := range report.Providers {
		if err := ctx.Err(); err != nil {
			return report, err
		}
		if checked[id] {
			continue
		}
		provider, ok := catalog.FindBuiltInProvider(id)
		if !ok {
			return report, fmt.Errorf("reviewed provider %q missing during update", id)
		}
		sources := catalogUpdateReferenceSources(provider)
		summary := catalogupdate.Outcome{ProviderID: id, Status: "checked"}
		if len(sources) == 0 {
			coverage := catalog.BuiltInAdvisoryEndpointCoverage(id)
			summary.Status = coverage.Status
			summary.Evidence, err = json.Marshal(coverage)
			if err != nil {
				return report, err
			}
		}
		providerPending := false
		for _, source := range sources {
			row, recorded := refsChecked[checkpoint{id, source.ReferenceID}]
			if recorded && row.Status != catalogupdate.RetryPendingStatus {
				continue
			}
			origin, block := catalogUpdateBlockForSource(source, blockedUntil)
			if block.OutsideWindow {
				details := catalogUpdateEvidence{Source: c.catalogUpdateEvidenceSource(source)}
				if recorded {
					if err := json.Unmarshal(row.Evidence, &details); err != nil {
						return report, err
					}
				}
				details.Diagnostic = "The publisher retry hint exceeds the 24-hour recovery window."
				if details.RetryBlocks == nil {
					details.RetryBlocks = make(map[string]catalogUpdateOriginBlock)
				}
				details.RetryBlocks[origin] = block
				data, err := json.Marshal(details)
				if err != nil {
					return report, err
				}
				deferred := catalogupdate.Outcome{ProviderID: id, ReferenceID: source.ReferenceID,
					Status: "publisher-retry-deferred", Evidence: data, Object: row.Object, PreviousRetryEvidence: row.Evidence}
				err = store.Record(ctx, report.RunID, deferred)
				if deferred.Object != nil && (errors.Is(err, fs.ErrNotExist) || errors.Is(err, catalogupdate.ErrObjectIntegrity)) {
					details.Diagnostic += " The cached bytes for the previously accepted artifact are unavailable."
					deferred.Object = nil // Keep the previous digest and adoption evidence without reopening the object.
					deferred.Evidence, err = json.Marshal(details)
					if err == nil {
						err = store.Record(ctx, report.RunID, deferred)
					}
				}
				if err != nil {
					return report, err
				}
				continue
			}
			now := time.Now().UTC()
			if now.Before(block.Until) ||
				(recorded && now.Before(readCatalogUpdateRetry(row).NextAttemptAt)) {
				providerPending = true
				continue
			}
			var prior *catalogupdate.Outcome
			if recorded {
				prior = &row
			}
			observed, err := c.checkCatalogUpdateSource(ctx, store, report.RunID, provider, source, config, prior)
			if err != nil {
				return report, fmt.Errorf("%s/%s: %w", id, source.ReferenceID, err)
			}
			mergeCatalogUpdateRetryBlocks(blockedUntil, observed, true)
			providerPending = providerPending || observed.Status == catalogupdate.RetryPendingStatus
		}
		if providerPending {
			pending = true
			continue
		}
		if err := store.Record(ctx, report.RunID, summary); err != nil {
			return report, err
		}
	}
	if !pending {
		if err := store.FinishRun(ctx, report.RunID); err != nil {
			return report, err
		}
		finished, err := store.ReadRun(ctx, report.RunID)
		if err != nil {
			return report, err
		}
		report.Completed, report.CompletedAt = finished.Completed, finished.CompletedAt
	}
	report.Outcomes, err = store.ReadOutcomes(ctx, report.RunID)
	return report, err
}

type catalogUpdateBaselineEvidence struct {
	Preferred catalogupdate.Version         `json:"version"`
	OpenAPI   *catalogUpdateOpenAPIEvidence `json:"openapi_validation,omitempty"`
	Native    *catalogUpdateNativeEvidence  `json:"native_validation,omitempty"`
	Artifact  *catalogUpdateArtifact        `json:"-"`
}

type catalogUpdateEvidence struct {
	Source           catalogUpdateSource                 `json:"source"`
	Artifact         *catalogUpdateArtifact              `json:"artifact,omitempty"`
	ReviewedBaseline *catalogUpdateBaselineEvidence      `json:"reviewed_baseline,omitempty"`
	OpenAPI          *catalogUpdateOpenAPIEvidence       `json:"openapi_validation,omitempty"`
	Native           *catalogUpdateNativeEvidence        `json:"native_validation,omitempty"`
	Inventory        *catalogUpdateInventory             `json:"inventory,omitempty"`
	Diagnostic       string                              `json:"diagnostic,omitempty"`
	Retry            *catalogUpdateRetry                 `json:"publisher_retry,omitempty"`
	RetryBlocks      map[string]catalogUpdateOriginBlock `json:"publisher_retry_blocks,omitempty"`
	AuthoredArtifact *catalogUpdateArtifact              `json:"authored_artifact,omitempty"`
	AttemptedAt      time.Time                           `json:"attempted_at,omitzero"`
	CheckedAt        time.Time                           `json:"checked_at,omitzero"`
}

func (c *Client) checkCatalogUpdateSource(ctx context.Context, store *catalogupdate.Store, runID string,
	provider catalog.Provider, source catalogUpdateSource, config CatalogUpdateConfig, prior *catalogupdate.Outcome,
) (catalogupdate.Outcome, error) {
	out := catalogupdate.Outcome{ProviderID: source.ProviderID, ReferenceID: source.ReferenceID}
	var priorEvidence catalogUpdateEvidence
	if prior != nil {
		out.PreviousRetryEvidence = prior.Evidence
		if err := json.Unmarshal(prior.Evidence, &priorEvidence); err != nil {
			return out, fmt.Errorf("read prior publisher retry evidence: %w", err)
		}
	}
	evidence := catalogUpdateEvidence{Source: c.catalogUpdateEvidenceSource(source)}
	var promotion, expected *catalogupdate.Version
	emit := func(status, diagnostic string) (catalogupdate.Outcome, error) {
		if status == "unchanged" && priorEvidence.ReviewedBaseline != nil {
			status = "reviewed-pin-reproduced"
		}
		out.Status = status
		evidence.Diagnostic = diagnostic
		evidence.CheckedAt = time.Now().UTC()
		evidence.AttemptedAt = evidence.Source.MetadataAttemptedAt
		for _, artifact := range []*catalogUpdateArtifact{evidence.Artifact, evidence.AuthoredArtifact} {
			if artifact != nil && !artifact.AttemptedAt.IsZero() &&
				(evidence.AttemptedAt.IsZero() || artifact.AttemptedAt.Before(evidence.AttemptedAt)) {
				evidence.AttemptedAt = artifact.AttemptedAt
			}
		}
		// A pending checkpoint may already have atomically accepted an
		// authored pin. Keep that proof separate from the latest head check.
		if evidence.ReviewedBaseline == nil && priorEvidence.ReviewedBaseline != nil {
			evidence.ReviewedBaseline = priorEvidence.ReviewedBaseline
			evidence.AuthoredArtifact = priorEvidence.AuthoredArtifact
		}
		var err error
		out.Evidence, err = json.Marshal(evidence)
		if err == nil {
			if promotion != nil {
				err = store.AcceptDirectVersionAndRecord(ctx, *promotion, expected, runID, out)
			} else {
				err = store.Record(ctx, runID, out)
			}
		}
		if err == nil && out.Object == nil && prior != nil {
			out.Object = prior.Object
		}
		return out, err
	}
	emitPublisher := func(status, diagnostic string, failures ...error) (catalogupdate.Outcome, error) {
		for _, failed := range failures {
			var publisher *catalogUpdatePublisherError
			if !errors.As(failed, &publisher) {
				continue
			}
			retry := publisher.nextAttempt(prior)
			if retry.Shared {
				if evidence.RetryBlocks == nil {
					evidence.RetryBlocks = make(map[string]catalogUpdateOriginBlock)
				}
				block := evidence.RetryBlocks[retry.Origin]
				if retry.OutsideWindow || !block.OutsideWindow && retry.NextAttemptAt.After(block.Until) {
					evidence.RetryBlocks[retry.Origin] = catalogUpdateOriginBlock{Until: retry.NextAttemptAt, OutsideWindow: retry.OutsideWindow}
				}
			}
			if evidence.Retry == nil || retry.NextAttemptAt.After(evidence.Retry.NextAttemptAt) {
				evidence.Retry = &retry
			}
		}
		if evidence.Retry != nil {
			status = catalogupdate.RetryPendingStatus
			if evidence.Retry.OutsideWindow || evidence.Retry.Attempt >= catalogUpdatePublisherAttempts {
				status = "publisher-retry-exhausted"
			}
		}
		return emit(status, diagnostic)
	}
	if source.Kind == catalog.SpecKindDropboxStone {
		return emit("publisher-source-review", "Check the official Dropbox OpenAPI before substituting Stone provenance.")
	}
	kind := OperationSourceKind(catalog.RefreshableSpecReference{Kind: source.Kind}.ProtocolClassification().UWSSourceType())
	var held catalogupdate.Version
	if source.Kind != catalog.SpecKindOpenAPIIndex {
		var err error
		held, err = store.PreferredVersion(ctx, source.ProviderID, source.ReferenceID, source.Service)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return out, err
		}
		if held.Object.SHA256 == "" {
			held, err = store.ReviewedReferencePreference(ctx, source.ProviderID, source.ReferenceID)
			switch {
			case errors.Is(err, sql.ErrNoRows):
			case errors.Is(err, catalogupdate.ErrAmbiguousReferencePreference):
				return emit("review-required", "The reviewed reference needs an explicit service match.")
			case err != nil:
				return out, err
			}
		}
	}
	previous := held
	catalogPin := catalogUpdateReviewedDigest(source.HeldRevision)
	authoredCommit, _ := catalogUpdateGitFile(source.URL)
	changeBaseline := held.Object.SHA256 == "" || held.CatalogPin != catalogPin || held.Locator != source.URL
	resourceMoved := held.Object.SHA256 != "" && held.Locator != source.URL &&
		!sameCatalogUpdateGitFile(held.Locator, source.URL) && !sameCatalogUpdateGitFile(held.URL, source.URL)
	_, heldGitFile := catalogUpdateGitFile(held.URL)
	canReproduceMovedPin := !resourceMoved || (authoredCommit != "" &&
		(catalogPin != held.CatalogPin || heldGitFile != "" && held.Object.SHA256 == catalogPin))
	setPromotion := func(version *catalogupdate.Version) {
		promotion = version
		if previous.Object.SHA256 != "" {
			expected = &previous
		}
	}
	// A failed head check does not invalidate separately observed reviewed
	// bytes. Keep the head status for retry and attach the accepted baseline.
	recordHeadFailure := func(status, diagnostic, target string, failed error) (catalogupdate.Outcome, error) {
		if ctx.Err() != nil {
			return out, ctx.Err()
		}
		var authoredFailure error
		if changeBaseline && catalogPin != "" && canReproduceMovedPin && validOperationSourceKind(kind) {
			var headPublisher *catalogUpdatePublisherError
			if errors.As(failed, &headPublisher) && headPublisher.retry.Shared &&
				catalogUpdateArtifactOrigin(source.URL) == headPublisher.retry.Origin {
				// Skip another request to the same artifact origin. Owned
				// cached bytes can still reproduce the authored locator.
				target = source.URL
			}
			anchored, err := c.catalogUpdateReviewedBaseline(ctx, store, source, held, target, catalogPin)
			evidence.AuthoredArtifact = anchored.Artifact
			var publisher *catalogUpdatePublisherError
			if errors.As(err, &publisher) {
				authoredFailure = err
			} else if err != nil && !errors.Is(err, catalogupdate.ErrObjectIntegrity) {
				return out, err
			}
			if anchored.Preferred.Object.SHA256 != "" {
				evidence.ReviewedBaseline = &anchored
				out.Object = &anchored.Preferred.Object
				setPromotion(&anchored.Preferred)
			}
		}
		return emitPublisher(status, diagnostic, failed, authoredFailure)
	}
	resolved, err := c.resolveCatalogUpdateSource(ctx, source)
	if held.Service != "" {
		resolved.Service = held.Service
	}
	evidence.Source = c.catalogUpdateEvidenceSource(resolved)
	if err != nil || resolved.Resolution == "unexamined-locator" {
		return recordHeadFailure("publisher-source-unexamined", "Could not establish the publisher artifact locator.", "", err)
	}
	target := firstNonEmpty(resolved.ResolvedURL, resolved.URL)
	if resolved.ResolvedURL != "" && target == held.URL && resolved.Revision == held.Revision &&
		held.Locator == source.URL && (catalogPin == "" || held.CatalogPin == catalogPin) {
		file, err := store.OpenObject(ctx, held.Object)
		if err == nil {
			if err := file.Close(); err != nil {
				return out, err
			}
			out.Object = &held.Object
			evidence.Artifact = &catalogUpdateArtifact{
				URL: target, Bytes: held.Object.Bytes, SHA256: held.Object.SHA256, CheckedAt: time.Now().UTC(),
				Status: "verified-cache-hit",
			}
			return emit("unchanged", "")
		}
	}
	var body bytes.Buffer
	artifact, err := c.acquireCatalogUpdateArtifact(ctx, target, DefaultMaxBytes, &body)
	evidence.Artifact = &artifact
	if ctx.Err() != nil {
		return out, ctx.Err()
	}
	if err != nil {
		return recordHeadFailure("publisher-check-failed", "Publisher did not return a matching downloadable artifact.", target, err)
	}
	if artifact.Status != catalogUpdateDownloaded {
		return recordHeadFailure(artifact.Status, "", target, nil)
	}
	if resolved.ResolvedURL != "" && artifact.URL != resolved.ResolvedURL {
		return recordHeadFailure("publisher-source-unexamined", "The artifact redirected away from the resolved immutable locator.", target, nil)
	}
	if source.Kind == catalog.SpecKindOpenAPIIndex {
		inventory, err := hubSpotUpdateInventory(artifact.URL, artifact.CheckedAt, body.Bytes())
		if err != nil {
			return emit("publisher-index-unexamined", "The publisher index did not match the implemented reviewed format.")
		}
		evidence.Inventory = &inventory
		return emit("publisher-index-inventoried", "")
	}
	if !validOperationSourceKind(kind) {
		return emit("publisher-source-review", "The current artifact kind needs a reviewed source recipe.")
	}
	if kind == OperationSourceOpenAPI {
		validation := catalogUpdateNativeOpenAPI(ctx, body.Bytes(), artifact.URL, source.refreshReference(), catalogPin)
		if ctx.Err() != nil {
			return out, ctx.Err()
		}
		evidence.OpenAPI = &validation
		if !validation.accepted() {
			return recordHeadFailure("publisher-spec-invalid", "Local validation established neither a valid OpenAPI document nor evidence for the exact reviewed artifact.", target, nil)
		}
		resolved.Version = firstNonEmpty(validation.Version, resolved.Version)
	} else {
		validation, ok := qualifyCatalogUpdateNativeSource(ctx, source, kind, artifact.URL, artifact.SHA256, body.Bytes())
		if !ok {
			return recordHeadFailure("publisher-spec-unexamined", "Local source validation did not establish a native operation inventory.", target, nil)
		}
		evidence.Native = validation
	}
	object := catalogupdate.Object{SHA256: artifact.SHA256, Bytes: artifact.Bytes}
	resolved.Revision = firstNonEmpty(resolved.Revision, "sha256:"+artifact.SHA256)
	if kind == OperationSourceGoogleDiscovery {
		native, err := decodeShapeDocument(ctx, body.Bytes())
		if err != nil {
			return recordHeadFailure("publisher-spec-unexamined", "Native source provenance could not be read within the validation budget.", target, nil)
		}
		resolved.Version = firstNonEmpty(stringValue(native["version"]), resolved.Version)
		resolved.NativeRevision = stringValue(native["revision"])
		if resolved.ResolvedURL == "" {
			resolved.Revision = firstNonEmpty(resolved.NativeRevision, resolved.Revision)
		}
	}
	evidence.Source = c.catalogUpdateEvidenceSource(resolved)
	matchesReviewedPin := catalogPin == object.SHA256
	if !canReproduceMovedPin && !matchesReviewedPin {
		return emit("review-required", "The publisher did not reproduce the reviewed pin at the new locator.")
	}
	if matchesReviewedPin && held.CatalogPin == catalogPin && held.Object != object {
		return emit("review-required", "Publisher returned the standing catalog pin after a newer local acceptance.")
	}
	if matchesReviewedPin {
		allowed, err := catalogUpdateCanReplaceWithPin(ctx, store, held, object, catalogPin, resolved.Version)
		if err != nil {
			return out, err
		}
		if !allowed {
			return emit("review-required", "The reviewed artifact would replay an earlier pin or replace a newer release.")
		}
	}
	adoptReviewedPin := matchesReviewedPin && changeBaseline && (authoredCommit == "" || artifact.URL == source.URL)
	key := source.ProviderID + "/" + source.ReferenceID
	var adoptedBaseline bool
	if !adoptReviewedPin && catalogPin != "" && changeBaseline {
		anchored, err := c.catalogUpdateReviewedBaseline(ctx, store, source, held, target, catalogPin)
		evidence.AuthoredArtifact = anchored.Artifact
		var publisher *catalogUpdatePublisherError
		switch {
		case errors.As(err, &publisher):
			return emitPublisher("review-required", "The authored immutable artifact could not be checked.", err)
		case errors.Is(err, catalogupdate.ErrObjectIntegrity):
			return emit("review-required", "The reviewed baseline cache entry could not be safely restored.")
		case err != nil:
			return out, err
		case anchored.Preferred.Object.SHA256 != "":
			held = anchored.Preferred
			adoptedBaseline = true
			evidence.ReviewedBaseline = &anchored
		case held.Object.SHA256 != "":
			return emit("review-required", "The publisher did not reproduce the new reviewed catalog pin.")
		}
	}
	// A reproduced baseline can be committed with this checkpoint even when
	// the separately observed publisher head cannot be accepted.
	reviewCandidate := func(reason string) (catalogupdate.Outcome, error) {
		if !adoptedBaseline {
			return emit("review-required", reason)
		}
		setPromotion(&held)
		out.Object = &held.Object
		return emit("reviewed-pin-reproduced", reason)
	}
	if config.RetainImmutable || resolved.ResolvedURL == "" || held.Object.SHA256 != "" || adoptReviewedPin {
		if err := store.PutObject(ctx, object, bytes.NewReader(body.Bytes())); err != nil {
			if errors.Is(err, catalogupdate.ErrObjectIntegrity) {
				return reviewCandidate("The candidate's cache entry could not be safely restored.")
			}
			return out, err
		}
		out.Object = &object
	}
	version := catalogupdate.Version{
		ProviderID: source.ProviderID, ReferenceID: source.ReferenceID, Service: resolved.Service,
		Locator: source.URL, URL: artifact.URL, Version: resolved.Version, Revision: resolved.Revision,
		CatalogPin: held.CatalogPin, CheckedAt: artifact.CheckedAt, Object: object,
	}
	if adoptReviewedPin {
		version.CatalogPin = catalogPin
		setPromotion(&version)
		return emit("reviewed-pin-reproduced", "")
	}
	if adoptedBaseline && held.Object == object {
		return reviewCandidate("")
	}
	if held.Object.SHA256 == "" {
		return emit("review-required", "No previously accepted local version is available for comparison.")
	}
	if held.Object == object && held.Locator == source.URL &&
		((held.URL == artifact.URL && held.Revision == resolved.Revision) ||
			(held.URL != artifact.URL && sameCatalogUpdateGitFile(held.URL, artifact.URL))) {
		// The outcome records the observed provenance. The stored
		// preference still describes the already accepted identical bytes.
		return emit("unchanged", "")
	}
	if previous.Object.SHA256 != "" && previous.Object != object {
		replay, err := store.HasAcceptedCatalogDigest(ctx, previous, object.SHA256)
		if err != nil {
			return out, err
		}
		if replay {
			return reviewCandidate("The publisher returned an earlier accepted artifact.")
		}
	}
	if !catalogUpdateVersionCanAdvance(held.Version, version.Version) {
		return reviewCandidate("Publisher version stability or order needs review.")
	}
	if order, known := catalogUpdateVersionOrder(previous.Version, version.Version); adoptedBaseline && known && order < 0 {
		return reviewCandidate("The publisher head is older than the previous local preference.")
	}
	if hasCatalogUpdateSecurityOverlay(provider.ID, source.ReferenceID) {
		return reviewCandidate("The source has advisory security evidence that needs separate review.")
	}
	comparison, err := c.proveCompatibleCatalogVersion(ctx, store, held, version, kind, config.ReviewedOperations[key])
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) || errors.Is(err, catalogupdate.ErrObjectIntegrity) {
			return reviewCandidate("The previously accepted artifact is unavailable for compatibility comparison.")
		}
		return out, err
	}
	if comparison == catalogUpdateCompatible {
		setPromotion(&version)
		return emit("accepted-compatible", "")
	}
	return reviewCandidate(string(comparison))
}

func catalogUpdateVersionCanAdvance(held, candidate string) bool {
	order, known := catalogUpdateVersionOrder(held, candidate)
	return known && order >= 0
}

func catalogUpdateVersionOrder(held, candidate string) (int, bool) {
	if held == "" || candidate == "" {
		return 0, false
	}
	if _, stable := numericAPIVersion(candidate); stable {
		return compareAPIVersions(candidate, held)
	}
	for _, layout := range []string{"2006-01-02", "2006-01"} {
		before, errBefore := time.Parse(layout, held)
		after, errAfter := time.Parse(layout, candidate)
		if errBefore == nil && errAfter == nil {
			return after.Compare(before), true
		}
	}
	return 0, false
}

func hasCatalogUpdateReviewedDigest(revision, digest string) bool {
	return digest != "" && catalogUpdateReviewedDigest(revision) == digest
}

func catalogUpdateReviewedDigest(revision string) string {
	var result string
	for _, claim := range strings.Fields(revision) {
		digest, ok := strings.CutPrefix(claim, "sha256:")
		if !ok {
			continue
		}
		raw, err := hex.DecodeString(digest)
		if err != nil || len(raw) != sha256.Size || digest != strings.ToLower(digest) ||
			(result != "" && result != digest) {
			return ""
		}
		result = digest
	}
	return result
}

// catalogUpdateReviewedBaseline verifies a changed authored pin without
// changing the preference. The caller commits it with the outcome, or uses
// it to prove the current publisher artifact and commits that one instead.
func (c *Client) catalogUpdateReviewedBaseline(ctx context.Context, store *catalogupdate.Store, source catalogUpdateSource,
	held catalogupdate.Version, target, pin string,
) (catalogUpdateBaselineEvidence, error) {
	if _, err := c.catalogUpdateURL(source.URL); err != nil {
		return catalogUpdateBaselineEvidence{}, nil
	}
	authoredCommit, _ := catalogUpdateGitFile(source.URL)
	limit := int64(DefaultMaxBytes)
	if configured := c.effective().MaxBytes; configured > 0 && configured < limit {
		limit = configured
	}
	if held.Object.SHA256 == pin && held.Object.Bytes <= limit &&
		((authoredCommit == "" && held.Locator == source.URL) ||
			(authoredCommit != "" && held.URL == source.URL && held.Revision == authoredCommit)) {
		file, err := store.OpenObject(ctx, held.Object)
		if err == nil {
			raw, err := io.ReadAll(io.LimitReader(file, held.Object.Bytes+1))
			closeErr := file.Close()
			if err != nil {
				return catalogUpdateBaselineEvidence{}, err
			}
			if closeErr != nil {
				return catalogUpdateBaselineEvidence{}, closeErr
			}
			sum := sha256.Sum256(raw)
			if int64(len(raw)) == held.Object.Bytes && hex.EncodeToString(sum[:]) == pin {
				held.Locator, held.CatalogPin = source.URL, pin
				evidence, valid := qualifyCatalogUpdateBaseline(ctx, source, held, raw)
				if ctx.Err() != nil {
					return catalogUpdateBaselineEvidence{}, ctx.Err()
				}
				if valid {
					return evidence, nil
				}
			}
		}
		if err != nil && !errors.Is(err, fs.ErrNotExist) && !errors.Is(err, catalogupdate.ErrObjectIntegrity) {
			return catalogUpdateBaselineEvidence{}, err
		}
	}
	if authoredCommit == "" || source.URL == target ||
		(held.CatalogPin == pin && held.Object.SHA256 != pin) {
		return catalogUpdateBaselineEvidence{}, nil
	}
	if held.Service != "" {
		source.Service = held.Service
	}
	authored, err := c.reconstructCatalogUpdateBaseline(ctx, store, source)
	if err != nil || authored.Preferred.Object.SHA256 == "" {
		return authored, err
	}
	allowed, err := catalogUpdateCanReplaceWithPin(ctx, store, held, authored.Preferred.Object, pin, authored.Preferred.Version)
	if err != nil || !allowed {
		return catalogUpdateBaselineEvidence{Artifact: authored.Artifact}, err
	}
	return authored, nil
}

func catalogUpdateCanReplaceWithPin(ctx context.Context, store *catalogupdate.Store, held catalogupdate.Version,
	authored catalogupdate.Object, pin, version string,
) (bool, error) {
	if held.Object.SHA256 == "" || held.Object == authored {
		return true, nil
	}
	order, known := catalogUpdateVersionOrder(held.Version, version)
	// Migrated v2 preferences have no reviewed anchor. Require evidence of
	// a strictly newer release before replacing their bytes with a new pin.
	if held.CatalogPin == "" {
		return known && order > 0, nil
	}
	if known && order < 0 {
		return false, nil
	}
	replay, err := store.HasAcceptedCatalogDigest(ctx, held, pin)
	return !replay, err
}

func (c *Client) reconstructCatalogUpdateBaseline(ctx context.Context, store *catalogupdate.Store, source catalogUpdateSource) (catalogUpdateBaselineEvidence, error) {
	commit, _ := catalogUpdateGitFile(source.URL)
	kind := OperationSourceKind(source.refreshReference().ProtocolClassification().UWSSourceType())
	if commit == "" || !validOperationSourceKind(kind) || catalogUpdateReviewedDigest(source.HeldRevision) == "" {
		return catalogUpdateBaselineEvidence{}, nil
	}
	var body bytes.Buffer
	artifact, err := c.acquireCatalogUpdateArtifact(ctx, source.URL, DefaultMaxBytes, &body)
	if ctx.Err() != nil {
		return catalogUpdateBaselineEvidence{Artifact: &artifact}, ctx.Err()
	}
	var publisher *catalogUpdatePublisherError
	if errors.As(err, &publisher) {
		return catalogUpdateBaselineEvidence{Artifact: &artifact}, err
	}
	if err != nil || artifact.Status != catalogUpdateDownloaded || artifact.URL != source.URL ||
		!hasCatalogUpdateReviewedDigest(source.HeldRevision, artifact.SHA256) {
		return catalogUpdateBaselineEvidence{Artifact: &artifact}, nil
	}
	object := catalogupdate.Object{SHA256: artifact.SHA256, Bytes: artifact.Bytes}
	baseline := catalogupdate.Version{
		ProviderID: source.ProviderID, ReferenceID: source.ReferenceID, Service: source.Service,
		Locator: source.URL, URL: source.URL, Version: firstNonEmpty(source.Version, source.HeldVersion),
		Revision: commit, CatalogPin: artifact.SHA256, CheckedAt: artifact.CheckedAt, Object: object,
	}
	evidence, valid := qualifyCatalogUpdateBaseline(ctx, source, baseline, body.Bytes())
	if ctx.Err() != nil {
		return catalogUpdateBaselineEvidence{Artifact: &artifact}, ctx.Err()
	}
	if !valid {
		return catalogUpdateBaselineEvidence{Artifact: &artifact}, nil
	}
	if err := store.PutObject(ctx, object, bytes.NewReader(body.Bytes())); err != nil {
		return catalogUpdateBaselineEvidence{Artifact: &artifact}, err
	}
	evidence.Artifact = &artifact
	return evidence, nil
}

func qualifyCatalogUpdateBaseline(ctx context.Context, source catalogUpdateSource, version catalogupdate.Version, body []byte) (catalogUpdateBaselineEvidence, bool) {
	evidence := catalogUpdateBaselineEvidence{Preferred: version}
	kind := OperationSourceKind(source.refreshReference().ProtocolClassification().UWSSourceType())
	if kind == OperationSourceOpenAPI {
		validation := catalogUpdateNativeOpenAPI(ctx, body, version.URL, source.refreshReference(), version.Object.SHA256)
		evidence.Preferred.Version = firstNonEmpty(validation.Version, version.Version)
		evidence.OpenAPI = &validation
		return evidence, validation.accepted()
	}
	validation, ok := qualifyCatalogUpdateNativeSource(ctx, source, kind, version.URL, version.Object.SHA256, body)
	if !ok {
		return evidence, false
	}
	evidence.Native = validation
	if kind == OperationSourceGoogleDiscovery {
		native, err := decodeShapeDocument(ctx, body)
		if err != nil {
			return evidence, false
		}
		evidence.Preferred.Version = firstNonEmpty(stringValue(native["version"]), version.Version)
	}
	return evidence, true
}

// Reports retain safe authored and resolved locators even when a check
// fails before a publisher returns an artifact.
func (c *Client) catalogUpdateEvidenceSource(source catalogUpdateSource) catalogUpdateSource {
	if source.URL != "" {
		if _, err := c.catalogUpdateURL(source.URL); err != nil {
			source.URL, source.Scope = "", ""
		}
	}
	if source.ResolvedURL != "" {
		if _, err := c.catalogUpdateURL(source.ResolvedURL); err != nil {
			source.ResolvedURL = ""
		}
	}
	return source
}

func hasCatalogUpdateSecurityOverlay(provider, reference string) bool {
	for _, overlay := range catalog.BuiltInSecurityOverlays() {
		if overlay.ProviderID == provider && (overlay.SpecRefID == "" || overlay.SpecRefID == reference) {
			return true
		}
	}
	return false
}
