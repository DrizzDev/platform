package bridge

import (
	"context"

	"github.com/DrizzDev/platform/internal/device/domain/device"
	"github.com/DrizzDev/platform/internal/device/domain/geo"
	"github.com/DrizzDev/platform/internal/device/domain/pinch"
	"github.com/DrizzDev/platform/internal/device/domain/press"
	"github.com/DrizzDev/platform/internal/device/domain/swipe"
	"github.com/DrizzDev/platform/internal/device/domain/text"
)

// erasure is how many characters a clear request deletes from the focused field, well above a typical field's length.
const erasure = 256

type spot struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type dragging struct {
	Platform string `json:"platform"`
	UDID     string `json:"udid"`
	Points   []spot `json:"points"`
	Duration int    `json:"durationMs"`
}

type squeeze struct {
	Platform string `json:"platform"`
	UDID     string `json:"udid"`
	X        int    `json:"centerX"`
	Y        int    `json:"centerY"`
	Inner    int    `json:"startRadius"`
	Outer    int    `json:"endRadius"`
}

type keypress struct {
	Platform string `json:"platform"`
	UDID     string `json:"udid"`
	Button   string `json:"button"`
}

type typing struct {
	Platform string `json:"platform"`
	UDID     string `json:"udid"`
	Text     string `json:"text"`
}

type erasing struct {
	Platform string `json:"platform"`
	UDID     string `json:"udid"`
	Length   int    `json:"length"`
}

type locating struct {
	Platform  string  `json:"platform"`
	UDID      string  `json:"udid"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// These are mutations, so — like Tap — they are never auto-retried: on an ambiguous transport failure they report
// unavailable rather than risk performing the gesture twice.
func (driver *Driver) Swipe(scope context.Context, drag swipe.Swipe) error {
	target := drag.Device()
	_, failure := driver.channel.Invoke(scope, Request{Method: "device.swipe", Params: dragging{
		Platform: forward[target.Platform()],
		UDID:     target.Serial().String(),
		Points:   []spot{{X: drag.From().X(), Y: drag.From().Y()}, {X: drag.To().X(), Y: drag.To().Y()}},
		Duration: int(drag.Span().Milliseconds()),
	}})
	if failure != nil {
		return driver.fault(failure)
	}
	return nil
}

func (driver *Driver) Pinch(scope context.Context, gesture pinch.Pinch) error {
	target := gesture.Device()
	_, failure := driver.channel.Invoke(scope, Request{Method: "device.pinch", Params: squeeze{
		Platform: forward[target.Platform()],
		UDID:     target.Serial().String(),
		X:        gesture.Centre().X(),
		Y:        gesture.Centre().Y(),
		Inner:    gesture.From(),
		Outer:    gesture.To(),
	}})
	if failure != nil {
		return driver.fault(failure)
	}
	return nil
}

func (driver *Driver) Press(scope context.Context, key press.Press) error {
	target := key.Device()
	_, failure := driver.channel.Invoke(scope, Request{Method: "device.press", Params: keypress{
		Platform: forward[target.Platform()],
		UDID:     target.Serial().String(),
		Button:   key.Button().String(),
	}})
	if failure != nil {
		return driver.fault(failure)
	}
	return nil
}

func (driver *Driver) Type(scope context.Context, entry text.Text) error {
	target := entry.Device()
	_, failure := driver.channel.Invoke(scope, Request{Method: "device.type", Params: typing{
		Platform: forward[target.Platform()],
		UDID:     target.Serial().String(),
		Text:     entry.Content(),
	}})
	if failure != nil {
		return driver.fault(failure)
	}
	return nil
}

func (driver *Driver) Clear(scope context.Context, target device.Device) error {
	_, failure := driver.channel.Invoke(scope, Request{Method: "device.clear", Params: erasing{
		Platform: forward[target.Platform()],
		UDID:     target.Serial().String(),
		Length:   erasure,
	}})
	if failure != nil {
		return driver.fault(failure)
	}
	return nil
}

func (driver *Driver) Back(scope context.Context, target device.Device) error {
	_, failure := driver.channel.Invoke(scope, Request{Method: "device.back", Params: locator{
		Platform: forward[target.Platform()],
		UDID:     target.Serial().String(),
	}})
	if failure != nil {
		return driver.fault(failure)
	}
	return nil
}

func (driver *Driver) Home(scope context.Context, target device.Device) error {
	_, failure := driver.channel.Invoke(scope, Request{Method: "device.home", Params: locator{
		Platform: forward[target.Platform()],
		UDID:     target.Serial().String(),
	}})
	if failure != nil {
		return driver.fault(failure)
	}
	return nil
}

func (driver *Driver) Locate(scope context.Context, fix geo.Fix) error {
	target := fix.Device()
	_, failure := driver.channel.Invoke(scope, Request{Method: "device.locate", Params: locating{
		Platform:  forward[target.Platform()],
		UDID:      target.Serial().String(),
		Latitude:  fix.Lat(),
		Longitude: fix.Lon(),
	}})
	if failure != nil {
		return driver.fault(failure)
	}
	return nil
}
