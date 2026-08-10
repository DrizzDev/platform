package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
	"github.com/DrizzDev/platform/internal/capture/domain/trace"
)

// Lease records that a trace is active until the claim's instant, protecting its records from reclaim. A later claim
// for the same trace extends it, so a running execution can keep its protection fresh.
func (store Store) Lease(scope context.Context, claim journal.Claim) error {
	return store.transact(scope, task{name: "lease", work: func(scope context.Context, transaction *sql.Tx) error {
		_, failure := transaction.ExecContext(scope,
			"INSERT INTO lease (trace, until) VALUES (?, ?) ON CONFLICT (trace) DO UPDATE SET until = excluded.until",
			claim.Trace.String(), claim.Until.Unix())
		return failure
	}})
}

// Leases returns every active claim, so a reclaim pass can skip the traces of executions still in flight.
func (store Store) Leases(scope context.Context) ([]journal.Claim, error) {
	var claims []journal.Claim

	failure := store.observe(scope, probe{name: "leases", work: func(scope context.Context) error {
		rows, failure := store.handle.QueryContext(scope, "SELECT trace, until FROM lease")
		if failure != nil {
			return failure
		}
		defer func() { _ = rows.Close() }()
		gathered, failure := store.claims(rows)
		if failure != nil {
			return failure
		}
		claims = gathered
		return nil
	}})
	return claims, failure
}

func (Store) claims(rows *sql.Rows) ([]journal.Claim, error) {
	var claims []journal.Claim

	for rows.Next() {
		var until int64
		var thread string

		if failure := rows.Scan(&thread, &until); failure != nil {
			return nil, failure
		}
		subject, failure := trace.New(thread)
		if failure != nil {
			return nil, Corrupt{}
		}
		claims = append(claims, journal.Claim{Trace: subject, Until: time.Unix(until, 0)})
	}
	return claims, rows.Err()
}
