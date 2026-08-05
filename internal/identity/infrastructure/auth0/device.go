package auth0

import (
	"context"
	"errors"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"github.com/DrizzDev/platform/internal/identity/application/device"
	"github.com/DrizzDev/platform/internal/identity/application/login"
)

var _ device.Provider = Authorizer{}

// Request issues a device-authorization challenge bound to the Platform API
// audience.
func (authorizer Authorizer) Request(scope context.Context) (device.Instruction, error) {
	response, failure := authorizer.config.DeviceAuth(
		oidc.ClientContext(scope, authorizer.agent),
		oauth2.SetAuthURLParam("audience", authorizer.audience),
	)
	if failure != nil {
		return device.Instruction{}, authorizer.classify(failure)
	}
	verification := response.VerificationURIComplete
	if verification == "" {
		verification = response.VerificationURI
	}
	return device.Instruction{
		Code:         response.DeviceCode,
		User:         response.UserCode,
		Verification: verification,
		Interval:     time.Duration(response.Interval) * time.Second,
		Expiry:       response.Expiry,
	}, nil
}

// Await polls the token endpoint until the user approves or the challenge
// resolves, then validates the resulting ID token.
func (authorizer Authorizer) Await(scope context.Context, instruction device.Instruction) (login.Token, error) {
	response := &oauth2.DeviceAuthResponse{
		DeviceCode: instruction.Code,
		Interval:   int64(instruction.Interval / time.Second),
		Expiry:     instruction.Expiry,
	}
	token, failure := authorizer.config.DeviceAccessToken(oidc.ClientContext(scope, authorizer.agent), response)
	if failure != nil {
		return login.Token{}, authorizer.classify(failure)
	}
	return authorizer.settle(scope, token)
}

// settle validates the ID token from a device exchange and builds the login
// token. The device grant carries no nonce, so only the signature, issuer,
// audience, and expiry are checked before the subject is trusted.
func (authorizer Authorizer) settle(scope context.Context, token *oauth2.Token) (login.Token, error) {
	raw, present := token.Extra("id_token").(string)
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
	return authorizer.respond(seal{token: token, subject: identity.Subject, issued: identity.IssuedAt})
}
