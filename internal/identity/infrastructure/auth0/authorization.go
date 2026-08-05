package auth0

import (
	"context"
	"errors"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/domain/method"
)

var _ login.Authorization = Authorizer{}

// Authorizer implements the login Authorization port over Auth0 using
// Authorization Code with PKCE S256 and OIDC ID-token validation.
type Authorizer struct {
	verifier *oidc.IDTokenVerifier
	agent    *http.Client
	issuer   string
	client   string
	audience string
	method   method.Method
	config   oauth2.Config
}

func (authorizer Authorizer) Begin(_ context.Context, secret login.Secret) (login.Redirect, error) {
	url := authorizer.config.AuthCodeURL(
		secret.State,
		oidc.Nonce(secret.Nonce),
		oauth2.S256ChallengeOption(secret.Verifier),
		oauth2.SetAuthURLParam("audience", authorizer.audience),
	)
	return login.Redirect{URL: url}, nil
}

func (authorizer Authorizer) Finish(scope context.Context, exchange login.Exchange) (login.Token, error) {
	scope = oidc.ClientContext(scope, authorizer.agent)
	token, failure := authorizer.config.Exchange(scope, exchange.Callback.Code, oauth2.VerifierOption(exchange.Secret.Verifier))
	if failure != nil {
		return login.Token{}, authorizer.classify(failure)
	}
	return authorizer.verify(scope, verification{token: token, secret: exchange.Secret, callback: exchange.Callback})
}

// classify maps a token-endpoint failure to a port outcome. Cancellation keeps
// its identity; a provider error response is a rejected sign-in; a transport
// failure is temporary.
func (Authorizer) classify(failure error) error {
	if errors.Is(failure, context.Canceled) || errors.Is(failure, context.DeadlineExceeded) {
		return failure
	}
	var retrieve *oauth2.RetrieveError
	if errors.As(failure, &retrieve) {
		return login.Rejected{}
	}
	return login.Unavailable{}
}
