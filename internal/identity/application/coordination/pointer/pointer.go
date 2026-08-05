package pointer

import (
	"errors"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
)

type Pointer struct {
	key      string
	revision uint64
	epoch    epoch.Epoch
}

type Input struct {
	Key      string
	Revision uint64
	Epoch    epoch.Epoch
}

func New(input Input) (Pointer, error) {
	pointer := Pointer{key: input.Key, revision: input.Revision, epoch: input.Epoch}
	if failure := pointer.validate(); failure != nil {
		return Pointer{}, failure
	}
	return pointer, nil
}

func (pointer Pointer) Key() string {
	return pointer.key
}

func (pointer Pointer) Revision() uint64 {
	return pointer.revision
}

func (pointer Pointer) Epoch() epoch.Epoch {
	return pointer.epoch
}

func (pointer Pointer) validate() error {
	switch {
	case pointer.key == "":
		return errors.New("pointer key is required")
	case pointer.revision < 1:
		return errors.New("pointer revision must be positive")
	default:
		return nil
	}
}
