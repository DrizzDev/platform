package operator

import (
	"context"
	"log/slog"

	"github.com/DrizzDev/platform/internal/capture/application/recording"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/fidelity"
	"github.com/DrizzDev/platform/internal/capture/domain/origin"
)

// entry is one execution record: the capability that produced it and the note to write.
type entry struct {
	capability string
	note       recording.Note
}

// action builds the record note for a performed interaction: a tool-category entry describing what was done, with no
// artifact of its own. The payload describes the action without capturing any typed text.
func (Operator) action(payload string) recording.Note {
	return recording.Note{
		Origin:   origin.Capability,
		Fidelity: fidelity.Exact,
		Category: category.Tool,
		Payload:  []byte(payload),
	}
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
