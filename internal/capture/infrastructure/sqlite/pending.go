package sqlite

import (
	"context"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
)

// Pending returns every entry not yet synchronized, in recorded order across all traces, so the courier can drain them.
func (store Store) Pending(scope context.Context) ([]journal.Entry, error) {
	var records []journal.Entry

	failure := store.observe(scope, probe{name: "pending", work: func(scope context.Context) error {
		rows, failure := store.handle.QueryContext(scope,
			"SELECT "+entries.columns()+" FROM journal WHERE state = ? ORDER BY id", journal.Pending.String())
		if failure != nil {
			return failure
		}
		defer func() { _ = rows.Close() }()
		gathered, failure := store.collect(rows)
		if failure != nil {
			return failure
		}
		records = gathered
		return nil
	}})
	return records, failure
}
