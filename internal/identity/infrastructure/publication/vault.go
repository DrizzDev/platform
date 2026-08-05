package publication

import (
	"context"

	"github.com/DrizzDev/platform/internal/identity/application/credential"
)

// Vault is the credential store the publisher depends on, satisfied by the
// operating-system vault adapter.
type Vault interface {
	Write(context.Context, credential.Record) error
	Read(context.Context, credential.Key) (credential.Record, error)
	Delete(context.Context, credential.Key) error
}
