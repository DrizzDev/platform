package sqlite

import (
	"context"
	"database/sql"

	"github.com/DrizzDev/platform/internal/capture/domain/digest"
)

// Digests returns the distinct set of artifact references still held by any entry, so a reclaim pass can delete only
// the stored objects that no remaining entry points to.
func (store Store) Digests(scope context.Context) ([]digest.Digest, error) {
	var references []digest.Digest

	failure := store.observe(scope, probe{name: "digests", work: func(scope context.Context) error {
		rows, failure := store.handle.QueryContext(scope,
			"SELECT DISTINCT artifact FROM journal WHERE artifact != ''")
		if failure != nil {
			return failure
		}
		defer func() { _ = rows.Close() }()
		gathered, failure := store.references(rows)
		if failure != nil {
			return failure
		}
		references = gathered
		return nil
	}})
	return references, failure
}

func (Store) references(rows *sql.Rows) ([]digest.Digest, error) {
	var references []digest.Digest
	for rows.Next() {
		var value string
		if failure := rows.Scan(&value); failure != nil {
			return nil, failure
		}
		reference, failure := digest.New(value)
		if failure != nil {
			return nil, Corrupt{}
		}
		references = append(references, reference)
	}
	return references, rows.Err()
}
