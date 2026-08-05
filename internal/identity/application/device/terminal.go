package device

import (
	"context"

	"github.com/DrizzDev/platform/internal/identity/application/login"
)

var _ login.Establishment = terminal{}

// terminal establishes a sign-in through the OAuth device-authorization grant:
// it requests a challenge, shows the user where to enter the code, and polls
// until the token is granted.
type terminal struct {
	provider Provider
	display  Display
}

func (terminal terminal) Establish(scope context.Context) (login.Token, error) {
	instruction, failure := terminal.provider.Request(scope)
	if failure != nil {
		return login.Token{}, failure
	}
	if failure := terminal.display.Show(scope, instruction); failure != nil {
		return login.Token{}, failure
	}
	return terminal.provider.Await(scope, instruction)
}
