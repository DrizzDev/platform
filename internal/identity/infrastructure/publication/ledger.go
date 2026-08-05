package publication

import (
	"context"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/marking"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/pointer"
	notice "github.com/DrizzDev/platform/internal/identity/application/coordination/publication"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/result"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/retraction"
	"github.com/DrizzDev/platform/internal/identity/domain/session"
)

// Ledger is the non-secret coordination store the publisher and renewer depend
// on, satisfied by the SQLite adapter.
type Ledger interface {
	Epoch(context.Context) (epoch.Epoch, error)
	Head(context.Context, session.Session) (pointer.Pointer, error)
	Fence(context.Context, marking.Marking) error
	Publish(context.Context, notice.Publication) (result.Result, error)
	Retract(context.Context, retraction.Retraction) error
	Backlog(context.Context) (int, error)
}
