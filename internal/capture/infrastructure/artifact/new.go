package artifact

import "go.opentelemetry.io/otel/metric"

func New(options Options) (Store, error) {
	if failure := options.validate(); failure != nil {
		return Store{}, failure
	}

	latency, failure := options.Meter.Float64Histogram(instrument, metric.WithUnit("s"))
	if failure != nil {
		return Store{}, failure
	}

	temp, failure := options.prepare()
	if failure != nil {
		return Store{}, failure
	}

	limit := options.Ceiling
	if limit == 0 {
		limit = ceiling
	}

	return Store{
		temp:    temp,
		ceiling: limit,
		latency: latency,
		root:    options.Root,
		logger:  options.Logger,
		tracer:  options.Tracer,
	}, nil
}
