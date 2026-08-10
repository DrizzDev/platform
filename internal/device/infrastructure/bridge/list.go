package bridge

import (
	"context"
	"encoding/json"

	"github.com/DrizzDev/platform/internal/device/application/control"
	"github.com/DrizzDev/platform/internal/device/domain/device"
)

func (driver *Driver) List(scope context.Context) ([]device.Device, error) {
	response, failure := driver.read(scope, Request{Method: "devices.list"})
	if failure != nil {
		return nil, driver.fault(failure)
	}
	var listing []listed
	if json.Unmarshal(response.Result, &listing) != nil || len(listing) > roster {
		return nil, control.Failed{}
	}
	devices := make([]device.Device, 0, len(listing))
	for _, entry := range listing {
		if target, valid := entry.device(); valid == nil {
			devices = append(devices, target)
		}
	}
	return devices, nil
}
