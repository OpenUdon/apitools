package catalogupdate

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	_ "modernc.org/sqlite"
)

const stateSchema = 3
const RetryPendingStatus = "publisher-retry-pending"

// Serialize the first creation and close of a state file in this process.
var stateFileCreation sync.Mutex

// Store keeps the durable run and artifact registry on the configured local
// state filesystem. The owner supplies and closes the artifact cache root.
type Store struct {
	db    *sql.Conn
	pool  *sql.DB
	cache *os.Root
	// Serialize access to the pinned connection across transaction lifetimes.
	mu sync.Mutex
}

type Run struct {
	ID        string
	Catalog   string
	Providers []string
	Completed bool
}

type Outcome struct {
	ProviderID  string          `json:"provider_id"`
	ReferenceID string          `json:"reference_id"`
	Status      string          `json:"status"`
	Evidence    json.RawMessage `json:"evidence,omitempty"`
	Object      *Object         `json:"object,omitempty"`
	// Replacing a pending retry is conditional on the exact observed
	// checkpoint. A competing or finalized outcome stays authoritative.
	PreviousRetryEvidence json.RawMessage `json:"-"`
}

type Object struct {
	SHA256 string `json:"sha256"`
	Bytes  int64  `json:"bytes"`
}

func OpenStore(ctx context.Context, statePath string, cache *os.Root) (*Store, error) {
	return openStore(ctx, statePath, cache, checkStateFilesystem)
}

func openStore(ctx context.Context, statePath string, cache *os.Root, check func(*os.File) error) (*Store, error) {
	if strings.TrimSpace(statePath) == "" || cache == nil {
		return nil, errors.New("catalog state path and artifact cache root are required")
	}
	abs, err := filepath.Abs(statePath)
	if err != nil {
		return nil, err
	}
	dir, err := openArtifactCache(filepath.Dir(abs), func(root *os.Root) error {
		dir, err := root.Open(".")
		if err != nil {
			return err
		}
		return errors.Join(check(dir), dir.Close())
	})
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	name := filepath.Base(abs)
	statePath = filepath.Join(dir.Name(), name)
	if err := createStateFile(dir, name); err != nil {
		return nil, err
	}
	info, err := dir.Lstat(name)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, errors.New("catalog state path must be a regular file")
	}
	dbFile, err := openStatePath(dir, name)
	if err != nil {
		return nil, err
	}
	defer dbFile.Close()
	if opened, err := dbFile.Stat(); err != nil || !os.SameFile(info, opened) {
		return nil, errors.New("catalog state file changed during open")
	}
	if err := check(dbFile); err != nil {
		return nil, err
	}
	// A previous process may have stopped between a mkdir and its parent
	// fsync. Flush the pinned filesystems before referring to either path.
	if err := syncArtifactFilesystem(dir); err != nil {
		return nil, err
	}
	if err := syncArtifactFilesystem(cache); err != nil {
		return nil, err
	}
	if err := checkStatePath(statePath, info); err != nil {
		return nil, err
	}
	dsn := url.URL{Scheme: "file", Path: filepath.ToSlash(statePath), RawQuery: "mode=rw&_txlock=immediate&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)"}
	pool, err := sql.Open("sqlite", dsn.String())
	if err != nil {
		return nil, err
	}
	pool.SetMaxOpenConns(1)
	db, err := pool.Conn(ctx)
	if err != nil {
		_ = pool.Close()
		return nil, err
	}
	store := &Store{db: db, pool: pool, cache: cache}
	if err := checkStatePath(statePath, info); err != nil {
		_ = store.Close()
		return nil, err
	}
	if err := store.initialize(ctx); err != nil {
		_ = store.Close()
		return nil, err
	}
	return store, nil
}

func createStateFile(dir *os.Root, name string) error {
	stateFileCreation.Lock()
	defer stateFileCreation.Unlock()
	probe, err := dir.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0o600)
	if errors.Is(err, fs.ErrExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return probe.Close()
}

func checkStatePath(path string, verified fs.FileInfo) error {
	current, err := os.Lstat(path)
	if err != nil || !os.SameFile(verified, current) {
		return errors.New("catalog state path changed before SQLite initialization")
	}
	return nil
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return errors.Join(s.db.Close(), s.pool.Close())
}

