package auth0

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/DrizzDev/platform/internal/identity/application/credential"
	"github.com/DrizzDev/platform/internal/identity/application/logout"
)

var _ logout.Revocation = Authorizer{}

// Revoke invalidates the refresh token at Auth0. The public client sends no
// secret; any non-success response is a failed revocation the caller reports as
// a partial logout.
func (authorizer Authorizer) Revoke(scope context.Context, record credential.Record) error {
	form := url.Values{"client_id": {authorizer.client}, "token": {string(record.Refresh())}}
	request, failure := http.NewRequestWithContext(scope, http.MethodPost,
		authorizer.issuer+"oauth/revoke", strings.NewReader(form.Encode()))
	if failure != nil {
		return failure
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, failure := authorizer.agent.Do(request)
	if failure != nil {
		return failure
	}
	defer func() { _ = response.Body.Close() }()
	_, _ = io.Copy(io.Discard, response.Body)
	if response.StatusCode != http.StatusOK {
		return errors.New("revocation was refused")
	}
	return nil
}
