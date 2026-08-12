package operator

import (
	"context"

	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
	"github.com/DrizzDev/platform/internal/capability/domain/outcome"
	"github.com/DrizzDev/platform/internal/device/application/control"
	"github.com/DrizzDev/platform/internal/device/domain/app"
	"github.com/DrizzDev/platform/internal/device/domain/bundle"
	"github.com/DrizzDev/platform/internal/device/domain/parcel"
)

// Install installs an application package on the target device and records the action.
func (operator Operator) Install(scope context.Context, target Package) (Ack, error) {
	subject, failure := operator.resolve(scope, target.Serial)
	if failure != nil {
		return Ack{}, failure
	}
	payload, failure := parcel.New(parcel.Input{Device: subject, Path: target.Path})
	if failure != nil {
		return Ack{}, Refusal{Code: outcome.Invalid}
	}
	if reason, failed := operator.flow.Install(scope, payload).Failure(); failed {
		return Ack{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{capability: catalog.Install, note: operator.action("install")})
	return Ack{}, nil
}

// Launch launches an application on the target device and records the action.
func (operator Operator) Launch(scope context.Context, target Application) (Ack, error) {
	return operator.address(scope, order{capability: catalog.Launch, serial: target.Serial, app: target.App, perform: operator.flow.Launch})
}

// Terminate stops a running application on the target device and records the action.
func (operator Operator) Terminate(scope context.Context, target Application) (Ack, error) {
	return operator.address(scope, order{capability: catalog.Terminate, serial: target.Serial, app: target.App, perform: operator.flow.Terminate})
}

// Wipe clears an application's stored data on the target device and records the action.
func (operator Operator) Wipe(scope context.Context, target Application) (Ack, error) {
	return operator.address(scope, order{capability: catalog.Wipe, serial: target.Serial, app: target.App, perform: operator.flow.Wipe})
}

// order names one application action: the capability, the device serial, the application identifier, and the flow
// method that performs it — so the app actions share one shape without dispatching by name.
type order struct {
	perform    func(context.Context, bundle.Bundle) control.Action
	capability string
	serial     string
	app        string
}

func (operator Operator) address(scope context.Context, spec order) (Ack, error) {
	subject, failure := operator.resolve(scope, spec.serial)
	if failure != nil {
		return Ack{}, failure
	}
	application, failure := bundle.New(bundle.Input{Device: subject, Id: spec.app})
	if failure != nil {
		return Ack{}, Refusal{Code: outcome.Invalid}
	}
	if reason, failed := spec.perform(scope, application).Failure(); failed {
		return Ack{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{capability: spec.capability, note: operator.action(spec.capability)})
	return Ack{}, nil
}

// Installed lists the applications installed on the target device and records the read.
func (operator Operator) Installed(scope context.Context, target Target) (Listing, error) {
	subject, failure := operator.resolve(scope, target.Serial)
	if failure != nil {
		return Listing{}, failure
	}
	inventory := operator.flow.Installed(scope, subject)
	if reason, failed := inventory.Failure(); failed {
		return Listing{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{capability: catalog.Installed, note: operator.action("installed")})
	return Listing{Apps: operator.catalogue(inventory.Apps())}, nil
}

// Running lists the applications currently running on the target device and records the read.
func (operator Operator) Running(scope context.Context, target Target) (Listing, error) {
	subject, failure := operator.resolve(scope, target.Serial)
	if failure != nil {
		return Listing{}, failure
	}
	inventory := operator.flow.Running(scope, subject)
	if reason, failed := inventory.Failure(); failed {
		return Listing{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{capability: catalog.Running, note: operator.action("running")})
	return Listing{Apps: operator.catalogue(inventory.Apps())}, nil
}

// Foreground reads the application in the foreground on the target device and records the read.
func (operator Operator) Foreground(scope context.Context, target Target) (Report, error) {
	subject, failure := operator.resolve(scope, target.Serial)
	if failure != nil {
		return Report{}, failure
	}
	reading := operator.flow.Foreground(scope, subject)
	if reason, failed := reading.Failure(); failed {
		return Report{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{capability: catalog.Foreground, note: operator.action("foreground")})
	return Report{Text: reading.Text()}, nil
}

// Url reads the current link open in the active application on the target device and records the read.
func (operator Operator) Url(scope context.Context, target Target) (Report, error) {
	subject, failure := operator.resolve(scope, target.Serial)
	if failure != nil {
		return Report{}, failure
	}
	reading := operator.flow.Url(scope, subject)
	if reason, failed := reading.Failure(); failed {
		return Report{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{capability: catalog.Url, note: operator.action("url")})
	return Report{Text: reading.Text()}, nil
}

// Disk reads the free disk space on the target device and records the read.
func (operator Operator) Disk(scope context.Context, target Target) (Measure, error) {
	subject, failure := operator.resolve(scope, target.Serial)
	if failure != nil {
		return Measure{}, failure
	}
	measure := operator.flow.Disk(scope, subject)
	if reason, failed := measure.Failure(); failed {
		return Measure{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{capability: catalog.Disk, note: operator.action("disk")})
	return Measure{Value: measure.Value()}, nil
}

func (Operator) catalogue(source []app.App) []App {
	listing := make([]App, 0, len(source))
	for _, item := range source {
		listing = append(listing, App{Id: item.Id(), Name: item.Name(), Note: item.Note()})
	}
	return listing
}