func (s *Store) initialize(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var version int
	if err := tx.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil {
		return err
	}
	switch version {
	case 1, 2, stateSchema:
		var ok int
		if err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM sqlite_schema WHERE type='table' AND name IN ('update_run','update_provider','update_outcome','update_object')").Scan(&ok); err != nil {
			return err
		}
		if ok != 4 {
			return errors.New("file is not a catalog maintenance database")
		}
		if version >= 2 {
			if err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM sqlite_schema WHERE type='table' AND name IN ('update_version','update_preferred')").Scan(&ok); err != nil {
				return err
			}
			if ok != 2 {
				return errors.New("catalog maintenance database lacks version tables")
			}
		}
	case 0:
		var tables int
		if err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM sqlite_schema WHERE type='table'").Scan(&tables); err != nil {
			return err
		}
		if tables != 0 {
			return errors.New("file is not an empty or initialized catalog maintenance database")
		}
		if _, err := tx.ExecContext(ctx, schemaSQL); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported catalog maintenance schema %d", version)
	}
	if version < 2 {
		if _, err := tx.ExecContext(ctx, versionSchemaSQL); err != nil {
			return err
		}
	} else if version == 2 {
		if _, err := tx.ExecContext(ctx, `
ALTER TABLE update_version ADD COLUMN catalog_pin TEXT NOT NULL DEFAULT '';
PRAGMA user_version = 3;`); err != nil {
			return err
		}
	}
	return tx.Commit()
}

const schemaSQL = `
CREATE TABLE update_run (
  id TEXT PRIMARY KEY,
  catalog TEXT NOT NULL,
  completed INTEGER NOT NULL DEFAULT 0
) STRICT;
CREATE TABLE update_provider (
  run_id TEXT NOT NULL REFERENCES update_run(id),
  provider_id TEXT NOT NULL,
  position INTEGER NOT NULL,
  PRIMARY KEY (run_id, provider_id),
  UNIQUE (run_id, position)
) STRICT;
CREATE TABLE update_object (
  digest TEXT PRIMARY KEY,
  bytes INTEGER NOT NULL CHECK (bytes > 0)
) STRICT;
CREATE TABLE update_outcome (
  run_id TEXT NOT NULL,
  provider_id TEXT NOT NULL,
  reference_id TEXT NOT NULL,
  status TEXT NOT NULL,
  evidence TEXT,
  digest TEXT REFERENCES update_object(digest),
  PRIMARY KEY (run_id, provider_id, reference_id),
  FOREIGN KEY (run_id, provider_id) REFERENCES update_provider(run_id, provider_id)
) STRICT;
PRAGMA user_version = 1;`

