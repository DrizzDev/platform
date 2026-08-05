package grant

import (
	"errors"
	"slices"
	"time"
)

type Credential struct {
	token  []byte
	expiry time.Time
}

type Input struct {
	Token  []byte
	Expiry time.Time
}

func New(input Input) (Credential, error) {
	credential := Credential{
		expiry: input.Expiry,
		token:  slices.Clone(input.Token),
	}
	if failure := credential.validate(); failure != nil {
		return Credential{}, failure
	}
	return credential, nil
}

func (credential Credential) Token() []byte {
	return slices.Clone(credential.token)
}

func (credential Credential) Expiry() time.Time {
	return credential.expiry
}

func (credential Credential) Expired(now time.Time) bool {
	return !now.Before(credential.expiry)
}

func (credential Credential) validate() error {
	switch {
	case len(credential.token) == 0:
		return errors.New("grant token is required")
	case credential.expiry.IsZero():
		return errors.New("grant expiry is required")
	default:
		return nil
	}
}
