package host

import (
	"context"

	"github.com/DrizzDev/platform/internal/capability/application/operator"
	"github.com/DrizzDev/platform/internal/capability/domain/outcome"
	"github.com/DrizzDev/platform/internal/transport/cli"
	"github.com/DrizzDev/platform/internal/transport/mcp"
)

var (
	_ cli.Perform = pilot{}
	_ mcp.Perform = pilot{}
)

// pilot is the host's device runtime. It is a cheap wrapper the surfaces hold, so the app always starts; the device
// connection and the recorder are assembled only when a capability actually runs. Until the device helper is present
// every capability reports the same typed, non-retryable message, and the very same runtime works end to end once the
// helper is in place — no change to the surfaces. Cloud capabilities will follow this identical shape.
type pilot struct {
	foundation
}

func (pilot pilot) Screenshot(scope context.Context, target operator.Target) (operator.Shot, error) {
	made, failure := pilot.assemble(scope)
	if failure != nil {
		return operator.Shot{}, failure
	}
	return made.Screenshot(scope, target)
}

func (pilot pilot) Devices(scope context.Context) (operator.Roster, error) {
	made, failure := pilot.assemble(scope)
	if failure != nil {
		return operator.Roster{}, failure
	}
	return made.Devices(scope)
}

// assemble builds the device operator on demand. The device helper is placed and pinned by installation; until it is
// present the capability is unprepared, which the surfaces render as one clear message. Installation completes this
// seam by resolving the helper here and assembling the device connection and the local recorder around it.
func (pilot pilot) assemble(context.Context) (operator.Operator, error) {
	return operator.Operator{}, operator.Refusal{Code: outcome.Unprepared}
}
