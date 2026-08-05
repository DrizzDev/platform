package console

import (
	"context"
	"fmt"
	"io"

	"github.com/DrizzDev/platform/internal/identity/application/device"
)

var _ device.Display = Display{}

// Display presents the device-authorization challenge on the terminal.
type Display struct {
	writer io.Writer
}

func (display Display) Show(_ context.Context, instruction device.Instruction) error {
	_, failure := fmt.Fprintf(display.writer,
		"To sign in, open %s and enter the code: %s\n", instruction.Verification, instruction.User)
	return failure
}
