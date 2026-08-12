package bridge

import "go.opentelemetry.io/otel/metric"

func New(options Options) (*Driver, error) {
	if failure := options.validate(); failure != nil {
		return nil, failure
	}
	duration, failure := options.Meter.Float64Histogram("drizz.bridge.duration", metric.WithUnit("s"))
	if failure != nil {
		return nil, failure
	}
	channel := &Channel{
		options:  options,
		done:     make(chan struct{}),
		slots:    make(chan struct{}, inflight),
		tracer:   options.Tracer,
		duration: duration,
	}
	channel.build = channel.spawn
	return &Driver{channel: channel}, nil
}
