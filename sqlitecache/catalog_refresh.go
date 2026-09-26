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

func insertRefreshSpec(ctx context.Context, tx *sql.Tx, registration catalogRefreshRegistration) error {
	now := time.Now().UTC().Unix()
	accessed := time.Now().UTC().UnixNano()
	for _, urlValue := range uniqueStrings(registration.originalURL, registration.finalURL) {
		if _, err := tx.ExecContext(ctx, `
INSERT INTO spec_documents (
  url, original_url, final_url, sha256, bytes, metadata_json, content_path, content, first_seen_at, updated_at, accessed_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(url) DO UPDATE SET
  original_url = excluded.original_url,
  final_url = excluded.final_url,
  sha256 = excluded.sha256,
  bytes = excluded.bytes,
  metadata_json = excluded.metadata_json,
  content_path = excluded.content_path,
  content = excluded.content,
	updated_at = excluded.updated_at,
	accessed_at = excluded.accessed_at`,
			urlValue, registration.originalURL, registration.finalURL, registration.sha256, registration.bytes,
			registration.specMetadataJSON, nullableString(registration.contentPath), nullableContent(registration.contentPath, nil), now, now, accessed); err != nil {
			return err
		}
	}
	return nil
}

func insertRefreshCatalogArtifact(ctx context.Context, tx *sql.Tx, registration catalogRefreshRegistration) error {
	now := time.Now().UTC().Unix()
	accessed := time.Now().UTC().UnixNano()
	_, err := tx.ExecContext(ctx, `
INSERT INTO catalog_artifacts (
  provider_id, artifact_id, kind, path, source_url, overlay_path, builder_path,
  sha256, bytes, metadata_json, first_seen_at, updated_at, accessed_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(provider_id, artifact_id) DO UPDATE SET
  kind = excluded.kind,
  path = excluded.path,
  source_url = excluded.source_url,
  overlay_path = excluded.overlay_path,
  builder_path = excluded.builder_path,
  sha256 = excluded.sha256,
  bytes = excluded.bytes,
  metadata_json = excluded.metadata_json,
	updated_at = excluded.updated_at,
	accessed_at = excluded.accessed_at`,
		registration.providerID,
		registration.artifactID,
		registration.kind,
		registration.contentPath,
		registration.sourceURL,
		"",
		"",
		registration.sha256,
		registration.bytes,
		registration.catalogArtifactMetadataJSON,
		now,
		now,
		accessed,
	)
	return err
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
