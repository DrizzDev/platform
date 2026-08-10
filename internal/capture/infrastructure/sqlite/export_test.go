package sqlite

import "context"

// Query and Exec expose raw SQL for the store's own qualification tests.

func (store Store) Query(statement string) (string, error) {
	var value string
	failure := store.handle.QueryRowContext(context.Background(), statement).Scan(&value)
	return value, failure
}

func (store Store) Exec(statement string) error {
	_, failure := store.handle.ExecContext(context.Background(), statement)
	return failure
}
