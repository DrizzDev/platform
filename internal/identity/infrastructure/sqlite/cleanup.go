package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/cleanup"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/deferral"
)

const (
	reconcile = 4
	retries   = 5
)

// Enqueue records a superseded or rejected credential for later removal.
func (store Store) Enqueue(scope context.Context, record cleanup.Record) error {
	return store.transact(scope, task{name: "enqueue", work: func(scope context.Context, transaction *sql.Tx) error {
		_, failure := transaction.ExecContext(scope,
			"INSERT INTO cleanup (name, reason, state, attempts, next, deadline, created) "+
				"VALUES (?, ?, ?, ?, ?, ?, ?) ON CONFLICT(name) DO NOTHING",
			record.Key(), string(record.Reason()), string(record.State()), record.Attempts(),
			record.Next().UnixNano(), record.Deadline().UnixNano(), record.Created().UnixNano())
		return failure
	}})
}

// Pending returns the due, unblocked cleanup records up to the reconcile batch.
func (store Store) Pending(scope context.Context, moment time.Time) ([]cleanup.Record, error) {
	var records []cleanup.Record
	failure := store.observe(scope, probe{name: "pending", work: func(scope context.Context) error {
		rows, failure := store.handle.QueryContext(scope,
			"SELECT name, reason, state, attempts, next, deadline, created FROM cleanup "+
				"WHERE state = ? AND next <= ? ORDER BY created LIMIT ?",
			string(cleanup.Pending), moment.UnixNano(), reconcile)
		if failure != nil {
			return failure
		}
		defer func() { _ = rows.Close() }()
		for rows.Next() {
			var attempts uint
			var name, reason, state string
			var next, deadline, created int64
			if failure := rows.Scan(&name, &reason, &state, &attempts, &next, &deadline, &created); failure != nil {
				return failure
			}
			record, failure := cleanup.New(cleanup.Input{
				Key:      name,
				Attempts: attempts,
				Next:     time.Unix(0, next),
				State:    cleanup.State(state),
				Created:  time.Unix(0, created),
				Reason:   cleanup.Reason(reason),
				Deadline: time.Unix(0, deadline),
			})
			if failure != nil {
				return failure
			}
			records = append(records, record)
		}
		return rows.Err()
	}})
	if failure != nil {
		return nil, failure
	}
	return records, nil
}

// Acknowledge removes a reconciled record from the queue.
func (store Store) Acknowledge(scope context.Context, key string) error {
	return store.transact(scope, task{name: "acknowledge", work: func(scope context.Context, transaction *sql.Tx) error {
		_, failure := transaction.ExecContext(scope, "DELETE FROM cleanup WHERE name = ?", key)
		return failure
	}})
}

// Defer reschedules a record after a failed attempt, blocking it once retries run out.
func (store Store) Defer(scope context.Context, deferral deferral.Deferral) error {
	return store.transact(scope, task{name: "defer", work: func(scope context.Context, transaction *sql.Tx) error {
		var attempts uint
		failure := transaction.QueryRowContext(scope,
			"SELECT attempts FROM cleanup WHERE name = ?", deferral.Key()).Scan(&attempts)
		switch {
		case errors.Is(failure, sql.ErrNoRows):
			return nil
		case failure != nil:
			return failure
		}
		attempts++
		if attempts >= retries {
			_, failure := transaction.ExecContext(scope,
				"UPDATE cleanup SET attempts = ?, state = ? WHERE name = ?",
				attempts, string(cleanup.Blocked), deferral.Key())
			return failure
		}
		_, failure = transaction.ExecContext(scope,
			"UPDATE cleanup SET attempts = ?, next = ? WHERE name = ?",
			attempts, deferral.Next().UnixNano(), deferral.Key())
		return failure
	}})
}

// Backlog reports the total number of queued records.
func (store Store) Backlog(scope context.Context) (int, error) {
	var count int
	failure := store.observe(scope, probe{name: "backlog", work: func(scope context.Context) error {
		return store.handle.QueryRowContext(scope, "SELECT count(*) FROM cleanup").Scan(&count)
	}})
	if failure != nil {
		return 0, failure
	}
	return count, nil
}
