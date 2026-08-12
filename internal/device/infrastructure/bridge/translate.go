package bridge

import (
	"encoding/base64"

	"github.com/DrizzDev/platform/internal/device/application/control"
	"github.com/DrizzDev/platform/internal/device/domain/capture"
	"github.com/DrizzDev/platform/internal/device/domain/device"
	"github.com/DrizzDev/platform/internal/device/domain/format"
	"github.com/DrizzDev/platform/internal/device/domain/image"
	"github.com/DrizzDev/platform/internal/device/domain/platform"
	"github.com/DrizzDev/platform/internal/device/domain/serial"
)

var forward = map[platform.Platform]string{
	platform.Android:   "android",
	platform.Simulator: "ios-simulator",
	platform.Handset:   "ios-device",
}

var backward = map[string]platform.Platform{
	"android":       platform.Android,
	"ios-simulator": platform.Simulator,
	"ios-device":    platform.Handset,
}

// kinds is the single registry mapping a sidecar error class to a device port failure; an unmapped kind falls through
// to Unavailable at the call site.
var kinds = map[string]error{
	"DeviceNotFoundError":         control.Missing{},
	"DeviceUnauthorizedError":     control.Unauthorized{},
	"SecureScreenCaptureError":    control.Protected{},
	"DeviceOperationTimeoutError": control.Timeout{},
	"IncompatibleDeviceError":     control.Incompatible{},
	"UnsupportedOperationError":   control.Incompatible{},
	"DeviceCommandError":          control.Failed{},
}

type locator struct {
	Platform string `json:"platform"`
	UDID     string `json:"udid"`
}

type gesture struct {
	Platform string `json:"platform"`
	UDID     string `json:"udid"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
}

type listed struct {
	Platform string `json:"platform"`
	UDID     string `json:"udid"`
}

type framed struct {
	Format string `json:"format"`
	Data   string `json:"data"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func (listed listed) device() (device.Device, error) {
	family, found := backward[listed.Platform]
	if !found {
		return device.Device{}, control.Failed{}
	}
	handle, failure := serial.New(listed.UDID)
	if failure != nil {
		return device.Device{}, control.Failed{}
	}
	target, failure := device.New(device.Input{Platform: family, Serial: handle})
	if failure != nil {
		return device.Device{}, control.Failed{}
	}
	return target, nil
}

func (framed framed) capture() (capture.Capture, error) {
	if framed.Format != "png" {
		return capture.Capture{}, control.Failed{}
	}
	raw, failure := base64.StdEncoding.DecodeString(framed.Data)
	if failure != nil {
		return capture.Capture{}, control.Failed{}
	}
	pixels, failure := image.New(raw)
	if failure != nil {
		return capture.Capture{}, control.Failed{}
	}
	shot, failure := capture.New(capture.Input{Format: format.Png, Image: pixels, Width: framed.Width, Height: framed.Height})
	if failure != nil {
		return capture.Capture{}, control.Failed{}
	}
	return shot, nil
}
