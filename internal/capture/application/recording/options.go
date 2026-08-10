package recording

import (
	"log/slog"
	"time"
)

type Options struct {
	Sink   sink
	Clock  clock
	Writer writer
	Keeper keeper
	Logger *slog.Logger
	Lease  time.Duration
}
