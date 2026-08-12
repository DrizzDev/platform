package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
)

func (command Command) devices(scope context.Context) *cobra.Command {
	entry, _ := catalog.New().Lookup(catalog.Devices)
	return &cobra.Command{
		Use:   command.slug(entry.Title()),
		Short: entry.Summary(),
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
