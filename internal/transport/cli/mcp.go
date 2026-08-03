package cli

import (
	"context"

	"github.com/spf13/cobra"
)

func (command Command) mcp(scope context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "mcp",
		Args:  cobra.NoArgs,
		Short: "Run the local MCP server over standard input and output",
		RunE: func(_ *cobra.Command, _ []string) error {
			return command.options.MCP.Run(scope)
		},
	}
}
