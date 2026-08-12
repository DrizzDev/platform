package emulator

import (
	"errors"

	"github.com/DrizzDev/platform/internal/device/domain/platform"
)

// Boot is an emulator image to start: the platform it runs on and the image name.
type Boot struct {
	platform platform.Platform
	name     string
}

type Input struct {
	Platform platform.Platform
	Name     string
}

func New(input Input) (Boot, error) {
	made := Boot{platform: input.Platform, name: input.Name}
	if failure := made.validate(); failure != nil {
		return Boot{}, failure
	}
	return made, nil
}

func (boot Boot) Platform() platform.Platform {
	return boot.platform
}

func (boot Boot) Name() string {
	return boot.name
}

func (boot Boot) validate() error {
	if boot.name == "" {
		return errors.New("emulator image name is required")
	}
	return nil
}
