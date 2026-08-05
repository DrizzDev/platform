package cli

import (
	"context"

	"github.com/spf13/cobra"
)

func (command Command) root(scope context.Context) *cobra.Command {
	root := &cobra.Command{
		Use:   "drizz",
		Short: "Use Drizz capabilities from agents and developer tools",

		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.SetIn(command.options.Streams.Input)
	root.SetOut(command.options.Streams.Output)
	root.SetErr(command.options.Streams.Failure)

	root.AddCommand(command.version())
	root.AddCommand(command.mcp(scope))
	root.AddCommand(command.login(scope))
	root.AddCommand(command.logout(scope))

	return root
}
