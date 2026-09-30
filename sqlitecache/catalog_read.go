package sqlitecache

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/OpenUdon/apitools/catalog"
)

// MaxCatalogRegistrationBytes bounds aggregate registry text before decoding.
const MaxCatalogRegistrationBytes = 32 << 20

// ReadCatalogSpecArtifacts adapts the same read-only snapshot for the root
// library's index builder, avoiding a root-package import of sqlitecache.
func ReadCatalogSpecArtifacts(ctx context.Context, options catalog.RootOptions) ([]catalog.CatalogSpecArtifact, error) {
	rows, err := ReadCatalogArtifacts(ctx, options)
	if err != nil {
		return nil, err
	}
	result := make([]catalog.CatalogSpecArtifact, 0, len(rows))
	for _, row := range rows {
		ref := row.Metadata["spec_ref_id"]
		if ref == "" {
			ref = row.ArtifactID
		}
		result = append(result, catalog.CatalogSpecArtifact{ProviderID: row.ProviderID, SpecRefID: ref, ArtifactID: row.ArtifactID, Kind: row.Kind, Path: row.Path, SourceURL: row.SourceURL, OverlayPath: row.OverlayPath, BuilderPath: row.BuilderPath, SHA256: row.SHA256, Bytes: row.Bytes, Metadata: row.Metadata})
	}
	return result, nil
}

// ReadCatalogArtifacts reads an existing catalog root's registrations without
// migration, pruning or access-time writes. A missing registry remains an
// os.ErrNotExist error so consumers can report insufficient evidence.
// No root returns no registrations; it never discovers an implicit directory.
func ReadCatalogArtifacts(ctx context.Context, options catalog.RootOptions) ([]CatalogArtifact, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	paths, err := catalog.ResolveRoot(options)
	if err != nil {
		return nil, err
	}
	if paths.Directory == "" {
		return nil, nil
	}
	info, err := os.Lstat(paths.RegistryPath)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("catalog registry is not a regular file")
	}
	uri := url.URL{Scheme: "file", Path: filepath.ToSlash(paths.RegistryPath), RawQuery: "mode=ro"}
	db, err := sql.Open("sqlite", uri.String())
	if err != nil {
		return nil, err
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var version, versionRows int
	if err := tx.QueryRowContext(ctx, "SELECT COUNT(*), MIN(version) FROM schema_meta").Scan(&versionRows, &version); err != nil {
		return nil, fmt.Errorf("read catalog registry schema: %w", err)
	}
	if versionRows != 1 || version < 1 || version > schemaVersion {
		return nil, fmt.Errorf("unsupported catalog registry schema %d", version)
	}
	var count, bytes int64
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*), COALESCE(SUM(
 COALESCE(length(CAST(provider_id AS BLOB)), 0) + COALESCE(length(CAST(artifact_id AS BLOB)), 0) +
 COALESCE(length(CAST(kind AS BLOB)), 0) + COALESCE(length(CAST(path AS BLOB)), 0) +
 COALESCE(length(CAST(source_url AS BLOB)), 0) + COALESCE(length(CAST(overlay_path AS BLOB)), 0) +
 COALESCE(length(CAST(builder_path AS BLOB)), 0) + COALESCE(length(CAST(sha256 AS BLOB)), 0) +
 COALESCE(length(CAST(metadata_json AS BLOB)), 0)), 0) FROM catalog_artifacts`).Scan(&count, &bytes); err != nil {
		return nil, err
	}
	if count > DefaultMaxCatalogArtifacts || bytes > MaxCatalogRegistrationBytes {
		return nil, fmt.Errorf("catalog registrations exceed snapshot bounds (%d rows, %d bytes)", count, bytes)
	}
	// The catalog table's read columns are compatible across schemas 1–3.
	// Snapshot the rows in the same transaction as the schema check.
	rows, err := tx.QueryContext(ctx, `SELECT provider_id, artifact_id, kind, path,
 source_url, overlay_path, builder_path, sha256, bytes, metadata_json, updated_at
 FROM catalog_artifacts ORDER BY provider_id, artifact_id LIMIT ?`, DefaultMaxCatalogArtifacts+1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return readCatalogArtifactRows(rows, normalizeOptions(Options{}))
}
