package bridge

import (
	"context"
	"encoding/json"

	"github.com/DrizzDev/platform/internal/device/application/control"
	"github.com/DrizzDev/platform/internal/device/domain/capture"
	"github.com/DrizzDev/platform/internal/device/domain/device"
)

type portrait struct {
	Format string `json:"format"`
	Data   string `json:"data"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	XML    string `json:"xml"`
}

type sizing struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type tree struct {
	XML string `json:"xml"`
}

// Snapshot returns the screen image together with the on-screen element tree in one round trip.
func (driver *Driver) Snapshot(scope context.Context, target device.Device) (capture.Capture, string, error) {
	response, failure := driver.read(scope, Request{Method: "device.snapshot", Params: locator{
		UDID:     target.Serial().String(),
		Platform: forward[target.Platform()],
	}})
	if failure != nil {
		return capture.Capture{}, "", driver.fault(failure)
	}
	var reply portrait
	if json.Unmarshal(response.Result, &reply) != nil {
		return capture.Capture{}, "", control.Failed{}
	}
	shot, failure := framed{Format: reply.Format, Data: reply.Data, Width: reply.Width, Height: reply.Height}.capture()
	if failure != nil {
		return capture.Capture{}, "", failure
	}
	return shot, reply.XML, nil
}

func (driver *Driver) Hierarchy(scope context.Context, target device.Device) (string, error) {
	response, failure := driver.read(scope, Request{Method: "device.hierarchy", Params: locator{
		UDID:     target.Serial().String(),
		Platform: forward[target.Platform()],
	}})
	if failure != nil {
		return "", driver.fault(failure)
	}
	var reply tree
	if json.Unmarshal(response.Result, &reply) != nil {
		return "", control.Failed{}
	}
	return reply.XML, nil
}

func (driver *Driver) Dimensions(scope context.Context, target device.Device) (int, int, error) {
	response, failure := driver.read(scope, Request{Method: "device.dimensions", Params: locator{
		UDID:     target.Serial().String(),
		Platform: forward[target.Platform()],
	}})
	if failure != nil {
		return 0, 0, driver.fault(failure)
	}
	var reply sizing
	if json.Unmarshal(response.Result, &reply) != nil {
		return 0, 0, control.Failed{}
	}
	return reply.Width, reply.Height, nil
}
