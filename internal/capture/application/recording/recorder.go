package recording

import (
	"context"
	"log/slog"
	"time"
)

// Recorder writes execution records. It is observational: a failed write is logged and swallowed,
// never surfaced to the capability being recorded, so recording can never break the operation it observes.
type Recorder struct {
	sink   sink
	clock  clock
	writer writer
	keeper keeper
	logger *slog.Logger
	lease  time.Duration
}

func (recorder Recorder) drop(scope context.Context, step string) {
	recorder.logger.WarnContext(scope, "capture.record.dropped", slog.String("step", step))
}
