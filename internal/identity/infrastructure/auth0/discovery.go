package auth0

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

func New(scope context.Context, options Options) (Authorizer, error) {
	if failure := options.validate(); failure != nil {
		return Authorizer{}, failure
	}
	provider, failure := oidc.NewProvider(oidc.ClientContext(scope, options.Agent), options.Issuer)
	if failure != nil {
		return Authorizer{}, failure
	}
	return Authorizer{
		config: oauth2.Config{
			ClientID:    options.Client,
			RedirectURL: options.Redirect,
			Scopes:      options.Scopes,
			Endpoint:    provider.Endpoint(),
		},
		verifier: provider.Verifier(&oidc.Config{ClientID: options.Client}),
		agent:    options.Agent,
		issuer:   options.Issuer,
		client:   options.Client,
		audience: options.Audience,
		method:   options.Method,
	}, nil
}
