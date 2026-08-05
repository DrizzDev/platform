package reconcile

import (
	"context"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/cleanup"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/deferral"
)

// Queue is the cleanup backlog the flow drains: it yields the due records,
// removes the ones deleted, and reschedules the ones that failed.
type Queue interface {
	Pending(context.Context, time.Time) ([]cleanup.Record, error)
	Acknowledge(context.Context, string) error
	Defer(context.Context, deferral.Deferral) error
}
