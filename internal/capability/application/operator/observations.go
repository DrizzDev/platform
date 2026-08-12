package operator

import (
	"context"

	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
	"github.com/DrizzDev/platform/internal/capture/application/recording"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/fidelity"
	"github.com/DrizzDev/platform/internal/capture/domain/origin"
)

// Snapshot captures the screen of the target device together with its on-screen element tree, and records the capture.
func (operator Operator) Snapshot(scope context.Context, target Target) (shot Snapshot, failure error) {
	scope, watch := operator.begin(scope, catalog.Snapshot)
	defer func() { watch.finish(scope, failure) }()

	subject, failure := operator.resolve(scope, target.Serial)
	if failure != nil {
		return Snapshot{}, failure
	}
	portrait := operator.flow.Snapshot(scope, subject)
	if reason, failed := portrait.Failure(); failed {
		return Snapshot{}, operator.refuse(reason)
	}
	frame := portrait.Capture()
	operator.record(scope, frame)
	return Snapshot{Image: frame.Image().Bytes(), Format: frame.Format().String(), Hierarchy: portrait.Hierarchy()}, nil
}

// Hierarchy reads the on-screen element tree of the target device and records it.
func (operator Operator) Hierarchy(scope context.Context, target Target) (tree Tree, failure error) {
	scope, watch := operator.begin(scope, catalog.Hierarchy)
	defer func() { watch.finish(scope, failure) }()

	subject, failure := operator.resolve(scope, target.Serial)
	if failure != nil {
		return Tree{}, failure
	}
	reading := operator.flow.Hierarchy(scope, subject)
	if reason, failed := reading.Failure(); failed {
		return Tree{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{
		capability: catalog.Hierarchy,
		note: recording.Note{
			Fidelity: fidelity.Exact,
			Origin:   origin.Capability,
			Category: category.Hierarchy,
			Payload:  []byte("hierarchy"),
			Artifact: []byte(reading.Text()),
		},
	})
	return Tree{Hierarchy: reading.Text()}, nil
}

// Dimensions reads the screen size of the target device and records the read.
func (operator Operator) Dimensions(scope context.Context, target Target) (size Extent, failure error) {
	scope, watch := operator.begin(scope, catalog.Dimensions)
	defer func() { watch.finish(scope, failure) }()

	subject, failure := operator.resolve(scope, target.Serial)
	if failure != nil {
		return Extent{}, failure
	}
	extent := operator.flow.Dimensions(scope, subject)
	if reason, failed := extent.Failure(); failed {
		return Extent{}, operator.refuse(reason)
	}
	operator.inscribe(scope, entry{capability: catalog.Dimensions, note: operator.action("dimensions")})
	return Extent{Width: extent.Width(), Height: extent.Height()}, nil
}
