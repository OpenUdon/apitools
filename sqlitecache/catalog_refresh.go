package sqlitecache

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/OpenUdon/apitools"
)

// RegisterCatalogRefreshResults records successfully saved refresh artifacts
// in the bounded SQLite cache. Artifact paths must remain under the cache file's
// directory so later reads retain the cache's path-containment guarantees.
func RegisterCatalogRefreshResults(ctx context.Context, cachePath, cacheDir string, report apitools.CatalogSpecRefreshReport) error {
	cachePath = strings.TrimSpace(cachePath)
	if cachePath == "" {
		return fmt.Errorf("cache path is required")
	}
	cache, err := Open(cachePath)
	if err != nil {
		return err
	}
	defer cache.Close()
	prepared := make([]catalogRefreshRegistration, 0, len(report.Results))
	for _, result := range report.Results {
		registration, err := prepareCatalogRefreshRegistration(cache, cacheDir, result)
		if err != nil {
			return fmt.Errorf("%s/%s: %w", result.ProviderID, result.SpecRefID, err)
		}
		prepared = append(prepared, registration)
	}
	if len(prepared) == 0 {
		return nil
	}
	tx, err := cache.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()
	for _, registration := range prepared {
		if err := insertRefreshSpec(ctx, tx, registration); err != nil {
			return err
		}
		if err := insertRefreshCatalogArtifact(ctx, tx, registration); err != nil {
			return err
		}
	}
	if _, err := cache.pruneTx(ctx, tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	committed = true
	return nil
}

type catalogRefreshRegistration struct {
	providerID                  string
	artifactID                  string
	kind                        string
	contentPath                 string
	sourceURL                   string
	originalURL                 string
	finalURL                    string
	sha256                      string
	bytes                       int64
	specMetadataJSON            string
	catalogArtifactMetadataJSON string
}

func prepareCatalogRefreshRegistration(cache *Cache, cacheDir string, result apitools.CatalogSpecRefreshResult) (catalogRefreshRegistration, error) {
	contentPath, err := cache.catalogRefreshContentPath(cacheDir, result)
	if err != nil {
		return catalogRefreshRegistration{}, err
	}
	file, err := cache.readArtifact(contentPath, result.SHA256, result.Bytes)
	if err != nil {
		return catalogRefreshRegistration{}, err
	}
	if file.Bytes <= 0 {
		return catalogRefreshRegistration{}, fmt.Errorf("catalog refresh artifact is empty")
	}
	originalURL := strings.TrimSpace(result.URL)
	finalURL := strings.TrimSpace(result.FinalURL)
	if originalURL == "" {
		originalURL = finalURL
	}
	if finalURL == "" {
		finalURL = originalURL
	}
	if originalURL == "" {
		return catalogRefreshRegistration{}, fmt.Errorf("spec URL is required")
	}
	providerID := strings.TrimSpace(result.ProviderID)
	artifactID := strings.TrimSpace(result.SpecRefID)
	kind := strings.TrimSpace(string(result.Kind))
	if providerID == "" || artifactID == "" || kind == "" {
		return catalogRefreshRegistration{}, fmt.Errorf("provider, spec reference, and artifact kind are required")
	}
	specMetadataJSON, err := json.Marshal(result.RawMetadata)
	if err != nil {
		return catalogRefreshRegistration{}, err
	}
	if int64(len(specMetadataJSON)) > cache.options.MaxMetadataBytes {
		return catalogRefreshRegistration{}, fmt.Errorf("spec metadata is %d bytes, over limit %d", len(specMetadataJSON), cache.options.MaxMetadataBytes)
	}
	artifactMetadata := map[string]string{
		"official":                    "true",
		"kind":                        kind,
		"raw_validation_status":       result.RawValidationStatus,
		"corrected_validation_status": result.CorrectedValidationStatus,
		"sha256":                      file.SHA256,
	}
	artifactMetadataJSON, err := json.Marshal(artifactMetadata)
	if err != nil {
		return catalogRefreshRegistration{}, err
	}
	if int64(len(artifactMetadataJSON)) > cache.options.MaxMetadataBytes {
		return catalogRefreshRegistration{}, fmt.Errorf("catalog artifact metadata is %d bytes, over limit %d", len(artifactMetadataJSON), cache.options.MaxMetadataBytes)
	}
	if err := validateRecordBudget("cached spec", cache.options.MaxMetadataBytes, originalURL, finalURL, file.SHA256, contentPath); err != nil {
		return catalogRefreshRegistration{}, err
	}
	if err := validateRecordBudget("catalog artifact", cache.options.MaxMetadataBytes, providerID, artifactID, kind, contentPath, originalURL, file.SHA256); err != nil {
		return catalogRefreshRegistration{}, err
	}
	return catalogRefreshRegistration{
		providerID:                  providerID,
		artifactID:                  artifactID,
		kind:                        kind,
		contentPath:                 contentPath,
		sourceURL:                   originalURL,
		originalURL:                 originalURL,
		finalURL:                    finalURL,
		sha256:                      file.SHA256,
		bytes:                       file.Bytes,
		specMetadataJSON:            string(specMetadataJSON),
		catalogArtifactMetadataJSON: string(artifactMetadataJSON),
	}, nil
}

// insertRefreshSpec and insertRefreshCatalogArtifact share their INSERT SQL
// (storeSpecDocumentTx, storeCatalogArtifactTx in cache.go) with StoreSpec and
// StoreCatalogArtifact, so a future migration or column change to either
// cannot silently diverge between the two registration paths.
func insertRefreshSpec(ctx context.Context, tx *sql.Tx, registration catalogRefreshRegistration) error {
	now := time.Now().UTC().Unix()
	accessed := time.Now().UTC().UnixNano()
	for _, urlValue := range uniqueStrings(registration.originalURL, registration.finalURL) {
		if err := storeSpecDocumentTx(ctx, tx, urlValue, registration.originalURL, registration.finalURL, registration.sha256, registration.bytes, registration.specMetadataJSON, registration.contentPath, nil, now, accessed); err != nil {
			return err
		}
	}
	return nil
}

func insertRefreshCatalogArtifact(ctx context.Context, tx *sql.Tx, registration catalogRefreshRegistration) error {
	now := time.Now().UTC().Unix()
	accessed := time.Now().UTC().UnixNano()
	return storeCatalogArtifactTx(ctx, tx, registration.providerID, registration.artifactID, registration.kind, registration.contentPath, registration.sourceURL, "", "", registration.sha256, registration.bytes, registration.catalogArtifactMetadataJSON, now, accessed)
}

func (c *Cache) catalogRefreshContentPath(cacheDir string, result apitools.CatalogSpecRefreshResult) (string, error) {
	if c == nil || c.db == nil {
		return "", fmt.Errorf("cache is closed")
	}
	if c.baseDir == "" {
		return "", fmt.Errorf("catalog refresh artifacts require a file-backed cache")
	}
	fullPath := strings.TrimSpace(result.SavedPath)
	if fullPath == "" {
		fullPath = filepath.Join(cacheDir, filepath.FromSlash(result.ArtifactPath))
	}
	fullPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(c.baseDir, fullPath)
	if err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel) {
		return filepath.ToSlash(rel), nil
	}
	return "", fmt.Errorf("refreshed artifact %q must stay under SQLite cache directory %q", fullPath, c.baseDir)
}
