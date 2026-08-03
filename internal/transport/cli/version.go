package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (command Command) version() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Args:  cobra.NoArgs,
		Short: "Print version information",
		RunE: func(root *cobra.Command, _ []string) error {
			identity := command.options.Release
			_, failure := fmt.Fprintf(
				root.OutOrStdout(),
				"%s %s (%s)\n",
				identity.Name(),
				identity.Version(),
				identity.Revision(),
			)
			return failure
		},
	}
}
