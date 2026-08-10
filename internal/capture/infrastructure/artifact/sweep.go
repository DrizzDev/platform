package artifact

import (
	"context"
	"os"
	"path/filepath"
)

// Sweep removes temporary files orphaned by a crashed write. Published objects live elsewhere and are never touched.
func (store Store) Sweep(scope context.Context) error {
	return store.observe(scope, probe{name: "sweep", work: func(context.Context) error {
		entries, failure := os.ReadDir(store.temp)
		if failure != nil {
			return failure
		}
		for _, entry := range entries {
			if failure := os.Remove(filepath.Join(store.temp, entry.Name())); failure != nil {
				return failure
			}
		}
		return nil
	}})
}
