package bridge

import (
	"context"
	"encoding/json"

	"github.com/DrizzDev/platform/internal/device/application/control"
	"github.com/DrizzDev/platform/internal/device/domain/device"
	"github.com/DrizzDev/platform/internal/device/domain/emulator"
	"github.com/DrizzDev/platform/internal/device/domain/platform"
)

type imaging struct {
	Platform string `json:"platform"`
}

type booting struct {
	Platform string `json:"platform"`
	Name     string `json:"name"`
}

type gallery struct {
	Images []string `json:"images"`
}

func (driver *Driver) Images(scope context.Context, family platform.Platform) ([]string, error) {
	response, failure := driver.read(scope, Request{Method: "device.images", Params: imaging{Platform: forward[family]}})
	if failure != nil {
		return nil, driver.fault(failure)
	}
	var reply gallery
	if json.Unmarshal(response.Result, &reply) != nil {
		return nil, control.Failed{}
	}
	return reply.Images, nil
}

func (driver *Driver) Boot(scope context.Context, image emulator.Boot) error {
	_, failure := driver.channel.Invoke(scope, Request{Method: "device.boot", Params: booting{
		Platform: forward[image.Platform()],
		Name:     image.Name(),
	}})
	if failure != nil {
		return driver.fault(failure)
	}
	return nil
}

func (driver *Driver) Pause(scope context.Context, target device.Device) error {
	_, failure := driver.channel.Invoke(scope, Request{Method: "device.pause", Params: locator{
		Platform: forward[target.Platform()],
		UDID:     target.Serial().String(),
	}})
	if failure != nil {
		return driver.fault(failure)
	}
	return nil
}

func (driver *Driver) Resume(scope context.Context, target device.Device) error {
	_, failure := driver.channel.Invoke(scope, Request{Method: "device.resume", Params: locator{
		Platform: forward[target.Platform()],
		UDID:     target.Serial().String(),
	}})
	if failure != nil {
		return driver.fault(failure)
	}
	return nil
}
