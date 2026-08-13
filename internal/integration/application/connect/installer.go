package connect

import (
	"context"
	"errors"
	"log/slog"

	"github.com/DrizzDev/platform/internal/integration/domain/agent"
	"github.com/DrizzDev/platform/internal/integration/domain/server"
)

const (
	// Name is the key the Drizz entry is filed under in every agent's configuration. It is exported so a store
	// adapter uses the same name when it looks the entry up.
	Name = "drizz"
	// Launch is the Drizz subcommand that runs the stdio MCP server an agent starts.
	Launch = "mcp"
)

// Installer wires supported agent applications to this installed Drizz. It reads the agent catalog, resolves the Drizz
// executable, and drives each agent's configuration through the store, recording every change. Both the command line
// and future surfaces share it, so an agent is connected the same way however the request arrives.
type Installer struct {
	catalog  agent.Catalog
	resolver resolver
	store    store
	recorder recorder
	monitor  monitor
	logger   *slog.Logger
}

// Survey reports every supported agent, whether it is installed, and whether Drizz is already connected, so a person
// can see what a change would touch before making it.
func (installer Installer) Survey(scope context.Context) Report {
	scope, watch := installer.begin(scope, "survey")
	report := Report{}
	defer func() { watch.finish(scope, installer.tally(summary{report: report})) }()

	var outcomes []Outcome
	for _, target := range installer.chosen(Selection{All: true}) {
		outcomes = append(outcomes, installer.inspect(scope, target))
	}
	report = Report{outcomes: outcomes}
	return report
}

// Enable connects Drizz to the selected agents. It resolves the executable and builds the entry once; a failure there
// aborts before any file is touched. Each agent's outcome is reported individually, so one unwritable configuration
// never stops the rest.
func (installer Installer) Enable(scope context.Context, selection Selection) (report Report, failure error) {
	scope, watch := installer.begin(scope, "enable")
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
		outcomes = append(outcomes, installer.wire(scope, Task{Agent: target, Server: entry}))
	}
	report = Report{outcomes: outcomes}
	return report, nil
}

// Disable removes the Drizz entry from the selected agents, leaving every other setting untouched.
func (installer Installer) Disable(scope context.Context, selection Selection) Report {
	scope, watch := installer.begin(scope, "disable")
	report := Report{}
	defer func() { watch.finish(scope, installer.tally(summary{report: report})) }()

	var outcomes []Outcome
	for _, target := range installer.chosen(selection) {
		outcomes = append(outcomes, installer.unwire(scope, target))
	}
	report = Report{outcomes: outcomes}
	return report
}

func (installer Installer) inspect(scope context.Context, target agent.Agent) Outcome {
	present, failure := installer.store.Detect(target)
	if failure != nil {
		installer.logger.WarnContext(scope, "integration.inspect.failed", slog.String("agent", target.Kind().String()))
		return Outcome{kind: target.Kind(), title: target.Title(), state: Failed}
	}
	if !present {
		return Outcome{kind: target.Kind(), title: target.Title(), state: Missing}
	}
	wired, failure := installer.store.Wired(target)
	if failure != nil {
		installer.logger.WarnContext(scope, "integration.inspect.failed", slog.String("agent", target.Kind().String()))
		return Outcome{kind: target.Kind(), title: target.Title(), state: Failed, detail: installer.detail(failure)}
	}
	state := Ready
	if wired {
		state = Connected
	}
	return Outcome{
		kind:      target.Kind(),
		title:     target.Title(),
		state:     state,
		restart:   target.Restart(),
		capturing: installer.capturing(scope, target),
	}
}

func (installer Installer) wire(scope context.Context, job Task) Outcome {
	present, failure := installer.store.Detect(job.Agent)
	if failure != nil || !present {
		return Outcome{kind: job.Agent.Kind(), title: job.Agent.Title(), state: Missing}
	}
	wired, _ := installer.store.Wired(job.Agent)
	if failure := installer.store.Connect(scope, job); failure != nil {
		return Outcome{kind: job.Agent.Kind(), title: job.Agent.Title(), state: installer.grade(failure), detail: installer.detail(failure)}
	}
	installer.inscribe(scope, mark{kind: job.Agent.Kind(), action: "connected"})
	state := Connected
	if wired {
		state = Updated
	}
	return Outcome{kind: job.Agent.Kind(), title: job.Agent.Title(), state: state, note: installer.command(scope, job), restart: job.Agent.Restart()}
}

// command installs the /author command for the agent and returns a short note when the person's own command was kept, or
// when the command could not be written. Installing the command never fails the connection it accompanies.
func (installer Installer) command(scope context.Context, job Task) string {
	conflict, failure := installer.store.Command(scope, job)
	if failure != nil {
		installer.logger.WarnContext(scope, "integration.command.failed", slog.String("agent", job.Agent.Kind().String()))
		return "the /author command could not be written"
	}
	if conflict {
		return "kept your existing /author command; Drizz's authoring command was not installed"
	}
	return ""
}

func (installer Installer) unwire(scope context.Context, target agent.Agent) Outcome {
	present, failure := installer.store.Detect(target)
	if failure != nil || !present {
		return Outcome{kind: target.Kind(), title: target.Title(), state: Missing}
	}
	if failure := installer.store.Disconnect(scope, target); failure != nil {
		return Outcome{kind: target.Kind(), title: target.Title(), state: installer.grade(failure), detail: installer.detail(failure)}
	}
	if failure := installer.store.Uncommand(scope, target); failure != nil {
		installer.logger.WarnContext(scope, "integration.uncommand.failed", slog.String("agent", target.Kind().String()))
	}
	installer.inscribe(scope, mark{kind: target.Kind(), action: "disconnected"})
	return Outcome{kind: target.Kind(), title: target.Title(), state: Removed, restart: target.Restart()}
}

func (installer Installer) chosen(selection Selection) []agent.Agent {
	if selection.All {
		return installer.catalog.List()
	}
	if target, found := installer.catalog.Lookup(selection.Kind); found {
		return []agent.Agent{target}
	}
	return nil
}

func (Installer) grade(failure error) State {
	var carrier interface{ state() State }
	if errors.As(failure, &carrier) {
		return carrier.state()
	}
	return Failed
}

func (Installer) detail(failure error) string {
	var carrier interface{ state() State }
	if errors.As(failure, &carrier) {
		return failure.Error()
	}
	return "the change could not be completed"
}
