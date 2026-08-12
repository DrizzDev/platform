package mcp

import (
	"context"

	"github.com/DrizzDev/platform/internal/capability/application/operator"
)

// Perform runs one device capability. The capability operator satisfies it; the MCP server owns this port.
type Perform interface {
	Screenshot(context.Context, operator.Target) (operator.Shot, error)
	Snapshot(context.Context, operator.Target) (operator.Snapshot, error)
	Hierarchy(context.Context, operator.Target) (operator.Tree, error)
	Dimensions(context.Context, operator.Target) (operator.Extent, error)
	Devices(context.Context) (operator.Roster, error)
	Disk(context.Context, operator.Target) (operator.Measure, error)
	Installed(context.Context, operator.Target) (operator.Listing, error)
	Running(context.Context, operator.Target) (operator.Listing, error)
	Foreground(context.Context, operator.Target) (operator.Report, error)
	Url(context.Context, operator.Target) (operator.Report, error)
	Images(context.Context) (operator.Images, error)
	Tap(context.Context, operator.Contact) (operator.Ack, error)
	Boot(context.Context, operator.Image) (operator.Ack, error)
	Pause(context.Context, operator.Target) (operator.Ack, error)
	Resume(context.Context, operator.Target) (operator.Ack, error)
	Install(context.Context, operator.Package) (operator.Ack, error)
	Launch(context.Context, operator.Application) (operator.Ack, error)
	Terminate(context.Context, operator.Application) (operator.Ack, error)
	Wipe(context.Context, operator.Application) (operator.Ack, error)
	Swipe(context.Context, operator.Drag) (operator.Ack, error)
	Pinch(context.Context, operator.Squeeze) (operator.Ack, error)
	Press(context.Context, operator.Key) (operator.Ack, error)
	Type(context.Context, operator.Entry) (operator.Ack, error)
	Clear(context.Context, operator.Target) (operator.Ack, error)
	Back(context.Context, operator.Target) (operator.Ack, error)
	Home(context.Context, operator.Target) (operator.Ack, error)
	Locate(context.Context, operator.Fix) (operator.Ack, error)
}
