package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/cleanup"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/publication"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/result"
)

const orphan = "INSERT INTO cleanup (name, reason, state, attempts, next, deadline, created) " +
	"VALUES (?, ?, ?, 0, ?, ?, ?) ON CONFLICT(name) DO NOTHING"

type transfer struct {
	transaction *sql.Tx
	notice      publication.Publication
	current     uint64
}

// Publish advances the fenced head pointer when the caller's expected epoch still
// matches the stored epoch. A stale epoch is Rejected; a failed transaction is
// Uncertain, since durable state cannot be assumed either way.
func (store Store) Publish(scope context.Context, notice publication.Publication) (result.Result, error) {
	resolution := result.Rejected
	failure := store.transact(scope, task{name: "publish", work: func(scope context.Context, transaction *sql.Tx) error {
		var current uint64
		if failure := transaction.QueryRowContext(scope, "SELECT value FROM epoch WHERE id = 1").Scan(&current); failure != nil {
			return failure
		}
		if current != uint64(notice.Expected()) {
			return store.reject(scope, transfer{transaction: transaction, notice: notice, current: current})
		}
		if failure := store.exchange(scope, transfer{transaction: transaction, notice: notice, current: current}); failure != nil {
			return failure
		}
		resolution = result.Published
		return nil
	}})
	if failure != nil {
		return result.Uncertain, failure
	}
	return resolution, nil
}

func (store Store) exchange(scope context.Context, transfer transfer) error {
	notice := transfer.notice
	handle := notice.Session().String()
	var prior string
	failure := transfer.transaction.QueryRowContext(scope, "SELECT name FROM pointer WHERE session = ?", handle).Scan(&prior)
	if failure != nil && !errors.Is(failure, sql.ErrNoRows) {
		return failure
	}
	if _, failure := transfer.transaction.ExecContext(scope, "UPDATE epoch SET value = value + 1 WHERE id = 1"); failure != nil {
		return failure
	}
	if _, failure := transfer.transaction.ExecContext(scope,
		"INSERT INTO pointer (session, name, revision, epoch) VALUES (?, ?, ?, ?) "+
			"ON CONFLICT(session) DO UPDATE SET name = excluded.name, revision = excluded.revision, epoch = excluded.epoch",
		handle, notice.Key(), notice.Revision(), transfer.current+1); failure != nil {
		return failure
	}
	if prior == "" || prior == notice.Key() {
		return nil
	}
	moment := notice.Moment().UnixNano()
	_, failure = transfer.transaction.ExecContext(scope, orphan,
		prior, string(cleanup.Superseded), string(cleanup.Pending), moment, moment, moment)
	return failure
}

// reject queues the losing candidate credential for removal so a stale publish
// does not leave its vault entry orphaned.
func (store Store) reject(scope context.Context, transfer transfer) error {
	moment := transfer.notice.Moment().UnixNano()
	_, failure := transfer.transaction.ExecContext(scope, orphan,
		transfer.notice.Key(), string(cleanup.Rejected), string(cleanup.Pending), moment, moment, moment)
	return failure
}
