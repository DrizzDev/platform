package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func (command Command) devices(scope context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "devices",
		Short: "List the connected devices",
		Args:  cobra.NoArgs,
		RunE: func(root *cobra.Command, arguments []string) error {
			roster, failure := command.options.Perform.Devices(scope)
			if failure != nil {
				return denied{message: command.explain(failure)}
			}
			for _, address := range roster.Serials {
				if _, failure := fmt.Fprintln(root.OutOrStdout(), address); failure != nil {
					return failure
				}
			}
			return nil
		},
	}
}
