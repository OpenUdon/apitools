package catalogupdate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"
)

const monthlyRunTable = `CREATE TABLE IF NOT EXISTS update_monthly (
    period TEXT PRIMARY KEY,
    run_id TEXT NOT NULL UNIQUE REFERENCES update_run(id)
) STRICT`

const supersededMonthlyRunTable = `CREATE TABLE IF NOT EXISTS update_monthly_superseded (
    run_id TEXT PRIMARY KEY REFERENCES update_run(id),
    period TEXT NOT NULL,
    superseded_at TEXT NOT NULL
) STRICT`

// MonthlyRun reuses an interrupted scheduled run, then a completed run for
// the current period. Create a new run and its period assignment atomically.
// The caller holds the monthly run lock until the run finishes or stops.
func (s *Store) MonthlyRun(ctx context.Context, period, identity, catalogSHA256 string, providers []string) (Run, error) {
	if _, err := time.Parse("2006-01", period); err != nil {
		return Run{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return Run{}, err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, monthlyRunTable); err != nil {
		return Run{}, err
	}
	var id string
	err = tx.QueryRowContext(ctx, `SELECT m.run_id FROM update_monthly m JOIN update_run r ON r.id=m.run_id
        WHERE (r.completed=0 AND m.period<=?) OR m.period=? ORDER BY r.completed, m.period LIMIT 1`, period, period).Scan(&id)
	if err == nil {
		return readRun(ctx, tx, id, stateSchema)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return Run{}, err
	}
	run, err := insertRunSnapshot(ctx, tx, identity, catalogSHA256, providers)
	if err != nil {
		return Run{}, err
	}
	if _, err := tx.ExecContext(ctx, "INSERT INTO update_monthly VALUES (?, ?)", period, run.ID); err != nil {
		return Run{}, err
	}
	if err := tx.Commit(); err != nil {
		return Run{}, err
	}
	run.MonthlyPeriod = period
	return run, nil
}

// SupersedeMonthlyRun releases an unfinished period assignment under the
// writer lock. The run, its checkpoints and accepted versions remain in place.
func SupersedeMonthlyRun(ctx context.Context, statePath, cachePath, expectedUUID, id string) error {
	if _, err := os.Lstat(statePath); err != nil {
		return err
	}
	cache, err := OpenArtifactCache(cachePath, expectedUUID, 0)
	if err != nil {
		return err
	}
	defer cache.Close()
	lock, err := LockMonthlyRun(ctx, statePath)
	if err != nil {
		return err
	}
	defer lock.Close()
	store, err := OpenStore(ctx, statePath, cache)
	if err != nil {
		return err
	}
	defer store.Close()
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	run, err := readRun(ctx, tx, id, stateSchema)
	if err != nil {
		return err
	}
	if !run.MonthlySupersededAt.IsZero() {
		return nil
	}
	if run.MonthlyPeriod == "" || run.Completed {
		return fmt.Errorf("run %s is not an unfinished scheduled run", id)
	}
	if _, err := tx.ExecContext(ctx, supersededMonthlyRunTable); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "INSERT INTO update_monthly_superseded VALUES (?, ?, ?)", id, run.MonthlyPeriod,
		time.Now().UTC().Format(time.RFC3339Nano)); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM update_monthly WHERE run_id=?", id); err != nil {
		return err
	}
	return tx.Commit()
}

func readMonthlyAssignment(ctx context.Context, db catalogRunQueries, run *Run) error {
	var active, superseded bool
	err := db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM sqlite_schema WHERE type='table' AND name='update_monthly'),
        EXISTS(SELECT 1 FROM sqlite_schema WHERE type='table' AND name='update_monthly_superseded')`).Scan(&active, &superseded)
	if err != nil {
		return err
	}
	if active {
		err := db.QueryRowContext(ctx, "SELECT period FROM update_monthly WHERE run_id=?", run.ID).Scan(&run.MonthlyPeriod)
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	if !superseded {
		return nil
	}
	var at string
	err = db.QueryRowContext(ctx, "SELECT period, superseded_at FROM update_monthly_superseded WHERE run_id=?", run.ID).Scan(&run.MonthlyPeriod, &at)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err == nil {
		run.MonthlySupersededAt, err = time.Parse(time.RFC3339Nano, at)
	}
	return err
}

// ReadLatestMonthlyRunID does not open the artifact cache or migrate state.
// A database predating scheduled runs has no monthly run.
func ReadLatestMonthlyRunID(ctx context.Context, statePath string) (_ string, err error) {
	defer func() { err = classifyRunReadError(err) }()
	if _, err := os.Stat(statePath); errors.Is(err, fs.ErrNotExist) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	db, err := openReadOnlyState(statePath)
	if err != nil {
		return "", err
	}
	defer db.Close()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	if _, err := checkCatalogStateSchema(ctx, tx); err != nil {
		return "", err
	}
	var exists bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM sqlite_schema
        WHERE type='table' AND name='update_monthly')`).Scan(&exists); err != nil || !exists {
		return "", err
	}
	var id string
	err = tx.QueryRowContext(ctx, "SELECT run_id FROM update_monthly ORDER BY period DESC LIMIT 1").Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return id, err
}
