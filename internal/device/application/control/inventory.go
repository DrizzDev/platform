package control

import "github.com/DrizzDev/platform/internal/device/domain/app"

// Inventory is a set of applications on a device, or a code-only failure.
type Inventory struct {
	outcome
	apps []app.App
}

func (inventory Inventory) Apps() []app.App {
	return append([]app.App(nil), inventory.apps...)
}
