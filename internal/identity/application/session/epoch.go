package session

import (
	"context"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
)

// Epoch reads the current publication epoch. The flow fences and publishes
// against this one value so the attempt and the compare-and-swap agree.
type Epoch interface {
	Read(context.Context) (epoch.Epoch, error)
}
