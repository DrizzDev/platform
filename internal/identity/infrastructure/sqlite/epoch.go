package sqlite

import (
	"context"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
)

func (store Store) Epoch(scope context.Context) (epoch.Epoch, error) {
	var value uint64
	failure := store.observe(scope, probe{name: "epoch", work: func(scope context.Context) error {
		return store.handle.QueryRowContext(scope, "SELECT value FROM epoch WHERE id = 1").Scan(&value)
	}})
	if failure != nil {
		return 0, failure
	}
	return epoch.Epoch(value), nil
}
