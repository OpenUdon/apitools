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
	CatalogSHA256 string                  `json:"catalog_sha256"`
	RunID         string                  `json:"run_id,omitempty"`
	Providers     []string                `json:"providers,omitempty"`
	DryRun        bool                    `json:"dry_run,omitempty"`
	Completed     bool                    `json:"completed,omitempty"`
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
	if strings.TrimSpace(config.StateDB) == "" || strings.TrimSpace(config.CacheDir) == "" {
		return report, errors.New("catalog update requires cache_dir and state_db")
	}
	cache, err := catalogupdate.OpenArtifactCache(config.CacheDir, config.ExpectedVolumeUUID, 0)
	if err != nil {
		return report, err
	}
	defer cache.Close()
	store, err := catalogupdate.OpenStore(ctx, config.StateDB, cache)
	if err != nil {
		return report, err
	}
	defer store.Close()
	var run catalogupdate.Run
	configBytes, err := json.Marshal(config)
	if err != nil {
		return report, err
	}
	identity := fmt.Sprintf("%x", sha256.Sum256(append([]byte(report.CatalogSHA256+"\n"), configBytes...)))
	if options.Resume != "" {
		run, err = store.ReadRun(ctx, options.Resume)
	} else {
		run, err = store.StartRun(ctx, identity, report.Providers)
	}
	if err != nil {
		return report, err
	}
	report.RunID, report.Providers = run.ID, run.Providers
	if run.Catalog != identity {
		return report, errors.New("catalog or maintenance configuration changed since the update run started")
	}
	if options.OnStart != nil && !run.Completed {
		options.OnStart(run.ID)
	}
	return c.continueCatalogUpdate(ctx, store, report, run.Completed, config)
}

func (c *Client) continueCatalogUpdate(ctx context.Context, store *catalogupdate.Store, report CatalogUpdateReport, completed bool, config CatalogUpdateConfig) (CatalogUpdateReport, error) {
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
	refsChecked := make(map[checkpoint]bool)
	for _, row := range outcomes {
		if row.ReferenceID == "" {
			checked[row.ProviderID] = true
		} else {
			refsChecked[checkpoint{row.ProviderID, row.ReferenceID}] = true
		}
	}
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
		status := "checked"
		if len(sources) == 0 {
			status = "needs-publisher-spec-or-overlay"
		}
		for _, source := range sources {
			if refsChecked[checkpoint{id, source.ReferenceID}] {
				continue
			}
			_, err := c.checkCatalogUpdateSource(ctx, store, report.RunID, provider, source, config)
			if err != nil {
				return report, fmt.Errorf("%s/%s: %w", id, source.ReferenceID, err)
			}
		}
		if err := store.Record(ctx, report.RunID, catalogupdate.Outcome{ProviderID: id, Status: status}); err != nil {
			return report, err
		}
	}
	if err := store.FinishRun(ctx, report.RunID); err != nil {
		return report, err
	}
	report.Completed = true
	report.Outcomes, err = store.ReadOutcomes(ctx, report.RunID)
	return report, err
}

type catalogUpdateBaselineEvidence struct {
	Preferred catalogupdate.Version         `json:"version"`
	OpenAPI   *catalogUpdateOpenAPIEvidence `json:"openapi_validation,omitempty"`
}

type catalogUpdateEvidence struct {
	Source           catalogUpdateSource            `json:"source"`
	Artifact         *catalogUpdateArtifact         `json:"artifact,omitempty"`
	ReviewedBaseline *catalogUpdateBaselineEvidence `json:"reviewed_baseline,omitempty"`
	OpenAPI          *catalogUpdateOpenAPIEvidence  `json:"openapi_validation,omitempty"`
	Inventory        *catalogUpdateInventory        `json:"inventory,omitempty"`
	Diagnostic       string                         `json:"diagnostic,omitempty"`
}

