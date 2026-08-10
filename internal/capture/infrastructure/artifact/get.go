package artifact

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"

	"github.com/DrizzDev/platform/internal/capture/domain/digest"
)

// Get reads an artifact and re-verifies its digest,
// so silent corruption or tampering fails rather than returning bad bytes.
func (store Store) Get(scope context.Context, key digest.Digest) ([]byte, error) {
	var payload []byte

	failure := store.observe(scope, probe{name: "get", work: func(context.Context) error {
		content, failure := os.ReadFile(store.locate(key))
		if errors.Is(failure, os.ErrNotExist) {
			return Absent{}
		}
		if failure != nil {
			return failure
		}
		sum := sha256.Sum256(content)
		if hex.EncodeToString(sum[:]) != key.String() {
			return Integrity{}
		}
		payload = content
		return nil
	}})
	return payload, failure
}
