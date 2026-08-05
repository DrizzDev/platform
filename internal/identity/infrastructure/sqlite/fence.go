package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/marking"
)

// Fence records a one-time attempt against the active revision. A stale revision
// or epoch is Contention; a replayed revision is Fenced.
func (store Store) Fence(scope context.Context, mark marking.Marking) error {
	return store.transact(scope, task{name: "fence", work: func(scope context.Context, transaction *sql.Tx) error {
		handle := mark.Session().String()
		var revision, current uint64
		failure := transaction.QueryRowContext(scope,
			"SELECT revision, epoch FROM pointer WHERE session = ?", handle).Scan(&revision, &current)
		switch {
		case errors.Is(failure, sql.ErrNoRows):
			return Contention{}
		case failure != nil:
			return failure
		case revision != mark.Attempt().Revision() || current != uint64(mark.Attempt().Epoch()):
			return Contention{}
		}
		outcome, failure := transaction.ExecContext(scope,
			"INSERT INTO attempt (session, revision, epoch) VALUES (?, ?, ?) ON CONFLICT DO NOTHING",
			handle, mark.Attempt().Revision(), uint64(mark.Attempt().Epoch()))
		if failure != nil {
			return failure
		}
		count, failure := outcome.RowsAffected()
		if failure != nil {
			return failure
		}
		if count == 0 {
			return Fenced{}
		}
		return nil
	}})
}
