package device

import (
	"context"

	"github.com/DrizzDev/platform/internal/identity/application/login"
)

// Provider is the device-owned port over the OAuth device-authorization grant.
// Request issues the challenge; Await polls the token endpoint until the user
// approves, denies, or the challenge expires, then returns the validated token.
type Provider interface {
	Request(context.Context) (Instruction, error)
	Await(context.Context, Instruction) (login.Token, error)
}
