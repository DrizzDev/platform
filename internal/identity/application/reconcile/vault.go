package reconcile

import (
	"context"

	"github.com/DrizzDev/platform/internal/identity/application/credential"
)

// Vault is the credential store the flow deletes reconciled candidates from.
type Vault interface {
	Delete(context.Context, credential.Key) error
}
