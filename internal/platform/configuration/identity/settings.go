package identity

import "slices"

// Settings carries the Auth0 Platform application facts the login flow needs.
// The values are provisioned per environment; an empty value fails the login at
// use, not at process start, so other commands remain available.
type Settings struct {
	issuer   string
	client   string
	audience string
	redirect string
	session  string
	cloud    string
	scopes   []string
}

func (settings Settings) Issuer() string {
	return settings.issuer
}

func (settings Settings) Client() string {
	return settings.client
}

func (settings Settings) Audience() string {
	return settings.audience
}

func (settings Settings) Redirect() string {
	return settings.redirect
}

func (settings Settings) Session() string {
	return settings.session
}

func (settings Settings) Cloud() string {
	return settings.cloud
}

func (settings Settings) Scopes() []string {
	return slices.Clone(settings.scopes)
}
