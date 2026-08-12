package cli

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
)

// Receiver records one inbound agent hook event. It is owned by the CLI; the host supplies the implementation. A hook
// must never disrupt the agent that fired it, so the command always exits successfully and the receiver's own failures
// are swallowed after being recorded internally.
type Receiver interface {
	Receive(context.Context, Signal) error
}

// Signal is one hook invocation from an agent: which agent fired it, which turn moment it marks, and the argument
// payload present when the agent passes the event as an argument rather than on standard input.
type Signal struct {
	Agent   string
	Slot    string
	Payload string
}

// hook is the hidden endpoint every agent hook runs. It is not part of the person-facing command set; agents invoke it
// with the agent name and slot Drizz registered, e.g. "drizz hook claude-code prompt".
func (command Command) hook(scope context.Context) *cobra.Command {
	return &cobra.Command{
		Use:    "hook [agent] [slot]",
		Short:  "Receive an agent hook event",
		Hidden: true,
		Args:   cobra.MinimumNArgs(2),
		RunE: func(_ *cobra.Command, arguments []string) error {
			signal := Signal{Agent: arguments[0], Slot: arguments[1]}
			if len(arguments) > 2 {
				signal.Payload = strings.Join(arguments[2:], " ")
			}
			// A hook must never block or fail the agent that fired it, so a receiver failure is swallowed here and the
			// command exits successfully; the receiver records its own drops internally.
			_ = command.options.Receiver.Receive(scope, signal)
			return nil
		},
	}
}
