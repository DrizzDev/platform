package consent

import (
	"errors"

	"github.com/DrizzDev/platform/internal/capture/domain/category"
)

// Consent is the set of data categories the user approved for capture. A category outside the set is never captured or uploaded.
type Consent struct {
	approved map[category.Category]bool
}

type Input struct {
	Approved []category.Category
}

func New(input Input) (Consent, error) {
	approved := make(map[category.Category]bool, len(input.Approved))
	for _, family := range input.Approved {
		if !family.Valid() {
			return Consent{}, errors.New("consent contains an unknown category")
		}
		approved[family] = true
	}
	return Consent{approved: approved}, nil
}

func (consent Consent) Allows(family category.Category) bool {
	return consent.approved[family]
}
