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
		link := entry.Correlation()
		_, failure := transaction.ExecContext(scope,
			"INSERT INTO journal (trace, span, parent, sequence, mark, origin, fidelity, category, payload, artifact, state, stamped) "+
				"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			link.Trace().String(), link.Span().String(), link.Parent().String(),
			link.Sequence(), link.Mark().String(),
			entry.Origin().String(), entry.Fidelity().String(), entry.Category().String(),
			entry.Payload(), entry.Artifact().String(), entry.State().String(), time.Now().Unix())
		return failure
	}})
}
