package bridge

import (
	"context"
	"encoding/json"

	"github.com/DrizzDev/platform/internal/device/application/control"
	"github.com/DrizzDev/platform/internal/device/domain/app"
	"github.com/DrizzDev/platform/internal/device/domain/bundle"
	"github.com/DrizzDev/platform/internal/device/domain/device"
	"github.com/DrizzDev/platform/internal/device/domain/parcel"
)

type installing struct {
	Platform string `json:"platform"`
	UDID     string `json:"udid"`
	Path     string `json:"path"`
}

type naming struct {
	Platform string `json:"platform"`
	UDID     string `json:"udid"`
	App      string `json:"app"`
}

type listing struct {
	Apps []named `json:"apps"`
}

type named struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Note string `json:"note"`
}

type valued struct {
	Value string `json:"value"`
}

type spaced struct {
	Megabytes int `json:"megabytes"`
}

func (driver *Driver) Install(scope context.Context, payload parcel.Parcel) error {
	target := payload.Device()
	_, failure := driver.channel.Invoke(scope, Request{Method: "device.install", Params: installing{
		Platform: forward[target.Platform()],
		UDID:     target.Serial().String(),
		Path:     payload.Path(),
	}})
	if failure != nil {
		return driver.fault(failure)
	}
	return nil
}

func (driver *Driver) Launch(scope context.Context, application bundle.Bundle) error {
	return driver.manage(scope, instruction{method: "device.launch", application: application})
}

func (driver *Driver) Terminate(scope context.Context, application bundle.Bundle) error {
	return driver.manage(scope, instruction{method: "device.terminate", application: application})
}

func (driver *Driver) Wipe(scope context.Context, application bundle.Bundle) error {
	return driver.manage(scope, instruction{method: "device.wipe", application: application})
}

// manage sends an application action naming the target app. Like every mutation it is never auto-retried.
func (driver *Driver) manage(scope context.Context, request instruction) error {
	target := request.application.Device()
	_, failure := driver.channel.Invoke(scope, Request{Method: request.method, Params: naming{
		Platform: forward[target.Platform()],
		UDID:     target.Serial().String(),
		App:      request.application.Id(),
	}})
	if failure != nil {
		return driver.fault(failure)
	}
	return nil
}

// instruction names an application action: the wire method and the target application.
type instruction struct {
	method      string
	application bundle.Bundle
}

func (driver *Driver) Installed(scope context.Context, target device.Device) ([]app.App, error) {
	response, failure := driver.read(scope, Request{Method: "device.installed", Params: locator{
		UDID:     target.Serial().String(),
		Platform: forward[target.Platform()],
	}})
	if failure != nil {
		return nil, driver.fault(failure)
	}
	return driver.apps(response)
}

func (driver *Driver) Running(scope context.Context, target device.Device) ([]app.App, error) {
	response, failure := driver.read(scope, Request{Method: "device.running", Params: locator{
		UDID:     target.Serial().String(),
		Platform: forward[target.Platform()],
	}})
	if failure != nil {
		return nil, driver.fault(failure)
	}
	return driver.apps(response)
}

func (driver *Driver) Foreground(scope context.Context, target device.Device) (string, error) {
	return driver.value(scope, probe{method: "device.foreground", target: target})
}

func (driver *Driver) Url(scope context.Context, target device.Device) (string, error) {
	return driver.value(scope, probe{method: "device.url", target: target})
}

// probe names a metadata read: the wire method and the target device.
type probe struct {
	method string
	target device.Device
}

func (driver *Driver) Disk(scope context.Context, target device.Device) (int, error) {
	response, failure := driver.read(scope, Request{Method: "device.disk", Params: locator{
		UDID:     target.Serial().String(),
		Platform: forward[target.Platform()],
	}})
	if failure != nil {
		return 0, driver.fault(failure)
	}
	var reply spaced
	if json.Unmarshal(response.Result, &reply) != nil {
		return 0, control.Failed{}
	}
	return reply.Megabytes, nil
}

func (driver *Driver) apps(response Response) ([]app.App, error) {
	var reply listing
	if json.Unmarshal(response.Result, &reply) != nil {
		return nil, control.Failed{}
	}
	apps := make([]app.App, 0, len(reply.Apps))
	for _, item := range reply.Apps {
		apps = append(apps, app.New(app.Input{Id: item.Id, Name: item.Name, Note: item.Note}))
	}
	return apps, nil
}

func (driver *Driver) value(scope context.Context, request probe) (string, error) {
	response, failure := driver.read(scope, Request{Method: request.method, Params: locator{
		UDID:     request.target.Serial().String(),
		Platform: forward[request.target.Platform()],
	}})
	if failure != nil {
		return "", driver.fault(failure)
	}
	var reply valued
	if json.Unmarshal(response.Result, &reply) != nil {
		return "", control.Failed{}
	}
	return reply.Value, nil
}
