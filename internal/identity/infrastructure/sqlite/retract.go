package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/cleanup"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/retraction"
)

// Retract clears the active pointer, advances the epoch, and enqueues the active
// key for cleanup in one short transaction. It reports Absent when nothing is
// active, so logout stays idempotent.
func (store Store) Retract(scope context.Context, request retraction.Retraction) error {
	return store.transact(scope, task{name: "retract", work: func(scope context.Context, transaction *sql.Tx) error {
		handle := request.Session().String()
		var name string
		failure := transaction.QueryRowContext(scope,
			"SELECT name FROM pointer WHERE session = ?", handle).Scan(&name)
		switch {
		case errors.Is(failure, sql.ErrNoRows):
			return Absent{}
		case failure != nil:
			return failure
		}
		if _, failure := transaction.ExecContext(scope, "DELETE FROM pointer WHERE session = ?", handle); failure != nil {
			return failure
		}
		if _, failure := transaction.ExecContext(scope, "UPDATE epoch SET value = value + 1 WHERE id = 1"); failure != nil {
			return failure
		}
		moment := request.Moment().UnixNano()
		_, failure = transaction.ExecContext(scope, orphan,
			name, string(cleanup.Logout), string(cleanup.Pending), moment, moment, moment)
		return failure
	}})
}
