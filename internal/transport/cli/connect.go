package cli

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/DrizzDev/platform/internal/integration/application/connect"
	"github.com/DrizzDev/platform/internal/integration/domain/agent"
)

// Connector wires this Drizz into agent applications and reports the result of each change. It is owned by the CLI;
// the host supplies the implementation.
type Connector interface {
	Survey(context.Context) connect.Report
	Enable(context.Context, connect.Selection) (connect.Report, error)
	Disable(context.Context, connect.Selection) connect.Report
	Capture(context.Context, connect.Selection) (connect.Report, error)
	Uncapture(context.Context, connect.Selection) connect.Report
}

func (command Command) connect(scope context.Context) *cobra.Command {
	group := &cobra.Command{
		Use:   "connect",
		Short: "Connect Drizz to agent applications such as Claude and Codex",
	}
	group.AddCommand(command.roster(scope))
	group.AddCommand(command.enable(scope))
	group.AddCommand(command.disable(scope))
	group.AddCommand(command.capture(scope))
	group.AddCommand(command.uncapture(scope))
	return group
}

func (command Command) capture(scope context.Context) *cobra.Command {
	var approved bool
	entry := &cobra.Command{
		Use:   "capture [agent]",
		Short: "Let Drizz record your prompts and responses for context, for one agent or all",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(root *cobra.Command, arguments []string) error {
			if !approved && !command.consented(root) {
				_, failure := fmt.Fprintln(root.OutOrStdout(), "Cancelled. Re-run with --yes to enable capture.")
				return failure
			}
			report, failure := command.options.Connect.Capture(scope, command.selection(arguments))
			if failure != nil {
				return denied{message: "Drizz could not locate its own program to enable capture."}
			}
			return command.emit(report)
		},
	}
	entry.Flags().BoolVar(&approved, "yes", false, "Enable capture without asking for confirmation")
	return entry
}

func (command Command) uncapture(scope context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "uncapture [agent]",
		Short: "Stop Drizz recording prompts and responses, for one agent or all",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, arguments []string) error {
			return command.emit(command.options.Connect.Uncapture(scope, command.selection(arguments)))
		},
	}
}

func (command Command) roster(scope context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Show which agents are installed and whether Drizz is connected",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return command.emit(command.options.Connect.Survey(scope))
		},
	}
}

func (command Command) enable(scope context.Context) *cobra.Command {
	var plain bool
	entry := &cobra.Command{
		Use:   "enable [agent]",
		Short: "Connect Drizz to one agent, or every detected agent, and record context unless --no-capture",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, arguments []string) error {
			selection := command.selection(arguments)
			report, failure := command.options.Connect.Enable(scope, selection)
			if failure != nil {
				return denied{message: "Drizz could not locate its own program to connect."}
			}
			if failure := command.emit(report); failure != nil {
				return failure
			}
			if plain {
				return nil
			}
			capture, failure := command.options.Connect.Capture(scope, selection)
			if failure != nil {
				return denied{message: "Drizz could not locate its own program to enable capture."}
			}
			return command.emit(capture)
		},
	}
	entry.Flags().BoolVar(&plain, "no-capture", false, "Connect without recording prompts and responses")
	return entry
}

func (command Command) disable(scope context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "disable [agent]",
		Short: "Remove Drizz from one agent, or from every detected agent",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, arguments []string) error {
			return command.emit(command.options.Connect.Disable(scope, command.selection(arguments)))
		},
	}
}

func (Command) selection(arguments []string) connect.Selection {
	if len(arguments) == 1 {
		return connect.Selection{Kind: agent.Kind(arguments[0])}
	}
	return connect.Selection{All: true}
}

// prompting is one yes/no confirmation: where to ask and what to say.
type prompting struct {
	root     *cobra.Command
	question string
}

func (command Command) consented(root *cobra.Command) bool {
	return command.confirm(prompting{
		root:     root,
		question: "Let Drizz record your prompts and responses for context? This captures what you type. [y/N]: ",
	})
}

func (Command) confirm(ask prompting) bool {
	_, _ = fmt.Fprint(ask.root.OutOrStdout(), ask.question)
	scanner := bufio.NewScanner(ask.root.InOrStdin())
	if !scanner.Scan() {
		return false
	}
	answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
	return answer == "y" || answer == "yes"
}

func (command Command) emit(report connect.Report) error {
	for _, outcome := range report.Outcomes() {
		if _, failure := fmt.Fprintln(command.options.Streams.Output, command.summarize(outcome)); failure != nil {
			return failure
		}
	}
	return nil
}

func (command Command) summarize(outcome connect.Outcome) string {
	line := outcome.Title() + ": " + command.phrase(outcome.State())
	if outcome.Capturing() && command.settled(outcome.State()) {
		line += ", capturing context"
	}
	if outcome.Detail() != "" {
		line += " (" + outcome.Detail() + ")"
	}
	if outcome.Restart() && command.changed(outcome.State()) {
		line += " — restart " + outcome.Title() + " to apply"
	}
	return line
}

func (Command) phrase(state connect.State) string {
	switch state {
	case connect.Connected:
		return "connected"
	case connect.Updated:
		return "updated"
	case connect.Removed:
		return "removed"
	case connect.Captured:
		return "capturing context"
	case connect.Cleared:
		return "capture stopped"
	case connect.Incapable:
		return "no hook support"
	case connect.Ready:
		return "installed, not connected"
	case connect.Missing:
		return "not installed"
	case connect.Failed:
		return "could not be changed"
	}
	return "unknown"
}

// settled reports whether a state describes an agent's standing in a survey — connected or merely installed — as
// opposed to the result of a just-performed change, so a capture note is added only where it makes sense.
func (Command) settled(state connect.State) bool {
	switch state {
	case connect.Connected, connect.Ready:
		return true
	case connect.Updated, connect.Removed, connect.Captured, connect.Cleared, connect.Incapable, connect.Missing, connect.Failed:
		return false
	default:
		return false
	}
}

// changed reports whether a state reflects a write that a restart-required agent must reload.
func (Command) changed(state connect.State) bool {
	switch state {
	case connect.Connected, connect.Updated, connect.Removed, connect.Captured, connect.Cleared:
		return true
	case connect.Ready, connect.Missing, connect.Incapable, connect.Failed:
		return false
	default:
		return false
	}
}
