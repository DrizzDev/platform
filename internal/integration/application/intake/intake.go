// Package intake records inbound agent hook events. It is the receiving half of the hook feature: an agent notifies
// Drizz that a person submitted a prompt or that a turn finished, and this service writes that as a host-origin
// capture so the context around a device action is preserved. Recording is observational — a hook that cannot be
// recorded never fails the agent that fired it.
package intake

import (
	"context"
	"log/slog"

	"github.com/DrizzDev/platform/internal/capture/application/recording"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/fidelity"
	"github.com/DrizzDev/platform/internal/capture/domain/origin"
	"github.com/DrizzDev/platform/internal/integration/domain/agent"
)

// Sink records one host observation; the capture recorder satisfies it.
type Sink interface {
	Begin() (*recording.Execution, error)
}

// Monitor is intake's port to operational observability. Its signatures are standard-library types, so the application
// core stays free of any telemetry vendor; an infrastructure adapter maps a scope to a trace span and a latency metric.
type Monitor interface {
	Watch(scope context.Context, operation string) (context.Context, func(string))
}

// Intake writes an inbound hook event as a host-origin capture.
type Intake struct {
	recorder Sink
	monitor  Monitor
	logger   *slog.Logger
}

type Options struct {
	Recorder Sink
	Monitor  Monitor
	Logger   *slog.Logger
}

func New(options Options) (Intake, error) {
	if failure := options.validate(); failure != nil {
		return Intake{}, failure
	}
	return Intake{recorder: options.Recorder, monitor: options.Monitor, logger: options.Logger}, nil
}

// Record writes the event as a host observation. An event with no text is still recorded as a turn marker, so the
// ordering of prompts, turns, and device actions is preserved even when an agent exposes no text for a moment. The
// hook path fires on every prompt and turn, so it is measured: a span and a latency metric label the slot and outcome
// only — never the agent identifier or the captured text.
func (intake Intake) Record(scope context.Context, event Event) {
	scope, close := intake.monitor.Watch(scope, "hook")
	outcome := "recorded"
	defer func() { close(outcome) }()

	execution, failure := intake.recorder.Begin()
	if failure != nil {
		outcome = "dropped"
		intake.logger.WarnContext(scope, "integration.hook.dropped",
			slog.String("slot", event.Slot.String()))
		return
	}
	execution.Record(scope, recording.Note{
		Origin:   origin.Host,
		Fidelity: fidelity.Inferred,
		Category: intake.classify(event.Slot),
		Payload:  []byte(event.Agent.String() + " " + event.Slot.String()),
		Artifact: []byte(event.Text),
	})
}

// classify maps a turn moment to the sensitive-data class of the text it carries: a prompt is a person's words, a
// completed turn is the agent's response.
func (Intake) classify(slot agent.Slot) category.Category {
	if slot == agent.Turn {
		return category.Response
	}
	return category.Prompt
}
