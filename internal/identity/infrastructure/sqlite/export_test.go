package sqlite

import "context"

func (store Store) Query(statement string) (string, error) {
	var value string
	if failure := store.handle.QueryRowContext(context.Background(), statement).Scan(&value); failure != nil {
		return "", failure
	}
	return value, nil
}

func (store Store) Exec(statement string) error {
	_, failure := store.handle.ExecContext(context.Background(), statement)
	return failure
}

func (store Store) Plan(statement string) (string, error) {
	rows, failure := store.handle.QueryContext(context.Background(), "EXPLAIN QUERY PLAN "+statement)
	if failure != nil {
		return "", failure
	}
	defer func() { _ = rows.Close() }()
	var plan string
	for rows.Next() {
		var node, parent, aux int
		var detail string
		if failure := rows.Scan(&node, &parent, &aux, &detail); failure != nil {
			return "", failure
		}
		plan += detail + "\n"
	}
	return plan, rows.Err()
}
