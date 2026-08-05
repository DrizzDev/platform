package attempt

import (
	"errors"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
)

type Attempt struct {
	revision uint64
	epoch    epoch.Epoch
}

type Input struct {
	Revision uint64
	Epoch    epoch.Epoch
}

func New(input Input) (Attempt, error) {
	attempt := Attempt{revision: input.Revision, epoch: input.Epoch}
	if failure := attempt.validate(); failure != nil {
		return Attempt{}, failure
	}
	return attempt, nil
}

func (attempt Attempt) Revision() uint64 {
	return attempt.revision
}

func (attempt Attempt) Epoch() epoch.Epoch {
	return attempt.epoch
}

func (attempt Attempt) validate() error {
	if attempt.revision < 1 {
		return errors.New("attempt revision must be positive")
	}
	return nil
}
