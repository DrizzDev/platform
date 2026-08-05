package sqlite

import (
	"context"
	"database/sql"
	_ "embed"
)

//go:embed migration/0001_identity.sql
var schema string

// migrate applies the embedded schema once, re-checking the version inside the
// transaction so a concurrent first migration cannot apply the schema twice.
func (store Store) migrate(scope context.Context) error {
	return store.transact(scope, task{name: "migrate", work: func(scope context.Context, transaction *sql.Tx) error {
		var version int
		if failure := transaction.QueryRowContext(scope, "PRAGMA user_version").Scan(&version); failure != nil {
			return failure
		}
		if version >= 1 {
			return nil
		}
		if _, failure := transaction.ExecContext(scope, schema); failure != nil {
			return failure
		}
		_, failure := transaction.ExecContext(scope, "PRAGMA user_version = 1")
		return failure
	}})
}
