package wiring

import (
	"context"

	"github.com/DrizzDev/platform/internal/integration/application/connect"
	"github.com/DrizzDev/platform/internal/integration/domain/agent"
)

// Capture registers Drizz for the agent's turn events, writing the hook block into the agent's hook configuration file
// — which may differ from where its MCP servers live — while preserving every other setting. The command each hook
// runs is carried on the task's server entry, so the registered hook launches exactly this installed Drizz.
func (store Store) Capture(scope context.Context, job connect.Task) (failure error) {
	scope, gauge := store.begin(scope, "capture")
	defer func() { gauge.close(scope, failure) }()

	hooking := job.Agent.Hooking()
	if !hooking.Supported() {
		return connect.Unsupported{}
	}
	document, failure := store.load(hooking)
	if failure != nil {
		return failure
	}
	writer, failure := store.stylist(hooking.Style())
	if failure != nil {
		return failure
	}
	if failure := writer.inscribe(editing{document: document, mark: brand{agent: job.Agent, command: job.Server.Command()}}); failure != nil {
		return failure
	}

	coder, failure := store.codec(hooking.Dialect())
	if failure != nil {
		return failure
	}
	raw, failure := coder.render(document)
	if failure != nil {
		return connect.Locked{}
	}
	return store.stamp(change{path: store.hookpath(hooking), raw: raw})
}

// Uncapture removes Drizz's hook registration from the agent, leaving every other setting untouched.
func (store Store) Uncapture(scope context.Context, target agent.Agent) (failure error) {
	scope, gauge := store.begin(scope, "uncapture")
	defer func() { gauge.close(scope, failure) }()

	hooking := target.Hooking()
	if !hooking.Supported() {
		return connect.Unsupported{}
	}
	document, failure := store.load(hooking)
	if failure != nil {
		return failure
	}
	writer, failure := store.stylist(hooking.Style())
	if failure != nil {
		return failure
	}
	writer.erase(document)

	coder, failure := store.codec(hooking.Dialect())
	if failure != nil {
		return failure
	}
	raw, failure := coder.render(document)
	if failure != nil {
		return connect.Locked{}
	}
	return store.stamp(change{path: store.hookpath(hooking), raw: raw})
}

// Captures reports whether Drizz is already registered for the agent's turn events.
func (store Store) Captures(target agent.Agent) (bool, error) {
	hooking := target.Hooking()
	if !hooking.Supported() {
		return false, nil
	}
	document, failure := store.load(hooking)
	if failure != nil {
		return false, failure
	}
	writer, failure := store.stylist(hooking.Style())
	if failure != nil {
		return false, failure
	}
	return writer.present(document), nil
}
