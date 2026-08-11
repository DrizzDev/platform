package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/DrizzDev/platform/internal/capability/application/operator"
	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
	"github.com/DrizzDev/platform/internal/platform/filesystem"
)

func (command Command) screenshot(scope context.Context) *cobra.Command {
	entry, _ := catalog.New().Lookup(catalog.Screenshot)
	return &cobra.Command{
		Use:   command.slug(entry.Title()) + " <serial>",
		Short: entry.Summary(),
		Args:  cobra.ExactArgs(1),
		RunE: func(root *cobra.Command, arguments []string) error {
			shot, failure := command.options.Perform.Screenshot(scope, operator.Target{Serial: arguments[0]})
			if failure != nil {
				return denied{message: command.explain(failure)}
			}
			path, failure := filesystem.New().Save(filesystem.File{
				Extension: strings.ToLower(shot.Format),
				Content:   shot.Image,
			})
			if failure != nil {
				return failure
			}
			_, failure = fmt.Fprintln(root.OutOrStdout(), path)
			return failure
		},
	}
}

// explain maps a capability failure to one safe, actionable message. A refusal carries the code's own detail; anything
// else is a composition failure the person can only retry.
func (Command) explain(failure error) string {
	var refusal operator.Refusal
	if errors.As(failure, &refusal) {
		return refusal.Code.Detail()
	}
	return "Drizz could not complete the request. Try again."
}
