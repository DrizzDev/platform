package vault

import (
	"context"
	"errors"

	system "github.com/zalando/go-keyring"
)

const namespace = "drizz.platform.identity"

var _ Store = Keyring{}

type Keyring struct{}

type reading struct {
	value   string
	failure error
}

func (keyring Keyring) Read(scope context.Context, key string) ([]byte, error) {
	done := make(chan reading, 1)
	go func() {
		value, failure := system.Get(namespace, key)
		done <- reading{value: value, failure: failure}
	}()
	select {
	case <-scope.Done():
		return nil, scope.Err()
	case outcome := <-done:
		if errors.Is(outcome.failure, system.ErrNotFound) {
			return nil, Missing{}
		}
		if outcome.failure != nil {
			return nil, outcome.failure
		}
		return []byte(outcome.value), nil
	}
}

func (keyring Keyring) Write(scope context.Context, entry Entry) error {
	return keyring.dispatch(scope, func() error {
		return system.Set(namespace, entry.Key, string(entry.Value))
	})
}

func (keyring Keyring) Delete(scope context.Context, key string) error {
	failure := keyring.dispatch(scope, func() error {
		return system.Delete(namespace, key)
	})
	if errors.Is(failure, system.ErrNotFound) {
		return nil
	}
	return failure
}

// dispatch bounds a blocking keychain call by the caller's context. The buffered
// channel lets the abandoned goroutine finish and exit even after a timeout.
func (keyring Keyring) dispatch(scope context.Context, work func() error) error {
	done := make(chan error, 1)
	go func() { done <- work() }()
	select {
	case <-scope.Done():
		return scope.Err()
	case failure := <-done:
		return failure
	}
}
