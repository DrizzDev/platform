package recording

import (
	"context"
	"io"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
	"github.com/DrizzDev/platform/internal/capture/domain/digest"
)

// writer is the journal port; the sqlite store satisfies it.
type writer interface {
	Append(context.Context, journal.Entry) error
}

// sink is the artifact port; the artifact store satisfies it.
type sink interface {
	Put(context.Context, io.Reader) (digest.Digest, error)
}