func (s *Store) StartRun(ctx context.Context, catalog string, providers []string) (Run, error) {
	if !validDigest(catalog) || len(providers) == 0 {
		return Run{}, errors.New("catalog digest and provider scope are required")
	}
	run := Run{ID: rand.Text(), Catalog: catalog, Providers: append([]string(nil), providers...)}
	s.mu.Lock()
	defer s.mu.Unlock()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return Run{}, err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, "INSERT INTO update_run(id, catalog) VALUES (?, ?)", run.ID, run.Catalog); err != nil {
		return Run{}, err
	}
	for position, provider := range run.Providers {
		if strings.TrimSpace(provider) == "" {
			return Run{}, errors.New("empty provider ID in catalog update scope")
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO update_provider VALUES (?, ?, ?)", run.ID, provider, position); err != nil {
			return Run{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return Run{}, err
	}
	return run, nil
}

func (s *Store) ReadRun(ctx context.Context, id string) (Run, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run := Run{ID: id}
	if err := s.db.QueryRowContext(ctx, "SELECT catalog, completed FROM update_run WHERE id = ?", id).Scan(&run.Catalog, &run.Completed); err != nil {
		return Run{}, err
	}
	rows, err := s.db.QueryContext(ctx, "SELECT provider_id FROM update_provider WHERE run_id = ? ORDER BY position", id)
	if err != nil {
		return Run{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return Run{}, err
		}
		run.Providers = append(run.Providers, id)
	}
	return run, rows.Err()
}

func (s *Store) Record(ctx context.Context, runID string, outcome Outcome) error {
	if outcome.Status == "" || (len(outcome.Evidence) > 0 && !json.Valid(outcome.Evidence)) {
		return errors.New("catalog update outcome needs a status and valid evidence")
	}
	var digest any
	if outcome.Object != nil {
		file, err := s.OpenObject(ctx, *outcome.Object)
		if err != nil {
			return err
		}
		if err := file.Close(); err != nil {
			return err
		}
		digest = outcome.Object.SHA256
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return recordUpdateOutcome(ctx, s.db, runID, outcome, digest)
}

func recordUpdateOutcome(ctx context.Context, executor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}, runID string, outcome Outcome, digest any) error {
	if len(outcome.PreviousRetryEvidence) != 0 {
		result, err := executor.ExecContext(ctx, `
UPDATE update_outcome SET status = ?, evidence = ?, digest = ?
WHERE run_id = ? AND provider_id = ? AND reference_id = ? AND status = ? AND evidence = ?
  AND EXISTS (SELECT 1 FROM update_run WHERE id = ? AND completed = 0)`,
			outcome.Status, string(outcome.Evidence), digest, runID, outcome.ProviderID, outcome.ReferenceID,
			RetryPendingStatus, string(outcome.PreviousRetryEvidence), runID)
		if err != nil {
			return err
		}
		return requireUpdatedRun(result)
	}
	result, err := executor.ExecContext(ctx, `
INSERT INTO update_outcome(run_id, provider_id, reference_id, status, evidence, digest)
SELECT id, ?, ?, ?, ?, ? FROM update_run WHERE id = ? AND completed = 0
ON CONFLICT(run_id, provider_id, reference_id) DO NOTHING`,
		outcome.ProviderID, outcome.ReferenceID, outcome.Status, string(outcome.Evidence), digest, runID)
	if err != nil {
		return err
	}
	return requireUpdatedRun(result)
}

func (s *Store) ReadOutcomes(ctx context.Context, runID string) ([]Outcome, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rows, err := s.db.QueryContext(ctx, `
SELECT o.provider_id, o.reference_id, o.status, o.evidence, d.digest, d.bytes
FROM update_outcome o JOIN update_provider p USING (run_id, provider_id)
LEFT JOIN update_object d ON o.digest=d.digest
WHERE o.run_id = ? ORDER BY p.position, o.reference_id`, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var outcomes []Outcome
	for rows.Next() {
		var out Outcome
		var evidence string
		var digest sql.NullString
		var size sql.NullInt64
		if err := rows.Scan(&out.ProviderID, &out.ReferenceID, &out.Status, &evidence, &digest, &size); err != nil {
			return nil, err
		}
		if evidence != "" {
			out.Evidence = json.RawMessage(evidence)
		}
		if digest.Valid && size.Valid {
			out.Object = &Object{SHA256: digest.String, Bytes: size.Int64}
		}
		outcomes = append(outcomes, out)
	}
	return outcomes, rows.Err()
}

func (s *Store) FinishRun(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// A provider-level checkpoint uses an empty reference ID.
	result, err := s.db.ExecContext(ctx, `
UPDATE update_run SET completed = 1 WHERE id = ? AND NOT EXISTS (
 SELECT 1 FROM update_provider p WHERE p.run_id = update_run.id AND NOT EXISTS (
  SELECT 1 FROM update_outcome o WHERE o.run_id = p.run_id AND
    o.provider_id = p.provider_id AND o.reference_id = ''
 )
) AND NOT EXISTS (
 SELECT 1 FROM update_outcome WHERE run_id = update_run.id AND status = ?
)`, id, RetryPendingStatus)
	if err != nil {
		return err
	}
	return requireUpdatedRun(result)
}

func requireUpdatedRun(result sql.Result) error {
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("catalog update run cannot be changed in its current state")
	}
	return nil
}

func validDigest(value string) bool {
	raw, err := hex.DecodeString(value)
	return err == nil && len(raw) == 32 && value == strings.ToLower(value)
}
