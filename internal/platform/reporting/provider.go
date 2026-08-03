package reporting

import (
	"context"
	"errors"
	"log/slog"
	"slices"
)

type Provider struct {
	sinks []Sink
}

func (provider Provider) Handler() slog.Handler {
	switch len(provider.sinks) {
	case 0:
		return nil
	case 1:
		return provider.sinks[0].Handler()
	default:
		handlers := make([]slog.Handler, 0, len(provider.sinks))
		for _, sink := range provider.sinks {
			handlers = append(handlers, sink.Handler())
		}
		return slog.NewMultiHandler(handlers...)
	}
}

func (provider Provider) Close(scope context.Context) error {
	var failures []error
	for _, v := range slices.Backward(provider.sinks) {
		failures = append(failures, v.Close(scope))
	}
	return errors.Join(failures...)
}
