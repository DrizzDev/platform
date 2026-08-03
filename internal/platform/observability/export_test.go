package observability

import (
	"context"
	"log/slog"
)

type Wiring struct {
	Diagnostics *slog.Logger
	Sink        slog.Handler
	Close       func(context.Context) error
}

func Wire(wiring Wiring) Provider {
	return Provider{
		diagnostics: wiring.Diagnostics,
		sink:        slog.New(wiring.Sink),
		closer:      wiring.Close,
	}
}
