package bridge

import (
	"context"
	"encoding/json"

	"github.com/DrizzDev/platform/internal/device/application/control"
	"github.com/DrizzDev/platform/internal/device/domain/capture"
	"github.com/DrizzDev/platform/internal/device/domain/device"
)

func (driver *Driver) Screenshot(scope context.Context, target device.Device) (capture.Capture, error) {
	response, failure := driver.read(scope, Request{Method: "device.screenshot", Params: locator{
		UDID:     target.Serial().String(),
		Platform: forward[target.Platform()],
	}})
	if failure != nil {
		return capture.Capture{}, driver.fault(failure)
	}
	var shot framed
	if json.Unmarshal(response.Result, &shot) != nil {
		return capture.Capture{}, control.Failed{}
	}
	return shot.capture()
}
