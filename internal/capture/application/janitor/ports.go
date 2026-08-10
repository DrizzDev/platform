package janitor

import (
	"context"
	"time"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
	"github.com/DrizzDev/platform/internal/capture/domain/digest"
	"github.com/DrizzDev/platform/internal/capture/domain/span"
)

// archive is the journal port: it lists synchronized reclaim candidates, deletes selected entries, reports which
// artifacts remain referenced, and reports active leases. The sqlite store satisfies it.
type archive interface {
	Settled(context.Context) ([]journal.Retained, error)
	Discard(context.Context, []span.Span) error
	Digests(context.Context) ([]digest.Digest, error)
	Leases(context.Context) ([]journal.Claim, error)
}

// vault is the artifact port: it prunes every object a predicate rejects, weighs its footprint, and raises or lowers
// the write gate. The artifact store satisfies it. Passing a predicate keeps the grace and reference policy here in
// the application while the walk stays in the store.
type vault interface {
	Prune(context.Context, func(digest.Digest, time.Time) bool) error
	Footprint(context.Context) (int64, error)
	Restrict(context.Context) error
	Admit(context.Context) error
}

// notifier delivers a disk-pressure alert to the user; the composition root implements it for the running surface.
type notifier interface {
	Alert(context.Context, Pressure) error
}

// clock supplies the current time so retention and lease expiry stay deterministic under test.
type clock interface {
	Now() time.Time
}
