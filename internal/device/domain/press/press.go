package press

import (
	"errors"

	"github.com/DrizzDev/platform/internal/device/domain/device"
)

// Button is a hardware or remote button that can be pressed on a device.
type Button string

const (
	Up        Button = "up"
	Down      Button = "down"
	Left      Button = "left"
	Right     Button = "right"
	Select    Button = "select"
	Menu      Button = "menu"
	Home      Button = "home"
	Playpause Button = "playpause"
)

func (button Button) Valid() bool {
	switch button {
	case Up, Down, Left, Right, Select, Menu, Home, Playpause:
		return true
	default:
		return false
	}
}

func (button Button) String() string {
	return string(button)
}

// Press is a single button press on a device.
type Press struct {
	device device.Device
	button Button
}

type Input struct {
	Device device.Device
	Button Button
}

func New(input Input) (Press, error) {
	made := Press{device: input.Device, button: input.Button}
	if failure := made.validate(); failure != nil {
		return Press{}, failure
	}
	return made, nil
}

func (press Press) Device() device.Device {
	return press.device
}

func (press Press) Button() Button {
	return press.button
}

func (press Press) validate() error {
	switch {
	case press.device.Serial().String() == "":
		return errors.New("press device is required")
	case !press.button.Valid():
		return errors.New("press button is not recognized")
	default:
		return nil
	}
}
