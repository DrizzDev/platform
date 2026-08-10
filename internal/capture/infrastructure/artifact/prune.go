package artifact

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/DrizzDev/platform/internal/capture/domain/digest"
)

// Prune deletes every published object the keep predicate rejects, so a reclaim pass can remove objects no entry
// references while sparing any the predicate protects. The store owns the walk; the caller owns the policy.
// Deletion is idempotent, so a prune interrupted mid-sweep safely resumes on the next pass.
func (store Store) Prune(scope context.Context, keep func(digest.Digest, time.Time) bool) error {
	return store.observe(scope, probe{name: "prune", work: func(context.Context) error {
		found, failure := store.walk()
		if failure != nil {
			return failure
		}
		for _, object := range found {
			if keep(object.digest, object.modified) {
				continue
			}
			if failure := store.remove(object.digest); failure != nil {
				return failure
			}
		}
		return nil
	}})
}

func (store Store) remove(key digest.Digest) error {
	failure := os.Remove(store.locate(key))
	if errors.Is(failure, os.ErrNotExist) {
		return nil
	}
	return failure
}
