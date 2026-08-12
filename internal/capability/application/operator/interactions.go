package operator

import (
	"context"
	"fmt"
	"time"

	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
	"github.com/DrizzDev/platform/internal/capability/domain/outcome"
	"github.com/DrizzDev/platform/internal/device/domain/geo"
	"github.com/DrizzDev/platform/internal/device/domain/pinch"
	"github.com/DrizzDev/platform/internal/device/domain/point"
	"github.com/DrizzDev/platform/internal/device/domain/press"
	"github.com/DrizzDev/platform/internal/device/domain/swipe"
	"github.com/DrizzDev/platform/internal/device/domain/text"
)

// Swipe drags across the target device from one point to another and records the action.
func (operator Operator) Swipe(scope context.Context, drag Drag) (ack Ack, failure error) {
	scope, watch := operator.begin(scope, catalog.Swipe)
	defer func() { watch.finish(scope, failure) }()

	subject, failure := operator.resolve(scope, drag.Serial)
	if failure != nil {
		return Ack{}, failure
	}
	from, failure := point.New(point.Input{X: drag.From.X, Y: drag.From.Y})
	if failure != nil {
		return Ack{}, Refusal{Code: outcome.Invalid}
	}
	to, failure := point.New(point.Input{X: drag.To.X, Y: drag.To.Y})
	if failure != nil {
		return Ack{}, Refusal{Code: outcome.Invalid}
	}
	gesture, failure := swipe.New(swipe.Input{Device: subject, From: from, To: to, Span: time.Duration(drag.Milliseconds) * time.Millisecond})
	if failure != nil {
		return Ack{}, Refusal{Code: outcome.Invalid}
	}
	if reason, failed := operator.flow.Swipe(scope, gesture).Failure(); failed {
		return Ack{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{
		capability: catalog.Swipe,
		note:       operator.action(fmt.Sprintf("swipe %d,%d to %d,%d", drag.From.X, drag.From.Y, drag.To.X, drag.To.Y)),
	})
	return Ack{}, nil
}

// Pinch zooms around a centre point on the target device and records the action.
func (operator Operator) Pinch(scope context.Context, squeeze Squeeze) (ack Ack, failure error) {
	scope, watch := operator.begin(scope, catalog.Pinch)
	defer func() { watch.finish(scope, failure) }()

	subject, failure := operator.resolve(scope, squeeze.Serial)
	if failure != nil {
		return Ack{}, failure
	}
	centre, failure := point.New(point.Input{X: squeeze.Centre.X, Y: squeeze.Centre.Y})
	if failure != nil {
		return Ack{}, Refusal{Code: outcome.Invalid}
	}
	gesture, failure := pinch.New(pinch.Input{Device: subject, Centre: centre, From: squeeze.From, To: squeeze.To})
	if failure != nil {
		return Ack{}, Refusal{Code: outcome.Invalid}
	}
	if reason, failed := operator.flow.Pinch(scope, gesture).Failure(); failed {
		return Ack{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{
		capability: catalog.Pinch,
		note:       operator.action(fmt.Sprintf("pinch %d,%d %d to %d", squeeze.Centre.X, squeeze.Centre.Y, squeeze.From, squeeze.To)),
	})
	return Ack{}, nil
}

// Press presses a hardware or remote button on the target device and records the action.
func (operator Operator) Press(scope context.Context, key Key) (ack Ack, failure error) {
	scope, watch := operator.begin(scope, catalog.Press)
	defer func() { watch.finish(scope, failure) }()

	subject, failure := operator.resolve(scope, key.Serial)
	if failure != nil {
		return Ack{}, failure
	}
	button, failure := press.New(press.Input{Device: subject, Button: press.Button(key.Button)})
	if failure != nil {
		return Ack{}, Refusal{Code: outcome.Invalid}
	}
	if reason, failed := operator.flow.Press(scope, button).Failure(); failed {
		return Ack{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{capability: catalog.Press, note: operator.action("press " + key.Button)})
	return Ack{}, nil
}

// Type types text into the focused field on the target device and records the length typed, never the text itself.
func (operator Operator) Type(scope context.Context, input Entry) (ack Ack, failure error) {
	scope, watch := operator.begin(scope, catalog.Type)
	defer func() { watch.finish(scope, failure) }()

	subject, failure := operator.resolve(scope, input.Serial)
	if failure != nil {
		return Ack{}, failure
	}
	typed, failure := text.New(text.Input{Device: subject, Content: input.Text})
	if failure != nil {
		return Ack{}, Refusal{Code: outcome.Invalid}
	}
	if reason, failed := operator.flow.Type(scope, typed).Failure(); failed {
		return Ack{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{
		capability: catalog.Type,
		note:       operator.action(fmt.Sprintf("type %d characters", len(input.Text))),
	})
	return Ack{}, nil
}

// Clear clears the focused text field on the target device and records the action.
func (operator Operator) Clear(scope context.Context, target Target) (ack Ack, failure error) {
	scope, watch := operator.begin(scope, catalog.Clear)
	defer func() { watch.finish(scope, failure) }()

	subject, failure := operator.resolve(scope, target.Serial)
	if failure != nil {
		return Ack{}, failure
	}
	if reason, failed := operator.flow.Clear(scope, subject).Failure(); failed {
		return Ack{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{capability: catalog.Clear, note: operator.action("clear")})
	return Ack{}, nil
}

// Back presses the back button on the target device and records the action.
func (operator Operator) Back(scope context.Context, target Target) (ack Ack, failure error) {
	scope, watch := operator.begin(scope, catalog.Back)
	defer func() { watch.finish(scope, failure) }()

	subject, failure := operator.resolve(scope, target.Serial)
	if failure != nil {
		return Ack{}, failure
	}
	if reason, failed := operator.flow.Back(scope, subject).Failure(); failed {
		return Ack{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{capability: catalog.Back, note: operator.action("back")})
	return Ack{}, nil
}

// Home presses the home button on the target device and records the action.
func (operator Operator) Home(scope context.Context, target Target) (ack Ack, failure error) {
	scope, watch := operator.begin(scope, catalog.Home)
	defer func() { watch.finish(scope, failure) }()

	subject, failure := operator.resolve(scope, target.Serial)
	if failure != nil {
		return Ack{}, failure
	}
	if reason, failed := operator.flow.Home(scope, subject).Failure(); failed {
		return Ack{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{capability: catalog.Home, note: operator.action("home")})
	return Ack{}, nil
}

// Locate sets the reported location of the target device and records the action.
func (operator Operator) Locate(scope context.Context, fix Fix) (ack Ack, failure error) {
	scope, watch := operator.begin(scope, catalog.Locate)
	defer func() { watch.finish(scope, failure) }()

	subject, failure := operator.resolve(scope, fix.Serial)
	if failure != nil {
		return Ack{}, failure
	}
	position, failure := geo.New(geo.Input{Device: subject, Lat: fix.Lat, Lon: fix.Lon})
	if failure != nil {
		return Ack{}, Refusal{Code: outcome.Invalid}
	}
	if reason, failed := operator.flow.Locate(scope, position).Failure(); failed {
		return Ack{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{
		capability: catalog.Locate,
		note:       operator.action(fmt.Sprintf("locate %.5f,%.5f", fix.Lat, fix.Lon)),
	})
	return Ack{}, nil
}
