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

// slip is a dropped record: which step failed and why. The recorder logs it and swallows it, so the reason a record
// was dropped is observable — a missing payload, for instance, is not mistaken for a storage fault — without ever
// surfacing to the capability being recorded.
type slip struct {
	reason error
	step   string
}

func (recorder Recorder) drop(scope context.Context, slip slip) {
	recorder.logger.WarnContext(scope, "capture.record.dropped",
		slog.String("step", slip.step), slog.Any("error", slip.reason))
}