func (c *Client) checkCatalogUpdateSource(ctx context.Context, store *catalogupdate.Store, runID string,
	provider catalog.Provider, source catalogUpdateSource, config CatalogUpdateConfig,
) (catalogupdate.Outcome, error) {
	out := catalogupdate.Outcome{ProviderID: source.ProviderID, ReferenceID: source.ReferenceID}
	evidence := catalogUpdateEvidence{Source: source}
	var promotion, expected *catalogupdate.Version
	emit := func(status, diagnostic string) (catalogupdate.Outcome, error) {
		out.Status = status
		evidence.Diagnostic = diagnostic
		var err error
		out.Evidence, err = json.Marshal(evidence)
		if err == nil {
			if promotion != nil {
				err = store.AcceptDirectVersionAndRecord(ctx, *promotion, expected, runID, out)
			} else {
				err = store.Record(ctx, runID, out)
			}
		}
		return out, err
	}
	if source.Kind == catalog.SpecKindDropboxStone {
		return emit("publisher-source-review", "Check the official Dropbox OpenAPI before substituting Stone provenance.")
	}
	kind := OperationSourceKind(catalog.RefreshableSpecReference{Kind: source.Kind}.ProtocolClassification().UWSSourceType())
	held, err := store.PreferredVersion(ctx, source.ProviderID, source.ReferenceID, source.Service)
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
	previous := held
	catalogPin := catalogUpdateReviewedDigest(source.HeldRevision)
	authoredCommit, _ := catalogUpdateGitFile(source.URL)
	changeBaseline := held.Object.SHA256 == "" || held.CatalogPin != catalogPin || held.Locator != source.URL
	resourceMoved := held.Object.SHA256 != "" && held.Locator != source.URL &&
		!sameCatalogUpdateGitFile(held.Locator, source.URL) && !sameCatalogUpdateGitFile(held.URL, source.URL)
	canReproduceMovedPin := !resourceMoved || (authoredCommit != "" && catalogPin != held.CatalogPin)
	setPromotion := func(version *catalogupdate.Version) {
		promotion = version
		if previous.Object.SHA256 != "" {
			expected = &previous
		}
	}
	// A failed head check does not invalidate separately observed reviewed
	// bytes. Keep the head status for retry and attach the accepted baseline.
	recordHeadFailure := func(status, diagnostic, target string) (catalogupdate.Outcome, error) {
		if ctx.Err() != nil {
			return out, ctx.Err()
		}
		if changeBaseline && catalogPin != "" && canReproduceMovedPin && validOperationSourceKind(kind) {
			anchored, err := c.catalogUpdateReviewedBaseline(ctx, store, source, held, target, catalogPin)
			if err != nil && !errors.Is(err, catalogupdate.ErrObjectIntegrity) {
				return out, err
			}
			if anchored.Preferred.Object.SHA256 != "" {
				evidence.ReviewedBaseline = &anchored
				out.Object = &anchored.Preferred.Object
				setPromotion(&anchored.Preferred)
			}
		}
		return emit(status, diagnostic)
	}
	resolved, err := c.resolveCatalogUpdateSource(ctx, source)
	if held.Service != "" {
		resolved.Service = held.Service
	}
	evidence.Source = resolved
	if err != nil || resolved.Resolution == "unexamined-locator" {
		return recordHeadFailure("publisher-source-unexamined", "Could not establish the publisher artifact locator.", "")
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
		return recordHeadFailure("publisher-check-failed", "Publisher did not return a matching downloadable artifact.", target)
	}
	if artifact.Status != catalogUpdateDownloaded {
		return recordHeadFailure(artifact.Status, "", target)
	}
	if resolved.ResolvedURL != "" && artifact.URL != resolved.ResolvedURL {
		return recordHeadFailure("publisher-source-unexamined", "The artifact redirected away from the resolved immutable locator.", target)
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
			return recordHeadFailure("publisher-spec-invalid", "Local validation established neither a valid OpenAPI document nor evidence for the exact reviewed artifact.", target)
		}
		resolved.Version = firstNonEmpty(validation.Version, resolved.Version)
	} else if _, err := BuildOperationShapeTable(ctx, OperationShapeOptions{
		Sources: []ShapeSourceInput{{
			ID: "catalog-update", OperationSourceInput: OperationSourceInput{Kind: kind, URL: artifact.URL, Content: body.Bytes()},
		}},
	}); err != nil {
		return recordHeadFailure("publisher-spec-unexamined", "Local source validation did not establish a native operation inventory.", target)
	}
	object := catalogupdate.Object{SHA256: artifact.SHA256, Bytes: artifact.Bytes}
	resolved.Revision = firstNonEmpty(resolved.Revision, "sha256:"+artifact.SHA256)
	if kind == OperationSourceGoogleDiscovery {
		native, err := decodeShapeDocument(ctx, body.Bytes())
		if err != nil {
			return recordHeadFailure("publisher-spec-unexamined", "Native source provenance could not be read within the validation budget.", target)
		}
		resolved.Version = firstNonEmpty(stringValue(native["version"]), resolved.Version)
		resolved.Revision = firstNonEmpty(stringValue(native["revision"]), resolved.Revision)
	}
	evidence.Source = resolved
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
		switch {
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
	if held.Object == object && held.URL == artifact.URL && held.Revision == resolved.Revision && held.Locator == source.URL {
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
	if _, err := c.versionURLSyntax(source.URL); err != nil {
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
		return catalogUpdateBaselineEvidence{}, err
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
		return catalogUpdateBaselineEvidence{}, ctx.Err()
	}
	if err != nil || artifact.Status != catalogUpdateDownloaded || artifact.URL != source.URL ||
		!hasCatalogUpdateReviewedDigest(source.HeldRevision, artifact.SHA256) {
		return catalogUpdateBaselineEvidence{}, nil
	}
	object := catalogupdate.Object{SHA256: artifact.SHA256, Bytes: artifact.Bytes}
	baseline := catalogupdate.Version{
		ProviderID: source.ProviderID, ReferenceID: source.ReferenceID, Service: source.Service,
		Locator: source.URL, URL: source.URL, Version: firstNonEmpty(source.Version, source.HeldVersion),
		Revision: commit, CatalogPin: artifact.SHA256, CheckedAt: artifact.CheckedAt, Object: object,
	}
	evidence, valid := qualifyCatalogUpdateBaseline(ctx, source, baseline, body.Bytes())
	if ctx.Err() != nil {
		return catalogUpdateBaselineEvidence{}, ctx.Err()
	}
	if !valid {
		return catalogUpdateBaselineEvidence{}, nil
	}
	if err := store.PutObject(ctx, object, bytes.NewReader(body.Bytes())); err != nil {
		return catalogUpdateBaselineEvidence{}, err
	}
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
	_, err := BuildOperationShapeTable(ctx, OperationShapeOptions{
		Sources: []ShapeSourceInput{{
			ID:                   "catalog-update",
			OperationSourceInput: OperationSourceInput{Kind: kind, URL: version.URL, Content: body},
		}},
	})
	if err != nil {
		return evidence, false
	}
	if kind == OperationSourceGoogleDiscovery {
		native, err := decodeShapeDocument(ctx, body)
		if err != nil {
			return evidence, false
		}
		evidence.Preferred.Version = firstNonEmpty(stringValue(native["version"]), version.Version)
	}
	return evidence, true
}

func hasCatalogUpdateSecurityOverlay(provider, reference string) bool {
	for _, overlay := range catalog.BuiltInSecurityOverlays() {
		if overlay.ProviderID == provider && (overlay.SpecRefID == "" || overlay.SpecRefID == reference) {
			return true
		}
	}
	return false
}
