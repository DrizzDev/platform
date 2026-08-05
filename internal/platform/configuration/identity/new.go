package identity

import "strings"

type Input struct {
	Issuer   string
	Client   string
	Audience string
	Redirect string
	Session  string
	Cloud    string
	Scopes   string
}

func New(input Input) Settings {
	var scopes []string
	if input.Scopes != "" {
		scopes = strings.Split(input.Scopes, ",")
	}
	return Settings{
		issuer:   input.Issuer,
		client:   input.Client,
		audience: input.Audience,
		redirect: input.Redirect,
		session:  input.Session,
		cloud:    input.Cloud,
		scopes:   scopes,
	}
}
