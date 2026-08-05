package cloud

import (
	"context"
	"net/http"

	"github.com/DrizzDev/platform/internal/identity/application/grant"
)

const route = "/api/organizations/me"

func (client Client) fetch(scope context.Context, credential grant.Credential) (*http.Response, error) {
	request, failure := http.NewRequestWithContext(scope, http.MethodGet, client.base+route, nil)
	if failure != nil {
		return nil, failure
	}
	request.Header.Set("Authorization", "Bearer "+string(credential.Token()))
	request.Header.Set("Accept", "application/json")
	return client.agent.Do(request)
}
