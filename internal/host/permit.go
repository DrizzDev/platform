package host

import (
	"context"

	"github.com/DrizzDev/platform/internal/identity/application/login"
)

var _ login.Authority = permit{}

// permit authorizes every sign-in. It stands in when no Drizz Cloud endpoint is
// configured, so login keeps working as authentication-only until the cloud
// authority is provisioned.
type permit struct{}

func (permit) Authorize(context.Context, login.Grant) (login.Tenant, error) {
	return login.Tenant{}, nil
}
