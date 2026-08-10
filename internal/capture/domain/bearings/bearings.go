package bearings

import "github.com/DrizzDev/platform/internal/capture/domain/identifier"

// Bearings are the cross-source identity dimensions that tie a hop to host and Drizz identities, used to match a hook
// with the MCP call it belongs to. Any dimension may be absent.
type Bearings struct {
	session    identifier.Identifier
	turn       identifier.Identifier
	call       identifier.Identifier
	connection identifier.Identifier
	execution  identifier.Identifier
	capability identifier.Identifier
}

type Input struct {
	Session    identifier.Identifier
	Turn       identifier.Identifier
	Call       identifier.Identifier
	Connection identifier.Identifier
	Execution  identifier.Identifier
	Capability identifier.Identifier
}

func New(input Input) Bearings {
	return Bearings{
		turn:       input.Turn,
		call:       input.Call,
		session:    input.Session,
		execution:  input.Execution,
		connection: input.Connection,
		capability: input.Capability,
	}
}

func (bearings Bearings) Session() identifier.Identifier    { return bearings.session }
func (bearings Bearings) Turn() identifier.Identifier       { return bearings.turn }
func (bearings Bearings) Call() identifier.Identifier       { return bearings.call }
func (bearings Bearings) Connection() identifier.Identifier { return bearings.connection }
func (bearings Bearings) Execution() identifier.Identifier  { return bearings.execution }
func (bearings Bearings) Capability() identifier.Identifier { return bearings.capability }
