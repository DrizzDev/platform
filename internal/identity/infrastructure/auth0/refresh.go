package auth0

import (
	"context"
	"errors"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"github.com/DrizzDev/platform/internal/identity/application/credential"
	"github.com/DrizzDev/platform/internal/identity/application/session"
)

var _ session.Refresh = Authorizer{}

// Renew exchanges the stored refresh token for a rotated credential. The refresh
// runs once and is never retried: a rotating token may already have been
// consumed, so any non-cancellation failure is reported as uncertain and the
// caller must sign in again.
func (authorizer Authorizer) Renew(scope context.Context, record credential.Record) (session.Renewal, error) {
	scope = oidc.ClientContext(scope, authorizer.agent)
	source := authorizer.config.TokenSource(scope, &oauth2.Token{RefreshToken: string(record.Refresh())})
	token, failure := source.Token()
	if failure != nil {
		return session.Renewal{}, authorizer.renewal(failure)
	}
	refresh := []byte(token.RefreshToken)
	access := []byte(token.AccessToken)
	if len(refresh) == 0 || len(access) == 0 || token.Expiry.IsZero() {
		return session.Renewal{}, session.Uncertain{}
	}
	return session.Renewal{Refresh: refresh, Access: access, Expiry: token.Expiry}, nil
}

// renewal maps a refresh failure to a port outcome. Cancellation keeps its
// identity; every other failure is uncertain, since the token may have been
// spent and must not be replayed.
func (Authorizer) renewal(failure error) error {
	if errors.Is(failure, context.Canceled) || errors.Is(failure, context.DeadlineExceeded) {
		return failure
	}
	return session.Uncertain{}
}
