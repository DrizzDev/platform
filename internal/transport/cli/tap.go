package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/DrizzDev/platform/internal/capability/application/operator"
	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
)

func (command Command) tap(scope context.Context) *cobra.Command {
	entry, _ := catalog.New().Lookup(catalog.Tap)
	return &cobra.Command{
		Use:   command.slug(entry.Title()) + " <serial> <x> <y>",
		Short: entry.Summary(),
		Args:  cobra.ExactArgs(3),
		RunE: func(root *cobra.Command, arguments []string) error {
			x, failure := strconv.Atoi(arguments[1])
			if failure != nil {
				return denied{message: "The horizontal coordinate must be a whole number."}
			}
			y, failure := strconv.Atoi(arguments[2])
			if failure != nil {
				return denied{message: "The vertical coordinate must be a whole number."}
			}
			if _, failure := command.options.Perform.Tap(scope, operator.Contact{Serial: arguments[0], X: x, Y: y}); failure != nil {
				return denied{message: command.explain(failure)}
			}
			_, failure = fmt.Fprintln(root.OutOrStdout(), "Tapped.")
			return failure
		},
	}
}
