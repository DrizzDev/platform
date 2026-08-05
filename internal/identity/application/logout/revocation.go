package logout

import (
	"context"

	"github.com/DrizzDev/platform/internal/identity/application/credential"
)

// Revocation invalidates the credential at the provider after the local clear.
// A failure is reported as a partial logout, never a lost local sign-out.
type Revocation interface {
	Revoke(context.Context, credential.Record) error
}
