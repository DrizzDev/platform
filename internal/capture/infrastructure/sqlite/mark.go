package sqlite

import (
	"context"
	"database/sql"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
)

// Advance moves one entry to a new synchronization state. The target state is carried in the transition, so a new
// state applies without changing this method.
func (store Store) Advance(scope context.Context, transition journal.Transition) error {
	return store.transact(scope, task{name: "advance", work: func(scope context.Context, transaction *sql.Tx) error {
		_, failure := transaction.ExecContext(scope, "UPDATE journal SET state = ? WHERE span = ?",
			transition.State.String(), transition.Entry.Correlation().Span().String())
		return failure
	}})
}
