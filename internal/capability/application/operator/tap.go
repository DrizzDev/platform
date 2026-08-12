package operator

import (
	"context"
	"fmt"

	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
	"github.com/DrizzDev/platform/internal/capability/domain/outcome"
	"github.com/DrizzDev/platform/internal/device/domain/point"
	"github.com/DrizzDev/platform/internal/device/domain/touch"
)

// Tap presses the target device at a point and records the action. It resolves the serial to a connected device first,
// so a missing or unauthorized device is reported as a stable code rather than a raw failure.
func (operator Operator) Tap(scope context.Context, contact Contact) (Ack, error) {
	subject, failure := operator.resolve(scope, contact.Serial)
	if failure != nil {
		return Ack{}, failure
	}
	spot, failure := point.New(point.Input{X: contact.X, Y: contact.Y})
	if failure != nil {
		return Ack{}, Refusal{Code: outcome.Invalid}
	}
	press, failure := touch.New(touch.Input{Device: subject, Point: spot})
	if failure != nil {
		return Ack{}, Refusal{Code: outcome.Invalid}
	}
	performed := operator.flow.Act(scope, press)
	if reason, failed := performed.Failure(); failed {
		return Ack{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{
		capability: catalog.Tap,
		note:       operator.action(fmt.Sprintf("tap %d,%d", contact.X, contact.Y)),
	})
	return Ack{}, nil
}
