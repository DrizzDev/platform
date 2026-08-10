package artifact

import (
	"context"
	"errors"
	"os"
	"path/filepath"
)

const sentinel = ".degraded"

func (store Store) barrier() string {
	return filepath.Join(store.root, sentinel)
}

// Restrict raises the back-pressure gate so new writes are refused while the disk budget is exceeded. It is durable and
// shared, so every process reading this store honours the same gate.
func (store Store) Restrict(scope context.Context) error {
	return store.observe(scope, probe{name: "restrict", work: func(context.Context) error {
		return os.WriteFile(store.barrier(), nil, 0o600)
	}})
}

// Admit lowers the back-pressure gate so new writes resume once reclaim has freed enough space.
func (store Store) Admit(scope context.Context) error {
	return store.observe(scope, probe{name: "admit", work: func(context.Context) error {
		failure := os.Remove(store.barrier())
		if errors.Is(failure, os.ErrNotExist) {
			return nil
		}
		return failure
	}})
}

func (store Store) restricted() bool {
	_, failure := os.Stat(store.barrier())
	return failure == nil
}
