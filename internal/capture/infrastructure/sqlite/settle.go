package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
)

// Settled returns every synchronized entry as a reclaim candidate, oldest first, so a reclaim pass can age out records
// whose retention has elapsed. Un-synchronized entries are never returned — they are never reclaimed.
func (store Store) Settled(scope context.Context) ([]journal.Retained, error) {
	var retained []journal.Retained

	failure := store.observe(scope, probe{name: "settled", work: func(scope context.Context) error {
		rows, failure := store.handle.QueryContext(scope,
			"SELECT "+columns+", stamped FROM journal WHERE state = ? ORDER BY id", journal.Synced.String())
		if failure != nil {
			return failure
		}
		defer func() { _ = rows.Close() }()
		gathered, failure := store.aged(rows)
		if failure != nil {
			return failure
		}
		retained = gathered
		return nil
	}})
	return retained, failure
}

func (Store) aged(rows *sql.Rows) ([]journal.Retained, error) {
	var retained []journal.Retained
	for rows.Next() {
		var line row
		var stamped int64
		if failure := rows.Scan(append(line.targets(), &stamped)...); failure != nil {
			return nil, failure
		}
		entry, failure := line.entry()
		if failure != nil {
			return nil, failure
		}
		link := entry.Correlation()
		retained = append(retained, journal.Retained{
			Trace:    link.Trace(),
			Span:     link.Span(),
			Category: entry.Category(),
			Recorded: time.Unix(stamped, 0),
		})
	}
	return retained, rows.Err()
}
