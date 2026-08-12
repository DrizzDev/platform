package host

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/DrizzDev/platform/internal/transport/cli"
)

type Host struct {
	command cli.Command
	failure io.Writer
	device  *station
}

func (host *Host) Run(scope context.Context) error {
	defer host.device.close()
	result := host.command.Run(scope)
	if result == nil {
		return nil
	}
	var marker handled
	if !errors.As(result, &marker) {
		_, _ = fmt.Fprintln(host.failure, result)
	}
	return result
}
