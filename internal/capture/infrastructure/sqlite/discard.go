package sqlite

import (
	"context"
	"database/sql"

	"github.com/DrizzDev/platform/internal/capture/domain/span"
)

// Discard deletes the named entries in one short transaction, so a reclaim pass selects and removes them atomically in
// this authoritative store before any artifact file is touched.
func (store Store) Discard(scope context.Context, steps []span.Span) error {
	if len(steps) == 0 {
		return nil
	}
	return store.transact(scope, task{name: "discard", work: func(scope context.Context, transaction *sql.Tx) error {
		for _, step := range steps {
			if _, failure := transaction.ExecContext(scope,
				"DELETE FROM journal WHERE span = ?", step.String()); failure != nil {
				return failure
			}
		}
		return nil
	}})
}
