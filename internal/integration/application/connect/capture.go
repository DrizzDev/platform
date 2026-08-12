package connect

import (
	"context"
	"log/slog"

	"github.com/DrizzDev/platform/internal/integration/domain/agent"
	"github.com/DrizzDev/platform/internal/integration/domain/server"
)

// Capture registers Drizz for the selected agents' turn events, so a person's prompts and the agent's responses are
// recorded as context around device actions. It is a separate, opt-in step from connecting the agent's tools, because
// it captures a person's own words. Like Enable, it resolves the executable once and reports each agent individually.
func (installer Installer) Capture(scope context.Context, selection Selection) (report Report, failure error) {
	scope, watch := installer.begin(scope, "capture")
	defer func() { watch.finish(scope, installer.tally(summary{report: report, failure: failure})) }()

	command, failure := installer.resolver.Locate()
	if failure != nil {
		return Report{}, failure
	}
	entry, failure := server.New(server.Input{Name: Name, Command: command, Args: []string{Launch}})
	if failure != nil {
		return Report{}, failure
	}
	var outcomes []Outcome
	for _, target := range installer.chosen(selection) {
		outcomes = append(outcomes, installer.watch(scope, Task{Agent: target, Server: entry}))
	}
	report = Report{outcomes: outcomes}
	return report, nil
}

// Uncapture removes Drizz's turn-event registration from the selected agents, leaving every other setting untouched.
func (installer Installer) Uncapture(scope context.Context, selection Selection) Report {
	scope, watch := installer.begin(scope, "uncapture")
	report := Report{}
	defer func() { watch.finish(scope, installer.tally(summary{report: report})) }()

	var outcomes []Outcome
	for _, target := range installer.chosen(selection) {
		outcomes = append(outcomes, installer.unwatch(scope, target))
	}
	report = Report{outcomes: outcomes}
	return report
}

func (installer Installer) watch(scope context.Context, job Task) Outcome {
	base := Outcome{kind: job.Agent.Kind(), title: job.Agent.Title(), restart: job.Agent.Restart()}
	present, failure := installer.store.Detect(job.Agent)
	if failure != nil || !present {
		base.state = Missing
		return base
	}
	if !job.Agent.Hooking().Supported() {
		base.state = Incapable
		return base
	}
	if failure := installer.store.Capture(scope, job); failure != nil {
		base.state = installer.grade(failure)
		base.detail = installer.detail(failure)
		return base
	}
	installer.inscribe(scope, mark{kind: job.Agent.Kind(), action: "captured"})
	base.state = Captured
	base.capturing = true
	return base
}

func (installer Installer) unwatch(scope context.Context, target agent.Agent) Outcome {
	base := Outcome{kind: target.Kind(), title: target.Title(), restart: target.Restart()}
	present, failure := installer.store.Detect(target)
	if failure != nil || !present {
		base.state = Missing
		return base
	}
	if !target.Hooking().Supported() {
		base.state = Incapable
		return base
	}
	if failure := installer.store.Uncapture(scope, target); failure != nil {
		base.state = installer.grade(failure)
		base.detail = installer.detail(failure)
		return base
	}
	installer.inscribe(scope, mark{kind: target.Kind(), action: "cleared"})
	base.state = Cleared
	return base
}

// capturing reports, for a survey, whether Drizz is registered for an agent's turn events, swallowing a read failure
// as "not capturing" so a survey never fails on one unreadable file.
func (installer Installer) capturing(scope context.Context, target agent.Agent) bool {
	if !target.Hooking().Supported() {
		return false
	}
	captures, failure := installer.store.Captures(target)
	if failure != nil {
		installer.logger.WarnContext(scope, "integration.inspect.failed", slog.String("agent", target.Kind().String()))
		return false
	}
	return captures
}
