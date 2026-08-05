package auth0

import (
	"context"
	"crypto/subtle"
	"errors"

	"golang.org/x/oauth2"

	"github.com/DrizzDev/platform/internal/identity/application/login"
)

type verification struct {
	token    *oauth2.Token
	secret   login.Secret
	callback login.Callback
}

type claims struct {
	Nonce string `json:"nonce"`
}

// validate confirms the callback state, ID-token signature, issuer, audience,
// expiry, and nonce before the response is trusted. Any mismatch is a rejected
// sign-in with the provider cause withheld.
func (authorizer Authorizer) verify(scope context.Context, verification verification) (login.Token, error) {
	if subtle.ConstantTimeCompare([]byte(verification.callback.State), []byte(verification.secret.State)) != 1 {
		return login.Token{}, login.Rejected{}
	}
	raw, present := verification.token.Extra("id_token").(string)
	if !present || raw == "" {
		return login.Token{}, login.Rejected{}
	}
	identity, failure := authorizer.verifier.Verify(scope, raw)
	if failure != nil {
		if errors.Is(failure, context.Canceled) || errors.Is(failure, context.DeadlineExceeded) {
			return login.Token{}, failure
		}
		return login.Token{}, login.Rejected{}
	}
	var payload claims
	if failure := identity.Claims(&payload); failure != nil {
		return login.Token{}, login.Rejected{}
	}
	if subtle.ConstantTimeCompare([]byte(payload.Nonce), []byte(verification.secret.Nonce)) != 1 {
		return login.Token{}, login.Rejected{}
	}
	return authorizer.respond(seal{token: verification.token, subject: identity.Subject, issued: identity.IssuedAt})
}
