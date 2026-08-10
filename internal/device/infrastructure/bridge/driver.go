package bridge

import "github.com/DrizzDev/platform/internal/device/application/control"

var _ control.Bridge = (*Driver)(nil)

// Driver implements the device port over the sidecar channel: it builds each wire
// request and validates every reply into a domain value object.
type Driver struct {
	channel *Channel
}

func (driver *Driver) Close() error {
	return driver.channel.Close()
}
