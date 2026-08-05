package organization

import "github.com/DrizzDev/platform/internal/identity/domain/identifier"

type Organization struct {
	identifier.Identifier
}

func New(value string) (Organization, error) {
	inner, failure := identifier.New(value)
	if failure != nil {
		return Organization{}, failure
	}
	return Organization{Identifier: inner}, nil
}
