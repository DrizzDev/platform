package mcp

import (
	"context"

	"github.com/DrizzDev/platform/internal/capability/application/operator"
)

// Perform runs one device capability. The capability operator satisfies it; the MCP server owns this port.
type Perform interface {
	Screenshot(context.Context, operator.Target) (operator.Shot, error)
	Devices(context.Context) (operator.Roster, error)
	Tap(context.Context, operator.Contact) (operator.Ack, error)
}
