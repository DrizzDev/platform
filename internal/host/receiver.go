package host

import (
	"context"

	"github.com/DrizzDev/platform/internal/integration/application/intake"
	"github.com/DrizzDev/platform/internal/integration/domain/agent"
	"github.com/DrizzDev/platform/internal/integration/infrastructure/inbound"
	"github.com/DrizzDev/platform/internal/integration/infrastructure/telemetry"
	"github.com/DrizzDev/platform/internal/transport/cli"
)

var _ cli.Receiver = receiver{}

// receiver is the host runtime behind the hidden hook endpoint. It records one inbound agent event as a host
// observation. It must never disrupt the agent that fired the hook, so an unknown agent, an unprovisionable runtime, or
// an unrecordable event all resolve quietly rather than as an error the agent would see.
type receiver struct {
	base foundation
}

func (receiver receiver) Receive(scope context.Context, signal cli.Signal) error {
	target, found := agent.New().Lookup(agent.Kind(signal.Agent))
	if !found {
		return nil
	}

	_, observer, failure := receiver.base.provision(scope)
	if failure != nil {
		return nil
	}
	defer session{observer: observer}.shutdown(scope)

	recorder, store, failure := receiver.base.recorder(scope, observer)
	if failure != nil {
		return nil
	}
	defer func() { _ = store.Close() }()

	monitor, failure := telemetry.New(telemetry.Options{Tracer: observer.Tracer(), Meter: observer.Meter()})
	if failure != nil {
		return nil
	}
	made, failure := intake.New(intake.Options{
		Recorder: recorder,
		Monitor:  monitor,
		Logger:   observer.Diagnostics(),
	})
	if failure != nil {
		return nil
	}
	event := inbound.New(receiver.base.streams.Input).Read(inbound.Request{
		Agent:   target,
		Slot:    agent.Slot(signal.Slot),
		Payload: signal.Payload,
	})
	made.Record(scope, event)
	return nil
}
