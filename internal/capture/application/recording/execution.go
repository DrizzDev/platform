package recording

import (
	"bytes"
	"context"
	"sync/atomic"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
	"github.com/DrizzDev/platform/internal/capture/domain/correlation"
	"github.com/DrizzDev/platform/internal/capture/domain/digest"
	"github.com/DrizzDev/platform/internal/capture/domain/mark"
	"github.com/DrizzDev/platform/internal/capture/domain/span"
	"github.com/DrizzDev/platform/internal/capture/domain/trace"
)

// Execution is one recorded run: a trace shared by every step, a root span, and a monotonic sequence.
// Each step is recorded as a child of the root, and concurrent steps are safe.
type Execution struct {
	recorder Recorder
	root     span.Span
	trace    trace.Trace
	sequence atomic.Int64
}

func (recorder Recorder) Begin() (*Execution, error) {
	seed, failure := recorder.mint()
	if failure != nil {
		return nil, failure
	}
	stem, failure := recorder.mint()
	if failure != nil {
		return nil, failure
	}
	thread, failure := trace.New(seed)
	if failure != nil {
		return nil, failure
	}
	root, failure := span.New(stem)
	if failure != nil {
		return nil, failure
	}
	return &Execution{recorder: recorder, trace: thread, root: root}, nil
}

func (execution *Execution) Trace() trace.Trace {
	return execution.trace
}

// Record writes one step of the execution. It never returns an error: a failure to store the artifact or append the
// entry is logged and dropped, so recording cannot break the capability it observes.
func (execution *Execution) Record(scope context.Context, note Note) {

	execution.renew(scope)

	value, failure := execution.recorder.mint()
	if failure != nil {
		execution.recorder.drop(scope, "span")
		return
	}

	hop, failure := span.New(value)
	if failure != nil {
		execution.recorder.drop(scope, "span")
		return
	}

	stamp := note.Mark
	if !stamp.Valid() {
		stamp = mark.Exact
	}

	link, failure := correlation.New(correlation.Input{
		Span:     hop,
		Mark:     stamp,
		Parent:   execution.root,
		Trace:    execution.trace,
		Sequence: execution.sequence.Add(1),
	})
	if failure != nil {
		execution.recorder.drop(scope, "correlation")
		return
	}

	reference, failure := execution.store(scope, note.Artifact)
	if failure != nil {
		execution.recorder.drop(scope, "artifact")
		return
	}

	entry, failure := journal.New(journal.Input{
		Correlation: link,
		Artifact:    reference,
		Origin:      note.Origin,
		Payload:     note.Payload,
		Fidelity:    note.Fidelity,
		Category:    note.Category,
		State:       journal.Pending,
	})
	if failure != nil {
		execution.recorder.drop(scope, "entry")
		return
	}
	if failure := execution.recorder.writer.Append(scope, entry); failure != nil {
		execution.recorder.drop(scope, "append")
	}
}

// renew extends the execution's lease so a reclaim pass leaves its records alone while it runs. A failure to lease is
// logged and dropped, never surfaced, so protection is best-effort and recording stays observational.
func (execution *Execution) renew(scope context.Context) {
	deadline := execution.recorder.clock.Now().Add(execution.recorder.lease)

	claim := journal.Claim{Trace: execution.trace, Until: deadline}
	if failure := execution.recorder.keeper.Lease(scope, claim); failure != nil {
		execution.recorder.drop(scope, "lease")
	}
}

func (execution *Execution) store(scope context.Context, blob []byte) (digest.Digest, error) {
	if len(blob) == 0 {
		return digest.Digest{}, nil
	}

	return execution.recorder.sink.Put(scope, bytes.NewReader(blob))
}
