package bridge

import (
	"bufio"
	"context"
	"io"
	"time"

	metricnoop "go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

// Patience is the wedge-recycle threshold, exposed for external tests.
const Patience = patience

// Verify exposes the integrity check for external tests.
func Verify(options Options) error {
	return options.verify()
}

// Wrap builds a Driver over a test channel.
func Wrap(channel *Channel) *Driver {
	return &Driver{channel: channel}
}

// Probe drives the bounded frame reader in isolation.
type Probe struct {
	Source io.Reader
	Bound  int
}

func Scan(probe Probe) ([]byte, error) {
	live := &session{bound: probe.Bound}
	return live.scan(bufio.NewReaderSize(probe.Source, 64))
}

// Attach builds a channel whose sessions run over caller-provided pipes instead of
// a spawned process, so the transport and supervision are testable without a binary.
func Attach(factory func(context.Context) (io.Writer, io.Reader, func() error)) *Channel {
	duration, _ := metricnoop.NewMeterProvider().Meter("test").Float64Histogram("drizz.bridge.duration")
	channel := &Channel{
		options:  Options{Timeout: 2 * time.Second},
		slots:    make(chan struct{}, inflight),
		done:     make(chan struct{}),
		tracer:   tracenoop.NewTracerProvider().Tracer("test"),
		duration: duration,
	}
	channel.build = func(scope context.Context) (*session, error) {
		input, output, stop := factory(scope)
		live := &session{
			input:   input,
			stop:    stop,
			outbox:  make(chan []byte, inflight),
			depart:  make(chan struct{}),
			pending: make(map[int64]chan answer),
			bound:   frame,
		}
		go live.pump()
		go live.drain(output)
		if failure := live.greet(scope); failure != nil {
			live.expire()
			return nil, failure
		}
		return live, nil
	}
	return channel
}
