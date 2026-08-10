package recording

import (
	"context"
	"log/slog"
)

// Recorder writes execution records. It is observational: a failed write is logged and swallowed,
// never surfaced to the capability being recorded, so recording can never break the operation it observes.
type Recorder struct {
	sink   sink
	writer writer
	logger *slog.Logger
}

func (recorder Recorder) drop(scope context.Context, step string) {
	recorder.logger.WarnContext(scope, "capture.record.dropped", slog.String("step", step))
}
