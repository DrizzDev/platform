package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/hold"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/lease"
)

// Acquire grants the session lease to the requesting owner unless a live lease is
// held by a competitor, in which case it reports Contention.
func (store Store) Acquire(scope context.Context, request hold.Hold) (lease.Lease, error) {
	var granted lease.Lease
	failure := store.transact(scope, task{name: "acquire", work: func(scope context.Context, transaction *sql.Tx) error {
		handle := request.Session().String()
		var owner string
		var expiry int64
		failure := transaction.QueryRowContext(scope,
			"SELECT owner, expiry FROM lease WHERE session = ?", handle).Scan(&owner, &expiry)
		switch {
		case errors.Is(failure, sql.ErrNoRows):
		case failure != nil:
			return failure
		case expiry > request.Moment().UnixNano() && owner != request.Owner():
			return Contention{}
		}
		deadline := request.Deadline()
		if _, failure := transaction.ExecContext(scope,
			"INSERT INTO lease (session, owner, expiry) VALUES (?, ?, ?) "+
				"ON CONFLICT(session) DO UPDATE SET owner = excluded.owner, expiry = excluded.expiry",
			handle, request.Owner(), deadline.UnixNano()); failure != nil {
			return failure
		}
		held, failure := lease.New(lease.Input{Owner: request.Owner(), Expiry: deadline})
		if failure != nil {
			return failure
		}
		granted = held
		return nil
	}})
	if failure != nil {
		return lease.Lease{}, failure
	}
	return granted, nil
}

// Renew extends the lease only for its current owner; a non-owner is Contention.
func (store Store) Renew(scope context.Context, request hold.Hold) error {
	return store.transact(scope, task{name: "renew", work: func(scope context.Context, transaction *sql.Tx) error {
		outcome, failure := transaction.ExecContext(scope,
			"UPDATE lease SET expiry = ? WHERE session = ? AND owner = ?",
			request.Deadline().UnixNano(), request.Session().String(), request.Owner())
		if failure != nil {
			return failure
		}
		count, failure := outcome.RowsAffected()
		if failure != nil {
			return failure
		}
		if count == 0 {
			return Contention{}
		}
		return nil
	}})
}

// Release drops the lease held by the owner.
func (store Store) Release(scope context.Context, request hold.Hold) error {
	return store.transact(scope, task{name: "release", work: func(scope context.Context, transaction *sql.Tx) error {
		_, failure := transaction.ExecContext(scope,
			"DELETE FROM lease WHERE session = ? AND owner = ?",
			request.Session().String(), request.Owner())
		return failure
	}})
}
