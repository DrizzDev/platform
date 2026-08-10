package artifact

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/DrizzDev/platform/internal/capture/domain/digest"
)

// Put streams an artifact to a bounded temporary file, then atomically publishes it at its content-addressed path.
// A crash before the rename leaves only a temp, never a half-published object.
func (store Store) Put(scope context.Context, source io.Reader) (digest.Digest, error) {
	var key digest.Digest

	failure := store.observe(scope, probe{name: "put", work: func(context.Context) error {
		temporary, failure := os.CreateTemp(store.temp, "artifact-*")
		if failure != nil {
			return failure
		}
		defer func() { _ = os.Remove(temporary.Name()) }()

		hasher := sha256.New()
		written, failure := io.Copy(io.MultiWriter(temporary, hasher), io.LimitReader(source, store.ceiling+1))
		if closing := temporary.Close(); failure == nil {
			failure = closing
		}
		if failure != nil {
			return failure
		}
		if written > store.ceiling {
			return Oversize{}
		}

		key, failure = digest.New(hex.EncodeToString(hasher.Sum(nil)))
		if failure != nil {
			return failure
		}

		final := store.locate(key)
		if failure := os.MkdirAll(filepath.Dir(final), 0o700); failure != nil {
			return failure
		}
		return os.Rename(temporary.Name(), final)
	}})

	return key, failure
}
