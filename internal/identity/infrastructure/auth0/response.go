package auth0

import (
	"time"

	"golang.org/x/oauth2"

	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/domain/subject"
)

type seal struct {
	token   *oauth2.Token
	subject string
	issued  time.Time
}

// respond builds the validated login token. A missing subject, refresh token,
// or access-token expiry is a rejected sign-in.
func (authorizer Authorizer) respond(seal seal) (login.Token, error) {
	account, failure := subject.New(seal.subject)
	if failure != nil {
		return login.Token{}, login.Rejected{}
	}
	refresh := []byte(seal.token.RefreshToken)
	if len(refresh) == 0 || seal.token.Expiry.IsZero() || seal.issued.IsZero() {
		return login.Token{}, login.Rejected{}
	}
	return login.Token{
		Issuer:  authorizer.issuer,
		Client:  authorizer.client,
		Subject: account,
		Method:  authorizer.method,
		Refresh: refresh,
		Access:  []byte(seal.token.AccessToken),
		Issued:  seal.issued,
		Expiry:  seal.token.Expiry,
	}, nil
}
