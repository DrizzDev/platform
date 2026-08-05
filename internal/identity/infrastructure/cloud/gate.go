package cloud

import (
	"context"
	"errors"

	"github.com/DrizzDev/platform/internal/identity/application/grant"
	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/application/organization"
)

var _ login.Authority = Client{}

// Authorize gates a sign-in on active organization membership, using the access
// token the sign-in just obtained — no refresh. A definitive denial blocks the
// login; an unresolved or unavailable outcome permits local sign-in, since every
// cloud operation re-authorizes at the cloud regardless. On success it returns
// the organization to surface to the user.
func (client Client) Authorize(scope context.Context, bearer login.Grant) (login.Tenant, error) {
	credential, failure := grant.New(grant.Input{Token: bearer.Token, Expiry: bearer.Expiry})
	if failure != nil {
		return login.Tenant{}, failure
	}
	member, failure := client.call(scope, credential)
	var forbidden organization.Forbidden
	if errors.As(failure, &forbidden) {
		return login.Tenant{}, login.Forbidden{}
	}
	if failure != nil {
		return login.Tenant{}, failure
	}
	return login.Tenant{Name: member.name}, nil
}
