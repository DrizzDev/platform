package host

import (
	"context"
	"time"

	"github.com/DrizzDev/platform/internal/platform/observability"
)

const shutdown = 10 * time.Second

type session struct {
	observer observability.Provider
}

func (session session) shutdown(scope context.Context) {
	stopping, cancel := context.WithTimeout(context.WithoutCancel(scope), shutdown)
	defer cancel()
	// Close records a code on failure. Shutdown is best-effort, so the returned
	// provider error is intentionally discarded rather than logged here.
	_ = session.observer.Close(stopping)
}
