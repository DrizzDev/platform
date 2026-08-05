package session

import (
	"context"

	"github.com/DrizzDev/platform/internal/identity/application/credential"
)

// Vault reads the active credential the flow renews. It reports Missing when
// nothing is signed in, which the flow resolves into a required sign-in.
type Vault interface {
	Active(context.Context) (credential.Record, error)
}
