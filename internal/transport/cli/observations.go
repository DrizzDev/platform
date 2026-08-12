package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/DrizzDev/platform/internal/capability/application/operator"
	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
	"github.com/DrizzDev/platform/internal/platform/filesystem"
)

func (command Command) snapshot(scope context.Context) *cobra.Command {
	entry, _ := catalog.New().Lookup(catalog.Snapshot)
	return &cobra.Command{
		Use:   command.slug(entry.Title()) + " <serial>",
		Short: entry.Summary(),
		Args:  cobra.ExactArgs(1),
		RunE: func(root *cobra.Command, arguments []string) error {
			shot, failure := command.options.Perform.Snapshot(scope, operator.Target{Serial: arguments[0]})
			if failure != nil {
				return denied{message: command.explain(failure)}
			}
			path, failure := filesystem.New().Save(filesystem.File{Extension: strings.ToLower(shot.Format), Content: shot.Image})
			if failure != nil {
				return failure
			}
			if _, failure := fmt.Fprintln(root.OutOrStdout(), path); failure != nil {
				return failure
			}
			_, failure = fmt.Fprintln(root.OutOrStdout(), shot.Hierarchy)
			return failure
		},
	}
}

func (command Command) hierarchy(scope context.Context) *cobra.Command {
	entry, _ := catalog.New().Lookup(catalog.Hierarchy)
	return &cobra.Command{
		Use:   command.slug(entry.Title()) + " <serial>",
		Short: entry.Summary(),
		Args:  cobra.ExactArgs(1),
		RunE: func(root *cobra.Command, arguments []string) error {
			tree, failure := command.options.Perform.Hierarchy(scope, operator.Target{Serial: arguments[0]})
			if failure != nil {
				return denied{message: command.explain(failure)}
			}
			_, failure = fmt.Fprintln(root.OutOrStdout(), tree.Hierarchy)
			return failure
		},
	}
}

func (command Command) dimensions(scope context.Context) *cobra.Command {
	entry, _ := catalog.New().Lookup(catalog.Dimensions)
	return &cobra.Command{
		Use:   command.slug(entry.Title()) + " <serial>",
		Short: entry.Summary(),
		Args:  cobra.ExactArgs(1),
		RunE: func(root *cobra.Command, arguments []string) error {
			extent, failure := command.options.Perform.Dimensions(scope, operator.Target{Serial: arguments[0]})
			if failure != nil {
				return denied{message: command.explain(failure)}
			}
			_, failure = fmt.Fprintf(root.OutOrStdout(), "%dx%d\n", extent.Width, extent.Height)
			return failure
		},
	}
}
