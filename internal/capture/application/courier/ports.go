package courier

import (
	"context"
	"io"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
	"github.com/DrizzDev/platform/internal/capture/domain/digest"
)

// ledger yields unsynchronized entries and records their new state; the sqlite store satisfies it.
type ledger interface {
	Pending(context.Context) ([]journal.Entry, error)
	Advance(context.Context, journal.Transition) error
}

// vault reads an artifact's bytes for upload; the artifact store satisfies it.
type vault interface {
	Get(context.Context, digest.Digest) ([]byte, error)
}

// uploader sends an artifact and an entry, each idempotent by identity, so an at-least-once resend yields one effect.
type uploader interface {
	Blob(context.Context, Cargo) error
	Record(context.Context, journal.Entry) error
}

// Cargo is one artifact to upload: its digest — the idempotency key — and a stream of its bytes.
type Cargo struct {
	Source io.Reader
	Digest digest.Digest
}
