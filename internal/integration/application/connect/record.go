package connect

import (
	"context"
	"log/slog"

	"github.com/DrizzDev/platform/internal/capture/application/recording"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/fidelity"
	"github.com/DrizzDev/platform/internal/capture/domain/origin"
	"github.com/DrizzDev/platform/internal/integration/domain/agent"
)

// mark is one installer action to record: which agent it touched and what was done to it.
type mark struct {
	kind   agent.Kind
	action string
}

// inscribe writes one durable record of an installer action: an authoritative Drizz event, tool category, exact
// fidelity, naming the agent and the action. Recording is observational — a record that cannot be opened is logged and
// swallowed, never fatal to the change it records — but a dropped record is never silent.
func (installer Installer) inscribe(scope context.Context, note mark) {
	execution, failure := installer.recorder.Begin()
	if failure != nil {
		installer.logger.WarnContext(scope, "integration.record.dropped", slog.String("agent", note.kind.String()))
		return
	}
	execution.Record(scope, recording.Note{
		Origin:   origin.Capability,
		Fidelity: fidelity.Exact,
		Category: category.Tool,
		Payload:  []byte(note.action + " " + note.kind.String()),
	})
}
