package vault

import "context"

type Store interface {
	Read(scope context.Context, key string) ([]byte, error)
	Write(scope context.Context, entry Entry) error
	Delete(scope context.Context, key string) error
}

type Entry struct {
	Key   string
	Value []byte
}

// Missing reports that a key is absent from the store, distinct from a store
// failure, so the caller can require a fresh sign-in.
type Missing struct{}

func (Missing) Error() string {
	return "credential is absent"
}
