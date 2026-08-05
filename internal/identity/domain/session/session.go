package session

import "github.com/DrizzDev/platform/internal/identity/domain/identifier"

type Session struct {
	identifier.Identifier
}

func New(value string) (Session, error) {
	inner, failure := identifier.New(value)
	if failure != nil {
		return Session{}, failure
	}
	return Session{Identifier: inner}, nil
}
