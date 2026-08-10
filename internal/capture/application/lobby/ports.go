package lobby

import (
	"context"
	"time"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
)

// register is the bounded pending window: it admits host observations, returns the live ones to match against a call,
// evicts the claimed, and expires the stale. The sqlite store satisfies it.
type register interface {
	Admit(context.Context, journal.Observation) error
	Waiting(context.Context, journal.Window) ([]journal.Held, error)
	Evict(context.Context, []int64) error
	Expire(context.Context, time.Time) error
}

// clock supplies the current time so the window's expiry and match horizon stay deterministic under test.
type clock interface {
	Now() time.Time
}
