package cloud

import (
	"context"
	"errors"
	"net/http"

	"github.com/DrizzDev/platform/internal/identity/application/grant"
	"github.com/DrizzDev/platform/internal/identity/application/organization"
	"github.com/DrizzDev/platform/internal/identity/application/session"
	tenant "github.com/DrizzDev/platform/internal/identity/domain/organization"
)

var _ organization.Resolver = Client{}

// Client resolves the current organization from Drizz Cloud. It maps the cloud
// outcome to the organization contract without leaking cloud detail.
type Client struct {
	agent  *http.Client
	source *source
	base   string
}

// Resolve renews a session-backed access token and resolves the current
// organization. It is the port for callers that hold no fresh grant of their own.
func (client Client) Resolve(scope context.Context) (tenant.Organization, error) {
	if client.source == nil {
		return tenant.Organization{}, organization.Unavailable{}
	}
	credential, failure := client.source.acquire(scope)
	if failure != nil {
		return tenant.Organization{}, client.reauth(failure)
	}
	member, failure := client.call(scope, credential)
	if failure != nil {
		return tenant.Organization{}, failure
	}
	return member.org, nil
}

// call performs one authenticated request and closes the response body.
func (client Client) call(scope context.Context, credential grant.Credential) (membership, error) {
	response, failure := client.fetch(scope, credential)
	if failure != nil {
		return membership{}, organization.Unavailable{}
	}
	defer func() { _ = response.Body.Close() }()
	return client.interpret(response)
}

// reauth maps a token-acquisition failure. A broken store is unavailable;
// anything else means the session must be re-established.
func (Client) reauth(failure error) error {
	var storage session.Storage
	if errors.As(failure, &storage) {
		return organization.Unavailable{}
	}
	return organization.Required{}
}
