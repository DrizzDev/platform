package subject

import "github.com/DrizzDev/platform/internal/identity/domain/identifier"

type Subject struct {
	identifier.Identifier
}

func New(value string) (Subject, error) {
	inner, failure := identifier.New(value)
	if failure != nil {
		return Subject{}, failure
	}
	return Subject{Identifier: inner}, nil
}
