package recording

import (
	"context"
	"io"
	"time"

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

// keeper leases the execution's trace so a reclaim pass leaves a running execution's records alone; the sqlite store
// satisfies it.
type keeper interface {
	Lease(context.Context, journal.Claim) error
}

// clock supplies the current time so the lease deadline stays deterministic and the application reads no wall clock.
type clock interface {
	Now() time.Time
}
