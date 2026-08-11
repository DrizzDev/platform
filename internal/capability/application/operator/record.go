package operator

import (
	"context"
	"log/slog"

	"github.com/DrizzDev/platform/internal/capture/application/recording"
)

// entry is one execution record: the capability that produced it and the note to write.
type entry struct {
	capability string
	note       recording.Note
}

// inscribe opens one execution and writes its note. Recording is observational, so a failure to open the record is
// logged and swallowed — it can never break the capability it records — but a dropped record is never silent.
func (operator Operator) inscribe(scope context.Context, record entry) {
	execution, failure := operator.recorder.Begin()
	if failure != nil {
		operator.logger.WarnContext(scope, "capability.record.dropped", slog.String("capability", record.capability))
		return
	}
	execution.Record(scope, record.note)
}
