package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
)

// Append writes one entry to the journal as a single immediate transaction,
// so the record and its synchronization state commit atomically.
func (store Store) Append(scope context.Context, entry journal.Entry) error {
	return store.transact(scope, task{name: "append", work: func(scope context.Context, transaction *sql.Tx) error {
		values := append(store.values(entry), time.Now().Unix())
		_, failure := transaction.ExecContext(scope,
			"INSERT INTO journal ("+columns+", stamped) VALUES ("+placeholders+", ?)", values...)
		return failure
	}})
}
