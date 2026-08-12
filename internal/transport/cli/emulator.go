package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/DrizzDev/platform/internal/capability/application/operator"
	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
)

func (command Command) images(scope context.Context) *cobra.Command {
	entry, _ := catalog.New().Lookup(catalog.Images)
	return &cobra.Command{
		Use:   command.slug(entry.Title()),
		Short: entry.Summary(),
		Args:  cobra.NoArgs,
		RunE: func(root *cobra.Command, arguments []string) error {
			images, failure := command.options.Perform.Images(scope)
			if failure != nil {
				return denied{message: command.explain(failure)}
			}
			for _, name := range images.Names {
				if _, failure := fmt.Fprintln(root.OutOrStdout(), name); failure != nil {
					return failure
				}
			}
			return nil
		},
	}
}

func (command Command) boot(scope context.Context) *cobra.Command {
	entry, _ := catalog.New().Lookup(catalog.Boot)
	return &cobra.Command{
		Use:   command.slug(entry.Title()) + " <image>",
		Short: entry.Summary(),
		Args:  cobra.ExactArgs(1),
		RunE: func(root *cobra.Command, arguments []string) error {
			if _, failure := command.options.Perform.Boot(scope, operator.Image{Name: arguments[0]}); failure != nil {
				return denied{message: command.explain(failure)}
			}
			_, failure := fmt.Fprintln(root.OutOrStdout(), "Booting.")
			return failure
		},
	}
}

func (command Command) pause(scope context.Context) *cobra.Command {
	return command.single(scope, gesture{name: catalog.Pause, done: "Paused.", perform: command.options.Perform.Pause})
}

func (command Command) resume(scope context.Context) *cobra.Command {
	return command.single(scope, gesture{name: catalog.Resume, done: "Resumed.", perform: command.options.Perform.Resume})
}
