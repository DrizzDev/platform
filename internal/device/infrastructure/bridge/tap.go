package bridge

import (
	"context"

	"github.com/DrizzDev/platform/internal/device/domain/touch"
)

// Tap is a mutation and is never retried: on an ambiguous transport failure it reports unavailable rather than risk
// a duplicate press (REL-012, REL-015).
func (driver *Driver) Tap(scope context.Context, contact touch.Touch) error {

	spot := contact.Point()
	target := contact.Device()

	_, failure := driver.channel.Invoke(scope, Request{Method: "device.tap", Params: gesture{
		X:        spot.X(),
		Y:        spot.Y(),
		UDID:     target.Serial().String(),
		Platform: forward[target.Platform()],
	}})

	if failure != nil {
		return driver.fault(failure)
	}

	return nil
}
