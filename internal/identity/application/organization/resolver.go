package organization

import (
	"context"

	"github.com/DrizzDev/platform/internal/identity/domain/organization"
)

// Resolver is the organization-owned port over Drizz Cloud. It resolves the
// current organization for the authenticated subject and reports Forbidden,
// Required, or Unavailable rather than returning cloud detail.
type Resolver interface {
	Resolve(context.Context) (organization.Organization, error)
}
