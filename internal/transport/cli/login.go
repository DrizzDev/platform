package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/code"
)

// Session runs one browser login. A returned error is a composition failure the
// transport maps to a temporary-unavailable outcome; otherwise the Result
// carries the trusted context or a code-only failure. It is owned by the CLI.
type Session interface {
	Run(context.Context) (login.Result, error)
}

// denied carries an already-mapped failure message so the process boundary
// prints it once and exits non-zero.
type denied struct {
	message string
}

func (denied denied) Error() string {
	return denied.message
}

func (command Command) login(scope context.Context) *cobra.Command {
	return &cobra.Command{
		Use:                "login",
		Short:              "Sign in to Drizz in the system browser",
		DisableFlagParsing: true,
		RunE: func(root *cobra.Command, arguments []string) error {
			session := command.options.Login
			switch {
			case len(arguments) == 1 && arguments[0] == "--device":
				session = command.options.Device
			case len(arguments) > 0:
				return denied{message: "Usage: drizz login [--device]"}
			}
			result, failure := session.Run(scope)
			if failure != nil {
				return denied{message: command.notify(code.Unavailable)}
			}
			if result.Failed() {
				value, _ := result.Failure()
				return denied{message: command.notify(value.Code())}
			}
			_, failure = fmt.Fprintln(root.OutOrStdout(), command.welcome(result.Organization()))
			return failure
		},
	}
}

func (Command) welcome(organization string) string {
	if organization == "" {
		return "Signed in to Drizz."
	}
	return "Signed in to Drizz (" + organization + ")."
}

func (Command) notify(kind code.Code) string {
	switch kind {
	case code.Required:
		return "Sign in to Drizz by running drizz login."
	case code.Cancelled:
		return "Drizz sign-in was cancelled."
	case code.Forbidden:
		return "Drizz access is not allowed."
	case code.Unavailable:
		return "Drizz authentication is temporarily unavailable. Try again."
	case code.Rejected:
		return "Drizz could not verify the sign-in. Try again."
	case code.Conflict:
		return "Another Drizz account is signed in. Run drizz logout first."
	case code.Storage:
		return "Secure credential storage is unavailable on this computer."
	case code.Partial:
		return "Signed out locally, but remote revocation could not be confirmed."
	case code.Failed:
		return "Drizz could not complete the authentication request."
	}
	return "Drizz could not complete the authentication request."
}
