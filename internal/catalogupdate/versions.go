package catalogupdate

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

const versionSchemaSQL = `
CREATE TABLE update_version (
  id INTEGER PRIMARY KEY,
  provider_id TEXT NOT NULL,
  reference_id TEXT NOT NULL,
  service TEXT NOT NULL,
  digest TEXT NOT NULL REFERENCES update_object(digest),
  locator TEXT NOT NULL,
  url TEXT NOT NULL,
  version TEXT NOT NULL,
  revision TEXT NOT NULL,
  checked_at TEXT NOT NULL,
  catalog_pin TEXT NOT NULL DEFAULT '',
  UNIQUE (provider_id, reference_id, service, digest, locator, url, version, revision)
) STRICT;
CREATE TABLE update_preferred (
  provider_id TEXT NOT NULL,
  reference_id TEXT NOT NULL,
  service TEXT NOT NULL,
  version_id INTEGER NOT NULL REFERENCES update_version(id),
  PRIMARY KEY (provider_id, reference_id, service)
) STRICT;
PRAGMA user_version = 3;`

// Version records a reviewed local preference and its exact cached bytes.
type Version struct {
	ProviderID  string    `json:"provider_id"`
	ReferenceID string    `json:"reference_id"`
	Service     string    `json:"service"`
	Locator     string    `json:"locator"`
	URL         string    `json:"url"`
	Version     string    `json:"version,omitempty"`
	Revision    string    `json:"revision,omitempty"`
	CatalogPin  string    `json:"catalog_pin,omitempty"`
	CheckedAt   time.Time `json:"checked_at"`
	Object      Object    `json:"object"`
}

var ErrPreferenceChanged = errors.New("preferred catalog version changed during comparison")
var ErrAmbiguousReferencePreference = errors.New("catalog reference has multiple preferred services")

// AcceptVersion advances a preference after the caller qualifies the version.
// expected binds the update to the version that the caller reviewed. Nil
// permits only the first acceptance of that reference.
func (s *Store) AcceptVersion(ctx context.Context, version Version, expected *Version) error {
	return s.acceptVersion(ctx, version, expected, "", Outcome{})
}

// AcceptVersionAndRecord commits a preference and its run checkpoint in the
// same SQLite transaction. Resume sees both or neither.
func (s *Store) AcceptVersionAndRecord(ctx context.Context, version Version, expected *Version, runID string, outcome Outcome) error {
	if runID == "" || outcome.ProviderID != version.ProviderID || outcome.ReferenceID != version.ReferenceID ||
		outcome.Object == nil || *outcome.Object != version.Object || outcome.Status == "" ||
		(len(outcome.Evidence) != 0 && !json.Valid(outcome.Evidence)) {
		return errors.New("catalog acceptance does not match its run checkpoint")
	}
	return s.acceptVersion(ctx, version, expected, runID, outcome)
}

