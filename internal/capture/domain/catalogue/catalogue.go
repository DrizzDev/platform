package catalogue

import (
	"errors"
	"maps"

	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/policy"
)

// Catalogue is the classification: each data category maps to one handling policy. A category with no policy fails
// closed at lookup rather than defaulting to something permissive (SEC-011).
type Catalogue struct {
	policies map[category.Category]policy.Policy
}

type Input struct {
	Policies map[category.Category]policy.Policy
}

func New(input Input) (Catalogue, error) {
	catalogue := Catalogue{policies: maps.Clone(input.Policies)}
	if failure := catalogue.validate(); failure != nil {
		return Catalogue{}, failure
	}
	return catalogue, nil
}

func (catalogue Catalogue) Policy(family category.Category) (policy.Policy, error) {
	rule, found := catalogue.policies[family]
	if !found {
		return policy.Policy{}, errors.New("category is not classified")
	}
	return rule, nil
}

func (catalogue Catalogue) validate() error {
	for family := range catalogue.policies {
		if !family.Valid() {
			return errors.New("catalogue contains an unknown category")
		}
	}
	return nil
}
