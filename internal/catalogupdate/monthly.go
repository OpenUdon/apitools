package catalogupdate

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

const monthlyRunTable = `CREATE TABLE IF NOT EXISTS update_monthly (
    period TEXT PRIMARY KEY,
    run_id TEXT NOT NULL UNIQUE REFERENCES update_run(id)
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
	var id, assigned string
	err = tx.QueryRowContext(ctx, `SELECT m.run_id, m.period FROM update_monthly m JOIN update_run r ON r.id=m.run_id
        WHERE r.completed=0 OR m.period=? ORDER BY r.completed, m.period LIMIT 1`, period).Scan(&id, &assigned)
	if err == nil {
		run, err := readRun(ctx, tx, id, stateSchema)
		run.MonthlyPeriod = assigned
		return run, err
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
