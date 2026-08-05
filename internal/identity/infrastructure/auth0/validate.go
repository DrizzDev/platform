package auth0

import (
	"errors"
	"slices"
)

func (options Options) validate() error {
	switch {
	case options.Agent == nil:
		return errors.New("auth0 http client is required")
	case options.Issuer == "":
		return errors.New("auth0 issuer is required")
	case options.Client == "":
		return errors.New("auth0 client is required")
	case options.Audience == "":
		return errors.New("auth0 audience is required")
	case options.Redirect == "":
		return errors.New("auth0 redirect is required")
	case !options.Method.Valid():
		return errors.New("auth0 method is invalid")
	case !slices.Contains(options.Scopes, "openid"):
		return errors.New("auth0 scopes must include openid")
	case !slices.Contains(options.Scopes, "offline_access"):
		return errors.New("auth0 scopes must include offline_access")
	default:
		return nil
	}
}