func (s *Store) acceptVersion(ctx context.Context, version Version, expected *Version, runID string, outcome Outcome) error {
	if version.ProviderID == "" || version.ReferenceID == "" || version.Service == "" || version.Locator == "" ||
		version.URL == "" || version.CheckedAt.IsZero() || (version.CatalogPin != "" && !validDigest(version.CatalogPin)) {
		return errors.New("catalog version is missing identity, provenance or check date")
	}
	object, err := s.OpenObject(ctx, version.Object)
	if err != nil {
		return err
	}
	if err := object.Close(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	current, err := scanVersion(tx.QueryRowContext(ctx, selectVersion+" WHERE p.provider_id=? AND p.reference_id=? AND p.service=?", version.ProviderID, version.ReferenceID, version.Service))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if (expected == nil && err == nil) || (expected != nil && (err != nil || !samePreferredVersion(current, *expected))) {
		return ErrPreferenceChanged
	}
	_, err = tx.ExecContext(ctx, `
INSERT INTO update_version(provider_id, reference_id, service, digest, locator, url, version, revision, checked_at, catalog_pin)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(provider_id, reference_id, service, digest, locator, url, version, revision)
DO UPDATE SET checked_at = excluded.checked_at, catalog_pin = excluded.catalog_pin`, version.ProviderID, version.ReferenceID, version.Service,
		version.Object.SHA256, version.Locator, version.URL, version.Version, version.Revision,
		version.CheckedAt.UTC().Format(time.RFC3339Nano), version.CatalogPin)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
INSERT INTO update_preferred(provider_id, reference_id, service, version_id)
SELECT provider_id, reference_id, service, id FROM update_version
WHERE provider_id = ? AND reference_id = ? AND service = ? AND digest = ? AND locator = ? AND url = ? AND version = ? AND revision = ?
ON CONFLICT(provider_id, reference_id, service) DO UPDATE SET version_id = excluded.version_id`,
		version.ProviderID, version.ReferenceID, version.Service, version.Object.SHA256, version.Locator, version.URL, version.Version, version.Revision)
	if err != nil {
		return err
	}
	if runID != "" {
		if err := recordUpdateOutcome(ctx, tx, runID, outcome, version.Object.SHA256); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func samePreferredVersion(a, b Version) bool {
	a.CheckedAt = a.CheckedAt.UTC()
	b.CheckedAt = b.CheckedAt.UTC()
	return a == b
}

const selectVersion = `
SELECT v.provider_id, v.reference_id, v.service, v.locator, v.url, v.version, v.revision,
 v.checked_at, v.catalog_pin, o.digest, o.bytes
FROM update_preferred p JOIN update_version v ON v.id = p.version_id
JOIN update_object o ON o.digest = v.digest`

func (s *Store) PreferredVersion(ctx context.Context, provider, reference, service string) (Version, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	row := s.db.QueryRowContext(ctx, selectVersion+" WHERE p.provider_id=? AND p.reference_id=? AND p.service=?", provider, reference, service)
	return scanVersion(row)
}

// ReviewedReferencePreference recovers the stable identity of a direct
// source after a reviewed locator change. An index may share a reference
// across several services. Those cases require an explicit service match.
func (s *Store) ReviewedReferencePreference(ctx context.Context, provider, reference string) (Version, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rows, err := s.db.QueryContext(ctx, selectVersion+" WHERE p.provider_id=? AND p.reference_id=? ORDER BY p.service LIMIT 2", provider, reference)
	if err != nil {
		return Version{}, err
	}
	defer rows.Close()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return Version{}, err
		}
		return Version{}, sql.ErrNoRows
	}
	version, err := scanVersion(rows)
	if err != nil {
		return Version{}, err
	}
	if rows.Next() {
		return Version{}, ErrAmbiguousReferencePreference
	}
	return version, rows.Err()
}

func scanVersion(row interface{ Scan(...any) error }) (Version, error) {
	var v Version
	var checked string
	if err := row.Scan(&v.ProviderID, &v.ReferenceID, &v.Service, &v.Locator, &v.URL, &v.Version,
		&v.Revision, &checked, &v.CatalogPin, &v.Object.SHA256, &v.Object.Bytes); err != nil {
		return Version{}, err
	}
	when, err := time.Parse(time.RFC3339Nano, checked)
	v.CheckedAt = when
	return v, err
}

// ExportVersions produces a deterministic metadata transfer for maintainer
// review. Raw publisher objects remain in the local artifact cache.
func (s *Store) ExportVersions(ctx context.Context) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return exportVersions(ctx, s.db, selectVersion)
}

// ReadVersionExport reads an existing maintenance database without creating
// state, checking an artifact volume, or running schema migrations.
func ReadVersionExport(ctx context.Context, statePath string) ([]byte, error) {
	path, err := filepath.Abs(statePath)
	if err != nil {
		return nil, err
	}
	dsn := url.URL{Scheme: "file", Path: filepath.ToSlash(path), RawQuery: "mode=ro&_pragma=busy_timeout(5000)"}
	db, err := sql.Open("sqlite", dsn.String())
	if err != nil {
		return nil, err
	}
	defer db.Close()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var version int
	if err := tx.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil {
		return nil, err
	}
	query := selectVersion
	if version == 2 {
		query = strings.Replace(query, "v.catalog_pin", "''", 1)
	} else if version != stateSchema {
		return nil, errors.New("catalog maintenance database version is unsupported for export")
	}
	return exportVersions(ctx, tx, query)
}

func exportVersions(ctx context.Context, db interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}, query string) ([]byte, error) {
	rows, err := db.QueryContext(ctx, query+" ORDER BY p.provider_id,p.reference_id,p.service")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	doc := struct {
		Schema   string    `json:"schema"`
		Versions []Version `json:"versions"`
	}{Schema: "apitools.catalog-update-export.v1", Versions: []Version{}}
	for rows.Next() {
		version, err := scanVersion(rows)
		if err != nil {
			return nil, err
		}
		doc.Versions = append(doc.Versions, version)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	result, err := json.MarshalIndent(doc, "", "  ")
	return append(result, '\n'), err
}
