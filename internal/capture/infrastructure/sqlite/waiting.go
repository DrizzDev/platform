package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
)

// Admit stores one host observation in the pending window. The window is never synchronized;
// it exists only to match a later capability call.
func (store Store) Admit(scope context.Context, observation journal.Observation) error {
	return store.transact(scope, task{name: "admit", work: func(scope context.Context, transaction *sql.Tx) error {
		_, failure := transaction.ExecContext(scope,
			"INSERT INTO pending ("+observations.columns()+") VALUES ("+observations.placeholders()+")",
			observations.values(observation)...)
		return failure
	}})
}

// Waiting returns the live pending observations, newest first and bounded, so a capability call can be matched against
// a bounded candidate set.
func (store Store) Waiting(scope context.Context, window journal.Window) ([]journal.Held, error) {
	var held []journal.Held

	failure := store.observe(scope, probe{name: "waiting", work: func(scope context.Context) error {
		rows, failure := store.handle.QueryContext(scope,
			"SELECT id, "+observations.columns()+" FROM pending WHERE moment >= ? ORDER BY id DESC LIMIT ?",
			window.Cutoff.Unix(), window.Capacity)
		if failure != nil {
			return failure
		}
		defer func() { _ = rows.Close() }()
		gathered, failure := store.observations(rows)
		if failure != nil {
			return failure
		}
		held = gathered
		return nil
	}})
	return held, failure
}

func (Store) observations(rows *sql.Rows) ([]journal.Held, error) {
	var held []journal.Held

	for rows.Next() {
		var line slot
		if failure := rows.Scan(append([]any{&line.reference}, observations.targets(&line)...)...); failure != nil {
			return nil, failure
		}
		observation, failure := line.held()
		if failure != nil {
			return nil, failure
		}
		held = append(held, observation)
	}
	return held, rows.Err()
}

// Evict removes the named observations once they are claimed by a call, in one short transaction.
func (store Store) Evict(scope context.Context, refs []int64) error {
	if len(refs) == 0 {
		return nil
	}
	return store.transact(scope, task{name: "evict", work: func(scope context.Context, transaction *sql.Tx) error {
		for _, ref := range refs {
			if _, failure := transaction.ExecContext(scope, "DELETE FROM pending WHERE id = ?", ref); failure != nil {
				return failure
			}
		}
		return nil
	}})
}

// Expire drops observations older than the cutoff, so unmatched host data for unrelated work never lingers or is kept.
func (store Store) Expire(scope context.Context, cutoff time.Time) error {
	return store.transact(scope, task{name: "expire", work: func(scope context.Context, transaction *sql.Tx) error {
		_, failure := transaction.ExecContext(scope, "DELETE FROM pending WHERE moment < ?", cutoff.Unix())
		return failure
	}})
}
