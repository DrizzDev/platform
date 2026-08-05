package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/pointer"
	"github.com/DrizzDev/platform/internal/identity/domain/session"
)

// Head returns the current published pointer for a session, or Absent when a
// session has never published a credential.
func (store Store) Head(scope context.Context, handle session.Session) (pointer.Pointer, error) {
	var name string
	var revision, mark uint64
	failure := store.observe(scope, probe{name: "head", work: func(scope context.Context) error {
		outcome := store.handle.QueryRowContext(scope,
			"SELECT name, revision, epoch FROM pointer WHERE session = ?", handle.String()).Scan(&name, &revision, &mark)
		if errors.Is(outcome, sql.ErrNoRows) {
			return Absent{}
		}
		return outcome
	}})
	if failure != nil {
		return pointer.Pointer{}, failure
	}
	return pointer.New(pointer.Input{Key: name, Revision: revision, Epoch: epoch.Epoch(mark)})
}
