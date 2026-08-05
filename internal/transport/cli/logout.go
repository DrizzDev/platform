package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/DrizzDev/platform/internal/identity/application/logout"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/code"
)

// Departure runs one logout. A returned error is a composition failure the
// transport maps to a temporary-unavailable outcome; otherwise the Result
// carries a clean or partial logout. It is owned by the CLI.
type Departure interface {
	Run(context.Context) (logout.Result, error)
}

func (command Command) logout(scope context.Context) *cobra.Command {
	return &cobra.Command{
		Use:                "logout",
		Short:              "Sign out of Drizz",
		DisableFlagParsing: true,
		RunE: func(root *cobra.Command, arguments []string) error {
			if len(arguments) > 0 {
				return denied{message: "Usage: drizz logout"}
			}
			result, failure := command.options.Logout.Run(scope)
			if failure != nil {
				return denied{message: command.notify(code.Unavailable)}
			}
			if result.Failed() {
				value, _ := result.Failure()
				return denied{message: command.notify(value.Code())}
			}
			_, failure = fmt.Fprintln(root.OutOrStdout(), "Signed out of Drizz.")
			return failure
		},
	}
}
