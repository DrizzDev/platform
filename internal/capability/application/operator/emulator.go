package operator

import (
	"context"

	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
	"github.com/DrizzDev/platform/internal/capability/domain/outcome"
	"github.com/DrizzDev/platform/internal/device/domain/emulator"
	"github.com/DrizzDev/platform/internal/device/domain/platform"
)

// Images lists the emulator images available to run and records the read.
func (operator Operator) Images(scope context.Context) (Images, error) {
	catalogue := operator.flow.Images(scope, platform.Android)
	if reason, failed := catalogue.Failure(); failed {
		return Images{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{capability: catalog.Images, note: operator.action("images")})
	return Images{Names: catalogue.Names()}, nil
}

// Boot starts an emulator from an image and records the action.
func (operator Operator) Boot(scope context.Context, target Image) (Ack, error) {
	image, failure := emulator.New(emulator.Input{Platform: platform.Android, Name: target.Name})
	if failure != nil {
		return Ack{}, Refusal{Code: outcome.Invalid}
	}
	if reason, failed := operator.flow.Boot(scope, image).Failure(); failed {
		return Ack{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{capability: catalog.Boot, note: operator.action("boot " + target.Name)})
	return Ack{}, nil
}

// Pause pauses a running emulator and records the action.
func (operator Operator) Pause(scope context.Context, target Target) (Ack, error) {
	subject, failure := operator.resolve(scope, target.Serial)
	if failure != nil {
		return Ack{}, failure
	}
	if reason, failed := operator.flow.Pause(scope, subject).Failure(); failed {
		return Ack{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{capability: catalog.Pause, note: operator.action("pause")})
	return Ack{}, nil
}

// Resume resumes a paused emulator and records the action.
func (operator Operator) Resume(scope context.Context, target Target) (Ack, error) {
	subject, failure := operator.resolve(scope, target.Serial)
	if failure != nil {
		return Ack{}, failure
	}
	if reason, failed := operator.flow.Resume(scope, subject).Failure(); failed {
		return Ack{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{capability: catalog.Resume, note: operator.action("resume")})
	return Ack{}, nil
}
