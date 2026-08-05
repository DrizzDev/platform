package logout

import (
	"context"

	"github.com/DrizzDev/platform/internal/identity/application/credential"
)

// Vault reads the active credential so logout can revoke it. It reports Missing
// when nothing is signed in, keeping logout idempotent.
type Vault interface {
	Active(context.Context) (credential.Record, error)
}
